package engine

import (
	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/help"
)

// buildFlagDescription creates the full help text for a flag, including metadata
// such as allowed values, default, environment variable, deprecation, and group info.
func buildFlagDescription(flag *core.BaseFlag, globalHideEnvs bool, name string) string {
	return help.BuildFlagDescription(flag, globalHideEnvs, name)
}
