import bcrypt from 'bcryptjs'
import type { FastifyInstance } from 'fastify'
import type { ZodTypeProvider } from 'fastify-type-provider-zod'
import z from 'zod'

import { security } from '@/config/security'
import { prisma } from '@/lib/prisma'

import { UnauthorizedError } from '../_errors/unauthorized-error'

export async function resetPassword(app: FastifyInstance) {
  app.withTypeProvider<ZodTypeProvider>().post(
    '/password/reset',
    {
      schema: {
        tags: ['auth'],
        summary: 'Get password recover',
        body: z.object({
          code: z.string(),
          password: z.string().min(6),
        }),
        response: {
          204: z.null(),
        },
      },
    },
    async (request, reply) => {
      const { code, password } = request.body

      const tokenByCode = await prisma.token.findUnique({
        where: {
          id: code,
        },
      })

      if (!tokenByCode) {
        throw new UnauthorizedError()
      }

      const expiryTime = new Date(tokenByCode.createdAt)
      expiryTime.setMinutes(
        expiryTime.getMinutes() + security.PASSWORD_RESET_TOKEN_EXPIRY_MINUTES
      )

      if (new Date() > expiryTime) {
        await prisma.token.delete({ where: { id: code } })
        throw new UnauthorizedError()
      }

      const passwordHash = await bcrypt.hash(
        password,
        security.BCRYPT_COST_FACTOR
      )

      await prisma.$transaction([
        prisma.user.update({
          where: {
            id: tokenByCode.userId,
          },
          data: {
            passwordHash,
          },
        }),
        prisma.token.delete({
          where: {
            id: code,
          },
        }),
      ])

      return reply.status(204).send()
    }
  )
}
