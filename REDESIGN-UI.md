# REDESIGN-UI.md — Ketuk.id menjadi platform undangan acara umum

> **Dokumen ini adalah prompt implementasi, bukan diskusi.** Agen yang
> membaca ini wajib mengimplementasikan apa yang ditulis di sini, dalam urutan
> yang ditentukan, dan memverifikasi hasilnya di §16 sebelum bilang selesai.
> Kalau ada bagian yang ambigu, ambil interpretasi paling masuk akal dan catat
> asumsinya di akhir laporan — **jangan tanya per langkah.**

## 0. Baca dulu, sebelum menyentuh kode

Sebelum mulai, baca file-file ini (wajib, urutannya ditentukan):

1. `AGENTS.md` — memori proyek. Fakta stack, layout monorepo, konvensi
   (Biome tab/single-quote, komentar Bahasa Indonesia menjelaskan *kenapa*,
   jangan `git commit`).
2. `.opencode/skills/antislop/SKILL.md` — filter inti.
3. `.opencode/skills/antislop-copywriting/SKILL.md` — semua teks yang ditulis.
4. `.opencode/skills/antislop-layoutmobile/SKILL.md` — semua layout.
5. `.opencode/skills/antislop-ui/SKILL.md` — semua visual.
6. `frontend/src/app.css` — token desain yang sudah ada. **Pakai, jangan
   tambah palet baru.**

**Mode antislop: DURING.** Terapkan aturannya sambil menulis, bukan saat
"review akhir". Konsekuensi konkret: tidak ada kalimat penjelas generik, tidak
ada fitur yang dipajang tanpa fungsi, tidak ada breakpoints yang hanya bagus
di desktop, tidak ada CTA kata-kata pemasaran kosong.

## 1. Apa yang berubah dan kenapa

Ketuk.id sekarang hanya terlihat sebagai **situs undangan pernikahan**.
Buktinya nyata di kode, bukan asumsi:

- `frontend/src/routes/(marketing)/+page.svelte:86` — judul dan meta
  description menyebut "Undangan Pernikahan Digital Eksklusif".
- Hero headline (`:116`) hanya menulis "Pernikahan dimulai dengan undangan."
- Tab katalog (`:199`) hanya dua: "Undangan Pernikahan" dan "Undangan Ulang
  Tahun & Acara Lain" — opsi kedua adalah sampah kategori.
- `(marketing)/+layout.svelte:28` — announcement bar mempromosikan "Seluruh
  Template Pernikahan". Footer `:155` menyebut "undangan digital pernikahan".
- `(app)/dashboard/+page.svelte` menampilkan 4 modul statis tanpa hubungan
  dengan acara user.
- `(auth)/daftar/+page.server.ts` hanya mengumpulkan nama, email, provinsi,
  password. Tidak ada pertanyaan "acara apa yang kamu siapkan?"

Padahal di `packages/shared/src/constants/event-types.ts` sudah ada **10 tipe
acara** (pernikahan, lamaran, ulang tahun, khitanan, aqiqah, reuni, acara
perusahaan, syukuran, wisuda, lainnya). Sayangnya daftar itu **datar** —
daftar 10–20 kartu tidak scalable saat jenis undangan bertambah.

**Tujuan:** reposisi Ketuk.id sebagai platform undangan digital untuk **acara
penting di Indonesia**, dengan taxonomy **dua level: kategori → sub-kategori**.

## 2. Aturan main (jangan dilanggar)

1. **Frontend dulu, backend nanti.** Fase ini hanya menyentuh `frontend/` dan
   `packages/shared/`. **Jangan sentuh `backend/`** (Go) kecuali diperintahkan
   di §15.
2. **Backend adalah Go, bukan Bun.** Jangan tambah `backend/` ke
   `workspaces` root, jangan buat script Bun untuk backend, jangan import
   apapun dari `backend/`.
3. **Jangan `git commit`, `push`, atau `tag`.** Pemilik repo commit manual.
4. **Pakai token di `app.css`.** Palet `coffee`, `terracotta`, `cream`,
   `champagne`, plus warna modul (`undangan`, `planner`, `vendor`, `hadiah`).
   Jangan perkenalkan warna keenam.
5. **Jangan hapus data atau rombak `TEMPLATES` secara spekulatif.** Tambah,
   jangan rombak.
6. **Komentar Bahasa Indonesia, jelaskan *kenapa***, bukan *apa*. Komentar
   yang hanya mengulangi kode dihapus.

## 3. Taxonomy dua level (inti arsitektur)

**Pola: kategori besar → sub-kategori.** Bukan 10–20 kartu datar. Alasannya
skalabilitas: penambahan jenis undangan baru hanya menambah sub-kategori di
bawah kategori yang already ada, tanpa merombak UI maupun storage.

```text
Pernikahan        → Pernikahan, Lamaran/Pertunangan, Anniversary
Ulang Tahun       → Ulang Tahun Anak, Ulang Tahun Dewasa
Acara Keluarga    → Aqiqah, Khitanan, Syukuran, Wisuda, Arisan
Keagamaan         → Pengajian/Kajian, Marawis/Rebana, Haul/Peringatan
Komunitas & Reuni → Reuni, Gathering, Pertemuan Komunitas
Bisnis            → Corporate Event, Seminar, Workshop, Grand Opening
Lainnya           → (user menulis sendiri)
```

Setiap sub-kategori memetakan ke **salah satu dari 10 `EventType` backend**
yang sudah ada (`packages/shared/src/types/event.ts`). Ini supaya API
`POST /v1/events` tidak perlu berubah sama sekali di fase ini. Sub-kategori
yang belum punya pasangan backend (`anniversary`, `gathering`, `seminar`, …)
dipetakan ke `other`; keterangannya hidup di `eventSubType`.

Pemetaan lengkap (10 backend `EventType` semuanya tercakup):

| Kategori | Sub-kategori | Backend `EventType` |
|---|---|---|
| pernikahan | Pernikahan | `wedding` |
| pernikahan | Lamaran/Pertunangan | `engagement` |
| pernikahan | Anniversary | `other` |
| ulang-tahun | Ulang Tahun Anak | `birthday` |
| ulang-tahun | Ulang Tahun Dewasa | `birthday` |
| keluarga | Aqiqah | `aqiqah` |
| keluarga | Khitanan | `khitanan` |
| keluarga | Syukuran | `syukuran` |
| keluarga | Wisuda | `graduation` |
| keluarga | Arisan | `other` |
| keagamaan | Pengajian/Kajian | `other` |
| keagamaan | Marawis/Rebana | `other` |
| keagamaan | Haul/Peringatan | `other` |
| komunitas | Reuni | `reunion` |
| komunitas | Gathering | `other` |
| komunitas | Pertemuan Komunitas | `other` |
| bisnis | Corporate Event | `corporate` |
| bisnis | Seminar | `corporate` |
| bisnis | Workshop | `corporate` |
| bisnis | Grand Opening | `corporate` |
| other | (teks bebas) | `other` |

### 3.1 Tipe baru di `@ketuk/shared`

`packages/shared/src/types/onboarding.ts`:

```ts
import type { EventType } from './event';

/** Kategori acara level pertama — yang tampil sebagai kartu di onboarding. */
export type EventCategoryId =
	| 'pernikahan'
	| 'ulang-tahun'
	| 'keluarga'
	| 'keagamaan'
	| 'komunitas'
	| 'bisnis'
	| 'other';

export interface EventSubType {
	id: string;
	label: string;
	/** Tipe acara backend yang dikirim ke API saat membuat event nyata. */
	eventType: EventType;
}

export interface EventCategory {
	id: EventCategoryId;
	label: string;
	emoji: string;
	shortDesc: string;
	subTypes: EventSubType[];
	/** Hanya 'other': tidak ada sub-kategori tetap, user menulis nama acaranya. */
	customLabel?: boolean;
}

export interface UserEventPreference {
	/** Kategori yang dipilih (level pertama). */
	eventType: EventCategoryId;
	/**
	 * Sub-kategori (level kedua). Untuk 'other' berisi teks bebas nama acara
	 * yang user ketik — lihat §7.3. Bisa kosong kalau kategori punya satu
	 * sub-kategori dan dipilih langsung.
	 */
	eventSubType?: string;
	/** Tujuan utama di platform, mis. 'undangan' | 'planner'. Bebas teks. */
	purpose?: string;
	selectedAt: string;
}
```

Export dari barrel `packages/shared/src/index.ts`.

### 3.2 Taxonomy

`packages/shared/src/constants/event-categories.ts` — `EVENT_CATEGORIES:
EventCategory[]` sesuai tabel di atas. Kategori 'other' punya
`customLabel: true` dan `subTypes: []`. Lengkapi juga helper:

```ts
export function findCategory(id: string): EventCategory | undefined;
export function findSubType(subTypeId: string): { category: EventCategory; subType: EventSubType } | undefined;
```

`findSubType` dipakai `/app/undangan/baru` untuk menerjemahkan prefill +
mengirim `eventType` yang benar ke API.

Jalankan `bun run build:shared` setelah ini — tanpa rebuild, `frontend` tidak
bisa import tipenya dan `bun run check` akan gagal.

### 3.3 Dua konsep yang tetap dipisah

```text
User
 ├── onboardingPreference   ← "lagi nyiapin acara keluarga (aqiqah)"
 │     eventType: 'keluarga', eventSubType: 'aqiqah'
 │
 └── Events / Invitations   ← objek nyata, type-nya backend EventType
       ├── Wedding #1
       ├── Birthday #1
       └── Reunion #1
```

`onboardingPreference` = konteks/tonjolan saat ini, bisa diganti, boleh
kosong. **Bukan** kolom `type` di event. Event nyata tetap punya
`Event.type: EventType` sendiri.

## 4. Arsitektur adapter persistence (fase bertahap)

Semua komponen **dilarang** membaca `localStorage` langsung. Semua akses
preferensi lewat satu port, supaya fase backend cukup mengganti implementasi.

### 4.1 Port

`frontend/src/lib/onboarding/types.ts`:

```ts
import type { UserEventPreference } from '@ketuk/shared';

export interface OnboardingStore {
	/** null = belum onboarding. Bukan throw, supaya UI bisa pakai default. */
	get(): UserEventPreference | null;
	set(preference: UserEventPreference): void;
	/** Catat bahwa user sengaja melewati onboarding. Lihat §7.4 kenapa ini wajib. */
	skip(): void;
	isSkipped(): boolean;
	clear(): void;
}
```

### 4.2 Implementasi sementara: localStorage

`frontend/src/lib/onboarding/local-store.ts`. Tangani **semua** kasus, bukan
happy path saja:

- **SSR** — pakai `import { browser } from '$app/environment'`. Di server,
  `get()` → `null`, method lain no-op.
- **Storage diblokir** (Safari private mode, permission ditolak) — tangkap
  `SecurityError`/`QuotaExceededError`, jangan sampai meledak ke UI.
- **JSON korup / versi lama** — `JSON.parse` gagal atau bentuk tidak cocok →
  anggap tidak ada preferensi, hapus key yang rusak.
- **`eventType` tidak dikenal** — nilai yang tidak ada di `EVENT_CATEGORIES`
  → anggap null, jangan dipakai begitu saja.
- **Bentuk setengah** — field wajib hilang → anggap tidak valid.

Dua key terpisah: `ketuk.onboardingPreference` (objek §3.1) dan
`ketuk.onboardingSkipped` (ISO timestamp). Skip dipisah dari preference karena
dua hal yang berbeda: preference = jawaban, skipped = "jangan tanya lagi".

### 4.3 Composable untuk UI

`frontend/src/lib/onboarding/preference.svelte.ts` — Svelte 5 runes, bukan
store lama. Satu-satunya jalan komponen ke preferensi:

```ts
import { browser } from '$app/environment';
import { localStore } from './local-store';
import type { UserEventPreference } from '@ketuk/shared';

export function createOnboardingPreference() {
	let current = $state<UserEventPreference | null>(
		browser ? localStore.get() : null,
	);
	let skipped = $state<boolean>(browser ? localStore.isSkipped() : false);

	function set(preference: UserEventPreference) {
		current = preference;
		localStore.set(preference);
	}
	function skip() {
		skipped = true;
		localStore.skip();
	}
	function clear() {
		current = null;
		skipped = false;
		localStore.clear();
	}

	return {
		get value() {
			return current;
		},
		get isSkipped() {
			return skipped;
		},
		set,
		skip,
		clear,
	};
}
```

Di-mount per komponen, bukan global singleton — Svelte 5 tidak butuh context
global, dan tiap halaman yang butuh preferensi cukup panggil fungsi ini.

### 4.4 Janji untuk fase backend (§15)

Urutan pemanggilan tidak boleh berubah saat implementasi diganti.
`supabase-store.ts` akan membaca kolom baru di `profiles` (atau tabel
`onboarding_preferences`), jatuh ke localStorage sebagai cache offline, dan
`set()` melakukan upsert. UI tidak boleh tahu perbedaan itu.

## 5. Konfigurasi personalisasi per kategori

`frontend/src/lib/onboarding/event-config.ts` — sumber kebenaran untuk
personalisasi dashboard. **Komponen membaca config, tidak memeriksa string
kategori.** Alasannya: kalau logika `{#if eventType === 'wedding'}` tersebar
di 15 komponen, tambah kategori baru jadi kerjaan mencari-jumpa; lewat config,
satu file selesai.

```ts
import type { EventCategoryId } from '@ketuk/shared';

export type ModuleId = 'undangan' | 'planner' | 'vendor' | 'hadiah';

export interface EventConfig {
	category: EventCategoryId;
	label: string;
	emoji: string;
	shortDesc: string;
	/** Modul yang tampil di dashboard + urutannya. Modul di luar ini = irrelevant, tidak ditampilkan. */
	modules: ModuleId[];
	/** Fitur utama modul Undangan untuk kategori ini — copy spesifik, bisa diverifikasi. */
	invitationFeatures: string[];
	/** Fitur sekunder yang tampil di blok "Fitur lainnya" (visual weight rendah). */
	optionalFeatures?: string[];
	/** Default copy untuk empty-state CTA, mis. "Buat Undangan Ulang Tahun". */
	ctaLabel: string;
}

export function getEventConfig(category: EventCategoryId | null | undefined): EventConfig;
```

Isi (keputusan pemilik produk — **dashboard membedakan primary / optional /
irrelevant**, bukan urutan yang sama untuk semua kategori):

| Kategori | modules | invitationFeatures | optionalFeatures |
|---|---|---|---|
| pernikahan | undangan, planner, vendor, hadiah | Data mempelai; Jadwal akad & resepsi; Lokasi & peta; RSVP; Buku tamu; Galeri; Amplop digital | — |
| ulang-tahun | undangan, planner, vendor, hadiah | Nama resi & tema; Tanggal & waktu; Lokasi & peta; RSVP; Galeri; Kado digital | — |
| keluarga | undangan, planner, vendor, hadiah | Informasi acara; Tanggal & waktu; Lokasi & peta; RSVP; Galeri | — |
| keagamaan | undangan, planner, vendor | Informasi acara; Tanggal & waktu; Lokasi; RSVP | Hadiah / Gift |
| komunitas | undangan, planner, vendor | Informasi acara; Tanggal & waktu; Lokasi; RSVP; Daftar tamu | Hadiah / Gift |
| bisnis | undangan, planner, vendor | Informasi acara; Tanggal & waktu; Lokasi; RSVP; Daftar tamu | Hadiah / Gift |
| **other** | undangan, planner, vendor, hadiah | Informasi acara; Tanggal & waktu; Lokasi; Deskripsi; RSVP; Daftar tamu; Galeri; Kontak; Bagikan undangan | — |

### 5.1 Aturan terpenting: `other` JANGAN fallback ke wedding

Ini keputusan eksplisit pemilik produk. **`other` wajib konfigurasi netral
generik** (baris terakhir tabel di atas). Fitur `Data mempelai`, `Akad`,
`Resepsi`, `Amplop pernikahan` **tidak boleh muncul** hanya karena fallback-nya
wedding. Alasannya: kalau `other` meniru wedding, arsitektur produk diam-diam
tetap wedding-centric — persis yang ingin dihilangkan dari redesign ini.

Eksekusinya: `getEventConfig(null)` dan `getEventConfig('other')` **sama-sama**
mengembalikan konfigurasi netral. Tidak ada cabang `if (category === 'other')
return weddingConfig` di mana pun.

`getEventConfig(null)` dipakai saat user skip onboarding — default sehat, bukan
kosong.

## 6. Route area aplikasi: `/app/*`

Area marketing/public dan area aplikasi terpisah. Semua route yang
ter-autentikasi dipindah ke bawah segment `/app/`:

```text
/app/onboarding     ← baru
/app/dashboard      ← dulunya /dashboard
/app/undangan       ← /undangan
/app/planner        ← /planner
/app/vendor         ← /vendor
/app/hadiah         ← /hadiah
```

**Kenapa pindah semua, bukan cuma onboarding:** onboarding di `/app/` sementara
dashboard di root jadi dua area aplikasi terpisah — persis yang ingin
dihindarkan. Nama modul (`undangan`, `planner`, …) **tidak diubah**;
menggantinya ke `/app/events`, `/app/templates`, `/app/settings` adalah
keputusan produk terpisah untuk fase lain.

### 6.1 Cara migrasi

1. Pindah folder `frontend/src/routes/(app)/` → `frontend/src/routes/app/`
   (group `(app)` hilang, jadi segment nyata). Struktur dalamnya tidak berubah.
2. Update **semua** referensi internal (hasil grep: ~35 lokasi, termasuk
   `(auth)/+layout.server.ts:9`, `(auth)/masuk/+page.server.ts:72`,
   `(auth)/callback/+server.ts:7`, `(auth)/daftar/+page.svelte:47` (OAuth
   `redirectTo`), `admin/+layout.svelte:265`, `EventCard.svelte:18`,
   `VendorCard.svelte:17`, dan semua `href`/`goto`/`redirect` di dalam `app/`).
3. Tambah **redirect permanen** dari path lama lewat `hooks.server.ts`
   (satu tempat, bukan puluhan stub route): prefix
   `['/dashboard','/undangan','/planner','/vendor','/hadiah']` →
   `throw redirect(308, \`/app${pathname}\`)`. Jalankan **hanya** untuk
   `BASE_HOSTNAMES` (jangan di subdomain undangan), dan taruh pertama dalam
   `sequence()` sebelum `handleSubdomain`. 308 dipilih karena ini redirect
   permanen untuk bookmark/link lama.

Verifikasi: buka `/dashboard` langsung → harus mendarat di `/app/dashboard`
dengan sidebar aktif.

## 7. Flow onboarding

```text
Register / login pertama
   ↓
/app/onboarding  ← pilih kategori → sub-kategori
   ↓
simpan preference (adapter §4)
   ↓
/app/dashboard  ← personalisasi sesuai kategori
   ↓
template + fitur + produk sesuai kategori
```

### 7.1 Kapan muncul

- **Setelah register berhasil** — `(auth)/daftar/+page.server.ts:176` default
  target berubah dari `/dashboard` ke `/app/onboarding`. Pertahankan logika
  `next` (open-redirect guard) — jangan hapus.
- **Gate di sisi klien** — server tidak bisa baca localStorage, jadi gate
  hidup di `app/+layout.svelte` (bukan `+layout.server.ts`).
  `app/+layout.server.ts` memanggil `listMyEvents` melalui adapter API Go dan
  mengembalikan `hasEvents` ke layout; layout men-cek: jika
  `pathname !== '/app/onboarding'` **dan**
  tidak ada preference **dan** belum skip **dan** `!hasEvents` →
  `goto('/app/onboarding')`. Query ini sementara — fase §15 menghilangkannya
  saat preferensi pindah ke server.
- **Login user lama yang sudah punya event** → `hasEvents` true, gate tidak
  memicu. User lama tanpa event sama sekali tetap diarahkan (memang belum
  punya konteks).
- **Bisa dibuka ulang** dari dashboard lewat tombol "Ganti jenis acara".

### 7.2 Spesifikasi UI

Route: `frontend/src/routes/app/onboarding/+page.svelte`.

- **Satu layar, bukan multi-step.** "Acara apa yang sedang kamu siapkan?" +
  grid pilihan + tombol simpan. Multi-step untuk satu pertanyaan adalah
  over-engineering.
- **Grid 2 kolang di mobile 360px, tanpa horizontal scroll.** 6 kartu kategori
  utama + 1 "Lainnya" = 7 kartu (4 baris di mobile). Kartu: emoji, label,
  deskripsi 1 baris. Tap target ≥44px kedua dimensi. Selected state jelas
  (border + background), bukan gelap 5%.
- **Drill-down sub-kategori.** Saat kategori dengan >1 sub-kategori dipilih,
  tampilkan sub-kategori sebagai pilihan keduan (chip/kartu kecil di bawah,
  atau ganti grid). Kategori dengan **tepat satu** sub-kategori → pilih
  langsung, tidak usah drill-down (drill-down 1 pilihan adalah sampah UI).
  Di taxonomy §3.2 hanya 'ulang-tahun' dan 'other' yang butuh perlakuan
  khusus: ulang-tahun punya 2 sub-kategori (drill-down), other lihat §7.3.
- **Simpan** → `set({ eventType, eventSubType, selectedAt: new
  Date().toISOString() })` → `goto('/app/dashboard')`.
- **Skip** wajib ada ("Lewati untuk sekarang") → `skip()` →
  `goto('/app/dashboard')` dengan konfigurasi netral. Tanpa skip, gate §7.1
  jadi loop tak terhindarkan.
- **Loading state** di tombol simpan. **Error state** kalau adapter gagal
  (storage penuh/diblokir): pesan singkat + tetap lanjut ke dashboard dengan
  preferensi in-memory sesi itu saja.
- Aksesibilitas: `<form>` dengan `<input type="radio">` asli (bukan div
  kustom) atau `aria-pressed`. Arrow key memilih, Enter/Space submit.

### 7.3 Kategori "Lainnya" — input nama acara

Saat `other` dipilih, tampilkan input teks **"Nama/Jenis Acara"** (mis.
"Tunangan adik", "Syukuran kantor", "Buka puasa bersama"). Nilainya disimpan
sebagai `eventSubType`, dan backend `eventType`-nya `other`.

Tidak ada saran otomatis / autocomplete yang mengarah ke wedding. Placeholder
boleh netral ("Tulis jenis acaramu").

### 7.4 Yang tidak boleh ada di onboarding

- Pertanyaan nama tanggapan (sudah di register).
- Upsell paket berbayar — onboarding yang langsung jualan terasa jebakan.
- Wajib isi tanggal acara — itu tugas `/app/undangan/baru`.

## 8. Dashboard personalisasi

File: `frontend/src/routes/app/dashboard/+page.svelte`.

1. **Header dinamis** sesuai kategori: "Persiapan pernikahanmu" (pernikahan),
   "Persiapan ulang tahun" (ulang-tahun), "Persiapan acara keluarga" (keluarga).
   Tanpa preference: "Mau melakukan apa?" (teks lama, sudah baik). Copy dari
   config, bukan `{#if}` per string di komponen.
2. **Kartu modul** = `config.modules` dengan copy spesifik kategori di tiap
   kartu (`invitationFeatures` untuk kartu undangan, `ctaLabel` untuk CTA).
   Modul yang tidak ada di `config.modules` **tidak ditampilkan** — ini
   "irrelevant" (§5). Hadiah untuk keagamaan/komunitas/bisnis tidak muncul
   sebagai kartu.
3. **Blok "Fitur lainnya"** di bawah grid untuk `optionalFeatures` — visual
   weight rendah (teks kecil, tanpa kartu besar). Hadiah untuk kategori
   tersebut hidup di sini, bukan sebagai kartu utama.
4. **Empty state pintar**: user baru tanpa event → CTA utama memakai
   `ctaLabel` ("Buat Undangan Ulang Tahun") yang membawa prefill ke
   `/app/undangan/baru?type=…`.
5. **Ganti jenis acara**: tombol kecil di header → balik ke `/app/onboarding`,
   pilihan lama ter-highlight.
6. **List event existing** di bawah modul (kalau user sudah punya), pakai
   `EventCard` yang sudah ada. Lebih berguni dari sekadar 4 ikon.
7. **Sidebar nav** `(app)/+layout.svelte:28` tetap menampilkan kelima modul
   (navigasi global, bukan tonjolan kategori). Modul yang "irrelevant" tetap
   bisa dibuka lewat sidebar — tidak dihapus dari sistem, hanya tidak
   ditonjolkan.

## 9. `/app/undangan/baru` — prefill & select berkelompok

- Terima query param `type`. Nilainya bisa **sub-type id** (`aqiqah`, `seminar`)
  atau **backend EventType** (`wedding`). Validasi lewat `findSubType()` /
  `EVENT_TYPES`; tolak nilai asal dengan diam (fallback ke preference atau
  kosong). Jangan pernah passing input user mentah ke API.
- Select "Jenis acara" diubah dari flat `EVENT_TYPE_CONFIGS` menjadi
  **grouped by category** pakai `EVENT_CATEGORIES` (sudah ada pola
  `SelectOptionGroup` di `lib/components/ui` — lihat
  `(auth)/daftar/+page.svelte:38` untuk pemakaian). Value = sub-type id.
- Saat submit, kirim `type: <subType.eventType>` ke API — payload backend
  tidak berubah.
- Label input "Judul acara" placeholder ikut kategori: "Pernikahan Budi &
  Sinta" untuk pernikahan, "Ulang Tahun Ananda ke-7" untuk ulang-tahun.
  Ambil dari config, jangan hardcode 10 placeholder.

## 10. Homepage

File: `frontend/src/routes/(marketing)/+page.svelte` (596 baris).

- **Hero**: copy mencakup acara umum, tapi **jangan jadi slogan tanpa
  kepribadian**. Tetap memimpin dengan pernikahan sebagai kategori terkuat,
  lalu tunjukkan yang lain. Contoh arah (finalkan saat menulis sesuai
  antislop-copywriting): "Setiap momen besar dimulai dengan undangan." +
  subcopy yang menyebut pernikahan, ulang tahun, aqiqah, syukuran, reuni.
- **Meta + `<svelte:head>`** (`:86`): dari "Undangan Pernikahan Digital
  Eksklusif" ke positioning umum.
- **Tab katalog** (`:199`): ganti 2 tab jadi filter kategori nyata — pakai
  `EVENT_CATEGORIES` (atau subset cerdas: Semua / Pernikahan / Ulang Tahun /
  Acara Keluarga / Keagamaan / Komunitas / Bisnis). Pola pill yang sudah ada
  di `:236` (`filterCategories`) bisa ditiru, jangan buat mekanisme baru.
- **`TEMPLATES`** (`frontend/src/lib/data/templates.ts`): tambah field
  opsional `eventTypes?: EventCategoryId[]` di `TemplateItem` — jangan ubah
  struktur yang ada. Isi hanya untuk template yang jelas non-pernikahan;
  template tanpa field dianggap cocok untuk "Semua". Copy
  `preview.coupleName` dan envelope script pernikahan-sentris: perbarui untuk
  template non-pernikahan saja.
- **FAQ** (`:61`): jawaban menjanjikan hal spesifik pernikahan ("nama
  pengantin", "pasangan bahagia"). Perbarui generalisasi yang salah, sisakan
  yang benar. **FAQ tetap section homepage, bukan halaman baru.**
- **PhoneMockup** (`:189`): `coupleNames="Sarah & Dimas"` → dinamis sesuai
  kategori aktif (pasangan untuk pernikahan, "Ulang Tahun ke-7 Ananda" untuk
  ulang-tahun). Tambah prop kalau perlu, jangan rombak komponen.
- **Social proof** (`:161`): "(500+ ulasan pasangan bahagia)" → bahasa yang
  tidak mengecualikan acara lain.

## 11. Navigasi & footer

File: `frontend/src/routes/(marketing)/+layout.svelte`.

- **Aturan utama pemilik produk: navbar tidak boleh menduplikasi section
  homepage sebagai halaman terpisah.** `navLinks` saat ini sudah anchor-based
  (`/#katalog`, `/#fitur`, `/#cara-kerja`, `/#faq`) — **ini sudah benar,
  pertahankan.** Hanya `/harga` dan `/tentang` yang berupa halaman.
- Label "Koleksi Desain" → yang tidak menyiratkan hanya pernikahan. **Jangan
  buat mega-menu** kalau tidak benar-benar dibutuhkan.
- **Announcement bar** (`:28`): hapus "Seluruh Template Pernikahan". Ganti
  benefit umum yang jujur (aktif selamanya, tanpa watermark) atau hapus
  sepenuhnya.
- **Footer** (`:154`): deskripsi "undangan digital pernikahan dan event
  modern" → perbaiki urutan jadi "undangan digital untuk setiap acara
  penting". Kolom "Koleksi Desain" (`:172`) berisi 5 link label
  pernikahan-sentris → ganti jadi kategori nyata.
- **Nav app** `app/+layout.svelte`: item pertama boleh label dinamis sesuai
  kategori ("Undangan" → "Undangan Pernikahan") **hanya** kalau tidak terlalu
  panjang di mobile drawer.

## 12. Modul lain — seberapa jauh personalisasinya

1. **`/app/undangan` list**: empty state user `ulang-tahun` tanpa event → CTA
   "Buat Undangan Ulang Tahun" dengan prefill, bukan generik.
2. **`/app/planner`** (`+page.svelte:25`): empty state "Belum ada undangan" →
   pakai istilah "acara" + CTA prefill sesuai preference.
3. **`/app/vendor`** & **`/app/hadiah`**: cek
   `packages/shared/src/constants/vendor-categories.ts` — kalau ada kategori
   yang jujur dipetakan per kategori acara (khitanan → catering & tenda),
   lakukan; kalau tidak ada mapping yang jujur, **jangan rekayasa**, lewati
   dan catat di laporan.
4. **Marketing `/template`**: halaman katalog publik, terima filter kategori
   yang sama dengan homepage.

Semua modul ini memanggil backend Go via `lib/api`. **Jangan perkirakan
endpoint baru** — hanya perubahan presentasi + query param yang frontend
olah sendiri.

## 13. Mobile & aksesibilitas

(dari antislop-layoutmobile & antislop-human — wajib, bukan saran)

- **Mobile-first**: semua layout baru ditulis untuk 360px dulu, baru
  scale-up. Grid onboarding: 2 kolom mobile, 3 sm, 4 lg.
- **Tidak ada horizontal scroll** untuk pilihan utama onboarding. Scroll
  horizontal hanya untuk daftar template (pola `overflow-x-auto` sudah ada di
  homepage `:237`).
- **Tap target ≥44px**. Cek ulang semua tombol/kartu.
- **Focus visible**: `app.css` punya `:focus-visible` rule. Jangan `outline:
  none` tanpa pengganti.
- **Kontras**: `coffee-500` di atas `cream-50` — cek ratio untuk teks kecil;
  kalau < 4.5:1 pakai `coffee-700`.
- **Reduced motion**: `app.css` handle global. Jangan animasi >200ms tanpa
  hormati preferensi.
- **`aria-current`** untuk filter pill aktif.
- **Tes keyboard**: tab melalui onboarding sampai submit tanpa mouse.

## 14. Copy

(dari antislop-copywriting)

- **Bahasa Indonesia.** UI app pakai "kamu" (sudah konsisten di daftar:
  "Isi lima kolom di bawah"). Marketing boleh "Anda" (sudah konsisten).
- **Dilarang**: "Momen tak terlupakan", "Wujudkan mimpi", "Solusi terbaik",
  "Lebih dari sekadar undangan", "Tak terbatas", emoji sebagai dekorasi.
- **CTA menjelaskan hasil**: "Buat Undangan" (baik), "Mulai perjalananmu"
  (buruk).
- **Deskripsi kartu harus bisa diverifikasi**: "Budget, checklist, dan
  timeline acaramu" (baik), "Semua yang kamu butuhkan" (buruk).
- Hilangkan "haute-couture" dari footer (`:156`) — janji yang tidak didukung
  produk.

## 15. Fase backend — batas yang tidak boleh dilewati sekarang

**Jangan implementasikan di fase ini.** Rancang frontend supaya ini jadi
drop-in:

1. Kolom `onboarding_preference jsonb` di `profiles`, atau tabel
   `onboarding_preferences(user_id uuid pk, …)`. Migrasi goose di
   `backend/internal/migrate/sql/` (format `NNN_nama.sql`, ada 16 file).
2. Endpoint Go: `GET` + `PUT /v1/me/onboarding` (ter-autentikasi `s.auth`).
   Ikuti konvensi `internal/identity` (hexagonal: port Repository + sentinel
   errors), handler di `internal/httpserver`.
3. Adapter frontend `supabase-store.ts` menggantikan `local-store.ts` dalam
   port yang sama. `localStorage` jadi cache offline. Gate §7.1 pindah ke
   `+layout.server.ts` (preferensi terbaca server, query `hasEvents` hilang).
4. **Backend hanya Go.** Jangan tambah `backend/` ke `workspaces` Bun; jalan
   via `make run` di `backend/`.

Kalau pekerjaan ini **terpaksa** menyentuh backend karena frontend tidak jalan
tanpanya: berhenti, catat di laporan, jangan modifikasi `backend/` tanpa
konfirmasi pemilik repo.

## 16. Verifikasi — harus lulus sebelum bilang selesai

Dari root repo:

```bash
bun run build:shared   # WAJIB pertama
bun run check          # tsc + svelte-check
bun run lint           # biome
```

Backend (hanya kalau §15 disentuh — normalnya tidak):

```bash
cd backend && go build ./... && go vet ./... && go test ./...
```

Tes manual wajib (ceklist di laporan):

- [ ] Register baru → mendarat di `/app/onboarding`, bukan dashboard.
- [ ] `/dashboard` lama → redirect 308 ke `/app/dashboard`; semua link internal
      (`EventCard`, `VendorCard`, sidebar, auth redirect) mengarah ke `/app/*`.
- [ ] Pilih "Acara Keluarga" → drill-down muncul → pilih "Aqiqah" → dashboard
      menampilkan copy keluarga, bukan wedding.
- [ ] Pilih "Ulang Tahun" → dashboard **tidak** menampilkan "Data mempelai"
      atau "Akad".
- [ ] Pilih "Lainnya" → input "Nama/Jenis Acara" muncul; disimpan sebagai
      `eventSubType`; dashboard memakai fitur netral (Informasi acara, Tanggal
      & waktu, Lokasi, Deskripsi, RSVP, Daftar tamu, Galeri, Kontak, Bagikan
      undangan) — tidak ada mempelai/akad/resepsi/amplop pernikahan.
- [ ] Onboarding → dashboard → "Ganti jenis acara" → pilihan lama
      ter-highlight.
- [ ] Skip onboarding → dashboard render dengan config netral, tidak crash,
      tidak di-redirect balik ke onboarding (skip flag bekerja).
- [ ] Refresh halaman → preferensi bertahan (localStorage).
- [ ] Tab private (storage diblokir) → tidak ada exception konsol; onboarding
      tetap bisa dipakai sesi itu.
- [ ] DevTools → hapus key storage → halaman masih render.
- [ ] `/app/undangan/baru?type=aqiqah` → select terisi & grouped dengan benar;
      `/app/undangan/baru?type=hack` → ditolak dengan diam.
- [ ] User lama yang sudah punya event → tidak dipaksa onboarding saat login.
- [ ] Homepage: filter kategori bekerja, pill aktif ada `aria-current`.
- [ ] Mobile 360px: onboarding 2 kolom, tidak ada horizontal scroll, tap
      target ≥44px.
- [ ] Keyboard: tab saja onboarding sampai submit; Enter memilih & mengirim.
- [ ] Lighthouse/Axe: tidak ada error kontras baru.

Kalau ada item gagal dan tidak diperbaiki dalam sesi ini, sebutkan **terus
terang**. Jangan tahan informasi agar terlihat selesai.

## 17. Urutan pengerjaan

1. `packages/shared`: tipe + taxonomy + helper, lalu `bun run build:shared`
   (§3).
2. Migrasi route `(app)` → `app` + hooks redirect + update semua referensi
   (§6).
3. `lib/onboarding/*`: port, local-store, composable, event-config (§4, §5).
4. Route `/app/onboarding` + gate di layout + redirect register (§7).
5. Dashboard personalisasi + prefill `/app/undangan/baru` (§8, §9).
6. Homepage (§10).
7. Nav & footer (§11).
8. Modul lain (§12).
9. Verifikasi §16.

Jalankan `bun run check` + `bun run lint` di akhir blok yang menyentuh banyak
file (minimal: blok 2, 3, 5, 8).

## 18. Dilarang keras

- Menghapus file lama tanpa konfirmasi (kecuali yang dibuat di sesi ini).
- `git commit`, `git push`, `git tag`, `git reset --hard`.
- Menyentuh `backend/` (Go) atau `.env`.
- Membuat halaman FAQ/Tentangan baru "supaya lengkap" — FAQ tetap section
  homepage.
- Menulis komentar yang mengulangi kode.
- Mengarang angka social proof baru tanpa sumber data.
- Meninggalkan kode mati atau TODO tanpa alasan yang ditulis di komentar.
- **Fallback `other` ke konfigurasi wedding** (lihat §5.1 — ini alasan utama
  redesign ini ada).
