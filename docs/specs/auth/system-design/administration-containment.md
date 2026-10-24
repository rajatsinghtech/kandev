---
status: draft
system: auth
requirements:
  - REQ-AUTH-ADMIN-CONTAINMENT-001
  - REQ-AUTH-ADMIN-CONTAINMENT-002
---

# Administration containment system design

## Purpose and boundaries

Auth owns caller identity and administration authorization. These requirements specialize existing role/scope and multi-tenancy contracts at account and installation administration boundaries; workspace reach remains unchanged.

## Requirement mapping

| Requirement | Design sections |
| --- | --- |
| `REQ-AUTH-ADMIN-CONTAINMENT-001` | Components, control flow, and verification below |
| `REQ-AUTH-ADMIN-CONTAINMENT-002` | Components, control flow, and verification below |

## Restricted delivery

TenantAccountRepository scopes users and serializes last-admin changes in one transaction. Invite operations use org predicates. ControlPlaneBoundary protects installation routes after authentication. New service/store and HTTP matrix regressions cover the restricted behavior; Organizations mode itself remains refused pending runtime isolation.

## Components and responsibilities

`internal/auth/service_users.go` owns account administration. Its repository already offers `ListUsersByOrg` in `internal/user/store/sqlite_tenancy.go`. `internal/auth/authn/identity.go` already carries `OrgID` and `Instance`. `internal/org/controller.go:RequireOperator` currently protects only organization-management routes. `internal/system/system.go:RegisterRoutes` currently uses an organization settings scope for installation mutations.

## Control flow and contracts

Enforce real caller, own-organization scope, and target organization inside service entry points, not only HTTP handlers. Use tenant-constrained updates and a transaction for the last-admin check plus mutation. Map absent/foreign targets to the same public error. Reuse a shared auth-owned installation-authority policy instead of importing the org controller from unrelated packages. In Organizations mode require `Instance`; in authenticated single-org mode require admin. Retain explicit synthetic compatibility only in local auth-disabled mode where the existing operation permits it; never permit synthetic account management.

## Security, persistence, and compatibility

Account list, role/status mutation, invite list/revoke, and account creation require tenant-consistent checks. Installation inventory must include backup downloads, logs that contain sensitive data, jobs, runtime flag endpoints, profile/global configuration, and non-HTTP transports. Do not infer tenancy from request-supplied organization IDs. Preserve existing last-operator and workspace privacy constraints. Full backup access is privileged recovery authority and is not evidence that operators are unable to access tenant data offline.

## Verification and failure behavior

Add `TestOrgAdminCannotListForeignAccounts`, `TestOrgAdminCannotMutateForeignAccount`, `TestLastAdminPerOrgConcurrent`, and `TestInstanceManagementRequiresOperator` beside current auth and system tests. Use two organizations, own/foreign targets, both SQL stores supported by the project, and failure injection. Preserve `TestManagementRoutesRejectSyntheticIdentity` and the existing workspace privacy tests.

Denial precedes externally observable side effects. Record sanitized decision
reasons without tokens, request bodies, private paths, or sensitive configuration.
Pending verification is not a successful implementation result.

## Related contracts and decisions

- [roles-and-scopes](../requirements/roles-and-scopes.md)
- [spec](../../multi-tenancy/spec.md)
- [Fork security policy](../../../decisions/2026-10-24-fork-security-boundaries.md)
- [Implementation plan](../../../plans/fork-security-hardening/plan.md)
