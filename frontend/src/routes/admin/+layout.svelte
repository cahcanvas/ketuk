<script lang="ts">
	import { page } from '$app/stores';
	import {
		LayoutDashboard,
		Sparkles,
		CreditCard,
		Users,
		Receipt,
		Store,
		Sliders,
		ExternalLink,
		ShieldCheck,
		Menu,
		X,
		Database,
		Cloud,
		Zap,
		Bell,
		ChevronRight,
		Palette,
	} from '@lucide/svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
	}

	let { children }: Props = $props();
	let mobileOpen = $state(false);

	const navItems = [
		{
			href: '/admin',
			exact: true,
			label: 'Dashboard Overview',
			icon: LayoutDashboard,
			badge: null,
		},
		{
			href: '/admin/produk-undangan',
			exact: false,
			label: 'Master Produk Undangan',
			icon: Sparkles,
			badge: 'Fokus',
			highlight: true,
		},
		{
			href: '/admin/paket',
			exact: false,
			label: 'Paket & Pricing Tiers',
			icon: CreditCard,
			badge: null,
		},
		{
			href: '/admin/pengguna',
			exact: false,
			label: 'Member & Undangan',
			icon: Users,
			badge: null,
		},
		{
			href: '/admin/transaksi',
			exact: false,
			label: 'Transaksi & Hadiah',
			icon: Receipt,
			badge: null,
		},
		{
			href: '/admin/vendor',
			exact: false,
			label: 'Marketplace Vendor',
			icon: Store,
			badge: null,
		},
		{
			href: '/admin/pengaturan',
			exact: false,
			label: 'Pengaturan & Integrasi',
			icon: Sliders,
			badge: null,
		},
	];

	function isCurrentActive(href: string, exact: boolean): boolean {
		if (exact) {
			return $page.url.pathname === href;
		}
		return $page.url.pathname.startsWith(href);
	}
</script>

<svelte:head>
	<title>Ketuk Admin | Master Studio & Control Center</title>
</svelte:head>

<div class="flex min-h-screen bg-navy-950 font-sans text-navy-100 antialiased selection:bg-coral-500 selection:text-white">
	<!-- Desktop Sidebar Navigation -->
	<aside class="hidden w-72 shrink-0 flex-col border-r border-navy-800/80 bg-navy-900/95 backdrop-blur-xl lg:flex">
		<!-- Brand & Environment Header -->
		<div class="p-6 border-b border-navy-800/60">
			<div class="flex items-center justify-between">
				<a href="/admin" class="group flex items-center gap-2.5">
					<div class="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-tr from-coral-600 via-coral-500 to-amber-400 text-white shadow-md shadow-coral-500/20 group-hover:scale-105 transition-transform">
						<Sparkles size={18} />
					</div>
					<div>
						<div class="flex items-center gap-1.5">
							<span class="font-display text-lg font-bold tracking-tight text-white">Ketuk<span class="text-coral-500">.id</span></span>
							<span class="rounded bg-navy-800 px-1.5 py-0.5 text-[9px] font-semibold uppercase tracking-wider text-champagne-300 border border-champagne-400/20">ADMIN</span>
						</div>
						<span class="text-[11px] text-navy-400 font-medium">Enterprise Control Hub</span>
					</div>
				</a>
			</div>

			<!-- System Infrastructure Status Chips -->
			<div class="mt-4 grid grid-cols-3 gap-1.5 rounded-lg bg-navy-950/60 p-2 border border-navy-800/60 text-[10px]">
				<div class="flex items-center gap-1 text-emerald-400" title="Supabase Database: Online">
					<Database size={11} class="shrink-0" />
					<span class="truncate">DB 99.9%</span>
				</div>
				<div class="flex items-center gap-1 text-cyan-400" title="Cloudflare Edge SSR: Active">
					<Cloud size={11} class="shrink-0" />
					<span class="truncate">Edge SSR</span>
				</div>
				<div class="flex items-center gap-1 text-amber-400" title="Duitku Payment Gateway: Ready">
					<Zap size={11} class="shrink-0" />
					<span class="truncate">Duitku</span>
				</div>
			</div>
		</div>

		<!-- Navigation Menu -->
		<div class="flex-1 overflow-y-auto px-4 py-5 scrollbar-thin scrollbar-thumb-navy-700">
			<div class="mb-2 px-3 text-[10px] font-semibold uppercase tracking-widest text-navy-400">
				Manajemen Sistem
			</div>

			<nav class="space-y-1">
				{#each navItems as item (item.href)}
					{@const active = isCurrentActive(item.href, item.exact)}
					<a
						href={item.href}
						class="group relative flex items-center justify-between rounded-xl px-3.5 py-2.5 text-xs font-medium transition-all duration-200 {active
							? 'bg-gradient-to-r from-coral-500/20 via-navy-800/80 to-navy-800 text-white shadow-sm border border-coral-500/30'
							: 'text-navy-300 hover:bg-navy-800/60 hover:text-white border border-transparent'}"
					>
						<div class="flex items-center gap-3">
							<item.icon
								size={17}
								class="transition-colors {active
									? 'text-coral-400'
									: 'text-navy-400 group-hover:text-navy-200'}"
							/>
							<span>{item.label}</span>
						</div>

						<div class="flex items-center gap-1.5">
							{#if item.badge}
								<span
									class="rounded-full px-2 py-0.5 text-[9px] font-bold uppercase tracking-wider {active
										? 'bg-coral-500 text-white shadow-xs'
										: 'bg-navy-800 text-coral-300 border border-coral-500/30'}"
								>
									{item.badge}
								</span>
							{/if}
							{#if active}
								<span class="h-1.5 w-1.5 rounded-full bg-coral-400 animate-pulse"></span>
							{/if}
						</div>
					</a>
				{/each}
			</nav>

			<!-- Quick Helper Note -->
			<div class="mt-8 rounded-xl border border-navy-800 bg-navy-950/70 p-3.5 text-xs">
				<div class="flex items-center gap-2 text-champagne-300 font-medium">
					<Palette size={14} class="text-champagne-400" />
					<span>Master Desain Aktif</span>
				</div>
				<p class="mt-1 text-[11px] text-navy-400 leading-relaxed">
					Konfigurasi tema, animasi, dan matriks hak kustomisasi member akan langsung terhubung ke studio pengantin.
				</p>
			</div>
		</div>

		<!-- Footer User Profile & Quick Actions -->
		<div class="border-t border-navy-800/80 p-4 bg-navy-950/40">
			<div class="flex items-center justify-between mb-3">
				<div class="flex items-center gap-2.5">
					<div class="flex h-8 w-8 items-center justify-center rounded-lg bg-wine-800 text-xs font-bold text-champagne-200 border border-champagne-400/30">
						SA
					</div>
					<div class="overflow-hidden">
						<p class="truncate text-xs font-semibold text-white">Super Administrator</p>
						<p class="truncate text-[10px] text-emerald-400 flex items-center gap-1">
							<ShieldCheck size={10} /> Full Privilege
						</p>
					</div>
				</div>
			</div>

			<a
				href="/template"
				target="_blank"
				class="flex items-center justify-center gap-1.5 rounded-lg border border-navy-700 bg-navy-800/80 px-3 py-2 text-[11px] font-medium text-navy-200 transition-colors hover:bg-navy-700 hover:text-white"
			>
				<span>Lihat Katalog Publik</span>
				<ExternalLink size={12} />
			</a>
		</div>
	</aside>

	<!-- Main Stage Content Area -->
	<div class="flex min-w-0 flex-1 flex-col overflow-x-hidden">
		<!-- Sticky Top Bar Header -->
		<header class="sticky top-0 z-40 flex h-16 items-center justify-between border-b border-navy-800 bg-navy-900/90 px-4 sm:px-6 backdrop-blur-md">
			<div class="flex items-center gap-3">
				<!-- Mobile Menu Toggle Button -->
				<button
					type="button"
					onclick={() => (mobileOpen = !mobileOpen)}
					aria-label={mobileOpen ? 'Tutup navigasi' : 'Buka navigasi'}
					class="inline-flex h-9 w-9 items-center justify-center rounded-lg border border-navy-700 bg-navy-800 text-navy-300 hover:text-white lg:hidden"
				>
					{#if mobileOpen}
						<X size={18} />
					{:else}
						<Menu size={18} />
					{/if}
				</button>

				<!-- Breadcrumb Indicator -->
				<div class="hidden sm:flex items-center gap-2 text-xs text-navy-400">
					<span class="font-medium text-navy-300">Admin Portal</span>
					<ChevronRight size={13} />
					<span class="font-semibold text-white">
						{#if $page.url.pathname === '/admin'}
							Dashboard Overview
						{:else if $page.url.pathname.includes('/admin/produk-undangan')}
							Master Pengelolaan Produk Undangan
						{:else if $page.url.pathname.includes('/admin/paket')}
							Paket & Pricing Tiers
						{:else if $page.url.pathname.includes('/admin/pengguna')}
							Member & Undangan
						{:else if $page.url.pathname.includes('/admin/transaksi')}
							Keuangan & Hadiah
						{:else}
							Pengaturan
						{/if}
					</span>
				</div>
			</div>

			<!-- Right Header Actions -->
			<div class="flex items-center gap-3">
				<!-- Live System Status Pill -->
				<div class="hidden md:flex items-center gap-2 rounded-full border border-emerald-500/30 bg-emerald-950/40 px-3 py-1 text-[11px] font-medium text-emerald-400">
					<span class="h-2 w-2 rounded-full bg-emerald-400 animate-pulse"></span>
					<span>Semua Server Aktif</span>
				</div>

				<a
					href="/dashboard"
					class="hidden sm:inline-flex items-center gap-1 rounded-lg border border-navy-700 bg-navy-800/80 px-3 py-1.5 text-xs font-medium text-navy-200 hover:bg-navy-700 hover:text-white transition-colors"
				>
					<span>Dashboard Member</span>
				</a>

				<button
					type="button"
					class="relative inline-flex h-9 w-9 items-center justify-center rounded-lg border border-navy-700 bg-navy-800 text-navy-300 hover:text-white transition-colors"
					aria-label="Notifikasi"
				>
					<Bell size={16} />
					<span class="absolute top-2 right-2 h-2 w-2 rounded-full bg-coral-500"></span>
				</button>
			</div>
		</header>

		<!-- Mobile Navigation Drawer -->
		{#if mobileOpen}
			<div class="border-b border-navy-800 bg-navy-900 px-4 py-4 lg:hidden">
				<nav class="space-y-1">
					{#each navItems as item (item.href)}
						{@const active = isCurrentActive(item.href, item.exact)}
						<a
							href={item.href}
							onclick={() => (mobileOpen = false)}
							class="flex items-center justify-between rounded-lg px-3 py-2 text-xs font-medium {active
								? 'bg-coral-500 text-white'
								: 'text-navy-300 hover:bg-navy-800 hover:text-white'}"
						>
							<div class="flex items-center gap-2.5">
								<item.icon size={16} />
								<span>{item.label}</span>
							</div>
							{#if item.badge}
								<span class="rounded bg-black/20 px-1.5 py-0.5 text-[9px] font-bold uppercase">
									{item.badge}
								</span>
							{/if}
						</a>
					{/each}
				</nav>
			</div>
		{/if}

		<!-- Page Content Body -->
		<main class="flex-1 p-4 sm:p-6 lg:p-8 max-w-7xl w-full mx-auto">
			{@render children()}
		</main>
	</div>
</div>
