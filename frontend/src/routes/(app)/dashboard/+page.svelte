<script lang="ts">
	import { Mail, ClipboardList, Store, Gift, ArrowRight, type Icon as IconType } from '@lucide/svelte';

	const actions: {
		href: string;
		icon: typeof IconType;
		iconClass: string;
		bgClass: string;
		title: string;
		desc: string;
	}[] = [
		{
			href: '/undangan',
			icon: Mail,
			iconClass: 'text-undangan-600',
			bgClass: 'bg-undangan-100',
			title: 'Kelola Undangan',
			desc: 'Buat undangan baru atau kelola yang sudah ada.',
		},
		{
			href: '/planner',
			icon: ClipboardList,
			iconClass: 'text-planner-600',
			bgClass: 'bg-planner-100',
			title: 'Atur Planner',
			desc: 'Budget, checklist, dan timeline acaramu.',
		},
		{
			href: '/vendor',
			icon: Store,
			iconClass: 'text-vendor-600',
			bgClass: 'bg-vendor-100',
			title: 'Cari Vendor',
			desc: 'Temukan katering, dekorasi, dan vendor lainnya.',
		},
		{
			href: '/hadiah',
			icon: Gift,
			iconClass: 'text-hadiah-600',
			bgClass: 'bg-hadiah-100',
			title: 'Kelola Hadiah',
			desc: 'Lihat hadiah yang masuk untuk acaramu.',
		},
	];
</script>

<svelte:head>
	<title>Dashboard | Ketuk.id</title>
</svelte:head>

<div class="mx-auto max-w-6xl">
	<div>
		<h1 class="font-display text-2xl font-bold text-coffee-900 sm:text-3xl">Mau melakukan apa?</h1>
		<p class="mt-1.5 text-coffee-500">Pilih salah satu untuk mulai.</p>
	</div>

	<div class="mt-8 grid gap-4 sm:gap-5 sm:grid-cols-2 lg:grid-cols-4">
		{#each actions as action (action.href)}
			<a
				href={action.href}
				class="group flex flex-col gap-4 rounded-2xl border border-coffee-100 bg-white p-5 transition-all hover:-translate-y-0.5 hover:border-coffee-200 hover:shadow-md sm:p-6"
			>
				<span
					class="inline-flex h-12 w-12 items-center justify-center rounded-xl {action.bgClass}"
				>
					<action.icon size={22} class={action.iconClass} />
				</span>
				<div class="flex-1">
					<h2 class="font-display text-base font-semibold text-coffee-900 sm:text-lg">
						{action.title}
					</h2>
					<p class="mt-1.5 text-sm text-coffee-500">{action.desc}</p>
				</div>
				<span
					class="inline-flex items-center gap-1 text-sm font-medium text-terracotta-500 opacity-0 transition-opacity group-hover:opacity-100"
				>
					Buka <ArrowRight size={14} />
				</span>
			</a>
		{/each}
	</div>
</div>
