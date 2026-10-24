package config

import "errors"

// ValidateExecutionIsolation also runs after persisted feature overrides resolve.
func ValidateExecutionIsolation(cfg *Config) error {
	if cfg.Features.MultiTenancy {
		return errors.New("Organizations mode requires credential isolation that is not supported by this fork; use separate OS accounts and instances")
	}
	return nil
}
