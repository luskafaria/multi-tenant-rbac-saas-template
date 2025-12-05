import { Plus } from 'lucide-react'
import Link from 'next/link'

import { ability, getCurrentOrg } from '@/components/auth/auth'
import { Button } from '@/components/ui/button'

import { ProjectList } from './project-list'

export default async function Home() {
  const currentOrg = await getCurrentOrg()
  const permissions = await ability()

  const canCreateProject = permissions?.can('create', 'Project')
  const canGetProject = permissions?.can('get', 'Project')

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Projects</h1>

        {canCreateProject && (
          <Button size="sm" asChild>
            <Link href={`/org/${currentOrg}/create-project`}>
              <Plus className="mr-2 size-4" />
              Create project
            </Link>
          </Button>
        )}
      </div>

      {canGetProject ? (
        <ProjectList />
      ) : (
        <p className="text-muted-foreground text-sm">
          You're not allowed to see organization projects.
        </p>
      )}
    </div>
  )
}
