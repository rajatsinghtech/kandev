---
created: 2026-10-24
status: implemented
requirements:
  - REQ-PLATFORM-LOCAL-LISTENER-001
  - REQ-AUTH-ADMIN-CONTAINMENT-001
  - REQ-AUTH-ADMIN-CONTAINMENT-002
  - REQ-AUTH-OPERATOR-CONTROL-001
  - REQ-AUTH-OPERATOR-CONTROL-002
  - REQ-EXECUTORS-CREDENTIAL-TRUST-001
  - REQ-RELEASE-DESKTOP-TRUST-001
system_design:
  - ../../specs/platform/system-design/local-listener-default.md
  - ../../specs/auth/system-design/administration-containment.md
  - ../../specs/auth/system-design/operator-control-boundary.md
  - ../../specs/executors/system-design/credential-trust-containment.md
  - ../../specs/release/system-design/trusted-desktop-publication.md
legacy_specs: []
---

# Implementation record: restricted fork security hardening

## Outcome and scope

The owner requested implementation on 2026-10-24 after the design handoff.
The fork now implements loopback defaults, tenant administration containment,
installation authority checks, human-only launcher configuration, redacted agent
discovery, and signed-desktop publication admission. Unsupported active-content
transports and Organizations mode fail closed. This is a restricted delivery,
not a claim that a new preview sandbox or tenant execution sandbox exists.

The original requirements retain draft status where they describe future isolated
execution or preview restoration. The conditional paths are unavailable in this
fork. No vulnerability-free certification, deployment, upstream disclosure,
certificate provisioning, or release publication is included.

## Delivered behavior

| Area | Implementation | Evidence and limit |
| --- | --- | --- |
| Listener | Implicit IPv4 loopback; explicit binds preserved | Actual listener and configuration tests |
| Tenant accounts | Org-scoped users/invites and transactional last-admin guard | Two-org and concurrent SQLite tests; live PostgreSQL not exercised |
| Installation authority | Central HTTP boundary after authentication | Route denial tests prove handlers have no side effects; no workspace grant added |
| Agent configuration | Human admin sessions required; PAT/runtime/synthetic identities rejected | HTTP identity matrix, catalog projection and MCP dispatcher regressions |
| Active content | Port proxy, embedded VS Code and canvas runtime denied; active API downloads sandboxed | HTTP tests; separate preview site and browser isolation tests deferred until restoration |
| Tenant execution | Organizations startup refused before runtime initialization, including persisted flags | Config and persisted startup tests; no OS/container isolation claimed |
| Desktop publication | Missing signing inputs or inherited updater identity stop Stable publication | Five policy tests and workflow contracts; real signing requires owner credentials |

The backend still trusts code executing as its OS account. Use separate OS accounts
and instances for distinct trust groups. Existing upstream npm and container builds
do not contain these fork changes. Build this branch from source until fork artifacts
are deliberately published. Existing privileged native extensions remain trusted.

## Work orders

- [x] [Task 01: Loopback default](task-01-loopback-default.md)
- [x] [Task 02: Tenant administration](task-02-tenant-account-administration.md)
- [x] [Task 03: Installation authority](task-03-instance-administration.md)
- [x] [Task 04: Deny unsafe active-content transports](task-04-untrusted-content-origin.md)
- [x] [Task 05: Human configuration principals](task-05-operator-settings-principals.md)
- [x] [Task 06: Refuse unsupported Organizations execution](task-06-credential-execution-boundary.md)
- [x] [Task 07: Signed publication admission](task-07-signed-desktop-publication.md)

## Validation

The new security regression tests were observed failing before their fixes.
Local Windows checks passed for full auth and user-store packages; selected
listener/launcher, agent DTO/controller/handler, MCP authority, and persisted
startup regressions; desktop policy; workflow action pinning; public docs;
and harness validation. Detailed commands are in the security regression workflow.

Full Windows launcher/config tests include POSIX executable and permission
assumptions; release contract tests require POSIX tools and GPG; one spec-linter
test requires symlink privileges unavailable here. A broad agent-controller test
run stalled on Windows; the targeted authority tests pass. None of those broader
Windows runs is reported as a pass. Linux CI runs the affected auth/listener suite
with race detection, release contracts, and complete spec-linter tests.

The full workflow scanner reports inherited findings. The new security regression
workflow has no scanner findings. Existing unrelated CI findings are not remediated
or concealed by this change. Real signed installers, browser preview isolation,
live PostgreSQL, and Docker/OS canary tests were not run.

## Restoration and external prerequisites

Restoring active previews requires a separate site, task capabilities and browser
adversarial tests. Restoring Organizations requires real execution isolation,
credential grants and cross-tenant canaries; owner metadata or a worktree is
insufficient. Existing multi-tenancy companion plans are not marked complete.
Trusted desktop distribution requires fork-owned signing credentials, updater keys
and protected release setup. The inherited public key is retained only as a rejected
placeholder; publication is blocked until the owner replaces it. No private keys
were created, requested, committed or configured.

See [source evidence](evidence.md), the linked designs, and the final GitHub PR for
the exact commit, hosted validation and unresolved findings. The policy favors
explicitly unavailable features over an unverified security boundary.
