package link

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ketuk.id/api/internal/apierr"
	"ketuk.id/api/internal/invitation"
	"ketuk.id/api/internal/planner"
)

type Service struct {
	pool        *pgxpool.Pool
	invitations *invitation.Service
	planners    *planner.Service
}

func New(pool *pgxpool.Pool, invitations *invitation.Service, planners *planner.Service) *Service {
	return &Service{pool: pool, invitations: invitations, planners: planners}
}

type Link struct {
	ID           uuid.UUID `json:"id"`
	InvitationID uuid.UUID `json:"invitation_id"`
	PlannerID    uuid.UUID `json:"planner_id"`
}

func (s *Service) Create(ctx context.Context, owner, invitationID, plannerID uuid.UUID) (Link, error) {
	if _, err := s.invitations.Get(ctx, owner, invitationID); err != nil {
		return Link{}, err
	}
	if _, err := s.planners.Get(ctx, owner, plannerID); err != nil {
		return Link{}, err
	}
	l := Link{ID: uuid.New(), InvitationID: invitationID, PlannerID: plannerID}
	_, err := s.pool.Exec(ctx, `INSERT INTO resource_links (id, invitation_id, planner_id) VALUES ($1,$2,$3)`,
		l.ID, l.InvitationID, l.PlannerID)
	if err != nil {
		return Link{}, err
	}
	return l, nil
}

// List returns links whose invitation and planner are both owned by the caller.
// Either filter may be nil; passing both narrows to a single pair.
func (s *Service) List(ctx context.Context, owner uuid.UUID, invitationID, plannerID *uuid.UUID) ([]Link, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT l.id, l.invitation_id, l.planner_id
		FROM resource_links l
		JOIN invitations i ON i.id = l.invitation_id AND i.owner_id = $1
		JOIN planners p ON p.id = l.planner_id AND p.owner_id = $1
		WHERE ($2::uuid IS NULL OR l.invitation_id = $2)
		  AND ($3::uuid IS NULL OR l.planner_id = $3)
		ORDER BY l.created_at DESC`, owner, invitationID, plannerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Link{}
	for rows.Next() {
		var l Link
		if err := rows.Scan(&l.ID, &l.InvitationID, &l.PlannerID); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (s *Service) Delete(ctx context.Context, owner, id uuid.UUID) error {
	var invitationID, plannerID uuid.UUID
	err := s.pool.QueryRow(ctx, `SELECT invitation_id, planner_id FROM resource_links WHERE id=$1`, id).
		Scan(&invitationID, &plannerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return apierr.NotFound("link")
	}
	if err != nil {
		return err
	}
	if _, err := s.invitations.Get(ctx, owner, invitationID); err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `DELETE FROM resource_links WHERE id=$1`, id)
	return err
}
