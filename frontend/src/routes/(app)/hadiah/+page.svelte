<script lang="ts">
	import { Badge, EmptyState } from '$lib/components/ui';
	import { formatRupiah } from '@ketuk/shared';
	import { Gift, ArrowRight } from '@lucide/svelte';
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const statusLabel: Record<string, string> = {
		pending: 'Menunggu Bayar',
		paid: 'Dibayar',
		processing: 'Diproses',
		shipped: 'Dikirim',
		delivered: 'Sampai',
		cancelled: 'Dibatalkan',
	};

	function statusTone(status: string): 'neutral' | 'success' | 'danger' {
		if (status === 'delivered') return 'success';
		if (status === 'cancelled') return 'danger';
		return 'neutral';
	}
</script>

<svelte:head>
	<title>Hadiah | Ketuk.id</title>
</svelte:head>

<div class="mx-auto max-w-4xl">
	<h1 class="font-display text-2xl font-bold text-coffee-900 sm:text-3xl">Hadiah Masuk</h1>
	<p class="mt-1.5 text-coffee-500">Semua hadiah yang dikirim tamu ke acaramu.</p>

	<div class="mt-8">
		{#if data.error}
			<p class="rounded-lg bg-red-50 p-4 text-sm text-red-600">{data.error}</p>
		{:else if data.orders.length === 0}
			<EmptyState
				icon={Gift}
				title="Belum ada hadiah masuk"
				description="Hadiah yang dikirim tamu akan muncul di sini."
			/>
		{:else}
			<ul class="flex flex-col gap-3">
				{#each data.orders as order (order.id)}
					<li>
						<a
							href="/hadiah/{order.id}"
							class="flex items-center justify-between gap-4 rounded-xl border border-coffee-100 bg-white p-4 transition-all hover:border-coffee-200 hover:shadow-sm"
						>
							<div class="min-w-0 flex-1">
								<p class="flex items-center gap-2 truncate font-medium text-coffee-900">
									<span>{order.sender_name}</span>
									<ArrowRight size={14} class="shrink-0 text-coffee-300" />
									<span>{order.recipient_name}</span>
								</p>
								<p class="mt-0.5 text-sm text-coffee-500">
									{formatRupiah(order.total_amount)} · {order.quantity}x
								</p>
							</div>
							<Badge tone={statusTone(order.status)}>
								{statusLabel[order.status] ?? order.status}
							</Badge>
						</a>
					</li>
				{/each}
			</ul>
		{/if}
	</div>
</div>
