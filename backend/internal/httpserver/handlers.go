package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"ketuk.id/api/internal/httputil"
	"ketuk.id/api/internal/invitation"
)

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	user, tokens, err := s.identity.Register(r.Context(), body.Email, body.Password, body.Name)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, map[string]any{"data": user, "tokens": tokens})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	user, tokens, err := s.identity.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]any{"data": user, "tokens": tokens})
}

func (s *Server) refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	tokens, err := s.identity.Refresh(r.Context(), body.RefreshToken)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]any{"tokens": tokens})
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	if err := s.identity.Logout(r.Context(), body.RefreshToken); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) changePassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	if err := s.identity.ChangePassword(r.Context(), userID(r), body.CurrentPassword, body.NewPassword); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	user, err := s.identity.Get(r.Context(), userID(r))
	if err != nil {
		httputil.Error(w, err)
		return
	}
	sub, err := s.billing.Subscription(r.Context(), user.ID)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	invs, err := s.invitations.List(r.Context(), user.ID)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	pls, err := s.planners.List(r.Context(), user.ID)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"user":         user,
			"subscription": sub,
			"invitations":  invs,
			"planners":     pls,
		},
	})
}

func (s *Server) patchMe(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.identity.UpdateProfile(r.Context(), userID(r), body.Name)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) eventTypes(w http.ResponseWriter, r *http.Request) {
	out, err := s.catalog.ListTypes(r.Context())
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) templates(w http.ResponseWriter, r *http.Request) {
	allowPremium := false
	if h := r.Header.Get("Authorization"); len(h) > 7 && h[:7] == "Bearer " {
		if uid, err := s.identity.ParseAccess(h[7:]); err == nil {
			if sub, err := s.billing.Subscription(r.Context(), uid); err == nil {
				allowPremium = sub.Entitlements["premium_templates"] == "true"
			}
		}
	}
	out, err := s.catalog.ListTemplates(r.Context(), r.URL.Query().Get("type"), allowPremium)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) plans(w http.ResponseWriter, r *http.Request) {
	out, err := s.billing.ListPlans(r.Context())
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) subscription(w http.ResponseWriter, r *http.Request) {
	out, err := s.billing.Subscription(r.Context(), userID(r))
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) listInvitations(w http.ResponseWriter, r *http.Request) {
	out, err := s.invitations.List(r.Context(), userID(r))
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) createInvitation(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title      string          `json:"title"`
		TemplateID uuid.UUID       `json:"template_id"`
		Content    json.RawMessage `json:"content"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.invitations.Create(r.Context(), userID(r), body.Title, body.TemplateID, body.Content)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) getInvitation(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	out, err := s.invitations.Get(r.Context(), userID(r), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) patchInvitation(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	var body struct {
		Title    *string         `json:"title"`
		Content  json.RawMessage `json:"content"`
		StartsAt *time.Time      `json:"starts_at"`
		EndsAt   *time.Time      `json:"ends_at"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.invitations.Patch(r.Context(), userID(r), id, body.Title, body.Content, body.StartsAt, body.EndsAt)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) archiveInvitation(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	if err := s.invitations.Archive(r.Context(), userID(r), id); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) publishInvitation(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	out, err := s.invitations.Publish(r.Context(), userID(r), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) setTemplate(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	var body struct {
		TemplateID uuid.UUID `json:"template_id"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.invitations.SetTemplate(r.Context(), userID(r), id, body.TemplateID)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) listLocations(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	if _, err := s.invitations.Get(r.Context(), userID(r), id); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.invitations.ListLocations(r.Context(), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) addLocation(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	var loc invitation.Location
	if err := httputil.Decode(r, &loc); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.invitations.AddLocation(r.Context(), userID(r), id, loc)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) patchLocation(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	locationID, ok := parseID(w, r, "locationId")
	if !ok {
		return
	}
	var body invitation.LocationPatch
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.invitations.PatchLocation(r.Context(), userID(r), id, locationID, body)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) deleteLocation(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	locationID, ok := parseID(w, r, "locationId")
	if !ok {
		return
	}
	if err := s.invitations.DeleteLocation(r.Context(), userID(r), id, locationID); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listGuests(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	out, err := s.invitations.ListGuests(r.Context(), userID(r), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) addGuest(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		Group string `json:"group"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.invitations.AddGuest(r.Context(), userID(r), id, body.Name, body.Email, body.Phone, body.Group)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) patchGuest(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	guestID, ok := parseID(w, r, "guestId")
	if !ok {
		return
	}
	var body invitation.GuestPatch
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.invitations.PatchGuest(r.Context(), userID(r), id, guestID, body)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) deleteGuest(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	guestID, ok := parseID(w, r, "guestId")
	if !ok {
		return
	}
	if err := s.invitations.DeleteGuest(r.Context(), userID(r), id, guestID); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) publicRSVPGet(w http.ResponseWriter, r *http.Request) {
	out, err := s.invitations.PublicGet(r.Context(), chiParam(r, "token"))
	if err != nil {
		httputil.Error(w, err)
		return
	}
	methods, err := s.gifts.ListPublicMethods(r.Context(), out.InvitationID)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(map[string]any{
		"guest_name":  out.GuestName,
		"title":       out.Title,
		"content":     out.Content,
		"status":      out.Status,
		"locations":   out.Locations,
		"inviters":    out.Inviters,
		"gallery":     out.Gallery,
		"music":       out.Music,
		"template_id": out.TemplateID,
		"guestbook":   out.Guestbook,
		"gifts":       methods,
	}))
}

func (s *Server) publicRSVPPost(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status   string `json:"status"`
		Message  string `json:"message"`
		PlusOnes int    `json:"plus_ones"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.invitations.PublicRSVP(r.Context(), chiParam(r, "token"), body.Status, body.Message, body.PlusOnes)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}
