# Service Guide

This guide captures the service conventions used by:
- `services/projects.go`
- `services/tags.go`
- `services/statuses.go`
- `services/priorities.go`
- `services/comments.go`
- `services/ticket_activity.go`
- `services/tickets.go`

## 1) Service Role

Services are the business layer between handlers and stores.

Services should own:
- normalization
- validation and defaults
- cross-store checks (for example project existence)
- orchestration across stores

Services should avoid:
- raw SQL
- HTTP-specific behavior (status codes, request parsing)

## 2) Construction Pattern

Use interface + private implementation:

```go
type ProjectService interface { ... }

type projectService struct {
    store store.ProjectStore
}

func NewProjectService(store store.ProjectStore) ProjectService {
    return &projectService{store: store}
}
```

## 3) Validation Semantics

### Projects
- trim `name`
- trim + uppercase `keyPrefix`
- validate name length using `models.MaxProjectNameLength`
- validate key prefix format
- reject zero `userID`

### Tickets
- trim `title`
- require non-empty `title`
- enforce title max length
- verify caller is a member of the project for list/create and of the ticket's project for get/update/delete
- verify status/priority/tag references belong to the ticket project
- validate assignee user existence when provided

### Tags / Statuses / Priorities
- require non-empty `label`
- validate `color` as `#RRGGBB` when provided
- default empty color to `#808080` on create
- verify target project exists before create

### Comments
- trim `body`
- require non-empty `body`
- verify caller is a member of the ticket's project before list/create/update/delete

### Ticket Activity
- verify caller is a member of the ticket's project before list

## 4) Error Semantics

- return service validation errors from `services/errors.go`
- propagate store errors unchanged for handlers to map
- use sentinel errors, not free-form error strings for business outcomes

## 5) Method Naming

Keep method names action-oriented and scope-aware:
- shared CRUD-style domains: `List`, `Create`, `Update`, `Delete`
- user-scoped domains: `ListForUser`, `GetForUser`, `CreateForUser`, `UpdateForUser`, `DeleteForUser`

## 6) Dependency Semantics

Inject only required stores per service:
- project service: `ProjectStore`
- tag/status/priority services: domain store + `ProjectStore`
- ticket service: `TicketStore` + `ProjectStore` + `ProjectMemberStore` + `StatusStore` + `PriorityStore` + `UserStore` + `TagStore`

Avoid hidden globals in service implementations.

## 7) Update Semantics

For update operations:
- service validates semantic constraints (for example color format)
- store handles partial update SQL execution
- handlers only enforce request-shape constraints (for example no fields provided)
