---
status: draft
system: executors
requirements:
  - REQ-EXECUTORS-CREDENTIAL-TRUST-001
---

# Credential trust containment system design

## Purpose and boundaries

Executors own process admission, injected credentials, and recovery. This extends per-user resource ownership with an execution boundary; an owner field, redirected HOME, or Git worktree alone is not a sandbox.

## Requirement mapping

| Requirement | Design sections |
| --- | --- |
| `REQ-EXECUTORS-CREDENTIAL-TRUST-001` | Components, control flow, and verification below |

## Restricted delivery

The shipped admission boundary is ValidateExecutionIsolation in configuration validation and after persisted runtime flag resolution, before services start. Organizations mode is always refused. No container, remote, Local or Worktree executor is admitted as an isolated tenant shape. The executor and canary design below is future restoration work; the existing companion plans remain incomplete.

## Components and responsibilities

Use `internal/agent/runtime/lifecycle` admission and `SetExecutorAdmission`, building on `executor_ownership_test.go`. Integrate `internal/agent/executor`, runtime executor backends and credential preparation, plus persisted execution metadata in `internal/task/repository/sqlite`. Existing `docs/plans/multi-tenancy/task-10-per-org-agent-credentials.md` and `task-11-executor-tenant-pinning.md` already identify the companion work; update those records when implementation changes their state.

## Control flow and contracts

Preserve trusted single-user local operation. For Organizations mode, default-deny host executors unless an OS identity boundary is configured and verified. Reuse tenant-scoped profiles and durable runtime ownership; fail closed on missing ownership, credential selection, or boundary evidence. An initial implementation may deny unsupported multi-org executor shapes rather than claim isolation. For approved container/remote shapes, require dedicated credential homes, bounded mounts, no host Docker socket or privileged host access, and organization-pinned resources.

## Security, persistence, and compatibility

Do not create a new trusted-standalone bypass to pass tests. Separate containers without safe mounts and runtime policy do not prove credential isolation. Windows support may initially use separate service accounts/instances or deny the unsupported multi-org shape. A shared Kandev backend remains trusted with credentials it injects. Revocation must occur at credential providers; deleting a setting or stopping an agent is insufficient.

## Verification and failure behavior

Extend `TestExecutorOwnershipLaunchFailsBeforeRuntimeAdmission`, `TestExecutorOwnershipTerminalCredentialsFailClosed`, and prompt-admission tests. Add `TestCredentialTrustBoundaryBeforeProvisioning` with failing dependencies and `TestForeignRuntimeCannotBeAdopted`. Use real container/OS identity tests to prove one tenant cannot read a canary token or file of another; test launch, resume, terminal, and cleanup.

Denial precedes externally observable side effects. Record sanitized decision
reasons without tokens, request bodies, private paths, or sensitive configuration.
Pending verification is not a successful implementation result.

## Related contracts and decisions

- [per-user-ownership](../requirements/per-user-ownership.md)
- [task-11-executor-tenant-pinning](../../../plans/multi-tenancy/task-11-executor-tenant-pinning.md)
- [Fork security policy](../../../decisions/2026-10-24-fork-security-boundaries.md)
- [Implementation plan](../../../plans/fork-security-hardening/plan.md)
