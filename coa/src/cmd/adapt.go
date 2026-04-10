package cmd

import (
	"coa/src/internal/engine"

	"github.com/spf13/cobra"
)

var adaptCmd = &cobra.Command{
	Use:   "adapt",
	Short: "Adapt monitor resolution for VMs",
	Run: func(cmd *cobra.Command, args []string) {
		CheckSudoRequirements(cmd.Name(), false)
		engine.HandleAdapt()
	},
}

func init() {
	rootCmd.AddCommand(adaptCmd)
}
