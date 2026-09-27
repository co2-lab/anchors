<!-- @anchors
  code: SPSCS
  updated_at: 2026-09-26
  layer: infra
-->
# SpecSections — the doctor tells when most specs of a layer lack a section a gate needs to see

> **Code**: `SPSCS`

## Overview

The section is Anchors' unit of confrontation: gates look for a NAMED section, and content outside a section
cannot be confronted. A decision can be written, correct and complete, and still no gate reaches it. The
catalogue of sections a new spec is born with only acts at creation: a spec already written does not start
failing because the catalogue changed, and the gates that charge a section only charge the sections the spec
already has. A project born with a poor catalogue has no path to correct it.

This unit closes that hole. It reads every spec in the map once, takes the layer of the unit it describes
and the titles of its sections, and confronts them with six recommended sections: navigation, data contract,
data states, test identifiers, domain and states. Some apply only to screens and components. A section is
reported only when it is the PATTERN, missing in more than half of the specs it applies to, and the report
is ONE line per section with the numbers, never a line per file, so the correction is done in one batch.

The severity follows the corrective action. When the gate that confronts the section is already declared,
the gate is BLIND on those specs, and the finding is a warning asking to edit them. When the gate is not
declared, the finding is informational and asks to adopt the section and declare the gate.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the map | a graph whose spec nodes point at files under the root | a nil graph, which gives nothing | this unit |
| the configuration | the declared gates, by name | a nil configuration, which gives nothing | this unit |
| the spec's layer | the `layer:` line of the spec's header, or the map's layer when there is none | — | this unit: the header wins, so a stale map does not hide the unit's layer |
| the section titles | headings in any supported language, compared without case, plus the project's own spellings | a title outside every translation and extra spelling, which counts as missing | the translation catalogue, read in every language |

## Effects

| Effect | Description |
| --- | --- |
| `SPSCS-B01` | Six sections are recommended: navigation, data contract, data states, test identifiers, domain and states; a screen spec with none of them and no gate declared is told about all six. |
| `SPSCS-B02` | A section is reported only when it is missing in MORE than half of the specs it applies to; exactly half, or a minority, is silent. |
| `SPSCS-B03` | When the gate that confronts the section is declared, the finding is a warning `secao-ausente`: the gate is declared and blind. |
| `SPSCS-B04` | When the section's gate is not declared, the finding is an informational `secao-recomendada`, and its text also asks to declare that gate. |
| `SPSCS-B05` | A section title counts in any supported language and in any case. |
| `SPSCS-B06` | The layer of a spec is the one its header declares, over the map's layer. |
| `SPSCS-B07` | A section scoped to other layers is not demanded: a spec of business logic is not asked for navigation, test identifiers, data contract, data states or states. |
| `SPSCS-B08` | The findings are one per section, never one per spec, each carrying how many specs lack it out of how many it applies to. |
| `SPSCS-B09` | A nil map or configuration gives no finding. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `SPSCS-I01` | Adding the section that was reported to the specs removes the finding: presence and absence are read by the same titles. | a spec lacking navigation is reported, the English title is added, and the finding disappears |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `SPSCS-X01` | The unit never reports a spec file by name: it aggregates per section. | 37 specs times six missing sections is 200 lines nobody reads; one line per section is what lets the team correct in one batch. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `SPSCS-E01` | A spec node's file cannot be read. | The spec is left out of the counts; with no readable spec there is no finding. | An unread spec has no titles to confront, and counting it as lacking every section would invent a pattern; the missing file is the doctor's `no-fantasma` finding. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | config — the declared gates |
| DEP2 | `internal/mapx/model.go` | `Graph`, `KindSpec` | mapa — the specs of the project |
| DEP3 | `internal/i18n/i18n.go` | `T`, `AllTranslations` | apoio — the section titles in every language and the finding texts |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
