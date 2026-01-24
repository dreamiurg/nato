# Quickstart: Automatic Dependency Updates

**Feature**: 002-auto-dependabot
**Validation Time**: ~5 minutes (manual verification after 24 hours for full cycle)

## Prerequisites

- GitHub repository with Go modules (`go.mod`)
- Existing CI workflow that runs tests
- Repository admin access (for settings changes)

## Quick Validation

### 1. Verify Dependabot Configuration

After merging, check that Dependabot is active:

```bash
# View Dependabot status in GitHub UI
# Go to: Insights → Dependency graph → Dependabot
```

Expected: Dependabot shows "gomod" ecosystem with daily schedule.

### 2. Verify Auto-Merge Workflow

```bash
# Check workflow exists
gh workflow list | grep -i dependabot
```

Expected output includes: `Dependabot Auto-Merge`

### 3. Verify Branch Protection

```bash
# Check branch protection rules
gh api repos/{owner}/{repo}/branches/main/protection --jq '.required_status_checks.contexts'
```

Expected: Your CI job name appears in the list.

### 4. Test Auto-Merge Flow (Optional - Simulated)

To test without waiting for a real update:

```bash
# Manually trigger Dependabot check
gh api -X POST /repos/{owner}/{repo}/dispatches \
  -f event_type=dependabot_check
```

Or wait 24 hours for the scheduled run.

## Full Validation Checklist

| Step | Expected Result | Verified |
|------|-----------------|----------|
| Dependabot creates PR for patch update | PR appears with `chore(deps):` prefix | [ ] |
| PR has correct labels | `dependencies`, `go` labels applied | [ ] |
| CI runs on Dependabot PR | Status checks appear on PR | [ ] |
| Auto-merge enables when CI passes | "Auto-merge enabled" badge on PR | [ ] |
| PR merges automatically | PR merged without manual approval | [ ] |
| Major update PR NOT auto-merged | Major version PR stays open | [ ] |

## Troubleshooting

### Dependabot PRs not appearing

1. Check `.github/dependabot.yml` syntax: `gh api /repos/{owner}/{repo}/vulnerability-alerts`
2. Verify Dependabot is enabled: Repository Settings → Security → Dependabot

### Auto-merge not triggering

1. Verify workflow file exists and has correct trigger
2. Check "Allow auto-merge" is enabled in repository settings
3. Ensure branch protection requires status checks

### CI not running on Dependabot PRs

1. Verify CI workflow triggers on `pull_request` event
2. Check Dependabot has permissions to trigger workflows

## Repository Settings Reference

### Enable Auto-Merge

1. Go to repository Settings
2. Under "Pull Requests", check "Allow auto-merge"

### Branch Protection Rule

1. Settings → Branches → Add rule
2. Branch name pattern: `main`
3. Enable:
   - [x] Require a pull request before merging
   - [x] Require status checks to pass before merging
   - [x] Require branches to be up to date before merging
4. Add required status check: (your CI job name)
