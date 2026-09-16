# ketuk-api

Backend Go untuk ketuk.id. Modular monolith, tanpa Docker. Checkout host dan amplop hadiah lewat adapter Duitku; Vendor belum ada.

Cetak biru first release + future: [`docs/blueprint.md`](docs/blueprint.md). Keputusan terkunci: [`AGENTS.md`](AGENTS.md). Peta paket: [`docs/knowledge.md`](docs/knowledge.md). Agen baru: `/knowledge`.

## Syarat

- Go 1.23+
- Postgres lokal (DBngin / Homebrew / Postgres.app)

```bash
createdb ketuk
cp .env.example .env
make run
```

Di mesin dev ini Postgres jalan lewat **DBngin 16.2** dan `psql`/`createdb` tidak ada di PATH — pakai `/Users/Shared/DBngin/postgresql/16.2/bin/createdb -h localhost -U postgres ketuk`.

- Monitoring: `GET http://localhost:8080/healthz` (Postgres mati → `503 database_unavailable`)
- Kontrak: `api/openapi.yaml`

Response sukses dibungkus `{"data": ...}` kecuali auth yang juga mengembalikan `tokens`.

## Endpoint utama

| Area | Path |
|---|---|
| Auth | `POST /v1/auth/register` `login` `refresh` `logout` · `POST /v1/auth/password` · `GET/PATCH /v1/me` |
| Catalog | `GET /v1/event-types` · `GET /v1/templates?type=` |
| Billing | `GET /v1/billing/plans` · `GET /v1/billing/subscription` · `POST /v1/billing/checkout` |
| Pesanan | `GET /v1/orders?kind=` · `POST /v1/orders/{id}/sync` |
| Undangan | `/v1/invitations` + locations, guests, inviters, gallery, music, guestbook, gift-methods, publish, template |
| RSVP tamu | `GET/POST /v1/public/rsvp/{token}` · `+/gifts` `+/guestbook` |
| Undangan publik | `GET /v1/public/invitations/{slug}` (preview) · `GET /v1/public/media/{id}` |
| Planner | `/v1/planners` + days, categories, budget items, checklist (PATCH/DELETE; `DELETE` planner = archive) |
| Link opsional | `GET /v1/links` · `POST /v1/links` · `DELETE /v1/links/{id}` |

Paket default saat register: `free`. `/v1/auth/*` dibatasi `AUTH_RATE_LIMIT` per IP per menit (default 10) → `429 rate_limited`. Adapter bayar: **Duitku** lewat `internal/pay`; tanpa kredensial di `.env` checkout membalas `503 payment_unavailable`. Pohon fitur undangan/planner/vendor: [`docs/blueprint.md`](docs/blueprint.md).
