package pilot

// Task rappresenta un'unità di lavoro atomica.
type Task struct {
	Name        string
	Files       map[string]string
	Commands    []string
	Chroot      bool
	Description string
}

// BrainProfile è il contenitore dei compiti.
type BrainProfile struct {
	Tasks []Task
}

// --- AGGIUNGI QUESTE STRUTTURE ---

type IdentityConfig struct {
	AdminGroup string   `yaml:"admin_group"`
	UserGroups []string `yaml:"user_groups"`
}

type InitrdConfig struct {
	Command string            `yaml:"command"`
	Files   map[string]string `yaml:"files"`
}

type BootConfig struct {
	Params string `yaml:"params"`
}

// links sull'iso per arch
type LayoutConfig struct {
	Links map[string]string `yaml:"links"` // "destinazione" -> "sorgente"
}
