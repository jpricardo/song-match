'use server';
import { refresh } from 'next/cache';
import { redirect } from 'next/navigation';

import { trackService } from '@/services/track';

export async function submitAudio(byteArray: Uint8Array<ArrayBuffer>) {
	return await trackService.findMatches(byteArray);
}

export async function addTrack(fd: FormData) {
	const url = fd.get('url');
	if (typeof url !== 'string') throw new Error('Invalid URL!');

	await trackService.addTrack({ url });

	refresh();
	redirect('/tracks');
}

export async function deleteTrack(fd: FormData) {
	const id = fd.get('id');
	if (typeof id !== 'string') throw new Error('Invalid ID!');

	await trackService.delete(id);

	refresh();
	redirect('/tracks');
}
