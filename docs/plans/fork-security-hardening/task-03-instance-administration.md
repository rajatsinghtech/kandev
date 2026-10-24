---
id: "03-instance-administration"
title: "Enforce installation operator authority"
status: implemented
wave: 2
depends_on: ["02-tenant-account-administration"]
plan: "plan.md"
requirements:
  - REQ-AUTH-ADMIN-CONTAINMENT-002
acceptance_criteria:
  - AC-AUTH-ADMIN-CONTAINMENT-002.1
  - AC-AUTH-ADMIN-CONTAINMENT-002.2
  - AC-AUTH-ADMIN-CONTAINMENT-002.3
system_design:
  - ../../specs/auth/system-design/administration-containment.md
---

# Task 03: Enforce installation operator authority

## Delivered scope

Installed ControlPlaneBoundary after authentication to deny installation-wide reads/mutations without the required installation authority. HTTP matrix tests check response and absence of handler side effects. Organizations mode remains unavailable; existing workspace authorizers are retained.

## Original design summary

Centralize installation authority and inventory backup, storage, database, runtime flags, update/restart, and global configuration routes and service entry points. Preserve workspace privacy.

## In scope

The linked acceptance criteria, implementation files below, focused regression
coverage, and directly affected public documentation and companion records.

## Out of scope

Unrelated refactors, new product UI layouts, automatic publication, upstream
messaging, dependency upgrades without evidence, and claiming all vulnerabilities
are fixed. No delegation is authorized.

## Acceptance

- Ordinary org admins cannot download full backups or invoke installation mutations.
- Real operators retain installation management; single-org authenticated admins retain existing access.
- Denied calls start no background jobs and emit no successful mutation events.

## Verification

Run from the repository root in Git Bash or a supported Linux shell. Before any
frontend command in a fresh checkout, run `(cd apps && pnpm install --frozen-lockfile)`.
The commands below preserve the original validation design. Actual delivery and verification limits are recorded in Results; unimplemented browser/container scenario names are restoration work.

```bash
(cd apps/backend && go test ./internal/auth/... ./internal/authz/... ./internal/org/... ./internal/system/... ./internal/runtimeflags/... -count=1)
```

Record explicit platform/service prerequisites and unsuccessful checks. For Go
race coverage use a supported compiler-enabled host and rerun the same packages
with `-race` when the change affects shared state.

## Files likely touched

- `apps/backend/internal/auth/authn/identity.go`
- `apps/backend/internal/authz/gin.go`
- `apps/backend/internal/system/system.go`
- `apps/backend/internal/system/system_authz_test.go`
- `apps/backend/internal/system/backups/handler.go`
- `apps/backend/internal/runtimeflags`
- `apps/backend/internal/org/controller.go`

## Dependencies

02-tenant-account-administration

## Risks

Inventory all HTTP/WS/MCP alternatives. An operator with full backup permission is privileged recovery personnel, not an adversary excluded from offline tenant data.

## Parallelism

`sequential`. Waves express prerequisites only and do not authorize subagents.

## Inputs

- The referenced requirements and system design in full.
- [Plan and compatibility matrix](plan.md).
- [Evidence and scope corrections](evidence.md).
- Applicable scoped `AGENTS.md` and existing tests beside the affected code.

## Results

Installed ControlPlaneBoundary after authentication to deny installation-wide reads/mutations without the required installation authority. HTTP matrix tests check response and absence of handler side effects. Organizations mode remains unavailable; existing workspace authorizers are retained.
