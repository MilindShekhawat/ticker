Consistency check: unchanged.

Below is a concise status list based on the earlier structural issues.

---

# ✅ Solved

### Authentication & Identity

- Removed client-controlled `created_by` and `author_id`
- Implemented bcrypt password hashing
- Added DB-backed sessions
- Added `sessions` table + migration
- Implemented HttpOnly Secure cookies
- Added `AuthRequired` middleware
- Added `/auth/me`
- Added `/auth/logout`
- Immediate session revocation supported
- Session expiry enforced in DB

### Security Improvements

- Eliminated identity spoofing
- Removed trust in client-supplied user IDs
- Centralized identity via middleware
- Enforced protected route grouping

### Specification Alignment

- Updated spec to reflect session-based auth
- Documented removal of JWT
- Updated DB schema section

---

# 🟡 Partially Solved / Needs Refactor

### Handlers

- Some handlers may still accept `CreatedBy`, `AuthorID` fields (needs full sweep)
- Activity logging not yet enforced consistently
- Tag operations still “best effort” (not transactional)

### Transactions

- Ticket creation + tag insertion not wrapped in transaction
- Ticket update + tag replacement not transactional
- Comment creation + activity log not atomic

---

# ❌ Still Unresolved Structural Issues

### Authorization (Not Authentication)

- No project membership enforcement
- No role-based access checks (owner/admin/member/viewer)
- No project-scoped validation on tickets, statuses, priorities, tags
- Users can still access resources across projects if IDs are known

### Data Integrity

- `validateStatus` and `validatePriority` do not verify project ownership
- Tags can be added across projects without restriction
- Status/priority updates do not ensure same project as ticket

### Concurrency

- Ticket number generation not transaction-safe
- Position updates not conflict-safe
- No optimistic locking

### Validation Gaps

- No max length checks at DB level
- No normalization of email (case-insensitive uniqueness)
- No rate limiting (intentionally out of scope, but noted)

### Session Management

- No background cleanup of expired sessions
- No “logout all sessions” support
- No session rotation

---

# 🔴 Architectural Gaps Remaining

- No authorization layer
- No RBAC middleware
- No project isolation enforcement
- No consistent activity logging integration
- No audit consistency guarantees (transactions)

---

# Summary

Authentication layer: **Correct and production-safe for monolith.**

Authorization layer: **Not implemented.**

Data integrity under concurrency: **Not hardened.**

Multi-step atomicity: **Not enforced.**

---

If you want next step prioritized, the most critical remaining issue is:

**Project-scoped authorization enforcement (membership + role middleware).**
