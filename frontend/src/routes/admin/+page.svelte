<script lang="ts">
	import {
		Sparkles,
		ArrowRight,
		CreditCard,
		Users,
		Receipt,
		Store,
		Sliders,
		CheckCircle2,
		Layers,
		Eye,
		Activity,
		Wand2,
	} from '@lucide/svelte';
	import { MASTER_TEMPLATES } from '$lib/data/master-templates';

	const activeTemplatesCount = MASTER_TEMPLATES.filter((t) => t.status === 'published').length;
</script>

<svelte:head>
	<title>Admin Dashboard Overview | Ketuk.id</title>
</svelte:head>

<div class="space-y-8">
	<!-- Top Welcome & System Status Bar -->
	<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-coffee-800/80 pb-6">
		<div>
			<div class="flex items-center gap-2">
				<h1 class="font-display text-2xl sm:text-3xl font-bold tracking-tight text-white">
					Admin Control Center
				</h1>
				<span class="rounded-md bg-terracotta-500/20 px-2 py-0.5 text-xs font-semibold text-terracotta-400 border border-terracotta-500/30">
					v1.2-core
				</span>
			</div>
			<p class="mt-1 text-xs sm:text-sm text-coffee-400">
				Pusat kendali konfigurasi produk, hak kustomisasi member, dan monitoring platform undangan Ketuk.id.
			</p>
		</div>

		<div class="flex items-center gap-3 shrink-0">
			<a
				href="/admin/produk-undangan"
				class="inline-flex items-center gap-2 rounded-xl bg-gradient-to-r from-terracotta-600 to-terracotta-500 px-4 py-2.5 text-xs font-bold uppercase tracking-wider text-white shadow-lg shadow-terracotta-500/25 transition-all hover:scale-102 hover:from-terracotta-500 hover:to-terracotta-400 active:scale-98"
			>
				<Sparkles size={15} />
				<span>Buka Master Produk</span>
				<ArrowRight size={14} />
			</a>
		</div>
	</div>

	<!-- Highlight KPI Metrics -->
	<div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
		<div class="rounded-2xl border border-coffee-800 bg-coffee-900/80 p-5 backdrop-blur-md">
			<div class="flex items-center justify-between">
				<span class="text-xs font-medium text-coffee-400">Master Template Aktif</span>
				<div class="rounded-lg bg-terracotta-500/10 p-2 text-terracotta-400">
					<Layers size={16} />
				</div>
			</div>
			<p class="mt-3 font-display text-2xl sm:text-3xl font-bold text-white">{activeTemplatesCount}</p>
			<p class="mt-1 text-[11px] text-emerald-400 flex items-center gap-1">
				<CheckCircle2 size={12} /> Semua tema siap dipesan
			</p>
		</div>

		<div class="rounded-2xl border border-coffee-800 bg-coffee-900/80 p-5 backdrop-blur-md">
			<div class="flex items-center justify-between">
				<span class="text-xs font-medium text-coffee-400">Total Undangan Pasangan</span>
				<div class="rounded-lg bg-undangan-500/10 p-2 text-undangan-400">
					<Users size={16} />
				</div>
			</div>
			<p class="mt-3 font-display text-2xl sm:text-3xl font-bold text-white">1,842</p>
			<p class="mt-1 text-[11px] text-emerald-400">
				+148 pasangan bulan ini
			</p>
		</div>

		<div class="rounded-2xl border border-coffee-800 bg-coffee-900/80 p-5 backdrop-blur-md">
			<div class="flex items-center justify-between">
				<span class="text-xs font-medium text-coffee-400">Penjualan Paket (GMV)</span>
				<div class="rounded-lg bg-amber-500/10 p-2 text-amber-400">
					<CreditCard size={16} />
				</div>
			</div>
			<p class="mt-3 font-display text-2xl sm:text-3xl font-bold text-white">Rp 142.8M</p>
			<p class="mt-1 text-[11px] text-coffee-400">
				Paket Pro & Lengkap 84%
			</p>
		</div>

		<div class="rounded-2xl border border-coffee-800 bg-coffee-900/80 p-5 backdrop-blur-md">
			<div class="flex items-center justify-between">
				<span class="text-xs font-medium text-coffee-400">Edge Cache Hit Ratio</span>
				<div class="rounded-lg bg-cyan-500/10 p-2 text-cyan-400">
					<Activity size={16} />
				</div>
			</div>
			<p class="mt-3 font-display text-2xl sm:text-3xl font-bold text-white">98.4%</p>
			<p class="mt-1 text-[11px] text-cyan-400">
				Cloudflare Cache Shield
			</p>
		</div>
	</div>

	<!-- Featured Hero: Focus to Master Produk Undangan -->
	<div class="relative overflow-hidden rounded-3xl border border-terracotta-500/40 bg-gradient-to-br from-coffee-900 via-coffee-950/40 to-coffee-900 p-6 sm:p-8 shadow-2xl">
		<div class="absolute -right-16 -top-16 h-64 w-64 rounded-full bg-terracotta-500/10 blur-3xl pointer-events-none"></div>
		<div class="absolute -left-16 -bottom-16 h-64 w-64 rounded-full bg-coffee-600/10 blur-3xl pointer-events-none"></div>

		<div class="relative z-10 max-w-2xl">
			<div class="inline-flex items-center gap-2 rounded-full border border-terracotta-500/30 bg-terracotta-500/10 px-3 py-1 text-xs font-semibold text-terracotta-400">
				<Wand2 size={13} />
				<span>FOKUS PENGEMBANGAN SEKARANG</span>
			</div>

			<h2 class="font-display mt-4 text-2xl sm:text-3xl font-bold text-white leading-snug">
				Master Studio Pengelolaan Desain & Animasi Undangan
			</h2>

			<p class="mt-2 text-xs sm:text-sm text-coffee-300 leading-relaxed">
				Kelola desain visual, palet warna hex, border frame smartphone, stempel lilin <span class="text-champagne-300 font-semibold">wax seal</span>, efek animasi buka amplop (pecah stempel, tirai beludru, kelopak bunga mekar), efek partikel melayang, preset musik, serta <strong>Matriks Batasan Kustomisasi Member per Paket</strong>.
			</p>

			<div class="mt-6 flex flex-wrap items-center gap-3">
				<a
					href="/admin/produk-undangan"
					class="inline-flex items-center gap-2 rounded-xl bg-terracotta-500 px-5 py-3 text-xs font-bold uppercase tracking-wider text-white shadow-lg shadow-terracotta-500/30 transition-all hover:bg-terracotta-600 hover:scale-102 active:scale-98"
				>
					<Sparkles size={16} />
					<span>Masuk ke Studio Master Produk</span>
					<ArrowRight size={15} />
				</a>

				<a
					href="/template"
					target="_blank"
					class="inline-flex items-center gap-2 rounded-xl border border-coffee-700 bg-coffee-800/80 px-4 py-3 text-xs font-medium text-coffee-200 hover:bg-coffee-700 hover:text-white transition-colors"
				>
					<Eye size={15} />
					<span>Lihat Hasil Live Publik</span>
				</a>
			</div>
		</div>
	</div>

	<!-- Complete Architecture Modules Grid -->
	<div>
		<div class="mb-4">
			<h3 class="font-display text-lg font-bold text-white">
				7 Pilar Modul Sistem Admin Ketuk.id
			</h3>
			<p class="text-xs text-coffee-400">
				Struktur lengkap arsitektur portal administrator untuk operasional bisnis SaaS undangan pernikahan.
			</p>
		</div>

		<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
			<!-- Module 1: Master Produk -->
			<div class="group relative flex flex-col justify-between rounded-2xl border border-terracotta-500/40 bg-coffee-900/90 p-5 transition-all hover:border-terracotta-400 hover:shadow-lg hover:shadow-terracotta-500/10">
				<div>
					<div class="flex items-center justify-between">
						<div class="flex h-10 w-10 items-center justify-center rounded-xl bg-terracotta-500/20 text-terracotta-400 border border-terracotta-500/30">
							<Sparkles size={20} />
						</div>
						<span class="rounded-full bg-terracotta-500/20 px-2 py-0.5 text-[10px] font-bold text-terracotta-400 uppercase">
							Fokus Utama
						</span>
					</div>
					<h4 class="font-display mt-3 text-base font-bold text-white">1. Master Produk & Template</h4>
					<p class="mt-1.5 text-xs text-coffee-300 leading-relaxed">
						Input konfigurasi desain visual, tipografi, animasi amplop, efek partikel, dan matriks hak kustomisasi member per paket (Gratis/Pro/Lengkap).
					</p>
				</div>
				<div class="mt-4 pt-3 border-t border-coffee-800 flex items-center justify-between">
					<span class="text-[11px] text-terracotta-400 font-medium">{activeTemplatesCount} Template Aktif</span>
					<a href="/admin/produk-undangan" class="text-xs font-semibold text-terracotta-400 group-hover:text-terracotta-300 flex items-center gap-1">
						Kelola <ArrowRight size={12} />
					</a>
				</div>
			</div>

			<!-- Module 2: Paket & Pricing -->
			<div class="flex flex-col justify-between rounded-2xl border border-coffee-800 bg-coffee-900/60 p-5">
				<div>
					<div class="flex items-center justify-between">
						<div class="flex h-10 w-10 items-center justify-center rounded-xl bg-coffee-800 text-champagne-400 border border-champagne-400/20">
							<CreditCard size={20} />
						</div>
						<span class="rounded-full bg-coffee-800 px-2 py-0.5 text-[10px] font-medium text-coffee-300">
							3 Tiers
						</span>
					</div>
					<h4 class="font-display mt-3 text-base font-bold text-white">2. Paket & Pricing Tiers</h4>
					<p class="mt-1.5 text-xs text-coffee-400 leading-relaxed">
						Aturan harga per acara (Gratis Rp 0, Pro Rp 99.000, Lengkap Rp 249.000), diskon voucher/promo, dan add-on kuota WhatsApp.
					</p>
				</div>
				<div class="mt-4 pt-3 border-t border-coffee-800/80 flex items-center justify-between">
					<span class="text-[11px] text-coffee-400">Gratis • Pro • Lengkap</span>
					<span class="text-xs text-coffee-400">Terkonfigurasi</span>
				</div>
			</div>

			<!-- Module 3: Member & Undangan -->
			<div class="flex flex-col justify-between rounded-2xl border border-coffee-800 bg-coffee-900/60 p-5">
				<div>
					<div class="flex items-center justify-between">
						<div class="flex h-10 w-10 items-center justify-center rounded-xl bg-coffee-800 text-undangan-400 border border-undangan-400/20">
							<Users size={20} />
						</div>
						<span class="rounded-full bg-coffee-800 px-2 py-0.5 text-[10px] font-medium text-coffee-300">
							1,842 Akun
						</span>
					</div>
					<h4 class="font-display mt-3 text-base font-bold text-white">3. Member & Acara</h4>
					<p class="mt-1.5 text-xs text-coffee-400 leading-relaxed">
						Daftar pasangan pengantin, status kepemilikan paket, tanggal acara, slug publik, serta fitur audit & bantuan support teknis.
					</p>
				</div>
				<div class="mt-4 pt-3 border-t border-coffee-800/80 flex items-center justify-between">
					<span class="text-[11px] text-coffee-400">Manajemen Tenant</span>
					<span class="text-xs text-coffee-400">RLS Secured</span>
				</div>
			</div>

			<!-- Module 4: Keuangan & Hadiah -->
			<div class="flex flex-col justify-between rounded-2xl border border-coffee-800 bg-coffee-900/60 p-5">
				<div>
					<div class="flex items-center justify-between">
						<div class="flex h-10 w-10 items-center justify-center rounded-xl bg-coffee-800 text-amber-400 border border-amber-400/20">
							<Receipt size={20} />
						</div>
						<span class="rounded-full bg-coffee-800 px-2 py-0.5 text-[10px] font-medium text-coffee-300">
							Finansial
						</span>
					</div>
					<h4 class="font-display mt-3 text-base font-bold text-white">4. Transaksi & Hadiah</h4>
					<p class="mt-1.5 text-xs text-coffee-400 leading-relaxed">
						Monitoring callback Duitku, verifikasi pembayaran paket, pencatatan amplop digital, dan pencairan saldo hadiah ke pengantin.
					</p>
				</div>
				<div class="mt-4 pt-3 border-t border-coffee-800/80 flex items-center justify-between">
					<span class="text-[11px] text-coffee-400">Duitku Webhook</span>
					<span class="text-xs text-emerald-400">Terhubung</span>
				</div>
			</div>

			<!-- Module 5: Marketplace Vendor -->
			<div class="flex flex-col justify-between rounded-2xl border border-coffee-800 bg-coffee-900/60 p-5">
				<div>
					<div class="flex items-center justify-between">
						<div class="flex h-10 w-10 items-center justify-center rounded-xl bg-coffee-800 text-vendor-500 border border-vendor-500/20">
							<Store size={20} />
						</div>
						<span class="rounded-full bg-coffee-800 px-2 py-0.5 text-[10px] font-medium text-coffee-300">
							Marketplace
						</span>
					</div>
					<h4 class="font-display mt-3 text-base font-bold text-white">5. Marketplace Vendor</h4>
					<p class="mt-1.5 text-xs text-coffee-400 leading-relaxed">
						Approval pendaftaran vendor (fotografer, MUA, dekorasi, catering, venue), verifikasi portofolio, dan bagi hasil komisi platform.
					</p>
				</div>
				<div class="mt-4 pt-3 border-t border-coffee-800/80 flex items-center justify-between">
					<span class="text-[11px] text-coffee-400">Kurasi Vendor</span>
					<span class="text-xs text-coffee-400">Kategori Aktif</span>
				</div>
			</div>

			<!-- Module 6: Pengaturan Sistem -->
			<div class="flex flex-col justify-between rounded-2xl border border-coffee-800 bg-coffee-900/60 p-5">
				<div>
					<div class="flex items-center justify-between">
						<div class="flex h-10 w-10 items-center justify-center rounded-xl bg-coffee-800 text-cyan-400 border border-cyan-400/20">
							<Sliders size={20} />
						</div>
						<span class="rounded-full bg-coffee-800 px-2 py-0.5 text-[10px] font-medium text-coffee-300">
							Sistem
						</span>
					</div>
					<h4 class="font-display mt-3 text-base font-bold text-white">6. Pengaturan & Integrasi</h4>
					<p class="mt-1.5 text-xs text-coffee-400 leading-relaxed">
						Konfigurasi WhatsApp Gateway (Fonnte/Waba), mailer transaksional, edge cache Cloudflare purge, dan pengaturan domain publik.
					</p>
				</div>
				<div class="mt-4 pt-3 border-t border-coffee-800/80 flex items-center justify-between">
					<span class="text-[11px] text-coffee-400">Integrasi API</span>
					<span class="text-xs text-cyan-400">Normal</span>
				</div>
			</div>
		</div>
	</div>
</div>
