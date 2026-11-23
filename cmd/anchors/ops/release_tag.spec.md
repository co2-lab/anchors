<!-- @anchors
  code: RLSTA
  updated_at: 2026-11-23
  layer: comando
-->
# ReleaseTag — the tag and the GitHub Release a stage that releases writes

> **Code**: `RLSTA`

## Overview

A stage marked `release: <product>` (PCSPP) writes the release's record itself when it passes: the
product's tag (`versioning.<product>.tag`) with the version its sha carries, on that sha, and the
tag's GitHub Release, whose body is the release's changelog. A tag the workflow wrote on its own
could name another version, or another commit, than the one the stage passed. Written here, the
record and the tag cannot disagree, and every tag is a version really released, on the exact commit
that went out.

The tag and the Release go through GitHub's API (`gh api`), so the workflow grants
`contents: write`; `pipeline-wired` diverges when it does not (PWSPP-B04).

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the version | the one the product's sources hold at the sha, else the one the branch's strategy computes (VRQRA-B07) | none, sources that disagree | this unit: the release fails (`RLSTA-E01`) |
| GitHub's answers | a tag ref, a release, or a 404 | any other failure | this unit: the release fails (`RLSTA-E01`) |
| the body's template | the template of the product's first output, else `changelog.template`, else the built-in one | a template that cannot be read | CLGCM-E04 |

## Effects

| Effect | Description |
| --- | --- |
| `RLSTA-B01` | When the stage passed, the tag is written on its sha, then the Release of the tag, named after it, whose body is the changelog of the commits since the product's highest release reachable from the sha, the one being written aside. (`releaseOf`, `releaseBody`, `previousRelease`) |
| `RLSTA-B02` | A tag already on the sha is a rerun: nothing is written when its Release exists, and only the Release when it does not. |
| `RLSTA-B03` | On a branch whose label goes into the tag alone, the release is the next prerelease of the version above those GitHub holds, and a rerun releases the one a tag already records on its sha. (`labelledOnRemote`) |
| `RLSTA-B04` | A changelog longer than the body GitHub keeps on a Release, 125000 characters, is cut at the last whole line that fits, and a last line says how many lines were left out: a project's first tagged release, which follows no tag, has the whole history for a changelog. (`fitRelease`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `RLSTA-E01` | The tag is already on another sha; the sha gives no version; GitHub answers neither the ref nor a 404. | The record fails naming the tag and the commit, or the cause, and the stage is recorded as failed. | The same version released from two commits is two releases under one name; a release that could not be written is a stage that did not pass. |
| `RLSTA-E02` | GitHub does not write the tag or the Release, answers a ref that cannot be read or neither a release nor a 404, or `gh` is not in the PATH. | The record fails naming the step — read or write, tag or release — and the stage is recorded as failed. | A record whose tag or Release was not written leaves a release no one can find. |
| `RLSTA-E03` | A tag of the product its scheme cannot read. | The record fails naming the tag. | Skipped, a release would go unseen, and the body would hold the commits of two releases. |
