import type { PageServerLoad } from './$types';

export interface GiftOrderSafeRow {
	id: string;
	event_id: string;
	product_id: string;
	sender_name: string;
	recipient_name: string;
	recipient_address: string;
	quantity: number;
	total_amount: number;
	message: string | null;
	status: string;
	created_at: string;
}

/**
 * Query langsung ke view `gift_orders_safe` (backend/sql/rls-policies.sql) lewat
 * Supabase client — bukan lewat backend Hono. View ini sudah menyaring kolom
 * sensitif (payment_id, nomor HP pengirim) dan baris ke event milik user ini
 * saja, jadi query di sini aman dilakukan langsung.
 */
export const load: PageServerLoad = async ({ locals }) => {
	const { data, error } = await locals.supabase
		.from('gift_orders_safe')
		.select('*')
		.order('created_at', { ascending: false })
		.limit(100);

	if (error) {
		return { orders: [] as GiftOrderSafeRow[], error: 'Gagal memuat data hadiah.' };
	}

	return { orders: (data ?? []) as GiftOrderSafeRow[], error: null };
};
