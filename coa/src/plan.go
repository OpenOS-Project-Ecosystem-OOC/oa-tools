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
	RunCommand      string   `json:"run_command,omitempty"`
	Args            []string `json:"args,omitempty"`	
}

// src/plan.go

type UserConfig struct {
	Login    string   `json:"login"`
	Password string   `json:"password"`
	Gecos    string   `json:"gecos"`
	Home     string   `json:"home"`
	Shell    string   `json:"shell"`
	Groups   []string `json:"groups"`
}

type FlightPlan struct {
	PathLiveFs      string       `json:"pathLiveFs"`
	Mode            string       `json:"mode"`
	InitrdCmd       string       `json:"initrd_cmd"`
	BootloadersPath string       `json:"bootloaders_path"`
	Users           []UserConfig `json:"users"` // Array globale degli utenti
	Plan            []Action     `json:"plan"`
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
		plan.InitrdCmd = "mkinitcpio -g {{out}} -k {{ver}}"
		plan.BootloadersPath = BootloaderRoot
	case "fedora", "opensuse":
		plan.InitrdCmd = "dracut --nomadas --force {{out}} {{ver}}"
		plan.BootloadersPath = BootloaderRoot
	default:
		plan.InitrdCmd = "mkinitramfs -o {{out}} {{ver}}"
		plan.BootloadersPath = "" 
	}

	// 2. Configurazione Utenti (Globale)
	if mode == "standard" {
		// Gestione dinamica dei gruppi admin
		adminGroup := "sudo"
		if d.FamilyID == "archlinux" || d.FamilyID == "fedora" {
			adminGroup = "wheel"
		}

		plan.Users = []UserConfig{
			{
				Login:    "live",
				Password: "$6$wM.wY0QtatvbQMHZ$QtIKXSpIsp2Sk57.Ny.JHk7hWDu.lxPtUYaTOiBnP4WBG5KS6JpUlpXj2kcSaaMje7fr01uiGmxZhE8kfZRqv.",
				Gecos:    "live,,,",
				Home:     "/home/live",
				Shell:    "/bin/bash",
				Groups:   []string{"cdrom", "audio", "video", "plugdev", "netdev", "autologin", adminGroup},
			},
		}
	} else {
		plan.Users = []UserConfig{}
	}

	// 3. Assemblaggio dinamico della catena di montaggio
	plan.Plan = []Action{
		{Command: "action_prepare"},
		{Command: "action_users"},
	}

	// --- Task di "Vestizione" (Patching configurazioni) ---
	// Se non siamo su Debian, iniettiamo i file estratti da coa nella liveroot
	if d.FamilyID != "debian" {
		configSrc := ""
		configDest := ""

		switch d.FamilyID {
		case "archlinux":
			configSrc = "/tmp/coa/configs/mkinitcpio/arch/."
			configDest = "/etc/mkinitcpio.d/"
		case "fedora":
			configSrc = "/tmp/coa/configs/dracut/."
			configDest = "/etc/dracut.conf.d/"
		}

		if configSrc != "" {
			plan.Plan = append(plan.Plan, Action{
				Command:    "action_run",
				RunCommand: "cp",
				Args:       []string{"-r", configSrc, configDest},
			})
		}
	}

	// Proseguiamo con il resto del piano standard
	plan.Plan = append(plan.Plan, 
		Action{Command: "action_initrd"},
		Action{Command: "action_livestruct"},
		Action{Command: "action_isolinux"},
		Action{Command: "action_uefi"},
		Action{Command: "action_squash"},
	)

	// Inserzione modulare per cifratura
	if mode == "crypted" {
		plan.Plan = append(plan.Plan, Action{
			Command:         "action_crypted",
			CryptedPassword: "evolution",
		})
	}

	// 4. Generazione ISO e chiusura
	isoName := fmt.Sprintf("egg-of_%s-%s-oa_amd64.iso", d.DistroID, d.CodenameID)
	
	plan.Plan = append(plan.Plan, Action{
		Command:   "action_iso",
		VolID:     "OA_LIVE",
		OutputISO: isoName,
	})
	
	plan.Plan = append(plan.Plan, Action{Command: "action_cleanup"})

	return plan
}
