package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"

	"github.com/fatih/color"

	"github.com/AlecAivazis/survey/v2"

	"github.com/redradrat/kable/pkg/concepts"
)

var (
	cursor *terminal.Cursor
	hl     = color.New(color.Bold, color.Underline).Sprintf
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
		PrintMsg(hl("\nMandatory Values"))
		keys := getSortedMapKeys(id.inputs.Mandatory)
		for _, key := range keys {
			value, err := getValue(key, id.inputs.Mandatory[key])
			if err != nil {
				return nil, err
			}
			values[key] = value
		}
	}

	if len(id.inputs.Optional) != 0 {

		optConfirm := false
		optPrompt := &survey.Confirm{
			Message: "Provide values for optional inputs?",
		}
		if err := survey.AskOne(optPrompt, &optConfirm); err != nil {
			return nil, err
		}
		erasePreviousLine()

		// Only go into optional values if the users want's to
		if optConfirm {
			PrintMsg(hl("\nOptional Values"))
			keys := getSortedMapKeys(id.inputs.Optional)
			for _, key := range keys {
				valConfirm := false
				valPrompt := &survey.Confirm{
					Message: fmt.Sprintf("Provide value for %s?", key),
					Help:    breakEvery60chars(id.inputs.Optional[key].Description),
				}
				if err := survey.AskOne(valPrompt, &valConfirm); err != nil {
					return nil, err
				}
				if valConfirm {
					erasePreviousLine()
					value, err := getValue(key, id.inputs.Optional[key])
					if err != nil {
						return nil, err
					}
					values[key] = value
				}
			}
		}
	}

	fmt.Println()
	return &values, nil
}

func getSortedMapKeys(values map[string]concepts.InputType) []string {
	var keys []string
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func getValue(name string, input concepts.InputType) (concepts.ValueType, error) {

	var helpText string
	if input.Description != "" {
		helpText = helpText + breakEvery60chars(input.Description) + "\n\n"
	}
	if input.Example != "" {
		helpText = helpText + "Example: " + input.Example
	}

	var value concepts.ValueType
	switch input.Type {
	case concepts.ConceptStringInputType:
		if helpText == "" {
			helpText = "example: test"
		}
		val := ""
		prompt := &survey.Input{
			Message: name + color.CyanString(" (string)"),
			Help:    helpText,
		}
		if err := survey.AskOne(prompt, &val); err != nil {
			return nil, err
		}

		value = concepts.StringValueType(val)
	case concepts.ConceptIntInputType:
		if helpText == "" {
			helpText = "example: 3"
		}
		var val int
		prompt := &survey.Input{
			Message: name + color.CyanString(" (integer)"),
			Help:    helpText,
		}
		if err := survey.AskOne(prompt, &val); err != nil {
			return nil, err
		}

		value = concepts.IntValueType(val)
	case concepts.ConceptBoolInputType:
		var val bool
		prompt := &survey.Confirm{
			Message: name + color.CyanString(" (boolean)"),
			Help:    helpText,
		}
		if err := survey.AskOne(prompt, &val); err != nil {
			return nil, err
		}

		value = concepts.BoolValueType(val)
	case concepts.ConceptSelectionInputType:
		var val string
		prompt := &survey.Select{
			Message: name + color.CyanString(" (select)"),
			Options: input.Options,
			Help:    helpText,
		}
		if err := survey.AskOne(prompt, &val); err != nil {
			return nil, err
		}

		value = concepts.StringValueType(val)
	case concepts.ConceptMapInputType:
		if helpText == "" {
			helpText = "example: {'foo':'bar'}"
		}
		val := ""
		prompt := &survey.Input{
			Message: name + color.CyanString(" (map)"),
			Help:    helpText,
		}
		if err := survey.AskOne(prompt, &val); err != nil {
			return nil, err
		}

		outmap := map[string]interface{}{}
		if err := json.Unmarshal([]byte(val), &outmap); err != nil {
			return nil, errors.New(fmt.Sprintf("unable to parse given map input: %s", err.Error()))
		}
		value = concepts.MapValueType(outmap)
	default:
		return nil, fmt.Errorf("input type not supported")
	}
	return value, nil
}

func breakEvery60chars(in string) string {
	if len(in) <= 60 {
		return in
	}
	var chunks []string
	chunk := make([]rune, 60)
	len := 0
	for _, r := range in {
		chunk[len] = r
		len++
		if len == 60 {
			chunks = append(chunks, string(chunk))
			len = 0
		}
	}
	if len > 0 {
		chunks = append(chunks, string(chunk[:len]))
	}
	return strings.Join(chunks, "\n")
}

func erasePreviousLine() {
	if cursor == nil {
		cursor = &terminal.Cursor{
			In:  os.Stdin,
			Out: os.Stdout,
		}
	}
	cursor.PreviousLine(1)
	terminal.EraseLine(os.Stdout, terminal.ERASE_LINE_ALL)
}
