# AGENTS.md

## Source Priority

1. `SPECIFICATION.md` (baseline product/behavior spec; improvable)
2. `AGENTS.md`
3. Existing patterns in `backend/` and `frontend/` (only when compatible with higher-priority sources)

If rules conflict, follow the higher-priority source and call out the conflict.

## Repository Context

- Backend: Go, Fiber v2, SQLite, migrations (`backend/`)
- Frontend: Preact + TypeScript + Vite + Tailwind (`frontend/`)

## Working Rules

- Check impacted behavior against `SPECIFICATION.md` before coding.
- Keep changes focused; avoid broad refactors unless requested.
- Do not invent unspecified behavior; choose the safest default and note it.
- Preserve backward compatibility unless the spec requires change.
- Update docs with non-trivial behavior changes.

## Spec Evolution

- You may challenge and improve `SPECIFICATION.md` when it improves correctness or clarity.
- Keep spec updates focused and traceable in the same change when applicable.
- If an improvement is valid but out of scope, record it under an `Out of scope` section in `SPECIFICATION.md`.

## Backend Rules (Go)

- Use idiomatic Go: small functions, explicit error handling, early returns.
- Keep handlers thin; move business logic to internal packages.
- Enforce auth and project membership on project-scoped operations.
- Use parameterized queries only.
- Keep data-access style consistent within each file/module.
- Keep response shapes and HTTP status usage consistent.
- Prefer composition/dependency injection over globals.
- Keep files domain-focused (`handlers`, `models`, `middleware`, `db`, etc.).
- Run:
  - `go test ./...`
  - `go vet ./...`

## Frontend Rules (Preact)

- Use typed functional components and typed props.
- Keep state local by default; lift only when shared.
- Prefer derived state over duplicated state.
- Reuse existing UI primitives in `frontend/src/components/ui/`.
- Handle loading, error, and empty states explicitly.
- Keep routes and page flows aligned with the spec.
- Maintain accessibility basics: semantics, labels, keyboard access, visible focus.
- Run:
  - `npm run build` (from `frontend/`)

## Consistency Rules

Applies to code style, user-facing copy, comments, docs, API messages, and task summaries.

- Match local conventions before introducing new patterns.
- Keep similar operations structurally similar.
- Use predictable naming (`rows`, `query`, `args`, etc.) and avoid ambiguous abbreviations.
- Keep wording aligned across backend responses and frontend UI.
- Tone: direct, neutral, concise.
- Voice: active, present tense.
- Canonical terms: `Project`, `Ticket`, `Tag`, `Status`, `Priority`, `Comment`, `Member`, `Owner/Admin/Member/Viewer`.
- Casing: sentence case for UI labels; technical identifiers follow code conventions.
- Errors should state failure, likely cause (if known), and next action; avoid vague text.

## Definition of Done

- Behavior matches the spec for the impacted area.
- Relevant checks pass for changed areas.
- No unrelated file churn.
- Summaries explain what changed, why, and any assumptions/gaps.
