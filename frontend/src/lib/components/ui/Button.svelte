<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		variant?: 'primary' | 'secondary' | 'ghost' | 'outline' | 'danger';
		size?: 'sm' | 'md' | 'lg';
		type?: 'button' | 'submit' | 'reset';
		href?: string;
		disabled?: boolean;
		loading?: boolean;
		fullWidth?: boolean;
		children: Snippet;
		onclick?: (event: MouseEvent) => void;
	}

	let {
		variant = 'primary',
		size = 'md',
		type = 'button',
		href,
		disabled = false,
		loading = false,
		fullWidth = false,
		children,
		onclick,
	}: Props = $props();

	const sizeClasses: Record<NonNullable<Props['size']>, string> = {
		sm: 'px-3 py-1.5 text-sm',
		md: 'px-4 py-2.5 text-sm',
		lg: 'px-6 py-3 text-base',
	};

	const variantClasses: Record<NonNullable<Props['variant']>, string> = {
		primary: 'bg-terracotta-500 text-white hover:bg-terracotta-600',
		secondary: 'bg-white border border-coffee-200 text-coffee-900 hover:bg-coffee-50',
		ghost: 'bg-transparent text-coffee-900 hover:bg-coffee-100',
		outline: 'bg-white/10 border border-white/25 text-white hover:bg-white/15 backdrop-blur',
		danger: 'bg-red-600 text-white hover:bg-red-700',
	};

	const classes = $derived(
		[
			'inline-flex items-center justify-center gap-2 rounded-lg font-medium transition-colors',
			'focus-visible:outline-2 disabled:cursor-not-allowed disabled:opacity-50',
			sizeClasses[size],
			variantClasses[variant],
			fullWidth ? 'w-full' : '',
		].join(' '),
	);

	const isDisabled = $derived(disabled || loading);
</script>

{#if href && !isDisabled}
	<a {href} class={classes}>
		{@render children()}
	</a>
{:else}
	<button {type} disabled={isDisabled} class={classes} {onclick}>
		{#if loading}
			<span
				class="h-4 w-4 animate-spin rounded-full border-2 border-current/30 border-t-current"
				aria-hidden="true"
			></span>
		{/if}
		{@render children()}
	</button>
{/if}
