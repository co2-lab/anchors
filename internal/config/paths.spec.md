<!-- @anchors
  code: PTHSA
  updated_at: 2026-11-06
  layer: config
-->
# Paths — the aliases a path of anchors.yaml may open with, translated in one place

> **Code**: `PTHSA`

## Overview

A path written in anchors.yaml is often long (`apps/mobile/src/screens/...`), and the same prefix repeats
across layers, guides, derived templates and gate exceptions. `aliases:` names it once, as the `paths`
of a tsconfig or the `moduleNameMapper` of Jest do:

```yaml
aliases:
  "@app": apps/mobile/src
  "@guides": docs/engineering/guides
layers:
  ui:
    pattern: "@app/screens/**/*.tsx"
governs:
  - from: "@guides/SCREEN_GUIDE.md"
    governs: tela
```

The fields that hold a path carry a `path:` tag in the configuration's types, and the load translates
every one of them before anything reads it: no command translates an alias on its own. A field named
for a path that carries no tag fails this unit's test, so a new path field cannot be born outside the
aliases.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| an alias's name | `@` followed by letters, digits, `_`, `.` or `-` | any other name | this unit: refuses the configuration (`PTHSA-E01`) |
| an alias's path | a path relative to the root, a trailing `/` ignored | empty, a URL, a path opening with another alias | this unit: refuses the configuration (`PTHSA-E01`) |
| a path of a field | one opening with an alias, or any other path | one opening with `@` that names no alias and no entry of the root | this unit: refuses the configuration (`PTHSA-E01`) |

## Effects

| Effect | Description |
| --- | --- |
| `PTHSA-B01` | A path whose first segment is a declared alias reads the alias's path in its place, in every field that holds one — a file read (`governs[].from`, `gates[].guide`, `changelog.template`, `packs`), a glob (`layers`, `frozen_paths`, `suites[].paths`, the keys of `no_signal`…) or a file written (`derived.files`, `docs.required[].path`, `changelog.path`, the suites' reports); the alias alone reads its path; a first segment opening with `@` that names no alias but an entry of the root is kept as written. (`Aliases.Expand`, `Config.ExpandPaths`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PTHSA-E01` | An alias's name does not open with `@`; its path is empty, a URL or opens with another alias; a path cites an alias `aliases:` does not declare. | The configuration does not load, naming the alias or the field. | A typo in an alias would make a glob match nothing, or a guide read nothing, without a word. |
