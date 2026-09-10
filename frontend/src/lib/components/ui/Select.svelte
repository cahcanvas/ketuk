<script module lang="ts">
	export interface SelectOption {
		value: string;
		label: string;
	}

	/**
	 * Opsi yang dikelompokkan, dirender sebagai `<optgroup>`. Dipakai kalau
	 * daftarnya terlalu panjang untuk di-scan sebagai satu daftar rata — mis.
	 * 38 provinsi Indonesia yang dikelompokkan per pulau.
	 */
	export interface SelectOptionGroup {
		label: string;
		options: SelectOption[];
	}
</script>

<script lang="ts">
	interface Props {
		id?: string;
		name?: string;
		label?: string;
		/** Daftar rata. Diabaikan kalau `groups` diisi. */
		options?: SelectOption[];
		groups?: SelectOptionGroup[];
		value?: string;
		placeholder?: string;
		error?: string;
		hint?: string;
		disabled?: boolean;
		required?: boolean;
	}

	let {
		id,
		name,
		label,
		options = [],
		groups,
		value = $bindable(''),
		placeholder = 'Pilih salah satu',
		error,
		hint,
		disabled = false,
		required = false,
	}: Props = $props();

	const inputId = $derived(id ?? name ?? label);
</script>

<div class="flex flex-col gap-1.5">
	{#if label}
		<label for={inputId} class="text-sm font-medium text-navy-800">
			{label}
			{#if required}<span class="text-coral-500">*</span>{/if}
		</label>
	{/if}
	<select
		id={inputId}
		{name}
		{disabled}
		{required}
		bind:value
		aria-invalid={error ? 'true' : undefined}
		aria-describedby={error ? `${inputId}-error` : hint ? `${inputId}-hint` : undefined}
		class="rounded-lg border bg-white px-3.5 py-2.5 text-sm text-navy-900
			focus-visible:outline-2 disabled:cursor-not-allowed disabled:bg-navy-50
			{error ? 'border-red-400' : 'border-navy-200'}"
	>
		<option value="" disabled selected={!value}>{placeholder}</option>
		{#if groups}
			{#each groups as group (group.label)}
				<optgroup label={group.label}>
					{#each group.options as option (option.value)}
						<option value={option.value}>{option.label}</option>
					{/each}
				</optgroup>
			{/each}
		{:else}
			{#each options as option (option.value)}
				<option value={option.value}>{option.label}</option>
			{/each}
		{/if}
	</select>
	{#if error}
		<p id="{inputId}-error" class="text-sm text-red-600">{error}</p>
	{:else if hint}
		<p id="{inputId}-hint" class="text-sm text-navy-400">{hint}</p>
	{/if}
</div>
