package helm

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/jsonnet-bundler/jsonnet-bundler/spec/v1/deps"

	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"

	"github.com/otiai10/copy"
	"github.com/redradrat/kable/kable/config"

	"github.com/redradrat/kable/kable/errors"
)

const jsonnetTpl = `local helm = (import "github.com/grafana/jsonnet-libs/helm-util/helm.libsonnet").new(std.thisFile);

{
  {{.Name}}: helm.template("{{.Name}}", "./charts/{{.Name}}", {
    namespace: "default",
    values: {
      foo: { bar: baz }
    }
  })
}
`

type HelmChart struct {
	URL    string
	Subdir *string
}

func ImportHelmChart(chart HelmChart, out string) error {
	basedir := config.CacheDir
	path, err := cloneAndCheckRepo(chart.URL, basedir)
	if err != nil {
		return err
	}

	var chartPath string
	if chart.Subdir != nil {
		chartPath = filepath.Join(path, *chart.Subdir)
	} else {
		chartPath = path
	}

	if !isHelmRepo(chartPath) {
		return errors.NotHelmChartError
	}

	chartName := filepath.Base(chartPath)
	if err := copy.Copy(chartPath, filepath.Join(out, "/charts/", chartName)); err != nil {
		return err
	}

	tpl, err := template.New("jsonnet").Parse(jsonnetTpl)
	if err != nil {
		return err
	}
	buf := bytes.Buffer{}
	if err := tpl.Execute(&buf, struct {
		Name string
	}{Name: chartName}); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(out, "/lib/", chartName+".jsonnet"), buf.Bytes(), os.ModePerm); err != nil {
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
		name := libDep.Name()
		if _, ok := bundle.Dependencies[name]; !ok {
			bundle.Dependencies[name] = libDep

			b, err := json.MarshalIndent(bundle, "", "  ")
			if err != nil {
				return err
			}
			b = append(b, []byte("\n")...)

			if err := ioutil.WriteFile(jsonnetfilepath, b, 0644); err != nil {
				return err
			}

			if err := exec.Command("jb", "install").Run(); err != nil {
				return err
			}
		}

		return nil
	}
}

func cloneAndCheckRepo(gitUrl, dir string) (string, error) {
	uri, err := url.Parse(gitUrl)
	if err != nil {
		return "", err
	}

	if !(uri.Scheme == "http" || uri.Scheme == "https") {
		return "", errors.UnsupportedURISchemeError
	}

	path := filepath.Join(dir, filepath.Base(uri.Path))
	if err := Checkout(gitUrl, path); err != nil {
		return "", err
	}

	return path, nil
}

func isHelmRepo(path string) bool {
	_, err := os.Stat(filepath.Join(path, "Chart.yaml"))
	return err == nil
}
