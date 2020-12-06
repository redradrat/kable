package api

import (
	"fmt"

	"github.com/redradrat/kable/pkg/concepts"
)

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

func NewConceptPayload() ConceptPayload {
	return ConceptPayload{
		Inputs: []ConceptInputsPayload{},
	}
}

type ConceptMetadataPayload struct {
	Maintainer ConceptMaintainerPayload `json:"maintainer,omitempty"`
	Tags       []string                 `json:"tags,omitempty"`
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

type RenderConceptInputPayload struct {
	TargetType     string                 `json:"type"`
	SingleManifest bool                   `json:"singleManifest"`
	Values         *concepts.RenderValues `json:"values"`
}

type RenderConceptResultPayload struct {
	Manifests     []string                `json:"manifests"`
	ManifestCount int                     `json:"manifestCount"`
	Origin        *concepts.ConceptOrigin `json:"origin"`
}

type ByID []ConceptInputsPayload

func (a ByID) Len() int           { return len(a) }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
