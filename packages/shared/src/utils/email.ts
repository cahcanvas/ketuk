import { DISPOSABLE_EMAIL_DOMAINS } from '../constants/disposable-email-domains';

const disposableDomains = new Set(DISPOSABLE_EMAIL_DOMAINS);

/**
 * Titik normalisasi tunggal untuk email, sejajar dengan normalizeIndonesianPhone
 * di schemas/gift.ts. Bagian domain di-lowercase karena memang tidak
 * case-sensitive menurut spesifikasi; bagian lokal juga di-lowercase karena
 * semua provider besar memperlakukannya begitu, dan membiarkannya apa adanya
 * hanya akan membuat "Budi@gmail.com" dan "budi@gmail.com" jadi dua akun berbeda.
 */
export function normalizeEmail(email: string): string {
	return email.trim().toLowerCase();
}

/** Bagian setelah "@" terakhir, sudah di-lowercase. String kosong kalau email tidak punya domain. */
export function getEmailDomain(email: string): string {
	const at = email.lastIndexOf('@');
	if (at === -1) return '';
	return email
		.slice(at + 1)
		.trim()
		.toLowerCase();
}

/**
 * Cocokkan domain beserta semua induknya, jadi satu entri `mailinator.com` di
 * blocklist otomatis menolak `team.mailinator.com` juga — layanan temp mail
 * lazim membagikan subdomain tak terbatas untuk menghindari daftar seperti ini.
 */
export function isDisposableEmailDomain(domain: string): boolean {
	const normalized = domain.trim().toLowerCase();
	if (!normalized) return false;

	const parts = normalized.split('.');
	for (let i = 0; i < parts.length - 1; i++) {
		if (disposableDomains.has(parts.slice(i).join('.'))) return true;
	}

	return false;
}

export function isDisposableEmail(email: string): boolean {
	return isDisposableEmailDomain(getEmailDomain(email));
}
