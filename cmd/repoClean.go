/*
Copyright © 2025 Bartholomaeuss
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	repoCleanAll  bool
	repoCleanDir  string
	repoCleanHost string
	repoCleanUser string
)

// repoCleanCmd implements `shipit repo clean`.
var repoCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove temporary shipit clone directories",
	Long: `Delete temporary directories created by ` + "`shipit repo clone`" + `.
Use --all to remove every directory matching shipit-repo-* under the OS temp directory,
or --specific-dir to delete a single directory. Provide --user (and optionally --host)
to delete the corresponding directories on your test host via SSH.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !repoCleanAll && repoCleanDir == "" {
			return errors.New("specify either --all or --specific-dir")
		}

		if cmd.Flags().Changed("host") && repoCleanHost == "" {
			return fmt.Errorf("--host cannot be empty")
		}

		if repoCleanUser == "" {
			return fmt.Errorf("--user cannot be empty")
		}

		if repoCleanAll {
			if err := deleteAllTempRepos(cmd); err != nil {
				return err
			}
		}

		if repoCleanDir != "" {
			if err := deleteSpecificDir(cmd, repoCleanDir); err != nil {
				return err
			}
		}

		if repoCleanUser != "" {
			if err := deleteRemoteRepos(cmd); err != nil {
				return err
			}
		}

		return nil
	},
}

func deleteAllTempRepos(cmd *cobra.Command) error {
	pattern := filepath.Join(os.TempDir(), "shipit-repo-*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("glob %s: %w", pattern, err)
	}

	if len(matches) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "No directories matching %s\n", pattern)
		return nil
	}

	for _, dir := range matches {
		if err := deleteSpecificDir(cmd, dir); err != nil {
			return err
		}
	}

	return nil
}

func deleteSpecificDir(cmd *cobra.Command, dir string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("resolve %s: %w", dir, err)
	}

	base := filepath.Base(abs)
	if !strings.HasPrefix(base, "shipit-repo-") {
		return fmt.Errorf("refusing to delete %s: not a shipit temp directory", abs)
	}

	info, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(cmd.OutOrStdout(), "Directory %s does not exist, skipping\n", abs)
			return nil
		}
		return fmt.Errorf("stat %s: %w", abs, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", abs)
	}

	if err := os.RemoveAll(abs); err != nil {
		return fmt.Errorf("remove %s: %w", abs, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Removed %s\n", abs)
	return nil
}

func deleteRemoteRepos(cmd *cobra.Command) error {
	target := repoCleanHost
	if target == "" {
		target = "test"
	}

	sshTarget := target
	if repoCleanUser != "" {
		sshTarget = fmt.Sprintf("%s@%s", repoCleanUser, target)
	}

	if repoCleanAll {
		if err := runSSH(cmd, sshTarget, "rm -rf ~/shipit-repo-*"); err != nil {
			return err
		}
	}

	if repoCleanDir != "" {
		base := filepath.Base(repoCleanDir)
		if !strings.HasPrefix(base, "shipit-repo-") {
			return fmt.Errorf("refusing remote delete for %s: not a shipit directory", repoCleanDir)
		}
		cmdStr := fmt.Sprintf("rm -rf ~/%s", base)
		if err := runSSH(cmd, sshTarget, cmdStr); err != nil {
			return err
		}
	}

	return nil
}

func runSSH(cmd *cobra.Command, target, remoteCmd string) error {
	sshCmd := exec.Command("ssh", target, remoteCmd)
	sshCmd.Stdout = cmd.OutOrStdout()
	sshCmd.Stderr = cmd.ErrOrStderr()

	if err := sshCmd.Run(); err != nil {
		return fmt.Errorf("ssh -v %s %q failed: %w", target, remoteCmd, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Remote cleanup via %s: %s\n", target, remoteCmd)
	return nil
}

func init() {
	repoCmd.AddCommand(repoCleanCmd)
	repoCleanCmd.Flags().BoolVar(&repoCleanAll, "all", false, "Delete all shipit-repo-* directories under the OS temp directory")
	repoCleanCmd.Flags().StringVar(&repoCleanDir, "specific-dir", "", "Delete a specific directory created by shipit repo clone")
	repoCleanCmd.Flags().StringVar(&repoCleanHost, "host", "test", "SSH host whose home directory should be cleaned up")
	repoCleanCmd.Flags().StringVar(&repoCleanUser, "user", "", "SSH username to use for remote cleanup")
}
