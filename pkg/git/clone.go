package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)


type Cloner struct {
	RepoURL string
	Auth    *http.BasicAuth
}

func (cl Cloner) Checkout(path string) error {
	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:  cl.RepoURL,
		Auth: cl.Auth,
	})
	if err != nil {
		return err
	}
	return nil
}

type ClonerAuth interface {
	FormURI(baseuri string) (string, error)
}

