package engine

import (
	"fmt"
	"sort"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
)

// parseEnv loads unset flags from environment variables.
func (f *FlagSet) parseEnv() error {
	if err := f.parseStaticEnv(); err != nil {
		return err
	}
	return f.parseDynamicEnv()
}

// parseStaticEnv loads unset static flags from exact or derived environment keys.
func (f *FlagSet) parseStaticEnv() error {
	for _, fl := range f.staticFlagsMap {
		envKey, ok := fl.LookupEnvKey(f.envPrefix, f.envKeyFunc)
		if !ok {
			continue
		}
		val := f.getEnv(envKey)
		if val == "" {
			continue
		}
		if err := fl.Value.Set(val); err != nil {
			if f.ignoreInvalidEnv {
				continue
			}
			return fmt.Errorf("invalid value for flag --%s from environment: %w", fl.Name, err)
		}
	}
	return nil
}

// parseDynamicEnv loads dynamic flags from APP_GROUP_ID_FIELD style keys.
func (f *FlagSet) parseDynamicEnv() error {
	if f.envPrefix == "" || len(f.dynamicGroupsMap) == 0 || f.getEnvVars == nil {
		return nil
	}

	for _, entry := range f.getEnvVars() {
		key, val, ok := strings.Cut(entry, "=")
		if !ok || val == "" {
			continue
		}
		if err := f.tryParseDynamicEnv(key, val); err != nil {
			if f.ignoreInvalidEnv {
				continue
			}
			return err
		}
	}
	return nil
}

// tryParseDynamicEnv applies one environment entry if it matches a dynamic flag.
func (f *FlagSet) tryParseDynamicEnv(key, val string) error {
	for _, group := range f.dynamicGroups() {
		for _, fl := range dynamicEnvFlags(group) {
			if fl == nil || fl.DisableEnv {
				continue
			}

			template := core.DynamicEnvKey(f.envPrefix, group.Name(), "{ID}", fl.Name)
			id, ok := matchDynamicEnvKey(key, template)
			if !ok {
				continue
			}

			item, ok := group.Items()[fl.Name]
			if !ok || item.Value == nil {
				continue
			}
			if _, changed := item.Value.GetAny(id); changed {
				return nil
			}
			if err := item.Value.Set(id, val); err != nil {
				return fmt.Errorf("invalid value for flag --%s.%s.%s from environment %s: %w", group.Name(), id, fl.Name, key, err)
			}
			return nil
		}
	}
	return nil
}

func dynamicEnvFlags(group interface{ Flags() []*core.BaseFlag }) []*core.BaseFlag {
	flags := append([]*core.BaseFlag(nil), group.Flags()...)
	sort.SliceStable(flags, func(i, j int) bool {
		left := core.NormalizeEnvKeyPart(flags[i].Name)
		right := core.NormalizeEnvKeyPart(flags[j].Name)
		return len(left) > len(right)
	})
	return flags
}

// matchDynamicEnvKey extracts the dynamic ID from an ENV key template.
func matchDynamicEnvKey(key, template string) (string, bool) {
	before, after, ok := strings.Cut(template, "{ID}")
	if !ok || before == "" || after == "" {
		return "", false
	}
	if !strings.HasPrefix(key, before) || !strings.HasSuffix(key, after) {
		return "", false
	}
	idPart := key[len(before) : len(key)-len(after)]
	if idPart == "" {
		return "", false
	}
	return core.DynamicEnvID(idPart), true
}
