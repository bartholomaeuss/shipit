/*
Copyright © 2025 Bartholomaeuss
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const defaultDeployRunHost = "test"

var (
	deployRunDir  string
	deployRunHost string
	deployRunUser string
)

var deployRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run repository docker deployment scripts on the remote host",
	Long: `Connect to your test host over SSH and execute the repository's docker scripts
located under scripts/docker: prerun.sh, run.sh, and postrun.sh (in that order).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(deployRunDir) == "" {
			return fmt.Errorf("--dir is required")
		}

		if cmd.Flags().Changed("host") && strings.TrimSpace(deployRunHost) == "" {
			return fmt.Errorf("--host cannot be empty")
		}

		host := deployRunHost
		if host == "" {
			host = defaultDeployRunHost
		}

		scripts := []string{
			"~/shipit-repo-987431033/scripts/docker/prerun.sh",
			"~/shipit-repo-987431033/scripts/docker/run.sh",
			"~/shipit-repo-987431033/scripts/docker/postrun.sh",
		}

		sshTarget := host
		if deployRunUser != "" {
			sshTarget = fmt.Sprintf("%s@%s", deployRunUser, host)
		}

		for _, script := range scripts {
			remoteCmd := fmt.Sprintf(script)
			if err := runSSH(cmd, sshTarget, remoteCmd); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	deployCmd.AddCommand(deployRunCmd)
	deployRunCmd.Flags().StringVar(&deployRunDir, "dir", "", "Remote path to the repository containing scripts/docker (required)")
	deployRunCmd.Flags().StringVar(&deployRunHost, "host", defaultDeployRunHost, "SSH host to execute the deployment scripts on")
	deployRunCmd.Flags().StringVar(&deployRunUser, "user", "", "SSH username to use for remote execution")
	_ = deployRunCmd.MarkFlagRequired("dir")
}
