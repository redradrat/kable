/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"

	"github.com/redradrat/kable/pkg/kable/concepts"

	"github.com/spf13/cobra"
)

var outpath string
var targetType string

// renderConceptCmd represents the create command
var renderConceptCmd = &cobra.Command{
	Use:   "render [CONCEPT@REPO] [NAME]",
	Short: "RenderMeta a given concept",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires exactly THREE arguments")
		}
		conceptIdentifier := args[0]
		name := args[1]

		if !concepts.RenderNameIsValid(name) {
			PrintError("invalid name given: %s", name)
		}

		if !concepts.IsValidConceptIdentifier(conceptIdentifier) {
			PrintError("invalid concept identifier given: %s", conceptIdentifier)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[1]
		conceptIdentifier := concepts.ConceptIdentifier(args[0])

		// First let's get our concept... maybe it doesn't even exist, hm? meow
		PrintMsg("Fetching Concept '%s'...", conceptIdentifier.String())
		cpt, err := concepts.GetConcept(conceptIdentifier)
		if err != nil {
			PrintError("unable to get specified concept: %s", err)
		}

		// Run dialog to get values for concept inputs
		avs, err := NewInputDialog(cpt.Inputs).RunInputDialog()
		if err != nil {
			PrintError("error processing concept inputs: %s", err)
		}

		// Now let's render our app
		PrintMsg("Rendering concept...")
		app, err := concepts.NewRenderV1(name, avs)
		if err != nil {
			PrintError("unable to render concept: %s", err)
		}

		bundle, err := concepts.RenderConcept(app, conceptIdentifier, outpath, concepts.TargetType(targetType))
		if err != nil {
			PrintError("unable to render concept: %s", err)
		}
		if err := bundle.Write(); err != nil {
			PrintError("unable to write rendered concept to file system: %s", err)
		}
		PrintSuccess("Successfully created concept!", outpath)
	},
}

func init() {
	conceptCmd.AddCommand(renderConceptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renderConceptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	renderConceptCmd.Flags().StringVarP(&outpath, "output", "o", ".", "The output directory this app will be placed in")
	renderConceptCmd.Flags().StringVarP(&targetType, "targetType", "t", string(concepts.YamlTargetType), "The target format, this ConceptRenderV1 will be rendered as")
}
