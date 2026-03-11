import { Button } from '@/components/ui/button';
import { parseTime } from '@/lib/time';

const secondsFormatter = new Intl.NumberFormat('en-US', {
	minimumFractionDigits: 2,
	maximumFractionDigits: 2,
	minimumIntegerDigits: 2,
});

const minutesFormatter = new Intl.NumberFormat('en-US', {
	minimumFractionDigits: 0,
	maximumFractionDigits: 0,
	minimumIntegerDigits: 2,
});

type Props = {
	currentTime: number;
	recording: boolean;
	loading: boolean;
	onClick: VoidFunction;
};
export function Timer({ currentTime, recording, loading, onClick }: Props) {
	const { minutes, seconds } = parseTime(currentTime);

	const timeStr = `${minutesFormatter.format(minutes)}:${secondsFormatter.format(seconds)}`;

	return (
		<Button
			variant='outline'
			className='w-60 h-60 rounded-full flex justify-center items-center cursor-pointer select-none font-mono text-4xl font-semibold'
			disabled={loading}
			onClick={onClick}
		>
			{recording ? <>{timeStr}</> : <>START</>}
		</Button>
	);
}
