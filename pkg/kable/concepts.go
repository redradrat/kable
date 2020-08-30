package kable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	ConceptFileName        = "concept.json"
	ConceptTypeUnsupported = "given concept type is not supported"
)

var (
	JsonnetMainTemplate = []byte(`local kausal = import "ksonnet-util/kausal.libsonnet";

local container = kausal.core.v1.container;
local port = kausal.core.v1.containerPort;
local service = kausal.core.v1.service;

local echoServDeployment = deployment.new(
    name=std.extVar("instanceName"), replicas=1,
    containers=[
      container.new("echoserver", "k8s.gcr.io/echoserver:1.4")
      + container.withPorts([port.new("ui", "8080")]),
    ],
  );

[
  echoServDeployment,
  kausal.util.serviceFor(echoServDeployment),
]
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
	jsonnet main.jsonnet -J lib -J vendor --ext-str instanceName="dummy" | yq r --prettyPrint -

install:
	jb install
`)
)

// Concept defines model for Concept.
type Concept struct {
	ApiVersion int           `json:"apiVersion"`
	Meta       ConceptMeta   `json:"metadata"`
	Inputs     ConceptInputs `json:"inputs,omitempty"`
}

// ConceptMeta defines model for ConceptMeta.
type ConceptMeta struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ConceptInputs defines model for ConceptInputs.
type ConceptInputs struct {
	Mandatory map[string]InputType `json:"mandatory,omitempty"`
	Optional  map[string]InputType `json:"optional,omitempty"`
}

type InputType string

func (it InputType) String() string {
	return string(it)
}

func parseConcept(path string, repoid uuid.UUID) (*Concept, error) {
	concept := Concept{}
	if !IsInitialized(repoid) {
		return nil, fmt.Errorf(RepositoryNotInitializedError)
	}
	content, err := ioutil.ReadFile(filepath.Join(MustGetCachePath(repoid), path))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(content, &concept); err != nil {
		return nil, err
	}
	return &concept, nil
}

type ConceptRepoPair struct {
	Concept Concept
	Path    string
	RepoId  uuid.UUID
}

func ListConceptsForRepo(repoid uuid.UUID) ([]ConceptRepoPair, error) {
	var concepts []ConceptRepoPair
	if !IsInitialized(repoid) {
		return concepts, nil
	}
	ri, err := GetRepoIndex(repoid)
	if err != nil {
		return nil, err
	}
	for _, path := range ri.ConceptEntries {
		c, err := parseConcept(filepath.Join(path.String(), ConceptFileName), repoid)
		if err != nil {
			return concepts, err
		}
		concepts = append(concepts, ConceptRepoPair{RepoId: repoid, Path: path.String(), Concept: *c})
	}
	return concepts, nil
}

func ListConcepts() ([]ConceptRepoPair, error) {
	var repoList []ConceptRepoPair
	for id, _ := range currentConfig.Repositories {
		concepts, err := ListConceptsForRepo(id)
		if err != nil {
			return nil, err
		}
		repoList = append(repoList, concepts...)
	}
	return repoList, nil
}

func InitConcept(name, conceptType string) error {
	cpt := Concept{
		ApiVersion: 1,
		Meta: ConceptMeta{
			Name:    name,
			Version: "0.1.0",
		},
		Inputs: ConceptInputs{
			Mandatory: map[string]InputType{},
			Optional:  nil,
		},
	}

	switch conceptType {
	case "jsonnet":
		if err := createFile(JsonnetMainTemplate, "./main.jsonnet"); err != nil {
			return err
		}
		if err := createFile(JsonnetDepTemplate, "./jsonnetfile.json"); err != nil {
			return err
		}
		if err := createFile(JsonnetMakeFile, "./Makefile"); err != nil {
			return err
		}
		if err := os.MkdirAll("./lib", os.ModePerm); err != nil {
			return err
		}
		if err := createFile(JsonnetLibTemplate, "./lib/k.libsonnet"); err != nil {
			return err
		}

		cmd := exec.Command("jb", "install")
		err := cmd.Run()
		if err != nil {
			return nil
		}
	case "yaml":

	default:
		return fmt.Errorf(ConceptTypeUnsupported)
	}

	if err := createJson(cpt, "./concept.json"); err != nil {
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
