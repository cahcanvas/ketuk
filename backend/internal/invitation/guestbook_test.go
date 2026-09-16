package invitation

import (
	"strings"
	"testing"
	"time"
)

func TestValidateGuestbook(t *testing.T) {
	if _, _, err := validateGuestbook("  ", "halo"); err == nil {
		t.Fatal("expected name required")
	}
	if _, _, err := validateGuestbook("Budi", "  "); err == nil {
		t.Fatal("expected message required")
	}
	name, msg, err := validateGuestbook("  Budi  ", "  selamat  ")
	if err != nil {
		t.Fatal(err)
	}
	if name != "Budi" || msg != "selamat" {
		t.Fatalf("got %q %q", name, msg)
	}
}

func TestApplyGuestbookPatch(t *testing.T) {
	created := time.Now().Add(-time.Hour)
	entry := GuestbookEntry{Name: "Budi", Message: "selamat", CreatedAt: created}

	same, err := applyGuestbookPatch(entry, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if same.Name != "Budi" || same.Message != "selamat" || !same.CreatedAt.Equal(created) {
		t.Fatalf("empty patch changed entry: %+v", same)
	}

	msg := "  turut berbahagia  "
	patched, err := applyGuestbookPatch(entry, nil, &msg)
	if err != nil {
		t.Fatal(err)
	}
	if patched.Message != "turut berbahagia" || patched.Name != "Budi" {
		t.Fatalf("got %+v", patched)
	}

	blank := "   "
	if _, err := applyGuestbookPatch(entry, &blank, nil); err == nil {
		t.Fatal("expected blank name rejected")
	}

	long := strings.Repeat("a", maxGuestbookMessage+1)
	if _, err := applyGuestbookPatch(entry, nil, &long); err == nil {
		t.Fatal("expected too long message rejected")
	}
}
