# 🐧 oa: Action Reference Manual

Every operation in **oa** is driven by a JSON "Plan."
This document defines the available actions, their parameters, and their expected behavior on the system.

---

## 🏗️ 1. action_prepare
**Purpose**: Initializes the Zero-Copy environment using OverlayFS and bind mounts, with built-in protections against infinite scanning loops.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `pathLiveFs` | String | The base directory for the remastering process. |
| `mode` | String | Operation mode: `""` (default), `"clone"`, or `"crypted"`. |

**Behavior**:
1. Creates the base directory structure: `liveroot/` and `.overlay/` (with `lowerdir`, `upperdir`, and `workdir` inside).
2. Performs a physical copy of `/etc` to the `liveroot`.
3. Bind-mounts root entries (e.g., `/bin`, `/sbin`, `/lib`) in read-only mode using `MS_PRIVATE` propagation.
4. **Smart `/home` Handling**: If mode is `"clone"` or `"crypted"`, `/home` is bind-mounted read-only. For `"standard"`, it is created as an empty directory.
5. Projects `/usr` and `/var` using **OverlayFS** to allow modifications without touching the host.
6. Bind-mounts kernel API filesystems: `/proc`, `/sys`, `/run`, `/dev` into `liveroot/`.
7. **Anti-Recursion Mask**: Mounts a `tmpfs` layer over the working directory path inside the `liveroot` to completely shield tools like `mksquashfs` or `nftw` from falling into infinite "Inception" loops.

---

## 👤 2. action_users
**Purpose**: Creates the Live user identity within the `liveroot` independently, without relying on host binaries.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `pathLiveFs` | String | The base directory for the remastering process. |
| `users` | Array of Objects | Contains user definitions (`login`, `password`, `home`, `shell`, `gecos`, `groups`). |
| `mode` | String | Operation mode: `""` (default), `"clone"`, or `"crypted"`. |

**Behavior**:
1. If `mode` is `"standard"`, purges host identities by sanitizing `/etc/passwd`, `/etc/shadow`, and `/etc/group` (removing UIDs between 1000 and 59999).
2. Opens `liveroot/etc/passwd` and `liveroot/etc/shadow` directly via C file streams.
3. Writes the new user identities and passwords natively using Yocto-inspired helper functions.
4. Injects secondary groups (e.g., `sudo`, `cdrom`) directly into `/etc/group`.
5. Creates the user's home directory, populates it from `/etc/skel` (including hidden files), and sets recursive ownership to the new UID/GID.

---

## ⚙️ 3. action_initrd
**Purpose**: Generates the Initial RAM Disk for the live session via template substitution.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `pathLiveFs` | String | The base directory for the remastering process. |
| `initrd_cmd` | String | Shell template to generate the initrd (e.g., `mkinitramfs -o {{out}} {{ver}}`). |

**Behavior**:
1. Detects the host's kernel version using the `uname` syscall.
2. Replaces the `{{out}}` placeholder with the target `initrd.img` path.
3. Replaces the `{{ver}}` placeholder with the detected kernel version.
4. Executes the finalized command to build the initramfs.

---

## 🗂️ 4. action_livestruct
**Purpose**: Prepares the core live directory structure and extracts the host kernel.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `pathLiveFs` | String | The base directory for the remastering process. |

**Behavior**:
1. Creates the `iso/live` directory.
2. Detects the host's running kernel version via `uname`.
3. Copies the corresponding `vmlinuz` from `/boot` into the live directory (with fallback to `/vmlinuz` symlink).

---

## 🖥️ 5. action_isolinux
**Purpose**: Populates legacy BIOS bootloader binaries and configuration.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `pathLiveFs` | String | The base directory for the remastering process. |

**Behavior**:
1. Creates the `iso/isolinux` directory.
2. Copies `isolinux.bin` and BIOS modules (`*.c32`) from `/usr/lib/`.
3. Generates a default `isolinux.cfg` boot menu if it does not already exist.

---

## 🛡️ 6. action_uefi
**Purpose**: Prepares the directory structure for UEFI booting.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `pathLiveFs` | String | The base directory for the remastering process. |

**Behavior**:
1. Creates the `iso/EFI/BOOT` directory.
2. Creates the `iso/boot/grub` directory.

---

## ⏸️ 7. action_sus