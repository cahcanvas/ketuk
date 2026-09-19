<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/stores';
	import { createGuest, deleteGuest, importGuestsCsv } from '$lib/api';
	import { GuestTable } from '$lib/components/domain';
	import { Button, Input, Modal, Textarea } from '$lib/components/ui';
	import { confirmDialog } from '$lib/stores/confirm.svelte';
	import { pushToast } from '$lib/stores/toast.svelte';
	import type { Guest } from '@ketuk/shared';
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	let addOpen = $state(false);
	let importOpen = $state(false);
	let name = $state('');
	let phone = $state('');
	let guestGroup = $state('');
	let csvText = $state('');
	let submitting = $state(false);

	const eventId = $derived($page.params.id ?? '');

	const statTotals = $derived(
		data.stats.reduce(
			(acc, row) => {
				acc.total += row.guestCount;
				if (row.status === 'attending') acc.attending = row.guestCount;
				if (row.status === 'not_attending') acc.notAttending = row.guestCount;
				return acc;
			},
			{ total: 0, attending: 0, notAttending: 0 },
		),
	);

	async function handleAdd(event: SubmitEvent) {
		event.preventDefault();
		submitting = true;
		try {
			await createGuest(
				eventId,
				{ name, phone: phone || undefined, guestGroup: guestGroup || undefined },
				{ accessToken: $page.data.accessToken },
			);
			pushToast('Tamu ditambahkan.', 'success');
			name = '';
			phone = '';
			guestGroup = '';
			addOpen = false;
			await invalidateAll();
		} catch {
			pushToast('Gagal menambah tamu.', 'error');
		} finally {
			submitting = false;
		}
	}

	async function handleImport(event: SubmitEvent) {
		event.preventDefault();
		submitting = true;
		try {
			const imported = await importGuestsCsv(eventId, csvText, { accessToken: $page.data.accessToken });
			pushToast(`${imported.length} tamu berhasil diimport.`, 'success');
			csvText = '';
			importOpen = false;
			await invalidateAll();
		} catch {
			pushToast('Gagal import CSV. Pastikan ada kolom "name".', 'error');
		} finally {
			submitting = false;
		}
	}

	async function handleDelete(guest: Guest) {
		const confirmed = await confirmDialog({
			title: 'Hapus tamu?',
			message: `${guest.name} akan dihapus dari daftar tamu.`,
			danger: true,
		});
		if (!confirmed) return;

		try {
			await deleteGuest(eventId, guest.id, { accessToken: $page.data.accessToken });
			pushToast('Tamu dihapus.', 'success');
			await invalidateAll();
		} catch {
			pushToast('Gagal menghapus tamu.', 'error');
		}
	}
</script>

<div class="flex flex-wrap items-center justify-between gap-3">
	<div class="flex gap-6 text-sm">
		<div>
			<p class="text-coffee-400">Total Tamu</p>
			<p class="font-display text-xl font-bold text-coffee-900">{statTotals.total}</p>
		</div>
		<div>
			<p class="text-coffee-400">Hadir</p>
			<p class="font-display text-xl font-bold text-vendor-600">{statTotals.attending}</p>
		</div>
		<div>
			<p class="text-coffee-400">Tidak Hadir</p>
			<p class="font-display text-xl font-bold text-red-500">{statTotals.notAttending}</p>
		</div>
	</div>
	<div class="flex gap-2">
		<Button variant="ghost" onclick={() => (importOpen = true)}>Import CSV</Button>
		<Button onclick={() => (addOpen = true)}>Tambah Tamu</Button>
	</div>
</div>

<div class="mt-6">
	{#if data.error}
		<p class="rounded-lg bg-red-50 p-4 text-sm text-red-600">{data.error}</p>
	{:else}
		<GuestTable guests={data.guests} onDelete={handleDelete} />
	{/if}
</div>

<Modal open={addOpen} title="Tambah Tamu" onclose={() => (addOpen = false)}>
	<form class="flex flex-col gap-4" onsubmit={handleAdd}>
		<Input label="Nama" bind:value={name} required />
		<Input label="Nomor HP (opsional)" bind:value={phone} />
		<Input label="Grup (opsional)" bind:value={guestGroup} placeholder="Keluarga Mempelai Pria" />
		<Button type="submit" loading={submitting}>Tambah</Button>
	</form>
</Modal>

<Modal open={importOpen} title="Import Tamu dari CSV" onclose={() => (importOpen = false)}>
	<form class="flex flex-col gap-4" onsubmit={handleImport}>
		<Textarea
			label="Isi CSV"
			bind:value={csvText}
			rows={8}
			hint={'Baris pertama header: name,phone,guestGroup'}
			required
		/>
		<Button type="submit" loading={submitting}>Import</Button>
	</form>
</Modal>
