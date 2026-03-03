'use client';
import { fetchFile } from '@ffmpeg/util';

import { useFfmpeg } from './use-ffmpeg';

export function useCompression() {
	const ffmpeg = useFfmpeg();

	const compress = async (blob: Blob | undefined) => {
		if (!ffmpeg) throw new Error('FFMPEG is not ready!');
		await ffmpeg.writeFile('input.wav', await fetchFile(blob));

		await ffmpeg.exec(['-i', 'input.wav', '-acodec', 'libmp3lame', '-ab', '128k', '-ar', '16000', 'output.mp3']);

		const compressedData = await ffmpeg.readFile('output.mp3');
		if (typeof compressedData === 'string') throw new Error('Error compressing audio file');
		const uint8Array = new Uint8Array(compressedData);

		await ffmpeg.deleteFile('input.wav');
		await ffmpeg.deleteFile('output.mp3');

		return uint8Array;
	};

	return { compress };
}
