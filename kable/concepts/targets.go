package concepts

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/ghodss/yaml"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/redradrat/kable/kable/errors"

	"github.com/google/go-jsonnet"
)

const (
	YamlTargetType TargetType = "yaml"
	CRDTargetType  TargetType = "crd"
)

type TargetType string

// Target is the interface for all Target implementations
type Target interface {
	TargetName() string
	Render(path string, vals *RenderValues, cpt ConceptType) (*Render, error)
}

type CRDTarget struct {
}

func (c CRDTarget) TargetName() string {
	return string(CRDTargetType)
}

func (c CRDTarget) Render(path string, vals *RenderValues, cpt ConceptType) (*Render, error) {
	panic("implement me")
}

type YamlTarget struct {
}

func (y YamlTarget) TargetName() string {
	return string(YamlTargetType)
}

func (y YamlTarget) Render(path string, vals *RenderValues, cpt ConceptType) (*Render, error) {
	var err error
	bundle := Render{}

	switch cpt {
	case ConceptJsonnetType:
		bundle.Files, err = renderJsonnetConcept(path, vals)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.ConceptTypeUnsupportedError
	}

	return &bundle, nil
}

func renderJsonnetConcept(path string, avs *RenderValues) ([]File, error) {
	var bundle []File

	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{
		JPaths: []string{
			filepath.Join(path),
			filepath.Join(path, ConceptLibDir),
			filepath.Join(path, ConceptVendorDir),
		},
	})

	if avs != nil {
		for id, val := range *avs {
			switch val.(type) {
			case StringValueType:
				vm.ExtVar(id, val.String())
			case MapValueType, IntValueType, BoolValueType:
				vm.ExtCode(id, val.String())
			default:
				return nil, errors.ValueTypeNotSupported
			}
		}
	}

	mainJsonnet, err := ioutil.ReadFile(filepath.Join(path, ConceptMainJsonnet))
	if err != nil {
		return nil, err
	}

	jsonnetout, err := vm.EvaluateSnippet(ConceptMainJsonnet, string(mainJsonnet))
	if err != nil {
		return nil, err
	}

	var objs unstructured.UnstructuredList

	// attempt to unmarshal either array or single object
	var jsonObjs []unstructured.Unstructured
	err = json.Unmarshal([]byte(jsonnetout), &jsonObjs)
	if err == nil {
		objs.Items = append(objs.Items, jsonObjs...)
	} else {
		var jsonObj unstructured.Unstructured
		err = json.Unmarshal([]byte(jsonnetout), &jsonObj)
		if err != nil {
			return nil, err
		}
		objs.Items = append(objs.Items, jsonObj)
	}

	objs.SetAPIVersion("v1")
	objs.SetKind("List")

	jsonout, err := objs.MarshalJSON()
	if err != nil {
		return nil, err
	}

	yamlout, err := yaml.JSONToYAML(jsonout)
	if err != nil {
		return nil, err
	}

	bundle = append(bundle, File{
		path:    "manifest.yaml",
		content: yamlout,
	})

	return bundle, nil
}
