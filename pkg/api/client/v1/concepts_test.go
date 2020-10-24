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

func TestConceptsClient_List_All(t *testing.T) {
	removeMod, client, err := addDemoHttps(t)
	assert.NoError(t, err)

	concepts, response, err := client.Concepts.List(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	exp := api.ConceptsPayload{Concepts: api.ConceptsMapPayload{"e2e-test/testconcept1@demo-https": api.ConceptPayload{Type: "jsonnet", Metadata: api.ConceptMetadataPayload{Maintainer: api.ConceptMaintainerPayload{MaintainerName: "Demo Maintainer", MaintainerEmail: "demo.maintainer@kab.le"}, Tags: []string(nil)}, Inputs: []api.ConceptInputsPayload{api.ConceptInputsPayload{ID: "instanceName", Type: "string", Mandatory: true}, api.ConceptInputsPayload{ID: "nameSelection", Type: "select", Mandatory: true}}}, "e2e-test/testconcept2@demo-https": api.ConceptPayload{Type: "jsonnet", Metadata: api.ConceptMetadataPayload{Maintainer: api.ConceptMaintainerPayload{MaintainerName: "Demo Maintainer", MaintainerEmail: "demo.maintainer@kab.le"}, Tags: []string(nil)}, Inputs: []api.ConceptInputsPayload{api.ConceptInputsPayload{ID: "instanceName", Type: "string", Mandatory: true}, api.ConceptInputsPayload{ID: "nameSelection", Type: "select", Mandatory: true}}}}}
	assert.Equal(t, exp, *concepts)

	assert.NoError(t, repositories.UpdateRegistry(removeMod))
}

func TestConceptsClient_ListFromRepository_All(t *testing.T) {
	removeMod, client, err := addDemoHttps(t)
	assert.NoError(t, err)

	concepts, response, err := client.Concepts.ListFromRepository(context.Background(), repositories.DemoHttpsRepository.Name, nil)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	exp := api.ConceptsPayload{Concepts: api.ConceptsMapPayload{"e2e-test/testconcept1@demo-https": api.ConceptPayload{Type: "jsonnet", Metadata: api.ConceptMetadataPayload{Maintainer: api.ConceptMaintainerPayload{MaintainerName: "Demo Maintainer", MaintainerEmail: "demo.maintainer@kab.le"}, Tags: []string(nil)}, Inputs: []api.ConceptInputsPayload{api.ConceptInputsPayload{ID: "instanceName", Type: "string", Mandatory: true}, api.ConceptInputsPayload{ID: "nameSelection", Type: "select", Mandatory: true}}}, "e2e-test/testconcept2@demo-https": api.ConceptPayload{Type: "jsonnet", Metadata: api.ConceptMetadataPayload{Maintainer: api.ConceptMaintainerPayload{MaintainerName: "Demo Maintainer", MaintainerEmail: "demo.maintainer@kab.le"}, Tags: []string(nil)}, Inputs: []api.ConceptInputsPayload{api.ConceptInputsPayload{ID: "instanceName", Type: "string", Mandatory: true}, api.ConceptInputsPayload{ID: "nameSelection", Type: "select", Mandatory: true}}}}}
	assert.Equal(t, exp, *concepts)

	concepts, response, err = client.Concepts.ListFromRepository(context.Background(), "dummyname", nil)
	assert.Error(t, err)
	assert.Equal(t, 404, response.StatusCode)

	assert.NoError(t, repositories.UpdateRegistry(removeMod))
}

func TestConceptsClient_GetFromRepository_All(t *testing.T) {
	removeMod, client, err := addDemoHttps(t)
	assert.NoError(t, err)

	concepts, response, err := client.Concepts.GetFromRepository(context.Background(), repositories.DemoHttpsRepository.Name, "e2e-test/testconcept1")
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	exp := api.ConceptPayload{Type: "jsonnet", Metadata: api.ConceptMetadataPayload{Maintainer: api.ConceptMaintainerPayload{MaintainerName: "Demo Maintainer", MaintainerEmail: "demo.maintainer@kab.le"}, Tags: []string(nil)}, Inputs: []api.ConceptInputsPayload{api.ConceptInputsPayload{ID: "instanceName", Type: "string", Mandatory: true}, api.ConceptInputsPayload{ID: "nameSelection", Type: "select", Mandatory: true}}}
	assert.Equal(t, exp, *concepts)

	concepts, response, err = client.Concepts.GetFromRepository(context.Background(), repositories.DemoHttpsRepository.Name, "e2e-test/testconcept")
	assert.Error(t, err)
	assert.Equal(t, 404, response.StatusCode)

	concepts, response, err = client.Concepts.GetFromRepository(context.Background(), "dummyname", "e2e-test/testconcept1")
	assert.Error(t, err)
	assert.Equal(t, 404, response.StatusCode)

	assert.NoError(t, repositories.UpdateRegistry(removeMod))
}

func TestConceptsClient_Get(t *testing.T) {
	removeMod, client, err := addDemoHttps(t)
	assert.NoError(t, err)

	concepts, response, err := client.Concepts.Get(context.Background(), "e2e-test/testconcept1@demo-https")
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	exp := api.ConceptPayload{Type: "jsonnet", Metadata: api.ConceptMetadataPayload{Maintainer: api.ConceptMaintainerPayload{MaintainerName: "Demo Maintainer", MaintainerEmail: "demo.maintainer@kab.le"}, Tags: []string(nil)}, Inputs: []api.ConceptInputsPayload{api.ConceptInputsPayload{ID: "instanceName", Type: "string", Mandatory: true}, api.ConceptInputsPayload{ID: "nameSelection", Type: "select", Mandatory: true}}}
	assert.Equal(t, exp, *concepts)

	concepts, response, err = client.Concepts.Get(context.Background(), "e2e-test/testconcept@demo-https")
	assert.Error(t, err)
	assert.Equal(t, 404, response.StatusCode)

	concepts, response, err = client.Concepts.Get(nil, "e2e-test/testconcept@demo-https")
	assert.Error(t, err)

	assert.NoError(t, repositories.UpdateRegistry(removeMod))
}

/////////////
// HELPERS //
/////////////

func addDemoHttps(t *testing.T) (repositories.RegistryModification, *Client, error) {
	viper.Set(repositories.StoreKey, repositories.MockStoreConfigMap().Map())
	client := NewClient(nil, &uri)

	// Add repo to mock store
	addMod, err := repositories.AddRepository(repositories.DemoHttpsRepository)
	assert.NoError(t, err)
	assert.NoError(t, repositories.UpdateRegistry(addMod))
	removeMod, err := repositories.RemoveRepository(repositories.DemoHttpsRepository.Name)
	assert.NoError(t, err)
	return removeMod, client, err
}
