/*
* src/actions/action_users.c
* Remastering core: User & Group Identity artisan
* oa: eggs in my dialect🥚🥚
*
* Author: Piero Proietti <piero.proietti@gmail.com>
* License: GPL-3.0-or-later
*/
#include "oa.h"
#include <pwd.h>
#include <grp.h>

/**
 * @brief action_users
 * Gestisce la pulizia degli utenti host e la creazione dell'identità live.
 */
int action_users(cJSON *json) {
    cJSON *pathLiveFs = cJSON_GetObjectItemCaseSensitive(json, "pathLiveFs");
    cJSON *mode_item = cJSON_GetObjectItemCaseSensitive(json, "mode"); 
    cJSON *users_array = cJSON_GetObjectItemCaseSensitive(json, "users");

    if (!cJSON_IsString(pathLiveFs)) {
        fprintf(stderr, "{\"error\": \"pathLiveFs missing in action_users\"}\n");
        return 1;
    }

    const char *mode = cJSON_IsString(mode_item) ? mode_item->valuestring : "standard";
    char liveroot[PATH_SAFE];
    snprintf(liveroot, PATH_SAFE, "%s/liveroot", pathLiveFs->valuestring);

    // 1. FASE DI PULIZIA (Se non siamo in modalità CLONE)
    if (strcmp(mode, "clone") != 0) { 
        printf("\033[1;34m[oa]\033[0m Mode %s: Purging host users from liveroot...\n", mode);
        
        char cmd[CMD_MAX];
        snprintf(cmd, sizeof(cmd), 
                 "chroot %s /bin/bash -c \"awk -F: '$3 >= %d && $3 <= %d {print $1}' /etc/passwd | xargs -r -n1 userdel -r -f\"", 
                 liveroot, OE_UID_HUMAN_MIN, OE_UID_HUMAN_MAX);
        system(cmd);
    }

    // 2. FASE DI CREAZIONE (Identità Live dal JSON)
    if (cJSON_IsArray(users_array)) {
        cJSON *u;
        cJSON_ArrayForEach(u, users_array) {
            cJSON *login_obj = cJSON_GetObjectItemCaseSensitive(u, "login");
            cJSON *home_obj = cJSON_GetObjectItemCaseSensitive(u, "home");
            cJSON *groups = cJSON_GetObjectItemCaseSensitive(u, "groups");

            if (!cJSON_IsString(login_obj) || !cJSON_IsString(home_obj)) continue;

            const char *login = login_obj->valuestring;
            const char *home = home_obj->valuestring;

            // Validazione Yocto-Style prima di procedere
            if (!yocto_is_human_user(OE_UID_HUMAN_MIN, home)) {
                printf("\033[1;33m[oa]\033[0m User %s ignored (Path not allowed or missing home)\n", login);
                continue;
            }

            printf("\033[1;32m[oa]\033[0m Crafting user identity: %s\n", login);

            // Gestione Gruppi: Rilevamento dinamico GID dall'host
            if (cJSON_IsArray(groups)) {
                cJSON *g_node;
                cJSON_ArrayForEach(g_node, groups) {
                    if (!cJSON_IsString(g_node)) continue;

                    struct group *gr = getgrnam(g_node->valuestring);
                    if (gr) {
                        char gr_cmd[CMD_MAX];
                        // Crea il gruppo se non esiste con lo stesso GID dell'host
                        snprintf(gr_cmd, sizeof(gr_cmd), "chroot %s groupadd -g %d %s || true", 
                                 liveroot, gr->gr_gid, g_node->valuestring);
                        system(gr_cmd);
                        
                        // Aggiunge l'utente al gruppo
                        snprintf(gr_cmd, sizeof(gr_cmd), "chroot %s usermod -aG %s %s", 
                                 liveroot, g_node->valuestring, login);
                        system(gr_cmd);
                    }
                }
            }
        }
    }

    printf("{\"status\": \"ok\", \"action\": \"users_complete\"}\n");
    return 0;
}