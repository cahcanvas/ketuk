<script lang="ts">
	import { enhance } from '$app/forms';
	import { page } from '$app/stores';
	import { Button, Input, PasswordInput } from '$lib/components/ui';
	import { createSupabaseBrowserClient } from '$lib/supabase/client';
	import { untrack } from 'svelte';
	import type { ActionData } from './$types';

	interface Props {
		form: ActionData;
	}

	let { form }: Props = $props();

	const supabase = createSupabaseBrowserClient();

	// Hanya nilai pertama yang dibutuhkan — lihat catatan yang sama di /daftar.
	let email = $state(untrack(() => form?.values.email) ?? '');
	let password = $state('');
	let submitting = $state(false);

	const next = $derived($page.url.searchParams.get('next') ?? '/dashboard');

	async function handleGoogle() {
		await supabase.auth.signInWithOAuth({
			provider: 'google',
			options: {
				redirectTo: `${window.location.origin}/callback?next=${encodeURIComponent(next)}`,
			},
		});
	}
</script>

<svelte:head>
	<title>Masuk | Ketuk.id</title>
</svelte:head>

<h1 class="font-display text-xl font-semibold text-navy-900 sm:text-2xl">Masuk ke Ketuk.id</h1>
<p class="mt-1.5 text-sm text-navy-500">Kelola undangan dan acaramu.</p>

<form
	method="POST"
	class="mt-6 flex flex-col gap-4"
	use:enhance={() => {
		submitting = true;
		return async ({ update }) => {
			await update({ reset: false });
			submitting = false;
			password = '';
		};
	}}
>
	<Input
		name="email"
		label="Email"
		type="email"
		bind:value={email}
		placeholder="nama@email.com"
		required
		autocomplete="email"
		error={form?.errors.email}
	/>

	<PasswordInput
		name="password"
		label="Password"
		bind:value={password}
		placeholder="Password kamu"
		required
		autocomplete="current-password"
		error={form?.errors.password}
	/>

	{#if form?.formError}
		<p class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-600" role="alert">
			{form.formError}
		</p>
	{/if}

	<Button type="submit" loading={submitting} fullWidth>Masuk</Button>
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
	Lanjutkan dengan Google
</Button>

<p class="mt-6 text-center text-sm text-navy-500">
	Belum punya akun?
	<a href="/daftar" class="font-medium text-coral-500 hover:text-coral-600">Daftar</a>
</p>
