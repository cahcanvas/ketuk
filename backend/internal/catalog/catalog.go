package catalog

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ketuk.id/api/internal/apierr"
)

type Service struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Service { return &Service{pool: pool} }

type EventType struct {
	ID                uuid.UUID       `json:"id"`
	Slug              string          `json:"slug"`
	Name              string          `json:"name"`
	InvitationSchema  json.RawMessage `json:"invitation_schema"`
	PlannerCategories json.RawMessage `json:"planner_categories"`
}

type Template struct {
	ID      uuid.UUID       `json:"id"`
	Slug    string          `json:"slug"`
	TypeID  uuid.UUID       `json:"type_id"`
	Name    string          `json:"name"`
	Tier    string          `json:"tier"`
	Version int             `json:"version"`
	Slots   json.RawMessage `json:"slots"`
	Locked  bool            `json:"locked"`
}

func (s *Service) ListTypes(ctx context.Context) ([]EventType, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, slug, name, invitation_schema, planner_categories FROM event_types ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []EventType
	for rows.Next() {
		var t EventType
		if err := rows.Scan(&t.ID, &t.Slug, &t.Name, &t.InvitationSchema, &t.PlannerCategories); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *Service) ListTemplates(ctx context.Context, typeSlug string, allowPremium bool) ([]Template, error) {
	q := `SELECT t.id, t.slug, t.type_id, t.name, t.tier, t.version, t.slots
		FROM invitation_templates t JOIN event_types e ON e.id = t.type_id`
	args := []any{}
	if typeSlug != "" {
		q += ` WHERE e.slug=$1`
		args = append(args, typeSlug)
	}
	q += ` ORDER BY t.tier, t.name`
	rows, err := s.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Template
	for rows.Next() {
		var t Template
		if err := rows.Scan(&t.ID, &t.Slug, &t.TypeID, &t.Name, &t.Tier, &t.Version, &t.Slots); err != nil {
			return nil, err
		}
		t.Locked = t.Tier == "premium" && !allowPremium
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *Service) GetTemplate(ctx context.Context, id uuid.UUID) (Template, error) {
	var t Template
	err := s.pool.QueryRow(ctx, `SELECT id, slug, type_id, name, tier, version, slots FROM invitation_templates WHERE id=$1`, id).
		Scan(&t.ID, &t.Slug, &t.TypeID, &t.Name, &t.Tier, &t.Version, &t.Slots)
	if errors.Is(err, pgx.ErrNoRows) {
		return Template{}, apierr.NotFound("template")
	}
	return t, err
}
