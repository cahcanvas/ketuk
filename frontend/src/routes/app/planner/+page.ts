import type { Event } from '@ketuk/shared';
import { redirect } from '@sveltejs/kit';
import type { PlannerSummary } from '$lib/api';
import { getPlannerSummary, listMyEvents } from '$lib/api';
import type { PageLoad } from './$types';

/**
 * Planner tidak nested di bawah /app/undangan/[id] — dipilih lewat query `?event=`.
 * Kalau user cuma punya satu event, langsung diarahkan ke situ supaya tidak
 * perlu memilih dari daftar berisi satu item.
 */
export const load: PageLoad = async ({ fetch, parent, url }) => {
	const { accessToken } = await parent();
	const eventId = url.searchParams.get('event');

	let events: Event[] = [];
	try {
		events = await listMyEvents({ fetch, accessToken });
	} catch {
		return {
			events: [] as Event[],
			eventId,
			summary: null as PlannerSummary | null,
			error: 'Koneksi terputus. Coba muat ulang.',
		};
	}

	if (!eventId) {
		if (events.length === 1 && events[0]) {
			throw redirect(303, `/app/planner?event=${events[0].id}`);
		}
		return { events, eventId: null, summary: null as PlannerSummary | null, error: null };
	}

	try {
		const summary = await getPlannerSummary(eventId, { fetch, accessToken });
		return { events, eventId, summary, error: null };
	} catch {
		return {
			events,
			eventId,
			summary: null as PlannerSummary | null,
			error: 'Gagal memuat ringkasan planner.',
		};
	}
};
