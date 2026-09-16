# API Ketuk — peta endpoint

Sukses JSON: `{"data": ...}` (auth register/login: `data` + `tokens`). Error: `{"error":{"code","message"}}`. Kuota: `402 entitlement_denied`. Terlalu sering: `429 rate_limited`. Checkout tanpa Duitku: `503 payment_unavailable`.

Callback Duitku **bukan** JSON `data` — form-urlencoded, response teks `SUCCESS`.

Upload galeri/musik: `multipart/form-data` field `file` (caption galeri opsional).

Vendor **belum** punya route. Rencana di akhir dokumen.

---

## Lintas modul (bukan salah satu dari tiga fitur utama)

| Method | Path | Auth | Keterangan |
|---|---|---|---|
| GET | `/healthz` | — | Monitoring sistem; Postgres mati → `503 database_unavailable` |
| POST | `/v1/auth/register` | — | Host, plan `free` |
| POST | `/v1/auth/login` | — | JWT |
| POST | `/v1/auth/refresh` | — | Rotasi; refresh lama langsung dicabut |
| POST | `/v1/auth/logout` | — | Body `{refresh_token}`; cabut satu sesi |
| POST | `/v1/auth/password` | JWT | Body `{current_password, new_password}`; cabut semua sesi |
| GET | `/v1/me` | JWT | Dashboard tipis |
| PATCH | `/v1/me` | JWT | Ubah `name` (email belum, butuh verifikasi) |
| GET | `/v1/event-types` | — | Katalog tipe acara |
| GET | `/v1/templates` | — | Template; `?type=` |
| GET | `/v1/billing/plans` | — | Paket SaaS |
| GET | `/v1/billing/subscription` | JWT | Paket + entitlement |
| POST | `/v1/billing/checkout` | JWT | Body `{plan}` → Duitku `payment_url` |
| GET | `/v1/orders` | JWT | Pesanan host; `?kind=subscription\|gift\|vendor` |
| POST | `/v1/orders/{id}/sync` | JWT | Rekonsiliasi order `pending` ke Duitku |
| POST | `/v1/public/payments/duitku/callback` | — | Webhook Duitku (kas SaaS + gift) |
| GET | `/v1/links` | JWT | List tautan; filter `invitation_id` / `planner_id` |
| POST | `/v1/links` | JWT | Taut undangan↔planner |
| DELETE | `/v1/links/{id}` | JWT | Hapus tautan |

`/v1/auth/*` dibatasi `AUTH_RATE_LIMIT` permintaan per IP per menit (default 10, `0` = mati) → `429 rate_limited`.

Reset password lewat email belum ada (notify masih `Noop`).

---

## 1. Undangan

| Method | Path | Auth | Keterangan |
|---|---|---|---|
| GET | `/v1/invitations` | JWT | List |
| POST | `/v1/invitations` | JWT | Buat (kuota + template) |
| GET | `/v1/invitations/{id}` | JWT | Detail |
| PATCH | `/v1/invitations/{id}` | JWT | Title/content/tanggal |
| DELETE | `/v1/invitations/{id}` | JWT | Archive |
| POST | `/v1/invitations/{id}/publish` | JWT | Isi slug |
| PATCH | `/v1/invitations/{id}/template` | JWT | Sebelum publish |
| GET/POST | `/v1/invitations/{id}/locations` | JWT | Acara & lokasi (`maps_url`) |
| PATCH/DELETE | `/v1/invitations/{id}/locations/{locationId}` | JWT | Body parsial; tanpa geocode. 404 `location` |
| GET/POST | `/v1/invitations/{id}/guests` | JWT | Daftar tamu + `rsvp_token` |
| PATCH/DELETE | `/v1/invitations/{id}/guests/{guestId}` | JWT | `status` host: `draft`\|`invited`; `rsvp_token` tetap; hapus = token mati |
| GET/POST | `/v1/invitations/{id}/inviters` | JWT | Turut mengundang |
| PATCH | `/v1/invitations/{id}/inviters/{inviterId}` | JWT | Body parsial `name`/`role`/`sort_order`; `role` kosong = hapus peran |
| DELETE | `/v1/invitations/{id}/inviters/{inviterId}` | JWT | Hapus pengundang |
| GET/POST | `/v1/invitations/{id}/gallery` | JWT | Galeri acara (upload) |
| PATCH | `/v1/invitations/{id}/gallery/{mediaId}` | JWT | Hanya `caption`; file tidak diganti; id musik → 404 `media` |
| DELETE | `/v1/invitations/{id}/gallery/{mediaId}` | JWT | Hapus foto |
| GET/PUT/DELETE | `/v1/invitations/{id}/music` | JWT | Musik saat undangan dibuka |
| GET/POST | `/v1/invitations/{id}/gift-methods` | JWT | Amplop: `bank` / `ewallet` / `duitku` |
| DELETE | `/v1/invitations/{id}/gift-methods/{methodId}` | JWT | Hapus metode |
| GET | `/v1/invitations/{id}/gifts` | JWT | Riwayat amplop masuk |
| GET | `/v1/invitations/{id}/guestbook` | JWT | Buku ucapan (moderation) |
| PATCH | `/v1/invitations/{id}/guestbook/{entryId}` | JWT | Ubah `name`/`message`; `created_at` tetap |
| DELETE | `/v1/invitations/{id}/guestbook/{entryId}` | JWT | Hapus ucapan |
| GET | `/v1/public/rsvp/{token}` | — | Halaman tamu + inviters, gallery, music, guestbook, gifts |
| POST | `/v1/public/rsvp/{token}` | — | RSVP `{status, message, plus_ones}` |
| POST | `/v1/public/rsvp/{token}/gifts` | — | Tamu kirim amplop |
| POST | `/v1/public/rsvp/{token}/guestbook` | — | Tamu isi buku ucapan |
| GET | `/v1/public/invitations/{slug}` | — | Preview undangan (bukan RSVP) |
| GET/POST | `/v1/public/invitations/{slug}/guestbook` | — | Baca / isi buku ucapan |
| GET | `/v1/public/media/{id}` | — | File galeri/musik |

---

## 2. Planner

| Method | Path | Auth | Keterangan |
|---|---|---|---|
| GET | `/v1/planners` | JWT | List `status=active` |
| POST | `/v1/planners` | JWT | Buat (kategori kegiatan opsional) |
| GET | `/v1/planners/{id}` | JWT | Detail (termasuk archived) |
| PATCH | `/v1/planners/{id}` | JWT | Parsial: `title`, `budget_mode`, `type_id`, `starts_at`, `ends_at` |
| DELETE | `/v1/planners/{id}` | JWT | Archive (bukan hard delete); idempoten |
| GET/POST | `/v1/planners/{id}/days` | JWT | Hari / rundown kasar; tanggal ganda → `409 conflict` |
| PATCH/DELETE | `/v1/planners/{id}/days/{dayId}` | JWT | Ubah / hapus hari (`day_id` anak jadi `null`; sesi ikut terhapus) |
| GET/POST | `/v1/planners/{id}/days/{dayId}/sessions` | JWT | Sesi rundown jam + PIC; urut `starts_at`, `sort_order` |
| PATCH/DELETE | `/v1/planners/{id}/days/{dayId}/sessions/{sessionId}` | JWT | Ubah (200 + body) / hapus sesi |
| GET/POST | `/v1/planners/{id}/categories` | JWT | Kategori budget |
| PATCH/DELETE | `/v1/planners/{id}/categories/{categoryId}` | JWT | Ubah / hapus kategori — **hapus ikut item di dalamnya** |
| GET | `/v1/planners/{id}/budget` | JWT | Ringkasan + item (tidak ada list item terpisah) |
| POST | `/v1/planners/{id}/budget/items` | JWT | Tambah item |
| PATCH/DELETE | `/v1/planners/{id}/budget/items/{itemId}` | JWT | Ubah / hapus item; `paid_at` diisi FE |
| GET/POST | `/v1/planners/{id}/checklist` | JWT | Daftar tugas (`due_at`, `sort_order` ikut dibaca) |
| PATCH/DELETE | `/v1/planners/{id}/checklist/{itemId}` | JWT | Ubah (200 + body) / hapus tugas |

Bentuk parsial: key absen = tidak diubah. `type_id`, `starts_at`, `ends_at`, `day_id`, `paid_at`, `due_at` menerima `null` untuk dikosongkan. 404 anak memakai nama resource: `day`, `session`, `category`, `budget_item`, `checklist`.

`day_id` pada budget item hanya boleh saat `budget_mode` = `daily` (mode `event` → 400); checklist boleh menempel hari di mode apa pun. Turun `daily` → `event` mengosongkan `day_id` budget item, checklist tetap.

Sesi rundown memakai **jam dinding**, bukan RFC3339: `starts_at` / `ends_at` dikirim `"HH:MM"` atau `"HH:MM:SS"` dan selalu dibalas `"HH:MM:SS"` (kolom `TIME`; tanggalnya ada di hari induk). `ends_at` boleh `null` (sesi tanpa jam selesai) dan lebih awal dari `starts_at` → 400. `pic_name` teks bebas — bukan id user atau vendor.

Vendor di dalam planner **belum** ada.

---

## 3. Vendor — rencana (belum diimplementasi)

Tidak ada route `/v1/vendors*` sekarang. Rencana produk:

| Method | Path (rencana) | Auth | Keterangan |
|---|---|---|---|
| GET | `/v1/vendors` | — | Cari / browse katalog jasa (`q`, `type`, kategori) |
| GET | `/v1/vendors/{id}` | — | Profil vendor |
| GET | `/v1/vendors/{id}/services` | — | Layanan + harga |
| POST | `/v1/vendors/{id}/requests` | JWT | Hubungi / request vendor |
| GET | `/v1/vendor-requests` | JWT | Inbox request milik host |
| GET | `/v1/vendor/me` | JWT vendor | Dashboard vendor (nanti) |
| POST | `/v1/vendor/services` | JWT vendor | Kelola katalog sendiri |

Filter kebutuhan/acara memakai `event_types`, **bukan** Event instance sebagai parent. Tautan ke undangan/planner opsional (pola `resource_links`). Pembayaran ke vendor ≠ checkout paket Ketuk ≠ amplop tamu.
