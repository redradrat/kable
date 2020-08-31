package kable

import (
	"encoding/json"
	"errors"
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

type RepoIndex struct {
	Version        int      `json:"version"`
	ConceptEntries []string `json:"concepts"`
}

func ListRepositories() ([][]interface{}, error) {
	idx, err := readCacheIndex()
	if err != nil {
		return nil, err
	}

	var repoSlices [][]interface{}
	for id, repo := range idx.Index {
		if repo.SoftDeleted {
			continue
		}
		if IsInitialized(id) {
			repoSlices = append(repoSlices, []interface{}{id, repo.URI, true})
		} else {
			repoSlices = append(repoSlices, []interface{}{id, repo.URI, false})
		}
	}
	return repoSlices, nil
}

type ClonerConfig struct {
	git.CloneOptions
	BaseDir string
}

type RepoCache struct {
	Index    map[string]RepoCacheEntry
	modCount uuid.UUID
}

func NewRepoCache() RepoCache {
	return RepoCache{
		Index:    map[string]RepoCacheEntry{},
		modCount: uuid.UUID{},
	}
}

type RepoCacheEntry struct {
	Path        string
	Branch      string
	URI         string
	SoftDeleted bool
}

func writeCacheIndex(ci RepoCache) error {
	return writeCacheIndexI(ci, false)
}

func initCacheIndex() error {
	rc := NewRepoCache()
	return writeCacheIndexI(rc, true)
}

func writeCacheIndexI(ci RepoCache, init bool) error {
	if !init {
		oldIdx, err := readCacheIndex()
		if err != nil {
			return err
		}

		if ci.modCount != oldIdx.modCount {
			return StaleRepoCacheIndexError
		}
	}

	// Before writing, we update the modCount.
	ci.modCount = uuid.New()

	out, err := json.MarshalIndent(ci, "", "	")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(userDir, "repo.cache"), out, 0666); err != nil {
		return err
	}

	return nil
}

func readCacheIndex() (RepoCache, error) {
	index := NewRepoCache()
	content, err := ioutil.ReadFile(filepath.Join(userDir, "repo.cache"))
	if err != nil {
		if !os.IsNotExist(err) {
			return index, err
		}
		if err := initCacheIndex(); err != nil {
			return index, err
		}
	}
	if err := json.Unmarshal(content, &index); err != nil {
		return index, err
	}
	return index, nil
}

func AddToCacheIndex(id string, url, path, branch string) error {
	cache, err := readCacheIndex()

	// If we get an "IsNotExist"-error here, we just assume we're initializing
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	cache.Index[id] = RepoCacheEntry{
		Branch: branch,
		Path:   path,
		URI:    url,
	}

	if err := writeCacheIndex(cache); err != nil {
		return err
	}

	return nil
}

func GetCacheInfo(id string) (*RepoCacheEntry, error) {
	ci, err := readCacheIndex()
	if err != nil {
		return nil, err
	}
	out := ci.Index[id]
	return &out, nil
}

func DeactivateRepoCache(id string) error {
	for {
		cache, err := readCacheIndex()
		if err != nil {
			return err
		}

		entry := cache.Index[id]
		entry.SoftDeleted = true
		cache.Index[id] = entry

		if err := writeCacheIndex(cache); err != nil {
			if errors.Is(err, StaleRepoCacheIndexError) {
				continue
			}
			return err
		}

		return nil
	}

}

func ActivateRepoCache(oldid, newid string) error {
	for {
		cache, err := readCacheIndex()
		if err != nil {
			return err
		}

		entry := cache.Index[oldid]
		entry.SoftDeleted = false
		cache.Index[newid] = entry

		if oldid != newid {
			delete(cache.Index, oldid)
			if err := os.Rename(filepath.Join(repoDir, oldid), filepath.Join(repoDir, newid)); err != nil {
				return err
			}
		}

		if err := writeCacheIndex(cache); err != nil {
			if errors.Is(err, StaleRepoCacheIndexError) {
				continue
			}
			return err
		}

		return nil
	}
}

func RemoveFromCacheIndex(id string) error {
	for {
		cache, err := readCacheIndex()
		if err != nil {
			return err
		}

		delete(cache.Index, id)
		if err := writeCacheIndex(cache); err != nil {
			if errors.Is(err, StaleRepoCacheIndexError) {
				continue
			}
			return err
		}

		return nil
	}
}

// Returns a UUID if cached, else just returns nil.
func IsCached(url, branch string) (*string, error) {
	var id string
	idx, err := readCacheIndex()
	if err != nil {
		return &id, err
	}
	for existingId, ref := range idx.Index {
		if ref.URI == url && ref.Branch == branch {
			id = existingId
			return &id, nil
		}
	}
	return nil, nil

}

// TidyRepositories cleans up all cached repositories that are not in use in the current config,
// and returns the names of the deleted ones.
func TidyRepositories() error {
	idx, err := readCacheIndex()
	if err != nil {
		return err
	}

	// Go through index, and delete all repositories, that are not used in the current config
	fi, err := ioutil.ReadDir(repoDir)
	if err != nil {
		return err
	}

	for _, ref := range fi {
		if _, exists := idx.Index[ref.Name()]; !exists || idx.Index[ref.Name()].SoftDeleted == true {
			if err := os.RemoveAll(filepath.Join(repoDir, ref.Name())); err != nil {
				return err
			}
			if err := RemoveFromCacheIndex(ref.Name()); err != nil {
				return err
			}
		}
	}

	return nil
}

func AddRepository(id, url, branch string) error {
	idx, err := readCacheIndex()
	if err != nil {
		return err
	}

	if _, exists := idx.Index[id]; exists && !idx.Index[id].SoftDeleted {
		return RepositoryAlreadyExistsError
	}

	// Now let's first check if we have a repo cached.
	cachedId, err := IsCached(url, branch)
	if err != nil {
		return err
	}

	if cachedId != nil {
		if err := ActivateRepoCache(*cachedId, id); err != nil {
			return err
		}
		// Exit without additional cloning.
		return nil
	}

	// If cachedId is nil, we need to clone.
	if err := cloneRepo(id, url, branch); err != nil {
		return err
	}

	return nil
}

func RemoveRepository(id string) error {
	return DeactivateRepoCache(id)
}

func UpdateRepositories() error {
	idx, err := readCacheIndex()
	if err != nil {
		return err
	}

	for id, ref := range idx.Index {
		if !IsInitialized(id) {
			if err := cloneRepo(id, ref.URI, ref.Branch); err != nil {
				return err
			}
			continue
		}

		repo, err := git.PlainOpen(MustGetCacheInfo(id).Path)
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

func cloneRepo(name, url, branch string) error {
	repopath := filepath.Join(repoDir, name)
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
		return err
	}

	// Get the name from the Index. If we get a PathError from GetRepoIndex here,
	// it means something is fishy with this repository. Probably not a kable repo.
	cont, err := ioutil.ReadFile(filepath.Join(repopath, RepoIndexFilename))
	if err != nil {
		if err := os.RemoveAll(repopath); err != nil {
			return err
		}
		if os.IsNotExist(err) {
			return RepositoryInvalidError
		}
		return err
	}

	ri := RepoIndex{}
	if err := json.Unmarshal(cont, &ri); err != nil {
		if err := os.RemoveAll(repopath); err != nil {
			return err
		}
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return RepositoryInvalidError
		}
		return err
	}

	// After a successful clone, we need to add the repo to the cache index.
	if err := AddToCacheIndex(name, url, repopath, branch); err != nil {
		return err
	}

	return nil
}

func IsInitialized(repoid string) bool {
	initialized := true
	_, err := GetRepoIndex(repoid)
	if err != nil {
		initialized = false
	}
	return initialized
}

func MustGetCacheInfo(id string) *RepoCacheEntry {
	out, _ := GetCacheInfo(id)
	return out
}

func GetRepoIndex(repoid string) (RepoIndex, error) {
	index := RepoIndex{}
	cacheInfo, err := GetCacheInfo(repoid)
	if err != nil {
		return index, err
	}

	// Read in the file
	content, err := ioutil.ReadFile(filepath.Join(cacheInfo.Path, RepoIndexFilename))
	if err != nil {
		return index, err
	}

	// Unmarshal the index file
	if err := json.Unmarshal(content, &index); err != nil {
		return index, err
	}
	return index, nil
}

func MustGetRepoIndex(repoid string) RepoIndex {
	out, _ := GetRepoIndex(repoid)
	return out
}
