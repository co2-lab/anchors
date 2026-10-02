<!-- @anchors
  code: GSRGH
  updated_at: 2026-10-01
  layer: apoio
-->
# GherkinScenarioReader — the scenarios of a unit's feature, with their steps, for the documentation

> **Code**: `GSRGH`

## Overview

The feature is the one artifact of the unit written to be read by people who do not program,
and it answers the question the spec does not: the spec states the rule in the abstract, the
scenario says what happens. Leaving scenarios out of the documentation would keep the system's
observable behaviour in files only the development team opens.

This unit reads the scenarios of the feature linked to a spec and hands them to the templates
with their title, their identity code, their other tags and their body — the steps as written.
The body is what separates this reading from a gate's: a gate confronts a scenario against a
test and needs only the code; documentation needs the steps, because a title alone is a
headline. The feature is found through the map's edge to the spec, never by file-name
convention, which would break in the first project that organised its files differently.

## Effects

| Effect | Description |
| --- | --- |
| `GSRGH-B01` | A line opening with any scenario keyword of any Gherkin dialect, outlines included, followed by a colon, opens a scenario titled by the rest of the line; an examples table heading does not open one. |
| `GSRGH-B02` | In the tag line above a scenario, the first tag shaped like a scenario code becomes the scenario's code, and every other tag goes to its tags without the `@`. |
| `GSRGH-B03` | A scenario's body is its non-blank lines, as written, until the next scenario opens; blank lines are dropped. |
| `GSRGH-B04` | A Feature, Background or Rule heading, in English or Portuguese, closes the current scenario, so its lines never join the scenario before it. |
| `GSRGH-B05` | A spec's features are the feature nodes linked to it by a map edge in either direction, whatever their file names. |
| `GSRGH-B06` | The scenarios of a selection are those of every spec the selection filter picks, and a filter the spec selection refuses is refused here with the same error. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the feature text | Gherkin in any dialect the configuration knows | text that is not Gherkin, which simply yields no scenario | this unit: only recognised title lines open a scenario |
| the spec | a spec the compiler loaded | — | the compiler (`DTCDC`) |
| the map's edges | edges of any type | — | this unit: only edges joining a spec and a feature count |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `GSRGH-I01` | Every scenario read for a spec carries that spec's code, so the documentation and the links know which unit it belongs to. | reads the scenarios of a spec and verifies each one names the spec's code |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `GSRGH-X01` | The reader does not interpret the steps: the body is kept verbatim, indentation included. | The documentation shows the scenario the team wrote; rewriting it would publish a text nobody reviewed. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `GSRGH-E01` | A feature the map links to the spec cannot be read from disk. | That feature contributes no scenario; the other features of the spec are still read, and no error is raised. | A missing feature file is the map's staleness, charged by the map's own gates; failing the whole documentation build here would block every page over one file. |
| `GSRGH-E02` | REF[GSRGH-B06]: a malformed or empty selection filter is refused with the spec selection's error, which fails the template that asked | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/dialect.go` | `GherkinScenarioAlternatives` | config — every scenario keyword of every dialect |
| DEP2 | `internal/doct/doct.go` | `Compiler`, `Spec` | apoio — the loaded specs and the selection filter |
| DEP3 | `internal/mapx/model.go` | `Graph`, `KindFeature`, `KindSpec` | mapa — the edges that link a spec to its features |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
