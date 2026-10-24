---
status: draft
system: auth
created: 2026-10-24
owners:
  - rajatsinghtech
---

# Administration containment requirements

## Overview

Auth owns caller identity and administration authorization. These requirements specialize existing role/scope and multi-tenancy contracts at account and installation administration boundaries; workspace reach remains unchanged.

The [implementation record](../../../plans/fork-security-hardening/plan.md)
describes the delivered restrictions and verification limits. Separate-site preview
and isolated tenant execution contracts remain conditional: those features are
currently refused. Draft status does not assert that all future capabilities ship.

## Requirements

### REQ-AUTH-ADMIN-CONTAINMENT-001: Organization account boundary

**Intent:** Organization account boundary preserves the owning system's trust boundary.

#### Acceptance criteria

- **AC-AUTH-ADMIN-CONTAINMENT-001.1:** With Organizations enabled, an organization administrator shall list and manage only accounts in their authenticated organization.
- **AC-AUTH-ADMIN-CONTAINMENT-001.2:** A foreign or missing account target shall return the same not-found outcome without changing role, status, sessions, or tokens.
- **AC-AUTH-ADMIN-CONTAINMENT-001.3:** Last-administrator protection shall be evaluated within the affected organization, including concurrent requests; unresolved tenant identity shall fail closed.

### REQ-AUTH-ADMIN-CONTAINMENT-002: Installation authority

**Intent:** Installation authority preserves the owning system's trust boundary.

#### Acceptance criteria

- **AC-AUTH-ADMIN-CONTAINMENT-002.1:** With Organizations enabled, only a real instance operator shall download installation-wide backups or mutate installation-wide database, storage, restart, update, feature-toggle, or shared launch configuration.
- **AC-AUTH-ADMIN-CONTAINMENT-002.2:** An organization administrator without instance authority shall receive a denied response before a sensitive read, job, write, or notification occurs.
- **AC-AUTH-ADMIN-CONTAINMENT-002.3:** Operator status shall not grant workspace, transcript, or workspace-secret access; single-organization authenticated administrators shall retain their existing installation-management access.

## Out of scope

A vulnerability-free certification, automatic deployment, upstream disclosure,
and weakening existing workspace authorization. Delivery and external
prerequisites are tracked in the linked implementation plan.
