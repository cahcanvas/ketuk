<script lang="ts">
	import type { GiftProduct } from '@ketuk/shared';
	import { formatRupiah } from '@ketuk/shared';
	import { Gift } from '@lucide/svelte';
	import Button from '../ui/Button.svelte';

	interface Props {
		product: GiftProduct;
		onSelect?: (product: GiftProduct) => void;
	}

	let { product, onSelect }: Props = $props();
</script>

<div class="flex flex-col gap-3 rounded-xl border border-coffee-100 bg-white p-4 transition-shadow hover:shadow-md">
	<div class="aspect-square overflow-hidden rounded-lg bg-coffee-50">
		{#if product.images[0]}
			<img
				src={product.images[0]}
				alt={product.name}
				loading="lazy"
				class="h-full w-full object-cover"
			/>
		{:else}
			<div class="flex h-full items-center justify-center" aria-hidden="true">
				<Gift size={40} class="text-coffee-300" strokeWidth={1.5} />
			</div>
		{/if}
	</div>
	<div class="flex-1">
		<h3 class="font-medium text-coffee-900">{product.name}</h3>
		<p class="mt-0.5 text-sm text-coffee-500">{formatRupiah(product.price)}</p>
	</div>
	<Button size="sm" fullWidth onclick={() => onSelect?.(product)}>Kirim Ini</Button>
</div>
