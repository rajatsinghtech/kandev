---
id: "07-signed-desktop-publication"
title: "Gate trusted fork desktop publication on signing"
status: implemented
wave: 1
depends_on: []
plan: "plan.md"
requirements:
  - REQ-RELEASE-DESKTOP-TRUST-001
acceptance_criteria:
  - AC-RELEASE-DESKTOP-TRUST-001.1
  - AC-RELEASE-DESKTOP-TRUST-001.2
  - AC-RELEASE-DESKTOP-TRUST-001.3
system_design:
  - ../../specs/release/system-design/trusted-desktop-publication.md
---

# Task 07: Gate trusted fork desktop publication on signing

## Delivered scope

Implemented signing prerequisites before Stable prepare/publication, macOS verification, mandatory updater signatures and fork endpoint checks. Five credential-free policy tests pass. The inherited updater public key is a rejected placeholder; fork keys, OS certificates, real artifact verification and publication require owner setup. No release was dispatched.

## Original design summary

Add trusted-desktop publication gates to existing signing automation. Preserve unsigned development inspection builds, separate OS and updater signatures, and replace upstream update identity before fork publication.

## In scope

The linked acceptance criteria, implementation files below, focused regression
coverage, and directly affected public documentation and companion records.

## Out of scope

Unrelated refactors, new product UI layouts, automatic publication, upstream
messaging, dependency upgrades without evidence, and claiming all vulnerabilities
are fixed. No delegation is authorized.

## Acceptance

- Missing or failed signatures cannot produce trusted stable desktop downloads or update-feed promotion.
- Windows Authenticode and macOS notarization verification gate final artifact checksums/publication.
- No private keys are committed or sent to PR jobs; fork identity and endpoints are verified.

## Verification

Run from the repository root in Git Bash or a supported Linux shell. Before any
frontend command in a fresh checkout, run `(cd apps && pnpm install --frozen-lockfile)`.
The commands below preserve the original validation design. Actual delivery and verification limits are recorded in Results; unimplemented browser/container scenario names are restoration work.

```bash
python .github/scripts/release-workflow-contract_test.py
python .github/scripts/lint-action-pinning_test.py
zizmor .github/workflows
git diff --check
```

Record explicit platform/service prerequisites and unsuccessful checks. For Go
race coverage use a supported compiler-enabled host and rerun the same packages
with `-race` when the change affects shared state.

## Files likely touched

- `.github/workflows/release.yml`
- `.github/scripts/release-workflow-contract_test.py`
- `apps/desktop/src-tauri/tauri.conf.json`
- `docs/desktop-tauri-signing.md`
- `docs/plans/desktop-tauri-app/task-05-signing-and-release-trust.md`

## Dependencies

None.

## Risks

Real signed-artifact acceptance is blocked until the owner supplies protected Apple/Windows certificate inputs and fork updater identity. No release dispatch, secret mutation, or certificate purchase is authorized by this plan.

## Parallelism

`sequential`. Waves express prerequisites only and do not authorize subagents.

## Inputs

- The referenced requirements and system design in full.
- [Plan and compatibility matrix](plan.md).
- [Evidence and scope corrections](evidence.md).
- Applicable scoped `AGENTS.md` and existing tests beside the affected code.

## Results

Implemented signing prerequisites before Stable prepare/publication, macOS verification, mandatory updater signatures and fork endpoint checks. Five credential-free policy tests pass. The inherited updater public key is a rejected placeholder; fork keys, OS certificates, real artifact verification and publication require owner setup. No release was dispatched.
