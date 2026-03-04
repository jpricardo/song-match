import { addTrack } from '@/app/actions';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Field, FieldGroup, FieldLabel } from '@/components/ui/field';
import { Input } from '@/components/ui/input';

export default function NewTrackPage() {
	return (
		<form action={addTrack}>
			<Card className='w-full max-w-md'>
				<CardHeader>
					<CardTitle>Add Track</CardTitle>
					<CardDescription>Manually add a track/song.</CardDescription>
				</CardHeader>

				<CardContent>
					<FieldGroup>
						<Field>
							<FieldLabel>URL (Youtube)</FieldLabel>
							<Input type='text' name='url' placeholder='Insert here' required />
						</Field>
					</FieldGroup>
				</CardContent>

				<CardFooter className='justify-end'>
					<Button type='submit'>Submit</Button>
				</CardFooter>
			</Card>
		</form>
	);
}
