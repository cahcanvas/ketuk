<script module lang="ts">
	export interface TabItem {
		value: string;
		label: string;
		href?: string;
	}
</script>

<script lang="ts">
	interface Props {
		tabs: TabItem[];
		active: string;
		onchange?: (value: string) => void;
	}

	let { tabs, active, onchange }: Props = $props();
</script>

<div class="flex gap-1 overflow-x-auto border-b border-coffee-100" role="tablist">
	{#each tabs as tab (tab.value)}
		{#if tab.href}
			<a
				href={tab.href}
				role="tab"
				aria-selected={active === tab.value}
				class="whitespace-nowrap border-b-2 px-4 py-2.5 text-sm font-medium transition-colors
					{active === tab.value
					? 'border-terracotta-500 text-coffee-900'
					: 'border-transparent text-coffee-400 hover:text-coffee-700'}"
			>
				{tab.label}
			</a>
		{:else}
			<button
				type="button"
				role="tab"
				aria-selected={active === tab.value}
				onclick={() => onchange?.(tab.value)}
				class="whitespace-nowrap border-b-2 px-4 py-2.5 text-sm font-medium transition-colors
					{active === tab.value
					? 'border-terracotta-500 text-coffee-900'
					: 'border-transparent text-coffee-400 hover:text-coffee-700'}"
			>
				{tab.label}
			</button>
		{/if}
	{/each}
</div>
