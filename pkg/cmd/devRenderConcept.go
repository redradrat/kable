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
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/redradrat/kable/pkg/kable/concepts"

	"github.com/spf13/cobra"
)

var devRenderTargetType string

// devRenderConceptCmd represents the create command
var devRenderConceptCmd = &cobra.Command{
	Use:   "render [PATH]",
	Short: "render a local concept",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires exactly ONE arguments")
		}
		path := args[0]

		f, err := os.Stat(path)
		if os.IsNotExist(err) {
			PrintError("given path does not exist")
		}
		if !f.IsDir() {
			PrintError("given path is not a directory")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		PrintMsg("Trying to render concept at '%s'", path)

		// First let's get our concept... maybe it doesn't even exist, hm? meow
		cpt, err := concepts.GetConcept(path)
		if err != nil {
			PrintError("unable to get specified concept: %s", err)
		}

		// Run dialog to get values for concept inputs
		avs, err := NewInputDialog(cpt.Inputs).RunInputDialog()
		if err != nil {
			PrintError("error processing concept inputs: %s", err)
		}

		// Let's try to render our concept... maybe it doesn't even exist, hm? meow
		bundle, err := concepts.RenderConcept(path, avs, concepts.TargetType(devRenderTargetType))
		if err != nil {
			PrintError("unable to render concept: %s", err)
		}

		boldNUnderline := color.New(color.FgCyan, color.Bold, color.Underline)
		fmt.Println()
		boldNUnderline.Println("Output:")
		for _, file := range bundle.Files {
			PrintMsg(file.String())
		}
	},
}

func init() {
	devCmd.AddCommand(devRenderConceptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// devRenderConceptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	devRenderConceptCmd.Flags().StringVarP(&devRenderTargetType, "targetType", "t", string(concepts.YamlTargetType), "The target format, this concept will be rendered as")
}
