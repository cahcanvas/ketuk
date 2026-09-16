package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"ketuk.id/api/internal/httputil"
	"ketuk.id/api/internal/planner"
)

func chiParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func (s *Server) listPlanners(w http.ResponseWriter, r *http.Request) {
	out, err := s.planners.List(r.Context(), userID(r))
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) createPlanner(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title      string          `json:"title"`
		BudgetMode string          `json:"budget_mode"`
		TypeID     *uuid.UUID      `json:"type_id"`
		Categories json.RawMessage `json:"categories"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.planners.Create(r.Context(), userID(r), body.Title, body.BudgetMode, body.TypeID, body.Categories)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) getPlanner(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	out, err := s.planners.Get(r.Context(), userID(r), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) patchPlanner(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	var body planner.PlannerPatch
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.planners.Patch(r.Context(), userID(r), id, body)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) archivePlanner(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	if err := s.planners.Archive(r.Context(), userID(r), id); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listDays(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	if _, err := s.planners.Get(r.Context(), userID(r), id); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.planners.ListDays(r.Context(), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) addDay(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	var body struct {
		Date  time.Time `json:"date"`
		Label string    `json:"label"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.planners.AddDay(r.Context(), userID(r), id, body.Date, body.Label)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) patchDay(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	dayID, ok := parseID(w, r, "dayId")
	if !ok {
		return
	}
	var body planner.DayPatch
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.planners.PatchDay(r.Context(), userID(r), id, dayID, body)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) deleteDay(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	dayID, ok := parseID(w, r, "dayId")
	if !ok {
		return
	}
	if err := s.planners.DeleteDay(r.Context(), userID(r), id, dayID); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listSessions(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	dayID, ok := parseID(w, r, "dayId")
	if !ok {
		return
	}
	out, err := s.planners.ListSessions(r.Context(), userID(r), id, dayID)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) addSession(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	dayID, ok := parseID(w, r, "dayId")
	if !ok {
		return
	}
	var body struct {
		Title     string  `json:"title"`
		StartsAt  string  `json:"starts_at"`
		EndsAt    *string `json:"ends_at"`
		PicName   string  `json:"pic_name"`
		SortOrder int     `json:"sort_order"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.planners.AddSession(r.Context(), userID(r), id, dayID, body.Title, body.StartsAt, body.EndsAt, body.PicName, body.SortOrder)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) patchSession(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	dayID, ok := parseID(w, r, "dayId")
	if !ok {
		return
	}
	sessionID, ok := parseID(w, r, "sessionId")
	if !ok {
		return
	}
	var body planner.SessionPatch
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.planners.PatchSession(r.Context(), userID(r), id, dayID, sessionID, body)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) deleteSession(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	dayID, ok := parseID(w, r, "dayId")
	if !ok {
		return
	}
	sessionID, ok := parseID(w, r, "sessionId")
	if !ok {
		return
	}
	if err := s.planners.DeleteSession(r.Context(), userID(r), id, dayID, sessionID); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listCategories(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	if _, err := s.planners.Get(r.Context(), userID(r), id); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.planners.ListCategories(r.Context(), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) addCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.planners.AddCategory(r.Context(), userID(r), id, body.Name)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) patchCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	categoryID, ok := parseID(w, r, "categoryId")
	if !ok {
		return
	}
	var body planner.CategoryPatch
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.planners.PatchCategory(r.Context(), userID(r), id, categoryID, body)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) deleteCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	categoryID, ok := parseID(w, r, "categoryId")
	if !ok {
		return
	}
	if err := s.planners.DeleteCategory(r.Context(), userID(r), id, categoryID); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) addBudgetItem(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	var body struct {
		CategoryID uuid.UUID  `json:"category_id"`
		DayID      *uuid.UUID `json:"day_id"`
		Name       string     `json:"name"`
		Amount     int64      `json:"amount"`
		Note       string     `json:"note"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.planners.AddItem(r.Context(), userID(r), id, body.CategoryID, body.Name, body.Amount, body.DayID, body.Note)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) patchBudgetItem(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	itemID, ok := parseID(w, r, "itemId")
	if !ok {
		return
	}
	var body planner.ItemPatch
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.planners.PatchItem(r.Context(), userID(r), id, itemID, body)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) deleteBudgetItem(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	itemID, ok := parseID(w, r, "itemId")
	if !ok {
		return
	}
	if err := s.planners.DeleteItem(r.Context(), userID(r), id, itemID); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) budget(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	out, err := s.planners.Budget(r.Context(), userID(r), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) listChecks(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	out, err := s.planners.ListChecks(r.Context(), userID(r), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) addCheck(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	var body struct {
		Title     string     `json:"title"`
		DayID     *uuid.UUID `json:"day_id"`
		DueAt     *time.Time `json:"due_at"`
		SortOrder int        `json:"sort_order"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.planners.AddCheck(r.Context(), userID(r), id, body.Title, body.DayID, body.DueAt, body.SortOrder)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) patchCheck(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	itemID, ok := parseID(w, r, "itemId")
	if !ok {
		return
	}
	var body planner.CheckPatch
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.planners.PatchCheck(r.Context(), userID(r), id, itemID, body)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) deleteCheck(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	itemID, ok := parseID(w, r, "itemId")
	if !ok {
		return
	}
	if err := s.planners.DeleteCheck(r.Context(), userID(r), id, itemID); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listLinks(w http.ResponseWriter, r *http.Request) {
	invitationID, ok := optionalID(w, r, "invitation_id")
	if !ok {
		return
	}
	plannerID, ok := optionalID(w, r, "planner_id")
	if !ok {
		return
	}
	out, err := s.links.List(r.Context(), userID(r), invitationID, plannerID)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) createLink(w http.ResponseWriter, r *http.Request) {
	var body struct {
		InvitationID uuid.UUID `json:"invitation_id"`
		PlannerID    uuid.UUID `json:"planner_id"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.links.Create(r.Context(), userID(r), body.InvitationID, body.PlannerID)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) deleteLink(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	if err := s.links.Delete(r.Context(), userID(r), id); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
