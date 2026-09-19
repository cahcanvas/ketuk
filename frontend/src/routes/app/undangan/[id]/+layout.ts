import { error } from '@sveltejs/kit';
import { getEvent } from '$lib/api';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ fetch, parent, params }) => {
	const { accessToken } = await parent();

	try {
		const event = await getEvent(params.id, { fetch, accessToken });
		return { event };
	} catch {
		throw error(404, 'Undangan tidak ditemukan');
	}
};
