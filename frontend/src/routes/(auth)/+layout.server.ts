import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

/** Sudah login? Tidak perlu lihat form masuk/daftar lagi. */
export const load: LayoutServerLoad = async ({ locals }) => {
	const { session } = await locals.safeGetSession();

	if (session) {
		throw redirect(303, '/app/dashboard');
	}

	return {};
};
