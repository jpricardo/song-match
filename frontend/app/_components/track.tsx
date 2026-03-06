import { ChevronRightIcon } from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';

import { Badge } from '@/components/ui/badge';
import { Item, ItemActions, ItemContent, ItemDescription, ItemHeader, ItemTitle } from '@/components/ui/item';
import { Skeleton } from '@/components/ui/skeleton';
import { TrackDTO } from '@/services/track';

const badgeVariants = {
	ready: 'default',
	failed: 'destructive',
	processing: 'outline',
} as const;

export function TrackSkeleton() {
	return (
		<Item variant='outline' className='w-full h-full'>
			<ItemHeader>
				<Skeleton className='w-full h-50' />
			</ItemHeader>

			<ItemContent>
				<ItemTitle>
					<Skeleton className='h-4 w-16' />
				</ItemTitle>
				<ItemDescription>
					<Skeleton className='h-4 w-8' />
				</ItemDescription>
			</ItemContent>
		</Item>
	);
}

type Props = { data: TrackDTO };
export default function Track({ data }: Props) {
	return (
		<Item variant='outline' asChild className='w-full h-full'>
			<Link href={`/tracks/${data.id}`}>
				<ItemHeader>
					{!data.thumbnail ? (
						<Skeleton className='w-full h-50' />
					) : (
						<Image src={data.thumbnail} alt='thumbnail' width={500} height={500} />
					)}
				</ItemHeader>

				<ItemContent>
					<ItemTitle>{data.name}</ItemTitle>
					<ItemDescription>
						<Badge className='capitalize' variant={badgeVariants[data.status]}>
							{data.status}
						</Badge>
					</ItemDescription>
				</ItemContent>
				<ItemActions>
					<ChevronRightIcon className='size-4' />
				</ItemActions>
			</Link>
		</Item>
	);
}
