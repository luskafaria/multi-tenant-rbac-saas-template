import { auth } from '@/auth/auth'

export default async function Home() {
  const { user } = await auth()

  return <div className="max-w-md">{JSON.stringify(user, null, 2)}</div>
}
