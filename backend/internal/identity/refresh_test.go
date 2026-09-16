package identity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func testService() *Service {
	return New(nil, "test-secret", time.Minute, time.Hour)
}

func TestParseRefreshRejectsTokenWithoutSession(t *testing.T) {
	s := testService()
	// Tokens issued before refresh sessions existed carry no jti.
	token, err := s.sign(uuid.New(), "refresh", time.Hour, uuid.Nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.parseRefresh(token); err == nil {
		t.Fatal("a refresh token without jti must be rejected")
	}
}

func TestParseRefreshRejectsAccessToken(t *testing.T) {
	s := testService()
	token, err := s.sign(uuid.New(), "access", time.Minute, uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.parseRefresh(token); err == nil {
		t.Fatal("an access token must not be usable as a refresh token")
	}
}

func TestParseRefreshAcceptsIssuedPair(t *testing.T) {
	s := testService()
	userID, sessionID := uuid.New(), uuid.New()
	token, err := s.sign(userID, "refresh", time.Hour, sessionID)
	if err != nil {
		t.Fatal(err)
	}
	gotUser, gotSession, err := s.parseRefresh(token)
	if err != nil {
		t.Fatal(err)
	}
	if gotUser != userID || gotSession != sessionID {
		t.Fatalf("got %s/%s, want %s/%s", gotUser, gotSession, userID, sessionID)
	}
}
