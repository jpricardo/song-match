'use client';
import type { FFmpeg } from '@ffmpeg/ffmpeg';
import { useEffect, useState } from 'react';

export function useFfmpeg() {
	const [ffmpeg, setFfmpeg] = useState<FFmpeg | null>(null);

	useEffect(() => {
		let mounted = true;

		import('@ffmpeg/ffmpeg').then(({ FFmpeg }) => {
			if (!mounted) return;
			const ff = new FFmpeg();
			if (!ff.loaded) {
				ff.load();
			}
			setFfmpeg(ff);
		});

		return () => {
			mounted = false;
		};
	}, []);

	return ffmpeg;
}
