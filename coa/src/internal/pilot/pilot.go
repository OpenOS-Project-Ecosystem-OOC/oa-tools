package pilot

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// src/internal/pilot/pilot.go
func findBrainDir() string {
	// Otteniamo la directory corrente (CWD)
	cwd, _ := os.Getwd()
	// Otteniamo la directory dove risiede l'eseguibile
	exePath, _ := os.Executable()
	baseDir := filepath.Dir(exePath)

	// Lista dei sospettati (dove potrebbe nascondersi il cervello)
	paths := []string{
		// 1. Percorso specifico indicato da te (relativo alla root del progetto)
		filepath.Join(cwd, "coa/conf/brain.d"),
		filepath.Join(baseDir, "conf/brain.d"),

		// 2. Percorsi standard (per sicurezza)
		filepath.Join(cwd, "brain.d"),
		filepath.Join(baseDir, "brain.d"),
		filepath.Join(baseDir, "..", "brain.d"),

		// 3. Percorsi di sistema
		"/usr/share/coa/brain.d",
		"/etc/coa/brain.d",
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			// Se troviamo la cartella, restituiamo il percorso assoluto
			absPath, _ := filepath.Abs(p)
			return absPath
		}
	}

	return ""
}

// getFamilyPath risolve la mappatura da familyID (es. "debian") alla cartella (es. "debian.d").
func getFamilyPath(familyID string) (string, error) {
	basePath := findBrainDir()
	if basePath == "" {
		return "", fmt.Errorf("brain directory not found")
	}

	mappingData, err := os.ReadFile(filepath.Join(basePath, "distro.yaml"))
	if err != nil {
		return "", err
	}

	var mapping struct {
		Families map[string]string `yaml:"families"`
	}
	if err := yaml.Unmarshal(mappingData, &mapping); err != nil {
		return "", err
	}

	folderName := mapping.Families[familyID]
	if folderName == "" {
		folderName = familyID + ".d"
	}

	return filepath.Join(basePath, folderName), nil
}

// readAreaConfig legge un file specifico (es. identity.yaml) e lo parsa nella struct fornita.
func readAreaConfig(familyID string, area string, out interface{}) error {
	familyPath, err := getFamilyPath(familyID)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(familyPath, area+".yaml")
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, out)
}
