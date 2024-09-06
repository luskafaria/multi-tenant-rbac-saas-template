import type { Role } from '@saas/auth'

import { api } from './api-client'

interface CreateInviteRequest {
  org: string
  email: string
  role: Role
}

export async function createInvite({
  org,
  email,
  role,
}: CreateInviteRequest): Promise<void> {
  await api.post(`organizations/${org}/invites`, {
    json: {
      email,
      role,
    },
  })
}
