import { getEmailDomain, signUpSchema } from '@ketuk/shared';
import { fail, redirect } from '@sveltejs/kit';
import { verifyMailDomain } from '$lib/server/email-domain';
import type { Actions } from './$types';

type SignUpField = 'name' | 'email' | 'province' | 'password' | 'confirmPassword';

type FieldErrors = Partial<Record<SignUpField, string>>;

/**
 * Field yang dikirim balik ke form saat validasi gagal, supaya user tidak perlu
 * mengisi ulang semuanya. Password SENGAJA tidak termasuk — mengirim password
 * kembali ke browser berarti ia ikut masuk ke HTML respons, riwayat cache, dan
 * potensi log perantara. Mengetik ulang password jauh lebih murah daripada itu.
 */
interface SubmittedValues {
	name: string;
	email: string;
	province: string;
}

/**
 * SEMUA cabang — gagal maupun berhasil — mengembalikan bentuk yang sama.
 *
 * Ini bukan sekadar kerapian: kalau bentuknya berbeda per cabang, tipe ActionData
 * yang digenerate SvelteKit menjadi union, dan komponen tidak bisa membaca
 * `form.errors.email` sebelum mempersempit union itu dulu. Satu bentuk seragam
 * menghilangkan seluruh kelas masalah itu di sisi komponen.
 */
interface SignUpActionResult {
	values: SubmittedValues;
	errors: FieldErrors;
	formError: string | null;
	emailConfirmationRequired: boolean;
}

function invalid(
	status: number,
	values: SubmittedValues,
	errors: FieldErrors = {},
	formError: string | null = null,
) {
	const result: SignUpActionResult = {
		values,
		errors,
		formError,
		emailConfirmationRequired: false,
	};
	return fail(status, result);
}

export const actions: Actions = {
	/**
	 * Pendaftaran dijalankan sebagai form action di server, bukan lewat Supabase
	 * client di browser, karena satu alasan yang menentukan: penolakan email
	 * sekali-pakai harus berjalan di tempat yang tidak bisa dilewati pengguna.
	 * Validasi yang hanya hidup di browser bisa diakali dengan satu request manual.
	 *
	 * Efek sampingnya sama-sama menguntungkan: form ini tetap berfungsi tanpa
	 * JavaScript, dan cookie session diset langsung dari respons server.
	 */
	default: async ({ request, locals, url }) => {
		const formData = await request.formData();

		const raw = {
			name: String(formData.get('name') ?? ''),
			email: String(formData.get('email') ?? ''),
			province: String(formData.get('province') ?? ''),
			password: String(formData.get('password') ?? ''),
			confirmPassword: String(formData.get('confirmPassword') ?? ''),
		};

		const values: SubmittedValues = {
			name: raw.name,
			email: raw.email,
			province: raw.province,
		};

		const parsed = signUpSchema.safeParse(raw);

		if (!parsed.success) {
			const errors: FieldErrors = {};
			for (const issue of parsed.error.issues) {
				const field = issue.path[0];
				// Pesan pertama per field saja. Satu field yang melanggar tiga aturan
				// sekaligus lebih mudah diperbaiki kalau ditunjukkan bertahap daripada
				// dihujani tiga pesan sekaligus.
				if (typeof field === 'string' && !(field in errors)) {
					errors[field as SignUpField] = issue.message;
				}
			}
			return invalid(400, values, errors);
		}

		const { name, email, province, password } = parsed.data;

		// Lapis kedua pemeriksaan email asli: domainnya benar-benar bisa menerima
		// email? Menangkap salah ketik (`gmial.com`) dan domain karangan yang lolos
		// blocklist. Lihat $lib/server/email-domain.ts untuk perilaku fail-open-nya.
		if ((await verifyMailDomain(getEmailDomain(email))) === 'undeliverable') {
			return invalid(400, values, {
				email: 'Domain email ini tidak bisa menerima email. Periksa lagi ejaannya.',
			});
		}

		const { data, error } = await locals.supabase.auth.signUp({
			email,
			password,
			options: {
				// Disimpan di user_metadata dulu, lalu dipindahkan ke tabel `profiles`
				// saat pertama kali membuka dashboard — lihat (app)/+layout.server.ts.
				data: { name, province },
			},
		});

		if (error) {
			// Menyebut "email ini sudah terdaftar" memang memungkinkan orang menguji
			// apakah sebuah email punya akun di sini. Itu diterima dengan sadar: form
			// login di sebelah sudah membocorkan hal yang sama, dan alternatifnya
			// (pesan samar) membuat orang yang lupa pernah mendaftar terjebak mencoba
			// berulang kali tanpa tahu bahwa yang ia butuhkan adalah halaman masuk.
			if (error.code === 'user_already_exists' || error.message.includes('already registered')) {
				return invalid(409, values, {
					email: 'Email ini sudah terdaftar. Silakan masuk, atau pakai email lain.',
				});
			}

			if (error.code === 'weak_password') {
				return invalid(400, values, {
					password: 'Password terlalu mudah ditebak. Coba yang lebih panjang.',
				});
			}

			if (error.code === 'email_address_invalid') {
				return invalid(400, values, { email: 'Email ini ditolak sistem. Coba email lain.' });
			}

			if (error.status === 429) {
				return invalid(
					429,
					values,
					{},
					'Terlalu banyak percobaan pendaftaran. Tunggu beberapa menit, lalu coba lagi.',
				);
			}

			// Detail error asli hanya ke log server — pesan mentah dari Supabase bisa
			// memuat informasi internal yang tidak perlu dilihat pengguna.
			console.error('[daftar] signUp gagal', { code: error.code, status: error.status });
			return invalid(500, values, {}, 'Gagal membuat akun. Coba lagi sebentar lagi.');
		}

		/**
		 * Konfirmasi email dimatikan di project Supabase, jadi signUp langsung
		 * mengembalikan session dan cookie-nya sudah diset oleh client SSR.
		 *
		 * Cabang `!data.session` tetap ada, dan bukan kode mati: kalau nanti
		 * konfirmasi email diaktifkan di dashboard Supabase, halaman ini langsung
		 * berperilaku benar (memberi tahu user untuk cek inbox) tanpa perlu diubah —
		 * bukan diam-diam mengarahkan ke dashboard yang lalu menendangnya balik ke
		 * halaman masuk.
		 */
		if (!data.session) {
			const result: SignUpActionResult = {
				values,
				errors: {},
				formError: null,
				emailConfirmationRequired: true,
			};
			return result;
		}

		const next = url.searchParams.get('next');
		// Hanya path internal. `next` datang dari URL yang bisa disusun siapa saja,
		// dan tanpa pemeriksaan ini ia jadi celah open redirect ke domain luar.
		const target = next?.startsWith('/') && !next.startsWith('//') ? next : '/dashboard';

		throw redirect(303, target);
	},
};
