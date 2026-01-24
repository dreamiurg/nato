# Implementation Plan: GitHub Repository Hygiene Settings

**Spec**: specs/004-repo-hygiene/spec.md
**Created**: 2026-01-24

## Phase 0: Documentation Discovery (Complete)

### Sources Consulted
1. `gh repo edit --help` - CLI flags for repository settings
2. `gh api repos/dreamiurg/nato/branches/main/protection` - Current protection state
3. `.github/workflows/ci.yml` - CI job name for status checks

### Allowed APIs

**gh repo edit flags** (verified from --help):
- `--enable-squash-merge` - Enable squash merge
- `--enable-merge-commit=false` - Disable merge commit
- `--enable-rebase-merge=false` - Disable rebase merge
- `--delete-branch-on-merge` - Auto-delete head branches
- `--enable-auto-merge` - Enable auto-merge for Dependabot

**gh api for branch protection** (REST endpoint):
- Endpoint: `repos/{owner}/{repo}/branches/main/protection`
- Method: PUT
- Required fields: `required_status_checks`, `enforce_admins`, `restrictions`, `required_pull_request_reviews`

### Current State (from API response)
- `required_status_checks.checks`: `[]` (MISSING - no checks required)
- `enforce_admins.enabled`: `false` (MISSING - admins can bypass)
- `required_conversation_resolution.enabled`: `false` (MISSING)
- `allow_force_pushes.enabled`: `false` (OK)
- `allow_deletions.enabled`: `false` (OK)

### CI Job Name
- File: `.github/workflows/ci.yml`
- Job: `test` (line 13)
- This is the status check context name

---

## Phase 1: Configure Repository Settings

### What to Implement
Configure repository-level settings using `gh repo edit` command.

### Commands
```bash
gh repo edit dreamiurg/nato \
  --enable-squash-merge \
  --enable-merge-commit=false \
  --enable-rebase-merge=false \
  --delete-branch-on-merge \
  --enable-auto-merge
```

### Verification
```bash
# Check repository settings via API
gh api repos/dreamiurg/nato --jq '{
  squash_merge: .allow_squash_merge,
  merge_commit: .allow_merge_commit,
  rebase_merge: .allow_rebase_merge,
  delete_branch_on_merge: .delete_branch_on_merge,
  auto_merge: .allow_auto_merge
}'
```

Expected output:
```json
{
  "squash_merge": true,
  "merge_commit": false,
  "rebase_merge": false,
  "delete_branch_on_merge": true,
  "auto_merge": true
}
```

### Anti-Pattern Guards
- Do NOT use `--disable-*` flags (they don't exist, use `--enable-*=false`)
- Do NOT omit the repository argument when running from a worktree

---

## Phase 2: Configure Branch Protection

### What to Implement
Set branch protection rules on main branch using REST API via `gh api`.

### Commands
```bash
gh api repos/dreamiurg/nato/branches/main/protection \
  --method PUT \
  -H "Accept: application/vnd.github+json" \
  -f required_status_checks='{"strict":true,"checks":[{"context":"test"}]}' \
  -f enforce_admins=true \
  -f required_pull_request_reviews=null \
  -f restrictions=null \
  -F allow_force_pushes=false \
  -F allow_deletions=false \
  -F required_conversation_resolution=true
```

### Field Explanations
- `required_status_checks.checks[].context`: Must match CI job name `test`
- `strict: true`: Branch must be up-to-date before merge
- `enforce_admins: true`: Admins cannot bypass protection
- `required_pull_request_reviews: null`: Not requiring reviews (single-maintainer)
- `restrictions: null`: No push restrictions on specific users/teams
- `required_conversation_resolution: true`: Must resolve all review comments

### Verification
```bash
gh api repos/dreamiurg/nato/branches/main/protection --jq '{
  status_checks: .required_status_checks.checks,
  strict: .required_status_checks.strict,
  enforce_admins: .enforce_admins.enabled,
  conversation_resolution: .required_conversation_resolution.enabled,
  force_pushes: .allow_force_pushes.enabled,
  deletions: .allow_deletions.enabled
}'
```

Expected output:
```json
{
  "status_checks": [{"context": "test", "app_id": null}],
  "strict": true,
  "enforce_admins": true,
  "conversation_resolution": true,
  "force_pushes": false,
  "deletions": false
}
```

### Anti-Pattern Guards
- Do NOT use `--field` for boolean values (use `-F` for typed values)
- Do NOT omit `restrictions` and `required_pull_request_reviews` (API requires them, set to null)
- Do NOT use deprecated `contexts` array (use `checks` array instead)

---

## Phase 3: Verification

### Verification Checklist

1. **Repository settings verified**
   ```bash
   gh api repos/dreamiurg/nato --jq '.allow_squash_merge and (not .allow_merge_commit) and (not .allow_rebase_merge) and .delete_branch_on_merge and .allow_auto_merge'
   # Expected: true
   ```

2. **Branch protection verified**
   ```bash
   gh api repos/dreamiurg/nato/branches/main/protection --jq '
     (.required_status_checks.checks | map(.context) | contains(["test"])) and
     .enforce_admins.enabled and
     .required_conversation_resolution.enabled and
     (not .allow_force_pushes.enabled) and
     (not .allow_deletions.enabled)
   '
   # Expected: true
   ```

3. **Idempotency check** - Run both commands again, should succeed without error

### Success Criteria Mapping
- SC-001 (PRs require CI): `required_status_checks.checks` contains "test" ✓
- SC-002 (No force pushes): `allow_force_pushes.enabled` is false ✓
- SC-003 (Admins enforced): `enforce_admins.enabled` is true ✓
- SC-004 (Squash only): `allow_squash_merge` true, others false ✓
- SC-005 (Auto-delete branches): `delete_branch_on_merge` is true ✓

---

## Summary

| Phase | Action | Verification |
|-------|--------|--------------|
| 1 | `gh repo edit` for merge settings | API query confirms settings |
| 2 | `gh api` PUT for branch protection | API query confirms protection |
| 3 | Full verification | All success criteria pass |

**Total Commands**: 2 configuration commands + verification queries
**Estimated Complexity**: Low - direct API calls with known parameters
