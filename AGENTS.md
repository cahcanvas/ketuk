# AGENTS.md — Memori Proyek Ketuk.id

> File ini adalah "otak" asisten untuk repo ini. Baca di awal setiap sesi;
> **jangan screening ulang seluruh codebase** — perbarui bagian terkait di sini
> saja saat ada perubahan. Struktur: ringkas, padat, path relatif terhadap root
> repo kecuali diberi catatan.

## 1. Identitas produk

**Ketuk.id** — platform modular untuk urusan acara di Indonesia (pernikahan,
khitanan, aqiqah, ulang tahun, wisuda, reuni, syukuran, corporate). 4 modul
yang dijual terpisah: **Undangan** (undangan digital + RSVP), **Planner**
(budget/checklist/timeline/tamu), **Vendor** (marketplace katering/dekor/foto/
WO/MUA), **Hadiah** (hampers, bouquet, kue ke penyelenggara acara).

## 2. Tech stack (yang SEBENARNYA jalan)

| Bagian | Teknologi |
|---|---|
| Monorepo | Bun 1.4.1 workspaces (`frontend`, `packages/*`) — lihat catatan backend |
| Frontend | SvelteKit 2.70 + Svelte 5.57 (runes), Vite 8.2, TypeScript 5.9, Tailwind CSS **v4** (`@tailwindcss/vite`, `@import "tailwindcss"`, tanpa config JS), adapter Vercel runtime `nodejs22.x` |
| Backend | **Go 1.23** modular monolith — chi v5, pgx v5 (raw SQL, **tanpa ORM**), goose v3, golang-jwt v5, godotenv, google/uuid, x/crypto |
| DB | Supabase Postgres (pgx pool + Supavisor pooler: `SimpleProtocol` sengaja) |
| Auth | Supabase Auth (frontend, SSR cookie) **+** JWT HS256 sendiri di backend Go (`internal/identity`) |
| Storage | Supabase Storage (adapter) atau lokal (`backend/storage/`) |
| Payment | Duitku (HMAC-SHA256 invoice + callback; `internal/pay/duitku.go`) |
| Shared | `@ketuk/shared` — Zod 4.5.4 schema + tipe domain + konstanta |
| Lint/format | Biome 2.5.12 — tab, single quote, semicolon, width 100 |
| CI | `.github/workflows/ci.yml` — lint + check + build, `--frozen-lockfile` |

Skala: backend **48 file Go / ~6.600 LOC / 16 migrasi goose**; frontend
**105 file (.svelte+.ts) di src/**, ~73 endpoint, OpenAPI 3.0.3 (54 path).

## 3. Layout monorepo

```
ketuk/
├── frontend/          # SvelteKit (deploy Vercel)
├── backend/           # Go API — TANPA package.json (bukan workspace bun lagi)
│   ├── cmd/api/       # main.go: app.Build + sweepjob.Loop(15m) + graceful shutdown
│   ├── internal/      # 20 package domain (lihat §5)
│   ├── api/           # openapi.yaml (3.0.3, v0.2.0, server localhost:8080), auth.yaml
│   ├── storage/       # media lokal ( Supabase bila kredensial ada)
│   └── dist/          # LEGACY: hasil compile JS/Bun+Hono lama — abaikan
├── packages/shared/   # @ketuk/shared (types/schemas/constants/utils)
├── docs/              # ARCHITECTURE.md, DEVELOPMENT.md, prompts/00–03
├── .kiro/steering/    # cara-kerja.md: gaya kerja (lihat §10)
└── .mimosa/           # hook-state sesi tool (generated, jangan sentuh)
```

## 4. Command cheat-sheet

Berfungsi (dari root):
- `bun install` · `bun run dev:fe` (frontend only) · `bun run build:shared`
- `bun run lint` (biome) · `bun run format` · `bun run check` (tsc + svelte-check)

Backend (dari `backend/`): `make run` (`go run ./cmd/api`) · `make tidy` ·
`make test` (`go test ./...`) · `make bundle-api` (redocly). Verifikasi cepat:
`go build ./... && go vet ./...`.

**Script root yang STALE (backend bukan workspace bun lagi):** `bun run dev`,
`dev:be`, `db:generate`, `db:migrate`, `db:studio` (Drizzle) — semua merujuk
`--filter backend` yang tak punya package.json. `docs/DEVELOPMENT.md` bab
"Migrasi database" dan script table-nya juga kedaluwarsa (Drizzle → sekarang
goose di `internal/migrate/sql/`).

Port: backend dengar `:8080` default (`PORT` diutamakan atas `HTTP_ADDR`), tapi
`.env.example` menulis `BACKEND_PORT=3000` + `PUBLIC_API_URL=:3000` — **inkonsisten**, perlu saataman saat konfigurasi lokal.

## 5. Backend Go — peta `internal/`

Composition root: `internal/app/app.go` (load config → `db.Connect` →
`migrate.Up` bila `MIGRATE_ON_START` → wiring service → `httpserver.New`).
Port/adapter dipilih di sini: `pay.Disabled` vs Duitku, `storage.Local` vs
Supabase, `notify.Noop`.

- `apierr` — error standar `{Status, Code, Message}` snake_case; konstruktor
  per status (BadRequest/Unauthorized/Forbidden/NotFound/Conflict/Entitlement
  402/RateLimited 429/Unavailable 503).
- `auth` — **DEAD CODE** (JWT + middleware lama); digantikan identity+authctx.
- `authctx` — context key `userID`; dipakai middleware `s.auth`.
- `identity` — satu-satunya pola hexagonal (port `Repository` + sentinel);
  Register/Login (email **atau** username)/Refresh (rotasi sesi di DB)/Logout/
  ChangePassword (revoke semua sesi). JWT HS256 claim sub/typ/exp/iat/jti.
- `billing` — plan/subscription/entitlements (quota `-1` = unlimited),
  `Checkout` (order `saas-*`), `ConfirmCallback` (`SELECT … FOR UPDATE`),
  `SyncOrder`, `ExpirePending`.
- `catalog` — read-only `event_types` + `invitation_templates`; template premium
  ditandai `Locked` tanpa entitlement.
- `config` — godotenv; `DATABASE_URL` wajib; `JWT_SECRET`, ACCESS_TTL 15m /
  REFRESH_TTL 168h; `CORS_ORIGINS` CSV; `AUTH_RATE_LIMIT` (0 = mati);
  `CRON_SECRET` kosong = cron tak tercapai.
- `db` — `pgxpool` + `SimpleProtocol` (pooler Supavisor, sengaja).
- `gift` — metode bank/ewallet (langsung `paid`) vs duitku (`pending`, invoice
  `gift-*`); `ConfirmCallback`, `ExpirePending`.
- `httpserver` — semua handler + router chi (server.go, handlers.go, planner.go,
  extra.go, ratelimit.go, identity_dto.go). Upload multipart 16 MiB; callback
  Duitku `errors.Join(bill, gift)`; media resolver redirect-vs-stream.
- `httputil` — `JSON`, `Decode`, `Error` (errors.As apierr → `{"error":{…}}`).
- `invitation` — agregat terbesar (invitation.go, media.go, inviters.go,
  guestbook.go): CRUD, publish (slug `judul-uuidprefix`), RSVP publik by token,
  preview by slug, gallery (gambar only) & music (audio only, 1 per undangan
  via partial unique index).
- `link` — `resource_links` many-to-many undangan↔planner.
- `migrate` — goose v3 + `//go:embed sql/*.sql` (16 file `001…016`).
- `notify` — port `Notifier`, hanya `Noop` (dipanggil saat RSVP).
- `pay` — interface `Gateway` + `duitku.go` + `Disabled`.
- `planner` — planner (budget_mode event|daily), days, budget_categories/items,
  checklist; `sessions.go` tipe `Clock` ("HH:MM[:SS]" ↔ pgtype.Time).
- `storage` — port `Store`; `local.go` (disk, URL "" → stream) vs `supabase.go`
  (REST, URL publik → redirect).
- `sweepjob` — `RunOnce` (expire order/gift >90 menit, purge refresh sessions)
  + `Loop` ticker 15 menit.

**Router** (`httpserver/server.go`, chi v5): RequestID → RealIP → Logger →
Recoverer → Timeout 60s → CORS. Top-level: `GET /healthz` (503
`database_unavailable` bila DB down), `GET+POST /api/cron/sweep` (Bearer
`CRON_SECRET`). Grup `/v1`: rate-limited (register/login/refresh/logout, fixed-
window per IP/menit) | publik (event-types, templates, billing/plans,
`public/rsvp/{token}`, `public/invitations/{slug}`, `public/media/{id}`,
`public/payments/duitku/callback`) | ter-autentikasi (`s.auth` → `/me`,
`/auth/password`, billing/subscription/checkout, orders, invitations +
sub-resource, planners + sub-resource, links).

**DB (16 migrasi):** `users` (username unik-lower; migrasi 015/016 name→username),
`plans`, `plan_entitlements`, `subscriptions`, `orders` (kind
subscription/gift/vendor), `refresh_sessions`, `event_types`, `invitation_templates`
(tier free|premium), `invitations` (slug unik, status draft/published/archived,
content JSONB), `invitation_locations`, `invitation_guests` (rsvp_token, status
draft/invited/opened/accepted/declined, plus_ones), `invitation_inviters`,
`invitation_media` (kind gallery|music), `invitation_guestbook`, `planners`,
`planner_days`, `budget_categories`, `budget_items` (amount BIGINT),
`checklist_items`, `planner_sessions` (TIME), `resource_links`. Seed (006/013):
3 plan (free / starter 79rb / pro 199rb) + 5 entitlement, 7 event_type,
5 template.

**Konvensi backend:** 1 package per agregat, file kecil per sub-domain, test
co-located `*_test.go`; layer httpserver (decode/validate/DTO) → Service domain
→ raw SQL pgx; envelope `{"data":…}` (refresh → `tokens`; register/login →
`data`+`tokens`; callback → teks `SUCCESS`); DELETE → 204; PATCH pakai pointer
+ `fieldSent`/`fieldNull`; tx kanonik `pool.Begin` + `defer tx.Rollback`;
ownership guard `Get(owner, id)` → NotFound; mapping PG `23505`→Conflict,
`23503`→invalid_input, `ErrNoRows`/`RowsAffected()==0`→NotFound; komentar
paragraf Bahasa Indonesia menjelaskan **kenapa**.

## 6. Frontend SvelteKit — peta `frontend/src`

- `hooks.server.ts` — `sequence(handleSubdomain, handleSupabase)`.
  `handleSubdomain`: rewrite `slug.ketuk.id` → `/[slug]` (whitelist
  `BASE_HOSTNAMES`: ketuk.id, www.ketuk.id, localhost, 127.0.0.1). Supabase SSR
  client per-request + `locals.safeGetSession()` (validasi ke server).
- `lib/api/` — **satu-satunya jalan** ke backend Go: `client.ts` (`apiFetch`,
  `ApiRequestError`, tipe `ApiSuccess`/`ApiErrorBody`, `FetchCtx` =
  `{fetch, accessToken}`) + modul events/guests/planner/vendors/gifts/payments.
  `PUBLIC_API_URL` dibaca dinamis.
- `lib/supabase/` — `server.ts` (SSR cookie) & `client.ts` (browser, session di
  **cookie** bukan localStorage).
- `lib/components/ui/` (15 generik: Badge, Button, Card, ConfirmDialog,
  EmptyState, Input, Modal, PasswordInput, Select, Skeleton, Tabs, Textarea,
  Toast) & `lib/components/domain/` (18 bisnis: BudgetTable, ChecklistList,
  CountdownTimer, Envelope*, EventCard, GiftCard, GuestTable, ImageGallery,
  InvitationView, PhoneMockup, PlanCard, RsvpForm, Template*, VendorCard,
  Wish*). Barrel `index.ts` di tiap folder.
- `lib/stores/` — `toast.svelte.ts`, `confirm.svelte.ts` (Svelte 5 runes).
- `lib/server/email-domain.ts` — cek MX record, fail-open, cache 500 domain.
- `lib/utils/ics.ts`, `lib/data/{templates,master-templates}.ts`,
  `lib/icons/index.ts`.
- `routes/` — 4 grup + admin:
  - `(marketing)` — landing, `harga/`, `template/` + `[slug]`, `tentang/`.
  - `(auth)` — `masuk` & `daftar` (+page.server.ts, Zod dari `@ketuk/shared`,
    form action, pesan identik untuk email/password salah), `callback`.
  - `(app)` — guard login (`+layout.server.ts` + JIT profile dari
    `user_metadata`); `dashboard`, `undangan/` + `[id]`/{edit,tamu,ucapan} +
    `baru`, `planner/` + {budget,checklist,timeline}, `vendor/` + `[slug]`,
    `hadiah/` + `[id]` (query **view `gift_orders_safe`** langsung via Supabase).
  - `(public)` — `[slug]` (+page.server.ts: `isPublished` + cache-control
    agresif `public, max-age=60, s-maxage=3600, swr=86400`, sengaja tak sentuh
    session agar CDN-aman) + `[slug]/tamu/[guestSlug]`.
  - `admin/` — statis (`+layout.svelte`, `produk-undangan/`).
- `app.css` — Tailwind v4 `@theme`: palet `navy`(basis gelap), `coral`, `wine`,
  `cream`, `champagne`, `espresso` + warna per modul (undangan ungu, planner
  biru langit, vendor hijau, hadiah oranye). Font: Cormorant Garamond (serif),
  Montserrat (display), Inter (body). `prefers-reduced-motion` reset animasi.
- `app.d.ts` — `App.Locals = {supabase, safeGetSession}`; `App.PageData` =
  `{session?, user?, accessToken?}` (opsional sengaja).

**Pola data:** `+page.ts` universal → backend Go lewat `{fetch, accessToken}`
(error ditangkap → `{error:'Koneksi terputus…'}`, page tetap render);
`+page.server.ts` bila butuh `locals.supabase`; mutasi dominan client-side via
`lib/api` + `pushToast` + `invalidateAll()`; sebagian operasi langsung Supabase
browser mengandalkan RLS (`profiles`, `invitations`, `wishes`,
`gift_orders_safe`).

## 7. `@ketuk/shared`

`packages/shared/src/`: `index.ts` barrel; `types/` (user, event, invitation,
guest, planner, vendor, gift, payment); `schemas/` (auth, event, guest, gift,
payment — Zod 4.5.4, dipakai form action frontend); `constants/` (plans,
event-types, vendor-categories, provinces, payment-methods,
disposable-email-domains); `utils/` (slug, format, email, date). Build:
`tsc` ke `dist/` — **wajib `bun run build:shared` sebelum `check`/`build`**.

## 8. Environment

Root `.env` (lihat `.env.example`): `PUBLIC_SUPABASE_URL`,
`PUBLIC_SUPABASE_ANON_KEY`, `SUPABASE_SERVICE_ROLE_KEY` (backend only),
`DATABASE_URL` (pgx), `DUITKU_{MERCHANT_CODE,API_KEY,ENV,CALLBACK_URL,RETURN_URL}`,
`PUBLIC_APP_URL`, `PUBLIC_API_URL`, `BACKEND_PORT`, `NODE_ENV`, opsional
`CLOUDFLARE_{ZONE_ID,API_TOKEN}`. Backend: `backend/.env` (godotenv), plus
`JWT_SECRET`, `HTTP_ADDR`/`PORT`, `CORS_ORIGINS`, `AUTH_RATE_LIMIT`,
`MIGRATE_ON_START`, `CRON_SECRET`. Tanpa kredensial Duitku → checkout membalas
`503 payment_unavailable`. **Jangan commit `.env`.**

## 9. Gotcha & inkonsistensi yang harus diingat

1. **Migrasi stack Hono→Go**: README badge, `docs/ARCHITECTURE.md`, dan
   `docs/prompts/*` masih menyebut **Hono + Drizzle + adapter-node**. Realita:
   **Go + chi + goose**, adapter Vercel. Perbarui dokumen itu saat kesempatan.
2. **Script root stale**: `dev`/`dev:be`/`db:*` broken (backend tanpa
   package.json). Jalankan backend lewat `make run` di `backend/`.
3. **Port mismatch**: backend default `:8080` vs `.env.example` `3000`.
4. **`internal/auth` dead code** — jangan ikuti polanya; pakai identity+authctx.
5. **`backend/dist/`** artefak JS lama; `backend/vercel.json` = `{}`.
6. Referensi `backend/README.md` ke `docs/blueprint.md`, `AGENTS.md`,
   `docs/knowledge.md`, perintah `/knowledge` — **file/command itu tidak ada**.
7. `docs/prompts/README.md` membahas Claude Code workflow lama (Drizzle/Hono).
8. Branch: `main` (aktif), `origin/backend`, `origin/develop`.
9. Biome mengabaikan `*.svelte` — format komponen manual; biome juga ignore
   `.mimosa`, `build`, `dist`, `.svelte-kit`, `migrations`.

## 10. Gaya kerja (dari `.kiro/steering/cara-kerja.md`)

- **Kerjakan, bukan usulkan.** Ambil pendekatan paling masuk akal, jalankan,
  lalu laporkan singkat. Tebak maksud yang paling mungkin bila perintah ambigu;
  sebutkan asumsi di akhir. Jangan minta review per langkah.
- **Verifikasi sebelum bilang selesai**: `go build ./... && go vet ./...` +
  `go test ./...` (backend, dari `backend/`); `bun run check` + `bun run lint`
  (frontend/shared, dari root). Kalau tak bisa diverifikasi, katakan terus
  terang.
- **Konfirmasi dulu untuk aksi berisiko**: hapus banyak file, operasi git yang
  menulis remote (`push`, `reset --hard`, `clean -f`, force push, hapus branch),
  ubah/hapus data DB atau migrasi yang sudah jalan, sentuh kredensial/infra
  production, ubah auth/otorisasi.
- **Jangan `git commit`/`push`/`tag`** — commit manual oleh pemilik repo
  (Conventional Commits: `feat`/`fix`/`chore`/`docs`/`refactor`/`test`).
- Ikuti konvensi: Biome (tab, single quote), komentar Bahasa Indonesia
  menjelaskan *kenapa*, pola `@ketuk/shared`, tanpa komentar asal/umpan.
- Ringkas di akhir: apa yang berubah, asumsi yang diambil, apa yang masih perlu
  pemilik repo lakukan.

## 11. Cara memperbarui file ini

Saat selesai tugas yang mengubah struktur/fakta di atas (tambah package, pindah
stack, ganti port, tambah modul, ubah konvensi), **edit hanya bagian terkait di
file ini** — jangan regenerate utuh. Pertahankan gaya: bullet pendek, path
relatif, tabel untuk daftar. Kalau ragu apakah fakta masih berlaku, verifikasi
lewat grep singkat, bukan baca ulang seluruh folder.

<!-- antislop:start -->
## antislop
For UI, copy, people, mobile layout, or code comments work, load the antislop skill for the task:
- Core filter, always on: `antislop`
- Copy & text: `antislop-copywriting`
- People: `antislop-human`
- Mobile / responsive: `antislop-layoutmobile`
- Code comments: `antislop-code`
- UI / visual: `antislop-ui`
Before starting, ask the user when antislop applies: during the work, or after it is done.
<!-- antislop:end -->
