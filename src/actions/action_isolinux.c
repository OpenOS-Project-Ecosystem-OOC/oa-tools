/*
* oa: remastering core
*
* Author: Piero Proietti <piero.proietti@gmail.com>
* License: GPL-3.0-or-later
*/
#include "oa.h"

int action_isolinux(OA_Context *ctx) {
    cJSON *pathLiveFs = cJSON_GetObjectItemCaseSensitive(ctx->task, "pathLiveFs");
    if (!pathLiveFs) pathLiveFs = cJSON_GetObjectItemCaseSensitive(ctx->root, "pathLiveFs");
    if (!cJSON_IsString(pathLiveFs)) return 1;

    char isolinux_dir[PATH_SAFE];
    snprintf(isolinux_dir, PATH_SAFE, "%s/iso/isolinux", pathLiveFs->valuestring);

    char cmd[CMD_MAX];
    snprintf(cmd, sizeof(cmd), "mkdir -p %s", isolinux_dir);
    system(cmd);

    printf("\033[1;34m[oa ISOLINUX]\033[0m Populating BIOS bootloader binaries...\n");
    snprintf(cmd, sizeof(cmd), "cp /usr/lib/ISOLINUX/isolinux.bin %s/ && "
                               "cp /usr/lib/syslinux/modules/bios/*.c32 %s/", 
             isolinux_dir, isolinux_dir);
    system(cmd);

    // Configurazione Isolinux
    char cfg_path[PATH_SAFE];
    snprintf(cfg_path, PATH_SAFE, "%s/isolinux.cfg", isolinux_dir);
    if (access(cfg_path, F_OK) != 0) {
        FILE *f = fopen(cfg_path, "w");
        if (f) {
            fprintf(f, "UI vesamenu.c32\n"
                       "PROMPT 0\n"
                       "TIMEOUT 50\n"
                       "DEFAULT live\n"
                       "MENU TITLE oa Boot Menu\n"
                       "LABEL live\n"
                       "  MENU LABEL oa Live (Standard)\n"
                       "  KERNEL /live/vmlinuz\n"
                       "  APPEND initrd=/live/initrd.img boot=live components quiet splash\n");
            fclose(f);
            printf("\033[1;32m[oa ISOLINUX]\033[0m isolinux.cfg generated.\n");
        }
    }

    return 0;
}
