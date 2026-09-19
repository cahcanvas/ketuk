<script lang="ts">
	import {
		PhoneMockup,
		TemplateCard,
		TemplatePreviewModal,
		type TemplateItem,
	} from '$lib/components/domain';
	import {
		Check,
		ArrowRight,
		Star,
		MessageCircle,
		MapPin,
		Music,
		Gift,
		Clock,
		Shield,
		ChevronDown,
	} from '@lucide/svelte';

	// State for event type tabs (Pernikahan vs Acara Lain)
	let eventTab = $state<'wedding' | 'other'>('wedding');

	// State for category filters
	let activeCategory = $state('Semua');

	// State for template preview modal
	let selectedTemplate = $state<TemplateItem | null>(null);
	let isPreviewOpen = $state(false);

	// FAQ Accordion State
	let openFaq = $state<number | null>(0);

	const filterCategories = [
		'Semua',
		'Klasik',
		'Minimalis',
		'Elegan',
		'Modern',
		'Botanical',
		'Luxury Gold',
	];

	import { TEMPLATES } from '$lib/data/templates';

	// Filter templates based on active category
	const filteredTemplates = $derived.by(() => {
		if (activeCategory === 'Semua') return TEMPLATES;
		return TEMPLATES.filter((t) => t.category === activeCategory);
	});

	function handlePreview(template: TemplateItem) {
		selectedTemplate = template;
		isPreviewOpen = true;
	}

	function closePreview() {
		isPreviewOpen = false;
	}

	const faqs = [
		{
			q: 'Berapa lama proses pembuatan undangan digital sampai siap dibagikan?',
			a: 'Hanya dalam hitungan 5 hingga 15 menit! Anda cukup memilih template favorit, melengkapi detail acara (nama pengantin, waktu, lokasi peta, dan rekening), lalu link undangan siap disebarkan langsung ke tamu melalui WhatsApp.',
		},
		{
			q: 'Apakah saya bisa mengubah detail acara jika ada perubahan tanggal atau waktu?',
			a: 'Tentu saja! Anda memiliki akses dashboard penuh untuk mengedit nama, waktu, denah lokasi, lagu pengiring, maupun foto galeri kapan saja secara real-time tanpa biaya tambahan.',
		},
		{
			q: 'Bagaimana cara kerja RSVP tamu via WhatsApp?',
			a: 'Ketika tamu membuka link undangan Anda dan menekan tombol RSVP (Konfirmasi Hadir / Berhalangan), data tamu otomatis tersimpan di dashboard Anda dan dapat langsung mengirimkan pesan konfirmasi personal ke WhatsApp pengantin.',
		},
		{
			q: 'Berapa lama masa aktif undangan digital saya?',
			a: 'Semua undangan di Ketuk.id aktif selamanya tanpa biaya perpanjangan berkala. Link undangan Anda tetap dapat diakses sebagai kenang-kenangan manis pernikahan Anda.',
		},
		{
			q: 'Apakah bisa menambahkan lagu sendiri dan amplop digital transfer?',
			a: 'Ya! Anda bisa memilih koleksi musik romantis yang tersedia atau menyematkan musik pilihan Anda. Fitur amplop digital mendukung nomor rekening bank (BCA, Mandiri, BRI, BNI) serta QRIS.',
		},
	];
</script>

<svelte:head>
	<title>Ketuk.id | Undangan Pernikahan Digital Eksklusif & RSVP</title>
	<meta
		name="description"
		content="Desain undangan pernikahan digital haute-couture dengan RSVP WhatsApp instan, peta Google Maps, amplop digital, dan musik pengiring. Momen indah dimulai dengan undangan berkelas."
	/>
</svelte:head>

<!-- Hero Section (Replicating The Digital Yes Layout & Aesthetic) -->
<section class="relative overflow-hidden bg-cream-50 pt-10 pb-16 sm:pt-16 sm:pb-24 lg:pt-20 lg:pb-32">
	<!-- Background subtle architectural radial warmth -->
	<div
		class="absolute top-0 right-1/4 -z-10 h-[600px] w-[600px] rounded-full bg-[radial-gradient(circle_at_center,rgba(214,188,152,0.25),transparent_70%)] blur-3xl"
	></div>
	<div
		class="absolute bottom-0 left-10 -z-10 h-[500px] w-[500px] rounded-full bg-[radial-gradient(circle_at_center,rgba(110,27,44,0.08),transparent_70%)] blur-3xl"
	></div>

	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="grid items-center gap-12 lg:grid-cols-12 lg:gap-8">
			<!-- Left Column: Haute-Couture Typography & CTAs -->
			<div class="text-center lg:col-span-7 lg:text-left">
				<!-- Eyebrow Tag -->
				<div class="inline-flex items-center gap-2 rounded-full border border-champagne-300/80 bg-cream-100/90 px-4 py-1.5 shadow-2xs backdrop-blur-xs">
					<span class="h-1.5 w-1.5 rounded-full bg-champagne-600"></span>
					<span class="text-[10px] sm:text-[11px] font-semibold tracking-[0.25em] uppercase text-coffee-700">
						Undangan Digital Pernikahan & Acara Eksklusif
					</span>
				</div>

				<!-- Monumental Serif Headline (Lora) -->
				<h1
					class="font-serif mt-6 text-4xl sm:text-6xl lg:text-7xl font-medium tracking-tight text-coffee-950 leading-[1.08]"
				>
					Pernikahan <br class="hidden sm:inline" />
					dimulai dengan <br class="hidden sm:inline" />
					<span class="italic text-coffee-800 font-normal">undangan.</span>
				</h1>

				<!-- Refined Descriptive Paragraph -->
				<p
					class="mx-auto mt-6 max-w-xl text-base sm:text-lg text-coffee-700 font-normal leading-relaxed lg:mx-0"
				>
					Desain mewah yang siap memukau setiap tamu Anda. Lengkap dengan RSVP instan via WhatsApp,
					peta lokasi akurat, musik pengiring, dan amplop digital dalam satu link eksklusif.
				</p>

				<!-- Dual Action Buttons -->
				<div class="mt-8 flex flex-col sm:flex-row items-center justify-center lg:justify-start gap-3.5">
					<a
						href="#katalog"
						class="w-full sm:w-auto inline-flex items-center justify-center gap-2 rounded-full bg-coffee-800 px-8 py-4 text-xs font-bold tracking-widest uppercase text-white shadow-md transition-all duration-200 hover:bg-coffee-900 hover:shadow-lg active:scale-98"
					>
						<span>Jelajah Desain</span>
						<ArrowRight size={15} />
					</a>
					<a
						href="https://wa.me/6281288889999?text=Halo%20Ketuk.id,%20saya%20ingin%20konsultasi%20undangan%20pernikahan"
						target="_blank"
						rel="noopener noreferrer"
						class="w-full sm:w-auto inline-flex items-center justify-center gap-2 rounded-full border border-coffee-300 bg-white/80 px-7 py-4 text-xs font-semibold tracking-wider uppercase text-coffee-800 transition-all hover:border-coffee-900 hover:bg-white active:scale-98"
					>
						<MessageCircle size={15} class="text-emerald-700" />
						<span>Hubungi Kami</span>
					</a>
				</div>

				<!-- Social Proof & Trust Badges -->
				<div class="mt-10 pt-6 border-t border-cream-200/80 max-w-lg mx-auto lg:mx-0">
					<!-- Star Rating -->
					<div class="flex items-center justify-center lg:justify-start gap-3">
						<div class="flex text-amber-500">
							{#each Array(5) as _}
								<Star size={16} class="fill-amber-400 text-amber-400" />
							{/each}
						</div>
						<span class="text-xs font-semibold text-coffee-900">
							Penilaian 4.9/5
						</span>
						<span class="text-xs text-coffee-500">
							(500+ ulasan pasangan bahagia)
						</span>
					</div>

					<!-- Checklist points -->
					<div class="mt-4 flex flex-wrap items-center justify-center lg:justify-start gap-x-5 gap-y-2 text-xs font-medium text-coffee-700">
						<span class="inline-flex items-center gap-1.5">
							<Check size={14} class="text-coffee-700" />
							100% Kustomisasi
						</span>
						<span class="inline-flex items-center gap-1.5">
							<Check size={14} class="text-coffee-700" />
							Tanpa Biaya Tersembunyi
						</span>
						<span class="inline-flex items-center gap-1.5">
							<Check size={14} class="text-coffee-700" />
							Aktif Selamanya
						</span>
					</div>
				</div>
			</div>

			<!-- Right Column: Interactive Realistic Smartphone Mockup -->
			<div class="lg:col-span-5 flex items-center justify-center">
				<PhoneMockup coupleNames="Sarah & Dimas" eventDate="Sabtu, 24 Oktober 2026" />
			</div>
		</div>
	</div>
</section>

<!-- Catalog Section ("Temukan undangan Anda.") -->
<section id="katalog" class="scroll-mt-14 py-16 sm:py-24 bg-cream-50 border-t border-cream-200/60">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<!-- Event Type Segmented Tabs (Undangan Pernikahan vs Ulang Tahun / Lainnya) -->
		<div class="flex justify-center">
			<div
				class="inline-flex items-center rounded-full border border-cream-200 bg-cream-100/90 p-1.5 shadow-2xs"
			>
				<button
					type="button"
					class="rounded-full px-5 py-2 text-xs font-semibold tracking-wide transition-all {eventTab ===
					'wedding'
						? 'bg-coffee-800 text-white shadow-xs'
						: 'text-coffee-700 hover:text-coffee-800'}"
					onclick={() => (eventTab = 'wedding')}
				>
					Undangan Pernikahan
				</button>
				<button
					type="button"
					class="rounded-full px-5 py-2 text-xs font-semibold tracking-wide transition-all {eventTab ===
					'other'
						? 'bg-coffee-800 text-white shadow-xs'
						: 'text-coffee-700 hover:text-coffee-800'}"
					onclick={() => (eventTab = 'other')}
				>
					Undangan Ulang Tahun & Acara Lain
				</button>
			</div>
		</div>

		<!-- Section Header -->
		<div class="mt-10 text-center max-w-2xl mx-auto">
			<h2 class="font-serif text-3xl sm:text-5xl font-medium tracking-tight text-coffee-950">
				Temukan undangan Anda.
			</h2>
			<p class="mt-3 text-sm sm:text-base text-coffee-600 leading-relaxed">
				Koleksi pernikahan kami dirancang dengan teliti untuk setiap gaya, kepribadian, dan cerita cinta istimewa Anda.
			</p>
		</div>

		<!-- Filter Pills (Scrollable with touch momentum on Android/Tablet) -->
		<div class="mt-8 flex items-center justify-start sm:justify-center overflow-x-auto pb-2 scrollbar-none px-2">
			<div class="inline-flex items-center gap-2">
				{#each filterCategories as cat (cat)}
					<button
						type="button"
						class="shrink-0 rounded-full px-4 py-2 text-xs font-medium tracking-wide transition-all {activeCategory ===
						cat
							? 'bg-coffee-800 text-white shadow-xs'
							: 'border border-cream-300/80 bg-white/80 text-coffee-700 hover:bg-cream-100 hover:border-cream-400'}"
						onclick={() => (activeCategory = cat)}
					>
						{cat}
					</button>
				{/each}
			</div>
		</div>

		<!-- Responsive Grid of Cards (1 col mobile, 2-3 col tablet, 4 col desktop) -->
		<div
			class="mt-10 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 sm:gap-6 lg:gap-7"
		>
			{#each filteredTemplates as template (template.id)}
				<TemplateCard {template} onPreview={handlePreview} />
			{/each}
		</div>

		<!-- Bottom Catalog Callout -->
		<div class="mt-16 text-center">
			<p class="font-serif text-lg text-coffee-800 italic">
				Ingin desain kustom yang sepenuhnya unik untuk tema pesta Anda?
			</p>
			<div class="mt-3">
				<a
					href="https://wa.me/6281288889999?text=Halo%20Ketuk.id,%20saya%20tertarik%20dengan%20custom%20undangan%20eksklusif"
					target="_blank"
					rel="noopener noreferrer"
					class="inline-flex items-center gap-2 rounded-full border border-coffee-800 bg-white px-6 py-2.5 text-xs font-semibold tracking-wider uppercase text-coffee-800 shadow-2xs hover:bg-coffee-50 transition-colors"
				>
					<span>Konsultasi Custom Design</span>
					<ArrowRight size={13} />
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Features Section ("Fitur Lengkap yang Elegan") -->
<section id="fitur" class="scroll-mt-14 py-16 sm:py-24 bg-cream-100/60 border-t border-cream-200/60">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="max-w-2xl mx-auto text-center">
			<span class="text-[11px] font-semibold tracking-[0.2em] uppercase text-coffee-800">
				Fitur Eksklusif
			</span>
			<h2 class="font-serif mt-2 text-3xl sm:text-4xl font-medium text-coffee-950">
				Bukan sekadar tautan undangan biasa
			</h2>
			<p class="mt-3 text-sm text-coffee-600">
				Segala yang Anda dan para tamu butuhkan dirancang secara intuitif dalam satu sentuhan berkelas.
			</p>
		</div>

		<div class="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
			<!-- Feature 1 -->
			<div class="rounded-2xl border border-cream-200 bg-white p-6 shadow-2xs transition-shadow hover:shadow-md">
				<div class="flex h-12 w-12 items-center justify-center rounded-xl bg-coffee-50 text-coffee-800">
					<MessageCircle size={22} />
				</div>
				<h3 class="font-serif mt-4 text-xl font-semibold text-coffee-900">
					Konfirmasi RSVP via WhatsApp
				</h3>
				<p class="mt-2 text-xs sm:text-sm text-coffee-600 leading-relaxed">
					Tamu mengonfirmasi kehadiran dalam hitungan detik. Data jumlah tamu dan ucapan doa otomatis tercatat di dashboard Anda.
				</p>
			</div>

			<!-- Feature 2 -->
			<div class="rounded-2xl border border-cream-200 bg-white p-6 shadow-2xs transition-shadow hover:shadow-md">
				<div class="flex h-12 w-12 items-center justify-center rounded-xl bg-champagne-100 text-champagne-700">
					<MapPin size={22} />
				</div>
				<h3 class="font-serif mt-4 text-xl font-semibold text-coffee-900">
					Petunjuk Arah Google Maps
				</h3>
				<p class="mt-2 text-xs sm:text-sm text-coffee-600 leading-relaxed">
					Tamu tidak akan tersasar berkat integrasi titik koordinat langsung ke aplikasi Google Maps & Waze dengan satu klik.
				</p>
			</div>

			<!-- Feature 3 -->
			<div class="rounded-2xl border border-cream-200 bg-white p-6 shadow-2xs transition-shadow hover:shadow-md">
				<div class="flex h-12 w-12 items-center justify-center rounded-xl bg-rose-50 text-rose-800">
					<Music size={22} />
				</div>
				<h3 class="font-serif mt-4 text-xl font-semibold text-coffee-900">
					Musik Latar & Lagu Romantis
				</h3>
				<p class="mt-2 text-xs sm:text-sm text-coffee-600 leading-relaxed">
					Alunan melodi favorit menyambut tamu saat undangan dibuka, menghadirkan atmosfer hangat dan menyentuh hati.
				</p>
			</div>

			<!-- Feature 4 -->
			<div class="rounded-2xl border border-cream-200 bg-white p-6 shadow-2xs transition-shadow hover:shadow-md">
				<div class="flex h-12 w-12 items-center justify-center rounded-xl bg-emerald-50 text-emerald-800">
					<Gift size={22} />
				</div>
				<h3 class="font-serif mt-4 text-xl font-semibold text-coffee-900">
					Amplop Digital & Kirim Kado
				</h3>
				<p class="mt-2 text-xs sm:text-sm text-coffee-600 leading-relaxed">
					Mudahkan kerabat dan sahabat mengirimkan tanda kasih melalui transfer nomor rekening resmi atau QRIS instan.
				</p>
			</div>

			<!-- Feature 5 -->
			<div class="rounded-2xl border border-cream-200 bg-white p-6 shadow-2xs transition-shadow hover:shadow-md">
				<div class="flex h-12 w-12 items-center justify-center rounded-xl bg-amber-50 text-amber-800">
					<Clock size={22} />
				</div>
				<h3 class="font-serif mt-4 text-xl font-semibold text-coffee-900">
					Countdown & Simpan ke Kalender
				</h3>
				<p class="mt-2 text-xs sm:text-sm text-coffee-600 leading-relaxed">
					Hitung mundur hari bahagia serta tombol "Add to Calendar" agar tamu tidak melewatkan hari spesial Anda.
				</p>
			</div>

			<!-- Feature 6 -->
			<div class="rounded-2xl border border-cream-200 bg-white p-6 shadow-2xs transition-shadow hover:shadow-md">
				<div class="flex h-12 w-12 items-center justify-center rounded-xl bg-blue-50 text-blue-800">
					<Shield size={22} />
				</div>
				<h3 class="font-serif mt-4 text-xl font-semibold text-coffee-900">
					Privasi & Masa Aktif Selamanya
				</h3>
				<p class="mt-2 text-xs sm:text-sm text-coffee-600 leading-relaxed">
					Data tamu dan kenangan foto Anda aman terjaga. Link undangan tidak akan kedaluwarsa setelah acara selesai.
				</p>
			</div>
		</div>
	</div>
</section>

<!-- How it works ("3 Langkah Mudah") -->
<section id="cara-kerja" class="scroll-mt-14 py-16 sm:py-24 bg-cream-50">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="text-center max-w-xl mx-auto">
			<span class="text-[11px] font-semibold tracking-[0.2em] uppercase text-coffee-800">
				Proses Mudah & Cepat
			</span>
			<h2 class="font-serif mt-2 text-3xl sm:text-4xl font-medium text-coffee-950">
				Tiga langkah menuju undangan impian
			</h2>
		</div>

		<div class="mt-14 grid gap-8 md:grid-cols-3">
			<div class="relative flex flex-col items-center text-center p-6">
				<div class="flex h-14 w-14 items-center justify-center rounded-full bg-coffee-800 text-white font-serif text-xl font-semibold shadow-md">
					1
				</div>
				<h3 class="font-serif mt-5 text-xl font-semibold text-coffee-900">Pilih Desain Favorit</h3>
				<p class="mt-2 text-xs sm:text-sm text-coffee-600 leading-relaxed">
					Eksplorasi katalog kurasi kami dengan beragam tema estetis dari klasik hingga modern.
				</p>
			</div>

			<div class="relative flex flex-col items-center text-center p-6">
				<div class="flex h-14 w-14 items-center justify-center rounded-full bg-coffee-800 text-white font-serif text-xl font-semibold shadow-md">
					2
				</div>
				<h3 class="font-serif mt-5 text-xl font-semibold text-coffee-900">Personalisasi Konten</h3>
				<p class="mt-2 text-xs sm:text-sm text-coffee-600 leading-relaxed">
					Isi identitas pasangan, waktu akad & resepsi, foto pre-wedding, dan musik pilihan Anda.
				</p>
			</div>

			<div class="relative flex flex-col items-center text-center p-6">
				<div class="flex h-14 w-14 items-center justify-center rounded-full bg-coffee-800 text-white font-serif text-xl font-semibold shadow-md">
					3
				</div>
				<h3 class="font-serif mt-5 text-xl font-semibold text-coffee-900">Sebarkan ke Tamu</h3>
				<p class="mt-2 text-xs sm:text-sm text-coffee-600 leading-relaxed">
					Dapatkan link eksklusif Anda dan bagikan langsung ke keluarga dan sahabat via WhatsApp.
				</p>
			</div>
		</div>

		<div class="mt-10 text-center">
			<a
				href="/daftar"
				class="inline-flex items-center gap-2 rounded-full bg-coffee-800 px-8 py-3.5 text-xs font-semibold uppercase tracking-wider text-white shadow-md hover:bg-coffee-900 transition-all"
			>
				<span>Mulai Bikin Undangan Sekarang</span>
				<ArrowRight size={14} />
			</a>
		</div>
	</div>
</section>

<!-- Testimonial Section -->
<section class="py-16 sm:py-24 bg-cream-100/70 border-t border-cream-200/60">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="text-center max-w-xl mx-auto">
			<span class="text-[11px] font-semibold tracking-[0.2em] uppercase text-coffee-800">
				Cerita Pengantin
			</span>
			<h2 class="font-serif mt-2 text-3xl sm:text-4xl font-medium text-coffee-950">
				Dipercaya ratusan pasangan bahagia
			</h2>
		</div>

		<div class="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
			<div class="rounded-2xl border border-cream-200 bg-white p-6 shadow-xs">
				<div class="flex text-amber-500 mb-3">
					{#each Array(5) as _}
						<Star size={14} class="fill-amber-400 text-amber-400" />
					{/each}
				</div>
				<p class="text-xs sm:text-sm text-coffee-700 italic leading-relaxed">
					"Banyak tamu yang memuji undangan pernikahan kami. Desainnya sangat mahal dan tidak pasaran, musiknya menyentuh, dan RSVP-nya sangat memudahkan kami merekap data katering."
				</p>
				<div class="mt-4 pt-3 border-t border-cream-100 flex items-center gap-3">
					<div class="h-9 w-9 rounded-full bg-coffee-100 text-coffee-800 flex items-center justify-center font-serif font-semibold text-sm">
						D & R
					</div>
					<div>
						<h4 class="font-serif text-sm font-semibold text-coffee-900">Dian & Raditya</h4>
						<p class="text-[10px] text-coffee-500">Jakarta • Tema Firenze</p>
					</div>
				</div>
			</div>

			<div class="rounded-2xl border border-cream-200 bg-white p-6 shadow-xs">
				<div class="flex text-amber-500 mb-3">
					{#each Array(5) as _}
						<Star size={14} class="fill-amber-400 text-amber-400" />
					{/each}
				</div>
				<p class="text-xs sm:text-sm text-coffee-700 italic leading-relaxed">
					"Proses editnya gampang banget bahkan dari HP Android saya. Dalam 10 menit undangan langsung jadi dan siap kirim ke grup keluarga. Terima kasih Ketuk.id!"
				</p>
				<div class="mt-4 pt-3 border-t border-cream-100 flex items-center gap-3">
					<div class="h-9 w-9 rounded-full bg-champagne-100 text-champagne-800 flex items-center justify-center font-serif font-semibold text-sm">
						N & A
					</div>
					<div>
						<h4 class="font-serif text-sm font-semibold text-coffee-900">Nadia & Adnan</h4>
						<p class="text-[10px] text-coffee-500">Bandung • Tema Capri Rosa</p>
					</div>
				</div>
			</div>

			<div class="rounded-2xl border border-cream-200 bg-white p-6 shadow-xs sm:col-span-2 lg:col-span-1">
				<div class="flex text-amber-500 mb-3">
					{#each Array(5) as _}
						<Star size={14} class="fill-amber-400 text-amber-400" />
					{/each}
				</div>
				<p class="text-xs sm:text-sm text-coffee-700 italic leading-relaxed">
					"Fitur amplop digital dan RSVP WhatsApp-nya sangat membantu. Teman-teman yang berhalangan hadir tetap bisa kirim hadiah dengan gampang via transfer bank dan QRIS."
				</p>
				<div class="mt-4 pt-3 border-t border-cream-100 flex items-center gap-3">
					<div class="h-9 w-9 rounded-full bg-rose-100 text-rose-800 flex items-center justify-center font-serif font-semibold text-sm">
						C & F
					</div>
					<div>
						<h4 class="font-serif text-sm font-semibold text-coffee-900">Clarissa & Fajar</h4>
						<p class="text-[10px] text-coffee-500">Bali • Tema Tuscan Garden</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- FAQ Section -->
<section id="faq" class="scroll-mt-14 py-16 sm:py-24 bg-cream-50">
	<div class="mx-auto max-w-4xl px-4 sm:px-6">
		<div class="text-center max-w-xl mx-auto">
			<span class="text-[11px] font-semibold tracking-[0.2em] uppercase text-coffee-800">
				Bantuan
			</span>
			<h2 class="font-serif mt-2 text-3xl sm:text-4xl font-medium text-coffee-950">
				Pertanyaan yang sering diajukan
			</h2>
		</div>

		<div class="mt-12 space-y-3">
			{#each faqs as faq, i (faq.q)}
				<div
					class="overflow-hidden rounded-2xl border border-cream-200 bg-white transition-all duration-200 {openFaq ===
					i
						? 'border-champagne-400/60 shadow-xs'
						: 'hover:border-cream-300'}"
				>
					<button
						type="button"
						class="flex w-full items-center justify-between p-5 text-left text-sm sm:text-base font-medium text-coffee-900"
						onclick={() => (openFaq = openFaq === i ? null : i)}
						aria-expanded={openFaq === i}
					>
						<span class="font-serif text-base sm:text-lg">{faq.q}</span>
						<ChevronDown
							size={18}
							class="shrink-0 text-coffee-500 transition-transform duration-200 {openFaq === i
								? 'rotate-180 text-coffee-800'
								: ''}"
						/>
					</button>

					{#if openFaq === i}
						<div class="px-5 pb-5 pt-0 text-xs sm:text-sm text-coffee-600 leading-relaxed border-t border-cream-100/60 mt-1">
							{faq.a}
						</div>
					{/if}
				</div>
			{/each}
		</div>
	</div>
</section>

<!-- Final Call to Action Banner -->
<section class="relative overflow-hidden bg-coffee-900 py-16 sm:py-24 text-white">
	<div
		class="absolute inset-0 bg-[radial-gradient(circle_at_center,rgba(197,164,120,0.2),transparent_70%)]"
	></div>

	<div class="relative mx-auto max-w-4xl px-4 text-center sm:px-6">
		<span class="font-display text-[10px] sm:text-xs font-semibold tracking-[0.25em] text-champagne-300 uppercase">
			Momen Bahagia Dimulai Di Sini
		</span>

		<h2 class="font-serif mt-4 text-3xl sm:text-5xl lg:text-6xl font-medium tracking-tight text-white leading-tight">
			Siap membuat undangan pernikahan impian Anda?
		</h2>

		<p class="mx-auto mt-4 max-w-xl text-sm sm:text-base text-champagne-100/80 font-normal leading-relaxed">
			Bergabunglah bersama ribuan pasangan pengantin lainnya. Mulai buat undangan Anda secara gratis, pilih desain favorit, dan bagikan cinta Anda hari ini.
		</p>

		<div class="mt-8 flex flex-col sm:flex-row items-center justify-center gap-3">
			<a
				href="/daftar"
				class="w-full sm:w-auto inline-flex items-center justify-center gap-2 rounded-full bg-white px-8 py-3.5 text-xs font-bold tracking-widest uppercase text-coffee-950 shadow-lg hover:bg-champagne-100 transition-all active:scale-98"
			>
				<span>Mulai Buat Undangan</span>
				<ArrowRight size={14} />
			</a>
			<a
				href="#katalog"
				class="w-full sm:w-auto inline-flex items-center justify-center gap-2 rounded-full border border-white/30 bg-white/10 px-7 py-3.5 text-xs font-semibold tracking-wider uppercase text-white hover:bg-white/15 backdrop-blur-xs transition-all active:scale-98"
			>
				<span>Lihat Semua Model</span>
			</a>
		</div>
	</div>
</section>

<!-- Interactive Live Preview Modal -->
<TemplatePreviewModal open={isPreviewOpen} template={selectedTemplate} onclose={closePreview} />
