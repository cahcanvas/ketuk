# Knowledge — ketuk-api

Peta kerja agen. Chat baru: `/knowledge` (command di `.cursor/commands/knowledge.md`, skill di `.cursor/skills/knowledge/`). Keputusan produk: `AGENTS.md`. Cetak biru: `docs/blueprint.md`. Peta endpoint: `docs/api.md`. Kontrak CRUD undangan (sudah dikode): `docs/contracts/undangan-crud.md`. Kontrak PATCH inviter + caption galeri (sudah dikode): `docs/contracts/undangan-patch-gallery-inviter.md`. Kontrak rapatkan planner bagian 1 + 2 (sudah dikode): `docs/contracts/planner-crud.md` — command `/planner-crud` dan `/planner-sessions` (keduanya historis), skill `.cursor/skills/planner-crud/` + `.cursor/skills/planner-sessions/`. Graph: `graphify-out/` (gitignored) — jika kosong: `graphify update .` lalu `graphify export wiki`. Query dulu: `graphify query "…"`.

## Lingkungan lokal & status verifikasi

Postgres lokal **sudah terpasang**: DBngin PostgreSQL 16.2, binary di `/Users/Shared/DBngin/postgresql/16.2/bin` (bukan Homebrew, bukan Postgres.app — `psql` tidak ada di PATH). Database `ketuk` sudah dibuat, `.env` sudah ada dari `.env.example`, goose di **versi 14**.

Alur penunjang sudah diuji end-to-end ke Postgres: register → rotasi refresh (replay token lama `401`) → logout → ganti password (semua sesi dicabut) → kuota planner `402` → `GET /links` + filter → checkout `503` → `GET /orders?kind=` → sync order asing `404` → rate limit `429` → sweep menandai order basi `expired`.

CRUD tipis undangan juga sudah diuji manual ke Postgres: PATCH lokasi parsial (`label` tetap saat hanya `venue_name` dikirim) → lokasi asing `404 location` → host lain `404 invitation` → DELETE `204` lalu `404` → PATCH tamu `invited` dengan `rsvp_token` tetap → `accepted`/`opened`/`declined`/`Invited` ditolak `400` → hapus tamu bikin `GET /v1/public/rsvp/{token}` jadi `404` → PATCH buku ucapan trim tanpa mengubah `created_at`, body `{}` no-op `200`, 2001 rune `400`.

Catatan uji manual: `GET /v1/templates?type=` memakai **slug** tipe acara (`wedding`), bukan UUID; template `free` pertama = `20000000-0000-0000-0000-000000000001`.

**Rundown sesi planner** sudah diuji manual ke Postgres: POST `"09:00"` + `"18:30:00"` → balasan `"HH:MM:SS"`, list urut jam, PATCH `{}` no-op `200`, `ends_at: null` + `pic_name: ""` mengosongkan keduanya, `ends_at` lebih awal → `400`, jam ngawur → `400 invalid_input`, hari asing → `404 day`, sesi asing → `404 session`, DELETE `204` lalu `404`, hapus hari menghapus sesinya (CASCADE).

**CRUD planner bagian 1** juga sudah diuji manual ke Postgres: PATCH parsial planner (title trim, `ends_at: null` tanpa menyentuh `starts_at`, `{}` no-op `200`, title kosong / `ends_at` < `starts_at` / `type_id` asing → `400`, `status`/`currency` diabaikan) → tanggal hari ganda `409` (POST dan PATCH) → label hari boleh kosong → item: `amount < 0` / kategori asing `404 category` / hari asing `404 day`, `paid_at` set → no-op → `null` tanpa mengubah `day_id`, `GET .../budget` total 30jt paid 5jt → checklist PATCH `200` + body, `due_at: null` → hapus kategori cascade itemnya (total turun ke 25jt) → hapus hari bikin `day_id` item **dan** checklist jadi `null` → `daily` → `event` mengosongkan `day_id` item saja (checklist tetap) dan `day_id` item ditolak `400` di mode `event` → kuota `402 max_active_planners`, DELETE archive `204` idempoten, hilang dari list tapi `GET /{id}` tetap terbaca `archived`, slot bebas untuk planner baru, `resource_links` tetap ada setelah archive → planner host lain `404 planner` (termasuk lewat jalur hari/checklist).

Catatan bentuk tanggal: `starts_at` / `ends_at` / `date` yang **baru dikirim** dibalas apa adanya (mis. `+07:00`), sedangkan yang dibaca dari kolom `DATE` keluar sebagai tengah malam `Z`. Komponen tanggalnya sama; hanya representasi offset yang beda antar response.

Belum diuji: `healthz` jalur `503` (butuh Postgres dimatikan) dan tes integrasi sebagai kode yang bisa diulang di CI.

## Wiring

`cmd/api/main.go`:

`config.Load` → `db.Connect` → (opsional `migrate.Up`) → `identity` / `pay` (Duitku atau `Disabled`) / `billing` / `catalog` / `invitation` / `planner` / `link` / `gift` → `httpserver.New`

`invitation.New(..., notify.Noop{}, storageDir)`. `link` satu-satunya yang kenal undangan **dan** planner. `gift` boleh import `invitation`; `invitation` tidak import `gift`/`pay`. Duitku hanya lewat `internal/pay`.

`internal/vendor` **belum ada**.

## Package map

| Package | Peran | Jangan |
|---|---|---|
| `internal/httpserver` | Chi routes, JWT, decode JSON, multipart, rate limit `/v1/auth/*` | Logic bisnis / SQL |
| `internal/identity` | Register, login, refresh (rotasi + `refresh_sessions`), logout, ganti password, `GET|PATCH /me` | Cek plan |
| `internal/auth` + `authctx` | JWT HS256, bcrypt, `UserID` di context | |
| `internal/catalog` | `event_types`, `invitation_templates` | |
| `internal/billing` | Plan, subscription, `Assert*`, checkout SaaS, `orders` (+`kind`, `SyncOrder`, `ExpirePending`) | Panggil Duitku langsung |
| `internal/pay` | Port `Gateway`; `Duitku` + `Disabled`; `CreateInvoice` / `VerifyCallback` / `CheckTransaction` | Dipakai dari undangan |
| `internal/invitation` | Undangan, lokasi, tamu, inviters, galeri, musik, buku ucapan, publish, preview slug, RSVP | Import `planner` / `gift` / `pay` |
| `internal/planner` | Budget event/daily, hari, kategori, item, checklist | Import `invitation` |
| `internal/link` | `resource_links` undangan↔planner | Cascade hapus produk |
| `internal/gift` | Metode amplop + pembayaran tamu | Import `planner` |
| `internal/notify` | Port `Notifier`; `Noop` | SMTP/WA dari domain |
| `internal/apierr` | Error JSON + `402 entitlement_denied` + `429 rate_limited` + `503 payment_unavailable` | |
| `internal/httputil` | `JSON` / `Error` / `Decode` | |
| `internal/migrate/sql` | Goose embed; seed `006_seed.sql`; `007` inviters/media, `008` gifts, `009` orders, `010` buku ucapan, `011` refresh_sessions, `012` orders.kind, `013` seed kuota planner, `014` planner_sessions | |

## Entitlement keys

Pakai `billing.Assert*` / `RequireInt` / `RequireFlag`. Jangan `if plan == "pro"`.

| Key | Fungsi |
|---|---|
| `max_active_invitations` | `AssertInvitationQuota` |
| `max_guests_per_invitation` | `AssertGuestQuota` |
| `premium_templates` | `AssertPremiumTemplate` (tier `premium` → 402) |
| `planner_daily` | `AssertPlannerDaily` (mode `daily`) |
| `max_active_planners` | `AssertPlannerQuota` (dipakai di `planner.Create`); seed di `013` |

Limit `< 0` = unlimited. Register selalu plan `free`. Tanpa kredensial Duitku, gateway = `pay.Disabled` → checkout/amplop Duitku `503 payment_unavailable`.

Tiga kas: `orders` subscription (SaaS) ≠ `gift_payments` (amplop) ≠ vendor (belum). Webhook: `POST /v1/public/payments/duitku/callback` (form, teks `SUCCESS`).

## HTTP

Kontrak: `api/openapi.yaml` (v0.2.0). Schema request/response **katalog + undangan digital + planner** sudah diisi dan **semua path undangan/planner punya handler** (termasuk PATCH lokasi/tamu/guestbook/inviter/caption galeri, dan PATCH/DELETE anak planner). Auth/billing masih path-only. Sukses `{"data": ...}` (auth: `data` + `tokens`). Callback Duitku bukan JSON `data`. Upload galeri/musik: multipart field `file`.

### Penunjang

- Monitoring: `GET /healthz` (API + Postgres, tanpa auth; DB mati → `503 database_unavailable`)
- Auth: `POST /v1/auth/register|login|refresh|logout`, `POST /v1/auth/password` (JWT), `GET|PATCH /v1/me`
- Sesi: refresh punya `jti` → baris `refresh_sessions`; refresh = rotasi (yang lama dicabut), logout cabut satu sesi, ganti password cabut semua
- Rate limit: `/v1/auth/*` dibatasi `AUTH_RATE_LIMIT`/IP/menit (default 10) → `429 rate_limited`
- Catalog: `GET /v1/event-types`, `GET /v1/templates?type=`
- Billing: `GET /v1/billing/plans`, `GET /v1/billing/subscription`, `POST /v1/billing/checkout` body `{plan}`
- Pesanan: `GET /v1/orders?kind=`, `POST /v1/orders/{id}/sync` (rekonsiliasi ke Duitku)
- Duitku: `POST /v1/public/payments/duitku/callback`
- Link: `GET /v1/links?invitation_id=&planner_id=`, `POST /v1/links`, `DELETE /v1/links/{id}`
- Sweep tiap 15 menit di `cmd/api`: order/amplop `pending` > 90 menit → `expired`, sesi kedaluwarsa dihapus
- Perilaku yang bikin FE bingung: `POST /v1/billing/checkout` menyisipkan baris `orders` **sebelum** memanggil Duitku, jadi checkout yang gagal (`503`) tetap muncul di `GET /v1/orders` berstatus `failed`

### Undangan

CRUD + publish + template + locations + guests + **inviters** (turut mengundang, patch/hapus) + **gallery** (patch caption) + **music** + **guestbook** (buku ucapan) + gift-methods + gifts. Publik: `GET/POST /v1/public/rsvp/{token}`, `POST .../gifts`, `POST .../guestbook`, `GET /v1/public/invitations/{slug}` (preview, bukan RSVP), `GET/POST .../guestbook`, `GET /v1/public/media/{id}`.

### Planner

CRUD penuh: planner (PATCH parsial + DELETE archive), hari, **sesi rundown**, kategori, item budget, checklist — semua punya PATCH/DELETE. `GET .../budget` = ringkasan **sekaligus** daftar item (tidak ada `GET .../budget/items`). Vendor di planner belum.

Sesi rundown ada di bawah hari: `GET/POST /v1/planners/{id}/days/{dayId}/sessions`, `PATCH/DELETE .../sessions/{sessionId}`. Jam dinding (`TIME`), bukan RFC3339 — request `"HH:MM"` atau `"HH:MM:SS"`, response selalu `"HH:MM:SS"` lewat tipe `planner.Clock` (punya `MarshalJSON`; jembatan DB `pgtype.Time` di `clockFromPg` / `Clock.pg`). Hapus hari menghapus sesinya (CASCADE), beda dari item/checklist yang hanya kehilangan `day_id`.

Bentuk parsial planner beda dari undangan: `type_id`, `starts_at`, `ends_at`, `day_id`, `paid_at`, `due_at` pakai `json.RawMessage` supaya key absen (tidak diubah) beda dari `null` (dikosongkan) — helper `fieldSent` / `fieldNull` / `patchUUID` / `patchTime` di `internal/planner`.

Aturan yang mudah terlewat: `day_id` budget item hanya boleh di mode `daily` (mode `event` → 400 `invalid_input`); checklist boleh menempel hari di mode apa pun; turun `daily` → `event` mengosongkan `day_id` budget item saja. Tanggal hari ganda → `409 conflict` (POST dan PATCH). Hapus kategori **cascade** ke itemnya. `PATCH .../checklist/{itemId}` sekarang `200` + body (dulu `204`).

Lokasi: `maps_url` tempelan. RSVP jalur utama = token per tamu, bukan slug terbuka. Preview: `GET /v1/public/invitations/{slug}`.

## CRUD yang sudah rapat vs sisa tipis

Undangan — kontrak `docs/contracts/undangan-crud.md` **sudah dikode** (OpenAPI tanpa `x-implemented`):

- `PATCH|DELETE /v1/invitations/{id}/locations/{locationId}` — `invitation.PatchLocation` / `DeleteLocation`, body parsial (`invitation.LocationPatch`), tanpa geocode
- `PATCH|DELETE /v1/invitations/{id}/guests/{guestId}` — `invitation.PatchGuest` / `DeleteGuest`; `status` host hanya `draft` | `invited` (`validateHostGuestStatus`), `rsvp_token` tidak diganti, tanpa notify
- `PATCH /v1/invitations/{id}/guestbook/{entryId}` — `invitation.PatchGuestbook` + `applyGuestbookPatch`, validasi sama `validateGuestbook`, `created_at` tetap

Kontrak `docs/contracts/undangan-patch-gallery-inviter.md` **sudah dikode** juga:

- `PATCH /v1/invitations/{id}/inviters/{inviterId}` — `invitation.PatchInviter` + `applyInviterPatch`; `name` kosong → 400, `role` kosong = hapus peran, `sort_order` ganti urutan
- `PATCH /v1/invitations/{id}/gallery/{mediaId}` — `invitation.PatchGallery` + `applyGalleryCaption`; hanya `caption` (trim, `""` = clear), SELECT/UPDATE difilter `kind='gallery'` sehingga id musik → 404 `media`; file di `./storage` tidak disentuh

Pola sama untuk semua: `Get(owner, invitationID)` dulu (bukan milik host → 404 `invitation`), baris anak difilter `id + invitation_id` (salah → 404 `location`/`guest`/`guestbook`/`inviter`/`media`), body `{}` = no-op 200, DELETE 204 tanpa body.

Undangan yang masih sengaja belum: PATCH gift-method, ganti file galeri, reorder galeri, caption musik, cerita pasangan.

Planner — kontrak `docs/contracts/planner-crud.md` **bagian 1 sudah dikode**:

- `PATCH /v1/planners/{id}` — `planner.Patch` + `applyPlannerPatch`; `SetMode` tetap yang memiliki aturan entitlement `daily` dan pembersihan `daily` → `event`; `status`/`currency` tidak lewat PATCH
- `DELETE /v1/planners/{id}` — `planner.Archive`, `status=archived`, idempoten, membebaskan slot `max_active_planners`; `resource_links` tetap
- `PATCH|DELETE .../days/{dayId}` — `PatchDay` / `DeleteDay`, unique `(planner_id, day_date)` → `409 conflict`
- `PATCH|DELETE .../categories/{categoryId}` — `PatchCategory` / `DeleteCategory` (cascade item)
- `PATCH|DELETE .../budget/items/{itemId}` — `PatchItem` / `DeleteItem`; SELECT/DELETE difilter lewat `budget_categories.planner_id`
- `PATCH|DELETE .../checklist/{itemId}` — `PatchCheck` (ganti `ToggleCheck`, kini `200` + body, 0 row → 404) / `DeleteCheck`

POST yang ikut dirapatkan: `AddDay` (`23505` → 409), `AddItem` (`amount >= 0`, `category_id` wajib milik planner, `day_id` sesuai mode), `AddCheck` (baca `due_at`/`sort_order`, `day_id` wajib milik planner).

Bagian 2 kontrak **sudah dikode** juga — migrasi `014_planner_sessions.sql`, domain `internal/planner/sessions.go` (`ListSessions` / `AddSession` / `PatchSession` / `DeleteSession` + `applySessionPatch`, tes `session_patch_test.go`), handler `listSessions` / `addSession` / `patchSession` / `deleteSession`. Ownership: `Get(owner, plannerID)` → `assertDay` → baris sesi difilter `id + day_id` (hari asing → 404 `day`, sesi beda hari → 404 `session`). `ends_at` tiga-state lewat `patchClock`; `pic_name` teks bebas (bukan vendor/user id).

Sisa planner: unarchive, mata uang selain IDR.

Penunjang yang sengaja belum: reset password lewat email dan verifikasi email (butuh notify nyata), hapus akun (perlu keputusan soft delete vs cascade ke `orders`), cancel/downgrade langganan, baris `gift`/`vendor` di `orders` (kolom `kind` sudah ada, penulisnya belum).

## Vendor

Future. Jangan kode sampai diminta. Rencana: `docs/blueprint.md` pohon fitur §3.

## Graphify

Lokal saja (`graphify-out/` di `.gitignore`). Setelah ubah kode: `graphify update .` lalu `graphify export wiki` (wiki ditimpa).
