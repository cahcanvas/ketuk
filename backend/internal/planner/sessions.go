package planner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"ketuk.id/api/internal/apierr"
)

// Clock is a wall clock without a date: a session hangs on a planner day, so
// the day already carries the calendar part. JSON accepts "HH:MM" or
// "HH:MM:SS" and always encodes "HH:MM:SS"; Postgres stores it as TIME.
type Clock struct {
	Hour   int
	Minute int
	Second int
}

type Session struct {
	ID        uuid.UUID `json:"id"`
	DayID     uuid.UUID `json:"day_id"`
	Title     string    `json:"title"`
	StartsAt  Clock     `json:"starts_at"`
	EndsAt    *Clock    `json:"ends_at"`
	PicName   string    `json:"pic_name"`
	SortOrder int       `json:"sort_order"`
}

func (s *Service) ListSessions(ctx context.Context, owner, plannerID, dayID uuid.UUID) ([]Session, error) {
	if _, err := s.Get(ctx, owner, plannerID); err != nil {
		return nil, err
	}
	if err := s.assertDay(ctx, plannerID, dayID); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, day_id, title, starts_at, ends_at, pic_name, sort_order
		FROM planner_sessions WHERE day_id=$1 ORDER BY starts_at, sort_order`, dayID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Session
	for rows.Next() {
		sess, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, sess)
	}
	return out, rows.Err()
}

// AddSession takes the clocks as strings so a malformed time reads as
// invalid_input from the domain instead of invalid_json from the decoder.
func (s *Service) AddSession(ctx context.Context, owner, plannerID, dayID uuid.UUID, title, startsAt string, endsAt *string, picName string, sortOrder int) (Session, error) {
	if _, err := s.Get(ctx, owner, plannerID); err != nil {
		return Session{}, err
	}
	if err := s.assertDay(ctx, plannerID, dayID); err != nil {
		return Session{}, err
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return Session{}, apierr.BadRequest("invalid_input", "title is required")
	}
	if strings.TrimSpace(startsAt) == "" {
		return Session{}, apierr.BadRequest("invalid_input", "starts_at is required")
	}
	starts, err := parseClock(startsAt, "starts_at")
	if err != nil {
		return Session{}, err
	}
	sess := Session{
		ID: uuid.New(), DayID: dayID, Title: title, StartsAt: starts,
		PicName: strings.TrimSpace(picName), SortOrder: sortOrder,
	}
	if endsAt != nil && strings.TrimSpace(*endsAt) != "" {
		ends, err := parseClock(*endsAt, "ends_at")
		if err != nil {
			return Session{}, err
		}
		sess.EndsAt = &ends
	}
	if err := assertSessionOrder(sess); err != nil {
		return Session{}, err
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO planner_sessions (id, day_id, title, starts_at, ends_at, pic_name, sort_order)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		sess.ID, sess.DayID, sess.Title, sess.StartsAt.pg(), clockParam(sess.EndsAt), sess.PicName, sess.SortOrder)
	return sess, err
}

// SessionPatch clears ends_at with null; an absent key changes nothing.
type SessionPatch struct {
	Title     *string         `json:"title"`
	StartsAt  *string         `json:"starts_at"`
	EndsAt    json.RawMessage `json:"ends_at"`
	PicName   *string         `json:"pic_name"`
	SortOrder *int            `json:"sort_order"`
}

func (s *Service) PatchSession(ctx context.Context, owner, plannerID, dayID, sessionID uuid.UUID, in SessionPatch) (Session, error) {
	if _, err := s.Get(ctx, owner, plannerID); err != nil {
		return Session{}, err
	}
	if err := s.assertDay(ctx, plannerID, dayID); err != nil {
		return Session{}, err
	}
	row := s.pool.QueryRow(ctx, `
		SELECT id, day_id, title, starts_at, ends_at, pic_name, sort_order
		FROM planner_sessions WHERE id=$1 AND day_id=$2`, sessionID, dayID)
	sess, err := scanSession(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Session{}, apierr.NotFound("session")
	}
	if err != nil {
		return Session{}, err
	}
	next, err := applySessionPatch(sess, in)
	if err != nil {
		return Session{}, err
	}
	_, err = s.pool.Exec(ctx, `
		UPDATE planner_sessions SET title=$3, starts_at=$4, ends_at=$5, pic_name=$6, sort_order=$7
		WHERE id=$1 AND day_id=$2`,
		next.ID, dayID, next.Title, next.StartsAt.pg(), clockParam(next.EndsAt), next.PicName, next.SortOrder)
	return next, err
}

func applySessionPatch(sess Session, in SessionPatch) (Session, error) {
	if in.Title != nil {
		title := strings.TrimSpace(*in.Title)
		if title == "" {
			return sess, apierr.BadRequest("invalid_input", "title is required")
		}
		sess.Title = title
	}
	if in.StartsAt != nil {
		starts, err := parseClock(*in.StartsAt, "starts_at")
		if err != nil {
			return sess, err
		}
		sess.StartsAt = starts
	}
	ends, err := patchClock(in.EndsAt, sess.EndsAt, "ends_at")
	if err != nil {
		return sess, err
	}
	sess.EndsAt = ends
	if in.PicName != nil {
		sess.PicName = strings.TrimSpace(*in.PicName)
	}
	if in.SortOrder != nil {
		sess.SortOrder = *in.SortOrder
	}
	if err := assertSessionOrder(sess); err != nil {
		return sess, err
	}
	return sess, nil
}

func (s *Service) DeleteSession(ctx context.Context, owner, plannerID, dayID, sessionID uuid.UUID) error {
	if _, err := s.Get(ctx, owner, plannerID); err != nil {
		return err
	}
	if err := s.assertDay(ctx, plannerID, dayID); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `DELETE FROM planner_sessions WHERE id=$1 AND day_id=$2`, sessionID, dayID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apierr.NotFound("session")
	}
	return nil
}

func assertSessionOrder(sess Session) error {
	if sess.EndsAt != nil && sess.EndsAt.seconds() < sess.StartsAt.seconds() {
		return apierr.BadRequest("invalid_input", "ends_at must be on or after starts_at")
	}
	return nil
}

func scanSession(row scanner) (Session, error) {
	var sess Session
	var starts, ends pgtype.Time
	if err := row.Scan(&sess.ID, &sess.DayID, &sess.Title, &starts, &ends, &sess.PicName, &sess.SortOrder); err != nil {
		return Session{}, err
	}
	if c := clockFromPg(starts); c != nil {
		sess.StartsAt = *c
	}
	sess.EndsAt = clockFromPg(ends)
	return sess, nil
}

func parseClock(raw, field string) (Clock, error) {
	raw = strings.TrimSpace(raw)
	for _, layout := range []string{"15:04:05", "15:04"} {
		if t, err := time.Parse(layout, raw); err == nil {
			return Clock{Hour: t.Hour(), Minute: t.Minute(), Second: t.Second()}, nil
		}
	}
	return Clock{}, apierr.BadRequest("invalid_input", field+` must be a time like "14:00" or "14:00:00"`)
}

func patchClock(raw json.RawMessage, cur *Clock, field string) (*Clock, error) {
	if !fieldSent(raw) {
		return cur, nil
	}
	if fieldNull(raw) {
		return nil, nil
	}
	var text string
	if err := json.Unmarshal(raw, &text); err != nil {
		return nil, apierr.BadRequest("invalid_input", field+` must be a time like "14:00" or null`)
	}
	c, err := parseClock(text, field)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (c Clock) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", c.Hour, c.Minute, c.Second)
}

func (c Clock) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c Clock) seconds() int {
	return c.Hour*3600 + c.Minute*60 + c.Second
}

func (c Clock) pg() pgtype.Time {
	return pgtype.Time{Microseconds: int64(c.seconds()) * 1_000_000, Valid: true}
}

func clockParam(c *Clock) pgtype.Time {
	if c == nil {
		return pgtype.Time{}
	}
	return c.pg()
}

func clockFromPg(t pgtype.Time) *Clock {
	if !t.Valid {
		return nil
	}
	secs := int(t.Microseconds / 1_000_000)
	return &Clock{Hour: secs / 3600, Minute: (secs % 3600) / 60, Second: secs % 60}
}
