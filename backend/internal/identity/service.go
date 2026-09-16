package identity

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"ketuk.id/api/internal/apierr"
)

// Service holds the identity use cases. It depends only on the Repository
// port, never on pgx/Postgres directly — persistence is an implementation
// detail injected at the composition root (cmd/api/main.go).
type Service struct {
	repo       Repository
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func New(repo Repository, secret string, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{repo: repo, secret: []byte(secret), accessTTL: accessTTL, refreshTTL: refreshTTL}
}

func (s *Service) Register(ctx context.Context, email, password, username string) (User, TokenPair, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	username = strings.TrimSpace(username)
	if email == "" || username == "" {
		return User{}, TokenPair{}, apierr.BadRequest("invalid_input", "email and username are required")
	}
	if len(password) < 8 {
		return User{}, TokenPair{}, apierr.BadRequest("weak_password", "password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return User{}, TokenPair{}, err
	}

	user := User{ID: uuid.New(), Email: email, Username: username}
	if err := s.repo.CreateUser(ctx, user.ID, user.Email, user.Username, string(hash)); err != nil {
		switch {
		case errors.Is(err, ErrUsernameTaken):
			return User{}, TokenPair{}, apierr.Conflict("username_taken", "username already registered")
		case errors.Is(err, ErrEmailTaken):
			return User{}, TokenPair{}, apierr.Conflict("email_taken", "email already registered")
		default:
			return User{}, TokenPair{}, err
		}
	}

	tokens, err := s.issue(ctx, user.ID)
	return user, tokens, err
}

// Login accepts either the account email or username as identifier, since the
// register payload no longer guarantees the caller knows which one they set.
func (s *Service) Login(ctx context.Context, identifier, password string) (User, TokenPair, error) {
	identifier = strings.ToLower(strings.TrimSpace(identifier))
	user, hash, err := s.repo.UserByIdentifier(ctx, identifier)
	if errors.Is(err, ErrUserNotFound) {
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
	if err := s.repo.ConsumeRefreshSession(ctx, sessionID, userID); err != nil {
		if errors.Is(err, ErrSessionInactive) {
			return TokenPair{}, apierr.Unauthorized("refresh token is no longer valid")
		}
		return TokenPair{}, err
	}
	return s.issue(ctx, userID)
}

func (s *Service) Logout(ctx context.Context, token string) error {
	userID, sessionID, err := s.parseRefresh(token)
	if err != nil {
		return err
	}
	return s.repo.RevokeRefreshSession(ctx, sessionID, userID)
}

func (s *Service) UpdateProfile(ctx context.Context, id uuid.UUID, name string) (User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return User{}, apierr.BadRequest("invalid_input", "name is required")
	}
	user, err := s.repo.UpdateUserName(ctx, id, name)
	if errors.Is(err, ErrUserNotFound) {
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
	hash, err := s.repo.PasswordHash(ctx, id)
	if errors.Is(err, ErrUserNotFound) {
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
	return s.repo.ChangePassword(ctx, id, string(fresh))
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
	user, err := s.repo.UserByID(ctx, id)
	if errors.Is(err, ErrUserNotFound) {
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
	if err := s.repo.CreateRefreshSession(ctx, sessionID, userID, time.Now().Add(s.refreshTTL)); err != nil {
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
	return s.repo.DeleteExpiredSessions(ctx)
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
