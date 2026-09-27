<!-- @anchors
  code: DCOXX
  updated_at: 2026-09-26
  layer: apoio
-->
# DocLayout — the one decision of whether a documentation page shows everything or summarizes

> **Code**: `DCOXX`

## Overview

A documentation page for a layer with three units can carry everything — rules, invariants and
scenarios. With thirty units the same page is a document nobody scrolls to the end. The choice
between a full page and a summary is simple; distributing it is not, because three places must
agree on it: the layer page (full or summary?), the rules index (does the link go to the rule or
to the unit that holds it?) and the behaviour index (the same question for a scenario).

The first version let each place decide on its own, and 483 of 812 generated links came out
broken: the page summarized, the index pointed at the rule's anchor, and that anchor did not
exist. This unit is where the decision is taken ONCE, before any template runs, so every page
and every link asks the same object and there is no second place to disagree.

A selection is "big" when either its unit count or its spec line count passes the cut-off. Both
measures exist because counting units alone misleads: five long specs make a longer page than
twenty short ones. The default cut-off is explicit and configurable, so the documentation format
never changes by itself on the day the next unit enters without anyone having decided it.

## Effects

| Effect | Description |
| --- | --- |
| `DCOXX-B01` | A selection is big when its unit count is above the unit cut-off OR its line count is above the line cut-off; either one passing is enough. |
| `DCOXX-B02` | `DefaultLayout`: The default cut-off, when the project declares none, is 20 units or 2000 lines. |
| `DCOXX-B03` | The summary sentence is empty for a selection that is not big; for a big one it states the unit count, the rule count and both cut-offs, so the reader knows the page is a summary and where the full text is. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the measured selection | a size with non-negative unit, rule and line counts | negative counts | the size measure of the link unit (`DCLND`), which only counts |
| the cut-off | any pair of unit and line limits | — | the caller: the default, or the value `docs build` receives |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DCOXX-I01` | The cut-off is strict: a selection exactly at a limit is not big, one past it is. | measures a selection at each limit and one past it, and verifies the answer turns only past the limit |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DCOXX-X01` | The layout only answers "is this selection big"; it never splits a layer into several pages. | Splitting when a limit is crossed would change the number of pages and break every external link on the day the next unit entered. |

## Errors

none — the decision is arithmetic over counts already measured; there is no input it can refuse and nothing it reads.

## Dependencies

(none)

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
