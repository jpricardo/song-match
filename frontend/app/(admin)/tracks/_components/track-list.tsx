import { ChevronRightIcon } from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';

import { Item, ItemActions, ItemContent, ItemHeader, ItemTitle } from '@/components/ui/item';
import { trackService } from '@/services/track';

export default async function TrackList() {
	const { tracks } = await trackService.getMany();

	return (
		<div className='flex flex-col gap-2'>
			<span>Total: {tracks.length}</span>

			<div className='flex flex-row gap-4 flex-wrap'>
				{tracks.map((t) => (
					<Item key={t.id} variant='outline' asChild className='max-w-xs w-full'>
						<Link href={`/tracks/${t.id}`}>
							<ItemHeader>
								{!!t.thumbnail && <Image src={t.thumbnail} alt='thumbnail' width={500} height={500} />}
							</ItemHeader>
							<ItemContent>
								<ItemTitle>{t.name}</ItemTitle>
							</ItemContent>
							<ItemActions>
								<ChevronRightIcon className='size-4' />
							</ItemActions>
						</Link>
					</Item>
				))}
			</div>
		</div>
	);
}
