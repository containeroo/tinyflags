# Detect platform for sed compatibility
SED := $(shell if [ "$(shell uname)" = "Darwin" ]; then echo gsed; else echo sed; fi)

# Find the latest tag (default to 0.0.0 if none found)
LATEST_TAG := $(shell git tag --list 'v*' --sort=-v:refname | head -n 1)
VERSION := $(shell echo $(LATEST_TAG) | sed 's/^v//' || echo "0.0.0")

.PHONY: test cover clean update patch minor major tag push

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: download
download: ## Download go packages
	go mod download

.PHONY: update-packages
update-packages: ## Update all Go packages to their latest versions
	go get -u ./...
	go mod tidy

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt .

.PHONY: vet
vet: ## Run go vet against code.
	go vet .

.PHONY: test
test: ## Run all unit tests
	go test . -count=1

.PHONY: cover
cover: ## Generate and display test coverage
	go test . -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: clean
clean: ## Clean up generated files
	find . -name '*.out' -delete

##@ Versioning

patch: ## Create a new patch release (x.y.Z+1)
	@NEW_VERSION=$$(echo "$(VERSION)" | awk -F. '{printf "%d.%d.%d", $$1, $$2, $$3+1}') && \
	git tag "v$${NEW_VERSION}" && \
	echo "Tagged v$${NEW_VERSION}"

minor: ## Create a new minor release (x.Y+1.0)
	@NEW_VERSION=$$(echo "$(VERSION)" | awk -F. '{printf "%d.%d.0", $$1, $$2+1}') && \
	git tag "v$${NEW_VERSION}" && \
	echo "Tagged v$${NEW_VERSION}"

major: ## Create a new major release (X+1.0.0)
	@NEW_VERSION=$$(echo "$(VERSION)" | awk -F. '{printf "%d.0.0", $$1+1}') && \
	git tag "v$${NEW_VERSION}" && \
	echo "Tagged v$${NEW_VERSION}"

tag: ## Show latest tag
	@echo "Latest version: $(LATEST_TAG)"

push: ## Push tags to remote
	git push origin --tags
