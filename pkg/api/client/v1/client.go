package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/coreos/pkg/httputil"

	"github.com/redradrat/kable/cmd"
)

var userAgent = "kable-client/" + cmd.CliVersion

// ListOptions defines paramaters to be used with List Operations
type ListOptions struct {
	// Which page to list.
	Page int `url:"page,omitempty"`

	// Results per page.
	Limit int `url:"limit,omitempty"`
}

type Response struct {
	*http.Response
}

type ErrorResponse struct {
	Response *http.Response

	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%v => %v: %d %v",
		e.Response.Request.Method, e.Response.Request.URL, e.Response.StatusCode, e.Message)
}

type Client struct {
	client  *http.Client
	BaseURL *url.URL

	Concepts     ConceptsService
	Repositories RepositoriesService
}

func NewClient(c *http.Client, baseURL *url.URL) *Client {
	if c == nil {
		c = http.DefaultClient
	}

	client := &Client{
		client:  c,
		BaseURL: baseURL,
	}
	client.Concepts = ConceptsClient{client}
	client.Repositories = RepositoriesClient{client}

	return client
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", httputil.JSONContentType)
	req.Header.Add("Accept", httputil.JSONContentType)
	req.Header.Set("User-Agent", userAgent)
	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, w interface{}) (*Response, error) {
	if ctx == nil {
		return nil, ErrorResponse{Message: fmt.Sprintf("context cannot be nil")}
	}

	req = req.WithContext(ctx)
	hresp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if w != nil {
		err = json.NewDecoder(hresp.Body).Decode(w)
		if err != nil {
			return nil, err
		}
	}

	return UnwrapResponse(hresp), UnwrapError(hresp)
}

func UnwrapResponse(r *http.Response) *Response {
	return &Response{r}
}

func UnwrapError(r *http.Response) error {
	s := r.StatusCode
	if s >= 200 && s < 300 {
		return nil
	}

	er := &ErrorResponse{
		Response: r,
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return er
	}
	er.Message = string(data)

	return er
}
