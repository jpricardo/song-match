import { Trash2 } from 'lucide-react';

import { deleteTrack } from '@/app/actions';
import { Button } from '@/components/ui/button';
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog';
import { Field, FieldGroup } from '@/components/ui/field';
import { Input } from '@/components/ui/input';
import { TrackDTO } from '@/services/track';

type Props = { track: TrackDTO };
export default function DeleteTrackDialog({ track }: Props) {
	return (
		<Dialog>
			<DialogTrigger asChild>
				<Button variant='destructive'>
					<Trash2 />
					Delete track
				</Button>
			</DialogTrigger>

			<DialogContent>
				<form action={deleteTrack}>
					<DialogHeader>
						<DialogTitle>Delete {track.name}?</DialogTitle>
						<DialogDescription>Are you sure?</DialogDescription>
					</DialogHeader>

					<FieldGroup>
						<Field>
							<Input name='id' type='hidden' value={track.id} />
						</Field>

						<Field>
							<DialogFooter>
								<Button variant='destructive' type='submit'>
									<Trash2 />
									Submit
								</Button>
							</DialogFooter>
						</Field>
					</FieldGroup>
				</form>
			</DialogContent>
		</Dialog>
	);
}
