---
id: "04-untrusted-content-origin"
title: "Separate generated content from the operator origin"
status: implemented
wave: 2
depends_on: ["01-loopback-default"]
plan: "plan.md"
requirements:
  - REQ-AUTH-OPERATOR-CONTROL-001
acceptance_criteria:
  - AC-AUTH-OPERATOR-CONTROL-001.1
  - AC-AUTH-OPERATOR-CONTROL-001.2
  - AC-AUTH-OPERATOR-CONTROL-001.3
system_design:
  - ../../specs/auth/system-design/operator-control-boundary.md
---

# Task 04: Separate generated content from the operator origin

## Delivered scope

Delivered the fail-closed option: port proxy, embedded VS Code and canvas runtime return 503 until a separate site exists. Active API file responses receive attachment, sandbox CSP and nosniff headers. HTTP regressions pass. A separate preview site and browser/mobile adversarial testing are future restoration work, not completed isolation.

## Original design summary

Move active repository previews and generated canvases to a distinct content origin with task-scoped capabilities. Strip control credentials at every hop; deny active rendering when isolation is unavailable.

## In scope

The linked acceptance criteria, implementation files below, focused regression
coverage, and directly affected public documentation and companion records.

## Out of scope

Unrelated refactors, new product UI layouts, automatic publication, upstream
messaging, dependency upgrades without evidence, and claiming all vulnerabilities
are fixed. No delegation is authorized.

## Acceptance

- Preview scripts cannot read operator DOM/storage/cookies or issue authenticated operator changes.
- Redirects, nested frames, and WebSocket upgrades preserve the boundary.
- No same-origin fallback occurs when the content endpoint is unavailable.

## Verification

Run from the repository root in Git Bash or a supported Linux shell. Before any
frontend command in a fresh checkout, run `(cd apps && pnpm install --frozen-lockfile)`.
The commands below preserve the original validation design. Actual delivery and verification limits are recorded in Results; unimplemented browser/container scenario names are restoration work.

```bash
(cd apps/backend && go test ./internal/agentctl/server/api/... ./internal/canvas/... -count=1)
(cd apps/web && pnpm e2e:run --project auth tests/auth/operator-control-boundary.spec.ts)
(cd apps/web && pnpm e2e:run --project mobile-chrome tests/auth/mobile-operator-control-boundary.spec.ts)
```

Record explicit platform/service prerequisites and unsuccessful checks. For Go
race coverage use a supported compiler-enabled host and rerun the same packages
with `-race` when the change affects shared state.

## Files likely touched

- `apps/backend/internal/backendapp/canvas_runtime_routes.go`
- `apps/backend/internal/agentctl/server/api/port_proxy.go`
- `apps/backend/internal/agentctl/server/api/vscode_proxy.go`
- `apps/backend/internal/agentctl/server/api/workspace_preview.go`
- `apps/web/e2e/tests/auth/operator-control-boundary.spec.ts`
- `apps/web/e2e/tests/auth/mobile-operator-control-boundary.spec.ts`

## Dependencies

01-loopback-default

## Risks

Requires a deployable content-origin endpoint and browser evidence. Use a separate hostname/site, not only a different port. No new rendered controls are proposed; existing unavailable/error surfaces handle refusal.

## Parallelism

`sequential`. Waves express prerequisites only and do not authorize subagents.

## Inputs

- The referenced requirements and system design in full.
- [Plan and compatibility matrix](plan.md).
- [Evidence and scope corrections](evidence.md).
- Applicable scoped `AGENTS.md` and existing tests beside the affected code.

## Results

Delivered the fail-closed option: port proxy, embedded VS Code and canvas runtime return 503 until a separate site exists. Active API file responses receive attachment, sandbox CSP and nosniff headers. HTTP regressions pass. A separate preview site and browser/mobile adversarial testing are future restoration work, not completed isolation.
