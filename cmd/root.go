/*
Copyright © 2025 Bartholomaeuss
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shipit",
	Short: "Clone, sync, clean, and deploy repos to your remote test hosts",
	Long: `Shipit streamlines the clone → tweak → ship loop for remote sandboxes.
Use the repo commands to create disposable clones and copy them to a host,
clean up both local and remote scratch directories with a single flag, and
run docker-based deployment scripts over SSH via the deploy namespace.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
