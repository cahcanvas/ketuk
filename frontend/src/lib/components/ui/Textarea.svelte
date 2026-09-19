<script lang="ts">
	interface Props {
		id?: string;
		name?: string;
		label?: string;
		placeholder?: string;
		value?: string;
		error?: string;
		hint?: string;
		rows?: number;
		maxlength?: number;
		disabled?: boolean;
		required?: boolean;
	}

	let {
		id,
		name,
		label,
		placeholder,
		value = $bindable(''),
		error,
		hint,
		rows = 4,
		maxlength,
		disabled = false,
		required = false,
	}: Props = $props();

	const inputId = $derived(id ?? name ?? label);
</script>

<div class="flex flex-col gap-1.5">
	{#if label}
		<label for={inputId} class="text-sm font-medium text-coffee-800">
			{label}
			{#if required}<span class="text-terracotta-500">*</span>{/if}
		</label>
	{/if}
	<textarea
		id={inputId}
		{name}
		{placeholder}
		{rows}
		{maxlength}
		{disabled}
		{required}
		bind:value
		aria-invalid={error ? 'true' : undefined}
		class="resize-y rounded-lg border px-3.5 py-2.5 text-sm text-coffee-900 placeholder:text-coffee-400
			focus-visible:outline-2 disabled:cursor-not-allowed disabled:bg-coffee-50
			{error ? 'border-red-400' : 'border-coffee-200'}"
	></textarea>
	{#if maxlength}
		<p class="text-right text-xs text-coffee-400">{value.length}/{maxlength}</p>
	{/if}
	{#if error}
		<p class="text-sm text-red-600">{error}</p>
	{:else if hint}
		<p class="text-sm text-coffee-400">{hint}</p>
	{/if}
</div>
