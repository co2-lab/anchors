<!-- @anchors
  code: OBPCB
  updated_at: 2026-09-26
  layer: infra
-->
# ObligationPack — distributable sets of obligations from a norm, resolved against the project

> **Code**: `OBPCB`

## Overview

The obligations mechanism is agnostic: the engine does not know what a privacy law is, only that there
are duties crossing units. Without packs, every project wrote its duties by hand, in its own vocabulary,
and the regulatory content leaked into the prose of its specs. A pack distributes those duties: a YAML
file with the norm's obligations and the metadata that points at the norm (its domain, its jurisdiction,
its authority, the article each duty comes from). A pack is declarative: it runs nothing.

What belongs to the pack is the duty; what belongs to the project is WHERE the duty lives in it. A pack
names that place with placeholders, and the project resolves them with its own values. A pack whose
placeholders the project has not resolved is refused as a whole, naming the missing values: a duty that
points nowhere would pass green without verifying anything.

A project declares the jurisdictions it operates in, and a pack of another jurisdiction is skipped with a
warning rather than loaded in silence, since declaring it is probably a mistake. A global pack, a pack
without jurisdiction, or a project that declares no jurisdictions loads everything.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the pack references | a short name under `packs/`, or a path ending in `.yaml` or `.yml` | — | the project's `anchors.yaml`; a reference to a missing file is an error |
| the project values | a map from placeholder name to a project path | — | the project's `pack_values:` |
| the jurisdictions | the codes the project declares, in any case and with surrounding spaces | — | the project's `jurisdictions:` |

## Effects

| Effect | Description |
| --- | --- |
| `OBPCB-B01` | A pack file gives its name, domain, jurisdiction, authority, required values and obligations, each obligation with the article that originates it. |
| `OBPCB-B02` | A reference ending in `.yaml` or `.yml` is a path (kept when absolute, joined to the project root when relative); any other reference names `packs/<reference>.yaml`. |
| `OBPCB-B03` | Every `{{ name }}` placeholder in an obligation's places is replaced by the project's value for it, spaces inside the braces allowed; the rest of the text is kept. (`LoadAll`) |
| `OBPCB-B04` | The loaded packs come back sorted by name. |
| `OBPCB-B05` | A pack of a jurisdiction the project does not declare is not loaded, and a warning names the pack and its jurisdiction; jurisdictions are compared ignoring case and surrounding spaces. |
| `OBPCB-B06` | A global pack, a pack without jurisdiction, or any pack of a project that declares no jurisdictions is loaded. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `OBPCB-E01` | A pack has no name. | Refused with "pack without `name`", naming the file. | A pack without a name cannot be reported, sorted or referenced. |
| `OBPCB-E02` | A pack has no obligations. | Refused with "pack without obligations". | An empty pack confronts nothing and would look like conformity. |
| `OBPCB-E03` | A pack file is not valid YAML, or does not exist. | Refused with the parse error naming the file, or with the read error. | Loading part of the packs would report fewer duties than the project declared. |
| `OBPCB-E04` | A required value or a placeholder of a pack is not declared by the project. | The whole load is refused, naming the pack and every missing value, sorted. | A duty that points nowhere passes green without verifying anything. |

## Dependencies

(none)

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
