package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kandev/kandev/internal/db/dialect"
	"github.com/kandev/kandev/internal/user/models"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrLastOrgAdmin    = errors.New("cannot remove the last organization administrator")
)

// TenantAccountRepository keeps tenant targeting and last-admin admission atomic.
type TenantAccountRepository interface {
	ListUsersByOrg(context.Context, string) ([]*models.User, error)
	UpdateUserRoleStatusInOrg(context.Context, string, string, string, string) (*models.User, error)
}

func (r *sqliteRepository) UpdateUserRoleStatusInOrg(ctx context.Context, id, orgID, role, status string) (*models.User, error) {
	if orgID == "" {
		return nil, ErrAccountNotFound
	}
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	if err := r.lockOrgAccounts(ctx, tx, orgID); err != nil {
		return nil, err
	}
	user, err := scanUser(tx.QueryRowxContext(ctx, tx.Rebind("SELECT "+userColumns+" FROM users WHERE id = ? AND org_id = ?"), id, orgID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	if role == "" {
		role = user.Role
	}
	if status == "" {
		status = user.Status
	}
	role = models.NormalizeRole(role)
	if status != models.StatusDisabled {
		status = models.StatusActive
	}
	if err := requireOtherOrgAdmin(ctx, tx, user, role, status); err != nil {
		return nil, err
	}
	user.Role, user.Status, user.UpdatedAt = role, status, time.Now().UTC()
	if _, err := tx.ExecContext(ctx, tx.Rebind("UPDATE users SET role = ?, status = ?, updated_at = ? WHERE id = ? AND org_id = ?"), role, status, user.UpdatedAt, id, orgID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return user, nil
}

func requireOtherOrgAdmin(ctx context.Context, tx *sqlx.Tx, user *models.User, role, status string) error {
	if user.Role != models.RoleAdmin || user.Status != models.StatusActive || (role == models.RoleAdmin && status == models.StatusActive) {
		return nil
	}
	var count int
	err := tx.GetContext(ctx, &count, tx.Rebind("SELECT COUNT(*) FROM users WHERE org_id = ? AND role = ? AND status = ? AND id <> ?"), user.OrgID, models.RoleAdmin, models.StatusActive, user.ID)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrLastOrgAdmin
	}
	return nil
}

func (r *sqliteRepository) lockOrgAccounts(ctx context.Context, tx *sqlx.Tx, orgID string) error {
	if dialect.IsPostgres(r.db.DriverName()) {
		var ids []string
		return tx.SelectContext(ctx, &ids, tx.Rebind("SELECT id FROM users WHERE org_id = ? ORDER BY id FOR UPDATE"), orgID)
	}
	// Acquire SQLite's writer before reading the administrator set.
	_, err := tx.ExecContext(ctx, tx.Rebind("UPDATE users SET updated_at = updated_at WHERE org_id = ?"), orgID)
	return err
}
