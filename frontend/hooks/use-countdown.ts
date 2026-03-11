import { useCallback, useEffect, useState } from 'react';

type CountdownStatus = 'idle' | 'running' | 'done';

type CountdownOptions = {
	step?: number;
	onStart?: () => void;
	onStop?: () => void;
	onUpdate?: (v: number) => void;
	onReset?: () => void;
};

type StartCountdownOptions = { reset: boolean };

export function useCountdown(duration: number, options?: CountdownOptions) {
	const [current, setCurrent] = useState(duration);
	const [status, setStatus] = useState<CountdownStatus>('idle');

	const countdownStep = options?.step ?? 100;

	const reset = () => {
		setStatus('idle');
		setCurrent(duration);
		options?.onReset?.();
	};

	const start = (o?: StartCountdownOptions) => {
		if (o?.reset) reset();

		setStatus('running');
		options?.onStart?.();
	};

	const stop = useCallback(() => {
		if (status === 'done') return;

		setStatus('done');
		options?.onStop?.();
	}, [options, status]);

	useEffect(() => {
		if (status !== 'running') return;

		const interval = setInterval(() => {
			setCurrent((p) => {
				const n = p - countdownStep;
				options?.onUpdate?.(n);
				if (n <= 0) stop();
				return n;
			});
		}, countdownStep);
		return () => clearInterval(interval);
	}, [current, options, status, stop, countdownStep]);

	return { start, stop, reset, current, status };
}
