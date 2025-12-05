import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import { CheckCircle, LogIn, LogOut } from 'lucide-react'
import { cookies } from 'next/headers'
import Link from 'next/link'
import { redirect } from 'next/navigation'

import { auth, isAuthenticated } from '@/components/auth/auth'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { acceptInvite } from '@/http/accept-invite'
import { getInvite } from '@/http/get-invite'

dayjs.extend(relativeTime)

interface InvitePageProps {
  readonly params: {
    id: string
  }
}

export default async function InvitePage({ params }: InvitePageProps) {
  const { id: inviteId } = params

  const { invite } = await getInvite(inviteId)
  const isUserAuthenticated = isAuthenticated()

  let currentUserEmail = null
  if (isUserAuthenticated) {
    const { user } = await auth()

    currentUserEmail = user.email
  }

  const userIsAuthenticatedWithSameEmailFromInvite =
    currentUserEmail === invite.email

  async function signInFromInvite() {
    'use server'

    cookies().set('inviteId', inviteId)
    redirect(`/auth/sign-in?email=${invite.email}`)
  }

  async function acceptInviteAction() {
    'use server'

    await acceptInvite(inviteId)
    redirect('/')
  }

  return (
    <main className="h-screen content-center items-center justify-center space-y-4 py-4">
      <Card className="mx-auto w-full max-w-md align-middle">
        <CardHeader className="p-6 text-foreground dark:bg-white">
          <div className="flex items-center gap-4">
            <Avatar className="h-12 w-12">
              {invite.author?.avatarUrl && (
                <AvatarImage
                  src={invite.author?.avatarUrl}
                  alt={`@${invite.author.name}`}
                />
              )}
              <AvatarFallback />
            </Avatar>
            <div className="grid gap-1 dark:text-black">
              <h2 className="text-balance text-lg">
                <span className="font-semibold">
                  {`${invite.author?.name ?? 'Someone'} `}
                </span>
                <span>invited you to join </span>
                <span className="font-semibold">
                  {`${invite.organization.name} `}
                </span>
                <span>{dayjs(invite.createdAt).fromNow()}</span>
              </h2>
            </div>
          </div>
        </CardHeader>
        <CardContent className="grid gap-4 p-6">
          {!isUserAuthenticated && (
            <form action={signInFromInvite}>
              <Button type="submit" variant="secondary" className="w-full">
                <LogIn className="mr-2 size-4" />
                Sign in to accept the invite
              </Button>
            </form>
          )}

          {userIsAuthenticatedWithSameEmailFromInvite && (
            <form action={acceptInviteAction}>
              <Button type="submit" variant="secondary" className="w-full">
                <CheckCircle className="mr-2 size-4" />
                Join {invite.organization.name}
              </Button>
            </form>
          )}

          {isUserAuthenticated &&
            !userIsAuthenticatedWithSameEmailFromInvite && (
              <div className="space-y-4">
                <p className="text-balance text-center text-sm leading-relaxed text-muted-foreground">
                  This invite was sent to{' '}
                  <span className="font-medium text-foreground">
                    {invite.email}
                  </span>{' '}
                  but you're authenticated as{' '}
                  <span className="font-medium text-foreground">
                    {currentUserEmail}
                  </span>
                </p>

                <div className="space-y-2">
                  <Button className="w-full" variant="secondary" asChild>
                    <a href="/api/auth/sign-out">
                      <LogOut className="mr-2 size-4" />
                      Sign out from {currentUserEmail}
                    </a>
                  </Button>

                  <Button className="w-full" variant="outline" asChild>
                    <Link href="/">Back to dashboard</Link>
                  </Button>
                </div>
              </div>
            )}
        </CardContent>
      </Card>
    </main>
  )
}
