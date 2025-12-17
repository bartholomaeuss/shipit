/*
Copyright © 2025 Bartholomaeuss
*/
package cmd

import (
	"shipit/core/repo"

	"github.com/spf13/cobra"
)

var (
	repoCloneURL  string
	repoCloneHost string
	repoCloneUser string
)

// repoCloneCmd implements `shipit repo clone`.
var repoCloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a GitHub repository into a temp directory",
	Long: `Create a temporary directory, clone the requested repository into it,
and print the cd command so you can jump into the clone immediately.

Example: shipit repo clone --url https://github.com/foo/bar`,
	RunE: func(cmd *cobra.Command, args []string) error {

		svc := repo.NewCloneService(repoCloneHost, repoCloneUser, cmd.OutOrStdout(), cmd.ErrOrStderr())

		if err := svc.ParseRepoUrl(repoCloneURL); err != nil {
			return err
		}

		return svc.Run()
	},
}

func init() {
	repoCmd.AddCommand(repoCloneCmd)
	repoCloneCmd.Flags().StringVar(&repoCloneURL, "url", "", "GitHub repository URL to clone")
	repoCloneCmd.Flags().StringVar(&repoCloneHost, "host", "test", "SSH host (from ~/.ssh/config) whose home directory receives the clone")
	repoCloneCmd.Flags().StringVar(&repoCloneUser, "user", "", "SSH username to use for the remote copy")
}
