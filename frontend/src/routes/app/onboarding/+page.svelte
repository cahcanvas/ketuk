<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { EVENT_CATEGORIES } from '@ketuk/shared';
	import type { EventCategoryId, UserEventPreference } from '@ketuk/shared';
	import { Button, Input } from '$lib/components/ui';
	import { createOnboardingPreference } from '$lib/onboarding/preference.svelte';

	const preference = createOnboardingPreference();

	// `next` dibawa dari gate layout (halaman yang dituju user sebelum diarahkan
	// ke sini). Hanya path internal yang dihormati — tanpa cek ini, `next` dari
	// URL yang disusun siapa saja bisa jadi open redirect ke domain luar.
	const rawNext = $page.url.searchParams.get('next');
	const next = rawNext?.startsWith('/') && !rawNext.startsWith('//') ? rawNext : '/app/dashboard';

	// Preselect pilihan lama saat user membuka kembali halaman ini ("Ganti jenis
	// acara" dari dashboard).
	let selectedCategory = $state<EventCategoryId | null>(preference.value?.eventType ?? null);
	let selectedSubType = $state<string | null>(preference.value?.eventSubType ?? null);
	let customLabel = $state<string>(
		preference.value?.eventType === 'other' ? (preference.value?.eventSubType ?? '') : '',
	);
	let submitting = $state(false);
	let error = $state('');

	const category = $derived(
		selectedCategory ? EVENT_CATEGORIES.find((c) => c.id === selectedCategory) : undefined,
	);

	/**
	 * Kategori dengan tepat satu sub-kategori dipilih langsung — drill-down untuk
	 * satu pilihan adalah sampah UI. Kategori dengan banyak sub-kategori
	 * memunculkan pilihan kedua. 'other' memunculkan input teks.
	 */
	function selectCategory(id: EventCategoryId): void {
		selectedCategory = id;
		const found = EVENT_CATEGORIES.find((c) => c.id === id);
		if (found && found.subTypes.length === 1) {
			const only = found.subTypes[0];
			selectedSubType = only ? only.id : null;
		} else {
			selectedSubType = null;
		}
		if (id !== 'other') {
			customLabel = '';
		}
		error = '';
	}

	function selectSubType(id: string): void {
		selectedSubType = id;
		error = '';
	}

	function buildPreference(): UserEventPreference | null {
		if (!selectedCategory) return null;
		if (selectedCategory === 'other') {
			const trimmed = customLabel.trim();
			if (!trimmed) return null;
			return {
				eventType: 'other',
				eventSubType: trimmed,
				selectedAt: new Date().toISOString(),
			};
		}
		// Kategori banyak sub-kategori wajib punya pilihan; kategori satu
		// sub-kategori sudah diisi otomatis di selectCategory.
		if (!selectedSubType) return null;
		return {
			eventType: selectedCategory,
			eventSubType: selectedSubType,
			selectedAt: new Date().toISOString(),
		};
	}

	function handleSubmit(event: SubmitEvent): void {
		event.preventDefault();
		const pref = buildPreference();
		if (!pref) {
			error =
				selectedCategory === 'other'
					? 'Tulis nama jenis acara kamu dulu.'
					: 'Pilih satu jenis acara di bawah ini.';
			return;
		}
		submitting = true;
		preference.set(pref);
		goto(next);
	}

	function handleSkip(): void {
		preference.skip();
		goto(next);
	}
</script>

<svelte:head>
	<title>Pilih Jenis Acara | Ketuk.id</title>
</svelte:head>

<div class="mx-auto max-w-4xl">
	<div>
		<h1 class="font-display text-2xl font-bold text-coffee-900 sm:text-3xl">
			Acara apa yang sedang kamu siapkan?
		</h1>
		<p class="mt-1.5 text-coffee-500">
			Pilih satu — dashboard dan rekomendasi template akan mengikutinya. Bisa diganti kapan saja.
		</p>
	</div>

	<form class="mt-8 flex flex-col gap-8" onsubmit={handleSubmit}>
		<!-- Kategori level pertama. Radio asli (bukan div kustom) supaya keyboard
		     berfungsi tanpa JS tambahan; kartu adalah label yang membungkusnya. -->
		<fieldset class="flex flex-col gap-0">
			<legend class="sr-only">Kategori acara</legend>
			<div class="grid grid-cols-2 gap-3 sm:grid-cols-3 sm:gap-4 lg:grid-cols-4">
				{#each EVENT_CATEGORIES as cat (cat.id)}
					<label
						class="flex min-h-[88px] cursor-pointer flex-col gap-1.5 rounded-2xl border bg-white p-4 transition-all has-[:checked]:border-terracotta-500 has-[:checked]:bg-terracotta-50 has-[:focus-visible]:ring-2 has-[:focus-visible]:ring-terracotta-500 has-[:focus-visible]:ring-offset-2 hover:border-coffee-300
							{selectedCategory === cat.id ? 'border-terracotta-500 bg-terracotta-50' : 'border-coffee-100'}"
					>
						<input
							type="radio"
							name="category"
							value={cat.id}
							class="sr-only"
							checked={selectedCategory === cat.id}
							onchange={() => selectCategory(cat.id)}
						/>
						<span class="text-2xl" aria-hidden="true">{cat.emoji}</span>
						<span class="font-display text-sm font-semibold text-coffee-900">{cat.label}</span>
						<span class="text-xs leading-snug text-coffee-500">{cat.shortDesc}</span>
					</label>
				{/each}
			</div>
		</fieldset>

		{#if category && category.subTypes.length > 1}
			<fieldset class="flex flex-col gap-3">
				<legend class="font-display text-base font-semibold text-coffee-900">
					Jenis {category.label}
				</legend>
				<div class="flex flex-wrap gap-2">
					{#each category.subTypes as sub (sub.id)}
						<label
							class="inline-flex cursor-pointer items-center rounded-full border bg-white px-4 py-2.5 text-sm font-medium transition-all has-[:checked]:border-terracotta-500 has-[:checked]:bg-terracotta-50 has-[:checked]:text-terracotta-700 has-[:focus-visible]:ring-2 has-[:focus-visible]:ring-terracotta-500 has-[:focus-visible]:ring-offset-2 hover:border-coffee-300
								{selectedSubType === sub.id ? 'border-terracotta-500 bg-terracotta-50 text-terracotta-700' : 'border-coffee-100 text-coffee-700'}"
						>
							<input
								type="radio"
								name="subtype"
								value={sub.id}
								class="sr-only"
								checked={selectedSubType === sub.id}
								onchange={() => selectSubType(sub.id)}
							/>
							{sub.label}
						</label>
					{/each}
				</div>
			</fieldset>
		{/if}

		{#if selectedCategory === 'other'}
			<!-- 'other' memakai input teks, bukan daftar. Apa yang user ketik
			     disimpan sebagai eventSubType, dan tetap menjadi satu-satunya
			     sumber kebenaran jenis acaranya — tidak ada menebak/memaksakan
			     kategori tetap. -->
			<Input
				label="Nama/Jenis Acara"
				name="customLabel"
				bind:value={customLabel}
				placeholder="Mis. Syukuran kantor, buka puasa bersama, tunangan adik"
				required
				error={error}
			/>
		{:else if error}
			<p class="text-sm text-red-600" role="alert">{error}</p>
		{/if}

		<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-end">
			<Button type="button" variant="ghost" onclick={handleSkip}>
				Lewati untuk sekarang
			</Button>
			<Button type="submit" loading={submitting}>Simpan & Lanjut</Button>
		</div>
	</form>
</div>
