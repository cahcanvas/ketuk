<script lang="ts">
	import { submitRsvp } from '$lib/api';
	import { pushToast } from '$lib/stores/toast.svelte';
	import Button from '../ui/Button.svelte';
	import Input from '../ui/Input.svelte';
	import Textarea from '../ui/Textarea.svelte';

	interface Props {
		eventSlug: string;
		guestSlug: string;
		defaultName?: string;
	}

	let { eventSlug, guestSlug, defaultName = '' }: Props = $props();

	let name = $state(defaultName);
	let status = $state<'attending' | 'not_attending' | null>(null);
	let pax = $state('1');
	let message = $state('');
	let submitting = $state(false);
	let submitted = $state(false);
	let error = $state('');

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (!status) {
			error = 'Pilih kehadiran dulu ya.';
			return;
		}
		error = '';
		submitting = true;
		try {
			await submitRsvp({
				eventSlug,
				guestSlug,
				name,
				status,
				pax: Number(pax) || 1,
				message: message || undefined,
			});
			submitted = true;
			pushToast('RSVP kamu sudah terkirim. Terima kasih!', 'success');
		} catch {
			error = 'Koneksi terputus. Coba muat ulang.';
		} finally {
			submitting = false;
		}
	}
</script>

{#if submitted}
	<div class="rounded-xl bg-white/10 p-6 text-center">
		<p class="text-lg font-medium text-white">Terima kasih, {name}!</p>
		<p class="mt-1 text-sm text-white/70">
			{status === 'attending'
				? 'Kami tunggu kehadiranmu.'
				: 'Terima kasih sudah memberi tahu kami.'}
		</p>
		<button type="button" class="mt-4 text-sm text-white underline underline-offset-4" onclick={() => (submitted = false)}>
			Ubah jawaban
		</button>
	</div>
{:else}
	<form class="flex flex-col gap-4" onsubmit={handleSubmit}>
		<Input label="Nama kamu" bind:value={name} required placeholder="Nama lengkap" />

		<div class="flex flex-col gap-1.5">
			<span class="text-sm font-medium text-white">Konfirmasi kehadiran</span>
			<div class="flex gap-3">
				<button
					type="button"
					onclick={() => (status = 'attending')}
					class="flex-1 rounded-lg border px-4 py-2.5 text-sm font-medium transition-colors
						{status === 'attending'
						? 'border-terracotta-400 bg-terracotta-500 text-white'
						: 'border-white/20 text-white hover:border-white/40'}"
				>
					Hadir
				</button>
				<button
					type="button"
					onclick={() => (status = 'not_attending')}
					class="flex-1 rounded-lg border px-4 py-2.5 text-sm font-medium transition-colors
						{status === 'not_attending'
						? 'border-terracotta-400 bg-terracotta-500 text-white'
						: 'border-white/20 text-white hover:border-white/40'}"
				>
					Tidak Bisa Hadir
				</button>
			</div>
		</div>

		{#if status === 'attending'}
			<Input label="Jumlah orang yang hadir" type="number" bind:value={pax} />
		{/if}

		<Textarea label="Pesan (opsional)" bind:value={message} rows={3} maxlength={500} />

		{#if error}
			<p class="text-sm text-terracotta-300">{error}</p>
		{/if}

		<Button type="submit" loading={submitting} fullWidth>Kirim RSVP</Button>
	</form>
{/if}
