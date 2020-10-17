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
	"net/url"
	"regexp"

	"github.com/redradrat/kable/pkg/repositories"

	"github.com/spf13/cobra"
)

var repoUser string
var repoPass string

// addRepoCmd represents the add command
var addRepoCmd = &cobra.Command{
	Use:   "add [ID] [URL]",
	Short: "Add a repository to current config",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires exactly TWO arguments")
		}
		name := args[0]
		repoUrl := args[1]

		rxp := regexp.MustCompile("^[a-zA-Z]+$").MatchString
		if !rxp(name) {
			PrintError("invalid name given: %s", args[0])
		}

		if _, err := url.Parse(repoUrl); err != nil {
			PrintError("invalid URL given: %s", args[1])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		PrintMsg("Fetching repository...")
		name := args[0]
		repoUrl := args[1]

		var mods []repositories.RegistryModification
		authExists, err := repositories.RepoAuthExists(repoUrl)
		if err != nil {
			PrintError("unable to check configured auths: %s", err)
		}
		if !authExists {
			newAuth, user, pw, err := RunAuthDialog(repoUser, repoPass)
			if err != nil {
				PrintError("unable to display authentication dialog: %s", err)
			}
			if newAuth {
				storemod, err := repositories.StoreRepoAuth(repoUrl, repositories.AuthPair{Username: user, Password: pw})
				if err != nil {
					PrintError("unable to store authentication data: %s", err)
				}
				mods = append(mods, storemod)
			}
		}

		mod := repositories.AddRepository(repositories.Repository{
			Name: name,
			GitRepository: repositories.GitRepository{
				URL:    repoUrl,
				GitRef: "refs/heads/master",
			},
		})
		mods = append(mods, mod)
		err = repositories.UpdateRegistry(mods...)
		if err != nil {
			PrintError("unable to update registry: %s", err)
		}
		PrintSuccess("Successfully added repository!")
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
	addRepoCmd.Flags().StringVarP(&repoUser, "username", "u", "", "The username for this repository.")
	addRepoCmd.Flags().StringVarP(&repoPass, "password", "p", "", "The password for this repository.")
}
