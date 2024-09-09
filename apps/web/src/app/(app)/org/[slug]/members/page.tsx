import { ability } from '@/components/auth/auth'

import { Invites } from './invites'
import { MembersList } from './members-list'

export default async function MembersPage() {
  const permissions = await ability()

  const canGetInvites = permissions?.can('get', 'Invite')
  const canGetUsers = permissions?.can('get', 'User')

  return (
    <div className="space-y-4">
      {canGetInvites && <Invites />}
      {canGetUsers && <MembersList />}
    </div>
  )
}
