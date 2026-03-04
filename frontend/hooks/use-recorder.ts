'use client';
import { useState } from 'react';

import { encodeWAV } from '@/lib/audio';

type MediaRecorderOptions = {
	onStart: VoidFunction;
	onStop: VoidFunction;
	onDataAvailable: (ev: BlobEvent) => void;
};
function createRecorder(stream: MediaStream, options: MediaRecorderOptions): MediaRecorder {
	const mr = new MediaRecorder(stream);

	mr.ondataavailable = options.onDataAvailable;
	mr.onstart = options.onStart;
	mr.onstop = options.onStop;

	return mr;
}

export function useRecorder() {
	const [mediaRecorder, setMediaRecorder] = useState<MediaRecorder>();
	const [chunks, setChunks] = useState<Blob[]>([]);

	const start = () => mediaRecorder?.start();
	const stop = () => mediaRecorder?.stop();

	const setup = async () => {
		const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
		const mr = createRecorder(stream, {
			onDataAvailable: (event) => event.data.size > 0 && setChunks((prev) => [...prev, event.data]),
			onStart: () => setChunks([]),
			onStop: () => console.log('Stop!'),
		});

		setMediaRecorder(mr);
	};

	const ready = !!mediaRecorder;
	const state = mediaRecorder?.state;

	const data = chunks.length ? new Blob(chunks, { type: 'audio/webm' }) : undefined;
	const url = data ? URL.createObjectURL(data) : undefined;

	const getData = async () => {
		if (!data) throw new Error('No data!');

		// Decode WebM to raw PCM samples
		const arrayBuffer = await data.arrayBuffer();
		const audioContext = new AudioContext();

		try {
			const audioBuffer = await audioContext.decodeAudioData(arrayBuffer);
			console.log('Decoded audio:', audioBuffer.sampleRate, audioBuffer.length);

			// Manually construct WAV file from PCM
			const wav = encodeWAV(audioBuffer);
			console.log('WAV buffer size:', wav.byteLength);
			return new Blob([wav], { type: 'audio/wav' });
		} catch (err) {
			console.error('Failed to decode audio:', err);
			throw err;
		}
	};

	return {
		ready,
		state,
		url,
		getData,
		start,
		stop,
		setup,
	};
}
