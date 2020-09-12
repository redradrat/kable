package concepts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/redradrat/kable/pkg/kable/repositories"

	"github.com/redradrat/kable/pkg/kable/errors"

	"github.com/go-git/go-git/v5"
	giturls "github.com/whilp/git-urls"
)

const (
	ConceptFileName                               = "concept.json"
	ConceptIdentifierRegex                        = "^([a-z/]+)@([a-z]+)$"
	ConceptStringInputType    InputTypeIdentifier = "string"
	ConceptSelectionInputType InputTypeIdentifier = "select"
	ConceptJsonnetType        ConceptType         = "jsonnet"
	ConceptJsonnetfile                            = "jsonnetfile.json"
	ConceptMainJsonnet                            = "main.jsonnet"
	ConceptMakefile                               = "Makefile"
	ConceptLibDir                                 = "lib/"
	ConceptVendorDir                              = "vendor/"
	ConceptKlibsonnet                             = "k.libsonnet"
	ConceptRenderFileName                         = "renderinfo.json"
)

var (
	IsValidConceptIdentifier = regexp.MustCompile(ConceptIdentifierRegex).MatchString
	JsonnetMainTemplate      = []byte(`local kausal = import "ksonnet-util/kausal.libsonnet";

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

type ConceptType string

func (ct ConceptType) IsSupported() bool {
	if ct == ConceptJsonnetType {
		return true
	}
	return false
}

// ConceptMeta defines model for ConceptMeta.
type ConceptMeta struct {
	Name       string         `json:"name"`
	Maintainer MaintainerInfo `json:"maintainer,omitempty"`
}

// ConceptInputs defines model for ConceptInputs.
type ConceptInputs struct {
	Mandatory map[string]InputType `json:"mandatory,omitempty"`
	Optional  map[string]InputType `json:"optional,omitempty"`
}

type InputType struct {
	Type    InputTypeIdentifier `json:"type"`
	Options string              `json:"options,omitempty"`
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

func GetConcept(cid ConceptIdentifier) (*Concept, error) {
	if !cid.IsValid() {
		return nil, errors.InvalidConceptIdentifierError
	}
	concept := Concept{}
	if !repositories.IsInitialized(cid.Repo()) {
		return nil, errors.RepositoryNotInitializedError
	}
	content, err := ioutil.ReadFile(filepath.Join(repositories.MustGetCacheInfo(cid.Repo()).Path, filepath.Join(cid.Concept(), ConceptFileName)))
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

func GetConceptOrigin(cid ConceptIdentifier) (*ConceptOrigin, error) {
	if !cid.IsValid() {
		return nil, errors.InvalidConceptIdentifierError
	}
	if !repositories.IsInitialized(cid.Repo()) {
		return nil, errors.RepositoryNotInitializedError
	}
	repo, err := git.PlainOpen(repositories.MustGetCacheInfo(cid.Repo()).Path)
	if err != nil {
		return nil, err
	}

	repoURL, err := repo.Remote(git.DefaultRemoteName)
	if err != nil {
		return nil, err
	}

	parsedURL, err := giturls.Parse(repoURL.Config().URLs[0])
	if err != nil {
		return nil, err
	}

	repoID := parsedURL.Host + parsedURL.Path

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	origin := ConceptOrigin{
		Repository: strings.TrimSuffix(repoID, ".git"),
		Ref:        ref.Hash().String(),
	}
	return &origin, nil
}

type ConceptRepoInfo struct {
	Concept Concept
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

func ListConceptsForRepo(repoid string) ([]ConceptRepoInfo, error) {
	var concepts []ConceptRepoInfo
	if !repositories.IsInitialized(repoid) {
		return concepts, nil
	}
	ri, err := repositories.GetRepoIndex(repoid)
	if err != nil {
		return nil, err
	}
	for _, entry := range ri.ConceptEntries {
		c, err := GetConcept(NewConceptIdentifier(entry, repoid))
		if err != nil {
			return concepts, err
		}
		concepts = append(concepts, ConceptRepoInfo{RepoId: repoid, Path: entry, Concept: *c})
	}
	return concepts, nil
}

func ListConcepts() ([]ConceptRepoInfo, error) {
	var repoList []ConceptRepoInfo
	repos, err := repositories.ListRepositories()
	if err != nil {
		return nil, err
	}
	for id, _ := range repos {
		concepts, err := ListConceptsForRepo(id)
		if err != nil {
			return nil, err
		}
		repoList = append(repoList, concepts...)
	}
	return repoList, nil
}

func InitConcept(name string, conceptType ConceptType) error {
	cpt := Concept{
		ApiVersion: 1,
		Type:       conceptType,
		Meta: ConceptMeta{
			Name: name,
		},
		Inputs: ConceptInputs{
			Mandatory: map[string]InputType{},
			Optional:  nil,
		},
	}

	switch conceptType {
	case ConceptJsonnetType:
		if err := createFile(JsonnetMainTemplate, ConceptMainJsonnet); err != nil {
			return err
		}
		if err := createFile(JsonnetDepTemplate, ConceptJsonnetfile); err != nil {
			return err
		}
		if err := createFile(JsonnetMakeFile, ConceptMakefile); err != nil {
			return err
		}
		if err := os.MkdirAll(ConceptLibDir, os.ModePerm); err != nil {
			return err
		}
		if err := createFile(JsonnetLibTemplate, filepath.Join(ConceptLibDir, ConceptKlibsonnet)); err != nil {
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
