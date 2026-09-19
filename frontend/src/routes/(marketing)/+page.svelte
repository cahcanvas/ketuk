<script lang="ts">
	import { ArrowDown, ArrowRight, CalendarDays, Check, MapPin, Users } from '@lucide/svelte';
	import { EVENT_CATEGORIES, type EventCategoryId } from '@ketuk/shared';
	import {
		TemplateCard,
		TemplatePreviewModal,
		type TemplateItem,
	} from '$lib/components/domain';
	import { TEMPLATES } from '$lib/data/templates';

	type EventFilter = EventCategoryId | 'all';

	const eventFilters: { value: EventFilter; label: string }[] = [
		{ value: 'all', label: 'Semua acara' },
		...EVENT_CATEGORIES.map((category) => ({ value: category.id, label: category.label })),
	];

	const styleFilters = ['Semua', 'Klasik', 'Minimalis', 'Elegan', 'Modern', 'Botanical', 'Luxury Gold'];

	let activeEvent = $state<EventFilter>('all');
	let activeStyle = $state('Semua');
	let selectedTemplate = $state<TemplateItem | null>(null);
	let isPreviewOpen = $state(false);
	let openFaq = $state<number | null>(0);

	const selectedCategory = $derived(
		activeEvent === 'all' ? undefined : EVENT_CATEGORIES.find((category) => category.id === activeEvent),
	);
	const selectedLabel = $derived(selectedCategory?.label ?? 'acara apa pun');
	const filteredTemplates = $derived.by(() =>
		TEMPLATES.filter((template) => {
			const matchesEvent = activeEvent === 'all' || template.eventTypes?.includes(activeEvent);
			const matchesStyle = activeStyle === 'Semua' || template.category === activeStyle;
			return matchesEvent && matchesStyle;
		}),
	);

	const faqs = [
		{
			question: 'Acara apa saja yang bisa dibuat di Ketuk.id?',
			answer:
				'Kamu bisa membuat undangan untuk pernikahan, ulang tahun, aqiqah, khitanan, syukuran, wisuda, reuni, acara komunitas, dan acara bisnis. Untuk jenis acara lain, pilih Lainnya saat onboarding.',
		},
		{
			question: 'Berapa lama undangan siap dibagikan?',
			answer:
				'Setelah memilih desain, isi judul acara, waktu, lokasi, dan detail yang dibutuhkan. Kamu bisa membagikan link setelah informasinya siap, lalu mengubahnya kapan saja dari dashboard.',
		},
		{
			question: 'Apakah tamu bisa mengonfirmasi kehadiran?',
			answer:
				'Bisa. RSVP tamu masuk ke dashboard acara, jadi kamu tidak perlu mengumpulkan jawaban dari banyak chat secara manual.',
		},
		{
			question: 'Apakah satu akun bisa mengelola beberapa acara?',
			answer:
				'Bisa. Preferensi onboarding membantu menata dashboard, sementara setiap undangan tetap menjadi acara tersendiri yang bisa dikelola dari akun yang sama.',
		},
	];

	function chooseEvent(value: EventFilter) {
		activeEvent = value;
		activeStyle = 'Semua';
		document.getElementById('katalog')?.scrollIntoView({ behavior: 'smooth' });
	}

	function handlePreview(template: TemplateItem) {
		selectedTemplate = template;
		isPreviewOpen = true;
	}
</script>

<svelte:head>
	<title>Ketuk.id | Undangan untuk Acara Kamu</title>
	<meta
		name="description"
		content="Buat undangan digital untuk pernikahan, ulang tahun, acara keluarga, komunitas, dan bisnis. Pilih desain, isi detail acara, lalu bagikan link-nya."
	/>
</svelte:head>

<!--
	Design Read: editorial stationery untuk produk acara Indonesia.
	ENERGY 2 / RHYTHM 3 / MOTION 1.

	Hero sengaja tidak memakai badge, glow, angka social proof, atau mockup
	telepon generik. Kartu di kanan berfungsi sebagai artefak produk: ia
	menunjukkan bentuk undangan yang akan dibuat user, bukan dekorasi kosong.
-->
<section class="overflow-hidden bg-cream-50">
	<div class="mx-auto grid max-w-7xl gap-12 px-5 py-14 sm:px-8 sm:py-20 lg:grid-cols-[1fr_0.82fr] lg:items-center lg:gap-20 lg:px-12 lg:py-24">
		<div class="max-w-2xl">
			<p class="text-sm font-medium text-terracotta-600">Ketuk.id</p>
			<h1 class="mt-5 max-w-xl font-serif text-5xl leading-[0.98] tracking-[-0.04em] text-coffee-950 sm:text-6xl lg:text-7xl">
				Undangan yang mengikuti acaramu.
			</h1>
			<p class="mt-6 max-w-lg text-base leading-7 text-coffee-700 sm:text-lg">
				Satu tempat untuk menyiapkan undangan, mengatur tamu, dan membagikan informasi acara. Mulai dari jenis acaranya, lalu isi seperlunya.
			</p>

			<div class="mt-8 flex flex-col gap-3 sm:flex-row">
				<a
					href="/daftar"
					class="inline-flex min-h-11 items-center justify-center gap-2 rounded-lg bg-coffee-900 px-5 py-3 text-sm font-semibold text-white transition-colors hover:bg-coffee-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-terracotta-500"
				>
					Buat undangan
					<ArrowRight size={16} />
				</a>
				<a
					href="#katalog"
					class="inline-flex min-h-11 items-center justify-center gap-2 rounded-lg border border-coffee-300 px-5 py-3 text-sm font-semibold text-coffee-800 transition-colors hover:bg-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-terracotta-500"
				>
					Lihat desain
				</a>
			</div>

			<div class="mt-12 border-t border-coffee-200 pt-5">
				<p class="text-xs font-semibold uppercase tracking-[0.16em] text-coffee-500">Mulai dari sini</p>
				<div class="mt-3 flex flex-wrap gap-x-5 gap-y-2 text-sm text-coffee-700">
					{#each ['Pernikahan', 'Ulang Tahun', 'Acara Keluarga', 'Komunitas'] as label (label)}
						<span class="inline-flex items-center gap-2">
							<span class="h-1.5 w-1.5 rounded-full bg-terracotta-500" aria-hidden="true"></span>
							{label}
						</span>
					{/each}
				</div>
			</div>
		</div>

		<div class="relative mx-auto w-full max-w-md lg:mr-0">
			<div class="absolute -right-8 -top-8 h-32 w-32 rounded-full border border-champagne-300/70" aria-hidden="true"></div>
			<div class="absolute -bottom-8 -left-8 h-24 w-24 border-l border-b border-terracotta-300" aria-hidden="true"></div>
			<div class="relative rotate-[-2deg] border border-coffee-300 bg-white p-3 shadow-[14px_18px_0_0_rgba(93,64,55,0.12)] transition-transform duration-300 hover:rotate-0">
				<div class="border border-champagne-300 px-6 py-8 text-center sm:px-10 sm:py-12">
					<p class="text-[10px] font-semibold uppercase tracking-[0.24em] text-terracotta-600">Undangan acara</p>
					<h2 class="mt-8 font-serif text-4xl leading-tight text-coffee-900 sm:text-5xl">Ruang untuk<br />berkumpul.</h2>
					<div class="mx-auto my-8 h-px w-20 bg-champagne-500"></div>
					<p class="font-serif text-lg italic text-coffee-700">Nama acara kamu</p>
					<div class="mt-8 grid grid-cols-2 gap-3 border-t border-coffee-100 pt-5 text-left text-xs text-coffee-600">
						<span class="inline-flex items-center gap-1.5"><CalendarDays size={14} />Tanggal</span>
						<span class="inline-flex items-center gap-1.5"><MapPin size={14} />Lokasi</span>
					</div>
				</div>
			</div>
			<p class="mt-5 text-center text-xs text-coffee-500">Contoh tampilan. Isi dan gaya berubah mengikuti acara yang kamu pilih.</p>
		</div>
	</div>
</section>

<section class="border-y border-coffee-200 bg-white">
	<div class="mx-auto grid max-w-7xl gap-10 px-5 py-14 sm:px-8 lg:grid-cols-[0.7fr_1.3fr] lg:px-12 lg:py-20">
		<div>
			<p class="text-sm font-medium text-terracotta-600">Pilih konteksnya</p>
			<h2 class="mt-3 max-w-sm font-serif text-3xl leading-tight text-coffee-950 sm:text-4xl">Dashboard yang tidak menyamaratakan semua acara.</h2>
			<p class="mt-4 max-w-sm text-sm leading-6 text-coffee-600">Pilih jenis acara saat mulai. Kami akan menata rekomendasi fitur dan langkah berikutnya berdasarkan pilihan itu.</p>
		</div>
		<div class="grid grid-cols-2 gap-px overflow-hidden border border-coffee-200 bg-coffee-200 sm:grid-cols-3">
			{#each EVENT_CATEGORIES as category (category.id)}
				<button
					type="button"
					class="min-h-28 bg-white px-4 py-4 text-left transition-colors hover:bg-cream-50 focus-visible:z-10 focus-visible:outline-2 focus-visible:outline-offset-[-2px] focus-visible:outline-terracotta-500 {activeEvent === category.id ? 'bg-cream-100' : ''}"
					onclick={() => chooseEvent(category.id)}
					aria-pressed={activeEvent === category.id}
				>
					<span class="block text-sm font-semibold text-coffee-900">{category.label}</span>
					<span class="mt-2 block text-xs leading-5 text-coffee-500">{category.shortDesc}</span>
				</button>
			{/each}
		</div>
	</div>
</section>

<section id="fitur" class="bg-coffee-900 text-cream-50">
	<div class="mx-auto grid max-w-7xl gap-12 px-5 py-16 sm:px-8 lg:grid-cols-[0.75fr_1.25fr] lg:gap-24 lg:px-12 lg:py-24">
		<div>
			<p class="text-sm font-medium text-champagne-300">Yang bisa kamu atur</p>
			<h2 class="mt-4 max-w-md font-serif text-4xl leading-tight sm:text-5xl">Informasi yang dibutuhkan tamu, di satu tempat.</h2>
			<p class="mt-5 max-w-md text-sm leading-6 text-cream-200">Tidak semua acara membutuhkan fitur yang sama. Pilih yang relevan, lengkapi detailnya, dan bagikan link yang mudah dibuka.</p>
		</div>
		<div class="divide-y divide-cream-50/15 border-y border-cream-50/20">
			<div class="grid gap-3 py-5 sm:grid-cols-[auto_1fr] sm:gap-6">
				<Users size={22} class="text-champagne-300" aria-hidden="true" />
				<div><h3 class="font-semibold">Daftar tamu dan RSVP</h3><p class="mt-1 text-sm leading-6 text-cream-200">Kumpulkan jawaban hadir dan catatan tamu tanpa mengandalkan spreadsheet terpisah.</p></div>
			</div>
			<div class="grid gap-3 py-5 sm:grid-cols-[auto_1fr] sm:gap-6">
				<MapPin size={22} class="text-champagne-300" aria-hidden="true" />
				<div><h3 class="font-semibold">Tanggal, waktu, dan lokasi</h3><p class="mt-1 text-sm leading-6 text-cream-200">Tampilkan informasi penting dengan format yang mudah dibaca dan link peta yang bisa langsung dibuka.</p></div>
			</div>
			<div class="grid gap-3 py-5 sm:grid-cols-[auto_1fr] sm:gap-6">
				<Check size={22} class="text-champagne-300" aria-hidden="true" />
				<div><h3 class="font-semibold">Galeri dan detail acara</h3><p class="mt-1 text-sm leading-6 text-cream-200">Tambahkan cerita, foto, kontak, dan informasi lain yang memang dibutuhkan oleh tamu kamu.</p></div>
			</div>
		</div>
	</div>
</section>

<section id="katalog" class="scroll-mt-8 bg-cream-50">
	<div class="mx-auto max-w-7xl px-5 py-16 sm:px-8 lg:px-12 lg:py-24">
		<div class="flex flex-col gap-5 border-b border-coffee-200 pb-8 sm:flex-row sm:items-end sm:justify-between">
			<div>
				<p class="text-sm font-medium text-terracotta-600">Koleksi desain</p>
				<h2 class="mt-3 font-serif text-4xl leading-tight text-coffee-950 sm:text-5xl">Untuk {selectedLabel}.</h2>
				<p class="mt-3 max-w-xl text-sm leading-6 text-coffee-600">Pilih berdasarkan acara, lalu saring lagi berdasarkan gaya visual.</p>
			</div>
			<a href="/template" class="inline-flex min-h-11 items-center gap-2 text-sm font-semibold text-coffee-800 hover:text-terracotta-600 focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-terracotta-500">
				Lihat semua desain <ArrowRight size={16} />
			</a>
		</div>

		<div class="mt-7 flex gap-2 overflow-x-auto pb-2" aria-label="Filter jenis acara">
			{#each eventFilters as filter (filter.value)}
				<button
					type="button"
					class="min-h-11 shrink-0 border-b-2 px-1 text-sm {activeEvent === filter.value ? 'border-terracotta-500 font-semibold text-coffee-900' : 'border-transparent text-coffee-500 hover:text-coffee-800'}"
					aria-current={activeEvent === filter.value ? 'true' : undefined}
					onclick={() => (activeEvent = filter.value)}
				>
					{filter.label}
				</button>
			{/each}
		</div>
		<div class="mt-3 flex gap-2 overflow-x-auto pb-2" aria-label="Filter gaya desain">
			{#each styleFilters as style (style)}
				<button
					type="button"
					class="min-h-10 shrink-0 rounded-full border px-3.5 text-xs {activeStyle === style ? 'border-coffee-900 bg-coffee-900 text-white' : 'border-coffee-200 text-coffee-600 hover:border-coffee-400'}"
					aria-current={activeStyle === style ? 'true' : undefined}
					onclick={() => (activeStyle = style)}
				>
					{style}
				</button>
			{/each}
		</div>

		{#if filteredTemplates.length > 0}
			<div class="mt-10 grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
				{#each filteredTemplates as template (template.id)}
					<TemplateCard {template} onPreview={handlePreview} />
				{/each}
			</div>
		{:else}
			<div class="mt-10 border border-dashed border-coffee-300 bg-white px-6 py-12 text-center">
				<h3 class="font-serif text-2xl text-coffee-900">Desain untuk kategori ini sedang kami siapkan.</h3>
				<p class="mx-auto mt-3 max-w-md text-sm leading-6 text-coffee-600">Lihat koleksi yang tersedia sekarang, atau mulai dari pilihan acara lain.</p>
				<button type="button" class="mt-5 min-h-11 rounded-lg bg-coffee-900 px-5 py-3 text-sm font-semibold text-white hover:bg-coffee-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-terracotta-500" onclick={() => { activeEvent = 'all'; activeStyle = 'Semua'; }}>Lihat koleksi yang tersedia</button>
			</div>
		{/if}
	</div>
</section>

<section id="cara-kerja" class="border-t border-coffee-200 bg-white">
	<div class="mx-auto max-w-7xl px-5 py-16 sm:px-8 lg:px-12 lg:py-24">
		<div class="max-w-xl">
			<p class="text-sm font-medium text-terracotta-600">Cara kerja</p>
			<h2 class="mt-3 font-serif text-4xl leading-tight text-coffee-950 sm:text-5xl">Dari jenis acara sampai link siap dibagikan.</h2>
		</div>
		<div class="mt-12 grid gap-8 border-t border-coffee-200 pt-8 md:grid-cols-3 md:gap-10">
			<div><p class="font-serif text-5xl text-terracotta-500">01</p><h3 class="mt-5 text-lg font-semibold text-coffee-900">Pilih acaranya</h3><p class="mt-2 text-sm leading-6 text-coffee-600">Pilih kategori dan jenis acara agar dashboard menampilkan hal yang relevan.</p></div>
			<div><p class="font-serif text-5xl text-terracotta-500">02</p><h3 class="mt-5 text-lg font-semibold text-coffee-900">Isi bagian penting</h3><p class="mt-2 text-sm leading-6 text-coffee-600">Masukkan judul, tanggal, lokasi, tamu, dan detail lain sesuai kebutuhan acara.</p></div>
			<div><p class="font-serif text-5xl text-terracotta-500">03</p><h3 class="mt-5 text-lg font-semibold text-coffee-900">Bagikan link-nya</h3><p class="mt-2 text-sm leading-6 text-coffee-600">Kirim link undangan ke keluarga, teman, komunitas, atau rekan kerja.</p></div>
		</div>
	</div>
</section>

<section id="faq" class="border-t border-coffee-200 bg-cream-50">
	<div class="mx-auto grid max-w-7xl gap-10 px-5 py-16 sm:px-8 lg:grid-cols-[0.65fr_1.35fr] lg:gap-24 lg:px-12 lg:py-24">
		<div>
			<p class="text-sm font-medium text-terracotta-600">Pertanyaan umum</p>
			<h2 class="mt-3 font-serif text-4xl leading-tight text-coffee-950">Sebelum kamu mulai.</h2>
		</div>
		<div class="divide-y divide-coffee-200 border-y border-coffee-200">
			{#each faqs as faq, index (faq.question)}
				<div>
					<button type="button" class="flex min-h-16 w-full items-center justify-between gap-5 text-left text-sm font-semibold text-coffee-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-terracotta-500" aria-expanded={openFaq === index} onclick={() => (openFaq = openFaq === index ? null : index)}>
						{faq.question}
						<ArrowDown size={16} class="shrink-0 transition-transform {openFaq === index ? 'rotate-180' : ''}" />
					</button>
					{#if openFaq === index}<p class="pb-5 pr-8 text-sm leading-6 text-coffee-600">{faq.answer}</p>{/if}
				</div>
			{/each}
		</div>
	</div>
</section>

<section class="bg-terracotta-600 text-white">
	<div class="mx-auto flex max-w-7xl flex-col gap-7 px-5 py-14 sm:px-8 lg:flex-row lg:items-center lg:justify-between lg:px-12 lg:py-16">
		<div>
			<p class="text-sm font-medium text-terracotta-100">Siap menyiapkan acara?</p>
			<h2 class="mt-2 max-w-2xl font-serif text-4xl leading-tight sm:text-5xl">Mulai dari undangannya.</h2>
		</div>
		<a href="/daftar" class="inline-flex min-h-11 items-center justify-center gap-2 self-start rounded-lg bg-white px-5 py-3 text-sm font-semibold text-coffee-900 hover:bg-cream-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white lg:self-center">Buat undangan <ArrowRight size={16} /></a>
	</div>
</section>

<TemplatePreviewModal open={isPreviewOpen} template={selectedTemplate} onclose={() => (isPreviewOpen = false)} />
