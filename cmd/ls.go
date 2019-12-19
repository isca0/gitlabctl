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
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	target string
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls arg",
	Short: "List groups and the inner projects",
	Long: `Without arguments list all groups and projects, if received a valid group will
	list only this group`,
	Example: `  gitlabctl ls --target sessionA
  gitlabctl ls group/subgroup --target sessionB`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := &http.Client{
			Timeout: time.Second * 30,
		}
		runList(args, viper.GetString(target), client)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.PersistentFlags().StringVarP(&target, "target", "t", "SESSIONA", "specifies the target to be listed.")
}

// runList executes the ls command by treating received arguments.
func runList(args []string, token string, client *http.Client) {
	g := new(Groups)
	g.list(args, token, client)
}
