package kable

import (
	"regexp"
	"time"
)

const (
	AppStringValueTypeIdentifier ValueTypeIdentifier = "string"
	AppSelectValueTypeIdentifier ValueTypeIdentifier = "select"
	AppNameRegexString                               = "^[a-z-_]+$"
)

var AppNameIsValid = regexp.MustCompile(AppNameRegexString).MatchString

// AppV1 defines model for AppV1.
type AppV1 struct {
	Version  int            `json:"version"`
	Meta     AppMeta        `json:"meta"`
	Origin   *ConceptOrigin `json:"origin"`
	Values   *AppValues     `json:"values,omitempty"`
	FileTree []string       `json:"files,omitempty"`
}

type ValueTypeIdentifier string

type ValueType interface {
	ValueTypeIdentifier() string
	String() string
}

type AppValues map[string]ValueType

type StringValueType string

func (vt StringValueType) ValueTypeIdentifier() string {
	return string(AppStringValueTypeIdentifier)
}

func (vt StringValueType) String() string {
	return string(vt)
}

type SelectValueType string

func (vt SelectValueType) ValueTypeIdentifier() string {
	return string(AppSelectValueTypeIdentifier)
}

func (vt SelectValueType) String() string {
	return string(vt)
}

// AppMeta defines model for AppMeta.
type AppMeta struct {
	Name        string `json:"name"`
	DateCreated string `json:"date"`
}

func NewAppV1(name string, avs *AppValues) (*AppV1, error) {
	if !AppNameIsValid(name) {
		return nil, InvalidAppNameError
	}

	app := AppV1{
		Version: 1,
		Meta: AppMeta{
			Name:        name,
			DateCreated: time.Now().Format(time.RFC822),
		},
	}
	app.Values = avs

	return &app, nil
}

func RenderApp(app *AppV1, ci ConceptIdentifier, output string, target Target) error {
	var err error

	app.Origin, err = GetConceptOrigin(ci)
	if err != nil {
		return err
	}

	bundle, err := target.RenderBundle(app, ci, output)
	if err != nil {
		return err
	}
	if err := bundle.Write(); err != nil {
		return err
	}

	return nil
}
