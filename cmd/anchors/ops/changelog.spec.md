<!-- @anchors
  code: CLGCM
  updated_at: 2026-11-18
  layer: comando
-->
# anchors changelog — the technical changelog, printed or written

> **Code**: `CLGCM`

## Overview

`anchors changelog` builds the TECHNICAL changelog from the commits (`CHNGL`): every feature,
fix and breaking change, in the commits' own words, for whoever works on the code. It is not the
product's release notes; `anchors guide changelog` recommends a product changelog an agent
synthesizes from it.

It prints the latest release by default; `--from`, `--all` and `--unreleased` widen the range.
`--write` writes where the project's `changelog:` block says: one incremental file whose top
receives the releases it does not hold, or one file per release. The headings follow the
project's `lang`, and `changelog.template` replaces the built-in template.

With a `versioning:` block (VRSNN), the changelog is a product's: `--product` names it when there
are several. Its releases lie between its own tags, ordered by its scheme, each headed by the version
its tag records. A branch that releases lists its releases; another lists the candidate, what came
after the latest release, headed by the version the branch's strategy computes. `--branch` picks the
branch, by default the first rule that releases, else the current branch. `--write` writes each of
the product's outputs, with its own path, mode and template, or the `changelog:` block when the
product declares none; only on a branch the changelog is written on.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| `--root` | a git work tree, with or without `anchors.yaml` | a directory git does not know | `CHNGL-E02`: the git error |
| `--from` | a tag reachable from `--to` | any other name | this unit: refused naming it |
| `--to` | any ref git resolves | — | git |
| `--product`, `--branch` | a product of `versioning:` (VRSNN-B10); any branch name | — | VRSNN |

## Effects

| Effect | Description |
| --- | --- |
| `CLGCM-B01` | With no range flag, only the latest release is printed. |
| `CLGCM-B02` | `--from` prints every release after that tag, and `--all` every release, newest first. |
| `CLGCM-B03` | `--unreleased` adds what came after the last tag on top; a history with no tag is printed as unreleased without the flag. |
| `CLGCM-B04` | `--write` in incremental mode writes the path the block names (`CHANGELOG.md` by default), titled in the project's language, and a second run reports the file already holds every release and leaves it as it was. |
| `CLGCM-B05` | `--write` in `per_version` mode writes one file per release into the directory; a release's file that exists is not rewritten, the unreleased one is. |
| `CLGCM-B06` | `changelog.template` replaces the built-in template — a file of the project, or one read from a URL through PTRFS —, and the headings of the built-in one follow the project's `lang`. |
| `CLGCM-B07` | A project with no `anchors.yaml` gets the defaults. |
| `CLGCM-B08` | A range with nothing to list prints nothing, and says so on the error stream. |
| `CLGCM-B09` | With `versioning:`, the releases are those between the product's tags, newest by its scheme first, each headed by the version its tag records; another product's tags do not split them, and `--from` takes a version or a tag. (`productReleases`) |
| `CLGCM-B10` | A branch that releases lists its releases; another lists only the candidate, headed by the version its strategy computes. With no `--branch`, the first branch that releases is read. |
| `CLGCM-B11` | `--write` writes every output of the product with its path, mode and template, and refuses on a branch the changelog is not written on. (`renderBlocks`) |
| `CLGCM-B12` | An output whose `source` is `intents` renders what is not released yet from the intents waiting, under the candidate, and a release from the intents archived under `released/<version>/`; a release whose intents were deleted renders nothing, its block being the one the file holds. The commits' outputs read no intent. (`renderIntentBlocks`, `typeHeading`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CLGCM-E01` | `--from` names no tag reachable from `--to`. | The command fails naming the tag. | An unknown start would silently print nothing, or everything. |
| `CLGCM-E02` | `changelog.template` names a file that cannot be read, or a URL that answers nothing (PTRFS-E01). | The command fails naming `changelog.template`. | Falling back to the built-in template would write a changelog the project did not ask for. |
| `CLGCM-E03` | `anchors.yaml` exists and does not load. | The command fails naming the file. | A broken configuration must not be read as the defaults. |
| `CLGCM-E04` | `--from` names no tag of the product, or an output's template cannot be read. | The command fails naming the product, or the output's `template` field. | A changelog of the wrong tags, or with a template the project did not ask for, would look right. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
