<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		padding?: 'none' | 'sm' | 'md' | 'lg';
		children: Snippet;
	}

	let { padding = 'md', children }: Props = $props();

	const paddingClasses: Record<NonNullable<Props['padding']>, string> = {
		none: '',
		sm: 'p-4',
		md: 'p-6',
		lg: 'p-8',
	};
</script>

<div class="rounded-xl border border-coffee-100 bg-white {paddingClasses[padding]}">
	{@render children()}
</div>
