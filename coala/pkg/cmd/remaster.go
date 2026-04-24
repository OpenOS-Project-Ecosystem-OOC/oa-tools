package cmd

import (
	"fmt"
	"os"
	"os/exec"

	// Sostituisci con i percorsi corretti del tuo nuovo branch coala
	"coala/pkg/distro" // Lo manteniamo per ottenere il FamilyID da passare al JSON
	"coala/pkg/engine"
	"coala/pkg/pilot"

	"github.com/spf13/cobra"
)

var (
	produceMode string
	producePath string
)

var remasterCmd = &cobra.Command{
	Use:   "remaster",
	Short: "Start a system remastering flight (ISO production)",
	Long: `The 'remaster' command orchestrates the creation of a bootable live ISO. 
It uses the new Coala architecture to read the agnostic Brain profile 
and generate a precise execution plan for the OA engine.`,
	Example: `  # Start a standard ISO remastering
  sudo coala remaster --mode standard`,
	Run: func(cmd *cobra.Command, args []string) {
		CheckSudoRequirements(cmd.Name(), true)

		fmt.Println("\033[1;36m[coala]\033[0m Avvio procedura di rimasterizzazione...")

		// 1. Identità: Chi siamo?
		myDistro := distro.NewDistro()

		// 2. PILOT: Carichiamo lo spartito dal Brain
		profile, err := pilot.DetectAndLoad()
		if err != nil {
			fmt.Printf("\033[1;31m[ERRORE CRITICO]\033[0m Impossibile caricare il Brain Profile: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\033[1;32m[coala]\033[0m Spartito caricato con successo.\n")

		// 3. ENGINE: Generiamo il piano JSON per oa usando la sezione 'Remaster'
		// Passiamo i Task, la famiglia, isRemaster=true, e il path di lavoro
		err = engine.GeneratePlan(profile.Remaster, myDistro.FamilyID, true, producePath)
		if err != nil {
			fmt.Printf("\033[1;31m[ERRORE CRITICO]\033[0m Impossibile generare il piano di volo: %v\n", err)
			os.Exit(1)
		}

		// 4. DECOLLO: Eseguiamo il motore C (oa) passandogli il JSON appena generato
		fmt.Println("\033[1;36m[coala]\033[0m Passaggio dei comandi al motore OA...")
		oaCmd := exec.Command("oa", "/tmp/coa/finalize-plan.json")

		// Colleghiamo l'output di oa direttamente al terminale dell'utente
		oaCmd.Stdout = os.Stdout
		oaCmd.Stderr = os.Stderr

		if err := oaCmd.Run(); err != nil {
			fmt.Printf("\n\033[1;31m[ERRORE]\033[0m L'esecuzione di oa è fallita: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n\033[1;32m[SUCCESSO]\033[0m Rimasterizzazione completata! L'uovo è pronto. 🐧🥚")
	},
}

func init() {
	remasterCmd.Flags().StringVar(&produceMode, "mode", "standard", "standard, clone, or crypted")
	remasterCmd.Flags().StringVar(&producePath, "path", "/home/eggs", "working directory")

	rootCmd.AddCommand(remasterCmd)
}
