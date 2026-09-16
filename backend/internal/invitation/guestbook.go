package invitation

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"ketuk.id/api/internal/apierr"
)

const (
	maxGuestbookName    = 80
	maxGuestbookMessage = 2000
)

type GuestbookEntry struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Service) ListGuestbook(ctx context.Context, owner, invitationID uuid.UUID) ([]GuestbookEntry, error) {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return nil, err
	}
	return s.listGuestbook(ctx, invitationID)
}

func (s *Service) PatchGuestbook(ctx context.Context, owner, invitationID, entryID uuid.UUID, name, message *string) (GuestbookEntry, error) {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return GuestbookEntry{}, err
	}
	var e GuestbookEntry
	err := s.pool.QueryRow(ctx, `
		SELECT id, name, message, created_at FROM invitation_guestbook
		WHERE id=$1 AND invitation_id=$2`, entryID, invitationID).
		Scan(&e.ID, &e.Name, &e.Message, &e.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return GuestbookEntry{}, apierr.NotFound("guestbook")
	}
	if err != nil {
		return GuestbookEntry{}, err
	}
	e, err = applyGuestbookPatch(e, name, message)
	if err != nil {
		return GuestbookEntry{}, err
	}
	_, err = s.pool.Exec(ctx, `UPDATE invitation_guestbook SET name=$3, message=$4 WHERE id=$1 AND invitation_id=$2`,
		e.ID, invitationID, e.Name, e.Message)
	return e, err
}

// applyGuestbookPatch keeps created_at untouched; nil fields stay as stored.
func applyGuestbookPatch(entry GuestbookEntry, name, message *string) (GuestbookEntry, error) {
	if name != nil {
		entry.Name = *name
	}
	if message != nil {
		entry.Message = *message
	}
	cleanName, cleanMessage, err := validateGuestbook(entry.Name, entry.Message)
	if err != nil {
		return GuestbookEntry{}, err
	}
	entry.Name = cleanName
	entry.Message = cleanMessage
	return entry, nil
}

func (s *Service) DeleteGuestbook(ctx context.Context, owner, invitationID, entryID uuid.UUID) error {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `DELETE FROM invitation_guestbook WHERE id=$1 AND invitation_id=$2`, entryID, invitationID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apierr.NotFound("guestbook")
	}
	return nil
}

func (s *Service) AddGuestbookPublic(ctx context.Context, slug, name, message string) (GuestbookEntry, error) {
	inv, err := s.publishedBySlug(ctx, slug)
	if err != nil {
		return GuestbookEntry{}, err
	}
	return s.insertGuestbook(ctx, inv.ID, nil, name, message)
}

func (s *Service) AddGuestbookByToken(ctx context.Context, token, name, message string) (GuestbookEntry, error) {
	var invitationID, guestID uuid.UUID
	var guestName string
	err := s.pool.QueryRow(ctx, `
		SELECT g.id, g.name, i.id
		FROM invitation_guests g JOIN invitations i ON i.id = g.invitation_id
		WHERE g.rsvp_token=$1 AND i.status='published'`, token).
		Scan(&guestID, &guestName, &invitationID)
	if errors.Is(err, pgx.ErrNoRows) {
		return GuestbookEntry{}, apierr.NotFound("invitation")
	}
	if err != nil {
		return GuestbookEntry{}, err
	}
	if strings.TrimSpace(name) == "" {
		name = guestName
	}
	return s.insertGuestbook(ctx, invitationID, &guestID, name, message)
}

func validateGuestbook(name, message string) (string, string, error) {
	name = strings.TrimSpace(name)
	message = strings.TrimSpace(message)
	if name == "" {
		return "", "", apierr.BadRequest("invalid_input", "name is required")
	}
	if message == "" {
		return "", "", apierr.BadRequest("invalid_input", "message is required")
	}
	if utf8.RuneCountInString(name) > maxGuestbookName {
		return "", "", apierr.BadRequest("invalid_input", "name is too long")
	}
	if utf8.RuneCountInString(message) > maxGuestbookMessage {
		return "", "", apierr.BadRequest("invalid_input", "message is too long")
	}
	return name, message, nil
}

func (s *Service) insertGuestbook(ctx context.Context, invitationID uuid.UUID, guestID *uuid.UUID, name, message string) (GuestbookEntry, error) {
	name, message, err := validateGuestbook(name, message)
	if err != nil {
		return GuestbookEntry{}, err
	}
	in := GuestbookEntry{ID: uuid.New(), Name: name, Message: message, CreatedAt: time.Now()}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO invitation_guestbook (id, invitation_id, guest_id, name, message, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		in.ID, invitationID, guestID, in.Name, in.Message, in.CreatedAt)
	return in, err
}

func (s *Service) listGuestbook(ctx context.Context, invitationID uuid.UUID) ([]GuestbookEntry, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, message, created_at FROM invitation_guestbook
		WHERE invitation_id=$1 ORDER BY created_at DESC`, invitationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []GuestbookEntry{}
	for rows.Next() {
		var e GuestbookEntry
		if err := rows.Scan(&e.ID, &e.Name, &e.Message, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
