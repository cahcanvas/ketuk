package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(cfg.StorageDir, 0o755); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if cfg.MigrateOnStart {
		if err := migrate.Up(cfg.DatabaseURL); err != nil {
			log.Fatalf("migrate: %v", err)
		}
	}

	idn := identity.New(pool, cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)
	var gw pay.Gateway = pay.Disabled{}
	if cfg.Duitku.MerchantCode != "" && cfg.Duitku.APIKey != "" {
		gw = pay.NewDuitku(cfg.Duitku.MerchantCode, cfg.Duitku.APIKey, cfg.Duitku.BaseURL)
	}
	callbackURL := cfg.PublicBaseURL + "/v1/public/payments/duitku/callback"
	bill := billing.New(pool, gw, callbackURL, cfg.Duitku.ReturnURL)
	cat := catalog.New(pool)
	inv := invitation.New(pool, bill, cat, notify.Noop{}, cfg.StorageDir)
	pl := planner.New(pool, bill)
	ln := link.New(pool, inv, pl)
	gf := gift.New(pool, inv, gw, callbackURL, cfg.Duitku.ReturnURL)
	h := httpserver.New(cfg, pool, idn, bill, cat, inv, pl, ln, gf, gw)

	bgCtx, stopBackground := context.WithCancel(ctx)
	defer stopBackground()
	go sweep(bgCtx, idn, bill, gf)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           h,
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		log.Printf("ketuk-api listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	stopBackground()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

// sweep closes out abandoned checkouts and prunes dead refresh sessions.
// Duitku invoices expire after 60 minutes, so 90 leaves room for a late callback.
func sweep(ctx context.Context, idn *identity.Service, bill *billing.Service, gf *gift.Service) {
	const (
		every    = 15 * time.Minute
		unpaidAt = 90 * time.Minute
	)
	ticker := time.NewTicker(every)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if n, err := bill.ExpirePending(ctx, unpaidAt); err != nil {
				log.Printf("sweep orders: %v", err)
			} else if n > 0 {
				log.Printf("sweep: %d order(s) expired", n)
			}
			if n, err := gf.ExpirePending(ctx, unpaidAt); err != nil {
				log.Printf("sweep gifts: %v", err)
			} else if n > 0 {
				log.Printf("sweep: %d gift payment(s) expired", n)
			}
			if _, err := idn.PurgeExpiredSessions(ctx); err != nil {
				log.Printf("sweep sessions: %v", err)
			}
		}
	}
}
