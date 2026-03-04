'use client';
import { fetchFile } from '@ffmpeg/util';

import { useFfmpeg } from './use-ffmpeg';

const OUTPUT_FILE = 'output.wav';
const INPUT_FILE = 'input.wav';

export function useCompression() {
	const ffmpeg = useFfmpeg();

	const compress = async (blob: Blob | undefined) => {
		if (!ffmpeg) throw new Error('FFMPEG is not ready!');
		await ffmpeg.writeFile(INPUT_FILE, await fetchFile(blob));

		await ffmpeg.exec(['-i', INPUT_FILE, '-acodec', 'pcm_s16le', '-ar', '16000', OUTPUT_FILE]);

		const compressedData = await ffmpeg.readFile(OUTPUT_FILE);
		if (typeof compressedData === 'string') throw new Error('Error compressing audio file');
		const uint8Array = new Uint8Array(compressedData);

		await ffmpeg.deleteFile(INPUT_FILE);
		await ffmpeg.deleteFile(OUTPUT_FILE);

		return uint8Array;
	};

	return { compress };
}
