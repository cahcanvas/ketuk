<script lang="ts">
	import { ArrowRight, Eye } from '@lucide/svelte';

	export interface TemplateItem {
		id: string;
		title: string;
		subtitle: string;
		category: string;
		badge?: string;
		price: string;
		slug: string;
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
		};
	}

	interface Props {
		template: TemplateItem;
		onPreview?: (template: TemplateItem) => void;
	}

	let { template, onPreview }: Props = $props();
</script>

<div
	class="group relative flex flex-col overflow-hidden rounded-2xl sm:rounded-3xl border border-cream-200/80 bg-cream-100/70 p-4 sm:p-5 transition-all duration-300 hover:-translate-y-1.5 hover:border-champagne-400/50 hover:bg-cream-100 hover:shadow-[0_20px_40px_-15px_rgba(92,20,33,0.12)]"
>
	<!-- Card Preview Stage (Phone Mockup + Envelope / Stationery Card) -->
	<div
		class="relative flex h-[240px] sm:h-[270px] w-full items-center justify-center overflow-hidden rounded-xl sm:rounded-2xl bg-gradient-to-b from-cream-150/90 to-cream-200/60 p-4"
	>
		<!-- Subtle ambient background shadow & pattern -->
		<div
			class="absolute inset-0 bg-[radial-gradient(circle_at_50%_30%,rgba(255,255,255,0.7),transparent_70%)]"
		></div>

		<!-- Category/Popular Badge -->
		{#if template.badge}
			<div class="absolute top-3 left-3 z-20">
				<span
					class="inline-block rounded-md bg-coffee-900/90 px-2.5 py-1 text-[9px] font-semibold tracking-[0.15em] text-champagne-200 uppercase shadow-xs backdrop-blur-xs border border-champagne-400/20"
				>
					{template.badge}
				</span>
			</div>
		{/if}

		<!-- Interactive Preview Hover Overlay -->
		<div
			class="absolute inset-0 z-20 flex items-center justify-center bg-coffee-950/25 opacity-0 backdrop-blur-[2px] transition-opacity duration-300 group-hover:opacity-100"
		>
			<button
				type="button"
				class="inline-flex items-center gap-1.5 rounded-full bg-white px-4 py-2 text-xs font-semibold text-coffee-900 shadow-md transition-transform duration-200 hover:scale-105 active:scale-95"
				onclick={() => onPreview?.(template)}
			>
				<Eye size={14} class="text-coffee-700" />
				<span>Coba Preview</span>
			</button>
		</div>

		<!-- Combined Stationery Suite Mockup -->
		<div class="relative flex h-full w-full items-center justify-center">
			<!-- Layer 1: Folded Stationery / Envelope in background (slightly angled) -->
			<div
				class="absolute right-4 sm:right-6 top-6 sm:top-5 h-[190px] sm:h-[210px] w-[110px] sm:w-[130px] -rotate-3 rounded-lg border border-cream-300/80 bg-white p-3 shadow-md transition-transform duration-300 group-hover:rotate-0 group-hover:scale-102"
				style="background: linear-gradient(135deg, #ffffff 0%, #faf8f5 100%);"
			>
				<!-- Subtle stationery border & botanical monogram -->
				<div class="flex h-full w-full flex-col justify-between border border-champagne-200/70 p-2 text-center">
					<div class="mx-auto flex h-6 w-6 items-center justify-center rounded-full border border-champagne-300">
						<span class="font-serif text-[10px] font-semibold text-champagne-700">
							{template.preview.coupleName.charAt(0)}
						</span>
					</div>

					<!-- Wax Seal imitation -->
					<div class="mx-auto flex items-center justify-center">
						<div
							class="h-7 w-7 rounded-full shadow-xs flex items-center justify-center text-white/90 font-serif text-[10px]"
							style="background-color: {template.palette.sealColor};"
						>
							♥
						</div>
					</div>

					<div class="space-y-0.5">
						<div class="h-1 w-8 mx-auto rounded-full bg-cream-300"></div>
						<div class="h-1 w-12 mx-auto rounded-full bg-cream-200"></div>
					</div>
				</div>
			</div>

			<!-- Layer 2: Smartphone Mockup in foreground (standing upright) -->
			<div
				class="relative z-10 left-[-15px] sm:left-[-20px] h-[200px] sm:h-[225px] w-[100px] sm:w-[115px] rotate-2 rounded-[22px] border-[4px] {template.palette.phoneBorder} bg-white shadow-xl transition-transform duration-300 group-hover:rotate-0 group-hover:scale-105"
			>
				<!-- Dynamic Island -->
				<div class="absolute top-1 left-1/2 -translate-x-1/2 h-1.5 w-7 rounded-full bg-black"></div>

				<!-- Screen Graphic -->
				<div
					class="flex h-full w-full flex-col items-center justify-between rounded-[18px] p-2.5 pt-4 text-center {template.palette.bgGradient}"
				>
					<span
						class="text-[7px] font-semibold tracking-widest uppercase"
						style="color: {template.palette.accentTextColor};"
					>
						Wedding
					</span>

					<div class="space-y-1">
						<!-- Couple Name -->
						<h5
							class="font-serif text-sm font-semibold leading-tight italic"
							style="color: {template.palette.accentColor};"
						>
							{template.preview.coupleName}
						</h5>
						<!-- Decorative divider -->
						<div class="mx-auto h-0.5 w-6 rounded-full bg-champagne-400/80"></div>
						<p class="text-[7px] text-coffee-600">{template.preview.date}</p>
					</div>

					<!-- Bottom RSVP pill -->
					<div
						class="w-full rounded-full py-0.5 text-[7px] font-medium text-white shadow-2xs"
						style="background-color: {template.palette.accentColor};"
					>
						Buka Undangan
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Card Information -->
	<div class="mt-4 flex flex-1 flex-col justify-between">
		<div>
			<div class="flex items-center justify-between">
				<a href="/template/{template.slug}" class="focus:outline-none">
					<h4 class="font-serif text-xl sm:text-2xl font-semibold text-coffee-950 transition-colors group-hover:text-coffee-800">
						{template.title}
					</h4>
				</a>
			</div>
			<p class="mt-1 text-xs text-coffee-600 line-clamp-2 leading-relaxed">
				{template.subtitle}
			</p>
		</div>

		<div class="mt-4 pt-3 border-t border-cream-200/60 flex items-center justify-between">
			<div>
				<span class="block text-[10px] uppercase tracking-wider text-coffee-500 font-medium">Harga</span>
				<span class="font-serif text-base sm:text-lg font-semibold text-coffee-900">
					{template.price}
				</span>
			</div>

			<a
				href="/template/{template.slug}"
				class="inline-flex items-center gap-1 text-xs font-semibold text-coffee-800 transition-all group-hover:translate-x-0.5 group-hover:text-coffee-900"
			>
				<span>Lihat Desain</span>
				<ArrowRight size={13} />
			</a>
		</div>
	</div>
</div>
