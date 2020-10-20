package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func RegisterHandlersV1(e *echo.Group, serv *Serv) {
	e.GET("/concepts", serv.GetConcepts)
	e.GET("/concepts/:id", serv.GetConcept)
	e.GET("/repositories", serv.GetRepositories)
	e.GET("/repositories/:id", serv.GetRepository)
	e.PUT("/repositories/:id", serv.PutRepository)
	e.GET("/repositories/:id/concepts", serv.GetRepositoryConcepts)
	e.GET("/repositories/:id/concepts/:path", serv.GetRepositoryConcept)
}

type MessagePayload struct {
	Message string `json:"message"`
}

func NewMessage(msg string, e ...interface{}) MessagePayload {
	return MessagePayload{Message: fmt.Sprintf(msg, e...)}
}

type RepositoriesPayload struct {
	Repositories RepositoriesMapPayload `json:"repositories"`
}

type RepositoriesMapPayload map[string]RepositoryPayload

func NewRepositoriesPayload() RepositoriesPayload {
	return RepositoriesPayload{
		Repositories: map[string]RepositoryPayload{},
	}
}

type RepositoryPayload struct {
	URL    string `json:"url"`
	GitRef string `json:"gitRef,omitempty"`
}

type ConceptsPayload struct {
	Concepts ConceptsMapPayload `json:"concepts"`
}

type ConceptsMapPayload map[string]ConceptPayload

func NewConceptsPayload() ConceptsPayload {
	return ConceptsPayload{Concepts: map[string]ConceptPayload{}}
}

type ConceptPayload struct {
	Type     string                 `json:"type"`
	Metadata ConceptMetadataPayload `json:"metadata"`
	Inputs   []ConceptInputsPayload `json:"inputs,omitempty"`
}

type ConceptMetadataPayload struct {
	Maintainer ConceptMaintainerPayload `json:"maintainer,omitempty"`
	Tags       map[string]string        `json:"tags,omitempty"`
}

type ConceptMaintainerPayload struct {
	MaintainerName  string `json:"name"`
	MaintainerEmail string `json:"email"`
}

type ConceptInputsPayload struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Mandatory bool   `json:"mandatory"`
}
