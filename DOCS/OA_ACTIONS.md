# 🐧 oa: Action Reference Manual

Every operation in **oa** is driven by a JSON "Plan." 

### `oa-format`
Viene utilizzato solo per l'installazione, per preparare il device target

### `oa-mount`
Prepara la chroot per per la rimasterizzazione

### `oa-shell`
Viene utilizzato da coa per creare azioni complesse: `coa-boot-menus`, `coa-identity`, `coa-initd`, `coa-layout`, `coa-live-bootloader`, `coa-live-struct`, `coa-squashfs`, `coa-xorriso`.

### `oa-umount`
Smonta la chroot per per la rimasterizzazione

### `oa-users`
Crea gli utenti e/o li rimuove modificando direttamente i files `/etc/passwd`, `/etc/shadow`, `/etc/users`


