/*
Copyright © 2025 Bartholomaeuss
*/
package cmd

import "github.com/spf13/cobra"

// repoCmd provides the namespace for repository-related subcommands.
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Namespace for repository operations",
	Long: `Group repository-related subcommands like repo clone, repo list, etc.
Invoke one of those subcommands to perform an action.`,
}

func init() {
	rootCmd.AddCommand(repoCmd)
}
