---
status: draft
system: executors
created: 2026-10-24
owners:
  - rajatsinghtech
---

# Credential trust containment requirements

## Overview

Executors own process admission, injected credentials, and recovery. This extends per-user resource ownership with an execution boundary; an owner field, redirected HOME, or Git worktree alone is not a sandbox.

The [implementation record](../../../plans/fork-security-hardening/plan.md)
describes the delivered restrictions and verification limits. Separate-site preview
and isolated tenant execution contracts remain conditional: those features are
currently refused. Draft status does not assert that all future capabilities ship.

## Requirements

### REQ-EXECUTORS-CREDENTIAL-TRUST-001: Credential trust admission

**Intent:** Credential trust admission preserves the owning system's trust boundary.

#### Acceptance criteria

- **AC-EXECUTORS-CREDENTIAL-TRUST-001.1:** With Organizations enabled, Local and Worktree launches shall be refused before provisioning, credential copying, or process creation unless an OS-enforced organization execution boundary is established.
- **AC-EXECUTORS-CREDENTIAL-TRUST-001.2:** An isolated executor shall receive only the selected task/profile credentials and filesystem grants for the authenticated organization, and shall not inherit the backend user credential home.
- **AC-EXECUTORS-CREDENTIAL-TRUST-001.3:** Resume, prompt, terminal, preview, adoption, and cleanup paths shall revalidate execution ownership; an ambiguous or foreign runtime shall not be adopted, stopped, or deleted.

## Out of scope

A vulnerability-free certification, automatic deployment, upstream disclosure,
and weakening existing workspace authorization. Delivery and external
prerequisites are tracked in the linked implementation plan.
