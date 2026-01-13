# Ticker Frontend Skeleton - Implementation Plan

> **Goal**: Build all UI pages and components with mock data. No backend integration yet.
> **Focus**: Layout, routing, component structure, styling, user flow.

---

## Phase 1: Foundation & Basic UI Components (Day 1)

### Setup

- [x] ✅ Vite + Preact project created
- [x] ✅ Tailwind CSS configured
- [x] ✅ Folder structure created
- [x] ✅ preact-iso router already installed

### Task 1: Basic UI Components

**Create reusable components in `src/components/ui/`**

- [ ] **Button.tsx** - Primary, secondary, danger variants
  - Props: `variant`, `size`, `onClick`, `disabled`, `children`
  - Styles: Tailwind classes for each variant

- [ ] **Input.tsx** - Text, email, password inputs
  - Props: `type`, `placeholder`, `value`, `onChange`, `error`
  - Show error message if provided

- [ ] **Card.tsx** - Container with shadow/border
  - Props: `children`, `className`
  - Used for project cards, ticket cards, etc.

- [ ] **Badge.tsx** - For statuses and tags
  - Props: `label`, `color`, `variant`
  - Colored backgrounds based on type

- [ ] **Modal.tsx** - Reusable modal/dialog
  - Props: `isOpen`, `onClose`, `title`, `children`
  - Backdrop, close button, centered layout

---

## Phase 2: Layout Components (Day 1-2)

### Task 2: App Layout Structure

**Create layout components in `src/components/layout/`**

- [ ] **Header.tsx** - Top navigation bar
  - Logo (left)
  - Search bar (center) - placeholder for now
  - User menu (right): avatar, settings icon, logout button
  - Responsive: hamburger menu on mobile

- [ ] **AppShell.tsx** - Main app wrapper
  - Renders Header
  - Content area below (renders children)
  - Consistent padding/spacing

- [ ] **ProtectedRoute.tsx** - Auth guard wrapper
  - For now: just renders children (no auth check)
  - Later: will check if user is logged in

---

## Phase 3: Authentication Pages (Day 2)

### Task 3: Login/Signup Pages

**Create auth pages in `src/pages/auth/`**

- [ ] **Login.tsx** - Login page
  - Centered card layout
  - Email input
  - Password input
  - "Login" button
  - Link to switch to signup
  - Mock: Clicking login redirects to `/dashboard/projects`

- [ ] **Signup.tsx** - Signup page
  - Same layout as login
  - Additional "Name" input
  - "Sign Up" button
  - Link to switch to login
  - Mock: Clicking signup redirects to `/dashboard/projects`

- [ ] **AuthLayout.tsx** - Shared layout for auth pages
  - Centered card on clean background
  - Tab switcher (Login / Sign Up)

---

## Phase 4: Projects Page (Day 2-3)

### Task 4: Projects List

**Create project components in `src/components/project/`**

- [ ] **ProjectCard.tsx** - Single project card
  - Props: `project` (mock data)
  - Show: name, description, stats (12 open, 45 closed)
  - Click handler (navigates to ticket view)
  - Hover effect

- [ ] **CreateProjectModal.tsx** - Create project modal
  - Form: Name, Description, Key Prefix inputs
  - Save/Cancel buttons
  - Mock: Logs data to console, closes modal

**Create page in `src/pages/dashboard/`**

- [ ] **ProjectsList.tsx** - Projects list page
  - Uses AppShell layout
  - "New Project" button (opens modal)
  - Grid of ProjectCard components
  - Mock data: 3-4 sample projects

---

## Phase 5: Ticket Table View (Day 3-4)

### Task 5: Ticket Components

**Create ticket components in `src/components/ticket/`**

- [ ] **TicketFilters.tsx** - Filter bar
  - Search input
  - Status dropdown (multi-select)
  - Assignee dropdown
  - Tags dropdown
  - "Clear filters" button
  - Mock: Just UI, no actual filtering yet

- [ ] **TicketTable.tsx** - Table view
  - Props: `tickets` (mock data)
  - Columns: Key, Title, Status, Assignee, Tags, Updated
  - Sortable column headers (just visual for now)
  - Click row → navigate to ticket detail
  - Pagination controls (visual only)

- [ ] **StatusBadge.tsx** - Status indicator
  - Props: `status`
  - Color-coded badges (backlog=gray, todo=blue, in_progress=yellow, review=purple, done=green)

- [ ] **CreateTicketModal.tsx** - Create ticket modal
  - Form: Title, Description (textarea), Status, Assignee, Tags
  - Mock: Logs to console, closes modal

**Create page in `src/pages/project/`**

- [ ] **TicketTableView.tsx** - Main ticket table page
  - Uses AppShell layout
  - Header: Project name, view switcher (Table/Kanban), "New Ticket" button
  - TicketFilters component
  - TicketTable component
  - Mock data: 10-15 sample tickets

---

## Phase 6: Kanban Board (Day 4-5)

### Task 6: Kanban Components

**Create kanban components in `src/components/kanban/`**

- [ ] **KanbanColumn.tsx** - Single column
  - Props: `status`, `tickets`
  - Header with status name and count
  - List of ticket cards
  - No drag-drop yet (just layout)

- [ ] **KanbanTicketCard.tsx** - Ticket card for kanban
  - Props: `ticket`
  - Show: Key, Title, Assignee avatar, Tags, Comment count
  - Click → navigate to ticket detail
  - Compact design

- [ ] **KanbanBoard.tsx** - Full board
  - Props: `tickets`
  - 5 columns: Backlog, To Do, In Progress, Review, Done
  - Horizontal scrolling on mobile
  - Mock: No drag-drop, just display

**Create page**

- [ ] **KanbanView.tsx** - Kanban page
  - Uses AppShell layout
  - Same header as table view (with view switcher)
  - TicketFilters component
  - KanbanBoard component
  - Mock data: Same tickets as table view

---

## Phase 7: Ticket Detail Page (Day 5-6)

### Task 7: Ticket Detail

**Create components in `src/components/ticket/`**

- [ ] **TicketDetailSidebar.tsx** - Right sidebar
  - Status dropdown
  - Assignee dropdown
  - Tags multi-select
  - Created/Updated dates
  - Delete button

- [ ] **CommentList.tsx** - List of comments
  - Props: `comments` (mock data)
  - Show: Author name, avatar, timestamp, comment body
  - Edit/delete buttons (own comments only)

- [ ] **CommentForm.tsx** - Add comment
  - Textarea input
  - Submit button
  - Mock: Logs to console

**Create page**

- [ ] **TicketDetail.tsx** - Ticket detail page
  - Two-column layout
  - Left: Title (editable input), Description (editable textarea), Comments
  - Right: TicketDetailSidebar
  - Mock: All edits log to console
  - Back button to return to table/kanban

---

## Phase 8: Settings Page (Day 6)

### Task 8: Settings

**Create page in `src/pages/dashboard/`**

- [ ] **Settings.tsx** - Settings page
  - Uses AppShell layout
  - Tabs: Profile, Preferences
  - Profile tab: Edit name, email inputs
  - Preferences tab: Default project dropdown, Default view radio buttons
  - Save buttons
  - Mock: Logs changes to console

---

## Phase 9: Router & Navigation (Day 6-7)

### Task 9: Wire Up Routing

**Update `src/app.tsx`**

- [ ] **Configure all routes** with preact-iso

  ```tsx
  <Route path="/" component={Login} />
  <Route path="/signup" component={Signup} />
  <Route path="/dashboard/projects" component={ProjectsList} />
  <Route path="/dashboard/projects/:projectId/ticket" component={TicketTableView} />
  <Route path="/dashboard/projects/:projectId/kanban" component={KanbanView} />
  <Route path="/dashboard/projects/:projectId/tickets/:ticketId" component={TicketDetail} />
  <Route path="/dashboard/settings" component={Settings} />
  <Route default component={NotFound} />
  ```

- [ ] **Test all navigation flows**
  - Login → Projects list
  - Click project → Ticket table
  - Switch to Kanban view
  - Click ticket → Ticket detail
  - Navigate to Settings

---

## Phase 10: Polish & Mock Data (Day 7)

### Task 10: Mock Data & States

- [ ] **Create mock data file** - `src/utils/mockData.ts`
  - Mock users (5-6 people)
  - Mock projects (3-4 projects)
  - Mock tickets (15-20 tickets across different statuses)
  - Mock tags (10 tags with different colors)
  - Mock comments (5-10 comments)

- [ ] **Add loading states**
  - Create `Spinner.tsx` component
  - Show spinners where data would load (2 second setTimeout to simulate)

- [ ] **Add empty states**
  - Empty project list: "Create your first project"
  - Empty ticket list: "No tickets yet"
  - No comments: "Be the first to comment"

- [ ] **Responsive design check**
  - Test all pages on mobile width
  - Adjust layouts as needed
  - Hamburger menu in header

---

## Phase 11: Final Touches (Day 7)

### Task 11: UX Polish

- [ ] **Hover effects** on cards, buttons, table rows
- [ ] **Focus states** on inputs, buttons
- [ ] **Transitions** - Smooth modal open/close, page transitions
- [ ] **Error states** - Show validation errors on forms
- [ ] **Keyboard navigation** - Tab through forms properly
- [ ] **Accessibility** - Add ARIA labels where needed

---

## Deliverables

After completing this plan, you'll have:

✅ **Fully designed UI** - All pages look complete and polished
✅ **Component library** - Reusable UI components ready
✅ **Complete routing** - All navigation works
✅ **Mock data** - Realistic test data throughout
✅ **Responsive design** - Works on mobile and desktop
✅ **Ready for backend** - Just swap mock data with API calls

---

## File Structure Summary

```
src/
├── components/
│   ├── ui/
│   │   ├── Button.tsx
│   │   ├── Input.tsx
│   │   ├── Card.tsx
│   │   ├── Badge.tsx
│   │   ├── Modal.tsx
│   │   └── Spinner.tsx
│   ├── layout/
│   │   ├── Header.tsx
│   │   ├── AppShell.tsx
│   │   └── ProtectedRoute.tsx
│   ├── project/
│   │   ├── ProjectCard.tsx
│   │   └── CreateProjectModal.tsx
│   ├── ticket/
│   │   ├── TicketFilters.tsx
│   │   ├── TicketTable.tsx
│   │   ├── StatusBadge.tsx
│   │   ├── CreateTicketModal.tsx
│   │   ├── TicketDetailSidebar.tsx
│   │   ├── CommentList.tsx
│   │   └── CommentForm.tsx
│   └── kanban/
│       ├── KanbanBoard.tsx
│       ├── KanbanColumn.tsx
│       └── KanbanTicketCard.tsx
├── pages/
│   ├── auth/
│   │   ├── Login.tsx
│   │   └── Signup.tsx
│   ├── dashboard/
│   │   ├── ProjectsList.tsx
│   │   └── Settings.tsx
│   └── project/
│       ├── TicketTableView.tsx
│       ├── KanbanView.tsx
│       └── TicketDetail.tsx
├── utils/
│   └── mockData.ts
├── app.tsx (router)
└── main.tsx (entry)
```

---

## Estimated Timeline

- **Days 1-2**: UI components + Layout + Auth pages
- **Days 3-4**: Projects + Ticket table view
- **Days 4-5**: Kanban board
- **Days 5-6**: Ticket detail + Settings
- **Days 6-7**: Routing + Polish + Mock data

**Total: ~1 week for complete UI skeleton**

---

## Next Step

**Ready to start?** I'll give you the first file to create:

**Task 1.1: Button Component** - Want me to provide the code?
