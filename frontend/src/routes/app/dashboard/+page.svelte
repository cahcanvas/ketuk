<script lang="ts">
	import {
		Mail,
		ClipboardList,
		Store,
		Gift,
		ArrowRight,
		type Icon as IconType,
	} from '@lucide/svelte';
	import { EventCard } from '$lib/components/domain';
	import { Button } from '$lib/components/ui';
	import { createOnboardingPreference } from '$lib/onboarding/preference.svelte';
	import {
		buildCtaLabel,
		buildDashboardTitle,
		getCategoryLabel,
		getEventConfig,
		type ModuleId,
	} from '$lib/onboarding/event-config';
	import { findCategory, findSubType } from '@ketuk/shared';
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const preference = createOnboardingPreference();
	const config = $derived(getEventConfig(preference.value?.eventType));

	// Sub-kategori yang dipilih (bisa id sub-type atau teks bebas untuk 'other').
	// Dipakai untuk prefill form buat undangan dan menyempurnakan copy.
	const subType = $derived(preference.value?.eventSubType ?? null);
	const knownSubType = $derived(subType ? findSubType(subType) : undefined);

	const title = $derived(buildDashboardTitle(preference.value?.eventType));
	const categoryLabel = $derived(getCategoryLabel(preference.value?.eventType));
	const ctaLabel = $derived(buildCtaLabel(preference.value?.eventType, subType));
	const categoryEmoji = $derived(findCategory(preference.value?.eventType)?.emoji ?? '✨');

	// Prefill hanya untuk sub-type yang dikenali taxonomy. Untuk 'other' teks
	// bebas tidak bisa dipetakan ke select, jadi form diisi manual user.
	const ctaHref = $derived(
		knownSubType ? `/app/undangan/baru?type=${encodeURIComponent(knownSubType.subType.id)}` : '/app/undangan/baru',
	);

	const MODULE_META: Record<ModuleId, { href: string; icon: typeof IconType; title: string }> = {
		undangan: { href: '/app/undangan', icon: Mail, title: 'Kelola Undangan' },
		planner: { href: '/app/planner', icon: ClipboardList, title: 'Atur Planner' },
		vendor: { href: '/app/vendor', icon: Store, title: 'Cari Vendor' },
		hadiah: { href: '/app/hadiah', icon: Gift, title: 'Kelola Hadiah' },
	};
</script>

<svelte:head>
	<title>Dashboard | Ketuk.id</title>
</svelte:head>

<div class="mx-auto max-w-6xl">
	<div class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
		<div>
			<h1 class="font-display text-2xl font-bold text-coffee-900 sm:text-3xl">{title}</h1>
			<p class="mt-1.5 text-coffee-500">
				{preference.value
					? `Konteks: ${categoryLabel}${knownSubType ? ` — ${knownSubType.subType.label}` : subType ? ` — ${subType}` : ''}. Pilih salah satu untuk mulai.`
					: 'Pilih salah satu untuk mulai.'}
			</p>
		</div>
		<Button href="/app/onboarding" variant="secondary" size="sm">
			{preference.value ? 'Ganti jenis acara' : 'Pilih jenis acara'}
		</Button>
	</div>

	{#if data.error}
		<p class="mt-6 rounded-lg bg-red-50 p-4 text-sm text-red-600">{data.error}</p>
	{/if}

	{#if data.events.length > 0}
		<section class="mt-8">
			<h2 class="font-display text-lg font-semibold text-coffee-900">Acara kamu</h2>
			<div class="mt-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{#each data.events as event (event.id)}
					<EventCard {event} />
				{/each}
			</div>
		</section>
	{:else}
		<section class="mt-8">
			<div
				class="flex flex-col items-center gap-4 rounded-2xl border border-dashed border-coffee-200 bg-white p-8 text-center sm:p-10"
			>
				<span class="text-3xl" aria-hidden="true">{categoryEmoji}</span>
				<div>
					<h2 class="font-display text-lg font-semibold text-coffee-900">Belum ada acara</h2>
					<p class="mt-1 text-sm text-coffee-500">
						{preference.value
							? `Buat undangan ${categoryLabel.toLowerCase()} pertama kamu sekarang.`
							: 'Buat undangan pertama kamu sekarang.'}
					</p>
				</div>
				<Button href={ctaHref}>
					{ctaLabel}
					<ArrowRight size={15} />
				</Button>
			</div>
		</section>
	{/if}

	<section class="mt-10">
		<h2 class="font-display text-lg font-semibold text-coffee-900">Modul</h2>
		<div class="mt-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
			{#each config.modules as moduleId (moduleId)}
				{@const meta = MODULE_META[moduleId]}
				<a
					href={meta.href}
					class="group flex flex-col gap-4 rounded-2xl border border-coffee-100 bg-white p-5 transition-all hover:-translate-y-0.5 hover:border-coffee-200 hover:shadow-md sm:p-6"
				>
					<span
						class="inline-flex h-12 w-12 items-center justify-center rounded-xl
							{moduleId === 'undangan'
								? 'bg-undangan-100'
								: moduleId === 'planner'
									? 'bg-planner-100'
									: moduleId === 'vendor'
										? 'bg-vendor-100'
										: 'bg-hadiah-100'}"
					>
						<meta.icon
							size={22}
							class={moduleId === 'undangan'
								? 'text-undangan-600'
								: moduleId === 'planner'
									? 'text-planner-600'
									: moduleId === 'vendor'
										? 'text-vendor-600'
										: 'text-hadiah-600'}
						/>
					</span>
					<div class="flex-1">
						<h3 class="font-display text-base font-semibold text-coffee-900 sm:text-lg">
							{meta.title}
						</h3>
						<p class="mt-1.5 text-sm text-coffee-500">
							{#if moduleId === 'undangan'}
								{config.invitationFeatures.slice(0, 3).join(' · ')}
							{:else if moduleId === 'planner'}
								Budget, checklist, dan timeline acaramu.
							{:else if moduleId === 'vendor'}
								Temukan katering, dekorasi, dan vendor lainnya.
							{:else}
								Lihat hadiah yang masuk untuk acaramu.
							{/if}
						</p>
					</div>
					<span
						class="inline-flex items-center gap-1 text-sm font-medium text-terracotta-500 opacity-0 transition-opacity group-hover:opacity-100"
					>
						Buka <ArrowRight size={14} />
					</span>
				</a>
			{/each}
		</div>

		{#if config.optionalFeatures.length > 0}
			<!-- Fitur sekunder: tersedia, tapi bukan alasan utama kategori ini. Tampil
			     ringan, bukan sebagai kartu sejajar modul utama. -->
			<div class="mt-5 rounded-xl border border-coffee-100 bg-coffee-50/60 p-4">
				<h3 class="text-xs font-semibold uppercase tracking-wider text-coffee-500">
					Fitur lainnya
				</h3>
				<div class="mt-2 flex flex-wrap gap-x-5 gap-y-1.5">
					{#each config.optionalFeatures as feature (feature)}
						<span class="inline-flex items-center gap-1.5 text-sm text-coffee-600">
							{feature}
						</span>
					{/each}
				</div>
			</div>
		{/if}
	</section>
</div>
