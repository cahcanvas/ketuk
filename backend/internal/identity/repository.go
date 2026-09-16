package identity

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Sentinel errors a Repository implementation must return so the use cases in
// service.go never need to know about storage-specific error types
// (pgconn.PgError, pgx.ErrNoRows, ...).
var (
	ErrEmailTaken      = errors.New("email already registered")
	ErrUsernameTaken   = errors.New("username already registered")
	ErrUserNotFound    = errors.New("user not found")
	ErrSessionInactive = errors.New("refresh session is not active")
)

// Repository is the persistence port the identity use cases depend on.
// postgres_repository.go is the only file allowed to know this is Postgres —
// swapping storage, or unit-testing a use case with a fake, only needs a new
// implementation of this interface.
type Repository interface {
	CreateUser(ctx context.Context, id uuid.UUID, email, username, passwordHash string) error
	UserByIdentifier(ctx context.Context, identifier string) (User, string, error)
	UserByID(ctx context.Context, id uuid.UUID) (User, error)
	UpdateUserName(ctx context.Context, id uuid.UUID, name string) (User, error)
	PasswordHash(ctx context.Context, id uuid.UUID) (string, error)
	ChangePassword(ctx context.Context, id uuid.UUID, newHash string) error

	CreateRefreshSession(ctx context.Context, sessionID, userID uuid.UUID, expiresAt time.Time) error
	ConsumeRefreshSession(ctx context.Context, sessionID, userID uuid.UUID) error
	RevokeRefreshSession(ctx context.Context, sessionID, userID uuid.UUID) error
	DeleteExpiredSessions(ctx context.Context) (int64, error)
}
