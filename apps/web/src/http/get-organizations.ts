import { api } from './api-client';

interface GetOrganizationsResponse {
  organizations: {
    slug: string;
    id: string;
    name: string;
    avatarUrl: string | null;
  }[]
}

export async function getOrganizations() {
  const result = await api.get('organizations').json<GetOrganizationsResponse>()

  return result
}
