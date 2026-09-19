<script lang="ts">
	import {
		TemplateCard,
		TemplatePreviewModal,
		type TemplateItem,
	} from '$lib/components/domain';
	import { ArrowRight } from '@lucide/svelte';

	let activeCategory = $state('Semua');
	let selectedTemplate = $state<TemplateItem | null>(null);
	let isPreviewOpen = $state(false);

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
</script>

<svelte:head>
	<title>Koleksi Desain Undangan Digital | Ketuk.id</title>
	<meta
		name="description"
		content="Jelajahi seluruh koleksi template undangan pernikahan digital eksklusif Ketuk.id dengan estetika haute-couture."
	/>
</svelte:head>

<section class="py-12 sm:py-20 bg-cream-50">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<!-- Header -->
		<div class="mx-auto max-w-2xl text-center">
			<span class="inline-flex items-center gap-2 rounded-full border border-champagne-300 bg-cream-100 px-4 py-1 text-[10px] sm:text-[11px] font-semibold tracking-[0.2em] text-coffee-700 uppercase">
				<span class="h-1.5 w-1.5 rounded-full bg-champagne-600"></span>
				Koleksi Eksklusif 2026
			</span>
			<h1 class="font-serif mt-4 text-3xl sm:text-5xl font-medium tracking-tight text-coffee-950">
				Template untuk Setiap Cerita Cinta
			</h1>
			<p class="mt-3 text-sm sm:text-base text-coffee-600 leading-relaxed">
				Setiap tema dirancang dengan presisi grafis tinggi, tipografi mewah, dan integrasi RSVP WhatsApp serta musik romantis.
			</p>
		</div>

		<!-- Filter Pills -->
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

		<!-- Grid of Cards -->
		<div
			class="mt-10 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 sm:gap-6 lg:gap-7"
		>
			{#each filteredTemplates as template (template.id)}
				<TemplateCard {template} onPreview={handlePreview} />
			{/each}
		</div>

		<!-- Bottom CTA -->
		<div class="mt-16 text-center rounded-3xl border border-cream-200 bg-cream-100/70 p-8 sm:p-12 max-w-3xl mx-auto">
			<h3 class="font-serif text-2xl sm:text-3xl font-medium text-coffee-950">
				Sudah menemukan desain impian Anda?
			</h3>
			<p class="mt-2 text-xs sm:text-sm text-coffee-600">
				Mulai buat undangan Anda sekarang secara instan dalam 10 menit.
			</p>
			<div class="mt-6 flex flex-col sm:flex-row items-center justify-center gap-3">
				<a
					href="/daftar"
					class="w-full sm:w-auto inline-flex items-center justify-center gap-2 rounded-full bg-coffee-800 px-8 py-3.5 text-xs font-semibold uppercase tracking-wider text-white shadow-md hover:bg-coffee-900 transition-all"
				>
					<span>Mulai Gratis Sekarang</span>
					<ArrowRight size={14} />
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Interactive Preview Modal -->
<TemplatePreviewModal open={isPreviewOpen} template={selectedTemplate} onclose={closePreview} />
