package pilot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RunBrainLint scansiona brain.d e aggiunge gli header mancanti
func RunBrainLint() {
	basePath := findBrainDir()
	if basePath == "" {
		fmt.Println("[ERRORE] Cartella brain.d non trovata nei percorsi standard.")
		return
	}

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Saltiamo le directory e i file che non sono YAML (o il file delle mappature)
		if info.IsDir() || filepath.Ext(path) != ".yaml" || filepath.Base(path) == "distro.yaml" {
			return nil
		}

		// 1. Estrazione metadati dal percorso
		// relPath sarà qualcosa come "archlinux.d/initrd.yaml"
		relPath, _ := filepath.Rel(basePath, path)
		dirName := filepath.Dir(relPath)
		fileName := filepath.Base(relPath)

		family := strings.TrimSuffix(dirName, ".d")
		area := strings.TrimSuffix(fileName, ".yaml")

		fmt.Printf("[LINT] Controllo unità: %s/%s\n", family, area)

		// 2. Lettura e verifica dell'header
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("impossibile leggere %s: %v", path, err)
		}

		strContent := string(content)

		// Cerchiamo il marcatore unico di coa
		if !strings.Contains(strContent, "[ coa brain unit ]") {
			// Costruzione del "cappello"
			header := "# [ coa brain unit ]\n"
			header += fmt.Sprintf("# family: %s\n", family)
			header += fmt.Sprintf("# area:   %s\n", area)
			header += "# --------------------------------------------------\n\n"

			// Uniamo l'header al contenuto originale (pulendo spazi bianchi extra in cima)
			newContent := header + strings.TrimSpace(strContent) + "\n"

			// 3. Scrittura del file "vestito"
			err = os.WriteFile(path, []byte(newContent), 0644)
			if err != nil {
				fmt.Printf("  └─ [ERRORE] Scrittura fallita: %v\n", err)
			} else {
				fmt.Println("  └─ Cappello aggiunto con successo!")
			}
		} else {
			fmt.Println("  └─ Unità già configurata correttamente.")
		}

		return nil
	})

	if err != nil {
		fmt.Printf("[ERRORE] Durante il lint del cervello: %v\n", err)
	} else {
		fmt.Println("\n[SUCCESSO] Il Cervello è ora ordinato e riconoscibile.")
	}
}
