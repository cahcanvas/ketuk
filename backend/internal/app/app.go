// Package app is the composition root: it wires config, the database, and
// every domain service into one http.Handler. cmd/api/main.go (a persistent
// process) and api/index.go (a Vercel serverless function) both call Build
// so the wiring — and any future change to it — lives in exactly one place.
package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"ketuk.id/api/internal/billing"
	"ketuk.id/api/internal/catalog"
	"ketuk.id/api/internal/config"
	"ketuk.id/api/internal/db"
	"ketuk.id/api/internal/gift"
	"ketuk.id/api/internal/httpserver"
	"ketuk.id/api/internal/identity"
	"ketuk.id/api/internal/invitation"
	"ketuk.id/api/internal/link"
	"ketuk.id/api/internal/migrate"
	"ketuk.id/api/internal/notify"
	"ketuk.id/api/internal/pay"
	"ketuk.id/api/internal/planner"
	"ketuk.id/api/internal/storage"
)

type App struct {
	Config   config.Config
	Pool     *pgxpool.Pool
	Handler  http.Handler
	Identity *identity.Service
	Billing  *billing.Service
	Gift     *gift.Service
}

func Build(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if cfg.MigrateOnStart {
		if err := migrate.Up(cfg.DatabaseURL); err != nil {
			pool.Close()
			return nil, fmt.Errorf("migrate: %w", err)
		}
	}

	idn := identity.New(identity.NewPostgresRepository(pool), cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)

	var gw pay.Gateway = pay.Disabled{}
	if cfg.Duitku.MerchantCode != "" && cfg.Duitku.APIKey != "" {
		gw = pay.NewDuitku(cfg.Duitku.MerchantCode, cfg.Duitku.APIKey, cfg.Duitku.BaseURL)
	}
	callbackURL := cfg.PublicBaseURL + "/v1/public/payments/duitku/callback"
	bill := billing.New(pool, gw, callbackURL, cfg.Duitku.ReturnURL)
	cat := catalog.New(pool)

	var store storage.Store
	if cfg.Supabase.URL != "" && cfg.Supabase.ServiceRoleKey != "" {
		store = storage.NewSupabase(cfg.Supabase.URL, cfg.Supabase.StorageBucket, cfg.Supabase.ServiceRoleKey)
	} else {
		if err := os.MkdirAll(cfg.StorageDir, 0o755); err != nil {
			pool.Close()
			return nil, err
		}
		store = storage.NewLocal(cfg.StorageDir)
	}

	inv := invitation.New(pool, bill, cat, notify.Noop{}, store)
	pl := planner.New(pool, bill)
	ln := link.New(pool, inv, pl)
	gf := gift.New(pool, inv, gw, callbackURL, cfg.Duitku.ReturnURL)
	h := httpserver.New(cfg, pool, idn, bill, cat, inv, pl, ln, gf, gw)

	return &App{
		Config:   cfg,
		Pool:     pool,
		Handler:  h,
		Identity: idn,
		Billing:  bill,
		Gift:     gf,
	}, nil
}
