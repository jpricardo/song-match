import { Check, Plus } from 'lucide-react';

import { addTrack } from '@/app/actions';
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
import { Field, FieldGroup, FieldLabel } from '@/components/ui/field';
import { Input } from '@/components/ui/input';

export default function AddTrackDialog() {
	return (
		<Dialog>
			<DialogTrigger asChild>
				<Button>
					<Plus />
					Add track
				</Button>
			</DialogTrigger>

			<DialogContent>
				<form action={addTrack}>
					<DialogHeader>
						<DialogTitle>Add Track</DialogTitle>
						<DialogDescription>Manually add a track/song.</DialogDescription>
					</DialogHeader>

					<FieldGroup>
						<Field>
							<FieldLabel>URL (Youtube)</FieldLabel>
							<Input type='text' name='url' placeholder='Insert here' required />
						</Field>

						<Field>
							<DialogFooter>
								<Button type='submit'>
									<Check />
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
