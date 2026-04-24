package pilot

// Task è l'atomo di coala: descrive un'azione da far compiere a OA
type Task struct {
	Name        string            `yaml:"name"`        // Es: "coa-initrd"
	Description string            `yaml:"description"` // Es: "Rigenerazione Initramfs"
	Command     string            `yaml:"command"`     // Il comando reale
	Chroot      bool              `yaml:"chroot"`      // Eseguire nel chroot? (solo per Installer)
	SetupFiles  map[string]string `yaml:"setup_files"` // File da scrivere prima del comando
}

// Profile è lo spartito completo di una distribuzione
type Profile struct {
	Remaster []Task `yaml:"remaster"`
	Install  []Task `yaml:"install"`
}
