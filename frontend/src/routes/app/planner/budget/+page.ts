import type { BudgetItem } from '@ketuk/shared';
import { redirect } from '@sveltejs/kit';
import { listBudgetItems } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, parent, url }) => {
	const { accessToken } = await parent();
	const eventId = url.searchParams.get('event');
	if (!eventId) throw redirect(303, '/app/planner');

	try {
		const items = await listBudgetItems(eventId, { fetch, accessToken });
		return { items, eventId, error: null };
	} catch {
		return { items: [] as BudgetItem[], eventId, error: 'Koneksi terputus. Coba muat ulang.' };
	}
};
