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
	assert.Equal(t, response.StatusCode, 200)
	assert.Equal(t, *concepts, api.NewConceptsPayload())
}

func TestConceptsClient_List_Some(t *testing.T) {
	viper.Set(repositories.StoreKey, repositories.MockStoreConfigMap().Map())
	client := NewClient(nil, &uri)

	// Add repo to mock store
	addMod, err := repositories.AddRepository(repositories.DemoHttpsRepository)
	assert.NoError(t, err)
	assert.NoError(t, repositories.UpdateRegistry(addMod))

	concepts, response, err := client.Concepts.List(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, 200)
	assert.Equal(t, *concepts, api.NewConceptsPayload())
}
