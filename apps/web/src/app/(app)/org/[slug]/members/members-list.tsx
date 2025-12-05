import { organizationSchema } from '@saas/auth'
import { ArrowLeftRight, Crown, UserMinus } from 'lucide-react'
import Image from 'next/image'

import { ability, getCurrentOrg } from '@/components/auth/auth'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Table, TableBody, TableCell, TableRow } from '@/components/ui/table'
import { getMembers } from '@/http/get-members'
import { getMembership } from '@/http/get-membership'
import { getOrganization } from '@/http/get-organization'

import { removeMemberAction } from './actions'
import { UpdateMemberRoleSelect } from './update-member-role-select'

export async function MembersList() {
  const currentOrg = await getCurrentOrg()

  const permissions = await ability()

  const [{ membership }, { members }, { organization }] = await Promise.all([
    getMembership(currentOrg!),
    getMembers(currentOrg!),
    getOrganization(currentOrg!),
  ])

  const authOrganization = organizationSchema.parse(organization)

  const canTransferOwnership = permissions?.can(
    'transfer_ownership',
    authOrganization
  )
  const canDeleteMember = permissions?.can('delete', 'User')
  const cannotUpdateMember = permissions?.cannot('update', 'User')

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-semibold">Members</h1>

      <div className="rounded border">
        <Table>
          <TableBody>
            {members.map((member) => {
              const isItMe =
                member.userId === membership.userId ||
                member.userId === organization.ownerId

              return (
                <TableRow key={member.id}>
                  <TableCell className="py-2.5" style={{ width: 48 }}>
                    <Avatar>
                      <AvatarFallback />
                      {member.avatarUrl && (
                        <Image
                          className="aspect-square size-full"
                          src={member.avatarUrl}
                          width={32}
                          height={32}
                          alt={member.name!}
                        />
                      )}
                    </Avatar>
                  </TableCell>
                  <TableCell className="py-2.5">
                    <div className="flex flex-col">
                      <span className="inline-flex items-center gap-2 font-medium">
                        {member.name}
                        {member.userId === membership.userId && ' (me)'}
                        {organization.ownerId === member.userId && (
                          <span className="text-muted-foreground inline-flex items-center gap-1 text-xs">
                            <Crown className="size-3 text-yellow-500" /> Owner
                          </span>
                        )}
                      </span>
                      <span className="text-muted-foreground text-xs">
                        {member.email}
                      </span>
                    </div>
                  </TableCell>

                  <TableCell className="py-2.5">
                    <div className="flex items-center justify-end gap-2">
                      {canTransferOwnership && (
                        <Button size="sm" variant="ghost">
                          <ArrowLeftRight className="mr-2 size-4" />
                          Transfer ownership
                        </Button>
                      )}

                      <UpdateMemberRoleSelect
                        memberId={member.id}
                        disabled={isItMe || cannotUpdateMember}
                        value={member.role}
                      />

                      {canDeleteMember && (
                        <form action={removeMemberAction.bind(null, member.id)}>
                          <Button
                            disabled={isItMe}
                            type="submit"
                            size="sm"
                            variant="destructive"
                          >
                            <UserMinus className="mr-2 size-4" />
                            Remove
                          </Button>
                        </form>
                      )}
                    </div>
                  </TableCell>
                </TableRow>
              )
            })}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
