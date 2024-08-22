import { redirect } from 'next/navigation'

import { ability, getCurrentOrg } from '@/auth/auth'
import { Header } from '@/components/header'

import { ProjectForm } from './project-form'

export default async function CreateProject() {
  const permissions = await ability()
  const currentOrg = getCurrentOrg()

  if (permissions?.cannot('create', 'Project')) {
    redirect(`/org/${currentOrg}`)
  }

  return (
    <div className="space-y-4 py-4">
      <Header />
      <main className="mx-auto w-full max-w-[1200px] space-y-4">
        <h1 className="text-2xl font-bold">Create project</h1>

        <ProjectForm />
      </main>
    </div>
  )
}
