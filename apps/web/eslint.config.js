import nextConfig from '@saas/eslint-config/next'

/** @type {import('eslint').Linter.Config[]} */
export default [
  ...nextConfig,
  {
    ignores: ['node_modules/', '.next/', 'out/'],
  },
]
