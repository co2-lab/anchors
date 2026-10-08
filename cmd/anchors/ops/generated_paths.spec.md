<!-- @anchors
  code: GNPTG
  updated_at: 2026-10-08
  layer: comando
-->
# GeneratedPaths — the product names the files it derives, so a conflict in them is rebuilt, not merged

> **Code**: `GNPTG`

## Overview

Whoever resolves a merge conflict has to tell two cases apart that look the same in the
working tree: a conflict in a GENERATED file is noise, because rebuilding it produces the
right answer, while a conflict in a work file is a divergence that only its authors can
settle. The list of generated files used to be written by hand in each project's conflict
script, as three literal paths repeated from the product.

`generated-paths` makes the product answer the question itself. It prints the patterns of
the three kinds of file Anchors derives: the map at the project root, the compiled
documentation directory, and the progress companion of every plan. Each pattern comes from
the same constant the product uses to write that file, so a path that changes in the product
changes here too, and no project script has to be hunted down.

The output has two forms: one pattern per line (the default), or all of them joined as a
single alternation ready to be handed to an extended grep. The command only answers inside a
governed project.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | a directory holding a loadable `anchors.yaml` | a directory with no config, or with one that does not load | this unit: the command fails with the load error |
| the output format | `re` for the alternation; any other value prints one pattern per line | — | this unit: only `re` is special-cased |

## Effects

| Effect | Description |
| --- | --- |
| `GNPTG-B01` | The patterns name exactly the map at the root, the compiled documentation directory at the root, and any path ending in the plan progress suffix; nothing else matches them. |
| `GNPTG-B02` | By default the command prints one pattern per line. |
| `GNPTG-B03` | With `--format re` it prints the same patterns joined by `|` on one line, a valid regular expression. |
| `GNPTG-B04` | The dot of a path is escaped, so it matches only a literal dot: `anchorsXgraph.yaml` is not a generated file. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `GNPTG-I01` | The line form and the alternation form carry the same patterns, in the same order. | runs both forms on the same project and compares the lines joined by `|` with the alternation |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `GNPTG-X01` | The command resolves no conflict and rebuilds nothing: it only names the derived files. | Deciding what to do with a conflict is the resolver's job; the product owns only the list. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `GNPTG-E01` | The root has no `anchors.yaml`, or it does not load. | The command fails with the load error and prints no pattern. | "What is generated?" only has an answer in a governed project; a guessed list in a foreign tree would mislead the resolver. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
