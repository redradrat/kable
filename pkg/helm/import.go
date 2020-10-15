package helm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/grafana/tanka/pkg/helm"

	"github.com/redradrat/kable/pkg/concepts"

	"github.com/jsonnet-bundler/jsonnet-bundler/spec/v1/deps"

	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"
)

const jsonnetTpl = `local helm = (import "github.com/grafana/jsonnet-libs/helm-util/helm.libsonnet").new(std.thisFile);

{
  {{.Name}}: helm.template("{{.Name}}", "../charts/{{.Name}}", {
    namespace: "default",
    values: {
      foo: { bar: baz }
    }
  })
}
`

const helmConceptLibTpl = `local helm = (import "github.com/grafana/jsonnet-libs/helm-util/helm.libsonnet").new(std.thisFile);

{
  _values:: {
    foo: error "missing value foo"
  },
  {{.Name}}: helm.template("{{.Name}}", "../charts/{{.Name}}", {
    namespace: "default",
    values: $._values
  })
}
`
const helmConceptMainTpl = `local lib = import "lib/{{.Name}}.libsonnet";

local values = {
  foo: "test"
};

lib + { _values+: values }
`

type HelmChart struct {
	Repo    string
	Name    string
	Version string
}

func (hc HelmChart) Requirement() string {
	return fmt.Sprintf("%s/%s@%s", hc.Repo, hc.Name, hc.Version)
}

func InitHelmConcept(chart HelmChart, out string) error {

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(out, os.ModePerm); err != nil {
		return err
	}
	if err := os.Chdir(out); err != nil {
		return err
	}

	// Initialize a concept
	if err := concepts.InitConcept(".", chart.Name, concepts.ConceptJsonnetType); err != nil {
		return err
	}
	libpath := filepath.Join(concepts.ConceptLibDir, concepts.ConceptMainlibsonnet)
	if err := os.Remove(libpath); err != nil && !os.IsNotExist(err) {
		return err
	}
	mainpath := filepath.Join(concepts.ConceptMainJsonnet)
	if err := os.Remove(mainpath); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := ImportHelmChart(chart, "."); err != nil {
		return err
	}

	librender, err := templateString(chart.Name, helmConceptLibTpl)
	mainrender, err := templateString(chart.Name, helmConceptMainTpl)

	if err := ioutil.WriteFile(filepath.Join(concepts.ConceptLibDir, chart.Name+".libsonnet"), librender, 0644); err != nil {
		return err
	}
	if err := ioutil.WriteFile(mainpath, mainrender, 0644); err != nil {
		return err
	}

	if err := os.Chdir(wd); err != nil {
		return err
	}

	return nil
}

func ImportHelmChart(helmChart HelmChart, out string) error {

	// TODO (redradrat): Try to make upstream's logging configurable (silent)
	cf, err := helm.InitChartfile(filepath.Join(out, "chartfile.yaml"))
	if err != nil {
		if os.IsExist(err) {
			cf, err = helm.LoadChartfile(filepath.Join(out, "chartfile.yaml"))
			if err != nil {
				return err
			}
		}
		return err
	}

	if err := cf.Add([]string{helmChart.Requirement()}); err != nil {
		return err
	}

	libsonnetPath := filepath.Join(out, "/lib/", helmChart.Name+".libsonnet")
	render, err := templateString(helmChart.Name, jsonnetTpl)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(libsonnetPath, render, 0644); err != nil {
		return err
	}

	jsonnetfilepath := filepath.Join(out, "jsonnetfile.json")
	if _, err := os.Stat(jsonnetfilepath); err != nil {
		return nil
	} else {

		bundle, err := jsonnetfile.Load(jsonnetfilepath)
		if err != nil {
			return err
		}

		libDep := deps.Dependency{
			Source: deps.Source{
				GitSource: &deps.Git{
					Scheme: "https://",
					Host:   "github.com",
					User:   "grafana",
					Repo:   "jsonnet-libs",
					Subdir: "/helm-util",
				},
			},
			Version: "master",
		}
		if _, ok := bundle.Dependencies[libDep.Name()]; !ok {
			bundle.Dependencies[libDep.Name()] = libDep

			b, err := json.MarshalIndent(bundle, "", "  ")
			if err != nil {
				return err
			}
			b = append(b, []byte("\n")...)

			if err := ioutil.WriteFile(jsonnetfilepath, b, 0644); err != nil {
				return err
			}

			if err := os.Chdir(out); err != nil {
				return err
			}
			if err := exec.Command("jb", "install").Run(); err != nil {
				return err
			}
		}

		return nil
	}
}

func templateString(chart, s string) ([]byte, error) {
	tpl, err := template.New("jsonnet").Parse(s)
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	if err := tpl.Execute(&buf, struct {
		Name string
	}{Name: chart}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
