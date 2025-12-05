import { env } from '@saas/env'
import 'dotenv/config'

import { defineConfig } from 'prisma/config'

export default defineConfig({
  schema: 'prisma/schema.prisma',
  migrations: {
    path: 'prisma/migrations',
    seed: 'tsx prisma/seeds.ts',
  },
  datasource: {
    url: env.DATABASE_URL,
  },
})
