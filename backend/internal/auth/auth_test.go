package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestIssueAndParse(t *testing.T) {
	secret := "test-secret"
	id := uuid.MustParse("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	tok, err := Issue(secret, id, time.Minute, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Parse(secret, tok.AccessToken, "access")
	if err != nil || got != id {
		t.Fatalf("access parse: %v %s", err, got)
	}
	if _, err := Parse(secret, tok.AccessToken, "refresh"); err == nil {
		t.Fatal("expected wrong type")
	}
	got, err = Parse(secret, tok.RefreshToken, "refresh")
	if err != nil || got != id {
		t.Fatalf("refresh parse: %v %s", err, got)
	}
}

func TestPassword(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatal(err)
	}
	if err := CheckPassword(hash, "secret123"); err != nil {
		t.Fatal(err)
	}
	if err := CheckPassword(hash, "nope"); err == nil {
		t.Fatal("expected mismatch")
	}
}
