import { signInSchema } from '@ketuk/shared';
import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

type FieldErrors = Partial<Record<'email' | 'password', string>>;

/** Bentuk seragam untuk semua cabang gagal — alasannya sama seperti di /daftar. */
interface SignInActionResult {
	values: { email: string };
	errors: FieldErrors;
	formError: string | null;
}

function invalid(
	status: number,
	email: string,
	errors: FieldErrors = {},
	formError: string | null = null,
) {
	const result: SignInActionResult = { values: { email }, errors, formError };
	return fail(status, result);
}

export const actions: Actions = {
	/**
	 * Login dengan email + password, sepasang dengan form pendaftaran di /daftar.
	 *
	 * Magic link yang sebelumnya dipakai di sini dilepas karena tidak ada lagi
	 * jalur pendaftaran yang menghasilkannya, dan karena pengiriman email memang
	 * belum dikonfigurasi — menyisakannya berarti menyisakan tombol yang tidak
	 * pernah menghasilkan email.
	 */
	default: async ({ request, locals, url }) => {
		const formData = await request.formData();
		const rawEmail = String(formData.get('email') ?? '');
		const rawPassword = String(formData.get('password') ?? '');

		const parsed = signInSchema.safeParse({ email: rawEmail, password: rawPassword });

		if (!parsed.success) {
			const errors: FieldErrors = {};
			for (const issue of parsed.error.issues) {
				const field = issue.path[0];
				if ((field === 'email' || field === 'password') && !(field in errors)) {
					errors[field] = issue.message;
				}
			}
			return invalid(400, rawEmail, errors);
		}

		const { error } = await locals.supabase.auth.signInWithPassword(parsed.data);

		if (error) {
			if (error.status === 429) {
				return invalid(
					429,
					rawEmail,
					{},
					'Terlalu banyak percobaan masuk. Tunggu beberapa menit, lalu coba lagi.',
				);
			}

			// Pesan tunggal untuk "email tidak ada" dan "password salah" — memisahkan
			// keduanya akan mengubah form ini menjadi alat untuk mengecek email siapa
			// saja yang punya akun di sini, dan itu tidak memberi manfaat apa pun ke
			// pengguna sah yang sudah tahu email dan password-nya sendiri.
			return invalid(400, rawEmail, {}, 'Email atau password salah.');
		}

		const next = url.searchParams.get('next');
		// Hanya path internal — tanpa pemeriksaan ini, `next` jadi celah open redirect.
		const target = next?.startsWith('/') && !next.startsWith('//') ? next : '/dashboard';

		throw redirect(303, target);
	},
};
