// Package srv provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

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

// Concept defines model for Concept.
type Concept struct {
	Meta   ConceptMeta `json:"meta"`
	Values *struct {
		Mandatory *ValueList `json:"mandatory,omitempty"`
		Optional  *ValueList `json:"optional,omitempty"`
	} `json:"values,omitempty"`
}

// ConceptMeta defines model for ConceptMeta.
type ConceptMeta struct {
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

// deleteAppsJSONBody defines parameters for DeleteApps.
type deleteAppsJSONBody []struct {
	Path *string `json:"path,omitempty"`
}

// DeleteAppsRequestBody defines body for DeleteApps for application/json ContentType.
type DeleteAppsJSONRequestBody deleteAppsJSONBody