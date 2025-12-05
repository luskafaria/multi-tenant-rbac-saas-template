import { cookies } from 'next/headers'
import { NextRequest, NextResponse } from 'next/server'

import { security } from '@/config/security'
import { acceptInvite } from '@/http/accept-invite'
import { signInWithGithub } from '@/http/sign-in-with-github'

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams

  const code = searchParams.get('code')
  const state = searchParams.get('state')

  if (!code) {
    return NextResponse.json(
      { message: 'Github OAuth code was not found.' },
      { status: 400 }
    )
  }

  const cookieStore = await cookies()
  const storedState = cookieStore.get('oauth_state')?.value

  if (!state || !storedState || state !== storedState) {
    return NextResponse.json(
      { message: 'Invalid OAuth state parameter.' },
      { status: 400 }
    )
  }

  cookieStore.delete('oauth_state')

  try {
    const { token } = await signInWithGithub({ code })

    cookieStore.set('token', token, {
      path: '/',
      maxAge: security.SESSION_EXPIRY_SECONDS,
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
    })

    const inviteId = cookieStore.get('inviteId')?.value

    if (inviteId) {
      try {
        await acceptInvite(inviteId)
        cookieStore.delete('inviteId')
      } catch {
        // Silently fail if accepting invite fails
      }
    }

    const redirectUrl = request.nextUrl.clone()

    redirectUrl.pathname = '/'
    redirectUrl.search = ''

    return NextResponse.redirect(redirectUrl)
  } catch (error) {
    console.warn(error)

    // TODO: Handle Github auth errors properly
  }
}
