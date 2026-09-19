import type { EventType } from './event';

/**
 * Preferensi onboarding: acara apa yang sedang disiapkan user saat ini.
 *
 * Ini BUKAN tipe event yang user buat. Satu user boleh punya banyak event
 * dengan tipe berbeda (lihat `Event.type`); preferensi ini hanya mengatur
 * urutan dan tonjolan yang ia lihat di dashboard & katalog. Karena itu
 * dipisahkan dari Event, bukan menjadi kolom `type` di profile.
 *
 * Taxonomy-nya dua level: kategori (level pertama, tampil sebagai kartu di
 * onboarding) → sub-kategori (level kedua, muncul saat kategori dipilih).
 * Alasannya skalabilitas: jenis undangan baru cukup menambah sub-kategori di
 * bawah kategori yang sudah ada, tanpa merombak UI maupun bentuk storage.
 */
export interface UserEventPreference {
	/** Kategori level pertama, mis. 'pernikahan' | 'keluarga' | 'other'. */
	eventType: EventCategoryId;
	/**
	 * Sub-kategori level kedua, mis. 'aqiqah'. Untuk kategori 'other' berisi
	 * teks bebas nama acara yang user ketik sendiri. Bisa kosong kalau kategori
	 * punya tepat satu sub-kategori dan dipilih langsung tanpa drill-down.
	 */
	eventSubType?: string;
	/**
	 * Tujuan utama user di platform, mis. 'undangan' | 'planner' | 'hadiah'.
	 * Bebas teks dan tidak divalidasi ketat — dipakai untuk mengatur tonjolan
	 * modul, bukan untuk logika wajib.
	 */
	purpose?: string;
	/** ISO 8601. Dipakai untuk mengetahui kebaruan konteks user. */
	selectedAt: string;
}

/**
 * Kategori acara level pertama. Dijadikan array (bukan cuma union type) supaya
 * bisa dipakai langsung sebagai sumber kebenaran kategori yang ditampilkan,
 * mengikuti pola `EVENT_TYPES` di types/event.ts tanpa duplikasi daftar nilai.
 */
export const EVENT_CATEGORY_IDS = [
	'pernikahan',
	'ulang-tahun',
	'keluarga',
	'keagamaan',
	'komunitas',
	'bisnis',
	'other',
] as const;

export type EventCategoryId = (typeof EVENT_CATEGORY_IDS)[number];

/**
 * Type guard untuk nilai yang datang dari luar (localStorage, query param).
 * Di sini, bukan di komponen, supaya satu tempat yang tahu bentuk validnya.
 */
export function isEventCategoryId(value: unknown): value is EventCategoryId {
	return typeof value === 'string' && (EVENT_CATEGORY_IDS as readonly string[]).includes(value);
}

export interface EventSubType {
	id: string;
	label: string;
	/**
	 * Tipe acara backend yang dikirim ke API saat membuat event nyata. Sub-kategori
	 * yang belum punya pasangan backend (mis. 'seminar', 'gathering') dipetakan ke
	 * 'other'; keterangannya hidup di `eventSubType` saja, jadi API tidak perlu
	 * berubah saat taxonomy bertambah.
	 */
	eventType: EventType;
}

export interface EventCategory {
	id: EventCategoryId;
	label: string;
	emoji: string;
	shortDesc: string;
	subTypes: EventSubType[];
	/**
	 * Hanya 'other': tidak ada sub-kategori tetap. User menulis nama acaranya
	 * sendiri dan disimpan sebagai `UserEventPreference.eventSubType`.
	 */
	customLabel?: boolean;
}
