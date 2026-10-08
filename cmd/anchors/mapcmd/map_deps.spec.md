<!-- @anchors
  code: MDCMP
  layer: comando
-->

# MapDeps — the dependency tree of a file

> **Code**: `MDCMP`

## Overview

`anchors map deps <CODE|file>` prints the dependency tree of a file by the map's
`depends-on` edges — the ones the code's dependency flags declare and the ones the specs'
Dependencies tables declare: down, what the file uses; with `--up`, who uses it. It is the
question the dependency chain exists to answer (DESIGN-dependencies-and-navigation.md).

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the file | a file's own code, a unit's code, or a path of the map | anything else | this unit: refuses naming it |
| the depth | a number of levels, 0 for all | — | this unit |

## Effects

| Effect | Description |
| --- | --- |
| `MDCMP-B01` | The tree is printed one file per line, `CODE path`, indented by level: down the files it uses, or with `--up` the files that use it, each level sorted, `--depth` limiting the levels; a file already on the branch is marked `↺` and not walked again, so a cycle ends. (`DepsTree`) |
| `MDCMP-B02` | The file is named by its own code first, then by its unit's code (its code file), then by its path; a name that is none of them is refused. (`resolveDepsStart`) |
| `MDCMP-B03` | `--kind <kind>` lists each resource of that kind the code reaches, by name, and under each the files that reach it, by their codes. (`ResourceUsers`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MDCMP-E01` | REF[MDCMP-B02]: a name that is no code nor file of the map is refused, naming it | — | — |
