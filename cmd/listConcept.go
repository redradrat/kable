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
	"github.com/redradrat/kable/pkg/concepts"
	"github.com/spf13/cobra"
)

// listConceptsCmd represents the list command
var listConceptsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available concepts",
	Run: func(cmd *cobra.Command, args []string) {
		cis, err := concepts.ListConcepts()
		if err != nil {
			PrintError("unable to list concepts: %s", err)
		}
		var outList [][]string
		for _, ci := range cis {
			c, err := concepts.GetRepoConcept(ci)
			if err != nil {
				PrintError("error getting concept '%s': %s", ci.String(), err)
			}
			outList = append(outList, []string{
				ci.Concept(),
				ci.Repo(),
				c.Meta.Maintainer.String(),
			})
		}
		PrintTable([]string{"ID", "Repository", "Maintainer"}, outList...)
	},
}

func init() {
	rootCmd.AddCommand(listConceptsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listReposCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listReposCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
