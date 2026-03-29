# ACTIONS
| Action | JSON Parameters | Debian/Devuan (Default) | Divergence (Arch/Fedora/Suse) |
| :--- | :--- | :--- | :--- |
| **`prepare`** | `pathLiveFs`, `mounts[]` | Mounts `/dev, /proc, /sys, /run` | Standard POSIX (Universal) |
| **`users`** | `name`, `pass`, `groups[]`, `uid` | Admin group: `sudo` | Admin group: `wheel` |
| **`skeleton`** | `kernel_path`, `initrd_template` | `/vmlinuz`, `mkinitramfs` | `vmlinuz-linux`, `dracut/mkinitcpio` |
| **`customize`**| `hooks_path`, `env[]` | Bash scripts in `hooks.d/` | Distro-specific config triggers |
| **`squash`** | `comp`, `level`, `exclude_list` | APT/DPKG exclusions | Pacman/DNF/Zypper exclusions |
| **`iso`** | `volid`, `filename`, `boot_mode` | Isolinux/Grub (BIOS/UEFI) | Systemd-boot or custom EFI paths |
| **`cleanup`** | - | Umount OverlayFS & Bind | Standard POSIX (Universal) |
