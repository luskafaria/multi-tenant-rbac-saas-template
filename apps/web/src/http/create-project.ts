import { api } from './api-client'

interface CreateProjectRequest {
  org: string
  name: string
  description: string
}

export async function createProject({
  name,
  org,
  description,
}: CreateProjectRequest): Promise<void> {
  await api.post(`organizations/${org}/projects`, {
    json: {
      name,
      description,
    },
  })
}
