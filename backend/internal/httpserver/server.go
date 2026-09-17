package httpserver

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"ketuk.id/api/internal/apierr"
	"ketuk.id/api/internal/authctx"
	"ketuk.id/api/internal/billing"
	"ketuk.id/api/internal/catalog"
	"ketuk.id/api/internal/config"
	"ketuk.id/api/internal/gift"
	"ketuk.id/api/internal/httputil"
	"ketuk.id/api/internal/identity"
	"ketuk.id/api/internal/invitation"
	"ketuk.id/api/internal/link"
	"ketuk.id/api/internal/pay"
	"ketuk.id/api/internal/planner"
	"ketuk.id/api/internal/sweepjob"
)

type Server struct {
	cfg         config.Config
	pool        *pgxpool.Pool
	identity    *identity.Service
	billing     *billing.Service
	catalog     *catalog.Service
	invitations *invitation.Service
	planners    *planner.Service
	links       *link.Service
	gifts       *gift.Service
	pay         pay.Gateway
}

func New(cfg config.Config, pool *pgxpool.Pool, idn *identity.Service, bill *billing.Service, cat *catalog.Service, inv *invitation.Service, pl *planner.Service, ln *link.Service, gf *gift.Service, gw pay.Gateway) http.Handler {
	s := &Server{cfg: cfg, pool: pool, identity: idn, billing: bill, catalog: cat, invitations: inv, planners: pl, links: ln, gifts: gf, pay: gw}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Get("/healthz", s.health)
	r.Get("/api/cron/sweep", s.cronSweep)
	r.Post("/api/cron/sweep", s.cronSweep)

	r.Route("/v1", func(r chi.Router) {
		// Credential endpoints are the only brute-force target worth throttling.
		r.Group(func(r chi.Router) {
			r.Use(newRateLimiter(cfg.AuthRateLimit, time.Minute).middleware)
			r.Post("/auth/register", s.register)
			r.Post("/auth/login", s.login)
			r.Post("/auth/refresh", s.refresh)
			r.Post("/auth/logout", s.logout)
		})

		r.Get("/event-types", s.eventTypes)
		r.Get("/templates", s.templates)
		r.Get("/billing/plans", s.plans)

		r.Get("/public/rsvp/{token}", s.publicRSVPGet)
		r.Post("/public/rsvp/{token}", s.publicRSVPPost)
		r.Post("/public/rsvp/{token}/gifts", s.publicGift)
		r.Post("/public/rsvp/{token}/guestbook", s.publicRSVPGuestbook)
		r.Get("/public/invitations/{slug}/guestbook", s.publicGuestbookList)
		r.Post("/public/invitations/{slug}/guestbook", s.publicGuestbookAdd)
		r.Get("/public/invitations/{slug}", s.publicInvitationGet)
		r.Get("/public/media/{id}", s.publicMedia)
		r.Post("/public/payments/duitku/callback", s.duitkuCallback)

		r.Group(func(r chi.Router) {
			r.Use(s.auth)
			r.Get("/me", s.me)
			r.Patch("/me", s.patchMe)
			r.Post("/auth/password", s.changePassword)
			r.Get("/billing/subscription", s.subscription)
			r.Post("/billing/checkout", s.checkout)
			r.Get("/orders", s.listOrders)
			r.Post("/orders/{id}/sync", s.syncOrder)

			r.Get("/invitations", s.listInvitations)
			r.Post("/invitations", s.createInvitation)
			r.Get("/invitations/{id}", s.getInvitation)
			r.Patch("/invitations/{id}", s.patchInvitation)
			r.Delete("/invitations/{id}", s.archiveInvitation)
			r.Post("/invitations/{id}/publish", s.publishInvitation)
			r.Patch("/invitations/{id}/template", s.setTemplate)
			r.Get("/invitations/{id}/locations", s.listLocations)
			r.Post("/invitations/{id}/locations", s.addLocation)
			r.Patch("/invitations/{id}/locations/{locationId}", s.patchLocation)
			r.Delete("/invitations/{id}/locations/{locationId}", s.deleteLocation)
			r.Get("/invitations/{id}/guests", s.listGuests)
			r.Post("/invitations/{id}/guests", s.addGuest)
			r.Patch("/invitations/{id}/guests/{guestId}", s.patchGuest)
			r.Delete("/invitations/{id}/guests/{guestId}", s.deleteGuest)
			r.Get("/invitations/{id}/inviters", s.listInviters)
			r.Post("/invitations/{id}/inviters", s.addInviter)
			r.Patch("/invitations/{id}/inviters/{inviterId}", s.patchInviter)
			r.Delete("/invitations/{id}/inviters/{inviterId}", s.deleteInviter)
			r.Get("/invitations/{id}/gallery", s.listGallery)
			r.Post("/invitations/{id}/gallery", s.addGallery)
			r.Patch("/invitations/{id}/gallery/{mediaId}", s.patchGallery)
			r.Delete("/invitations/{id}/gallery/{mediaId}", s.deleteGallery)
			r.Get("/invitations/{id}/music", s.getMusic)
			r.Put("/invitations/{id}/music", s.setMusic)
			r.Delete("/invitations/{id}/music", s.deleteMusic)
			r.Get("/invitations/{id}/gift-methods", s.listGiftMethods)
			r.Post("/invitations/{id}/gift-methods", s.addGiftMethod)
			r.Delete("/invitations/{id}/gift-methods/{methodId}", s.deleteGiftMethod)
			r.Get("/invitations/{id}/gifts", s.listGifts)
			r.Get("/invitations/{id}/guestbook", s.listGuestbook)
			r.Patch("/invitations/{id}/guestbook/{entryId}", s.patchGuestbook)
			r.Delete("/invitations/{id}/guestbook/{entryId}", s.deleteGuestbook)

			r.Get("/planners", s.listPlanners)
			r.Post("/planners", s.createPlanner)
			r.Get("/planners/{id}", s.getPlanner)
			r.Patch("/planners/{id}", s.patchPlanner)
			r.Delete("/planners/{id}", s.archivePlanner)
			r.Get("/planners/{id}/days", s.listDays)
			r.Post("/planners/{id}/days", s.addDay)
			r.Patch("/planners/{id}/days/{dayId}", s.patchDay)
			r.Delete("/planners/{id}/days/{dayId}", s.deleteDay)
			r.Get("/planners/{id}/days/{dayId}/sessions", s.listSessions)
			r.Post("/planners/{id}/days/{dayId}/sessions", s.addSession)
			r.Patch("/planners/{id}/days/{dayId}/sessions/{sessionId}", s.patchSession)
			r.Delete("/planners/{id}/days/{dayId}/sessions/{sessionId}", s.deleteSession)
			r.Get("/planners/{id}/categories", s.listCategories)
			r.Post("/planners/{id}/categories", s.addCategory)
			r.Patch("/planners/{id}/categories/{categoryId}", s.patchCategory)
			r.Delete("/planners/{id}/categories/{categoryId}", s.deleteCategory)
			r.Post("/planners/{id}/budget/items", s.addBudgetItem)
			r.Patch("/planners/{id}/budget/items/{itemId}", s.patchBudgetItem)
			r.Delete("/planners/{id}/budget/items/{itemId}", s.deleteBudgetItem)
			r.Get("/planners/{id}/budget", s.budget)
			r.Get("/planners/{id}/checklist", s.listChecks)
			r.Post("/planners/{id}/checklist", s.addCheck)
			r.Patch("/planners/{id}/checklist/{itemId}", s.patchCheck)
			r.Delete("/planners/{id}/checklist/{itemId}", s.deleteCheck)

			r.Get("/links", s.listLinks)
			r.Post("/links", s.createLink)
			r.Delete("/links/{id}", s.deleteLink)
		})
	})
	return r
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	if err := s.pool.Ping(r.Context()); err != nil {
		// Monitoring needs "dependency down", not a generic 500.
		httputil.Error(w, apierr.Unavailable("database_unavailable", "cannot reach Postgres"))
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// cronSweep is hit by an external scheduler (Vercel's own Cron Jobs are
// capped at once/day on the Hobby plan, too coarse for this). CronSecret
// unset means the route is unreachable rather than silently open.
func (s *Server) cronSweep(w http.ResponseWriter, r *http.Request) {
	secret := s.cfg.CronSecret
	if secret == "" || r.Header.Get("Authorization") != "Bearer "+secret {
		httputil.Error(w, apierr.Unauthorized("invalid cron secret"))
		return
	}
	sweepjob.RunOnce(r.Context(), s.identity, s.billing, s.gifts)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			httputil.Error(w, apierr.Unauthorized("missing bearer token"))
			return
		}
		id, err := s.identity.ParseAccess(strings.TrimPrefix(h, "Bearer "))
		if err != nil {
			httputil.Error(w, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(authctx.WithUserID(r.Context(), id)))
	})
}

func userID(r *http.Request) uuid.UUID {
	id, _ := authctx.UserID(r.Context())
	return id
}

func parseID(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, key))
	if err != nil {
		httputil.Error(w, apierr.BadRequest("invalid_id", "invalid "+key))
		return uuid.Nil, false
	}
	return id, true
}

// optionalID reads a uuid from a query parameter that may be absent.
func optionalID(w http.ResponseWriter, r *http.Request, key string) (*uuid.UUID, bool) {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return nil, true
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		httputil.Error(w, apierr.BadRequest("invalid_id", "invalid "+key))
		return nil, false
	}
	return &id, true
}

func data(v any) map[string]any { return map[string]any{"data": v} }
