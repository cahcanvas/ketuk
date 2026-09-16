package identity

import "github.com/google/uuid"

// User and TokenPair are domain entities: no transport concerns (JSON tags,
// HTTP status codes) belong here. internal/httpserver owns the wire format
// and maps these to its own response DTOs.
type User struct {
	ID       uuid.UUID
	Email    string
	Username string
	Name     *string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}
