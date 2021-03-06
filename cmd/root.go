/*
Copyright © 2019 Igor Brandao <igorsca at protonmail dot com>

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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gitlabctl",
		Short: "A CLI tool to manage your Gitlab project and groups",
		Long: `gitlabctl is a gitlab CLI tool. 
It can operate gitlab to copy groups or migrate projects from one gitlab account to another.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AutomaticEnv()

	if !viper.IsSet("sessiona") {
		fmt.Println(`Ups... You must declare at least one valid session token. Please declare that
by exporting: export SESSIONA="myToken"`)
		os.Exit(1)
	}

}
