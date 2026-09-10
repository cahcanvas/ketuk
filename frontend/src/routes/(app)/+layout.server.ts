import { isProvinceCode } from '@ketuk/shared';
import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

/** Ambil string yang tidak kosong dari user_metadata, kalau bukan string kembalikan null. */
function readMetadataString(value: unknown): string | null {
	return typeof value === 'string' && value.trim() ? value.trim() : null;
}

/**
 * Guard dashboard: redirect ke /masuk kalau belum login, bawa `?next=` supaya
 * setelah login kembali ke halaman yang dituju.
 *
 * Juga memastikan baris `profiles` ada untuk user ini (just-in-time). Policy
 * RLS `profiles_insert_own` (backend/sql/rls-policies.sql) mengizinkan user
 * insert profilnya sendiri, jadi ini aman dilakukan langsung dari sini lewat
 * Supabase client tanpa perlu endpoint backend khusus.
 *
 * Kenapa nama & asal provinsi menempuh jalan memutar lewat user_metadata dan
 * baru mendarat di `profiles` di sini, bukan langsung di-insert saat mendaftar:
 * tabel `profiles` punya foreign key ke `auth.users`, jadi barisnya tidak bisa
 * ada sebelum user-nya ada. Menuliskannya di sini — satu tempat yang dilewati
 * semua jalur masuk — berarti akun dari form pendaftaran maupun dari Google
 * OAuth sama-sama berakhir dengan profil yang lengkap, tanpa logika yang
 * digandakan di dua tempat.
 */
export const load: LayoutServerLoad = async ({ locals, url }) => {
	const { session, user } = await locals.safeGetSession();

	if (!session || !user) {
		throw redirect(303, `/masuk?next=${encodeURIComponent(url.pathname)}`);
	}

	const { data: profile } = await locals.supabase
		.from('profiles')
		.select('id, name, province')
		.eq('id', user.id)
		.maybeSingle();

	const metadataName = readMetadataString(user.user_metadata?.name);
	const metadataProvince = readMetadataString(user.user_metadata?.province);
	// Kode yang tidak dikenal diabaikan, bukan disimpan apa adanya — user_metadata
	// bisa diisi dari mana saja, dan kolom ini tidak dijaga enum di database.
	const province = isProvinceCode(metadataProvince) ? metadataProvince : null;

	if (!profile) {
		await locals.supabase.from('profiles').insert({
			id: user.id,
			name: metadataName ?? user.email ?? 'Pengguna Ketuk',
			province,
		});
	} else if (!profile.province && province) {
		// Akun yang mendaftar sebelum kolom `province` ada, atau yang baru mengisinya
		// belakangan lewat OAuth — diisi sekali di sini, lalu tidak diutak-utik lagi.
		// Sengaja tidak menimpa nilai yang sudah ada: kalau user mengubah asalnya di
		// halaman profil nanti, user_metadata lama tidak boleh mengembalikannya.
		await locals.supabase.from('profiles').update({ province }).eq('id', user.id);
	}

	return {
		session,
		user,
		accessToken: session.access_token,
	};
};
