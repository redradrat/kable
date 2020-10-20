package v1

import (
	"context"

	"github.com/redradrat/kable/pkg/api"
	"github.com/redradrat/kable/pkg/concepts"
)

const repositoriessBasePath = "v1/repositories"

type RepositoriesService interface {
	List(context.Context, *ListOptions) (api.ConceptsPayload, *Response, error)
	ListByTag(context.Context, string, *ListOptions) (api.ConceptsPayload, *Response, error)
	Get(context.Context, concepts.ConceptIdentifier) (api.ConceptPayload, *Response, error)
	GetFromRepository(context.Context, string, string) (api.ConceptPayload, *Response, error)
	Delete(context.Context, int) (*Response, error)
	DeleteByTag(context.Context, string) (*Response, error)
}

var _ RepositoriesService = &RepositoriesClient{}

type RepositoriesClient struct {
	*Client
}

func (r RepositoriesClient) List(ctx context.Context, options *ListOptions) (api.ConceptsPayload, *Response, error) {
	panic("implement me")
}

func (r RepositoriesClient) ListByTag(ctx context.Context, s string, options *ListOptions) (api.ConceptsPayload, *Response, error) {
	panic("implement me")
}

func (r RepositoriesClient) Get(ctx context.Context, identifier concepts.ConceptIdentifier) (api.ConceptPayload, *Response, error) {
	panic("implement me")
}

func (r RepositoriesClient) GetFromRepository(ctx context.Context, s string, s2 string) (api.ConceptPayload, *Response, error) {
	panic("implement me")
}

func (r RepositoriesClient) Delete(ctx context.Context, i int) (*Response, error) {
	panic("implement me")
}

func (r RepositoriesClient) DeleteByTag(ctx context.Context, s string) (*Response, error) {
	panic("implement me")
}
