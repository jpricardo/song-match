import { Mic } from 'lucide-react';

import { Button } from '@/components/ui/button';

type Props = { onClick: VoidFunction };
export default function AllowMic({ onClick }: Props) {
	return (
		<div className='max-w-xs w-full'>
			<Button className='w-full' onClick={onClick}>
				<Mic />
				Allow Microphone Access
			</Button>
		</div>
	);
}
