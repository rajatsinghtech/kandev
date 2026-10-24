package store

import "context"

func (s *Store) ListInvitesInOrg(ctx context.Context, orgID string) ([]*Invite, error) {
	invites := []*Invite{}
	err := s.ro.SelectContext(ctx, &invites, s.ro.Rebind("SELECT * FROM auth_invites WHERE org_id = ? ORDER BY created_at DESC"), orgID)
	return invites, err
}

func (s *Store) DeleteInviteInOrg(ctx context.Context, id, orgID string) error {
	result, err := s.db.ExecContext(ctx, s.db.Rebind("DELETE FROM auth_invites WHERE id = ? AND org_id = ?"), id, orgID)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
