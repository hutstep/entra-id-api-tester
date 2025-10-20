# GitHub Actions Workflows

This document explains the GitHub Actions workflows used in this project.

## Workflow Overview

```
┌─────────────────┐
│  Pull Request   │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  PR Validation Workflow             │
├─────────────────────────────────────┤
│  • Lint (fmt, vet, golangci-lint)   │
│  • Test (Go 1.25 + race detector)   │
│  • Build (Linux, macOS, Windows)    │
│  • Security (gosec, govulncheck)    │
│  • Dependency Review                │
└────────┬────────────────────────────┘
         │
         │ ✅ All checks pass
         ▼
┌─────────────────┐
│  Merge to Main  │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Release Workflow                   │
├─────────────────────────────────────┤
│  1. Validate                        │
│     • Lint & Test                   │
│     • Security Scans                │
│                                     │
│  2. Version Calculation             │
│     • Analyze commits               │
│     • Calculate next version        │
│                                     │
│  3. Create Tag                      │
│     • Create Git tag                │
│                                     │
│  4. Build & Release                 │
│     • GoReleaser build              │
│     • Multi-platform binaries       │
│     • GitHub Release                │
│     • Auto-generated changelog      │
└─────────────────────────────────────┘
```

## Workflows

### 1. PR Validation (`.github/workflows/pr-validation.yml`)

**Trigger:** Pull requests to `main`

**Purpose:** Ensure code quality before merging

**Jobs:**

#### Lint

- Code formatting check (`gofmt`)
- Static analysis (`go vet`)
- Advanced linting (`golangci-lint`)

#### Test

- Run tests with race detector
- Generate coverage reports
- Upload coverage artifacts

#### Build

- Cross-platform builds (Linux, macOS, Windows)
- Verify binary execution

#### Security

- Security scanning (`gosec`)
- Vulnerability check (`govulncheck`)
- Upload SARIF reports to GitHub Security

#### Dependency Review

- Check for vulnerable dependencies
- Fail on moderate+ severity issues

#### Summary

- Aggregate results from all jobs
- Show pass/fail status

**Duration:** ~3-5 minutes

---

### 2. Release (`.github/workflows/release.yml`)

**Trigger:** Push to `main` (after PR merge)

**Purpose:** Automated releases with semantic versioning

**Jobs:**

#### Validate (runs first)

- Quick validation before release
- Lint, test, and security checks
- Prevents bad releases

#### Release (runs after validate)

1. **Version Calculation**

   - Analyzes commit messages since last tag
   - Determines next version using semver
   - Based on conventional commits

2. **Tag Creation**

   - Creates Git tag (e.g., `v1.2.3`)
   - Skips if tag already exists

3. **Build & Release**

   - GoReleaser builds for all platforms
   - Creates GitHub Release
   - Uploads binaries and checksums
   - Generates changelog from commits

4. **Summary**
   - Posts release info to GitHub summary
   - Links to new release

**Duration:** ~5-10 minutes

---

## Key Differences

| Aspect                | PR Validation     | Release                        |
| --------------------- | ----------------- | ------------------------------ |
| **Trigger**           | Pull requests     | Push to main                   |
| **Purpose**           | Quality gate      | Create release                 |
| **Test Matrix**       | Single Go version | Single Go version              |
| **Build**             | All platforms     | All platforms (via GoReleaser) |
| **Security**          | Full scans        | Full scans                     |
| **Artifacts**         | Coverage reports  | Release binaries               |
| **Dependency Review** | Yes               | No (already done in PR)        |
| **Duration**          | 3-5 min           | 5-10 min                       |

---

## Optimization Strategy

### Why PR Validation Doesn't Run on Main Anymore

**Before:**

```yaml
on:
  pull_request:
    branches: [main]
  push:
    branches: [main] # ❌ Redundant!
```

**After:**

```yaml
on:
  pull_request:
    branches: [main] # ✅ Only PRs
```

**Reason:** The Release workflow now includes validation steps, so we don't need to run the full PR validation again after merge. This:

- ✅ Saves CI/CD minutes
- ✅ Faster feedback loop
- ✅ Still maintains quality (validation in release)
- ✅ Prevents redundant work

### Flow Example

```
Developer creates PR
  ↓
PR Validation runs (lint, test, build, security)
  ↓
PR approved & merged to main
  ↓
Release workflow runs
  ↓
Validate job runs (quick checks)
  ↓
Release job calculates version
  ↓
Creates tag & GitHub Release
```

---

## Workflow Configuration

### PR Validation

- **File:** `.github/workflows/pr-validation.yml`
- **Runs on:** Pull requests only
- **Can merge if:** All jobs pass
- **Blocks merge if:** Any job fails

### Release

- **File:** `.github/workflows/release.yml`
- **Runs on:** Push to main only
- **Creates release if:** Version bump detected
- **Skips if:** No version bump needed

---

## Troubleshooting

### "PR validation didn't run"

**Cause:** PR is not targeting `main` branch  
**Solution:** Change PR base branch to `main`

### "Release didn't trigger"

**Cause:** No conventional commits since last tag  
**Solution:** Ensure commits follow conventional format (feat:, fix:, etc.)

### "Both workflows running on main"

**Check:** This should NOT happen anymore after optimization  
**Solution:** Verify PR validation doesn't have `push: branches: [main]`

### "Tests run twice"

**Old behavior:** PR validation ran on both PR and push to main  
**New behavior:** PR validation only runs on PRs, release includes validation

---

## Best Practices

1. **Always create PRs** - Don't push directly to main
2. **Wait for PR validation** - Don't merge until all checks pass
3. **Use conventional commits** - Required for automatic versioning
4. **Review security alerts** - Check SARIF uploads in Security tab
5. **Monitor releases** - Check Actions tab after merge

---

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [golangci-lint](https://golangci-lint.run/)
- [GoReleaser](https://goreleaser.com/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
