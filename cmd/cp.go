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

	"github.com/spf13/cobra"
)

var (

	//used for flags
	from, to string

	// cpCmd represents the cp command
	cpCmd = &cobra.Command{
		Use:       "cp args",
		Short:     "Copy groups and projects",
		Long:      `Copy group and projects from one session to another.`,
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"group", "proj"},
		Example: `  gitlabctl cp group --from sessionA:fullPath --to sessionB:fullPath
  gitlabctl cp proj --from sessionA:fullPath --to sessionB:fullPath

Args:
  group	- will copy groups.
  proj	- will copy projects.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("cp called")
		},
	}
)

func init() {
	rootCmd.AddCommand(cpCmd)
	cpCmd.SetArgs([]string{"group", "proj"})
	cpCmd.Flags().StringVarP(&from, "from", "f", "", "specifies the session token + full group path as the source. (required)")
	cpCmd.MarkFlagRequired("from")
	cpCmd.Flags().StringVarP(&to, "to", "t", "", "specifies the session token + full group path as the destination. (required)")
	cpCmd.MarkFlagRequired("to")

}

// runCopy executes the cp command by treating received arguments(group or projects).
func runCopy(args []string, token string, client *http.Client) {
	switch {
	case args[0] == "group":
		//groupCopy(args, tokenA, client)
		fmt.Println("groupCopy")
	case args[0] == "proj":
		fmt.Println("projectCopy")
	}
}
