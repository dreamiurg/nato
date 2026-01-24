# Implementation Plan: Automatic Dependency Updates

**Branch**: `002-auto-dependabot` | **Date**: 2026-01-24 | **Spec**: [spec.md](./spec.md)

## Summary

Configure GitHub Dependabot with automatic PR creation for Go module updates, combined with a secure GitHub Actions workflow that auto-merges patch and minor updates when CI passes. Major updates require manual review. Security hardening includes SHA-pinned actions, proper bot identity verification, and branch protection integration.

---

## Phase 0: Documentation Discovery (Complete)

### Allowed APIs & Patterns

**Dependabot Configuration** (`dependabot.yml`):
- Source: [GitHub Dependabot Options Reference](https://docs.github.com/en/code-security/dependabot/working-with-dependabot/dependabot-options-reference)
- `package-ecosystem: "gomod"` - Go modules ecosystem
- `schedule.interval: "daily"` - Daily scanning (runs Mon-Fri)
- `schedule.time: "03:00"` - UTC time for check
- `groups.<name>.patterns` - Glob patterns for grouping
- `groups.<name>.update-types` - Filter by: `major`, `minor`, `patch`
- `commit-message.prefix` - Conventional commit prefix
- `commit-message.include: "scope"` - Adds `deps` scope
- `open-pull-requests-limit: 10` - Max concurrent PRs
- `labels` - Array of labels to apply

**Fetch-Metadata Action**:
- Source: [dependabot/fetch-metadata](https://github.com/dependabot/fetch-metadata)
- **SHA for v2**: `21025c705c08248db411dc16f3619e6b5f9ea21a`
- Outputs:
  - `update-type` - Values: `version-update:semver-major`, `version-update:semver-minor`, `version-update:semver-patch`
  - `dependency-names` - Comma-separated package names
  - `dependency-type` - `direct:production`, `direct:development`, `indirect`

**Auto-Merge Command**:
- Source: [GitHub CLI Manual](https://cli.github.com/manual/gh_pr_merge)
- `gh pr merge --auto --squash` - Enable auto-merge with squash
- Requires: Repository setting "Allow auto-merge" enabled
- Requires: Branch protection with required status checks

### Anti-Patterns to Avoid

| DO NOT | DO INSTEAD |
|--------|------------|
| Use `pull_request_target` trigger | Use `pull_request` trigger |
| Check `github.actor` for bot identity | Check `github.event.pull_request.user.login` |
| Pin actions to version tags (`@v2`) | Pin to SHA (`@21025c705...`) |
| Use `--admin` flag to bypass checks | Let branch protection enforce rules |

---

## Phase 1: Dependabot Configuration

### Task 1.1: Create `.github/dependabot.yml`

**Copy from documentation pattern**, adapting for this project:

```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "daily"
      time: "03:00"
      timezone: "UTC"
    commit-message:
      prefix: "chore"
      include: "scope"
    labels:
      - "dependencies"
      - "go"
    open-pull-requests-limit: 10
    groups:
      golang-x:
        patterns:
          - "golang.org/x/*"
        update-types:
          - "minor"
          - "patch"
      all-minor-patch:
        patterns:
          - "*"
        exclude-patterns:
          - "golang.org/x/*"
        update-types:
          - "minor"
          - "patch"
```

**Verification**:
```bash
# Validate YAML syntax
python3 -c "import yaml; yaml.safe_load(open('.github/dependabot.yml'))"

# Confirm file structure
cat .github/dependabot.yml
```

---

## Phase 2: Auto-Merge Workflow

### Task 2.1: Create `.github/workflows/dependabot-auto-merge.yml`

**Copy from documentation pattern** with security hardening:

```yaml
name: Dependabot Auto-Merge

on: pull_request

permissions:
  contents: write
  pull-requests: write

jobs:
  auto-merge:
    runs-on: ubuntu-latest
    if: github.event.pull_request.user.login == 'dependabot[bot]'
    steps:
      - name: Fetch Dependabot metadata
        id: metadata
        uses: dependabot/fetch-metadata@21025c705c08248db411dc16f3619e6b5f9ea21a
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}

      - name: Auto-merge patch and minor updates
        if: contains(fromJSON('["version-update:semver-patch", "version-update:semver-minor"]'), steps.metadata.outputs.update-type)
        run: gh pr merge --auto --squash "$PR_URL"
        env:
          PR_URL: ${{ github.event.pull_request.html_url }}
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Log skipped major update
        if: steps.metadata.outputs.update-type == 'version-update:semver-major'
        run: |
          echo "::notice::Major update detected for ${{ steps.metadata.outputs.dependency-names }}. Manual review required."
```

**Security checklist**:
- [ ] Uses `pull_request` trigger (not `pull_request_target`)
- [ ] Checks `github.event.pull_request.user.login` (not `github.actor`)
- [ ] Action pinned to SHA `21025c705c08248db411dc16f3619e6b5f9ea21a`
- [ ] Uses `--squash` for clean history
- [ ] Excludes major updates from auto-merge

**Verification**:
```bash
# Validate YAML syntax
python3 -c "import yaml; yaml.safe_load(open('.github/workflows/dependabot-auto-merge.yml'))"

# Verify SHA pinning
grep -q "21025c705c08248db411dc16f3619e6b5f9ea21a" .github/workflows/dependabot-auto-merge.yml && echo "SHA pinned correctly"

# Verify trigger is pull_request (not pull_request_target)
grep -q "^on: pull_request$" .github/workflows/dependabot-auto-merge.yml && echo "Secure trigger"
```

---

## Phase 3: Repository Settings (Manual)

### Task 3.1: Enable Auto-Merge

**Location**: GitHub Repository Settings → Pull Requests

**Required settings**:
1. Check "Allow auto-merge"

### Task 3.2: Configure Branch Protection

**Location**: GitHub Repository Settings → Branches → Branch protection rules → `main`

**Required settings**:
1. Check "Require a pull request before merging"
2. Check "Require status checks to pass before merging"
   - Add required status check (if CI exists)
3. Check "Require branches to be up to date before merging"

**Verification command**:
```bash
gh api repos/{owner}/{repo} --jq '.allow_auto_merge'
# Should return: true

gh api repos/{owner}/{repo}/branches/main/protection --jq '.required_status_checks'
# Should show configured checks
```

---

## Phase 4: Verification

### Task 4.1: Validate Configuration Files

```bash
# Check both files exist
ls -la .github/dependabot.yml .github/workflows/dependabot-auto-merge.yml

# Validate YAML syntax
for f in .github/dependabot.yml .github/workflows/dependabot-auto-merge.yml; do
  python3 -c "import yaml; yaml.safe_load(open('$f'))" && echo "$f: valid"
done
```

### Task 4.2: Security Audit

```bash
# Verify no pull_request_target usage
! grep -r "pull_request_target" .github/workflows/ && echo "No pull_request_target: PASS"

# Verify SHA pinning on fetch-metadata
grep "dependabot/fetch-metadata@21025c705c08248db411dc16f3619e6b5f9ea21a" .github/workflows/dependabot-auto-merge.yml && echo "SHA pinned: PASS"

# Verify proper bot identity check
grep "github.event.pull_request.user.login" .github/workflows/dependabot-auto-merge.yml && echo "Bot identity check: PASS"
```

### Task 4.3: Functional Test

After pushing to main:
1. Wait for Dependabot's first scan (up to 24 hours, or trigger manually via GitHub UI)
2. Verify PR is created with correct labels
3. Verify auto-merge triggers when CI passes
4. Verify major updates are NOT auto-merged

---

## Complexity Summary

| Metric | Value | Notes |
|--------|-------|-------|
| Files added | 2 | Both YAML configuration |
| Lines of code | ~50 | Configuration only |
| Manual steps | 2 | Repository settings |
| Security controls | 4 | SHA pinning, bot verification, secure trigger, branch protection |

---

## Requirements Traceability

| Requirement | Implementation |
|-------------|----------------|
| FR-001: Daily detection | `schedule.interval: "daily"` |
| FR-002: Separate PRs | Default Dependabot behavior |
| FR-003: Auto-merge patch | `update-type` check in workflow |
| FR-004: Auto-merge minor | `update-type` check in workflow |
| FR-005: Block major auto-merge | Excluded from `fromJSON()` array |
| FR-006: Label by type | `labels` in dependabot.yml |
| FR-009: Minimal permissions | `contents: write, pull-requests: write` only |
| FR-012: Grouping | `groups` configuration |
| FR-013: Respect branch protection | `gh pr merge --auto` honors protection |
