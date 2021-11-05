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
	"fmt"

	"github.com/redradrat/kable/pkg/ui"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

const uiAddressKey = "address"
const uiPortKey = "port"

// uiCmd represents the serve command
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Run kable UI",
	Long:  `Runs kable UI as a server and connects to a kable server.`,
	Example: `kable ui --ui-address 127.0.0.1 --ui-port 2020
	KABLE_UIADDRESS=127.0.0.1 KABLE_UIPORT kable ui`,
	Run: func(cmd *cobra.Command, args []string) {
		err := ui.StartUp(fmt.Sprintf("%s:%s", viper.Get(uiAddressKey), viper.Get(uiPortKey)))
		if err != nil {
			PrintError("error occurred while running server: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)

	uiCmd.Flags().StringP("ui-address", "a", "0.0.0.0", "The adress to bind the UI server to")
	uiCmd.Flags().StringP("ui-port", "p", "1323", "The port for the kable UI to listen on")
	viper.BindPFlag(uiAddressKey, uiCmd.Flags().Lookup("ui-address"))
	viper.BindPFlag(uiPortKey, uiCmd.Flags().Lookup("ui-port"))
}
