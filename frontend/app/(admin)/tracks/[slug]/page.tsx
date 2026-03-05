import Image from 'next/image';

import { trackService } from '@/services/track';

import DeleteTrackDialog from './_components/delete-track-dialog';
import FingerprintChart from './_components/fingerprint-chart';

type Props = { params: Promise<{ slug: string }> };
export default async function TrackDetails({ params }: Props) {
	const { slug } = await params;
	const track = await trackService.getOne(slug);

	if (!track) return <>Track not found</>;

	const flatPoints: number[] = [];

	track.fingerprints.forEach((fp) => {
		if (!fp.peaks) return;
		fp.peaks.forEach((peak) => {
			flatPoints.push(Number(fp.timestamp.toFixed(3)), peak);
		});
	});

	const maxTime = track.fingerprints.length > 0 ? track.fingerprints[track.fingerprints.length - 1].timestamp : 1;
	const maxFreq = 4096;

	return (
		<div className='flex flex-col gap-2'>
			<div className='flex flex-row gap-8 justify-between'>
				<h1>{track.name}</h1>

				<DeleteTrackDialog track={track} />
			</div>

			{!!track.thumbnail && <Image src={track.thumbnail} alt='thumbnail' width={500} height={500} />}

			<FingerprintChart flatPoints={flatPoints} maxTime={maxTime} maxFreq={maxFreq} />
		</div>
	);
}
