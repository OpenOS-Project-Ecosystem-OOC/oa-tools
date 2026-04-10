package cmd

import (
	"coa/src/internal/krill"

	"github.com/spf13/cobra"
)

var krillCmd = &cobra.Command{
	Use:   "krill",
	Short: "Start the system installation (The Hatching)",
	Run: func(cmd *cobra.Command, args []string) {
		// Krill formatterà dischi, serve root assoluto
		CheckSudoRequirements(cmd.Name(), true)

		// Avvia l'interfaccia utente (TUI) e l'installazione
		krill.HandleKrill()
	},
}

func init() {
	rootCmd.AddCommand(krillCmd)
}
