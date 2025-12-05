import { createEnv } from '@t3-oss/env-nextjs'
import { z } from 'zod'

export const env = createEnv({
  server: {
    DATABASE_URL: z.string().url(),
    SERVER_PORT: z.coerce.number().default(3333),
    JWT_SECRET: z.string(),
    GH_OAUTH_CLIENT_ID: z.string(),
    GH_OAUTH_CLIENT_SECRET: z.string(),
    GH_OAUTH_CLIENT_REDIRECT_URI: z.string().url(),
    WEB_URL: z.url(),
  },
  client: {},
  shared: {
    NEXT_PUBLIC_API_URL: z.string().url(),
  },
  runtimeEnv: {
    DATABASE_URL: process.env.DATABASE_URL,
    SERVER_PORT: process.env.SERVER_PORT,
    JWT_SECRET: process.env.JWT_SECRET,
    GH_OAUTH_CLIENT_ID: process.env.GH_OAUTH_CLIENT_ID,
    GH_OAUTH_CLIENT_SECRET: process.env.GH_OAUTH_CLIENT_SECRET,
    GH_OAUTH_CLIENT_REDIRECT_URI: process.env.GH_OAUTH_CLIENT_REDIRECT_URI,
    WEB_URL: process.env.WEB_URL,
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL,
  },
  emptyStringAsUndefined: true,
  skipValidation: !!process.env.CI || !!process.env.SKIP_ENV_VALIDATION,
})
