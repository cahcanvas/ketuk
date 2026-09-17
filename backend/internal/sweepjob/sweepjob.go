// Package sweepjob closes out abandoned checkouts and prunes dead refresh
// sessions. RunOnce is the unit of work; Loop drives it on a ticker for a
// persistent process (cmd/api), while a serverless host (Vercel) instead
// calls RunOnce once per Cron invocation — see api/cron/sweep.go.
package sweepjob

import (
	"context"
	"log"
	"time"

	"ketuk.id/api/internal/billing"
	"ketuk.id/api/internal/gift"
	"ketuk.id/api/internal/identity"
)

// Duitku invoices expire after 60 minutes, so 90 leaves room for a late
// callback.
const UnpaidAfter = 90 * time.Minute

func RunOnce(ctx context.Context, idn *identity.Service, bill *billing.Service, gf *gift.Service) {
	if n, err := bill.ExpirePending(ctx, UnpaidAfter); err != nil {
		log.Printf("sweep orders: %v", err)
	} else if n > 0 {
		log.Printf("sweep: %d order(s) expired", n)
	}
	if n, err := gf.ExpirePending(ctx, UnpaidAfter); err != nil {
		log.Printf("sweep gifts: %v", err)
	} else if n > 0 {
		log.Printf("sweep: %d gift payment(s) expired", n)
	}
	if _, err := idn.PurgeExpiredSessions(ctx); err != nil {
		log.Printf("sweep sessions: %v", err)
	}
}

func Loop(ctx context.Context, every time.Duration, idn *identity.Service, bill *billing.Service, gf *gift.Service) {
	ticker := time.NewTicker(every)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			RunOnce(ctx, idn, bill, gf)
		}
	}
}
