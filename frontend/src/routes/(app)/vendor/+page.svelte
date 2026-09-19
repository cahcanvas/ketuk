<script lang="ts">
	import { goto } from '$app/navigation';
	import { VendorCard } from '$lib/components/domain';
	import { EmptyState, Select } from '$lib/components/ui';
	import type { SelectOption } from '$lib/components/ui';
	import { VENDOR_CATEGORY_CONFIGS } from '@ketuk/shared';
	import { Store } from '@lucide/svelte';
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const categoryOptions: SelectOption[] = VENDOR_CATEGORY_CONFIGS.map((c) => ({
		value: c.value,
		label: c.label,
	}));

	let selectedCategory = $state(data.category ?? '');

	$effect(() => {
		const url = new URL(window.location.href);
		if (selectedCategory) url.searchParams.set('category', selectedCategory);
		else url.searchParams.delete('category');
		const next = `${url.pathname}${url.search}`;
		if (next !== `${window.location.pathname}${window.location.search}`) {
			goto(next, { keepFocus: true, noScroll: true });
		}
	});
</script>

<svelte:head>
	<title>Vendor | Ketuk.id</title>
</svelte:head>

<div class="mx-auto max-w-6xl">
	<h1 class="font-display text-2xl font-bold text-coffee-900 sm:text-3xl">Vendor</h1>
	<p class="mt-1.5 text-coffee-500">Cari katering, dekorasi, fotografer, dan vendor lainnya.</p>

	<div class="mt-6 max-w-xs">
		<Select
			label="Kategori"
			bind:value={selectedCategory}
			options={categoryOptions}
			placeholder="Semua kategori"
		/>
	</div>

	<div class="mt-8">
		{#if data.error}
			<p class="rounded-lg bg-red-50 p-4 text-sm text-red-600">{data.error}</p>
		{:else if data.result.items.length === 0}
			<EmptyState
				icon={Store}
				title="Belum ada vendor"
				description="Coba ganti kategori atau kembali lagi nanti."
			/>
		{:else}
			<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{#each data.result.items as vendor (vendor.id)}
					<VendorCard {vendor} />
				{/each}
			</div>
		{/if}
	</div>
</div>
