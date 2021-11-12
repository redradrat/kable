package repositories

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/go-git/go-git/v5"

	"github.com/redradrat/kable/pkg/errors"
)

const (
	RegistryFileName          = "kableconfig.json"
	RepoIndexFileName         = "kable.json"
	KableDirName              = ".kable"
	CacheDirName              = "cache"
	masterGitRef              = "refs/heads/master"
	RepositoryIdentifierRegex = "^([a-z\\-]+)$"
)

var IsValidRepositoryName = regexp.MustCompile(RepositoryIdentifierRegex).MatchString

func homeDir() string {
	out, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}
	return out
}

var KableDir = filepath.Join(homeDir(), KableDirName)
var RepoRegistryPath = filepath.Join(KableDir, RegistryFileName)
var CacheDir = filepath.Join(KableDir, CacheDirName)

type RepoRegistry struct {
	Repositories Repositories `json:"repositories"`
	Auths        Auths        `json:"auths,omitempty"`
}

func (registry RepoRegistry) String() string {
	var repostrings []string
	for name, repo := range registry.Repositories {
		repostrings = append(repostrings, fmt.Sprintf("'%s => %s@%s'", name, repo.GitRef, repo.URL))
	}
	return fmt.Sprintf("Registry[Repos: %s, Auths: %s]", strings.Join(repostrings, ", "), strconv.Itoa(len(registry.Auths)))
}

type Repositories map[string]Repository
type Auths map[string]Auth

func (repos Repositories) List() []Repository {
	out := make([]Repository, 0, len(repos))
	for _, v := range repos {
		out = append(out, v)
	}
	return out
}

func Registry() (*RepoRegistry, error) {
	store, err := GetStoreFromConfig()
	if err != nil {
		return nil, err
	}
	return store.ReadRegistry()
}

type RegistryModification func(RepoRegistry) RepoRegistry

func AddRepository(repo Repository) (RegistryModification, error) {
	if repo.URL == "" {
		return nil, fmt.Errorf("repository must have a valid url")
	}
	if !IsValidRepositoryName(repo.Name) {
		return nil, fmt.Errorf("repository name can only be lowercase letters (a-z) and '-'")
	}
	addrepofunc := func(registry RepoRegistry) RepoRegistry {
		if registry.Repositories == nil {
			registry.Repositories = Repositories{}
		}
		if repo.GitRef == "" {
			repo.GitRef = masterGitRef
		}
		repo.URL = trimUrl(repo.URL)
		registry.Repositories[repo.Name] = repo
		return registry
	}
	return addrepofunc, nil
}

func RemoveRepository(name string) (RegistryModification, error) {
	if !IsValidRepositoryName(name) {
		return nil, fmt.Errorf("repository name can only be lowercase letters (a-z) and '-'")
	}
	removeFunc := func(registry RepoRegistry) RepoRegistry {
		if registry.Repositories == nil {
			return registry
		}
		delete(registry.Repositories, name)
		return registry
	}
	return removeFunc, nil
}

func UpdateRegistry(updates ...RegistryModification) error {
	var registry RepoRegistry
	ptrRegistry, err := Registry()
	if err != nil {
		return err
	}
	for _, update := range updates {
		registry = update(*ptrRegistry)
	}

	store, err := GetStoreFromConfig()
	if err != nil {
		return err
	}
	if err := store.WriteRegistry(registry); err != nil {
		return err
	}

	return nil
}

type Auth struct {
	Basic string `json:"basic,omitempty"`
}

type Repository struct {
	GitRepository
	Name string `json:"name"`
}

func safeDelete(path string) error {
	if err := os.RemoveAll(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

// Clones the given repository, or, in case the repo already is checked
// out, pulls the upstream changes.
func maybeClone(r Repository, path string, pull bool) error {
	reg, err := Registry()
	if err != nil {
		return err
	}
	var auth *http.BasicAuth
	if val, ok := reg.Auths[r.URL]; ok {
		var pair AuthPair
		b, err := base64.StdEncoding.DecodeString(val.Basic)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(b, &pair); err != nil {
			return err
		}
		auth = &http.BasicAuth{
			Username: pair.Username,
			Password: pair.Password,
		}
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			cloneopts := &git.CloneOptions{
				URL:           r.URL,
				ReferenceName: plumbing.ReferenceName(r.GitRef),
				SingleBranch:  true,
				Auth:          auth,
			}

			_, err = git.PlainClone(path, false, cloneopts)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	if pull {
		wt, err := repo.Worktree()
		if err != nil {
			return err
		}

		if err := wt.Pull(&git.PullOptions{
			ReferenceName: plumbing.ReferenceName(r.GitRef),
			SingleBranch:  true,
			Auth:          auth,
		}); err != nil && err != git.NoErrAlreadyUpToDate {
			return err
		}
	}

	return nil
}

func (r Repository) AbsolutePath() (string, error) {
	path := computePath(r.URL)
	if err := maybeClone(r, path, false); err != nil {
		return "", err
	}
	return path, nil
}

func computePath(url string) string {
	a := strings.Split(trimUrl(url), "/")
	return filepath.Join(CacheDir, a[len(a)-1])
}

type AuthPair struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r Repository) RepoIndex() (*RepoIndex, error) {
	ri := RepoIndex{}
	path, err := r.AbsolutePath()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(filepath.Join(path, RepoIndexFileName))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &ri); err != nil {
		return nil, err
	}

	return &ri, nil
}

type GitRepository struct {
	URL    string `json:"url"`
	GitRef string `json:"gitRef"`
}

func UpdateRepositories() error {
	repos, err := ListRepositories()
	if err != nil {
		return err
	}
	for _, repo := range repos {
		if err := maybeClone(repo, computePath(repo.URL), true); err != nil {
			return err
		}
	}
	return nil
}

func ListRepositories() ([]Repository, error) {
	r, err := Registry()
	if err != nil {
		return nil, err
	}
	return r.Repositories.List(), nil
}

func TidyCache() error {
	files, err := ioutil.ReadDir(CacheDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		_, err := GetRepository(file.Name())
		if err != nil {
			if err := os.RemoveAll(filepath.Join(CacheDir, file.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}

func GetRepository(name string) (Repository, error) {
	r, err := Registry()
	if err != nil {
		return Repository{}, err
	}
	repo, ok := r.Repositories[name]
	if !ok {
		return Repository{}, errors.RepositoryUnknownError
	}

	return repo, nil
}

type RepoIndex struct {
	Version        int      `json:"version"`
	ConceptEntries []string `json:"concepts"`
}

func StoreRepoAuth(url string, pair AuthPair) (RegistryModification, error) {
	b, err := json.Marshal(pair)
	if err != nil {
		return nil, err
	}
	return func(registry RepoRegistry) RepoRegistry {
		if registry.Auths == nil {
			registry.Auths = Auths{}
		}
		registry.Auths[trimUrl(url)] = Auth{Basic: base64.StdEncoding.EncodeToString(b)}
		return registry
	}, nil
}

func RepoAuthExists(url string) (bool, error) {
	reg, err := Registry()
	if err != nil {
		return false, err
	}
	if _, ok := reg.Auths[trimUrl(url)]; ok {
		return true, nil
	}
	return false, nil
}

func trimUrl(url string) string {
	return strings.TrimSuffix(url, ".git")
}
