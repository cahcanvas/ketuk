<script lang="ts">
	import {
		Sparkles,
		Layers,
		Palette,
		Eye,
		ArrowRight,
		Check,
		X,
		Search,
		Plus,
		Copy,
		Trash2,
		RefreshCw,
		Music,
		Volume2,
		VolumeX,
		Smartphone,
		Mail,
		Settings,
		ShieldCheck,
		Heart,
		Sliders,
		CheckCircle2,
		Wand2,
		Calendar,
		MapPin,
		Gift,
		AlertCircle,
		Play,
		Share2,
		Info,
		Zap,
	} from '@lucide/svelte';
	import {
		MASTER_TEMPLATES,
		DEFAULT_PLAN_PERMISSIONS,
		type MasterTemplateProduct,
		type ProductCategory,
		type ProductBadge,
		type OpeningEffectType,
		type ParticleEffectType,
		type ParticleDensity,
		type ScrollMotionType,
	} from '$lib/data/master-templates';

	// Local mutable copy of master templates for the admin studio session
	let templates = $state<MasterTemplateProduct[]>(JSON.parse(JSON.stringify(MASTER_TEMPLATES)));
	let selectedTemplateId = $state<string>(MASTER_TEMPLATES[0]?.id ?? 'maison-doree');

	// Active tab inside the studio editor
	let activeTab = $state<'visual' | 'animations' | 'permissions'>('visual');

	// Filter and search state
	let searchQuery = $state('');
	let selectedCategory = $state<string>('Semua');
	let selectedStatus = $state<string>('all');

	// Simulator state
	let previewMode = $state<'envelope' | 'phone'>('envelope');
	let isEnvelopeOpen = $state(false);
	let isPlayingAnimation = $state(false);
	let simulatedMusicPlaying = $state(false);
	let saveFeedback = $state<string | null>(null);

	// Derived currently edited template
	const currentTemplate = $derived(
		templates.find((t) => t.id === selectedTemplateId) ?? templates[0],
	);

	// Derived filtered list of templates
	const filteredTemplates = $derived(
		templates.filter((item) => {
			const matchesSearch =
				item.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
				item.slug.toLowerCase().includes(searchQuery.toLowerCase()) ||
				item.subtitle.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesCategory =
				selectedCategory === 'Semua' || item.category === selectedCategory;
			const matchesStatus =
				selectedStatus === 'all' || item.status === selectedStatus;
			return matchesSearch && matchesCategory && matchesStatus;
		}),
	);

	// Color palette quick presets for admin convenience
	const WAX_SEAL_PRESETS = [
		{ name: 'Sage Green', color: '#7a9a7a' },
		{ name: 'Bordeaux Wine', color: '#882235' },
		{ name: 'Royal Gold', color: '#a98350' },
		{ name: 'Olive Green', color: '#2e5c3e' },
		{ name: 'Deep Burgundy', color: '#5c1421' },
		{ name: 'Imperial Emerald', color: '#064e3b' },
		{ name: 'Slate Charcoal', color: '#334155' },
		{ name: 'Blush Rose', color: '#d9777f' },
	];

	const BG_GRADIENT_PRESETS = [
		{ name: 'Stone & Champagne Warm', value: 'bg-gradient-to-b from-stone-50 via-cream-50 to-amber-50/40' },
		{ name: 'Rose & Warm Cream', value: 'bg-gradient-to-b from-rose-50 via-cream-50 to-cream-100' },
		{ name: 'Botanical Emerald Soft', value: 'bg-gradient-to-b from-emerald-50/70 via-cream-50 to-cream-100' },
		{ name: 'Wine Velvet Elegance', value: 'bg-gradient-to-b from-wine-50 via-cream-50 to-wine-100/50' },
		{ name: 'Royal Starlight Gold', value: 'bg-gradient-to-b from-zinc-50 via-cream-50 to-emerald-50/30' },
		{ name: 'Minimalist Santorini Pure', value: 'bg-gradient-to-b from-sky-50/50 via-white to-cream-50' },
	];

	const AUDIO_PRESETS = [
		{ name: 'A Thousand Years (Acoustic Piano)', url: '/audio/presets/thousand-years-piano.mp3' },
		{ name: 'Canon in D (Strings & Harp)', url: '/audio/presets/canon-in-d.mp3' },
		{ name: 'Renaissance Romance (Orchestral)', url: '/audio/presets/renaissance-waltz.mp3' },
		{ name: 'Acoustic Morning Song (Folk Guitar)', url: '/audio/presets/acoustic-garden.mp3' },
		{ name: 'Royal Overture (Symphonic)', url: '/audio/presets/royal-symphony.mp3' },
		{ name: 'Serenade of Breeze (Lo-Fi Piano)', url: '/audio/presets/lofi-piano.mp3' },
	];

	// Actions
	function selectTemplate(id: string) {
		selectedTemplateId = id;
		isEnvelopeOpen = false;
		simulatedMusicPlaying = false;
	}

	function handleCreateNewTemplate() {
		const newId = `custom-template-${Date.now().toString().slice(-4)}`;
		const newTemplate: MasterTemplateProduct = {
			id: newId,
			title: 'Desain Baru Kustom',
			slug: newId,
			subtitle: 'Kemewahan visual kontemporer dengan palet elegan dan animasi modern.',
			tagline: 'Sentuhan personal terbaik untuk momen paling berharga dalam hidup Anda.',
			category: 'Elegan',
			badge: 'BARU',
			status: 'draft',
			basePrice: 159000,
			rating: 5.0,
			reviewCount: 0,
			palette: {
				bgGradient: 'bg-gradient-to-b from-rose-50 via-cream-50 to-cream-100',
				accentColor: '#5c1421',
				accentTextColor: '#882235',
				sealColor: '#a98350',
				phoneBorder: 'border-wine-900',
				fontFamily: 'Cormorant Garamond',
			},
			preview: {
				coupleName: 'Aurel & Revan',
				date: 'Minggu, 12 Desember 2026',
				style: 'Modern Couture',
				location: 'The Glass House, Jakarta',
				monogram: 'A & R',
			},
			envelope: {
				envelopeScript: 'Warmly invite you to celebrate our love',
				paperTexture: 'classic_embossed',
				sealSymbol: 'heart',
			},
			animations: {
				openingEffect: 'wax_seal_crack',
				particleEffect: 'falling_petals',
				particleDensity: 'medium',
				scrollMotion: 'parallax',
				openingDurationMs: 850,
			},
			audio: {
				presetTrackName: 'A Thousand Years (Acoustic Piano)',
				presetTrackUrl: '/audio/presets/thousand-years-piano.mp3',
				autoplay: true,
				showVisualizer: true,
				fadeDurationMs: 1500,
			},
			customizationMatrix: JSON.parse(JSON.stringify(DEFAULT_PLAN_PERMISSIONS)),
		};

		templates = [newTemplate, ...templates];
		selectedTemplateId = newId;
		showNotification('Template baru berhasil dibuat dalam mode Draft.');
	}

	function handleDuplicateTemplate() {
		if (!currentTemplate) return;
		const cloned: MasterTemplateProduct = JSON.parse(JSON.stringify(currentTemplate));
		const cloneId = `${currentTemplate.id}-copy-${Date.now().toString().slice(-3)}`;
		cloned.id = cloneId;
		cloned.title = `${currentTemplate.title} (Salinan)`;
		cloned.slug = `${currentTemplate.slug}-salinan`;
		cloned.status = 'draft';
		cloned.badge = 'BARU';

		templates = [cloned, ...templates];
		selectedTemplateId = cloneId;
		showNotification(`Berhasil menduplikasi "${currentTemplate.title}".`);
	}

	function handleTogglePublishStatus() {
		if (!currentTemplate) return;
		currentTemplate.status =
			currentTemplate.status === 'published' ? 'draft' : 'published';
		showNotification(
			`Status template diubah menjadi: ${currentTemplate.status.toUpperCase()}`,
		);
	}

	function handleSaveAll() {
		showNotification('Perubahan Master Produk & Matriks Kustomisasi berhasil disimpan ke sistem.');
	}

	function showNotification(msg: string) {
		saveFeedback = msg;
		setTimeout(() => {
			if (saveFeedback === msg) saveFeedback = null;
		}, 3500);
	}

	function triggerTestAnimation() {
		isPlayingAnimation = true;
		isEnvelopeOpen = !isEnvelopeOpen;
		if (isEnvelopeOpen && currentTemplate.audio.autoplay) {
			simulatedMusicPlaying = true;
		} else if (!isEnvelopeOpen) {
			simulatedMusicPlaying = false;
		}

		setTimeout(() => {
			isPlayingAnimation = false;
		}, currentTemplate.animations.openingDurationMs + 300);
	}
</script>

<svelte:head>
	<title>Master Pengelolaan Produk Undangan | Ketuk Admin</title>
</svelte:head>

<div class="space-y-6">
	<!-- Toast Feedback Notification -->
	{#if saveFeedback}
		<div class="fixed bottom-6 right-6 z-50 flex items-center gap-3 rounded-xl border border-emerald-500/40 bg-emerald-950/90 px-4 py-3 text-xs font-semibold text-emerald-300 shadow-2xl backdrop-blur-md animate-in slide-in-from-bottom-3 duration-300">
			<CheckCircle2 size={18} class="text-emerald-400 shrink-0" />
			<span>{saveFeedback}</span>
			<button
				type="button"
				onclick={() => (saveFeedback = null)}
				class="ml-2 text-emerald-400 hover:text-white"
			>
				<X size={14} />
			</button>
		</div>
	{/if}

	<!-- Header with Action Buttons -->
	<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between border-b border-navy-800/80 pb-6">
		<div>
			<div class="flex items-center gap-2">
				<h1 class="font-display text-2xl sm:text-3xl font-bold tracking-tight text-white">
					Master Produk & Template Undangan
				</h1>
				<span class="rounded-full bg-coral-500/20 px-2.5 py-0.5 text-xs font-semibold text-coral-400 border border-coral-500/30">
					Visual & Animation Studio
				</span>
			</div>
			<p class="mt-1 text-xs sm:text-sm text-navy-400 max-w-3xl">
				Konfigurasi desain visual, tipografi, efek animasi buka amplop, partikel, audio latar, serta matriks hak kustomisasi member untuk setiap paket (Gratis, Pro, Lengkap).
			</p>
		</div>

		<div class="flex flex-wrap items-center gap-2.5 shrink-0">
			<button
				type="button"
				onclick={handleCreateNewTemplate}
				class="inline-flex items-center gap-1.5 rounded-xl border border-navy-700 bg-navy-800/90 px-3.5 py-2 text-xs font-semibold text-white shadow-sm hover:bg-navy-700 hover:text-white transition-all active:scale-95"
			>
				<Plus size={15} class="text-coral-400" />
				<span>Tambah Template</span>
			</button>

			<button
				type="button"
				onclick={handleDuplicateTemplate}
				class="inline-flex items-center gap-1.5 rounded-xl border border-navy-700 bg-navy-800/90 px-3.5 py-2 text-xs font-semibold text-navy-200 hover:bg-navy-700 hover:text-white transition-all active:scale-95"
				title="Duplikat template yang sedang dipilih"
			>
				<Copy size={14} />
				<span>Duplikat</span>
			</button>

			<button
				type="button"
				onclick={handleSaveAll}
				class="inline-flex items-center gap-1.5 rounded-xl bg-coral-500 px-4 py-2 text-xs font-bold uppercase tracking-wider text-white shadow-lg shadow-coral-500/30 hover:bg-coral-600 transition-all hover:scale-102 active:scale-98"
			>
				<Check size={15} />
				<span>Simpan Perubahan</span>
			</button>
		</div>
	</div>

	<!-- Template Ribbon / Interactive Selector Bar -->
	<div class="space-y-3">
		<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
			<div class="flex items-center gap-2">
				<span class="text-xs font-bold uppercase tracking-wider text-navy-300">PILIH TEMPLATE MASTER:</span>
				<span class="text-xs text-navy-400">({filteredTemplates.length} dari {templates.length} desain)</span>
			</div>

			<!-- Search & Filter Controls -->
			<div class="flex flex-wrap items-center gap-2">
				<!-- Search -->
				<div class="relative">
					<Search size={14} class="absolute left-2.5 top-1/2 -translate-y-1/2 text-navy-400" />
					<input
						type="text"
						bind:value={searchQuery}
						placeholder="Cari nama / kategori..."
						class="w-48 sm:w-56 rounded-lg border border-navy-700 bg-navy-900/90 pl-8 pr-3 py-1 text-xs text-white placeholder-navy-500 focus:border-coral-500 focus:outline-none"
					/>
				</div>

				<!-- Category Filter -->
				<select
					bind:value={selectedCategory}
					class="rounded-lg border border-navy-700 bg-navy-900/90 px-2.5 py-1 text-xs text-navy-200 focus:border-coral-500 focus:outline-none"
				>
					<option value="Semua">Semua Kategori</option>
					<option value="Klasik">Klasik</option>
					<option value="Elegan">Elegan</option>
					<option value="Botanical">Botanical</option>
					<option value="Luxury Gold">Luxury Gold</option>
					<option value="Minimalis">Minimalis</option>
					<option value="Modern">Modern</option>
				</select>
			</div>
		</div>

		<!-- Scrollable Template Cards Ribbon -->
		<div class="flex gap-3 overflow-x-auto pb-2 scrollbar-thin scrollbar-thumb-navy-700">
			{#each filteredTemplates as tmpl (tmpl.id)}
				{@const isSelected = tmpl.id === selectedTemplateId}
				<button
					type="button"
					onclick={() => selectTemplate(tmpl.id)}
					class="group relative flex w-60 shrink-0 flex-col justify-between rounded-2xl border p-3.5 text-left transition-all duration-200 {isSelected
						? 'border-coral-500 bg-navy-900 shadow-lg shadow-coral-500/10 ring-2 ring-coral-500/30'
						: 'border-navy-800 bg-navy-900/60 hover:border-navy-700 hover:bg-navy-900'}"
				>
					<div>
						<!-- Top badge & Wax Seal Dot -->
						<div class="flex items-center justify-between mb-2">
							<div class="flex items-center gap-1.5">
								<!-- Mini Wax Seal Dot with template's real sealColor -->
								<div
									class="h-4 w-4 rounded-full border border-black/20 shadow-xs flex items-center justify-center text-[7px] text-white font-serif"
									style="background-color: {tmpl.palette.sealColor};"
								>
									♥
								</div>
								<span class="text-[10px] font-medium text-navy-400 uppercase">{tmpl.category}</span>
							</div>

							<div class="flex items-center gap-1">
								{#if tmpl.badge}
									<span class="rounded bg-wine-900/80 px-1.5 py-0.5 text-[8px] font-semibold text-champagne-300 uppercase border border-champagne-400/20">
										{tmpl.badge}
									</span>
								{/if}
								<span class="rounded px-1.5 py-0.5 text-[8px] font-semibold uppercase {tmpl.status === 'published' ? 'bg-emerald-500/20 text-emerald-400' : 'bg-amber-500/20 text-amber-400'}">
									{tmpl.status}
								</span>
							</div>
						</div>

						<h4 class="font-serif text-base font-bold text-white transition-colors group-hover:text-coral-300">
							{tmpl.title}
						</h4>
						<p class="mt-0.5 text-[11px] text-navy-400 line-clamp-1">
							{tmpl.subtitle}
						</p>
					</div>

					<div class="mt-3 pt-2.5 border-t border-navy-800/80 flex items-center justify-between text-[11px]">
						<span class="font-serif text-champagne-300 font-semibold">
							Rp {tmpl.basePrice.toLocaleString('id-ID')}
						</span>
						{#if isSelected}
							<span class="flex items-center gap-1 text-coral-400 font-medium">
								<Check size={12} /> Aktif Diedit
							</span>
						{:else}
							<span class="text-navy-500 group-hover:text-navy-300">Pilih</span>
						{/if}
					</div>
				</button>
			{/each}
		</div>
	</div>

	<!-- Main Two-Column Workspace: Left Editor vs Right Real-Time Simulator -->
	<div class="grid gap-6 lg:grid-cols-12 items-start">
		<!-- LEFT COLUMN (7 Cols): Studio Editor Tabs -->
		<div class="lg:col-span-7 space-y-4">
			<!-- Studio Navigation Tabs -->
			<div class="flex rounded-xl border border-navy-800 bg-navy-900/90 p-1">
				<button
					type="button"
					onclick={() => (activeTab = 'visual')}
					class="flex-1 flex items-center justify-center gap-2 rounded-lg py-2.5 text-xs font-semibold transition-all {activeTab ===
					'visual'
						? 'bg-coral-500 text-white shadow-md'
						: 'text-navy-400 hover:text-white'}"
				>
					<Palette size={15} />
					<span>1. Desain &amp; Visual</span>
				</button>

				<button
					type="button"
					onclick={() => (activeTab = 'animations')}
					class="flex-1 flex items-center justify-center gap-2 rounded-lg py-2.5 text-xs font-semibold transition-all {activeTab ===
					'animations'
						? 'bg-coral-500 text-white shadow-md'
						: 'text-navy-400 hover:text-white'}"
				>
					<Sparkles size={15} />
					<span>2. Efek &amp; Animasi</span>
				</button>

				<button
					type="button"
					onclick={() => (activeTab = 'permissions')}
					class="flex-1 flex items-center justify-center gap-2 rounded-lg py-2.5 text-xs font-semibold transition-all {activeTab ===
					'permissions'
						? 'bg-coral-500 text-white shadow-md'
						: 'text-navy-400 hover:text-white'}"
				>
					<ShieldCheck size={15} />
					<span>3. Hak Kustomisasi Member</span>
				</button>
			</div>

			<!-- TAB 1: DESAIN & VISUAL ENGINE -->
			{#if activeTab === 'visual'}
				<div class="space-y-4 animate-in fade-in duration-200">
					<!-- Section 1.1: Metadata Produk -->
					<div class="rounded-2xl border border-navy-800 bg-navy-900/80 p-5 backdrop-blur-md space-y-4">
						<div class="flex items-center justify-between border-b border-navy-800/80 pb-3">
							<h3 class="font-display text-sm font-bold text-white flex items-center gap-2">
								<Layers size={16} class="text-coral-400" />
								<span>Informasi &amp; Identitas Template</span>
							</h3>
							<button
								type="button"
								onclick={handleTogglePublishStatus}
								class="rounded-full px-2.5 py-1 text-[10px] font-bold uppercase transition-colors {currentTemplate.status ===
								'published'
									? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 hover:bg-emerald-500/30'
									: 'bg-amber-500/20 text-amber-400 border border-amber-500/30 hover:bg-amber-500/30'}"
							>
								Status: {currentTemplate.status} (Klik utk ubah)
							</button>
						</div>

						<div class="grid gap-4 sm:grid-cols-2">
							<div>
								<label for="template-title" class="block text-xs font-semibold text-navy-300">Nama Template</label>
								<input
									id="template-title"
									type="text"
									bind:value={currentTemplate.title}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3.5 py-2 text-xs text-white focus:border-coral-500 focus:outline-none font-medium"
								/>
							</div>

							<div>
								<label for="template-slug" class="block text-xs font-semibold text-navy-300">Slug URL</label>
								<input
									id="template-slug"
									type="text"
									bind:value={currentTemplate.slug}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3.5 py-2 text-xs text-white font-mono focus:border-coral-500 focus:outline-none"
								/>
							</div>

							<div>
								<label for="template-category" class="block text-xs font-semibold text-navy-300">Kategori Desain</label>
								<select
									id="template-category"
									bind:value={currentTemplate.category}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3.5 py-2 text-xs text-white focus:border-coral-500 focus:outline-none"
								>
									<option value="Klasik">Klasik</option>
									<option value="Elegan">Elegan</option>
									<option value="Botanical">Botanical</option>
									<option value="Luxury Gold">Luxury Gold</option>
									<option value="Minimalis">Minimalis</option>
									<option value="Modern">Modern</option>
								</select>
							</div>

							<div>
								<label for="template-badge" class="block text-xs font-semibold text-navy-300">Badge Penjualan</label>
								<select
									id="template-badge"
									bind:value={currentTemplate.badge}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3.5 py-2 text-xs text-white focus:border-coral-500 focus:outline-none"
								>
									<option value="">Tanpa Badge</option>
									<option value="BESTSELLER">BESTSELLER</option>
									<option value="POPULER">POPULER</option>
									<option value="BARU">BARU</option>
									<option value="EKSKLUSIF">EKSKLUSIF</option>
								</select>
							</div>

							<div class="sm:col-span-2">
								<label for="template-subtitle" class="block text-xs font-semibold text-navy-300">Deskripsi Singkat / Subtitle</label>
								<input
									id="template-subtitle"
									type="text"
									bind:value={currentTemplate.subtitle}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3.5 py-2 text-xs text-white focus:border-coral-500 focus:outline-none"
								/>
							</div>

							<div class="sm:col-span-2">
								<label for="template-tagline" class="block text-xs font-semibold text-navy-300">Tagline Pemasaran</label>
								<input
									id="template-tagline"
									type="text"
									bind:value={currentTemplate.tagline}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3.5 py-2 text-xs text-white focus:border-coral-500 focus:outline-none"
								/>
							</div>
						</div>
					</div>

					<!-- Section 1.2: Palet Warna & Visual Token Engine -->
					<div class="rounded-2xl border border-navy-800 bg-navy-900/80 p-5 backdrop-blur-md space-y-4">
						<div class="border-b border-navy-800/80 pb-3">
							<h3 class="font-display text-sm font-bold text-white flex items-center gap-2">
								<Palette size={16} class="text-champagne-400" />
								<span>Palet Warna &amp; Styling Engine</span>
							</h3>
							<p class="text-[11px] text-navy-400 mt-0.5">
								Atur kode warna hex untuk aksen, stempel lilin wax seal, dan border frame smartphone.
							</p>
						</div>

						<div class="grid gap-4 sm:grid-cols-2">
							<!-- Accent Color -->
							<div>
								<label for="accent-color-hex" class="block text-xs font-semibold text-navy-300">Warna Aksen Utama</label>
								<div class="mt-1.5 flex items-center gap-2.5">
									<input
										type="color"
										bind:value={currentTemplate.palette.accentColor}
										class="h-9 w-12 cursor-pointer rounded-lg border border-navy-700 bg-transparent p-1"
									/>
									<input
										id="accent-color-hex"
										type="text"
										bind:value={currentTemplate.palette.accentColor}
										class="flex-1 rounded-xl border border-navy-700 bg-navy-950/80 px-3 py-2 text-xs font-mono text-white focus:border-coral-500 focus:outline-none"
									/>
								</div>
							</div>

							<!-- Wax Seal Color with Quick Preset Swatches -->
							<div>
								<label for="seal-color-hex" class="block text-xs font-semibold text-navy-300">Warna Stempel Lilin (Wax Seal)</label>
								<div class="mt-1.5 flex items-center gap-2.5">
									<input
										type="color"
										bind:value={currentTemplate.palette.sealColor}
										class="h-9 w-12 cursor-pointer rounded-lg border border-navy-700 bg-transparent p-1"
									/>
									<input
										id="seal-color-hex"
										type="text"
										bind:value={currentTemplate.palette.sealColor}
										class="flex-1 rounded-xl border border-navy-700 bg-navy-950/80 px-3 py-2 text-xs font-mono text-white focus:border-coral-500 focus:outline-none"
									/>
								</div>

								<!-- Preset Swatches for Wax Seal -->
								<div class="mt-2 flex flex-wrap gap-1.5">
									{#each WAX_SEAL_PRESETS as preset (preset.name)}
										<button
											type="button"
											onclick={() => (currentTemplate.palette.sealColor = preset.color)}
											class="h-5 w-5 rounded-full border border-white/20 transition-transform hover:scale-125 focus:outline-none"
											style="background-color: {preset.color};"
											title={preset.name}
										></button>
									{/each}
								</div>
							</div>

							<!-- Accent Text Color -->
							<div>
								<label for="accent-text-hex" class="block text-xs font-semibold text-navy-300">Warna Teks Aksen</label>
								<div class="mt-1.5 flex items-center gap-2.5">
									<input
										type="color"
										bind:value={currentTemplate.palette.accentTextColor}
										class="h-9 w-12 cursor-pointer rounded-lg border border-navy-700 bg-transparent p-1"
									/>
									<input
										id="accent-text-hex"
										type="text"
										bind:value={currentTemplate.palette.accentTextColor}
										class="flex-1 rounded-xl border border-navy-700 bg-navy-950/80 px-3 py-2 text-xs font-mono text-white focus:border-coral-500 focus:outline-none"
									/>
								</div>
							</div>

							<!-- Phone Border Class -->
							<div>
								<label for="phone-border-class" class="block text-xs font-semibold text-navy-300">Border Frame Mockup HP</label>
								<select
									id="phone-border-class"
									bind:value={currentTemplate.palette.phoneBorder}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3.5 py-2 text-xs text-white focus:border-coral-500 focus:outline-none"
								>
									<option value="border-stone-800">Stone Black (Mewah Klasik)</option>
									<option value="border-wine-900">Wine Bordeaux (Bordeaux Elegan)</option>
									<option value="border-emerald-950">Emerald Deep (Botanical)</option>
									<option value="border-zinc-900">Titanium Zinc (Luxury Gold)</option>
									<option value="border-slate-800">Slate Charcoal (Modern Minimalis)</option>
									<option value="border-rose-900">Rose Deep (Blush Pink)</option>
								</select>
							</div>

							<!-- Background Gradient Presets -->
							<div class="sm:col-span-2">
								<label for="bg-gradient-preset" class="block text-xs font-semibold text-navy-300">Gradasi Latar Belakang Kartu &amp; Layar</label>
								<select
									id="bg-gradient-preset"
									bind:value={currentTemplate.palette.bgGradient}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3.5 py-2 text-xs text-white focus:border-coral-500 focus:outline-none"
								>
									{#each BG_GRADIENT_PRESETS as grad (grad.name)}
										<option value={grad.value}>{grad.name}</option>
									{/each}
								</select>
							</div>
						</div>
					</div>

					<!-- Section 1.3: Tipografi & Monogram -->
					<div class="rounded-2xl border border-navy-800 bg-navy-900/80 p-5 backdrop-blur-md space-y-4">
						<div class="border-b border-navy-800/80 pb-3">
							<h3 class="font-display text-sm font-bold text-white flex items-center gap-2">
								<Wand2 size={16} class="text-coral-400" />
								<span>Tipografi &amp; Pratinjau Monogram</span>
							</h3>
						</div>

						<div class="grid gap-4 sm:grid-cols-3">
							<div>
								<label for="font-family-select" class="block text-xs font-semibold text-navy-300">Font Family Serif</label>
								<select
									id="font-family-select"
									bind:value={currentTemplate.palette.fontFamily}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3 py-2 text-xs text-white focus:border-coral-500 focus:outline-none"
								>
									<option value="Cormorant Garamond">Cormorant Garamond</option>
									<option value="Playfair Display">Playfair Display</option>
									<option value="Cinzel">Cinzel Decorative</option>
									<option value="Montserrat">Montserrat Editorial</option>
								</select>
							</div>

							<div>
								<label for="monogram-initials" class="block text-xs font-semibold text-navy-300">Inisial Monogram Segel</label>
								<input
									id="monogram-initials"
									type="text"
									bind:value={currentTemplate.preview.monogram}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3 py-2 text-xs text-white text-center font-serif text-base focus:border-coral-500 focus:outline-none font-bold"
								/>
							</div>

							<div>
								<label for="seal-symbol" class="block text-xs font-semibold text-navy-300">Simbol Segel Lilin</label>
								<select
									id="seal-symbol"
									bind:value={currentTemplate.envelope.sealSymbol}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3 py-2 text-xs text-white focus:border-coral-500 focus:outline-none"
								>
									<option value="heart">♥ Simbol Hati Romansa</option>
									<option value="monogram">Inisial Monogram (Huruf)</option>
									<option value="botanical">Daun Botanical</option>
									<option value="rings">Cincin Pernikahan</option>
								</select>
							</div>

							<div class="sm:col-span-3">
								<label for="envelope-script" class="block text-xs font-semibold text-navy-300">Teks Kaligrafi Luar Amplop</label>
								<input
									id="envelope-script"
									type="text"
									bind:value={currentTemplate.envelope.envelopeScript}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3.5 py-2 text-xs text-white italic font-serif focus:border-coral-500 focus:outline-none"
								/>
							</div>

							<!-- Dummy Couple preview data -->
							<div>
								<label for="dummy-couple-name" class="block text-xs font-semibold text-navy-300">Nama Pasangan Preview</label>
								<input
									id="dummy-couple-name"
									type="text"
									bind:value={currentTemplate.preview.coupleName}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3 py-2 text-xs text-white focus:border-coral-500 focus:outline-none"
								/>
							</div>

							<div>
								<label for="dummy-event-date" class="block text-xs font-semibold text-navy-300">Tanggal Acara Preview</label>
								<input
									id="dummy-event-date"
									type="text"
									bind:value={currentTemplate.preview.date}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3 py-2 text-xs text-white focus:border-coral-500 focus:outline-none"
								/>
							</div>

							<div>
								<label for="dummy-venue-loc" class="block text-xs font-semibold text-navy-300">Lokasi Venue Preview</label>
								<input
									id="dummy-venue-loc"
									type="text"
									bind:value={currentTemplate.preview.location}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3 py-2 text-xs text-white focus:border-coral-500 focus:outline-none"
								/>
							</div>
						</div>
					</div>
				</div>

			<!-- TAB 2: ANIMASI & INTERAKTIVITAS ENGINE -->
			{:else if activeTab === 'animations'}
				<div class="space-y-4 animate-in fade-in duration-200">
					<!-- Section 2.1: Opening Effect Selector -->
					<div class="rounded-2xl border border-navy-800 bg-navy-900/80 p-5 backdrop-blur-md space-y-4">
						<div class="border-b border-navy-800/80 pb-3">
							<h3 class="font-display text-sm font-bold text-white flex items-center gap-2">
								<Sparkles size={16} class="text-coral-400" />
								<span>Animasi Pembuka Undangan (Opening Transition)</span>
							</h3>
							<p class="text-[11px] text-navy-400 mt-0.5">
								Tentukan animasi transisi ketika tamu pertama kali mengetuk stempel lilin wax seal atau tombol pembuka undangan.
							</p>
						</div>

						<div class="grid gap-3 sm:grid-cols-2">
							<!-- Effect 1 -->
							<button
								type="button"
								onclick={() => (currentTemplate.animations.openingEffect = 'wax_seal_crack')}
								class="flex items-start gap-3 rounded-xl border p-3.5 text-left transition-all {currentTemplate.animations.openingEffect === 'wax_seal_crack'
									? 'border-coral-500 bg-coral-500/10 ring-1 ring-coral-500'
									: 'border-navy-700 bg-navy-950/60 hover:border-navy-600'}"
							>
								<div class="mt-0.5 flex h-7 w-7 items-center justify-center rounded-lg bg-navy-800 text-coral-400 shrink-0">
									<Heart size={14} />
								</div>
								<div>
									<h4 class="text-xs font-bold text-white">Wax Seal Crack &amp; Flap</h4>
									<p class="text-[11px] text-navy-400 mt-0.5 leading-relaxed">
										Stempel lilin bergetar realistis, retak secara 3D dan tutup amplop terbuka ke atas.
									</p>
								</div>
							</button>

							<!-- Effect 2 -->
							<button
								type="button"
								onclick={() => (currentTemplate.animations.openingEffect = 'velvet_curtain')}
								class="flex items-start gap-3 rounded-xl border p-3.5 text-left transition-all {currentTemplate.animations.openingEffect === 'velvet_curtain'
									? 'border-coral-500 bg-coral-500/10 ring-1 ring-coral-500'
									: 'border-navy-700 bg-navy-950/60 hover:border-navy-600'}"
							>
								<div class="mt-0.5 flex h-7 w-7 items-center justify-center rounded-lg bg-navy-800 text-amber-400 shrink-0">
									<Layers size={14} />
								</div>
								<div>
									<h4 class="text-xs font-bold text-white">Velvet Curtain Pull</h4>
									<p class="text-[11px] text-navy-400 mt-0.5 leading-relaxed">
										Tirai beludru mewah bergeser perlahan ke samping menyingkap isi undangan kerajaan.
									</p>
								</div>
							</button>

							<!-- Effect 3 -->
							<button
								type="button"
								onclick={() => (currentTemplate.animations.openingEffect = 'blooming_rose')}
								class="flex items-start gap-3 rounded-xl border p-3.5 text-left transition-all {currentTemplate.animations.openingEffect === 'blooming_rose'
									? 'border-coral-500 bg-coral-500/10 ring-1 ring-coral-500'
									: 'border-navy-700 bg-navy-950/60 hover:border-navy-600'}"
							>
								<div class="mt-0.5 flex h-7 w-7 items-center justify-center rounded-lg bg-navy-800 text-rose-400 shrink-0">
									<Sparkles size={14} />
								</div>
								<div>
									<h4 class="text-xs font-bold text-white">Blooming Rose Unfold</h4>
									<p class="text-[11px] text-navy-400 mt-0.5 leading-relaxed">
										Kelopak bunga bermekaran memutar perlahan sebelum halaman surat undangan muncul.
									</p>
								</div>
							</button>

							<!-- Effect 4 -->
							<button
								type="button"
								onclick={() => (currentTemplate.animations.openingEffect = 'gate_fold')}
								class="flex items-start gap-3 rounded-xl border p-3.5 text-left transition-all {currentTemplate.animations.openingEffect === 'gate_fold'
									? 'border-coral-500 bg-coral-500/10 ring-1 ring-coral-500'
									: 'border-navy-700 bg-navy-950/60 hover:border-navy-600'}"
							>
								<div class="mt-0.5 flex h-7 w-7 items-center justify-center rounded-lg bg-navy-800 text-emerald-400 shrink-0">
									<ShieldCheck size={14} />
								</div>
								<div>
									<h4 class="text-xs font-bold text-white">Royal Gate Opening</h4>
									<p class="text-[11px] text-navy-400 mt-0.5 leading-relaxed">
										Pintu gerbang ornamen vintage terbelah ganda dengan transisi perspektif 3D.
									</p>
								</div>
							</button>
						</div>

						<!-- Opening Duration Slider -->
						<div class="pt-3 border-t border-navy-800/80">
							<div class="flex items-center justify-between text-xs">
								<span class="font-semibold text-navy-300">Durasi Transisi Buka</span>
								<span class="font-mono text-coral-400 font-bold">{currentTemplate.animations.openingDurationMs} ms</span>
							</div>
							<input
								type="range"
								min="400"
								max="1500"
								step="50"
								bind:value={currentTemplate.animations.openingDurationMs}
								class="mt-2 w-full accent-coral-500"
							/>
						</div>
					</div>

					<!-- Section 2.2: Particle Effect Engine -->
					<div class="rounded-2xl border border-navy-800 bg-navy-900/80 p-5 backdrop-blur-md space-y-4">
						<div class="border-b border-navy-800/80 pb-3">
							<h3 class="font-display text-sm font-bold text-white flex items-center gap-2">
								<Sparkles size={16} class="text-champagne-400" />
								<span>Efek Partikel Visual Melayang (Floating Particles)</span>
							</h3>
							<p class="text-[11px] text-navy-400 mt-0.5">
								Animasi latar belakang interaktif saat tamu membaca isi undangan.
							</p>
						</div>

						<div class="grid gap-3 sm:grid-cols-3">
							<button
								type="button"
								onclick={() => (currentTemplate.animations.particleEffect = 'falling_petals')}
								class="rounded-xl border p-3 text-center transition-all {currentTemplate.animations.particleEffect === 'falling_petals'
									? 'border-coral-500 bg-coral-500/10 ring-1 ring-coral-500'
									: 'border-navy-700 bg-navy-950/60 hover:border-navy-600'}"
							>
								<span class="text-xl">🌸</span>
								<h5 class="text-xs font-bold text-white mt-1">Kelopak Bunga Gugur</h5>
								<p class="text-[10px] text-navy-400 mt-0.5">Petals sakura &amp; mawar</p>
							</button>

							<button
								type="button"
								onclick={() => (currentTemplate.animations.particleEffect = 'gold_stardust')}
								class="rounded-xl border p-3 text-center transition-all {currentTemplate.animations.particleEffect === 'gold_stardust'
									? 'border-coral-500 bg-coral-500/10 ring-1 ring-coral-500'
									: 'border-navy-700 bg-navy-950/60 hover:border-navy-600'}"
							>
								<span class="text-xl">✨</span>
								<h5 class="text-xs font-bold text-white mt-1">Debu Bintang Emas</h5>
								<p class="text-[10px] text-navy-400 mt-0.5">Gold Stardust Glow</p>
							</button>

							<button
								type="button"
								onclick={() => (currentTemplate.animations.particleEffect = 'confetti')}
								class="rounded-xl border p-3 text-center transition-all {currentTemplate.animations.particleEffect === 'confetti'
									? 'border-coral-500 bg-coral-500/10 ring-1 ring-coral-500'
									: 'border-navy-700 bg-navy-950/60 hover:border-navy-600'}"
							>
								<span class="text-xl">🎉</span>
								<h5 class="text-xs font-bold text-white mt-1">Confetti Pesta</h5>
								<p class="text-[10px] text-navy-400 mt-0.5">Kertas emas &amp; blush</p>
							</button>

							<button
								type="button"
								onclick={() => (currentTemplate.animations.particleEffect = 'floating_lanterns')}
								class="rounded-xl border p-3 text-center transition-all {currentTemplate.animations.particleEffect === 'floating_lanterns'
									? 'border-coral-500 bg-coral-500/10 ring-1 ring-coral-500'
									: 'border-navy-700 bg-navy-950/60 hover:border-navy-600'}"
							>
								<span class="text-xl">🏮</span>
								<h5 class="text-xs font-bold text-white mt-1">Lentera Mengambang</h5>
								<p class="text-[10px] text-navy-400 mt-0.5">Warm Glowing Orbs</p>
							</button>

							<button
								type="button"
								onclick={() => (currentTemplate.animations.particleEffect = 'none')}
								class="sm:col-span-2 rounded-xl border p-3 text-center transition-all {currentTemplate.animations.particleEffect === 'none'
									? 'border-coral-500 bg-coral-500/10 ring-1 ring-coral-500'
									: 'border-navy-700 bg-navy-950/60 hover:border-navy-600'}"
							>
								<span class="text-xl">🚫</span>
								<h5 class="text-xs font-bold text-white mt-1">Tanpa Partikel (Clean &amp; Minimal)</h5>
								<p class="text-[10px] text-navy-400 mt-0.5">Fokus penuh pada tipografi</p>
							</button>
						</div>

						<!-- Particle Density -->
						<div class="pt-3 border-t border-navy-800/80 flex items-center justify-between">
							<span class="text-xs font-semibold text-navy-300">Kepadatan Partikel:</span>
							<div class="flex gap-2">
								{#each ['low', 'medium', 'high'] as den}
									<button
										type="button"
										onclick={() => (currentTemplate.animations.particleDensity = den as ParticleDensity)}
										class="rounded-lg px-3 py-1 text-xs font-semibold uppercase transition-colors {currentTemplate.animations.particleDensity === den
											? 'bg-coral-500 text-white'
											: 'bg-navy-800 text-navy-400 hover:text-white'}"
									>
										{den}
									</button>
								{/each}
							</div>
						</div>
					</div>

					<!-- Section 2.3: Audio & Soundscape Presets -->
					<div class="rounded-2xl border border-navy-800 bg-navy-900/80 p-5 backdrop-blur-md space-y-4">
						<div class="border-b border-navy-800/80 pb-3">
							<h3 class="font-display text-sm font-bold text-white flex items-center gap-2">
								<Music size={16} class="text-undangan-400" />
								<span>Preset Musik &amp; Soundscape Latar</span>
							</h3>
							<p class="text-[11px] text-navy-400 mt-0.5">
								Pilihan lagu default yang menyala saat tamu membuka undangan digital.
							</p>
						</div>

						<div class="space-y-3">
							<div>
								<label for="track-preset-select" class="block text-xs font-semibold text-navy-300">Lagu Preset Bawaan</label>
								<select
									id="track-preset-select"
									bind:value={currentTemplate.audio.presetTrackName}
									class="mt-1.5 w-full rounded-xl border border-navy-700 bg-navy-950/80 px-3.5 py-2 text-xs text-white focus:border-coral-500 focus:outline-none"
								>
									{#each AUDIO_PRESETS as track (track.name)}
										<option value={track.name}>{track.name}</option>
									{/each}
								</select>
							</div>

							<div class="grid gap-3 sm:grid-cols-2 pt-2">
								<label class="flex items-center gap-3 rounded-xl border border-navy-700 bg-navy-950/60 p-3 cursor-pointer">
									<input
										type="checkbox"
										bind:checked={currentTemplate.audio.autoplay}
										class="h-4 w-4 rounded accent-coral-500"
									/>
									<div>
										<p class="text-xs font-bold text-white">Autoplay Saat Dibuka</p>
										<p class="text-[10px] text-navy-400">Musik otomatis memudar masuk</p>
									</div>
								</label>

								<label class="flex items-center gap-3 rounded-xl border border-navy-700 bg-navy-950/60 p-3 cursor-pointer">
									<input
										type="checkbox"
										bind:checked={currentTemplate.audio.showVisualizer}
										class="h-4 w-4 rounded accent-coral-500"
									/>
									<div>
										<p class="text-xs font-bold text-white">Audio Wave Visualizer</p>
										<p class="text-[10px] text-navy-400">Tampilkan gelombang bergetar</p>
									</div>
								</label>
							</div>
						</div>
					</div>
				</div>

			<!-- TAB 3: MATRIKS HAK KUSTOMISASI MEMBER PER PAKET -->
			{:else if activeTab === 'permissions'}
				<div class="space-y-4 animate-in fade-in duration-200">
					<!-- Explanation Banner -->
					<div class="rounded-2xl border border-champagne-400/30 bg-gradient-to-r from-wine-950/60 via-navy-900 to-navy-900 p-4 sm:p-5">
						<div class="flex items-start gap-3">
							<div class="rounded-lg bg-champagne-400/20 p-2 text-champagne-300 shrink-0">
								<ShieldCheck size={20} />
							</div>
							<div>
								<h3 class="font-display text-sm font-bold text-white">
									Matriks Hak Kustomisasi Member (Package Feature Gates)
								</h3>
								<p class="mt-1 text-xs text-navy-300 leading-relaxed">
									Atur secara fleksibel hal apa saja yang bisa diedit oleh pengguna/member yang membeli jasa paket pada tema ini. Batasan ini akan langsung diterapkan di dashboard editor pengantin.
								</p>
							</div>
						</div>
					</div>

					<!-- Comparison Matrix Table -->
					<div class="overflow-x-auto rounded-2xl border border-navy-800 bg-navy-900/90 shadow-xl">
						<table class="w-full text-left text-xs border-collapse">
							<thead>
								<tr class="border-b border-navy-800 bg-navy-950/70 text-navy-300">
									<th class="p-3.5 font-bold uppercase tracking-wider text-[10px] w-5/12">
										Fitur Kustomisasi Member
									</th>
									<th class="p-3.5 text-center font-bold uppercase tracking-wider text-[10px] w-7/36 bg-navy-900/40">
										<span class="block text-navy-300">Paket Gratis</span>
										<span class="text-[9px] text-navy-400 font-normal">Rp 0</span>
									</th>
									<th class="p-3.5 text-center font-bold uppercase tracking-wider text-[10px] w-7/36 bg-coral-950/30 text-coral-300 border-x border-coral-500/20">
										<span class="block">Paket Pro</span>
										<span class="text-[9px] text-coral-400 font-normal">Rp 99.000</span>
									</th>
									<th class="p-3.5 text-center font-bold uppercase tracking-wider text-[10px] w-7/36 bg-champagne-950/30 text-champagne-300">
										<span class="block">Paket Lengkap</span>
										<span class="text-[9px] text-champagne-400 font-normal">Rp 249.000</span>
									</th>
								</tr>
							</thead>

							<tbody class="divide-y divide-navy-800/60 text-navy-200">
								<!-- Row 1: Kuota Tamu -->
								<tr class="hover:bg-navy-800/30">
									<td class="p-3.5 font-medium">
										<div class="flex items-center gap-2">
											<span class="font-bold text-white">Batas Kuota Tamu RSVP</span>
										</div>
										<span class="text-[10px] text-navy-400 block">Maksimal nama tamu yang dapat di-generate tautan personal</span>
									</td>
									<td class="p-3.5 text-center bg-navy-900/40">
										<input
											type="number"
											bind:value={currentTemplate.customizationMatrix.gratis.guestLimit}
											class="w-16 rounded border border-navy-700 bg-navy-950 px-2 py-1 text-center text-xs text-white"
										/>
									</td>
									<td class="p-3.5 text-center bg-coral-950/20 border-x border-coral-500/20 font-bold text-emerald-400">
										Unlimited
									</td>
									<td class="p-3.5 text-center bg-champagne-950/20 font-bold text-emerald-400">
										Unlimited
									</td>
								</tr>

								<!-- Row 2: Galeri Foto -->
								<tr class="hover:bg-navy-800/30">
									<td class="p-3.5 font-medium">
										<span class="font-bold text-white">Batas Foto Galeri Momen</span>
										<span class="text-[10px] text-navy-400 block">Jumlah foto prewedding resolusi tinggi</span>
									</td>
									<td class="p-3.5 text-center bg-navy-900/40">
										<input
											type="number"
											bind:value={currentTemplate.customizationMatrix.gratis.galleryPhotoLimit}
											class="w-16 rounded border border-navy-700 bg-navy-950 px-2 py-1 text-center text-xs text-white"
										/>
									</td>
									<td class="p-3.5 text-center bg-coral-950/20 border-x border-coral-500/20">
										<input
											type="number"
											bind:value={currentTemplate.customizationMatrix.pro.galleryPhotoLimit}
											class="w-16 rounded border border-coral-500/40 bg-navy-950 px-2 py-1 text-center text-xs text-coral-300 font-bold"
										/>
									</td>
									<td class="p-3.5 text-center bg-champagne-950/20">
										<input
											type="number"
											bind:value={currentTemplate.customizationMatrix.lengkap.galleryPhotoLimit}
											class="w-16 rounded border border-champagne-400/40 bg-navy-950 px-2 py-1 text-center text-xs text-champagne-300 font-bold"
										/>
									</td>
								</tr>

								<!-- Row 3: Upload Musik Sendiri -->
								<tr class="hover:bg-navy-800/30">
									<td class="p-3.5 font-medium">
										<span class="font-bold text-white">Upload Musik Kustom (MP3 Sendiri)</span>
										<span class="text-[10px] text-navy-400 block">Jika nonaktif, pengantin hanya boleh memakai preset audio bawaan</span>
									</td>
									<td class="p-3.5 text-center bg-navy-900/40">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.gratis.allowCustomMusic}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
									<td class="p-3.5 text-center bg-coral-950/20 border-x border-coral-500/20">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.pro.allowCustomMusic}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
									<td class="p-3.5 text-center bg-champagne-950/20">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.lengkap.allowCustomMusic}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
								</tr>

								<!-- Row 4: Amplop Digital & QRIS -->
								<tr class="hover:bg-navy-800/30">
									<td class="p-3.5 font-medium">
										<span class="font-bold text-white">Amplop Digital &amp; Rekening QRIS</span>
										<span class="text-[10px] text-navy-400 block">Kirim kado pernikahan transfer langsung BCA/Mandiri/QRIS</span>
									</td>
									<td class="p-3.5 text-center bg-navy-900/40">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.gratis.allowDigitalEnvelope}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
									<td class="p-3.5 text-center bg-coral-950/20 border-x border-coral-500/20">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.pro.allowDigitalEnvelope}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
									<td class="p-3.5 text-center bg-champagne-950/20">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.lengkap.allowDigitalEnvelope}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
								</tr>

								<!-- Row 5: Cerita Cinta / Love Story Timeline -->
								<tr class="hover:bg-navy-800/30">
									<td class="p-3.5 font-medium">
										<span class="font-bold text-white">Timeline Cerita Cinta (Love Story)</span>
										<span class="text-[10px] text-navy-400 block">Kisah awal bertemu, lamaran, hingga menuju pernikahan</span>
									</td>
									<td class="p-3.5 text-center bg-navy-900/40">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.gratis.allowLoveStory}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
									<td class="p-3.5 text-center bg-coral-950/20 border-x border-coral-500/20">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.pro.allowLoveStory}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
									<td class="p-3.5 text-center bg-champagne-950/20">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.lengkap.allowLoveStory}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
								</tr>

								<!-- Row 6: Custom Domain -->
								<tr class="hover:bg-navy-800/30">
									<td class="p-3.5 font-medium">
										<span class="font-bold text-white">Custom Domain Pribadi (.com / .id)</span>
										<span class="text-[10px] text-navy-400 block">Contoh: julian-elena.com tanpa embel-embel ketuk.id</span>
									</td>
									<td class="p-3.5 text-center bg-navy-900/40">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.gratis.allowCustomDomain}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
									<td class="p-3.5 text-center bg-coral-950/20 border-x border-coral-500/20">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.pro.allowCustomDomain}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
									<td class="p-3.5 text-center bg-champagne-950/20">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.lengkap.allowCustomDomain}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
								</tr>

								<!-- Row 7: Override Palet Warna oleh Member -->
								<tr class="hover:bg-navy-800/30">
									<td class="p-3.5 font-medium">
										<span class="font-bold text-white">Member Boleh Ganti Warna Aksen &amp; Stempel</span>
										<span class="text-[10px] text-navy-400 block">Pengantin bebas memilih warna tema sesuai konsep busana</span>
									</td>
									<td class="p-3.5 text-center bg-navy-900/40">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.gratis.allowCustomPaletteOverride}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
									<td class="p-3.5 text-center bg-coral-950/20 border-x border-coral-500/20">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.pro.allowCustomPaletteOverride}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
									<td class="p-3.5 text-center bg-champagne-950/20">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.lengkap.allowCustomPaletteOverride}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
								</tr>

								<!-- Row 8: Hilangkan Watermark -->
								<tr class="hover:bg-navy-800/30">
									<td class="p-3.5 font-medium">
										<span class="font-bold text-white">Hilangkan Watermark "Powered by Ketuk.id"</span>
										<span class="text-[10px] text-navy-400 block">Tampilan 100% white-label eksklusif nama pengantin</span>
									</td>
									<td class="p-3.5 text-center bg-navy-900/40">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.gratis.removeKetukWatermark}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
									<td class="p-3.5 text-center bg-coral-950/20 border-x border-coral-500/20">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.pro.removeKetukWatermark}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
									<td class="p-3.5 text-center bg-champagne-950/20">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.lengkap.removeKetukWatermark}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
								</tr>

								<!-- Row 9: WhatsApp Blast Generator -->
								<tr class="hover:bg-navy-800/30">
									<td class="p-3.5 font-medium">
										<span class="font-bold text-white">Generator WhatsApp Blast Personal</span>
										<span class="text-[10px] text-navy-400 block">Tombol kirim satu-satu pesan sopan otomatis ke kontak WA tamu</span>
									</td>
									<td class="p-3.5 text-center bg-navy-900/40">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.gratis.allowWhatsAppBlast}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
									<td class="p-3.5 text-center bg-coral-950/20 border-x border-coral-500/20">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.pro.allowWhatsAppBlast}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
									<td class="p-3.5 text-center bg-champagne-950/20">
										<input
											type="checkbox"
											bind:checked={currentTemplate.customizationMatrix.lengkap.allowWhatsAppBlast}
											class="h-4 w-4 rounded accent-coral-500 cursor-pointer"
										/>
									</td>
								</tr>
							</tbody>
						</table>
					</div>

					<!-- Bottom helper action for matrix -->
					<div class="flex items-center justify-between pt-2">
						<button
							type="button"
							onclick={() => {
								currentTemplate.customizationMatrix = JSON.parse(
									JSON.stringify(DEFAULT_PLAN_PERMISSIONS),
								);
								showNotification('Matriks kustomisasi direset ke standar paket.');
							}}
							class="text-xs text-navy-400 hover:text-white underline"
						>
							Reset Matriks ke Nilai Standar
						</button>

						<button
							type="button"
							onclick={handleSaveAll}
							class="inline-flex items-center gap-1.5 rounded-xl bg-coral-500 px-4 py-2 text-xs font-bold uppercase text-white shadow-md hover:bg-coral-600"
						>
							<Check size={14} />
							<span>Simpan Matriks Paket</span>
						</button>
					</div>
				</div>
			{/if}
		</div>

		<!-- RIGHT COLUMN (5 Cols): Real-Time Interactive Live Simulator -->
		<div class="lg:col-span-5 sticky top-20 space-y-4">
			<div class="rounded-3xl border border-navy-800 bg-navy-900/90 p-5 backdrop-blur-xl shadow-2xl">
				<!-- Simulator Stage Controls Header -->
				<div class="flex items-center justify-between border-b border-navy-800/80 pb-3">
					<div>
						<span class="text-[10px] font-semibold uppercase tracking-widest text-navy-400">
							LIVE SIMULATOR
						</span>
						<h4 class="font-display text-sm font-bold text-white truncate max-w-[200px]">
							{currentTemplate.title}
						</h4>
					</div>

					<!-- Switcher: Envelope View vs Phone View -->
					<div class="flex rounded-lg border border-navy-700 bg-navy-950 p-0.5 text-xs">
						<button
							type="button"
							onclick={() => (previewMode = 'envelope')}
							class="flex items-center gap-1.5 rounded-md px-2.5 py-1 font-medium transition-all {previewMode ===
							'envelope'
								? 'bg-coral-500 text-white shadow-xs'
								: 'text-navy-400 hover:text-white'}"
							title="Pratinjau Amplop Segel Lilin"
						>
							<Mail size={13} />
							<span>Amplop</span>
						</button>
						<button
							type="button"
							onclick={() => (previewMode = 'phone')}
							class="flex items-center gap-1.5 rounded-md px-2.5 py-1 font-medium transition-all {previewMode ===
							'phone'
								? 'bg-coral-500 text-white shadow-xs'
								: 'text-navy-400 hover:text-white'}"
							title="Pratinjau Layar Smartphone"
						>
							<Smartphone size={13} />
							<span>Ponsel HP</span>
						</button>
					</div>
				</div>

				<!-- Live Preview Area -->
				<div class="py-4 flex flex-col items-center justify-center min-h-[460px]">
					{#if previewMode === 'envelope'}
						<!-- ENVELOPE STAGE SIMULATOR -->
						<div class="relative w-full max-w-[320px] aspect-[9/14] rounded-2xl border border-champagne-300/50 shadow-2xl overflow-hidden select-none transition-all duration-300 bg-gradient-to-b from-stone-50 via-cream-50 to-amber-50/40">
							<!-- Subtle paper pattern -->
							<div
								class="absolute inset-0 opacity-40 mix-blend-multiply pointer-events-none"
								style="background-image: radial-gradient(#d6bc98 0.75px, transparent 0.75px); background-size: 20px 20px;"
							></div>

							<!-- Floating Particles Layer if enabled -->
							{#if currentTemplate.animations.particleEffect !== 'none'}
								<div class="absolute inset-0 pointer-events-none z-30 overflow-hidden">
									{#if currentTemplate.animations.particleEffect === 'falling_petals'}
										<div class="absolute top-4 left-6 text-sm animate-bounce opacity-70">🌸</div>
										<div class="absolute top-16 right-8 text-xs animate-pulse opacity-60">🌸</div>
										<div class="absolute bottom-12 left-10 text-base animate-pulse opacity-75">🌸</div>
									{:else if currentTemplate.animations.particleEffect === 'gold_stardust'}
										<div class="absolute top-6 right-6 text-amber-400 text-xs animate-ping">✨</div>
										<div class="absolute bottom-20 left-8 text-amber-300 text-xs animate-pulse">✨</div>
										<div class="absolute top-28 left-12 text-champagne-400 text-sm animate-pulse">✨</div>
									{:else if currentTemplate.animations.particleEffect === 'confetti'}
										<div class="absolute top-8 left-8 text-xs animate-bounce">🎉</div>
										<div class="absolute top-20 right-10 text-xs animate-spin">🎊</div>
									{/if}
								</div>
							{/if}

							<!-- Triangular Flap (Envelope Flap) -->
							<div
								class="absolute top-0 left-0 right-0 h-[44%] transition-transform duration-700 origin-top z-10 {isEnvelopeOpen
									? '-rotate-x-180 opacity-0 pointer-events-none'
									: ''}"
								style="perspective: 1000px;"
							>
								<svg class="w-full h-full drop-shadow-sm" viewBox="0 0 400 240" preserveAspectRatio="none">
									<path d="M0,0 L400,0 L200,225 Z" fill="#FAF7F2" stroke="#EADCC8" stroke-width="1.5" />
								</svg>
							</div>

							{#if !isEnvelopeOpen}
								<!-- CLOSED ENVELOPE: Wax Seal & Script Text -->
								<div class="absolute inset-0 flex flex-col items-center justify-between p-6 z-20">
									<span class="text-[8px] font-semibold uppercase tracking-[0.25em] text-espresso-400 pt-2">
										The Wedding Collection
									</span>

									<!-- Center Interactive Wax Seal -->
									<div class="my-auto flex flex-col items-center">
										<button
											type="button"
											onclick={triggerTestAnimation}
											class="group relative flex items-center justify-center transition-transform hover:scale-105 active:scale-95 focus:outline-none"
											title="Klik untuk membuka segel"
										>
											<!-- Wax Seal 3D Blob rendered with exact sealColor from editor -->
											<div
												class="relative flex h-20 w-20 items-center justify-center rounded-full shadow-[0_10px_20px_-4px_rgba(40,20,10,0.35),inset_0_2px_4px_rgba(255,255,255,0.4),inset_0_-3px_5px_rgba(0,0,0,0.35)] transition-colors duration-300"
												style="background-color: {currentTemplate.palette.sealColor};"
											>
												<div class="flex h-15 w-15 items-center justify-center rounded-full border border-black/15 shadow-inner bg-black/5">
													{#if currentTemplate.envelope.sealSymbol === 'heart'}
														<span class="text-white text-lg drop-shadow-xs">♥</span>
													{:else if currentTemplate.envelope.sealSymbol === 'botanical'}
														<span class="text-white text-base drop-shadow-xs">🌿</span>
													{:else if currentTemplate.envelope.sealSymbol === 'rings'}
														<span class="text-white text-base drop-shadow-xs">💍</span>
													{:else}
														<span class="font-serif text-lg font-bold text-white drop-shadow-xs italic">
															{currentTemplate.preview.monogram}
														</span>
													{/if}
												</div>
											</div>

											<!-- Tooltip Pulse -->
											<span class="absolute -bottom-6 whitespace-nowrap rounded bg-espresso-950/90 px-2 py-0.5 text-[8px] text-white opacity-0 group-hover:opacity-100 transition-opacity">
												Klik Buka
											</span>
										</button>

										<!-- Script quote from editor -->
										<div class="mt-6 text-center px-2">
											<p
												class="font-serif text-base italic tracking-wide font-normal leading-relaxed"
												style="color: {currentTemplate.palette.accentColor}; font-family: {currentTemplate.palette.fontFamily}, serif;"
											>
												"{currentTemplate.envelope.envelopeScript}"
											</p>
											<p class="font-serif mt-1 text-xs text-espresso-600 italic">
												{currentTemplate.preview.coupleName}
											</p>
										</div>
									</div>

									<span class="text-[9px] font-medium tracking-wider text-espresso-500 pb-1">
										{currentTemplate.preview.date}
									</span>
								</div>
							{:else}
								<!-- OPENED ENVELOPE: Letter Contents -->
								<div class="absolute inset-0 z-20 flex flex-col justify-between p-5 bg-gradient-to-b from-cream-50 via-white to-cream-50 text-center animate-in fade-in duration-300">
									<div class="flex items-center justify-between pb-2 border-b border-cream-200 text-[10px]">
										<span class="text-espresso-600 font-medium">
											Animasi: {currentTemplate.animations.openingEffect}
										</span>
										<button
											type="button"
											onclick={triggerTestAnimation}
											class="text-wine-800 font-bold hover:underline"
										>
											Tutup Amplop
										</button>
									</div>

									<div class="py-2 space-y-3">
										<span class="text-[8px] font-semibold uppercase tracking-[0.2em] text-champagne-600">
											Pernikahan Suci
										</span>
										<h4
											class="font-serif text-2xl font-bold italic"
											style="color: {currentTemplate.palette.accentColor}; font-family: {currentTemplate.palette.fontFamily}, serif;"
										>
											{currentTemplate.preview.coupleName}
										</h4>
										<div class="mx-auto h-0.5 w-8 rounded-full bg-champagne-400"></div>

										<div class="rounded-xl border border-cream-200 bg-cream-50/80 p-3 text-left space-y-1.5 text-[10px] text-espresso-800">
											<p class="flex items-center gap-1.5 font-serif font-semibold">
												<Calendar size={11} class="text-wine-700" />
												<span>{currentTemplate.preview.date}</span>
											</p>
											<p class="flex items-center gap-1.5 font-serif">
												<MapPin size={11} class="text-wine-700" />
												<span class="truncate">{currentTemplate.preview.location}</span>
											</p>
										</div>

										<div
											class="w-full rounded-full py-1.5 text-[10px] font-bold text-white shadow-sm"
											style="background-color: {currentTemplate.palette.accentColor};"
										>
											Buka Undangan Lengkap
										</div>
									</div>

									<span class="font-serif text-[10px] italic text-espresso-500">
										Simulasi Tampilan Surat Selesai
									</span>
								</div>
							{/if}
						</div>
					{:else}
						<!-- SMARTPHONE PHONE VIEW SIMULATOR -->
						<div class="relative w-full max-w-[270px] aspect-[9/18.5] rounded-[42px] border-[8px] {currentTemplate.palette.phoneBorder} bg-white shadow-2xl overflow-hidden select-none transition-all duration-300">
							<!-- Dynamic Island -->
							<div class="absolute top-2 left-1/2 -translate-x-1/2 z-30 h-4 w-20 rounded-full bg-black flex items-center justify-between px-2">
								<span class="h-1.5 w-1.5 rounded-full bg-zinc-800"></span>
								<span class="h-1.5 w-1.5 rounded-full bg-zinc-700"></span>
							</div>

							<!-- Phone Screen Graphic -->
							<div class="flex h-full w-full flex-col justify-between p-4 pt-8 text-center {currentTemplate.palette.bgGradient}">
								<div>
									<span
										class="text-[8px] font-semibold tracking-widest uppercase"
										style="color: {currentTemplate.palette.accentTextColor};"
									>
										The Wedding of
									</span>

									<!-- Couple Name -->
									<h4
										class="font-serif text-xl font-bold mt-2 leading-tight"
										style="color: {currentTemplate.palette.accentColor}; font-family: {currentTemplate.palette.fontFamily}, serif;"
									>
										{currentTemplate.preview.coupleName}
									</h4>
									<div class="mx-auto mt-1.5 h-0.5 w-8 rounded-full bg-champagne-400"></div>
								</div>

								<!-- Center Card inside Phone -->
								<div class="rounded-2xl border border-champagne-300/50 bg-white/70 p-3 shadow-xs space-y-2">
									<!-- Monogram Circle -->
									<div
										class="mx-auto flex h-12 w-12 items-center justify-center rounded-full text-white font-serif text-sm font-bold shadow-xs"
										style="background-color: {currentTemplate.palette.sealColor};"
									>
										{currentTemplate.preview.monogram}
									</div>

									<div class="space-y-0.5">
										<p class="text-[9px] font-semibold text-espresso-800">{currentTemplate.preview.date}</p>
										<p class="text-[8px] text-espresso-500 line-clamp-1">{currentTemplate.preview.location}</p>
									</div>

									<!-- Countdown Mini Pill -->
									<div class="grid grid-cols-3 gap-1 pt-1 text-[8px] font-semibold text-espresso-700">
										<div class="rounded bg-white p-1 shadow-2xs">24 Hari</div>
										<div class="rounded bg-white p-1 shadow-2xs">12 Jam</div>
										<div class="rounded bg-white p-1 shadow-2xs">45 Mnt</div>
									</div>
								</div>

								<!-- Bottom CTA & Audio Indicator -->
								<div class="space-y-2">
									<div
										class="w-full rounded-full py-2 text-[10px] font-bold text-white shadow-md"
										style="background-color: {currentTemplate.palette.accentColor};"
									>
										Konfirmasi Kehadiran (RSVP)
									</div>

									<!-- Floating Audio Indicator if Autoplay is on -->
									<div class="flex items-center justify-center gap-1.5 text-[8px] text-espresso-600">
										<Music size={10} class="text-wine-800 animate-pulse" />
										<span class="truncate max-w-[160px]">{currentTemplate.audio.presetTrackName}</span>
									</div>
								</div>

								<!-- Phone Home Indicator Bar -->
								<div class="h-1 w-20 mx-auto rounded-full bg-black/30 mt-1"></div>
							</div>
						</div>
					{/if}
				</div>

				<!-- Simulator Action Controls Footer -->
				<div class="mt-4 pt-3 border-t border-navy-800/80 flex items-center justify-between">
					<button
						type="button"
						onclick={triggerTestAnimation}
						class="flex-1 flex items-center justify-center gap-1.5 rounded-xl bg-gradient-to-r from-coral-600 to-coral-500 py-2.5 text-xs font-bold uppercase tracking-wider text-white shadow-lg shadow-coral-500/20 hover:scale-102 active:scale-98 transition-all"
					>
						<Zap size={14} />
						<span>{isEnvelopeOpen ? 'Tutup & Reset Amplop' : '⚡ Uji Buka Animasi'}</span>
					</button>

					<button
						type="button"
						onclick={() => (simulatedMusicPlaying = !simulatedMusicPlaying)}
						class="ml-2 flex h-9 w-9 items-center justify-center rounded-xl border border-navy-700 bg-navy-800 text-navy-200 hover:text-white transition-colors"
						title={simulatedMusicPlaying ? 'Jeda Audio Simulasi' : 'Uji Putar Musik'}
					>
						{#if simulatedMusicPlaying}
							<Volume2 size={16} class="text-coral-400 animate-pulse" />
						{:else}
							<VolumeX size={16} class="text-navy-400" />
						{/if}
					</button>
				</div>
			</div>

			<!-- Quick Summary Card of Permissions configured -->
			<div class="rounded-2xl border border-navy-800 bg-navy-900/60 p-4 text-xs space-y-2">
				<div class="flex items-center justify-between text-navy-300 font-semibold">
					<span>Ringkasan Hak Member (Pro Tier)</span>
					<span class="text-coral-400">Rp 99.000</span>
				</div>
				<ul class="space-y-1 text-[11px] text-navy-400">
					<li class="flex items-center gap-1.5">
						<Check size={12} class="text-emerald-400" />
						<span>Tamu RSVP: {currentTemplate.customizationMatrix.pro.guestLimit === -1 ? 'Unlimited' : `${currentTemplate.customizationMatrix.pro.guestLimit} Tamu`}</span>
					</li>
					<li class="flex items-center gap-1.5">
						<Check size={12} class="text-emerald-400" />
						<span>Maksimal {currentTemplate.customizationMatrix.pro.galleryPhotoLimit} Foto Galeri</span>
					</li>
					<li class="flex items-center gap-1.5">
						{#if currentTemplate.customizationMatrix.pro.allowCustomMusic}
							<Check size={12} class="text-emerald-400" />
							<span>Upload Lagu MP3 Sendiri: <strong>Diizinkan</strong></span>
						{:else}
							<X size={12} class="text-amber-400" />
							<span>Upload Lagu MP3: Hanya Preset</span>
						{/if}
					</li>
					<li class="flex items-center gap-1.5">
						{#if currentTemplate.customizationMatrix.pro.allowDigitalEnvelope}
							<Check size={12} class="text-emerald-400" />
							<span>Amplop Digital &amp; QRIS: <strong>Aktif</strong></span>
						{:else}
							<X size={12} class="text-navy-500" />
							<span>Amplop Digital: Nonaktif</span>
						{/if}
					</li>
				</ul>
			</div>
		</div>
	</div>
</div>
