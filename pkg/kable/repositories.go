package kable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/go-git/go-git/v5/plumbing"

	"github.com/go-git/go-git/v5"
)

const (
	RepoIndexFilename      = "kable.json"
	RepositoryInvalidError = "given repository is not a valid kable repository"
)

type Repositories map[uuid.UUID]Repository

type Repository struct {
	RepoURL string
	Branch  string
}

type RepoIndex struct {
	Version        int            `json:"version"`
	Name           string         `json:"name"`
	ConceptEntries []ConceptEntry `json:"concepts"`
}

type ConceptEntry string

func (ce ConceptEntry) String() string {
	return string(ce)
}

func (repos Repositories) ToArray() ([][]interface{}, error) {
	var repoSlices [][]interface{}
	for id, repo := range repos {
		if IsInitialized(id) {
			repoSlices = append(repoSlices, []interface{}{MustGetRepoIndex(id).Name, repo.RepoURL, id, true})
		} else {
			repoSlices = append(repoSlices, []interface{}{"", repo.RepoURL, id, false})
		}
	}
	return repoSlices, nil
}

type Cloner struct {
	Config       ClonerConfig
	Repositories Repositories
}

type ClonerConfig struct {
	git.CloneOptions
	BaseDir string
}

type CloneIndex map[uuid.UUID]CloneRef

type CloneRef struct {
	Path   string
	Branch string
	URI    string
}

func AddToCacheIndex(id uuid.UUID, url, branch string) error {
	ci := CloneIndex{}
	cipath := filepath.Join(repoDir, "cloneindex.json")
	content, _ := ioutil.ReadFile(cipath)
	_ = json.Unmarshal(content, &ci)

	path := MustGetCachePath(id)

	ci[id] = CloneRef{
		Branch: branch,
		Path:   path,
		URI:    url,
	}
	out, err := json.MarshalIndent(ci, "", "	")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(cipath, out, 0666); err != nil {
		return err
	}
	return nil
}

func GetCacheIndex() (CloneIndex, error) {
	ci := CloneIndex{}
	cipath := filepath.Join(repoDir, "cloneindex.json")
	content, err := ioutil.ReadFile(cipath)
	if err == nil {
		if err := json.Unmarshal(content, &ci); err != nil {
			return ci, err
		}
	}

	return ci, nil
}

func GetCacheInfo(id uuid.UUID) (*CloneRef, error) {
	ci, err := GetCacheIndex()
	if err != nil {
		return nil, err
	}
	out := ci[id]
	return &out, nil
}

func RemoveFromCacheIndex(id uuid.UUID) error {
	ci := CloneIndex{}
	cipath := filepath.Join(repoDir, "cloneindex.json")
	content, _ := ioutil.ReadFile(cipath)
	_ = json.Unmarshal(content, &ci)

	delete(ci, id)
	out, err := json.MarshalIndent(ci, "", "	")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(cipath, out, 0666); err != nil {
		return err
	}
	return nil
}

// Returns a UUID if cached, else just returns nil.
func IsCached(url, branch string) (*uuid.UUID, error) {
	var uid uuid.UUID
	idx, err := GetCacheIndex()
	if err != nil {
		return &uid, err
	}
	for existingId, ref := range idx {
		if ref.URI == url && ref.Branch == branch {
			uid = existingId
			return &uid, nil
		}
	}
	return nil, nil

}

func AddRepository(url, branch string) (string, error) {
	return currentConfig.Repositories.AddRepository(url, branch)
}

func UpdateRepositories() error {
	return currentConfig.Repositories.UpdateRepositories()
}

func RemoveRepository(name string) error {
	return currentConfig.Repositories.RemoveRepository(name)
}

// TidyRepositories cleans up all cached repositories that are not in use in the current config,
// and returns the names of the deleted ones.
func TidyRepositories() ([]string, error) {
	return currentConfig.Repositories.TidyRepositories()
}

func (repos Repositories) TidyRepositories() ([]string, error) {
	var repolist []string

	idx, err := GetCacheIndex()
	if err != nil {
		return repolist, err
	}

	// Go through index, and delete all repositories, that are not used in the current config
	for id, ref := range idx {
		if _, exists := currentConfig.Repositories[id]; !exists {
			if err := os.RemoveAll(ref.Path); err != nil {
				return repolist, err
			}
			if err := RemoveFromCacheIndex(id); err != nil {
				return repolist, err
			}
			repolist = append(repolist, MustGetRepoIndex(id).Name)
		}
	}

	return repolist, nil
}

func (repos Repositories) AddRepository(url, branch string) (string, error) {
	if !configSet() {
		return "", fmt.Errorf(ConfigNotInitializedError)
	}

	for _, repo := range currentConfig.Repositories {
		if repo.RepoURL == url && repo.Branch == branch {
			return "", fmt.Errorf(RepositoryAlreadyExistsError)
		}
	}

	var uid uuid.UUID
	// Now let's first check if we have a repo cached.
	cachedId, err := IsCached(url, branch)
	if err != nil {
		return "", err
	}
	// If cachedId is nil, we need to clone.
	if cachedId != nil {
		uid = *cachedId
	} else {
		uid, err = cloneRepo(url, branch)
		if err != nil {
			return "", err
		}
	}

	currentConfig.Repositories[uid] = Repository{
		RepoURL: url,
		Branch:  branch,
	}

	name := MustGetRepoIndex(uid).Name
	if err := writeConfig(configPath); err != nil {
		return name, err
	}
	return name, nil
}

func cloneRepo(url, branch string) (uuid.UUID, error) {
	return cloneRepoWithExistingId(url, branch, nil)
}

func cloneRepoWithExistingId(url, branch string, id *uuid.UUID) (uuid.UUID, error) {
	var uid uuid.UUID
	if id != nil {
		uid = *id
	} else {
		uid = uuid.New()
	}

	repopath := filepath.Join(repoDir, uid.String())
	refName := plumbing.NewBranchReferenceName(branch)
	err := clone(ClonerConfig{
		CloneOptions: git.CloneOptions{
			URL:           url,
			ReferenceName: refName,
			SingleBranch:  true,
		},
		BaseDir: repopath,
	})
	if err != nil {
		return uid, err
	}

	// Get the name from the Index. If we get a PathError from GetRepoIndex here,
	// it means something is fishy with this repository. Probably not a kable repo.
	_, err = GetRepoIndex(uid)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.RemoveAll(repopath); err != nil {
				return uid, err
			}
			return uid, fmt.Errorf(RepositoryInvalidError)
		}
		return uid, err
	}

	// After a successful clone, we need to add the repo to the cache index.
	if err := AddToCacheIndex(uid, url, branch); err != nil {
		return uid, err
	}

	return uid, err
}

func (repos Repositories) RemoveRepository(name string) error {
	if !configSet() {
		return fmt.Errorf(ConfigNotInitializedError)
	}

	id := getRepoIdForName(name)

	if _, exists := currentConfig.Repositories[*id]; !exists {
		return fmt.Errorf(RepositoryNotExistsError)
	}

	delete(currentConfig.Repositories, *id)

	if err := writeConfig(configPath); err != nil {
		return err
	}
	return nil
}

func getRepoIdForName(name string) *uuid.UUID {
	for repoid, _ := range currentConfig.Repositories {
		if IsInitialized(repoid) {
			if MustGetRepoIndex(repoid).Name == name {
				return &repoid
			}
		}
	}
	return nil
}

func (repos Repositories) UpdateRepositories() error {
	for id, ref := range currentConfig.Repositories {

		if !IsInitialized(id) {
			_, err := cloneRepoWithExistingId(ref.RepoURL, ref.Branch, &id)
			if err != nil {
				return err
			}
			if err := AddToCacheIndex(id, ref.RepoURL, ref.Branch); err != nil {
				return err
			}

			continue
		}

		repo, err := git.PlainOpen(MustGetCachePath(id))
		if err != nil {
			return err
		}
		wt, err := repo.Worktree()
		if err != nil {
			return err
		}
		if err := wt.Pull(&git.PullOptions{
			SingleBranch: true,
		}); err != nil {
			if err == git.NoErrAlreadyUpToDate {
				continue
			}
			return err
		}
	}
	return nil
}

func IsInitialized(repoid uuid.UUID) bool {
	initialized := true
	_, err := GetRepoIndex(repoid)
	if err != nil {
		initialized = false
	}
	return initialized
}

func GetCachePath(repo uuid.UUID) (string, error) {
	ref, err := GetCacheInfo(repo)
	if err != nil {
		return "", err
	}
	return ref.Path, nil
}

func MustGetCachePath(repo uuid.UUID) string {
	out, _ := GetCachePath(repo)
	return out
}

func GetRepoIndex(repoid uuid.UUID) (RepoIndex, error) {
	index := RepoIndex{}
	repoPath, err := GetCachePath(repoid)
	if err != nil {
		return index, err
	}

	// Read in the file
	content, err := ioutil.ReadFile(filepath.Join(repoPath, RepoIndexFilename))
	if err != nil {
		return index, err
	}

	// Unmarshal the index file
	if err := json.Unmarshal(content, &index); err != nil {
		return index, err
	}
	return index, nil
}

func MustGetRepoIndex(repoid uuid.UUID) RepoIndex {
	out, _ := GetRepoIndex(repoid)
	return out
}
