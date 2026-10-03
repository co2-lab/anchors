<!-- @anchors
  code: VRSTC
  updated_at: 2026-10-03
  layer: gate
-->
# VRStatesCovered — each state of a visual unit tied to its visual regression, both ways

> **Code**: `VRSTC`

## Overview

A state is what a screen or a component looks like under a condition, and a visual capture is what
keeps it looking so. An agent setting up a project with screens wrote specs with states, features and
unit tests, and no visual regression at all — every gate green, and no state of any screen protected
against a visual change. `vr-baseline` only asked a VR scenario that exists for its image.

Four gates ask the four questions that tie a state to its capture, both ways:

| Gate | Question |
| --- | --- |
| `vr-states-covered` | Does every state of the spec have a VR scenario in the feature? |
| `vr-scenarios-tested` | Does every VR scenario have a VR test, and a baseline image? |
| `vr-scenarios-of-states` | Is every VR scenario of a state the spec registers? |
| `vr-tests-of-scenarios` | Is every VR test of a VR scenario the feature declares? |

Every state of a screen or a component is asked for a capture — only visual units are confronted at
all. The exception is written where the state is: `@no-vr: <reason>` on its heading or its row, for a
state with no visual value of its own, such as a transient loading or a state that looks like another.
An exemption with no reason does not exempt.

They run on a visual unit's main code file — the screen or the component, which is what a project tags
as visual — and read the unit's spec, feature and tests from it (`Button.tsx` → `Button.spec.md`,
`Button.feature`). A VR scenario is a scenario tagged with the project's visual regime and the code of
the state it captures (`@BUTTN-S01` or `@BUTTN-VR-S01`). A VR test names `BUTTN-VR-S01` in its path or
its text. A baseline is `<Unit>.BUTTN-VR-S01[-variant].<ext>` beside the unit.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the node | a code file | any other kind | this unit: leaves without a verdict |
| the spec | `<Unit>.spec.md` beside the code file, with an `@anchors` header `code:` | none beside it, or one with no code | this unit: leaves without a verdict, saying so |
| the states | codes `{CODE}-<letter><NN>` of the State letter | none registered | `vr-states-covered` leaves without a verdict |
| the feature | `<Unit>.feature` beside the spec | missing | this unit: it declares no VR scenario |
| the tests | the map's test nodes of the unit | none | this unit: no VR test is found |

## Effects

| Effect | Description |
| --- | --- |
| `VRSTC-B01` | A node that is not code, a code file with no spec beside it (a part of the unit) and a spec with no code leave every gate without a verdict. |
| `VRSTC-B02` | `vr-states-covered` fails naming each state of the spec with no VR scenario and no exemption, and the regime tag to use; with no state it leaves without a verdict. |
| `VRSTC-B03` | `vr-scenarios-tested` fails naming each VR scenario that no VR test names; with no VR scenario it leaves without a verdict. |
| `VRSTC-B04` | `vr-scenarios-tested` also fails naming each VR scenario with no image beside the unit, `<Unit>.{CODE}-VR-<state>` with an optional variant, as png, jpg, jpeg, webp, gif or svg, and the name to save it under. |
| `VRSTC-B05` | `vr-scenarios-of-states` fails naming each VR scenario of a state the spec does not register, and each VR scenario of a state the spec exempts with `@no-vr`. |
| `VRSTC-B06` | `vr-scenarios-of-states` fails naming each VR scenario that carries no state's code — a scenario capturing the whole unit at once. |
| `VRSTC-B07` | `vr-tests-of-scenarios` fails naming each VR state a test names with no VR scenario in the feature, and the tests that name it; with no VR test it leaves without a verdict. |
| `VRSTC-B08` | A VR scenario is a feature line carrying the visual-regime tag; its state is a `{CODE}-<state>` or `{CODE}-VR-<state>` code on that line. |
| `VRSTC-B09` | A test is of the unit when its path names the unit's code, when it sits beside the unit under the unit's name, or when a folder of its path is named after the unit; an image is never a test. |
| `VRSTC-B10` | The State letter is the one of the project's rule type whose term or a section starts with "state" or "estado"; without one, `S`. |
| `VRSTC-B11` | `vr-baseline` accepts the same image formats for a VR scenario's baseline. |
| `VRSTC-B13` | Every message the spec catalogs (the codes of its User Messages section) is captured like a state: a VR scenario `{CODE}-VR-M01`, a VR test and a baseline image, `@no-vr: <reason>` on its row exempting it — an error shows on the screen as its message, the state is the same. |
| `VRSTC-B14` | The states a spec registers are the codes in its States section — whose title may carry a note in parentheses — when it has one; a state code cited elsewhere registers nothing. |
| `VRSTC-B12` | A state is exempted from visual regression by `@no-vr: <reason>` on a line that declares it — its heading or its row in a states table. An exemption with no reason does not exempt: the state is still asked, and the failure names it as an exemption with no reason. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `VRSTC-I01` | One state's scenario, test or image never answers for another state. | a unit whose S01 has scenario, test and image and whose S02 has none fails naming S02 alone |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `VRSTC-E01` | The spec beside the code file cannot be read, or there is none. | No verdict, saying there is no spec beside the unit's main file. | A part of the unit has no spec; failing it would charge every part of a screen. |
| `VRSTC-E02` | The unit's feature cannot be read, or there is none. | It declares no VR scenario: every state lacks one. | A missing feature is the same gap as a feature without the scenarios. |
| `VRSTC-E03` | Looking for a baseline image fails on the file system. | Counted as no image. | Unseen is not found; passing it would approve a state nobody captured. |
| `VRSTC-E04` | A test of the unit cannot be read. | Only its path is read. | A capture flow named by its code still says what it captures; its text cannot be shown to. |
