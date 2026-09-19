<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/stores';
	import { publishEvent } from '$lib/api';
	import { Badge, Button } from '$lib/components/ui';
	import { confirmDialog } from '$lib/stores/confirm.svelte';
	import { pushToast } from '$lib/stores/toast.svelte';
	import { formatEventDate } from '@ketuk/shared';
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	let publishing = $state(false);

	async function handlePublish() {
		const confirmed = await confirmDialog({
			title: 'Publikasikan undangan?',
			message: 'Undangan akan bisa diakses publik lewat link-nya begitu dipublikasikan.',
		});
		if (!confirmed) return;

		publishing = true;
		try {
			await publishEvent(data.event.id, { accessToken: $page.data.accessToken });
			pushToast('Undangan dipublikasikan.', 'success');
			await invalidateAll();
		} catch {
			pushToast('Gagal publikasikan. Lengkapi tanggal dan lokasi/venue dulu.', 'error');
		} finally {
			publishing = false;
		}
	}
</script>

<svelte:head>
	<title>{data.event.title} | Ketuk.id</title>
</svelte:head>

<div class="grid gap-6 lg:grid-cols-3">
	<div class="flex flex-col gap-4 rounded-xl border border-coffee-100 bg-white p-6 lg:col-span-2">
		<div class="flex items-center justify-between">
			<h2 class="font-display text-lg font-semibold text-coffee-900">Info Acara</h2>
			<Badge tone={data.event.isPublished ? 'success' : 'neutral'}>
				{data.event.isPublished ? 'Terbit' : 'Draf'}
			</Badge>
		</div>
		<dl class="grid gap-4 text-sm sm:grid-cols-2">
			<div>
				<dt class="text-coffee-400">Tanggal</dt>
				<dd class="text-coffee-900">{data.event.date ? formatEventDate(data.event.date) : 'Belum diatur'}</dd>
			</div>
			<div>
				<dt class="text-coffee-400">Lokasi</dt>
				<dd class="text-coffee-900">{data.event.venue ?? data.event.location ?? 'Belum diatur'}</dd>
			</div>
			<div>
				<dt class="text-coffee-400">Alamat undangan</dt>
				<dd class="text-coffee-900">ketuk.id/{data.event.slug}</dd>
			</div>
			<div>
				<dt class="text-coffee-400">Dilihat</dt>
				<dd class="text-coffee-900">{data.event.viewCount}x</dd>
			</div>
		</dl>
	</div>

	<div class="flex flex-col gap-3 rounded-xl border border-coffee-100 bg-white p-6">
		<h2 class="font-display text-lg font-semibold text-coffee-900">Aksi</h2>
		{#if !data.event.isPublished}
			<Button onclick={handlePublish} loading={publishing} fullWidth>Publikasikan</Button>
		{:else}
			<Button href="/{data.event.slug}" variant="secondary" fullWidth>Lihat Undangan</Button>
		{/if}
		<Button href="/app/undangan/{data.event.id}/edit" variant="ghost" fullWidth>Edit Undangan</Button>
	</div>
</div>
