<script lang="ts">
	import { X, Check, Music, MapPin, Calendar } from '@lucide/svelte';
	import type { TemplateItem } from './TemplateCard.svelte';

	interface Props {
		open: boolean;
		template: TemplateItem | null;
		onclose: () => void;
	}

	let { open, template, onclose }: Props = $props();

	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) {
			onclose();
		}
	}

	function handleKeyDown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			onclose();
		}
	}
</script>

<svelte:window onkeydown={handleKeyDown} />

{#if open && template}
	<!-- Modal Backdrop -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-coffee-950/70 p-3 sm:p-6 backdrop-blur-sm transition-opacity"
		onclick={handleBackdropClick}
		onkeydown={(e) => {
			if (e.key === 'Enter' || e.key === ' ') {
				onclose();
			}
		}}
		role="dialog"
		aria-modal="true"
		aria-label="Preview Template {template.title}"
		tabindex="-1"
	>
		<!-- Modal Box -->
		<div
			class="relative flex max-h-[92vh] w-full max-w-4xl flex-col overflow-hidden rounded-2xl sm:rounded-3xl border border-cream-200 bg-cream-50 shadow-2xl md:flex-row"
		>
			<!-- Close Button -->
			<button
				type="button"
				class="absolute top-4 right-4 z-30 flex h-9 w-9 items-center justify-center rounded-full bg-cream-200/80 text-coffee-800 transition-colors hover:bg-coffee-800 hover:text-white"
				onclick={onclose}
				aria-label="Tutup preview"
			>
				<X size={18} />
			</button>

			<!-- Left Side: Realistic Interactive Phone Simulation -->
			<div
				class="flex flex-1 items-center justify-center bg-gradient-to-b from-cream-150 to-cream-200/70 p-6 sm:p-8"
			>
				<!-- Simulated Phone Container -->
				<div
					class="relative h-[440px] sm:h-[490px] w-[240px] sm:w-[260px] overflow-hidden rounded-[36px] border-[6px] {template.palette.phoneBorder} bg-white shadow-2xl"
				>
					<!-- Dynamic Island -->
					<div class="absolute top-1.5 left-1/2 -translate-x-1/2 h-3.5 w-18 rounded-full bg-black z-20"></div>

					<!-- Screen View -->
					<div
						class="flex h-full w-full flex-col justify-between overflow-y-auto p-4 pt-8 text-center scrollbar-none {template.palette.bgGradient}"
					>
						<div class="space-y-2">
							<span
								class="text-[9px] font-semibold tracking-widest uppercase"
								style="color: {template.palette.accentTextColor};"
							>
								Wedding Celebration
							</span>
							<h3
								class="font-serif text-2xl font-medium leading-tight italic"
								style="color: {template.palette.accentColor};"
							>
								{template.preview.coupleName}
							</h3>
							<div class="mx-auto h-0.5 w-8 rounded-full bg-champagne-400"></div>
							<p class="text-[10px] text-coffee-700 flex items-center justify-center gap-1">
								<Calendar size={11} class="text-champagne-600" />
								{template.preview.date}
							</p>
						</div>

						<!-- Mock Card Detail -->
						<div class="rounded-xl border border-white/60 bg-white/80 p-3 shadow-xs space-y-2">
							<p class="font-serif text-xs italic text-coffee-800">
								"Bersama dalam cinta dan harapan abadi"
							</p>
							<div class="text-[10px] text-coffee-600">
								<MapPin size={11} class="inline text-coffee-700 mr-0.5" />
								Hotel Mulia Senayan, Jakarta
							</div>
							<div
								class="rounded-lg py-1.5 text-[10px] font-semibold text-white shadow-xs"
								style="background-color: {template.palette.accentColor};"
							>
								Konfirmasi Kehadiran (RSVP)
							</div>
						</div>

						<!-- Mini Audio Simulation -->
						<div class="flex items-center justify-between rounded-full bg-white/80 px-3 py-1 text-[10px] text-coffee-700 border border-white/60">
							<div class="flex items-center gap-1.5">
								<Music size={12} class="animate-pulse text-coffee-700" />
								<span class="truncate max-w-[120px]">Can't Help Falling in Love</span>
							</div>
							<span class="text-[9px] font-medium text-emerald-600">Playing</span>
						</div>
					</div>
				</div>
			</div>

			<!-- Right Side: Details & Action -->
			<div class="flex flex-1 flex-col justify-between p-6 sm:p-8 overflow-y-auto">
				<div>
					<div class="flex items-center gap-2">
						{#if template.badge}
							<span
								class="inline-block rounded-md bg-coffee-900 px-2 py-0.5 text-[9px] font-semibold tracking-[0.15em] text-champagne-200 uppercase"
							>
								{template.badge}
							</span>
						{/if}
						<span class="text-xs text-coffee-500 font-medium">Koleksi {template.category}</span>
					</div>

					<h3 class="font-serif mt-2 text-2xl sm:text-3xl font-semibold text-coffee-900">
						{template.title}
					</h3>
					<p class="mt-2 text-sm text-coffee-600 leading-relaxed">
						{template.subtitle}
					</p>

					<div class="mt-5 rounded-2xl border border-cream-200 bg-white p-4">
						<span class="block text-xs uppercase tracking-wider text-coffee-500 font-medium">Harga Promo</span>
						<div class="flex items-baseline gap-2 mt-0.5">
							<span class="font-serif text-2xl font-bold text-coffee-800">{template.price}</span>
							<span class="text-xs text-coffee-400 line-through">Rp 299.000</span>
							<span class="text-xs font-semibold text-emerald-600">Hemat 50%</span>
						</div>
					</div>

					<!-- Features list -->
					<div class="mt-5 space-y-2">
						<h4 class="text-xs font-semibold uppercase tracking-wider text-coffee-700">Fitur Termasuk:</h4>
						<ul class="space-y-1.5 text-xs text-coffee-700">
							<li class="flex items-center gap-2">
								<Check size={14} class="text-emerald-600 shrink-0" />
								<span>RSVP instan terhubung otomatis ke WhatsApp</span>
							</li>
							<li class="flex items-center gap-2">
								<Check size={14} class="text-emerald-600 shrink-0" />
								<span>Petunjuk arah & navigasi terintegrasi Google Maps</span>
							</li>
							<li class="flex items-center gap-2">
								<Check size={14} class="text-emerald-600 shrink-0" />
								<span>Amplop digital (BCA, Mandiri, BRI, QRIS & Kado)</span>
							</li>
							<li class="flex items-center gap-2">
								<Check size={14} class="text-emerald-600 shrink-0" />
								<span>Galeri hingga 15 foto & background musik pilihan</span>
							</li>
							<li class="flex items-center gap-2">
								<Check size={14} class="text-emerald-600 shrink-0" />
								<span>Aktif selamanya tanpa biaya langganan berkala</span>
							</li>
						</ul>
					</div>
				</div>

				<!-- Actions -->
				<div class="mt-6 pt-4 border-t border-cream-200/80 flex flex-col sm:flex-row gap-2.5">
					<a
						href="/daftar"
						class="flex-1 rounded-xl bg-coffee-800 px-5 py-3 text-center text-sm font-semibold text-white shadow-md transition-all hover:bg-coffee-900 active:scale-98"
					>
						Gunakan Desain Ini
					</a>
					<button
						type="button"
						class="rounded-xl border border-cream-300 bg-white px-4 py-3 text-xs font-medium text-coffee-700 hover:bg-cream-100 transition-colors"
						onclick={onclose}
					>
						Tutup
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}
