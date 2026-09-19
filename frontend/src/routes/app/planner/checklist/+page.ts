import type { ChecklistItem } from '@ketuk/shared';
import { redirect } from '@sveltejs/kit';
import { listChecklistItems } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, parent, url }) => {
	const { accessToken } = await parent();
	const eventId = url.searchParams.get('event');
	if (!eventId) throw redirect(303, '/app/planner');

	try {
		const items = await listChecklistItems(eventId, { fetch, accessToken });
		return { items, eventId, error: null };
	} catch {
		return { items: [] as ChecklistItem[], eventId, error: 'Koneksi terputus. Coba muat ulang.' };
	}
};
