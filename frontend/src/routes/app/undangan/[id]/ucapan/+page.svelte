<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import type { WishItem } from '$lib/api';
	import { confirmDialog } from '$lib/stores/confirm.svelte';
	import { pushToast } from '$lib/stores/toast.svelte';
	import { createSupabaseBrowserClient } from '$lib/supabase/client';
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const supabase = createSupabaseBrowserClient();

	async function handleDelete(wish: WishItem) {
		const confirmed = await confirmDialog({
			title: 'Hapus ucapan?',
			message: `Ucapan dari ${wish.name} akan dihapus permanen.`,
			danger: true,
		});
		if (!confirmed) return;

		const { error } = await supabase.from('wishes').delete().eq('id', wish.id);
		if (error) {
			pushToast('Gagal menghapus ucapan.', 'error');
			return;
		}
		pushToast('Ucapan dihapus.', 'success');
		await invalidateAll();
	}
</script>

<h2 class="font-display text-lg font-semibold text-coffee-900">Moderasi Ucapan</h2>
<p class="mt-1 text-sm text-coffee-500">
	Hapus ucapan yang tidak pantas. Tamu tetap bisa melihat sisanya di halaman undangan.
</p>

<div class="mt-6">
	{#if data.error}
		<p class="rounded-lg bg-red-50 p-4 text-sm text-red-600">{data.error}</p>
	{:else if data.wishes.length === 0}
		<p class="text-sm text-coffee-400">Belum ada ucapan masuk.</p>
	{:else}
		<ul class="flex flex-col gap-3">
			{#each data.wishes as wish (wish.id)}
				<li class="flex items-start justify-between gap-4 rounded-xl border border-coffee-100 bg-white p-4">
					<div>
						<p class="font-medium text-coffee-900">{wish.name}</p>
						<p class="mt-1 text-sm text-coffee-600">{wish.message}</p>
					</div>
					<button
						type="button"
						class="shrink-0 text-sm text-red-500 hover:text-red-700"
						onclick={() => handleDelete(wish)}
					>
						Hapus
					</button>
				</li>
			{/each}
		</ul>
	{/if}
</div>
