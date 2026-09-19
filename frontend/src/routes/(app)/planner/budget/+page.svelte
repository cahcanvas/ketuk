<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/stores';
	import { createBudgetItem, deleteBudgetItem, updateBudgetItem } from '$lib/api';
	import { BudgetTable } from '$lib/components/domain';
	import { Button, Input, Modal } from '$lib/components/ui';
	import { confirmDialog } from '$lib/stores/confirm.svelte';
	import { pushToast } from '$lib/stores/toast.svelte';
	import type { BudgetItem } from '@ketuk/shared';
	import { ArrowLeft, Plus } from '@lucide/svelte';
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	let modalOpen = $state(false);
	let editing = $state<BudgetItem | null>(null);
	let category = $state('');
	let name = $state('');
	let estimated = $state('0');
	let actual = $state('');
	let submitting = $state(false);

	function openCreate() {
		editing = null;
		category = '';
		name = '';
		estimated = '0';
		actual = '';
		modalOpen = true;
	}

	function openEdit(item: BudgetItem) {
		editing = item;
		category = item.category;
		name = item.name;
		estimated = String(item.estimated);
		actual = item.actual !== null ? String(item.actual) : '';
		modalOpen = true;
	}

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		submitting = true;
		const payload = {
			category,
			name,
			estimated: Number(estimated) || 0,
			actual: actual ? Number(actual) : null,
		};
		try {
			if (editing) {
				await updateBudgetItem(data.eventId, editing.id, payload, { accessToken: $page.data.accessToken });
			} else {
				await createBudgetItem(data.eventId, payload, { accessToken: $page.data.accessToken });
			}
			pushToast('Budget disimpan.', 'success');
			modalOpen = false;
			await invalidateAll();
		} catch {
			pushToast('Gagal menyimpan budget.', 'error');
		} finally {
			submitting = false;
		}
	}

	async function handleDelete(item: BudgetItem) {
		const confirmed = await confirmDialog({
			title: 'Hapus item budget?',
			message: `"${item.name}" akan dihapus.`,
			danger: true,
		});
		if (!confirmed) return;

		try {
			await deleteBudgetItem(data.eventId, item.id, { accessToken: $page.data.accessToken });
			pushToast('Item budget dihapus.', 'success');
			await invalidateAll();
		} catch {
			pushToast('Gagal menghapus.', 'error');
		}
	}
</script>

<svelte:head>
	<title>Budget | Ketuk.id</title>
</svelte:head>

<div class="mx-auto max-w-5xl">
<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
	<div>
		<a
			href="/planner?event={data.eventId}"
			class="inline-flex items-center gap-1.5 text-sm text-coffee-500 hover:text-coffee-900"
		>
			<ArrowLeft size={14} />
			Kembali ke Planner
		</a>
		<h1 class="mt-1 font-display text-2xl font-bold text-coffee-900 sm:text-3xl">Budget</h1>
	</div>
	<Button onclick={openCreate}>
		<Plus size={16} />
		Tambah Item
	</Button>
</div>

<div class="mt-6">
	{#if data.error}
		<p class="rounded-lg bg-red-50 p-4 text-sm text-red-600">{data.error}</p>
	{:else}
		<BudgetTable items={data.items} onEdit={openEdit} onDelete={handleDelete} />
	{/if}
</div>

<Modal
	open={modalOpen}
	title={editing ? 'Ubah Item Budget' : 'Tambah Item Budget'}
	onclose={() => (modalOpen = false)}
>
	<form class="flex flex-col gap-4" onsubmit={handleSubmit}>
		<Input label="Nama item" bind:value={name} required />
		<Input label="Kategori" bind:value={category} required placeholder="Katering" />
		<div class="grid grid-cols-2 gap-3">
			<Input label="Estimasi (Rp)" type="number" bind:value={estimated} />
			<Input label="Aktual (Rp, opsional)" type="number" bind:value={actual} />
		</div>
		<Button type="submit" loading={submitting}>Simpan</Button>
	</form>
</Modal>
</div>
