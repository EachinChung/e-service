# Build all by default, even if it's not first
.DEFAULT_GOAL := all

.PHONY: all
all: tidy build

# ==============================================================================
# Build options

ROOT_PACKAGE=github.com/eachinchung/e-service
VERSION_PACKAGE=github.com/eachinchung/component-base/version

# ==============================================================================
# Includes

# make sure include common.mk at the first include line
include scripts/make-rules/common.mk
include scripts/make-rules/golang.mk

# ==============================================================================
# Targets

.PHONY: tidy
tidy:
	@$(GO) mod tidy

## build: Build source code for host platform.
.PHONY: build
build:
	@$(MAKE) go.build

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)