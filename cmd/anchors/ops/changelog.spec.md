<!-- @anchors
  code: CLGCM
  updated_at: 2026-09-27
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

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| `--root` | a git work tree, with or without `anchors.yaml` | a directory git does not know | `CHNGL-E02`: the git error |
| `--from` | a tag reachable from `--to` | any other name | this unit: refused naming it |
| `--to` | any ref git resolves | — | git |

## Effects

| Effect | Description |
| --- | --- |
| `CLGCM-B01` | With no range flag, only the latest release is printed. |
| `CLGCM-B02` | `--from` prints every release after that tag, and `--all` every release, newest first. |
| `CLGCM-B03` | `--unreleased` adds what came after the last tag on top; a history with no tag is printed as unreleased without the flag. |
| `CLGCM-B04` | `--write` in incremental mode writes the path the block names (`CHANGELOG.md` by default), titled in the project's language, and a second run reports the file already holds every release and leaves it as it was. |
| `CLGCM-B05` | `--write` in `per_version` mode writes one file per release into the directory; a release's file that exists is not rewritten, the unreleased one is. |
| `CLGCM-B06` | `changelog.template` replaces the built-in template, and the headings of the built-in one follow the project's `lang`. |
| `CLGCM-B07` | A project with no `anchors.yaml` gets the defaults. |
| `CLGCM-B08` | A range with nothing to list prints nothing, and says so on the error stream. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CLGCM-E01` | `--from` names no tag reachable from `--to`. | The command fails naming the tag. | An unknown start would silently print nothing, or everything. |
| `CLGCM-E02` | `changelog.template` names a file that cannot be read. | The command fails naming `changelog.template`. | Falling back to the built-in template would write a changelog the project did not ask for. |
| `CLGCM-E03` | `anchors.yaml` exists and does not load. | The command fails naming the file. | A broken configuration must not be read as the defaults. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/changelog/changelog.go` | `Git`, `Render`, `Prepend` | infra — the history and the file |
| DEP2 | `internal/config/config.go` | `Load`, `Changelog` | core — the `changelog:` block |
| DEP3 | `internal/i18n/i18n.go` | `T` | infra — the headings' language |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
