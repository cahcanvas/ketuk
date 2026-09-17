// Package handler is the Vercel Go entrypoint: one function serves every
// route via the same httpserver.Handler cmd/api/main.go runs. vercel.json
// rewrites /healthz and /v1/* here; everything else about the app is wired
// once in internal/app.Build.
package handler

import (
	"context"
	"log"
	"net/http"

	"ketuk.id/api/internal/app"
)

var built *app.App

func init() {
	a, err := app.Build(context.Background())
	if err != nil {
		log.Fatalf("app build: %v", err)
	}
	built = a
}

func Handler(w http.ResponseWriter, r *http.Request) {
	built.Handler.ServeHTTP(w, r)
}
