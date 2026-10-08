<!-- @anchors
  code: RCDEO
  updated_at: 2026-10-08
  layer: comando
-->
# Recode — renames an identity code and carries the change to every textual surface of the project

> **Code**: `RCDEO`

## Overview

A unit's identity code is meant to be stable, and stable must still mean reversible: renaming a code by
hand leaves residue — a header changed here, a scenario code forgotten there — and the reference
application still carries a half-finished rename that proves it. This command is the rename done in one
step: it swaps the old code for the new one in the identity header of every file, in the scenario codes
derived from it, and in bare mentions of the code in other units' cross-references.

It is a dry run by default. It shows the plan — which files change and how many occurrences of each kind
— and writes nothing. With the apply switch it writes the files and then rebuilds the map from the
headers, which are the source of truth, instead of editing the map's text — keeping, as the map build
does, the stamps, judgments and flow graph the headers do not carry. Writing is done in bulk and
stops at the first file it cannot write; when that happens after some files were already rewritten, the
command says the project is half converted, because a user who does not know it would run it again over
a project that is partly in the new form.

The plan itself — which occurrences are found and how they are classified — belongs to the recode
engine; this unit is the command that shows it, applies it and puts the map back in step.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the arguments | exactly two codes, the old and the new, in any letter case | any other number of arguments | this unit: refuses any other count |
| the old code | a code that appears in some file of the project | a code no file carries | the recode engine: refuses to plan it |
| the project | a root with a configuration file | a root with no configuration | this unit: fails naming the configuration file |

## Effects

| Effect | Description |
| --- | --- |
| `RCDEO-B01` | Both codes are upper-cased before planning, and the report names the rename in that form with the number of files and content substitutions. |
| `RCDEO-B02` | The plan lists each file with its occurrences counted by kind, in the fixed order header, scenario code, bare reference, omitting a kind with none. |
| `RCDEO-B03` | Without the apply switch nothing is written, and the report ends saying it was a dry run. |
| `RCDEO-B04` | With the apply switch every surface is rewritten — the header, the scenario codes of the spec and of the test — and the report says how many files were rewritten. |
| `RCDEO-B05` | After applying, the map is rebuilt from the rewritten headers, so its node carries the new code, and the report gives the node count. |
| `RCDEO-B06` | The rebuild after applying keeps what the previous map knew beyond the headers, as the map build does: the stamps and judgments of every edge that survives, and the flow graph; a judgment lost with an edge that did not survive is reported. |
| `RCDEO-B07` | When the project declares a recode dialect, the report also counts the testIDs and the files to rename and lists each rename from its old name to its new one; without them, no testID or rename line appears. |
| `RCDEO-B08` | When the dialect's testID prefix appears nowhere but the files that carry the code hold testIDs with another prefix, the report warns about the divergent prefix; with nothing to warn about, no warning appears. |
| `RCDEO-B09` | A write failure before any file changed fails the apply without saying the project is half converted. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RCDEO-I01` | After an apply, the old code is found neither in the rewritten spec nor in the rebuilt map: the files and the map describe the same identity. | applies a rename and reads the spec and the map afterwards |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RCDEO-X01` | Does not edit the map's text; it rebuilds the map from the files. | Editing the map by pattern would perpetuate any divergence between the headers and the map. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `RCDEO-E01` | The project has no configuration file. | The command fails naming the configuration file, before planning. | Without the layers the project cannot be walked. |
| `RCDEO-E02` | The old code appears in no file of the project. | The command fails and nothing is written. | There is nothing to rename, and a silent success would hide a typo. |
| `RCDEO-E03` | A file cannot be written during the apply, after other files were already rewritten. | The command fails and first reports how many files had already changed, saying the project is half converted and the map still describes the old form. | The disk is left in an intermediate state, and running again without knowing it would work over a mixed project. |
| `RCDEO-E04` | The command receives other than two codes. | It is refused before reading anything. | A rename needs exactly a source and a target. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
