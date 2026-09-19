<!-- @anchors
  code: VRBSV
  updated_at: 2026-09-19
  layer: gate
-->
# VRBaseline — ensures visual regression scenarios have captured reference baseline images

> **Code**: `VRBSV`

## Overview

Verifies that visual regression scenarios declared in feature files have corresponding reference baseline capture images on disk. Visual regression represents a distinct proof surface where verification occurs through visual screen capture rather than unit test assertions. Without a baseline image, a visual scenario exists in the feature and is counted as covered, yet no comparison image exists against which changes can be evaluated, leaving promised visual proofs unverified. In the originating repository audit, out of 105 baselines on disk and 4 features declaring visual regression scenarios, 3 scenarios were declared without any baseline image.

## Signature

| Parameter | Type | Description |
| --- | --- | --- |
| `content` | `string` | Text content of the feature artifact |
| `n` | `mapx.Node` | Node representation of the feature artifact |
| `root` | `string` | Repository root directory path where baseline images are searched |
| `g` | `*mapx.Graph` | Graph model representing project artifact relationships |
| `cfg` | `*config.Config` | Project configuration specifying derived visual regime tags |

**Returns**: `(Verdict, string)` returning `Pass` with empty detail when all visual scenarios have baseline images, `Skip` when the artifact is not a feature or lacks visual scenarios, or `Fail` naming missing baseline scenario codes.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| `n` | any `mapx.Node` | nil node | the gate runner, which passes the current graph node |
| `n.Kind` | `mapx.KindFeature` | non-feature nodes | this unit: non-feature nodes leave with `Skip` |
| `content` | valid feature text | binary content | the gate runner, reading UTF-8 feature files |
| `root` | valid directory path | inaccessible directory | the caller, providing the repository root |
| `cfg` | project `config.Config` or nil | unparsed configuration | this unit, falling back to default visual regime tag `vr-level` |

## Effects

| Effect | Description |
| --- | --- |
| `VRBSV-B01` | Artifacts that are not feature files leave with verdict `Skip`. |
| `VRBSV-B02` | Feature files declaring no visual regression scenarios leave with verdict `Skip`. |
| `VRBSV-B03` | Visual regression scenarios having matching baseline images on disk pass with verdict `Pass`. |
| `VRBSV-B04` | Visual regression scenarios lacking matching baseline images fail with verdict `Fail` naming the missing scenario codes. |
| `VRBSV-B05` | Baseline image matching supports naming variants where variant suffixes are appended after the scenario identifier. |
| `VRBSV-B06` | Reads visual regime tag from project configuration `derived.regimes` mapping tag keys to regime names. |
| `VRBSV-B07` | Falls back to default tag `vr-level` when project configuration is nil or defines no visual regime. |
| `VRBSV-B08` | Multiple missing baseline scenario codes are sorted alphabetically in the failure diagnostic message. |
| `VRBSV-B09` | Only scenario codes containing the visual regression suffix `-VR` are matched for baseline existence. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `VRBSV-I01` | Every visual regression scenario declared in a feature must correspond to at least one baseline image on disk. | evaluates features with visual scenarios against present and missing image files and verifies exact verdict mapping |
| `VRBSV-I02` | The visual regime tag mapping treats the map key as the tag in the feature and the map value as the regime name. | tests configuration parsing with custom regime mappings verifying correct tag resolution |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `VRBSV-X01` | Does not evaluate baseline staleness using disk modification timestamps. | Disk mtime changes on repository cloning, which accused 105 of 105 baselines in audit without providing true staleness signal. |
| `VRBSV-X02` | Does not fail commits based on git commit dates of baseline images. | Accusing 103 of 105 baselines (98% of the repository) would cause the gate to be disabled immediately, neutralizing functioning protections. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | core — project configuration containing derived visual regime mappings |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized messages for skips and failure verdicts |
| DEP3 | `internal/mapx/model.go` | `Graph`, `KindFeature`, `Node` | core — graph model and feature node representations |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
