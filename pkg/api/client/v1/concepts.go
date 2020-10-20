package v1

import (
	"context"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/redradrat/kable/pkg/api"
	"github.com/redradrat/kable/pkg/concepts"
)

const conceptsBasePath = "v1/concepts"

type ConceptsService interface {
	List(context.Context, *ListOptions) (*api.ConceptsPayload, *Response, error)
	ListByTag(context.Context, string, *ListOptions) (*api.ConceptsPayload, *Response, error)
	Get(context.Context, concepts.ConceptIdentifier) (*api.ConceptPayload, *Response, error)
	GetFromRepository(context.Context, string, string) (*api.ConceptPayload, *Response, error)
}

var _ ConceptsService = &ConceptsClient{}

type ConceptsClient struct {
	client *Client
}

func (c ConceptsClient) List(ctx context.Context, options *ListOptions) (*api.ConceptsPayload, *Response, error) {
	uri, err := url.Parse(conceptsBasePath)
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

func (c ConceptsClient) ListByTag(ctx context.Context, s string, options *ListOptions) (*api.ConceptsPayload, *Response, error) {
	panic("implement me")
}

func (c ConceptsClient) Get(ctx context.Context, identifier concepts.ConceptIdentifier) (*api.ConceptPayload, *Response, error) {
	panic("implement me")
}

func (c ConceptsClient) GetFromRepository(ctx context.Context, s string, s2 string) (*api.ConceptPayload, *Response, error) {
	panic("implement me")
}
