<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		tone?: 'neutral' | 'success' | 'warning' | 'danger' | 'undangan' | 'planner' | 'vendor' | 'hadiah';
		children: Snippet;
	}

	let { tone = 'neutral', children }: Props = $props();

	const toneClasses: Record<NonNullable<Props['tone']>, string> = {
		neutral: 'bg-coffee-100 text-coffee-700',
		success: 'bg-vendor-100 text-vendor-600',
		warning: 'bg-hadiah-100 text-hadiah-600',
		danger: 'bg-red-100 text-red-700',
		undangan: 'bg-undangan-100 text-undangan-600',
		planner: 'bg-planner-100 text-planner-600',
		vendor: 'bg-vendor-100 text-vendor-600',
		hadiah: 'bg-hadiah-100 text-hadiah-600',
	};
</script>

<span class="inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium {toneClasses[tone]}">
	{@render children()}
</span>
