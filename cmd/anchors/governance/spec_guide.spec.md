<!-- @anchors
  code: SPGDS
  updated_at: 2026-09-26
  layer: comando
-->
# SpecGuide — the project's own spec guide, instantiated with its dialect and a complete example

> **Code**: `SPGDS`

## Overview

This unit renders the `SPEC_GUIDE.md` that `anchors init` seeds into a project: the ruler for writing a spec, placed in the repository and not only behind `anchors guide spec`.

The project needs the file for three reasons. The built-in guide belongs to the framework, while each project has its own dialect — section titles by language, rule letters declared in `rule_types`, code lengths declared in `code_lengths` — and an agent that reads only the generic text does not know what THIS project requires. Anchors also asks every artifact to have a guide in the repository, and not distributing the spec one was the framework demanding what it did not give. And a file in the repository is findable: listed by an agent, opened by a person, linked to its targets by `governs`.

The central difference from the built-in guide is a COMPLETE EXAMPLE, because nobody writes conforming markdown from a description. The guide starts from the command that generates a conforming skeleton, shows the three forms of a catalogued rule, and gives a whole spec to copy — with section titles taken from the same catalogue the generator uses, so the example never contradicts the tool it teaches.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the configuration | a loaded configuration, with or without `rule_types` and `code_lengths` | — | the caller loads it |
| the example code | a code for the example, or empty | — | this unit: an empty code falls back to a default |

## Effects

| Effect | Description |
| --- | --- |
| `SPGDS-B01` | `RenderSpecGuide` renders the guide from the configuration and an example code; an empty example code defaults to `LOGI`. |
| `SPGDS-B02` | The guide shows the three forms of a catalogued rule — heading, table row, bold bullet — and a complete example spec with its `@anchors` header and code. |
| `SPGDS-B03` | The example's section titles are the ones the catalogue gives for the current language, the same catalogue `anchors new spec` uses. |
| `SPGDS-B04` | When the project declares `rule_types`, the guide lists them as a letter and nature table; otherwise it lists the framework's canonical letters. |
| `SPGDS-B05` | The code length sentence ("N or M character(s)") appears only when the project declares `code_lengths`. |
| `SPGDS-B06` | The guide starts from the command that generates the skeleton (`anchors new spec`, and `--list-sections`), before the format. |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `SPGDS-X01` | A project that declares its rule types is not offered the canonical letters. | The canonical list would contradict the project's own, and an agent following it would write codes the `rule-types` gate rejects. |

## Errors

none — the unit renders text from a configuration already loaded and has no failure to handle: an empty example code is normal flow answered by `SPGDS-B01`.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/i18n/i18n.go` | `T`, `TIn` | the section titles and labels of the current language |
| DEP2 | `internal/config/config.go` | `Config` | the project's `rule_types` and `code_lengths` |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
