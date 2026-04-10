package cmd

import (
	"coa/src/internal/engine"

	"github.com/spf13/cobra"
)

var docsCmd = &cobra.Command{
	Use:    "docs",
	Short:  "Generate man pages, markdown wiki, and completion scripts",
	Hidden: true, // Lo manteniamo nascosto agli utenti finali
	Run: func(cmd *cobra.Command, args []string) {
		CheckSudoRequirements(cmd.Name(), false)

		// Passiamo rootCmd all'engine in modo che Cobra possa analizzare
		// l'albero di tutti i comandi registrati e generare il manuale.
		engine.HandleDocs(rootCmd)
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}
