import type { FastifyInstance } from 'fastify'
import { ZodError } from 'zod'
import { BadRequestError } from './routes/_errors/bad-request-error'
import { UnauthorizedError } from './routes/_errors/unauthorized-error'

type FastifyErrorHandler = FastifyInstance['errorHandler']

export const errorHandler: FastifyErrorHandler = (error, req, rep) => {
  if (error instanceof ZodError) {
    rep.status(400).send({
      message: 'Validation error',
      errors: error.flatten().fieldErrors,
    })
  }

  if (error instanceof BadRequestError) {
    return rep.status(400).send({
      message: error.message,
    })
  }

  if (error instanceof UnauthorizedError) {
    return rep.status(401).send({
      message: error.message,
    })
  }

  // TODO: integrate w/ a watch tool 🕵️
  console.error(error)

  return rep.status(500).send({ message: 'Internal server error' })
}
