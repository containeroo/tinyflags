# Shared GNU Make core for dev-tools.
# Upstream: https://github.com/gi8lino/dev-tools
# Version: 0.9.0

DEV_TOOLS_ROOT := $(patsubst %/,%,$(dir $(lastword $(MAKEFILE_LIST))))
DEV_TOOLS_CACHE := $(if $(strip $(DEV_TOOLS_VERSION)),$(DEV_TOOLS_ROOT)/.dev-tools/$(DEV_TOOLS_VERSION))
DEV_TOOLS_BIN ?= $(if $(DEV_TOOLS_CACHE),$(DEV_TOOLS_CACHE),$(DEV_TOOLS_ROOT))

GO_INSTALL_TOOL := $(DEV_TOOLS_BIN)/go-install-tool
GITHUB_RELEASE_INSTALL := $(DEV_TOOLS_BIN)/github-release-install

dev-tools-module = $(DEV_TOOLS_BIN)/dev-tools-$(1).mk

# GitHub's release URL can time out while its release asset API is available.
define fetch-dev-asset
curl --fail --silent --show-error --location --retry 2 \
	"https://github.com/gi8lino/dev-tools/releases/download/$(DEV_TOOLS_VERSION)/$(1)" \
	-o "$$tmp" || gh release download "$(DEV_TOOLS_VERSION)" --repo gi8lino/dev-tools \
	--pattern "$(1)" --output - > "$$tmp"
endef

# Download a release file atomically.
define download-dev-file
@set -eu; \
if [ -z "$(DEV_TOOLS_VERSION)" ]; then \
	echo "DEV_TOOLS_VERSION must be set before downloading dev-tools assets" >&2; \
	exit 2; \
fi; \
tmp="$(2).tmp"; \
trap 'rm -f "$$tmp"' EXIT INT TERM; \
printf '%s\n' "Downloading gi8lino/dev-tools $(DEV_TOOLS_VERSION) $(1)"; \
$(call fetch-dev-asset,$(1)); \
mv "$$tmp" "$(2)"; \
trap - EXIT INT TERM
endef

# Download an executable release asset atomically.
define download-dev-tool
@set -eu; \
if [ -z "$(DEV_TOOLS_VERSION)" ]; then \
	echo "DEV_TOOLS_VERSION must be set before downloading dev-tools assets" >&2; \
	exit 2; \
fi; \
tmp="$(2).tmp"; \
trap 'rm -f "$$tmp"' EXIT INT TERM; \
printf '%s\n' "Downloading gi8lino/dev-tools $(DEV_TOOLS_VERSION) $(1)"; \
$(call fetch-dev-asset,$(1)); \
chmod +x "$$tmp"; \
mv "$$tmp" "$(2)"; \
trap - EXIT INT TERM
endef

# Run a local tool while displaying only its executable name.
define run-tool
@printf '%s\n' '$(notdir $(1)) $(strip $(2))'
@$(1) $(2)
endef

ifneq ($(DEV_TOOLS_CACHE),)
$(DEV_TOOLS_BIN):
	@mkdir -p "$@"

# Feature modules are cached under their pinned release version. A version
# change therefore selects a different path instead of relying on mtimes.
$(DEV_TOOLS_BIN)/dev-tools-%.mk: | $(DEV_TOOLS_BIN)
	$(call download-dev-file,dev-tools-$*.mk,$@)

# Generic installers are part of the core because projects use them for their own tools.
$(GO_INSTALL_TOOL): | $(DEV_TOOLS_BIN)
	$(call download-dev-tool,go-install-tool,$@)

$(GITHUB_RELEASE_INSTALL): | $(DEV_TOOLS_BIN)
	$(call download-dev-tool,github-release-install,$@)
endif
