type Props = Readonly<{ children: React.ReactNode }>;
export default function Layout({ children }: Props) {
	return <section>{children}</section>;
}
