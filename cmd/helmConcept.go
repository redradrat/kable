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

	"github.com/redradrat/kable/pkg/helm"

	"github.com/spf13/cobra"
)

var chartVersion string
var chartRepoName string
var chartRepoURL string
var dir string

// helmImportCmd represents the import command
var helmConceptCmd = &cobra.Command{
	Use:     "concept",
	Short:   "Create a concept, wrapping a helm chart from a git repo",
	Example: "kable helm concept --directory sentry sentry --repoName stable --version 4.3.0",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires exactly ONE argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		PrintMsg("Creating concept from helm chart '%s'...", args[0])
		if err := helm.InitHelmConcept(helm.HelmChart{Name: args[0], Version: chartVersion, Repo: helm.HelmRepo{Name: chartRepoName, URL: chartRepoURL}}, dir); err != nil {
			PrintError("unable to import helm chart: %s", err)
		}
	},
}

func init() {
	helmCmd.AddCommand(helmConceptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helmImportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	helmConceptCmd.Flags().StringVarP(&chartVersion, "version", "v", "", "The version of the helm chart.")
	helmConceptCmd.Flags().StringVar(&chartRepoName, "repo", "stable", "The name of the repository where the helm chart resides. (stable: https://charts.helm.sh/stable)")
	helmConceptCmd.Flags().StringVar(&chartRepoURL, "repoURL", "", "The URL of the repository where the helm chart resides.")
	helmConceptCmd.Flags().StringVarP(&dir, "directory", "d", ".", "The directory to create the concept in.")
}
