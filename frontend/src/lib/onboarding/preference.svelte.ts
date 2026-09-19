import type { UserEventPreference } from '@ketuk/shared';
import { browser } from '$app/environment';
import { localStore } from './local-store';
import type { OnboardingStore } from './types';

/**
 * Satu-satunya jalan komponen UI ke preferensi onboarding.
 *
 * Berbasis Svelte 5 runes (bukan store lama) supaya reaktifitasnya mengikuti
 * komponen yang memakainya. Tidak dibuat singleton global: setiap halaman yang
 * butuh preferensi cukup panggil fungsi ini, karena Svelte 5 tidak butuh
 * context global untuk berbagi state reaktif antar komponen.
 *
 * Pembacaan awal di-guard `browser` karena `localStore` selalu mengembalikan
 * null di SSR — tanpa ini, hydration mismatch terjadi pada render pertama.
 */
export function createOnboardingPreference(store: OnboardingStore = localStore) {
	let current = $state<UserEventPreference | null>(browser ? store.get() : null);
	let skipped = $state<boolean>(browser ? store.isSkipped() : false);

	function set(preference: UserEventPreference): void {
		current = preference;
		store.set(preference);
	}

	function skip(): void {
		skipped = true;
		store.skip();
	}

	function clear(): void {
		current = null;
		skipped = false;
		store.clear();
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
