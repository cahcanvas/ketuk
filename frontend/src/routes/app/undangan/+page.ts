import type { Event } from '@ketuk/shared';
import { listMyEvents } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, parent }) => {
	const { accessToken } = await parent();

	try {
		const events: Event[] = await listMyEvents({ fetch, accessToken });
		return { events, error: null };
	} catch {
		return { events: [] as Event[], error: 'Koneksi terputus. Coba muat ulang.' };
	}
};
