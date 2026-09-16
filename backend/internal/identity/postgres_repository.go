package identity

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const freePlanID = "00000000-0000-0000-0000-000000000001"

// PostgresRepository is the only piece of the identity domain that knows it
// talks to Postgres via pgx; it implements Repository.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, id uuid.UUID, email, username, passwordHash string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `INSERT INTO users (id, email, password_hash, username) VALUES ($1,$2,$3,$4)`,
		id, email, passwordHash, username)
	if err != nil {
		var pg *pgconn.PgError
		if errors.As(err, &pg) && pg.Code == "23505" {
			if pg.ConstraintName == "users_username_lower_idx" {
				return ErrUsernameTaken
			}
			return ErrEmailTaken
		}
		return err
	}
	// Register always starts a host on the free plan.
	_, err = tx.Exec(ctx, `INSERT INTO subscriptions (id, user_id, plan_id, status, current_period_end)
		VALUES ($1,$2,$3,'active', now() + interval '100 years')`,
		uuid.New(), id, freePlanID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *PostgresRepository) UserByIdentifier(ctx context.Context, identifier string) (User, string, error) {
	var user User
	var hash string
	err := r.pool.QueryRow(ctx, `
		SELECT id, email, username, name, password_hash FROM users
		WHERE lower(email)=$1 OR lower(username)=$1
		LIMIT 1`, identifier).
		Scan(&user.ID, &user.Email, &user.Username, &user.Name, &hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, "", ErrUserNotFound
	}
	return user, hash, err
}

func (r *PostgresRepository) UserByID(ctx context.Context, id uuid.UUID) (User, error) {
	var user User
	err := r.pool.QueryRow(ctx, `SELECT id, email, username, name FROM users WHERE id=$1`, id).
		Scan(&user.ID, &user.Email, &user.Username, &user.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	return user, err
}

func (r *PostgresRepository) UpdateUserName(ctx context.Context, id uuid.UUID, name string) (User, error) {
	var user User
	err := r.pool.QueryRow(ctx, `
		UPDATE users SET name=$2, updated_at=now() WHERE id=$1
		RETURNING id, email, username, name`, id, name).
		Scan(&user.ID, &user.Email, &user.Username, &user.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	return user, err
}

func (r *PostgresRepository) PasswordHash(ctx context.Context, id uuid.UUID) (string, error) {
	var hash string
	err := r.pool.QueryRow(ctx, `SELECT password_hash FROM users WHERE id=$1`, id).Scan(&hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrUserNotFound
	}
	return hash, err
}

func (r *PostgresRepository) ChangePassword(ctx context.Context, id uuid.UUID, newHash string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `UPDATE users SET password_hash=$2, updated_at=now() WHERE id=$1`, id, newHash); err != nil {
		return err
	}
	// A password change is also the way to kick out a device whose refresh
	// token was lost, so every session dies with it.
	if _, err := tx.Exec(ctx, `
		UPDATE refresh_sessions SET revoked_at=now()
		WHERE user_id=$1 AND revoked_at IS NULL`, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *PostgresRepository) CreateRefreshSession(ctx context.Context, sessionID, userID uuid.UUID, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO refresh_sessions (id, user_id, expires_at) VALUES ($1,$2,$3)`,
		sessionID, userID, expiresAt)
	return err
}

// ConsumeRefreshSession rotates the pair: the presented session is spent, so
// a refresh token that leaks and is replayed after the owner refreshed is
// rejected via ErrSessionInactive.
func (r *PostgresRepository) ConsumeRefreshSession(ctx context.Context, sessionID, userID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE refresh_sessions SET revoked_at=now()
		WHERE id=$1 AND user_id=$2 AND revoked_at IS NULL AND expires_at > now()`, sessionID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrSessionInactive
	}
	return nil
}

// RevokeRefreshSession is used by logout: unlike ConsumeRefreshSession it is
// lenient about an already-gone session, matching prior logout behaviour.
func (r *PostgresRepository) RevokeRefreshSession(ctx context.Context, sessionID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE refresh_sessions SET revoked_at=now()
		WHERE id=$1 AND user_id=$2 AND revoked_at IS NULL`, sessionID, userID)
	return err
}

// DeleteExpiredSessions drops rows that can no longer authorise anything.
func (r *PostgresRepository) DeleteExpiredSessions(ctx context.Context) (int64, error) {
	tag, err := r.pool.Exec(ctx, `DELETE FROM refresh_sessions WHERE expires_at < now()`)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}
