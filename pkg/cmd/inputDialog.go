package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/redradrat/kable/pkg/kable"
)

type InputDialog struct {
	inputs kable.ConceptInputs
}

func NewInputDialog(inputs kable.ConceptInputs) InputDialog {
	return InputDialog{inputs: inputs}
}

func (id InputDialog) RunInputDialog() (*kable.AppValues, error) {
	values := kable.AppValues{}
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

func getValue(name string, input kable.InputType) (kable.ValueType, error) {
	var value kable.ValueType
	switch input.Type {
	case kable.ConceptStringInputType:
		prompt := promptui.Prompt{
			Label: name,
		}

		result, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		value = kable.StringValueType(result)
	case kable.ConceptSelectionInputType:
		prompt := promptui.Select{
			Label: name,
			Items: input.Options,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		value = kable.SelectValueType(result)
	default:
		return nil, fmt.Errorf("input type not supported")
	}
	return value, nil
}
