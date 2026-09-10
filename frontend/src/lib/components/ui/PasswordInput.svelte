<script lang="ts">
	import { getPasswordStrength } from '@ketuk/shared';
	import { Eye, EyeOff } from '@lucide/svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';

	interface Props {
		id?: string;
		name?: string;
		label?: string;
		placeholder?: string;
		value?: string;
		error?: string;
		hint?: string;
		disabled?: boolean;
		required?: boolean;
		autocomplete?: HTMLInputAttributes['autocomplete'];
		/** Tampilkan indikator kekuatan password — hanya berguna saat MEMBUAT password, bukan saat login. */
		showStrength?: boolean;
	}

	let {
		id,
		name,
		label,
		placeholder,
		value = $bindable(''),
		error,
		hint,
		disabled = false,
		required = false,
		autocomplete,
		showStrength = false,
	}: Props = $props();

	const inputId = $derived(id ?? name ?? label);
	let revealed = $state(false);

	/**
	 * Kelas warna ditulis penuh, bukan disusun dari string (`bg-${x}-500`), karena
	 * Tailwind memindai kode sebagai teks — nama kelas hasil interpolasi tidak akan
	 * ikut ter-generate di CSS build.
	 *
	 * Indeks 0 tidak akan pernah terpakai (indikator disembunyikan saat field masih
	 * kosong) tapi tetap didefinisikan supaya lookup di template tidak butuh cast.
	 */
	const strengthStyles = [
		{ bar: 'w-0', text: 'text-navy-400' },
		{ bar: 'w-1/4 bg-red-500', text: 'text-red-600' },
		{ bar: 'w-2/4 bg-amber-500', text: 'text-amber-600' },
		{ bar: 'w-3/4 bg-lime-500', text: 'text-lime-700' },
		{ bar: 'w-full bg-vendor-500', text: 'text-vendor-600' },
	] as const;

	const strength = $derived(showStrength && value ? getPasswordStrength(value) : null);
	const strengthStyle = $derived(strength ? strengthStyles[strength.score] : null);
</script>

<div class="flex flex-col gap-1.5">
	{#if label}
		<label for={inputId} class="text-sm font-medium text-navy-800">
			{label}
			{#if required}<span class="text-coral-500">*</span>{/if}
		</label>
	{/if}

	<div class="relative">
		<input
			id={inputId}
			{name}
			type={revealed ? 'text' : 'password'}
			{placeholder}
			{disabled}
			{required}
			{autocomplete}
			bind:value
			aria-invalid={error ? 'true' : undefined}
			aria-describedby={error
				? `${inputId}-error`
				: hint
					? `${inputId}-hint`
					: strength
						? `${inputId}-strength`
						: undefined}
			class="w-full rounded-lg border py-2.5 pr-11 pl-3.5 text-sm text-navy-900 placeholder:text-navy-400
				focus-visible:outline-2 disabled:cursor-not-allowed disabled:bg-navy-50
				{error ? 'border-red-400' : 'border-navy-200'}"
		/>
		<button
			type="button"
			{disabled}
			onclick={() => {
				revealed = !revealed;
			}}
			class="absolute inset-y-0 right-0 flex w-11 items-center justify-center rounded-r-lg
				text-navy-400 transition-colors hover:text-navy-700 focus-visible:outline-2
				disabled:cursor-not-allowed"
			aria-label={revealed ? 'Sembunyikan password' : 'Tampilkan password'}
			aria-pressed={revealed}
		>
			{#if revealed}
				<EyeOff size={18} aria-hidden="true" />
			{:else}
				<Eye size={18} aria-hidden="true" />
			{/if}
		</button>
	</div>

	{#if error}
		<p id="{inputId}-error" class="text-sm text-red-600">{error}</p>
	{:else if strength && strengthStyle}
		<!--
			aria-live supaya pembaca layar mengumumkan perubahan kekuatan password saat
			user mengetik. Bar-nya sendiri aria-hidden — informasinya sudah disampaikan
			lewat teks di sebelahnya, tidak perlu diumumkan dua kali.
		-->
		<div id="{inputId}-strength" class="flex items-center gap-2" aria-live="polite">
			<span class="h-1 flex-1 overflow-hidden rounded-full bg-navy-100" aria-hidden="true">
				<span
					class="block h-full rounded-full transition-all duration-200 {strengthStyle.bar}"
				></span>
			</span>
			<span class="text-xs font-medium {strengthStyle.text}">{strength.label}</span>
		</div>
	{:else if hint}
		<p id="{inputId}-hint" class="text-sm text-navy-400">{hint}</p>
	{/if}
</div>
