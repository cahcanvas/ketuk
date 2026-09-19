import type { EventCategory, EventCategoryId, EventSubType } from '../types/onboarding';

/**
 * Taxonomy acara dua level: kategori → sub-kategori.
 *
 * Kenapa dua level dan bukan daftar datar 10-20 item: daftar datar tidak
 * scalable. Penambahan jenis undangan baru (mis. "Tasyakuran 7 bulanan"
 * besok) cukup menambah satu sub-kategori di bawah kategori yang sudah ada,
 * tanpa menyentuh UI onboarding, bentuk storage, maupun API.
 *
 * Setiap sub-kategori memetakan ke salah satu `EventType` backend yang sudah
 * ada. Sub-kategori tanpa pasangan backend dipetakan ke 'other'; seluk-beluk
 * spesifiknya hidup di `eventSubType`. Konsekuensinya endpoint pembuatan event
 * tidak berubah sama sekali di fase ini.
 */
export const EVENT_CATEGORIES: EventCategory[] = [
	{
		id: 'pernikahan',
		label: 'Pernikahan',
		emoji: '💍',
		shortDesc: 'Akad, resepsi, lamaran, hingga anniversary',
		subTypes: [
			{ id: 'wedding', label: 'Pernikahan', eventType: 'wedding' },
			{ id: 'engagement', label: 'Lamaran / Pertunangan', eventType: 'engagement' },
			{ id: 'anniversary', label: 'Anniversary', eventType: 'other' },
		],
	},
	{
		id: 'ulang-tahun',
		label: 'Ulang Tahun',
		emoji: '🎂',
		shortDesc: 'Ulang tahun anak, remaja, hingga dewasa',
		subTypes: [
			{ id: 'birthday-kid', label: 'Ulang Tahun Anak', eventType: 'birthday' },
			{ id: 'birthday-adult', label: 'Ulang Tahun Dewasa', eventType: 'birthday' },
		],
	},
	{
		id: 'keluarga',
		label: 'Acara Keluarga',
		emoji: '🏠',
		shortDesc: 'Aqiqah, khitanan, syukuran, wisuda, arisan',
		subTypes: [
			{ id: 'aqiqah', label: 'Aqiqah', eventType: 'aqiqah' },
			{ id: 'khitanan', label: 'Khitanan', eventType: 'khitanan' },
			{ id: 'syukuran', label: 'Syukuran', eventType: 'syukuran' },
			{ id: 'wisuda', label: 'Wisuda', eventType: 'graduation' },
			{ id: 'arisan', label: 'Arisan', eventType: 'other' },
		],
	},
	{
		id: 'keagamaan',
		label: 'Keagamaan',
		emoji: '🕌',
		shortDesc: 'Pengajian, kajian, marawis, peringatan',
		subTypes: [
			{ id: 'pengajian', label: 'Pengajian / Kajian', eventType: 'other' },
			{ id: 'marawis', label: 'Marawis / Rebana', eventType: 'other' },
			{ id: 'haul', label: 'Haul / Peringatan', eventType: 'other' },
		],
	},
	{
		id: 'komunitas',
		label: 'Komunitas & Reuni',
		emoji: '🤝',
		shortDesc: 'Reuni angkatan, gathering, pertemuan komunitas',
		subTypes: [
			{ id: 'reuni', label: 'Reuni', eventType: 'reunion' },
			{ id: 'gathering', label: 'Gathering', eventType: 'other' },
			{ id: 'pertemuan-komunitas', label: 'Pertemuan Komunitas', eventType: 'other' },
		],
	},
	{
		id: 'bisnis',
		label: 'Bisnis & Profesional',
		emoji: '🏢',
		shortDesc: 'Corporate event, seminar, workshop, grand opening',
		subTypes: [
			{ id: 'corporate', label: 'Corporate Event', eventType: 'corporate' },
			{ id: 'seminar', label: 'Seminar', eventType: 'corporate' },
			{ id: 'workshop', label: 'Workshop', eventType: 'corporate' },
			{ id: 'grand-opening', label: 'Grand Opening', eventType: 'corporate' },
		],
	},
	{
		id: 'other',
		label: 'Lainnya',
		emoji: '✨',
		shortDesc: 'Jenis acara lain yang kamu tentukan sendiri',
		subTypes: [],
		customLabel: true,
	},
];

export function findCategory(id: string | undefined | null): EventCategory | undefined {
	if (!id) return undefined;
	return EVENT_CATEGORIES.find((category) => category.id === id);
}

/**
 * Cari sub-kategori berdasarkan id-nya. Mengembalikan kategori dan sub-kategori
 * sekaligus karena UI perlu keduanya (label kategori untuk copy, sub-kategori
 * untuk `eventType` backend).
 */
export function findSubType(
	subTypeId: string | undefined | null,
): { category: EventCategory; subType: EventSubType } | undefined {
	if (!subTypeId) return undefined;
	for (const category of EVENT_CATEGORIES) {
		const subType = category.subTypes.find((item) => item.id === subTypeId);
		if (subType) return { category, subType };
	}
	return undefined;
}

export function isEventCategory(id: string | undefined | null): id is EventCategoryId {
	return !!findCategory(id);
}
