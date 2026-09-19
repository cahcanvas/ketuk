import { type EventCategoryId, findCategory } from '@ketuk/shared';

/**
 * Sumber kebenaran personalisasi dashboard per kategori acara.
 *
 * Komponen membaca config ini, TIDAK memeriksa string kategori langsung
 * (`{#if eventType === 'wedding'}`). Alasannya: saat kategori baru ditambahkan,
 * satu file ini selesai; kalau logika tersebar di banyak komponen, tambah
 * kategori jadi kerjaan mencari-jumpa di seluruh codebase.
 *
 * Tiga tingkat relevansi modul:
 *  - `modules`: modul yang ditampilkan sebagai kartu utama, urutan = urutan tampil.
 *    Modul di luar daftar ini tidak ditampilkan (tidak relevan untuk kategori ini).
 *  - `invitationFeatures`: fitur modul Undangan yang ditonjolkan, copy spesifik
 *    kategori dan bisa diverifikasi.
 *  - `optionalFeatures`: fitur sekunder, tampil di blok "Fitur lainnya" dengan
 *    visual weight rendah.
 */

export type ModuleId = 'undangan' | 'planner' | 'vendor' | 'hadiah';

export interface EventConfig {
	category: EventCategoryId;
	modules: ModuleId[];
	invitationFeatures: string[];
	optionalFeatures: string[];
}

const CONFIGS: Record<EventCategoryId, EventConfig> = {
	pernikahan: {
		category: 'pernikahan',
		modules: ['undangan', 'planner', 'vendor', 'hadiah'],
		invitationFeatures: [
			'Data mempelai',
			'Jadwal akad & resepsi',
			'Lokasi & peta',
			'RSVP',
			'Buku tamu',
			'Galeri',
			'Amplop digital',
		],
		optionalFeatures: [],
	},
	'ulang-tahun': {
		category: 'ulang-tahun',
		modules: ['undangan', 'planner', 'vendor', 'hadiah'],
		invitationFeatures: [
			'Nama resi & tema',
			'Tanggal & waktu',
			'Lokasi & peta',
			'RSVP',
			'Galeri',
			'Kado digital',
		],
		optionalFeatures: [],
	},
	keluarga: {
		category: 'keluarga',
		modules: ['undangan', 'planner', 'vendor', 'hadiah'],
		invitationFeatures: ['Informasi acara', 'Tanggal & waktu', 'Lokasi & peta', 'RSVP', 'Galeri'],
		optionalFeatures: [],
	},
	keagamaan: {
		category: 'keagamaan',
		modules: ['undangan', 'planner', 'vendor'],
		invitationFeatures: ['Informasi acara', 'Tanggal & waktu', 'Lokasi', 'RSVP'],
		optionalFeatures: ['Hadiah / Gift'],
	},
	komunitas: {
		category: 'komunitas',
		modules: ['undangan', 'planner', 'vendor'],
		invitationFeatures: ['Informasi acara', 'Tanggal & waktu', 'Lokasi', 'RSVP', 'Daftar tamu'],
		optionalFeatures: ['Hadiah / Gift'],
	},
	bisnis: {
		category: 'bisnis',
		modules: ['undangan', 'planner', 'vendor'],
		invitationFeatures: ['Informasi acara', 'Tanggal & waktu', 'Lokasi', 'RSVP', 'Daftar tamu'],
		optionalFeatures: ['Hadiah / Gift'],
	},
	/**
	 * Konfigurasi netral generik. ATURAN PENTING: 'other' dan null (user skip
	 * onboarding) SAMPAI menggunakan konfigurasi ini. Tidak boleh ada cabang yang
	 * mengembalikan konfigurasi pernikahan untuk 'other' — itu diam-diam membuat
	 * arsitektur produk tetap wedding-centric, persis yang ingin dihindari dari
	 * redesign ini. Lihat REDESIGN-UI.md §5.1.
	 */
	other: {
		category: 'other',
		modules: ['undangan', 'planner', 'vendor', 'hadiah'],
		invitationFeatures: [
			'Informasi acara',
			'Tanggal & waktu',
			'Lokasi',
			'Deskripsi',
			'RSVP',
			'Daftar tamu',
			'Galeri',
			'Kontak',
			'Bagikan undangan',
		],
		optionalFeatures: [],
	},
};

/**
 * Konfigurasi netral, dipakai untuk kategori 'other' DAN untuk user yang skip
 * onboarding. Dipisah jadi konstanta supaya keduanya jelas-jelas merujuk ke
 * satu objek yang sama, bukan dua cara menulis 'other'.
 */
const NEUTRAL_CONFIG = CONFIGS.other;

export function getEventConfig(category: EventCategoryId | null | undefined): EventConfig {
	// Null (user skip onboarding) diperlakukan sama dengan 'other'. Keduanya
	// tidak boleh jatuh ke konfigurasi pernikahan — lihat komentar di CONFIGS.other.
	return (category && CONFIGS[category]) ?? NEUTRAL_CONFIG;
}

/**
 * Label kategori untuk copy UI. Melalui taxonomy (bukan hardcode di sini)
 * supaya penambahan kategori otomatis terbawa.
 */
export function getCategoryLabel(category: EventCategoryId | null | undefined): string {
	if (!category) return 'Acara';
	return findCategory(category)?.label ?? 'Acara';
}

/**
 * Judul dashboard. Pola "Persiapan {label}" dipakai seragam untuk semua kategori
 * supaya tidak ada bentuk kepunyaan yang aneh ("acara keluargamu").
 */
export function buildDashboardTitle(category: EventCategoryId | null | undefined): string {
	if (!category) return 'Mau melakukan apa?';
	return `Persiapan ${getCategoryLabel(category).toLowerCase()}`;
}

/**
 * Label CTA empty-state. Pakai label sub-kategori bila ada (lebih spesifik:
 * "Buat Undangan Aqiqah"), atau teks bebas kategori 'other', atau label kategori.
 */
export function buildCtaLabel(
	category: EventCategoryId | null | undefined,
	eventSubType?: string | null,
): string {
	if (!category) return 'Buat Undangan';
	if (category === 'other' && eventSubType) return `Buat Undangan ${eventSubType}`;
	if (eventSubType) {
		const found = findCategory(category)?.subTypes.find((s) => s.id === eventSubType);
		if (found) return `Buat Undangan ${found.label}`;
	}
	return `Buat Undangan ${getCategoryLabel(category)}`;
}

/**
 * Placeholder input judul acara di /app/undangan/baru. Spesifik per kategori
 * supaya user tahu format yang diharapkan, bukan placeholder generik.
 */
export function getTitlePlaceholder(
	category: EventCategoryId | null | undefined,
	eventSubType?: string | null,
): string {
	if (category === 'pernikahan') {
		return eventSubType === 'engagement' ? 'Lamaran Budi & Sinta' : 'Pernikahan Budi & Sinta';
	}
	if (category === 'ulang-tahun') return 'Ulang Tahun Ananda ke-7';
	if (category === 'other' && eventSubType) return eventSubType;
	if (!category) return 'Judul acara kamu';
	return getCategoryLabel(category);
}
