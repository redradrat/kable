package concepts

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/redradrat/kable/pkg/kable/repositories"

	"github.com/redradrat/kable/pkg/kable/errors"

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
	RenderBundle(app *ConceptRenderV1, ci ConceptIdentifier, outpath string) (*Bundle, error)
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

func (y YamlTarget) RenderBundle(cr *ConceptRenderV1, ci ConceptIdentifier, outpath string) (*Bundle, error) {
	bundle := Bundle{
		baseDir: filepath.Join(outpath, cr.Meta.Name),
	}

	cpt, err := GetConcept(ci)
	if err != nil {
		return nil, err
	}

	// As the initialization check has been done via GetConcept
	cache := repositories.MustGetCacheInfo(ci.Repo())

	switch cpt.Type {
	case ConceptJsonnetType:
		bundle.files, err = renderJsonnetConcept(cr.Meta.Name, filepath.Join(cache.Path, ci.Concept()), cr.Values)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.ConceptTypeUnsupportedError

	}

	appFile, err := json.MarshalIndent(cr, "", "	")
	if err != nil {
		return nil, err
	}

	bundle.files = append(bundle.files, file{
		path:    ConceptRenderFileName,
		content: appFile,
	})

	return &bundle, nil
}

func renderJsonnetConcept(name, path string, avs *RenderValues) ([]file, error) {
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
