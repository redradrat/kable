package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/fatih/structs"
	"github.com/google/logger"
	"github.com/mitchellh/mapstructure"
	"github.com/redradrat/kable/pkg/errors"
	"github.com/spf13/viper"
)

const (
	StoreKey                       = "store"
	StoreTypeKey                   = "type"
	StoreConfigKey                 = "config"
	LocalStoreType StoreConfigType = "local"
	EtcdStoreType  StoreConfigType = "etcd"
)

type StoreConfigType string

func (t StoreConfigType) String() string {
	return string(t)
}

type StoreConfigMap map[string]interface{}

func GetStoreFromConfig() (Store, error) {
	m := viper.GetStringMap(StoreKey)
	return StoreConfigMap(m).GetStore()
}

func (cm StoreConfigMap) Map() map[string]interface{} {
	return cm
}

func (cm StoreConfigMap) GetStore() (Store, error) {
	var out Store
	switch cm[StoreTypeKey] {
	case LocalStoreType.String():
		out = LocalStore{}
	case EtcdStoreType.String():
		e := EtcdStore{}
		err := mapstructure.Decode(cm[StoreConfigKey], &e)
		if err != nil {
			return nil, err
		}
		out = e
	default:
		return nil, errors.InvalidStoreType
	}
	return out, nil
}

func LocalStoreConfigMap() StoreConfigMap {
	out := StoreConfigMap{}
	out[StoreTypeKey] = LocalStoreType.String()
	return out
}

// EtcdStoreConfigMap returns a config for Etcd with given endpoints and timeout in milliseconds
func EtcdStoreConfigMap(endpoints []string, timeout time.Duration) StoreConfigMap {
	out := StoreConfigMap{}
	out[StoreTypeKey] = EtcdStoreType.String()
	store := EtcdStore{}

	store.Endpoints = endpoints
	store.Timeout = timeout * time.Millisecond
	store.DialTimeout = 10000 * time.Millisecond
	out[StoreConfigKey] = structs.Map(store)
	return out
}

type Store interface {
	WriteRegistry(registry RepoRegistry) error
	ReadRegistry() (*RepoRegistry, error)
}

type LocalStore struct{}

func (l LocalStore) WriteRegistry(registry RepoRegistry) error {
	b, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(RepoRegistryPath, b, 0644); err != nil {
		return err
	}
	return nil
}

func (l LocalStore) ReadRegistry() (*RepoRegistry, error) {
	file, err := ioutil.ReadFile(RepoRegistryPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := l.WriteRegistry(RepoRegistry{}); err != nil {
				return nil, err
			}
			file, err = ioutil.ReadFile(RepoRegistryPath)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	r := RepoRegistry{}
	if err := json.Unmarshal(file, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

type EtcdStore struct {
	clientv3.Config
	Timeout time.Duration
}

const EtcdRegistryKey = "/kable/registry"

func (e EtcdStore) WriteRegistry(registry RepoRegistry) error {
	c, err := clientv3.New(e.Config)
	if err != nil {
		logger.Errorf("unable to create new etcd client: %v", err)
		return err
	}
	defer func() {
		err = c.Close()
	}()

	return e.writeRegistry(registry, c)
}

func (e EtcdStore) writeRegistry(registry RepoRegistry, c *clientv3.Client) error {
	marshalledRegistry, err := json.Marshal(registry)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout)
	resp, err := c.Put(ctx, EtcdRegistryKey, string(marshalledRegistry), clientv3.WithPrevKV())
	cancel()
	if err != nil {
		return err
	}

	logger.V(2).Info("Updated registry")
	logger.V(3).Info(fmt.Sprintf("New Registry: %s", registry.String()))
	if resp.PrevKv != nil {
		oldReg := RepoRegistry{}
		if err := json.Unmarshal(resp.PrevKv.Value, &oldReg); err != nil {
			return err
		}
		logger.V(3).Info(fmt.Sprintf("Previous Registry: %s", oldReg.String()))
	}

	return nil
}

func (e EtcdStore) ReadRegistry() (*RepoRegistry, error) {
	c, err := clientv3.New(e.Config)
	if err != nil {
		logger.Errorf("unable to create new etcd client: %v", err)
		return nil, err
	}
	defer func() {
		err = c.Close()
	}()

	return e.readRegistry(c)
}

func (e EtcdStore) readRegistry(c *clientv3.Client) (*RepoRegistry, error) {
	out := RepoRegistry{}
	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout)
	resp, err := c.Get(ctx, EtcdRegistryKey)
	cancel()
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) > 1 {
		logger.Errorf("detected multiple registries in store")
		return nil, errors.MultipleRegistriesInStoreError
	}
	for _, ev := range resp.Kvs {
		r := RepoRegistry{}
		if err := json.Unmarshal(ev.Value, &r); err != nil {
			logger.Errorf("unable to unmarshal registry from store: %v", err)
			return nil, err
		}
		out = r
	}

	return &out, nil
}
