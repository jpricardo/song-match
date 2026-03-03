'use client';
import { useState } from 'react';

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

	return {
		ready,
		state,
		url,
		data,
		start,
		stop,
		setup,
	};
}
