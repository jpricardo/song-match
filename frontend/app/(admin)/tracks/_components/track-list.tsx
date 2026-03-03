import { trackService } from '@/services/track';

export async function TrackList() {
	const { tracks } = await trackService.getMany();

	return <>There are {tracks.length} tracks!</>;
}
