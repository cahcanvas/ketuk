<script lang="ts">
	import { Maximize2, Volume2, VolumeX, MapPin, Calendar, Check, X } from '@lucide/svelte';
	import type { TemplateItem } from '$lib/data/templates';

	interface Props {
		template: TemplateItem;
		onExpand?: () => void;
		interactive?: boolean;
		isModal?: boolean;
	}

	let { template, onExpand, interactive = true, isModal = false }: Props = $props();

	let isOpen = $state(false);
	let musicPlaying = $state(false);
	let rsvpAttending = $state<boolean | null>(null);

	function toggleOpen() {
		if (!interactive) return;
		isOpen = !isOpen;
		if (isOpen) {
			musicPlaying = true;
		} else {
			musicPlaying = false;
		}
	}
</script>

<div class="flex flex-col items-center w-full">
	<!-- Envelope Outer Card Container -->
	<div
		class="relative aspect-[9/15.5] shrink-0 rounded-2xl sm:rounded-3xl shadow-[0_20px_50px_-15px_rgba(70,13,23,0.2),0_10px_25px_-5px_rgba(0,0,0,0.08)] border border-champagne-300/60 overflow-hidden select-none transition-all duration-300 {isModal
			? 'w-[min(100%,360px)] sm:w-[min(100%,420px)] md:w-[460px]'
			: 'w-[min(100%,320px)] sm:w-[min(100%,380px)] md:w-[400px]'}"
		style="background: linear-gradient(160deg, #FAF7F2 0%, #F5EFE6 50%, #EDE4D8 100%);"
	>
		<!-- Tactile Embossed Floral Texture / Pattern Overlay -->
		<div
			class="absolute inset-0 opacity-40 mix-blend-multiply pointer-events-none"
			style="background-image: radial-gradient(#d6bc98 0.75px, transparent 0.75px), radial-gradient(#d6bc98 0.75px, #FAF7F2 0.75px); background-size: 24px 24px; background-position: 0 0, 12px 12px;"
		></div>

		<!-- Subtle Embossed Corner Flourish -->
		<svg
			class="absolute top-4 left-4 h-16 w-16 text-champagne-400/40 pointer-events-none"
			viewBox="0 0 100 100"
			fill="none"
			stroke="currentColor"
			stroke-width="1"
		>
			<path d="M10,90 Q10,10 90,10 M25,90 Q25,25 90,25 M40,90 Q40,40 90,40" />
		</svg>
		<svg
			class="absolute bottom-4 right-4 h-16 w-16 text-champagne-400/40 pointer-events-none rotate-180"
			viewBox="0 0 100 100"
			fill="none"
			stroke="currentColor"
			stroke-width="1"
		>
			<path d="M10,90 Q10,10 90,10 M25,90 Q25,25 90,25 M40,90 Q40,40 90,40" />
		</svg>

		<!-- Envelope Flap (Triangle with soft realistic shadow) -->
		<div
			class="absolute top-0 left-0 right-0 h-[46%] transition-transform duration-700 origin-top z-10 {isOpen
				? '-rotate-x-180 opacity-0 pointer-events-none'
				: ''}"
			style="perspective: 1000px;"
		>
			<svg
				class="w-full h-full drop-shadow-[0_8px_12px_rgba(0,0,0,0.06)]"
				viewBox="0 0 400 240"
				preserveAspectRatio="none"
			>
				<path
					d="M0,0 L400,0 L200,225 Z"
					fill="#FAF7F2"
					stroke="#EADCC8"
					stroke-width="1.5"
				/>
				<!-- Inner flap contour -->
				<path
					d="M15,5 L385,5 L200,210 Z"
					fill="none"
					stroke="#DFC5A4"
					stroke-width="0.8"
					stroke-dasharray="3,3"
					opacity="0.6"
				/>
			</svg>
		</div>

		<!-- When Envelope is CLOSED: Show Wax Seal & Script Invitation Line -->
		{#if !isOpen}
			<div class="absolute inset-0 flex flex-col items-center justify-between p-8 z-20">
				<!-- Top branding watermark -->
				<div class="pt-3 text-center">
					<span
						class="text-[9px] font-semibold uppercase tracking-[0.3em] text-coffee-400"
					>
						The Wedding Collection
					</span>
				</div>

				<!-- Center: Realistic Wax Seal (Clickable to open) -->
				<div class="relative my-auto flex flex-col items-center">
					<button
						type="button"
						class="group relative flex items-center justify-center transition-transform duration-300 hover:scale-105 active:scale-95 focus:outline-none"
						onclick={toggleOpen}
						aria-label="Buka undangan"
					>
						<!-- Outer Wax Seal 3D Blob with realistic edge irregularities -->
						<div
							class="relative flex h-24 w-24 sm:h-28 sm:w-28 items-center justify-center rounded-full shadow-[0_12px_24px_-6px_rgba(40,20,10,0.35),inset_0_2px_4px_rgba(255,255,255,0.4),inset_0_-3px_6px_rgba(0,0,0,0.35)]"
							style="background-color: {template.palette.sealColor};"
						>
							<!-- Inner Stamp Ring -->
							<div
								class="flex h-18 w-18 sm:h-21 sm:w-21 items-center justify-center rounded-full border-2 border-black/15 shadow-[inset_0_2px_5px_rgba(0,0,0,0.3)] bg-black/5"
							>
								<!-- Monogram Initials -->
								<span
									class="font-serif text-2xl sm:text-3xl font-normal text-white/90 drop-shadow-[0_1px_2px_rgba(0,0,0,0.5)] italic select-none"
								>
									{template.preview.monogram}
								</span>
							</div>

							<!-- Subtle wax seal sheen -->
							<div
								class="absolute top-2 left-4 h-6 w-8 rounded-full bg-white/20 blur-xs rotate-[-30deg]"
							></div>
						</div>

						<!-- Pulse hint on hover -->
						<span
							class="absolute -bottom-8 whitespace-nowrap rounded-full bg-coffee-950/80 px-3 py-1 text-[10px] font-medium tracking-wider text-white opacity-0 transition-opacity group-hover:opacity-100"
						>
							Ketuk untuk membuka
						</span>
					</button>

					<!-- Elegant Script Typography on the Envelope -->
					<div class="mt-8 text-center px-4">
						<p
							class="font-serif text-xl sm:text-2xl italic tracking-wide text-coffee-800 leading-relaxed font-normal"
							style="color: {template.palette.accentColor};"
						>
							"{template.envelopeScript}"
						</p>
						<p class="font-serif mt-2 text-sm text-coffee-600 italic">
							{template.preview.coupleName}
						</p>
					</div>
				</div>

				<!-- Bottom subtle date -->
				<div class="pb-2 text-center">
					<span
						class="text-[10px] font-medium tracking-[0.2em] uppercase text-coffee-500"
					>
						{template.preview.date}
					</span>
				</div>
			</div>
		{:else}
			<!-- When Envelope is OPENED: Beautiful Invitation Letter Content -->
			<div
				class="absolute inset-0 z-20 flex flex-col justify-between overflow-y-auto p-6 sm:p-8 bg-gradient-to-b from-cream-50 via-white to-cream-50 text-center animate-in fade-in duration-500 scrollbar-none"
			>
				<!-- Top Audio & Close Button -->
				<div class="flex items-center justify-between pb-2 border-b border-cream-200">
					<button
						type="button"
						class="flex items-center gap-1.5 text-[11px] font-medium text-coffee-700 hover:text-coffee-800"
						onclick={() => (musicPlaying = !musicPlaying)}
					>
						{#if musicPlaying}
							<Volume2 size={14} class="text-coffee-700 animate-pulse" />
							<span>Musik: Romansa</span>
						{:else}
							<VolumeX size={14} class="text-coffee-400" />
							<span>Musik Dijeda</span>
						{/if}
					</button>

					<button
						type="button"
						class="inline-flex items-center gap-1 rounded-full bg-cream-200/80 px-2.5 py-1 text-[11px] font-semibold text-coffee-900 hover:bg-cream-300 transition-colors"
						onclick={toggleOpen}
					>
						<span>Tutup Surat</span>
						<X size={12} />
					</button>
				</div>

				<!-- Main Invitation Content -->
				<div class="py-4 space-y-4">
					<div class="space-y-1">
						<span
							class="text-[9px] font-semibold uppercase tracking-[0.25em] text-champagne-600"
						>
							Pernikahan Suci
						</span>
						<h3
							class="font-serif text-3xl sm:text-4xl font-medium italic text-coffee-950"
							style="color: {template.palette.accentColor};"
						>
							{template.preview.coupleName}
						</h3>
						<div class="mx-auto h-0.5 w-12 rounded-full bg-champagne-400"></div>
					</div>

					<p class="font-serif text-xs italic text-coffee-600 leading-relaxed px-2">
						"Dengan memohon rahmat dan ridho Tuhan Yang Maha Esa, kami bermaksud merayakan ikatan suci janji suci kami."
					</p>

					<!-- Event Detail Card -->
					<div
						class="rounded-xl border border-cream-200 bg-cream-50/70 p-4 text-left space-y-3"
					>
						<div>
							<div class="flex items-center gap-1.5 text-xs font-semibold text-coffee-900 font-serif">
								<Calendar size={13} class="text-coffee-700" />
								<span>Waktu & Tanggal</span>
							</div>
							<p class="text-xs text-coffee-700 mt-0.5">{template.preview.date}</p>
							<p class="text-[11px] text-coffee-500">Akad: 08.00 WIB | Resepsi: 19.00 WIB</p>
						</div>

						<div class="border-t border-cream-200 pt-2">
							<div class="flex items-center gap-1.5 text-xs font-semibold text-coffee-900 font-serif">
								<MapPin size={13} class="text-coffee-700" />
								<span>Lokasi Acara</span>
							</div>
							<p class="text-xs text-coffee-700 mt-0.5">{template.preview.location}</p>
						</div>
					</div>

					<!-- Quick RSVP Interaction -->
					<div class="rounded-xl border border-cream-200 bg-white p-3.5 shadow-2xs">
						<p class="text-xs font-medium text-coffee-800">Konfirmasi Kehadiran Tamu</p>
						<div class="mt-2.5 flex gap-2">
							<button
								type="button"
								class="flex-1 rounded-lg py-2 text-xs font-medium transition-all {rsvpAttending ===
								true
									? 'bg-emerald-600 text-white'
									: 'bg-cream-100 text-coffee-800 hover:bg-cream-200'}"
								onclick={() => (rsvpAttending = true)}
							>
								{#if rsvpAttending === true}
									<Check size={12} class="inline mr-1" />
								{/if}
								Hadir
							</button>
							<button
								type="button"
								class="flex-1 rounded-lg py-2 text-xs font-medium transition-all {rsvpAttending ===
								false
									? 'bg-coffee-800 text-white'
									: 'bg-cream-100 text-coffee-800 hover:bg-cream-200'}"
								onclick={() => (rsvpAttending = false)}
							>
								Maaf, Berhalangan
							</button>
						</div>
						{#if rsvpAttending !== null}
							<p class="mt-2 text-[10px] text-emerald-700 font-medium">
								Respon kehadiran Anda tersimpan otomatis via WhatsApp!
							</p>
						{/if}
					</div>
				</div>

				<div class="pt-2 text-center">
					<span class="font-serif text-xs italic text-coffee-500">
						Terima kasih atas doa & restu Anda
					</span>
				</div>
			</div>
		{/if}
	</div>

	<!-- Controls Bar below the Card (Interactive Demo & Expand button) -->
	<div
		class="mt-4 flex w-full items-center justify-between px-2 text-xs text-coffee-600 {isModal
			? 'max-w-[min(100%,360px)] sm:max-w-[min(100%,420px)] md:max-w-[460px]'
			: 'max-w-[min(100%,320px)] sm:max-w-[min(100%,380px)] md:max-w-[400px]'}"
	>
		<button
			type="button"
			class="flex items-center gap-1.5 text-[11px] font-medium uppercase tracking-wider text-coffee-500 hover:text-coffee-800 transition-colors"
			onclick={toggleOpen}
		>
			<span class="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
			<span>{isOpen ? 'SURAT TERBUKA • KETUK UTK TUTUP' : 'DEMO INTERAKTIF • KETUK SEGEL'}</span>
		</button>

		{#if onExpand}
			<button
				type="button"
				class="inline-flex items-center gap-1.5 rounded-full border border-cream-300 bg-white/90 px-3.5 py-1 text-[11px] font-semibold uppercase tracking-wider text-coffee-800 shadow-2xs hover:bg-white hover:border-coffee-900 transition-all active:scale-95"
				onclick={onExpand}
			>
				<Maximize2 size={11} class="text-coffee-800" />
				<span>EXPAND</span>
			</button>
		{/if}
	</div>
</div>
