// Package handler is the Vercel Cron target for the same job cmd/api/main.go
// runs on an in-process ticker (see internal/sweepjob). Vercel invokes this
// once per schedule tick (see vercel.json) instead of keeping a goroutine
// alive, since serverless instances aren't guaranteed to stay up between
// requests.
package handler

import (
	"context"
	"log"
	"net/http"
	"os"

	"ketuk.id/api/internal/app"
	"ketuk.id/api/internal/sweepjob"
)

var built *app.App

func init() {
	a, err := app.Build(context.Background())
	if err != nil {
		log.Fatalf("app build: %v", err)
	}
	built = a
}

// Handler requires the CRON_SECRET Vercel sends as a bearer token on
// scheduled invocations, so the endpoint can't be triggered by anyone who
// finds the URL. Set CRON_SECRET in the Vercel project's environment
// variables — Vercel attaches it to Cron requests automatically once set.
func Handler(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("CRON_SECRET")
	if secret == "" || r.Header.Get("Authorization") != "Bearer "+secret {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	sweepjob.RunOnce(r.Context(), built.Identity, built.Billing, built.Gift)
	w.WriteHeader(http.StatusNoContent)
}
