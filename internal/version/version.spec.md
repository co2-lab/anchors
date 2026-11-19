<!-- @anchors
  code: VRSNE
  updated_at: 2026-11-18
  layer: infra
-->
# Version — a product's version read, ordered and moved, and the next one decided

> **Code**: `VRSNE`

## Overview

A product's version is a string its scheme reads. The scheme orders two versions and raises a
part of one; nothing here knows a language, a package manager or a store. Four schemes cover what
projects write: `semver` (`1.2.0`, `2.0.0-rc.1`), `calver` with a declared format
(`2026.11.2` under `YYYY.0M.MICRO`), `number` (`412`) and `pattern`, a regular expression whose
capture groups are the version's numbers (`r14-3` under `^r(\d+)-(\d+)$`).

The next version is decided by the strategy a branch resolves to:

| Strategy | The next version |
| --- | --- |
| `manual` | what the sources hold; a person or an agent writes it |
| `increment` | the latest release with one part raised and the parts after it reset |
| `semantic` | the latest release raised by the commits since it: a breaking change the major, a `feat` the minor, a `fix` the patch |
| `none` | no version: the branch never versions |

A `Version:` footer on a commit since the latest release forces the result, whatever the
strategy. A version the scheme cannot read is an error, never "smaller" nor "equal": a guess
there would let anything through.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the scheme | `semver`, `calver` with a format of `YYYY`, `YY`, `0M`, `MM`, `0D`, `DD` and a last `MICRO`, `number`, `pattern` with at least one capture group | anything else | this unit: `NewScheme` refuses it (`VRSNE-E01`) |
| a version | any string | — | this unit: one the scheme cannot read is an error (`VRSNE-E02`) |
| the commits | any subject and body, oldest first | — | the changelog's classification (CHNGL) |
| the strategy | `manual`, `increment` with a part the scheme has, `semantic`, `none` | anything else | the configuration refuses it at load; this unit refuses it too (`VRSNE-E04`) |

## Effects

### Schemes

| Effect | Description |
| --- | --- |
| `VRSNE-B01` | Two semver versions are ordered as SemVer 2.0 says: major, minor and patch as numbers, a pre-release below its release, its identifiers numeric below alphanumeric, and build metadata ignored. (`Compare`) |
| `VRSNE-B02` | A calver version is read field by field as its format says, each field as a number in the width the format gives it (`0M` two digits, `MM` with no leading zero), and ordered left to right. |
| `VRSNE-B03` | A `number` version is a non-negative integer, ordered as one. |
| `VRSNE-B04` | A `pattern` version is ordered by its capture groups, as numbers, left to right. |

### Increment

| Effect | Description |
| --- | --- |
| `VRSNE-B05` | Semver raises `major`, `minor` or `patch` and resets the parts after it; a pre-release is already the release it precedes, so `1.1.0-rc.1` raised in its minor or its patch is `1.1.0`. (`Increment`) |
| `VRSNE-B06` | Calver writes the date fields from today, and its `MICRO` counter goes up while the date is the same and starts at 0 on a new date. |
| `VRSNE-B07` | A pattern raises the group its index or its name points to, resets the later groups to 0, and keeps the text around the groups. |

### The next version

| Effect | Description |
| --- | --- |
| `VRSNE-B08` | `manual` is the version the sources hold. (`Next`) |
| `VRSNE-B09` | `increment` raises the part of the latest release; with no commit since the release, the version stays the release. |
| `VRSNE-B10` | `semantic` raises the latest release by the largest change of the commits since it — a breaking change the major, a `feat` the minor, a `fix` the patch —, and says what asked for it (`2 feat, 1 fix`); with nothing that bumps, the version stays the release. A commit's type and its `!` are read from its subject. (`CommitType`) |
| `VRSNE-B11` | `semantic` with `initial: true` bumps the minor for a breaking change below 1.0.0, and `types:` gives a commit type its level (`perf: patch`, `feat: none`) among `major`, `minor`, `patch` and `none`. (`ValidLevel`) |
| `VRSNE-B12` | Before any release, `increment` and `semantic` start from what the sources hold, or, with no source, from the scheme's first version: `1.0.0` (`0.1.0` with `initial: true`) or `1`. |
| `VRSNE-B13` | A `Version:` footer on a commit since the latest release forces the next version, the newest one winning, and the result names the commit. |
| `VRSNE-B14` | `none` gives no version and says the branch never versions; a footer does not change it. |
| `VRSNE-B15` | `intents` raises the latest release by the largest bump among the intents waiting — the first of the scheme's parts any of them names —, saying how many intents moved it; with none, or only `none`, there is nothing to release; before the first release, it takes what the sources hold. (`largestBump`) |
| `VRSNE-B16` | `promote` and `script` take the version decided outside them — the one a branch carries, the one a command printed —, the reason saying by what; with none decided they fail with why, and the version decided is still confronted with the latest release. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `VRSNE-E01` | A scheme that is not one of the four, a calver format with an unknown field or a `MICRO` before the end, a pattern that does not compile or has no group. | `NewScheme` fails, naming the scheme and what is wrong. | A scheme that reads nothing would order nothing. |
| `VRSNE-E02` | A version the scheme cannot read: the sources', the release's, the computed one or a forced one. | An error naming the version and the scheme. | Read as smaller or equal, a broken version would pass every check. |
| `VRSNE-E03` | The next version — forced or computed — is below the latest release. | An error naming both. | A version that goes back would publish over a release. |
| `VRSNE-E04` | A part the scheme does not have; `semantic` on a scheme that is not semver; `manual` with no source; `increment` or `semantic` with no release, no source and no first version for the scheme; a calver with no counter raised on the date it already holds; an unknown strategy. | An error naming the strategy and what it lacks. | A strategy that guessed would publish a version nobody chose. |
