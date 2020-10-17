package repositories

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/go-git/go-git/v5"

	"github.com/redradrat/kable/pkg/errors"
)

const (
	RegistryFileName = "kableconfig.json"
	KableDirName     = ".kable"
	CacheDirName     = "cache"
)

func homeDir() string {
	out, err := os.UserHomeDir()
	if err != nil {
		panic(err)
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

type Repositories map[string]Repository
type Auths map[string]Auth

func (repos Repositories) List() []Repository {
	out := make([]Repository, 0, len(repos))
	for _, v := range repos {
		out = append(out, v)
	}
	return out
}

func Registry() (RepoRegistry, error) {
	file, err := ioutil.ReadFile(RepoRegistryPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := writeRegistry(RepoRegistry{}); err != nil {
				return RepoRegistry{}, err
			}
			file, err = ioutil.ReadFile(RepoRegistryPath)
			if err != nil {
				return RepoRegistry{}, err
			}
		} else {
			return RepoRegistry{}, err
		}
	}

	r := RepoRegistry{}
	if err := json.Unmarshal(file, &r); err != nil {
		return RepoRegistry{}, err
	}

	return r, nil
}

type RegistryModification func(RepoRegistry) RepoRegistry

func AddRepository(repo Repository) RegistryModification {
	return func(registry RepoRegistry) RepoRegistry {
		if registry.Repositories == nil {
			registry.Repositories = Repositories{}
		}
		registry.Repositories[repo.Name] = repo
		repo.URL = strings.TrimSuffix(repo.URL, ".git")
		return registry
	}
}

func RemoveRepository(name string) RegistryModification {
	return func(registry RepoRegistry) RepoRegistry {
		if registry.Repositories == nil {
			return registry
		}
		delete(registry.Repositories, name)
		return registry
	}
}

func UpdateRegistry(updates ...RegistryModification) error {
	registry, err := Registry()
	if err != nil {
		return err
	}
	for _, update := range updates {
		registry = update(registry)
	}

	if err := writeRegistry(registry); err != nil {
		return err
	}

	return nil
}

func writeRegistry(r RepoRegistry) error {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(RepoRegistryPath, b, 0644); err != nil {
		return err
	}
	return nil
}

type Auth struct {
	Basic string `json:"basic,omitempty"`
}

type Repository struct {
	GitRepository
	Name string
}

func maybeDelete(path string) error {
	if err := os.RemoveAll(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

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
		if err != nil {
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
	a := strings.Split(strings.TrimSuffix(url, ".git"), "/")
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

	b, err := ioutil.ReadFile(filepath.Join(path, RegistryFileName))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &ri); err != nil {
		return nil, err
	}

	return &ri, nil
}

type GitRepository struct {
	URL    string
	GitRef string
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
		registry.Auths[url] = Auth{Basic: base64.StdEncoding.EncodeToString(b)}
		return registry
	}, nil
}
