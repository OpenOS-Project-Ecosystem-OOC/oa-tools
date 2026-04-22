package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// BlockDevice mappa un singolo disco restituito da lsblk.
type BlockDevice struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Model string `json:"model"`
	Type  string `json:"type"`
}

type LsblkOutput struct {
	Blockdevices []BlockDevice `json:"blockdevices"`
}

// GetHostDiskName risale l'albero dei mountpoint per trovare il disco fisico che ospita la root "/"
func GetHostDiskName() string {
	// 1. Troviamo il nodo montato su "/" (es. /dev/sda2, /dev/mapper/vg-root, overlay)
	cmd := exec.Command("findmnt", "-n", "-o", "SOURCE", "/")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	currentDev := strings.TrimSpace(string(out))

	// 2. Risaliamo l'albero dei blocchi finché non troviamo il disco "padre"
	for i := 0; i < 5; i++ { // Limite di sicurezza a 5 cicli per evitare loop infiniti
		cmd := exec.Command("lsblk", "-n", "-d", "-o", "PKNAME", currentDev)
		out, err := cmd.Output()
		if err != nil {
			break
		}

		parent := strings.TrimSpace(string(out))
		if parent == "" {
			break // Non ha più genitori, abbiamo raggiunto il disco fisico di base
		}

		// Il prossimo giro verificherà il genitore
		currentDev = "/dev/" + parent
	}

	// Restituiamo il nome pulito (es. "sda" o "nvme0n1")
	return strings.TrimPrefix(currentDev, "/dev/")
}

// GetAvailableDisks interroga il sistema usando lsblk, filtrando il disco di boot host.
func GetAvailableDisks() ([]BlockDevice, error) {
	cmd := exec.Command("lsblk", "-J", "-b", "-d", "-o", "NAME,SIZE,MODEL,TYPE")

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("errore durante l'esecuzione di lsblk: %v", err)
	}

	var lsblkData LsblkOutput
	if err := json.Unmarshal(out.Bytes(), &lsblkData); err != nil {
		return nil, fmt.Errorf("errore nel parsing del JSON di lsblk: %v", err)
	}

	bootDisk := GetHostDiskName()

	var realDisks []BlockDevice
	for _, dev := range lsblkData.Blockdevices {
		// Seleziona solo i dischi veri e propri ed ESCLUDE rigorosamente il disco di boot
		if dev.Type == "disk" && dev.Name != bootDisk {
			realDisks = append(realDisks, dev)
		}
	}

	return realDisks, nil
}

// FormatSize rende leggibili le dimensioni in byte (GB, TB).
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
