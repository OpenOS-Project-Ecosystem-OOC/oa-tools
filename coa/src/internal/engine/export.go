package engine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
)

const (
	remoteUserHost = "root@192.168.1.2"
	remoteIsoPath  = "/var/lib/vz/template/iso/"
	remotePkgPath  = "/eggs/"
	isoSrcDir      = "/home/eggs"
)

// HandleExportIso copia in remoto la ISO utilizzando SSH/SCP con Multiplexing
func HandleExportIso(clean bool) {
	allFiles, _ := filepath.Glob(filepath.Join(isoSrcDir, "egg-of_*.iso"))
	if len(allFiles) == 0 {
		fmt.Println("\033[1;31m[coa]\033[0m Nest is empty. No ISOs found in", isoSrcDir)
		return
	}

	latestFiles := make(map[string]string)
	re := regexp.MustCompile(`_\d{4}-\d{2}-\d{2}_\d{4}\.iso$`)

	for _, path := range allFiles {
		fileName := filepath.Base(path)
		prefix := re.ReplaceAllString(fileName, "")

		if info, err := os.Stat(path); err == nil {
			if current, exists := latestFiles[prefix]; exists {
				cInfo, _ := os.Stat(current)
				if info.ModTime().After(cInfo.ModTime()) {
					latestFiles[prefix] = path
				}
			} else {
				latestFiles[prefix] = path
			}
		}
	}

	// SSH Multiplexing setup
	socketPath := "/tmp/coa-ssh-mux"
	muxArgs := []string{
		"-o", "ControlMaster=auto",
		"-o", "ControlPath=" + socketPath,
		"-o", "ControlPersist=2m",
	}

	defer func() {
		exec.Command("ssh", "-O", "exit", "-o", "ControlPath="+socketPath, remoteUserHost).Run()
		os.Remove(socketPath)
	}()

	for prefix, localPath := range latestFiles {
		targetFileName := filepath.Base(localPath)
		fmt.Printf("\n\033[1;35m[PROCESS]\033[0m Family: %s\n", prefix)

		if clean {
			fmt.Printf("\033[1;34m[CLEAN]\033[0m Removing old versions on Proxmox...\n")
			rmCmdStr := fmt.Sprintf("rm -f %s%s*", remoteIsoPath, prefix)
			sshArgs := append(muxArgs, remoteUserHost, rmCmdStr)
			sshCmd := exec.Command("ssh", sshArgs...)
			sshCmd.Stdout, sshCmd.Stderr, sshCmd.Stdin = os.Stdout, os.Stderr, os.Stdin
			if err := sshCmd.Run(); err != nil {
				fmt.Printf("\033[1;33m[WARNING]\033[0m Remote cleanup failed or no old files found.\n")
			} else {
				fmt.Printf("\033[1;32m[CLEAN]\033[0m Old versions removed.\n")
			}
		}

		fmt.Printf("\033[1;32m[COPY]\033[0m Sending %s to Proxmox...\n", targetFileName)
		dstStr := fmt.Sprintf("%s:%s", remoteUserHost, remoteIsoPath)
		scpArgs := append(muxArgs, localPath, dstStr)
		scpCmd := exec.Command("scp", scpArgs...)
		scpCmd.Stdout, scpCmd.Stderr, scpCmd.Stdin = os.Stdout, os.Stderr, os.Stdin

		if err := scpCmd.Run(); err != nil {
			fmt.Printf("\033[1;31m[ERROR]\033[0m Copy failed: %v\n", err)
		} else {
			fmt.Printf("\033[1;32m[SUCCESS]\033[0m %s is now on Proxmox.\n", targetFileName)
		}
	}
}

// HandleExportPkg esporta i pacchetti nativi (DEB o Arch)
// HandleExportPkg esporta i pacchetti nativi (DEB o Arch)
func HandleExportPkg(clean bool) {
	fmt.Println("\033[1;34m[PROCESS]\033[0m Searching for native packages...")

	debFiles, _ := filepath.Glob("*.deb")
	archFiles, _ := filepath.Glob("*.pkg.tar.zst")
	var allPackages []string
	allPackages = append(append(allPackages, debFiles...), archFiles...)

	if len(allPackages) == 0 {
		fmt.Println("\033[1;31m[ERROR]\033[0m No native packages found.")
		return
	}

	sort.Slice(allPackages, func(i, j int) bool {
		infoI, _ := os.Stat(allPackages[i])
		infoJ, _ := os.Stat(allPackages[j])
		return infoI.ModTime().Before(infoJ.ModTime())
	})
	latestPkg := allPackages[len(allPackages)-1]

	// --- INIZIO SETUP SSH MULTIPLEXING ---
	socketPath := "/tmp/coa-ssh-mux-pkg"
	muxArgs := []string{
		"-o", "ControlMaster=auto",
		"-o", "ControlPath=" + socketPath,
		"-o", "ControlPersist=2m",
	}

	// Chiude la connessione master alla fine della funzione
	defer func() {
		exec.Command("ssh", "-O", "exit", "-o", "ControlPath="+socketPath, remoteUserHost).Run()
		os.Remove(socketPath)
	}()
	// --- FINE SETUP ---

	if clean {
		fmt.Printf("\033[1;33m[CLEAN]\033[0m Removing old packages on %s...\n", remoteUserHost)
		cleanCmdStr := "rm -f " + remotePkgPath + "*.deb " + remotePkgPath + "*.pkg.tar.zst"
		sshArgs := append(muxArgs, remoteUserHost, cleanCmdStr)

		cleanCmd := exec.Command("ssh", sshArgs...)
		// È fondamentale passare lo Stdin affinché il prompt della password sia visibile se clean=true
		cleanCmd.Stdout, cleanCmd.Stderr, cleanCmd.Stdin = os.Stdout, os.Stderr, os.Stdin

		if err := cleanCmd.Run(); err != nil {
			fmt.Printf("\033[1;31m[ERROR]\033[0m Remote cleanup failed: %v\n", err)
		} else {
			fmt.Println("\033[1;32m[CLEAN]\033[0m Old packages removed.")
		}
	}

	fmt.Printf("\033[1;34m[COPY]\033[0m Sending \033[1m%s\033[0m to Proxmox...\n", latestPkg)

	dstStr := fmt.Sprintf("%s:%s", remoteUserHost, remotePkgPath)
	scpArgs := append(muxArgs, latestPkg, dstStr)

	scpCmd := exec.Command("scp", scpArgs...)
	// Passiamo lo Stdin: se clean=false, sarà questo comando a chiedere la password per primo
	scpCmd.Stdout, scpCmd.Stderr, scpCmd.Stdin = os.Stdout, os.Stderr, os.Stdin

	if err := scpCmd.Run(); err != nil {
		fmt.Printf("\033[1;31m[ERROR]\033[0m SCP transfer failed: %v\n", err)
		return
	}
	fmt.Printf("\033[1;32m[SUCCESS]\033[0m %s successfully exported to Proxmox.\n", latestPkg)
}
