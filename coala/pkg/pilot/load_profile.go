package pilot

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadProfile legge direttamente il file appiattito (es: "debian.yaml")
func LoadProfile(filename string) (*Profile, error) {
	path := fmt.Sprintf("/etc/penguins-eggs/brain.d/%s", filename)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("impossibile leggere lo spartito %s: %w", path, err)
	}

	var profile Profile
	err = yaml.Unmarshal(data, &profile)
	if err != nil {
		return nil, fmt.Errorf("errore di sintassi nello YAML: %w", err)
	}

	return &profile, nil
}
