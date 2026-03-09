'use client';
import { Check, CircleStop, Info, Mic } from 'lucide-react';
import { useState } from 'react';

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Spinner } from '@/components/ui/spinner';

type Props = {
	state: RecordingState;
	url: string | undefined;
	onStart: VoidFunction;
	onStop: VoidFunction;
	onSubmit: () => Promise<void>;
};
export default function AudioRecorder({ state, url, onStart, onStop, onSubmit }: Props) {
	const [error, setError] = useState<Error>();
	const [loading, setLoading] = useState(false);

	const handleSubmit = async () => {
		setLoading(true);
		setError(undefined);

		try {
			await onSubmit();
		} catch (err) {
			setError(err as Error);
		} finally {
			setLoading(false);
		}
	};

	return (
		<div className='flex flex-col gap-2 justify-center items-center w-full max-w-xs'>
			{error && (
				<Alert variant='destructive' className='w-full'>
					<Info />
					<AlertTitle>Something went wrong!</AlertTitle>
					<AlertDescription>{error.message}</AlertDescription>
				</Alert>
			)}

			<div className='flex flex-col gap-4 w-full'>
				{state === 'recording' && <span>Recording...</span>}

				<div className='flex gap-2 w-full flex-col'>
					{state === 'recording' ? (
						<Button onClick={onStop}>
							<CircleStop />
							Stop
						</Button>
					) : (
						<Button onClick={onStart}>
							<Mic />
							Start Recording
						</Button>
					)}

					{url && (
						<Button onClick={handleSubmit} disabled={loading}>
							{loading ? <Spinner /> : <Check />}
							Submit
						</Button>
					)}
				</div>
			</div>
		</div>
	);
}
