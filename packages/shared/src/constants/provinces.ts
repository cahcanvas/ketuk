/**
 * 38 provinsi Indonesia (kondisi setelah pemekaran Papua 2022: Papua Selatan,
 * Papua Tengah, Papua Pegunungan, dan Papua Barat Daya).
 *
 * `code` yang disimpan di database, bukan `label`. Nama provinsi bisa berubah
 * secara administratif (mis. DKI Jakarta yang akan menjadi Daerah Khusus
 * Jakarta) — kalau yang tersimpan adalah label, satu perubahan nama akan
 * memecah data lama jadi dua nilai berbeda yang sebenarnya provinsi yang sama.
 * Dengan `code` stabil, ganti nama cuma perlu mengubah label di file ini.
 */
export const PROVINCE_ISLANDS = [
	'Sumatera',
	'Jawa',
	'Bali & Nusa Tenggara',
	'Kalimantan',
	'Sulawesi',
	'Maluku',
	'Papua',
] as const;

export type ProvinceIsland = (typeof PROVINCE_ISLANDS)[number];

export interface ProvinceConfig {
	code: string;
	label: string;
	/** Pulau/kepulauan induk — dipakai untuk mengelompokkan opsi di dropdown. */
	island: ProvinceIsland;
}

/**
 * Urutan sengaja geografis barat ke timur (Sumatera → Papua), bukan alfabetis.
 * Ini urutan yang sama dengan yang dipakai BPS dan yang paling familiar untuk
 * orang Indonesia saat mencari daerahnya sendiri di daftar panjang.
 *
 * `as const satisfies` — bukan cuma anotasi tipe biasa — supaya ProvinceCode di
 * bawah bisa diturunkan sebagai union literal dari daftar ini. Efeknya: daftar
 * ini jadi satu-satunya sumber kebenaran, dan salah tulis kode di tempat lain
 * (mis. seed data atau test) ketahuan saat type-check, bukan saat runtime.
 */
export const INDONESIA_PROVINCES = [
	// ── Sumatera ──
	{ code: 'aceh', label: 'Aceh', island: 'Sumatera' },
	{ code: 'sumatera-utara', label: 'Sumatera Utara', island: 'Sumatera' },
	{ code: 'sumatera-barat', label: 'Sumatera Barat', island: 'Sumatera' },
	{ code: 'riau', label: 'Riau', island: 'Sumatera' },
	{ code: 'kepulauan-riau', label: 'Kepulauan Riau', island: 'Sumatera' },
	{ code: 'jambi', label: 'Jambi', island: 'Sumatera' },
	{ code: 'bengkulu', label: 'Bengkulu', island: 'Sumatera' },
	{ code: 'sumatera-selatan', label: 'Sumatera Selatan', island: 'Sumatera' },
	{ code: 'kepulauan-bangka-belitung', label: 'Kepulauan Bangka Belitung', island: 'Sumatera' },
	{ code: 'lampung', label: 'Lampung', island: 'Sumatera' },

	// ── Jawa ──
	{ code: 'banten', label: 'Banten', island: 'Jawa' },
	{ code: 'dki-jakarta', label: 'DKI Jakarta', island: 'Jawa' },
	{ code: 'jawa-barat', label: 'Jawa Barat', island: 'Jawa' },
	{ code: 'jawa-tengah', label: 'Jawa Tengah', island: 'Jawa' },
	{ code: 'di-yogyakarta', label: 'DI Yogyakarta', island: 'Jawa' },
	{ code: 'jawa-timur', label: 'Jawa Timur', island: 'Jawa' },

	// ── Bali & Nusa Tenggara ──
	{ code: 'bali', label: 'Bali', island: 'Bali & Nusa Tenggara' },
	{ code: 'nusa-tenggara-barat', label: 'Nusa Tenggara Barat', island: 'Bali & Nusa Tenggara' },
	{ code: 'nusa-tenggara-timur', label: 'Nusa Tenggara Timur', island: 'Bali & Nusa Tenggara' },

	// ── Kalimantan ──
	{ code: 'kalimantan-barat', label: 'Kalimantan Barat', island: 'Kalimantan' },
	{ code: 'kalimantan-tengah', label: 'Kalimantan Tengah', island: 'Kalimantan' },
	{ code: 'kalimantan-selatan', label: 'Kalimantan Selatan', island: 'Kalimantan' },
	{ code: 'kalimantan-timur', label: 'Kalimantan Timur', island: 'Kalimantan' },
	{ code: 'kalimantan-utara', label: 'Kalimantan Utara', island: 'Kalimantan' },

	// ── Sulawesi ──
	{ code: 'sulawesi-utara', label: 'Sulawesi Utara', island: 'Sulawesi' },
	{ code: 'gorontalo', label: 'Gorontalo', island: 'Sulawesi' },
	{ code: 'sulawesi-tengah', label: 'Sulawesi Tengah', island: 'Sulawesi' },
	{ code: 'sulawesi-barat', label: 'Sulawesi Barat', island: 'Sulawesi' },
	{ code: 'sulawesi-selatan', label: 'Sulawesi Selatan', island: 'Sulawesi' },
	{ code: 'sulawesi-tenggara', label: 'Sulawesi Tenggara', island: 'Sulawesi' },

	// ── Maluku ──
	{ code: 'maluku', label: 'Maluku', island: 'Maluku' },
	{ code: 'maluku-utara', label: 'Maluku Utara', island: 'Maluku' },

	// ── Papua ──
	{ code: 'papua-barat', label: 'Papua Barat', island: 'Papua' },
	{ code: 'papua-barat-daya', label: 'Papua Barat Daya', island: 'Papua' },
	{ code: 'papua', label: 'Papua', island: 'Papua' },
	{ code: 'papua-selatan', label: 'Papua Selatan', island: 'Papua' },
	{ code: 'papua-tengah', label: 'Papua Tengah', island: 'Papua' },
	{ code: 'papua-pegunungan', label: 'Papua Pegunungan', island: 'Papua' },
] as const satisfies readonly ProvinceConfig[];

export type ProvinceCode = (typeof INDONESIA_PROVINCES)[number]['code'];

/**
 * Tuple non-kosong berisi kode saja, supaya bisa langsung dipakai `z.enum()` di
 * schemas/auth.ts. Diturunkan dari INDONESIA_PROVINCES — tidak ada daftar kedua
 * yang bisa jadi tidak sinkron.
 */
export const PROVINCE_CODES = INDONESIA_PROVINCES.map((province) => province.code) as unknown as [
	ProvinceCode,
	...ProvinceCode[],
];

const provinceByCode = new Map<string, ProvinceConfig>(
	INDONESIA_PROVINCES.map((province) => [province.code, province]),
);

export function isProvinceCode(value: string | null | undefined): value is ProvinceCode {
	return !!value && provinceByCode.has(value);
}

/**
 * Kode tidak dikenal (mis. data lama dari sebelum daftar ini ada) dikembalikan
 * apa adanya, bukan dilempar error — halaman profil tidak boleh gagal render
 * hanya karena satu nilai lama yang tidak lagi ada di daftar.
 */
export function getProvinceLabel(code: string | null | undefined): string {
	if (!code) return '';
	return provinceByCode.get(code)?.label ?? code;
}

/**
 * Provinsi dikelompokkan per pulau, siap dipakai untuk `<optgroup>`. 38 opsi
 * dalam satu daftar rata terlalu panjang untuk di-scan mata; pengelompokan per
 * pulau memotong ruang pencarian jadi 2-10 opsi begitu orang menemukan pulaunya.
 */
export function getProvincesByIsland(): {
	island: ProvinceIsland;
	provinces: ProvinceConfig[];
}[] {
	return PROVINCE_ISLANDS.map((island) => ({
		island,
		provinces: INDONESIA_PROVINCES.filter((province) => province.island === island),
	}));
}
