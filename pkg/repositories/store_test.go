package repositories

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"

	"github.com/spf13/viper"

	"github.com/labstack/gommon/random"
)

func TestGetRepository(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    Repository
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRepository(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRepository() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStoreFromConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    Store
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStoreFromConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStoreFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStoreFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListRepositories(t *testing.T) {
	tests := []struct {
		name    string
		want    []Repository
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListRepositories()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListRepositories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListRepositories() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalStoreConfigMap(t *testing.T) {
	tests := []struct {
		name string
		want StoreConfigMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LocalStoreConfigMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocalStoreConfigMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalStore_ReadRegistry(t *testing.T) {
	tests := []struct {
		name    string
		want    *RepoRegistry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := LocalStore{}
			got, err := l.ReadRegistry()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadRegistry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadRegistry() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalStore_WriteRegistry(t *testing.T) {
	type args struct {
		registry RepoRegistry
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := LocalStore{}
			if err := l.WriteRegistry(tt.args.registry); (err != nil) != tt.wantErr {
				t.Errorf("WriteRegistry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegistry(t *testing.T) {
	tests := []struct {
		name    string
		want    *RepoRegistry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Registry()
			if (err != nil) != tt.wantErr {
				t.Errorf("Registry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Registry() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepoAuthExists(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RepoAuthExists(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("RepoAuthExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RepoAuthExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepoRegistry_String(t *testing.T) {
	type fields struct {
		Repositories Repositories
		Auths        Auths
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := RepoRegistry{
				Repositories: tt.fields.Repositories,
				Auths:        tt.fields.Auths,
			}
			if got := registry.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepositories_List(t *testing.T) {
	tests := []struct {
		name  string
		repos Repositories
		want  []Repository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.repos.List(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_AbsolutePath(t *testing.T) {
	type fields struct {
		GitRepository GitRepository
		Name          string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{
				GitRepository: tt.fields.GitRepository,
				Name:          tt.fields.Name,
			}
			got, err := r.AbsolutePath()
			if (err != nil) != tt.wantErr {
				t.Errorf("AbsolutePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AbsolutePath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_RepoIndex(t *testing.T) {
	type fields struct {
		GitRepository GitRepository
		Name          string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *RepoIndex
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{
				GitRepository: tt.fields.GitRepository,
				Name:          tt.fields.Name,
			}
			got, err := r.RepoIndex()
			if (err != nil) != tt.wantErr {
				t.Errorf("RepoIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RepoIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreConfigMap_GetStore(t *testing.T) {
	tests := []struct {
		name    string
		cm      StoreConfigMap
		want    Store
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cm.GetStore()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStore() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreConfigMap_Map(t *testing.T) {
	tests := []struct {
		name string
		cm   StoreConfigMap
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cm.Map(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreConfigType_String(t *testing.T) {
	tests := []struct {
		name string
		t    StoreConfigType
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreRepoAuth(t *testing.T) {
	type args struct {
		url  string
		pair AuthPair
	}
	tests := []struct {
		name    string
		args    args
		want    RegistryModification
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StoreRepoAuth(tt.args.url, tt.args.pair)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreRepoAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StoreRepoAuth() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTidyCache(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := TidyCache(); (err != nil) != tt.wantErr {
				t.Errorf("TidyCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateRegistry(t *testing.T) {
	type args struct {
		updates []RegistryModification
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateRegistry(tt.args.updates...); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRegistry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateRepositories(t *testing.T) {
	tests := []struct {
		name       string
		localstore bool
		wantErr    bool
	}{
		{name: "nostore", localstore: false, wantErr: true},
		{name: "local store", localstore: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.localstore {
				viper.Set(StoreKey, LocalStoreConfigMap().Map())
			}
			if err := UpdateRepositories(); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRepositories() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_computePath(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "https", args: struct{ url string }{url: DemoHttpsUrl}, want: filepath.Join(CacheDir, "kable")},
		{name: "ssh", args: struct{ url string }{url: DemoSshUrl}, want: filepath.Join(CacheDir, "kable")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := computePath(tt.args.url); got != tt.want {
				t.Errorf("computePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_homeDir(t *testing.T) {
	dir, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name string
		want string
	}{
		{name: "get homdir", want: dir},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := homeDir(); got != tt.want {
				t.Errorf("homeDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_maybeClone(t *testing.T) {
	type args struct {
		r          Repository
		path       string
		localstore bool
		pull       bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "nostore", args: struct {
			r          Repository
			path       string
			localstore bool
			pull       bool
		}{r: DemoHttpsRepository, path: randomPath(), localstore: false, pull: true}, wantErr: true},
		{name: "https", args: struct {
			r          Repository
			path       string
			localstore bool
			pull       bool
		}{r: DemoHttpsRepository, path: randomPath(), localstore: true, pull: true}, wantErr: false},
		{name: "ssh", args: struct {
			r          Repository
			path       string
			localstore bool
			pull       bool
		}{r: DemoSshRepository, path: randomPath(), localstore: true, pull: true}, wantErr: false},
		{name: "not allowed", args: struct {
			r          Repository
			path       string
			localstore bool
			pull       bool
		}{r: DemoHttpsRepository, path: filepath.Join("/var", random.String(24, random.Lowercase)), localstore: true, pull: true}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.localstore {
				viper.Set(StoreKey, LocalStoreConfigMap().Map())
			} else {
				viper.Set(StoreKey, "")
			}
			if err := maybeClone(tt.args.r, tt.args.path, tt.args.pull); (err != nil) != tt.wantErr {
				t.Errorf("maybeClone() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := os.RemoveAll(tt.args.path); err != nil {
				t.Error(err)
			}
		})
	}
}

func randomPath() string {
	return filepath.Join(os.TempDir(), random.String(24, random.Lowercase))
}

func Test_safeDelete(t *testing.T) {
	pathGood := randomPath()
	if err := os.MkdirAll(pathGood, os.ModePerm); err != nil {
		t.Error(err)
	}
	pathdir := randomPath()
	if err := os.MkdirAll(pathdir, os.ModeDir); err != nil {
		t.Error(err)
	}

	pathdelete := func(path string) {
		if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
	}
	defer pathdelete(pathGood)
	defer pathdelete(pathdir)

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "delete", args: struct{ path string }{path: randomPath()}, wantErr: false},
		{name: "delete", args: struct{ path string }{path: pathGood}, wantErr: false},
		{name: "delete", args: struct{ path string }{path: pathdir}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := safeDelete(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("safeDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var testTrimUrl = "https://github.com/redradrat/demo-concepts.git"

func Test_trimUrl(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "trimcorrect", args: struct{ url string }{url: testTrimUrl}, want: strings.TrimSuffix(testTrimUrl, ".git")},
		{name: "trimunneccessary", args: struct{ url string }{url: strings.TrimSuffix(testTrimUrl, ".git")}, want: strings.TrimSuffix(testTrimUrl, ".git")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimUrl(tt.args.url); got != tt.want {
				t.Errorf("trimUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unpointerify(t *testing.T) {
	type args struct {
		i *RepoRegistry
		e error
	}
	tests := []struct {
		name    string
		args    args
		want    RepoRegistry
		wantErr bool
	}{
		{name: "noerror", args: struct {
			i *RepoRegistry
			e error
		}{i: &DemoRegistry, e: nil}, want: DemoRegistry, wantErr: false},
		{name: "error", args: struct {
			i *RepoRegistry
			e error
		}{i: &DemoRegistry, e: testPtrError}, want: DemoRegistry, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unpointerify(tt.args.i, tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpointerify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpointerify() got = %v, want %v", got, tt.want)
			}
		})
	}
}

var testPtrError = errors.New("test")

func TestEtcdStore_readRegistry(t *testing.T) {

	type fields struct {
		Timeout time.Duration
		Config  clientv3.Config
	}
	type args struct {
		c *clientv3.Client
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *RepoRegistry
		wantErr bool
	}{
		// Tests TBD
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EtcdStore{
				Timeout: tt.fields.Timeout,
				Config:  tt.fields.Config,
			}
			got, err := e.readRegistry(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("readRegistry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readRegistry() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEtcdStore_writeRegistry(t *testing.T) {
	type fields struct {
		Timeout time.Duration
		Config  clientv3.Config
	}
	type args struct {
		registry RepoRegistry
		c        *clientv3.Client
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EtcdStore{
				Timeout: tt.fields.Timeout,
				Config:  tt.fields.Config,
			}
			if err := e.writeRegistry(tt.args.registry, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("writeRegistry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
