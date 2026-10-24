---
id: "02-tenant-account-administration"
title: "Constrain organization account administration"
status: implemented
wave: 1
depends_on: []
plan: "plan.md"
requirements:
  - REQ-AUTH-ADMIN-CONTAINMENT-001
acceptance_criteria:
  - AC-AUTH-ADMIN-CONTAINMENT-001.1
  - AC-AUTH-ADMIN-CONTAINMENT-001.2
  - AC-AUTH-ADMIN-CONTAINMENT-001.3
system_design:
  - ../../specs/auth/system-design/administration-containment.md
---

# Task 02: Constrain organization account administration

## Delivered scope

Implemented org-scoped user/invite operations and a transactional org-local last-administrator guard. Full auth and user-store tests, including concurrent SQLite demotions and denied foreign targets, pass locally. Live PostgreSQL concurrency is not verified.

## Original design summary

Implement account and invite service authorization using the authenticated tenant, tenant-scoped repository operations, and transactional last-admin protection.

## In scope

The linked acceptance criteria, implementation files below, focused regression
coverage, and directly affected public documentation and companion records.

## Out of scope

Unrelated refactors, new product UI layouts, automatic publication, upstream
messaging, dependency upgrades without evidence, and claiming all vulnerabilities
are fixed. No delegation is authorized.

## Acceptance

- Two-org tests deny foreign list/update/invite access without revoking foreign credentials.
- Own-org management works; last-admin checks cannot be bypassed by concurrent requests.

## Verification

Run from the repository root in Git Bash or a supported Linux shell. Before any
frontend command in a fresh checkout, run `(cd apps && pnpm install --frozen-lockfile)`.
The commands below preserve the original validation design. Actual delivery and verification limits are recorded in Results; unimplemented browser/container scenario names are restoration work.

```bash
(cd apps/backend && go test ./internal/auth/... ./internal/user/store/... ./internal/org/... -count=1)
```

Record explicit platform/service prerequisites and unsuccessful checks. For Go
race coverage use a supported compiler-enabled host and rerun the same packages
with `-race` when the change affects shared state.

## Files likely touched

- `apps/backend/internal/auth/service_users.go`
- `apps/backend/internal/auth/service_invites.go`
- `apps/backend/internal/auth/httpapi/handlers_admin.go`
- `apps/backend/internal/user/store/sqlite_tenancy.go`
- `apps/backend/internal/auth/service_test.go`
- `apps/backend/internal/auth/httpapi/handlers_test.go`

## Dependencies

None.

## Risks

SQLite and PostgreSQL transaction semantics must both be tested. Existing global admin-count behavior must not be silently retained for tenant accounts.

## Parallelism

`sequential`. Waves express prerequisites only and do not authorize subagents.

## Inputs

- The referenced requirements and system design in full.
- [Plan and compatibility matrix](plan.md).
- [Evidence and scope corrections](evidence.md).
- Applicable scoped `AGENTS.md` and existing tests beside the affected code.

## Results

Implemented org-scoped user/invite operations and a transactional org-local last-administrator guard. Full auth and user-store tests, including concurrent SQLite demotions and denied foreign targets, pass locally. Live PostgreSQL concurrency is not verified.
