# Implementation Plan: Automatic Dependency Updates

**Feature**: 002-auto-dependabot
**Created**: 2026-01-24
**Status**: Ready for Implementation

## Overview

Configure GitHub Dependabot with automatic PR creation for Go module updates, combined with a secure GitHub Actions workflow that auto-merges patch and minor updates when CI passes. Major updates require manual review.

## Documentation Sources

| Component | Documentation | SHA/Version |
|-----------|--------------|-------------|
| Dependabot config schema | [GitHub Docs: Dependabot Options Reference](https://docs.github.com/en/code-security/dependabot/working-with-dependabot/dependabot-options-reference) | v2 schema |
| fetch-metadata action | [dependabot/fetch-metadata](https://github.com/dependabot/fetch-metadata) | v2.5.0 / `21025c705c08248db411dc16f3619e6b5f9ea21a` |
| gh CLI merge command | `gh pr merge --help` | Latest |
| Security patterns | [specs/002-auto-dependabot/research.md](./research.md) | Local |

---

## Phase 0: Prerequisites

**Goal**: Ensure CI infrastructure exists before Dependabot automation.

### Task 0.1: Create CI Workflow

The repository has tests but no CI workflow. Dependabot auto-merge requires status checks to gate merges.

**Create**: `.github/workflows/ci.yml`

```yaml
name: CI

on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

permissions:
  contents: read

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - name: Run tests
        run: go test -v ./...
      - name: Run vet
        run: go vet ./...
```

**Verification**:
```bash
# Push to branch, verify workflow runs
gh run list --workflow=ci.yml
```

---

## Phase 1: Dependabot Configuration

**Goal**: Configure Dependabot to create PRs for Go module updates.

### Task 1.1: Create Dependabot Configuration

**Create**: `.github/dependabot.yml`

**Documentation Reference**: Dependabot Options Reference - `package-ecosystem`, `groups`, `schedule`, `commit-message`

```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "daily"
      time: "03:00"
      timezone: "UTC"
    open-pull-requests-limit: 10

    commit-message:
      prefix: "chore(deps)"
      include: "scope"

    labels:
      - "dependencies"
      - "go"

    groups:
      golang-org-x:
        patterns:
          - "golang.org/x/*"

      minor-and-patch:
        patterns:
          - "*"
        exclude-patterns:
          - "golang.org/x/*"
        update-types:
          - "minor"
          - "patch"
```

**Anti-Patterns to Avoid**:
- Do NOT use `update-types` at the top level (only valid inside `groups`)
- Do NOT include `major` in grouped `update-types` (defeats isolation purpose)

**Verification**:
```bash
# After merge to main, check Dependabot status
# GitHub UI: Insights → Dependency graph → Dependabot
gh api repos/dreamiurg/nato/vulnerability-alerts --silent && echo "Dependabot enabled"
```

---

## Phase 2: Auto-Merge Workflow

**Goal**: Create secure workflow to auto-merge patch/minor updates.

### Task 2.1: Create Auto-Merge Workflow

**Create**: `.github/workflows/dependabot-auto-merge.yml`

**Documentation References**:
- `pull_request` trigger (NOT `pull_request_target` - security)
- `github.event.pull_request.user.login` (NOT `github.actor` - security)
- fetch-metadata SHA: `21025c705c08248db411dc16f3619e6b5f9ea21a`
- Output: `update-type` with values `version-update:semver-{patch,minor,major}`

```yaml
name: Dependabot Auto-Merge

on: pull_request

permissions:
  contents: write
  pull-requests: write

jobs:
  dependabot-auto-merge:
    runs-on: ubuntu-latest
    if: github.event.pull_request.user.login == 'dependabot[bot]'
    steps:
      - name: Fetch Dependabot metadata
        id: dependabot-metadata
        uses: dependabot/fetch-metadata@21025c705c08248db411dc16f3619e6b5f9ea21a
        with:
          github-token: "${{ secrets.GITHUB_TOKEN }}"

      - name: Enable auto-merge for patch updates
        if: steps.dependabot-metadata.outputs.update-type == 'version-update:semver-patch'
        run: gh pr merge --auto --merge "${{ github.event.pull_request.html_url }}"
        env:
          GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"

      - name: Enable auto-merge for minor updates
        if: steps.dependabot-metadata.outputs.update-type == 'version-update:semver-minor'
        run: gh pr merge --auto --merge "${{ github.event.pull_request.html_url }}"
        env:
          GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"

      - name: Log major update (requires manual review)
        if: steps.dependabot-metadata.outputs.update-type == 'version-update:semver-major'
        run: |
          echo "::notice::Major version update detected for ${{ steps.dependabot-metadata.outputs.dependency-names }}"
          echo "This PR requires manual review before merging."
```

**Security Checklist**:
- [x] Uses `pull_request` trigger (isolated context)
- [x] Verifies `github.event.pull_request.user.login == 'dependabot[bot]'`
- [x] Action pinned to SHA (`21025c705c08248db411dc16f3619e6b5f9ea21a`)
- [x] Uses `gh pr merge --auto` (respects branch protection)
- [x] Minimal permissions (`contents: write`, `pull-requests: write`)

**Verification**:
```bash
# Check workflow exists
gh workflow list | grep -i dependabot
```

---

## Phase 3: Repository Settings (Manual)

**Goal**: Configure GitHub repository settings for auto-merge.

### Task 3.1: Enable Auto-Merge

**Location**: Repository Settings → General → Pull Requests

- [x] Allow auto-merge

### Task 3.2: Configure Branch Protection

**Location**: Repository Settings → Branches → Add rule

**Branch name pattern**: `main`

- [x] Require a pull request before merging
- [x] Require status checks to pass before merging
  - Add required check: `test` (from ci.yml)
- [x] Require branches to be up to date before merging

**Verification**:
```bash
gh api repos/dreamiurg/nato/branches/main/protection --jq '.required_status_checks.contexts'
# Expected: ["test"]
```

---

## Phase 4: Verification

**Goal**: Validate the complete implementation.

### Task 4.1: Verify Configuration Files

```bash
# Check dependabot.yml syntax (no output = valid)
python3 -c "import yaml; yaml.safe_load(open('.github/dependabot.yml'))"

# Check workflow syntax
gh workflow view dependabot-auto-merge.yml
gh workflow view ci.yml
```

### Task 4.2: Verify Security Controls

```bash
# Grep for insecure patterns
grep -r "pull_request_target" .github/workflows/ && echo "FAIL: Insecure trigger found" || echo "PASS"
grep -r "github.actor" .github/workflows/dependabot*.yml && echo "FAIL: Insecure actor check" || echo "PASS"
grep -E "uses:.*@v[0-9]" .github/workflows/dependabot*.yml && echo "FAIL: Tag-based pinning" || echo "PASS"
```

### Task 4.3: End-to-End Test

Wait 24 hours or manually trigger Dependabot:

1. Check Dependabot creates PR with correct labels
2. Verify CI runs on Dependabot PR
3. Confirm auto-merge enables when CI passes
4. Verify major update PR does NOT auto-merge

---

## Implementation Artifacts Summary

| File | Action | Purpose |
|------|--------|---------|
| `.github/workflows/ci.yml` | CREATE | Run tests on PRs (prerequisite for auto-merge) |
| `.github/dependabot.yml` | CREATE | Configure Dependabot for Go modules |
| `.github/workflows/dependabot-auto-merge.yml` | CREATE | Auto-merge patch/minor updates |
| Repository Settings | CONFIGURE | Enable auto-merge, branch protection |

---

## Risk Mitigations

| Risk | Mitigation |
|------|------------|
| Pwn request attack | Use `pull_request` trigger, verify PR author |
| Confused deputy attack | Check `github.event.pull_request.user.login` |
| Supply chain compromise | Pin actions to SHA |
| Untested code merge | Require CI via branch protection |
| Breaking changes | Block auto-merge for major versions |

---

## Success Criteria Mapping

| Requirement | Implementation |
|-------------|----------------|
| FR-001: Daily detection | `schedule.interval: "daily"` |
| FR-002: Separate PRs | Default Dependabot behavior |
| FR-003: Auto-merge patch | Workflow condition on `semver-patch` |
| FR-004: Auto-merge minor | Workflow condition on `semver-minor` |
| FR-005: Block major | Workflow logs but skips auto-merge |
| FR-006: Label PRs | `labels: ["dependencies", "go"]` |
| FR-007: Run tests | CI workflow + branch protection |
| FR-012: Group deps | `groups.golang-org-x` pattern |
| FR-013: Respect protection | `gh pr merge --auto` |
