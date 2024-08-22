/** @type {import('eslint').Linter.Config} */
module.exports = {
  parser: '@typescript-eslint/parser',
  extends: ['@rocketseat/eslint-config/next', 'plugin:@next/next/recommended'],
  plugins: ['simple-import-sort', '@typescript-eslint'],
  rules: {
    'simple-import-sort/imports': 'error',
  },
}
