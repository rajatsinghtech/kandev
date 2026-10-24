package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/kandev/kandev/internal/auth/authn"
	usermodels "github.com/kandev/kandev/internal/user/models"
)

func tenantAdministrationFixture(t *testing.T) (*serviceFixture, context.Context, context.Context) {
	t.Helper()
	f := newServiceFixture(t, true)
	for _, u := range []*usermodels.User{
		{ID: "admin-a", Email: "admin-a@example.test", Role: "admin", Status: "active", OrgID: "a"},
		{ID: "member-a", Email: "member-a@example.test", Role: "member", Status: "active", OrgID: "a"},
		{ID: "admin-b", Email: "admin-b@example.test", Role: "admin", Status: "active", OrgID: "b"},
	} {
		if err := f.svc.users.CreateUser(context.Background(), u); err != nil {
			t.Fatal(err)
		}
	}
	f.svc.cfg.Features.MultiTenancy = true
	a := authn.WithIdentity(context.Background(), authn.Identity{UserID: "admin-a", OrgID: "a", Role: authn.RoleAdmin})
	b := authn.WithIdentity(context.Background(), authn.Identity{UserID: "admin-b", OrgID: "b", Role: authn.RoleAdmin})
	return f, a, b
}

func TestLastAdminPerOrgConcurrent(t *testing.T) {
	f, a, _ := tenantAdministrationFixture(t)
	f.db.SetMaxOpenConns(1)
	if err := f.svc.users.CreateUser(context.Background(), &usermodels.User{
		ID: "admin-a2", Email: "admin-a2@example.test", Role: "admin", Status: "active", OrgID: "a",
	}); err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	results := make(chan error, 2)
	for _, id := range []string{"admin-a", "admin-a2"} {
		go func() {
			<-start
			_, err := f.svc.AdminSetRoleStatus(a, id, "member", "")
			results <- err
		}()
	}
	close(start)
	allowed, denied := 0, 0
	for range 2 {
		err := <-results
		if err == nil {
			allowed++
		} else if errors.Is(err, ErrLastAdmin) {
			denied++
		} else {
			t.Fatal(err)
		}
	}
	if allowed != 1 || denied != 1 {
		t.Fatalf("concurrent demotions: allowed=%d denied=%d", allowed, denied)
	}
}

// @covers AC-AUTH-ADMIN-CONTAINMENT-001.1
func TestOrgAdminCannotListForeignAccounts(t *testing.T) {
	f, a, _ := tenantAdministrationFixture(t)
	users, err := f.svc.ListUsers(a)
	if err != nil || len(users) != 2 {
		t.Fatalf("own org list = %v, %v", users, err)
	}
	for _, user := range users {
		if user.OrgID != "a" {
			t.Fatalf("foreign account disclosed: %s", user.ID)
		}
	}
	if _, err := f.svc.ListUsers(context.Background()); err == nil {
		t.Fatal("missing tenant identity granted account visibility")
	}
}

// @covers AC-AUTH-ADMIN-CONTAINMENT-001.2
func TestOrgAdminCannotMutateForeignAccount(t *testing.T) {
	f, a, _ := tenantAdministrationFixture(t)
	ctx := context.Background()
	_, _, err := f.svc.MintToken(ctx, "admin-b", "foreign-canary", 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.svc.createSession(ctx, "admin-b", "canary", ""); err != nil {
		t.Fatal(err)
	}
	_, foreignErr := f.svc.AdminSetRoleStatus(a, "admin-b", "member", "disabled")
	_, missingErr := f.svc.AdminSetRoleStatus(a, "missing", "member", "disabled")
	if foreignErr == nil || missingErr == nil || foreignErr.Error() != missingErr.Error() {
		t.Fatalf("foreign/missing outcomes = %v / %v", foreignErr, missingErr)
	}
	user, err := f.svc.users.GetUser(ctx, "admin-b")
	if err != nil || user.Role != "admin" || user.Status != "active" {
		t.Fatalf("foreign account changed: %+v, %v", user, err)
	}
	tokens, err := f.svc.ListTokens(ctx, "admin-b")
	if err != nil || len(tokens) != 1 {
		t.Fatalf("foreign tokens changed: %v, %v", tokens, err)
	}
	sessions, err := f.svc.ListSessions(ctx, "admin-b")
	if err != nil || len(sessions) != 1 {
		t.Fatalf("foreign sessions changed: %v, %v", sessions, err)
	}
	if _, err := f.svc.AdminSetRoleStatus(a, "member-a", "", "disabled"); err != nil {
		t.Fatalf("own member mutation denied: %v", err)
	}
}

// @covers AC-AUTH-ADMIN-CONTAINMENT-001.3
func TestLastAdminIsProtectedWithinOrganization(t *testing.T) {
	f, a, _ := tenantAdministrationFixture(t)
	if _, err := f.svc.AdminSetRoleStatus(a, "admin-a", "member", ""); err == nil {
		t.Fatal("another organization's administrator satisfied the last-admin guard")
	}
}

func TestOrgAdminCannotReadOrRevokeForeignInvites(t *testing.T) {
	f, a, b := tenantAdministrationFixture(t)
	foreign, _, err := f.svc.CreateInvite(b, "admin-b", "", "member", 0)
	if err != nil {
		t.Fatal(err)
	}
	own, _, err := f.svc.CreateInvite(a, "admin-a", "", "member", 0)
	if err != nil {
		t.Fatal(err)
	}
	list, err := f.svc.ListInvites(a)
	if err != nil || len(list) != 1 || list[0].ID != own.ID {
		t.Fatalf("invite list disclosed another org: %+v, %v", list, err)
	}
	if err := f.svc.DeleteInvite(a, foreign.ID); err == nil {
		t.Fatal("foreign invite revoked")
	}
	if err := f.svc.DeleteInvite(a, own.ID); err != nil {
		t.Fatalf("own invite revoke denied: %v", err)
	}
	list, err = f.svc.ListInvites(b)
	if err != nil || len(list) != 1 || list[0].ID != foreign.ID {
		t.Fatalf("foreign invite changed: %+v, %v", list, err)
	}
}
