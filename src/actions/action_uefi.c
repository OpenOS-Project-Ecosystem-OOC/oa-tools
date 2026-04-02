/*
* oa: remastering core
*
* Author: Piero Proietti <piero.proietti@gmail.com>
* License: GPL-3.0-or-later
*/
#include "oa.h"

int action_uefi(OA_Context *ctx) {
    cJSON *pathLiveFs = cJSON_GetObjectItemCaseSensitive(ctx->task, "pathLiveFs");
    if (!pathLiveFs) pathLiveFs = cJSON_GetObjectItemCaseSensitive(ctx->root, "pathLiveFs");
    if (!cJSON_IsString(pathLiveFs)) return 1;

    char efi_dir[PATH_SAFE], grub_dir[PATH_SAFE];
    snprintf(efi_dir, PATH_SAFE, "%s/iso/EFI/BOOT", pathLiveFs->valuestring);
    snprintf(grub_dir, PATH_SAFE, "%s/iso/boot/grub", pathLiveFs->valuestring);

    char cmd[CMD_MAX];
    snprintf(cmd, sizeof(cmd), "mkdir -p %s %s", efi_dir, grub_dir);
    system(cmd);

    printf("\033[1;34m[oa UEFI]\033[0m Preparing UEFI boot directories...\n");
    
    // TODO: Aggiungere l'estrazione di grubx64.efi, bootx64.efi e la generazione di grub.cfg

    return 0;
}