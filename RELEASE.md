# Release Process

This project uses automated releases powered by [GoReleaser](https://goreleaser.com/) and semantic versioning based on [Conventional Commits](https://www.conventionalcommits.org/).

## How It Works

### 1. Conventional Commits

All commits to the `main` branch should follow the Conventional Commits specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**

- `feat`: A new feature (triggers MINOR version bump)
- `fix`: A bug fix (triggers PATCH version bump)
- `perf`: Performance improvement (triggers PATCH version bump)
- `refactor`: Code refactoring (triggers PATCH version bump)
- `docs`: Documentation changes (no version bump)
- `test`: Test changes (no version bump)
- `ci`: CI/CD changes (no version bump)
- `chore`: Maintenance tasks (no version bump)
- `style`: Code style changes (no version bump)

**Breaking Changes:** Add `BREAKING CHANGE:` in the footer to trigger a MAJOR version bump:

```
feat: add new authentication method

BREAKING CHANGE: The old authentication method has been removed.
Users must migrate to the new method.
```

### 2. Automatic Release on Merge to Main

When you push or merge to `main`:

1. **Version Calculation**: The workflow analyzes commit messages since the last tag
2. **Tag Creation**: If a version bump is needed, a new tag is created (e.g., `v1.2.3`)
3. **GoReleaser**: Builds binaries for multiple platforms and creates a GitHub release
4. **Changelog**: Automatically generated from commit messages

### 3. Release Assets

Each release includes:

- **Binaries** for:
  - Linux (amd64, arm64, arm v6, arm v7)
  - macOS (amd64, arm64)
  - Windows (amd64)
- **Checksums** (SHA256)
- **Documentation**: README.md, QUICKSTART.md, config.example.json
- **Changelog**: Auto-generated from commits

## Version Examples

**Starting version: v1.0.0**

| Commit | New Version | Reason |
| --- | --- | --- |
| `fix: resolve timeout issue` | v1.0.1 | Patch release (bug fix) |
| `feat: add retry logic` | v1.1.0 | Minor release (new feature) |
| `feat: redesign config format`<br>`BREAKING CHANGE: Old format no longer supported` | v2.0.0 | Major release (breaking change) |

## Manual Release (Emergency)

If you need to create a release manually:

1. **Create a tag:**

   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```

2. **Run GoReleaser locally:**
   ```bash
   goreleaser release --clean
   ```

## Testing the Release Process

### Local Build (No Release)

Test the build process without creating a release:

```bash
goreleaser build --snapshot --clean
```

This creates binaries in `./dist/` directory.

### Local Release (No Push)

Test the full release process locally:

```bash
goreleaser release --snapshot --clean
```

This simulates a release but doesn't push to GitHub.

## Commit Message Best Practices

### Good Examples

✅ **Feature:**

```
feat(auth): add support for certificate authentication

Implements certificate-based authentication for Azure AD.
This allows using client certificates instead of secrets.

Closes #42
```

✅ **Bug Fix:**

```
fix(client): prevent panic on nil response body

Added nil check before reading response body to prevent
potential panic when server returns empty response.
```

✅ **Breaking Change:**

```
feat(config): redesign configuration file format

BREAKING CHANGE: The configuration file format has changed.
Old format is no longer supported. See migration guide in docs.

Migration guide:
- Rename `clientId` to `client_id`
- Move `auth` settings to root level
```

### Bad Examples

❌ **Too vague:**

```
fix: bug fix
```

❌ **No type:**

```
resolved timeout issue
```

❌ **Mixed concerns:**

```
feat: add retry logic and fix timeout bug and update docs
```

Should be 3 separate commits!

## Release Workflow Details

The `.github/workflows/release.yml` workflow:

1. **Triggers** on push to `main`
2. **Calculates** next version from commits
3. **Creates** Git tag if needed
4. **Builds** with Go 1.25
5. **Releases** using GoReleaser v2
6. **Publishes** GitHub Release with assets

## Troubleshooting

### Release Didn't Trigger

**Check:**

- Did you push to `main` branch?
- Do your commits follow conventional format?
- Check GitHub Actions tab for workflow runs

### Version Not Bumped

**Reasons:**

- Commits were `docs:`, `test:`, `ci:`, or `chore:` (no bump)
- Tag already exists for that version
- No commits since last tag

### Build Failed

**Common issues:**

- Tests failing? Check PR validation first
- Dependencies issue? Run `go mod tidy`
- Platform-specific? Check build matrix in GoReleaser config

## Configuration Files

### `.goreleaser.yml`

Main GoReleaser configuration. Customize:

- **Builds**: Platforms, architectures
- **Archives**: File formats, included files
- **Changelog**: Grouping, filtering
- **Release**: GitHub settings

### `.github/workflows/release.yml`

Release automation workflow. Customize:

- **Triggers**: Branches, tags
- **Version bump rules**: Patch types
- **Permissions**: GitHub token scopes

## Resources

- [Conventional Commits](https://www.conventionalcommits.org/)
- [GoReleaser Documentation](https://goreleaser.com/)
- [Semantic Versioning](https://semver.org/)
- [GitHub Actions](https://docs.github.com/en/actions)
