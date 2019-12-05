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
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:       "rm args",
	Short:     "Deletes a group or a project",
	Long:      `Delete a group or a project on a received target.`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"group", "proj"},
	Example: `  gitlabctl rm group --from sessionA:/some/group
  gitlabctl rm proj --from sessionB:/full/path/project

Args:
  proj:  - will delete projects.
  group: - will delete groups.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := &http.Client{
			Timeout: time.Second * 30,
		}
		runRm(args, from, client)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.SetArgs([]string{"group", "proj"})
	rmCmd.Flags().StringVarP(&from, "from", "f", "", "sepecifies the session token and full path as the source.(required)")
	rmCmd.MarkFlagRequired("from")
}

// runRm executes the rm command.
func runRm(args []string, from string, client *http.Client) {
	switch {
	case args[0] == "proj":
		fmt.Println("will remove a project in a future release.")
	case args[0] == "group":
		g := new(Groups)
		g.rm(from, client)
	}
}
