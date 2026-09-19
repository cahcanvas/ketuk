<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/stores';
	import { createTimelineItem, deleteTimelineItem } from '$lib/api';
	import { Button, EmptyState, Input, Modal } from '$lib/components/ui';
	import { pushToast } from '$lib/stores/toast.svelte';
	import { ArrowLeft, Plus, Clock, X } from '@lucide/svelte';
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	let modalOpen = $state(false);
	let title = $state('');
	let time = $state('08:00');
	let pic = $state('');
	let submitting = $state(false);

	async function handleAdd(event: SubmitEvent) {
		event.preventDefault();
		submitting = true;
		try {
			await createTimelineItem(
				data.eventId,
				{ title, time, pic: pic || undefined },
				{ accessToken: $page.data.accessToken },
			);
			pushToast('Item rundown ditambahkan.', 'success');
			title = '';
			pic = '';
			modalOpen = false;
			await invalidateAll();
		} catch {
			pushToast('Gagal menambah item.', 'error');
		} finally {
			submitting = false;
		}
	}

	async function handleDelete(id: string) {
		try {
			await deleteTimelineItem(data.eventId, id, { accessToken: $page.data.accessToken });
			await invalidateAll();
		} catch {
			pushToast('Gagal menghapus item.', 'error');
		}
	}

	const sortedItems = $derived([...data.items].sort((a, b) => a.time.localeCompare(b.time)));
</script>

<svelte:head>
	<title>Timeline | Ketuk.id</title>
</svelte:head>

<div class="mx-auto max-w-4xl">
	<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
		<div>
			<a
				href="/app/planner?event={data.eventId}"
				class="inline-flex items-center gap-1.5 text-sm text-coffee-500 hover:text-coffee-900"
			>
				<ArrowLeft size={14} />
				Kembali ke Planner
			</a>
			<h1 class="mt-1 font-display text-2xl font-bold text-coffee-900 sm:text-3xl">
				Timeline Acara
			</h1>
		</div>
		<Button onclick={() => (modalOpen = true)}>
			<Plus size={16} />
			Tambah Item
		</Button>
	</div>

	<div class="mt-6">
		{#if data.error}
			<p class="rounded-lg bg-red-50 p-4 text-sm text-red-600">{data.error}</p>
		{:else if sortedItems.length === 0}
			<EmptyState icon={Clock} title="Belum ada rundown" description="Susun urutan acara hari-H." />
		{:else}
			<ol
				class="flex flex-col divide-y divide-coffee-100 overflow-hidden rounded-xl border border-coffee-100 bg-white"
			>
				{#each sortedItems as item (item.id)}
					<li class="flex items-center gap-4 px-4 py-3">
						<span class="w-14 shrink-0 font-mono text-sm font-medium text-terracotta-500">
							{item.time}
						</span>
						<span class="flex-1 text-sm text-coffee-900">{item.title}</span>
						{#if item.pic}
							<span class="hidden text-xs text-coffee-400 sm:inline">{item.pic}</span>
						{/if}
						<button
							type="button"
							class="inline-flex h-7 w-7 items-center justify-center rounded-md text-coffee-300 hover:bg-red-50 hover:text-red-600"
							onclick={() => handleDelete(item.id)}
							aria-label="Hapus item rundown"
						>
							<X size={14} />
						</button>
					</li>
				{/each}
			</ol>
		{/if}
	</div>

	<Modal open={modalOpen} title="Tambah Item Rundown" onclose={() => (modalOpen = false)}>
		<form class="flex flex-col gap-4" onsubmit={handleAdd}>
			<Input label="Judul kegiatan" bind:value={title} required />
			<Input label="Jam" type="time" bind:value={time} required />
			<Input label="Penanggung jawab (opsional)" bind:value={pic} />
			<Button type="submit" loading={submitting}>Tambah</Button>
		</form>
	</Modal>
</div>
