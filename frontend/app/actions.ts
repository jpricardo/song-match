'use server';

import { trackService } from '@/services/track';

export async function submitAudio(d: Blob) {
	return await trackService.findMatches({ content: d });
}
