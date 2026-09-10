import { resolveMx } from 'node:dns/promises';

/**
 * Hasil pemeriksaan DNS domain email. Sengaja tiga nilai, bukan boolean:
 * "tidak bisa dipastikan" adalah kondisi nyata yang harus diperlakukan berbeda
 * dari "pasti tidak bisa menerima email".
 */
export type MailDomainVerdict = 'deliverable' | 'undeliverable' | 'unknown';

/**
 * Cache per-instance. Domain email sangat terkonsentrasi (gmail.com, yahoo.com,
 * dan beberapa domain kantor menutup mayoritas pendaftaran), jadi cache kecil ini
 * menghapus hampir semua DNS lookup berulang tanpa infrastruktur tambahan. Hilang
 * saat instance serverless didaur — itu bukan masalah, ini optimasi, bukan sumber
 * kebenaran.
 */
const verdictCache = new Map<string, MailDomainVerdict>();
const CACHE_LIMIT = 500;

/** Lookup DNS yang macet tidak boleh menggantung form pendaftaran. */
const DNS_TIMEOUT_MS = 2500;

/**
 * Kode error yang berarti "kita gagal bertanya ke DNS", bukan "DNS menjawab tidak
 * ada". Dicatat sekali per proses supaya kalau lapis pemeriksaan ini ternyata
 * tidak pernah berfungsi di production, itu terlihat di log — bukan diam-diam
 * lolos sebagai lapisan keamanan yang kita sangka aktif padahal tidak.
 */
let resolverFailureLogged = false;

function logResolverFailure(code: string | undefined) {
	if (resolverFailureLogged) return;
	resolverFailureLogged = true;
	console.warn(
		`[email-domain] Resolver DNS tidak bisa dipakai (${code ?? 'timeout'}). Verifikasi MX ` +
			'dilewati untuk semua pendaftaran di instance ini — perlindungan yang tersisa ' +
			'hanya blocklist domain di @ketuk/shared.',
	);
}

function withTimeout<T>(promise: Promise<T>, ms: number): Promise<T> {
	return Promise.race([
		promise,
		new Promise<never>((_, reject) =>
			setTimeout(() => reject(new Error('DNS lookup timeout')), ms),
		),
	]);
}

/**
 * Lapis kedua pemeriksaan "email asli", setelah blocklist domain sekali-pakai di
 * @ketuk/shared. Blocklist hanya tahu domain yang sudah kita daftarkan; ini
 * menangkap kelas masalah yang berbeda dan lebih umum: domain yang memang tidak
 * bisa menerima email sama sekali — salah ketik (`gmial.com`, `gmail.con`) dan
 * domain karangan yang tidak pernah ada.
 *
 * Kenapa MX record dan bukan sekadar keberadaan domain: domain bisa terdaftar dan
 * punya website tanpa pernah dikonfigurasi menerima email. Yang relevan di sini
 * adalah bisa-tidaknya menerima email, dan itu tepatnya yang dijawab MX.
 *
 * FAIL OPEN kalau lookup gagal karena alasan selain "domain/MX tidak ada":
 * timeout, resolver mati, atau jaringan bermasalah adalah masalah kita, bukan
 * masalah pengguna. Menolak pendaftaran karena resolver DNS sedang bermasalah
 * berarti memblokir pengguna sah demi memblokir pengguna nakal — pertukaran yang
 * salah arah untuk sebuah form pendaftaran.
 */
export async function verifyMailDomain(domain: string): Promise<MailDomainVerdict> {
	const normalized = domain.trim().toLowerCase();
	if (!normalized) return 'undeliverable';

	const cached = verdictCache.get(normalized);
	if (cached) return cached;

	let verdict: MailDomainVerdict;

	try {
		const records = await withTimeout(resolveMx(normalized), DNS_TIMEOUT_MS);
		// Array kosong atau MX "null" (exchange "." per RFC 7505) sama-sama berarti
		// domain ini secara eksplisit menyatakan tidak menerima email.
		const usable = records.filter((record) => record.exchange && record.exchange !== '.');
		verdict = usable.length > 0 ? 'deliverable' : 'undeliverable';
	} catch (error) {
		const code = (error as NodeJS.ErrnoException).code;
		// Hanya dua kode ini yang merupakan JAWABAN dari DNS, bukan kegagalan
		// menghubungi DNS: ENOTFOUND (nama domainnya tidak ada) dan ENODATA
		// (domainnya ada, tapi tanpa MX). Sisanya — ESERVFAIL, ETIMEOUT,
		// ECONNREFUSED, EAI_AGAIN, atau timeout kita sendiri di atas — berarti kita
		// yang gagal bertanya, dan itu tidak boleh berujung pada penolakan pengguna.
		if (code === 'ENOTFOUND' || code === 'ENODATA') {
			verdict = 'undeliverable';
		} else {
			logResolverFailure(code);
			verdict = 'unknown';
		}
	}

	// Verdict 'unknown' tidak di-cache: itu kegagalan sementara, bukan fakta tentang
	// domainnya. Meng-cache-nya berarti satu gangguan DNS sesaat ikut mempengaruhi
	// pendaftar berikutnya.
	if (verdict !== 'unknown') {
		if (verdictCache.size >= CACHE_LIMIT) {
			const oldest = verdictCache.keys().next().value;
			if (oldest) verdictCache.delete(oldest);
		}
		verdictCache.set(normalized, verdict);
	}

	return verdict;
}
