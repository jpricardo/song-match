import { Suspense } from 'react';

import { TrackList } from './_components/track-list';

export const dynamic = 'force-dynamic';

export default function Employees() {
	return (
		<div className='flex flex-col gap-2'>
			<h2>Tracks</h2>

			<Suspense fallback={<>Loading...</>}>
				<TrackList />
			</Suspense>
		</div>
	);
}
