import { resolve, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = dirname(fileURLToPath(import.meta.url))

/** @type {import('next').NextConfig} */
const nextConfig = {
  turbopack: {
    root: resolve(__dirname, '../..'),
  },
  images: {
    remotePatterns: [
      {
        hostname: 'github.com',
      },
      { hostname: 'avatars.githubusercontent.com' },
    ],
  },
}

export default nextConfig
