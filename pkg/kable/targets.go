package kable

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-jsonnet"
)

const YamlTargetIdentifier = "yaml"

type file struct {
	path    string
	content []byte
}

type Bundle struct {
	files   []file
	baseDir string
}

// Target is the interface for all Target implementations
type Target interface {
	TargetName() string
	RenderBundle(app *AppV1, ci ConceptIdentifier, outpath string) (*Bundle, error)
}

func (b Bundle) Write() error {
	for _, file := range b.files {
		path := filepath.Join(b.baseDir, file.path)
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}
		if err := ioutil.WriteFile(path, file.content, 0666); err != nil {
			return err
		}
	}
	return nil
}

type YamlTarget struct {
}

func (y YamlTarget) TargetName() string {
	return YamlTargetIdentifier
}

func (y YamlTarget) RenderBundle(app *AppV1, ci ConceptIdentifier, outpath string) (*Bundle, error) {
	bundle := Bundle{
		baseDir: filepath.Join(outpath, app.Meta.Name),
	}

	cpt, err := GetConcept(ci)
	if err != nil {
		return nil, err
	}

	// As the initialization check has been done via GetConcept
	cache := MustGetCacheInfo(ci.Repo())

	switch cpt.Type {
	case ConceptJsonnetType:
		bundle.files, err = renderJsonnetConcept(app.Meta.Name, filepath.Join(cache.Path, ci.Concept()), app.Values)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ConceptTypeUnsupportedError

	}

	appFile, err := json.MarshalIndent(app, "", "	")
	if err != nil {
		return nil, err
	}

	bundle.files = append(bundle.files, file{
		path:    "App.json",
		content: appFile,
	})

	return &bundle, nil
}

func renderJsonnetConcept(name, path string, avs *AppValues) ([]file, error) {
	var bundle []file

	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{
		JPaths: []string{
			filepath.Join(path, ConceptLibDir),
			filepath.Join(path, ConceptVendorDir),
		},
	})

	for id, val := range *avs {
		vm.ExtVar(id, val.String())
	}

	mainJsonnet, err := ioutil.ReadFile(filepath.Join(path, ConceptMainJsonnet))
	if err != nil {
		return nil, err
	}

	out, err := vm.EvaluateSnippet(ConceptMainJsonnet, string(mainJsonnet))
	if err != nil {
		return nil, err
	}

	bundle = append(bundle, file{
		path:    name + ".yaml",
		content: []byte(out),
	})

	return bundle, nil
}
