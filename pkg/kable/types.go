package kable

type APIVersion string
type Kind string

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
