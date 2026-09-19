<script lang="ts">
	import type { EventWithInvitation, WishItem } from '$lib/api';
	import { createGiftOrder, listWishes } from '$lib/api';
	import { downloadIcsFile, generateIcs } from '$lib/utils/ics';
	import { EVENT_TYPE_CONFIGS, formatEventDate, formatEventTime, formatRupiah, PAYMENT_METHODS } from '@ketuk/shared';
	import type { GiftProduct } from '@ketuk/shared';
	import { getEventIcon } from '$lib/icons';
	import { Volume2, VolumeX } from '@lucide/svelte';
	import Button from '../ui/Button.svelte';
	import Input from '../ui/Input.svelte';
	import Modal from '../ui/Modal.svelte';
	import Select from '../ui/Select.svelte';
	import type { SelectOption } from '../ui/Select.svelte';
	import Textarea from '../ui/Textarea.svelte';
	import CountdownTimer from './CountdownTimer.svelte';
	import ImageGallery from './ImageGallery.svelte';
	import RsvpForm from './RsvpForm.svelte';
	import WishForm from './WishForm.svelte';
	import WishList from './WishList.svelte';

	interface BankAccountInfo {
		bank: string;
		accountNumber: string;
		accountName: string;
	}

	interface Props {
		event: EventWithInvitation;
		giftProducts?: GiftProduct[];
		guestSlug?: string;
		guestName?: string;
	}

	let { event, giftProducts = [], guestSlug, guestName }: Props = $props();

	const invitation = $derived(event.invitation);
	const isWeddingLike = $derived(event.type === 'wedding' || event.type === 'engagement');
	const config = $derived(EVENT_TYPE_CONFIGS.find((c) => c.value === event.type));
	const EventIcon = $derived(getEventIcon(event.type));

	const displayNames = $derived.by(() => {
		if (invitation) {
			if (isWeddingLike && 'groomName' in invitation && invitation.groomName && invitation.brideName) {
				return `${invitation.brideName} & ${invitation.groomName}`;
			}
			if ('hostName' in invitation && invitation.hostName) return invitation.hostName;
		}
		return event.title;
	});

	const bankAccount = $derived.by((): BankAccountInfo | null => {
		const raw = invitation?.customData?.['bankAccount'];
		if (
			raw &&
			typeof raw === 'object' &&
			'bank' in raw &&
			'accountNumber' in raw &&
			'accountName' in raw
		) {
			return raw as BankAccountInfo;
		}
		return null;
	});

	let coverOpen = $state(true);
	let audioEl: HTMLAudioElement | null = $state(null);
	let musicPlaying = $state(false);

	function openInvitation() {
		coverOpen = false;
		if (invitation?.musicUrl && audioEl) {
			audioEl
				.play()
				.then(() => {
					musicPlaying = true;
				})
				.catch(() => {});
		}
	}

	function toggleMusic() {
		if (!audioEl) return;
		if (musicPlaying) {
			audioEl.pause();
		} else {
			audioEl.play().catch(() => {});
		}
		musicPlaying = !musicPlaying;
	}

	function handleAddToCalendar() {
		if (!event.date) return;
		const start = event.timeStart ? `${event.date}T${event.timeStart}:00` : `${event.date}T00:00:00`;
		const ics = generateIcs({
			title: event.title,
			location: event.venue ?? event.location ?? undefined,
			start,
		});
		downloadIcsFile(`${event.slug}.ics`, ics);
	}

	let copiedAccount = $state(false);
	function copyAccountNumber(number: string) {
		navigator.clipboard
			.writeText(number)
			.then(() => {
				copiedAccount = true;
				setTimeout(() => (copiedAccount = false), 2000);
			})
			.catch(() => {});
	}

	// Ucapan diambil client-side setelah halaman ter-render, bukan di SSR —
	// supaya HTML halaman tetap satu bentuk untuk semua orang dan bisa di-cache bersama.
	let wishes = $state<WishItem[]>([]);
	let wishesLoading = $state(true);

	async function loadWishes() {
		wishesLoading = true;
		try {
			wishes = await listWishes(event.slug);
		} catch {
			wishes = [];
		} finally {
			wishesLoading = false;
		}
	}

	$effect(() => {
		loadWishes();
	});

	function handleWishSent() {
		loadWishes();
	}

	// Kirim hadiah
	let orderModalProduct = $state<GiftProduct | null>(null);
	let orderResult = $state<{ paymentUrl: string | null; vaNumber: string | null; qrString: string | null } | null>(
		null,
	);
	let senderName = $state('');
	let senderPhone = $state('');
	let recipientName = $state(displayNames);
	let recipientAddress = $state(event.venue ?? event.location ?? '');
	let quantity = $state('1');
	let orderMessage = $state('');
	let paymentMethod = $state('');
	let ordering = $state(false);
	let orderError = $state('');

	const paymentMethodOptions: SelectOption[] = PAYMENT_METHODS.map((m) => ({ value: m.code, label: m.label }));

	function openOrderModal(product: GiftProduct) {
		orderModalProduct = product;
		orderResult = null;
		orderError = '';
	}

	async function handleOrderSubmit(submitEvent: SubmitEvent) {
		submitEvent.preventDefault();
		if (!orderModalProduct || !paymentMethod) {
			orderError = 'Pilih metode pembayaran dulu ya.';
			return;
		}
		orderError = '';
		ordering = true;
		try {
			const result = await createGiftOrder(event.slug, {
				productId: orderModalProduct.id,
				senderName,
				senderPhone,
				recipientName,
				recipientAddress,
				quantity: Number(quantity) || 1,
				message: orderMessage || undefined,
				paymentMethod,
			});
			orderResult = {
				paymentUrl: result.payment.paymentUrl,
				vaNumber: result.payment.vaNumber,
				qrString: result.payment.qrString,
			};
		} catch {
			orderError = 'Gagal membuat pesanan. Koneksi terputus, coba lagi.';
		} finally {
			ordering = false;
		}
	}
</script>

{#if invitation?.musicUrl}
	<audio bind:this={audioEl} src={invitation.musicUrl} loop></audio>
{/if}

{#if coverOpen}
	<div class="fixed inset-0 z-50 flex flex-col items-center justify-center gap-6 bg-coffee-900 px-6 text-center text-white">
		<p class="inline-flex items-center gap-2 text-sm text-white/60">
			<EventIcon size={16} />
			{config?.label ?? ''}
		</p>
		<h1 class="font-display text-3xl font-bold sm:text-4xl">{displayNames}</h1>
		{#if event.date}
			<p class="text-white/70">{formatEventDate(event.date)}</p>
		{/if}
		{#if guestName}
			<p class="mt-4 text-sm text-white/60">Kepada Yth.</p>
			<p class="font-display text-xl font-semibold">{guestName}</p>
		{/if}
		<Button size="lg" onclick={openInvitation}>Buka Undangan</Button>
	</div>
{:else}
	<div class="bg-coffee-900 text-white">
		<!-- Kontrol musik, selalu terlihat -->
		{#if invitation?.musicUrl}
			<button
				type="button"
				onclick={toggleMusic}
				class="fixed top-4 right-4 z-40 flex h-10 w-10 items-center justify-center rounded-full bg-white/10 backdrop-blur"
				aria-label={musicPlaying ? 'Matikan musik' : 'Putar musik'}
			>
				{#if musicPlaying}
					<Volume2 size={18} />
				{:else}
					<VolumeX size={18} />
				{/if}
			</button>
		{/if}

		<!-- Hero -->
		<section class="flex flex-col items-center justify-center gap-3 px-6 py-24 text-center">
			<p class="inline-flex items-center gap-2 text-sm text-white/60">
				<EventIcon size={16} />
				{config?.label ?? ''}
			</p>
			<h1 class="font-display text-4xl font-bold sm:text-5xl">{displayNames}</h1>
			{#if event.date}
				<p class="text-lg text-white/70">{formatEventDate(event.date)}</p>
			{/if}
		</section>

		<!-- Profil -->
		{#if invitation}
			<section class="px-6 py-16">
				<div class="mx-auto max-w-xl text-center">
					{#if isWeddingLike && 'groomName' in invitation}
						<div class="grid gap-8 sm:grid-cols-2">
							<div>
								<h2 class="font-display text-2xl font-semibold">{invitation.brideName}</h2>
								{#if invitation.brideParents}
									<p class="mt-1 text-sm text-white/60">Putri dari {invitation.brideParents}</p>
								{/if}
							</div>
							<div>
								<h2 class="font-display text-2xl font-semibold">{invitation.groomName}</h2>
								{#if invitation.groomParents}
									<p class="mt-1 text-sm text-white/60">Putra dari {invitation.groomParents}</p>
								{/if}
							</div>
						</div>
					{:else if 'hostName' in invitation && invitation.hostName}
						<h2 class="font-display text-2xl font-semibold">{invitation.hostName}</h2>
					{/if}
					{#if invitation.openingText}
						<p class="mt-8 text-white/80">{invitation.openingText}</p>
					{/if}
				</div>
			</section>
		{/if}

		<!-- Waktu & Tempat -->
		<section class="bg-white/5 px-6 py-16">
			<div class="mx-auto flex max-w-md flex-col items-center gap-4 text-center">
				<h2 class="font-display text-2xl font-semibold">Waktu &amp; Tempat</h2>
				{#if event.date}
					<p>{formatEventDate(event.date)}</p>
				{/if}
				{#if event.timeStart}
					<p class="text-white/70">
						{formatEventTime(event.timeStart)}{event.timeEnd ? ` – ${formatEventTime(event.timeEnd)}` : ''}
					</p>
				{/if}
				{#if event.venue}
					<p class="font-medium">{event.venue}</p>
				{/if}
				{#if event.location}
					<p class="text-sm text-white/70">{event.location}</p>
				{/if}
				<div class="mt-2 flex flex-wrap justify-center gap-3">
					{#if event.locationUrl}
						<Button href={event.locationUrl} variant="secondary">Buka Google Maps</Button>
					{/if}
					{#if event.date}
						<Button variant="ghost" onclick={handleAddToCalendar}>Tambah ke Kalender</Button>
					{/if}
				</div>
			</div>
		</section>

		<!-- Countdown -->
		{#if event.date}
			<section class="px-6 py-16 text-center">
				<h2 class="font-display text-xl font-semibold text-white/80">Menuju Hari Bahagia</h2>
				<div class="mt-6">
					<CountdownTimer date={event.timeStart ? `${event.date}T${event.timeStart}:00` : event.date} />
				</div>
			</section>
		{/if}

		<!-- Galeri -->
		{#if invitation && invitation.gallery.length > 0}
			<section class="bg-white/5 px-6 py-16">
				<div class="mx-auto max-w-3xl">
					<h2 class="mb-6 text-center font-display text-2xl font-semibold">Galeri</h2>
					<ImageGallery images={invitation.gallery} alt={displayNames} />
				</div>
			</section>
		{/if}

		<!-- Cerita -->
		{#if isWeddingLike && invitation?.loveStory && invitation.loveStory.length > 0}
			<section class="px-6 py-16">
				<div class="mx-auto flex max-w-xl flex-col gap-8">
					<h2 class="text-center font-display text-2xl font-semibold">Kisah Kami</h2>
					{#each invitation.loveStory as chapter (chapter.title)}
						<div>
							<h3 class="font-display font-semibold text-terracotta-300">{chapter.title}</h3>
							<p class="mt-1 text-white/80">{chapter.description}</p>
						</div>
					{/each}
				</div>
			</section>
		{/if}

		<!-- RSVP -->
		<section class="bg-white/5 px-6 py-16">
			<div class="mx-auto max-w-md">
				<h2 class="mb-6 text-center font-display text-2xl font-semibold">RSVP</h2>
				{#if guestSlug}
					<RsvpForm eventSlug={event.slug} {guestSlug} defaultName={guestName ?? ''} />
				{:else}
					<p class="text-center text-sm text-white/70">
						RSVP dilakukan lewat link undangan personal yang dikirim host. Cek pesan undangan yang kamu
						terima.
					</p>
				{/if}
			</div>
		</section>

		<!-- Ucapan -->
		<section class="px-6 py-16">
			<div class="mx-auto max-w-md">
				<h2 class="mb-6 text-center font-display text-2xl font-semibold">Ucapan</h2>
				<WishForm eventSlug={event.slug} {guestSlug} onSent={handleWishSent} />
				<div class="mt-8">
					<WishList {wishes} loading={wishesLoading} />
				</div>
			</div>
		</section>

		<!-- Kirim Hadiah -->
		<section class="bg-white/5 px-6 py-16">
			<div class="mx-auto max-w-2xl">
				<h2 class="mb-6 text-center font-display text-2xl font-semibold">Kirim Hadiah</h2>

				{#if bankAccount}
					<div class="mb-8 rounded-xl border border-white/10 p-5 text-center">
						<p class="text-sm text-white/60">Amplop Digital</p>
						<p class="mt-1 font-medium">{bankAccount.bank} · {bankAccount.accountName}</p>
						<div class="mt-2 flex items-center justify-center gap-2">
							<span class="font-mono text-lg">{bankAccount.accountNumber}</span>
							<button
								type="button"
								onclick={() => bankAccount && copyAccountNumber(bankAccount.accountNumber)}
								class="rounded-lg bg-white/10 px-3 py-1 text-xs hover:bg-white/20"
							>
								{copiedAccount ? 'Tersalin!' : 'Salin'}
							</button>
						</div>
					</div>
				{/if}

				{#if giftProducts.length > 0}
					<div class="grid gap-4 sm:grid-cols-2">
						{#each giftProducts as product (product.id)}
							<div class="rounded-xl border border-white/10 p-4">
								<p class="font-medium">{product.name}</p>
								<p class="mt-1 text-sm text-white/70">{formatRupiah(product.price)}</p>
								<Button size="sm" onclick={() => openOrderModal(product)}>Kirim Ini</Button>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</section>

		<!-- Footer -->
		<footer class="px-6 py-10 text-center text-xs text-white/40">
			Dibuat dengan <a href="https://ketuk.id" class="underline">Ketuk.id</a>
		</footer>
	</div>
{/if}

<Modal
	open={orderModalProduct !== null}
	title={orderModalProduct?.name}
	onclose={() => (orderModalProduct = null)}
>
	{#if orderResult}
		<div class="flex flex-col gap-3 text-center">
			<p class="font-medium text-coffee-900">Pesanan dibuat! Selesaikan pembayaran:</p>
			{#if orderResult.paymentUrl}
				<Button href={orderResult.paymentUrl}>Bayar Sekarang</Button>
			{:else if orderResult.vaNumber}
				<p class="text-sm text-coffee-500">Nomor Virtual Account</p>
				<p class="font-mono text-xl text-coffee-900">{orderResult.vaNumber}</p>
			{:else if orderResult.qrString}
				<p class="text-sm text-coffee-500">Scan QRIS di aplikasi pembayaranmu.</p>
			{/if}
		</div>
	{:else}
		<form class="flex flex-col gap-4" onsubmit={handleOrderSubmit}>
			<Input label="Nama kamu" bind:value={senderName} required />
			<Input label="Nomor HP" bind:value={senderPhone} placeholder="08xxxxxxxxxx" required />
			<Input label="Nama penerima" bind:value={recipientName} required />
			<Textarea label="Alamat pengiriman" bind:value={recipientAddress} rows={2} required />
			<Input label="Jumlah" type="number" bind:value={quantity} />
			<Textarea label="Pesan (opsional)" bind:value={orderMessage} rows={2} />
			<Select label="Metode pembayaran" bind:value={paymentMethod} options={paymentMethodOptions} required />
			{#if orderError}
				<p class="text-sm text-red-600">{orderError}</p>
			{/if}
			<Button type="submit" loading={ordering}>Buat Pesanan</Button>
		</form>
	{/if}
</Modal>
