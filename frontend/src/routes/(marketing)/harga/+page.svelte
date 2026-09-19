<script lang="ts">
	import { PLANS } from '@ketuk/shared';
	import { PlanCard } from '$lib/components/domain';
</script>

<svelte:head>
	<title>Harga | Ketuk.id</title>
	<meta
		name="description"
		content="Harga Ketuk.id per acara, bukan langganan bulanan. Mulai gratis, upgrade kapan saja."
	/>
</svelte:head>

<section class="py-16 sm:py-24">
	<div class="mx-auto max-w-6xl px-4 sm:px-6">
		<div class="mx-auto max-w-2xl text-center">
			<h1 class="font-display text-3xl font-bold text-coffee-900 sm:text-4xl">
				Harga per acara, bukan langganan
			</h1>
			<p class="mt-3 text-coffee-500">
				Sekali bayar untuk satu acara. Tidak ada biaya bulanan yang jalan terus meski acaranya
				sudah lewat.
			</p>
		</div>
		<div class="mt-10 grid gap-5 sm:mt-12 sm:gap-6 md:grid-cols-3">
			{#each PLANS as plan (plan.id)}
				<PlanCard {plan} highlighted={plan.id === 'pro'} />
			{/each}
		</div>
	</div>
</section>
