'use client';
import { Check, CircleStop, Info, Mic } from 'lucide-react';
import { useState } from 'react';

import { submitAudio } from '@/app/actions';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Spinner } from '@/components/ui/spinner';
import { useCompression } from '@/hooks/use-compression';
import { useRecorder } from '@/hooks/use-recorder';
import { TrackDTO } from '@/services/track';

export default function AudioRecorder() {
	const compression = useCompression();
	const recorder = useRecorder();
	const [matches, setMatches] = useState<TrackDTO[]>();
	const [error, setError] = useState<Error>();
	const [loading, setLoading] = useState(false);

	const handleSubmit = async () => {
		setLoading(true);
		setError(undefined);
		setMatches(undefined);

		try {
			const audioData = await recorder.getData();
			const compressed = await compression.compress(audioData);
			const res = await submitAudio(compressed);
			setMatches(res.matches);
		} catch (err) {
			setError(err as Error);
		} finally {
			setLoading(false);
		}
	};

	return (
		<div className='flex flex-col gap-2 justify-center items-center w-full'>
			{loading && (
				<Alert className='w-full'>
					<Spinner />
					<AlertTitle>Loading...</AlertTitle>
				</Alert>
			)}

			{error && (
				<Alert variant='destructive' className='w-full'>
					<Info />
					<AlertTitle>Something went wrong!</AlertTitle>
					<AlertDescription>{error.message}</AlertDescription>
				</Alert>
			)}

			{matches && (
				<Alert className='w-full'>
					<Check />
					<AlertTitle>{matches.length} matches found!</AlertTitle>
				</Alert>
			)}

			{recorder.state ? (
				<div className='flex flex-col gap-4'>
					<audio className='w-full min-w-xs' src={recorder.url} controls />

					<span>{recorder.state === 'recording' && <>Recording...</>}</span>

					<div className='flex gap-2 w-full flex-col'>
						{recorder.state === 'recording' ? (
							<Button onClick={() => recorder.stop()}>
								<CircleStop />
								Stop
							</Button>
						) : (
							<Button onClick={() => recorder.start()}>
								<Mic />
								Start recording
							</Button>
						)}

						{recorder.url && (
							<Button onClick={handleSubmit}>
								{loading ? <Spinner /> : <Check />}
								Submit
							</Button>
						)}
					</div>
				</div>
			) : (
				<Button className='max-w-xs' onClick={async () => await recorder.setup()}>
					<Mic />
					Allow Microphone Access
				</Button>
			)}
		</div>
	);
}
