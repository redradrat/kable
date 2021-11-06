package concepts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/redradrat/kable/pkg/repositories"

	"github.com/redradrat/kable/pkg/errors"
)

const (
	ConceptFileName = "concept.json"
	// _ (underscore) is specifically not part of this list, as this will be our
	// replacement character for forming URLs
	ConceptIdentifierRegex                        = "^([a-z/\\-123456789]+)@([a-z\\-]+)$"
	ConceptStringInputType    InputTypeIdentifier = "string"
	ConceptSelectionInputType InputTypeIdentifier = "select"
	ConceptMapInputType       InputTypeIdentifier = "map"
	ConceptIntInputType       InputTypeIdentifier = "int"
	ConceptBoolInputType      InputTypeIdentifier = "bool"
	ConceptJsonnetType        ConceptType         = "jsonnet"
	ConceptJsonnetfile                            = "jsonnetfile.json"
	ConceptMainJsonnet                            = "main.jsonnet"
	ConceptMakefile                               = "Makefile"
	ConceptGitignorefile                          = ".gitignore"
	ConceptLibDir                                 = "lib/"
	ConceptVendorDir                              = "vendor/"
	ConceptKlibsonnet                             = "k.libsonnet"
	ConceptMainlibsonnet                          = "main.libsonnet"
	ConceptRenderFileName                         = "renderinfo.json"
)

var (
	IsValidConceptIdentifier = regexp.MustCompile(ConceptIdentifierRegex).MatchString
	JsonnetMainTemplate      = []byte(`local lib = import 'lib/main.libsonnet';

// Final JSON Output
lib.new(std.extVar("instanceName"))
`)
	JsonnetMainLibTemplate = []byte(`local kausal = import "ksonnet-util/kausal.libsonnet";

local deployment = kausal.apps.v1.deployment;
local container = kausal.core.v1.container;
local port = kausal.core.v1.containerPort;
local service = kausal.core.v1.service;

local grafanaDeploy(name) = deployment.new(
        name=name, replicas=2,
        containers=[
          container.new("grafana", "grafana/grafana")
          + container.withPorts([port.new("ui", 10330)]),
        ],
      );


// Final JSON Object
{
  new(name):: [
    grafanaDeploy(name),
    kausal.util.serviceFor(grafanaDeploy(name))
  ]
}
`)

	JsonnetLibTemplate = []byte(`(import "github.com/jsonnet-libs/k8s-alpha/1.14/main.libsonnet")
+ (import "github.com/jsonnet-libs/k8s-alpha/1.14/extensions/kausal-shim.libsonnet")
`)
	JsonnetDepTemplate = []byte(`{
  "version": 1,
  "dependencies": [
    {
      "source": {
        "git": {
          "remote": "https://github.com/grafana/jsonnet-libs.git",
          "subdir": "ksonnet-util"
        }
      },
      "version": "master"
    },
    {
      "source": {
        "git": {
          "remote": "https://github.com/jsonnet-libs/k8s-alpha.git",
          "subdir": "1.14"
        }
      },
      "version": "master"
    }
  ],
  "legacyImports": true
}
`)
	JsonnetMakeFile = []byte(`render:
	kable render -l . -o out/

install:
	jb install
`)
	JsonnetGitignoreFile = []byte(`out/
`)
)

type ConceptIdentifier string

func (ci ConceptIdentifier) String() string {
	return string(ci)
}

func (ci ConceptIdentifier) IsValid() bool {
	return IsValidConceptIdentifier(ci.String())
}

func (ci ConceptIdentifier) Concept() string {
	getStrings := regexp.MustCompile(ConceptIdentifierRegex).FindStringSubmatch
	matches := getStrings(ci.String())
	return matches[1]
}

func (ci ConceptIdentifier) Repo() string {
	getStrings := regexp.MustCompile(ConceptIdentifierRegex).FindStringSubmatch
	matches := getStrings(ci.String())
	return matches[2]
}

func NewConceptIdentifier(path, repoid string) ConceptIdentifier {
	return ConceptIdentifier(path + "@" + repoid)
}

// Concept defines model for Concept.
type Concept struct {
	ApiVersion int           `json:"apiVersion"`
	Type       ConceptType   `json:"type"`
	Meta       ConceptMeta   `json:"metadata"`
	Inputs     ConceptInputs `json:"inputs,omitempty"`
}

func (c *Concept) UnmarshalJSON(bytes []byte) error {
	type conceptCopy Concept
	inter := conceptCopy{}
	if err := json.Unmarshal(bytes, &inter); err != nil {
		return err
	}
	sort.Strings(inter.Meta.Tags)

	c.ApiVersion = inter.ApiVersion
	c.Meta = inter.Meta
	c.Inputs = inter.Inputs
	c.Type = inter.Type

	return nil
}

type ConceptType string

func (ct ConceptType) String() string {
	return string(ct)
}

func (ct ConceptType) IsSupported() bool {
	if ct == ConceptJsonnetType {
		return true
	}
	return false
}

// ConceptMeta defines model for ConceptMeta.
type ConceptMeta struct {
	Name       string         `json:"name"`
	Tags       Tags           `json:"tags,omitempty"`
	Maintainer MaintainerInfo `json:"maintainer,omitempty"`
}

type Tags []string

// ConceptInputs defines model for ConceptInputs.
type ConceptInputs struct {
	Mandatory map[string]InputType `json:"mandatory,omitempty"`
	Optional  map[string]InputType `json:"optional,omitempty"`
}

func (ci ConceptInputs) All() map[string]InputType {
	outmap := ci.Optional
	for k, v := range ci.Mandatory {
		outmap[k] = v
	}
	return outmap
}

type InputType struct {
	Type    InputTypeIdentifier `json:"type"`
	Options []string            `json:"options,omitempty"`
}

type InputTypeIdentifier string

func (iti InputTypeIdentifier) IsValid() bool {
	if iti == ConceptSelectionInputType || iti == ConceptStringInputType {
		return true
	}
	return false
}

func (iti InputTypeIdentifier) String() string {
	return string(iti)
}

func GetRepoConcept(cid ConceptIdentifier) (*Concept, error) {
	r, err := repositories.GetRepository(cid.Repo())
	if err != nil {
		return nil, err
	}
	path, err := r.AbsolutePath()
	if err != nil {
		return nil, err
	}
	return GetConcept(filepath.Join(path, cid.Concept()))
}

func GetConcept(path string) (*Concept, error) {
	concept := Concept{}
	content, err := ioutil.ReadFile(filepath.Join(path, ConceptFileName))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(content, &concept); err != nil {
		return nil, err
	}
	return &concept, nil
}

// Origin defines the git source of origin
type ConceptOrigin struct {
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
}

func GetConceptOriginFromRepository(repositoryName string) (*ConceptOrigin, error) {
	r, err := repositories.GetRepository(repositoryName)
	if err != nil {
		return nil, err
	}

	// Get the origin of the the concept
	origin := &ConceptOrigin{
		Repository: r.URL,
		Ref:        r.GitRef,
	}

	return origin, nil
}

type ConceptRepoInfo struct {
	Concept string
	Path    string
	RepoId  string
}

type MaintainerInfo struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempt"`
}

func (mi MaintainerInfo) String() string {
	if mi.Name == "" {
		return ""
	}
	elements := []string{mi.Name}
	if mi.Email != "" {
		elements = append(elements, fmt.Sprintf("<%s>", mi.Email))
	}
	return strings.Join(elements, " ")
}

func ListConcepts() ([]ConceptIdentifier, error) {
	var cis []ConceptIdentifier
	repos, err := repositories.ListRepositories()
	if err != nil {
		return nil, err
	}
	for _, repo := range repos {
		idx, err := repo.RepoIndex()
		if err != nil {
			return nil, err
		}
		for _, path := range idx.ConceptEntries {
			cis = append(cis, NewConceptIdentifier(path, repo.Name))
		}
	}
	return cis, nil
}

func InitConcept(workdir, name string, conceptType ConceptType) error {
	cpt := Concept{
		ApiVersion: 1,
		Type:       conceptType,
		Meta: ConceptMeta{
			Name: name,
		},
		Inputs: ConceptInputs{
			Mandatory: map[string]InputType{
				"instanceName": {
					Type: ConceptStringInputType,
				},
				"nameSelection": {
					Type:    ConceptSelectionInputType,
					Options: []string{"Option 1", "Option 2"},
				},
			},
			Optional: nil,
		},
	}

	if workdir != "." {
		if err := os.MkdirAll(workdir, os.ModePerm); err != nil {
			return err
		}
	}

	switch conceptType {
	case ConceptJsonnetType:
		if err := createFile(JsonnetMainTemplate, filepath.Join(workdir, ConceptMainJsonnet)); err != nil {
			return err
		}
		if err := createFile(JsonnetDepTemplate, filepath.Join(workdir, ConceptJsonnetfile)); err != nil {
			return err
		}
		if err := createFile(JsonnetMakeFile, filepath.Join(workdir, ConceptMakefile)); err != nil {
			return err
		}
		if err := createFile(JsonnetGitignoreFile, filepath.Join(workdir, ConceptGitignorefile)); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join(workdir, ConceptLibDir), os.ModePerm); err != nil {
			return err
		}
		if err := createFile(JsonnetMainLibTemplate, filepath.Join(filepath.Join(workdir, ConceptLibDir), ConceptMainlibsonnet)); err != nil {
			return err
		}
		if err := createFile(JsonnetLibTemplate, filepath.Join(filepath.Join(workdir, ConceptLibDir), ConceptKlibsonnet)); err != nil {
			return err
		}

		cmd := exec.Command("jb", "install")
		err := cmd.Run()
		if err != nil {
			return nil
		}
	default:
		return errors.ConceptTypeUnsupportedError
	}

	if err := createJson(cpt, filepath.Join(workdir, "concept.json")); err != nil {
		return err
	}
	return nil
}

func createFile(out []byte, path string) error {
	if err := ioutil.WriteFile(path, out, 0666); err != nil {
		return err
	}
	return nil
}

func createJson(content interface{}, path string) error {
	out, err := json.MarshalIndent(content, "", "	")
	if err != nil {
		return err
	}
	return createFile(out, path)
}
