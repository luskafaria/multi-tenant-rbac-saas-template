import { redirect } from 'next/navigation'

import { ability, getCurrentOrg } from '@/components/auth/auth'

import { ProjectForm } from './project-form'

export default async function CreateProject() {
  const permissions = await ability()
  const currentOrg = getCurrentOrg()

  if (permissions?.cannot('create', 'Project')) {
    redirect(`/org/${currentOrg}`)
  }

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Create project</h1>

      <ProjectForm />
    </div>
  )
}
