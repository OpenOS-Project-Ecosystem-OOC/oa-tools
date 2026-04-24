package engine

// Sostituisci "coala" con il nome reale del tuo modulo nel go.mod

// OaPlan è la struttura radice che il motore C si aspetta di leggere
type OaPlan struct {
	PathLiveFs      string     `json:"pathLiveFs"`
	Mode            string     `json:"mode"`
	Family          string     `json:"family"`
	InitrdCmd       string     `json:"initrd_cmd"`
	BootloadersPath string     `json:"bootloaders_path"`
	Plan            []OaAction `json:"plan"`
}

// OaAction è la singola azione che oa dovrà eseguire
type OaAction struct {
	Command    string `json:"command"`               // "oa_shell", "oa_users", o "oa_umount"
	Info       string `json:"info"`                  // La descrizione da mostrare a video
	RunCommand string `json:"run_command,omitempty"` // Il comando bash vero e proprio (solo per oa_shell)
	Chroot     bool   `json:"chroot,omitempty"`      // Se true, esegue in chroot
	Users      []any  `json:"users,omitempty"`       // Lista utenti (solo per oa_users)
}

const planPath = "/tmp/coa/finalize-plan.json" // Manteniamo /tmp/coa per non rompere Calamares
