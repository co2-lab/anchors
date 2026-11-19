<!-- @anchors
  code: VAPLA
  updated_at: 2026-11-19
  layer: comando
-->
# VersionApply — the version the branch asks for, brought into the repository

> **Code**: `VAPLA`

## Overview

A computed version is a version no one wrote yet. `anchors version apply` brings it into the
repository the way the product's `apply:` says (VRSNN-B16):

- `release_pr`, the default. The command commits the release on `anchors/release/<product>`,
  pushes it, and opens or updates its pull request into the branch. The commit carries the version
  in every source, the release's block in every changelog output and the intents it consumed
  (INTNA). Merging it is the release. Nobody writes the version by hand, and the CI never pushes to
  a protected branch.
- `commit`: the sources already hold the version, bumped by hand with `anchors version bump`, and
  the command confirms it.
- `build`: the version is written into the sources of the workspace, nothing is committed, and the
  tag alone records it. The command prints `version=` and `tag=`, and appends them to
  `$GITHUB_OUTPUT`.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the product and the branch | as `anchors version` reads them (VRSNI) | — | VRSNI |
| the remote | `origin`, which takes a push of `anchors/release/<product>` | — | the CI's checkout and token |
| the pull request | what the `gh` CLI answers | no `gh` | this unit (`VAPLA-E01`) |
| `$GITHUB_OUTPUT` | a file path, or unset | — | the CI |

## Effects

| Effect | Description |
| --- | --- |
| `VAPLA-B01` | `release_pr` commits `chore(release): <product> <version>` on `anchors/release/<product>`, cut from the branch. The commit holds the version in every source, the release's block in every output, rendered from the commits since the latest release or from the product's intents, and the intents consumed. The command pushes the commit, opens its pull request into the branch, titled `release <product> <version>` with the first block as its body, and leaves the checkout on the branch. (`applyReleasePR`, `writeRelease`, `pendingIntentBlock`, `consumeIntents`) |
| `VAPLA-B02` | A rerun rewrites the release over the branch's new head and updates the open pull request. Once the branch carries the release, there is nothing to propose: nothing is pushed or opened. |
| `VAPLA-B03` | With nothing to release since the latest release, nothing is opened. |
| `VAPLA-B04` | `commit` passes when the sources hold the version the branch asks for, and fails naming both and pointing to `anchors version bump`. |
| `VAPLA-B05` | `build` writes the version into the workspace's sources (the labelled one when the label goes into the sources) and commits nothing. It prints `version=` and `tag=` (the labelled tag when the label goes into the tag) and appends them to `$GITHUB_OUTPUT` when set. (`applyBuild`) |
| `VAPLA-B06` | A `per_version` output gets the release's file, named by `format.file`, in place of `unreleased.md`. Intents consumed by `delete` are removed, and by `archive` moved under `released/<version>/`. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `VAPLA-E01` | `release_pr` of a product with no source, a branch that never versions, or `gh` that cannot be run. | The command fails naming why, and the checkout is back on the branch. | A release pull request with nothing to bump, or one no one could open, would leave a branch behind and no release. |
