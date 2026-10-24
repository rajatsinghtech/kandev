package backendapp

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/kandev/kandev/internal/common/config"
	"github.com/kandev/kandev/internal/common/logger"
	"github.com/kandev/kandev/internal/runtimeflags"
)

func TestStartupRejectsPersistedUnisolatedTenancy(t *testing.T) {
	conn, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	flags, err := runtimeflags.NewSQLiteStore(conn, conn)
	if err != nil {
		t.Fatal(err)
	}
	if err := flags.SetOverride(context.Background(), "features.multiTenancy", true); err != nil {
		t.Fatal(err)
	}
	log, err := logger.NewLogger(logger.LoggingConfig{Level: "error", Format: "json"})
	if err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{}
	cfg.Features.Auth = true
	if applyStartupRuntimeFlags(context.Background(), cfg, &Repositories{RuntimeFlags: flags}, log) {
		t.Fatal("persisted Organizations override bypassed the credential-isolation boundary")
	}
}
