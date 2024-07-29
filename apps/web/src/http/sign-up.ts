import { api } from './api-client'

interface SignUpRequest {
  name: string
  email: string
  password: string
}

type SignUpResponse = Promise<void>

export async function signUp({
  name,
  email,
  password,
}: SignUpRequest): SignUpResponse {
  await api.post('users', {
    json: { name, email, password },
  })
}
