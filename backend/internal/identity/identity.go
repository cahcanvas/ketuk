package identity

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"ketuk.id/api/internal/apierr"
)

const freePlanID = "00000000-0000-0000-0000-000000000001"

type Service struct {
	pool       *pgxpool.Pool
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func New(pool *pgxpool.Pool, secret string, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{pool: pool, secret: []byte(secret), accessTTL: accessTTL, refreshTTL: refreshTTL}
}

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func (s *Service) Register(ctx context.Context, email, password, name string) (User, TokenPair, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	name = strings.TrimSpace(name)
	if email == "" || name == "" {
		return User{}, TokenPair{}, apierr.BadRequest("invalid_input", "email and name are required")
	}
	if len(password) < 8 {
		return User{}, TokenPair{}, apierr.BadRequest("weak_password", "password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return User{}, TokenPair{}, err
	}

	user := User{ID: uuid.New(), Email: email, Name: name}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return User{}, TokenPair{}, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `INSERT INTO users (id, email, password_hash, name) VALUES ($1,$2,$3,$4)`,
		user.ID, user.Email, string(hash), user.Name)
	if err != nil {
		if isUnique(err) {
			return User{}, TokenPair{}, apierr.Conflict("email_taken", "email already registered")
		}
		return User{}, TokenPair{}, err
	}
	_, err = tx.Exec(ctx, `INSERT INTO subscriptions (id, user_id, plan_id, status, current_period_end)
		VALUES ($1,$2,$3,'active', now() + interval '100 years')`,
		uuid.New(), user.ID, freePlanID)
	if err != nil {
		return User{}, TokenPair{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return User{}, TokenPair{}, err
	}

	tokens, err := s.issue(ctx, user.ID)
	return user, tokens, err
}

func (s *Service) Login(ctx context.Context, email, password string) (User, TokenPair, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	var user User
	var hash string
	err := s.pool.QueryRow(ctx, `SELECT id, email, name, password_hash FROM users WHERE lower(email)=$1`, email).
		Scan(&user.ID, &user.Email, &user.Name, &hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, TokenPair{}, apierr.Unauthorized("invalid credentials")
	}
	if err != nil {
		return User{}, TokenPair{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return User{}, TokenPair{}, apierr.Unauthorized("invalid credentials")
	}
	tokens, err := s.issue(ctx, user.ID)
	return user, tokens, err
}

// Refresh rotates the pair: the presented session is spent, so a refresh token
// that leaks and is replayed after the owner refreshed is rejected.
func (s *Service) Refresh(ctx context.Context, token string) (TokenPair, error) {
	userID, sessionID, err := s.parseRefresh(token)
	if err != nil {
		return TokenPair{}, err
	}
	tag, err := s.pool.Exec(ctx, `
		UPDATE refresh_sessions SET revoked_at=now()
		WHERE id=$1 AND user_id=$2 AND revoked_at IS NULL AND expires_at > now()`, sessionID, userID)
	if err != nil {
		return TokenPair{}, err
	}
	if tag.RowsAffected() == 0 {
		return TokenPair{}, apierr.Unauthorized("refresh token is no longer valid")
	}
	return s.issue(ctx, userID)
}

func (s *Service) Logout(ctx context.Context, token string) error {
	userID, sessionID, err := s.parseRefresh(token)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `
		UPDATE refresh_sessions SET revoked_at=now()
		WHERE id=$1 AND user_id=$2 AND revoked_at IS NULL`, sessionID, userID)
	return err
}

func (s *Service) UpdateProfile(ctx context.Context, id uuid.UUID, name string) (User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return User{}, apierr.BadRequest("invalid_input", "name is required")
	}
	var user User
	err := s.pool.QueryRow(ctx, `
		UPDATE users SET name=$2, updated_at=now() WHERE id=$1
		RETURNING id, email, name`, id, name).Scan(&user.ID, &user.Email, &user.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, apierr.NotFound("user")
	}
	return user, err
}

// ChangePassword revokes every session: a password change is also the way to
// kick out a device whose refresh token was lost.
func (s *Service) ChangePassword(ctx context.Context, id uuid.UUID, current, next string) error {
	if len(next) < 8 {
		return apierr.BadRequest("weak_password", "password must be at least 8 characters")
	}
	var hash string
	err := s.pool.QueryRow(ctx, `SELECT password_hash FROM users WHERE id=$1`, id).Scan(&hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return apierr.NotFound("user")
	}
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(current)) != nil {
		return apierr.Unauthorized("current password is incorrect")
	}
	fresh, err := bcrypt.GenerateFromPassword([]byte(next), 12)
	if err != nil {
		return err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `UPDATE users SET password_hash=$2, updated_at=now() WHERE id=$1`, id, string(fresh)); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		UPDATE refresh_sessions SET revoked_at=now()
		WHERE user_id=$1 AND revoked_at IS NULL`, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Service) ParseAccess(token string) (uuid.UUID, error) {
	claims, err := s.parse(token)
	if err != nil || claims["typ"] != "access" {
		return uuid.Nil, apierr.Unauthorized("invalid access token")
	}
	sub, _ := claims["sub"].(string)
	id, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, apierr.Unauthorized("invalid access token")
	}
	return id, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (User, error) {
	var user User
	err := s.pool.QueryRow(ctx, `SELECT id, email, name FROM users WHERE id=$1`, id).
		Scan(&user.ID, &user.Email, &user.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, apierr.NotFound("user")
	}
	return user, err
}

func (s *Service) issue(ctx context.Context, userID uuid.UUID) (TokenPair, error) {
	access, err := s.sign(userID, "access", s.accessTTL, uuid.Nil)
	if err != nil {
		return TokenPair{}, err
	}
	sessionID := uuid.New()
	refresh, err := s.sign(userID, "refresh", s.refreshTTL, sessionID)
	if err != nil {
		return TokenPair{}, err
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO refresh_sessions (id, user_id, expires_at) VALUES ($1,$2,$3)`,
		sessionID, userID, time.Now().Add(s.refreshTTL))
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int(s.accessTTL.Seconds()),
	}, nil
}

func (s *Service) sign(userID uuid.UUID, typ string, ttl time.Duration, sessionID uuid.UUID) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"typ": typ,
		"iat": now.Unix(),
		"exp": now.Add(ttl).Unix(),
	}
	if sessionID != uuid.Nil {
		claims["jti"] = sessionID.String()
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *Service) parseRefresh(token string) (userID, sessionID uuid.UUID, err error) {
	claims, err := s.parse(token)
	if err != nil || claims["typ"] != "refresh" {
		return uuid.Nil, uuid.Nil, apierr.Unauthorized("invalid refresh token")
	}
	sub, _ := claims["sub"].(string)
	userID, err = uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, uuid.Nil, apierr.Unauthorized("invalid refresh token")
	}
	jti, _ := claims["jti"].(string)
	sessionID, err = uuid.Parse(jti)
	if err != nil {
		return uuid.Nil, uuid.Nil, apierr.Unauthorized("invalid refresh token")
	}
	return userID, sessionID, nil
}

// PurgeExpiredSessions drops rows that can no longer authorise anything.
func (s *Service) PurgeExpiredSessions(ctx context.Context) (int64, error) {
	tag, err := s.pool.Exec(ctx, `DELETE FROM refresh_sessions WHERE expires_at < now()`)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (s *Service) parse(token string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, apierr.Unauthorized("invalid token")
		}
		return s.secret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, apierr.Unauthorized("invalid token")
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, apierr.Unauthorized("invalid token")
	}
	return claims, nil
}

func isUnique(err error) bool {
	var pg *pgconn.PgError
	return errors.As(err, &pg) && pg.Code == "23505"
}
