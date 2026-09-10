<script lang="ts">
	import { enhance } from '$app/forms';
	import { Button, Input, PasswordInput, Select } from '$lib/components/ui';
	import type { SelectOptionGroup } from '$lib/components/ui';
	import { createSupabaseBrowserClient } from '$lib/supabase/client';
	import { getProvincesByIsland } from '@ketuk/shared';
	import { CheckCircle2 } from '@lucide/svelte';
	import { untrack } from 'svelte';
	import type { ActionData } from './$types';

	interface Props {
		form: ActionData;
	}

	let { form }: Props = $props();

	const supabase = createSupabaseBrowserClient();

	/**
	 * Nilai awal diambil dari hasil submit sebelumnya kalau ada, supaya validasi
	 * yang gagal tidak menghapus isian yang sudah benar. Password tidak ikut
	 * dikembalikan dari server (lihat +page.server.ts), jadi selalu mulai kosong.
	 *
	 * `untrack` di sini bukan formalitas: yang dibutuhkan memang HANYA nilai
	 * pertama. Setelah komponen hidup, sumber kebenaran isian adalah apa yang
	 * sedang diketik user — kalau `form` ikut dilacak, hasil submit yang gagal
	 * akan menimpa perbaikan yang sudah user ketik ulang.
	 */
	const initial = untrack(() => form?.values);

	let name = $state(initial?.name ?? '');
	let email = $state(initial?.email ?? '');
	let province = $state(initial?.province ?? '');
	let password = $state('');
	let confirmPassword = $state('');
	let submitting = $state(false);

	const provinceGroups: SelectOptionGroup[] = getProvincesByIsland().map((group) => ({
		label: group.island,
		options: group.provinces.map((item) => ({ value: item.code, label: item.label })),
	}));

	async function handleGoogle() {
		await supabase.auth.signInWithOAuth({
			provider: 'google',
			options: {
				redirectTo: `${window.location.origin}/callback?next=${encodeURIComponent('/dashboard')}`,
			},
		});
	}
</script>

<svelte:head>
	<title>Daftar | Ketuk.id</title>
</svelte:head>

{#if form?.emailConfirmationRequired}
	<h1 class="font-display text-xl font-semibold text-navy-900 sm:text-2xl">Cek email kamu</h1>
	<div
		class="mt-6 flex items-start gap-3 rounded-xl border border-vendor-200 bg-vendor-50 p-4 text-sm text-vendor-700"
	>
		<CheckCircle2 size={18} class="mt-0.5 shrink-0 text-vendor-600" />
		<p>
			Akun untuk <strong>{form.values.email}</strong> sudah dibuat. Klik link konfirmasi yang kami
			kirim ke email itu untuk mulai memakainya.
		</p>
	</div>
{:else}
	<h1 class="font-display text-xl font-semibold text-navy-900 sm:text-2xl">Mulai gratis</h1>
	<p class="mt-1.5 text-sm text-navy-500">
		Isi lima kolom di bawah, langsung bisa dipakai. Tidak perlu kartu kredit.
	</p>

	<!--
		Tanpa `use:enhance` form ini tetap berfungsi (submit biasa ke server action).
		`enhance` cuma menambah: tidak ada reload halaman, dan tombol bisa
		menampilkan status loading.
	-->
	<form
		method="POST"
		class="mt-6 flex flex-col gap-4"
		use:enhance={() => {
			submitting = true;
			return async ({ update }) => {
				// `reset: false` supaya isian tidak dikosongkan saat validasi gagal —
				// perilaku default SvelteKit mengosongkan form setelah setiap submit.
				await update({ reset: false });
				submitting = false;
				// Password sudah tidak dipercaya lagi setelah satu kali submit gagal
				// (server tidak mengembalikannya), jadi dikosongkan agar konsisten
				// dengan apa yang akan terkirim di percobaan berikutnya.
				password = '';
				confirmPassword = '';
			};
		}}
	>
		<Input
			name="name"
			label="Nama lengkap"
			bind:value={name}
			placeholder="Budi Santoso"
			required
			autocomplete="name"
			error={form?.errors.name}
			hint="Nama asli kamu — ini yang muncul sebagai penyelenggara di undangan."
		/>

		<Input
			name="email"
			label="Email"
			type="email"
			bind:value={email}
			placeholder="nama@email.com"
			required
			autocomplete="email"
			error={form?.errors.email}
			hint="Pakai email utama yang aktif. Email sekali-pakai tidak bisa dipakai."
		/>

		<Select
			name="province"
			label="Asal"
			bind:value={province}
			groups={provinceGroups}
			placeholder="Pilih provinsi"
			required
			error={form?.errors.province}
		/>

		<PasswordInput
			name="password"
			label="Password"
			bind:value={password}
			placeholder="Minimal 8 karakter"
			required
			autocomplete="new-password"
			showStrength
			error={form?.errors.password}
			hint="Minimal 8 karakter, gabungan huruf dan angka."
		/>

		<PasswordInput
			name="confirmPassword"
			label="Konfirmasi password"
			bind:value={confirmPassword}
			placeholder="Ulangi password"
			required
			autocomplete="new-password"
			error={form?.errors.confirmPassword}
		/>

		{#if form?.formError}
			<p class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-600" role="alert">
				{form.formError}
			</p>
		{/if}

		<Button type="submit" loading={submitting} fullWidth>Daftar Gratis</Button>
	</form>

	<div class="my-6 flex items-center gap-3 text-xs text-navy-400">
		<span class="h-px flex-1 bg-navy-100"></span>
		atau
		<span class="h-px flex-1 bg-navy-100"></span>
	</div>

	<Button variant="secondary" fullWidth onclick={handleGoogle}>
		<svg
			width="18"
			height="18"
			viewBox="0 0 24 24"
			xmlns="http://www.w3.org/2000/svg"
			aria-hidden="true"
		>
			<path
				fill="#4285F4"
				d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
			/>
			<path
				fill="#34A853"
				d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
			/>
			<path
				fill="#FBBC05"
				d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
			/>
			<path
				fill="#EA4335"
				d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
			/>
		</svg>
		Daftar dengan Google
	</Button>
{/if}

<p class="mt-6 text-center text-sm text-navy-500">
	Sudah punya akun?
	<a href="/masuk" class="font-medium text-coral-500 hover:text-coral-600">Masuk</a>
</p>
