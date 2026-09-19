<script lang="ts">
	import { page } from '$app/stores';
</script>

<svelte:head>
	<title>{$page.status} | Ketuk.id</title>
</svelte:head>

<div class="flex min-h-screen flex-col items-center justify-center gap-4 bg-coffee-50 px-6 text-center">
	<p class="font-display text-6xl font-bold text-coffee-200">{$page.status}</p>
	<h1 class="font-display text-2xl font-semibold text-coffee-900">
		{$page.status === 404 ? 'Halaman tidak ditemukan' : 'Terjadi kesalahan'}
	</h1>
	<p class="max-w-sm text-coffee-500">
		{$page.error?.message ?? 'Sesuatu tidak berjalan seperti seharusnya. Coba muat ulang halaman.'}
	</p>
	<a href="/" class="mt-2 text-sm font-medium text-terracotta-500 hover:underline">Kembali ke beranda</a>
</div>
