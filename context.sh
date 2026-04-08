#!/bin/bash

# Abilita nullglob per l'espansione corretta degli array
shopt -s nullglob

# 1. Genera un SINGOLO suffisso comune per l'intera esecuzione
RAND_SUFFIX=$(printf "%03d" $((RANDOM % 1000)))
FILE_OA="CONTEXT_OA_${RAND_SUFFIX}.txt"
FILE_COA="CONTEXT_COA_${RAND_SUFFIX}.txt"

echo -e "\033[1;34m[Context Builder]\033[0m Generazione sessione: \033[1m$RAND_SUFFIX\033[0m"

# 2. Funzione universale per estrarre e formattare i file
build_context() {
    local out_file=$1
    shift
    local files=("$@")

    echo -e " -> Assembling \033[1;33m$out_file\033[0m..."

    (
        echo '````'
        for f in "${files[@]}"; do
            if [ -f "$f" ]; then
                echo "### 📄 FILE: $f"
                
                filename=$(basename "$f")
                ext="${filename##*.}"

                # Mappatura estensioni -> linguaggio markdown
                case "$ext" in
                    c|h)      lang="c" ;;
                    go)       lang="go" ;;      # <-- Aggiunto Go!
                    sh)       lang="bash" ;;
                    json)     lang="json" ;;
                    md)       lang="markdown" ;;
                    yaml|yml) lang="yaml" ;;    # <-- Aggiunto Yaml!
                    *)        lang="text" ;;
                esac

                # Casi speciali per nomi senza estensione
                if [[ "$filename" == "Makefile" || "$filename" == "m" ]]; then lang="make"; fi

                echo '```'"$lang"
                cat "$f"
                echo '```'
                echo ""
            fi
        done
        echo '````'
    ) > "$out_file"
}

# 3. Definisci i percorsi per OA
FILES_OA=(
    oa/CHANGELOG.md
    oa/Makefile
    oa/MANIFESTUM.md
    oa/README.md
    oa/docs/*.md
    oa/include/*.h
    oa/json/*.json
    oa/src/*.c
    oa/src/actions/*.c
    oa/src/vendors/*.c
)

# 4. Definisci i percorsi per COA
FILES_COA=(
    coa/m
    coa/go.mod
    coa/src/*.go
    coa/conf/*.yaml
    coa/docs/ROADMAP.md
    coa/README.md
)

# 5. Costruisci i due file
build_context "$FILE_OA" "${FILES_OA[@]}"
build_context "$FILE_COA" "${FILES_COA[@]}"

# Disabilita nullglob
shopt -u nullglob

# 6. Trasferimento unificato (una sola connessione SSH per tutti e due i file)
echo -e "\033[1;32m[SCP]\033[0m Trasferimento verso artisan@192.168.1.2..."
scp "$FILE_OA" "$FILE_COA" artisan@192.168.1.2:/home/artisan/

# 7. Pulizia
rm "$FILE_OA" "$FILE_COA"

echo -e "\033[1;32m[OK]\033[0m Entrambi i contesti ($RAND_SUFFIX) sincronizzati con successo!"
