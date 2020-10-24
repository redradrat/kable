package v1

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/redradrat/kable/pkg/repositories"
	"github.com/spf13/viper"

	"github.com/redradrat/kable/pkg/api"
)

const base = "127.0.0.1:31111"

var uri = url.URL{
	Scheme: "http",
	Host:   base,
}

func init() {
	fmt.Println("running mock server for test!")
	go api.StartUp(base)
	time.Sleep(1 * time.Second)
}

func TestConceptsClient_List_None(t *testing.T) {
	viper.Set(repositories.StoreKey, repositories.MockStoreConfigMap().Map())
	client := NewClient(nil, &uri)

	concepts, response, err := client.Concepts.List(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, api.NewConceptsPayload(), *concepts)
}

func TestConceptsClient_List_Some(t *testing.T) {
	client, err := addDemoHttps(t)
	assert.NoError(t, err)

	concepts, response, err := client.Concepts.List(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	act := api.ConceptsPayload{Concepts: api.ConceptsMapPayload{"apps/grafana@demo-https": api.ConceptPayload{Type: "jsonnet", Metadata: api.ConceptMetadataPayload{Maintainer: api.ConceptMaintainerPayload{MaintainerName: "Ralph Kühnert", MaintainerEmail: "kuehnert.ralph@gmail.com"}, Tags: []string(nil)}, Inputs: []api.ConceptInputsPayload{api.ConceptInputsPayload{ID: "instanceName", Type: "string", Mandatory: true}, api.ConceptInputsPayload{ID: "nameSelection", Type: "select", Mandatory: true}}}, "apps/sentry@demo-https": api.ConceptPayload{Type: "jsonnet", Metadata: api.ConceptMetadataPayload{Maintainer: api.ConceptMaintainerPayload{MaintainerName: "Ralph Kühnert", MaintainerEmail: "kuehnert.ralph@gmail.com"}, Tags: []string(nil)}, Inputs: []api.ConceptInputsPayload{api.ConceptInputsPayload{ID: "instanceName", Type: "string", Mandatory: true}, api.ConceptInputsPayload{ID: "nameSelection", Type: "select", Mandatory: true}}}}}
	assert.Equal(t, act, *concepts)
}

func TestConceptsClient_ListFromRepository_Some(t *testing.T) {
	client, err := addDemoHttps(t)
	assert.NoError(t, err)

	concepts, response, err := client.Concepts.ListFromRepository(context.Background(), repositories.DemoHttpsRepository.Name, nil)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	act := api.ConceptsPayload{Concepts: api.ConceptsMapPayload{"apps/grafana@demo-https": api.ConceptPayload{Type: "jsonnet", Metadata: api.ConceptMetadataPayload{Maintainer: api.ConceptMaintainerPayload{MaintainerName: "Ralph Kühnert", MaintainerEmail: "kuehnert.ralph@gmail.com"}, Tags: []string(nil)}, Inputs: []api.ConceptInputsPayload{api.ConceptInputsPayload{ID: "instanceName", Type: "string", Mandatory: true}, api.ConceptInputsPayload{ID: "nameSelection", Type: "select", Mandatory: true}}}, "apps/sentry@demo-https": api.ConceptPayload{Type: "jsonnet", Metadata: api.ConceptMetadataPayload{Maintainer: api.ConceptMaintainerPayload{MaintainerName: "Ralph Kühnert", MaintainerEmail: "kuehnert.ralph@gmail.com"}, Tags: []string(nil)}, Inputs: []api.ConceptInputsPayload{api.ConceptInputsPayload{ID: "instanceName", Type: "string", Mandatory: true}, api.ConceptInputsPayload{ID: "nameSelection", Type: "select", Mandatory: true}}}}}
	assert.Equal(t, act, *concepts)
}

func TestConceptsClient_GetFromRepository_Some(t *testing.T) {
	client, err := addDemoHttps(t)
	assert.NoError(t, err)

	concepts, response, err := client.Concepts.GetFromRepository(context.Background(), repositories.DemoHttpsRepository.Name, "apps/sentry")
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	act := api.ConceptPayload(api.ConceptPayload{Type: "jsonnet", Metadata: api.ConceptMetadataPayload{Maintainer: api.ConceptMaintainerPayload{MaintainerName: "Ralph Kühnert", MaintainerEmail: "kuehnert.ralph@gmail.com"}, Tags: []string(nil)}, Inputs: []api.ConceptInputsPayload{api.ConceptInputsPayload{ID: "instanceName", Type: "string", Mandatory: true}, api.ConceptInputsPayload{ID: "nameSelection", Type: "select", Mandatory: true}}})
	assert.Equal(t, act, *concepts)
}

/////////////
// HELPERS //
/////////////

func addDemoHttps(t *testing.T) (*Client, error) {
	viper.Set(repositories.StoreKey, repositories.MockStoreConfigMap().Map())
	client := NewClient(nil, &uri)

	// Add repo to mock store
	addMod, err := repositories.AddRepository(repositories.DemoHttpsRepository)
	assert.NoError(t, err)
	assert.NoError(t, repositories.UpdateRegistry(addMod))
	return client, err
}
