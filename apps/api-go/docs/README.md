# Go API Migration Plan

This document outlines the plan for migrating the existing Node.js/TypeScript backend to Go.

## Current Status

### ✅ Infrastructure & Setup
- A new Go service has been created in `apps/api-go`.
- The basic project structure is in place.
- Environment variables are loaded using `viper`.
- A database connection is established using the standard `database/sql` package.
- A health check endpoint (`/v1/health`) has been created using the `chi` router.
- Docker integration (Dockerfile and docker-compose.yml) has been added for the Go service and its dedicated PostgreSQL database.
- A `Makefile` has been added for development, including `dev` (with Air for hot-reloading), `db-up`, and `db-seed` commands.
- Swagger API documentation has been integrated, accessible at `/docs`.
- `sqlc` has been set up to generate a type-safe data access layer from the Prisma-generated SQL schema.
- Coverage reports are accessible at `/coverage`.
- **Database seeding** - Comprehensive seed script (`cmd/seed/main.go`) creates test data matching Node.js seeds.

### ✅ Authentication Endpoints Migrated
- **POST `/v1/users`** - Create user account with auto-join organization logic ✅
- **POST `/v1/sessions/password`** - Password-based authentication ✅
- **POST `/v1/sessions/github`** - GitHub OAuth authentication ✅
- **GET `/v1/profile`** - Get authenticated user profile ✅

## Development Workflow

### **Quick Start:**
```bash
# Start database
make db-up

# Seed with test data
make db-seed

# Start development server
make dev
```

### **Test Data Available:**
After seeding, you can test with these accounts:
- **john@gmail.com** / `123456` (Admin in Acme Admin org)
- **jane@example.com** / `123456` (Admin in Acme Member org) 
- **bob@example.com** / `123456` (Member in various orgs)

### **API Documentation:**
- **Swagger UI:** http://localhost:3334/docs
- **Coverage Report:** http://localhost:3334/coverage

## Remaining Endpoints to Migrate

### ❌ Authentication (Missing 2 endpoints)
- **POST `/password/recover`** - Request password recovery token
- **POST `/password/reset`** - Reset password using token

### ❌ Organization Management (Missing 8 endpoints)
- **POST `/v1/organizations`** - Create new organization with RBAC checks
- **GET `/v1/organizations`** - List user's organizations with roles
- **GET `/v1/organizations/:slug`** - Get organization details
- **PUT `/v1/organizations/:slug`** - Update organization (name, domain, auto-join)
- **DELETE `/v1/organizations/:slug`** - Shutdown/delete organization
- **GET `/v1/organizations/:slug/membership`** - Get user's membership details
- **PATCH `/v1/organizations/:slug/owner`** - Transfer organization ownership
- **GET `/v1/organizations/:slug/billing`** - Get billing information (seats/projects)

### ❌ Member Management (Missing 3 endpoints)
- **GET `/v1/organizations/:slug/members`** - List organization members with roles
- **PUT `/v1/organizations/:slug/members/:memberId`** - Update member role
- **DELETE `/v1/organizations/:slug/members/:memberId`** - Remove member from organization

### ❌ Invite Management (Missing 7 endpoints)
- **POST `/v1/organizations/:slug/invites`** - Create organization invite
- **GET `/v1/organizations/:slug/invites`** - List organization invites
- **DELETE `/v1/organizations/:slug/invites/:inviteId`** - Revoke invite
- **GET `/v1/invites/:inviteId`** - Get invite details
- **POST `/v1/invites/:inviteId/accept`** - Accept invite
- **POST `/v1/invites/:inviteId/reject`** - Reject invite
- **GET `/v1/pending-invites`** - List user's pending invites

### ❌ Project Management (Missing 5 endpoints)
- **POST `/v1/organizations/:slug/projects`** - Create project within organization
- **GET `/v1/organizations/:slug/projects`** - List organization projects
- **GET `/v1/organizations/:slug/projects/:projectSlug`** - Get project details
- **PATCH `/v1/organizations/:slug/projects/:projectId`** - Update project
- **DELETE `/v1/organizations/:slug/projects/:projectId`** - Delete project

## Migration Strategy

### Phase 1: Core Infrastructure (HIGH PRIORITY)
1. **Implement RBAC system** - Port the Node.js RBAC logic using either `casbin` or custom implementation
2. **Password recovery/reset** - Complete the authentication system
3. **Organization management** - Core organization CRUD operations

### Phase 2: Member & Invite System (MEDIUM PRIORITY)
1. **Member management** - List, update, and remove members
2. **Invite system** - Create, list, and revoke invites
3. **Invite workflow** - Accept/reject invite functionality

### Phase 3: Project Management (MEDIUM PRIORITY)
1. **Project CRUD** - Create, read, update, delete projects
2. **Project listing** - Organization-scoped project management

### Phase 4: Advanced Features (LOW PRIORITY)
1. **Billing calculations** - Seats and project-based billing
2. **Organization transfer** - Ownership management
3. **Organization shutdown** - Complete organization lifecycle

## Technical Requirements

### RBAC Implementation
The Node.js app uses a sophisticated RBAC system with:
- **Roles:** ADMIN, MEMBER, BILLING
- **Resources:** User, Project, Invite, Organization, Billing
- **Permissions:** create, read, update, delete, manage, get, export

### Database Operations
Many endpoints require complex transactions:
- Organization creation with automatic membership
- Invite acceptance (create member + delete invite)
- Organization transfer (update roles + ownership)
- Member management with authorization checks

### Business Logic
- Auto-join organizations by email domain
- Slug generation for organizations and projects
- Billing calculations: seats × $10 + projects × $20
- Email validation and conflict detection

**Total Missing Endpoints: 25**  
**Current Progress: 4/29 endpoints (14% complete)**