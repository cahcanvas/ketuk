export type ProductCategory =
	| 'Klasik'
	| 'Elegan'
	| 'Botanical'
	| 'Luxury Gold'
	| 'Minimalis'
	| 'Modern';

export type ProductBadge = 'BESTSELLER' | 'POPULER' | 'BARU' | 'EKSKLUSIF' | '';
export type ProductStatus = 'published' | 'draft' | 'archived';

export type OpeningEffectType =
	| 'wax_seal_crack'
	| 'velvet_curtain'
	| 'blooming_rose'
	| 'sliding_luxe'
	| 'gate_fold';

export type ParticleEffectType =
	| 'falling_petals'
	| 'gold_stardust'
	| 'confetti'
	| 'floating_lanterns'
	| 'none';

export type ParticleDensity = 'low' | 'medium' | 'high';
export type ScrollMotionType = 'parallax' | 'fade_up' | 'blur_zoom';

export interface PlanPermissions {
	guestLimit: number; // -1 = Unlimited
	galleryPhotoLimit: number; // e.g. 5, 20, 50
	allowVideoStreaming: boolean;
	allowLoveStory: boolean;
	allowCustomMusic: boolean; // Upload MP3 sendiri vs hanya pilihan preset
	allowDigitalEnvelope: boolean; // QRIS / Rekening Bank
	allowCustomDomain: boolean; // nama-pasangan.com
	allowCustomPaletteOverride: boolean; // Pasangan boleh modifikasi warna aksen/stempel
	allowCustomFont: boolean; // Pasangan boleh ganti font
	removeKetukWatermark: boolean; // Hilangkan branding "Powered by Ketuk.id"
	allowWhatsAppBlast: boolean; // Generator pesan blast WhatsApp
	allowLiveGuestbook: boolean; // Buku tamu ucapan interaktif
	allowCountdownTimer: boolean; // Widget hitung mundur & Add to Calendar
}

export interface MasterTemplateProduct {
	id: string;
	title: string;
	slug: string;
	subtitle: string;
	tagline: string;
	category: ProductCategory;
	badge?: ProductBadge;
	status: ProductStatus;
	basePrice: number;
	rating: number;
	reviewCount: number;

	// 1. Visual & Styling Engine
	palette: {
		bgGradient: string;
		accentColor: string;
		accentTextColor: string;
		sealColor: string;
		phoneBorder: string;
		fontFamily: string;
	};

	preview: {
		coupleName: string;
		date: string;
		style: string;
		location: string;
		monogram: string;
	};

	envelope: {
		envelopeScript: string;
		paperTexture: 'classic_embossed' | 'handmade_deckle' | 'smooth_linen' | 'silk_matte';
		sealSymbol: 'heart' | 'monogram' | 'botanical' | 'rings';
	};

	// 2. Engine Animasi & Interaktivitas
	animations: {
		openingEffect: OpeningEffectType;
		particleEffect: ParticleEffectType;
		particleDensity: ParticleDensity;
		scrollMotion: ScrollMotionType;
		openingDurationMs: number;
	};

	// 3. Audio & Soundscape
	audio: {
		presetTrackName: string;
		presetTrackUrl: string;
		autoplay: boolean;
		showVisualizer: boolean;
		fadeDurationMs: number;
	};

	// 4. Matriks Kustomisasi Member per Paket
	customizationMatrix: {
		gratis: PlanPermissions;
		pro: PlanPermissions;
		lengkap: PlanPermissions;
	};
}

export const DEFAULT_PLAN_PERMISSIONS: {
	gratis: PlanPermissions;
	pro: PlanPermissions;
	lengkap: PlanPermissions;
} = {
	gratis: {
		guestLimit: 100,
		galleryPhotoLimit: 4,
		allowVideoStreaming: false,
		allowLoveStory: false,
		allowCustomMusic: false,
		allowDigitalEnvelope: false,
		allowCustomDomain: false,
		allowCustomPaletteOverride: false,
		allowCustomFont: false,
		removeKetukWatermark: false,
		allowWhatsAppBlast: false,
		allowLiveGuestbook: true,
		allowCountdownTimer: true,
	},
	pro: {
		guestLimit: -1, // Unlimited
		galleryPhotoLimit: 20,
		allowVideoStreaming: false,
		allowLoveStory: true,
		allowCustomMusic: true,
		allowDigitalEnvelope: true,
		allowCustomDomain: false,
		allowCustomPaletteOverride: true,
		allowCustomFont: false,
		removeKetukWatermark: false,
		allowWhatsAppBlast: true,
		allowLiveGuestbook: true,
		allowCountdownTimer: true,
	},
	lengkap: {
		guestLimit: -1, // Unlimited
		galleryPhotoLimit: 60,
		allowVideoStreaming: true,
		allowLoveStory: true,
		allowCustomMusic: true,
		allowDigitalEnvelope: true,
		allowCustomDomain: true,
		allowCustomPaletteOverride: true,
		allowCustomFont: true,
		removeKetukWatermark: true,
		allowWhatsAppBlast: true,
		allowLiveGuestbook: true,
		allowCountdownTimer: true,
	},
};

export const MASTER_TEMPLATES: MasterTemplateProduct[] = [
	{
		id: 'maison-doree',
		title: 'La Maison Dorée',
		slug: 'maison-doree',
		subtitle:
			'Kemewahan klasik Prancis dengan ornamen floral timbul dan stempel lilin sage green eksklusif.',
		tagline: 'Undangan yang membuat para tamu terpukau bahkan sebelum hari pernikahan tiba.',
		category: 'Klasik',
		badge: 'BESTSELLER',
		status: 'published',
		basePrice: 175000,
		rating: 4.9,
		reviewCount: 248,
		palette: {
			bgGradient: 'bg-gradient-to-b from-stone-50 via-cream-50 to-amber-50/40',
			accentColor: '#5c1421',
			accentTextColor: '#7e6f65',
			sealColor: '#7a9a7a',
			phoneBorder: 'border-stone-800',
			fontFamily: 'Cormorant Garamond',
		},
		preview: {
			coupleName: 'Elena & Julian',
			date: 'Sabtu, 24 Oktober 2026',
			style: 'French Renaissance Couture',
			location: 'Château de Montfort / Plataran Menteng, Jakarta',
			monogram: 'E & J',
		},
		envelope: {
			envelopeScript: 'Requests the pleasure of your company',
			paperTexture: 'classic_embossed',
			sealSymbol: 'heart',
		},
		animations: {
			openingEffect: 'wax_seal_crack',
			particleEffect: 'falling_petals',
			particleDensity: 'medium',
			scrollMotion: 'parallax',
			openingDurationMs: 800,
		},
		audio: {
			presetTrackName: 'A Thousand Years (Acoustic Piano)',
			presetTrackUrl: '/audio/presets/thousand-years-piano.mp3',
			autoplay: true,
			showVisualizer: true,
			fadeDurationMs: 1500,
		},
		customizationMatrix: { ...DEFAULT_PLAN_PERMISSIONS },
	},
	{
		id: 'capri-rosa',
		title: 'Capri Rosa',
		slug: 'capri-rosa',
		subtitle:
			'Estetika memikat terinspirasi nuansa keanggunan abadi pesisir Capri dengan sentuhan blush lembut.',
		tagline: 'Pesona romansa pesisir Mediterania dalam sentuhan modern yang memikat hati.',
		category: 'Elegan',
		badge: 'POPULER',
		status: 'published',
		basePrice: 149000,
		rating: 4.9,
		reviewCount: 312,
		palette: {
			bgGradient: 'bg-gradient-to-b from-rose-50 via-cream-50 to-cream-100',
			accentColor: '#882235',
			accentTextColor: '#882235',
			sealColor: '#882235',
			phoneBorder: 'border-rose-900',
			fontFamily: 'Playfair Display',
		},
		preview: {
			coupleName: 'Capri & Rosa',
			date: 'Sabtu, 14 November 2026',
			style: 'Romantic Floral',
			location: 'The Hermitage, Jakarta Pusat',
			monogram: 'C & R',
		},
		envelope: {
			envelopeScript: 'Requests the honor of your presence',
			paperTexture: 'handmade_deckle',
			sealSymbol: 'botanical',
		},
		animations: {
			openingEffect: 'blooming_rose',
			particleEffect: 'falling_petals',
			particleDensity: 'high',
			scrollMotion: 'fade_up',
			openingDurationMs: 900,
		},
		audio: {
			presetTrackName: 'Canon in D (Strings & Harp)',
			presetTrackUrl: '/audio/presets/canon-in-d.mp3',
			autoplay: true,
			showVisualizer: true,
			fadeDurationMs: 1200,
		},
		customizationMatrix: { ...DEFAULT_PLAN_PERMISSIONS },
	},
	{
		id: 'firenze',
		title: 'Firenze',
		slug: 'firenze',
		subtitle:
			'Kemegahan arsitektur Renaissance Italia dengan tipografi klasik dan ornamen bingkai halus emas.',
		tagline:
			'Keanggunan masa lalu yang dihadirkan kembali dengan teknologi undangan digital masa kini.',
		category: 'Klasik',
		badge: 'BESTSELLER',
		status: 'published',
		basePrice: 169000,
		rating: 5.0,
		reviewCount: 184,
		palette: {
			bgGradient: 'bg-gradient-to-b from-stone-50 via-cream-50 to-cream-100',
			accentColor: '#460d17',
			accentTextColor: '#8c6a3e',
			sealColor: '#a98350',
			phoneBorder: 'border-stone-800',
			fontFamily: 'Cinzel',
		},
		preview: {
			coupleName: 'Fiona & Lorenzo',
			date: 'Minggu, 20 Desember 2026',
			style: 'Classic Renaissance',
			location: 'Hotel Mulia Senayan, Jakarta',
			monogram: 'F & L',
		},
		envelope: {
			envelopeScript: 'Together with their families',
			paperTexture: 'smooth_linen',
			sealSymbol: 'monogram',
		},
		animations: {
			openingEffect: 'wax_seal_crack',
			particleEffect: 'gold_stardust',
			particleDensity: 'medium',
			scrollMotion: 'parallax',
			openingDurationMs: 850,
		},
		audio: {
			presetTrackName: 'Renaissance Romance (Orchestral)',
			presetTrackUrl: '/audio/presets/renaissance-waltz.mp3',
			autoplay: true,
			showVisualizer: false,
			fadeDurationMs: 2000,
		},
		customizationMatrix: { ...DEFAULT_PLAN_PERMISSIONS },
	},
	{
		id: 'tuscan-garden',
		title: 'Tuscan Garden',
		slug: 'tuscan-garden',
		subtitle:
			'Nuansa alam perbukitan Tuscany dengan ilustrasi olive botanical yang tenang dan teduh.',
		tagline:
			'Ketenangan dedaunan zaitun untuk pernikahan outdoor yang penuh kedamaian dan kehangatan.',
		category: 'Botanical',
		badge: 'BARU',
		status: 'published',
		basePrice: 149000,
		rating: 4.8,
		reviewCount: 156,
		palette: {
			bgGradient: 'bg-gradient-to-b from-emerald-50/70 via-cream-50 to-cream-100',
			accentColor: '#166534',
			accentTextColor: '#15803d',
			sealColor: '#2e5c3e',
			phoneBorder: 'border-emerald-950',
			fontFamily: 'Cormorant Garamond',
		},
		preview: {
			coupleName: 'Tara & Julian',
			date: 'Sabtu, 05 September 2026',
			style: 'Botanical Olive',
			location: 'Pine Hill Cibodas, Bandung',
			monogram: 'T & J',
		},
		envelope: {
			envelopeScript: 'Cordially invite you to celebrate',
			paperTexture: 'handmade_deckle',
			sealSymbol: 'botanical',
		},
		animations: {
			openingEffect: 'gate_fold',
			particleEffect: 'falling_petals',
			particleDensity: 'low',
			scrollMotion: 'fade_up',
			openingDurationMs: 800,
		},
		audio: {
			presetTrackName: 'Acoustic Morning Song (Folk Guitar)',
			presetTrackUrl: '/audio/presets/acoustic-garden.mp3',
			autoplay: true,
			showVisualizer: true,
			fadeDurationMs: 1500,
		},
		customizationMatrix: { ...DEFAULT_PLAN_PERMISSIONS },
	},
	{
		id: 'monaco-royal',
		title: 'Monaco Royal',
		slug: 'monaco-royal',
		subtitle:
			'Kemewahan aristokrat sejati dengan aksen dark emerald dan sentuhan emas berkilau memukau.',
		tagline: 'Simbol kemewahan absolut untuk momen perayaan cinta paling prestisius.',
		category: 'Luxury Gold',
		badge: 'EKSKLUSIF',
		status: 'published',
		basePrice: 199000,
		rating: 5.0,
		reviewCount: 64,
		palette: {
			bgGradient: 'bg-gradient-to-b from-zinc-50 via-cream-50 to-emerald-50/30',
			accentColor: '#064e3b',
			accentTextColor: '#8c6a3e',
			sealColor: '#a98350',
			phoneBorder: 'border-zinc-900',
			fontFamily: 'Cinzel',
		},
		preview: {
			coupleName: 'Melody & Richard',
			date: 'Sabtu, 28 November 2026',
			style: 'Royal Emerald & Gold',
			location: 'The Ritz-Carlton Pacific Place, Jakarta',
			monogram: 'M & R',
		},
		envelope: {
			envelopeScript: 'The honor of your presence is requested',
			paperTexture: 'silk_matte',
			sealSymbol: 'monogram',
		},
		animations: {
			openingEffect: 'velvet_curtain',
			particleEffect: 'gold_stardust',
			particleDensity: 'high',
			scrollMotion: 'parallax',
			openingDurationMs: 1000,
		},
		audio: {
			presetTrackName: 'Royal Overture (Symphonic)',
			presetTrackUrl: '/audio/presets/royal-symphony.mp3',
			autoplay: true,
			showVisualizer: true,
			fadeDurationMs: 2000,
		},
		customizationMatrix: {
			gratis: { ...DEFAULT_PLAN_PERMISSIONS.gratis },
			pro: { ...DEFAULT_PLAN_PERMISSIONS.pro },
			lengkap: { ...DEFAULT_PLAN_PERMISSIONS.lengkap },
		},
	},
	{
		id: 'santorini-white',
		title: 'Santorini White',
		slug: 'santorini-white',
		subtitle:
			'Kesederhanaan minimalis modern dengan ruang bernapas lapang dan tipografi monokrom bersih.',
		tagline: 'Keindahan dalam kesederhanaan, bersih, rapi, dan memukau dalam segala sudut.',
		category: 'Minimalis',
		badge: 'POPULER',
		status: 'published',
		basePrice: 139000,
		rating: 4.9,
		reviewCount: 142,
		palette: {
			bgGradient: 'bg-gradient-to-b from-sky-50/50 via-white to-cream-50',
			accentColor: '#1e293b',
			accentTextColor: '#0284c7',
			sealColor: '#334155',
			phoneBorder: 'border-slate-800',
			fontFamily: 'Montserrat',
		},
		preview: {
			coupleName: 'Stella & Andre',
			date: 'Minggu, 04 Oktober 2026',
			style: 'Pure Minimalist',
			location: 'Plataran Canggu, Bali',
			monogram: 'S & A',
		},
		envelope: {
			envelopeScript: 'Joyfully invite you to their wedding',
			paperTexture: 'smooth_linen',
			sealSymbol: 'rings',
		},
		animations: {
			openingEffect: 'sliding_luxe',
			particleEffect: 'none',
			particleDensity: 'low',
			scrollMotion: 'blur_zoom',
			openingDurationMs: 650,
		},
		audio: {
			presetTrackName: 'Serenade of Breeze (Lo-Fi Piano)',
			presetTrackUrl: '/audio/presets/lofi-piano.mp3',
			autoplay: false,
			showVisualizer: true,
			fadeDurationMs: 1000,
		},
		customizationMatrix: { ...DEFAULT_PLAN_PERMISSIONS },
	},
];
