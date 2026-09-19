<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
	}

	let { children }: Props = $props();
</script>

<div class="flex min-h-screen items-center justify-center bg-coffee-50 px-4 py-8 sm:px-6">
	<div class="w-full max-w-sm">
		<a
			href="/"
			class="mb-8 block text-center font-display text-2xl font-bold text-coffee-900 transition-opacity hover:opacity-80"
		>
			Ketuk<span class="text-terracotta-500">.id</span>
		</a>
		<div class="rounded-2xl border border-coffee-100 bg-white p-6 shadow-sm sm:p-8">
			{@render children()}
		</div>
	</div>
</div>
