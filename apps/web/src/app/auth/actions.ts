'use server'

import crypto from 'node:crypto'

import { env } from '@saas/env'
import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'

import { security } from '@/config/security'

export async function signInWithGithub() {
  const githubSignInURL = new URL('login/oauth/authorize', 'https://github.com')

  const state = crypto.randomBytes(32).toString('hex')
  const cookieStore = await cookies()
  cookieStore.set('oauth_state', state, {
    path: '/',
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'lax',
    maxAge: security.OAUTH_STATE_EXPIRY_SECONDS,
  })

  githubSignInURL.searchParams.set('client_id', env.GH_OAUTH_CLIENT_ID)
  githubSignInURL.searchParams.set(
    'redirect_uri',
    env.GH_OAUTH_CLIENT_REDIRECT_URI
  )
  githubSignInURL.searchParams.set('scope', 'user')
  githubSignInURL.searchParams.set('state', state)

  redirect(githubSignInURL.toString())
}
