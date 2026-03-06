'use client';
import { ArrowLeft } from 'lucide-react';

import Track from '@/app/_components/track';
import { Button } from '@/components/ui/button';
import { TrackDTO } from '@/services/track';

type TopMatchProps = { data: TrackDTO };
function TopMatch({ data }: TopMatchProps) {
	return <Track data={data} />;
}

type Props = { data: TrackDTO[]; onRetry: VoidFunction };
export default function Result({ data, onRetry }: Props) {
	const maxMatches = 4;
	const [topMatch, ...matches] = data.slice(0, maxMatches);

	return (
		<div className='flex flex-col gap-4'>
			{topMatch ? (
				<div className='flex flex-col gap-4'>
					<div className='flex flex-col gap-2'>
						<span>Top Match</span>
						<TopMatch data={topMatch} />
					</div>

					<div className='flex flex-col gap-2'>
						<span>Related Tracks</span>
						<div className='flex flex-row gap-4 max-w-1/3'>
							{matches.map((t) => (
								<div key={t.id}>
									<Track data={t} />
								</div>
							))}
						</div>
					</div>
				</div>
			) : (
				<>{"Sorry, we couldn't find your song"}</>
			)}

			<div className='flex justify-end'>
				<Button onClick={onRetry}>
					<ArrowLeft />
					Try again
				</Button>
			</div>
		</div>
	);
}
