'use client';
import { useState } from 'react';

import { submitAudio } from '@/app/actions';
import { useCompression } from '@/hooks/use-compression';
import { useRecorder } from '@/hooks/use-recorder';
import { TrackDTO } from '@/services/track';

import AllowMic from './allow-mic';
import AudioRecorder from './audio-recorder';
import Result from './result';

type FormStep = 'start' | 'record' | 'result';

export default function MatcherForm() {
	const recorder = useRecorder();
	const compression = useCompression();

	const [step, setStep] = useState<FormStep>('start');
	const [data, setData] = useState<TrackDTO[]>();

	const steps: Record<FormStep, React.ReactNode> = {
		start: (
			<AllowMic
				onClick={async () => {
					await recorder.setup();
					setStep('record');
				}}
			/>
		),
		record: (
			<AudioRecorder
				url={recorder.url}
				state={recorder.state ?? 'inactive'}
				onStart={() => recorder.start()}
				onStop={() => recorder.stop()}
				onSubmit={async () => {
					const audioData = await recorder.getData();
					const compressed = await compression.compress(audioData);
					const res = await submitAudio(compressed);
					setData(res.matches);
					setStep('result');
				}}
			/>
		),
		result: (
			<Result
				data={data ?? []}
				onRetry={() => {
					setData([]);
					setStep('start');
				}}
			/>
		),
	};

	return steps[step];
}
