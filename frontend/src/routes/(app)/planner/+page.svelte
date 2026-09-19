<script lang="ts">
	import { EventCard } from '$lib/components/domain';
	import { Button, EmptyState } from '$lib/components/ui';
	import { formatRupiah } from '@ketuk/shared';
	import { ClipboardList, ArrowRight, Wallet, ListChecks, Clock } from '@lucide/svelte';
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const currentEvent = $derived(data.events.find((e) => e.id === data.eventId));
</script>

<svelte:head>
	<title>Planner | Ketuk.id</title>
</svelte:head>

<div class="mx-auto max-w-5xl">
	{#if !data.eventId}
		<h1 class="font-display text-2xl font-bold text-coffee-900 sm:text-3xl">Planner</h1>
		<p class="mt-1.5 text-coffee-500">Pilih acara yang mau dikelola.</p>
		{#if data.events.length === 0}
			<div class="mt-8">
				<EmptyState
					icon={ClipboardList}
					title="Belum ada undangan"
					description="Buat undangan dulu sebelum mengatur planner."
				>
					{#snippet action()}
						<Button href="/undangan/baru">Buat Undangan</Button>
					{/snippet}
				</EmptyState>
			</div>
		{:else}
			<div class="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{#each data.events as event (event.id)}
					<a href="/planner?event={event.id}" class="block">
						<EventCard {event} />
					</a>
				{/each}
			</div>
		{/if}
	{:else}
		<div class="flex flex-col gap-1">
			<a href="/planner" class="text-sm text-coffee-500 hover:text-coffee-900">
				<span class="inline-flex items-center gap-1.5">Ganti acara</span>
			</a>
			<h1 class="font-display text-2xl font-bold text-coffee-900 sm:text-3xl">
				Planner: {currentEvent?.title ?? ''}
			</h1>
		</div>

		{#if data.error}
			<p class="mt-6 rounded-lg bg-red-50 p-4 text-sm text-red-600">{data.error}</p>
		{:else if data.summary}
			<div class="mt-6 grid gap-4 sm:grid-cols-2">
				<a
					href="/planner/budget?event={data.eventId}"
					class="group flex flex-col gap-3 rounded-2xl border border-coffee-100 bg-white p-6 transition-all hover:-translate-y-0.5 hover:border-coffee-200 hover:shadow-md"
				>
					<span class="inline-flex h-10 w-10 items-center justify-center rounded-lg bg-planner-100">
						<Wallet size={18} class="text-planner-600" />
					</span>
					<div>
						<p class="text-sm text-coffee-500">Budget</p>
						<p class="mt-0.5 font-display text-2xl font-bold text-coffee-900">
							{formatRupiah(data.summary.budget.totalActual)}
						</p>
						<p class="text-sm text-coffee-500">
							dari estimasi {formatRupiah(data.summary.budget.totalEstimated)}
						</p>
					</div>
				</a>
				<a
					href="/planner/checklist?event={data.eventId}"
					class="group flex flex-col gap-3 rounded-2xl border border-coffee-100 bg-white p-6 transition-all hover:-translate-y-0.5 hover:border-coffee-200 hover:shadow-md"
				>
					<span class="inline-flex h-10 w-10 items-center justify-center rounded-lg bg-vendor-100">
						<ListChecks size={18} class="text-vendor-600" />
					</span>
					<div>
						<p class="text-sm text-coffee-500">Checklist</p>
						<p class="mt-0.5 font-display text-2xl font-bold text-coffee-900">
							{data.summary.checklist.done}/{data.summary.checklist.total}
						</p>
						<p class="text-sm text-coffee-500">item selesai</p>
					</div>
				</a>
			</div>
			<div class="mt-4">
				<a
					href="/planner/timeline?event={data.eventId}"
					class="group flex items-center justify-between gap-4 rounded-2xl border border-coffee-100 bg-white p-6 transition-all hover:-translate-y-0.5 hover:border-coffee-200 hover:shadow-md"
				>
					<div class="flex items-center gap-3">
						<span
							class="inline-flex h-10 w-10 items-center justify-center rounded-lg bg-undangan-100"
						>
							<Clock size={18} class="text-undangan-600" />
						</span>
						<p class="font-display font-semibold text-coffee-900">Timeline / Rundown Acara</p>
					</div>
					<ArrowRight size={18} class="text-coffee-400 transition-transform group-hover:translate-x-0.5" />
				</a>
			</div>
		{/if}
	{/if}
</div>
