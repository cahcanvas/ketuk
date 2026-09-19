<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { createEvent } from '$lib/api';
	import { Button, Input, Select } from '$lib/components/ui';
	import type { SelectOption } from '$lib/components/ui';
	import { pushToast } from '$lib/stores/toast.svelte';
	import { EVENT_TYPE_CONFIGS, type EventType, generateSlug } from '@ketuk/shared';

	let title = $state('');
	let type = $state('');
	let slug = $state('');
	let slugTouched = $state(false);
	let submitting = $state(false);
	let error = $state('');

	const typeOptions: SelectOption[] = EVENT_TYPE_CONFIGS.map((c) => ({
		value: c.value,
		label: c.label,
	}));

	$effect(() => {
		if (!slugTouched) slug = generateSlug(title);
	});

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (!title.trim() || !type || !slug.trim()) {
			error = 'Judul, jenis acara, dan alamat undangan wajib diisi.';
			return;
		}
		error = '';
		submitting = true;
		try {
			const created = await createEvent(
				{ title, type: type as EventType, slug },
				{ accessToken: $page.data.accessToken },
			);
			pushToast('Undangan berhasil dibuat.', 'success');
			goto(`/undangan/${created.id}`);
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
<p class="mt-1.5 text-coffee-500">Isi info dasar dulu. Tanggal, lokasi, dan detail lain bisa dilengkapi setelahnya.</p>

<form class="mt-8 flex flex-col gap-5" onsubmit={handleSubmit}>
	<Select label="Jenis acara" bind:value={type} options={typeOptions} required />
	<Input label="Judul acara" bind:value={title} placeholder="Pernikahan Budi & Sinta" required />
	<Input
		label="Alamat undangan"
		bind:value={slug}
		oninput={() => (slugTouched = true)}
		hint="ketuk.id/{slug || '...'}"
		required
	/>
	{#if error}
		<p class="text-sm text-red-600">{error}</p>
	{/if}
	<Button type="submit" loading={submitting}>Buat Undangan</Button>
</form>
</div>
