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
	"net/url"
	"os"

	"github.com/redradrat/kable/pkg/kable"

	"github.com/spf13/cobra"
)

// addRepoCmd represents the add command
var addRepoCmd = &cobra.Command{
	Use:   "add [URL]",
	Short: "Add a repository to current config",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires exactly one arguments")
		}
		repoUrl := args[0]

		if _, err := url.Parse(repoUrl); err != nil {
			return fmt.Errorf("invalid URL given: %s", args[1])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		PrintMsg("Fetching repository...")
		repoUrl := args[0]
		name, err := kable.AddRepository(repoUrl, "master")
		if err != nil {
			if err.Error() != kable.RepositoryAlreadyExistsError {
				PrintError("unable to add repository: %s", err)
				os.Exit(1)
			} else {
				PrintSuccess("Repository already configured! [%s]", name)
				os.Exit(0)
			}
		}
		PrintSuccess("Successfully added repository! [%s]", name)
	},
}

func init() {
	repoCmd.AddCommand(addRepoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
