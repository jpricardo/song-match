type Props = Readonly<{ children: React.ReactNode }>;
export default function Layout({ children }: Props) {
	return <section className='m-8'>{children}</section>;
}
