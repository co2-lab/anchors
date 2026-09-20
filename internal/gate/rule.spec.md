<!-- @anchors
  code: RLUEX
  updated_at: 2026-09-19
  layer: gate
-->
# Rule — the identity of a verification INSIDE a gate, and the waiver that names it

> **Code**: `RLUEX`

## Overview

A gate does not verify ONE thing. The one that confronts a spec charges two — that no
placeholder survives, and that at least one rule is catalogued. The one that confronts a
header charges several. Until this unit existed, the verdict said only WHICH GATE failed,
and that had two consequences, both measured:

- whoever waives a deliberate case — the spec that is born before the feature — could
  only waive the ENTIRE gate, and every other verification it did well went down with it;
- two different defects inside the same gate are indistinguishable in a report. The
  reader needs the MESSAGE to know which of the two happened, and a message changes with
  the next rewrite or the next translation.

So the verification gets an identity of its own: `<gate>/<rule>`. It is the same
principle the doctrine already applies to the artifact — a stable identifier is what
lets one speak of a decision without describing it again.

**What separates this unit from its neighbours**: the engine decides WHO runs and
aggregates the answers; the registry decides WHICH function answers. This unit decides
NOTHING about running — it only knows how a verification is named, and whether somebody
declared, in writing and with a reason, that this particular verification should not be
confronted this time. It reads no file, consults no map, and returns no verdict.

**The reason is part of the datum, not a comment beside it.** A waiver without a written
justification is indistinguishable from somebody fleeing a gate that found a defect. The
second measured failure is finer: a waiver restricted to targets that is asked about
WITHOUT a target must answer "not waived" — otherwise the gate is filtered out of the
list entirely, and waiving four new specs would erase the gate for the whole repository.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the rule identifier | `<gate>/<rule>`, or the bare gate name when the gate has a single verification | — (an unknown name is a case, not an error: it simply is not waived) | this unit: a name it never saw is answered "not waived" |
| the textual waiver form | `rule=reason`, optionally `rule@CODE=reason`, comma-separated | an entry with no reason, and a target that is a path | this unit: it collects an error per malformed entry instead of accepting it in silence |
| the commit-message form | the bracketed marker carrying rule, optional code and reason | a marker whose reason is blank after the colon | this unit: the marker is collected as an error |
| the waiver target | the artifact CODE, and only it | a path, a glob, a file name | this unit: a target carrying a separator, a dot or a star is refused with an error that names what to use instead |
| the waived artifact's code | the code of the node being confronted, or nothing | — (no code is a case: it is answered "not waived") | the caller, which reads the code off the node |

## Effects

| Effect | Description |
| --- | --- |
| `RLUEX-B01` | `NewRuleID` joins gate and rule with a separator, so that short rule names cannot collide between gates. |
| `RLUEX-B02` | `NewRuleID` with an empty rule returns the gate name alone: a gate with a single verification gains no separator it does not need. |
| `RLUEX-B03` | `Gate` returns the gate half of the identifier, and `Rule` the rule half. |
| `RLUEX-B04` | `Rule` is empty when the identifier carries no rule half — the verdict then belongs to the gate as a whole. |
| `RLUEX-B05` | `ParseWaiver` refuses an entry with no reason, and an entry whose reason is blank, because the reason is the only thing separating a deliberate waiver from an ignored gate. |
| `RLUEX-B06` | `ParseWaiver` refuses an entry with no rule name. |
| `RLUEX-B07` | `ParseWaiver` refuses a target that looks like a PATH, and the error names the artifact code as what to use instead. |
| `RLUEX-B08` | `ParseWaiver` refuses a target marker with nothing after it, which is a typo that would otherwise produce a waiver that waives nothing. |
| `RLUEX-B09` | `Waived` accepts the two granularities: waiving the gate covers every rule inside it, and waiving one rule preserves the rest of the gate. |
| `RLUEX-B10` | `Waived` answers "not waived" for a waiver that DECLARES targets, so the caller cannot drop the gate from the list and erase it for the whole repository. |
| `RLUEX-B11` | `WaivedTarget` waives the named codes and confronts every code that was not named. |
| `RLUEX-B12` | `WaivedTarget` with no declared target waives every code, which is the coarse exit a freshly declared gate still needs. |
| `RLUEX-B13` | `WaivedTarget` answers "not waived" for an artifact with no code: without identity there is no specific target to waive. |
| `RLUEX-B14` | Each target carries ITS OWN reason, so two waivers of the same rule no longer make the second overwrite the first. |
| `RLUEX-B15` | `WaiverFromMessage` reads the waivers declared in the commit message, which is the form that survives in the history beside the why of the change. |
| `RLUEX-B16` | `WaiverFromMessage` refuses a marker whose reason is blank, by the same guarantee as the textual form. |
| `RLUEX-B17` | `Merge` joins two waivers, so a pipeline hook using the environment variable and an author writing the marker can coexist. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RLUEX-I01` | A waiver never reaches what nobody waived. The error that would cost most here is one leaking onto a gate that was never named. | waives one gate and verifies that another one, and the empty waiver, waive nothing |
| `RLUEX-I02` | Every accepted waiver carries a written reason, in both input forms. A waiver with no why is indistinguishable from an ignored gate, and the report would show it waived without saying why. | parses both forms with a blank reason and verifies both are refused |
| `RLUEX-I03` | A refusal always produces an error that reaches the caller, never a silent acceptance with an empty reason. A waiver that does not waive would fail the next commit with no visible explanation. | collects the errors of each malformed form and verifies the count |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RLUEX-X01` | Does not accept a PATH as the waiver target. | A path is not identity. It changes when somebody reorganises folders and the waiver stops holding in SILENCE; it is ambiguous by nature and invites broad forms that rebuild the very problem per-target granularity came to solve; and an artifact with no code is an earlier problem that another gate charges — a side exit here would hide the cause. |
| `RLUEX-X02` | Does not decide whether a rule PASSES. | The verdict belongs to the checker. This unit only knows how a verification is named and whether it was waived; mixing the two would put the escape hatch inside the ruler. |
| `RLUEX-X03` | Does not read files, the map, or the project structure. | Identity and waiver are pure text arriving by flag, by environment variable or by commit message. A unit that went looking at disk would make the waiver depend on the state of the repository, and the same waiver would mean different things in two checkouts. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none — identity and waiver are pure text, and depending on the map or the structure would
make the same waiver mean different things in two checkouts

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
