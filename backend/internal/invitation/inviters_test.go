package invitation

import (
	"testing"

	"github.com/google/uuid"
)

func TestApplyInviterPatch(t *testing.T) {
	current := Inviter{ID: uuid.New(), Name: "Bpk. Slamet", Role: "Ayah mempelai", SortOrder: 2}

	same, err := applyInviterPatch(current, InviterPatch{})
	if err != nil {
		t.Fatal(err)
	}
	if same != current {
		t.Fatalf("empty patch changed inviter: %+v", same)
	}

	name, role, order := "  Ibu Sri  ", "  Ibu mempelai  ", 1
	patched, err := applyInviterPatch(current, InviterPatch{Name: &name, Role: &role, SortOrder: &order})
	if err != nil {
		t.Fatal(err)
	}
	if patched.Name != "Ibu Sri" || patched.Role != "Ibu mempelai" || patched.SortOrder != 1 {
		t.Fatalf("got %+v", patched)
	}
	if patched.ID != current.ID {
		t.Fatal("id must not change")
	}

	// An empty role clears it; an empty name is still rejected like AddInviter.
	blank := ""
	cleared, err := applyInviterPatch(current, InviterPatch{Role: &blank})
	if err != nil {
		t.Fatal(err)
	}
	if cleared.Role != "" || cleared.Name != current.Name {
		t.Fatalf("got %+v", cleared)
	}
	for _, bad := range []string{"", "   "} {
		if _, err := applyInviterPatch(current, InviterPatch{Name: &bad}); err == nil {
			t.Fatalf("expected name %q rejected", bad)
		}
	}
}
