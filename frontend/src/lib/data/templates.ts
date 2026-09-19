import type { EventCategoryId } from '@ketuk/shared';

export interface PricingTier {
	id: string;
	name: string;
	price: string;
	description: string;
	popular?: boolean;
}

export interface TemplateItem {
	id: string;
	title: string;
	subtitle: string;
	tagline?: string;
	category: string;
	/**
	 * Kategori acara yang cocok untuk template ini (dari taxonomy dua level).
	 * Opsional: template tanpa field ini dianggap cocok hanya untuk filter "Semua"
	 * — bukan untuk semua kategori — supaya desain pernikahan tidak muncul saat
	 * user mencari contoh undangan ulang tahun.
	 */
	eventTypes?: EventCategoryId[];
	badge?: string;
	price: string;
	slug: string;
	rating: number;
	reviewCount: number;
	envelopeScript: string;
	palette: {
		bgGradient: string;
		accentColor: string;
		accentTextColor: string;
		sealColor: string;
		phoneBorder: string;
	};
	preview: {
		coupleName: string;
		date: string;
		style: string;
		location: string;
		monogram: string;
	};
	tiers: PricingTier[];
	features: string[];
}

export const TEMPLATES: TemplateItem[] = [
	{
		id: 'maison-doree',
		eventTypes: ['pernikahan'],
		title: 'La Maison Dorée',
		subtitle:
			'Kemewahan klasik Prancis dengan ornamen floral timbul dan stempel lilin sage green eksklusif.',
		tagline: 'Undangan yang membuat para tamu terpukau bahkan sebelum hari pernikahan tiba.',
		category: 'Klasik',
		badge: 'BESTSELLER',
		price: 'Rp 175.000',
		slug: 'maison-doree',
		rating: 4.9,
		reviewCount: 248,
		envelopeScript: 'Requests the pleasure of your company',
		palette: {
			bgGradient: 'bg-gradient-to-b from-stone-50 via-cream-50 to-amber-50/40',
			accentColor: '#5c1421',
			accentTextColor: '#7e6f65',
			sealColor: '#7a9a7a', // sage green wax seal as seen in the screenshot!
			phoneBorder: 'border-stone-800',
		},
		preview: {
			coupleName: 'Elena & Julian',
			date: 'Sabtu, 24 Oktober 2026',
			style: 'French Renaissance Couture',
			location: 'Château de Montfort / Plataran Menteng, Jakarta',
			monogram: 'E & J',
		},
		tiers: [
			{
				id: 'essential',
				name: 'Essential',
				price: 'Rp 149.000',
				description:
					'Desain digital esensial, RSVP otomatis via WhatsApp, dan petunjuk Google Maps.',
			},
			{
				id: 'premium',
				name: 'Premium',
				price: 'Rp 229.000',
				description:
					'Desain interaktif lengkap, musik romantis latar, amplop digital QRIS & galeri foto.',
				popular: true,
			},
			{
				id: 'exclusive',
				name: 'Exclusive',
				price: 'Rp 349.000',
				description:
					'Kustomisasi tanpa batas, WhatsApp blast personal untuk tamu VIP, dan domain nama pengantin.',
			},
		],
		features: [
			'Animasi pembuka amplop dengan stempel lilin interaktif',
			'RSVP instan terhubung otomatis ke WhatsApp pengantin',
			'Integrasi petunjuk arah akurat Google Maps & Waze',
			'Musik latar romantis dengan kontrol audio autoplay',
			'Galeri foto momen terindah hingga 20 foto resolusi tinggi',
			'Amplop digital resmi (BCA, Mandiri, BRI & QRIS instan)',
			'Buku tamu digital & rekap otomatis kehadiran',
			'Masa aktif selamanya tanpa biaya perpanjangan',
		],
	},
	{
		id: 'capri-rosa',
		eventTypes: ['pernikahan'],
		title: 'Capri Rosa',
		subtitle:
			'Estetika memikat terinspirasi nuansa keanggunan abadi pesisir Capri dengan sentuhan blush lembut.',
		tagline: 'Pesona romansa pesisir Mediterania dalam sentuhan modern yang memikat hati.',
		category: 'Elegan',
		badge: 'POPULER',
		price: 'Rp 149.000',
		slug: 'capri-rosa',
		rating: 4.9,
		reviewCount: 312,
		envelopeScript: 'Requests the honor of your presence',
		palette: {
			bgGradient: 'bg-gradient-to-b from-rose-50 via-cream-50 to-cream-100',
			accentColor: '#882235',
			accentTextColor: '#882235',
			sealColor: '#882235',
			phoneBorder: 'border-rose-900',
		},
		preview: {
			coupleName: 'Capri & Rosa',
			date: 'Sabtu, 14 November 2026',
			style: 'Romantic Floral',
			location: 'The Hermitage, Jakarta Pusat',
			monogram: 'C & R',
		},
		tiers: [
			{
				id: 'essential',
				name: 'Essential',
				price: 'Rp 149.000',
				description:
					'Desain digital esensial, RSVP otomatis via WhatsApp, dan petunjuk Google Maps.',
			},
			{
				id: 'premium',
				name: 'Premium',
				price: 'Rp 229.000',
				description:
					'Desain interaktif lengkap, musik romantis latar, amplop digital QRIS & galeri foto.',
				popular: true,
			},
			{
				id: 'exclusive',
				name: 'Exclusive',
				price: 'Rp 349.000',
				description:
					'Kustomisasi tanpa batas, WhatsApp blast personal untuk tamu VIP, dan domain nama pengantin.',
			},
		],
		features: [
			'Animasi pembuka amplop dengan stempel lilin bordeaux',
			'RSVP instan terhubung otomatis ke WhatsApp pengantin',
			'Integrasi petunjuk arah akurat Google Maps & Waze',
			'Musik latar romantis dengan kontrol audio autoplay',
			'Galeri foto momen terindah hingga 20 foto resolusi tinggi',
			'Amplop digital resmi (BCA, Mandiri, BRI & QRIS instan)',
			'Buku tamu digital & rekap otomatis kehadiran',
			'Masa aktif selamanya tanpa biaya perpanjangan',
		],
	},
	{
		id: 'firenze',
		eventTypes: ['pernikahan'],
		title: 'Firenze',
		subtitle:
			'Kemegahan arsitektur Renaissance Italia dengan tipografi klasik dan ornamen bingkai halus.',
		tagline:
			'Keanggunan masa lalu yang dihadirkan kembali dengan teknologi undangan digital masa kini.',
		category: 'Klasik',
		badge: 'BESTSELLER',
		price: 'Rp 169.000',
		slug: 'firenze',
		rating: 5.0,
		reviewCount: 184,
		envelopeScript: 'Together with their families',
		palette: {
			bgGradient: 'bg-gradient-to-b from-stone-50 via-cream-50 to-cream-100',
			accentColor: '#460d17',
			accentTextColor: '#8c6a3e',
			sealColor: '#a98350',
			phoneBorder: 'border-stone-800',
		},
		preview: {
			coupleName: 'Fiona & Lorenzo',
			date: 'Minggu, 20 Desember 2026',
			style: 'Classic Renaissance',
			location: 'Hotel Mulia Senayan, Jakarta',
			monogram: 'F & L',
		},
		tiers: [
			{
				id: 'essential',
				name: 'Essential',
				price: 'Rp 149.000',
				description:
					'Desain digital esensial, RSVP otomatis via WhatsApp, dan petunjuk Google Maps.',
			},
			{
				id: 'premium',
				name: 'Premium',
				price: 'Rp 229.000',
				description:
					'Desain interaktif lengkap, musik romantis latar, amplop digital QRIS & galeri foto.',
				popular: true,
			},
			{
				id: 'exclusive',
				name: 'Exclusive',
				price: 'Rp 349.000',
				description:
					'Kustomisasi tanpa batas, WhatsApp blast personal untuk tamu VIP, dan domain nama pengantin.',
			},
		],
		features: [
			'Animasi pembuka amplop dengan stempel lilin emas murni',
			'RSVP instan terhubung otomatis ke WhatsApp pengantin',
			'Integrasi petunjuk arah akurat Google Maps & Waze',
			'Musik latar romantis dengan kontrol audio autoplay',
			'Galeri foto momen terindah hingga 20 foto resolusi tinggi',
			'Amplop digital resmi (BCA, Mandiri, BRI & QRIS instan)',
			'Buku tamu digital & rekap otomatis kehadiran',
			'Masa aktif selamanya tanpa biaya perpanjangan',
		],
	},
	{
		id: 'tuscan-garden',
		eventTypes: ['pernikahan'],
		title: 'Tuscan Garden',
		subtitle:
			'Nuansa alam perbukitan Tuscany dengan ilustrasi olive botanical yang tenang dan teduh.',
		tagline:
			'Ketenangan dedaunan zaitun untuk pernikahan outdoor yang penuh kedamaian dan kehangatan.',
		category: 'Botanical',
		badge: 'BARU',
		price: 'Rp 149.000',
		slug: 'tuscan-garden',
		rating: 4.8,
		reviewCount: 156,
		envelopeScript: 'Cordially invite you to celebrate',
		palette: {
			bgGradient: 'bg-gradient-to-b from-emerald-50/70 via-cream-50 to-cream-100',
			accentColor: '#166534',
			accentTextColor: '#15803d',
			sealColor: '#2e5c3e',
			phoneBorder: 'border-emerald-950',
		},
		preview: {
			coupleName: 'Tara & Julian',
			date: 'Sabtu, 05 September 2026',
			style: 'Botanical Olive',
			location: 'Pine Hill Cibodas, Bandung',
			monogram: 'T & J',
		},
		tiers: [
			{
				id: 'essential',
				name: 'Essential',
				price: 'Rp 149.000',
				description:
					'Desain digital esensial, RSVP otomatis via WhatsApp, dan petunjuk Google Maps.',
			},
			{
				id: 'premium',
				name: 'Premium',
				price: 'Rp 229.000',
				description:
					'Desain interaktif lengkap, musik romantis latar, amplop digital QRIS & galeri foto.',
				popular: true,
			},
			{
				id: 'exclusive',
				name: 'Exclusive',
				price: 'Rp 349.000',
				description:
					'Kustomisasi tanpa batas, WhatsApp blast personal untuk tamu VIP, dan domain nama pengantin.',
			},
		],
		features: [
			'Animasi pembuka amplop dengan stempel lilin olive botanical',
			'RSVP instan terhubung otomatis ke WhatsApp pengantin',
			'Integrasi petunjuk arah akurat Google Maps & Waze',
			'Musik latar romantis dengan kontrol audio autoplay',
			'Galeri foto momen terindah hingga 20 foto resolusi tinggi',
			'Amplop digital resmi (BCA, Mandiri, BRI & QRIS instan)',
			'Buku tamu digital & rekap otomatis kehadiran',
			'Masa aktif selamanya tanpa biaya perpanjangan',
		],
	},
	{
		id: 'veneto-romantis',
		eventTypes: ['pernikahan'],
		title: 'Veneto Romantis',
		subtitle: 'Kemewahan warna bordeaux wine dengan sentuhan foil emas untuk pesta malam berkelas.',
		tagline: 'Dramatis, elegan, dan memesona untuk pesta resepsi malam yang tak terlupakan.',
		category: 'Elegan',
		badge: 'EKSKLUSIF',
		price: 'Rp 179.000',
		slug: 'veneto-romantis',
		rating: 4.9,
		reviewCount: 98,
		envelopeScript: 'Requests the honor of your company',
		palette: {
			bgGradient: 'bg-gradient-to-b from-coffee-50 via-cream-50 to-coffee-100/50',
			accentColor: '#5c1421',
			accentTextColor: '#c5a478',
			sealColor: '#6f1929',
			phoneBorder: 'border-coffee-950',
		},
		preview: {
			coupleName: 'Valerie & Nathan',
			date: 'Sabtu, 12 Desember 2026',
			style: 'Haute Bordeaux',
			location: 'The Dharmawangsa, Jakarta Selatan',
			monogram: 'V & N',
		},
		tiers: [
			{
				id: 'essential',
				name: 'Essential',
				price: 'Rp 149.000',
				description:
					'Desain digital esensial, RSVP otomatis via WhatsApp, dan petunjuk Google Maps.',
			},
			{
				id: 'premium',
				name: 'Premium',
				price: 'Rp 229.000',
				description:
					'Desain interaktif lengkap, musik romantis latar, amplop digital QRIS & galeri foto.',
				popular: true,
			},
			{
				id: 'exclusive',
				name: 'Exclusive',
				price: 'Rp 349.000',
				description:
					'Kustomisasi tanpa batas, WhatsApp blast personal untuk tamu VIP, dan domain nama pengantin.',
			},
		],
		features: [
			'Animasi pembuka amplop dengan stempel lilin bordeaux wine',
			'RSVP instan terhubung otomatis ke WhatsApp pengantin',
			'Integrasi petunjuk arah akurat Google Maps & Waze',
			'Musik latar romantis dengan kontrol audio autoplay',
			'Galeri foto momen terindah hingga 20 foto resolusi tinggi',
			'Amplop digital resmi (BCA, Mandiri, BRI & QRIS instan)',
			'Buku tamu digital & rekap otomatis kehadiran',
			'Masa aktif selamanya tanpa biaya perpanjangan',
		],
	},
	{
		id: 'amalfi-sun',
		eventTypes: ['pernikahan'],
		title: 'Amalfi Sun',
		subtitle:
			'Kehangatan hangat mediterania dengan aksen lemon blossom dan latar kertas bertekstur.',
		tagline: 'Keceriaan sinar matahari Mediterania untuk cinta yang berseri dan penuh tawa.',
		category: 'Modern',
		badge: 'BARU',
		price: 'Rp 149.000',
		slug: 'amalfi-sun',
		rating: 4.8,
		reviewCount: 76,
		envelopeScript: 'Invite you to celebrate their love',
		palette: {
			bgGradient: 'bg-gradient-to-b from-amber-50 via-cream-50 to-orange-50/50',
			accentColor: '#a98350',
			accentTextColor: '#c2410c',
			sealColor: '#c5a478',
			phoneBorder: 'border-amber-950',
		},
		preview: {
			coupleName: 'Alya & Marcel',
			date: 'Minggu, 18 Oktober 2026',
			style: 'Warm Terracotta',
			location: 'Ayana Resort, Jimbaran Bali',
			monogram: 'A & M',
		},
		tiers: [
			{
				id: 'essential',
				name: 'Essential',
				price: 'Rp 149.000',
				description:
					'Desain digital esensial, RSVP otomatis via WhatsApp, dan petunjuk Google Maps.',
			},
			{
				id: 'premium',
				name: 'Premium',
				price: 'Rp 229.000',
				description:
					'Desain interaktif lengkap, musik romantis latar, amplop digital QRIS & galeri foto.',
				popular: true,
			},
			{
				id: 'exclusive',
				name: 'Exclusive',
				price: 'Rp 349.000',
				description:
					'Kustomisasi tanpa batas, WhatsApp blast personal untuk tamu VIP, dan domain nama pengantin.',
			},
		],
		features: [
			'Animasi pembuka amplop dengan stempel lilin terracotta emas',
			'RSVP instan terhubung otomatis ke WhatsApp pengantin',
			'Integrasi petunjuk arah akurat Google Maps & Waze',
			'Musik latar romantis dengan kontrol audio autoplay',
			'Galeri foto momen terindah hingga 20 foto resolusi tinggi',
			'Amplop digital resmi (BCA, Mandiri, BRI & QRIS instan)',
			'Buku tamu digital & rekap otomatis kehadiran',
			'Masa aktif selamanya tanpa biaya perpanjangan',
		],
	},
	{
		id: 'monaco-royal',
		eventTypes: ['pernikahan'],
		title: 'Monaco Royal',
		subtitle: 'Simbol kemewahan absolut dengan aksen dark emerald dan sentuhan emas berkilau.',
		tagline: 'Kemewahan aristokrat sejati untuk momen perayaan cinta paling prestisius.',
		category: 'Luxury Gold',
		badge: 'EKSKLUSIF',
		price: 'Rp 199.000',
		slug: 'monaco-royal',
		rating: 5.0,
		reviewCount: 64,
		envelopeScript: 'The honor of your presence is requested',
		palette: {
			bgGradient: 'bg-gradient-to-b from-zinc-50 via-cream-50 to-emerald-50/30',
			accentColor: '#064e3b',
			accentTextColor: '#8c6a3e',
			sealColor: '#a98350',
			phoneBorder: 'border-zinc-900',
		},
		preview: {
			coupleName: 'Melody & Richard',
			date: 'Sabtu, 28 November 2026',
			style: 'Royal Emerald',
			location: 'The Ritz-Carlton Pacific Place, Jakarta',
			monogram: 'M & R',
		},
		tiers: [
			{
				id: 'essential',
				name: 'Essential',
				price: 'Rp 149.000',
				description:
					'Desain digital esensial, RSVP otomatis via WhatsApp, dan petunjuk Google Maps.',
			},
			{
				id: 'premium',
				name: 'Premium',
				price: 'Rp 229.000',
				description:
					'Desain interaktif lengkap, musik romantis latar, amplop digital QRIS & galeri foto.',
				popular: true,
			},
			{
				id: 'exclusive',
				name: 'Exclusive',
				price: 'Rp 349.000',
				description:
					'Kustomisasi tanpa batas, WhatsApp blast personal untuk tamu VIP, dan domain nama pengantin.',
			},
		],
		features: [
			'Animasi pembuka amplop dengan stempel lilin royal emerald',
			'RSVP instan terhubung otomatis ke WhatsApp pengantin',
			'Integrasi petunjuk arah akurat Google Maps & Waze',
			'Musik latar romantis dengan kontrol audio autoplay',
			'Galeri foto momen terindah hingga 20 foto resolusi tinggi',
			'Amplop digital resmi (BCA, Mandiri, BRI & QRIS instan)',
			'Buku tamu digital & rekap otomatis kehadiran',
			'Masa aktif selamanya tanpa biaya perpanjangan',
		],
	},
	{
		id: 'santorini-white',
		eventTypes: ['pernikahan'],
		title: 'Santorini White',
		subtitle: 'Kesederhanaan minimalis modern dengan ruang bernapas lapang dan tipografi monokrom.',
		tagline: 'Keindahan dalam kesederhanaan, bersih, rapi, dan memukau dalam segala sudut.',
		category: 'Minimalis',
		badge: 'POPULER',
		price: 'Rp 139.000',
		slug: 'santorini-white',
		rating: 4.9,
		reviewCount: 142,
		envelopeScript: 'Joyfully invite you to their wedding',
		palette: {
			bgGradient: 'bg-gradient-to-b from-sky-50/50 via-white to-cream-50',
			accentColor: '#1e293b',
			accentTextColor: '#0284c7',
			sealColor: '#334155',
			phoneBorder: 'border-slate-800',
		},
		preview: {
			coupleName: 'Stella & Andre',
			date: 'Minggu, 04 Oktober 2026',
			style: 'Pure Minimalist',
			location: 'Plataran Canggu, Bali',
			monogram: 'S & A',
		},
		tiers: [
			{
				id: 'essential',
				name: 'Essential',
				price: 'Rp 149.000',
				description:
					'Desain digital esensial, RSVP otomatis via WhatsApp, dan petunjuk Google Maps.',
			},
			{
				id: 'premium',
				name: 'Premium',
				price: 'Rp 229.000',
				description:
					'Desain interaktif lengkap, musik romantis latar, amplop digital QRIS & galeri foto.',
				popular: true,
			},
			{
				id: 'exclusive',
				name: 'Exclusive',
				price: 'Rp 349.000',
				description:
					'Kustomisasi tanpa batas, WhatsApp blast personal untuk tamu VIP, dan domain nama pengantin.',
			},
		],
		features: [
			'Animasi pembuka amplop dengan stempel lilin minimalis slate',
			'RSVP instan terhubung otomatis ke WhatsApp pengantin',
			'Integrasi petunjuk arah akurat Google Maps & Waze',
			'Musik latar romantis dengan kontrol audio autoplay',
			'Galeri foto momen terindah hingga 20 foto resolusi tinggi',
			'Amplop digital resmi (BCA, Mandiri, BRI & QRIS instan)',
			'Buku tamu digital & rekap otomatis kehadiran',
			'Masa aktif selamanya tanpa biaya perpanjangan',
		],
	},
	{
		id: 'kyoto-blossom',
		eventTypes: ['pernikahan'],
		title: 'Kyoto Blossom',
		subtitle:
			'Harmoni ketenangan zen dengan dedaunan sakura lembut dan tata letak elegan kontemporer.',
		tagline: 'Kedamaian estetika Jepang dengan sentuhan modern kontemporer yang abadi.',
		category: 'Botanical',
		badge: '',
		price: 'Rp 149.000',
		slug: 'kyoto-blossom',
		rating: 4.8,
		reviewCount: 95,
		envelopeScript: 'Warmly invite you to share in their joy',
		palette: {
			bgGradient: 'bg-gradient-to-b from-pink-50/70 via-cream-50 to-cream-100',
			accentColor: '#9d174d',
			accentTextColor: '#be185d',
			sealColor: '#882235',
			phoneBorder: 'border-pink-950',
		},
		preview: {
			coupleName: 'Kezia & Bryan',
			date: 'Sabtu, 19 September 2026',
			style: 'Zen Floral',
			location: 'Alila Seminyak, Bali',
			monogram: 'K & B',
		},
		tiers: [
			{
				id: 'essential',
				name: 'Essential',
				price: 'Rp 149.000',
				description:
					'Desain digital esensial, RSVP otomatis via WhatsApp, dan petunjuk Google Maps.',
			},
			{
				id: 'premium',
				name: 'Premium',
				price: 'Rp 229.000',
				description:
					'Desain interaktif lengkap, musik romantis latar, amplop digital QRIS & galeri foto.',
				popular: true,
			},
			{
				id: 'exclusive',
				name: 'Exclusive',
				price: 'Rp 349.000',
				description:
					'Kustomisasi tanpa batas, WhatsApp blast personal untuk tamu VIP, dan domain nama pengantin.',
			},
		],
		features: [
			'Animasi pembuka amplop dengan stempel lilin sakura',
			'RSVP instan terhubung otomatis ke WhatsApp pengantin',
			'Integrasi petunjuk arah akurat Google Maps & Waze',
			'Musik latar romantis dengan kontrol audio autoplay',
			'Galeri foto momen terindah hingga 20 foto resolusi tinggi',
			'Amplop digital resmi (BCA, Mandiri, BRI & QRIS instan)',
			'Buku tamu digital & rekap otomatis kehadiran',
			'Masa aktif selamanya tanpa biaya perpanjangan',
		],
	},
];
