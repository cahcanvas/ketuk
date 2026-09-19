<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/stores';
	import { updateEvent } from '$lib/api';
	import { Button, Input, Textarea } from '$lib/components/ui';
	import { pushToast } from '$lib/stores/toast.svelte';
	import { createSupabaseBrowserClient } from '$lib/supabase/client';
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const supabase = createSupabaseBrowserClient();
	const isWeddingLike = $derived(data.event.type === 'wedding' || data.event.type === 'engagement');
	const invitation = data.event.invitation;

	// Detail acara — disimpan lewat backend (tabel events).
	let date = $state(data.event.date ?? '');
	let timeStart = $state(data.event.timeStart ?? '');
	let timeEnd = $state(data.event.timeEnd ?? '');
	let venue = $state(data.event.venue ?? '');
	let location = $state(data.event.location ?? '');
	let locationUrl = $state(data.event.locationUrl ?? '');
	let savingEvent = $state(false);

	// Isi undangan — disimpan langsung lewat Supabase client (RLS invitations_*_own
	// sudah mengizinkan owner kelola penuh, tidak butuh endpoint backend khusus).
	let openingText = $state(invitation?.openingText ?? '');
	let closingText = $state(invitation?.closingText ?? '');
	let groomName = $state(invitation && 'groomName' in invitation ? (invitation.groomName ?? '') : '');
	let brideName = $state(invitation && 'brideName' in invitation ? (invitation.brideName ?? '') : '');
	let hostName = $state(invitation && 'hostName' in invitation ? (invitation.hostName ?? '') : '');
	let savingInvitation = $state(false);

	async function handleSaveEvent(event: SubmitEvent) {
		event.preventDefault();
		savingEvent = true;
		try {
			await updateEvent(
				data.event.id,
				{
					date: date || undefined,
					timeStart: timeStart || undefined,
					timeEnd: timeEnd || undefined,
					venue: venue || undefined,
					location: location || undefined,
					locationUrl: locationUrl || undefined,
				},
				{ accessToken: $page.data.accessToken },
			);
			pushToast('Detail acara disimpan.', 'success');
			await invalidateAll();
		} catch {
			pushToast('Gagal menyimpan. Koneksi terputus, coba lagi.', 'error');
		} finally {
			savingEvent = false;
		}
	}

	async function handleSaveInvitation(event: SubmitEvent) {
		event.preventDefault();
		savingInvitation = true;

		const payload: Record<string, unknown> = {
			event_id: data.event.id,
			opening_text: openingText || null,
			closing_text: closingText || null,
		};
		if (isWeddingLike) {
			payload.groom_name = groomName || null;
			payload.bride_name = brideName || null;
		} else {
			payload.host_name = hostName || null;
		}

		const { error } = await supabase.from('invitations').upsert(payload, { onConflict: 'event_id' });
		savingInvitation = false;

		if (error) {
			pushToast('Gagal menyimpan undangan.', 'error');
			return;
		}
		pushToast('Undangan disimpan.', 'success');
		await invalidateAll();
	}
</script>

<svelte:head>
	<title>Edit {data.event.title} | Ketuk.id</title>
</svelte:head>

<div class="flex flex-col gap-8">
	<form
		class="flex flex-col gap-4 rounded-xl border border-coffee-100 bg-white p-6"
		onsubmit={handleSaveEvent}
	>
		<h2 class="font-display text-lg font-semibold text-coffee-900">Waktu &amp; Tempat</h2>
		<div class="grid gap-4 sm:grid-cols-2">
			<Input label="Tanggal" type="date" bind:value={date} />
			<div class="grid grid-cols-2 gap-3">
				<Input label="Jam mulai" type="time" bind:value={timeStart} />
				<Input label="Jam selesai" type="time" bind:value={timeEnd} />
			</div>
			<Input label="Nama venue" bind:value={venue} placeholder="Gedung Serbaguna Grand Puri" />
			<Input label="Alamat lengkap" bind:value={location} placeholder="Jl. Merdeka No. 1, Jakarta" />
		</div>
		<Input label="Link Google Maps" bind:value={locationUrl} placeholder="https://maps.google.com/..." />
		<div>
			<Button type="submit" loading={savingEvent}>Simpan Waktu &amp; Tempat</Button>
		</div>
	</form>

	<form
		class="flex flex-col gap-4 rounded-xl border border-coffee-100 bg-white p-6"
		onsubmit={handleSaveInvitation}
	>
		<h2 class="font-display text-lg font-semibold text-coffee-900">Isi Undangan</h2>
		{#if isWeddingLike}
			<div class="grid gap-4 sm:grid-cols-2">
				<Input label="Nama mempelai pria" bind:value={groomName} />
				<Input label="Nama mempelai wanita" bind:value={brideName} />
			</div>
		{:else}
			<Input label="Nama tuan rumah / penyelenggara" bind:value={hostName} />
		{/if}
		<Textarea label="Kata pembuka" bind:value={openingText} rows={3} />
		<Textarea label="Kata penutup" bind:value={closingText} rows={3} />
		<div>
			<Button type="submit" loading={savingInvitation}>Simpan Isi Undangan</Button>
		</div>
	</form>
</div>
