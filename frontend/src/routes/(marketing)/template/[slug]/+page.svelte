<script lang="ts">
	import {
		ArrowLeft,
		ArrowRight,
		Check,
		Star,
		ChevronDown,
	} from '@lucide/svelte';
	import { TEMPLATES, type PricingTier } from '$lib/data/templates';
	import EnvelopeCard from '$lib/components/domain/EnvelopeCard.svelte';
	import EnvelopeModal from '$lib/components/domain/EnvelopeModal.svelte';
	import TemplateCard from '$lib/components/domain/TemplateCard.svelte';
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const template = $derived(data.template);

	const defaultTier: PricingTier = {
		id: 'essential',
		name: 'Essential',
		price: 'Rp 149.000',
		description: 'Desain digital esensial',
	};

	let selectedTier = $state<PricingTier>(
		data.template.tiers.find((t) => t.popular) ?? data.template.tiers[0] ?? defaultTier,
	);

	// Update selectedTier when template changes
	$effect(() => {
		selectedTier = template.tiers.find((t) => t.popular) ?? template.tiers[0] ?? defaultTier;
	});

	let isExpandedOpen = $state(false);
	let openFaq = $state<number | null>(0);

	const otherTemplates = $derived(TEMPLATES.filter((t) => t.id !== template.id).slice(0, 4));

	const faqs = [
		{
			q: 'Kapan undangan pernikahan saya siap setelah melakukan pemesanan?',
			a: 'Undangan Anda langsung aktif dan dapat diisi detailnya secara instan! Rata-rata pasangan menyelesaikan kustomisasi dalam waktu 10 hingga 15 menit.',
		},
		{
			q: 'Apakah saya bisa mengubah tanggal, waktu, atau lokasi setelah undangan selesai?',
			a: 'Bisa, tanpa batasan! Anda memiliki dashboard pribadi untuk mengedit informasi acara kapan saja tanpa biaya tambahan.',
		},
		{
			q: 'Bagaimana tamu mengisi RSVP dan konfirmasi kehadiran?',
			a: 'Tamu cukup menekan tombol RSVP di undangan. Pilihan kehadiran langsung tercatat di dashboard Anda dan dapat langsung mengirimkan pesan konfirmasi ke WhatsApp Anda.',
		},
		{
			q: 'Apakah ada batasan jumlah tamu yang bisa menerima undangan ini?',
			a: 'Tidak ada batasan kuota sama sekali. Anda dapat membagikan tautan undangan ke 100, 500, hingga ribuan tamu.',
		},
	];
</script>

<svelte:head>
	<title>{template.title} | Undangan Pernikahan Eksklusif Ketuk.id</title>
	<meta
		name="description"
		content="{template.title} - {template.subtitle}. Dilengkapi RSVP WhatsApp, peta Google Maps, amplop digital, dan musik pengiring."
	/>
</svelte:head>

<div class="bg-cream-50 min-h-screen py-6 sm:py-10">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<!-- Back to Collection Breadcrumb -->
		<div class="mb-6 sm:mb-8">
			<a
				href="/template"
				class="inline-flex items-center gap-2 text-xs font-semibold uppercase tracking-widest text-coffee-600 hover:text-coffee-800 transition-colors"
			>
				<ArrowLeft size={14} />
				<span>KEMBALI KE SEMUA DESAIN</span>
			</a>
		</div>

		<!-- Main Hero Stage: Left Envelope / Right Product Spec (Matching Screenshot 1) -->
		<div class="grid gap-10 lg:grid-cols-12 lg:gap-12 items-start">
			<!-- Left Column: Envelope Presentation & Interactive Demo -->
			<div class="lg:col-span-6 flex flex-col items-center justify-center w-full">
				<EnvelopeCard {template} onExpand={() => (isExpandedOpen = true)} interactive={true} />
			</div>

			<!-- Right Column: Product Detail, Tier Selector & CTAs -->
			<div class="lg:col-span-6 flex flex-col justify-between">
				<div>
					<!-- Category and Badge -->
					<div class="flex items-center gap-2">
						{#if template.badge}
							<span
								class="inline-block rounded-md bg-coffee-900 px-2.5 py-1 text-[9px] font-semibold tracking-[0.15em] text-champagne-200 uppercase shadow-xs border border-champagne-400/20"
							>
								{template.badge}
							</span>
						{/if}
						<span class="text-xs text-coffee-500 font-medium">Koleksi {template.category}</span>
					</div>

					<!-- Title in Lora -->
					<h1
						class="font-serif mt-3 text-4xl sm:text-5xl lg:text-6xl font-medium tracking-tight text-coffee-950 leading-tight"
					>
						{template.title}
					</h1>

					<!-- Tagline / Subtitle -->
					<p class="mt-3 text-sm sm:text-base text-coffee-700 leading-relaxed max-w-xl font-normal">
						{template.tagline || template.subtitle}
					</p>

					<!-- Rating and Reviews -->
					<div class="mt-4 flex items-center gap-2 text-xs">
						<div class="flex text-amber-500">
							{#each Array(5) as _}
								<Star size={15} class="fill-amber-400 text-amber-400" />
							{/each}
						</div>
						<span class="font-semibold text-coffee-900">{template.rating}</span>
						<span class="text-coffee-500">({template.reviewCount}+ pasangan merekomendasikan)</span>
					</div>

					<!-- Pricing Tier Selector Cards -->
					<div class="mt-8 space-y-3">
						<span class="block text-xs font-semibold uppercase tracking-wider text-coffee-600">
							PILIH PAKET UNDANGAN:
						</span>

						<div class="space-y-2.5">
							{#each template.tiers as tier (tier.id)}
								<button
									type="button"
									class="w-full text-left rounded-2xl border p-4 transition-all duration-200 {selectedTier.id ===
									tier.id
										? 'border-coffee-800 bg-white shadow-md ring-1 ring-coffee-800'
										: 'border-cream-300/80 bg-cream-100/60 hover:bg-cream-100 hover:border-cream-400'}"
									onclick={() => (selectedTier = tier)}
								>
									<div class="flex items-center justify-between">
										<div class="flex items-center gap-2.5">
											<div
												class="h-4 w-4 rounded-full border flex items-center justify-center {selectedTier.id ===
												tier.id
													? 'border-coffee-800 bg-coffee-800'
													: 'border-coffee-400 bg-white'}"
											>
												{#if selectedTier.id === tier.id}
													<span class="h-1.5 w-1.5 rounded-full bg-white"></span>
												{/if}
											</div>
											<div>
												<div class="flex items-center gap-2">
													<span class="font-serif text-lg font-semibold text-coffee-950">
														{tier.name}
													</span>
													{#if tier.popular}
														<span
															class="rounded-full bg-coffee-100 px-2 py-0.5 text-[9px] font-semibold text-coffee-800 uppercase"
														>
															REKOMENDASI
														</span>
													{/if}
												</div>
												<p class="text-xs text-coffee-600 mt-0.5 leading-normal">
													{tier.description}
												</p>
											</div>
										</div>

										<div class="text-right shrink-0 pl-3">
											<span class="font-serif text-lg font-bold text-coffee-900">
												{tier.price}
											</span>
										</div>
									</div>
								</button>
							{/each}
						</div>
					</div>

					<!-- Value Checklist -->
					<div class="mt-6 space-y-2 text-xs text-coffee-700">
						<div class="flex items-center gap-2">
							<Check size={14} class="text-emerald-700 shrink-0" />
							<span>Website undangan lengkap siap kirim dalam satu tautan eksklusif</span>
						</div>
						<div class="flex items-center gap-2">
							<Check size={14} class="text-emerald-700 shrink-0" />
							<span>Dashboard RSVP otomatis & rekap daftar kehadiran tamu</span>
						</div>
						<div class="flex items-center gap-2">
							<Check size={14} class="text-emerald-700 shrink-0" />
							<span>Akses instan dan masa aktif selamanya tanpa biaya perpanjangan</span>
						</div>
					</div>

					<!-- Primary Order Action Button -->
					<div class="mt-8 space-y-3">
						<a
							href="/daftar?template={template.slug}&tier={selectedTier.id}"
							class="w-full flex items-center justify-center gap-2 rounded-full bg-coffee-800 py-4 text-xs font-bold uppercase tracking-widest text-white shadow-lg transition-all hover:bg-coffee-900 active:scale-98"
						>
							<span>PESAN PAKET {selectedTier.name.toUpperCase()} ({selectedTier.price})</span>
							<ArrowRight size={15} />
						</a>

						<p class="text-center text-xs text-coffee-500">
							Ada pertanyaan khusus?{' '}
							<a
								href="https://wa.me/6281288889999?text=Halo%20Ketuk.id,%20saya%20ingin%20tanya%20tentang%20desain%20{template.title}"
								target="_blank"
								rel="noopener noreferrer"
								class="text-coffee-800 font-semibold underline hover:text-coffee-900"
							>
								Konsultasi langsung via WhatsApp
							</a>
						</p>
					</div>
				</div>
			</div>
		</div>

		<!-- What's Included / Feature Details -->
		<div class="mt-20 pt-12 border-t border-cream-200/80">
			<div class="text-center max-w-xl mx-auto mb-10">
				<span class="text-[10px] font-semibold tracking-[0.25em] uppercase text-coffee-500">
					KELENGKAPAN FITUR
				</span>
				<h2 class="font-serif mt-2 text-3xl sm:text-4xl font-medium text-coffee-950">
					Segala yang Anda Dapatkan
				</h2>
			</div>

			<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
				{#each template.features as feature, idx (feature)}
					<div class="rounded-2xl border border-cream-200 bg-white p-5 shadow-2xs">
						<div class="flex h-8 w-8 items-center justify-center rounded-full bg-coffee-50 text-coffee-800 font-serif text-xs font-semibold mb-3">
							0{idx + 1}
						</div>
						<p class="text-xs text-coffee-800 leading-relaxed font-medium">
							{feature}
						</p>
					</div>
				{/each}
			</div>
		</div>

		<!-- Timeline: 3 Easy Steps -->
		<div class="mt-20 rounded-3xl border border-cream-200 bg-cream-100/80 p-8 sm:p-12 text-center">
			<span class="text-[10px] font-semibold tracking-[0.25em] uppercase text-coffee-500">
				ALUR PENGERJAAN
			</span>
			<h2 class="font-serif mt-2 text-3xl font-medium text-coffee-950">
				Proses Mudah dan Terarah
			</h2>

			<div class="mt-10 grid gap-6 md:grid-cols-3">
				<div class="flex flex-col items-center p-4">
					<span class="h-10 w-10 rounded-full bg-coffee-800 text-white font-serif flex items-center justify-center text-sm font-semibold shadow-xs">
						1
					</span>
					<h3 class="font-serif mt-3 text-lg font-semibold text-coffee-900">Pilih Desain</h3>
					<p class="text-xs text-coffee-600 mt-1">
						Pilih paket {template.title} dan mulai dalam 1 menit.
					</p>
				</div>
				<div class="flex flex-col items-center p-4">
					<span class="h-10 w-10 rounded-full bg-coffee-800 text-white font-serif flex items-center justify-center text-sm font-semibold shadow-xs">
						2
					</span>
					<h3 class="font-serif mt-3 text-lg font-semibold text-coffee-900">Kustomisasi Isi</h3>
					<p class="text-xs text-coffee-600 mt-1">
						Lengkapi tanggal, denah peta, musik, dan nomor rekening kado.
					</p>
				</div>
				<div class="flex flex-col items-center p-4">
					<span class="h-10 w-10 rounded-full bg-coffee-800 text-white font-serif flex items-center justify-center text-sm font-semibold shadow-xs">
						3
					</span>
					<h3 class="font-serif mt-3 text-lg font-semibold text-coffee-900">Sebar ke Tamu</h3>
					<p class="text-xs text-coffee-600 mt-1">
						Undangan siap dibagikan ke WhatsApp keluarga dan sahabat.
					</p>
				</div>
			</div>
		</div>

		<!-- FAQ Accordion -->
		<div class="mt-20 max-w-3xl mx-auto">
			<div class="text-center mb-10">
				<span class="text-[10px] font-semibold tracking-[0.25em] uppercase text-coffee-500">
					PERTANYAAN UMUM
				</span>
				<h2 class="font-serif mt-2 text-3xl font-medium text-coffee-950">
					Pertanyaan Seputar Desain Ini
				</h2>
			</div>

			<div class="space-y-3">
				{#each faqs as faq, i (faq.q)}
					<div
						class="overflow-hidden rounded-2xl border border-cream-200 bg-white transition-all {openFaq ===
						i
							? 'border-champagne-400 shadow-xs'
							: 'hover:border-cream-300'}"
					>
						<button
							type="button"
							class="flex w-full items-center justify-between p-5 text-left text-sm font-serif sm:text-base font-semibold text-coffee-900"
							onclick={() => (openFaq = openFaq === i ? null : i)}
							aria-expanded={openFaq === i}
						>
							<span>{faq.q}</span>
							<ChevronDown
								size={16}
								class="shrink-0 text-coffee-500 transition-transform {openFaq === i
									? 'rotate-180 text-coffee-800'
									: ''}"
							/>
						</button>
						{#if openFaq === i}
							<div class="px-5 pb-5 pt-0 text-xs sm:text-sm text-coffee-600 leading-relaxed border-t border-cream-100/70 mt-1">
								{faq.a}
							</div>
						{/if}
					</div>
				{/each}
			</div>
		</div>

		<!-- Other Designs You Might Love -->
		<div class="mt-24 pt-12 border-t border-cream-200/80">
			<div class="flex items-center justify-between mb-8">
				<div>
					<span class="text-[10px] font-semibold tracking-[0.25em] uppercase text-coffee-500">
						EKSPLORASI LAINNYA
					</span>
					<h2 class="font-serif text-2xl sm:text-3xl font-medium text-coffee-950">
						Desain Lain yang Mungkin Anda Sukai
					</h2>
				</div>
				<a
					href="/template"
					class="hidden sm:inline-flex items-center gap-1 text-xs font-semibold text-coffee-800 hover:text-coffee-900 transition-colors"
				>
					<span>Lihat Semua</span>
					<ArrowRight size={13} />
				</a>
			</div>

			<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
				{#each otherTemplates as other (other.id)}
					<TemplateCard template={other} />
				{/each}
			</div>
		</div>
	</div>
</div>

<!-- Mobile Sticky Bottom Action Bar -->
<div
	class="fixed bottom-0 left-0 right-0 z-40 border-t border-cream-200 bg-cream-50/95 p-3.5 shadow-2xl backdrop-blur-md sm:hidden"
>
	<div class="flex items-center justify-between gap-3">
		<div>
			<span class="block text-[10px] uppercase tracking-wider text-coffee-500">Paket {selectedTier.name}</span>
			<span class="font-serif text-lg font-bold text-coffee-900">{selectedTier.price}</span>
		</div>
		<a
			href="/daftar?template={template.slug}&tier={selectedTier.id}"
			class="flex-1 text-center rounded-full bg-coffee-800 py-3 text-xs font-bold uppercase tracking-wider text-white shadow-md active:bg-coffee-900"
		>
			Pesan Sekarang
		</a>
	</div>
</div>

<!-- Fullscreen Lightbox Modal (Triggered by clicking EXPAND, matching Screenshot 2) -->
<EnvelopeModal open={isExpandedOpen} {template} onclose={() => (isExpandedOpen = false)} />
