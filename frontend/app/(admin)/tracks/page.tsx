import { Plus } from 'lucide-react';
import Link from 'next/link';
import { Suspense } from 'react';

import { Button } from '@/components/ui/button';

import { TrackList } from './_components/track-list';

export const dynamic = 'force-dynamic';

export default function Tracks() {
	return (
		<div className='flex flex-col gap-2'>
			<div className='flex flex-row gap-8 justify-between'>
				<h1>Tracks</h1>

				<Button asChild>
					<Link href='/tracks/new'>
						<Plus />
						Add track
					</Link>
				</Button>
			</div>

			<Suspense fallback={<>Loading...</>}>
				<TrackList />
			</Suspense>
		</div>
	);
}
