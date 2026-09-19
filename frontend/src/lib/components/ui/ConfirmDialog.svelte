<script lang="ts">
	import { getConfirmState, resolveConfirm } from '$lib/stores/confirm.svelte';
	import Button from './Button.svelte';

	const state = $derived(getConfirmState());
</script>

{#if state.open}
	<div class="fixed inset-0 z-[110] flex items-center justify-center p-4">
		<button
			type="button"
			class="absolute inset-0 bg-coffee-950/50"
			aria-label="Batal"
			onclick={() => resolveConfirm(false)}
		></button>
		<div role="alertdialog" aria-modal="true" class="relative z-10 w-full max-w-sm rounded-2xl bg-white p-6 shadow-xl">
			<h2 class="mb-2 font-display text-lg font-semibold text-coffee-900">{state.title}</h2>
			<p class="mb-6 text-sm text-coffee-600">{state.message}</p>
			<div class="flex justify-end gap-3">
				<Button variant="ghost" onclick={() => resolveConfirm(false)}>{state.cancelLabel}</Button>
				<Button variant={state.danger ? 'danger' : 'primary'} onclick={() => resolveConfirm(true)}>
					{state.confirmLabel}
				</Button>
			</div>
		</div>
	</div>
{/if}
