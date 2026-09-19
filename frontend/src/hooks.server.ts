import { type Handle, redirect } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';
import { createSupabaseServerClient } from '$lib/supabase/server';

const BASE_HOSTNAMES = new Set(['ketuk.id', 'www.ketuk.id', 'localhost', '127.0.0.1']);
const APEX_SUFFIX = '.ketuk.id';

/**
 * Path aplikasi lama (sebelum area app dipindah ke bawah segment `/app/`).
 * Link yang sudah tersebar di WhatsApp/email tidak boleh mati, jadi semua
 * path di bawah ini dialihkan permanen (308) ke pasangannya di `/app/`.
 *
 * Hanya berjalan untuk base hostname — di subdomain undangan
 * (`budi-sinta.ketuk.id`) pathname milik tamu, bukan area aplikasi.
 */
const LEGACY_APP_PREFIXES = ['/dashboard', '/undangan', '/planner', '/vendor', '/hadiah'];

const handleLegacyAppPaths: Handle = async ({ event, resolve }) => {
	const { pathname } = event.url;

	if (BASE_HOSTNAMES.has(event.url.hostname)) {
		const matched = LEGACY_APP_PREFIXES.some(
			(prefix) => pathname === prefix || pathname.startsWith(`${prefix}/`),
		);
		if (matched) {
			throw redirect(308, `/app${pathname}`);
		}
	}

	return resolve(event);
};

/**
 * Dua bentuk URL undangan didukung: path-based (ketuk.id/budi-sinta, default,
 * semua plan) dan subdomain (budi-sinta.ketuk.id, khusus plan Lengkap).
 * Subdomain di-rewrite ke route publik `/[slug]` di sini — browser tetap
 * melihat subdomain aslinya di address bar, cuma routing internalnya yang diarahkan.
 */
const handleSubdomain: Handle = async ({ event, resolve }) => {
	const hostname = event.url.hostname;

	if (!BASE_HOSTNAMES.has(hostname) && hostname.endsWith(APEX_SUFFIX)) {
		const slug = hostname.slice(0, -APEX_SUFFIX.length);
		if (slug && slug !== 'www') {
			const suffix = event.url.pathname === '/' ? '' : event.url.pathname;
			event.url.pathname = `/${slug}${suffix}`;
		}
	}

	return resolve(event);
};

/**
 * Setup standar @supabase/ssr: buat client per-request dan expose safeGetSession.
 * `safeGetSession` memvalidasi JWT ke server Supabase (bukan cuma decode cookie),
 * jadi aman dipakai sebagai sumber kebenaran soal siapa yang login.
 */
const handleSupabase: Handle = async ({ event, resolve }) => {
	event.locals.supabase = createSupabaseServerClient(event.cookies);

	event.locals.safeGetSession = async () => {
		const {
			data: { session },
		} = await event.locals.supabase.auth.getSession();

		if (!session) {
			return { session: null, user: null };
		}

		const {
			data: { user },
			error,
		} = await event.locals.supabase.auth.getUser();

		if (error) {
			return { session: null, user: null };
		}

		return { session, user };
	};

	return resolve(event, {
		filterSerializedResponseHeaders: (name) =>
			name === 'content-range' || name === 'x-supabase-api-version',
	});
};

export const handle: Handle = sequence(handleLegacyAppPaths, handleSubdomain, handleSupabase);
