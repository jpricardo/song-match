'use server';
import { trackService } from '@/services/track';

export async function submitAudio(byteArray: Uint8Array<ArrayBuffer>) {
	return await trackService.findMatches(byteArray);
}
