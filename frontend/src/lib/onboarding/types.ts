import type { UserEventPreference } from '@ketuk/shared';

/**
 * Port tunggal akses preferensi onboarding.
 *
 * Komponen UI dilarang membaca `localStorage` langsung. Semua baca/tulis harus
 * lewat sini, supaya implementasi penyimpanan bisa diganti (localStorage →
 * Supabase) tanpa menyentuh satu komponen pun — lihat local-store.ts dan
 * janji migrasi di REDESIGN-UI.md §4.4.
 */
export interface OnboardingStore {
	/**
	 * null = belum onboarding, atau storage tidak terbaca (SSR, diblokir, korup).
	 * Mengembalikan null, bukan throw, karena UI harus tetap bisa render dengan
	 * konfigurasi default — error di lapisan storage tidak boleh sampai ke user.
	 */
	get(): UserEventPreference | null;

	set(preference: UserEventPreference): void;

	/**
	 * Catat bahwa user sengaja melewati onboarding. Dipisah dari `set` karena dua
	 * hal yang berbeda: preference adalah jawaban, skipped adalah "jangan tanya
	 * lagi". Tanpa penanda ini, gate onboarding di layout akan menendang user
	 * kembali ke onboarding tanpa akhir.
	 */
	skip(): void;

	isSkipped(): boolean;

	clear(): void;
}
