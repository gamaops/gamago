/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"os"
	"time"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var esAddresses []string
var esUsername string
var esPassword string
var esMaxIdleConns int
var esDialTimeout time.Duration
var esResponseHeaderTimeout time.Duration
var prettyLog bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "elastion",
	Short: "elastion is a command line utility to manage Elasticsearch opertions",
	Long:  `elastion is a command line utility to manage Elasticsearch opertions`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.elastion.yaml)")

	rootCmd.PersistentFlags().StringSliceVar(&esAddresses, "addresses", []string{"http://localhost:9200"}, "addresses address1,address2 (default is http://localhost:9200)")
	rootCmd.PersistentFlags().StringVar(&esUsername, "username", "", "username myUserName")
	rootCmd.PersistentFlags().StringVar(&esPassword, "password", "", "password myPassword")
	rootCmd.PersistentFlags().IntVar(&esMaxIdleConns, "max-idle-conns", 10, "max-idle-conns 10")
	rootCmd.PersistentFlags().DurationVar(&esDialTimeout, "dial-timeout", 5*time.Second, "dial-timeout 5s")
	rootCmd.PersistentFlags().DurationVar(&esResponseHeaderTimeout, "res-header-timeout", 5*time.Second, "res-header-timeout 5s")
	rootCmd.PersistentFlags().BoolVar(&prettyLog, "pretty-log", true, "pretty-log true")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".elastion" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".elastion")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	setup()
}