# Makefile

.DEFAULT_GOAL := help

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint

## Tool Versions
# renovate: datasource=github-releases depName=golangci/golangci-lint
GOLANGCI_LINT_VERSION ?= v2.15.0

# renovate: datasource=github-releases depName=gi8lino/dev-tools
DEV_TOOLS_VERSION ?= v0.9.0

# renovate: datasource=github-releases depName=kumbuka-me/cli
KUMBUKA_CLI_VERSION ?= v0.18.1
KUMBUKA_CLI := $(LOCALBIN)/.kumbuka/$(KUMBUKA_CLI_VERSION)/kumbuka-cli
KUMBUKA_CLI_ASSET ?= kumbuka-cli_{version}_{os}_{arch}.tar.gz

include bin/dev-tools.mk

## Documentation Configuration
SITE_CONFIG ?= docs/site.toml
SITE_PORT ?= 8000
SITE_URL ?= http://127.0.0.1:$(SITE_PORT)/
BUILD_ARGS ?=

# Default: no prefix. Can be overridden via `make patch VERSION_PREFIX=v`
VERSION_PREFIX ?= "v"

##@ Tagging

# Find the latest tag (with prefix filter if defined, default to 0.0.0 if none found)
# Lazy evaluation ensures fresh values on every run
LATEST_TAG = $(shell git tag --list "$(VERSION_PREFIX)*" --sort=-v:refname | head -n 1)
VERSION = $(shell [ -n "$(LATEST_TAG)" ] && echo $(LATEST_TAG) | sed "s/^$(VERSION_PREFIX)//" || echo "0.0.0")

patch: ## Create a new patch release (x.y.Z+1)
	@NEW_VERSION=$$(echo "$(VERSION)" | awk -F. '{printf "%d.%d.%d", $$1, $$2, $$3+1}') && \
	git tag "$(VERSION_PREFIX)$${NEW_VERSION}" && \
	echo "Tagged $(VERSION_PREFIX)$${NEW_VERSION}"

minor: ## Create a new minor release (x.Y+1.0)
	@NEW_VERSION=$$(echo "$(VERSION)" | awk -F. '{printf "%d.%d.0", $$1, $$2+1}') && \
	git tag "$(VERSION_PREFIX)$${NEW_VERSION}" && \
	echo "Tagged $(VERSION_PREFIX)$${NEW_VERSION}"

major: ## Create a new major release (X+1.0.0)
	@NEW_VERSION=$$(echo "$(VERSION)" | awk -F. '{printf "%d.0.0", $$1+1}') && \
	git tag "$(VERSION_PREFIX)$${NEW_VERSION}" && \
	echo "Tagged $(VERSION_PREFIX)$${NEW_VERSION}"

tag: ## Show latest tag
	@echo "Latest version: $(LATEST_TAG)"

push: ## Push tags to remote
	git push --tags

##@ Development

.PHONY: download
download: ## Download go packages
	go mod download

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run unit tests.
	go test -covermode=atomic -count=1 -parallel=4 -timeout=5m ./...

.PHONY: cover
cover: ## Display test coverage
	go test -coverprofile=coverage.out -covermode=atomic -count=1 -parallel=4 -timeout=5m ./...
	go tool cover -html=coverage.out

.PHONY: clean
clean: ## Clean up generated files
	rm -f coverage.out coverage.html
	rm -rf docs/site

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter.
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes.
	$(GOLANGCI_LINT) run --fix

##@ Documentation

.PHONY: build-docs
build-docs: cli ## Build the production Kumbuka documentation site.
	$(call run-tool,$(KUMBUKA_CLI),build --config "$(SITE_CONFIG)" $(BUILD_ARGS))

.PHONY: serve-docs
serve-docs: cli ## Serve and watch documentation locally.
	$(call run-tool,$(KUMBUKA_CLI),serve --config "$(SITE_CONFIG)" --site-url "$(SITE_URL)" --listen-address "127.0.0.1:$(SITE_PORT)")

##@ Dependencies

$(KUMBUKA_CLI): $(GITHUB_RELEASE_INSTALL)
	@mkdir -p "$(dir $@)"
	@$(GITHUB_RELEASE_INSTALL) \
		--repo kumbuka-me/cli \
		--tag "$(KUMBUKA_CLI_VERSION)" \
		--asset "$(KUMBUKA_CLI_ASSET)" \
		--binary kumbuka-cli \
		--target "$@"

.PHONY: cli
cli: $(KUMBUKA_CLI) ## Download the pinned Kumbuka CLI release.

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
