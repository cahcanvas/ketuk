<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { createSupabaseBrowserClient } from '$lib/supabase/client';
	import { createOnboardingPreference } from '$lib/onboarding/preference.svelte';
	import {
		Home,
		Mail,
		ClipboardList,
		Store,
		Gift,
		Menu,
		X,
		LogOut,
		type Icon as IconType,
	} from '@lucide/svelte';
	import type { Snippet } from 'svelte';
	import type { LayoutData } from './$types';

	interface Props {
		children: Snippet;
		data: LayoutData;
	}

	let { children, data }: Props = $props();

	let mobileOpen = $state(false);

	const navItems: { href: string; label: string; icon: typeof IconType }[] = [
		{ href: '/app/dashboard', label: 'Beranda', icon: Home },
		{ href: '/app/undangan', label: 'Undangan', icon: Mail },
		{ href: '/app/planner', label: 'Planner', icon: ClipboardList },
		{ href: '/app/vendor', label: 'Vendor', icon: Store },
		{ href: '/app/hadiah', label: 'Hadiah', icon: Gift },
	];

	const supabase = createSupabaseBrowserClient();
	const preference = createOnboardingPreference();

	function isActive(href: string): boolean {
		return $page.url.pathname === href || $page.url.pathname.startsWith(`${href}/`);
	}

	/**
	 * Gate onboarding: user yang belum memilih kategori acara (dan belum pernah
	 * skip, dan belum punya event sama sekali) diarahkan ke /app/onboarding.
	 *
	 * Berjalan di klien karena localStorage tidak terbaca di server. Path asli
	 * dibawa sebagai `?next=` supaya setelah onboarding user kembali ke halaman
	 * yang dituju — mis. CTA template yang mengarah ke form buat undangan.
	 *
	 * Ditaruh di $effect (bukan saat init) supaya re-evaluasi setiap navigasi:
	 * setelah user memilih, preference terisi dan gate berhenti memicu.
	 */
	const needsOnboarding = $derived(
		browser &&
			$page.url.pathname !== '/app/onboarding' &&
			!preference.value &&
			!preference.isSkipped &&
			!data.hasEvents,
	);

	$effect(() => {
		if (needsOnboarding) {
			const target = `${$page.url.pathname}${$page.url.search}`;
			goto(`/app/onboarding?next=${encodeURIComponent(target)}`);
		}
	});

	async function handleLogout() {
		await supabase.auth.signOut();
		goto('/');
	}
</script>

<div class="flex min-h-screen bg-coffee-50">
	<!-- Desktop sidebar -->
	<aside class="hidden w-64 shrink-0 flex-col border-r border-coffee-100 bg-white p-5 lg:flex">
		<a href="/app/dashboard" class="font-display text-xl font-bold text-coffee-900">
			Ketuk<span class="text-terracotta-500">.id</span>
		</a>
		<nav class="mt-8 flex flex-col gap-1">
			{#each navItems as item (item.href)}
				<a
					href={item.href}
					class="flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors
						{isActive(item.href)
						? 'bg-terracotta-100 text-terracotta-600'
						: 'text-coffee-600 hover:bg-coffee-50 hover:text-coffee-900'}"
				>
					<item.icon size={18} />
					{item.label}
				</a>
			{/each}
		</nav>
		<div class="mt-auto flex flex-col gap-3 border-t border-coffee-100 pt-4">
			<p class="truncate text-xs text-coffee-400" title={data.user.email}>{data.user.email}</p>
			<button
				type="button"
				onclick={handleLogout}
				class="inline-flex items-center gap-2 rounded-lg px-2 py-1.5 text-left text-sm font-medium text-coffee-500 transition-colors hover:bg-coffee-50 hover:text-coffee-900"
			>
				<LogOut size={16} />
				Keluar
			</button>
		</div>
	</aside>

	<div class="flex min-w-0 flex-1 flex-col">
		<!-- Mobile header -->
		<header
			class="sticky top-0 z-30 flex items-center justify-between border-b border-coffee-100 bg-white px-4 py-3 lg:hidden"
		>
			<a href="/app/dashboard" class="font-display text-lg font-bold text-coffee-900">
				Ketuk<span class="text-terracotta-500">.id</span>
			</a>
			<button
				type="button"
				onclick={() => (mobileOpen = !mobileOpen)}
				aria-label={mobileOpen ? 'Tutup menu' : 'Buka menu'}
				aria-expanded={mobileOpen}
				class="inline-flex h-10 w-10 items-center justify-center rounded-lg text-coffee-700 hover:bg-coffee-50"
			>
				{#if mobileOpen}
					<X size={22} />
				{:else}
					<Menu size={22} />
				{/if}
			</button>
		</header>
		{#if mobileOpen}
			<nav
				class="flex flex-col gap-1 border-b border-coffee-100 bg-white p-3 shadow-sm lg:hidden"
			>
				{#each navItems as item (item.href)}
					<a
						href={item.href}
						class="flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium
							{isActive(item.href) ? 'bg-terracotta-100 text-terracotta-600' : 'text-coffee-600 hover:bg-coffee-50'}"
						onclick={() => (mobileOpen = false)}
					>
						<item.icon size={18} />
						{item.label}
					</a>
				{/each}
				<div class="mt-2 border-t border-coffee-100 pt-2">
					<p class="truncate px-3 py-1 text-xs text-coffee-400">{data.user.email}</p>
					<button
						type="button"
						onclick={handleLogout}
						class="flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-left text-sm font-medium text-coffee-500 hover:bg-coffee-50"
					>
						<LogOut size={16} />
						Keluar
					</button>
				</div>
			</nav>
		{/if}
	<main class="flex-1 p-4 sm:p-6 lg:p-8">
		{#if needsOnboarding}
			<!-- Gate onboarding sedang mengarahkan ke /app/onboarding. Halaman tujuan
			     sengaja tidak dirender dulu supaya load function-nya tidak jalan. -->
			<div class="flex min-h-[40vh] items-center justify-center" aria-busy="true">
				<div class="text-sm text-coffee-500">Mempersiapkan halaman…</div>
			</div>
		{:else}
			{@render children()}
		{/if}
	</main>
	</div>
</div>
