# Tahap 4 — Frontend Ketuk.id

## Aturan wajib

**JANGAN menjalankan `git commit`, `git push`, atau operasi git yang mengubah history.** Saya commit manual setelah review.

**JANGAN membuat file `.env` asli.** Hanya `.env.example`.

**Pakai Bun** untuk semua perintah package management.

**Jangan mendefinisikan ulang tipe yang sudah ada di `@ketuk/shared`.** Import dari sana.

---

## Konteks

Monorepo Ketuk.id sudah berdiri. Package `@ketuk/shared` berisi types, Zod schemas, dan konstanta. Backend Hono sudah jalan di `backend/` dengan endpoint yang terdokumentasi. Sekarang bangun `frontend/`.

Baca dulu:
- `docs/ARCHITECTURE.md`
- `packages/shared/src/index.ts`
- `backend/src/routes/index.ts` untuk melihat endpoint yang tersedia

---

## Prinsip arsitektur yang harus dipegang

**Halaman undangan publik terpisah total dari aplikasi.** Ini keputusan paling penting di frontend. Alasannya konkret: setelah host mem-blast link undangan ke grup WhatsApp, ratusan tamu bisa membuka dalam hitungan menit. Kalau setiap pembukaan menyentuh database, itu spike yang tidak perlu.

Konsekuensi teknisnya:

- Route `(public)/[slug]` di-render server-side dan mengirim header cache yang agresif
- HTML yang di-cache hanya berisi data statis undangan: nama, tanggal, lokasi, galeri, cerita
- Bagian dinamis — daftar ucapan, jumlah RSVP — diambil client-side setelah halaman ter-render, atau lewat form action
- Route ini **tidak** meng-import Supabase client, tidak butuh session, tidak menyentuh auth sama sekali

**Halaman undangan harus ringan.** Target: bisa dibuka nyaman di HP mid-range dengan koneksi 3G. Ini bukan optimasi kosmetik. Sebagian besar tamu di Indonesia membuka undangan dari WhatsApp, di perangkat yang tidak baru, dengan kuota terbatas. Kalau loading tiga detik, mereka menutup tab.

Praktiknya: minimalkan JavaScript di route publik, jangan import library berat, lazy-load galeri foto, dan jangan autoplay musik tanpa interaksi user.

**Dashboard boleh lebih berat.** Host yang sedang mengatur acaranya biasanya di desktop dengan koneksi bagus. Di sini interaktivitas lebih penting daripada ukuran bundle.

---

## Yang harus dibuat

### Dependencies

`bun add` di workspace frontend: `@supabase/supabase-js`, `@supabase/ssr`. Dev: `@sveltejs/kit`, `@sveltejs/adapter-node`, `@sveltejs/vite-plugin-svelte`, `svelte`, `vite`, `tailwindcss`, `@tailwindcss/vite`, `typescript`, `svelte-check`.

Pakai Svelte 5 dengan runes (`$state`, `$derived`, `$effect`, `$props`). Jangan pakai sintaks Svelte 4 (`export let`, store `$:`).

### Struktur route

```
src/routes/
├── +layout.svelte                    # root, import CSS global
├── +layout.server.ts                 # session Supabase
├── (marketing)/
│   ├── +layout.svelte
│   ├── +page.svelte                  # landing page
│   ├── harga/+page.svelte
│   ├── tentang/+page.svelte
│   └── template/+page.svelte         # galeri template undangan
├── (auth)/
│   ├── +layout.svelte
│   ├── masuk/+page.svelte
│   ├── daftar/+page.svelte
│   └── callback/+server.ts           # OAuth callback Supabase
├── (app)/
│   ├── +layout.svelte                # sidebar dashboard
│   ├── +layout.server.ts             # guard auth
│   ├── +page.svelte                  # "Mau melakukan apa?"
│   ├── undangan/
│   │   ├── +page.svelte              # daftar event
│   │   ├── baru/+page.svelte         # wizard buat event
│   │   └── [id]/
│   │       ├── +layout.svelte        # tab navigasi
│   │       ├── +page.svelte          # ringkasan
│   │       ├── edit/+page.svelte     # editor undangan
│   │       ├── tamu/+page.svelte     # kelola tamu
│   │       └── ucapan/+page.svelte   # moderasi ucapan
│   ├── planner/
│   │   ├── +page.svelte
│   │   ├── budget/+page.svelte
│   │   ├── checklist/+page.svelte
│   │   └── timeline/+page.svelte
│   ├── vendor/
│   │   ├── +page.svelte
│   │   └── [slug]/+page.svelte
│   └── hadiah/
│       ├── +page.svelte
│       └── [id]/+page.svelte
└── (public)/
    └── [slug]/
        ├── +page.server.ts           # SSR, cache agresif
        ├── +page.svelte              # halaman undangan
        └── tamu/[guestSlug]/
            ├── +page.server.ts
            └── +page.svelte          # undangan personal per tamu
```

### Caching halaman publik

Di `(public)/[slug]/+page.server.ts`:

```ts
setHeaders({
  'cache-control': 'public, max-age=60, s-maxage=3600, stale-while-revalidate=86400'
});
```

Angka ini disengaja: browser cache pendek (60 detik) supaya perubahan cepat terlihat kalau host membuka undangannya sendiri, CDN cache panjang (1 jam) untuk menahan spike, dan `stale-while-revalidate` supaya CDN menyajikan versi lama sambil mengambil yang baru di belakang.

Pastikan `+page.server.ts` ini tidak memanggil `locals.getSession()` atau apapun yang bergantung pada cookie — kalau iya, response tidak bisa di-cache bersama antar user.

### Handling subdomain

Di `src/hooks.server.ts`, tangani dua bentuk URL:
- `ketuk.id/budi-sinta` — path-based, default
- `budi-sinta.ketuk.id` — subdomain, untuk paket Lengkap

Deteksi subdomain dari header `host`, lalu rewrite ke route publik. Daftar domain dasar (`ketuk.id`, `www.ketuk.id`, `localhost:5173`) dikecualikan.

### Landing page

Ini halaman yang menentukan apakah orang mendaftar atau pergi. Yang harus ada:

**Hero** yang menjelaskan bahwa Ketuk bukan cuma undangan nikah. Jangan buka dengan "Buat Undangan" — buka dengan positioning yang lebih luas. Tampilkan variasi jenis acara secara visual.

**Bagian "Mau melakukan apa?"** dengan empat kartu modul. Ini menyampaikan sifat modular produk lebih baik daripada daftar fitur.

**Bagian dua sisi ekosistem** yang menjelaskan bahwa Ketuk melayani penyelenggara dan tamu. Ini pembeda utama dari kompetitor yang hanya melayani penyelenggara.

**Bagian jenis acara** yang menampilkan delapan atau lebih kategori dengan deskripsi singkat. Penting untuk SEO dan untuk meyakinkan orang yang bukan mau nikahan.

**Harga** yang mengambil data dari konstanta `PLANS` di shared, bukan hardcode.

**CTA penutup.**

Tulis copy dalam bahasa Indonesia yang wajar. Hindari terjemahan kaku dari template SaaS Inggris. "Mulai Gratis" bukan "Memulai Secara Gratis". "Bikin acara jadi mudah" bukan "Membuat Acara Menjadi Mudah".

### Halaman undangan publik

Ini produk yang sebenarnya dilihat ratusan orang. Struktur yang harus ada:

**Cover/amplop** — halaman pembuka dengan nama, tanggal, dan tombol "Buka Undangan". Kalau diakses lewat link personal tamu, tampilkan nama tamunya. Ini momen yang membuat undangan terasa dipersonalisasi.

**Hero** dengan nama dan foto utama.

**Profil mempelai atau tuan rumah** — untuk pernikahan tampilkan kedua mempelai beserta nama orang tua. Untuk jenis acara lain, sesuaikan. Jangan paksa struktur pernikahan ke ulang tahun.

**Waktu dan tempat** dengan tombol ke Google Maps dan tombol "Tambah ke Kalender" yang menghasilkan file `.ics`.

**Countdown** ke hari H.

**Galeri** dengan lazy loading. Jangan muat semua foto sekaligus.

**Cerita** (opsional, untuk pernikahan).

**RSVP** — form yang mengirim ke backend. Tampilkan konfirmasi setelah terkirim, dan izinkan mengubah jawaban.

**Ucapan** — form kirim ucapan dan daftar ucapan yang sudah masuk. Daftar ini diambil client-side setelah halaman ter-render, bukan di SSR, supaya HTML tetap bisa di-cache.

**Kirim hadiah** — daftar produk hadiah dan amplop digital dengan nomor rekening. Tombol salin nomor rekening yang memberi umpan balik visual saat berhasil.

**Footer** dengan atribusi Ketuk yang halus, bukan mencolok.

Musik latar boleh ada tapi harus dimulai dari interaksi user (tombol buka undangan), tidak autoplay. Sediakan kontrol untuk mematikan yang selalu terlihat.

### Desain visual

Buat sistem token di `src/app.css` pakai Tailwind v4 `@theme`. Arahan warna:

Basis coklat hangat (`#5D4037` dan variasinya, token `coffee`) dengan aksen terakota (`#9F4E2D`, token `terracotta`). Empat warna sekunder untuk masing-masing modul supaya user bisa mengenalinya secara visual: ungu untuk undangan, biru langit untuk planner, hijau untuk vendor, oranye untuk hadiah. Target user usia 25-50: kontras teks minimal AA (4.5:1), hindari warna saturasi tinggi di area besar.

Tipografi dua keluarga: Lora (serif, display/heading) dan Plus Jakarta Sans (sans, body). Kedua-duanya mendukung karakter Indonesia dan tersedia di Google Fonts. Jangan pakai lebih dari dua.

Halaman undangan boleh punya identitas visual sendiri yang berbeda dari dashboard — undangan harus terasa personal dan elegan, dashboard harus terasa efisien.

Hindari yang berikut karena membuat desain terasa generik: semua elemen dibungkus kartu dengan radius dan shadow yang sama; label ALL CAPS di atas setiap heading; animasi fade-up di setiap section saat scroll; tanda panah "→" di setiap tombol.

### Komponen

Buat komponen yang benar-benar dipakai ulang, jangan membuat abstraksi prematur. Yang jelas dibutuhkan:

`Button`, `Input`, `Select`, `Textarea`, `Modal`, `Toast`, `Card`, `Badge`, `Tabs`, `EmptyState`, `Skeleton`, `ConfirmDialog`.

Untuk domain: `EventCard`, `GuestTable`, `BudgetTable`, `ChecklistList`, `VendorCard`, `GiftCard`, `PlanCard`, `CountdownTimer`, `RsvpForm`, `WishForm`, `WishList`, `ImageGallery`.

Setiap komponen menerima props bertipe. Tidak ada `any`.

### State dan data fetching

Untuk data yang di-load per halaman, pakai `+page.server.ts` atau `+page.ts` sesuai kebutuhan. Untuk mutasi, pakai form actions supaya tetap jalan tanpa JavaScript kalau memungkinkan.

Untuk state UI global (toast, modal, session user), buat store sederhana di `src/lib/stores/`. Jangan pasang state management library.

Panggilan ke backend dibungkus di `src/lib/api/` dengan fungsi bertipe per domain. Jangan ada `fetch()` mentah tersebar di komponen.

### Auth

Pakai `@supabase/ssr` untuk session yang bekerja di server dan client. Metode login: magic link email dan Google OAuth. Simpan session di cookie, bukan localStorage.

Guard dashboard di `(app)/+layout.server.ts` — redirect ke `/masuk` kalau tidak ada session, dengan parameter `?next=` supaya setelah login kembali ke halaman yang dituju.

Route publik dan marketing tidak boleh memanggil session sama sekali.

### Accessibility dan mobile

Semua interaktif bisa diakses keyboard dengan focus ring yang terlihat. Kontras warna memenuhi WCAG AA. `prefers-reduced-motion` dihormati — kalau user mematikan animasi, matikan.

Semua halaman responsif dari 360px. Dashboard punya navigasi mobile yang layak, bukan sidebar yang dipaksa mengecil.

Halaman undangan diuji khusus di viewport kecil karena itu tempat mayoritas tamu membukanya.

### SEO dan sharing

Setiap halaman undangan punya meta tag Open Graph yang benar: judul dengan nama acara, deskripsi yang mengundang, dan gambar. Ini yang muncul saat link di-share di WhatsApp, dan preview yang bagus meningkatkan tingkat pembukaan.

Generate gambar OG dinamis kalau memungkinkan, atau pakai cover image event sebagai fallback.

---

## Standar kualitas

Setiap route yang mengambil data punya loading state dan error state yang ditulis dengan sengaja. Halaman kosong bukan kegagalan — itu kesempatan mengarahkan user melakukan sesuatu. "Belum ada undangan. Buat yang pertama." lebih baik daripada layar kosong.

Pesan error menjelaskan apa yang terjadi dan apa yang bisa dilakukan, bukan meminta maaf. "Koneksi terputus. Coba muat ulang." bukan "Maaf, terjadi kesalahan."

Nama tombol konsisten dengan hasilnya. Tombol "Publikasikan" menghasilkan notifikasi "Undangan dipublikasikan", bukan "Berhasil disimpan".

---

## Setelah selesai

1. `bun run check` — tidak boleh ada error TypeScript atau svelte-check
2. `bun run lint` — harus bersih
3. `bun run build` — build harus berhasil
4. Tampilkan `git status`

Lalu berhenti dan laporkan:
- Halaman apa saja yang sudah jadi dan mana yang masih placeholder
- Keputusan desain yang kamu ambil sendiri
- Apakah ada endpoint backend yang kamu butuhkan tapi belum ada

**Jangan commit.**
