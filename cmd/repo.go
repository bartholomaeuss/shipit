/*
Copyright © 2025 Bartholomaeuss
*/
package cmd

import "github.com/spf13/cobra"

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Namespace for repository operations",
	Long: `Group repository-related subcommands like repo clone, repo list, etc.
Invoke one of those subcommands to perform an action.`,
}

func init() {
	rootCmd.AddCommand(repoCmd)
}
