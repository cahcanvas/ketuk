<script lang="ts">
	import { Music, MapPin, Calendar, Check, Gift } from '@lucide/svelte';

	interface Props {
		coupleNames?: string;
		eventDate?: string;
		locationName?: string;
	}

	let {
		coupleNames = 'Sarah & Dimas',
		eventDate = 'Sabtu, 24 Oktober 2026',
		locationName = 'Plataran Dharmawangsa, Jakarta',
	}: Props = $props();

	let musicPlaying = $state(true);
	let activeTab = $state<'invitation' | 'details' | 'gift'>('invitation');
	let rsvpStatus = $state<'idle' | 'attending' | 'not_attending'>('idle');
</script>

<!-- iPhone Outer Frame -->
<div class="relative mx-auto w-full max-w-[290px] sm:max-w-[320px] md:max-w-[340px] select-none">
	<!-- Ambient Glow behind the phone -->
	<div
		class="absolute -inset-4 -z-10 rounded-[54px] bg-gradient-to-tr from-coffee-600/20 via-champagne-400/20 to-coffee-900/10 blur-2xl"
	></div>

	<!-- Phone Shell -->
	<div
		class="relative overflow-hidden rounded-[48px] border-[10px] border-coffee-900 bg-cream-50 shadow-[0_25px_60px_-15px_rgba(70,13,23,0.35),0_10px_20px_-5px_rgba(0,0,0,0.15)] ring-1 ring-champagne-400/40"
	>
		<!-- Dynamic Island & Camera cutout -->
		<div class="absolute top-2.5 left-1/2 z-30 -translate-x-1/2">
			<div class="flex h-5 w-24 items-center justify-between rounded-full bg-black px-2 py-0.5">
				<span class="h-2 w-2 rounded-full bg-zinc-900"></span>
				<span class="h-2.5 w-2.5 rounded-full bg-zinc-800 ring-1 ring-zinc-700"></span>
			</div>
		</div>

		<!-- Status Bar -->
		<div
			class="relative z-20 flex h-10 items-center justify-between px-6 pt-1 text-[11px] font-semibold text-coffee-700"
		>
			<span>09:41</span>
			<div class="flex items-center gap-1.5 text-[10px]">
				<span class="inline-block h-2 w-2 rounded-full bg-emerald-500"></span>
				<span>5G</span>
				<div class="h-2.5 w-5 rounded-sm border border-coffee-700 p-0.5">
					<div class="h-full w-full rounded-2xs bg-coffee-700"></div>
				</div>
			</div>
		</div>

		<!-- Phone Screen Content with Custom Scroll -->
		<div class="relative h-[530px] sm:h-[580px] overflow-y-auto bg-cream-50 scrollbar-none">
			<!-- Invitation Header / Hero inside Phone -->
			<div class="relative px-5 pt-4 pb-6 text-center">
				<!-- Botanical Wreath / Motif -->
				<div class="mx-auto mb-2 flex items-center justify-center">
					<svg
						class="h-8 w-32 text-champagne-600 opacity-80"
						viewBox="0 0 120 30"
						fill="none"
						stroke="currentColor"
						stroke-width="1.2"
					>
						<path
							d="M10,15 C25,5 35,25 50,15 C60,8 60,8 70,15 C85,25 95,5 110,15"
							stroke-linecap="round"
						/>
						<circle cx="60" cy="15" r="2.5" fill="currentColor" />
						<path d="M45,12 C40,9 38,5 42,4 C45,6 48,10 45,12" fill="currentColor" opacity="0.6" />
						<path d="M75,12 C80,9 82,5 78,4 C75,6 72,10 75,12" fill="currentColor" opacity="0.6" />
					</svg>
				</div>

				<span
					class="font-display text-[10px] font-semibold tracking-[0.25em] text-coffee-700 uppercase"
				>
					The Wedding Celebration
				</span>

				<!-- Couple Names in Editorial Serif -->
				<h3 class="font-serif mt-2 text-3xl font-medium tracking-wide text-coffee-900 italic">
					{coupleNames}
				</h3>

				<p class="font-serif mt-1 text-xs text-coffee-600 italic">
					Kami mengundang Anda untuk merayakan cinta kami
				</p>

				<!-- Main Arch Card / Illustration -->
				<div
					class="relative mx-auto mt-4 overflow-hidden rounded-t-[70px] rounded-b-2xl border border-champagne-300/60 bg-gradient-to-b from-cream-100 to-cream-50 p-4 shadow-sm"
				>
					<div
						class="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-champagne-200/30 via-transparent to-transparent"
					></div>

					<!-- Couple Silhouette or Monogram -->
					<div
						class="mx-auto flex h-20 w-20 items-center justify-center rounded-full border border-champagne-400/50 bg-cream-50/80 shadow-inner"
					>
						<span class="font-serif text-2xl font-normal text-coffee-800">S & D</span>
					</div>

					<div class="mt-4 space-y-2">
						<div
							class="inline-flex items-center gap-1.5 rounded-full bg-coffee-50 px-3 py-1 text-[11px] font-medium text-coffee-800 border border-coffee-100"
						>
							<Calendar size={12} class="text-coffee-600" />
							<span>{eventDate}</span>
						</div>

						<p class="text-[11px] text-coffee-600 flex items-center justify-center gap-1">
							<MapPin size={11} class="text-champagne-600 shrink-0" />
							<span class="truncate">{locationName}</span>
						</p>
					</div>

					<!-- Countdown Badge -->
					<div
						class="mt-4 grid grid-cols-4 gap-1.5 rounded-xl border border-cream-200 bg-white/70 p-2 text-center text-coffee-800"
					>
						<div>
							<span class="font-serif block text-sm font-semibold text-coffee-900">28</span>
							<span class="text-[9px] uppercase tracking-wider text-coffee-500">Hari</span>
						</div>
						<div>
							<span class="font-serif block text-sm font-semibold text-coffee-900">14</span>
							<span class="text-[9px] uppercase tracking-wider text-coffee-500">Jam</span>
						</div>
						<div>
							<span class="font-serif block text-sm font-semibold text-coffee-900">32</span>
							<span class="text-[9px] uppercase tracking-wider text-coffee-500">Mnt</span>
						</div>
						<div>
							<span class="font-serif block text-sm font-semibold text-coffee-900">45</span>
							<span class="text-[9px] uppercase tracking-wider text-coffee-500">Dtk</span>
						</div>
					</div>
				</div>

				<!-- Tabs for Mini-Interactive Experience -->
				<div class="mt-4 flex rounded-full border border-cream-200 bg-cream-100/80 p-1">
					<button
						type="button"
						class="flex-1 rounded-full py-1 text-[11px] font-medium transition-all {activeTab ===
						'invitation'
							? 'bg-coffee-800 text-white shadow-xs'
							: 'text-coffee-700 hover:text-coffee-900'}"
						onclick={() => (activeTab = 'invitation')}
					>
						Akad & Resepsi
					</button>
					<button
						type="button"
						class="flex-1 rounded-full py-1 text-[11px] font-medium transition-all {activeTab ===
						'details'
							? 'bg-coffee-800 text-white shadow-xs'
							: 'text-coffee-700 hover:text-coffee-900'}"
						onclick={() => (activeTab = 'details')}
					>
						RSVP Cepat
					</button>
					<button
						type="button"
						class="flex-1 rounded-full py-1 text-[11px] font-medium transition-all {activeTab ===
						'gift'
							? 'bg-coffee-800 text-white shadow-xs'
							: 'text-coffee-700 hover:text-coffee-900'}"
						onclick={() => (activeTab = 'gift')}
					>
						Amplop Digital
					</button>
				</div>

				<!-- Tab Content -->
				<div class="mt-3 text-left">
					{#if activeTab === 'invitation'}
						<div class="rounded-xl border border-cream-200 bg-white p-3.5 shadow-xs space-y-2">
							<div class="flex items-start justify-between">
								<div>
									<h4 class="font-serif text-sm font-semibold text-coffee-900">Akad Nikah</h4>
									<p class="text-[10px] text-coffee-500">Pukul 08.00 - 10.00 WIB</p>
								</div>
								<span
									class="rounded-full bg-cream-100 px-2 py-0.5 text-[9px] font-medium text-coffee-700"
								>
									Khidmat
								</span>
							</div>
							<div class="border-t border-cream-100 pt-2 flex items-start justify-between">
								<div>
									<h4 class="font-serif text-sm font-semibold text-coffee-900">Resepsi Malam</h4>
									<p class="text-[10px] text-coffee-500">Pukul 19.00 - 21.30 WIB</p>
								</div>
								<span
									class="rounded-full bg-coffee-50 px-2 py-0.5 text-[9px] font-medium text-coffee-800"
								>
									Gala Dinner
								</span>
							</div>
						</div>
					{:else if activeTab === 'details'}
						<div class="rounded-xl border border-cream-200 bg-white p-3.5 shadow-xs text-center">
							<p class="text-[11px] text-coffee-700">Apakah Anda bersedia hadir?</p>
							<div class="mt-2 flex gap-2">
								<button
									type="button"
									class="flex-1 rounded-lg border py-1.5 text-[10px] font-medium transition-all {rsvpStatus ===
									'attending'
										? 'border-emerald-500 bg-emerald-50 text-emerald-700'
										: 'border-cream-300 bg-cream-50 text-coffee-800 hover:bg-cream-100'}"
									onclick={() => (rsvpStatus = 'attending')}
								>
									{#if rsvpStatus === 'attending'}
										<Check size={11} class="inline mr-1" />
									{/if}
									Ya, Hadir
								</button>
								<button
									type="button"
									class="flex-1 rounded-lg border py-1.5 text-[10px] font-medium transition-all {rsvpStatus ===
									'not_attending'
										? 'border-coffee-500 bg-coffee-50 text-coffee-700'
										: 'border-cream-300 bg-cream-50 text-coffee-800 hover:bg-cream-100'}"
									onclick={() => (rsvpStatus = 'not_attending')}
								>
									Maaf, Belum Bisa
								</button>
							</div>
							{#if rsvpStatus !== 'idle'}
								<p class="mt-2 text-[10px] text-emerald-600 font-medium">
									Terima kasih, respon Anda tercatat!
								</p>
							{/if}
						</div>
					{:else}
						<div class="rounded-xl border border-cream-200 bg-white p-3.5 shadow-xs">
							<div class="flex items-center gap-2">
								<div class="rounded-lg bg-champagne-100 p-1.5 text-champagne-700">
									<Gift size={14} />
								</div>
								<div>
									<h4 class="font-serif text-xs font-semibold text-coffee-900">
										BCA Transfer Digital
									</h4>
									<p class="text-[10px] text-coffee-500 font-mono">8720 1928 331 (Sarah)</p>
								</div>
							</div>
							<p class="mt-2 text-[10px] text-coffee-600 italic">
								Bisa langsung salin nomor rekening atau kirim kado via alamat.
							</p>
						</div>
					{/if}
				</div>

				<div class="mt-4 pb-2">
					<span class="font-serif text-xs text-coffee-800 italic">
						"Dua jiwa, satu takdir indah."
					</span>
				</div>
			</div>

			<!-- Floating Audio Player Button inside Phone -->
			<button
				type="button"
				class="absolute bottom-4 right-4 z-30 flex h-10 w-10 items-center justify-center rounded-full bg-coffee-800 text-white shadow-lg transition-transform active:scale-95 hover:bg-coffee-900"
				onclick={() => (musicPlaying = !musicPlaying)}
				title={musicPlaying ? 'Jeda Musik' : 'Putar Musik'}
			>
				<Music size={16} class={musicPlaying ? 'animate-pulse' : 'opacity-70'} />
				{#if musicPlaying}
					<span
						class="absolute -top-1 -right-1 flex h-3 w-3 items-center justify-center rounded-full bg-champagne-400"
					>
						<span class="h-1.5 w-1.5 rounded-full bg-white"></span>
					</span>
				{/if}
			</button>
		</div>

		<!-- Home Indicator Bar -->
		<div class="absolute bottom-1.5 left-1/2 -translate-x-1/2">
			<div class="h-1 w-32 rounded-full bg-coffee-400/40"></div>
		</div>
	</div>
</div>
