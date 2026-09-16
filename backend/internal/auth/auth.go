package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"ketuk.id/api/internal/apierr"
)

type ctxKey int

const userIDKey ctxKey = 1

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(b), err
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func Issue(secret string, userID uuid.UUID, accessTTL, refreshTTL time.Duration) (Tokens, error) {
	now := time.Now()
	access, err := sign(secret, userID, "access", now.Add(accessTTL))
	if err != nil {
		return Tokens{}, err
	}
	refresh, err := sign(secret, userID, "refresh", now.Add(refreshTTL))
	if err != nil {
		return Tokens{}, err
	}
	return Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int64(accessTTL.Seconds()),
	}, nil
}

func sign(secret string, userID uuid.UUID, typ string, exp time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"typ": typ,
		"exp": exp.Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func Parse(secret, token, wantTyp string) (uuid.UUID, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		return uuid.Nil, apierr.Unauthorized("invalid token")
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, apierr.Unauthorized("invalid token")
	}
	if claims["typ"] != wantTyp {
		return uuid.Nil, apierr.Unauthorized("wrong token type")
	}
	sub, _ := claims["sub"].(string)
	id, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, apierr.Unauthorized("invalid token")
	}
	return id, nil
}

func Middleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, `{"error":{"code":"unauthorized","message":"missing bearer token"}}`, http.StatusUnauthorized)
				return
			}
			id, err := Parse(secret, strings.TrimPrefix(header, "Bearer "), "access")
			if err != nil {
				http.Error(w, `{"error":{"code":"unauthorized","message":"invalid token"}}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserID(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(userIDKey).(uuid.UUID)
	return id
}
