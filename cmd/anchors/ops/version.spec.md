<!-- @anchors
  code: VRSNI
  updated_at: 2026-11-18
  layer: comando
-->
# VersionCommand — a product's version, its latest release and the next version its branch asks for

> **Code**: `VRSNI`

## Overview

`anchors version` reads the `versioning:` block (VRSNN) and shows, for one product, the version its
sources hold, its latest release, the strategy its branch resolves to and the version that strategy
computes (VRSNE). `--json` prints the same as data, with the commits since the release already
classified: the context a project's own script reads to decide a version. `anchors version bump`
writes a version into every source, after confronting it.

```
cli (semver)
  sources    — (the version is the tag alone)
  released   0.2.0      v0.2.0
  branch     main → increment (default)
  next       0.2.1      increment patch of 0.2.0
  commits    1 since v0.2.0
```

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| `--product` | a product of `versioning:`; none when there is one | a name the block does not hold | this unit (`VRSNI-E01`) |
| `--branch` | any branch name; none: the current branch | a detached HEAD with no `--branch` | this unit (`VRSNI-E02`) |
| the tags | every tag reachable from HEAD; those of the product are read by its scheme | — | this unit |
| `bump`'s argument | none, a part of the scheme, or a version | — | this unit: the version is confronted (`VRSNI-E04`) |

## Effects

| Effect | Description |
| --- | --- |
| `VRSNI-B01` | The command shows the version each source holds, the latest release and its tag, the branch with the strategy it resolves to and the rule that decided (`default`, or the pattern), the next version with its reason, and how many commits came since the release. (`printVersionState`) |
| `VRSNI-B02` | The latest release is the highest version, by the product's scheme, among the product's tags reachable from HEAD: not git's order of names, and never another product's tag. (`latestRelease`) |
| `VRSNI-B03` | `--json` prints the product, scheme, branch, rule, strategy, current and released versions, the sources, the next version with its reason and whether a footer forced it, and each commit since the release with its type, whether it breaks, and its footers. (`describeCommit`) |
| `VRSNI-B04` | `bump` with no argument writes the version the branch's strategy computes into every source, changing only the version's bytes, and names the files. (`bumpVersion`) |
| `VRSNI-B05` | `bump <part>` raises that part of the version the sources hold; `bump <version>` writes that version. |
| `VRSNI-B06` | `bump` to the version the sources already hold says so and writes nothing. |
| `VRSNI-B07` | A project with a single product names none, and the current branch decides when `--branch` is not given. |
| `VRSNI-B08` | A branch with a range (VRSNN-B08) shows it with the rule that set it, and `bump` refuses a version outside it, writing nothing. (`inRange`) |
| `VRSNI-B09` | With the `intents` strategy, the next version reads the product's intents waiting — in the working tree, or at the commit a stage confronts —, and an intent that does not parse stops the command naming the file. (`intentsAt`) |
| `VRSNI-B10` | `promote` takes the version the `from:` branch carries at its head — on `origin` when the checkout has it, else the local branch —, and for a product with no source its latest release there; a branch it cannot find fails naming it. (`promoted`) |
| `VRSNI-B11` | `script` runs its command in the root with the product's state, as `--json` prints it, on its stdin, and takes the first line it prints; a command that fails, or prints nothing, fails naming it. (`scripted`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `VRSNI-E01` | No `versioning:` block, a `--product` the block does not hold, or several products and no `--product`. | Error naming the products there are. | A version shown for a guessed product is another product's. |
| `VRSNI-E02` | A detached HEAD, or no repository, and no `--branch`. | Error asking for `--branch`. | The strategy depends on the branch; a guessed one decides wrongly. |
| `VRSNI-E03` | Sources that disagree or cannot be read (VRSNN-E05), a tag of the product its scheme cannot read, or a strategy that cannot decide (VRSNE-E04). | Error naming the file, the tag or the strategy. | A release skipped or a version guessed would let a wrong version through. |
| `VRSNI-E04` | `bump` to a version the scheme cannot read, below what the sources hold, not above the latest release, on a product with no source, or on a branch that never versions. | Error naming the versions; no file is written. | A version only moves forward, and a release is never written over. |
