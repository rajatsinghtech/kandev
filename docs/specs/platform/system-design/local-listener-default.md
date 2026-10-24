---
status: draft
system: platform
requirements:
  - REQ-PLATFORM-LOCAL-LISTENER-001
---

# Local listener default system design

## Purpose and boundaries

Platform owns resolved startup configuration and listener selection. This fork changes the implicit network default while retaining explicit deployment configuration and the opt-in login contract.

## Requirement mapping

| Requirement | Design sections |
| --- | --- |
| `REQ-PLATFORM-LOCAL-LISTENER-001` | Components, control flow, and verification below |

## Components and responsibilities

`apps/backend/internal/common/config/config.go` owns `defaultServerHost` and `ServerConfig.ResolvedBinds`; `catalog.go` owns its catalog default. `internal/launcher/network.go` resolves access and readiness endpoints. `internal/launcher/startup_failure.go` reports provenance. Change these together; retain explicit IPv6 and wildcard support.

## Control flow and contracts

Resolve explicit environment, then the selected configuration file, then `127.0.0.1`. Keep launcher health probing consistent with the effective bind. Update wildcard-specific default expectations in config and launcher tests and public documentation. Explicit remote binds remain deliberate operator configuration and require authentication or a protected access gateway.

## Security, persistence, and compatibility

No schema or feature-flag migration. Reconcile the obsolete default-listener prose in `docs/specs/auth/requirements/auth.md` without changing its opt-in authentication acceptance criteria. Docker examples that require a wildcard inside the container must set it explicitly and publish only host loopback by default.

## Verification and failure behavior

`TestResolvedBinds`, `TestServerHostEnvOverridesConfigHosts`, `TestResolveBackendEndpointsPreservesDefaultBrowserOrigin`, and new `TestDefaultServerHostIsLoopback` cover default and explicit values. A temporary real-listener test must show loopback reachability and refusal through a non-loopback interface; skip with an explicit reason when the host has no such interface.

Denial precedes externally observable side effects. Record sanitized decision
reasons without tokens, request bodies, private paths, or sensitive configuration.
Pending verification is not a successful implementation result.

## Related contracts and decisions

- [auth](../../auth/requirements/auth.md)
- [startup-configuration-parity](../requirements/startup-configuration-parity.md)
- [Fork security policy](../../../decisions/2026-10-24-fork-security-boundaries.md)
- [Implementation plan](../../../plans/fork-security-hardening/plan.md)
