<!-- @anchors
  code: VRSTC
  updated_at: 2026-10-03
  layer: gate
-->
# VRStatesCovered — every state of a visual unit is proven by visual regression

> **Code**: `VRSTC`

## Overview

A state is what a screen or a component looks like under a condition, and a visual capture is what
keeps it looking so. `vr-baseline` asked, from the feature, whether a visual-regression scenario that
exists has its baseline image. Nothing asked whether it exists: an agent setting up a project with
screens wrote specs with states, features and unit tests, and no visual regression at all — every gate
green, and no state of any screen protected against a visual change.

This gate confronts a visual unit with its proof. It runs on the unit's main code file — the screen or
the component, which is what a project tags as visual — and reads the spec beside it
(`Button.tsx` → `Button.spec.md`). For every state the spec registers (the codes of the State letter,
`{CODE}-S01`, `{CODE}-S02`…) it asks for three things, all part of the unit:
the visual-regression scenario `{CODE}-VR` in the unit's feature, tagged with the project's visual
regime; a baseline image per state beside the unit, `<Unit>.{CODE}-VR-<state>.<ext>` in any common
image format; and a test that captures them — a capture flow whose path names `{CODE}-VR`, or a
screenshot test beside the unit whose text does.

Which units are visual is the project's to say: the catalog scopes the gate to code layers tagged
`screen` or `component`. A part of the unit (`Button.styles.ts`) has no spec of its own and is left
alone, so each gap is reported once.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the node | a code file | any other kind | this unit: leaves without a verdict |
| the spec | `<Unit>.spec.md` beside the code file, with an `@anchors` header `code:` | none beside it, or one with no code | this unit: skips, saying so |
| the states | codes `{CODE}-<letter><NN>` of the State letter: the project's rule type whose term or section names a state, else `S` | none registered | this unit: skips, saying so |
| the feature | `<Unit>.feature` beside the spec | missing | this unit: counts as no VR scenario |
| the map | the project's graph | none | this unit: no capture is found |

## Effects

| Effect | Description |
| --- | --- |
| `VRSTC-B01` | A node that is not code, a code file with no spec beside it (a part of the unit), a spec with no code, and a spec that registers no state leave without a verdict. |
| `VRSTC-B02` | The feature beside the spec must declare the scenario `{CODE}-VR` with the project's visual-regime tag; otherwise the failure names the code and the tag. |
| `VRSTC-B03` | Every state needs an image beside the unit named `<Unit>.{CODE}-VR-<state>` with an optional variant, as png, jpg, jpeg, webp, gif or svg; the failure counts and names the states without one and the name to save them under. |
| `VRSTC-B04` | A test node captures the VR when its path names `{CODE}-VR`, or when it sits beside the unit (its name starts with the unit's) and its text names `{CODE}-VR`; an image is never the capture. Without one, the failure says no test captures it. |
| `VRSTC-B05` | With the scenario, every state's image and a capture, the gate passes; otherwise it fails with every gap at once. |
| `VRSTC-B06` | The State letter is the one of the project's rule type whose term or a section starts with "state" or "estado"; without one, `S`. |
| `VRSTC-B07` | `vr-baseline` accepts the same image formats for a VR scenario's baseline. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `VRSTC-I01` | A state's image only counts for that state: `S01`'s image never covers `S02`. | a spec with two states and the image of one fails naming the other |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `VRSTC-E01` | The spec beside the code file cannot be read, or there is none. | No verdict, saying there is no spec beside the unit's main file. | A part of the unit has no spec; failing it would charge every part of a screen. |
| `VRSTC-E02` | The unit's feature cannot be read, or there is none. | Counted as no VR scenario, and named among the gaps. | A missing feature is the same gap as a feature without the scenario: nothing declares the capture. |
| `VRSTC-E03` | Looking for a state's image fails on the file system. | Counted as no image for that state. | Unseen is not found; passing it would approve a state nobody captured. |
| `VRSTC-E04` | A test file beside the unit cannot be read. | It is not counted as the capture. | A capture that cannot be read cannot be shown to name the VR. |

