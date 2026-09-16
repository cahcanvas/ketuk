package invitation

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"ketuk.id/api/internal/apierr"
)

type Inviter struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	SortOrder int       `json:"sort_order"`
}

func (s *Service) AddInviter(ctx context.Context, owner, invitationID uuid.UUID, name, role string, sortOrder int) (Inviter, error) {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return Inviter{}, err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return Inviter{}, apierr.BadRequest("invalid_input", "name is required")
	}
	in := Inviter{ID: uuid.New(), Name: name, Role: strings.TrimSpace(role), SortOrder: sortOrder}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO invitation_inviters (id, invitation_id, name, role, sort_order)
		VALUES ($1,$2,$3,$4,$5)`,
		in.ID, invitationID, in.Name, in.Role, in.SortOrder)
	return in, err
}

func (s *Service) ListInviters(ctx context.Context, owner, invitationID uuid.UUID) ([]Inviter, error) {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return nil, err
	}
	return s.listInviters(ctx, invitationID)
}

// InviterPatch carries only the fields the host sent; nil means "leave as is".
type InviterPatch struct {
	Name      *string `json:"name"`
	Role      *string `json:"role"`
	SortOrder *int    `json:"sort_order"`
}

func (s *Service) PatchInviter(ctx context.Context, owner, invitationID, inviterID uuid.UUID, in InviterPatch) (Inviter, error) {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return Inviter{}, err
	}
	var current Inviter
	err := s.pool.QueryRow(ctx, `
		SELECT id, name, role, sort_order FROM invitation_inviters
		WHERE id=$1 AND invitation_id=$2`, inviterID, invitationID).
		Scan(&current.ID, &current.Name, &current.Role, &current.SortOrder)
	if errors.Is(err, pgx.ErrNoRows) {
		return Inviter{}, apierr.NotFound("inviter")
	}
	if err != nil {
		return Inviter{}, err
	}
	patched, err := applyInviterPatch(current, in)
	if err != nil {
		return Inviter{}, err
	}
	_, err = s.pool.Exec(ctx, `
		UPDATE invitation_inviters SET name=$3, role=$4, sort_order=$5
		WHERE id=$1 AND invitation_id=$2`,
		patched.ID, invitationID, patched.Name, patched.Role, patched.SortOrder)
	return patched, err
}

// applyInviterPatch trims what was sent; an empty role clears it, an empty name is rejected.
func applyInviterPatch(in Inviter, patch InviterPatch) (Inviter, error) {
	if patch.Name != nil {
		name := strings.TrimSpace(*patch.Name)
		if name == "" {
			return Inviter{}, apierr.BadRequest("invalid_input", "name is required")
		}
		in.Name = name
	}
	if patch.Role != nil {
		in.Role = strings.TrimSpace(*patch.Role)
	}
	if patch.SortOrder != nil {
		in.SortOrder = *patch.SortOrder
	}
	return in, nil
}

func (s *Service) DeleteInviter(ctx context.Context, owner, invitationID, inviterID uuid.UUID) error {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `DELETE FROM invitation_inviters WHERE id=$1 AND invitation_id=$2`, inviterID, invitationID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apierr.NotFound("inviter")
	}
	return nil
}

func (s *Service) listInviters(ctx context.Context, invitationID uuid.UUID) ([]Inviter, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, role, sort_order FROM invitation_inviters
		WHERE invitation_id=$1 ORDER BY sort_order, name`, invitationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Inviter
	for rows.Next() {
		var in Inviter
		if err := rows.Scan(&in.ID, &in.Name, &in.Role, &in.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, in)
	}
	if out == nil {
		out = []Inviter{}
	}
	return out, rows.Err()
}
