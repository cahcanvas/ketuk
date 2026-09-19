<script lang="ts">
	import { submitWish } from '$lib/api';
	import { pushToast } from '$lib/stores/toast.svelte';
	import Button from '../ui/Button.svelte';
	import Input from '../ui/Input.svelte';
	import Textarea from '../ui/Textarea.svelte';

	interface Props {
		eventSlug: string;
		guestSlug?: string;
		onSent?: () => void;
	}

	let { eventSlug, guestSlug, onSent }: Props = $props();

	let name = $state('');
	let message = $state('');
	let submitting = $state(false);
	let error = $state('');

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (!name.trim() || !message.trim()) {
			error = 'Nama dan ucapan wajib diisi.';
			return;
		}
		error = '';
		submitting = true;
		try {
			await submitWish({ eventSlug, guestSlug, name, message });
			pushToast('Ucapan terkirim, terima kasih!', 'success');
			name = '';
			message = '';
			onSent?.();
		} catch {
			error = 'Koneksi terputus. Coba muat ulang.';
		} finally {
			submitting = false;
		}
	}
</script>

<form class="flex flex-col gap-3" onsubmit={handleSubmit}>
	<Input label="Nama" bind:value={name} placeholder="Nama kamu" required />
	<Textarea label="Ucapan" bind:value={message} rows={3} maxlength={500} required />
	{#if error}
		<p class="text-sm text-terracotta-300">{error}</p>
	{/if}
	<Button type="submit" loading={submitting}>Kirim Ucapan</Button>
</form>
