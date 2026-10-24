package auth

import (
	"context"
	"errors"

	"github.com/kandev/kandev/internal/auth/authn"
	userstore "github.com/kandev/kandev/internal/user/store"
)

var ErrAdministrationDenied = errors.New("administration access denied")

func (s *Service) tenantAdministrator(ctx context.Context) (authn.Identity, error) {
	id, ok := authn.IdentityFromContext(ctx)
	if !s.cfg.Features.Auth || !ok || id.Synthetic || id.UserID == "" || id.OrgID == "" || (!id.IsAdmin() && id.Role != "owner") {
		return authn.Identity{}, ErrAdministrationDenied
	}
	return id, nil
}

func (s *Service) tenantAccounts(ctx context.Context) (authn.Identity, userstore.TenantAccountRepository, error) {
	id, err := s.tenantAdministrator(ctx)
	if err != nil {
		return id, nil, err
	}
	repo, ok := s.users.(userstore.TenantAccountRepository)
	if !ok {
		return id, nil, ErrAdministrationDenied
	}
	return id, repo, nil
}
