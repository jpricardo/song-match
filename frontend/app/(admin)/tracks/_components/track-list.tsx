import Track, { TrackSkeleton } from '@/app/_components/track';
import { trackService } from '@/services/track';

export function TrackListSkeleton() {
	const amt = 24;
	const items = Array.from({ length: amt });

	return (
		<div className='flex flex-row gap-4 flex-wrap'>
			{items.map((_, i) => (
				<div key={i} className='min-w-xs sm:w-full md:w-xs'>
					<TrackSkeleton key={i} />
				</div>
			))}
		</div>
	);
}

export default async function TrackList() {
	const { tracks } = await trackService.getMany();

	return (
		<div className='flex flex-col gap-2'>
			<span>Total: {tracks.length}</span>

			<div className='flex flex-row gap-4 flex-wrap'>
				{tracks.map((t) => (
					<div key={t.id} className='min-w-xs sm:w-full md:w-xs'>
						<Track data={t} />
					</div>
				))}
			</div>
		</div>
	);
}
