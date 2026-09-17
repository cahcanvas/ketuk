package httpserver

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"ketuk.id/api/internal/apierr"
	"ketuk.id/api/internal/gift"
	"ketuk.id/api/internal/httputil"
	"ketuk.id/api/internal/invitation"
	"ketuk.id/api/internal/pay"
)

const maxUpload = 16 << 20

func (s *Server) listInviters(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	out, err := s.invitations.ListInviters(r.Context(), userID(r), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) addInviter(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	var body struct {
		Name      string `json:"name"`
		Role      string `json:"role"`
		SortOrder int    `json:"sort_order"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.invitations.AddInviter(r.Context(), userID(r), id, body.Name, body.Role, body.SortOrder)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) patchInviter(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	inviterID, ok := parseID(w, r, "inviterId")
	if !ok {
		return
	}
	var body invitation.InviterPatch
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.invitations.PatchInviter(r.Context(), userID(r), id, inviterID, body)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) deleteInviter(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	inviterID, ok := parseID(w, r, "inviterId")
	if !ok {
		return
	}
	if err := s.invitations.DeleteInviter(r.Context(), userID(r), id, inviterID); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listGallery(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	out, err := s.invitations.ListGallery(r.Context(), userID(r), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) addGallery(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	file, name, contentType, err := readUpload(r, "file")
	if err != nil {
		httputil.Error(w, err)
		return
	}
	defer file.Close()
	out, err := s.invitations.AddGallery(r.Context(), userID(r), id, name, contentType, r.FormValue("caption"), file)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) patchGallery(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	mediaID, ok := parseID(w, r, "mediaId")
	if !ok {
		return
	}
	var body invitation.GalleryPatch
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.invitations.PatchGallery(r.Context(), userID(r), id, mediaID, body)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) deleteGallery(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	mediaID, ok := parseID(w, r, "mediaId")
	if !ok {
		return
	}
	if err := s.invitations.DeleteMedia(r.Context(), userID(r), id, mediaID); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getMusic(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	out, err := s.invitations.GetMusic(r.Context(), userID(r), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) setMusic(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	file, name, contentType, err := readUpload(r, "file")
	if err != nil {
		httputil.Error(w, err)
		return
	}
	defer file.Close()
	out, err := s.invitations.SetMusic(r.Context(), userID(r), id, name, contentType, file)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) deleteMusic(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	music, err := s.invitations.GetMusic(r.Context(), userID(r), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	if music == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := s.invitations.DeleteMedia(r.Context(), userID(r), id, music.ID); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) publicMedia(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	url, contentType, original, body, err := s.invitations.ResolveMedia(r.Context(), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	if url != "" {
		http.Redirect(w, r, url, http.StatusFound)
		return
	}
	defer body.Close()
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	if original != "" {
		w.Header().Set("Content-Disposition", "inline; filename="+filepath.Base(original))
	}
	io.Copy(w, body)
}

func (s *Server) listGiftMethods(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	out, err := s.gifts.ListMethods(r.Context(), userID(r), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) addGiftMethod(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	var m gift.Method
	if err := httputil.Decode(r, &m); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.gifts.AddMethod(r.Context(), userID(r), id, m)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) deleteGiftMethod(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	methodID, ok := parseID(w, r, "methodId")
	if !ok {
		return
	}
	if err := s.gifts.DeleteMethod(r.Context(), userID(r), id, methodID); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listGifts(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	out, err := s.gifts.ListPayments(r.Context(), userID(r), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) listGuestbook(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	out, err := s.invitations.ListGuestbook(r.Context(), userID(r), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) patchGuestbook(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	entryID, ok := parseID(w, r, "entryId")
	if !ok {
		return
	}
	var body struct {
		Name    *string `json:"name"`
		Message *string `json:"message"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.invitations.PatchGuestbook(r.Context(), userID(r), id, entryID, body.Name, body.Message)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) deleteGuestbook(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	entryID, ok := parseID(w, r, "entryId")
	if !ok {
		return
	}
	if err := s.invitations.DeleteGuestbook(r.Context(), userID(r), id, entryID); err != nil {
		httputil.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) publicInvitationGet(w http.ResponseWriter, r *http.Request) {
	out, err := s.invitations.PublicPreview(r.Context(), chiParam(r, "slug"))
	if err != nil {
		httputil.Error(w, err)
		return
	}
	methods, err := s.gifts.ListPublicMethods(r.Context(), out.ID)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(map[string]any{
		"title":       out.Title,
		"content":     out.Content,
		"slug":        out.Slug,
		"locations":   out.Locations,
		"inviters":    out.Inviters,
		"gallery":     out.Gallery,
		"music":       out.Music,
		"template_id": out.TemplateID,
		"guestbook":   out.Guestbook,
		"gifts":       methods,
	}))
}

func (s *Server) publicGuestbookList(w http.ResponseWriter, r *http.Request) {
	out, err := s.invitations.PublicPreview(r.Context(), chiParam(r, "slug"))
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out.Guestbook))
}

func (s *Server) publicGuestbookAdd(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.invitations.AddGuestbookPublic(r.Context(), chiParam(r, "slug"), body.Name, body.Message)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) publicRSVPGuestbook(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.invitations.AddGuestbookByToken(r.Context(), chiParam(r, "token"), body.Name, body.Message)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) publicGift(w http.ResponseWriter, r *http.Request) {
	page, err := s.invitations.PublicGet(r.Context(), chiParam(r, "token"))
	if err != nil {
		httputil.Error(w, err)
		return
	}
	var body struct {
		MethodID   uuid.UUID `json:"method_id"`
		AmountIDR  int64     `json:"amount_idr"`
		SenderName string    `json:"sender_name"`
		Message    string    `json:"message"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	out, err := s.gifts.CreatePublic(r.Context(), page.InvitationID, page.GuestID, page.GuestName, body.MethodID, body.AmountIDR, body.SenderName, body.Message)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) checkout(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Plan string `json:"plan"`
	}
	if err := httputil.Decode(r, &body); err != nil {
		httputil.Error(w, err)
		return
	}
	user, err := s.identity.Get(r.Context(), userID(r))
	if err != nil {
		httputil.Error(w, err)
		return
	}
	customerName := user.Username
	if user.Name != nil && *user.Name != "" {
		customerName = *user.Name
	}
	out, err := s.billing.Checkout(r.Context(), user.ID, body.Plan, user.Email, customerName)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, data(out))
}

func (s *Server) listOrders(w http.ResponseWriter, r *http.Request) {
	out, err := s.billing.ListOrders(r.Context(), userID(r), r.URL.Query().Get("kind"))
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) syncOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	out, err := s.billing.SyncOrder(r.Context(), userID(r), id)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, data(out))
}

func (s *Server) duitkuCallback(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	merchantCode := r.FormValue("merchantCode")
	amount := r.FormValue("amount")
	orderID := r.FormValue("merchantOrderId")
	signature := r.FormValue("signature")
	if s.pay == nil || !s.pay.VerifyCallback(merchantCode, amount, orderID, signature) {
		httputil.Error(w, apierr.Unauthorized("invalid duitku signature"))
		return
	}
	cb := pay.Callback{
		MerchantOrderID: orderID,
		Amount:          amount,
		ResultCode:      r.FormValue("resultCode"),
		Reference:       r.FormValue("reference"),
		AdditionalParam: r.FormValue("additionalParam"),
		Success:         r.FormValue("resultCode") == "00",
	}
	// Both kas must be offered the callback: each ignores order ids that are not its own.
	billErr := s.billing.ConfirmCallback(r.Context(), cb)
	giftErr := s.gifts.ConfirmCallback(r.Context(), cb)
	if err := errors.Join(billErr, giftErr); err != nil {
		httputil.Error(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("SUCCESS"))
}

func readUpload(r *http.Request, field string) (io.ReadCloser, string, string, error) {
	if err := r.ParseMultipartForm(maxUpload); err != nil {
		return nil, "", "", apierr.BadRequest("invalid_input", "multipart file is required")
	}
	file, header, err := r.FormFile(field)
	if err != nil {
		return nil, "", "", apierr.BadRequest("invalid_input", "file is required")
	}
	ct := header.Header.Get("Content-Type")
	if strings.TrimSpace(ct) == "" {
		ct = "application/octet-stream"
	}
	return file, header.Filename, ct, nil
}
