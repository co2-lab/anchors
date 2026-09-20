# Spec guide — how to write a `.spec.md` in this project

> Seeded by `anchors init`. It is the built-in ruler (`anchors guide spec`) instantiated
> with THIS project's dialect. Read it before writing any spec.

## Start from the command, not from the text

Do not write the spec from scratch. The CLI emits a skeleton that already conforms:

```sh
anchors new spec <Name> --out <path>/<Name>.spec.md   # generates the skeleton
anchors new spec --list-sections                       # the sections and WHEN to use each
```

The command solves what is easiest to get wrong: it generates the identity code, writes
the `@anchors` header in the right dialect, and uses the exact format of a catalogued
rule. Then fill it in and confront it with `anchors check --changed <file>`.

## The format the gate demands

A rule is **catalogued** when it has a code AND a structured place. Three forms count, and
a loose mention in prose does NOT:

```md
### ANCH-B01 — the rule's description       <- heading (preferred)
| `ANCH-B02` | description |                    <- table row
- **ANCH-B03** description                      <- bold bullet
```

## Full example (copy and adapt)

```md
<!-- @anchors
  code: ANCH
  updated_at: 2026-01-15
  layer: screen
-->
# Login — authenticates the user and takes them into the app

> **Code**: `ANCH`

## Overview

The entry screen: it takes an e-mail and a password, authenticates, and navigates Home.

## Rules

### ANCH-S01 — Initial state
Empty fields, the sign-in button disabled.

### ANCH-A01 — Sign in with valid credentials
Authenticates and navigates Home.

### ANCH-V01 — Invalid e-mail
The field shows the message and the submit does not fire.

### ANCH-R01 — Only the anonymous reach it
An active session is redirected Home.

## Open Decisions

| Question | Who decides | Becomes |
| --- | --- | --- |

none
```

## The code's letter states the rule's NATURE

This project does not declare `rule_types`, so the framework's canonical letters hold:
`S` state, `R` permission, `V` validation, `A` action, `X` constraint, `B` behaviour,
`N` navigation, `M` message, `D` data. Declaring your own in `rule_types` makes the team's
vocabulary count in place of the generic one.

## The sections

Three are mandatory — header, overview and rules — plus the open decisions. The rest come
in with `--with <key>` when the unit calls for them. Run
`anchors new spec --list-sections` for the list with each one's choice criterion; it
includes the mutually exclusive ALTERNATIVES (`contract` or `signature`, `rules` or
`effects`), which is where the wrong choice costs a rewrite.

## What NOT to do

- **A rule with no code.** With no identity, the feature and the test have nothing to
  cite — the triad does not close and the relational gates are left without a target.
- **Describing implementation.** The spec states the BEHAVIOUR; the function's name and
  the library change without the rule changing.
- **Repeating the copy.** The text shown to the user lives once (in the messages
  section); the other sections reference its code.
- **Guessing at what is ambiguous.** It becomes a line in *Open Decisions* — the
  `open-questions-resolved` gate charges someone to decide, and that is the point.

## Compliance points

- CK1: every rule carries a code AND sits in one of the three structured forms (heading,
  table row, bold bullet) — a rule mentioned only in prose is not catalogued, and the
  feature and the test have nothing to cite.
- CK2: the code's letter matches the rule's NATURE according to the declared vocabulary
  (or the canonical letters, when the project declares none).
- CK3: the spec states behaviour, not implementation — no function name, library or
  framework detail stands where the rule should be.
- CK4: the text shown to the user appears in ONE place, and the other sections reference
  its code rather than repeating it.
- CK5: what is ambiguous became a line in *Open Decisions* instead of being guessed at.
- CK6: the mandatory sections are present — header, overview, rules, and open decisions.
- CK7: when the unit calls for mutually exclusive alternatives, exactly one of each pair
  is present (`contract` or `signature`, `rules` or `effects`), never both.

## Specialise this file

It is born generic. As the project settles conventions (spec profiles by unit type, real
examples, RIGHT/WRONG cases taken from the repo itself), edit it here — it is THIS
project's ruler, and the `governs` in `anchors.yaml` links this guide to the targets it
rules.
