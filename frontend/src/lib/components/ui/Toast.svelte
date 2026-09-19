<script lang="ts">
	import { dismissToast, getToasts } from '$lib/stores/toast.svelte';
	import { X, CheckCircle2, AlertCircle, Info } from '@lucide/svelte';

	const toasts = $derived(getToasts());

	const toneClasses = {
		success: 'bg-vendor-600',
		error: 'bg-red-600',
		info: 'bg-coffee-900',
	} as const;

	const toneIcons = {
		success: CheckCircle2,
		error: AlertCircle,
		info: Info,
	} as const;
</script>

<div class="pointer-events-none fixed inset-x-0 bottom-4 z-[100] flex flex-col items-center gap-2 px-4">
	{#each toasts as toast (toast.id)}
		{@const Icon = toneIcons[toast.type]}
		<div
			class="pointer-events-auto flex items-center gap-3 rounded-lg px-4 py-3 text-sm text-white shadow-lg {toneClasses[
				toast.type
			]}"
			role="status"
		>
			<Icon size={16} class="shrink-0" />
			<span>{toast.message}</span>
			<button
				type="button"
				onclick={() => dismissToast(toast.id)}
				class="inline-flex h-5 w-5 items-center justify-center text-white/70 hover:text-white"
				aria-label="Tutup notifikasi"
			>
				<X size={14} />
			</button>
		</div>
	{/each}
</div>
