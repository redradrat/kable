package v1

import (
	"context"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/redradrat/kable/pkg/api"
)

const repositoriessBasePath = "v1/repositories"

type RepositoriesService interface {
	List(ctx context.Context, options *ListOptions) (*api.RepositoriesPayload, *Response, error)
	Get(ctx context.Context, name string) (*api.RepositoryPayload, *Response, error)
	Put(ctx context.Context, name, giturl, gitref string) (*Response, error)
	Delete(ctx context.Context, name string) (*Response, error)
}

var _ RepositoriesService = &RepositoriesClient{}

type RepositoriesClient struct {
	client *Client
}

func (r RepositoriesClient) List(ctx context.Context, options *ListOptions) (*api.RepositoriesPayload, *Response, error) {
	uri, err := url.Parse(repositoriessBasePath)
	if err != nil {
		return nil, nil, err
	}

	vals, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	uri.RawQuery = vals.Encode()

	req, err := r.client.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	payload := api.NewRepositoriesPayload()
	resp, err := r.client.Do(ctx, req, &payload)
	if err != nil {
		return nil, resp, err
	}

	return &payload, resp, err
}

func (r RepositoriesClient) Get(ctx context.Context, name string) (*api.RepositoryPayload, *Response, error) {
	uri, err := url.Parse(repositoriessBasePath + "/")
	if err != nil {
		return nil, nil, err
	}

	getUri, err := uri.Parse(name)
	if err != nil {
		return nil, nil, err
	}

	req, err := r.client.NewRequest(http.MethodGet, getUri.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	payload := api.RepositoryPayload{}
	resp, err := r.client.Do(ctx, req, &payload)
	if err != nil {
		return nil, resp, err
	}

	return &payload, resp, err
}

func (r RepositoriesClient) Delete(ctx context.Context, name string) (*Response, error) {
	uri, err := url.Parse(repositoriessBasePath + "/")
	if err != nil {
		return nil, err
	}

	getUri, err := uri.Parse(name)
	if err != nil {
		return nil, err
	}

	req, err := r.client.NewRequest(http.MethodDelete, getUri.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (r RepositoriesClient) Put(ctx context.Context, name, giturl, gitref string) (*Response, error) {
	uri, err := url.Parse(repositoriessBasePath + "/")
	if err != nil {
		return nil, err
	}

	getUri, err := uri.Parse(name)
	if err != nil {
		return nil, err
	}

	payload := api.RepositoryPayload{
		URL:    giturl,
		GitRef: gitref,
	}

	req, err := r.client.NewRequest(http.MethodPut, getUri.String(), payload)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
