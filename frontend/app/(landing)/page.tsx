import MatcherForm from './_components/matcher-form';

export default function Home() {
	return (
		<div className='flex min-h-screen items-center justify-center bg-background font-sans'>
			<main className='flex min-h-screen w-full max-w-3xl flex-col gap-28 py-32 px-16 sm:items-start'>
				<div className='flex flex-col items-center gap-6 text-center sm:items-start sm:text-left'>
					<h1 className='uppercase max-w-xs text-4xl font-bold text-foreground'>Song Match</h1>
					<p className='max-w-md text-lg leading-8 text-zinc-600 dark:text-zinc-400'>
						Let us listen to a song, <b>we&apos;ll</b> (probably) <b>find it.</b>
					</p>
				</div>

				<section className='w-full flex justify-center align-center'>
					<MatcherForm />
				</section>
			</main>
		</div>
	);
}
