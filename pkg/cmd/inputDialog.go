package cmd

import (
	"fmt"

	"github.com/redradrat/kable/pkg/kable/concepts"

	"github.com/manifoldco/promptui"
)

type InputDialog struct {
	inputs concepts.ConceptInputs
}

func NewInputDialog(inputs concepts.ConceptInputs) InputDialog {
	return InputDialog{inputs: inputs}
}

func (id InputDialog) RunInputDialog() (*concepts.RenderValues, error) {
	values := concepts.RenderValues{}
	PrintMsg("Mandatory Values")
	for id, input := range id.inputs.Mandatory {
		value, err := getValue(id, input)
		if err != nil {
			return nil, err
		}
		values[id] = value
	}

	PrintMsg("Optional Values")
	for id, input := range id.inputs.Optional {
		value, err := getValue(id, input)
		if err != nil {
			return nil, err
		}
		values[id] = value
	}

	return &values, nil
}

func getValue(name string, input concepts.InputType) (concepts.ValueType, error) {
	var value concepts.ValueType
	switch input.Type {
	case concepts.ConceptStringInputType:
		prompt := promptui.Prompt{
			Label: name,
		}

		result, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		value = concepts.StringValueType(result)
	case concepts.ConceptSelectionInputType:
		prompt := promptui.Select{
			Label: name,
			Items: input.Options,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		value = concepts.SelectValueType(result)
	default:
		return nil, fmt.Errorf("input type not supported")
	}
	return value, nil
}
