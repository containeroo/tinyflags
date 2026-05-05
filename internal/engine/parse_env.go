package engine

import "fmt"

// parseEnv loads unset flags from environment variables.
func (f *FlagSet) parseEnv() error {
	for _, fl := range f.staticFlagsMap {
		if fl.Value == nil {
			// dynamically‐registered flags aren’t loaded from ENV
			continue
		}
		if fl.DisableEnv {
			continue
		}
		if fl.Value.Changed() {
			continue
		}
		envKey := fl.EnvKey
		if envKey == "" {
			envKey = f.envKeyFunc(f.envPrefix, fl.Name)
		}
		val := f.getEnv(envKey)
		if val == "" {
			continue
		}
		if err := fl.Value.Set(val); err != nil {
			if f.ignoreInvalidEnv {
				continue
			}
			return fmt.Errorf("invalid environment value for %s: %w", fl.Name, err)
		}
	}
	return nil
}
