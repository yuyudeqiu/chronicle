# Release process

Chronicle uses Semantic Versioning while the project is still in the `v0.x`
development phase. Releases are prepared by Release Please and built by
GoReleaser.

## Version policy

| Change | Commit prefix | Version change before v1 |
| --- | --- | --- |
| Bug fix | `fix:` | Patch: `v0.1.0` → `v0.1.1` |
| Backward-compatible feature | `feat:` | Minor: `v0.1.0` → `v0.2.0` |
| Breaking change | `feat!:` or `BREAKING CHANGE:` | Minor while `< v1.0.0` |
| Documentation, tests, or CI | `docs:`, `test:`, `ci:`, `chore:` | No release by itself |

Pull request titles should follow Conventional Commits because squash merges use
the PR title as the commit subject. Examples:

```text
fix: preserve task title when updating progress
feat: add task labels
feat!: change the exported task schema
```

## Publishing a release

1. Merge normal pull requests into `main`. CI continues to lint, test, and build
   them as usual.
2. The `Release` workflow creates or updates one Release PR containing the next
   version and `CHANGELOG.md` changes.
3. Review and merge the Release PR when the accumulated changes are ready to
   publish.
4. The same workflow creates the immutable `vX.Y.Z` tag and GitHub Release, then
   GoReleaser uploads macOS, Linux, and Windows archives plus `checksums.txt`.

Do not move or reuse a published version tag. If a release has a problem, fix it
on `main` and publish the next patch version.

## One-time repository setting

In GitHub, open **Settings → Actions → General → Workflow permissions** and
select **Read and write permissions**, then enable **Allow GitHub Actions to
create and approve pull requests**. The workflow uses the repository
`GITHUB_TOKEN`; no personal access token or extra secret is required.

Pull requests created with the repository token do not start another workflow
run. This is normally fine because a Release PR only updates release metadata and
the underlying commits have already passed CI. If branch protection requires a
fresh CI run on the Release PR, configure a GitHub App token or fine-grained PAT
for Release Please instead.

## Development and prerelease installs

Install the latest released version:

```bash
go install github.com/yuyudeqiu/chronicle@latest
```

Install the current `main` branch for testing:

```bash
go install github.com/yuyudeqiu/chronicle@main
```

Only create an alpha tag for a version that needs an explicit testing milestone,
for example `v0.3.0-alpha.1`. Go does not normally select prereleases for
`@latest`, so testers must request that version explicitly.

## Local verification

Before changing the release configuration, run:

```bash
go test ./...
npm --prefix frontend ci
npm --prefix frontend run build
goreleaser check
goreleaser release --snapshot --clean --skip=publish
```
