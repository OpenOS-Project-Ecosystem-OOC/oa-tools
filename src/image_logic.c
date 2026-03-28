// src/image_logic.c (VERSIONE "PIERO PROIETTI" - 10 anni di Eggs)

#include "image_logic.h"
#include "cJSON.h"
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

// Helper per aggiungere esclusioni stile Eggs (stripping leading slash)
static void append_eggs_exclusion(char *buffer, size_t buf_size,
                                  const char *path) {
  const char *p = (path[0] == '/') ? path + 1 : path;
  // Aggiungiamo lo spazio e il path tra apici singoli per la shell
  strncat(buffer, " '", buf_size - strlen(buffer) - 1);
  strncat(buffer, p, buf_size - strlen(buffer) - 1);
  strncat(buffer, "'", buf_size - strlen(buffer) - 1);
}

int action_squash(cJSON *json) {
  cJSON *pathLiveFs = cJSON_GetObjectItemCaseSensitive(json, "pathLiveFs");
  cJSON *comp = cJSON_GetObjectItemCaseSensitive(json, "compression");
  cJSON *exclude_file = cJSON_GetObjectItemCaseSensitive(json, "exclude_list");
  cJSON *include_root_home =
      cJSON_GetObjectItemCaseSensitive(json, "include_root_home");

  if (!cJSON_IsString(pathLiveFs))
    return 1;

  char liveroot[1024], squash_out[1024], live_dir[1024];
  snprintf(liveroot, 1024, "%s/liveroot", pathLiveFs->valuestring);
  snprintf(live_dir, 1024, "%s/iso/live", pathLiveFs->valuestring);
  snprintf(squash_out, 1024, "%s/filesystem.squashfs", live_dir);

  // Creazione directory di output
  char cmd_tmp[2048];
  snprintf(cmd_tmp, sizeof(cmd_tmp), "mkdir -p %s", live_dir);
  system(cmd_tmp);

  printf(
      "{\"status\": \"imaging\", \"step\": \"squashfs\", \"output\": \"%s\"}\n",
      squash_out);

  // --- COSTRUZIONE ESCLUSIONI "EGGS STYLE" ---
  char session_excludes[4096] = "";

  // Esclusioni standard di Eggs (senza slash iniziale!)
  const char *fexcludes[] = {"boot/efi/EFI",
                             "boot/loader/entries/",
                             "etc/fstab",
                             "var/lib/containers/",
                             "var/lib/docker/",
                             "etc/mtab",
                             "etc/udev/rules.d/70-persistent-cd.rules",
                             "etc/udev/rules.d/70-persistent-net.rules",
                             "proc/*",
                             "sys/*",
                             "dev/*",
                             "run/*",
                             "tmp/*"};

  for (size_t i = 0; i < sizeof(fexcludes) / sizeof(fexcludes[0]); i++) {
    append_eggs_exclusion(session_excludes, sizeof(session_excludes),
                          fexcludes[i]);
  }

  // Gestione root home (Eggs style)
  if (!cJSON_IsBool(include_root_home) || !include_root_home->valueint) {
    append_eggs_exclusion(session_excludes, sizeof(session_excludes), "root/*");
    append_eggs_exclusion(session_excludes, sizeof(session_excludes),
                          "root/.*");
  }

  // --- PREPARAZIONE COMANDO FINALE ---
  char cmd[8192];
  const char *comp_str = cJSON_IsString(comp) ? comp->valuestring : "zstd";

  // Costruiamo il comando base
  snprintf(cmd, sizeof(cmd),
           "mksquashfs %s %s -comp %s -b 1M -noappend -wildcards", liveroot,
           squash_out, comp_str);

  // Aggiungiamo la exclude.list se presente (il flag -ef di Cesare)
  if (cJSON_IsString(exclude_file)) {
    snprintf(cmd + strlen(cmd), sizeof(cmd) - strlen(cmd), " -ef %s",
             exclude_file->valuestring);
  }

  // Aggiungiamo le esclusioni di sessione (-e ...)
  if (strlen(session_excludes) > 0) {
    snprintf(cmd + strlen(cmd), sizeof(cmd) - strlen(cmd), " -e%s",
             session_excludes);
  }

  printf("{\"status\": \"info\", \"msg\": \"Executing Eggs-optimized "
         "mksquashfs...\"}\n");

  int res = system(cmd);

  if (res == 0) {
    printf("{\"status\": \"ok\", \"action\": \"squash_complete\"}\n");
    return 0;
  } else {
    fprintf(stderr, "{\"status\": \"error\", \"action\": \"squash_failed\"}\n");
    return 1;
  }
}

/**
 * @brief Crea l'immagine ISO avviabile dal filesystem preparato.
 */
int action_iso(cJSON *json) {
  cJSON *pathLiveFs = cJSON_GetObjectItemCaseSensitive(json, "pathLiveFs");
  cJSON *volid = cJSON_GetObjectItemCaseSensitive(json, "volume_id");
  cJSON *iso_name = cJSON_GetObjectItemCaseSensitive(json, "filename");

  if (!cJSON_IsString(pathLiveFs))
    return 1;

  char iso_root[1024], output_iso[1024];
  snprintf(iso_root, 1024, "%s/iso", pathLiveFs->valuestring);
  snprintf(output_iso, 1024, "%s/%s", pathLiveFs->valuestring,
           cJSON_IsString(iso_name) ? iso_name->valuestring
                                    : "live-system.iso");

  printf("{\"status\": \"imaging\", \"step\": \"iso\", \"output\": \"%s\"}\n",
         output_iso);

  char cmd[8192];
  // Usiamo il comando xorriso "pesante" dal log di Eggs per compatibilità
  // BIOS/EFI
  snprintf(cmd, 8192,
           "xorriso -as mkisofs -J -joliet-long -r -l -iso-level 3 "
           "-isohybrid-mbr /usr/lib/ISOLINUX/isohdpfx.bin "
           "-partition_offset 16 -V %s "
           "-b isolinux/isolinux.bin -c isolinux/boot.cat "
           "-no-emul-boot -boot-load-size 4 -boot-info-table "
           "-eltorito-alt-boot -e boot/grub/efi.img -isohybrid-gpt-basdat "
           "-no-emul-boot -o %s %s/",
           cJSON_IsString(volid) ? volid->valuestring : "VITELIUS_LIVE",
           output_iso, iso_root);

  return system(cmd);
}
