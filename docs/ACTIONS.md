# 🐧 oa: Action Reference Manual

Every operation in **oa** is driven by a JSON "Plan." This document defines the available actions, their parameters, and their expected behavior on the system.

---

## 🏗️ 1. action_prepare
**Purpose**: Initializes the Zero-Copy environment using OverlayFS.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `pathLiveFs` | String | The base directory for the remastering process (e.g., `/home/eggs`). |

**Behavior**:
1. Creates the directory structure: `iso/`, `liveroot/`, and `ovfs/` (upper/work).
2. Bind-mounts kernel interfaces: `/dev`, `/proc`, `/sys`, `/run` into `liveroot/`.
3. Performs the **OverlayFS** mount, projecting the host `/` into `liveroot/` as a writable layer.
4. Uses `MS_PRIVATE` propagation to isolate mounts from the host.

---

## 👤 2. action_users
**Purpose**: Creates the Live user identity within the `liveroot`.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `name` | String | The username for the Live session (e.g., `oa`). |
| `password` | String | The password for the Live user. |
| `groups` | Array | List of system groups (e.g., `["sudo", "audio", "video"]`). |

**Behavior**:
1. Enters the `liveroot` via `chroot`.
2. Executes `useradd` with UID 1000.
3. Sets the password via `chpasswd`.
4. Grants passwordless sudo privileges by creating `/etc/sudoers.d/oa-user`.

---

## 💀 3. action_skeleton
**Purpose**: Prepares the boot environment and manages user persistence based on the selected mode.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `kernel_path` | String | Source path of the kernel (e.g., `/vmlinuz`). |
| `initrd_cmd` | String | Shell template to generate the initrd (e.g., `mkinitramfs -o {{out}} {{ver}}`). |
| `groups` | Array | Default groups for the live user. |
| `mode` | String | Operation mode: `""` (default), `"clone"`, or `"crypted"`. |

**Behavior**:
1. Copies the kernel image to `iso/live/vmlinuz`.
2. Detects the kernel version and executes the `initrd_cmd` (injecting `{{out}}` and `{{ver}}`).
3. Populates `iso/isolinux/` with bootloader binaries.
4. **Logic by Mode**:
   * `""` (Default): Removes all users with UID $\ge 1000$ and creates a fresh live user.
   * `"crypted"`: Removes users with UID $\ge 1000$, encrypts system users, and creates a live user.
   * `"clone"`: Keeps all existing users and configurations unchanged.

---

## 📦 4. action_squash
**Purpose**: Compresses the `liveroot` into a high-performance SquashFS image.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `compression` | String | Algorithm (`zstd`, `xz`, `gzip`). Default: `zstd`. |
| `compression_level` | Integer | Level (e.g., 1-22 for zstd). |
| `exclude_list` | String | Path to patterns to exclude from compression. |
| `mode` | String | `""`, `"clone"`, or `"crypted"`. |

**Behavior**:
1. Detects CPU cores for multi-threaded compression.
2. Applies the `exclude_list` and hardcoded session excludes (`/proc`, `/sys`, etc.).
3. **Logic by Mode**:
   * If `mode == ""`, the `/home` directory is automatically added to the exclusion list.
   * If `mode == "clone"` or `"crypted"`, the `/home` directory is included (unless specified in `exclude_list`).
4. Generates the final `iso/live/filesystem.squashfs`.

---

## 💿 5. action_iso
**Purpose**: Masters the final bootable ISO image.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `volume_id` | String | The label of the ISO (e.g., `OA_LIVE`). |
| `filename` | String | The output filename (e.g., `custom-live.iso`). |

**Behavior**:
1. Invokes `xorriso` with hybrid boot support (BIOS/UEFI).
2. Generates the ISO in the `pathLiveFs` root.

---

## 🧹 6. action_cleanup
**Purpose**: Safely dismantles the Zero-Copy environment.

**Behavior**:
1. Recursively unmounts all bind-mounts and OverlayFS layers.
2. Uses **lazy unmount** (`MNT_DETACH`) if resources are busy to ensure host stability.
3. Deletes temporary workspace directories but preserves the generated ISO.

---
*Blueprint for oa v0.2 - Developed by Piero Proietti*