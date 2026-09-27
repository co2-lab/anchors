<!-- @anchors
  code: DCTRN
  updated_at: 2026-09-26
  layer: gate
-->
# Doctrine — the vertical axis: product doctrine exists, is realized, and is never copied

> **Code**: `DCTRN`

## Overview

Co-location ties a triad to one directory, but a business rule that holds for three screens belongs
to none of them. It lives in a product doctrine file, and each spec that concretises it points at the
doctrine rule with a `@realizes` tag. This unit holds the five gates of that vertical axis; each one
catches a silence the others cannot see.

- `plan-doctrine-exists` — a plan that seeds a doctrine promises it will exist. Without this gate the
  promise carries no charge, because the specs the plan also seeds pass every gate with nothing to
  realize.
- `doctrine-realized` — a doctrine rule that no spec realizes is a decision that never reached the
  code, and the doctrine file itself is well-formed, so nothing else looks at that far end.
- `spec-doctrine-exists` — the axis's reference check. A `@realizes` that does not resolve looks
  like traceability while the map gains no edge, typically after a doctrine rule was renamed.
- `doctrine-not-duplicated` — the defect the axis exists to eliminate: a spec that copies the
  doctrine text instead of pointing at it. The ruler is similarity, not equality, because whoever
  copies almost always changes a word.
- `spec-realizes-doctrine` — a layer may declare that every rule of its specs must say which product
  decision it concretises. It is off by default: most rules of a tool like this one are local
  mechanics, while in a product application the proportion inverts, and only the project knows which
  case it is.

The axis has one declared way out, `@TBD: <reason>`, and it means "not yet", never "never": on a
line it turns a charge into debt, which keeps showing up (as Pending) until somebody pays it. A bare
marker, or a marker quoted in backticks by a text that explains it, defers nothing.

An open question of a doctrine (a `-Q` item) names a decision nobody has taken, so it is never a
rule to realize.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | a plan, a product doctrine, or a spec, each for its own gate | any other kind | this unit: each gate skips what is not its kind |
| the map | the graph whose realizes edges carry, in their method, the doctrine rule realized | no map loaded | this unit: the gates that read edges are Pending |
| the doctrine files | files on disk under the paths the plan cites or the edges reach | a citation in prose, a template, a file that cannot be read | this unit: prose and templates seed nothing; an unreadable file confirms nothing |
| the layers | the project's layers, some declaring that they require doctrine | a project with no configuration | this unit: spec-realizes-doctrine skips |

## Effects

### Shared by the axis

| Effect | Description |
| --- | --- |
| `DCTRN-B01` | A `@TBD` marker defers what its line declares only when it carries a written reason and is not inside backticks. |
| `DCTRN-B02` | Without a map, `doctrine-realized`, `spec-doctrine-exists` and `doctrine-not-duplicated` are Pending: they read edges, and there are none to read. |

### plan-doctrine-exists

| Effect | Description |
| --- | --- |
| `DCTRN-B03` | An artifact that is not a plan is skipped. |
| `DCTRN-B04` | A plan whose cited doctrines all exist on disk passes. |
| `DCTRN-B05` | A plan citing doctrines that do not exist fails, naming each missing doctrine once, in order, with their count. |
| `DCTRN-B06` | A missing doctrine cited on a line deferred with `@TBD` is Pending, naming it, instead of failing. |
| `DCTRN-B07` | Only a doctrine path with a directory seeds a doctrine; a bare file name in prose, or a template file, seeds nothing, and a plan that seeds nothing is skipped. |

### doctrine-realized

| Effect | Description |
| --- | --- |
| `DCTRN-B08` | An artifact that is not a product doctrine is skipped. |
| `DCTRN-B09` | A doctrine that catalogues no rule is skipped. |
| `DCTRN-B10` | A doctrine whose every rule has an incoming realizes edge naming it passes. |
| `DCTRN-B11` | A doctrine with unrealized rules fails, naming only the unrealized ones. |
| `DCTRN-B12` | When every unrealized rule is deferred with `@TBD` on its own line, the doctrine is Pending, naming them. |
| `DCTRN-B13` | An open question (`-Q`) is not a rule, and is never charged for a realizer. |

### spec-doctrine-exists

| Effect | Description |
| --- | --- |
| `DCTRN-B14` | An artifact that is not a spec is skipped. |
| `DCTRN-B15` | A spec that declares no `@realizes` outside the lines deferred with `@TBD` is skipped. |
| `DCTRN-B16` | A declared rule passes when a realizes edge of the spec names it, or when a doctrine the spec's edges reach catalogues it. |
| `DCTRN-B17` | A declared rule that no reached doctrine catalogues fails, naming only the unresolved rules. |

### doctrine-not-duplicated

| Effect | Description |
| --- | --- |
| `DCTRN-B18` | An artifact that is not a spec is skipped. |
| `DCTRN-B19` | A spec with no readable realized doctrine text is skipped. |
| `DCTRN-B20` | When the rules of both sides together are fewer than four, the gate is Pending: the similarity weights cannot tell shared words apart; from four rules on, the pair is measured. |
| `DCTRN-B21` | A spec rule whose text copies the doctrine rule it realizes fails, naming the spec rule, the doctrine rule and the similarity as a percentage. |
| `DCTRN-B22` | A near copy, with a word changed or dropped, fails as well. |
| `DCTRN-B23` | A spec rule whose text says what is specific to its unit passes. |
| `DCTRN-B24` | A copy on a line deferred with `@TBD` is not charged: the wording is still being worked out. |

### spec-realizes-doctrine

| Effect | Description |
| --- | --- |
| `DCTRN-B25` | An artifact that is not a spec, or a project with no configuration, is skipped. |
| `DCTRN-B26` | A spec whose layer does not require doctrine is skipped. |
| `DCTRN-B27` | Where the layer requires doctrine, a rule with no `@realizes` fails, naming it. |
| `DCTRN-B28` | A rule that declares what it realizes passes, whether the tag sits on the rule's line or on the lines right below it. |
| `DCTRN-B29` | A rule deferred with `@TBD` is Pending, not failure. |
| `DCTRN-B30` | A blank line ends the rule a tag belongs to: a `@realizes` after a blank line declares nothing for the rule above. |
| `DCTRN-B31` | The demanding layer is resolved from the target the spec describes: the layers of the code its specifies edges reach, or, before that code exists, the layer the target's path would have. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DCTRN-I01` | Only a realizes edge realizes a doctrine rule: an edge of any other type that carries a rule code never counts. | adds a specifies edge carrying a rule code, and verifies that rule is still charged as unrealized |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DCTRN-X01` | Does not call a pair a copy because the two share a rare word: a pair judged similar must also reach a similarity of one half. | A spec that realizes a rule is expected to speak its vocabulary; without the floor, honest realizers that share one distinctive word were accused at scores under a fifth. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DCTRN-E01` | A doctrine that a realizes edge of the spec reaches cannot be read. | The rule the edge names still resolves; any other rule the spec declares for that doctrine fails as unresolved. | The edge proves only that the file was found when the map was built; without reading it, no other rule can be confirmed, and assuming it would certify a renamed rule. |
| `DCTRN-E02` | REF[DCTRN-B19]: a realized doctrine that cannot be read contributes no text, and B19 answers the spec left with none | — | — |
| `DCTRN-E03` | REF[DCTRN-B05]: a cited doctrine whose file cannot be found is the missing doctrine B05 charges | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config`, `Layers` | core — the layers that require doctrine |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized verdict messages |
| DEP3 | `internal/mapx/model.go` | `Graph`, `Node`, `EdgeRealizes`, `EdgeSpecifies` | core — the realizes edges of the axis |
| DEP4 | `internal/scan/scan.go` | `Classify` | scan — the layer a target path would have |
| DEP5 | `internal/similarity/similarity.go` | `Weights`, `Classify` | apoio — the similarity ruler |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
