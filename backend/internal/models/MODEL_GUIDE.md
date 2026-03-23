# Model Guide

This guide reflects the current refactored model style used by:
- `projects`
- `tags`
- `statuses`
- `priorities`

## 1) Model Layer Role

`internal/models` represents domain data contracts, not persistence logic.

Model files should primarily contain:
- domain structs
- domain-level constants shared across layers (for example `MaxProjectNameLength`)

Model files should avoid:
- SQL queries
- scan helpers
- store/business orchestration methods

## 2) File Structure

Preferred order:
1. imports
2. constants
3. struct definitions

Keep files compact and data-focused.

## 3) Struct Semantics

- Use singular type names (`Project`, `Tag`, `Status`, `Priority`).
- Keep JSON tags explicit and stable.
- Use pointer fields only for nullable data (`DeletedAt *time.Time`).
- Keep field names consistent with API and scan target expectations in store.

## 4) Constants and Shared Domain Rules

Constants that define stable domain limits belong in models when shared across layers.

Examples:
- `MaxProjectNameLength`

Validation execution still belongs in service layer.

## 5) Error Ownership

Current ownership split:
- store/persistence errors -> `internal/store`
- service validation/business errors -> `internal/services`
- models generally avoid operation-specific errors in refactored domains

## 6) Compatibility Rules

- Do not rename model fields or JSON tags without API review.
- Keep model changes backward-compatible with store scan order and handler responses.
- If a non-refactored domain is migrated, update its model file to this thin pattern in the same change.
