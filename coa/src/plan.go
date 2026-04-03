package main

import (
	"fmt"
)

// Action rappresenta un singolo blocco "command" nell'array "plan"
type Action struct {
	Command         string `json:"command"`
	VolID           string `json:"volid,omitempty"`
	OutputISO       string `json:"output_iso,omitempty"`
	CryptedPassword string `json:"crypted_password,omitempty"`
}

// FlightPlan rappresenta l'intero piano da passare a oa
type FlightPlan struct {
	PathLiveFs      string   `json:"pathLiveFs"`
	Mode            string   `json:"mode"`
	InitrdCmd       string   `json:"initrd_cmd"`
	BootloadersPath string   `json:"bootloaders_path"` // Questo è il nome corretto 
	Plan            []Action `json:"plan"`
}

func GeneratePlan(d *Distro, mode string, workPath string) FlightPlan {
	plan := FlightPlan{
		PathLiveFs: workPath,
		Mode:       mode,
	}

	// 1. Astrazione Initramfs e Bootloaders (Il Terzo Pilastro)
	switch d.FamilyID {
	case "debian":
		plan.InitrdCmd = "mkinitramfs -o {{out}} {{ver}}"
		plan.BootloadersPath = "" // Su Debian usiamo quelli di sistema
	case "archlinux":
		// Su Arch usiamo mkinitcpio
		plan.InitrdCmd = "mkinitcpio -g {{out}} -k {{ver}}"
		// Usiamo la costante BootloaderRoot definita in utils.go 
		plan.BootloadersPath = BootloaderRoot
	case "fedora", "opensuse":
		plan.InitrdCmd = "dracut --nomadas --force {{out}} {{ver}}"
		plan.BootloadersPath = BootloaderRoot
	default:
		plan.InitrdCmd = "mkinitramfs -o {{out}} {{ver}}"
		plan.BootloadersPath = "" // Fallback vuoto
	}

	// 3. Assemblaggio dinamico della catena di montaggio [cite: 28]
	plan.Plan = []Action{
		{Command: "action_prepare"},
		{Command: "action_users"},
		{Command: "action_initrd"},
		{Command: "action_livestruct"},
		{Command: "action_isolinux"},
		{Command: "action_uefi"},
		{Command: "action_squash"},
	}

	// Inserzione modulare per cifratura [cite: 29]
	if mode == "crypted" {
		plan.Plan = append(plan.Plan, Action{
			Command:         "action_crypted",
			CryptedPassword: "evolution",
		})
	}

	// 4. Generazione ISO e chiusura [cite: 30]
	isoName := fmt.Sprintf("egg-of_%s-%s-oa_amd64.iso", d.DistroID, d.CodenameID)
	
	plan.Plan = append(plan.Plan, Action{
		Command:   "action_iso",
		VolID:     "OA_LIVE",
		OutputISO: isoName,
	})
	
	plan.Plan = append(plan.Plan, Action{Command: "action_cleanup"})

	return plan
}
