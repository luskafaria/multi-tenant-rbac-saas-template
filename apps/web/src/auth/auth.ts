import { getProfile } from '@/http/get-profile'
import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'

export function isAuthenticated() {
  return !!cookies().get('token')?.value
}

export async function auth() {
  const token = cookies().get('token')?.value

  if (!token) {
    redirect('/api/auth/sign-out')
  }

  try {
    const { user } = await getProfile()

    return { user }
  } catch (err) {
    console.warn(err)
  }

  redirect('/auth/sign-in')
}
