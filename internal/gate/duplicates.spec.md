<!-- @anchors
  code: GTDPG
  layer: gate
-->
# Duplicates — each gate confronts the repeats of what it declares

> **Code**: `GTDPG`

## Overview

A gate that controls a kind of declaration registers an occurrence reader: the key of each
declaration it controls, with its line. After the gate measures a node, the engine counts the
keys, and a key declared more than once is a finding of that gate, naming the key and every
line (DESIGN-duplicate-declarations.md). Only declarations count — a citation repeats by
nature —, and each kind of declaration has one owner, so a repeat is reported once. On by
default; `duplicates: false` on the gate switches it off.

The first reader is `rule-types`'s: a rule code defined twice in one file; then the spec catalogue's.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the gate | a gate with a check; its `duplicates` key | a gate with no occurrence reader | this unit: nothing is counted |
| the verdict | the gate's verdict on the node | a skipped node | this unit: a skipped node is not counted |
| the file | the node's text | a file that cannot be read | this unit: the gate's verdict stands as it was |

## Effects

| Effect | Description |
| --- | --- |
| `GTDPG-B01` | A key the gate's reader finds more than once in a node turns the verdict into the reader's — a failure, unless the gate measures repeats as a divergence of its own, and never softer than a failure the gate already gave —, naming each repeated key, how many times and the lines, beside the gate's own finding. A node the gate skipped, a gate with no reader, and a gate with `duplicates: false` are left as they were. (`Occurrence`, `HasDuplicateReader`, `confrontDuplicates`) |
| `GTDPG-B02` | `rule-types` counts each rule code a file defines — a heading, the first cell of a table row, a bold or bare bullet —, where it is defined; not inside a section whose rows cite codes (what a rule uses, the open decisions, the navigation, the state flow, the events a unit emits, a change history, what a plan revises) nor its subsections — known by the catalog's titles and by the titles the project declares for them in any layer; not on an alias or a retired line, nor a list item opening with the code in backticks; and a heading with the rows under it that open with its own code once. (`definedRuleOccurrences`) |
| `GTDPG-B03` | The spec catalogue's gates count their own declarations, each by the row that declares it: `env-declared` a variable of the environment table, `dependency-honored` a `DEPn` opening a row, `domain-declared` a Domain entry (case and backticks aside), `open-questions-resolved` an open question's code, answered rows included, `revision-orphans` a revision code opening a line, and `spec-sections` a catalog section — by its catalog title in any language or the title the project gave it — under the same parent heading. (`envOccurrences`, `depOccurrences`, `domainOccurrences`, `openQuestionOccurrences`, `revisionOccurrences`, `sectionOccurrences`) |

## Errors

none — the unit reads text the gate already read; a file it cannot read leaves the verdict as it was (Domain).
