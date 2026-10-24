---
status: draft
system: auth
created: 2026-10-24
owners:
  - rajatsinghtech
---

# Operator control boundary requirements

## Overview

Auth owns the distinction between a human administration principal and an agent/data-plane principal. This implements the accepted operator-owned launcher ADR. It does not claim to defend credentials against a malicious process already running as the operator OS account.

The [implementation record](../../../plans/fork-security-hardening/plan.md)
describes the delivered restrictions and verification limits. Separate-site preview
and isolated tenant execution contracts remain conditional: those features are
currently refused. Draft status does not assert that all future capabilities ship.

## Requirements

### REQ-AUTH-OPERATOR-CONTROL-001: Untrusted content separation

**Intent:** Untrusted content separation preserves the owning system's trust boundary.

#### Acceptance criteria

- **AC-AUTH-OPERATOR-CONTROL-001.1:** Agent-controlled previews, repository HTML, and generated canvases shall not execute scripts with the operator application origin or receive its cookies, authorization headers, boot data, or storage.
- **AC-AUTH-OPERATOR-CONTROL-001.2:** When a separate content origin cannot be established, active content shall fail closed rather than fall back to operator-origin rendering.
- **AC-AUTH-OPERATOR-CONTROL-001.3:** Authorized preview streaming and WebSocket application traffic shall retain task-bound authorization and shall never disclose executor or operator control credentials.

### REQ-AUTH-OPERATOR-CONTROL-002: Authenticated configuration authority

**Intent:** Authenticated configuration authority preserves the owning system's trust boundary.

#### Acceptance criteria

- **AC-AUTH-OPERATOR-CONTROL-002.1:** Launch-affecting mutations and sensitive full-detail configuration reads shall require a real human administration identity authorized for the resource scope.
- **AC-AUTH-OPERATOR-CONTROL-002.2:** Agent runtime credentials, synthetic identities, missing credentials, and a replayed boot interlock alone shall not authorize those operations.
- **AC-AUTH-OPERATOR-CONTROL-002.3:** Agents shall retain a redacted discovery catalog containing only authorized identifiers, names, and capability metadata; removing sensitive values shall apply to HTTP, WebSocket, MCP, and boot projections.

## Out of scope

A vulnerability-free certification, automatic deployment, upstream disclosure,
and weakening existing workspace authorization. Delivery and external
prerequisites are tracked in the linked implementation plan.
