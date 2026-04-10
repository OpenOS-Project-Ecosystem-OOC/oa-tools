package cmd

import (
	"coa/src/internal/engine"

	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Free the nest and unmount filesystems",
	Run: func(cmd *cobra.Command, args []string) {
		CheckSudoRequirements(cmd.Name(), true)
		engine.HandleKill()
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
}
