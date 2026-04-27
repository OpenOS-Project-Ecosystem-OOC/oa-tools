package tailor

import (
	"coa/pkg/distro"
	"coa/pkg/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// findCompatibleYaml cerca il file YAML più adatto basandosi sulla distro attuale.
func findCompatibleYaml(costumePath string) string {
	// Usiamo la sintassi corretta del tuo pacchetto distro
	d := distro.NewDistro()
	distroLike := d.DistroLike // Es: "Debian", "Arch", "Fedora"

	// Mappa di fallback basata sulla logica del sarto originale
	fallbacks := map[string][]string{
		"Ubuntu":   {"ubuntu.yaml", "debian.yaml", "devuan.yaml"},
		"Debian":   {"debian.yaml", "devuan.yaml", "ubuntu.yaml"},
		"Devuan":   {"devuan.yaml", "debian.yaml", "ubuntu.yaml"},
		"Arch":     {"arch.yaml", "debian.yaml"},
		"Fedora":   {"fedora.yaml", "debian.yaml"},
		"Alpine":   {"alpine.yaml", "debian.yaml"},
		"Opensuse": {"opensuse.yaml", "debian.yaml"},
	}

	filesToTry, exists := fallbacks[distroLike]
	if !exists {
		// Se la distro non è mappata, proviamo debian.yaml come fallback universale
		filesToTry = []string{"debian.yaml"}
	}

	for _, file := range filesToTry {
		fullPath := filepath.Join(costumePath, file)
		if _, err := os.Stat(fullPath); err == nil {
			// utils.LogCoala("Trovato cartamodello compatibile: %s", file)
			return fullPath
		}
	}

	return ""
}

// loadSuit trasforma il file YAML fisico nella struttura Suit
func loadSuit(yamlFile string) (*Suit, error) {
	if yamlFile == "" {
		return nil, fmt.Errorf("nessun file di definizione costume trovato")
	}

	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}

	var suit Suit
	if err := yaml.Unmarshal(data, &suit); err != nil {
		return nil, err
	}

	return &suit, nil
}

// getAvailablePackages interroga il sistema per sapere quali pacchetti esistono nei repo
func getAvailablePackages() map[string]struct{} {
	available := make(map[string]struct{})

	// Usiamo 'apt-cache pkgnames'
	out, err := utils.ExecCapture("apt-cache pkgnames")
	if err != nil {
		utils.LogError("Errore durante l'esecuzione di apt-cache: %v", err)
		return available
	}

	// Sostituiamo ogni possibile ritorno a capo (\r\n o \n) con uno spazio
	// e poi dividiamo per campi per essere sicuri di avere solo i nomi puliti
	fields := strings.Fields(out)
	for _, pkg := range fields {
		available[pkg] = struct{}{}
	}

	utils.LogCoala("Database pacchetti caricato: %d nomi trovati.", len(available))
	return available
}

// installWithRetries gestisce l'installazione con i tentativi (retry)
func installWithRetries(packages []string, retries int) {
	if len(packages) == 0 {
		return
	}

	pkgString := strings.Join(packages, " ")
	// Usiamo DEBIAN_FRONTEND per evitare prompt interattivi
	cmd := fmt.Sprintf("DEBIAN_FRONTEND=noninteractive apt-get install -yqq %s", pkgString)

	for i := 1; i <= retries; i++ {
		utils.LogCoala("Tentativo di installazione %d di %d...", i, retries)

		err := utils.Exec(cmd)
		if err == nil {
			utils.LogCoala("✅ Installazione riuscita!")
			break
		}

		if i < retries {
			utils.LogCoala("⚠️ Tentativo fallito, attesa 3 secondi...")
			time.Sleep(3 * time.Second)
		} else {
			utils.LogError("❌ Impossibile installare i pacchetti dopo %d tentativi.", retries)
		}
	}
}
