# Dependency Update Plan

This document outlines the update strategy for the RBAC project dependencies, organized by scope.

---

## 1. Global Updates

Updates that affect multiple packages or the entire monorepo.

### 1.1 Monorepo Tooling

- [x] Update `turbo` from `1.13.3` to `2.6.3`
- [x] Update `pnpm` from `9.1.1` to latest (in `package.json` packageManager field) (we choose to update it to 9.15.x for now since pnpm 10 has to many breaking changes.)
- [ ] Update Node.js engine requirement if needed

### 1.2 TypeScript & Linting

- [ ] Update `typescript` from `5.4.5` to `5.9.3`
- [ ] Update `@typescript-eslint/eslint-plugin` from `7.17.0` to `8.48.1`
- [ ] Update `@typescript-eslint/parser` from `7.17.0` to `8.48.1`
- [ ] Update `eslint-plugin-simple-import-sort` from `12.1.0` to `12.1.1`

### 1.3 Prettier

- [ ] Update `prettier` from `3.2.5` to `3.7.4`
- [ ] Update `prettier-plugin-tailwindcss` from `0.5.14` to `0.7.2`

### 1.4 Zod (Used in all packages)

- [x] Update `zod` from `3.23.8` to `4.1.13` in `@saas/api`
- [x] Update `zod` from `3.23.8` to `4.1.13` in `@saas/auth`
- [x] Update `zod` from `3.23.8` to `4.1.13` in `@saas/env`
- [x] Update `zod` from `3.23.8` to `4.1.13` in `@saas/web`
- [x] Fix any breaking changes in Zod schemas (added `"type": "module"` to `@saas/auth` for ESM compatibility)

### 1.5 Environment Package

- [ ] Update `@t3-oss/env-nextjs` from `0.10.1` to `0.13.8`
- [ ] Update `@types/node` from `20.14.2` to `24.10.1` in `@saas/env`

### 1.6 Auth Package

- [ ] Update `@casl/ability` from `6.7.1` to `6.7.3`

### 1.7 Shared Dev Dependencies

- [ ] Update `dotenv-cli` from `7.4.2` to `11.0.0`

---

## 2. Frontend Updates (@saas/web)

Updates specific to the Next.js frontend application.

### 2.1 React 19 Stable Release

- [ ] Update `react` from `19.0.0-rc-f994737d14-20240522` to `19.2.1`
- [ ] Update `react-dom` from `19.0.0-rc-f994737d14-20240522` to `19.2.1`
- [ ] Update `@types/react` from `18.3.3` to `19.2.7`
- [ ] Update `@types/react-dom` from `18.3.0` to `19.2.3`
- [ ] Test all React components for compatibility

### 2.2 Next.js 16

- [ ] Update `next` from `15.0.0-rc.0` to `16.0.7`
- [ ] Update `@next/eslint-plugin-next` from `14.2.6` to `16.0.7`
- [ ] Review Next.js 16 migration guide
- [ ] Update any deprecated APIs
- [ ] Test SSR/SSG functionality
- [ ] Test middleware behavior

### 2.3 Tailwind CSS 4

- [ ] Update `tailwindcss` from `3.4.10` to `4.1.17`
- [ ] Update `tailwind-merge` from `2.5.2` to `3.4.0`
- [ ] Migrate `tailwind.config.ts` to new format
- [ ] Update PostCSS configuration if needed
- [ ] Update `postcss` from `8.4.39` to `8.5.6`
- [ ] Test all UI components for styling issues

### 2.4 Radix UI Components

- [ ] Update `@radix-ui/react-avatar` from `1.1.0` to `1.1.11`
- [ ] Update `@radix-ui/react-checkbox` from `1.1.1` to `1.3.3`
- [ ] Update `@radix-ui/react-dialog` from `1.1.1` to `1.1.15`
- [ ] Update `@radix-ui/react-dropdown-menu` from `2.1.1` to `2.1.16`
- [ ] Update `@radix-ui/react-icons` from `1.3.0` to `1.3.2`
- [ ] Update `@radix-ui/react-label` from `2.1.0` to `2.1.8`
- [ ] Update `@radix-ui/react-popover` from `1.1.1` to `1.1.15`
- [ ] Update `@radix-ui/react-select` from `2.1.1` to `2.2.6`
- [ ] Update `@radix-ui/react-separator` from `1.1.0` to `1.1.8`
- [ ] Update `@radix-ui/react-slot` from `1.1.0` to `1.2.4`

### 2.5 Data Fetching & HTTP

- [ ] Update `@tanstack/react-query` from `5.55.0` to `5.90.12`
- [ ] Update `ky` from `1.7.2` to `1.14.0`

### 2.6 Utilities

- [ ] Update `cookies-next` from `4.2.1` to `6.1.1`
- [ ] Update `next-themes` from `0.3.0` to `0.4.6`
- [ ] Update `lucide-react` from `0.439.0` to `0.555.0`
- [ ] Update `dayjs` from `1.11.13` to `1.11.19`
- [ ] Update `class-variance-authority` from `0.7.0` to `0.7.1`
- [ ] Update `clsx` (check if update available)

---

## 3. Backend Updates (@saas/api)

Updates specific to the Fastify API backend.

### 3.1 Fastify 5 Migration

- [x] Update `fastify` from `4.28.1` to `5.6.2`
- [x] Review Fastify 5 migration guide
- [x] Update `fastify-plugin` from `4.5.1` to `5.1.0`
- [x] Update `fastify-type-provider-zod` from `2.0.0` to `6.1.0`
- [x] Fix any breaking changes in route definitions (fixed `bcryptjs` ESM imports)
- [x] Test all API endpoints (server starts, Swagger docs verified)

### 3.2 Fastify Plugins

- [x] Update `@fastify/cors` from `9.0.1` to `11.1.0`
- [x] Update `@fastify/jwt` from `8.0.1` to `10.0.0`
- [x] Update `@fastify/swagger` from `8.15.0` to `9.6.1`
- [x] Update `@fastify/swagger-ui` from `4.1.0` to `5.2.3`
- [ ] Test JWT authentication flow (requires database)
- [ ] Test CORS configuration (requires frontend testing)
- [x] Verify Swagger documentation (OpenAPI spec verified at `/docs/json`)

### 3.3 Prisma 7

- [ ] Update `prisma` from `5.19.1` to `7.1.0`
- [ ] Update `@prisma/client` from `5.19.1` to `7.1.0`
- [ ] Review Prisma 7 migration guide
- [ ] Run `prisma generate` after update
- [ ] Test database operations
- [ ] Verify migrations work correctly

### 3.4 Security & Utilities

- [ ] Update `bcryptjs` from `2.4.3` to `3.0.3`
- [ ] Update `@types/bcryptjs` from `2.4.6` to `3.0.0`
- [ ] Test password hashing/verification

### 3.5 Dev Dependencies

- [ ] Update `@faker-js/faker` from `9.0.0` to `10.1.0`
- [ ] Update `tsx` from `4.19.0` to `4.21.0`
- [ ] Update seed scripts if Faker API changed

---

## 4. Verification Checklist

After completing updates, verify the following:

### 4.1 Build & Lint

- [ ] Run `pnpm install` successfully
- [ ] Run `pnpm build` without errors
- [ ] Run `pnpm lint` without errors
- [ ] Run `pnpm dev` and verify hot reload works

### 4.2 Database

- [ ] Run `pnpm db:migrate` successfully
- [ ] Run `pnpm db:seed` successfully
- [ ] Verify database operations in Prisma Studio

### 4.3 API Testing

- [ ] Test user authentication (password & GitHub)
- [ ] Test organization CRUD operations
- [ ] Test project CRUD operations
- [ ] Test invite flow (create, accept, reject, revoke)
- [ ] Test member management
- [ ] Test billing endpoint
- [ ] Verify Swagger UI at `/docs`

### 4.4 Frontend Testing

- [ ] Test sign-in flow
- [ ] Test sign-up flow
- [ ] Test organization creation and management
- [ ] Test project creation
- [ ] Test member invitations
- [ ] Test theme switching
- [ ] Verify responsive design
- [ ] Test all form validations

---

## 5. Rollback Plan

If critical issues arise:

1. Revert `pnpm-lock.yaml` to previous version
2. Revert changed `package.json` files
3. Run `pnpm install`
4. Document the issue for future reference

---

## Notes

- **Zod 4** is a significant update used across all packages - update and test thoroughly
- **Fastify 5** has breaking changes in plugin registration and hooks
- **Next.js 16** may have changes in App Router behavior
- **Tailwind CSS 4** has a completely new configuration format
- **Prisma 7** may require schema adjustments

## Estimated Risk Levels

| Update Group         | Risk Level | Reason                       |
| -------------------- | ---------- | ---------------------------- |
| Global: Tooling      | Low        | Mostly dev tools             |
| Global: Zod          | High       | Used everywhere, API changes |
| Frontend: React 19   | Medium     | Moving from RC to stable     |
| Frontend: Next.js 16 | High       | Major version jump           |
| Frontend: Tailwind 4 | High       | New config format            |
| Frontend: Radix UI   | Low        | Minor updates                |
| Backend: Fastify 5   | High       | Breaking changes expected    |
| Backend: Prisma 7    | High       | Database layer changes       |
