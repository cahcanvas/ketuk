# Ketuk.id

Platform modular untuk segala urusan acara di Indonesia seperti pernikahan, khitanan, aqiqah, ulang tahun, wisuda, reuni, syukuran, hingga corporate event. *Satu tempat untuk segala urusan acara.*

Empat modul yang bisa dibeli terpisah: **Undangan** (undangan digital & RSVP), **Planner** (budget, checklist, timeline, daftar tamu), **Vendor** (marketplace katering, dekorasi, fotografer, WO, MUA), dan **Hadiah** (kirim hampers, bouquet, kue ke penyelenggara acara).

![Go](https://img.shields.io/badge/backend-Go-00ADD8?logo=go)
![SvelteKit](https://img.shields.io/badge/frontend-SvelteKit-FF3E00?logo=svelte)
![Bun](https://img.shields.io/badge/runtime-Bun-000000?logo=bun)
![Supabase](https://img.shields.io/badge/database-Supabase-3ECF8E?logo=supabase)
![TypeScript](https://img.shields.io/badge/lang-TypeScript-3178C6?logo=typescript)
![License: MIT](https://img.shields.io/badge/license-MIT-blue)

## Quick start

Frontend (Bun + SvelteKit):

```bash
bun install
cp .env.example .env   # isi kredensial Supabase & Duitku sandbox
bun run dev
```

Backend (Go API):

```bash
cd backend
cp .env.example .env   # JWT_SECRET, DATABASE_URL, kredensial Duitku
make run               # go run ./cmd/api, dengar :8080
```

Dokumentasi lebih lanjut: [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) untuk keputusan arsitektur, dan [`docs/DEVELOPMENT.md`](docs/DEVELOPMENT.md) untuk panduan setup dan kontribusi.

## Struktur folder

```
ketuk/
├── frontend/       # SvelteKit app (Bun)
├── backend/        # Go API server (chi + pgx + goose)
├── packages/       # Package bersama (types & konstanta)
├── docs/           # Dokumentasi & prompt tahap pengerjaan
└── .github/        # CI workflows
```

## Lisensi

[MIT](LICENSE)
