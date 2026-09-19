import type { Event } from '@ketuk/shared';
import { listMyEvents } from '$lib/api';
import type { PageLoad } from './$types';

/**
 * Dashboard butuh daftar acara user: untuk menampilkan kartunya, dan sebagai
 * sinyal "sudah punya konteks sendiri" (empty state vs. list). Ambil lewat
 * backend Go seperti halaman lain — bukan Supabase langsung — supaya otorisasi
 * tetap satu jalan.
 */
export const load: PageLoad = async ({ fetch, parent }) => {
	const { accessToken } = await parent();

	try {
		const events: Event[] = await listMyEvents({ fetch, accessToken });
		return { events, error: null };
	} catch {
		return { events: [] as Event[], error: 'Koneksi terputus. Coba muat ulang.' };
	}
};
