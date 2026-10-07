<!-- @anchors
  code: VRCPT
  updated_at: 2026-10-07
  layer: mapx
-->
# Captures — a visual-regression test is tied to the unit it captures

> **Code**: `VRCPT`

## Overview

A capture flow reaches a screen through the running app, not through an import, so nothing in the map
linked the two: the screen changed and the evidence of its capture stayed fresh, green over a look the
screen no longer had. A baseline image swapped without a new capture went unseen the same way.

This unit adds a `captures` edge from each visual-regression test to what it captures — the main code
file of the unit whose `{CODE}-VR` it names, and that unit's baseline images — so the capture's
evidence goes stale when either changes, and `anchors test` knows to run it again.

It is ONE level by decision. A component a screen uses has states and a capture of its own; when the
component changes, its own capture goes stale, not every screen that uses it. The evidence closure
takes a `captures` target and does not descend past it.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the files | the scanned files of the project | — | the map build |
| the unit | a spec `<stem>.spec.md` with a declared header `code:`, and a code file `<stem>.<ext>` | a unit with no declared code, or no code file | this unit: no edge to a code file that is not there |
| a VR test | a test file, not an image, whose path names `{CODE}-VR` or whose text carries a `{CODE}-VR…` code | a test naming no VR code | this unit: no edge |
| a baseline | an image whose path starts `<stem>.{CODE}-VR` | an image named otherwise | this unit: not a baseline of the unit |

## Effects

| Effect | Description |
| --- | --- |
| `VRCPT-B01` | A test whose path or text names `{CODE}-VR` (with or without a state, `-S01`) gets a `captures` edge to the main code file of the unit that declares `{CODE}`, and one to each baseline image of that unit. |
| `VRCPT-B02` | A test naming no VR code, an image, and a VR code no spec declares get no `captures` edge. |
| `VRCPT-B03` | The evidence closure of a VR test holds the unit's code file and images, and does not descend past them: a component the screen depends on is not in it. |
| `VRCPT-B06` | A capture's closure also holds what its unit depends on that no test captures — the `depends-on` of the spec that specifies the unit's code and of the code file itself (its header's `dep:`), transitively through the files reached — and stops at any file a test captures: a hook or a store is in it, a component with its own capture is not. |
| `VRCPT-B07` | A spec's Parts Used names become `composes` edges to the code file of that name (its file name without extension); a name no code file carries ties nothing. They carry no change down. |
| `VRCPT-B08` | `CapturesReaching` are the capture tests whose closure holds one of the given files: a hook a captured screen depends on reaches the screen's captures; a component with its own capture reaches only its own. |
| `VRCPT-B09` | A dependency the code's flags declare — the dependency flag on each import line, one file to the next — reaches a capture's closure and the impact of a change like one a spec declares, transitively: a change to a store a hook imports stales the capture of the screen that imports the hook, and the impact of that change climbs to the hook and the screen. (`EvidenceStaleFor`, `AnalyzeImpact`) |
| `VRCPT-B05` | A contract test — its path or text names `{CODE}-CT` — gets a `captures` edge to the API unit's code file, to its spec (the OpenAPI is compiled from it) and to every OpenAPI document of the project (a file named `*openapi*.yaml`, `.yml` or `.json`), one level like a capture. |
| `VRCPT-B04` | A visual-regression code is read whole with its state — `BUTTN-VR-S01` — by the scan and by the test-signal reader, as the gates read it. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `VRCPT-I01` | A change to the captured unit's code file or to one of its images stales the capture's evidence; a change to a component it uses does not. | a capture ingested, then the screen changed, then only the component changed: stale, then fresh |
