package httpserver

import (
	"github.com/google/uuid"

	"ketuk.id/api/internal/identity"
)

// userResponse and tokenPairResponse are this layer's wire format for
// identity.User / identity.TokenPair. The domain entities carry no JSON tags
// on purpose, so the mapping — and therefore the response shape — lives here.
type userResponse struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Name     *string   `json:"name"`
}

func newUserResponse(u identity.User) userResponse {
	return userResponse{ID: u.ID, Email: u.Email, Username: u.Username, Name: u.Name}
}

type tokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func newTokenPairResponse(t identity.TokenPair) tokenPairResponse {
	return tokenPairResponse{AccessToken: t.AccessToken, RefreshToken: t.RefreshToken, ExpiresIn: t.ExpiresIn}
}
