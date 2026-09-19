<script lang="ts">
	interface Props {
		width?: string;
		height?: string;
		rounded?: string;
	}

	let { width = '100%', height = '1rem', rounded = '0.5rem' }: Props = $props();
</script>

<span
	class="block animate-pulse bg-coffee-100"
	style="width: {width}; height: {height}; border-radius: {rounded};"
	aria-hidden="true"
></span>
