'use client';
import { CircleStop, Mic, Play } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { useRecorder } from '@/hooks/use-recorder';

export function AudioRecorder() {
	const recorder = useRecorder();

	return (
		<div className='flex flex-col gap-2 justify-center items-center max-w-md'>
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
								<Play />
								Start
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
