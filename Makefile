# --- Configurazione del Compilatore ---
CC      := gcc
CFLAGS  := -Wall -Wextra -O2 -g -D_GNU_SOURCE
LDFLAGS := 

# --- Directory di Progetto ---
SRC_DIR := src
OBJ_DIR := obj
BIN_NAME := oa

# --- Ricerca Sorgenti e Oggetti ---
# Include cJSON.c se è nella cartella src
SRCS := $(wildcard $(SRC_DIR)/*.c)
OBJS := $(SRCS:$(SRC_DIR)/%.c=$(OBJ_DIR)/%.o)

# --- Regole Principali ---
.PHONY: all clean prepare_dirs

all: prepare_dirs $(BIN_NAME)

# Creazione dell'eseguibile
$(BIN_NAME): $(OBJS)
	@echo "  LD    $@"
	@$(CC) $(OBJS) -o $@ $(LDFLAGS)
	@echo "Build completata con successo: ./"$@

# Compilazione dei singoli file oggetto
$(OBJ_DIR)/%.o: $(SRC_DIR)/%.c
	@echo "  CC    $<"
	@$(CC) $(CFLAGS) -c $< -o $@

# Preparazione directory degli oggetti
prepare_dirs:
	@mkdir -p $(OBJ_DIR)

# Pulizia
clean:
	@echo "Pulizia in corso..."
	@rm -rf $(OBJ_DIR) $(BIN_NAME)
	@echo "Tutto pulito!"

# Regola per un rebuild veloce
re: clean all
