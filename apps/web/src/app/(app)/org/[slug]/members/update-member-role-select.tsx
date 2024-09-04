'use client'

import type { Role } from '@saas/auth'
import type { ComponentProps } from 'react'

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

import { updateMemberAction } from './actions'

interface UpdateMemberRoleSelectProps extends ComponentProps<typeof Select> {
  memberId: string
}

export function UpdateMemberRoleSelect({
  memberId,
  ...props
}: UpdateMemberRoleSelectProps) {
  const roleArr: Role[] = ['ADMIN', 'MEMBER', 'BILLING']

  async function updateMemberRole(role: Role) {
    await updateMemberAction(memberId, role)
  }

  return (
    <Select onValueChange={updateMemberRole} {...props}>
      <SelectTrigger className="h-8 w-32 capitalize">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        {roleArr.map((item) => {
          return (
            <SelectItem className="capitalize" value={item} key={item}>
              {item.toLowerCase()}
            </SelectItem>
          )
        })}
      </SelectContent>
    </Select>
  )
}
