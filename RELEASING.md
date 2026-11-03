# Releases

This document describes the versioning and release process of the OpenFGA Terraform provider. This document will be considered a living document. Scheduled releases, supported timelines, and stability guarantees will be updated here as they change.

## Release Requirements
OpenFGA endeavours to rapidly release new features to the community.  To support that goal releases are performed based on the impact and quantity of updates to the codebase in a release-ready state.  Releases may be initiated based on a number of factors.

### Functional Changes
Changes ready for use in production environments which add new functionality or performance improvements may individually trigger an new release.

### Feature Previews
New features which are ready for testing by the community but are marked as experimental may also trigger a new release.

### Maintenance Updates
Minor updates, performance improvements, code refinements, or other items which do not introduce signifigant changes may not trigger a release individually, however when several minor changes have been merged a release may be initiated if no new version has been released in two weeks.

### Security Updates and Bugfixes
Changes made to address a known vulnerability or security issue or fixes to resolve a bug impacting users in production will trigger a release as soon as the changes have been tested and are considered production-ready.  A CVE will also be issued if relevant.

## Versioning
Releases of the OpenFGA Terraform provider will be versioned using dotted triples following the [Semantic Versioning](https://semver.org/) guidelines. For the purposes of this document, we will refer to the respective components of this triple as \<major\>.\<minor\>.\<patch\>. When release candidates are published, they will include a suffix, such as "-alpha" for early alpha builds, "-beta" for beta builds and "-rc" for release candidates, followed by a "." separator and then a number (e.g. `v1.0.0-beta.1`) which identifies them as "pre-releases".

### Major and Minor Releases
Changes in the major and minor version components of the triple indicate changes in the functionality or performance of the OpenFGA Terraform provider which may include breaking changes requiring those integrating the provider in their project to make changes in their use of it when adopting the new version.  It is recommended to fully review the published release notes accompanying a new release before upgrading.

While under v1.0.0, you may expect more frequent breaking changes between minor releases. We will aim to avoid breaking changes between patch releases

### Patch Releases
Releases which increment only the "patch" component of the version triple may introduce new functionality, performance enhancements, bugfixes, or preview features while remaining backwards compatible with the previous release.  

### Changelog
Each release should be combined with a changelog entry in CHANGELOG.md, the CHANGELOG.md file should follow the [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) convention.

## Release Process

Releases are automated with [release-please](https://github.com/googleapis/release-please).

1. Go to **Actions → release-please → Run workflow**. Pick a bump type: `auto` follows conventional commits, `patch`/`minor`/`major` force that bump, and `explicit` uses the exact version entered in `release-version`.
2. release-please opens a PR titled `release: vX.Y.Z` with the version bump and updated `CHANGELOG.md` — review it.
3. Merge it. CI pushes a signed `vX.Y.Z` tag and creates a draft release. The tag triggers the `Release` workflow, which builds and signs the provider with GoReleaser and publishes the draft. Then confirm the version is live on the Terraform and OpenTofu registries.

### Version bumps

Bumps follow [conventional commits](https://www.conventionalcommits.org/), with `bump-minor-pre-major` and `bump-patch-for-minor-pre-major` enabled:

| Commit | Before v1.0.0 | v1.0.0 and later |
|---|---|---|
| `feat:` | patch | minor |
| `fix:`, `perf:`, `refactor:`, `revert:`, `docs:` | patch | patch |
| `type!:` or `BREAKING CHANGE:` footer | minor | major |

### Changelog entries

Built from commit messages since the last release: `feat:` → Added, `fix:` → Fixed, `perf:`/`refactor:` → Changed, `revert:` → Removed, `docs:` → Documentation. `test:`, `ci:` and `chore:` commits are hidden. Hidden commits don't block a release — forced bumps (`patch`/`minor`/`major`/`explicit`) always open a release PR — they just aren't listed in the changelog. Only `auto` opens no PR when every commit since the last release is hidden; use an override to surface one (see below).

#### Overriding an entry

If a commit message isn't user-facing, add this block to the PR description, before or after merging. release-please uses it instead of the commit message:

```
BEGIN_COMMIT_OVERRIDE
feat: user-facing description of the change
END_COMMIT_OVERRIDE
```

#### Breaking changes

Mark the type with `!` and add a `BREAKING CHANGE:` footer. PRs are squash-merged, so check the merge dialog: the squash title must keep the `!` and the body must contain the footer. On single-commit PRs GitHub defaults the squash title to the commit title, not the PR title.

```
feat!: remove the Foo option

BREAKING CHANGE: the Foo option has been removed. Use Bar instead.
```

### Pre-releases

Pre-releases (alpha/beta/rc) are cut from the `develop` branch (create it from `main` if it doesn't exist yet): **Actions → release-please-develop → Run workflow**, supplying an explicit pre-release version (e.g. `0.6.0-beta.1`). Stable releases only go through `main`.

## Release History
Current and previous releases can be found [here](https://github.com/openfga/terraform-provider-openfga/releases).

## Roles

### Release Manager
One of the OpenFGA maintainers is designated as the release manager.  The release manager is tasked with:
- Communicating the release status with other maintainers 
- Reviewing the release PR generated by release-please and ensuring the changelog is meaningful and easy to understand for humans and that contributors to the release are acknowledged
- Discussing the level of testing that is needed and create a test plan if sensible including any areas where test coverage should be expanded
- Running the release workflow and merging the generated release PR (see Release Process above)
- Verifying the release's publication to GitHub, the Terraform registry, and the OpenTofu registry.
- Publishing release announcements or coordinating with other maintainers and community managers to publish release announcements to the openfga.dev website, the #openfga channel in the CNCF slack, and on other channels

In addition to these duties, if the release is a security fix the release manager will coordinate with relevant maintainers and contributors to create and publish a CVE detailing the identified vulnerabilites and their resolution.

## Artifacts

### Binaries

#### MacOS
- terraform-provider-openfga_X.Y.Z_darwin_amd64.tar.gz
- terraform-provider-openfga_X.Y.Z_darwin_arm64.tar.gz

#### Linux
- terraform-provider-openfga_X.Y.Z_linux_386.tar.gz
- terraform-provider-openfga_X.Y.Z_linux_amd64.tar.gz
- terraform-provider-openfga_X.Y.Z_linux_arm.tar.gz
- terraform-provider-openfga_X.Y.Z_linux_arm64.tar.gz

#### Windows
- terraform-provider-openfga_X.Y.Z_windows_386.tar.gz
- terraform-provider-openfga_X.Y.Z_windows_amd64.tar.gz
- terraform-provider-openfga_X.Y.Z_windows_arm.tar.gz
- terraform-provider-openfga_X.Y.Z_windows_arm64.tar.gz

#### FreeBSD
- terraform-provider-openfga_X.Y.Z_freebsd_386.tar.gz
- terraform-provider-openfga_X.Y.Z_freebsd_amd64.tar.gz
- terraform-provider-openfga_X.Y.Z_freebsd_arm.tar.gz
- terraform-provider-openfga_X.Y.Z_freebsd_arm64.tar.gz
