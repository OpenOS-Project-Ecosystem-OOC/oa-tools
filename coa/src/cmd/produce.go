package cmd

import (
	"coa/src/internal/distro"
	"coa/src/internal/engine"

	"github.com/spf13/cobra"
)

var (
	produceMode string
	producePath string
)

var produceCmd = &cobra.Command{
	Use:   "produce",
	Short: "Start a system remastering production flight",
	Run: func(cmd *cobra.Command, args []string) {
		// Richiede i privilegi di root
		CheckSudoRequirements(cmd.Name(), true)

		// Rileva la distribuzione host
		myDistro := distro.NewDistro()

		// Delega l'esecuzione al motore interno
		engine.HandleProduce(produceMode, producePath, myDistro)
	},
}

func init() {
	produceCmd.Flags().StringVar(&produceMode, "mode", "standard", "standard, clone, or crypted")
	produceCmd.Flags().StringVar(&producePath, "path", "/home/eggs", "working directory")

	rootCmd.AddCommand(produceCmd)
}
