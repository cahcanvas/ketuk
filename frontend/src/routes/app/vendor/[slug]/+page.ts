import { error } from '@sveltejs/kit';
import { getVendor } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
	try {
		const vendor = await getVendor(params.slug, { fetch });
		return { vendor };
	} catch {
		throw error(404, 'Vendor tidak ditemukan');
	}
};
