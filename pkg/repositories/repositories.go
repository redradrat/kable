package repositories

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/go-git/go-git/v5"

	"github.com/redradrat/kable/pkg/errors"

	"github.com/zalando/go-keyring"
)

const (
	RegistryFileName = "kable.json"
	KableDirName     = ".kable"
	CacheDirName     = "cache"
	KableKeychainId  = "kablerepo_"
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

//var kableKeyringConfig = keyring.Config{
//	ServiceName: "kable",
//}

type RepoRegistry struct {
	Repositories Repositories `json:"repositories"`
}

type Repositories map[string]Repository

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

type Repository struct {
	GitRepository
	Name     string
	Username *string
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
	var auth *http.BasicAuth
	if r.Username != nil {
		a, err := getAuth(r.Name, *r.Username)
		if err != nil {
			return err
		}
		auth = &http.BasicAuth{
			Username: a.Username,
			Password: a.Password,
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
	Username string
	Password string
}

func getAuth(name, user string) (*AuthPair, error) {

	pass, err := keyring.Get(KableKeychainId+name, user)
	if err != nil {
		return nil, err
	}

	//ring, err := getKeyring()
	//if err != nil {
	//	return nil, err
	//}
	//var pass keyring.Item
	//user, err := ring.Get(name + "-user")
	//if err != nil {
	//	if err == keyring.ErrKeyNotFound {
	//		return nil, err
	//	}
	//	return nil, err
	//} else {
	//	pass, err = ring.Get(name + "-pass")
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	return &AuthPair{Username: user, Password: pass}, nil
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

func StoreRepoAuth(name string, pair AuthPair) error {
	// set password
	err := keyring.Set(KableKeychainId+name, pair.Username, pair.Password)
	if err != nil {
		return err
	}

	// store credentials in local keyring
	//err = ring.Set(keyring.Item{
	//	Key:   name + "-user",
	//	Data:  []byte(pair.Username),
	//	Label: "kable",
	//})
	//err = ring.Set(keyring.Item{
	//	Key:   name + "-pass",
	//	Data:  []byte(pair.Password),
	//	Label: "kable",
	//})
	//if err != nil {
	//	return err
	//}
	return nil
}

//var ring *keyring.Keyring
//
//func getKeyring() (keyring.Keyring, error) {
//	if ring != nil {
//		return *ring, nil
//	}
//	r, err := keyring.Open(kableKeyringConfig)
//	if err != nil {
//		return nil, err
//	}
//	ring = &r
//	return *ring, nil
//}

//
//type ClonerConfig struct {
//	git.CloneOptions
//	BaseDir string
//}
//
//type RepoCache struct {
//	Index    map[string]RepoCacheEntry
//	modCount uuid.UUID
//}
//
//func NewRepoCache() RepoCache {
//	return RepoCache{
//		Index:    map[string]RepoCacheEntry{},
//		modCount: uuid.UUID{},
//	}
//}
//
//type RepoCacheEntry struct {
//	RepoDirPath string
//	Branch      string
//	URI         string
//	SoftDeleted bool
//}
//
//func (rce RepoCacheEntry) AbsolutePath() string {
//	return filepath.Join(config.RepoDir, rce.RepoDirPath)
//}
//
//func writeCacheIndex(ci RepoCache) error {
//	return writeCacheIndexI(ci, false)
//}
//
//func initCacheIndex() error {
//	rc := NewRepoCache()
//	return writeCacheIndexI(rc, true)
//}
//
//func writeCacheIndexI(ci RepoCache, init bool) error {
//	if !init {
//		oldIdx, err := readCacheIndex()
//		if err != nil {
//			return err
//		}
//
//		if ci.modCount != oldIdx.modCount {
//			return errors.StaleRepoCacheIndexError
//		}
//	}
//
//	// Before writing, we update the modCount.
//	ci.modCount = uuid.New()
//
//	out, err := json.MarshalIndent(ci, "", "	")
//	if err != nil {
//		return err
//	}
//	if err := ioutil.WriteFile(filepath.Join(config.UserDir, "repo.cache"), out, 0666); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func readCacheIndex() (RepoCache, error) {
//	index := NewRepoCache()
//	content, err := ioutil.ReadFile(filepath.Join(config.UserDir, "repo.cache"))
//	if err != nil {
//		if !os.IsNotExist(err) {
//			return index, err
//		}
//		if err := initCacheIndex(); err != nil {
//			return index, err
//		}
//	}
//	if err := json.Unmarshal(content, &index); err != nil {
//		return index, err
//	}
//	return index, nil
//}
//
//func AddToCacheIndex(id string, url, branch string) error {
//	cache, err := readCacheIndex()
//
//	// If we get an "IsNotExist"-error here, we just assume we're initializing
//	if err != nil && !os.IsNotExist(err) {
//		return err
//	}
//
//	cache.Index[id] = RepoCacheEntry{
//		Branch:      branch,
//		RepoDirPath: id,
//		URI:         url,
//	}
//
//	if err := writeCacheIndex(cache); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func GetCacheInfo(id string) (*RepoCacheEntry, error) {
//	ci, err := readCacheIndex()
//	if err != nil {
//		return nil, err
//	}
//	out := ci.Index[id]
//	return &out, nil
//}
//
//func DeactivateRepoCache(id string) error {
//	for {
//		cache, err := readCacheIndex()
//		if err != nil {
//			return err
//		}
//
//		entry := cache.Index[id]
//		entry.SoftDeleted = true
//		cache.Index[id] = entry
//
//		if err := writeCacheIndex(cache); err != nil {
//			if errors.Is(err, errors.StaleRepoCacheIndexError) {
//				continue
//			}
//			return err
//		}
//
//		return nil
//	}
//
//}
//
//func ActivateRepoCache(oldid, newid string) error {
//	for {
//		cache, err := readCacheIndex()
//		if err != nil {
//			return err
//		}
//
//		entry := cache.Index[oldid]
//		entry.SoftDeleted = false
//		cache.Index[newid] = entry
//
//		if oldid != newid {
//			delete(cache.Index, oldid)
//			if err := os.Rename(filepath.Join(config.RepoDir, oldid), filepath.Join(config.RepoDir, newid)); err != nil {
//				return err
//			}
//		}
//
//		if err := writeCacheIndex(cache); err != nil {
//			if errors.Is(err, errors.StaleRepoCacheIndexError) {
//				continue
//			}
//			return err
//		}
//
//		return nil
//	}
//}
//
//func RemoveFromCacheIndex(id string) error {
//	for {
//		cache, err := readCacheIndex()
//		if err != nil {
//			return err
//		}
//
//		delete(cache.Index, id)
//		if err := writeCacheIndex(cache); err != nil {
//			if errors.Is(err, errors.StaleRepoCacheIndexError) {
//				continue
//			}
//			return err
//		}
//
//		return nil
//	}
//}
//
//// Returns a UUID if cached, else just returns nil.
//func IsCached(url, branch string) (*string, error) {
//	var id string
//	idx, err := readCacheIndex()
//	if err != nil {
//		return &id, err
//	}
//	for existingId, ref := range idx.Index {
//		if ref.URI == url && ref.Branch == branch {
//			id = existingId
//			return &id, nil
//		}
//	}
//	return nil, nil
//
//}
//
//// TidyRepositories cleans up all cached repositories that are not in use in the current config,
//// and returns the names of the deleted ones.
//func TidyRepositories() error {
//	idx, err := readCacheIndex()
//	if err != nil {
//		return err
//	}
//
//	// Go through index, and delete all repositories, that are not used in the current config
//	fi, err := ioutil.ReadDir(config.RepoDir)
//	if err != nil {
//		return err
//	}
//
//	for _, ref := range fi {
//		if _, exists := idx.Index[ref.Name()]; !exists || idx.Index[ref.Name()].SoftDeleted == true {
//			if err := os.RemoveAll(filepath.Join(config.RepoDir, ref.Name())); err != nil {
//				return err
//			}
//			if err := RemoveFromCacheIndex(ref.Name()); err != nil {
//				return err
//			}
//		}
//	}
//
//	return nil
//}
//
//func AddRepository(id, url, branch string) error {
//	return AddAuthRepository(id, url, "", "", branch)
//}
//
//func AddAuthRepository(id, url, user, pass, branch string) error {
//	idx, err := readCacheIndex()
//	if err != nil {
//		return err
//	}
//
//	if _, exists := idx.Index[id]; exists && !idx.Index[id].SoftDeleted {
//		return errors.RepositoryAlreadyExistsError
//	}
//
//	// Now let's first check if we have a repo cached.
//	cachedId, err := IsCached(url, branch)
//	if err != nil {
//		return err
//	}
//
//	if cachedId != nil {
//		if err := ActivateRepoCache(*cachedId, id); err != nil {
//			return err
//		}
//		// Exit without additional cloning.
//		return nil
//	}
//
//	// If cachedId is nil, we need to maybeClone.
//	if err := cloneRepo(id, url, user, pass, branch); err != nil {
//		return err
//	}
//
//	ring, err := keyring.Open(kableKeyringConfig)
//	if err != nil {
//		return err
//	}
//
//	// store credentials in local keyring
//	err = ring.Set(keyring.Item{
//		Key:   id + "-user",
//		Data:  []byte(user),
//		Label: "kable",
//	})
//	err = ring.Set(keyring.Item{
//		Key:   id + "-pass",
//		Data:  []byte(pass),
//		Label: "kable",
//	})
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func RemoveRepository(id string) error {
//	return DeactivateRepoCache(id)
//}
//
//func UpdateRepositories() error {
//	idx, err := readCacheIndex()
//	if err != nil {
//		return err
//	}
//
//	for id, _ := range idx.Index {
//		if !IsInitialized(id) {
//			continue
//		}
//
//		repo, err := git.PlainOpen(MustGetCacheInfo(id).AbsolutePath())
//		if err != nil {
//			return err
//		}
//		wt, err := repo.Worktree()
//		if err != nil {
//			return err
//		}
//		if err := wt.Pull(&git.PullOptions{
//			SingleBranch: true,
//		}); err != nil {
//			if err.Error() == "authentication required" {
//				// try to get auth from keychain
//				ring, err := keyring.Open(kableKeyringConfig)
//				if err != nil {
//					return err
//				}
//				user, err := ring.Get(id + "-user")
//				if err != nil {
//					return err
//				}
//				pass, err := ring.Get(id + "-pass")
//				if err != nil {
//					return err
//				}
//				err = wt.Pull(&git.PullOptions{
//					SingleBranch: true,
//					Auth: &http.BasicAuth{
//						Username: string(user.Data),
//						Password: string(pass.Data),
//					},
//				})
//				if err != nil {
//					if err == git.NoErrAlreadyUpToDate {
//						continue
//					}
//					return err
//				}
//			}
//		}
//	}
//	return nil
//}
//
//func cloneRepo(name, url, user, pass, branch string) error {
//	repopath := filepath.Join(config.RepoDir, name)
//	refName := plumbing.NewBranchReferenceName(branch)
//	var auth transport.AuthMethod
//	if user != "" && pass != "" {
//		auth = &http.BasicAuth{
//			Username: user,
//			Password: pass,
//		}
//	}
//
//	err := maybeClone(ClonerConfig{
//		CloneOptions: git.CloneOptions{
//			Auth:          auth,
//			URL:           url,
//			ReferenceName: refName,
//			SingleBranch:  true,
//		},
//		BaseDir: repopath,
//	})
//	if err != nil {
//		return err
//	}
//
//	// Get the name from the Index. If we get a PathError from GetRepoIndex here,
//	// it means something is fishy with this repository. Probably not a kable repo.
//	cont, err := ioutil.ReadFile(filepath.Join(repopath, RegistryFileName))
//	if err != nil {
//		if err := os.RemoveAll(repopath); err != nil {
//			return err
//		}
//		if os.IsNotExist(err) {
//			return errors.RepositoryInvalidError
//		}
//		return err
//	}
//
//	ri := RepoIndex{}
//	if err := json.Unmarshal(cont, &ri); err != nil {
//		if err := os.RemoveAll(repopath); err != nil {
//			return err
//		}
//		if _, ok := err.(*json.UnmarshalTypeError); ok {
//			return errors.RepositoryInvalidError
//		}
//		return err
//	}
//
//	// After a successful maybeClone, we need to add the repo to the cache index.
//	if err := AddToCacheIndex(name, url, branch); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func IsInitialized(repoid string) bool {
//	initialized := true
//	_, err := GetRepoIndex(repoid)
//	if err != nil {
//		initialized = false
//	}
//	return initialized
//}
//
//func MustGetCacheInfo(id string) *RepoCacheEntry {
//	out, _ := GetCacheInfo(id)
//	return out
//}
//
//func GetRepoIndex(repoid string) (RepoIndex, error) {
//	index := RepoIndex{}
//	cacheInfo, err := GetCacheInfo(repoid)
//	if err != nil {
//		return index, err
//	}
//
//	// Read in the file
//	content, err := ioutil.ReadFile(filepath.Join(cacheInfo.AbsolutePath(), RegistryFileName))
//	if err != nil {
//		return index, err
//	}
//
//	// Unmarshal the index file
//	if err := json.Unmarshal(content, &index); err != nil {
//		return index, err
//	}
//	return index, nil
//}
//
//func MustGetRepoIndex(repoid string) RepoIndex {
//	out, _ := GetRepoIndex(repoid)
//	return out
//}
