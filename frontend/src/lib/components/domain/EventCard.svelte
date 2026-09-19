<script lang="ts">
	import type { Event } from '@ketuk/shared';
	import { EVENT_TYPE_CONFIGS, formatEventDate } from '@ketuk/shared';
	import { getEventIcon } from '$lib/icons';
	import Badge from '../ui/Badge.svelte';

	interface Props {
		event: Event;
	}

	let { event }: Props = $props();

	const config = $derived(EVENT_TYPE_CONFIGS.find((c) => c.value === event.type));
	const Icon = $derived(getEventIcon(event.type));
</script>

<a
	href="/undangan/{event.id}"
	class="group flex flex-col gap-3 rounded-xl border border-coffee-100 bg-white p-5 transition-all hover:-translate-y-0.5 hover:border-coffee-200 hover:shadow-md"
>
	<div class="flex items-start justify-between gap-3">
		<div class="min-w-0 flex-1">
			<span class="inline-flex h-10 w-10 items-center justify-center rounded-lg bg-coffee-50">
				<Icon size={18} class="text-coffee-700" />
			</span>
			<h3 class="mt-3 truncate font-display text-lg font-semibold text-coffee-900">{event.title}</h3>
			{#if config?.label}
				<p class="text-xs text-coffee-400">{config.label}</p>
			{/if}
		</div>
		<Badge tone={event.isPublished ? 'success' : 'neutral'}>
			{event.isPublished ? 'Terbit' : 'Draf'}
		</Badge>
	</div>
	<p class="text-sm text-coffee-500">
		{event.date ? formatEventDate(event.date) : 'Tanggal belum diatur'}
	</p>
	<p class="truncate text-xs text-coffee-400">ketuk.id/{event.slug}</p>
</a>
