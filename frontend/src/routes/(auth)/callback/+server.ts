import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

/** Tujuan redirect Supabase (magic link & OAuth) setelah user klik link/setuju di provider. */
export const GET: RequestHandler = async ({ url, locals }) => {
	const code = url.searchParams.get('code');
	const next = url.searchParams.get('next') ?? '/app/dashboard';

	if (code) {
		const { error } = await locals.supabase.auth.exchangeCodeForSession(code);
		if (!error) {
			throw redirect(303, next);
		}
	}

	throw redirect(303, '/masuk?error=auth');
};
