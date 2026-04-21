CC := gcc
CFLAGS ?= -g -Wall -Wextra -std=c11

SRC := $(wildcard *.c)
BIN := $(SRC:.c=)

# Collect directories referenced by quoted includes, then add them to include path.
INCLUDE_HINTS := $(shell sed -n 's/^[[:space:]]*#include[[:space:]]*"\([^"]*\)".*/\1/p' $(SRC) 2>/dev/null)
INCLUDE_DIRS := $(sort $(dir $(INCLUDE_HINTS)))
CPPFLAGS += -I. $(addprefix -I,$(filter-out ./,$(INCLUDE_DIRS)))

.PHONY: all clean list

all: $(BIN)

%: %.c
	$(CC) $(CPPFLAGS) $(CFLAGS) $< -o $@

list:
	@echo "Sources: $(SRC)"
	@echo "Binaries: $(BIN)"
	@echo "Include dirs: . $(filter-out ./,$(INCLUDE_DIRS))"

clean:
	rm -f $(BIN)
