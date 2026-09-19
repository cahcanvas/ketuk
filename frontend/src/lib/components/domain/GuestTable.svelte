<script lang="ts">
	import type { Guest } from '@ketuk/shared';
	import { Users } from '@lucide/svelte';
	import Badge from '../ui/Badge.svelte';
	import EmptyState from '../ui/EmptyState.svelte';

	interface Props {
		guests: Guest[];
		onEdit?: (guest: Guest) => void;
		onDelete?: (guest: Guest) => void;
	}

	let { guests, onEdit, onDelete }: Props = $props();

	const statusLabel: Record<Guest['rsvpStatus'], string> = {
		pending: 'Menunggu',
		attending: 'Hadir',
		not_attending: 'Tidak Hadir',
	};

	const statusTone: Record<Guest['rsvpStatus'], 'neutral' | 'success' | 'danger'> = {
		pending: 'neutral',
		attending: 'success',
		not_attending: 'danger',
	};
</script>

{#if guests.length === 0}
	<EmptyState
		icon={Users}
		title="Belum ada tamu"
		description="Tambah tamu satu per satu atau import dari CSV."
	/>
{:else}
	<div class="overflow-x-auto rounded-xl border border-coffee-100">
		<table class="w-full min-w-[560px] text-sm">
			<thead class="bg-coffee-50 text-left text-coffee-500">
				<tr>
					<th class="px-4 py-3 font-medium">Nama</th>
					<th class="px-4 py-3 font-medium">Grup</th>
					<th class="px-4 py-3 font-medium">Status</th>
					<th class="px-4 py-3 font-medium">Pax</th>
					<th class="px-4 py-3 font-medium"><span class="sr-only">Aksi</span></th>
				</tr>
			</thead>
			<tbody class="divide-y divide-coffee-100">
				{#each guests as guest (guest.id)}
					<tr>
						<td class="px-4 py-3 font-medium text-coffee-900">{guest.name}</td>
						<td class="px-4 py-3 text-coffee-500">{guest.guestGroup ?? '–'}</td>
						<td class="px-4 py-3">
							<Badge tone={statusTone[guest.rsvpStatus]}>{statusLabel[guest.rsvpStatus]}</Badge>
						</td>
						<td class="px-4 py-3 text-coffee-500">{guest.pax}</td>
						<td class="px-4 py-3 text-right whitespace-nowrap">
							<button type="button" class="text-coffee-400 hover:text-coffee-700" onclick={() => onEdit?.(guest)}>
								Ubah
							</button>
							<button
								type="button"
								class="ml-3 text-coffee-400 hover:text-red-600"
								onclick={() => onDelete?.(guest)}
							>
								Hapus
							</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}
