<!-- @anchors
  code: VRQRA
  updated_at: 2026-11-19
  layer: comando
-->
# VersionRequirement — what a stage asks of the version its sha carries, and of the fixes it contains

> **Code**: `VRQRA`

## Overview

The commit gate (`version-forward`, VRFRA) cannot see what was released on another branch, nor which
file sends each version where. A stage can, and it confronts the commit it builds: the requirements
`version:` and `contains:` (PCSPP) are answered at the sha, never in the working tree.

- `fresh`: the version is above every released one. A version is released when a tag of the product
  records it, wherever the tag is reachable from, or when the product's `released:` lists it at the
  sha. A version still in review is not released, so resubmitting it passes.
- `listed`: the listing (VRSNN-B09), read at the sha, has an entry for the version. A review build
  whose version no file lists falls to whatever default the project's code has.
- `same_as_passed: <stage>`: the version is the one the commit that passed that stage carried.
  Production releases what passed review, not a version bumped after it.
- `expected`: the version is the one the branch's strategy computes (VRSNE). `manual` asks it be
  `fresh`, and `none` refuses: that branch never versions.
- The branch's `range:` (VRSNN-B08) bounds the version whatever is asked.
- `contains: <branch>`: the head of that branch on the remote is an ancestor of the sha. A branch
  that goes on without the fixes made on production releases the bug again.

`anchors pipeline` (PCSPA) answers them on the stage's sha; `anchors branch check` (BCSBR) on a
pull request's head, for the stages of its base.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the sha | a commit the repository has; the commit that passed a stage is fetched from `origin` when the checkout lacks it | — | the CI's checkout |
| the tags | every tag of the repository, read by the product's tag and scheme | a shallow checkout with no tags | the CI (`fetch-depth: 0`) |
| the branch | the branch the stage runs on; a pull request's base | — | the caller |
| `contains:`'s branch | a branch of `origin`, fetched before it is read | a branch the remote does not have | this unit: the requirement fails (`VRQRA-E01`) |
| the product's sources and listings at the sha | what VRSNN reads | a version or a listing not found | this unit: the requirement fails (`VRQRA-E01`) |

## Effects

| Effect | Description |
| --- | --- |
| `VRQRA-B01` | `fresh` holds when the version is above the highest release, among the product's tags and the versions `released:` lists at the sha, and says which release; it fails for a version released — naming the tag or the listing — and for one below, asking for a greater version. (`isFresh`) |
| `VRQRA-B02` | `listed` holds when the listing at the sha — its own file, or `released:`'s file and map — has the version with every field of `where`; it fails naming the version and the listing. |
| `VRQRA-B03` | `same_as_passed` holds when the version at the sha is the one the sources held at the commit that passed the stage, and fails naming both versions and that commit; when no commit passed the stage it fails with the reason. (`versionAt`) |
| `VRQRA-B04` | `expected` holds when the version is the one the branch's strategy computes, and fails naming the branch, the strategy, the version asked with its reason and the one the source holds; `manual` holds when `fresh` does, `none` fails saying the branch never versions. (`versionHolds`) |
| `VRQRA-B05` | A branch with a range fails any `version:` requirement for a version outside it, naming the bound and the rule; `from_base` reads the latest release reachable from the sha. |
| `VRQRA-B06` | `contains` fetches the branch from `origin` and holds when its head is an ancestor of the sha; it fails naming how many of its commits the sha lacks and the first five. (`containsHolds`) |
| `VRQRA-B07` | The version, the listings and the release reachable are read at the sha: a working tree that holds another version changes nothing. (`readAt`, `stateAt`) |
| `VRQRA-B08` | On a branch whose label goes into the sources, `expected` holds for a prerelease of the version the strategy computes. |
| `VRQRA-B09` | A key of `released:` the scheme cannot read stops `fresh` naming it and pointing to the listing's `skip:`; skipped, it is no release, and the listing still finds the version. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `VRQRA-E01` | `origin` has no such branch, the sources at the sha hold no version or disagree, a listing or a tag cannot be read. | The requirement fails naming the branch, the file or the tag. | A requirement no one could answer would let the stage run unconfronted. |
