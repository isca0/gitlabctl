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
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	target string
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls args",
	Short: "List projects or groups",
	Long:  `List projects or groups from the received target.`,
	Example: `  gitlabctl ls group --target $SESSIONA
  gitlabctl ls proj --target $SESSIONB

Args:
  group	- will list groups.
  proj	- will list projects.`,
	Args:      cobra.MaximumNArgs(2),
	ValidArgs: []string{"group", "proj"},
	Run: func(cmd *cobra.Command, args []string) {
		client := &http.Client{
			Timeout: time.Second * 30,
		}
		runList(args, target, client)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.PersistentFlags().StringVarP(&target, "target", "t", os.Getenv("SESSIONA"), "specifies the target to be listed.")
	lsCmd.SetArgs([]string{"group", "proj"})
}

// runList executes the ls command by treating received arguments.
func runList(arg []string, token string, client *http.Client) {
	switch {
	case arg[0] == "group":
		groupList(arg[1], token, client)
	case arg[0] == "proj":
		p := projectPages{}
		p.list(client, getProj, token)
		//for _, prj := range projects.Project {
		//	fmt.Println(prj[0].Name)
		//}

	}
}
