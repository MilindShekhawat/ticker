# Handler Guide

This guide defines the handler conventions used by the refactored domains:
- `projects`
- `tags`
- `statuses`
- `priorities`
- `comments`
- `ticket_activity`
- `tickets`

It is intentionally scoped. Other domains may still follow older patterns until migrated.

## 1) Purpose and Boundaries

Handlers should:
- parse HTTP params/body
- read auth context when endpoint semantics require user identity
- call service methods
- map domain/store/service errors to HTTP responses

Handlers should not:
- execute SQL
- contain cross-entity business rules

## 2) Construction Pattern

Use struct-based handlers with injected service interfaces.

```go
type TagHandler struct {
    service services.TagService
}

func NewTagHandler(service services.TagService) *TagHandler {
    return &TagHandler{service: service}
}
```

## 3) File Layout

Keep this order:
1. imports
2. handler struct + constructor
3. request DTOs
4. primary handlers (`List`, `Get` if present, `Create`, `Update`, `Delete`)
5. auxiliary handlers (for example `UpdatePosition`)
6. private error-mapping helper at bottom

## 4) Parsing and Validation Semantics

- Parse integer params with `ParamsInt`.
- Parse JSON body with `BodyParser`.
- Return `400` for malformed params/body.
- For patch endpoints, reject empty update payloads in handler.

Example:

```go
if req.Label == nil && req.Color == nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No fields to update"})
}
```

## 5) Auth Semantics

Read `c.Locals("user")` only where user identity is required by service contract:
- required in `projects` handlers (`List/Get/Create/Update/Delete` for user scope)
- required in `comments` handlers (`List/Create/Update/Delete` for membership checks)
- required in `ticket_activity` handlers (`GetTicketActivity` for membership checks)
- required in `tickets` handlers (`List/Get/Create/Update/Delete` for membership checks)
- not required in current `tags/statuses/priorities` handlers (project-scoped route inputs drive operations)

## 6) Error Mapping Semantics

Use one domain-local error mapper (`projectErrorResponse`, `tagErrorResponse`, etc.) and `errors.Is`.

Recommended mapping:
- validation/input errors -> `400`
- not found (`project/tag/status/priority`) -> `404`
- unauthorized auth context -> `401`
- unexpected errors -> `500`

Prefer stable, explicit error messages in JSON responses when returning bodies.

## 7) Response Semantics

- `List/Get/Update`: `200` + JSON payload
- `Create`: `201` + JSON payload
- `Delete`: `204` with empty body

Return collection responses as `[]` (not `null`). This is enforced by service/store return shape.

## 8) Naming Semantics

- handler types: `ProjectHandler`, `TagHandler`, `StatusHandler`, `PriorityHandler`, `TicketHandler`
- constructors: `NewXHandler`
- DTOs: `CreateXRequest`, `UpdateXRequest`
- path params: `projectID`, `ticketID`, `tagID`, `statusID`, `priorityID`
- mapper helper: `xErrorResponse`

## 9) Route Contract Alignment

Handler param keys must match route definitions exactly.

Current aligned patterns:
- `/projects/:projectID/tickets`
- `/tickets/:ticketID`
- `/projects/:projectID/tags`
- `/projects/:projectID/statuses`
- `/projects/:projectID/priorities`
- `/tags/:tagID`
- `/statuses/:statusID`
- `/priorities/:priorityID`

If route tokens change, update handler `ParamsInt(...)` keys in the same change.
