---
id: "06-credential-execution-boundary"
title: "Fail closed on unisolated tenant execution"
status: implemented
wave: 2
depends_on: ["02-tenant-account-administration", "03-instance-administration"]
plan: "plan.md"
requirements:
  - REQ-EXECUTORS-CREDENTIAL-TRUST-001
acceptance_criteria:
  - AC-EXECUTORS-CREDENTIAL-TRUST-001.1
  - AC-EXECUTORS-CREDENTIAL-TRUST-001.2
  - AC-EXECUTORS-CREDENTIAL-TRUST-001.3
system_design:
  - ../../specs/executors/system-design/credential-trust-containment.md
---

# Task 06: Fail closed on unisolated tenant execution

## Delivered scope

Delivered the conservative admission option: all Organizations startup is refused, both direct config and persisted feature overrides, before agent services start. Config and backend-startup regressions pass. No multi-org executor is admitted; dedicated credential roots, isolated executor restoration and OS/container canaries remain future work.

## Original design summary

Reconcile multi-tenancy credential and executor work orders, enforce admission before provisioning, and prove isolation with cross-tenant canary credentials. Refuse unsupported host execution rather than labeling worktrees as sandboxes.

## In scope

The linked acceptance criteria, implementation files below, focused regression
coverage, and directly affected public documentation and companion records.

## Out of scope

Unrelated refactors, new product UI layouts, automatic publication, upstream
messaging, dependency upgrades without evidence, and claiming all vulnerabilities
are fixed. No delegation is authorized.

## Acceptance

- Organizations mode denies unisolated Local/Worktree launches before any credential or process side effect.
- Approved isolated execution receives only its own grants and rejects foreign runtime recovery.
- Canary tests cover launch, resume, prompt, terminal, and cleanup without host credential inheritance.

## Verification

Run from the repository root in Git Bash or a supported Linux shell. Before any
frontend command in a fresh checkout, run `(cd apps && pnpm install --frozen-lockfile)`.
The commands below preserve the original validation design. Actual delivery and verification limits are recorded in Results; unimplemented browser/container scenario names are restoration work.

```bash
(cd apps/backend && go test ./internal/agent/executor/... ./internal/agent/credentials/... ./internal/agent/runtime/lifecycle/... -count=1)
(cd apps/web && KANDEV_E2E_CONTAINERS=1 pnpm e2e:run --project containers tests/docker/credential-trust-boundary.spec.ts)
```

Record explicit platform/service prerequisites and unsuccessful checks. For Go
race coverage use a supported compiler-enabled host and rerun the same packages
with `-race` when the change affects shared state.

## Files likely touched

- `apps/backend/internal/agent/runtime/lifecycle/executor_admission.go`
- `apps/backend/internal/agent/runtime/lifecycle/executor_ownership_test.go`
- `apps/backend/internal/agent/runtime/lifecycle/executor_backend.go`
- `apps/backend/internal/agent/credentials`
- `apps/backend/internal/agent/executor`
- `apps/backend/internal/task/repository/sqlite`
- `docs/plans/multi-tenancy/task-10-per-org-agent-credentials.md`
- `docs/plans/multi-tenancy/task-11-executor-tenant-pinning.md`

## Dependencies

02-tenant-account-administration, 03-instance-administration

## Risks

Real Docker/OS isolation tests need an appropriate host. Initial Windows behavior may deny unsupported multi-org executors. Resource ownership metadata alone is insufficient.

## Parallelism

`sequential`. Waves express prerequisites only and do not authorize subagents.

## Inputs

- The referenced requirements and system design in full.
- [Plan and compatibility matrix](plan.md).
- [Evidence and scope corrections](evidence.md).
- Applicable scoped `AGENTS.md` and existing tests beside the affected code.

## Results

Delivered the conservative admission option: all Organizations startup is refused, both direct config and persisted feature overrides, before agent services start. Config and backend-startup regressions pass. No multi-org executor is admitted; dedicated credential roots, isolated executor restoration and OS/container canaries remain future work.
