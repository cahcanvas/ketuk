package planner

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"

	"ketuk.id/api/internal/apierr"
)

func date(t *testing.T, s string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("bad fixture %q: %v", s, err)
	}
	return parsed
}

func code(t *testing.T, err error) string {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error")
	}
	ae, ok := err.(*apierr.Error)
	if !ok {
		t.Fatalf("expected apierr, got %T", err)
	}
	return ae.Message
}

func TestApplyPlannerPatchNoOp(t *testing.T) {
	typeID := uuid.New()
	starts := date(t, "2026-11-01T00:00:00+07:00")
	current := Planner{
		ID: uuid.New(), OwnerID: uuid.New(), TypeID: &typeID, Title: "Nikahan Sep",
		BudgetMode: "daily", Currency: "IDR", Status: "active", StartsAt: &starts,
	}

	got, err := applyPlannerPatch(current, PlannerPatch{})
	if err != nil {
		t.Fatalf("empty patch failed: %v", err)
	}
	if got != current {
		t.Fatalf("empty patch changed planner: %+v", got)
	}
}

func TestApplyPlannerPatchFields(t *testing.T) {
	typeID := uuid.New()
	starts := date(t, "2026-11-01T00:00:00+07:00")
	ends := date(t, "2026-11-03T00:00:00+07:00")
	current := Planner{
		Title: "Lama", BudgetMode: "event", Currency: "IDR", Status: "active",
		TypeID: &typeID, StartsAt: &starts, EndsAt: &ends,
	}

	title := "  Nikahan Sep  "
	mode := "daily"
	got, err := applyPlannerPatch(current, PlannerPatch{
		Title:      &title,
		BudgetMode: &mode,
		TypeID:     json.RawMessage(`null`),
		EndsAt:     json.RawMessage(`null`),
	})
	if err != nil {
		t.Fatalf("patch failed: %v", err)
	}
	if got.Title != "Nikahan Sep" {
		t.Fatalf("got title %q", got.Title)
	}
	if got.BudgetMode != "daily" {
		t.Fatalf("got mode %q", got.BudgetMode)
	}
	if got.TypeID != nil {
		t.Fatalf("expected type_id cleared, got %v", got.TypeID)
	}
	if got.EndsAt != nil {
		t.Fatalf("expected ends_at cleared, got %v", got.EndsAt)
	}
	if got.StartsAt == nil || !got.StartsAt.Equal(starts) {
		t.Fatalf("absent starts_at changed: %v", got.StartsAt)
	}
	if got.Status != "active" || got.Currency != "IDR" {
		t.Fatalf("patch touched status or currency: %+v", got)
	}
}

func TestApplyPlannerPatchRejects(t *testing.T) {
	starts := date(t, "2026-11-05T00:00:00+07:00")
	current := Planner{Title: "Lama", BudgetMode: "event", StartsAt: &starts}

	blank := "   "
	_, err := applyPlannerPatch(current, PlannerPatch{Title: &blank})
	if msg := code(t, err); msg != "title is required" {
		t.Fatalf("got %q", msg)
	}

	weekly := "weekly"
	_, err = applyPlannerPatch(current, PlannerPatch{BudgetMode: &weekly})
	if msg := code(t, err); msg != "budget_mode must be event or daily" {
		t.Fatalf("got %q", msg)
	}

	_, err = applyPlannerPatch(current, PlannerPatch{EndsAt: json.RawMessage(`"2026-11-04T00:00:00+07:00"`)})
	if msg := code(t, err); msg != "ends_at must not be before starts_at" {
		t.Fatalf("got %q", msg)
	}

	sameDay, err := applyPlannerPatch(current, PlannerPatch{EndsAt: json.RawMessage(`"2026-11-05T23:30:00+07:00"`)})
	if err != nil {
		t.Fatalf("same calendar day rejected: %v", err)
	}
	if sameDay.EndsAt == nil {
		t.Fatal("expected ends_at set")
	}

	_, err = applyPlannerPatch(current, PlannerPatch{TypeID: json.RawMessage(`"nope"`)})
	if msg := code(t, err); msg != "type_id must be a uuid or null" {
		t.Fatalf("got %q", msg)
	}
}

func TestApplyItemPatchThreeState(t *testing.T) {
	dayID := uuid.New()
	paid := date(t, "2026-10-02T09:00:00+07:00")
	current := Item{
		ID: uuid.New(), CategoryID: uuid.New(), DayID: &dayID,
		Name: "Katering", Amount: 25_000_000, PaidAt: &paid, Note: "DP 50%",
	}

	kept, err := applyItemPatch(current, ItemPatch{})
	if err != nil {
		t.Fatalf("empty patch failed: %v", err)
	}
	if kept != current {
		t.Fatalf("empty patch changed item: %+v", kept)
	}

	unpaid, err := applyItemPatch(current, ItemPatch{PaidAt: json.RawMessage(`null`)})
	if err != nil {
		t.Fatalf("clearing paid_at failed: %v", err)
	}
	if unpaid.PaidAt != nil {
		t.Fatalf("expected paid_at cleared, got %v", unpaid.PaidAt)
	}
	if unpaid.DayID == nil || *unpaid.DayID != dayID {
		t.Fatalf("absent day_id changed: %v", unpaid.DayID)
	}

	loose, err := applyItemPatch(current, ItemPatch{DayID: json.RawMessage(`null`)})
	if err != nil {
		t.Fatalf("clearing day_id failed: %v", err)
	}
	if loose.DayID != nil {
		t.Fatalf("expected day_id cleared, got %v", loose.DayID)
	}
	if loose.PaidAt == nil || !loose.PaidAt.Equal(paid) {
		t.Fatalf("absent paid_at changed: %v", loose.PaidAt)
	}
}

func TestApplyItemPatchRejects(t *testing.T) {
	current := Item{Name: "Katering", Amount: 100}

	blank := ""
	_, err := applyItemPatch(current, ItemPatch{Name: &blank})
	if msg := code(t, err); msg != "name is required" {
		t.Fatalf("got %q", msg)
	}

	negative := int64(-1)
	_, err = applyItemPatch(current, ItemPatch{Amount: &negative})
	if msg := code(t, err); msg != "amount must not be negative" {
		t.Fatalf("got %q", msg)
	}

	zero := int64(0)
	if _, err := applyItemPatch(current, ItemPatch{Amount: &zero}); err != nil {
		t.Fatalf("zero amount rejected: %v", err)
	}
}

func TestApplyCheckPatch(t *testing.T) {
	dayID := uuid.New()
	due := date(t, "2026-10-20T17:00:00+07:00")
	current := CheckItem{ID: uuid.New(), DayID: &dayID, Title: "Fitting jas", Done: false, DueAt: &due, SortOrder: 2}

	kept, err := applyCheckPatch(current, CheckPatch{})
	if err != nil {
		t.Fatalf("empty patch failed: %v", err)
	}
	if kept != current {
		t.Fatalf("empty patch changed item: %+v", kept)
	}

	done := true
	order := 5
	got, err := applyCheckPatch(current, CheckPatch{Done: &done, SortOrder: &order, DueAt: json.RawMessage(`null`)})
	if err != nil {
		t.Fatalf("patch failed: %v", err)
	}
	if !got.Done || got.SortOrder != 5 || got.DueAt != nil {
		t.Fatalf("unexpected item: %+v", got)
	}
	if got.DayID == nil || *got.DayID != dayID {
		t.Fatalf("absent day_id changed: %v", got.DayID)
	}

	blank := " "
	_, err = applyCheckPatch(current, CheckPatch{Title: &blank})
	if msg := code(t, err); msg != "title is required" {
		t.Fatalf("got %q", msg)
	}
}

func TestApplyCategoryPatch(t *testing.T) {
	current := Category{ID: uuid.New(), Name: "Katering", SortOrder: 1}

	kept, err := applyCategoryPatch(current, CategoryPatch{})
	if err != nil {
		t.Fatalf("empty patch failed: %v", err)
	}
	if kept != current {
		t.Fatalf("empty patch changed category: %+v", kept)
	}

	name := "  Dekorasi  "
	got, err := applyCategoryPatch(current, CategoryPatch{Name: &name})
	if err != nil {
		t.Fatalf("patch failed: %v", err)
	}
	if got.Name != "Dekorasi" || got.SortOrder != 1 {
		t.Fatalf("unexpected category: %+v", got)
	}

	blank := ""
	_, err = applyCategoryPatch(current, CategoryPatch{Name: &blank})
	if msg := code(t, err); msg != "name is required" {
		t.Fatalf("got %q", msg)
	}
}

func TestApplyDayPatch(t *testing.T) {
	current := Day{ID: uuid.New(), Date: date(t, "2026-11-01T00:00:00+07:00"), Label: "Akad", SortOrder: 0}

	if got := applyDayPatch(current, DayPatch{}); got != current {
		t.Fatalf("empty patch changed day: %+v", got)
	}

	blank := ""
	if got := applyDayPatch(current, DayPatch{Label: &blank}); got.Label != "" {
		t.Fatalf("expected label cleared, got %q", got.Label)
	}

	next := date(t, "2026-11-02T00:00:00+07:00")
	order := 3
	got := applyDayPatch(current, DayPatch{Date: &next, SortOrder: &order})
	if !got.Date.Equal(next) || got.SortOrder != 3 || got.Label != "Akad" {
		t.Fatalf("unexpected day: %+v", got)
	}
}

// A budget row may only hang on a day in daily mode; the guard rejects before
// touching the database, so no pool is needed here.
func TestAssertItemDayEventMode(t *testing.T) {
	err := (&Service{}).assertItemDay(context.Background(), "event", uuid.New(), uuid.New())
	if msg := code(t, err); msg != "day_id only allowed when budget_mode is daily" {
		t.Fatalf("got %q", msg)
	}
}

// Locks the decoder assumption the three-state fields rely on: an absent key
// leaves the RawMessage empty, while an explicit null arrives as "null".
func TestThreeStateDecoding(t *testing.T) {
	var absent ItemPatch
	if err := json.Unmarshal([]byte(`{"name":"Katering"}`), &absent); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if fieldSent(absent.PaidAt) || fieldSent(absent.DayID) {
		t.Fatalf("absent keys look sent: %q %q", absent.PaidAt, absent.DayID)
	}

	var cleared ItemPatch
	if err := json.Unmarshal([]byte(`{"paid_at":null}`), &cleared); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if !fieldSent(cleared.PaidAt) || !fieldNull(cleared.PaidAt) {
		t.Fatalf("null paid_at decoded as %q", cleared.PaidAt)
	}

	id := uuid.New()
	got, err := patchUUID(json.RawMessage(`"`+id.String()+`"`), nil, "day_id")
	if err != nil {
		t.Fatalf("patchUUID failed: %v", err)
	}
	if got == nil || *got != id {
		t.Fatalf("got %v want %v", got, id)
	}
}
