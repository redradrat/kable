package kable

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	YamlTargetIdentifier = "yaml"
)

// App defines model for App.
type App struct {
	Version  int      `json:"version"`
	Meta     AppMeta  `json:"meta"`
	Origin   Origin   `json:"origin"`
	FileTree []string `json:"files,omitempty"`
}

// AppMeta defines model for AppMeta.
type AppMeta struct {
	Name        string `json:"name"`
	DateCreated string `json:"created"`
}

// Origin defines the git source of origin
type Origin struct {
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
}

type file struct {
	path    string
	content []byte
}

type bundle struct {
	files   []file
	baseDir string
}

func (b bundle) Write() error {
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

// Target is the interface for all Target implementations
type Target interface {
	TargetName() string
	RenderBundle(concept Concept, outpath string) bundle
}

type YamlTarget struct {
}

func (y YamlTarget) TargetName() string {
	return YamlTargetIdentifier
}

func (y YamlTarget) RenderBundle(concept Concept, outpath string) bundle {
	panic("not yet implemented")
}

func RenderConcept(ci ConceptIdentifier, output string, target Target) error {
	app := App{}
	cache, err := GetCacheInfo(ci.Repo())
	if err != nil {
		return err
	}
	app.Origin.Repository = cache.URI
	app.Origin.Ref = cache.Branch

	concept, err := GetConcept(ci.Concept(), ci.Repo())
	if err != nil {
		return err
	}

	b := target.RenderBundle(*concept, output)

	return b.Write()
}
