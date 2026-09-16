package planner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"ketuk.id/api/internal/apierr"
	"ketuk.id/api/internal/billing"
)

type Service struct {
	pool    *pgxpool.Pool
	billing *billing.Service
}

func New(pool *pgxpool.Pool, billing *billing.Service) *Service {
	return &Service{pool: pool, billing: billing}
}

type Planner struct {
	ID         uuid.UUID  `json:"id"`
	OwnerID    uuid.UUID  `json:"owner_id"`
	TypeID     *uuid.UUID `json:"type_id"`
	Title      string     `json:"title"`
	BudgetMode string     `json:"budget_mode"`
	Currency   string     `json:"currency"`
	Status     string     `json:"status"`
	StartsAt   *time.Time `json:"starts_at"`
	EndsAt     *time.Time `json:"ends_at"`
}

type Day struct {
	ID        uuid.UUID `json:"id"`
	Date      time.Time `json:"date"`
	Label     string    `json:"label"`
	SortOrder int       `json:"sort_order"`
}

type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
}

type Item struct {
	ID         uuid.UUID  `json:"id"`
	CategoryID uuid.UUID  `json:"category_id"`
	DayID      *uuid.UUID `json:"day_id"`
	Name       string     `json:"name"`
	Amount     int64      `json:"amount"`
	PaidAt     *time.Time `json:"paid_at"`
	Note       string     `json:"note"`
}

type CheckItem struct {
	ID        uuid.UUID  `json:"id"`
	DayID     *uuid.UUID `json:"day_id"`
	Title     string     `json:"title"`
	Done      bool       `json:"done"`
	DueAt     *time.Time `json:"due_at"`
	SortOrder int        `json:"sort_order"`
}

func (s *Service) List(ctx context.Context, owner uuid.UUID) ([]Planner, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, owner_id, type_id, title, budget_mode, currency, status, starts_at, ends_at
		FROM planners WHERE owner_id=$1 AND status='active' ORDER BY created_at DESC`, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Planner
	for rows.Next() {
		p, err := scanPlanner(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Service) Get(ctx context.Context, owner, id uuid.UUID) (Planner, error) {
	p, err := s.get(ctx, id)
	if err != nil {
		return Planner{}, err
	}
	if p.OwnerID != owner {
		return Planner{}, apierr.NotFound("planner")
	}
	return p, nil
}

func (s *Service) Create(ctx context.Context, owner uuid.UUID, title, mode string, typeID *uuid.UUID, categories json.RawMessage) (Planner, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Planner{}, apierr.BadRequest("invalid_input", "title is required")
	}
	if mode == "" {
		mode = "event"
	}
	if mode != "event" && mode != "daily" {
		return Planner{}, apierr.BadRequest("invalid_input", "budget_mode must be event or daily")
	}
	if err := s.billing.AssertPlannerQuota(ctx, owner); err != nil {
		return Planner{}, err
	}
	if mode == "daily" {
		if err := s.billing.AssertPlannerDaily(ctx, owner); err != nil {
			return Planner{}, err
		}
	}
	p := Planner{
		ID: uuid.New(), OwnerID: owner, TypeID: typeID, Title: title,
		BudgetMode: mode, Currency: "IDR", Status: "active",
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Planner{}, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `INSERT INTO planners (id, owner_id, type_id, title, budget_mode, currency, status) VALUES ($1,$2,$3,$4,$5,'IDR','active')`,
		p.ID, p.OwnerID, p.TypeID, p.Title, p.BudgetMode)
	if err != nil {
		return Planner{}, err
	}
	names := categoryNames(categories)
	if len(names) == 0 && typeID != nil {
		var raw []byte
		if err := tx.QueryRow(ctx, `SELECT planner_categories FROM event_types WHERE id=$1`, *typeID).Scan(&raw); err == nil {
			names = categoryNames(raw)
		}
	}
	for i, name := range names {
		_, err = tx.Exec(ctx, `INSERT INTO budget_categories (id, planner_id, name, sort_order) VALUES ($1,$2,$3,$4)`,
			uuid.New(), p.ID, name, i)
		if err != nil {
			return Planner{}, err
		}
	}
	return p, tx.Commit(ctx)
}

func (s *Service) SetMode(ctx context.Context, owner, id uuid.UUID, mode string) (Planner, error) {
	p, err := s.Get(ctx, owner, id)
	if err != nil {
		return Planner{}, err
	}
	if mode != "event" && mode != "daily" {
		return Planner{}, apierr.BadRequest("invalid_input", "budget_mode must be event or daily")
	}
	if mode == "daily" {
		if err := s.billing.AssertPlannerDaily(ctx, owner); err != nil {
			return Planner{}, err
		}
	}
	if mode == "event" && p.BudgetMode == "daily" {
		_, err = s.pool.Exec(ctx, `
			UPDATE budget_items SET day_id=NULL WHERE category_id IN (SELECT id FROM budget_categories WHERE planner_id=$1)`, id)
		if err != nil {
			return Planner{}, err
		}
	}
	_, err = s.pool.Exec(ctx, `UPDATE planners SET budget_mode=$2, updated_at=now() WHERE id=$1`, id, mode)
	p.BudgetMode = mode
	return p, err
}

// PlannerPatch carries only the fields the host sent. Clearable columns use
// json.RawMessage so an absent key ("leave as is") differs from null ("clear").
type PlannerPatch struct {
	Title      *string         `json:"title"`
	BudgetMode *string         `json:"budget_mode"`
	TypeID     json.RawMessage `json:"type_id"`
	StartsAt   json.RawMessage `json:"starts_at"`
	EndsAt     json.RawMessage `json:"ends_at"`
}

func (s *Service) Patch(ctx context.Context, owner, id uuid.UUID, in PlannerPatch) (Planner, error) {
	p, err := s.Get(ctx, owner, id)
	if err != nil {
		return Planner{}, err
	}
	next, err := applyPlannerPatch(p, in)
	if err != nil {
		return Planner{}, err
	}
	if in.BudgetMode != nil {
		// SetMode owns the daily entitlement and the daily -> event cleanup.
		if _, err := s.SetMode(ctx, owner, id, next.BudgetMode); err != nil {
			return Planner{}, err
		}
	}
	_, err = s.pool.Exec(ctx, `
		UPDATE planners SET title=$3, type_id=$4, starts_at=$5, ends_at=$6, updated_at=now()
		WHERE id=$1 AND owner_id=$2`, id, owner, next.Title, next.TypeID, next.StartsAt, next.EndsAt)
	if isForeignKey(err) {
		return Planner{}, apierr.BadRequest("invalid_input", "type_id is not a known event type")
	}
	if err != nil {
		return Planner{}, err
	}
	return next, nil
}

func applyPlannerPatch(p Planner, in PlannerPatch) (Planner, error) {
	if in.Title != nil {
		title := strings.TrimSpace(*in.Title)
		if title == "" {
			return p, apierr.BadRequest("invalid_input", "title is required")
		}
		p.Title = title
	}
	if in.BudgetMode != nil {
		if *in.BudgetMode != "event" && *in.BudgetMode != "daily" {
			return p, apierr.BadRequest("invalid_input", "budget_mode must be event or daily")
		}
		p.BudgetMode = *in.BudgetMode
	}
	typeID, err := patchUUID(in.TypeID, p.TypeID, "type_id")
	if err != nil {
		return p, err
	}
	starts, err := patchTime(in.StartsAt, p.StartsAt, "starts_at")
	if err != nil {
		return p, err
	}
	ends, err := patchTime(in.EndsAt, p.EndsAt, "ends_at")
	if err != nil {
		return p, err
	}
	if starts != nil && ends != nil && dayOnly(*ends).Before(dayOnly(*starts)) {
		return p, apierr.BadRequest("invalid_input", "ends_at must not be before starts_at")
	}
	p.TypeID = typeID
	p.StartsAt = starts
	p.EndsAt = ends
	return p, nil
}

// Archive frees the quota slot without deleting budget or checklist rows.
// Links stay; remove those through DELETE /v1/links/{id}.
func (s *Service) Archive(ctx context.Context, owner, id uuid.UUID) error {
	p, err := s.Get(ctx, owner, id)
	if err != nil {
		return err
	}
	if p.Status == "archived" {
		return nil
	}
	_, err = s.pool.Exec(ctx, `UPDATE planners SET status='archived', updated_at=now() WHERE id=$1 AND owner_id=$2`, id, owner)
	return err
}

func (s *Service) AddDay(ctx context.Context, owner, plannerID uuid.UUID, date time.Time, label string) (Day, error) {
	if _, err := s.Get(ctx, owner, plannerID); err != nil {
		return Day{}, err
	}
	d := Day{ID: uuid.New(), Date: date, Label: strings.TrimSpace(label)}
	_, err := s.pool.Exec(ctx, `INSERT INTO planner_days (id, planner_id, day_date, label) VALUES ($1,$2,$3,$4)`,
		d.ID, plannerID, d.Date, d.Label)
	if isUnique(err) {
		return Day{}, apierr.Conflict("conflict", "day date already exists")
	}
	return d, err
}

func (s *Service) ListDays(ctx context.Context, plannerID uuid.UUID) ([]Day, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, day_date, label, sort_order FROM planner_days WHERE planner_id=$1 ORDER BY day_date`, plannerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Day
	for rows.Next() {
		var d Day
		if err := rows.Scan(&d.ID, &d.Date, &d.Label, &d.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// DayPatch has no clearable column: day_date is NOT NULL and label uses "".
type DayPatch struct {
	Date      *time.Time `json:"date"`
	Label     *string    `json:"label"`
	SortOrder *int       `json:"sort_order"`
}

func (s *Service) PatchDay(ctx context.Context, owner, plannerID, dayID uuid.UUID, in DayPatch) (Day, error) {
	if _, err := s.Get(ctx, owner, plannerID); err != nil {
		return Day{}, err
	}
	var d Day
	err := s.pool.QueryRow(ctx, `
		SELECT id, day_date, label, sort_order FROM planner_days WHERE id=$1 AND planner_id=$2`, dayID, plannerID).
		Scan(&d.ID, &d.Date, &d.Label, &d.SortOrder)
	if errors.Is(err, pgx.ErrNoRows) {
		return Day{}, apierr.NotFound("day")
	}
	if err != nil {
		return Day{}, err
	}
	d = applyDayPatch(d, in)
	_, err = s.pool.Exec(ctx, `
		UPDATE planner_days SET day_date=$3, label=$4, sort_order=$5 WHERE id=$1 AND planner_id=$2`,
		d.ID, plannerID, d.Date, d.Label, d.SortOrder)
	if isUnique(err) {
		return Day{}, apierr.Conflict("conflict", "day date already exists")
	}
	return d, err
}

func applyDayPatch(d Day, in DayPatch) Day {
	if in.Date != nil {
		d.Date = *in.Date
	}
	if in.Label != nil {
		d.Label = strings.TrimSpace(*in.Label)
	}
	if in.SortOrder != nil {
		d.SortOrder = *in.SortOrder
	}
	return d
}

// DeleteDay leaves budget items and checklist rows in place; both day_id
// columns are ON DELETE SET NULL.
func (s *Service) DeleteDay(ctx context.Context, owner, plannerID, dayID uuid.UUID) error {
	if _, err := s.Get(ctx, owner, plannerID); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `DELETE FROM planner_days WHERE id=$1 AND planner_id=$2`, dayID, plannerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apierr.NotFound("day")
	}
	return nil
}

func (s *Service) AddCategory(ctx context.Context, owner, plannerID uuid.UUID, name string) (Category, error) {
	if _, err := s.Get(ctx, owner, plannerID); err != nil {
		return Category{}, err
	}
	c := Category{ID: uuid.New(), Name: strings.TrimSpace(name)}
	_, err := s.pool.Exec(ctx, `INSERT INTO budget_categories (id, planner_id, name) VALUES ($1,$2,$3)`, c.ID, plannerID, c.Name)
	return c, err
}

func (s *Service) ListCategories(ctx context.Context, plannerID uuid.UUID) ([]Category, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name, sort_order FROM budget_categories WHERE planner_id=$1 ORDER BY sort_order, name`, plannerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

type CategoryPatch struct {
	Name      *string `json:"name"`
	SortOrder *int    `json:"sort_order"`
}

func (s *Service) PatchCategory(ctx context.Context, owner, plannerID, categoryID uuid.UUID, in CategoryPatch) (Category, error) {
	if _, err := s.Get(ctx, owner, plannerID); err != nil {
		return Category{}, err
	}
	var c Category
	err := s.pool.QueryRow(ctx, `
		SELECT id, name, sort_order FROM budget_categories WHERE id=$1 AND planner_id=$2`, categoryID, plannerID).
		Scan(&c.ID, &c.Name, &c.SortOrder)
	if errors.Is(err, pgx.ErrNoRows) {
		return Category{}, apierr.NotFound("category")
	}
	if err != nil {
		return Category{}, err
	}
	c, err = applyCategoryPatch(c, in)
	if err != nil {
		return Category{}, err
	}
	_, err = s.pool.Exec(ctx, `
		UPDATE budget_categories SET name=$3, sort_order=$4 WHERE id=$1 AND planner_id=$2`,
		c.ID, plannerID, c.Name, c.SortOrder)
	return c, err
}

func applyCategoryPatch(c Category, in CategoryPatch) (Category, error) {
	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return c, apierr.BadRequest("invalid_input", "name is required")
		}
		c.Name = name
	}
	if in.SortOrder != nil {
		c.SortOrder = *in.SortOrder
	}
	return c, nil
}

// DeleteCategory also drops every budget item inside it (FK ON DELETE CASCADE).
func (s *Service) DeleteCategory(ctx context.Context, owner, plannerID, categoryID uuid.UUID) error {
	if _, err := s.Get(ctx, owner, plannerID); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `DELETE FROM budget_categories WHERE id=$1 AND planner_id=$2`, categoryID, plannerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apierr.NotFound("category")
	}
	return nil
}

func (s *Service) AddItem(ctx context.Context, owner, plannerID, categoryID uuid.UUID, name string, amount int64, dayID *uuid.UUID, note string) (Item, error) {
	p, err := s.Get(ctx, owner, plannerID)
	if err != nil {
		return Item{}, err
	}
	if amount < 0 {
		return Item{}, apierr.BadRequest("invalid_input", "amount must not be negative")
	}
	if err := s.assertCategory(ctx, plannerID, categoryID); err != nil {
		return Item{}, err
	}
	if dayID != nil {
		if err := s.assertItemDay(ctx, p.BudgetMode, plannerID, *dayID); err != nil {
			return Item{}, err
		}
	}
	item := Item{ID: uuid.New(), CategoryID: categoryID, DayID: dayID, Name: strings.TrimSpace(name), Amount: amount, Note: note}
	_, err = s.pool.Exec(ctx, `INSERT INTO budget_items (id, category_id, day_id, name, amount, note) VALUES ($1,$2,$3,$4,$5,$6)`,
		item.ID, item.CategoryID, item.DayID, item.Name, item.Amount, item.Note)
	return item, err
}

// ItemPatch clears day_id or paid_at with null; an absent key changes nothing.
type ItemPatch struct {
	CategoryID *uuid.UUID      `json:"category_id"`
	DayID      json.RawMessage `json:"day_id"`
	Name       *string         `json:"name"`
	Amount     *int64          `json:"amount"`
	PaidAt     json.RawMessage `json:"paid_at"`
	Note       *string         `json:"note"`
}

func (s *Service) PatchItem(ctx context.Context, owner, plannerID, itemID uuid.UUID, in ItemPatch) (Item, error) {
	p, err := s.Get(ctx, owner, plannerID)
	if err != nil {
		return Item{}, err
	}
	var item Item
	err = s.pool.QueryRow(ctx, `
		SELECT i.id, i.category_id, i.day_id, i.name, i.amount, i.paid_at, i.note
		FROM budget_items i
		JOIN budget_categories c ON c.id = i.category_id
		WHERE i.id=$1 AND c.planner_id=$2`, itemID, plannerID).
		Scan(&item.ID, &item.CategoryID, &item.DayID, &item.Name, &item.Amount, &item.PaidAt, &item.Note)
	if errors.Is(err, pgx.ErrNoRows) {
		return Item{}, apierr.NotFound("budget_item")
	}
	if err != nil {
		return Item{}, err
	}
	next, err := applyItemPatch(item, in)
	if err != nil {
		return Item{}, err
	}
	if in.CategoryID != nil {
		if err := s.assertCategory(ctx, plannerID, next.CategoryID); err != nil {
			return Item{}, err
		}
	}
	if fieldSent(in.DayID) && next.DayID != nil {
		if err := s.assertItemDay(ctx, p.BudgetMode, plannerID, *next.DayID); err != nil {
			return Item{}, err
		}
	}
	_, err = s.pool.Exec(ctx, `
		UPDATE budget_items SET category_id=$2, day_id=$3, name=$4, amount=$5, paid_at=$6, note=$7 WHERE id=$1`,
		next.ID, next.CategoryID, next.DayID, next.Name, next.Amount, next.PaidAt, next.Note)
	return next, err
}

func applyItemPatch(item Item, in ItemPatch) (Item, error) {
	if in.CategoryID != nil {
		item.CategoryID = *in.CategoryID
	}
	dayID, err := patchUUID(in.DayID, item.DayID, "day_id")
	if err != nil {
		return item, err
	}
	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return item, apierr.BadRequest("invalid_input", "name is required")
		}
		item.Name = name
	}
	if in.Amount != nil {
		if *in.Amount < 0 {
			return item, apierr.BadRequest("invalid_input", "amount must not be negative")
		}
		item.Amount = *in.Amount
	}
	paidAt, err := patchTime(in.PaidAt, item.PaidAt, "paid_at")
	if err != nil {
		return item, err
	}
	if in.Note != nil {
		item.Note = strings.TrimSpace(*in.Note)
	}
	item.DayID = dayID
	item.PaidAt = paidAt
	return item, nil
}

func (s *Service) DeleteItem(ctx context.Context, owner, plannerID, itemID uuid.UUID) error {
	if _, err := s.Get(ctx, owner, plannerID); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `
		DELETE FROM budget_items
		WHERE id=$1 AND category_id IN (SELECT id FROM budget_categories WHERE planner_id=$2)`, itemID, plannerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apierr.NotFound("budget_item")
	}
	return nil
}

func (s *Service) Budget(ctx context.Context, owner, plannerID uuid.UUID) (map[string]any, error) {
	p, err := s.Get(ctx, owner, plannerID)
	if err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT c.id, c.name, i.id, i.day_id, i.name, i.amount, i.paid_at, i.note
		FROM budget_categories c
		LEFT JOIN budget_items i ON i.category_id = c.id
		WHERE c.planner_id=$1
		ORDER BY c.sort_order, c.name`, plannerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type catView struct {
		Category
		Items []Item `json:"items"`
		Total int64  `json:"total"`
		Paid  int64  `json:"paid"`
	}
	order := []uuid.UUID{}
	byID := map[uuid.UUID]*catView{}
	var grand, paidGrand int64
	for rows.Next() {
		var cid uuid.UUID
		var cname string
		var iid *uuid.UUID
		var dayID *uuid.UUID
		var iname *string
		var amount *int64
		var paidAt *time.Time
		var note *string
		if err := rows.Scan(&cid, &cname, &iid, &dayID, &iname, &amount, &paidAt, &note); err != nil {
			return nil, err
		}
		cv, ok := byID[cid]
		if !ok {
			cv = &catView{Category: Category{ID: cid, Name: cname}}
			byID[cid] = cv
			order = append(order, cid)
		}
		if iid == nil {
			continue
		}
		item := Item{ID: *iid, CategoryID: cid, DayID: dayID, Name: *iname, Amount: *amount, PaidAt: paidAt}
		if note != nil {
			item.Note = *note
		}
		cv.Items = append(cv.Items, item)
		cv.Total += item.Amount
		grand += item.Amount
		if item.PaidAt != nil {
			cv.Paid += item.Amount
			paidGrand += item.Amount
		}
	}
	cats := make([]catView, 0, len(order))
	for _, id := range order {
		cats = append(cats, *byID[id])
	}
	return map[string]any{
		"budget_mode": p.BudgetMode,
		"currency":    p.Currency,
		"total":       grand,
		"paid":        paidGrand,
		"categories":  cats,
	}, rows.Err()
}

// AddCheck accepts a day even in event mode: a task may hang on a day without
// being a budget row.
func (s *Service) AddCheck(ctx context.Context, owner, plannerID uuid.UUID, title string, dayID *uuid.UUID, dueAt *time.Time, sortOrder int) (CheckItem, error) {
	if _, err := s.Get(ctx, owner, plannerID); err != nil {
		return CheckItem{}, err
	}
	if dayID != nil {
		if err := s.assertDay(ctx, plannerID, *dayID); err != nil {
			return CheckItem{}, err
		}
	}
	c := CheckItem{ID: uuid.New(), Title: strings.TrimSpace(title), DayID: dayID, DueAt: dueAt, SortOrder: sortOrder}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO checklist_items (id, planner_id, day_id, title, due_at, sort_order) VALUES ($1,$2,$3,$4,$5,$6)`,
		c.ID, plannerID, c.DayID, c.Title, c.DueAt, c.SortOrder)
	return c, err
}

func (s *Service) ListChecks(ctx context.Context, owner, plannerID uuid.UUID) ([]CheckItem, error) {
	if _, err := s.Get(ctx, owner, plannerID); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, day_id, title, done, due_at, sort_order FROM checklist_items WHERE planner_id=$1 ORDER BY sort_order, title`, plannerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CheckItem
	for rows.Next() {
		var c CheckItem
		if err := rows.Scan(&c.ID, &c.DayID, &c.Title, &c.Done, &c.DueAt, &c.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// CheckPatch clears day_id or due_at with null; an absent key changes nothing.
type CheckPatch struct {
	Title     *string         `json:"title"`
	DayID     json.RawMessage `json:"day_id"`
	Done      *bool           `json:"done"`
	DueAt     json.RawMessage `json:"due_at"`
	SortOrder *int            `json:"sort_order"`
}

func (s *Service) PatchCheck(ctx context.Context, owner, plannerID, itemID uuid.UUID, in CheckPatch) (CheckItem, error) {
	if _, err := s.Get(ctx, owner, plannerID); err != nil {
		return CheckItem{}, err
	}
	var c CheckItem
	err := s.pool.QueryRow(ctx, `
		SELECT id, day_id, title, done, due_at, sort_order FROM checklist_items WHERE id=$1 AND planner_id=$2`, itemID, plannerID).
		Scan(&c.ID, &c.DayID, &c.Title, &c.Done, &c.DueAt, &c.SortOrder)
	if errors.Is(err, pgx.ErrNoRows) {
		return CheckItem{}, apierr.NotFound("checklist")
	}
	if err != nil {
		return CheckItem{}, err
	}
	next, err := applyCheckPatch(c, in)
	if err != nil {
		return CheckItem{}, err
	}
	if fieldSent(in.DayID) && next.DayID != nil {
		if err := s.assertDay(ctx, plannerID, *next.DayID); err != nil {
			return CheckItem{}, err
		}
	}
	_, err = s.pool.Exec(ctx, `
		UPDATE checklist_items SET day_id=$3, title=$4, done=$5, due_at=$6, sort_order=$7 WHERE id=$1 AND planner_id=$2`,
		next.ID, plannerID, next.DayID, next.Title, next.Done, next.DueAt, next.SortOrder)
	return next, err
}

func applyCheckPatch(c CheckItem, in CheckPatch) (CheckItem, error) {
	if in.Title != nil {
		title := strings.TrimSpace(*in.Title)
		if title == "" {
			return c, apierr.BadRequest("invalid_input", "title is required")
		}
		c.Title = title
	}
	dayID, err := patchUUID(in.DayID, c.DayID, "day_id")
	if err != nil {
		return c, err
	}
	dueAt, err := patchTime(in.DueAt, c.DueAt, "due_at")
	if err != nil {
		return c, err
	}
	if in.Done != nil {
		c.Done = *in.Done
	}
	if in.SortOrder != nil {
		c.SortOrder = *in.SortOrder
	}
	c.DayID = dayID
	c.DueAt = dueAt
	return c, nil
}

func (s *Service) DeleteCheck(ctx context.Context, owner, plannerID, itemID uuid.UUID) error {
	if _, err := s.Get(ctx, owner, plannerID); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `DELETE FROM checklist_items WHERE id=$1 AND planner_id=$2`, itemID, plannerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apierr.NotFound("checklist")
	}
	return nil
}

func (s *Service) assertCategory(ctx context.Context, plannerID, categoryID uuid.UUID) error {
	var found bool
	err := s.pool.QueryRow(ctx, `SELECT true FROM budget_categories WHERE id=$1 AND planner_id=$2`, categoryID, plannerID).Scan(&found)
	if errors.Is(err, pgx.ErrNoRows) {
		return apierr.NotFound("category")
	}
	return err
}

func (s *Service) assertDay(ctx context.Context, plannerID, dayID uuid.UUID) error {
	var found bool
	err := s.pool.QueryRow(ctx, `SELECT true FROM planner_days WHERE id=$1 AND planner_id=$2`, dayID, plannerID).Scan(&found)
	if errors.Is(err, pgx.ErrNoRows) {
		return apierr.NotFound("day")
	}
	return err
}

// assertItemDay guards budget rows only: pinning money to a day is what the
// daily mode buys.
func (s *Service) assertItemDay(ctx context.Context, mode string, plannerID, dayID uuid.UUID) error {
	if mode != "daily" {
		return apierr.BadRequest("invalid_input", "day_id only allowed when budget_mode is daily")
	}
	return s.assertDay(ctx, plannerID, dayID)
}

func (s *Service) get(ctx context.Context, id uuid.UUID) (Planner, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT id, owner_id, type_id, title, budget_mode, currency, status, starts_at, ends_at FROM planners WHERE id=$1`, id)
	p, err := scanPlanner(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Planner{}, apierr.NotFound("planner")
	}
	return p, err
}

type scanner interface{ Scan(dest ...any) error }

func scanPlanner(row scanner) (Planner, error) {
	var p Planner
	err := row.Scan(&p.ID, &p.OwnerID, &p.TypeID, &p.Title, &p.BudgetMode, &p.Currency, &p.Status, &p.StartsAt, &p.EndsAt)
	return p, err
}

// fieldSent reports whether the host included the key at all.
func fieldSent(raw json.RawMessage) bool {
	return len(bytes.TrimSpace(raw)) > 0
}

// fieldNull reports whether the host sent the key as JSON null, i.e. "clear it".
func fieldNull(raw json.RawMessage) bool {
	return string(bytes.TrimSpace(raw)) == "null"
}

func patchUUID(raw json.RawMessage, cur *uuid.UUID, field string) (*uuid.UUID, error) {
	if !fieldSent(raw) {
		return cur, nil
	}
	if fieldNull(raw) {
		return nil, nil
	}
	var id uuid.UUID
	if err := json.Unmarshal(raw, &id); err != nil {
		return nil, apierr.BadRequest("invalid_input", field+" must be a uuid or null")
	}
	return &id, nil
}

func patchTime(raw json.RawMessage, cur *time.Time, field string) (*time.Time, error) {
	if !fieldSent(raw) {
		return cur, nil
	}
	if fieldNull(raw) {
		return nil, nil
	}
	var t time.Time
	if err := json.Unmarshal(raw, &t); err != nil {
		return nil, apierr.BadRequest("invalid_input", field+" must be an RFC3339 timestamp or null")
	}
	return &t, nil
}

// dayOnly drops the clock so planner dates compare as calendar days, matching
// the DATE columns behind them.
func dayOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func isUnique(err error) bool {
	var pg *pgconn.PgError
	return errors.As(err, &pg) && pg.Code == "23505"
}

func isForeignKey(err error) bool {
	var pg *pgconn.PgError
	return errors.As(err, &pg) && pg.Code == "23503"
}

func categoryNames(raw json.RawMessage) []string {
	if len(raw) == 0 {
		return nil
	}
	var names []string
	if json.Unmarshal(raw, &names) == nil {
		return names
	}
	return nil
}
