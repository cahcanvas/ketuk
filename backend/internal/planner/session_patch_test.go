package planner

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func session(t *testing.T) Session {
	t.Helper()
	ends := Clock{Hour: 11, Minute: 30}
	return Session{
		ID: uuid.New(), DayID: uuid.New(), Title: "Akad",
		StartsAt: Clock{Hour: 9}, EndsAt: &ends, PicName: "Bu Rina", SortOrder: 1,
	}
}

func TestApplySessionPatchNoOp(t *testing.T) {
	current := session(t)

	got, err := applySessionPatch(current, SessionPatch{})
	if err != nil {
		t.Fatalf("empty patch failed: %v", err)
	}
	if got.Title != current.Title || got.StartsAt != current.StartsAt || got.PicName != current.PicName || got.SortOrder != current.SortOrder {
		t.Fatalf("empty patch changed session: %+v", got)
	}
	if got.EndsAt == nil || *got.EndsAt != *current.EndsAt {
		t.Fatalf("empty patch changed ends_at: %v", got.EndsAt)
	}
}

func TestApplySessionPatchFields(t *testing.T) {
	current := session(t)

	title := "  Resepsi  "
	starts := "13:00"
	pic := "  "
	order := 4
	got, err := applySessionPatch(current, SessionPatch{
		Title:     &title,
		StartsAt:  &starts,
		EndsAt:    json.RawMessage(`"16:30:00"`),
		PicName:   &pic,
		SortOrder: &order,
	})
	if err != nil {
		t.Fatalf("patch failed: %v", err)
	}
	if got.Title != "Resepsi" {
		t.Fatalf("got title %q", got.Title)
	}
	if got.StartsAt != (Clock{Hour: 13}) {
		t.Fatalf("got starts_at %v", got.StartsAt)
	}
	if got.EndsAt == nil || *got.EndsAt != (Clock{Hour: 16, Minute: 30}) {
		t.Fatalf("got ends_at %v", got.EndsAt)
	}
	if got.PicName != "" {
		t.Fatalf("expected pic_name cleared, got %q", got.PicName)
	}
	if got.SortOrder != 4 {
		t.Fatalf("got sort_order %d", got.SortOrder)
	}
}

// ends_at is three-state: absent keeps the current clock, null drops it.
func TestApplySessionPatchEndsAtThreeState(t *testing.T) {
	current := session(t)

	kept, err := applySessionPatch(current, SessionPatch{})
	if err != nil {
		t.Fatalf("absent ends_at failed: %v", err)
	}
	if kept.EndsAt == nil {
		t.Fatal("absent ends_at cleared the clock")
	}

	open, err := applySessionPatch(current, SessionPatch{EndsAt: json.RawMessage(`null`)})
	if err != nil {
		t.Fatalf("clearing ends_at failed: %v", err)
	}
	if open.EndsAt != nil {
		t.Fatalf("expected ends_at cleared, got %v", open.EndsAt)
	}
}

func TestApplySessionPatchRejects(t *testing.T) {
	current := session(t)

	blank := " "
	_, err := applySessionPatch(current, SessionPatch{Title: &blank})
	if msg := code(t, err); msg != "title is required" {
		t.Fatalf("got %q", msg)
	}

	_, err = applySessionPatch(current, SessionPatch{EndsAt: json.RawMessage(`"08:00"`)})
	if msg := code(t, err); msg != "ends_at must be on or after starts_at" {
		t.Fatalf("got %q", msg)
	}

	same, err := applySessionPatch(current, SessionPatch{EndsAt: json.RawMessage(`"09:00"`)})
	if err != nil {
		t.Fatalf("ends_at equal to starts_at rejected: %v", err)
	}
	if same.EndsAt == nil {
		t.Fatal("expected ends_at set")
	}

	noon := "noon"
	_, err = applySessionPatch(current, SessionPatch{StartsAt: &noon})
	if msg := code(t, err); msg != `starts_at must be a time like "14:00" or "14:00:00"` {
		t.Fatalf("got %q", msg)
	}

	_, err = applySessionPatch(current, SessionPatch{EndsAt: json.RawMessage(`14`)})
	if msg := code(t, err); msg != `ends_at must be a time like "14:00" or null` {
		t.Fatalf("got %q", msg)
	}
}

// Sessions ride on a day row, so JSON is a clock string, never RFC3339.
func TestClockJSON(t *testing.T) {
	for _, in := range []string{"14:00", "14:00:00"} {
		got, err := parseClock(in, "starts_at")
		if err != nil {
			t.Fatalf("parse %q failed: %v", in, err)
		}
		if got != (Clock{Hour: 14}) {
			t.Fatalf("parse %q gave %v", in, got)
		}
	}

	if _, err := parseClock("24:00", "starts_at"); err == nil {
		t.Fatal("expected 24:00 to be rejected")
	}

	raw, err := json.Marshal(Clock{Hour: 7, Minute: 5, Second: 9})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(raw) != `"07:05:09"` {
		t.Fatalf("got %s", raw)
	}
}

// Round-trips the TIME bridge so a stored session reads back as the same clock.
func TestClockPgRoundTrip(t *testing.T) {
	want := Clock{Hour: 21, Minute: 45, Second: 30}
	got := clockFromPg(want.pg())
	if got == nil || *got != want {
		t.Fatalf("got %v want %v", got, want)
	}
	if clockFromPg(clockParam(nil)) != nil {
		t.Fatal("expected NULL time to read back as nil")
	}
}
