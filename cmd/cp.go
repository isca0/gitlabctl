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
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var (

	//used for flags
	from, to    string
	proj, group bool

	// cpCmd represents the cp command
	cpCmd = &cobra.Command{
		Use:   "cp",
		Short: "Copy groups and projects",
		Long:  `Copy group and projects from one session to another.`,
		Args:  cobra.NoArgs,
		Example: `  gitlabctl cp -p --from sessionA:fullPath --to sessionB:fullPath
  gitlabctl cp -g --from sessionA:fullPath --to sessionB:fullPath`,

		Run: func(cmd *cobra.Command, args []string) {
			client := &http.Client{
				Timeout: time.Second * 30,
			}
			runCopy(from, to, client)
		},
	}
)

func init() {
	rootCmd.AddCommand(cpCmd)
	cpCmd.Flags().BoolP("proj", "p", false, "use to copy projects.")
	cpCmd.Flags().BoolP("group", "g", false, "use to copy groups.")
	cpCmd.Flags().StringVarP(&from, "from", "f", "", "specifies the session token + full group path as the source. (required)")
	cpCmd.MarkFlagRequired("from")
	cpCmd.Flags().StringVarP(&to, "to", "t", "", "specifies the session token + full group path as the destination. (required)")
	cpCmd.MarkFlagRequired("to")

}

// runCopy executes the cp command by treating received arguments(group or projects).
func runCopy(from, to string, client *http.Client) {
	switch {
	case proj:
		g := groupPages{}
		g.copy(from, to, client)
	case group:
		p := projectPages{}
		p.copy(from, to, client)
	}
}
