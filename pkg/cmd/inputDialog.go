package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"

	"github.com/redradrat/kable/pkg/kable/concepts"
)

type InputDialog struct {
	inputs concepts.ConceptInputs
}

func NewInputDialog(inputs concepts.ConceptInputs) InputDialog {
	return InputDialog{inputs: inputs}
}

func (id InputDialog) RunInputDialog() (*concepts.RenderValues, error) {
	values := concepts.RenderValues{}
	if len(id.inputs.Mandatory) != 0 {
		PrintMsg("Mandatory Values")
		for id, input := range id.inputs.Mandatory {
			value, err := getValue(id, input)
			if err != nil {
				return nil, err
			}
			values[id] = value
		}
	}

	if len(id.inputs.Optional) != 0 {
		PrintMsg("Optional Values")
		for id, input := range id.inputs.Optional {
			value, err := getValue(id, input)
			if err != nil {
				return nil, err
			}
			values[id] = value
		}
	}

	return &values, nil
}

func getValue(name string, input concepts.InputType) (concepts.ValueType, error) {

	var value concepts.ValueType
	switch input.Type {
	case concepts.ConceptStringInputType:

		val := ""
		prompt := &survey.Input{
			Message: name,
		}
		if err := survey.AskOne(prompt, &val); err != nil {
			return nil, err
		}

		value = concepts.StringValueType(val)
	case concepts.ConceptSelectionInputType:

		val := ""
		prompt := &survey.Select{
			Message: name,
			Options: input.Options,
		}
		if err := survey.AskOne(prompt, &val); err != nil {
			return nil, err
		}

		value = concepts.SelectValueType(val)
	default:
		return nil, fmt.Errorf("input type not supported")
	}
	return value, nil
}
