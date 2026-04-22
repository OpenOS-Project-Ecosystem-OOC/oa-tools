package cmd

import (
	"coa/src/internal/engine"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var sysinstallCmd = &cobra.Command{
	Use:   "sysinstall",
	Short: "Install the live system to a physical disk",
	Long: `The 'sysinstall' command is the interactive system installer.
It gathers user preferences via TUI and orchestrates the physical installation 
on the target disk.
WARNING: This operation is destructive.`,
	Example: `  # Launch the system installer
  sudo coa sysinstall`,
	Run: func(cmd *cobra.Command, args []string) {
		CheckSudoRequirements(cmd.Name(), true)

		// 1. Chiediamo all'utente quale disco usare
		targetDisk, err := engine.SelectTargetDisk()
		if err != nil {
			fmt.Printf("Errore critico: %v\n", err)
			os.Exit(1)
		}

		// 2. Ora che abbiamo il disco, generiamo il piano!
		err = engine.GenerateInstallPlan(targetDisk, "artisan", "password123")
		if err != nil {
			fmt.Printf("Errore generazione piano: %v\n", err)
			os.Exit(1)
		}

		// Qui in futuro chiameremo il Chirurgo C per eseguire il piano
	},
}

func init() {
	rootCmd.AddCommand(sysinstallCmd)
}
