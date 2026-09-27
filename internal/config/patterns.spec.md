<!-- @anchors
  code: DRPTD
  updated_at: 2026-09-26
  layer: config
-->
# DerivedPatterns — a derived file declared as one pattern or as a list of them

> **Code**: `DRPTD`

## Overview

A layer's derived files tie a spec to the files it governs, one pattern per kind of file.
One pattern covers the ordinary unit: a screen, a handler, a component. It does not cover a
configuration spec, which describes several files spread across the packages (six project
configuration files, a workspace file and four package manifests). With a single pattern the
triad of such a spec never closes, and the gate fails forever on work that is done,
teaching people to waive by habit.

This unit is the value of such a field in the configuration file: it accepts either one
pattern written as text — the usual shape, kept valid because most specs govern one file —
or a list of patterns. An empty list is refused, because a derived kind with no pattern
declares nothing; so is any other shape. When the configuration is written back, one
pattern goes out as text and several as a list, so that rewriting the file never turns the
user's shape into a more verbose one and the diff shows no change where there was none.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the field's value in the configuration file | one text, or a non-empty list of texts | an empty list, a mapping, a list of mappings | this unit: the load refuses them with the cause |
| the patterns themselves | any text | — | the configuration load: whether each pattern is a usable pattern is checked after this unit reads the shape |

## Effects

| Effect | Description |
| --- | --- |
| `DRPTD-B01` | A single pattern written as text becomes a list of one (`UnmarshalYAML`). |
| `DRPTD-B02` | A list of patterns is kept whole and in order. |
| `DRPTD-B03` | An empty list is refused, and the refusal names the cause and the two ways out: declare one, or remove the layer. |
| `DRPTD-B04` | Any shape other than text or a list of texts is refused. |
| `DRPTD-B05` | Written back, one pattern is text and several are a list (`MarshalYAML`). |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DRPTD-I01` | What is written back reads back as the same patterns, in the same order. | writes one pattern and several, reads each back, and compares |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DRPTD-X01` | The patterns are kept exactly as written: this unit neither expands them, trims them, nor checks them as patterns. | This unit reads the SHAPE of the value; whether a pattern is valid is the configuration load's check, and judging it here would fail a file for a reason reported in the wrong place. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DRPTD-E01` | REF[DRPTD-B03]: an empty list is the shape failure B03 refuses, naming the cause | — | — |
| `DRPTD-E02` | REF[DRPTD-B04]: a mapping or a list of mappings is the shape failure B04 refuses | — | — |

## Dependencies

none — `Padroes` is a type config.go uses; it reads nothing from any other file.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
