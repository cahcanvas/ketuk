import type { WishItem } from '$lib/api';
import { listWishes } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, parent }) => {
	const { event } = await parent();

	try {
		const wishes = await listWishes(event.slug, { fetch });
		return { wishes, error: null };
	} catch {
		return { wishes: [] as WishItem[], error: 'Koneksi terputus. Coba muat ulang.' };
	}
};
