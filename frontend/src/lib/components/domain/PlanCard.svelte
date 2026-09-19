<script lang="ts">
	import type { Plan } from '@ketuk/shared';
	import { formatRupiah } from '@ketuk/shared';
	import { Check, Minus } from '@lucide/svelte';
	import Button from '../ui/Button.svelte';

	interface Props {
		plan: Plan;
		highlighted?: boolean;
	}

	let { plan, highlighted = false }: Props = $props();
</script>

<div
	class="flex flex-col gap-6 rounded-2xl border p-6 sm:p-8
		{highlighted ? 'border-terracotta-500 bg-coffee-900 text-white shadow-xl' : 'border-coffee-100 bg-white'}"
>
	<div>
		<h3 class="font-display text-xl font-bold">{plan.name}</h3>
		<p class="mt-1 text-sm {highlighted ? 'text-white/70' : 'text-coffee-500'}">{plan.description}</p>
	</div>
	<p class="font-display text-3xl font-bold">
		{plan.price === 0 ? 'Gratis' : formatRupiah(plan.price)}
		{#if plan.price > 0}
			<span class="text-sm font-normal {highlighted ? 'text-white/60' : 'text-coffee-400'}">/acara</span>
		{/if}
	</p>
	<ul class="flex flex-1 flex-col gap-2.5 text-sm">
		{#each plan.features as feature (feature.label)}
			<li
				class="flex items-start gap-2 {feature.included
					? ''
					: highlighted
						? 'text-white/40'
						: 'text-coffee-300'}"
			>
				<span class="mt-0.5 shrink-0" aria-hidden="true">
					{#if feature.included}
						<Check size={16} class={highlighted ? 'text-terracotta-400' : 'text-vendor-600'} />
					{:else}
						<Minus size={16} />
					{/if}
				</span>
				{feature.label}
			</li>
		{/each}
	</ul>
	<Button variant={highlighted ? 'primary' : 'secondary'} href="/daftar" fullWidth>
		{plan.price === 0 ? 'Mulai Gratis' : 'Pilih Paket'}
	</Button>
</div>
