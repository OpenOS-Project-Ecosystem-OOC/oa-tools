/*
 * oa: eggs in my dialect🥚🥚
 *
 * src/vendors/oe-yocto.c
 * Logica di classificazione utenti basata su OpenEmbedded-Core
 * e sulla filosofia di penguins-eggs.
 */
#include "oa.h"

/**
 * @brief Verifica se il percorso della home è in una whitelist di sistema.
 */
static bool is_path_allowed(const char *home) {
    if (home == NULL || strlen(home) < 2) {
        return false;
    }

    const char *whitelist[] = {"home", "opt", "srv", "usr", "var", NULL};
    char path_tmp[PATH_SAFE];
    
    strncpy(path_tmp, home, sizeof(path_tmp));
    path_tmp[sizeof(path_tmp) - 1] = '\0'; // Sicurezza extra per il terminatore

    char *fLevel = strtok(path_tmp, "/");
    if (fLevel == NULL) return false;

    for (int i = 0; whitelist[i] != NULL; i++) {
        if (strcmp(fLevel, whitelist[i]) == 0) {
            return true;
        }
    }

    return false;
}

/**
 * @brief yocto_is_human_user
 * Decide se un utente dell'host deve essere processato.
 */
bool yocto_is_human_user(uint32_t uid, const char *home) {
    // 1. Filtro UID basato su OE-Core (1000-59999)
    if (uid < OE_UID_HUMAN_MIN || uid > OE_UID_HUMAN_MAX) {
        return false;
    }

    // 2. Controllo Whitelist dei percorsi
    if (!is_path_allowed(home)) {
        return false;
    }

    // 3. Verifica fisica
    struct stat st;
    if (stat(home, &st) != 0 || !S_ISDIR(st.st_mode)) {
        return false;
    }

    // 4. Analisi sottocartelle vietate
    if (strstr(home, "/cache") || strstr(home, "/run") || strstr(home, "/spool")) {
        return false;
    }

    return true;
}