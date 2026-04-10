package cmd

import (
	"coa/src/internal/builder"
	"coa/src/internal/distro"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Compile binaries and generate native distribution packages (.deb, PKGBUILD)",
	Run: func(cmd *cobra.Command, args []string) {
		// Non richiede i privilegi di root
		CheckSudoRequirements(cmd.Name(), false)

		// Rileva la distribuzione host (i Sensi)
		myDistro := distro.NewDistro()

		// Passa la palla al motore di build, includendo la versione di Git
		builder.HandleBuild(myDistro, AppVersion)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
