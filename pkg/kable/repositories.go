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
	RepoIndexFilename = "kable.json"
)

type CloneIndex map[uuid.UUID]CloneRef

type CloneRef struct {
	Path   string
	Branch string
	URI    string
}

func AddToIndex(id uuid.UUID, path, branch, url string) error {
	ci := CloneIndex{}
	cipath := filepath.Join(repoDir, "cloneindex.json")
	content, _ := ioutil.ReadFile(cipath)
	_ = json.Unmarshal(content, &ci)

	ci[id] = CloneRef{
		Branch: branch,
		Path:   path,
		URI:    url,
	}
	out, err := json.MarshalIndent(ci, "", "    ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(cipath, out, os.ModePerm); err != nil {
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

func GetFromIndex(id uuid.UUID) (*CloneRef, error) {
	ci, err := GetCacheIndex()
	if err != nil {
		return nil, err
	}
	out := ci[id]
	return &out, nil
}

func RemoveFromIndex(id uuid.UUID) error {
	ci := CloneIndex{}
	cipath := filepath.Join(repoDir, "cloneindex.json")
	content, _ := ioutil.ReadFile(cipath)
	_ = json.Unmarshal(content, &ci)

	delete(ci, id)
	out, err := json.MarshalIndent(ci, "", "    ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(cipath, out, os.ModePerm); err != nil {
		return err
	}
	return nil
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
		if repo, exists := currentConfig.Repositories[id]; !exists {
			if err := os.RemoveAll(ref.Path); err != nil {
				return repolist, err
			}
			if err := RemoveFromIndex(id); err != nil {
				return repolist, err
			}
			repolist = append(repolist, repo.MustGetIndex().Name)
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
			name, err := repo.GetRepoIndex()
			if err != nil {
				return "", err
			}
			return name.Name, fmt.Errorf(RepositoryAlreadyExistsError)
		}
	}

	uid, repopath, err := cloneRepo(url, branch)
	if err != nil {
		return "", err
	}

	repoIndex := RepoIndex{}
	content, err := ioutil.ReadFile(filepath.Join(repopath, RepoIndexFilename))
	if err != nil {
		return repoIndex.Name, err
	}
	if err := json.Unmarshal(content, &repoIndex); err != nil {
		return repoIndex.Name, err
	}
	if err := AddToIndex(uid, repopath, branch, url); err != nil {
		return repoIndex.Name, err
	}

	currentConfig.Repositories[uid] = Repository{
		RepoURL: url,
		Branch:  branch,
	}
	if err := writeConfig(configPath); err != nil {
		return repoIndex.Name, err
	}
	return repoIndex.Name, nil
}

func cloneRepo(url, branch string) (uuid.UUID, string, error) {
	return cloneRepoWithId(url, branch, nil)
}

func cloneRepoWithId(url, branch string, id *uuid.UUID) (uuid.UUID, string, error) {
	var uid uuid.UUID
	if id != nil {
		uid = *id
	} else {
		uid = uuid.New()
	}
	idx, err := GetCacheIndex()
	if err != nil {
		return uid, "", err
	}
	for existingId, ref := range idx {
		if ref.URI == url && ref.Branch == branch {
			return existingId, ref.Path, nil
		}
	}

	repopath := filepath.Join(repoDir, uid.String())
	refName := plumbing.NewBranchReferenceName(branch)
	err = clone(ClonerConfig{
		CloneOptions: git.CloneOptions{
			URL:           url,
			ReferenceName: refName,
			SingleBranch:  true,
		},
		BaseDir: repopath,
	})
	if err != nil {
		return uid, repopath, err
	}
	return uid, repopath, err
}

func (repos Repositories) RemoveRepository(name string) error {
	if !configSet() {
		return fmt.Errorf(ConfigNotInitializedError)
	}

	var id uuid.UUID
	for repoid, ref := range currentConfig.Repositories {
		if ref.MustGetIndex().Name == name {
			id = repoid
		}
	}

	if _, exists := currentConfig.Repositories[id]; !exists {
		return fmt.Errorf(RepositoryNotExistsError)
	}

	delete(currentConfig.Repositories, id)

	if err := writeConfig(configPath); err != nil {
		return err
	}
	return nil
}

func (repos Repositories) UpdateRepositories() error {
	for id, ref := range currentConfig.Repositories {
		cloneRef, err := GetFromIndex(id)
		if err != nil {
			return err
		}
		repo, err := git.PlainOpen(cloneRef.Path)
		if err != nil {
			// If the repo is not cached yet, let's clone fresh
			if err == git.ErrRepositoryNotExists {
				_, repopath, err := cloneRepoWithId(ref.RepoURL, ref.Branch, &id)
				if err != nil {
					return err
				}
				if err := AddToIndex(id, repopath, ref.Branch, ref.RepoURL); err != nil {
					return err
				}

				continue
			}
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

func (r Repository) GetRepoIndex() (RepoIndex, error) {
	var repoPath string
	index := RepoIndex{}
	idx, err := GetCacheIndex()
	if err != nil {
		return RepoIndex{}, err
	}
	for _, ref := range idx {
		if ref.URI == r.RepoURL && ref.Branch == r.Branch {
			repoPath = ref.Path
		}
	}
	content, err := ioutil.ReadFile(filepath.Join(repoPath, RepoIndexFilename))
	if err != nil {
		return index, err
	}
	if err := json.Unmarshal(content, &index); err != nil {
		return index, err
	}
	return index, nil
}

func (r Repository) MustGetIndex() RepoIndex {
	out, _ := r.GetRepoIndex()
	return out
}
