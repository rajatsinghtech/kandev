# Security hardening evidence and assumptions

## Checkout and authorization

- Upstream: `https://github.com/kdlbs/kandev`; user fork: `https://github.com/rajatsinghtech/kandev`.
- Source baseline: `0b81ae48d4e2c358f80ac489e39388dfd8d7e6ab`, upstream main fetched 2026-10-24.
- Local branch: `security/hardening-design`. `origin` is the user fork; `upstream` is the original project.
- Confirmed user intent: fork Kandev and begin addressing the five previously discussed security concerns.
- The repo's root AGENTS.md and spec-driven-development skill require this complete design handoff before production or permanent test edits. This package is unstaged/uncommitted for review.
- No publishing, automatic merge/deployment, release dispatch, secret mutation, certificate purchase, or upstream report has occurred. No subagents were used.

## Source-confirmed observations

| Prior concern | Source evidence at baseline | Classification and response |
| --- | --- | --- |
| Listener and login defaults | `internal/common/config/config.go:defaultServerHost` is `0.0.0.0`; `catalog.go` duplicates that default. Auth middleware has an explicit synthetic local identity. | Insecure network default; task 01 changes the implicit listener. Explicit remote deployments still need a protected boundary. |
| Cross-org user administration | `internal/auth/service_users.go:ListUsers` calls global `ListUsers`; `AdminSetRoleStatus` reads a target by ID without checking OrgID. `internal/user/store/sqlite.go` queries those rows without tenant predicates. `ListUsersByOrg` already exists in `sqlite_tenancy.go`. | Confirmed missing tenant check by static inspection; task 02 adds service and query boundaries plus two-org regressions. No exploit test has been run. |
| Installation authority | `internal/system/system.go:RegisterRoutes` gates install mutations with `ScopeOrgSettingsManage`. `authn.Identity.Instance` and `org.RequireOperator` already exist for org administration. | Confirmed mismatch between org capability and installation authority; task 03 reuses existing identity concepts. |
| Agent settings boundary | `internal/common/httpmw/interim_settings_interlock.go` explicitly describes replayability; the boot projection in `internal/backendapp/helpers.go` exposes the token; settings mutations use org config scope and the interlock. | Intentional interim defense, not proof of human operator authentication. Tasks 04 and 05 jointly address origin and principal separation. CORS or loopback alone is insufficient. |
| Shared host credentials | Existing runtime admission/ownership tests protect logical access, while Local/Worktree processes still run with the host account. Multi-tenancy task 10/11 already describe separate credential roots and executor pinning. | Architectural trust limitation; task 06 must enforce real isolation or refuse unsupported shapes. Do not call a worktree or owner field a sandbox. |
| Unsigned desktop artifacts | `docs/desktop-tauri-signing.md` and `.github/workflows/release.yml` explicitly permit unsigned development packages when inputs are absent. | Supply-chain policy concern, not evidence of malware. Task 07 gates trusted publication. Real signatures require owner credentials. |

All `internal/...` source paths in the table are relative to `apps/backend/`.

## Corrections to the preliminary review

The desktop launcher already chooses loopback; the wildcard concern concerns
default CLI/backend/service paths. Authentication and operator fields exist;
the issue is incomplete enforcement, not absence of all access control. The
installer signing machinery also exists. No CVE, exploit success, dependency
scan result, or complete audit is claimed. Published documentation tracks main
and can differ from a released binary; work is pinned to the commit above.

## Assumption check and scope

The existing product owners remain Platform for startup configuration, Auth
for administration principals and trust boundaries, Executors for process and
credential admission, and Release for distribution identity. No new system or
UI layout is needed. Preserve trusted single-user local development; deny
unsafe multi-org/privileged operations when their prerequisite is absent.
This is a conservative security choice within the requested remediation, not
a user assertion that a particular deployment shape is already configured.

Not yet established: a deployable distinct preview site, Windows OS isolation
support, the fork's signing certificates and updater keys, or container test
infrastructure. These are explicit task prerequisites; their absence cannot
be counted as passing implementation verification.

## Companion packages

Reconcile, do not duplicate or declare completed, the existing multi-tenancy
task 08 (operator), task 10 (per-org credentials), task 11 (executor pinning),
desktop-tauri-app task 05 (signing), and startup-configuration-parity records.
The existing Windows runtime signing contract remains unchanged because it
specifically excludes desktop installers. The operator-owned launcher ADR
already requires origin separation and authenticated human authority. Task 05
must amend the old byte-identical-auth-disabled wording in
AC-AUTH-AUTH-001.1 when the new protected-settings contract is implemented;
the draft replacement policy is not silently promoted to active here.

## Implementation follow-through

The table records upstream source at the audited base revision. The fork now
contains fixes and fail-closed restrictions described in the [implementation
record](plan.md). Two-org and authority regression tests reproduced the relevant
failures before remediation. No exploit against a deployed third-party instance
was attempted. Companion multi-tenancy plans remain incomplete because the fork
denies Organizations entirely; trusted publication remains blocked on owner keys.
