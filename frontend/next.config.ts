import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
	output: 'standalone',
	images: {
		remotePatterns: [new URL('https://i.ytimg.com/**')],
	},
};

export default nextConfig;
