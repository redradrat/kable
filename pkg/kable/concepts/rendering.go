package concepts

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/redradrat/kable/pkg/kable/errors"
)

const (
	RenderStringValueTypeIdentifier ValueTypeIdentifier = "string"
	RenderMapValueTypeIdentifier    ValueTypeIdentifier = "map"
	RenderSelectValueTypeIdentifier ValueTypeIdentifier = "select"
	RenderNameRegexString                               = "^[a-z-_]+$"
)

var RenderNameIsValid = regexp.MustCompile(RenderNameRegexString).MatchString

type File struct {
	path    string
	content []byte
}

type Render struct {
	Info  File
	Files []File
}

func (f File) String() string {
	return string(f.content)
}

func (b Render) Write(baseDir string) error {
	writeFile := func(file File) error {
		path := filepath.Join(baseDir, file.path)
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}
		if err := ioutil.WriteFile(path, file.content, 0666); err != nil {
			return err
		}
		return nil
	}

	for _, file := range b.Files {
		if err := writeFile(file); err != nil {
			return err
		}
	}
	if err := writeFile(b.Info); err != nil {
		return err
	}

	return nil
}

// RenderInfoV1 defines model for RenderInfoV1.
type RenderInfoV1 struct {
	Version  int            `json:"version"`
	Meta     RenderMeta     `json:"meta"`
	Origin   *ConceptOrigin `json:"origin,omitempty"`
	Values   *RenderValues  `json:"values,omitempty"`
	FileTree []string       `json:"files,omitempty"`
}

type ValueTypeIdentifier string

type ValueType interface {
	ValueTypeIdentifier() string
	String() string
}

type RenderValues map[string]ValueType

type MapValueType map[string]string

func (vt MapValueType) ValueTypeIdentifier() string {
	return string(RenderMapValueTypeIdentifier)
}

func (vt MapValueType) String() string {
	outstring, _ := json.Marshal(map[string]string(vt))
	return string(outstring)
}

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
	DateCreated string `json:"date"`
}

func NewRenderV1(avs *RenderValues, origin *ConceptOrigin) (*RenderInfoV1, error) {
	render := RenderInfoV1{
		Version: 1,
		Meta: RenderMeta{
			DateCreated: time.Now().Format(time.RFC822),
		},
		Origin: origin,
	}
	render.Values = avs

	return &render, nil
}

func RenderConcept(path string, avs *RenderValues, ttype TargetType) (*Render, error) {
	return renderConcept(path, nil, avs, ttype)
}

func renderConcept(path string, origin *ConceptOrigin, avs *RenderValues, ttype TargetType) (*Render, error) {
	var target Target
	switch ttype {
	case YamlTargetType:
		target = YamlTarget{}
	case CRDTargetType:
		target = CRDTarget{}
	default:
		return nil, errors.RenderTargetUnsupportedError
	}

	cpt, err := GetConcept(path)
	if err != nil {
		return nil, err
	}

	render, err := target.Render(path, avs, cpt.Type)
	if err != nil {
		return nil, err
	}

	cr, err := NewRenderV1(avs, origin)
	if err != nil {
		return nil, err
	}

	appFile, err := json.MarshalIndent(cr, "", "	")
	if err != nil {
		return nil, err
	}

	render.Info = File{
		path:    ConceptRenderFileName,
		content: appFile,
	}

	return render, nil
}

func RenderRepoConcept(avs *RenderValues, ci ConceptIdentifier, ttype TargetType) (*Render, error) {
	// As the initialization check has been done via GetRepoConcept
	conceptPath, err := GetRepoConceptPath(ci)
	if err != nil {
		return nil, err
	}

	// Compile ConceptRender File
	origin, err := GetConceptOrigin(ci)
	if err != nil {
		return nil, err
	}

	render, err := renderConcept(conceptPath, origin, avs, ttype)

	return render, nil
}
