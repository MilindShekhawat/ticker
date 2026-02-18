# Ticker - Product Specification

## 1: Product Overview

**What is Ticker?**

A lightweight, self deployable, project management system. Organize work into projects, track tickets, and maintain documentation.

**Core Capabilities**

- **Projects**: Organize work into separate projects with unique ticket prefixes
- **Tickets**: Create tickets, add customizable tags and statuses
- **Views**: View tickets as tables or kanban boards
- **Docs**: Write and maintain markdown documentation linked to tickets
- **Comments**: Discuss tickets through threaded comments

**Deployment Models**

- **Local**: Single-user setup, runs on localhost with database in home directory (`~/.ticker/ticker.db`)
- **Server**: Setup for a team with shared database on network location or company server, accessible via network URL

## 2: User Flow

**First Time User**

1. Opens application
2. Creates account (email + password)
3. Lands on empty dashboard
4. Creates first project (must select/create at least 1 status and 1 priority)
5. Adds tickets to project
6. Chooses view (table or kanban)

**Returning User**

1. Logs in
2. Sees project list on dashboard
3. Selects a project
4. Views tickets in last-used view (table/kanban)
5. Creates/edits tickets, adds comments, writes docs
6. Switches between projects as needed
   **Added route for independent docs:**

## Section 3: Pages & Routes

```
/                                          → Login/Signup
/dashboard                                 → Redirects to /dashboard/projects
/dashboard/projects                        → Project list
/dashboard/projects/:id/tickets            → Ticket table view
/dashboard/projects/:id/kanban             → Kanban board view
/dashboard/projects/:id/tickets/:ticketId  → Ticket detail
/dashboard/projects/:id/docs               → Project docs list
/dashboard/projects/:id/docs/:docId        → Project doc viewer/editor
/dashboard/settings                        → User settings
```

## Section 4: Page Details

### Login/Signup (`/`)

**Layout:** Centered card on clean background

**Components:**

- Toggle between Login/Signup tabs
- Email input
- Password input
- Name input (signup only)
- Submit button

**Actions:**

- Validate inputs
- Call auth API
- Store token
- Redirect to projects list

---

### Projects List (`/dashboard/projects`)

**Layout:** Header + grid of project cards

**Components:**

- Header (logo, search, user menu)
- "New Project" button
- Project cards showing:
  - Name
  - Description
  - Ticket stats (open/closed counts)

**Actions:**

- Click card → open project ticket view
- Create new project (modal)
- Search projects

---

### Ticket Table View (`/dashboard/projects/:id/tickets`)

**Layout:** Header + filters + table

**Components:**

- Project name
- View switcher (Table ↔ Kanban)
- "New Ticket" button
- Filters (search, status, assignee, tags)
- Table columns: Key, Title, Status, Assignee, Tags, Updated date

**Actions:**

- Click row → ticket detail
- Sort by columns
- Filter tickets
- Create ticket

---

### Kanban Board (`/dashboard/projects/:id/kanban`)

**Layout:** Header + columns

**Components:**

- Status columns (customizable per project)
- Ticket cards showing: Key, Title, Assignee, Tags

**Actions:**

- Drag cards between columns (updates status)
- Click card → ticket detail
- Filter tickets

---

### Ticket Detail (`/dashboard/projects/:id/tickets/:ticketId`)

**Layout:** Two-column (content + metadata)

**Left Column:**

- Title (editable)
- Description (editable)
- Comments section

**Right Sidebar:**

- Status dropdown
- Assignee dropdown
- Tags selector
- Created/updated dates
- Delete button

**Actions:**

- Edit fields (auto-save)
- Add comments
- Delete ticket

---

### Project Docs List (`/dashboard/projects/:id/docs`)

**Layout:** Header + doc list

**Components:**

- Project name
- Tab switcher (Tickets | Docs)
- "New Doc" button
- Doc list showing:
  - Title
  - Last updated date

**Actions:**

- Click doc → open doc viewer
- Create new doc
- Search docs

---

### Doc Viewer/Editor (`/dashboard/projects/:id/docs/:docId`)

**Layout:** Two-column (editor + metadata)

**Left Column:**

- Title (editable)
- Markdown editor with preview

**Right Sidebar:**

- Linked tickets
- "Link to Ticket" button
- Created/updated dates
- Delete button

**Actions:**

- Edit content (auto-save)
- Link to existing ticket
- Toggle preview
- Delete doc

---

### Settings (`/dashboard/settings`)

**Tabs:**

- Profile (name, email)
- Preferences (default view)

**Actions:**

- Update profile
- Save preferences

---

## Section 5: Database Schema

### Tables Overview

**users**
```
id                INTEGER  [PK, AUTO_INCREMENT]
email             TEXT     [UNIQUE, NOT NULL]
password_hash     TEXT     [NOT NULL]
name              TEXT     [NOT NULL]
created_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]
updated_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]
```

**user_preferences**
```
user_id            INTEGER  [PK, FK → users.id]
default_project_id INTEGER  [FK → projects.id, NULL]
default_view       INTEGER  [DEFAULT 1]  // 1=table, 2=kanban
theme              INTEGER  [DEFAULT 3]  // 1=light, 2=dark, 3=system
created_at         DATETIME [DEFAULT CURRENT_TIMESTAMP]
updated_at         DATETIME [DEFAULT CURRENT_TIMESTAMP]
```

**projects**
```
id                INTEGER  [PK, AUTO_INCREMENT]
name              TEXT     [NOT NULL]
description       TEXT     [NULL]
key_prefix        TEXT     [UNIQUE, NOT NULL]
created_by        INTEGER  [FK → users.id]
created_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]
updated_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]
deleted_at        DATETIME [NULL]

Indexes:
  idx_projects_key_prefix (key_prefix) UNIQUE
  idx_projects_creator (created_by)
```

**project_members**
```
project_id        INTEGER  [PK, FK → projects.id]
user_id           INTEGER  [PK, FK → users.id]
role              INTEGER  [NOT NULL, DEFAULT 3]  // 1=owner, 2=admin, 3=member, 4=viewer
created_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]

Indexes:
  idx_project_members_user (user_id)
```

**statuses**
```
id                INTEGER  [PK, AUTO_INCREMENT]
project_id        INTEGER  [FK → projects.id, NULL]  // NULL = global
label             TEXT     [NOT NULL]
color             TEXT     [DEFAULT '#808080']
position          INTEGER  [DEFAULT 0]
created_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]
updated_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]
deleted_at        DATETIME [NULL]

Indexes:
  idx_statuses_project_label (project_id, label) UNIQUE
  idx_statuses_project (project_id)
  idx_statuses_position (position)
```

**priorities**
```
id                INTEGER  [PK, AUTO_INCREMENT]
project_id        INTEGER  [FK → projects.id, NULL]  // NULL = global
label             TEXT     [NOT NULL]
color             TEXT     [DEFAULT '#808080']
position          INTEGER  [DEFAULT 0]
created_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]
updated_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]
deleted_at        DATETIME [NULL]

Indexes:
  idx_priorities_project_label (project_id, label) UNIQUE
  idx_priorities_project (project_id)
  idx_priorities_position (position)
```

**tags**
```
id                INTEGER  [PK, AUTO_INCREMENT]
project_id        INTEGER  [FK → projects.id, NULL]  // NULL = global
label             TEXT     [NOT NULL]
color             TEXT     [DEFAULT '#808080']
created_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]
updated_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]
deleted_at        DATETIME [NULL]

Indexes:
  idx_tags_project_label (project_id, label) UNIQUE
  idx_tags_project (project_id)
```

**tickets**
```
id                INTEGER  [PK, AUTO_INCREMENT]
project_id        INTEGER  [FK → projects.id, NOT NULL]
ticket_number     INTEGER  [NOT NULL]
title             TEXT     [NOT NULL]
description       TEXT     [NULL]
status_id         INTEGER  [FK → statuses.id, NOT NULL]
priority_id       INTEGER  [FK → priorities.id, NOT NULL]
position          INTEGER  [DEFAULT 0]
assignee_id       INTEGER  [FK → users.id, NULL]
created_by        INTEGER  [FK → users.id, NOT NULL]
created_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]
updated_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]
deleted_at        DATETIME [NULL]

Indexes:
  idx_tickets_project_number (project_id, ticket_number) UNIQUE
  idx_tickets_project (project_id)
  idx_tickets_status (status_id)
  idx_tickets_priority (priority_id)
  idx_tickets_assignee (assignee_id)
  idx_tickets_creator (created_by)
  idx_tickets_project_status (project_id, status_id)
  idx_tickets_created (created_at)
```

**ticket_tags**
```
ticket_id         INTEGER  [PK, FK → tickets.id]
tag_id            INTEGER  [PK, FK → tags.id]

Indexes:
  idx_ticket_tags_tag (tag_id)
```

**comments**
```
id                INTEGER  [PK, AUTO_INCREMENT]
ticket_id         INTEGER  [FK → tickets.id, NOT NULL]
author_id         INTEGER  [FK → users.id, NOT NULL]
body              TEXT     [NOT NULL]
created_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]
updated_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]
deleted_at        DATETIME [NULL]

Indexes:
  idx_comments_ticket (ticket_id)
  idx_comments_author (author_id)
  idx_comments_created (created_at)
```

**ticket_activity**
```
id                INTEGER  [PK, AUTO_INCREMENT]
ticket_id         INTEGER  [FK → tickets.id, NOT NULL]
actor_id          INTEGER  [FK → users.id, NOT NULL]
activity_type     INTEGER  [NOT NULL]  // 1=created, 2=status_changed, 3=assignee_changed, 
                                       // 4=commented, 5=title_changed, 6=description_changed, 
                                       // 7=tag_added, 8=tag_removed
from_status_id    INTEGER  [FK → statuses.id, NULL]
to_status_id      INTEGER  [FK → statuses.id, NULL]
from_assignee_id  INTEGER  [FK → users.id, NULL]
to_assignee_id    INTEGER  [FK → users.id, NULL]
comment_id        INTEGER  [FK → comments.id, NULL]
field_name        TEXT     [NULL]
old_value         TEXT     [NULL]
new_value         TEXT     [NULL]
created_at        DATETIME [DEFAULT CURRENT_TIMESTAMP]

Indexes:
  idx_activity_ticket (ticket_id)
  idx_activity_actor (actor_id)
  idx_activity_type (activity_type)
  idx_activity_created (created_at)
```

### Key Relationships

- Users can belong to multiple projects via `project_members`
- Projects have unique key prefixes (e.g., "TICK", "BUG")
- Tickets are numbered sequentially per project (not globally)
- Global items (statuses, priorities, tags) are templates only, projects create independent copies
- Tickets can have multiple tags via `ticket_tags` junction table
- All activity on tickets is logged in `ticket_activity`
- Soft deletes via `deleted_at` for tickets, comments, projects, statuses, priorities, and tags

---

## Section 7: API Endpoints

### Authentication

```
POST   /api/v1/auth/signup       → Create account
POST   /api/v1/auth/login        → Get JWT token
GET    /api/v1/auth/me           → Get current user
POST   /api/v1/auth/logout       → Invalidate token
```

### Projects

```
GET    /api/v1/projects          → List user's projects
POST   /api/v1/projects          → Create project
GET    /api/v1/projects/:id      → Get project details
PATCH  /api/v1/projects/:id      → Update project
DELETE /api/v1/projects/:id      → Soft delete project
```

### Project Members

```
GET    /api/v1/projects/:id/members    → List project members
POST   /api/v1/projects/:id/members    → Add member
PATCH  /api/v1/project-members/:id     → Update member role
DELETE /api/v1/project-members/:id     → Remove member
```

### Tickets

```
GET    /api/v1/projects/:id/tickets       → List tickets (with filters)
POST   /api/v1/projects/:id/tickets       → Create ticket
GET    /api/v1/tickets/:id                → Get ticket details
PATCH  /api/v1/tickets/:id                → Update ticket
DELETE /api/v1/tickets/:id                → Soft delete ticket
PATCH  /api/v1/projects/:id/tickets/bulk  → Bulk update (drag-drop)
```

### Statuses

```
GET    /api/v1/statuses                 → List global statuses
POST   /api/v1/statuses                 → Create global status
PATCH  /api/v1/statuses/:id             → Update status
DELETE /api/v1/statuses/:id             → Delete status
GET    /api/v1/projects/:id/statuses    → List project statuses
POST   /api/v1/projects/:id/statuses    → Create project status
PATCH  /api/v1/statuses/:id/position    → Update status position
```

### Priorities

```
GET    /api/v1/priorities               → List global priorities
POST   /api/v1/priorities               → Create global priority
PATCH  /api/v1/priorities/:id           → Update priority
DELETE /api/v1/priorities/:id           → Delete priority
GET    /api/v1/projects/:id/priorities  → List project priorities
POST   /api/v1/projects/:id/priorities  → Create project priority
PATCH  /api/v1/priorities/:id/position  → Update priority position
```

### Tags

```
GET    /api/v1/tags                     → List global tags
POST   /api/v1/tags                     → Create global tag
PATCH  /api/v1/tags/:id                 → Update tag
DELETE /api/v1/tags/:id                 → Delete tag
GET    /api/v1/projects/:id/tags        → List project tags
POST   /api/v1/projects/:id/tags        → Create project tag
GET    /api/v1/tickets/:id/tags         → List tags assigned to ticket
POST   /api/v1/tickets/:id/tags         → Add tag to ticket
DELETE /api/v1/tickets/:id/tags/:tag_id → Remove tag from ticket
```

### Comments

```
GET    /api/v1/tickets/:id/comments    → List ticket comments
POST   /api/v1/tickets/:id/comments    → Add comment
PATCH  /api/v1/comments/:id            → Update comment
DELETE /api/v1/comments/:id            → Soft delete comment
```

### Activity

```
GET    /api/v1/tickets/:id/activity    → Get ticket activity log
```

---

## Section 8: Out of Scope

The following features are intentionally excluded to keep the project focused and manageable:

- Real-time collaboration (WebSockets, live cursors)
- Email notifications
- File attachments
- Advanced permissions (granular role-based access)
- Third-party integrations (Slack, GitHub, Jira)
- Mobile native apps (iOS/Android)
- Offline mode with sync
- Custom workflows and automation
- Time tracking
- Gantt charts
- Roadmap views
- Advanced reporting and analytics
- Two-factor authentication
- SSO/SAML
- API rate limiting
- Webhooks

---
