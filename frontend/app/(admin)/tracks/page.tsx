import { Suspense } from 'react';

import AddTrackDialog from './_components/add-track-dialog';
import TrackList from './_components/track-list';

export const dynamic = 'force-dynamic';

export default function Tracks() {
	return (
		<div className='flex flex-col gap-2'>
			<div className='flex flex-row gap-8 justify-between'>
				<h1>Tracks</h1>

				<AddTrackDialog />
			</div>

			<Suspense fallback={<>Loading...</>}>
				<TrackList />
			</Suspense>
		</div>
	);
}
