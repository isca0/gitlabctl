package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the gitlabctl version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gitlabctl version " + Version)
		fmt.Println("Author: Igor Brandao aka Isca")
	},
}
