# ADR-2026-10-24-fork-security-boundaries: Fork security boundaries

**Status:** accepted
**Date:** 2026-10-24
**Area:** backend, infra, workflow

## Context

The fork owner requested remediation of documented network exposure,
administration, operator/agent separation, credential isolation, and desktop
publisher-trust concerns. Upstream prioritizes login-free local use and
opportunistic desktop signing. Existing identity and runtime controls provide
a base, but they do not establish every promised boundary.

## Decision

This fork uses loopback as the implicit listener while preserving explicit
deployment configuration and opt-in login. Privileged launch configuration
requires an authenticated human administration principal, with generated
active content on a separate site before that boundary is claimed. Follow
the existing operator-owned launcher ADR rather than using a boot token as
authentication. Multi-org execution is denied unless an OS/runtime credential
boundary is established. Trusted stable desktop publication requires verified
publisher signing and a fork-controlled update identity; unsigned artifacts
remain development builds. Actual publication and credential configuration
are separate owner operations.

The restricted implementation denies unsafe active-content transports and all
Organizations startup. It does not provide a new preview or executor sandbox.
Implementation and validation limits are tracked by the [implementation record](../plans/fork-security-hardening/plan.md) and its
requirements and designs.

## Consequences

Some remote launchers require explicit bind configuration. Protected settings
require login, and unsupported multi-org executors/previews can become
unavailable rather than silently retain an unsafe fallback. Desktop trusted
distribution cannot complete before signing prerequisites exist. Workspace
authorization remains independent of operator administration; backup recovery
authority is explicitly privileged.

## Alternatives Considered

- Warning-only network defaults retain accidental LAN exposure.
- Enabling experimental authentication everywhere changes local onboarding
  without closing filesystem or same-origin script access.
- A replayable SPA token or CORS policy does not authenticate direct clients.
- A Git worktree or separate HOME without OS isolation does not contain a
  process running under the same account.
- Publishing unsigned builds as trusted stable downloads retains the original
  publisher-identity concern. Preserve them only as development artifacts.
