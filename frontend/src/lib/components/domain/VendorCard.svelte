<script lang="ts">
	import type { VendorListItem } from '$lib/api';
	import { formatRupiah, VENDOR_CATEGORY_CONFIGS } from '@ketuk/shared';
	import { getVendorIcon } from '$lib/icons';
	import { BadgeCheck, MapPin, Star } from '@lucide/svelte';

	interface Props {
		vendor: VendorListItem;
	}

	let { vendor }: Props = $props();
	const config = $derived(VENDOR_CATEGORY_CONFIGS.find((c) => c.value === vendor.category));
	const Icon = $derived(getVendorIcon(vendor.category));
</script>

<a
	href="/vendor/{vendor.slug}"
	class="group flex flex-col gap-3 rounded-xl border border-coffee-100 bg-white p-5 transition-all hover:-translate-y-0.5 hover:border-coffee-200 hover:shadow-md"
>
	<div class="flex items-start justify-between gap-3">
		<div class="min-w-0 flex-1">
			<span class="inline-flex h-10 w-10 items-center justify-center rounded-lg bg-vendor-100">
				<Icon size={18} class="text-vendor-600" />
			</span>
			<h3 class="mt-3 truncate font-display text-lg font-semibold text-coffee-900">{vendor.name}</h3>
			<p class="mt-0.5 flex items-center gap-1.5 text-sm text-coffee-500">
				<span>{config?.label ?? vendor.category}</span>
				<span class="text-coffee-300">·</span>
				<MapPin size={12} class="shrink-0" />
				<span class="truncate">{vendor.city ?? 'Lokasi belum diatur'}</span>
			</p>
		</div>
		{#if vendor.isVerified}
			<span class="inline-flex shrink-0 items-center gap-1 rounded-full bg-vendor-100 px-2 py-0.5 text-xs font-medium text-vendor-700">
				<BadgeCheck size={12} />
				Terverifikasi
			</span>
		{/if}
	</div>
	<p class="text-sm font-medium text-coffee-700">{formatRupiah(vendor.priceMin)} – {formatRupiah(vendor.priceMax)}</p>
	<p class="flex items-center gap-1 text-xs text-coffee-400">
		<Star size={12} class="fill-amber-400 text-amber-400" />
		{vendor.rating.toFixed(1)} ({vendor.reviewCount} ulasan)
	</p>
</a>
