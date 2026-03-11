'use client';
import { useState } from 'react';

import { submitAudio } from '@/app/actions';
import { useCompression } from '@/hooks/use-compression';
import { useCountdown } from '@/hooks/use-countdown';
import { useRecorder } from '@/hooks/use-recorder';
import { TrackDTO } from '@/services/track';

import AudioRecorder from './audio-recorder';
import Result from './result';

type FormStep = 'record' | 'result';

export default function MatcherForm() {
	const [step, setStep] = useState<FormStep>('record');
	const [data, setData] = useState<TrackDTO[]>();
	const [error, setError] = useState<Error>();
	const [loading, setLoading] = useState(false);

	const compression = useCompression();
	const submitData = async (data: Blob) => {
		setLoading(true);

		try {
			const compressed = await compression.compress(data);
			const res = await submitAudio(compressed);
			setData(res.matches);
			setStep('result');
		} catch (err) {
			setError(err as Error);
		} finally {
			setLoading(false);
		}
	};

	const recorder = useRecorder({
		onStart: () => setError(undefined),
		onData: (b) => submitData(b),
	});

	const countdownDuration = 5 * 1000;
	const countdown = useCountdown(countdownDuration, {
		step: 10,
		onStart: () => recorder.start(),
		onStop: () => recorder.stop(),
	});

	const currentTime = countdownDuration - countdown.current;

	const steps: Record<FormStep, React.ReactNode> = {
		record: (
			<AudioRecorder
				currentTime={currentTime}
				state={recorder.state ?? 'inactive'}
				error={error}
				loading={loading}
				onStart={() => countdown.start({ reset: true })}
				onStop={() => countdown.stop()}
			/>
		),
		result: (
			<Result
				data={data ?? []}
				onRetry={() => {
					setData([]);
					setStep('record');
				}}
			/>
		),
	};

	return steps[step];
}
