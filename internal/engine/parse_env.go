package engine

import "fmt"

// parseEnv loads unset flags from environment variables.
func (f *FlagSet) parseEnv() error {
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
