<script lang="ts">
	import { Menu, X, ArrowRight, Heart, ShieldCheck, Mail, Phone } from '@lucide/svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
	}

	let { children }: Props = $props();

	let mobileOpen = $state(false);

	const navLinks = [
		{ href: '/#katalog', label: 'Koleksi Desain' },
		{ href: '/#fitur', label: 'Fitur' },
		{ href: '/#cara-kerja', label: 'Cara Kerja' },
		{ href: '/harga', label: 'Harga' },
		{ href: '/#faq', label: 'FAQ' },
		{ href: '/tentang', label: 'Tentang Kami' },
	];
</script>

<div class="flex min-h-screen flex-col bg-cream-50 text-coffee-900 selection:bg-coffee-100 selection:text-coffee-900 font-sans">
	<!-- Top Announcement / Guarantee Bar -->
	<div class="bg-coffee-900 py-2 px-4 text-center text-[10px] sm:text-[11px] font-medium tracking-[0.18em] uppercase text-champagne-200">
		<span class="inline-flex items-center gap-2">
			<span class="h-1.5 w-1.5 rounded-full bg-champagne-400"></span>
			<span>Undangan Aktif Selamanya · Tanpa Biaya Perpanjangan</span>
		</span>
	</div>

	<!-- Header -->
	<header class="sticky top-0 z-40 border-b border-cream-200/80 bg-cream-50/95 backdrop-blur-md transition-colors">
		<div class="mx-auto flex max-w-7xl items-center justify-between px-4 py-3.5 sm:px-6 sm:py-4">
			<!-- Logo -->
			<a href="/" class="group flex items-center gap-2.5">
				<div
					class="flex h-9 w-9 items-center justify-center rounded-full bg-coffee-800 text-champagne-200 shadow-xs transition-transform group-hover:scale-105"
				>
					<span class="font-serif text-lg font-normal">K</span>
				</div>
				<div class="flex flex-col">
					<span class="font-serif text-xl sm:text-2xl font-semibold tracking-wide text-coffee-950">
						Ketuk<span class="text-coffee-700 italic font-normal">.id</span>
					</span>
					<span class="text-[9px] font-medium tracking-[0.2em] uppercase text-coffee-500 -mt-1">
						Digital Invitations
					</span>
				</div>
			</a>

			<!-- Desktop Navigation -->
			<nav class="hidden items-center gap-7 lg:flex">
				{#each navLinks as link (link.href)}
					<a
						href={link.href}
						class="text-xs font-medium uppercase tracking-widest text-coffee-700 transition-colors hover:text-coffee-800"
					>
						{link.label}
					</a>
				{/each}
			</nav>

			<!-- Desktop Actions -->
			<div class="hidden items-center gap-4 lg:flex">
				<a
					href="/masuk"
					class="text-xs font-medium uppercase tracking-widest text-coffee-700 transition-colors hover:text-coffee-800"
				>
					Masuk
				</a>
				<a
					href="/daftar"
					class="inline-flex items-center gap-2 rounded-full bg-coffee-800 px-5 py-2.5 text-xs font-semibold uppercase tracking-wider text-white shadow-xs transition-all hover:bg-coffee-900 active:scale-98 hover:shadow-md"
				>
					<span>Buat Undangan</span>
					<ArrowRight size={13} />
				</a>
			</div>

			<!-- Mobile & Tablet Menu Toggle Button -->
			<button
				type="button"
				class="inline-flex h-11 w-11 items-center justify-center rounded-xl border border-cream-200 bg-cream-100/80 text-coffee-800 transition-colors hover:bg-cream-200/60 active:scale-95 lg:hidden"
				onclick={() => (mobileOpen = !mobileOpen)}
				aria-label={mobileOpen ? 'Tutup menu' : 'Buka menu'}
				aria-expanded={mobileOpen}
			>
				{#if mobileOpen}
					<X size={20} />
				{:else}
					<Menu size={20} />
				{/if}
			</button>
		</div>

		<!-- Mobile / Tablet Navigation Drawer -->
		{#if mobileOpen}
			<div
				class="border-b border-cream-200 bg-cream-50/98 px-5 py-5 shadow-xl backdrop-blur-lg lg:hidden transition-all"
			>
				<nav class="flex flex-col gap-1.5">
					{#each navLinks as link (link.href)}
						<a
							href={link.href}
							class="rounded-xl px-4 py-3 text-sm font-medium tracking-wide text-coffee-800 transition-colors hover:bg-cream-100 hover:text-coffee-800 active:bg-cream-200"
							onclick={() => (mobileOpen = false)}
						>
							{link.label}
						</a>
					{/each}

					<div class="mt-3 flex flex-col gap-2.5 border-t border-cream-200 pt-4">
						<a
							href="/masuk"
							class="rounded-xl border border-cream-300 py-3 text-center text-xs font-semibold uppercase tracking-wider text-coffee-800 hover:bg-cream-100"
							onclick={() => (mobileOpen = false)}
						>
							Masuk ke Akun
						</a>
						<a
							href="/daftar"
							class="flex items-center justify-center gap-2 rounded-xl bg-coffee-800 py-3.5 text-center text-xs font-semibold uppercase tracking-wider text-white shadow-md active:bg-coffee-900"
							onclick={() => (mobileOpen = false)}
						>
							<span>Buat Undangan Sekarang</span>
							<ArrowRight size={14} />
						</a>
					</div>
				</nav>
			</div>
		{/if}
	</header>

	<!-- Page Main Content -->
	<main class="flex-1">
		{@render children()}
	</main>

	<!-- Editorial Luxury Footer -->
	<footer class="border-t border-cream-200 bg-cream-100/90 pt-16 pb-12 text-coffee-700">
		<div class="mx-auto max-w-7xl px-4 sm:px-6">
			<div class="grid gap-10 sm:grid-cols-2 lg:grid-cols-5">
				<!-- Brand Column -->
				<div class="lg:col-span-2 space-y-4">
					<div class="flex items-center gap-2.5">
						<div class="flex h-10 w-10 items-center justify-center rounded-full bg-coffee-800 text-champagne-200 shadow-xs">
							<span class="font-serif text-xl font-normal">K</span>
						</div>
						<span class="font-serif text-2xl font-bold tracking-wide text-coffee-950">
							Ketuk<span class="text-coffee-700 italic font-normal">.id</span>
						</span>
					</div>
				<p class="text-sm text-coffee-600 leading-relaxed max-w-sm">
					Platform undangan digital untuk setiap acara penting — pernikahan, ulang tahun,
					aqiqah, syukuran, wisuda, hingga reuni. RSVP real-time, peta lokasi
					terintegrasi, dan amplop digital dalam satu link.
				</p>
					<div class="flex items-center gap-4 text-xs text-coffee-500 pt-1">
						<span class="inline-flex items-center gap-1">
							<ShieldCheck size={14} class="text-emerald-700" />
							Privasi Terjamin
						</span>
						<span>•</span>
						<span>WhatsApp Direct RSVP</span>
						<span>•</span>
						<span>Aktif Selamanya</span>
					</div>
				</div>

				<!-- Nav Column 1 -->
				<div>
					<h4 class="font-serif text-base font-semibold text-coffee-950 tracking-wide mb-3">
						Jenis Acara
					</h4>
					<ul class="space-y-2.5 text-xs text-coffee-600">
						<li><a href="/#katalog" class="hover:text-coffee-800 transition-colors">Undangan Pernikahan</a></li>
						<li><a href="/#katalog" class="hover:text-coffee-800 transition-colors">Undangan Ulang Tahun</a></li>
						<li><a href="/#katalog" class="hover:text-coffee-800 transition-colors">Aqiqah & Khitanan</a></li>
						<li><a href="/#katalog" class="hover:text-coffee-800 transition-colors">Syukuran & Wisuda</a></li>
						<li><a href="/#katalog" class="hover:text-coffee-800 transition-colors">Reuni & Acara Lainnya</a></li>
					</ul>
				</div>

				<!-- Nav Column 2 -->
				<div>
					<h4 class="font-serif text-base font-semibold text-coffee-950 tracking-wide mb-3">
						Fitur Eksklusif
					</h4>
					<ul class="space-y-2.5 text-xs text-coffee-600">
						<li><a href="/#fitur" class="hover:text-coffee-800 transition-colors">RSVP Otomatis WhatsApp</a></li>
						<li><a href="/#fitur" class="hover:text-coffee-800 transition-colors">Amplop & Kado Digital</a></li>
						<li><a href="/#fitur" class="hover:text-coffee-800 transition-colors">Navigasi Google Maps</a></li>
						<li><a href="/#fitur" class="hover:text-coffee-800 transition-colors">Galeri Foto & Background Musik</a></li>
						<li><a href="/#fitur" class="hover:text-coffee-800 transition-colors">Hitung Mundur Acara</a></li>
					</ul>
				</div>

				<!-- Nav Column 3 -->
				<div>
					<h4 class="font-serif text-base font-semibold text-coffee-950 tracking-wide mb-3">
						Bantuan & Kontak
					</h4>
					<ul class="space-y-2.5 text-xs text-coffee-600">
						<li class="flex items-center gap-1.5">
							<Mail size={13} class="text-champagne-600" />
							<a href="mailto:halo@ketuk.id" class="hover:text-coffee-800 transition-colors">halo@ketuk.id</a>
						</li>
						<li class="flex items-center gap-1.5">
							<Phone size={13} class="text-champagne-600" />
							<span>WhatsApp CS: 0812-8888-9999</span>
						</li>
						<li><a href="/tentang" class="hover:text-coffee-800 transition-colors">Tentang Kami</a></li>
						<li><a href="/#faq" class="hover:text-coffee-800 transition-colors">Pertanyaan Umum (FAQ)</a></li>
					</ul>
				</div>
			</div>

			<!-- Footer Bottom -->
			<div class="mt-12 pt-8 border-t border-cream-200 flex flex-col sm:flex-row items-center justify-between gap-4 text-xs text-coffee-500">
				<p>© {new Date().getFullYear()} Ketuk.id. Hak cipta dilindungi undang-undang.</p>
				<p class="flex items-center gap-1">
					Dibuat dengan rasa cinta <Heart size={12} class="text-coffee-700 fill-coffee-700" /> untuk momen terindah Anda.
				</p>
			</div>
		</div>
	</footer>
</div>
