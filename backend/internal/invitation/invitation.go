package invitation

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ketuk.id/api/internal/apierr"
	"ketuk.id/api/internal/billing"
	"ketuk.id/api/internal/catalog"
	"ketuk.id/api/internal/notify"
	"ketuk.id/api/internal/storage"
)

type Service struct {
	pool    *pgxpool.Pool
	billing *billing.Service
	catalog *catalog.Service
	notify  notify.Notifier
	store   storage.Store
}

func New(pool *pgxpool.Pool, billing *billing.Service, catalog *catalog.Service, n notify.Notifier, store storage.Store) *Service {
	if n == nil {
		n = notify.Noop{}
	}
	return &Service{pool: pool, billing: billing, catalog: catalog, notify: n, store: store}
}

type Invitation struct {
	ID              uuid.UUID       `json:"id"`
	OwnerID         uuid.UUID       `json:"owner_id"`
	TypeID          uuid.UUID       `json:"type_id"`
	TemplateID      uuid.UUID       `json:"template_id"`
	TemplateVersion int             `json:"template_version"`
	Slug            *string         `json:"slug"`
	Status          string          `json:"status"`
	Title           string          `json:"title"`
	Content         json.RawMessage `json:"content"`
	StartsAt        *time.Time      `json:"starts_at"`
	EndsAt          *time.Time      `json:"ends_at"`
}

type Location struct {
	ID        uuid.UUID  `json:"id"`
	Label     string     `json:"label"`
	VenueName string     `json:"venue_name"`
	Address   string     `json:"address"`
	MapsURL   string     `json:"maps_url"`
	Lat       *float64   `json:"lat"`
	Lng       *float64   `json:"lng"`
	StartsAt  *time.Time `json:"starts_at"`
	SortOrder int        `json:"sort_order"`
}

type Guest struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Group     string     `json:"group"`
	RSVPToken string     `json:"rsvp_token"`
	Status    string     `json:"status"`
	PlusOnes  int        `json:"plus_ones"`
	Message   string     `json:"message"`
	RSVPAt    *time.Time `json:"rsvp_at"`
}

func (s *Service) List(ctx context.Context, owner uuid.UUID) ([]Invitation, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, owner_id, type_id, template_id, template_version, slug, status, title, content, starts_at, ends_at
		FROM invitations WHERE owner_id=$1 AND status <> 'archived' ORDER BY created_at DESC`, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Invitation
	for rows.Next() {
		inv, err := scanInvitation(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, inv)
	}
	return out, rows.Err()
}

func (s *Service) Get(ctx context.Context, owner, id uuid.UUID) (Invitation, error) {
	inv, err := s.get(ctx, id)
	if err != nil {
		return Invitation{}, err
	}
	if inv.OwnerID != owner {
		return Invitation{}, apierr.NotFound("invitation")
	}
	return inv, nil
}

func (s *Service) Create(ctx context.Context, owner uuid.UUID, title string, templateID uuid.UUID, content json.RawMessage) (Invitation, error) {
	if strings.TrimSpace(title) == "" {
		return Invitation{}, apierr.BadRequest("invalid_input", "title is required")
	}
	if err := s.billing.AssertInvitationQuota(ctx, owner); err != nil {
		return Invitation{}, err
	}
	tpl, err := s.catalog.GetTemplate(ctx, templateID)
	if err != nil {
		return Invitation{}, err
	}
	if err := s.billing.AssertPremiumTemplate(ctx, owner, tpl.Tier); err != nil {
		return Invitation{}, err
	}
	if len(content) == 0 {
		content = json.RawMessage(`{}`)
	}
	inv := Invitation{
		ID:              uuid.New(),
		OwnerID:         owner,
		TypeID:          tpl.TypeID,
		TemplateID:      tpl.ID,
		TemplateVersion: tpl.Version,
		Status:          "draft",
		Title:           strings.TrimSpace(title),
		Content:         content,
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO invitations (id, owner_id, type_id, template_id, template_version, status, title, content)
		VALUES ($1,$2,$3,$4,$5,'draft',$6,$7)`,
		inv.ID, inv.OwnerID, inv.TypeID, inv.TemplateID, inv.TemplateVersion, inv.Title, []byte(inv.Content))
	return inv, err
}

func (s *Service) Patch(ctx context.Context, owner, id uuid.UUID, title *string, content json.RawMessage, starts, ends *time.Time) (Invitation, error) {
	inv, err := s.Get(ctx, owner, id)
	if err != nil {
		return Invitation{}, err
	}
	if title != nil {
		inv.Title = strings.TrimSpace(*title)
	}
	if content != nil {
		inv.Content = content
	}
	if starts != nil {
		inv.StartsAt = starts
	}
	if ends != nil {
		inv.EndsAt = ends
	}
	_, err = s.pool.Exec(ctx, `UPDATE invitations SET title=$2, content=$3, starts_at=$4, ends_at=$5, updated_at=now() WHERE id=$1`,
		inv.ID, inv.Title, []byte(inv.Content), inv.StartsAt, inv.EndsAt)
	return inv, err
}

func (s *Service) SetTemplate(ctx context.Context, owner, id, templateID uuid.UUID) (Invitation, error) {
	inv, err := s.Get(ctx, owner, id)
	if err != nil {
		return Invitation{}, err
	}
	if inv.Status == "published" {
		return Invitation{}, apierr.Conflict("immutable", "cannot change template after publish")
	}
	tpl, err := s.catalog.GetTemplate(ctx, templateID)
	if err != nil {
		return Invitation{}, err
	}
	if err := s.billing.AssertPremiumTemplate(ctx, owner, tpl.Tier); err != nil {
		return Invitation{}, err
	}
	_, err = s.pool.Exec(ctx, `UPDATE invitations SET template_id=$2, template_version=$3, type_id=$4, updated_at=now() WHERE id=$1`,
		id, tpl.ID, tpl.Version, tpl.TypeID)
	if err != nil {
		return Invitation{}, err
	}
	inv.TemplateID = tpl.ID
	inv.TemplateVersion = tpl.Version
	inv.TypeID = tpl.TypeID
	return inv, nil
}

func (s *Service) Publish(ctx context.Context, owner, id uuid.UUID) (Invitation, error) {
	inv, err := s.Get(ctx, owner, id)
	if err != nil {
		return Invitation{}, err
	}
	if inv.Status == "published" {
		return inv, nil
	}
	slug := uniqueSlug(inv.Title, inv.ID)
	_, err = s.pool.Exec(ctx, `UPDATE invitations SET status='published', slug=$2, updated_at=now() WHERE id=$1`, id, slug)
	if err != nil {
		return Invitation{}, err
	}
	inv.Status = "published"
	inv.Slug = &slug
	return inv, nil
}

func (s *Service) Archive(ctx context.Context, owner, id uuid.UUID) error {
	inv, err := s.Get(ctx, owner, id)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `UPDATE invitations SET status='archived', updated_at=now() WHERE id=$1`, inv.ID)
	return err
}

func (s *Service) AddLocation(ctx context.Context, owner, invitationID uuid.UUID, loc Location) (Location, error) {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return Location{}, err
	}
	loc.ID = uuid.New()
	_, err := s.pool.Exec(ctx, `
		INSERT INTO invitation_locations (id, invitation_id, label, venue_name, address, maps_url, lat, lng, starts_at, sort_order)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		loc.ID, invitationID, loc.Label, loc.VenueName, loc.Address, loc.MapsURL, loc.Lat, loc.Lng, loc.StartsAt, loc.SortOrder)
	return loc, err
}

func (s *Service) ListLocations(ctx context.Context, invitationID uuid.UUID) ([]Location, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, label, venue_name, address, maps_url, lat, lng, starts_at, sort_order
		FROM invitation_locations WHERE invitation_id=$1 ORDER BY sort_order, venue_name`, invitationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Location
	for rows.Next() {
		var loc Location
		if err := rows.Scan(&loc.ID, &loc.Label, &loc.VenueName, &loc.Address, &loc.MapsURL, &loc.Lat, &loc.Lng, &loc.StartsAt, &loc.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, loc)
	}
	return out, rows.Err()
}

// LocationPatch carries only the fields the host sent; nil means "leave as is".
type LocationPatch struct {
	Label     *string    `json:"label"`
	VenueName *string    `json:"venue_name"`
	Address   *string    `json:"address"`
	MapsURL   *string    `json:"maps_url"`
	Lat       *float64   `json:"lat"`
	Lng       *float64   `json:"lng"`
	StartsAt  *time.Time `json:"starts_at"`
	SortOrder *int       `json:"sort_order"`
}

func (s *Service) PatchLocation(ctx context.Context, owner, invitationID, locationID uuid.UUID, in LocationPatch) (Location, error) {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return Location{}, err
	}
	var loc Location
	err := s.pool.QueryRow(ctx, `
		SELECT id, label, venue_name, address, maps_url, lat, lng, starts_at, sort_order
		FROM invitation_locations WHERE id=$1 AND invitation_id=$2`, locationID, invitationID).
		Scan(&loc.ID, &loc.Label, &loc.VenueName, &loc.Address, &loc.MapsURL, &loc.Lat, &loc.Lng, &loc.StartsAt, &loc.SortOrder)
	if errors.Is(err, pgx.ErrNoRows) {
		return Location{}, apierr.NotFound("location")
	}
	if err != nil {
		return Location{}, err
	}
	if in.Label != nil {
		loc.Label = strings.TrimSpace(*in.Label)
	}
	if in.VenueName != nil {
		loc.VenueName = strings.TrimSpace(*in.VenueName)
	}
	if in.Address != nil {
		loc.Address = strings.TrimSpace(*in.Address)
	}
	if in.MapsURL != nil {
		loc.MapsURL = strings.TrimSpace(*in.MapsURL)
	}
	if in.Lat != nil {
		loc.Lat = in.Lat
	}
	if in.Lng != nil {
		loc.Lng = in.Lng
	}
	if in.StartsAt != nil {
		loc.StartsAt = in.StartsAt
	}
	if in.SortOrder != nil {
		loc.SortOrder = *in.SortOrder
	}
	_, err = s.pool.Exec(ctx, `
		UPDATE invitation_locations
		SET label=$3, venue_name=$4, address=$5, maps_url=$6, lat=$7, lng=$8, starts_at=$9, sort_order=$10
		WHERE id=$1 AND invitation_id=$2`,
		loc.ID, invitationID, loc.Label, loc.VenueName, loc.Address, loc.MapsURL, loc.Lat, loc.Lng, loc.StartsAt, loc.SortOrder)
	return loc, err
}

func (s *Service) DeleteLocation(ctx context.Context, owner, invitationID, locationID uuid.UUID) error {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `DELETE FROM invitation_locations WHERE id=$1 AND invitation_id=$2`, locationID, invitationID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apierr.NotFound("location")
	}
	return nil
}

func (s *Service) AddGuest(ctx context.Context, owner, invitationID uuid.UUID, name, email, phone, group string) (Guest, error) {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return Guest{}, err
	}
	if strings.TrimSpace(name) == "" {
		return Guest{}, apierr.BadRequest("invalid_input", "name is required")
	}
	if err := s.billing.AssertGuestQuota(ctx, owner, invitationID); err != nil {
		return Guest{}, err
	}
	g := Guest{
		ID:        uuid.New(),
		Name:      strings.TrimSpace(name),
		Email:     strings.TrimSpace(email),
		Phone:     strings.TrimSpace(phone),
		Group:     strings.TrimSpace(group),
		RSVPToken: uuid.NewString(),
		Status:    "draft",
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO invitation_guests (id, invitation_id, name, email, phone, guest_group, rsvp_token, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,'draft')`,
		g.ID, invitationID, g.Name, g.Email, g.Phone, g.Group, g.RSVPToken)
	return g, err
}

func (s *Service) ListGuests(ctx context.Context, owner, invitationID uuid.UUID) ([]Guest, error) {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, email, phone, guest_group, rsvp_token, status, plus_ones, message, rsvp_at
		FROM invitation_guests WHERE invitation_id=$1 ORDER BY created_at`, invitationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Guest
	for rows.Next() {
		var g Guest
		if err := rows.Scan(&g.ID, &g.Name, &g.Email, &g.Phone, &g.Group, &g.RSVPToken, &g.Status, &g.PlusOnes, &g.Message, &g.RSVPAt); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

// GuestPatch carries only the fields the host sent; nil means "leave as is".
// rsvp_token and the guest's own RSVP answer are not editable by the host.
type GuestPatch struct {
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	Phone  *string `json:"phone"`
	Group  *string `json:"group"`
	Status *string `json:"status"`
}

// validateHostGuestStatus keeps RSVP answers (opened/accepted/declined) out of host edits.
func validateHostGuestStatus(status string) (string, error) {
	status = strings.TrimSpace(status)
	if status != "draft" && status != "invited" {
		return "", apierr.BadRequest("invalid_input", "status must be draft or invited")
	}
	return status, nil
}

func (s *Service) PatchGuest(ctx context.Context, owner, invitationID, guestID uuid.UUID, in GuestPatch) (Guest, error) {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return Guest{}, err
	}
	var g Guest
	err := s.pool.QueryRow(ctx, `
		SELECT id, name, email, phone, guest_group, rsvp_token, status, plus_ones, message, rsvp_at
		FROM invitation_guests WHERE id=$1 AND invitation_id=$2`, guestID, invitationID).
		Scan(&g.ID, &g.Name, &g.Email, &g.Phone, &g.Group, &g.RSVPToken, &g.Status, &g.PlusOnes, &g.Message, &g.RSVPAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Guest{}, apierr.NotFound("guest")
	}
	if err != nil {
		return Guest{}, err
	}
	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return Guest{}, apierr.BadRequest("invalid_input", "name is required")
		}
		g.Name = name
	}
	if in.Email != nil {
		g.Email = strings.TrimSpace(*in.Email)
	}
	if in.Phone != nil {
		g.Phone = strings.TrimSpace(*in.Phone)
	}
	if in.Group != nil {
		g.Group = strings.TrimSpace(*in.Group)
	}
	if in.Status != nil {
		status, err := validateHostGuestStatus(*in.Status)
		if err != nil {
			return Guest{}, err
		}
		g.Status = status
	}
	_, err = s.pool.Exec(ctx, `
		UPDATE invitation_guests
		SET name=$3, email=$4, phone=$5, guest_group=$6, status=$7, updated_at=now()
		WHERE id=$1 AND invitation_id=$2`,
		g.ID, invitationID, g.Name, g.Email, g.Phone, g.Group, g.Status)
	return g, err
}

func (s *Service) DeleteGuest(ctx context.Context, owner, invitationID, guestID uuid.UUID) error {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `DELETE FROM invitation_guests WHERE id=$1 AND invitation_id=$2`, guestID, invitationID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apierr.NotFound("guest")
	}
	return nil
}

type PublicInvitation struct {
	ID         uuid.UUID        `json:"-"`
	Title      string           `json:"title"`
	Content    json.RawMessage  `json:"content"`
	Slug       string           `json:"slug"`
	Locations  []Location       `json:"locations"`
	Inviters   []Inviter        `json:"inviters"`
	Gallery    []Media          `json:"gallery"`
	Music      *Media           `json:"music"`
	TemplateID uuid.UUID        `json:"template_id"`
	Guestbook  []GuestbookEntry `json:"guestbook"`
}

type PublicRSVP struct {
	InvitationID uuid.UUID        `json:"-"`
	GuestID      uuid.UUID        `json:"-"`
	GuestName    string           `json:"guest_name"`
	Title        string           `json:"title"`
	Content      json.RawMessage  `json:"content"`
	Status       string           `json:"status"`
	Locations    []Location       `json:"locations"`
	Inviters     []Inviter        `json:"inviters"`
	Gallery      []Media          `json:"gallery"`
	Music        *Media           `json:"music"`
	TemplateID   uuid.UUID        `json:"template_id"`
	Guestbook    []GuestbookEntry `json:"guestbook"`
}

func (s *Service) PublicGet(ctx context.Context, token string) (PublicRSVP, error) {
	var inv Invitation
	var guestID uuid.UUID
	var guestName, status string
	err := s.pool.QueryRow(ctx, `
		SELECT g.id, g.name, g.status, i.id, i.owner_id, i.type_id, i.template_id, i.template_version, i.slug, i.status, i.title, i.content, i.starts_at, i.ends_at
		FROM invitation_guests g JOIN invitations i ON i.id = g.invitation_id
		WHERE g.rsvp_token=$1 AND i.status='published'`, token).
		Scan(&guestID, &guestName, &status, &inv.ID, &inv.OwnerID, &inv.TypeID, &inv.TemplateID, &inv.TemplateVersion, &inv.Slug, &inv.Status, &inv.Title, &inv.Content, &inv.StartsAt, &inv.EndsAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return PublicRSVP{}, apierr.NotFound("invitation")
	}
	if err != nil {
		return PublicRSVP{}, err
	}
	if status == "draft" || status == "invited" {
		_, _ = s.pool.Exec(ctx, `UPDATE invitation_guests SET status='opened', updated_at=now() WHERE rsvp_token=$1 AND status IN ('draft','invited')`, token)
		status = "opened"
	}
	page, err := s.publicDecorations(ctx, inv)
	if err != nil {
		return PublicRSVP{}, err
	}
	return PublicRSVP{
		InvitationID: inv.ID,
		GuestID:      guestID,
		GuestName:    guestName,
		Title:        page.Title,
		Content:      page.Content,
		Status:       status,
		Locations:    page.Locations,
		Inviters:     page.Inviters,
		Gallery:      page.Gallery,
		Music:        page.Music,
		TemplateID:   page.TemplateID,
		Guestbook:    page.Guestbook,
	}, nil
}

func (s *Service) PublicPreview(ctx context.Context, slug string) (PublicInvitation, error) {
	inv, err := s.publishedBySlug(ctx, slug)
	if err != nil {
		return PublicInvitation{}, err
	}
	return s.publicDecorations(ctx, inv)
}

func (s *Service) publishedBySlug(ctx context.Context, slug string) (Invitation, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return Invitation{}, apierr.NotFound("invitation")
	}
	inv, err := scanInvitation(s.pool.QueryRow(ctx, `
		SELECT id, owner_id, type_id, template_id, template_version, slug, status, title, content, starts_at, ends_at
		FROM invitations WHERE slug=$1 AND status='published'`, slug))
	if errors.Is(err, pgx.ErrNoRows) {
		return Invitation{}, apierr.NotFound("invitation")
	}
	return inv, err
}

func (s *Service) publicDecorations(ctx context.Context, inv Invitation) (PublicInvitation, error) {
	locs, err := s.ListLocations(ctx, inv.ID)
	if err != nil {
		return PublicInvitation{}, err
	}
	inviters, err := s.listInviters(ctx, inv.ID)
	if err != nil {
		return PublicInvitation{}, err
	}
	gallery, music, err := s.listMedia(ctx, inv.ID)
	if err != nil {
		return PublicInvitation{}, err
	}
	book, err := s.listGuestbook(ctx, inv.ID)
	if err != nil {
		return PublicInvitation{}, err
	}
	slug := ""
	if inv.Slug != nil {
		slug = *inv.Slug
	}
	return PublicInvitation{
		ID:         inv.ID,
		Title:      inv.Title,
		Content:    inv.Content,
		Slug:       slug,
		Locations:  locs,
		Inviters:   inviters,
		Gallery:    gallery,
		Music:      music,
		TemplateID: inv.TemplateID,
		Guestbook:  book,
	}, nil
}

func (s *Service) PublicRSVP(ctx context.Context, token, decision, message string, plusOnes int) (Guest, error) {
	if decision != "accepted" && decision != "declined" {
		return Guest{}, apierr.BadRequest("invalid_input", "status must be accepted or declined")
	}
	var g Guest
	var invitationID uuid.UUID
	err := s.pool.QueryRow(ctx, `
		SELECT g.id, g.name, g.email, g.phone, g.guest_group, g.rsvp_token, g.status, g.plus_ones, g.message, g.rsvp_at, g.invitation_id
		FROM invitation_guests g JOIN invitations i ON i.id = g.invitation_id
		WHERE g.rsvp_token=$1 AND i.status='published'`, token).
		Scan(&g.ID, &g.Name, &g.Email, &g.Phone, &g.Group, &g.RSVPToken, &g.Status, &g.PlusOnes, &g.Message, &g.RSVPAt, &invitationID)
	if errors.Is(err, pgx.ErrNoRows) {
		return Guest{}, apierr.NotFound("guest")
	}
	if err != nil {
		return Guest{}, err
	}
	now := time.Now()
	g.Status = decision
	g.Message = message
	g.PlusOnes = plusOnes
	g.RSVPAt = &now
	_, err = s.pool.Exec(ctx, `UPDATE invitation_guests SET status=$2, message=$3, plus_ones=$4, rsvp_at=$5, updated_at=now() WHERE id=$1`,
		g.ID, g.Status, g.Message, g.PlusOnes, g.RSVPAt)
	if err != nil {
		return Guest{}, err
	}
	_ = s.notify.Send(ctx, notify.Message{
		Channel:  "noop",
		Template: "rsvp.submitted",
		Payload:  map[string]any{"guest_id": g.ID.String(), "status": decision},
	})
	return g, nil
}

func (s *Service) get(ctx context.Context, id uuid.UUID) (Invitation, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT id, owner_id, type_id, template_id, template_version, slug, status, title, content, starts_at, ends_at
		FROM invitations WHERE id=$1`, id)
	inv, err := scanInvitation(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Invitation{}, apierr.NotFound("invitation")
	}
	return inv, err
}

type scanner interface {
	Scan(dest ...any) error
}

func scanInvitation(row scanner) (Invitation, error) {
	var inv Invitation
	err := row.Scan(&inv.ID, &inv.OwnerID, &inv.TypeID, &inv.TemplateID, &inv.TemplateVersion, &inv.Slug, &inv.Status, &inv.Title, &inv.Content, &inv.StartsAt, &inv.EndsAt)
	return inv, err
}

func uniqueSlug(title string, id uuid.UUID) string {
	return slugify(title) + "-" + strings.Split(id.String(), "-")[0]
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			prevDash = false
			continue
		}
		if !prevDash {
			b.WriteByte('-')
			prevDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "undangan"
	}
	return out
}
