package concepts

import (
	"fmt"
	"sort"
	"strings"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/process"
	"github.com/grafana/tanka/pkg/tanka"

	"github.com/redradrat/kable/pkg/errors"
)

const (
	YamlTargetType TargetType = "yaml"
	CRDTargetType  TargetType = "crd"
)

type TargetType string

// Target is the interface for all Target implementations
type Target interface {
	TargetName() string
	Render(path string, vals *RenderValues, cpt ConceptType, single bool) (*Render, error)
}

type CRDTarget struct {
}

func (c CRDTarget) TargetName() string {
	return string(CRDTargetType)
}

func (c CRDTarget) Render(path string, vals *RenderValues, cpt ConceptType, single bool) (*Render, error) {
	panic("implement me")
}

type YamlTarget struct {
}

func (y YamlTarget) TargetName() string {
	return string(YamlTargetType)
}

func (y YamlTarget) Render(path string, vals *RenderValues, cpt ConceptType, single bool) (*Render, error) {
	var err error
	bundle := Render{}

	switch cpt {
	case ConceptJsonnetType:
		bundle.Files, err = renderJsonnetConcept(path, vals, single)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.ConceptTypeUnsupportedError
	}

	return &bundle, nil
}

func renderJsonnetConcept(path string, avs *RenderValues, single bool) ([]File, error) {
	opts := tanka.Opts{}

	if avs != nil {
		if opts.ExtCode == nil {
			opts.ExtCode = make(jsonnet.InjectedCode)
		}
		if opts.TLACode == nil {
			opts.TLACode = make(jsonnet.InjectedCode)
		}
		for id, val := range *avs {
			switch val.(type) {
			case StringValueType:
				opts.ExtCode[id] = fmt.Sprintf(`"%s"`, val.String())
				opts.TLACode[id] = fmt.Sprintf(`"%s"`, val.String())
			case MapValueType, IntValueType, BoolValueType:
				opts.ExtCode[id] = val.String()
				opts.TLACode[id] = val.String()
			default:
				return nil, errors.ValueTypeNotSupported
			}
		}
	}

	raw, err := tanka.Eval(path, opts)
	if err != nil {
		return nil, err
	}

	// Use Tanka's extract
	extract, err := process.Extract(raw)
	if err != nil {
		return nil, err
	}
	if err := process.Unwrap(extract); err != nil {
		return nil, err
	}

	var bundle []File
	out := make(manifest.List, 0, len(extract))
	for _, m := range extract {
		out = append(out, m)
		bundle = append(bundle, File{
			path:    fmt.Sprintf("%s_%s_%s.yaml", strings.ReplaceAll(m.APIVersion(), "/", "-"), m.Kind(), m.Metadata().Name()),
			content: []byte(m.String()),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return bundle[i].path < bundle[j].path
	})

	singlebundle := []File{{
		path:    "manifest.yaml",
		content: []byte(out.String()),
	}}

	if single {
		return singlebundle, nil
	}

	return bundle, nil
}
