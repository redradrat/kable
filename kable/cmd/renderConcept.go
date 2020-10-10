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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/redradrat/kable/kable/concepts"

	"github.com/spf13/cobra"
)

var outpath string
var conceptRenderTargetType string
var local bool
var single bool
var renderinfo string
var printOnly bool

// renderConceptCmd represents the create command
var renderConceptCmd = &cobra.Command{
	Use:   "render [PATH]",
	Short: "Render a concept",
	Example: `
kable render my/concept@myrepo
kable render -l . -o out/
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires exactly ONE argument")
		}

		conceptIdentifier := args[0]
		if local {
			f, err := os.Stat(conceptIdentifier)
			if os.IsNotExist(err) {
				PrintError("given path does not exist")
			}
			if !f.IsDir() {
				PrintError("given path is not a directory")
			}
		} else {
			if conceptIdentifier != "." && !concepts.IsValidConceptIdentifier(conceptIdentifier) {
				PrintError("invalid concept identifier given: %s", conceptIdentifier)
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		conceptIdentifier := concepts.ConceptIdentifier(args[0])

		silent = printOnly
		// We only want to print in singlefile mode
		single = printOnly

		// If local let's get our concept from here, otherwise from the cache
		var cpt *concepts.Concept
		var err error
		if !local {
			// ... maybe it doesn't even exist, hm? meow
			PrintMsg("Fetching Concept '%s'...", conceptIdentifier.String())
			cpt, err = concepts.GetRepoConcept(conceptIdentifier)
		} else {
			cpt, err = concepts.GetConcept(conceptIdentifier.String())
		}
		if err != nil {
			PrintError("unable to get specified concept: %s", err)
		}

		// check if existing RenderInfo exists, or run dialog to get values for concept inputs
		var avs *concepts.RenderValues
		existingRenderInfo := true
		outdatedValues := false
		var ri *concepts.RenderInfoV1
		if renderinfo != "" {
			ri, err = concepts.ParseRenderInfoV1FromFile(renderinfo)
		} else {
			ri, err = concepts.ParseRenderInfoV1FromFile(filepath.Join(outpath, concepts.ConceptRenderFileName))
		}
		if err != nil {
			if os.IsNotExist(err) {
				existingRenderInfo = false
			} else {
				PrintError("error parsing existing renderinfo: %s", err)
			}
		} else {
			vals := *ri.Values
			for k, _ := range cpt.Inputs.Mandatory {
				if _, ok := vals[k]; !ok {
					outdatedValues = true
				}
			}
			avs = &vals
		}

		// Ask for values if renderinfo does not exist or values are outdated
		if !existingRenderInfo || outdatedValues {
			PrintMsg("No existing render or outdated values detected...")
			avs, err = NewInputDialog(cpt.Inputs).RunInputDialog()
			if err != nil {
				PrintError("error processing concept inputs: %s", err)
			}
		}

		// Now let's render our app
		PrintMsg("Rendering concept...")
		var bundle *concepts.Render
		bundle, err = concepts.RenderConcept(conceptIdentifier.String(), avs, concepts.TargetType(conceptRenderTargetType), concepts.RenderOpts{Single: single, Local: local, WriteRenderInfo: renderinfo == ""})
		if err != nil {
			PrintError("unable to render concept: %s", err)
		}

		if !printOnly {
			_, err := ioutil.ReadDir(outpath)
			if err != nil && !os.IsNotExist(err) {
				PrintError("unable to read directory '%s' for rendering: %s", outpath, err)
			}
		}

		if printOnly {
			fmt.Print(bundle.PrintFiles())
		} else {
			if renderinfo == "" {
				if err := bundle.WriteInfo(outpath); err != nil {
					PrintError("unable to write renderinfo to file system: %s", err)
				}
			}
			if err := bundle.WriteFiles(outpath); err != nil {
				PrintError("unable to write rendered concept to file system: %s", err)
			}
			PrintSuccess("Successfully created concept!")
		}

	},
}

func init() {
	rootCmd.AddCommand(renderConceptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renderConceptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	renderConceptCmd.Flags().StringVarP(&outpath, "output", "o", ".", "The output directory this app will be placed in")
	renderConceptCmd.Flags().BoolVarP(&local, "local", "l", false, "Whether to read the concept from a local path")
	renderConceptCmd.Flags().BoolVarP(&single, "single", "s", false, "Render into a single manifest.yaml file")
	renderConceptCmd.Flags().StringVarP(&renderinfo, "renderinfo", "r", "", "Path to an existing renderinfo to use. (skips writing renderinfo.json)")
	renderConceptCmd.Flags().StringVarP(&conceptRenderTargetType, "targetType", "t", string(concepts.YamlTargetType), "The target format, this concept will be rendered as")
	renderConceptCmd.Flags().BoolVarP(&printOnly, "print", "p", false, "Runs silent and prints manifests to stdout. (renderinfo.json needs to exist)")
}
