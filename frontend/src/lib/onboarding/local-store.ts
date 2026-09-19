import { isEventCategoryId, type UserEventPreference } from '@ketuk/shared';
import { browser } from '$app/environment';
import type { OnboardingStore } from './types';

/**
 * Implementasi sementara: browser localStorage. Dipilih sebagai langkah pertama
 * karena tidak butuh perubahan backend, tetapi semua pemanggilan lewat port
 * `OnboardingStore` — saat nanti dipindah ke Supabase, file ini diganti
 * implementasinya saja, bukan kontraknya.
 */

const PREFERENCE_KEY = 'ketuk.onboardingPreference';
const SKIPPED_KEY = 'ketuk.onboardingSkipped';

/**
 * Bentuk mentah yang mungkin tersimpan di localStorage. Ditulis longgar karena
 * isinya bisa datang dari versi aplikasi sebelumnya — bukan hanya dari versi
 * sekarang. Tugas fungsi ini menolak apa pun yang tidak terbukti valid.
 */
interface StoredPreference {
	eventType?: unknown;
	eventSubType?: unknown;
	purpose?: unknown;
	selectedAt?: unknown;
}

function readRaw(key: string): string | null {
	try {
		return localStorage.getItem(key);
	} catch {
		// Safari private mode bisa throw saat storage diblokir. Anggap saja kosong
		// — UX yang benar di sini adalah "bertindak seolah belum onboarding",
		// bukan meledakkan halaman.
		return null;
	}
}

function writeRaw(key: string, value: string): boolean {
	try {
		localStorage.setItem(key, value);
		return true;
	} catch {
		return false;
	}
}

function removeRaw(key: string): void {
	try {
		localStorage.removeItem(key);
	} catch {
		// Sama seperti di atas: tidak ada yang bisa dilakukan, diam saja.
	}
}

/**
 * Parse + validasi. Mengembalikan null untuk apa pun yang tidak dikenali:
 * JSON rusak, bentuk setengah, atau `eventType` dari versi lama yang sudah
 * tidak ada di taxonomy. Lebih aman "tidak tahu" daripada memakai nilai usang.
 */
function parsePreference(raw: string | null): UserEventPreference | null {
	if (!raw) return null;

	let parsed: StoredPreference;
	try {
		parsed = JSON.parse(raw) as StoredPreference;
	} catch {
		return null;
	}

	if (typeof parsed !== 'object' || parsed === null) return null;

	const { eventType, eventSubType, purpose, selectedAt } = parsed;

	if (!isEventCategoryId(eventType)) return null;
	if (eventSubType !== undefined && typeof eventSubType !== 'string') return null;
	if (purpose !== undefined && typeof purpose !== 'string') return null;
	if (typeof selectedAt !== 'string' || !selectedAt) return null;

	return {
		eventType,
		eventSubType: typeof eventSubType === 'string' ? eventSubType : undefined,
		purpose: typeof purpose === 'string' ? purpose : undefined,
		selectedAt,
	};
}

export const localStore: OnboardingStore = {
	get(): UserEventPreference | null {
		if (!browser) return null;
		return parsePreference(readRaw(PREFERENCE_KEY));
	},

	set(preference: UserEventPreference): void {
		if (!browser) return;
		writeRaw(PREFERENCE_KEY, JSON.stringify(preference));
	},

	skip(): void {
		if (!browser) return;
		writeRaw(SKIPPED_KEY, new Date().toISOString());
	},

	isSkipped(): boolean {
		if (!browser) return false;
		// Timestamp hanya penanda "pernah skip"; tidak ada logika yang
		// bergantung pada nilainya, asal bukan string kosong.
		return readRaw(SKIPPED_KEY) !== null;
	},

	clear(): void {
		if (!browser) return;
		removeRaw(PREFERENCE_KEY);
		removeRaw(SKIPPED_KEY);
	},
};
