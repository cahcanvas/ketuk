import type { Guest } from '@ketuk/shared';
import type { GuestStat } from '$lib/api';
import { getGuestStats, listGuests } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, parent, params }) => {
	const { accessToken } = await parent();

	try {
		const [guests, stats] = await Promise.all([
			listGuests(params.id, { fetch, accessToken }),
			getGuestStats(params.id, { fetch, accessToken }),
		]);
		return { guests, stats, error: null };
	} catch {
		return {
			guests: [] as Guest[],
			stats: [] as GuestStat[],
			error: 'Koneksi terputus. Coba muat ulang.',
		};
	}
};
