# Ticker - Product Specification

> **Goal**: Build a lightweight project and ticket management system inspired by GitHub Projects.  
> **Focus**: Clean data modeling, solid CRUD operations, simple permissions.

---

## 1. What We're Building

### Core Features

- **Projects**: Create and manage multiple projects
- **Tickets**: Track work items with statuses and assignments
- **Views**: Switch between table view and kanban board
- **Tags**: Organize tickets with custom labels
- **Comments**: Discuss tickets with team members

### User Experience

- Clean, modern interface (GitHub Projects style)
- Fast and responsive
- Works on desktop and mobile
- No login required for demo/testing mode

---

## 2. User Flow

### First Time User

1. Lands on login page
2. Signs up with email/password
3. Sees empty dashboard with "Create Project" prompt
4. Creates first project (e.g., "Personal Tasks")
5. Adds tickets to project
6. Views tickets in table or kanban mode

### Returning User

1. Logs in
2. Sees dashboard with project list
3. Clicks project → goes to ticket view (last used view: table/kanban)
4. Creates/edits tickets
5. Comments on tickets
6. Switches between projects

---

## 3. Pages & Routes

```
/                              → Login/Signup page
/dashboard                     → Redirect to /dashboard/projects
/dashboard/projects            → Project list (grid of cards)
/dashboard/projects/:id/ticket → Table view of tickets
/dashboard/projects/:id/kanban → Kanban board view
/dashboard/projects/:id/tickets/:ticketId → Ticket detail page
/dashboard/settings            → User settings
```

**Note**: Routes explicitly include `projectId` - no ambiguity, URLs are shareable.

---

## 4. Page Details

### Login/Signup Page (`/`)

- **Layout**: Centered card on clean background
- **Components**:
  - Toggle tabs (Login / Sign Up)
  - Email input
  - Password input
  - Name input (signup only)
  - Submit button
  - Error messages
- **Actions**:
  - Validate inputs
  - Call API
  - Store token
  - Redirect to `/dashboard/projects`

---

### Projects List (`/dashboard/projects`)

- **Layout**: Header + grid of project cards
- **Header**:
  - Logo
  - Search bar
  - User menu (avatar, settings, logout)
- **Content**:
  - "New Project" button (opens modal)
  - Project cards showing:
    - Project name
    - Description
    - Stats (12 open, 45 closed)
    - Members count
- **Actions**:
  - Click card → go to `/dashboard/projects/:id/ticket`
  - Search projects
  - Create new project

---

### Ticket Table View (`/dashboard/projects/:id/ticket`)

- **Layout**: Header + filters + table + pagination
- **Header**:
  - Project name
  - View switcher (Table ↔ Kanban)
  - "New Ticket" button
- **Filters**:
  - Search box
  - Status dropdown (multi-select)
  - Assignee dropdown
  - Tags dropdown (multi-select)
  - Clear filters button
- **Table Columns**:
  - Key (TCK-101)
  - Title
  - Status badge
  - Assignee avatar
  - Tags (colored badges)
  - Updated date
- **Actions**:
  - Click row → go to ticket detail
  - Sort by columns
  - Filter
  - Paginate

---

### Kanban Board (`/dashboard/projects/:id/kanban`)

- **Layout**: Header + filters + columns
- **Columns**: Backlog | To Do | In Progress | Review | Done
- **Cards** (per ticket):
  - Key + Title
  - Assignee avatar
  - Tags
  - Comment count
- **Actions**:
  - Drag card between columns (updates status)
  - Click card → go to ticket detail
  - Filter (same as table view)

---

### Ticket Detail (`/dashboard/projects/:id/tickets/:ticketId`)

- **Layout**: Two-column layout
  - **Left**: Ticket content
    - Title (editable)
    - Description (editable)
    - Comments section
  - **Right**: Metadata sidebar
    - Status dropdown
    - Assignee dropdown
    - Tags multi-select
    - Dates (created, updated)
    - Delete button
- **Actions**:
  - Edit any field (auto-save)
  - Add comment
  - Delete ticket

---

### Settings (`/dashboard/settings`)

- **Tabs**:
  - Profile (edit name, email)
  - Preferences (default project, default view)
- **Actions**:
  - Update profile
  - Save preferences

---

## 5. Components Breakdown

### Reusable UI Components

```
Button          → Primary, secondary, danger variants
Input           → Text, email, password, textarea
Select          → Single and multi-select dropdowns
Modal           → For create/edit forms
Card            → Container with shadow/border
Badge           → For statuses and tags
Avatar          → User profile picture or initials
Spinner         → Loading indicator
```

### Feature Components

```
ProjectCard           → Shows project info on grid
TicketCard            → Ticket card for kanban board
TicketTable           → Table view of tickets
TicketFilters         → Filter bar component
KanbanBoard           → Drag-drop board
KanbanColumn          → Single status column
CommentList           → List of comments
CommentForm           → Add comment form
TagSelector           → Multi-select tag picker
StatusBadge           → Colored status indicator
```

### Layout Components

```
AppShell              → Main layout with header
Header                → Top navigation bar
ProtectedRoute        → Auth guard for routes
```

---

## 6. Data Structure (Simplified)

### User

```
id, email, name, createdAt
```

### Project

```
id, name, description, keyPrefix (e.g., "TCK")
stats: { open, closed, members }
```

### Ticket

```
id, key (TCK-101), title, description
status, assigneeId, tagIds[]
createdAt, updatedAt
```

### Tag

```
id, projectId, label, color
```

### Comment

```
id, ticketId, authorId, body, createdAt
```

---

## 7. API Endpoints (What Frontend Needs)

### Auth

```
POST   /auth/signup       → Create account
POST   /auth/login        → Get token
GET    /auth/me           → Get current user
POST   /auth/logout       → Clear session
```

### Projects

```
GET    /projects                → List user's projects
POST   /projects                → Create project
GET    /projects/:id            → Get project details
PATCH  /projects/:id            → Update project
GET    /projects/:id/members    → Get project members (for assignee dropdown)
```

### Tickets

```
GET    /projects/:id/tickets          → List tickets (with filters)
POST   /projects/:id/tickets          → Create ticket
GET    /tickets/:id                   → Get ticket detail
PATCH  /tickets/:id                   → Update ticket
DELETE /tickets/:id                   → Delete ticket
PATCH  /projects/:id/tickets/bulk     → Bulk update (for drag-drop)
```

### Tags

```
GET    /projects/:id/tags    → List project tags
POST   /projects/:id/tags    → Create tag
PATCH  /tags/:id             → Update tag
DELETE /tags/:id             → Delete tag
```

### Comments

```
GET    /tickets/:id/comments  → List comments (or use ?include=comments on ticket detail)
POST   /tickets/:id/comments  → Add comment
```

---

## 8. State Management Strategy

### Global State (React Context)

- **AuthContext**: Current user, login/logout functions, loading state
- **ProjectContext**: Current project, switch project function

### Local State (Component State)

- Form inputs
- Modal open/close
- Loading indicators
- Error messages

### Server State (React Query or SWR - optional)

- Tickets list
- Project list
- Comments
- Auto-refresh, caching, optimistic updates

**For v1**: Just use React Context + useState. Keep it simple.

---

## 9. Feature Priority

### Phase 1: Authentication (Week 1)

- Login/Signup pages
- Token storage
- Protected routes
- Auth context

### Phase 2: Projects (Week 1-2)

- Projects list page
- Create project modal
- Project cards

### Phase 3: Tickets - Table View (Week 2-3)

- Table view page
- Filters
- Create ticket modal
- Ticket detail page

### Phase 4: Tickets - Kanban (Week 3-4)

- Kanban board
- Drag and drop
- Status updates

### Phase 5: Tags & Polish (Week 4-5)

- Tag selector
- Tag management
- Comments
- Loading states
- Error handling

---

## 10. Tech Stack

### Frontend

- **Framework**: Vite + Preact + TypeScript
- **Styling**: Tailwind CSS
- **Routing**: wouter (lightweight router)
- **Forms**: Controlled components with validation
- **Drag-Drop**: @dnd-kit/core (for kanban)
- **HTTP**: fetch API (wrapped in client)

### Backend (for reference)

- **Framework**: Go + Gorilla Mux
- **Database**: SQLite
- **Auth**: JWT tokens

---

## 11. Design Principles

### Keep It Simple

- No over-engineering
- Standard patterns
- Readable code over clever code

### User First

- Fast page loads
- Clear feedback (loading, errors)
- No surprises

### Maintainable

- Component reuse
- Clear file structure
- Type safety with TypeScript

---

## 12. Out of Scope (Not Building)

- Real-time collaboration (WebSockets)
- Email notifications
- File attachments
- Advanced permissions (just admin/member)
- Integrations (Slack, GitHub, etc.)
- Mobile apps
- Offline mode
- Custom workflows
- Time tracking
- Gantt charts

---

## Next Steps

Start with **Phase 1: Authentication**

1. Set up Vite project
2. Install dependencies
3. Create folder structure
4. Build Login/Signup page
5. Set up routing
6. Implement auth context

Then move to Phase 2, 3, 4, 5 sequentially.
