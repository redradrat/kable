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
	"github.com/redradrat/kable/pkg/kable"
	"github.com/spf13/cobra"
)

// listReposCmd represents the list command
var listReposCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all concept repositories in the current config",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := kable.GetConfig()
		if err != nil {
			PrintError("cannot list repositories: %s \n", err)
		}
		repoArray, err := cfg.Repositories.ToArray()
		if err != nil {
			PrintError("cannot list repositories: %s \n", err)
		}
		PrintTable([]string{"Name", "URL", "ID", "Initialized"}, repoArray...)
	},
}

func init() {
	repoCmd.AddCommand(listReposCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listReposCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listReposCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}