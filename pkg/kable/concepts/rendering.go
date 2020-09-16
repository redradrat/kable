package concepts

import (
	"regexp"
	"time"

	"github.com/redradrat/kable/pkg/kable/errors"
)

const (
	RenderStringValueTypeIdentifier ValueTypeIdentifier = "string"
	RenderSelectValueTypeIdentifier ValueTypeIdentifier = "select"
	RenderNameRegexString                               = "^[a-z-_]+$"
)

var RenderNameIsValid = regexp.MustCompile(RenderNameRegexString).MatchString

// ConceptRenderV1 defines model for ConceptRenderV1.
type ConceptRenderV1 struct {
	Version  int            `json:"version"`
	Meta     RenderMeta     `json:"meta"`
	Origin   *ConceptOrigin `json:"origin"`
	Values   *RenderValues  `json:"values,omitempty"`
	FileTree []string       `json:"files,omitempty"`
}

type ValueTypeIdentifier string

type ValueType interface {
	ValueTypeIdentifier() string
	String() string
}

type RenderValues map[string]ValueType

type StringValueType string

func (vt StringValueType) ValueTypeIdentifier() string {
	return string(RenderStringValueTypeIdentifier)
}

func (vt StringValueType) String() string {
	return string(vt)
}

type SelectValueType string

func (vt SelectValueType) ValueTypeIdentifier() string {
	return string(RenderSelectValueTypeIdentifier)
}

func (vt SelectValueType) String() string {
	return string(vt)
}

// RenderMeta defines model for RenderMeta.
type RenderMeta struct {
	Name        string `json:"name"`
	DateCreated string `json:"date"`
}

func NewRenderV1(name string, avs *RenderValues) (*ConceptRenderV1, error) {
	if !RenderNameIsValid(name) {
		return nil, errors.InvalidRenderNameError
	}

	app := ConceptRenderV1{
		Version: 1,
		Meta: RenderMeta{
			Name:        name,
			DateCreated: time.Now().Format(time.RFC822),
		},
	}
	app.Values = avs

	return &app, nil
}

func RenderConcept(cr *ConceptRenderV1, ci ConceptIdentifier, output string, target Target) (*Bundle, error) {
	var err error

	cr.Origin, err = GetConceptOrigin(ci)
	if err != nil {
		return nil, err
	}

	return target.RenderBundle(cr, ci, output)
}
