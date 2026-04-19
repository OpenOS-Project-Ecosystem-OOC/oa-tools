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
	os.MkdirAll(tmpDir, 0755)

	// 1. DEFAULT (Debian Style)
	bootParams := "boot=live components"
	kernelPath := "/live/vmlinuz"
	initrdPath := "/live/initrd.img"

	// 2. ADATTAMENTO PER ARCH
	if familyID == "archlinux" {
		bootParams = "archisobasedir=arch archisolabel=OA_LIVE rw"
	}

	// 3. SOVRASCRITTURA DA YAML (Se presente nel Brain)
	for _, t := range profile.Tasks {
		if t.Name == "boot" && len(t.Commands) > 0 {
			bootParams = t.Commands[0]
		}
	}

	// 4. COSTRUZIONE DEL TEMPLATE
	grubContent := fmt.Sprintf(`
set timeout=5
set default=0

# Ricerca della partizione tramite Label per evitare il rescue
search --no-floppy --set=root --label OA_LIVE

menuentry "coa Live (%s)" {
    linux %s %s
    initrd %s
}

menuentry "coa Live (%s) - RAM mode" {
    linux %s %s toram
    initrd %s
}
`, familyID, kernelPath, bootParams, initrdPath, familyID, kernelPath, bootParams, initrdPath)

	return os.WriteFile(targetFile, []byte(strings.TrimSpace(grubContent)), 0644)
}
