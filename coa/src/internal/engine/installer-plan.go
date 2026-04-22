package engine

import (
	"encoding/json"
	"fmt"
	"os"
)

// IsUEFI controlla se il sistema è avviato in modalità UEFI
func IsUEFI() bool {
	_, err := os.Stat("/sys/firmware/efi")
	return !os.IsNotExist(err)
}

// GenerateInstallPlan crea il JSON completo per il Chirurgo C
func GenerateInstallPlan(targetDisk *BlockDevice, newUsername, newPassword string) error {
	fmt.Println("\033[1;36m[oa-installer]\033[0m Generazione piano di volo in corso...")

	// 1. Variabili di base
	devPath := fmt.Sprintf("/dev/%s", targetDisk.Name)
	isUEFI := IsUEFI()

	// Dichiarazione delle variabili che useremo per le azioni fisiche
	var partitions []map[string]string
	var formatActions []map[string]string
	var mountCmds string
	var grubCmd string
	var diskLabel string

	// 2. Logica di partizionamento dinamico
	if isUEFI {
		fmt.Println(" -> Rilevato sistema UEFI. Configurazione tabella GPT + ESP.")
		diskLabel = "gpt"
		partitions = []map[string]string{
			{"name": "ESP", "size": "512M", "type": "EF00"},     // EFI System Partition
			{"name": "OA_ROOT", "size": "100%", "type": "8300"}, // Linux Root
		}
		formatActions = []map[string]string{
			{"device": devPath + "1", "fs": "vfat", "label": "EFI"},
			{"device": devPath + "2", "fs": "ext4", "label": "ROOT"},
		}
		mountCmds = fmt.Sprintf("mkdir -p /mnt/target && mount %s2 /mnt/target && mkdir -p /mnt/target/boot/efi && mount %s1 /mnt/target/boot/efi", devPath, devPath)
		grubCmd = "grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=OA && grub-mkconfig -o /boot/grub/grub.cfg"
	} else {
		fmt.Println(" -> Rilevato sistema BIOS/Legacy. Configurazione tabella MBR.")
		diskLabel = "dos"
		partitions = []map[string]string{
			{"name": "OA_ROOT", "size": "100%", "type": "8300"},
		}
		formatActions = []map[string]string{
			{"device": devPath + "1", "fs": "ext4", "label": "ROOT"},
		}
		mountCmds = fmt.Sprintf("mkdir -p /mnt/target && mount %s1 /mnt/target", devPath)
		grubCmd = fmt.Sprintf("grub-install %s && grub-mkconfig -o /boot/grub/grub.cfg", devPath)
	}

	// 3. Costruzione dell'Azione per il Nuovo Utente
	liveUser := UserDef{
		Login:    newUsername,
		Password: newPassword, // In produzione questo dovrà essere un hash!
		Home:     "/home/" + newUsername,
		Shell:    "/bin/bash",
		Gecos:    newUsername + ",,,",
		Uid:      1000,
		Gid:      1000,
		Groups:   []string{"wheel", "audio", "video"},
	}

	// 4. Stesura del Piano (FlightPlan)
	plan := FlightPlan{
		Mode:       "install",
		PathLiveFs: "/mnt/target", // Il Chirurgo lavorerà qui dentro!
		Plan: []Action{
			// Fase 1: Partizionamento (Nativo C)
			{
				Command:    "oa_partition",
				Info:       fmt.Sprintf("Partizionamento disco %s", devPath),
				Device:     devPath,
				Label:      diskLabel,
				Partitions: partitions, // Variabile iniettata!
			},
			// Fase 2: Formattazione (Nativo C)
			{
				Command: "oa_format",
				Info:    "Formattazione partizioni",
				Actions: formatActions, // Variabile iniettata!
			},
			// Fase 3: Montaggio (Bash)
			{
				Command:    "oa_shell",
				Info:       "Montaggio partizioni di destinazione",
				RunCommand: mountCmds,
				Chroot:     false,
			},
			// Fase 4: Travaso Dati (Unsquashfs)
			// Per test usiamo un rsync dalla root attuale. In prod sarà l'estrazione di squashfs.
			{
				Command:    "oa_shell",
				Info:       "Copia del sistema sul nuovo disco...",
				RunCommand: "rsync -aAXv --exclude={'/dev/*','/proc/*','/sys/*','/tmp/*','/run/*','/mnt/*','/media/*','/lost+found'} / /mnt/target/",
				Chroot:     false,
			},
			// Fase 5: Identità (Nativo C)
			{
				Command: "oa_users",
				Info:    "Configurazione identità utente finale",
				Mode:    "standard",
				Users:   []UserDef{liveUser},
				Chroot:  false, // oa_users aggiunge automaticamente PathLiveFs per operare nel mount
			},
			// Fase 6: Bootloader (Bash in Chroot)
			{
				Command:    "oa_shell",
				Info:       "Installazione bootloader (GRUB)",
				RunCommand: grubCmd,
				Chroot:     true,
			},
		},
	}

	planJSON, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return fmt.Errorf("errore serializzazione JSON: %v", err)
	}

	// 5. Scrittura su disco
	os.MkdirAll("/tmp/coa", 0755)
	planPath := "/tmp/coa/install-plan.json"
	err = os.WriteFile(planPath, planJSON, 0644)
	if err != nil {
		return fmt.Errorf("errore scrittura file %s: %v", planPath, err)
	}

	fmt.Printf("\033[1;32m[OK]\033[0m Piano di installazione generato con successo in: %s\n", planPath)
	return nil
}
