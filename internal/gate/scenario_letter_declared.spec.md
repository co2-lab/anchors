<!-- @anchors
  code: SCLTR
  updated_at: 2026-09-20
  layer: gate
-->
# ScenarioLetterDeclared — the letter of a scenario code exists in the vocabulary

> **Code**: `SCLTR`

## Overview

Confronts a scenario code against the project's own vocabulary: **does the letter it
carries mean anything?**

The project declares which letters it recognises, and a sibling gate charges the spec's
SECTIONS for using the right ones. Nobody charged the same of the code a SCENARIO carries
— and an invented letter passes every other gate. The code matches itself across feature
and test, the triad is complete, the relational gates find both ends of the edge, and
nothing notices that the letter means nothing.

Measured in the project that originated this gate: **18 codes carrying five undeclared
letters**, always accompanied by tags equally outside the vocabulary. The pattern is
recognisable once seen — someone needed a nature the project did not have and invented it
instead of declaring it.

Both repairs are legitimate, and the choice belongs to the project: declare the letter (if
the nature is genuinely missing) or remap the scenario onto a letter that already exists
(if it was only an alias for one). **The gate does not choose — it shows what is outside.**

That is why the verdict is UNDETERMINED and not a failure. Discovering that a nature was
used without registration is information; deciding between adopting it and remapping it is
work for whoever knows the domain.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the vocabulary | the letters the project declares as rule types | — (an undeclared vocabulary is a case, not an error) | the project's Structure, never a list invented here |
| the code length | what the project declares | — | the Structure, read at each call and not frozen at load |

## Effects

| Effect | Description |
| --- | --- |
| `SCLTR-B01` | An artifact that is not a feature leaves without a verdict: scenario codes live in features. |
| `SCLTR-B02` | With no declared vocabulary the gate leaves without a verdict — with nothing to compare against, every letter would be either all valid or all invented, and both answers are noise. |
| `SCLTR-B03` | A feature carrying no scenario code at all leaves without a verdict: there is nothing to judge. |
| `SCLTR-B04` | Every letter found inside the declared vocabulary passes. |
| `SCLTR-B05` | A letter outside the vocabulary is reported as UNDETERMINED, never as a failure — the nature may deserve declaring, and that decision is not the gate's. |
| `SCLTR-B06` | The verdict NAMES the letters that are outside and the codes that carry them, because that is what the reader needs in order to choose between declaring and remapping. |
| `SCLTR-B07` | Codes sharing one unknown letter are grouped into a SINGLE line — eight scenarios with the same invented letter are one finding, not eight. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `SCLTR-I01` | The scan is over the SHAPE of a code, never over the vocabulary. Reusing the shared scenario parser would make the gate blind to exactly what it exists to find: that parser builds its pattern FROM the declared letters, so a code with an invented letter is precisely what it discards. The vocabulary enters afterwards, when judging each letter found. | confronts a feature whose code carries an undeclared letter and verifies it is seen at all |
| `SCLTR-I02` | The code-length pattern is read at EVERY call, not frozen once. It comes from the project's Structure, which loads after the package globals — a pattern built once at load time would freeze the default and silently ignore the project's declaration. | confronts a project declaring a non-default code length and verifies its codes are recognised |
| `SCLTR-I03` | A valid letter is never named in the verdict, even when the same feature carries invented ones. Reporting what is correct alongside what is not would make the reader hunt for the real finding. | confronts a feature mixing valid and invented letters and verifies only the invented ones appear |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `SCLTR-X01` | Does not decide whether the invented letter should be declared or remapped. | Both repairs are legitimate and the choice depends on whether the nature is genuinely missing from the project — domain knowledge the gate does not have. Choosing for the author would impose one of two valid answers. |
| `SCLTR-X02` | Does not judge the TAGS accompanying the code, only the letter of the code itself. | The measured pattern shows invented letters travel with tags equally outside the vocabulary, but a tag is free vocabulary by design. Charging it here would turn a precise instrument into a style opinion. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `CodeLengthPattern` | core — the code shape is declared by the project |
| DEP2 | `internal/mapx/model.go` | `KindFeature` | core — the kind is what routes the jurisdiction |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
