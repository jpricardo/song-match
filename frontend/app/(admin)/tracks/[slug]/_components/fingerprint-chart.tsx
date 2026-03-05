'use client';

import { Card, CardContent, CardTitle } from '@/components/ui/card';
import { useEffect, useRef } from 'react';

type Props = {
	flatPoints: number[];
	maxTime: number;
	maxFreq: number;
};

export default function FingerprintChart({ flatPoints, maxTime, maxFreq }: Props) {
	const canvasRef = useRef<HTMLCanvasElement>(null);
	const containerRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		const canvas = canvasRef.current;
		const container = containerRef.current;
		if (!canvas || !container || flatPoints.length === 0) return;

		const ctx = canvas.getContext('2d');
		if (!ctx) return;

		const draw = () => {
			const rect = container.getBoundingClientRect();
			if (rect.width === 0 || rect.height === 0) return;

			const dpr = window.devicePixelRatio || 1;

			canvas.width = rect.width * dpr;
			canvas.height = rect.height * dpr;
			canvas.style.width = `${rect.width}px`;
			canvas.style.height = `${rect.height}px`;

			ctx.scale(dpr, dpr);

			ctx.clearRect(0, 0, rect.width, rect.height);
			ctx.fillStyle = '#09090b';
			ctx.fillRect(0, 0, rect.width, rect.height);

			ctx.fillStyle = '#2563eb';

			// Loop through the 1D array in pairs! (i = time, i+1 = frequency)
			for (let i = 0; i < flatPoints.length; i += 2) {
				const time = flatPoints[i];
				const freq = flatPoints[i + 1];

				const x = (time / maxTime) * rect.width;
				const y = rect.height - (freq / maxFreq) * rect.height;

				ctx.fillRect(x, y, 2, 2);
			}
		};

		const observer = new ResizeObserver(() => {
			requestAnimationFrame(draw);
		});

		observer.observe(container);

		return () => observer.disconnect();
	}, [flatPoints, maxTime, maxFreq]);

	return (
		<Card className='w-full'>
			<CardContent className='flex flex-col gap-2'>
				<CardTitle>
					<h2>Frequencies ({flatPoints.length / 2})</h2>
				</CardTitle>
				<div ref={containerRef} className='w-full h-100 border rounded-lg overflow-hidden bg-zinc-950'>
					<canvas ref={canvasRef} className='block' />
				</div>
			</CardContent>
		</Card>
	);
}
