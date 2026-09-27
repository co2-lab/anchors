<!-- @anchors
  code: TLVCD
  updated_at: 2026-09-27
  layer: gate
-->
# TestLevelCodes — each scenario references only codes its test level accepts

> **Code**: `TLVCD`

## Overview

A feature scenario is tagged with its test level (`@unit-level`, `@vr-level`, or the project's own
names) and references rules of the spec by code. The project declares, per level, which codes it
accepts, in the `levels` of this gate's own entry under `gates:`: a list of patterns to `allow`, a list to
`exclude`, or neither. This gate confronts every scenario with the filter of each level it is tagged
with, and names the codes the level does not accept.

It is the general form of a naming rule between level and code, and every rule of that kind is
configuration, not a gate of its own. In the reference app, 9 scenarios tagged with its visual level
had no `-VR` code, so the gate that asks for a baseline image never asked, and one `-VR` code sat
under the integration level. With the visual level allowing only `-VR$` and the other levels
excluding it, both show up here.

## Signature

| Parameter | Type | Description |
| --- | --- | --- |
| `content` | `string` | Text of the feature |
| `n` | `mapx.Node` | The feature node |
| `root` | `string` | Project root (not read) |
| `g` | `*mapx.Graph` | The map (not read) |
| `cfg` | `*config.Config` | Project configuration; carries the gate entries and their `levels` |

**Returns**: `(Verdict, string)` — `Pass` when every confronted code is accepted, `Skip` when there
is nothing to confront, `Fail` naming each refused code with its level.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| `n.Kind` | `mapx.KindFeature` | any other kind | this unit: leaves with `Skip` |
| `content` | feature text with each scenario's tags on one line above its title | tags spread over several lines | the feature guide |
| the `levels` of the gate entries | a filter per level tag, with patterns that compile | a pattern that does not compile | the configuration load, which refuses it naming the field |

## Effects

| Effect | Description |
| --- | --- |
| `TLVCD-B01` | A node that is not a feature leaves with `Skip`. |
| `TLVCD-B02` | A project that declares no filter per level leaves with `Skip`: every level accepts every code. |
| `TLVCD-B03` | A feature with no scenario tagged with a level that declares a filter leaves with `Skip`. |
| `TLVCD-B04` | A level that declares `allow` accepts only the codes that match one of its patterns; any other code of its scenarios fails. |
| `TLVCD-B05` | A level that declares `exclude` refuses the codes that match one of its patterns, even when `allow` accepts them. |
| `TLVCD-B06` | Every code of the scenario is confronted, not only the first, and without the `#NN` scenario suffix. |
| `TLVCD-B07` | The failure names each refused code with the level that refused it, sorted. |
| `TLVCD-B08` | When every confronted code is accepted, the gate passes. |
| `TLVCD-B09` | The levels are read from every gate entry that runs this check; two entries declaring the same level add their lists together, and a `levels` on an entry of another check is not read. |

## Errors

none — the gate reads only the text and the configuration it is given; a pattern that does not compile never reaches it, because the load refuses it

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config`, `TestLevel`, `Accepts` | core — the filter per level and what it accepts |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — the verdict messages |
| DEP3 | `internal/mapx/model.go` | `Node`, `KindFeature` | core — the feature node |
