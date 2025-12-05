# Multi-Tenant RBAC SaaS Template

This project contains all the necessary boilerplate to setup a multi-tenant SaaS with Next.js including authentication and RBAC authorization.

## Architecture

### System Overview

```mermaid
graph TB
    subgraph Client
        Browser[Browser]
    end

    subgraph "Frontend (apps/web)"
        NextJS[Next.js 15]
        ServerActions[Server Actions]
        ServerComponents[Server Components]
    end

    subgraph "Backend (apps/api)"
        Fastify[Fastify Server]
        Routes[API Routes]
        AuthMiddleware[Auth Middleware]
        Swagger[Swagger Docs]
    end

    subgraph "Shared Packages"
        AuthPkg["@saas/auth<br/>(CASL RBAC)"]
        EnvPkg["@saas/env<br/>(Env Config)"]
    end

    subgraph "Data Layer"
        Prisma[Prisma ORM]
        PostgreSQL[(PostgreSQL)]
    end

    subgraph "External"
        GitHub[GitHub OAuth]
    end

    Browser --> NextJS
    NextJS --> ServerActions
    NextJS --> ServerComponents
    ServerActions --> Fastify
    ServerComponents --> Fastify
    Fastify --> Routes
    Routes --> AuthMiddleware
    AuthMiddleware --> AuthPkg
    Routes --> Prisma
    Prisma --> PostgreSQL
    Fastify --> GitHub
    Fastify --> Swagger
    NextJS -.-> AuthPkg
    NextJS -.-> EnvPkg
    Fastify -.-> EnvPkg
```

### Data Model

```mermaid
erDiagram
    User ||--o{ Member : "has memberships"
    User ||--o{ Account : "OAuth accounts"
    User ||--o{ Token : "password recovery"
    User ||--o{ Invite : "creates invites"
    User ||--o{ Organization : "owns"
    User ||--o{ Project : "owns"

    Organization ||--o{ Member : "has members"
    Organization ||--o{ Project : "has projects"
    Organization ||--o{ Invite : "has invites"

    User {
        uuid id PK
        string email UK
        string name
        string passwordHash
        string avatarUrl
    }

    Organization {
        uuid id PK
        string name
        string slug UK
        string domain UK
        uuid ownerId FK
    }

    Member {
        uuid id PK
        enum role "ADMIN|MEMBER|BILLING"
        uuid organizationId FK
        uuid userId FK
    }

    Project {
        uuid id PK
        string name
        string slug UK
        string description
        uuid organizationId FK
        uuid ownerId FK
    }

    Invite {
        uuid id PK
        string email
        enum role
        uuid organizationId FK
        uuid authorId FK
    }

    Account {
        uuid id PK
        enum provider "GITHUB"
        string providerAccountId UK
        uuid userId FK
    }

    Token {
        uuid id PK
        enum type "PASSWORD_RECOVER"
        uuid userId FK
    }
```

### Authentication Flow

```mermaid
sequenceDiagram
    participant U as User
    participant W as Web (Next.js)
    participant A as API (Fastify)
    participant DB as PostgreSQL
    participant GH as GitHub

    rect rgb(40, 40, 40)
        Note over U,DB: Email/Password Authentication
        U->>W: Submit credentials
        W->>A: POST /sessions/password
        A->>DB: Find user by email
        DB-->>A: User record
        A->>A: Verify password (bcrypt)
        A-->>W: JWT Token (7 days)
        W->>W: Store in cookies
        W-->>U: Redirect to dashboard
    end

    rect rgb(40, 40, 40)
        Note over U,GH: GitHub OAuth Authentication
        U->>W: Click "Sign in with GitHub"
        W->>GH: Redirect to OAuth
        GH-->>W: Authorization code
        W->>A: POST /sessions/github
        A->>GH: Exchange code for token
        GH-->>A: Access token
        A->>GH: GET /user
        GH-->>A: User profile
        A->>DB: Upsert user + account
        A-->>W: JWT Token (7 days)
        W-->>U: Redirect to dashboard
    end
```

### Authorization Flow (RBAC)

```mermaid
flowchart TD
    A[API Request] --> B{Has JWT?}
    B -->|No| C[401 Unauthorized]
    B -->|Yes| D[Verify JWT]
    D -->|Invalid| C
    D -->|Valid| E[Extract User ID]
    E --> F[Get Membership for Org]
    F -->|Not Member| G[403 Forbidden]
    F -->|Is Member| H[Build CASL Ability]
    H --> I{Check Permission}
    I -->|Denied| G
    I -->|Allowed| J[Execute Action]
    J --> K[Return Response]

    subgraph "CASL Ability Builder"
        H --> L[Get User Role]
        L --> M[Apply Role Permissions]
        M --> N[Apply Conditions]
    end
```

### Multi-Tenancy Model

```mermaid
flowchart TB
    subgraph "Tenant Isolation"
        U1[User] --> M1[Membership]
        U1 --> M2[Membership]

        M1 --> O1[Org: Acme Corp]
        M2 --> O2[Org: Beta Inc]

        O1 --> P1[Project A]
        O1 --> P2[Project B]
        O2 --> P3[Project C]
    end

    subgraph "Role-Based Access"
        M1 -.- R1["Role: ADMIN<br/>Full access"]
        M2 -.- R2["Role: MEMBER<br/>Limited access"]
    end

    subgraph "Ownership"
        U1 -.->|owns| O1
        U1 -.->|owns| P1
    end
```

## Project Structure

```
├── apps/
│   ├── web/          # Next.js frontend
│   └── api/          # Fastify backend
├── packages/
│   ├── auth/         # RBAC utilities (CASL)
│   └── env/          # Environment configuration
└── config/
    ├── eslint-config/
    ├── prettier/
    └── typescript-config/
```

## Getting Started

### Prerequisites

- Node.js >= 18
- pnpm 9.15.9
- Docker (for PostgreSQL)

### Setup

```bash
# Install dependencies
pnpm install

# Start PostgreSQL
docker compose up -d

# Run database migrations
pnpm --filter @saas/api db:migrate

# Seed the database (optional)
pnpm --filter @saas/api db:seed

# Start development servers
pnpm dev
```

## Scripts

| Command          | Description                        |
| ---------------- | ---------------------------------- |
| `pnpm dev`       | Start all apps in development mode |
| `pnpm build`     | Build all apps                     |
| `pnpm lint`      | Run ESLint across all apps         |
| `pnpm typecheck` | Run TypeScript type checking       |

## CI/CD

GitHub Actions workflows run on every push and PR to `main`:

| Workflow     | Checks                  |
| ------------ | ----------------------- |
| **CI - Web** | Lint, Type Check, Build |
| **CI - API** | Lint, Type Check        |

Workflows use path filtering to only run when relevant files change.

## RBAC

Roles & permissions.

### Roles

- Owner (count as administrator)
- Administrator
- Member
- Billing (one per organization)
- Anonymous

### Permissions table

|                        | Administrator | Member | Billing | Anonymous |
| ---------------------- | ------------- | ------ | ------- | --------- |
| Update organization    | ✅            | ❌     | ❌      | ❌        |
| Delete organization    | ✅            | ❌     | ❌      | ❌        |
| Invite a member        | ✅            | ❌     | ❌      | ❌        |
| Revoke an invite       | ✅            | ❌     | ❌      | ❌        |
| List members           | ✅            | ✅     | ✅      | ❌        |
| Transfer ownership     | ⚠️            | ❌     | ❌      | ❌        |
| Update member role     | ✅            | ❌     | ❌      | ❌        |
| Delete member          | ✅            | ⚠️     | ❌      | ❌        |
| List projects          | ✅            | ✅     | ✅      | ❌        |
| Create a new project   | ✅            | ✅     | ❌      | ❌        |
| Update a project       | ✅            | ⚠️     | ❌      | ❌        |
| Delete a project       | ✅            | ⚠️     | ❌      | ❌        |
| Get billing details    | ✅            | ❌     | ✅      | ❌        |
| Export billing details | ✅            | ❌     | ✅      | ❌        |

> ✅ = allowed
> ❌ = not allowed
> ⚠️ = allowed w/ conditions

#### Conditions

- Only owners may transfer organization ownership;
- Only administrators and project authors may update/delete the project;
- Members can leave their own organization;

## Features

### Authentication

- [x] It should be able to authenticate using e-mail & password;
- [x] It should be able to authenticate using Github account;
- [x] It should be able to recover password using e-mail;
- [x] It should be able to create an account (e-mail, name and password);

### Organizations

- [x] It should be able to create a new organization;
- [x] It should be able to get organizations to which the user belongs;
- [x] It should be able to update an organization;
- [x] It should be able to shutdown an organization;
- [x] It should be able to transfer organization ownership;

### Invites

- [x] It should be able to invite a new member (e-mail, role);
- [x] It should be able to accept an invite;
- [x] It should be able to revoke a pending invite;

### Members

- [x] It should be able to get organization members;
- [x] It should be able to update a member role;

### Projects

- [x] It should be able to get projects within a organization;
- [x] It should be able to create a new project (name, url, description);
- [x] It should be able to update a project (name, url, description);
- [x] It should be able to delete a project;

### Billing

- [x] It should be able to get billing details for organization ($20 per project / $10 per member excluding billing role);
