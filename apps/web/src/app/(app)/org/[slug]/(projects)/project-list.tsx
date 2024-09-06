import { AvatarFallback } from '@radix-ui/react-avatar'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import { ArrowRight } from 'lucide-react'

import { getCurrentOrg } from '@/components/auth/auth'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { getProjects } from '@/http/get-projects'

dayjs.extend(relativeTime)

export async function ProjectList() {
  const currentOrg = await getCurrentOrg()
  const { projects } = await getProjects(currentOrg!)

  return (
    <div className="grid grid-cols-3 gap-4">
      {projects.map((project) => (
        <Card key={project.id} className="flex flex-col justify-between">
          <CardHeader>
            <CardTitle className="text-xl font-medium">
              {project.name}
            </CardTitle>
            <CardDescription className="line-clamp-2 leading-relaxed">
              {project.description}
            </CardDescription>
          </CardHeader>
          <CardFooter className="flex items-center gap-1.5">
            <Avatar className="size-4">
              {project.owner.avatarUrl && (
                <AvatarImage src={project.owner.avatarUrl} />
              )}
              <AvatarFallback />
            </Avatar>

            <span className="text-xs text-muted-foreground">
              <span className="font-medium text-foreground">
                {project.owner.name}
              </span>{' '}
              {dayjs(project.createdAt).fromNow()}
            </span>

            <Button size="xs" className="ml-auto" variant="outline">
              View <ArrowRight className="ml-2 size-3" />
            </Button>
          </CardFooter>
        </Card>
      ))}

      {(projects.length === 0 || !projects) && (
        <span className="text-muted-foreground">No project found!</span>
      )}
    </div>
  )
}
