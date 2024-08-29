import { api } from './api-client'

export async function shutdownOrganization(org: string) {
  await api.delete(`organizations/${org}`)
}
