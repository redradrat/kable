package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/redradrat/kable/pkg/api"
	"github.com/redradrat/kable/pkg/concepts"
)

const (
	conceptsBasePath       = "v1" + api.ConceptsApiPath
	conceptsFromFormatPath = repositoriessBasePath + "/%s" + api.ConceptsApiPath
)

type ConceptsService interface {
	Get(context.Context, concepts.ConceptIdentifier) (*api.ConceptPayload, *Response, error)
	GetFromRepository(context.Context, string, concepts.ConceptIdentifier) (*api.ConceptPayload, *Response, error)
	List(context.Context, *ListOptions) (*api.ConceptsPayload, *Response, error)
	ListFromRepository(context.Context, string, *ListOptions) (*api.ConceptsPayload, *Response, error)
	ListByTag(context.Context, string, *ListOptions) (*api.ConceptsPayload, *Response, error)
}

var _ ConceptsService = &ConceptsClient{}

type ConceptsClient struct {
	client *Client
}

func (c ConceptsClient) Get(ctx context.Context, identifier concepts.ConceptIdentifier) (*api.ConceptPayload, *Response, error) {
	basePath := conceptsBasePath
	return c.getWithBasepath(basePath, ctx, identifier)
}

func (c ConceptsClient) GetFromRepository(ctx context.Context, repo string, identifier concepts.ConceptIdentifier) (*api.ConceptPayload, *Response, error) {
	basePath := fmt.Sprintf(conceptsFromFormatPath, repo)
	return c.getWithBasepath(basePath, ctx, identifier)
}

func (c ConceptsClient) List(ctx context.Context, options *ListOptions) (*api.ConceptsPayload, *Response, error) {
	basePath := conceptsBasePath
	return c.listWithBasepath(basePath, ctx, options)
}

func (c ConceptsClient) ListByTag(ctx context.Context, s string, options *ListOptions) (*api.ConceptsPayload, *Response, error) {
	payload, resp, err := c.List(ctx, options)
	if err != nil {
		return nil, resp, err
	}

	filteredPayload := *payload

	for path, concept := range filteredPayload.Concepts {
		if !hasEntry(concept.Metadata.Tags, s) {
			delete(filteredPayload.Concepts, path)
		}
	}

	return &filteredPayload, resp, err
}

func hasEntry(l []string, s string) bool {
	hasEntry := false
	for _, tag := range l {
		if tag == s {
			hasEntry = true
		}
	}
	return hasEntry
}

func (c ConceptsClient) ListFromRepository(ctx context.Context, repo string, options *ListOptions) (*api.ConceptsPayload, *Response, error) {
	basePath := fmt.Sprintf(conceptsFromFormatPath, repo)
	return c.listWithBasepath(basePath, ctx, options)
}

func (c ConceptsClient) getWithBasepath(basePath string, ctx context.Context, identifier concepts.ConceptIdentifier) (*api.ConceptPayload, *Response, error) {
	uri, err := url.Parse(basePath + "/")
	if err != nil {
		return nil, nil, err
	}

	getUri, err := uri.Parse(api.MarshalId(identifier.String()))
	if err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest(http.MethodGet, getUri.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	payload := api.NewConceptPayload()
	resp, err := c.client.Do(ctx, req, &payload)
	if err != nil {
		return nil, resp, err
	}

	return &payload, resp, err
}

func (c ConceptsClient) listWithBasepath(basePath string, ctx context.Context, options *ListOptions) (*api.ConceptsPayload, *Response, error) {
	uri, err := url.Parse(basePath)
	if err != nil {
		return nil, nil, err
	}

	vals, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	uri.RawQuery = vals.Encode()

	req, err := c.client.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	payload := api.NewConceptsPayload()
	resp, err := c.client.Do(ctx, req, &payload)
	if err != nil {
		return nil, resp, err
	}

	return &payload, resp, err
}
