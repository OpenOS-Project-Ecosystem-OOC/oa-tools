/* yocto_ids ispirati a OpenEmbedded-Core  */
#define OE_UID_ROOT           0
#define OE_UID_SYSTEM_MAX     999
#define OE_UID_HUMAN_MIN      1000
#define OE_UID_HUMAN_MAX      59999

// Prototipi per la gestione utenti
int action_users(cJSON *json);
bool yocto_is_human_user(uint32_t uid, const char *home);