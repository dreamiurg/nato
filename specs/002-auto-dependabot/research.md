# Research: Automatic Dependency Updates

**Feature**: 002-auto-dependabot
**Date**: 2026-01-24

## Key Decisions

### 1. Dependabot Configuration Strategy

**Decision**: Use `dependabot.yml` with daily schedule, dependency grouping, and commit message prefixes compatible with release-please.

**Rationale**:
- Daily scanning ensures security patches are detected within 24 hours
- Grouping golang.org/x/* packages reduces PR noise (these are often updated together)
- `chore(deps):` prefix aligns with Conventional Commits used by release-please
- Setting `open-pull-requests-limit: 10` prevents PR flood

**Alternatives Considered**:
- Weekly scanning: Rejected - delays security patches unacceptably
- No grouping: Rejected - creates excessive PR noise for related packages
- Major version updates included: Rejected - breaks "hands-off" requirement

### 2. Auto-Merge Workflow Trigger

**Decision**: Use `pull_request` event trigger, NOT `pull_request_target`.

**Rationale**:
- `pull_request` runs in an isolated context with limited permissions
- `pull_request_target` has write permissions and is vulnerable to "pwn request" attacks
- Combined with `gh pr merge --auto`, this respects all branch protection rules

**Alternatives Considered**:
- `pull_request_target`: Rejected - security vulnerability (allows malicious forks to run privileged workflows)
- Manual merge only: Rejected - violates hands-off requirement

### 3. Bot Identity Verification

**Decision**: Check `github.event.pull_request.user.login == 'dependabot[bot]'` instead of `github.actor`.

**Rationale**:
- `github.actor` can be manipulated via "confused deputy" attacks
- `github.event.pull_request.user.login` identifies the actual PR author
- This prevents attackers from triggering auto-merge on malicious PRs

**Alternatives Considered**:
- `github.actor` check: Rejected - exploitable via `@dependabot recreate` on forked repos
- No verification: Rejected - critical security gap

### 4. Action Version Pinning

**Decision**: Pin GitHub Actions to full commit SHA, not version tags.

**Rationale**:
- Version tags can be moved by action maintainers (or attackers with access)
- SHA pinning ensures immutable, verifiable action code
- Eliminates supply chain attack vector

**Example**:
```yaml
# Secure (pinned to SHA)
uses: dependabot/fetch-metadata@5e5f99653a5b510e8555840e80cbf1514ad4af38

# Insecure (tag can be moved)
uses: dependabot/fetch-metadata@v2
```

### 5. Branch Protection Requirements

**Decision**: Require status checks to pass before merge; use `gh pr merge --auto` which respects these rules.

**Rationale**:
- `--auto` flag enables GitHub's native auto-merge, which waits for all required checks
- This ensures CI must pass before any merge occurs
- No need for custom wait loops or polling

**Required Settings**:
- Require pull request before merging: Enabled
- Require status checks to pass: Enabled (with CI job as required)
- Require branches to be up to date: Enabled

### 6. Update Type Handling

**Decision**: Auto-merge patch AND minor updates; block major updates.

**Rationale**:
- User explicitly requested "at least patch level" auto-merge
- Minor updates are generally safe when tests pass
- Major updates require manual review (potential breaking changes)
- Minor auto-merge further reduces maintenance burden

**Configuration**:
- Patch: Auto-merge when CI passes
- Minor: Auto-merge when CI passes
- Major: Create PR but require manual review (excluded from auto-merge workflow)

## Security Hardening Summary

| Layer | Protection | Implementation |
|-------|-----------|----------------|
| PR Source | Verify Dependabot identity | `github.event.pull_request.user.login` check |
| Code Quality | Require tests pass | Branch protection + required status checks |
| Supply Chain | Pin action versions | SHA-based `uses:` references |
| Workflow Injection | Prevent untrusted code execution | `pull_request` trigger (not `pull_request_target`) |
| Audit Trail | Track all merges | GitHub's native PR merge history |

## References

- [GitHub Docs: Dependabot Options Reference](https://docs.github.com/en/code-security/dependabot/working-with-dependabot/dependabot-options-reference)
- [GitHub Docs: Automating Dependabot with GitHub Actions](https://docs.github.com/en/code-security/dependabot/working-with-dependabot/automating-dependabot-with-github-actions)
- [GitHub Security Lab: Preventing Pwn Requests](https://securitylab.github.com/resources/github-actions-preventing-pwn-requests/)
- [BoostSecurity: Weaponizing Dependabot](https://boostsecurity.io/blog/weaponizing-dependabot-pwn-request-at-its-finest)
