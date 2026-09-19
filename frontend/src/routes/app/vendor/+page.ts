import type { VendorCategory } from '@ketuk/shared';
import type { ListVendorsResult } from '$lib/api';
import { listVendors } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, url }) => {
	const category = url.searchParams.get('category') ?? undefined;

	try {
		const result = await listVendors(
			{ category: category as VendorCategory | undefined },
			{ fetch },
		);
		return { result, category, error: null };
	} catch {
		return {
			result: { items: [], nextCursor: null } as ListVendorsResult,
			category,
			error: 'Koneksi terputus. Coba muat ulang.',
		};
	}
};
