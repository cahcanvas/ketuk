package httputil

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"ketuk.id/api/internal/apierr"
)

type errorBody struct {
	Error *apierr.Error `json:"error"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Decode(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(dst); err != nil {
		return apierr.BadRequest("invalid_json", "invalid JSON body")
	}
	return nil
}

func Error(w http.ResponseWriter, err error) {
	var ae *apierr.Error
	if errors.As(err, &ae) {
		JSON(w, ae.Status, errorBody{Error: ae})
		return
	}
	log.Printf("request failed: %v", err)
	JSON(w, http.StatusInternalServerError, errorBody{
		Error: apierr.New(500, "internal", "internal server error"),
	})
}
