package api

import "fmt"

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
	Metadata ConceptMetadataPayload `json:"maintainer,omitempty"`
	Inputs   []ConceptInputsPayload `json:"inputs,omitempty"`
}

type ConceptMetadataPayload struct {
	MaintainerName  string `json:"name"`
	MaintainerEmail string `json:"email"`
}

type ConceptInputsPayload struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Mandatory bool   `json:"mandatory"`
}
