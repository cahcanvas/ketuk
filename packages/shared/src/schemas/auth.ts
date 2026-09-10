import { z } from 'zod';
import { PROVINCE_CODES } from '../constants/provinces';
import { getEmailDomain, isDisposableEmailDomain, normalizeEmail } from '../utils/email';

/**
 * Nama asli, bukan username. Yang divalidasi di sini adalah BENTUK-nya, bukan
 * kebenarannya — tidak ada cara memverifikasi seseorang benar bernama "Budi"
 * tanpa dokumen identitas. Yang bisa dilakukan: menolak bentuk yang jelas bukan
 * nama orang (angka, alamat email, URL, simbol acak), karena nama ini akan
 * muncul di undangan yang dibagikan ke ratusan tamu.
 *
 * Tidak mewajibkan dua kata: banyak orang Indonesia memang bernama satu kata.
 */
export const fullNameSchema = z
	.string()
	.trim()
	// Rapikan spasi ganda/tab sebelum divalidasi panjangnya, supaya "Budi    Santoso"
	// tidak lolos sebagai nama yang "lebih panjang" dari yang sebenarnya.
	.transform((value) => value.replace(/\s+/g, ' '))
	.pipe(
		z
			.string()
			.min(3, 'Nama minimal 3 karakter')
			.max(100, 'Nama maksimal 100 karakter')
			.regex(
				/^[\p{L}][\p{L}\s'.,-]*$/u,
				'Nama hanya boleh berisi huruf, spasi, dan tanda petik/titik/koma/hubung',
			)
			.refine(
				(value) => (value.match(/\p{L}/gu) ?? []).length >= 3,
				'Nama minimal 3 huruf, bukan hanya tanda baca',
			),
	);

/**
 * Email asli. Tiga lapis pemeriksaan, dari yang paling murah ke yang paling mahal:
 *
 * 1. Bentuk email valid (di sini).
 * 2. Domain bukan penyedia email sekali-pakai (di sini, lewat blocklist).
 * 3. Domain benar-benar bisa menerima email — verifikasi MX record, dilakukan
 *    di sisi server saat pendaftaran karena butuh DNS lookup dan tidak boleh
 *    dikerjakan di browser.
 *
 * Lapis 2 di sini penting untuk dijalankan ULANG di server, bukan cuma di
 * browser: validasi client bisa dilewati siapa saja dengan mengedit request.
 */
export const emailSchema = z
	.string()
	.trim()
	.min(1, 'Email wajib diisi')
	.max(254, 'Email terlalu panjang')
	.email('Format email tidak valid')
	.transform(normalizeEmail)
	.refine(
		(value) => !isDisposableEmailDomain(getEmailDomain(value)),
		'Email sekali-pakai (temporary email) tidak bisa dipakai. Gunakan email utama kamu — email ini yang dipakai untuk memulihkan akun dan mengurus pembayaran.',
	);

/**
 * Batas 72 karakter bukan pilihan desain kami: bcrypt (yang dipakai Supabase Auth
 * untuk hashing) memotong input di byte ke-72 secara diam-diam. Membiarkan
 * password lebih panjang berarti membiarkan pengguna percaya bagian akhir
 * password-nya berarti padahal diabaikan.
 *
 * Syaratnya sengaja ringan (huruf + angka, 8 karakter) supaya pendaftaran tetap
 * cepat. Kekuatan password yang lebih tinggi didorong lewat indikator visual di
 * form, bukan lewat aturan yang membuat orang menyerah di tengah jalan.
 */
export const passwordSchema = z
	.string()
	.min(8, 'Password minimal 8 karakter')
	.max(72, 'Password maksimal 72 karakter')
	.regex(/\p{L}/u, 'Password harus mengandung minimal 1 huruf')
	.regex(/\d/, 'Password harus mengandung minimal 1 angka');

export const provinceSchema = z.enum(PROVINCE_CODES, {
	message: 'Pilih asal provinsi kamu',
});

export const signUpSchema = z
	.object({
		name: fullNameSchema,
		email: emailSchema,
		/** Asal provinsi host — disimpan sebagai kode stabil, lihat constants/provinces.ts. */
		province: provinceSchema,
		password: passwordSchema,
		confirmPassword: z.string().min(1, 'Konfirmasi password wajib diisi'),
	})
	.superRefine((data, ctx) => {
		if (data.password !== data.confirmPassword) {
			ctx.addIssue({
				code: 'custom',
				path: ['confirmPassword'],
				message: 'Konfirmasi password tidak sama',
			});
		}

		// Password yang isinya cuma bagian depan email-nya sendiri praktis sama
		// dengan tidak punya password — itu bagian yang paling mudah ditebak orang
		// yang sudah tahu emailnya (dan emailnya memang dibagikan ke tamu).
		const localPart = data.email.split('@')[0] ?? '';
		if (localPart.length >= 4 && data.password.toLowerCase().includes(localPart.toLowerCase())) {
			ctx.addIssue({
				code: 'custom',
				path: ['password'],
				message: 'Password tidak boleh mengandung bagian dari email kamu',
			});
		}

		// Ambang 4 karakter supaya nama pendek tidak memicu false positive — tanpa
		// itu, nama "Ali" akan menolak password yang kebetulan memuat "kualitas".
		const compactName = data.name.toLowerCase().replace(/\s/g, '');
		if (compactName.length >= 4 && data.password.toLowerCase().includes(compactName)) {
			ctx.addIssue({
				code: 'custom',
				path: ['password'],
				message: 'Password tidak boleh mengandung nama kamu',
			});
		}
	});

export type SignUpInput = z.infer<typeof signUpSchema>;

export const signInSchema = z.object({
	email: z
		.string()
		.trim()
		.min(1, 'Email wajib diisi')
		.email('Format email tidak valid')
		.transform(normalizeEmail),
	/**
	 * Sengaja TIDAK memakai passwordSchema. Aturan password bisa berubah/diperketat
	 * nanti; kalau login ikut memakai aturan yang sama, akun lama yang password-nya
	 * dibuat sebelum aturan itu akan ditolak di form login sebelum sempat
	 * diverifikasi ke Supabase — dan pesan errornya akan membocorkan aturan
	 * password ke penebak. Login cuma perlu tahu field-nya tidak kosong.
	 */
	password: z.string().min(1, 'Password wajib diisi'),
});

export type SignInInput = z.infer<typeof signInSchema>;

export interface PasswordStrength {
	/** 0-4. 0 = kosong/sangat lemah, 4 = kuat. */
	score: 0 | 1 | 2 | 3 | 4;
	label: string;
}

/**
 * Indikator kekuatan password untuk umpan balik visual di form pendaftaran.
 * Ini murni panduan untuk pengguna — BUKAN gerbang validasi. Yang menentukan
 * password diterima atau tidak tetap passwordSchema di atas, supaya tidak ada
 * dua sumber kebenaran soal "password ini boleh atau tidak".
 */
export function getPasswordStrength(password: string): PasswordStrength {
	if (!password) return { score: 0, label: 'Belum diisi' };

	let points = 0;
	if (password.length >= 8) points++;
	if (password.length >= 12) points++;
	if (/\p{Ll}/u.test(password) && /\p{Lu}/u.test(password)) points++;
	if (/\d/.test(password)) points++;
	if (/[^\p{L}\d]/u.test(password)) points++;

	const score = Math.min(4, Math.max(1, points - 1)) as 1 | 2 | 3 | 4;
	const labels: Record<1 | 2 | 3 | 4, string> = {
		1: 'Lemah',
		2: 'Cukup',
		3: 'Bagus',
		4: 'Kuat',
	};

	return { score, label: labels[score] };
}
