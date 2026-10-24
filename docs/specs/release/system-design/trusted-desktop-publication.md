---
status: draft
system: release
requirements:
  - REQ-RELEASE-DESKTOP-TRUST-001
  - REQ-RELEASE-NATIVE-DEPENDENCY-001
---

# Trusted desktop publication system design

## Purpose and boundaries

Release owns artifact publication and publisher identity. This fork will separate development packages from trusted stable desktop distribution without changing the existing optional runtime-archive signing contract.

## Requirement mapping

| Requirement | Design sections |
| --- | --- |
| `REQ-RELEASE-DESKTOP-TRUST-001` | Components, control flow, and verification below |
| `REQ-RELEASE-NATIVE-DEPENDENCY-001` | Native dependency admission below |

## Restricted delivery

trusted-desktop-policy.py checks secret presence and fork updater configuration before Stable release preparation and again before publication. Native Windows verification is retained, macOS codesign/spctl/stapler verification is required, and unsigned updater manifests are refused. Credential-free policy tests do not prove a real certificate works. Owner signing setup and real artifact verification remain external prerequisites.

## Components and responsibilities

`.github/workflows/release.yml`, `.github/scripts/release-workflow-contract_test.py`, `docs/desktop-tauri-signing.md`, and desktop Tauri configuration own existing opportunistic signing and updater verification. Read `.github/AGENTS.md`, `apps/desktop/AGENTS.md`, and the release skill before implementation. Preserve isolation of secret-bearing publication jobs from untrusted PR code.

## Control flow and contracts

Separate eligibility for stable desktop publication from ordinary unsigned development packaging. Evaluate readiness before requesting protected signing secrets; require Authenticode verification on Windows and signature/notarization verification on macOS before uploads become public stable downloads. Verify Tauri updater payload signatures and fork public key separately from OS publisher signatures. Generate checksums after final signing. Do not reuse upstream update URLs or signing identity for fork artifacts.

## Security, persistence, and compatibility

Do not dispatch a release or purchase/configure certificates as part of this design. Owner-provided Apple/Windows signing credentials and protected environment setup are external prerequisites. Missing credentials must produce a clearly incomplete desktop publication result, not a successful trusted release. Existing development inspection artifacts remain available. Reconcile ADR 0026 optional desktop signing for the fork through the linked fork policy decision; do not modify the separate Windows runtime archive fallback contract.

## Verification and failure behavior

Extend `.github/scripts/release-workflow-contract_test.py` for absent/partial/mismatched signing inputs, verification failure, development-only mode, upstream endpoint rejection, checksum ordering, and no stable updater promotion. Use platform-native signature verification on real signed artifacts only after owner credentials are configured.

Denial precedes externally observable side effects. Record sanitized decision
reasons without tokens, request bodies, private paths, or sensitive configuration.
Pending verification is not a successful implementation result.

## Related contracts and decisions

### Native dependency admission

Cargo Audit on the first PR commit found five vulnerabilities. Lock
`rustls-webpki` 0.103.16 and `webpki-root-certs` 1.0.10 with checksums from the
official sparse registry; both patch releases preserve their dependency edges.
Stable preparation runs cargo-audit 0.22.2 before version or tag mutations.

The updater pins ZIP major 4. The fixes for RUSTSEC-2026-0367, -0368 and -0369
require ZIP 9; a lock-only major upgrade does not satisfy the updater constraint.
Keep those findings visible. Updater replacement or a reviewed extraction
backport, with adversarial archive and native updater tests, is separate work.
No advisory ignore or major dependency substitution is included here. Source
use is not a clean native security audit, and trusted publication stays blocked.

- [windows-runtime-signing](../requirements/windows-runtime-signing.md)
- [0026-tauri-desktop-shell](../../../decisions/0026-tauri-desktop-shell.md)
- [Fork security policy](../../../decisions/2026-10-24-fork-security-boundaries.md)
- [Implementation plan](../../../plans/fork-security-hardening/plan.md)
