package concepts

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/redradrat/kable/pkg/repositories"

	"github.com/redradrat/kable/pkg/errors"
)

const (
	RenderStringValueTypeIdentifier ValueTypeIdentifier = "string"
	RenderMapValueTypeIdentifier    ValueTypeIdentifier = "map"
	RenderIntValueTypeIdentifier    ValueTypeIdentifier = "int"
	RenderBoolValueTypeIdentifier   ValueTypeIdentifier = "bool"
	RenderNameRegexString                               = "^[a-z-_]+$"
)

var RenderNameIsValid = regexp.MustCompile(RenderNameRegexString).MatchString

type File struct {
	path    string
	content []byte
}

type Render struct {
	Info  *File
	Files []File
}

func (f File) String() string {
	return string(f.content)
}

func (r Render) PrintFiles() string {
	var out []byte

	sort.Slice(r.Files, func(i, j int) bool {
		return r.Files[i].path < r.Files[j].path
	})

	for _, file := range r.Files {
		out = append(out, file.content...)
	}
	return string(out)
}

func writeFile(file File, baseDir string) error {
	path := filepath.Join(baseDir, file.path)
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, file.content, 0666); err != nil {
		return err
	}
	return nil
}

func (r Render) WriteFiles(baseDir string) error {
	for _, file := range r.Files {
		if err := writeFile(file, baseDir); err != nil {
			return err
		}
	}

	return nil
}

func (r Render) Write(baseDir string) error {
	err := r.WriteFiles(baseDir)
	if err != nil {
		return err
	}
	return r.WriteInfo(baseDir)
}

func (r Render) WriteInfo(baseDir string) error {
	if r.Info != nil {
		if err := writeFile(*r.Info, baseDir); err != nil {
			return err
		}
	}
	return nil
}

// RenderInfoV1 defines model for RenderInfoV1.
type RenderInfoV1 struct {
	Version int            `json:"version"`
	Meta    RenderMeta     `json:"meta"`
	Origin  *ConceptOrigin `json:"origin,omitempty"`
	Values  *RenderValues  `json:"values,omitempty"`
}

func ParseRenderInfoV1FromFile(path string) (*RenderInfoV1, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ri := &RenderInfoV1{}
	if err := json.Unmarshal(f, &ri); err != nil {
		return nil, err
	}

	return ri, nil
}

type ValueTypeIdentifier string

type ValueType interface {
	ValueTypeIdentifier() string
	String() string
}

type RenderValues map[string]ValueType

func (rv RenderValues) UnmarshalJSON(bytes []byte) error {
	inter := map[string]interface{}{}
	if err := json.Unmarshal(bytes, &inter); err != nil {
		return err
	}

	for k, v := range inter {
		switch assertedValue := v.(type) {
		case string:
			rv[k] = StringValueType(assertedValue)
		case map[string]interface{}:
			rv[k] = MapValueType(assertedValue)
		case int:
			rv[k] = IntValueType(assertedValue)
		case bool:
			rv[k] = BoolValueType(assertedValue)
		}
	}
	return nil
}

func (rv RenderValues) Map() map[string]ValueType {
	return map[string]ValueType(rv)
}

type MapValueType map[string]interface{}

func (vt MapValueType) ValueTypeIdentifier() string {
	return string(RenderMapValueTypeIdentifier)
}

func (vt MapValueType) String() string {
	outstring, _ := json.Marshal(map[string]interface{}(vt))
	return string(outstring)
}

type StringValueType string

func (vt StringValueType) ValueTypeIdentifier() string {
	return string(RenderStringValueTypeIdentifier)
}

func (vt StringValueType) String() string {
	return string(vt)
}

type IntValueType int

func (vt IntValueType) ValueTypeIdentifier() string {
	return string(RenderIntValueTypeIdentifier)
}

func (vt IntValueType) String() string {

	return strconv.Itoa(int(vt))
}

type BoolValueType bool

func (vt BoolValueType) ValueTypeIdentifier() string {
	return string(RenderBoolValueTypeIdentifier)
}

func (vt BoolValueType) String() string {
	if vt == false {
		return "false"
	} else {
		return "true"
	}
}

// RenderMeta defines model for RenderMeta.
type RenderMeta struct {
	DateCreated string `json:"date"`
}

type RenderOpts struct {
	Local           bool
	WriteRenderInfo bool
	Single          bool
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

func RenderConcept(path string, avs *RenderValues, ttype TargetType, opts RenderOpts) (*Render, error) {
	return renderConcept(path, avs, ttype, opts)
}

func renderConcept(id string, avs *RenderValues, ttype TargetType, opts RenderOpts) (*Render, error) {
	var err error
	var origin *ConceptOrigin
	path := id
	if !opts.Local {
		// Check if the identifier is correct
		if !IsValidConceptIdentifier(id) {
			return nil, errors.InvalidConceptIdentifierError
		}

		ci := ConceptIdentifier(id)
		r, err := repositories.GetRepository(ci.Repo())
		if err != nil {
			return nil, err
		}

		// Get the repo path
		repopath, err := r.AbsolutePath()
		if err != nil {
			return nil, err
		}

		path = filepath.Join(repopath, ci.Concept())

		// Get the origin of the the concept
		origin = &ConceptOrigin{
			Repository: r.URL,
			Ref:        r.GitRef,
		}
	}

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

	render, err := target.Render(path, avs, cpt.Type, opts.Single)
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

	render.Info = &File{
		path:    ConceptRenderFileName,
		content: appFile,
	}

	return render, nil
}
