#!/bin/sh
# ~/penguins-eggs/dracut/modules.d/00debug-shell/debug-hook.sh

# Controlla se il parametro "startdebug" è stato passato al kernel
if getargbool 0 startdebug; then
    echo "############################################"
    echo "### FORZANDO LA SHELL DI DEBUG"
    echo "############################################"
    
    # Lancia una shell interattiva
    /bin/sh
fi
