---
status: draft
system: auth
requirements:
  - REQ-AUTH-OPERATOR-CONTROL-001
  - REQ-AUTH-OPERATOR-CONTROL-002
---

# Operator control boundary system design

## Purpose and boundaries

Auth owns the distinction between a human administration principal and an agent/data-plane principal. This implements the accepted operator-owned launcher ADR. It does not claim to defend credentials against a malicious process already running as the operator OS account.

## Requirement mapping

| Requirement | Design sections |
| --- | --- |
| `REQ-AUTH-OPERATOR-CONTROL-001` | Components, control flow, and verification below |
| `REQ-AUTH-OPERATOR-CONTROL-002` | Components, control flow, and verification below |

## Restricted delivery

The shipped boundary is ControlPlaneBoundary after HTTP authentication. It denies unsafe preview, VS Code and canvas paths; apiContentWriter forces active API responses into sandboxed downloads. No separate content site is provisioned. HumanAdministrationSession requires a real admin browser session. The agent catalog uses an allowlist and sensitive MCP actions are refused by guardedMCPDispatcher. Separate-site and browser/mobile tests below remain restoration requirements.

## Components and responsibilities

`internal/agent/settings/handlers/handlers.go:registerHTTP` exposes profile mutations and full-detail reads. `internal/common/httpmw/interim_settings_interlock.go` is a replayable CSRF interlock. `internal/backendapp/helpers.go` projects its value in boot state. Preview forwarding is in `internal/agentctl/server/api/port_proxy.go`, `vscode_proxy.go`, `workspace_preview.go`; canvas routes are in `internal/backendapp/canvas_runtime_routes.go`. These require integration with the browser preview clients and auth session middleware.

## Control flow and contracts

First isolate generated content from the control-plane origin. Serve it through a dedicated content endpoint on a separate hostname/site, with narrowly scoped short-lived task capabilities and no operator cookie scope. Strip consumed control credentials in both forwarding hops; allow only approved parent messaging with exact origin/source checks. A deployment without a usable content endpoint must return a stable disabled/error result, not a same-origin fallback. Then enforce authenticated administration on mutation and secret-bearing read routes, retaining the boot interlock only as CSRF protection. Keep agent JWT and PAT identities distinct from the human browser session used for privileged settings.

## Security, persistence, and compatibility

Reuse native authentication and existing roles; do not introduce another boot-visible secret or trust Origin, CORS, source IP, user-agent, or headers as authentication. Auth-disabled installations must enable login before editing protected launch settings after this change. Global launch definitions require installation authority; workspace profiles also require workspace configuration authority. Existing privileged native plugins are trusted extensions, not an untrusted-content sandbox. Local/Worktree host access remains a separate executor limitation.

## Verification and failure behavior

Add `operator-control-boundary.spec.ts` and a mobile counterpart using the existing auth fixtures: preview scripts attempt parent DOM/storage reads and authenticated settings writes, including redirect and WebSocket paths. Add `TestOperatorSettingsRejectBootTokenReplay`, `TestSensitiveProfilesUseRedactedAgentCatalog`, and full route-table coverage. Supply actual HTTP and browser evidence before claiming this boundary is fixed.

Denial precedes externally observable side effects. Record sanitized decision
reasons without tokens, request bodies, private paths, or sensitive configuration.
Pending verification is not a successful implementation result.

## Related contracts and decisions

- [2026-07-24-operator-owned-agent-launcher-settings](../../../decisions/2026-07-24-operator-owned-agent-launcher-settings.md)
- [2026-09-05-trusted-browser-html-preview](../../../decisions/2026-09-05-trusted-browser-html-preview.md)
- [Fork security policy](../../../decisions/2026-10-24-fork-security-boundaries.md)
- [Implementation plan](../../../plans/fork-security-hardening/plan.md)
