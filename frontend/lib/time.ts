export function parseTime(ms: number) {
	let seconds = ms / 1000;
	let minutes = 0;
	let hours = 0;

	while (seconds >= 60) {
		minutes++;
		seconds -= 60;
	}

	while (minutes >= 60) {
		hours++;
		minutes -= 60;
	}

	return { seconds, minutes, hours };
}
