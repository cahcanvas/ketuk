<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { createEvent } from '$lib/api';
	import { Button, Input, Select } from '$lib/components/ui';
	import type { SelectOptionGroup } from '$lib/components/ui';
	import { pushToast } from '$lib/stores/toast.svelte';
	import { createOnboardingPreference } from '$lib/onboarding/preference.svelte';
	import { getTitlePlaceholder } from '$lib/onboarding/event-config';
	import {
		EVENT_CATEGORIES,
		findSubType,
		generateSlug,
		type EventType,
	} from '@ketuk/shared';

	const preference = createOnboardingPreference();

	/**
	 * Nilai select adalah id sub-kategori (mis. 'aqiqah'), bukan EventType backend.
	 * Saat submit, id itu diterjemahkan ke EventType lewat findSubType. Cara ini
	 * menjaga payload API tidak berubah, sementara UI memakai taxonomy dua level.
	 */
	let type = $state('');
	let title = $state('');
	let slug = $state('');
	let slugTouched = $state(false);
	let submitting = $state(false);
	let error = $state('');

	const typeGroups: SelectOptionGroup[] = EVENT_CATEGORIES.filter((c) => c.subTypes.length > 0).map(
		(category) => ({
			label: category.label,
			options: category.subTypes.map((sub) => ({ value: sub.id, label: sub.label })),
		}),
	);

	/**
	 * Prefill urutannya: query param `?type=` (dari dashboard/marketing) dulu,
	 * baru preferensi onboarding. Query param diterima apa adanya dari URL,
	 * jadi harus divalidasi — nilai asal ditolak diam-diam, bukan dikirim ke API.
	 */
	function resolveInitialType(): string {
		const param = $page.url.searchParams.get('type');
		if (param) {
			if (findSubType(param)) return param;
			// Kompatibilitas: nilai EventType backend lama (mis. link lama pakai
			// ?type=wedding) tetap diterima dengan mencari sub-type pertama yang
			// memetakan ke EventType itu.
			for (const category of EVENT_CATEGORIES) {
				const found = category.subTypes.find((sub) => sub.eventType === param);
				if (found) return found.id;
			}
		}
		const prefSub = preference.value?.eventSubType;
		if (prefSub && findSubType(prefSub)) return prefSub;
		return '';
	}

	$effect(() => {
		const resolved = resolveInitialType();
		if (resolved && !type) {
			type = resolved;
		}
	});

	const titlePlaceholder = $derived(
		getTitlePlaceholder(preference.value?.eventType, preference.value?.eventSubType),
	);

	$effect(() => {
		if (!slugTouched) slug = generateSlug(title);
	});

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		// EventType backend yang dikirim ke API, dapat dari sub-kategori yang
		// dipilih. Tipe tidak dikenali = tidak pernah terjadi dari UI ini, tapi
		// tetap dijaga agar tidak ada nilai asal yang lolos ke API.
		const resolved = findSubType(type);
		if (!title.trim() || !resolved || !slug.trim()) {
			error = 'Jenis acara, judul, dan alamat undangan wajib diisi.';
			return;
		}
		error = '';
		submitting = true;
		try {
			const created = await createEvent(
				{ title, type: resolved.subType.eventType as EventType, slug },
				{ accessToken: $page.data.accessToken },
			);
			pushToast('Undangan berhasil dibuat.', 'success');
			goto(`/app/undangan/${created.id}`);
		} catch {
			error = 'Gagal membuat undangan. Coba ganti alamat undangan kalau sudah dipakai orang lain.';
		} finally {
			submitting = false;
		}
	}
</script>

<svelte:head>
	<title>Buat Undangan | Ketuk.id</title>
</svelte:head>

<div class="mx-auto max-w-2xl">
	<h1 class="font-display text-2xl font-bold text-coffee-900 sm:text-3xl">Buat undangan baru</h1>
	<p class="mt-1.5 text-coffee-500">
		Isi info dasar dulu. Tanggal, lokasi, dan detail lain bisa dilengkapi setelahnya.
	</p>

	<form class="mt-8 flex flex-col gap-5" onsubmit={handleSubmit}>
		<Select
			label="Jenis acara"
			bind:value={type}
			groups={typeGroups}
			placeholder="Pilih jenis acara"
			required
		/>
		<Input
			label="Judul acara"
			bind:value={title}
			placeholder={titlePlaceholder}
			required
		/>
		<Input
			label="Alamat undangan"
			bind:value={slug}
			oninput={() => (slugTouched = true)}
			hint="ketuk.id/{slug || '...'}"
			required
		/>
		{#if error}
			<p class="text-sm text-red-600" role="alert">{error}</p>
		{/if}
		<Button type="submit" loading={submitting}>Buat Undangan</Button>
	</form>
</div>
