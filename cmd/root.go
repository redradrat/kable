// Copyright (c) 2020 Ralph Kühnert
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/redradrat/kable/pkg/repositories"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{}
var cfgFile string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kable/userconfig.json)")
	cobra.OnInitialize(initConfig)
	bindFlags(rootCmd, viper.GetViper())
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(repositories.KableDir)
		viper.SetConfigName("settings")
		viper.SetConfigType("json")
		viper.Set(repositories.StoreKey, repositories.LocalStoreConfigMap().Map())

		fullpath := filepath.Join(repositories.KableDir, "settings.json")
		if _, err := os.Stat(fullpath); err != nil {
			if os.IsNotExist(err) {
				err = viper.WriteConfigAs(fullpath)
				if err != nil {
					fmt.Println("Error writing config file:", err)
					os.Exit(1)
				}
			}
		}
	}

	viper.SetEnvPrefix("KABLE")

	viper.AutomaticEnv()

}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	//cmd.Flags().VisitAll(func(f *pflag.Flag) {
	//	// Environment variables can't have dashes in them, so bind them to their equivalent
	//	// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
	//	if strings.Contains(f.Name, "-") {
	//		envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
	//		v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
	//	}
	//
	//	// Apply the viper config value to the flag when the flag is not set and viper has a value
	//	if !f.Changed && v.IsSet(f.Name) {
	//		val := v.Get(f.Name)
	//		cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
	//	}
	//})
}
