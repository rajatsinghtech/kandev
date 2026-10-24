---
id: "01-loopback-default"
title: "Default implicit listeners to loopback"
status: implemented
wave: 1
depends_on: []
plan: "plan.md"
requirements:
  - REQ-PLATFORM-LOCAL-LISTENER-001
acceptance_criteria:
  - AC-PLATFORM-LOCAL-LISTENER-001.1
  - AC-PLATFORM-LOCAL-LISTENER-001.2
  - AC-PLATFORM-LOCAL-LISTENER-001.3
system_design:
  - ../../specs/platform/system-design/local-listener-default.md
---

# Task 01: Default implicit listeners to loopback

## Delivered scope

Implemented the resolved and catalog defaults, listener test, startup provenance, launcher browser compatibility and deployment documentation. Explicit wildcard binds remain supported. Targeted listener and launcher regressions pass locally.

## Original design summary

Change the resolved default, catalog, provenance text, regression tests, and deployment docs together. Keep container-internal wildcard binds explicit and host port publication private.

## In scope

The linked acceptance criteria, implementation files below, focused regression
coverage, and directly affected public documentation and companion records.

## Out of scope

Unrelated refactors, new product UI layouts, automatic publication, upstream
messaging, dependency upgrades without evidence, and claiming all vulnerabilities
are fixed. No delegation is authorized.

## Acceptance

- Implicit startup is loopback-only; explicit YAML/environment addresses still work.
- Existing local launch, health checks, and authentication mode remain usable.

## Verification

Run from the repository root in Git Bash or a supported Linux shell. Before any
frontend command in a fresh checkout, run `(cd apps && pnpm install --frozen-lockfile)`.
The commands below preserve the original validation design. Actual delivery and verification limits are recorded in Results; unimplemented browser/container scenario names are restoration work.

```bash
(cd apps/backend && go test ./internal/common/config ./internal/launcher)
(cd apps/backend && go test ./internal/auth -run TestModeDisabledByDefault -count=1)
```

Record explicit platform/service prerequisites and unsuccessful checks. For Go
race coverage use a supported compiler-enabled host and rerun the same packages
with `-race` when the change affects shared state.

## Files likely touched

- `apps/backend/internal/common/config/config.go`
- `apps/backend/internal/common/config/catalog.go`
- `apps/backend/internal/common/config/config_test.go`
- `apps/backend/internal/launcher/network_test.go`
- `apps/backend/internal/launcher/startup_failure.go`
- `docs/public/configuration.md`
- `docs/public/security.md`
- `docs/public/docker.md`

## Dependencies

None.

## Risks

Services and Docker deployment files that relied on an implicit wildcard require an explicit bind. Reconcile startup-configuration-parity companion records during implementation.

## Parallelism

`sequential`. Waves express prerequisites only and do not authorize subagents.

## Inputs

- The referenced requirements and system design in full.
- [Plan and compatibility matrix](plan.md).
- [Evidence and scope corrections](evidence.md).
- Applicable scoped `AGENTS.md` and existing tests beside the affected code.

## Results

Implemented the resolved and catalog defaults, listener test, startup provenance, launcher browser compatibility and deployment documentation. Explicit wildcard binds remain supported. Targeted listener and launcher regressions pass locally.
