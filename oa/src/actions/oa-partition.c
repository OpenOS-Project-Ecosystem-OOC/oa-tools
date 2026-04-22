#include "oa.h"
// #include <stdio.h>
// #include <stdlib.h>

int oa_partition(OA_Context *ctx) {
    cJSON *device_obj = cJSON_GetObjectItemCaseSensitive(ctx->task, "device");
    cJSON *label_obj  = cJSON_GetObjectItemCaseSensitive(ctx->task, "label");
    cJSON *parts_obj  = cJSON_GetObjectItemCaseSensitive(ctx->task, "partitions");

    if (!cJSON_IsString(device_obj) || !cJSON_IsString(label_obj) || !cJSON_IsArray(parts_obj)) {
        LOG_ERR("oa_partition requires 'device', 'label', and a 'partitions' array.");
        return 1;
    }

    const char *device = device_obj->valuestring;
    const char *label = label_obj->valuestring;

    printf("\033[1;35m[oa]\033[0m Inizializzazione disco %s con tabella %s...\n", device, label);
    LOG_INFO("Partitioning device %s with label %s", device, label);

    // Iteriamo sull'array delle partizioni per preparare i comandi
    cJSON *part;
    cJSON_ArrayForEach(part, parts_obj) {
        cJSON *size_obj = cJSON_GetObjectItemCaseSensitive(part, "size");
        cJSON *type_obj = cJSON_GetObjectItemCaseSensitive(part, "type");
        cJSON *name_obj = cJSON_GetObjectItemCaseSensitive(part, "name");

        const char *size = cJSON_IsString(size_obj) ? size_obj->valuestring : "";
        const char *type = cJSON_IsString(type_obj) ? type_obj->valuestring : "";
        const char *name = cJSON_IsString(name_obj) ? name_obj->valuestring : "part";

        LOG_INFO(" -> Planned partition: %s, Size: %s, Type: %s", name, size, type);
    }

    // TODO: Qui costruiremo la stringa di input per sfdisk e la eseguiremo via pipe (popen)
    
    return 0;
}