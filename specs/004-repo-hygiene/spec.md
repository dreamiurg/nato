# Feature Specification: GitHub Repository Hygiene Settings

**Feature Branch**: `004-repo-hygiene`
**Created**: 2026-01-24
**Status**: Complete
**Input**: Configure GitHub repo hygiene settings via gh CLI including branch protection, merge strategies, and auto-merge configuration

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Required CI Status Checks (Priority: P1)

As a repository owner, I want all pull requests to require passing CI checks before merge, so that broken code cannot be merged to main.

**Why this priority**: CI validation is the fundamental quality gate. Without this, any code can reach main regardless of test failures.

**Independent Test**: Create a PR with failing tests - it should be blocked from merging until tests pass.

**Acceptance Scenarios**:

1. **Given** a PR is opened against main, **When** CI checks have not completed, **Then** the merge button is disabled
2. **Given** a PR has failing CI checks, **When** attempting to merge, **Then** the merge is blocked with clear explanation
3. **Given** a PR has all CI checks passing, **When** other requirements are met, **Then** the merge is allowed

---

### User Story 2 - Enforce Rules for Administrators (Priority: P1)

As a repository owner, I want branch protection rules to apply to all users including administrators, so that no one can bypass quality gates.

**Why this priority**: Admin bypass undermines the entire protection system. If admins can skip checks, protection is advisory only.

**Independent Test**: As an admin, attempt to merge a PR with failing checks - it should still be blocked.

**Acceptance Scenarios**:

1. **Given** branch protection is configured, **When** an admin attempts to force push to main, **Then** the push is rejected
2. **Given** a PR has failing CI checks, **When** an admin attempts to merge, **Then** the merge is blocked
3. **Given** protection is enforced for admins, **When** reviewing settings, **Then** enforce_admins is enabled

---

### User Story 3 - Required Conversation Resolution (Priority: P2)

As a code reviewer, I want unresolved review comments to block merge, so that feedback is addressed before code is merged.

**Why this priority**: Prevents merging code where review comments were ignored or forgotten.

**Independent Test**: Leave a review comment without resolving it - PR should be unmergeable.

**Acceptance Scenarios**:

1. **Given** a PR has unresolved review comments, **When** attempting to merge, **Then** the merge is blocked
2. **Given** all review comments are resolved, **When** other requirements are met, **Then** the merge is allowed
3. **Given** a new comment is added after resolution, **When** attempting to merge, **Then** the merge is blocked again

---

### User Story 4 - Consistent Merge Strategy (Priority: P2)

As a repository owner, I want a consistent merge strategy (squash) for all PRs, so that the commit history stays clean and linear.

**Why this priority**: Squash merge keeps main history clean with one commit per PR. This aids bisecting and changelog generation.

**Independent Test**: Merge a multi-commit PR - it should result in a single squash commit on main.

**Acceptance Scenarios**:

1. **Given** the repository settings, **When** merging a PR, **Then** only squash merge is available
2. **Given** a PR with multiple commits, **When** squash merged, **Then** main receives a single commit
3. **Given** standard merge or rebase is attempted, **When** through UI, **Then** those options are disabled

---

### User Story 5 - Auto-Delete Head Branches (Priority: P3)

As a repository owner, I want feature branches to be automatically deleted after merge, so that stale branches don't accumulate.

**Why this priority**: Reduces clutter and makes it clear which branches are active.

**Independent Test**: Merge a PR - the source branch should be automatically deleted.

**Acceptance Scenarios**:

1. **Given** a PR is merged, **When** merge completes, **Then** the head branch is automatically deleted
2. **Given** the repository settings, **When** reviewing configuration, **Then** delete_branch_on_merge is enabled

---

### Edge Cases

- What if CI workflow is renamed? (Status check name in protection must be updated)
- What if an admin needs to make an emergency fix? (They can temporarily disable protection via API)
- What happens to existing PRs when protection is enabled? (They must meet new requirements before merge)
- What if the test job is split into multiple jobs? (Each required job must be added to status checks)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Branch protection MUST require the "test" CI status check to pass before merge
- **FR-002**: Branch protection MUST enforce rules for administrators (no bypass)
- **FR-003**: Branch protection MUST require conversation resolution before merge
- **FR-004**: Branch protection MUST block force pushes to main
- **FR-005**: Branch protection MUST block branch deletion for main
- **FR-006**: Repository MUST use squash merge as the only allowed merge strategy
- **FR-007**: Repository MUST auto-delete head branches after PR merge
- **FR-008**: Repository MUST have auto-merge enabled to work with Dependabot automation
- **FR-009**: All settings MUST be configured via gh CLI for reproducibility
- **FR-010**: Configuration MUST be idempotent (running multiple times produces same result)

### Key Entities

- **Branch Protection Rule**: GitHub configuration controlling what can be pushed/merged to a branch
- **Status Check**: CI job result that gates merge (pass/fail)
- **Merge Strategy**: How commits are combined when merging (squash/merge/rebase)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of PRs require passing CI before merge (verified by attempting merge with failing tests)
- **SC-002**: 0% of force pushes succeed on main branch (verified by attempting force push)
- **SC-003**: Administrators cannot bypass protection rules (verified by admin merge attempt with failing CI)
- **SC-004**: All merged PRs result in single squash commits (verified by checking commit history)
- **SC-005**: Feature branches are deleted automatically after merge (verified by checking branches after merge)

## Assumptions

- The repository already has a working CI pipeline (`test` job in `.github/workflows/ci.yml`)
- GitHub CLI (gh) is available and authenticated with appropriate permissions
- The user has admin access to the repository
- The repository is hosted on GitHub

## Out of Scope

- Required pull request reviews (single-maintainer project)
- CODEOWNERS configuration
- Protected tags
- Ruleset configuration (using classic branch protection)
- Webhook configuration

## Implementation Notes

GitHub CLI (`gh`) does not have native commands for branch protection. Configuration must use:

```bash
# Repository settings (merge strategy, auto-delete, auto-merge)
gh repo edit --enable-squash-merge --disable-merge-commit --disable-rebase-merge \
  --delete-branch-on-merge --enable-auto-merge

# Branch protection (requires REST API via gh api)
gh api repos/{owner}/{repo}/branches/main/protection --method PUT \
  --field required_status_checks='{"strict":true,"checks":[{"context":"test"}]}' \
  --field enforce_admins=true \
  --field required_conversation_resolution=true \
  --field restrictions=null \
  --field required_pull_request_reviews=null \
  --field allow_force_pushes=false \
  --field allow_deletions=false
```
