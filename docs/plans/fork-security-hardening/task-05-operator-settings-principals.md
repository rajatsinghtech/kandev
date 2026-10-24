---
id: "05-operator-settings-principals"
title: "Require real principals for privileged agent settings"
status: implemented
wave: 3
depends_on: ["03-instance-administration", "04-untrusted-content-origin"]
plan: "plan.md"
requirements:
  - REQ-AUTH-OPERATOR-CONTROL-002
acceptance_criteria:
  - AC-AUTH-OPERATOR-CONTROL-002.1
  - AC-AUTH-OPERATOR-CONTROL-002.2
  - AC-AUTH-OPERATOR-CONTROL-002.3
system_design:
  - ../../specs/auth/system-design/operator-control-boundary.md
---

# Task 05: Require real principals for privileged agent settings

## Delivered scope

Protected launcher and MCP configuration routes with real human administration sessions; rejected boot replay, agent/PAT and synthetic authority. Agent catalog projections are allowlisted and configuration tunnel actions are denied at dispatch. Targeted controller/handler, DTO and MCP regressions pass. Authentication must be enabled for privileged edits.

## Original design summary

Protect launch-affecting mutations and sensitive reads with real administration identity and resource scope. Retain redacted agent discovery; the boot token remains only a CSRF interlock.

## In scope

The linked acceptance criteria, implementation files below, focused regression
coverage, and directly affected public documentation and companion records.

## Out of scope

Unrelated refactors, new product UI layouts, automatic publication, upstream
messaging, dependency upgrades without evidence, and claiming all vulnerabilities
are fixed. No delegation is authorized.

## Acceptance

- Replaying boot state without a real human session cannot mutate settings or read sensitive configuration.
- Agent runtime JWTs and omitted credentials cannot acquire human authority.
- Authorized settings edits work; auth-disabled installations must enable login before privileged edits.

## Verification

Run from the repository root in Git Bash or a supported Linux shell. Before any
frontend command in a fresh checkout, run `(cd apps && pnpm install --frozen-lockfile)`.
The commands below preserve the original validation design. Actual delivery and verification limits are recorded in Results; unimplemented browser/container scenario names are restoration work.

```bash
(cd apps/backend && go test ./internal/agent/settings/... ./internal/common/httpmw/... ./internal/auth/... ./internal/mcp/... -count=1)
(cd apps/web && pnpm e2e:run --project auth tests/auth/operator-control-boundary.spec.ts)
(cd apps/web && pnpm e2e:run --project mobile-chrome tests/auth/mobile-operator-control-boundary.spec.ts)
```

Record explicit platform/service prerequisites and unsuccessful checks. For Go
race coverage use a supported compiler-enabled host and rerun the same packages
with `-race` when the change affects shared state.

## Files likely touched

- `apps/backend/internal/agent/settings/handlers/handlers.go`
- `apps/backend/internal/agent/settings/handlers/interim_settings_interlock_test.go`
- `apps/backend/internal/common/httpmw/interim_settings_interlock.go`
- `apps/backend/internal/backendapp/helpers.go`
- `apps/backend/internal/mcp`
- `apps/web/e2e/tests/auth/operator-control-boundary.spec.ts`
- `docs/specs/auth/requirements/auth.md`
- `docs/public/security.md`

## Dependencies

03-instance-administration, 04-untrusted-content-origin

## Risks

Do not ship or claim this boundary before task 04. Existing PAT settings clients need an explicit restricted API policy; do not simply remove the current bearer rejection globally.

## Parallelism

`sequential`. Waves express prerequisites only and do not authorize subagents.

## Inputs

- The referenced requirements and system design in full.
- [Plan and compatibility matrix](plan.md).
- [Evidence and scope corrections](evidence.md).
- Applicable scoped `AGENTS.md` and existing tests beside the affected code.

## Results

Protected launcher and MCP configuration routes with real human administration sessions; rejected boot replay, agent/PAT and synthetic authority. Agent catalog projections are allowlisted and configuration tunnel actions are denied at dispatch. Targeted controller/handler, DTO and MCP regressions pass. Authentication must be enabled for privileged edits.
