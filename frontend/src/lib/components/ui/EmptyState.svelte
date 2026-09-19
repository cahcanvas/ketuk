<script lang="ts">
	import { Inbox } from '@lucide/svelte';
	import type { Component, Snippet } from 'svelte';

	interface Props {
		icon?: Component;
		title: string;
		description?: string;
		action?: Snippet;
	}

	let { icon, title, description, action }: Props = $props();

	const Icon = $derived(icon ?? Inbox);
</script>

<div
	class="flex flex-col items-center gap-4 rounded-xl border border-dashed border-coffee-200 bg-white px-6 py-12 text-center"
>
	<span class="inline-flex h-14 w-14 items-center justify-center rounded-full bg-coffee-50">
		<Icon size={26} class="text-coffee-400" />
	</span>
	<div class="flex flex-col gap-1.5">
		<p class="font-display text-lg font-semibold text-coffee-900">{title}</p>
		{#if description}
			<p class="max-w-sm text-sm text-coffee-500">{description}</p>
		{/if}
	</div>
	{#if action}
		<div class="mt-2">{@render action()}</div>
	{/if}
</div>
