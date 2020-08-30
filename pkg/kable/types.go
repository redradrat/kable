package kable

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
)

type APIVersion string
type Kind string

// Concept defines model for Concept.
type Concept struct {
	APIVersion APIVersion  `json:"apiVersion"`
	Meta       ConceptMeta `json:"meta"`
	Values     *struct {
		Mandatory *ValueList `json:"mandatory,omitempty"`
		Optional  *ValueList `json:"optional,omitempty"`
	} `json:"values,omitempty"`
}

// ConceptMeta defines model for ConceptMeta.
type ConceptMeta struct {
	Group *string `json:"group,omitempty"`
	Name  string  `json:"name"`
}

// App defines model for App.
type App struct {
	Concept ConceptMeta `json:"concept"`
	Meta    AppMeta     `json:"meta"`
	Values  *[]struct {
		Id    *string `json:"id,omitempty"`
		Value *string `json:"value,omitempty"`
	} `json:"values,omitempty"`
}

// AppMeta defines model for AppMeta.
type AppMeta struct {
	Group *string `json:"group,omitempty"`
	Name  string  `json:"name"`
}

// Value defines model for Value.
type Value struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

// ValueList defines model for ValueList.
type ValueList []Value

type Labels map[string]string

type RepoIdentifier string
type Version int

type RepoIndex struct {
	Version  Version        `json:"version"`
	Name     string         `json:"name"`
	Concepts []ConceptEntry `json:"concepts"`
}

type ConceptEntry string

type Meta struct {
	Name   string `json:"name"`
	Labels Labels `json:"labels"`
}

// deleteAppsJSONBody defines parameters for DeleteApps.
type deleteAppsJSONBody []struct {
	Path *string `json:"path,omitempty"`
}

// DeleteAppsRequestBody defines body for DeleteApps for application/json ContentType.
type DeleteAppsJSONRequestBody deleteAppsJSONBody

type Cloner struct {
	Config       ClonerConfig
	Repositories Repositories
}

type ClonerConfig struct {
	git.CloneOptions
	BaseDir string
}

type Repositories map[uuid.UUID]Repository

type Repository struct {
	RepoURL string
	Branch  string
}

func (repos Repositories) ToArray() ([][]interface{}, error) {
	var repoSlices [][]interface{}
	for _, repo := range repos {
		initialized := true
		index, err := repo.GetRepoIndex()
		if err != nil {
			// If index returns a file not exists error,
			// we're okay with that, as the repo is just not initialized.
			if !os.IsNotExist(err) {
				return nil, err
			}
			initialized = false
		}
		repoSlices = append(repoSlices, []interface{}{index.Name, repo.RepoURL, initialized})
	}
	return repoSlices, nil
}

// Target is the interface for all Target implementations
type Target interface {
	TargetName() string
	RenderBundle(outpath string) error
}
