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

	"github.com/redradrat/kable/kable/repositories"

	"github.com/redradrat/kable/kable/errors"

	"github.com/go-git/go-git/v5"
	giturls "github.com/whilp/git-urls"
)

const (
	ConceptFileName                               = "concept.json"
	ConceptIdentifierRegex                        = "^([a-z/]+)@([a-z]+)$"
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
	path, err := GetRepoConceptPath(cid)
	if err != nil {
		return nil, err
	}
	return GetConcept(path)
}

func GetRepoConceptPath(cid ConceptIdentifier) (string, error) {
	if !cid.IsValid() {
		return "", errors.InvalidConceptIdentifierError
	}
	if !repositories.IsInitialized(cid.Repo()) {
		return "", errors.RepositoryNotInitializedError
	}
	return filepath.Join(repositories.MustGetCacheInfo(cid.Repo()).AbsolutePath(), filepath.Join(cid.Concept())), nil
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

func GetConceptOrigin(cid ConceptIdentifier) (*ConceptOrigin, error) {
	if !cid.IsValid() {
		return nil, errors.InvalidConceptIdentifierError
	}
	if !repositories.IsInitialized(cid.Repo()) {
		return nil, errors.RepositoryNotInitializedError
	}
	repo, err := git.PlainOpen(repositories.MustGetCacheInfo(cid.Repo()).AbsolutePath())
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
		c, err := GetRepoConcept(NewConceptIdentifier(entry, repoid))
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
