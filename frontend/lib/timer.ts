export async function timer(ms: number) {
	return new Promise<void>((resolve) => setTimeout(resolve, ms));
}
