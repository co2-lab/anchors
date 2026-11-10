<!-- @anchors
  code: VRPRS
  updated_at: 2026-11-10
  layer: gate
-->
# VariantsProven — a rule that names its variants is proven once per variant

> **Code**: `VRPRS`

## Overview

A rule may speak of several cases of one discriminant — an item that is an audit or a
notification, a payment by card or by pix. `spec-feature-match` asks one scenario per rule, and
one scenario proving one case satisfies it. Measured in MIF (Sentry MIF-1S): the alerts rule
spoke of audits and notifications, its one scenario tapped an audit, and the notification path
reached the store review broken, with every gate green — mutation at 100%, because the code both
cases share has no branch per case, and every mutant died by the audit's test.

The rule names its variants on its line, `@variants: audit, notification`, and each variant asks
a scenario of the rule that names it: a `@variant:<value>` tag beside the rule's tag, or a row of
the scenario's Examples. A rule that follows another's variants — an error rule shared by every
kind — writes `@variants: B02`, and asks the new kind as soon as the rule it follows gains it.

A spec that names no variants is not this gate's subject: the gate is in the catalog, and asks
nothing until a rule declares them.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the variants | `@variants:` on a rule's line (a table row, a list item, a heading that opens with the code), a comma-separated list, or one rule of the same unit (`B02`, `UNIT-B02`) | a marker on a line that declares no rule of the unit | this unit: it is read as no declaration |
| a variant's proof | a scenario tagged with the rule (with or without a `#nn` suffix) carrying `@variant:<value>`, or whose Examples have a row naming the value, compared without case | a scenario of another rule | this unit |
| the features | those the map links by `covered-by` | — | the map; with none, `unit-complete` accuses |

## Effects

| Effect | Description |
| --- | --- |
| `VRPRS-B01` | A rule's own variant that no scenario of the rule names — by its tag or a row of its Examples — fails, naming the rule and the variant. (`checkVariantsProven`, `declaredVariants`) |
| `VRPRS-B02` | A rule that follows another's variants asks each of them; one with no scenario of the follower diverges, saying the rule is proven only without it. |
| `VRPRS-B03` | A spec whose rules name no variants, a spec with no linked feature, or an artifact that is no spec, is skipped. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `VRPRS-E01` | A rule follows the variants of a rule that names none of its own. | Fails, naming both. | A declaration that stands for nothing would ask nothing, and read as if it did. |
| `VRPRS-E02` | No map is built. | Pending, saying there is no map. | The features are found by the map's edges; without them, Pass would claim a proof nobody checked. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
