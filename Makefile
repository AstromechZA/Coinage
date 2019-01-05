# ----------------------------------------------------------------------------------------------------------------------
# variables
# ----------------------------------------------------------------------------------------------------------------------

VERSION := $(shell git describe --always --tags --dirty)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
GIT_DATE := $(shell git log -1 --format=%cI)

ALL_GO_FILES = $(shell find . -type f -name '*.go')

BINARY_NAME := coinage
DIST_DIR := dist
DISTRIBUTABLES = $(DIST_DIR)/$(BINARY_NAME).linux.amd64 $(DIST_DIR)/$(BINARY_NAME).darwin.amd64 $(DIST_DIR)/$(BINARY_NAME).windows.amd64
GOOS_FUNC = $(shell echo $@ | grep amd64 | rev | cut -d. -f 2 | rev)
GOARCH_FUNC = $(shell echo $@ | grep amd64 | rev | cut -d. -f 1 | rev)

COVERAGE_DIR := .coverage
COVER_PACKAGES = $(shell go list ./... | grep -v vendor | grep -v coinage/cmd)
COVERAGE_FILES = $(addprefix $(COVERAGE_DIR)/,$(addsuffix .cover,$(COVER_PACKAGES)))

# ----------------------------------------------------------------------------------------------------------------------
# phony targets
# ----------------------------------------------------------------------------------------------------------------------

# help menu showing the documented targets
.PHONY: help
help:
	@echo "Choose a valid make target:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version: ## print current version information
	@echo "VERSION:    '$(VERSION)'"
	@echo "GIT_COMMIT: '$(GIT_COMMIT)'"
	@echo "GIT_DATE:   '$(GIT_DATE)'"

.PHONY: test
test: ## run unit tests and generate coverage reports
	@echo $(shell date): Running unit tests..
	go test -v ./...
	@echo
	@echo $(shell date): Vetting code..
	go vet ./...
	@echo

.PHONY: precoverage coverage
precoverage: ; @echo $(shell date): Running coverage tests..
coverage: precoverage $(COVERAGE_FILES) ## run tests and output test coverage information

.PHONY: dev
dev: $(BINARY_NAME) ## build development binary for the local platform

.PHONY: clean
clean: ## remove development files, builds, reports (does not delete IDEA files)
	@echo $(shell date): Removing old files..
	@rm -rfv $(DIST_DIR) $(BINARY_NAME) $(COVERAGE_DIR)
	@echo

.PHONY: distributables
distributables: $(DISTRIBUTABLES) $(DIST_DIR)/SHA256SUMS ## build distributable binaries and hash list

# ----------------------------------------------------------------------------------------------------------------------
# real targets
# ----------------------------------------------------------------------------------------------------------------------

# tasks to create dev binary
$(BINARY_NAME): $(ALL_GO_FILES)
	go build -o $@ ./cmd/$(BINARY_NAME)

# official distributable version
$(DISTRIBUTABLES): $(ALL_GO_FILES)
	@echo $(shell date): Building $@..
	@mkdir -p $(dir $@)
	@CGO_ENABLED=0 GOFLAGS=-mod=vendor GOOS=$(GOOS_FUNC) GOARCH=$(GOARCH_FUNC) time go build \
		-o $@ \
		-ldflags "-X main.gitHash=$(GIT_COMMIT) -X main.gitDate=$(GIT_DATE) -X main.version=$(VERSION)" \
		./cmd/$(BINARY_NAME)
	@ls -lh $@
	@echo

# shasums file for dist
$(DIST_DIR)/SHA256SUMS: $(DISTRIBUTABLES)
	@echo $(shell date): Writing $@..
ifneq "$(shell which sha256sum)" ""
	@cd $(dir $@) && sha256sum $(BINARY_NAME).* > $(notdir $@)
else
	@cd $(dir $@) && shasum -a 256 $(BINARY_NAME).* > $(notdir $@)
endif
	@cat $@
	@echo

# generate coverage for golang package
$(COVERAGE_DIR)/%.cover: $(ALL_GO_FILES)
	@mkdir -p $(dir $@)
	@go test -covermode=count -coverprofile $@ $*
