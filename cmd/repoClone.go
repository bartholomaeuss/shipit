/*
Copyright © 2025 Bartholomaeuss
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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
		url := repoCloneURL

		if cmd.Flags().Changed("host") && repoCloneHost == "" {
			return fmt.Errorf("--host cannot be empty; omit the flag to use the default value")
		}

		if repoCloneUser == "" {
			return fmt.Errorf("--user cannot be empty")
		}

		tempDir, err := os.MkdirTemp("", "shipit-repo-")
		if err != nil {
			return fmt.Errorf("failed to create temp directory: %w", err)
		}

		// git clone expects the destination path not to exist
		if err := os.Remove(tempDir); err != nil {
			return fmt.Errorf("failed to prepare temp directory: %w", err)
		}

		gitCmd := exec.Command("git", "clone", url, tempDir)
		gitCmd.Stdout = cmd.OutOrStdout()
		gitCmd.Stderr = cmd.ErrOrStderr()

		if err := gitCmd.Run(); err != nil {
			return fmt.Errorf("git clone failed: %w", err)
		}

		absPath, err := filepath.Abs(tempDir)
		if err != nil {
			absPath = tempDir
		}

		remoteDir := fmt.Sprintf("~/%s", filepath.Base(absPath))
		targetHost := repoCloneHost
		targetHost = fmt.Sprintf("%s@%s", repoCloneUser, repoCloneHost)

		scpTarget := fmt.Sprintf("%s:~", targetHost)
		scpCmd := exec.Command("scp", "-r", absPath, scpTarget)
		scpCmd.Stdout = cmd.OutOrStdout()
		scpCmd.Stderr = cmd.ErrOrStderr()

		if err := scpCmd.Run(); err != nil {
			return fmt.Errorf("failed to copy repository to %s: %w", scpTarget, err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "\nRepository copied to %s:%s\n", targetHost, remoteDir)

		fmt.Fprintf(cmd.OutOrStdout(), "\nRun:\n  cd %s\n", absPath)

		return nil
	},
}

func init() {
	repoCmd.AddCommand(repoCloneCmd)
	repoCloneCmd.Flags().StringVar(&repoCloneURL, "url", "", "GitHub repository URL to clone")
	repoCloneCmd.Flags().StringVar(&repoCloneHost, "host", "test", "SSH host (from ~/.ssh/config) whose home directory receives the clone")
	repoCloneCmd.Flags().StringVar(&repoCloneUser, "user", "", "SSH username to use for the remote copy")
}
