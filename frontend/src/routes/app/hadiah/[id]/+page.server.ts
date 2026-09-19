import { error } from '@sveltejs/kit';
import type { GiftOrderSafeRow } from '../+page.server';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, params }) => {
	const { data, error: dbError } = await locals.supabase
		.from('gift_orders_safe')
		.select('*')
		.eq('id', params.id)
		.maybeSingle();

	if (dbError || !data) {
		throw error(404, 'Hadiah tidak ditemukan');
	}

	return { order: data as GiftOrderSafeRow };
};
