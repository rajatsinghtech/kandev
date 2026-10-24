---
status: draft
system: release
created: 2026-10-24
owners:
  - rajatsinghtech
---

# Trusted desktop publication requirements

## Overview

Release owns artifact publication and publisher identity. This fork will separate development packages from trusted stable desktop distribution without changing the existing optional runtime-archive signing contract.

The [implementation record](../../../plans/fork-security-hardening/plan.md)
describes the delivered restrictions and verification limits. Separate-site preview
and isolated tenant execution contracts remain conditional: those features are
currently refused. Draft status does not assert that all future capabilities ship.

## Requirements

### REQ-RELEASE-DESKTOP-TRUST-001: Stable desktop identity

**Intent:** Stable desktop identity preserves the owning system's trust boundary.

#### Acceptance criteria

- **AC-RELEASE-DESKTOP-TRUST-001.1:** A stable Windows or macOS desktop installer shall be published as a trusted release only after publisher signature verification succeeds; macOS shall additionally pass notarization verification.
- **AC-RELEASE-DESKTOP-TRUST-001.2:** Missing, invalid, expired, mismatched, or failed signing inputs shall prevent trusted desktop publication and update-feed promotion; unsigned inspection builds shall remain explicitly development artifacts.
- **AC-RELEASE-DESKTOP-TRUST-001.3:** Signing private keys shall remain in protected release secrets, never in source, logs, artifacts, agents, or pull-request workflows; fork releases shall use the fork owner identity and fork update endpoints.

## Out of scope

A vulnerability-free certification, automatic deployment, upstream disclosure,
and weakening existing workspace authorization. Delivery and external
prerequisites are tracked in the linked implementation plan.
