<!-- @anchors
  code: SCRQR
  updated_at: 2026-10-19
  layer: gate
-->
# SectionRequired — the spec sections a unit's code calls for exist in its spec

> **Code**: `SCRQR`

## Overview

The gates that read a spec section confront it once it exists. When the code calls for a
section and the spec has none, they have nothing to confront and say nothing. Measured at
MIF, 67 component specs exposed handles with no Test IDs section.

This gate makes that absence a finding. Each layer declares, in the project's own
vocabulary, the sections its code calls for and the trigger that calls for each one
(`layers.<layer>.requires_sections`, DESIGN-section-required.md). The trigger is read from
what the core already measures:
- a handle the code exposes (`dialect.test_handles.exposed`);
- a dependency on a node of a layer;
- a flag the code carries;
- a pattern of the project's that the code matches.

No framework is known here.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the node | a spec | any other kind | this unit: leaves without a verdict |
| the requirements | the `requires_sections` of the layers of the files the spec specifies | an entry with no section, no trigger or two, an unknown trigger, a layer that does not exist, a pattern that does not compile | the config's load, which refuses it naming the entry |
| the spec's sections | its headings, compared ignoring case and surrounding spaces | — | this unit |
| the code | the files the spec `specifies`, as the map holds them | — | the map |

## Effects

| Effect | Description |
| --- | --- |
| `SCRQR-B01` | A node that is not a spec leaves without a verdict; a spec with no code leaves without one too, the missing code being `unit-complete`'s; so does a spec whose code's layers declare no requirement. |
| `SCRQR-B02` | A section a requirement names, triggered by the code and absent from the spec, fails the gate, naming the section, the file and what triggered it; each section is named once, whatever triggered it. |
| `SCRQR-B03` | A section present — its heading at any level, in any case — satisfies every requirement that names it, triggered or not. |
| `SCRQR-B04` | The triggers read: `exposes: test_handles`, a handle the code exposes by the dialect's `exposed` patterns; `depends_on_layer`, a `depends-on` edge of the code to a node of a named layer; `flag`, the flag inside a comment of the code; `code_matches`, the project's pattern over the code with its comments blanked. A requirement whose trigger finds nothing calls for nothing. |
| `SCRQR-B05` | A trigger that cannot be measured — `exposes` with no `dialect.test_handles.exposed`, or with a handle built from values alone (TICTS-B26), a `flag` in a file whose comments are not declared — leaves the gate pending, naming what is missing, when no section is missing. |
| `SCRQR-B06` | `depends_on_layer` counts a dependency on the layer's values; one on its types only — an edge of kind `type`, from `@dep[type]` — counts when the entry's `kinds` names `type`, and the message then says it depends on the types. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `SCRQR-E01` | A specified code file cannot be read. | Its triggers call for nothing. | A file that cannot be read is the map's and `unit-complete`'s to notice; this gate names no section from what it could not read. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none
