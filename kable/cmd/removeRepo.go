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
	"regexp"

	"github.com/redradrat/kable/pkg/repositories"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove [ID]",
	Aliases: []string{"delete"},
	Short:   "Removes a repository from current config",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires exactly one arguments")
		}

		name := args[0]

		rxp := regexp.MustCompile("^[a-zA-Z]+$").MatchString
		if !rxp(name) {
			PrintError("invalid name given: %s", args[0])
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		PrintMsg("Removing repository...")
		if err := repositories.RemoveRepository(args[0]); err != nil {
			PrintError("unable to remove repository: %s \n", err)
		}
		PrintSuccess("Successfully removed repository!")
	}}

func init() {
	repoCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
