package engine

import (
	"encoding/json"
	"fmt"
	"os"
)

// GeneratePlan prende lo spartito del Pilot e crea il file JSON per oa
func GeneratePlan(tasks []pilot.Task, family string, isRemaster bool, liveFsPath string) error {
	// 1. Inizializziamo il piano principale
	plan := OaPlan{
		PathLiveFs:      liveFsPath, // Es: "/home/eggs" o "/tmp/coa/calamares-root"
		Mode:            "standard",
		Family:          family,
		BootloadersPath: "/tmp/coa/bootloaders",
		Plan:            make([]OaAction, 0),
	}

	// 2. Traduciamo ogni Task di Coala in una Action per oa
	for _, t := range tasks {
		action := OaAction{
			Info: t.Description,
		}

		// Riconosciamo i comandi nativi di oa, altrimenti usiamo oa_shell
		switch t.Command {
		case "oa_users":
			action.Command = "oa_users"
			// Qui potremmo iniettare la logica per generare la struct degli utenti,
			// che magari possiamo passare al GeneratePlan in futuro.
			// Per ora mettiamo un array vuoto o la logica predefinita.
			action.Users = []any{}

		case "oa_umount":
			action.Command = "oa_umount"

		default:
			// Tutto il resto è uno script shell!
			action.Command = "oa_shell"
			action.RunCommand = t.Command
			action.Chroot = t.Chroot
		}

		plan.Plan = append(plan.Plan, action)
	}

	// 3. Creiamo la directory se non esiste
	if err := os.MkdirAll("/tmp/coa", 0755); err != nil {
		return fmt.Errorf("errore creazione cartella /tmp/coa: %w", err)
	}

	// 4. Scriviamo il file JSON formattato (Indent per debug facile)
	jsonData, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return fmt.Errorf("errore marshalling JSON: %w", err)
	}

	if err := os.WriteFile(planPath, jsonData, 0644); err != nil {
		return fmt.Errorf("errore scrittura %s: %w", planPath, err)
	}

	fmt.Printf("\033[1;32m[coala-engine]\033[0m Piano di volo generato con successo in %s\n", planPath)
	return nil
}
