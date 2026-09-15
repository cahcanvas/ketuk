# Panduan Development Ketuk.id

Panduan praktis untuk developer baru yang ingin menjalankan proyek ini secara lokal.

## Prasyarat

Sebelum mulai, pastikan sudah terpasang:

- **Bun** — package manager dan runtime development. Install dari [bun.sh](https://bun.sh).
- **Node.js 20+** — runtime yang dipakai saat production dan oleh beberapa tooling. Versi yang disarankan tercantum di `.nvmrc`.
- **Akun Supabase** — untuk database Postgres, autentikasi, dan storage. Buat project baru di [supabase.com](https://supabase.com) dan catat URL project, anon key, service role key, serta connection string database-nya.
- **Akun Duitku sandbox** — untuk integrasi pembayaran. Daftar di [duitku.com](https://duitku.com), aktifkan mode sandbox, dan catat merchant code serta API key sandbox.

## Setup dari clone sampai jalan

1. Clone repository dan masuk ke foldernya.
2. Jalankan `bun install` di root. Ini akan menginstal dependency untuk semua workspace (`frontend`, `backend`, `packages/shared`) sekaligus.
3. Salin `.env.example` menjadi `.env`, lalu isi semua variabel dengan kredensial Supabase dan Duitku sandbox milikmu. Jangan pernah commit file `.env`.
4. Jalankan migrasi database (lihat bagian [Migrasi database](#migrasi-database) di bawah).
5. Jalankan `bun run dev` untuk menjalankan frontend dan backend sekaligus, atau `bun run dev:fe` / `bun run dev:be` untuk menjalankan salah satunya saja.

## Penjelasan script

Semua script berikut dijalankan dari root repository dan didelegasikan ke sub-package lewat `bun run --filter`.

| Script | Fungsi |
|---|---|
| `bun run dev` | Menjalankan mode development untuk semua package sekaligus (frontend + backend). |
| `bun run dev:fe` | Menjalankan mode development khusus `frontend`. |
| `bun run dev:be` | Menjalankan mode development khusus `backend`. |
| `bun run build` | Build production untuk semua package. |
| `bun run check` | Menjalankan type-check (`tsc --noEmit`) di semua package. |
| `bun run lint` | Menjalankan Biome untuk memeriksa lint dan format di seluruh repo. |
| `bun run format` | Menjalankan Biome untuk merapikan format kode secara otomatis. |
| `bun run db:generate` | Menghasilkan file migrasi Drizzle dari perubahan schema di `backend`. |
| `bun run db:migrate` | Menjalankan migrasi yang sudah digenerate ke database Supabase. |
| `bun run db:studio` | Membuka Drizzle Studio untuk melihat dan mengedit data secara visual. |

## Menjaga dependency konsisten dengan CI

Gunakan versi Bun pada `packageManager` di `package.json` root (saat ini `bun@1.4.1`). GitHub Actions membaca versi dari file yang sama.

Dependency yang sebelumnya memakai `latest` sudah dikunci ke versi spesifik. TypeScript disamakan ke `5.9.3` di semua workspace. Saat menambah atau memperbarui dependency:

1. Ubah versi dependency secara eksplisit pada package yang terkait.
2. Jalankan `bun install` dari root untuk memperbarui `bun.lock`.
3. Verifikasi dengan `bun install --frozen-lockfile`.
4. Commit perubahan `package.json` dan `bun.lock` bersama-sama.

CI tetap menggunakan `--frozen-lockfile` agar dependency yang dipasang mengikuti lockfile yang sudah direview.

## Migrasi database

Schema database dikelola dengan Drizzle ORM di dalam `backend`. Setelah mengubah schema:

1. Jalankan `bun run db:generate` untuk menghasilkan file migrasi baru berdasarkan perubahan schema.
2. Review file migrasi yang dihasilkan sebelum menjalankannya.
3. Jalankan `bun run db:migrate` untuk menerapkan migrasi ke database yang alamatnya tercantum di `DATABASE_URL`.

Detail schema dan RLS policy dijelaskan lebih lanjut saat `backend` diimplementasikan (lihat `docs/prompts/02-BACKEND.md`).

## Konvensi commit

Repository ini memakai [Conventional Commits](https://www.conventionalcommits.org/). Format pesan commit:

```
<type>(<scope opsional>): <deskripsi singkat>
```

Tipe yang umum dipakai:

- `feat` — fitur baru
- `fix` — perbaikan bug
- `chore` — perubahan tooling, konfigurasi, atau housekeeping yang tidak mengubah perilaku aplikasi
- `docs` — perubahan dokumentasi saja
- `refactor` — perubahan kode yang tidak menambah fitur atau memperbaiki bug
- `test` — menambah atau memperbaiki test

Contoh: `feat(backend): add RSVP endpoint`, `fix(frontend): correct invitation slug validation`.

## Alur Git

Commit dan push dilakukan **manual** oleh developer, tidak lewat AI assistant. Claude Code atau tool sejenis hanya menulis file, menjalankan build, dan menjalankan test — ia tidak pernah menjalankan `git commit`, `git push`, atau `git tag`.

Alur kerja standar setelah selesai mengerjakan sebuah tahap:

```bash
git status
git diff
bun run check          # pastikan tidak ada error TypeScript
git add .
git commit -m "..."     # pesan mengikuti Conventional Commits
git push origin main
```

Selalu review `git status` dan `git diff` sebelum commit untuk memastikan tidak ada perubahan yang tidak diinginkan ikut ter-commit.
