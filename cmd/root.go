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

	switch {
	case !viper.IsSet("sessiona"):
		fmt.Println(`Ups... You must declare at least one valid session token. Please declare that
by exporting: export SESSIONA="myToken"`)
		os.Exit(1)
	}
}
