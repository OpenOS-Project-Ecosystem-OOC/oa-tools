// src/internal/pilot/boot.go

package pilot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateBootConfig(familyID string, profile *BrainProfile) error {
	tmpDir := "/tmp/coa"
	targetFile := filepath.Join(tmpDir, "grub.cfg.final")

	// Cerchiamo i parametri nel profilo
	bootParams := "boot=live quiet splash" // Default di sicurezza
	for _, t := range profile.Tasks {
		if t.Name == "boot" && len(t.Commands) > 0 {
			bootParams = t.Commands[0]
			break
		}
	}

	// Costruiamo il contenuto reale del GRUB
	// Qui puoi sbizzarrirti con la logica specifica per Arch, Debian, etc.
	grubContent := fmt.Sprintf(`
set timeout=5
set default=0

menuentry "coa Live (%s)" {
    linux /live/vmlinuz %s
    initrd /live/initrd.img
}

menuentry "coa Live (RAM mode)" {
    linux /live/vmlinuz %s toram
    initrd /live/initrd.img
}
`, familyID, bootParams, bootParams)

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(targetFile, []byte(strings.TrimSpace(grubContent)), 0644)
}
