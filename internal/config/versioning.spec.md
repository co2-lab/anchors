<!-- @anchors
  code: VRSNN
  updated_at: 2026-11-18
  layer: config
-->
# Versioning — the products a project versions, where each writes its version, and how each branch moves it

> **Code**: `VRSNN`

## Overview

`versioning:` maps a name to a product: one versioned thing with its own release line. A
repository with a single product still names it, so a second product is not a breaking change.
A product says how its versions are ordered (`scheme`, VRSNE), which files hold the version
(`source`), what the release record is called (`tag`), how a branch decides its next version
(`next`), the releases made before the project wrote tags (`released`), and the branches whose
rule differs from that default (`branches`), with the versions each may produce (`range`).

```yaml
versioning:
  app:
    scheme: semver
    source:
      - file: app/manifest.json
        path: version
      - file: app/build.txt
        regex: 'VERSION=(\S+)'
    tag: "app-v{version}"
    released: { file: release-policy.json, keys: versions, where: { target: production } }
    next: { strategy: semantic }
    branches:
      "hotfix/*":
        next: { strategy: increment, part: patch }
        range: from_base
      "feat/*":
        next: none
```

`released:` and a stage's `listed:` (PCSPP) are a *listing*: a file, the map in it whose keys are
versions, and the fields an entry must hold. It interprets nothing.

Nothing here is a default for a kind of project: the files, the paths, the tag and the rules are
the project's. A project with no `versioning:` block sees no change.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| a product | `scheme` (with `format` or `pattern` as the scheme asks), `source`, `tag`, `next`, `released`, `branches` | any other key | the strict load: the configuration does not load |
| `source` | a list of `file` with `path:` (a dot path into JSON or YAML, a number for a list item) or `regex:` (one group) | both, neither, no file, a regex with another number of groups | this unit (`VRSNN-E02`) |
| `tag` | a name with `{version}` once | none, or twice | this unit (`VRSNN-E03`) |
| `next` | a strategy's name, or a map of `strategy`, `part`, `initial`, `types` | any other key, a strategy that does not fit | this unit (`VRSNN-E04`) |
| `branches` | a map of branch pattern (BFSBR-B03) → `next:`, `range:` | any other key, a pattern that is neither a glob nor `re:` | this unit (`VRSNN-E04`) |
| a branch rule's `changelog`, `release` | `true` or `false` | — | the strict load |
| `changelog.outputs` | a list of `path`, `mode` (`incremental`, `per_version`), `template`, `format` (`heading`, `file`, `max_chars`, `required_sections`, `required`) | no path, another mode or key, a heading or file with no `{version}`, `file` off `per_version`, `required` other than `changed`, a negative limit | this unit (`VRSNN-E08`) |
| `range` | `from_base`, or numbers and an `x` (`1.0.x`, `2.x`) | anything else | this unit (`VRSNN-E06`) |
| a listing | `file` (JSON or YAML), `keys` (a dot path to a map; empty: the document), `where` (dot path → text) | `released:` with no file, any other key | this unit (`VRSNN-E06`) |
| a listing's file | a document with a map at `keys` | a file missing, unreadable, or without a map there | this unit (`VRSNN-E07`) |
| a source file | a JSON, YAML or any text file | a file missing, unreadable, or without the version where the source says | this unit (`VRSNN-E05`) |

## Effects

| Effect | Description |
| --- | --- |
| `VRSNN-B01` | `versioning:` loads its products, each with its scheme, sources, tag, default strategy and branch rules; a configuration without it loads as before and versions nothing. |
| `VRSNN-B02` | `next:` is a strategy's name (`next: none`) or a map of it and its options; a product that writes no `next:` resolves to `manual`. (`Resolved`) |
| `VRSNN-B03` | A branch resolves to the rule of the first pattern of `branches:` it matches, in the order written, and to the product's `next:` when none matches; the resolution names the pattern that decided. (`NextFor`) |
| `VRSNN-B04` | The tag is `v{version}` when the product names none; a tag is read back as a version of the product only when it has the tag's text around the version. (`TagFor`, `VersionOfTag`) |
| `VRSNN-B05` | A source's version is read at its path in a JSON or YAML document, or in its regex's group; two or more sources must hold the same version, and a disagreement names each file and its version. (`ReadSources`, `Agreed`) |
| `VRSNN-B06` | A source is rewritten with a new version and nothing else of the file changes: the bytes around the version — quotes, spacing, the other keys — are kept as they are. (`Rewrite`) |
| `VRSNN-B07` | `released:` loads its file, map and fields; a branch takes the rule of the first pattern it matches, and that rule's range, or none; a rule that writes no `next:` takes the product's. (`RangeFor`, `NextFor`) |
| `VRSNN-B08` | A version is inside an explicit range when its first parts are the numbers written before the `x`, and inside `from_base` when its parts but the last are those of the release the branch was cut from; with no such release it is outside. The reason names the bound. (`InRange`) |
| `VRSNN-B09` | A listing is the keys of its map whose entries hold every field of `where`, compared as text; a requirement's listing with no file takes the file and the map of `released:`, keeping its own fields. A listing writes itself as `file keys, where field = value`. (`Listing.Versions`, `ListingOf`, `Listing.String`) |
| `VRSNN-B10` | A product is the one named, or the only one; several and no name is an error that lists them and says where the name is written. (`Pick`) |
| `VRSNN-B11` | A product's changelog outputs load with their format; a branch releases when its rule says `release: true`, and the changelog is written on it when its rule says `changelog: true` — a product whose rules name no such branch releases, and writes, on every branch —; the branch a changelog reads by default is the first rule that releases. A heading pattern reads the version and the date of a line. (`Releases`, `WritesChangelog`, `ReleaseBranch`, `HeadingRE`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `VRSNN-E01` | A product whose scheme does not build (VRSNE-E01), or that is empty. | The configuration does not load, naming `versioning.<product>` and the scheme. | A scheme that reads nothing orders nothing. |
| `VRSNN-E02` | A source with no `file`, with both or neither of `path:` and `regex:`, or with a regex that does not compile or does not have exactly one group. | The configuration does not load, naming the source. | A source that names no place holds no version. |
| `VRSNN-E03` | A tag without `{version}`, or with it twice. | The configuration does not load, naming the tag. | A tag with no version in it records every release under one name. |
| `VRSNN-E04` | A strategy that does not exist, a part the scheme lacks, `semantic` off semver, a level of `types:` that is not one, `manual` on a product with no source, an option of another strategy, a key `next:` or a branch rule does not know, a branch pattern that is neither a glob nor `re:`. | The configuration does not load, naming the field. | A rule that cannot decide would leave the version to a guess. |
| `VRSNN-E05` | A source file missing or unreadable, not JSON nor YAML for a `path:`, without a scalar at the path, matched by no regex, or holding its version in a form that cannot be rewritten in place. | An error naming the file and the source. | A version not found is not "no version": read as absent, any version would pass. |
| `VRSNN-E06` | A range that is neither `from_base` nor numbers and an `x`, or `released:` with no file or with a key it does not know. | The configuration does not load, naming the branch or `released`. | A bound that reads nothing would let every version through. |
| `VRSNN-E07` | A listing's file missing or unreadable, not JSON nor YAML, or without a map at `keys`. | An error naming the file and the path. | A listing not found is not an empty one: read as empty, every version would be fresh and none listed. |
| `VRSNN-E08` | An output with no path, a mode other than `incremental` or `per_version`, a key it does not know, a heading or file with no `{version}` or a heading that does not compile, a `file` on an incremental output, a `required` other than `changed`, or a negative `max_chars`. | The configuration does not load, naming the output and the field. | A format that reads nothing would hold the changelog to nothing. |
