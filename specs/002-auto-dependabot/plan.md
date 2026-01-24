# Implementation Plan: Automatic Dependency Updates

**Branch**: `002-auto-dependabot` | **Date**: 2026-01-24 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-auto-dependabot/spec.md`

## Summary

Configure GitHub Dependabot with automatic PR creation for Go module updates, combined with a secure GitHub Actions workflow that auto-merges patch and minor updates when CI passes. Major updates require manual review. Security hardening includes SHA-pinned actions, proper bot identity verification, and branch protection integration.

## Technical Context

**Language/Version**: YAML configuration files (GitHub Actions, Dependabot)
**Primary Dependencies**: GitHub Dependabot, GitHub Actions, `dependabot/fetch-metadata` action
**Storage**: N/A (configuration only)
**Testing**: Validation via actual dependency update cycle; no unit tests applicable
**Target Platform**: GitHub repository (github.com/dreamiurg/nato)
**Project Type**: Configuration addition to existing Go project
**Performance Goals**: Updates detected within 24 hours; merged within 1 hour of CI passing
**Constraints**: Must not require manual intervention for patch/minor updates
**Scale/Scope**: Single repository with ~6 direct Go dependencies

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitution is a template (not project-specific). No violations to check.

**Post-Design Check**: N/A - no code being written, only configuration files.

## Project Structure

### Documentation (this feature)

```text
specs/002-auto-dependabot/
├── spec.md              # Feature specification
├── plan.md              # This file
├── research.md          # Research findings and decisions
├── quickstart.md        # Validation guide
└── checklists/
    └── requirements.md  # Spec validation checklist
```

### Source Code (repository root)

```text
.github/
├── dependabot.yml           # NEW: Dependabot configuration
└── workflows/
    ├── release-please.yml   # Existing
    ├── release.yml          # Existing
    └── dependabot-auto-merge.yml  # NEW: Auto-merge workflow
```

**Structure Decision**: Configuration-only feature. All changes are YAML files in `.github/` directory. No source code modifications required.

## Implementation Artifacts

### File 1: `.github/dependabot.yml`

Dependabot configuration for Go modules with:
- Daily scanning at 03:00 UTC
- Dependency grouping (golang.org/x/*, other minor/patch)
- Conventional commit message format (`chore(deps):`)
- PR limit of 10 open PRs
- Labels for filtering

### File 2: `.github/workflows/dependabot-auto-merge.yml`

Secure auto-merge workflow with:
- `pull_request` trigger (not `pull_request_target`)
- Bot identity verification via `github.event.pull_request.user.login`
- SHA-pinned `dependabot/fetch-metadata` action
- Conditional auto-merge for patch/minor updates
- Logging for skipped major updates

### Repository Settings (Manual)

Required configuration in GitHub repository settings:
1. Enable "Allow auto-merge" in repository settings
2. Branch protection rule for `main`:
   - Require pull request before merging
   - Require status checks to pass (add existing CI job)
   - Require branches to be up to date before merging

## Complexity Tracking

No complexity violations. This feature adds 2 configuration files with no code changes.

## Security Measures

| Threat | Mitigation |
|--------|------------|
| Pwn request (malicious fork) | Use `pull_request` trigger, verify PR author |
| Confused deputy attack | Check `github.event.pull_request.user.login`, not `github.actor` |
| Supply chain (action compromise) | Pin actions to commit SHA |
| Untested code merge | Require CI status checks via branch protection |
| Breaking changes | Block auto-merge for major version updates |
