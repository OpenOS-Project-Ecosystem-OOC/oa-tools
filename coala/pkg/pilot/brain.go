package pilot

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const brainPath = "/etc/penguins-eggs/brain.d"

// IndexMap rappresenta la struttura del file index.yaml
type IndexMap struct {
	Distributions []struct {
		ID   string   `yaml:"id"`
		Like []string `yaml:"like"`
		File string   `yaml:"file"`
	} `yaml:"distributions"`
}

// DetectAndLoad è la funzione magica: capisce chi sei e ti dà il tuo profilo
func DetectAndLoad() (*Profile, error) {
	// 1. Identifichiamo la distribuzione corrente da /etc/os-release
	id, likes, err := identifySystem()
	if err != nil {
		return nil, fmt.Errorf("impossibile identificare il sistema: %w", err)
	}

	// 2. Leggiamo l'indice del Brain
	indexData, err := os.ReadFile(fmt.Sprintf("%s/index.yaml", brainPath))
	if err != nil {
		return nil, fmt.Errorf("indice del brain non trovato: %w", err)
	}

	var index IndexMap
	if err := yaml.Unmarshal(indexData, &index); err != nil {
		return nil, fmt.Errorf("errore sintassi index.yaml: %w", err)
	}

	// 3. Cerchiamo la rima giusta nell'indice
	targetFile := ""
	for _, dist := range index.Distributions {
		// Controllo ID diretto (es: "debian")
		if dist.ID == id {
			targetFile = dist.File
			break
		}
		// Controllo "Like" (es: se siamo su "ubuntu" e l'indice dice che ubuntu è like debian)
		for _, l := range likes {
			if dist.ID == l {
				targetFile = dist.File
				break
			}
		}
		if targetFile != "" {
			break
		}
	}

	if targetFile == "" {
		return nil, fmt.Errorf("nessun profilo trovato nel brain per ID=%s (Likes=%v)", id, likes)
	}

	// 4. Carichiamo lo spartito finale
	return loadProfile(targetFile)
}

// loadProfile legge il file specifico (es: "debian.yaml")
func loadProfile(filename string) (*Profile, error) {
	fullPath := fmt.Sprintf("%s/%s", brainPath, filename)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	var profile Profile
	if err := yaml.Unmarshal(data, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// identifySystem estrae ID e ID_LIKE da /etc/os-release
func identifySystem() (string, []string, error) {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	var id string
	var likes []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ID=") {
			id = strings.Trim(strings.Split(line, "=")[1], "\"")
		}
		if strings.HasPrefix(line, "ID_LIKE=") {
			likeStr := strings.Trim(strings.Split(line, "=")[1], "\"")
			likes = strings.Fields(likeStr)
		}
	}

	return id, likes, nil
}
