'use client';
import { Info } from 'lucide-react';

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';

import { Timer } from './timer';

type Props = {
	currentTime: number;
	state: RecordingState;
	error: Error | undefined;
	loading: boolean;
	onStart: VoidFunction;
	onStop: VoidFunction;
};
export default function AudioRecorder({ currentTime, state, error, loading, onStart, onStop }: Props) {
	return (
		<div className='flex flex-col gap-2 justify-center items-center w-full max-w-xs'>
			{error && (
				<Alert variant='destructive' className='w-full'>
					<Info />
					<AlertTitle>Something went wrong!</AlertTitle>
					<AlertDescription>{error.message}</AlertDescription>
				</Alert>
			)}

			<Timer
				currentTime={currentTime}
				recording={state === 'recording'}
				loading={loading}
				onClick={() => {
					if (state === 'recording') return onStop();
					onStart();
				}}
			/>
		</div>
	);
}
