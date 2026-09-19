<script lang="ts">
	import { getVendorIcon } from '$lib/icons';
	import { formatRupiah, VENDOR_CATEGORY_CONFIGS } from '@ketuk/shared';
	import { ArrowLeft, BadgeCheck, MapPin, Phone, Star } from '@lucide/svelte';
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const config = $derived(VENDOR_CATEGORY_CONFIGS.find((c) => c.value === data.vendor.category));
	const Icon = $derived(getVendorIcon(data.vendor.category));
</script>

<svelte:head>
	<title>{data.vendor.name} | Ketuk.id</title>
</svelte:head>

<div class="mx-auto max-w-3xl">
	<a
		href="/vendor"
		class="inline-flex items-center gap-1.5 text-sm text-coffee-500 hover:text-coffee-900"
	>
		<ArrowLeft size={14} />
		Semua Vendor
	</a>

	<div class="mt-4 flex flex-col gap-6 rounded-2xl border border-coffee-100 bg-white p-6 sm:p-8">
		<div class="flex items-start justify-between gap-4">
			<div class="min-w-0 flex-1">
				<span class="inline-flex h-12 w-12 items-center justify-center rounded-xl bg-vendor-100">
					<Icon size={24} class="text-vendor-600" />
				</span>
				<h1 class="mt-3 font-display text-2xl font-bold text-coffee-900 sm:text-3xl">
					{data.vendor.name}
				</h1>
				<p class="mt-1 flex flex-wrap items-center gap-x-2 gap-y-1 text-coffee-500">
					<span>{config?.label ?? data.vendor.category}</span>
					<span class="text-coffee-300">·</span>
					<span class="inline-flex items-center gap-1">
						<MapPin size={13} />
						{data.vendor.city ?? 'Lokasi belum diatur'}
					</span>
				</p>
			</div>
			{#if data.vendor.isVerified}
				<span
					class="inline-flex shrink-0 items-center gap-1 rounded-full bg-vendor-100 px-2.5 py-1 text-xs font-medium text-vendor-700"
				>
					<BadgeCheck size={12} />
					Terverifikasi
				</span>
			{/if}
		</div>

		{#if data.vendor.description}
			<p class="leading-relaxed text-coffee-700">{data.vendor.description}</p>
		{/if}

		<div class="rounded-xl bg-coffee-50 p-4">
			<p class="text-xs font-medium tracking-wide text-coffee-500 uppercase">Kisaran harga</p>
			<p class="mt-1 font-display text-xl font-semibold text-coffee-900">
				{formatRupiah(data.vendor.priceMin)} – {formatRupiah(data.vendor.priceMax)}
			</p>
		</div>

		<div class="flex flex-wrap items-center gap-4 border-t border-coffee-100 pt-4 text-sm">
			<p class="inline-flex items-center gap-1.5 text-coffee-500">
				<Star size={14} class="fill-amber-400 text-amber-400" />
				{data.vendor.rating.toFixed(1)} ({data.vendor.reviewCount} ulasan)
			</p>
			{#if data.vendor.phone}
				<p class="inline-flex items-center gap-1.5 text-coffee-700">
					<Phone size={14} />
					{data.vendor.phone}
				</p>
			{/if}
		</div>
	</div>
</div>
