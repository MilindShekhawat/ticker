# Store Guide

This guide describes the persistence conventions used by:
- `store/projects.go`
- `store/tags.go`
- `store/statuses.go`
- `store/priorities.go`
- `store/comments.go`
- `store/ticket_activity.go`
- `store/tickets.go`

## 1) Store Role

Stores are responsible for persistence concerns only:
- SQL queries and commands
- transaction handling
- scan/mapping from DB rows into `models` structs
- mapping DB outcomes (`no rows`, constraint failures) to store-level errors

Stores should not implement HTTP behavior or business validation policy.

## 2) Construction Pattern

Use interface + private SQL-backed implementation:

```go
type ProjectStore interface { ... }

type projectStore struct {
    db *sql.DB
}

func NewProjectStore() ProjectStore {
    return &projectStore{db: db.DB}
}
```

## 3) Scan Helper Pattern

Use one private scan helper per domain file:
- `scanProject`
- `scanTag`
- `scanStatus`
- `scanPriority`
- `scanComment`
- `scanTicketActivity`
- `scanTicket`

Selected SQL columns must match scan field order exactly.

## 4) Query Semantics

- Always use parameterized SQL (`?` placeholders).
- Use raw multiline SQL strings for readability.
- For list methods:
  - initialize slice as non-nil (`make([]models.X, 0)` or `[]models.X{}`)
  - close rows with `defer rows.Close()`
  - return `rows.Err()` after iteration

## 5) Update/Delete Semantics

- For partial updates, build `updates` + `args` dynamically.
- For update/delete operations, always check `RowsAffected()`.
- Zero affected rows maps to domain store not-found errors.

## 6) Transaction Semantics

Use transactions for multi-step writes:
- begin transaction
- rollback on failure path
- commit only after all dependent writes succeed

Example: project create + project owner membership insert.

## 7) Error Semantics

- Store-level sentinel errors live in `store/errors.go` and domain store files.
- Convert `sql.ErrNoRows` to store not-found errors.
- Keep not-found and constraint errors stable for handler mapping.

For unexpected DB errors:
- include context when helpful
- preserve original error via wrapping when wrapping is used

## 8) Naming Semantics

- interfaces: `ProjectStore`, `TagStore`, `StatusStore`, `PriorityStore`, `CommentStore`, `TicketActivityStore`, `TicketStore`
- constructors: `NewProjectStore`, `NewTagStore`, `NewCommentStore`, `NewTicketStore`, ...
- core methods: `List`, `GetByID`, `Create`, `Update`, `Delete`
- scope-specific methods use explicit suffixes (`GetByUser`, `UpdateByUser`)
