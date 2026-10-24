package config

import (
	"strings"
	"testing"
)

// @covers AC-EXECUTORS-CREDENTIAL-TRUST-001.1
func TestUnisolatedMultiTenantStartupFailsClosed(t *testing.T) {
	cfg := minimalValidConfig()
	cfg.Features.Auth = true
	cfg.Features.MultiTenancy = true
	err := validate(cfg)
	if err == nil || !strings.Contains(err.Error(), "credential isolation") {
		t.Fatalf("multi-tenant startup accepted without isolation: %v", err)
	}
	cfg.Features.MultiTenancy = false
	if err := validate(cfg); err != nil {
		t.Fatalf("single-user startup denied: %v", err)
	}
}
