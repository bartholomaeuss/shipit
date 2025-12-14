/*
Copyright © 2025 Bartholomaeuss
*/
package cmd

import "github.com/spf13/cobra"

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Namespace for deployment-related commands",
	Long:  `Group deployment subcommands such as scripts to build, push, or release artifacts.`,
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
