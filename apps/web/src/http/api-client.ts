import { env } from '@saas/env'
import { getCookie } from 'cookies-next/client'
import { getCookie as getServerCookie } from 'cookies-next/server'
import ky from 'ky'

export const api = ky.create({
  prefixUrl: env.NEXT_PUBLIC_API_URL,
  hooks: {
    beforeRequest: [
      async (request) => {
        let token: string | undefined

        if (typeof window === 'undefined') {
          const { cookies } = await import('next/headers')
          token = await getServerCookie('token', { cookies })
        } else {
          token = getCookie('token')
        }

        if (token) {
          request.headers.set('Authorization', `Bearer ${token}`)
        }
      },
    ],
  },
})
