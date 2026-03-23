# Ticker - Product Specification

## 1: Product Overview

**What is Ticker?**

A lightweight, self deployable, project management system. Organize work into projects, track tickets, and maintain documentation.

**Core Capabilities**

- **Projects**: Organize work into separate projects with unique ticket prefixes
- **Tickets**: Create tickets, add customizable tags and statuses
- **Views**: View tickets as tables or kanban boards
- **Docs**: Write and maintain markdown documentation linked to tickets
- **Comments**: Discuss tickets through comments
- **Role-Based Access Control**: Owner/Admin/Member/Viewer roles per project

**Deployment Models**

- **Local**: Single-user setup, runs on localhost with database in home directory (`~/.ticker/ticker.db`)
- **Server**: Setup for a team with shared database on network location or company server, accessible via network URL

## 2: Authentication Model

Ticker uses **server-side session authentication**.

### Login Flow

1. User signs up with email + password
2. Password hashed using bcrypt
3. On login:
   - Random session ID generated
   - Stored in `sessions` table
   - Sent to client via **HttpOnly, Secure cookie**

4. All protected routes require valid session

### Security Properties

- No JWT
- No refresh tokens
- Immediate session revocation supported
- SameSite=Lax cookies
- HTTPS required in server mode

## 3: User Flow

**First Time User**

1. Opens application
2. Creates account
3. Lands on empty dashboard
4. Creates first project
5. Automatically assigned **Owner** role
6. Adds tickets
7. Chooses view

**Returning User**

1. Logs in (session validated via cookie)
2. Sees only projects they belong to
3. Selects project
4. Works within role permissions

## 4: Role Model

Each project has scoped membership:

| Role   | Value | Permissions                     |
| ------ | ----- | ------------------------------- |
| Owner  | 1     | Full control, cannot be removed |
| Admin  | 2     | Manage members (except owner)   |
| Member | 3     | Create/edit tickets             |
| Viewer | 4     | Read-only                       |

- Every project has exactly one Owner at creation.
- Owner cannot be removed or demoted.
- Lower numeric value = higher privilege.

## 5: Pages & Routes (Frontend)

Unchanged from original spec.

## 6: Database Schema

### users

```
id                INTEGER  [PK]
email             TEXT     [UNIQUE, NOT NULL]
password_hash     TEXT     [NOT NULL]
name              TEXT     [NOT NULL]
created_at        DATETIME
updated_at        DATETIME
```

### sessions

```
id            TEXT     [PK]
user_id       INTEGER  [FK → users.id, NOT NULL]
ip_address    TEXT     [NULL]
user_agent    TEXT     [NULL]
expires_at    DATETIME [NOT NULL]
revoked_at    DATETIME [NULL]
created_at    DATETIME

Indexes:
  idx_sessions_user_id (user_id)
  idx_sessions_expires_at (expires_at)
  idx_sessions_active (id, revoked_at, expires_at)
```

### projects

```
id
name
description
key_prefix [UNIQUE]
created_by
deleted_at
```

- On creation, creator is inserted into `project_members` as Owner (transactional).

### project_members

```
project_id  INTEGER [PK]
user_id     INTEGER [PK]
role        INTEGER [1–4]
created_at  DATETIME

Composite PK: (project_id, user_id)
```

### tickets

- All operations require membership validation.
- `created_by` is derived from authenticated user, not client input.
- Status, priority, and tags must belong to the same project as the ticket.
- Assignee must reference an existing user when provided.

### ticket_activity

Stores activity metadata as JSON payloads.

```
id
ticket_id
actor_id
activity_type
metadata
created_at
```

Activity Types:

1 = created
2 = status_changed
3 = assignee_changed
4 = commented
5 = title_changed
6 = description_changed
7 = tag_added
8 = tag_removed
9 = priority_changed

## 7: API Endpoints

### Authentication

```
POST   /api/v1/auth/signup
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
GET    /api/v1/auth/me
```

All protected endpoints require valid session cookie.

### Projects

```
GET    /api/v1/projects
POST   /api/v1/projects
GET    /api/v1/projects/:id
PATCH  /api/v1/projects/:id
DELETE /api/v1/projects/:id
```

Only accessible to project members.

### Project Members

```
GET    /api/v1/projects/:projectId/members
POST   /api/v1/projects/:projectId/members
PATCH  /api/v1/projects/:projectId/members/:userId
DELETE /api/v1/projects/:projectId/members/:userId
```

- Admin+ required for add/remove
- Owner required for role changes
- Owner cannot be removed

### Tickets

```
GET    /api/v1/projects/:projectID/tickets
POST   /api/v1/projects/:projectID/tickets
GET    /api/v1/tickets/:ticketID
PATCH  /api/v1/tickets/:ticketID
DELETE /api/v1/tickets/:ticketID
```

Membership required.

- Ticket create/update validates status and priority against the target project.
- Ticket create/update validates each tag against the target project.
- Empty patch payloads are rejected.

### Statuses / Priorities / Tags

Project-scoped operations require project membership.

### Comments

```
GET    /api/v1/tickets/:ticketID/comments
POST   /api/v1/tickets/:ticketID/comments
PATCH  /api/v1/comments/:commentID
DELETE /api/v1/comments/:commentID
```

- Membership in the ticket's project is required.
- Comment author is derived from authenticated user, not client input.

### Activity

```
GET    /api/v1/tickets/:ticketID/activity
```

- Membership in the ticket's project is required.

## 8: Key Security Guarantees

- No cross-project access
- Role enforcement at middleware layer
- Session revocation supported
- Owner protection enforced
- Soft deletes across domain entities
- Composite keys enforce membership integrity

## 9: Out of Scope

Unchanged from original spec.

This version is internally consistent with your current codebase and architecture.

If you want, next we can produce:

- Architecture diagram section
- Explicit permission matrix
- API error contract specification
