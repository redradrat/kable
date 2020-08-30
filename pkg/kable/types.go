package kable

type APIVersion string
type Kind string

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

type Labels map[string]string

type RepoIdentifier string

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

// Target is the interface for all Target implementations
type Target interface {
	TargetName() string
	RenderBundle(outpath string) error
}
