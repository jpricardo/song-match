'use client';
import { useCallback, useEffect, useState } from 'react';

import { encodeWAV } from '@/lib/audio';

async function getBlob(dataBlob: Blob) {
	// Decode WebM to raw PCM samples
	const arrayBuffer = await dataBlob.arrayBuffer();
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
}

type RecorderOptions = { onStart?: VoidFunction; onStop?: VoidFunction; onData?: (data: Blob) => void };
export function useRecorder(options?: RecorderOptions) {
	const [mediaRecorder, setMediaRecorder] = useState<MediaRecorder>();

	const state = mediaRecorder?.state;

	const start = () => mediaRecorder?.start();
	const stop = () => mediaRecorder?.stop();

	const setup = useCallback(async () => {
		const stream = await window.navigator.mediaDevices.getUserMedia({ audio: { noiseSuppression: { ideal: false } } });
		const mr = new MediaRecorder(stream);

		mr.addEventListener('start', () => {
			console.log('Starting recorder...');
			options?.onStart?.();
		});
		mr.addEventListener('stop', () => {
			console.log('Stopping recorder...');
			options?.onStop?.();
		});
		mr.addEventListener('dataavailable', (ev) => {
			if (ev.data.size <= 0) return;
			return getBlob(new Blob([ev.data], { type: 'audio/webm' })).then((b) => options?.onData?.(b));
		});

		return mr;
	}, [options]);

	useEffect(() => {
		if (!mediaRecorder) setup().then((mr) => setMediaRecorder(mr));
	}, [mediaRecorder, setup]);

	return { state, start, stop };
}
