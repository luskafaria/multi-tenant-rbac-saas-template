import nodeConfig from '@saas/eslint-config/node'

/** @type {import('eslint').Linter.Config[]} */
export default [
  ...nodeConfig,
  {
    ignores: ['node_modules/', 'dist/', '.prisma/'],
  },
]
