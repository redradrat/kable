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
	"strings"

	"github.com/redradrat/kable/pkg/repositories"

	"github.com/spf13/viper"

	"github.com/redradrat/kable/pkg/api"
	"github.com/spf13/cobra"
)

const serverAddressKey = "address"
const serverPortKey = "port"
const etcdEndpoints = "etcdEndpoints"
const etcdTimeout = "etcdTimeout"

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run kable as a server",
	Long:  `Runs kable as a server expecting payloads via a REST interface.`,
	Example: `kable serve --server-address 127.0.0.1 --server-port 2020
	KABLE_SERVERADDRESS=127.0.0.1 KABLE_SERVERPORT kable serve`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set(repositories.StoreKey, repositories.EtcdStoreConfigMap(strings.Split(viper.GetString(etcdEndpoints), ","), viper.GetDuration(etcdTimeout)).Map())
		api.StartUp(fmt.Sprintf("%s:%s", viper.Get(serverAddressKey), viper.Get(serverPortKey)))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("server-address", "a", "0.0.0.0", "The adress to bind the server to")
	serveCmd.Flags().StringP("server-port", "p", "1323", "The port for the kable api to listen on")
	serveCmd.Flags().String("etcd-endpoints", "", "The etcd endpoints to use")
	serveCmd.Flags().String("etcd-timeout", "5000", "The timeout for etcd interactions in milliseconds")
	viper.BindPFlag(serverAddressKey, serveCmd.Flags().Lookup("server-address"))
	viper.BindPFlag(serverPortKey, serveCmd.Flags().Lookup("server-port"))
	viper.BindPFlag(etcdEndpoints, serveCmd.Flags().Lookup("etcd-endpoints"))
	viper.BindPFlag(etcdTimeout, serveCmd.Flags().Lookup("etcd-timeout"))
}
