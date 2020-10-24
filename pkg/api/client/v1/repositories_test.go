package v1

import (
	"context"
	"strings"
	"testing"

	"github.com/redradrat/kable/pkg/api"
	"github.com/redradrat/kable/pkg/repositories"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRepositoriesClient_List_None(t *testing.T) {
	viper.Set(repositories.StoreKey, repositories.MockStoreConfigMap().Map())
	client := NewClient(nil, &uri)

	repos, response, err := client.Repositories.List(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	exp := api.NewRepositoriesPayload()
	assert.Equal(t, exp, *repos)
}

func TestRepositoriesClient_List_All(t *testing.T) {
	removeMod, client, err := addDemoHttps(t)
	assert.NoError(t, err)

	repos, response, err := client.Repositories.List(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	exp := api.RepositoriesPayload(api.RepositoriesPayload{Repositories: api.RepositoriesMapPayload{"demo-https": api.RepositoryPayload{URL: "https://github.com/redradrat/kable", GitRef: "refs/heads/master"}}})
	assert.Equal(t, exp, *repos)

	assert.NoError(t, repositories.UpdateRegistry(removeMod))
}

func TestRepositoriesClient_Get(t *testing.T) {
	removeMod, client, err := addDemoHttps(t)
	assert.NoError(t, err)

	concepts, response, err := client.Repositories.Get(context.Background(), repositories.DemoHttpsRepository.Name)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	exp := api.RepositoryPayload(api.RepositoryPayload{URL: strings.TrimSuffix(repositories.DemoHttpsRepository.URL, ".git"), GitRef: repositories.DemoHttpsRepository.GitRef})
	assert.Equal(t, exp, *concepts)

	concepts, response, err = client.Repositories.Get(context.Background(), "dummyname")
	assert.Error(t, err)
	assert.Equal(t, 404, response.StatusCode)

	concepts, response, err = client.Repositories.Get(nil, repositories.DemoHttpsRepository.Name)
	assert.Error(t, err)

	assert.NoError(t, repositories.UpdateRegistry(removeMod))
}

func TestRepositoriesClient_Delete(t *testing.T) {
	_, client, err := addDemoHttps(t)
	assert.NoError(t, err)

	response, err := client.Repositories.Delete(context.Background(), repositories.DemoHttpsRepository.Name)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)

	_, response, err = client.Repositories.Get(context.Background(), repositories.DemoHttpsRepository.Name)
	assert.Error(t, err)
	assert.Equal(t, 404, response.StatusCode)

}

func TestRepositoriesClient_Put(t *testing.T) {
	viper.Set(repositories.StoreKey, repositories.MockStoreConfigMap().Map())
	client := NewClient(nil, &uri)

	response, err := client.Repositories.Put(context.Background(), repositories.DemoHttpsRepository.Name, repositories.DemoHttpsRepository.URL, repositories.DemoHttpsRepository.GitRef)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)

	concepts, response, err := client.Repositories.Get(context.Background(), repositories.DemoHttpsRepository.Name)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	exp := api.RepositoryPayload(api.RepositoryPayload{URL: strings.TrimSuffix(repositories.DemoHttpsRepository.URL, ".git"), GitRef: repositories.DemoHttpsRepository.GitRef})
	assert.Equal(t, exp, *concepts)
}
