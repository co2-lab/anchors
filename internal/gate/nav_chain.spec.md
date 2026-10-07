<!-- @anchors
  code: NCGNV
  layer: gate
-->

# NavigationChain — every navigation flagged with the screen it leads to, and the screens' tables confronted with it

> **Code**: `NCGNV`

## Overview

A screen spec declares its route and its navigation tables, and nothing tied its Out table
to the navigation calls of its code, nor one screen's Out to the other's In: a screen could
lead where its spec did not say, with every gate green, and the app had no navigation map.
This unit confronts the navigation flags the agent writes beside each call with the code and
with the specs (DESIGN-dependencies-and-navigation.md):

- `nav-annotated` — every navigation call, back and reset included, carries the navigation
  flag naming the screen it leads to, or a waiver with its reason;
- `nav-matches-spec` — a screen's Out table and its code's flags name the same screens;
- `nav-symmetric` — a screen's Out leads to another exactly when that one's In comes from it;
- `nav-reachable` — every screen is reachable from the app's entry routes.

A call is read by a pattern — the dialect's `navigation_call`, its destination in the group
`route`, or the family's —, a screen by its spec's stem and its declared route.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted node | a code file (nav-annotated), a screen spec (the others) | any other node | this unit, which skips it |
| the navigation pattern | the dialect's `navigation_call`, or a family that knows one | a dialect with neither | this unit: pending |
| the entry routes | `navigation.entry` | none declared | this unit: nav-reachable pending |

## Effects

| Effect | Description |
| --- | --- |
| `NCGNV-B01` | `nav-annotated` fails naming each navigation call with no flag on its line or alone on the line above, by line and destination, and each flag whose screens do not include the one its route leads to; a waived call, a back navigation and a route no screen declares are not charged for the screen. (`checkNavAnnotated`, `NavCall`, `navCallsOf`) |
| `NCGNV-B02` | `nav-matches-spec` fails naming each Out row no flag of the screen's code answers, each Out row that names no screen of the app, and each flag the Out table does not declare; a spec that is no screen is skipped. (`checkNavMatchesSpec`, `Screen`, `NavRow`) |
| `NCGNV-B03` | `nav-symmetric` fails naming each Out whose destination's In does not come from the screen, and each In whose origin's Out does not lead to it. (`checkNavSymmetric`) |
| `NCGNV-B04` | `nav-reachable` fails a screen no path reaches from the entry routes — through the Out tables and the code's navigation flags (`NavEdges`) —, and is pending with no entry declared. (`checkNavReachable`) |
| `NCGNV-B05` | The fixer writes the navigation flag of each unflagged call whose route names one screen of the app, and leaves a back navigation and a dynamic route to the author. (`fixNavFlags`) |

## Errors

none — an unreadable spec or file is skipped by the index it belongs to, and a pattern that does not compile reads no call; neither is a failure of this unit.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/scan/scan.go` | `NavigatesIn` | scan — the navigation flags |
| DEP2 | `internal/config/dialect.go` | `DialectFor` | config — the navigation pattern |
