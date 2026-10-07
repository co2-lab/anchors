<!-- @anchors
  code: MNCMP
  layer: comando
-->

# MapNav — the app's navigation, screen by screen

> **Code**: `MNCMP`

## Overview

`anchors map nav [screen]` prints the navigation between the app's screens — the screens' Out
tables and the code's navigation flags (DESIGN-dependencies-and-navigation.md): every screen
with where it leads, or one screen with where it comes from and where it leads.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the screen | a screen's spec stem, or the stem without `Screen` | any other name | this unit: refuses naming it |

## Effects

| Effect | Description |
| --- | --- |
| `MNCMP-B01` | With no screen named, each screen that leads somewhere is printed as `A → B, C`, the screens and their destinations sorted. (`NavLines`) |
| `MNCMP-B02` | With a screen named, its origins (`←`) and its destinations (`→`) are printed; a name that is no screen of the navigation is refused. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MNCMP-E01` | REF[MNCMP-B02]: a name that is no screen is refused, naming it | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/gate/nav_chain.go` | `NavEdges` | gate — the navigation between screens |
