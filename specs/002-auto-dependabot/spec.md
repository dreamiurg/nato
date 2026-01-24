# Feature Specification: Automatic Dependency Updates

**Feature Branch**: `002-auto-dependabot`
**Created**: 2026-01-24
**Status**: Draft
**Input**: User description: "Set up automatic dependabot on nato repo with all bells and whistles in CI so that dependency updates are fully automatic (at least patch level updates when deps are updated). I want this to be secure as fuck and very hands off"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Hands-Off Patch Updates (Priority: P1)

As a repository owner, I want patch-level dependency updates to be automatically merged without any manual intervention, so that security patches and bug fixes are applied immediately without requiring my attention.

**Why this priority**: Patch updates represent the lowest-risk changes (bug fixes, security patches) and are the most frequent. Automating these eliminates the majority of dependency maintenance burden while maintaining safety.

**Independent Test**: When a dependency releases a patch version (e.g., 1.2.3 → 1.2.4), the update should be automatically detected, tested, and merged within 24 hours without any human action.

**Acceptance Scenarios**:

1. **Given** a Go dependency has a new patch release, **When** the dependency scanner runs, **Then** a pull request is automatically created with the update
2. **Given** a patch update PR is created, **When** all CI checks pass, **Then** the PR is automatically merged without human approval
3. **Given** a patch update PR is created, **When** any CI check fails, **Then** the PR remains open and the owner is notified

---

### User Story 2 - Secure Minor Version Updates (Priority: P2)

As a repository owner, I want minor version updates to be automatically proposed but require CI validation before auto-merge, so that new features from dependencies are adopted safely.

**Why this priority**: Minor updates may introduce new features or behavior changes. They should be tested but can still be auto-merged if tests pass, reducing manual review burden.

**Independent Test**: When a dependency releases a minor version (e.g., 1.2.0 → 1.3.0), a PR is created, and if all tests pass, it auto-merges. If tests fail, it awaits manual review.

**Acceptance Scenarios**:

1. **Given** a Go dependency has a new minor release, **When** the dependency scanner runs, **Then** a pull request is created with the update
2. **Given** a minor update PR is created, **When** all CI checks pass including full test suite, **Then** the PR is automatically merged
3. **Given** a minor update PR fails any CI check, **When** the failure is detected, **Then** the PR is labeled for manual review and owner is notified

---

### User Story 3 - Major Version Awareness (Priority: P3)

As a repository owner, I want to be notified about major version updates but have them require manual approval, so that breaking changes are reviewed before adoption.

**Why this priority**: Major updates may contain breaking changes that require code modifications. These should never be auto-merged but should be tracked and surfaced.

**Independent Test**: When a dependency releases a major version (e.g., 1.x → 2.x), a PR is created but remains open for manual review, clearly labeled as requiring attention.

**Acceptance Scenarios**:

1. **Given** a Go dependency has a new major release, **When** the dependency scanner runs, **Then** a pull request is created labeled as "major update"
2. **Given** a major update PR is created, **When** any automation runs, **Then** the PR is NOT auto-merged regardless of CI status
3. **Given** a major update PR exists, **When** viewing the repository, **Then** it is clearly distinguishable from patch/minor updates

---

### User Story 4 - Security-First Configuration (Priority: P1)

As a repository owner, I want the automation to follow security best practices with minimal permissions and protection against supply chain attacks, so that the automated process cannot be exploited.

**Why this priority**: Security is critical - the automation should not introduce vulnerabilities or attack vectors. This is equal priority with basic functionality.

**Independent Test**: Review the configuration to verify it uses minimal required permissions, validates update sources, and cannot be manipulated by malicious PRs.

**Acceptance Scenarios**:

1. **Given** the automation is configured, **When** reviewing permissions, **Then** only minimum necessary permissions are granted
2. **Given** an update PR is created, **When** the update is processed, **Then** it originates from verified official sources only
3. **Given** a malicious actor attempts to trigger auto-merge, **When** the PR doesn't meet all security criteria, **Then** the automation refuses to merge

---

### Edge Cases

- What happens when multiple dependency updates are available simultaneously? (Each should be handled as separate PRs)
- How does the system handle a dependency that is yanked/removed? (PR should be created to remove or alert)
- What if CI is temporarily broken for unrelated reasons? (Updates should queue but not merge until CI passes)
- What if a patch update actually breaks tests? (Should not auto-merge; owner notified)
- What happens during rate limiting from package registries? (Graceful retry with backoff)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST automatically detect new versions of Go module dependencies at least daily
- **FR-002**: System MUST create separate pull requests for each dependency update
- **FR-003**: System MUST automatically merge patch-level updates when all CI checks pass
- **FR-004**: System MUST automatically merge minor-level updates when all CI checks pass
- **FR-005**: System MUST NOT automatically merge major-level updates regardless of CI status
- **FR-006**: System MUST label PRs by update type (patch, minor, major) for easy identification
- **FR-007**: System MUST run the full test suite on all dependency update PRs before any merge decision
- **FR-008**: System MUST notify repository owner when any auto-merge fails or when manual review is needed
- **FR-009**: System MUST use minimal permissions required for its operation (no write access beyond what's needed)
- **FR-010**: System MUST only process updates from official Go module proxy sources
- **FR-011**: System MUST include security vulnerability information when available in PR descriptions
- **FR-012**: System MUST support grouping of related dependencies to reduce PR noise (e.g., all golang.org/x/* together)
- **FR-013**: System MUST respect branch protection rules and not bypass required checks
- **FR-014**: System MUST maintain an audit trail of all automated merges

### Key Entities

- **Dependency Update**: Represents a version change for a Go module (source, current version, new version, update type)
- **Update Policy**: Rules defining how different update types should be handled (auto-merge criteria, notification rules)
- **CI Status**: The pass/fail state of automated checks that gate merge decisions

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of patch-level updates are merged automatically within 24 hours when tests pass, with zero manual intervention
- **SC-002**: 100% of minor-level updates are merged automatically within 24 hours when tests pass, with zero manual intervention
- **SC-003**: 0% of major-level updates are auto-merged (all require manual approval)
- **SC-004**: Repository owner spends less than 5 minutes per week on routine dependency maintenance
- **SC-005**: Zero security incidents caused by the automation itself (no privilege escalation, no unauthorized merges)
- **SC-006**: All automated merges have full CI test suite passing as prerequisite

## Assumptions

- The repository already has a working CI pipeline with comprehensive tests
- Branch protection rules are enabled on the main branch
- The repository uses Go modules for dependency management
- GitHub is the hosting platform (Dependabot is a GitHub feature)
- The owner accepts the risk of auto-merging patch and minor updates that pass tests
- Notifications will be delivered via GitHub's default notification system

## Out of Scope

- Updates to non-Go dependencies (GitHub Actions versions, etc.) - can be added later
- Custom approval workflows beyond GitHub's built-in mechanisms
- Integration with external security scanning tools beyond what GitHub provides
- Rollback automation if a merged update later causes issues in production
