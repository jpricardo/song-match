import { ChevronRightIcon } from 'lucide-react';
import Link from 'next/link';

import { Item, ItemActions, ItemContent, ItemFooter, ItemTitle } from '@/components/ui/item';
import { trackService } from '@/services/track';

export async function TrackList() {
	const { tracks } = await trackService.getMany();

	return (
		<div className='flex flex-col gap-2'>
			<span>Total: {tracks.length}</span>

			{tracks.map((t) => (
				<Item key={t.id} variant='outline' asChild>
					<Link href={t.url} target='_blank'>
						<ItemContent>
							<ItemTitle>{t.name}</ItemTitle>
							<ItemFooter>{t.fingerprints.length} fingerprints</ItemFooter>
						</ItemContent>
						<ItemActions>
							<ChevronRightIcon className='size-4' />
						</ItemActions>
					</Link>
				</Item>
			))}
		</div>
	);
}
