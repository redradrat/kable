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

	"github.com/redradrat/kable/pkg/kable"

	"github.com/spf13/cobra"
)

var outpath string
var targetType string

// createAppCmd represents the create command
var createAppCmd = &cobra.Command{
	Use:   "create [NAME] [CONCEPT@REPO]",
	Short: "A brief description of your command",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires exactly TWO arguments")
		}
		name := args[0]
		conceptIdentifier := args[1]

		if !kable.AppNameIsValid(name) {
			PrintError("invalid name given: %s", args[0])
		}

		if !kable.IsValidConceptIdentifier(conceptIdentifier) {
			PrintError("invalid concept identifier given: %s", args[1])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		PrintMsg("Creating AppV1...")

		name := args[0]

		// First let's get our concept... maybe it doesn't even exist, hm? meow
		conceptIdentifier := kable.ConceptIdentifier(args[1])
		PrintMsg("Fetching Concept '%s'...", args[0])
		cpt, err := kable.GetConcept(conceptIdentifier)
		if err != nil {
			PrintError("unable to get specified concept: %s", err)
		}

		// Run dialog to get values for concept inputs
		avs, err := NewInputDialog(cpt.Inputs).RunInputDialog()
		if err != nil {
			PrintError("error processing concept inputs: %s", err)
		}

		// Now let's render our app
		PrintMsg("Rendering AppV1...")
		app, err := kable.NewAppV1(name, avs)
		if err != nil {
			PrintError("unable to render app: %s", err)
		}

		if err := kable.RenderApp(app, conceptIdentifier, outpath, kable.YamlTarget{}); err != nil {
			PrintError("unable to render app: %s", err)
		}

		PrintSuccess("Successfully created app at: %s", outpath)
	},
}

func init() {
	appCmd.AddCommand(createAppCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createAppCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	createAppCmd.Flags().StringVarP(&outpath, "output", "o", ".", "The output directory this app will be placed in")
	createAppCmd.Flags().StringVarP(&targetType, "targetType", "t", "yaml", "The target format, this AppV1 will be rendered as")
}
