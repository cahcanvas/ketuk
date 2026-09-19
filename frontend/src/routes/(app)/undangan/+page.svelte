<script lang="ts">
	import { EventCard } from '$lib/components/domain';
	import { Button, EmptyState } from '$lib/components/ui';
	import { Mail, Plus } from '@lucide/svelte';
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();
</script>

<svelte:head>
	<title>Undangan | Ketuk.id</title>
</svelte:head>

<div class="mx-auto max-w-6xl">
	<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
		<div>
			<h1 class="font-display text-2xl font-bold text-coffee-900 sm:text-3xl">Undangan</h1>
			<p class="mt-1 text-coffee-500">Semua acara yang kamu kelola.</p>
		</div>
		<Button href="/undangan/baru">
			<Plus size={16} />
			Buat Undangan
		</Button>
	</div>

	<div class="mt-8">
		{#if data.error}
			<p class="rounded-lg bg-red-50 p-4 text-sm text-red-600">{data.error}</p>
		{:else if data.events.length === 0}
			<EmptyState icon={Mail} title="Belum ada undangan" description="Buat yang pertama sekarang.">
				{#snippet action()}
					<Button href="/undangan/baru">Buat Undangan Pertama</Button>
				{/snippet}
			</EmptyState>
		{:else}
			<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{#each data.events as event (event.id)}
					<EventCard {event} />
				{/each}
			</div>
		{/if}
	</div>
</div>
