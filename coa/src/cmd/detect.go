package cmd

import (
	"coa/src/internal/distro"
	"coa/src/internal/engine"

	"github.com/spf13/cobra"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Show host distribution discovery info",
	Run: func(cmd *cobra.Command, args []string) {
		// Non richiede permessi di root
		CheckSudoRequirements(cmd.Name(), false)

		myDistro := distro.NewDistro()
		engine.HandleDetect(myDistro)
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)
}
