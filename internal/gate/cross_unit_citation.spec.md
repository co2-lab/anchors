<!-- @anchors
  code: CRUCT
  layer: gate
-->
# CrossUnitCitation — a rule of another unit a spec cites lives in the product

> **Code**: `CRUCT`

## Overview

A spec may **reference** another unit — point at its identity, as a screen's Navigation names
the screen it leads to — and may **cite** another unit's rule by its code. What it cites is
content two units share, though, and a shared rule lives in the product doctrine, which comes
first: the cited rule realizes a doctrine rule, so neither spec watches the other — both
follow the product above them (DESIGN-dependencies-out-of-the-spec.md, DOOSD-D02, D04, D05).
The gate names each cited rule of another unit that realizes no product rule; informative by
default. A prose naming another screen without its code is a judgment, asked in the review
guide.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the spec | its text | an artifact that is not a spec | this unit: it leaves without a verdict |
| the units | the specs of the map, by their codes | a product doctrine, a plan | this unit: their codes are no other unit's content |

## Effects

| Effect | Description |
| --- | --- |
| `CRUCT-B01` | A spec citing a rule of another unit's spec that realizes no product rule — no `@realizes` on the rule's line — fails, naming each code and the lines it is cited on; a cited rule that realizes the product passes; the navigation sections and their subsections — a reference —, a line realizing the product, an alias, a revision and a retired rule are no citation, and a spec citing none passes. (`checkCrossUnitCitation`) |

## Errors

none — the gate reads the spec and the map; with no map it is pending, not failing.
