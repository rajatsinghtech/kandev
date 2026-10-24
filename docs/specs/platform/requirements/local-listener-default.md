---
status: draft
system: platform
created: 2026-10-24
owners:
  - rajatsinghtech
---

# Local listener default requirements

## Overview

Platform owns resolved startup configuration and listener selection. This fork changes the implicit network default while retaining explicit deployment configuration and the opt-in login contract.

The [implementation record](../../../plans/fork-security-hardening/plan.md)
describes the delivered restrictions and verification limits. Separate-site preview
and isolated tenant execution contracts remain conditional: those features are
currently refused. Draft status does not assert that all future capabilities ship.

## Requirements

### REQ-PLATFORM-LOCAL-LISTENER-001: Implicit local listener

**Intent:** Implicit local listener preserves the owning system's trust boundary.

#### Acceptance criteria

- **AC-PLATFORM-LOCAL-LISTENER-001.1:** Without an explicit listener setting, CLI, service, and backend startup shall accept connections only on IPv4 loopback.
- **AC-PLATFORM-LOCAL-LISTENER-001.2:** An explicit valid listener setting shall retain configuration precedence, including environment overrides and configured address lists.
- **AC-PLATFORM-LOCAL-LISTENER-001.3:** Changing the listener default shall not silently enable or disable authentication or change the database, repositories, or credentials.

## Out of scope

A vulnerability-free certification, automatic deployment, upstream disclosure,
and weakening existing workspace authorization. Delivery and external
prerequisites are tracked in the linked implementation plan.
