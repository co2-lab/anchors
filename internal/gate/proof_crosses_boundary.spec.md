<!-- @anchors
  code: PCBPR
  updated_at: 2026-09-26
  layer: gate
-->
# ProofCrossesBoundary — when a rule claims a relation, the proof must reach the other side

> **Code**: `PCBPR`

## Overview

Confronts a spec against the question the whole triad leaves open: **a rule says it mirrors
another unit — does the governed code actually IMPORT that unit, or is the claim prose while
the proof stays local?**

The measurement comes from an audit of 51 spec×code divergences in a reference app
(2026-08). The three GRAVEST findings had the same shape: two sides defined the same
thing, each side had its own test, and each test confronted its OWN copy.

| finding | the divergence |
| --- | --- |
| RFB codes | `12` = "Terreno" in the app, `12` = "Casa" in the backend — and the document filed with the tax authority carried the wrong heading |
| balance filter | `=== 'statement'` in the app, `!== 'invoice'` in the backend |
| minimum boundary | `<=` on two screens, `<` on the third and in the audit |

In every one of them 53 gates stayed green — correctly, by the rulers they had. The triad
was complete: rule declared, scenario written, test whose title matched. What no gate
asked was whether the PROOF reaches the other side.

The case that motivated the design is ALIVE in the repository and has not yet diverged:
the seat-price rule declares "the price mirrors the backend `orgBilling.ts` — diverging here lies
about the billing", the scenario says "the two values are the same the backend charges",
and the test does `expect(SEAT_PRICE.individual).toBe(15)`. It proves it is 15; it does
not prove it is the same the backend charges. Changing the backend to 18 keeps everything
green.

**What separates it from its neighbours:** `triad-complete` asks whether the test exists,
`feature-test-match` confronts scenario against test case, `rule-implemented` asks whether
the rule reached the code — all three answer "yes" about a rule whose proof never leaves
its own file. And like `dependency-honored`, this gate charges only what the spec DECLARED:
it does not go hunting duplicated concepts across the project.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the artifact's kind | only `mapx.KindSpec` is judged | code, test, feature — they leave without a verdict | this unit: a kind that is not a spec returns Skip |
| the rule lines | catalogued rule rows (a table row opening with the rule identity in backticks), in the code length the project declares | prose paragraphs around the table | this unit, via `ruleLineRE()`, which reads `config.CodeLengthPattern()` per call |
| the declared target | a rule CODE in backticks (preferred) or a cited FILE on the same line | a target named in a neighbouring line or paragraph | this unit: the claim and the target must share the line |
| the map | a built graph, or none | — (no graph is a case, not an error) | this unit: with no graph the verdict is Pending, never Pass |
| the governed code | the files reached by `mapx.EdgeSpecifies` out of the spec | the test file, which may legitimately not import | this unit: the charge lands on the code, not on the proof |

## Effects

| Effect | Description |
| --- | --- |
| `PCBPR-B01` | An artifact that is not a spec leaves without a verdict: the gate has no jurisdiction over code, test or feature. |
| `PCBPR-B02` | A rule with the declared single-source mark and a cited file whose governed code does NOT import it fails — the claim is prose and the proof is local. |
| `PCBPR-B03` | A citation that lives only in a COMMENT does not satisfy the charge: comments are stripped before the code is read, because counting them would approve exactly the case that motivated the gate. |
| `PCBPR-B04` | Governed code that really imports the cited unit passes. |
| `PCBPR-B05` | An import written through an alias different from the path in the spec still matches: the comparison is by module base name without extension, because alias and extension vary per project and the module name does not. |
| `PCBPR-B06` | A rule carrying the declared waiver with a written reason is not charged: the relation exists and is not importable (network contract, generated file, value living in an external provider). |
| `PCBPR-B07` | A rule carrying the OWNER stamp is not charged: the owner does not mirror anybody — it IS the source, and there is nothing to import. |
| `PCBPR-B08` | A relation claimed in PROSE, with no declared mark, is reported as a suspicion — it teaches the convention instead of barring the delivery on the first encounter. |
| `PCBPR-B09` | The declared mark WITHOUT a target on the line is reported too: the mark says "I mirror something", and with no something there is nothing to confront. |
| `PCBPR-B10` | A target declared by RULE CODE is resolved through the map to the files of that triad, so the rule carries stable identity instead of a path that moves. |
| `PCBPR-B11` | A rule code that resolves to no unit is reported as unresolved, not silently dropped: a rule pointing at nothing warns nobody. |
| `PCBPR-B12` | A rule code of the unit ITSELF is not a target: self-reference is not the other side of a boundary. |
| `PCBPR-B13` | A rule with no relation claim at all leaves without a verdict — there was nothing to charge. |
| `PCBPR-B14` | Imports in other language shapes (`require`, `from `, `use `, `using `, `#include`, or the project's configured pattern) satisfy the charge the same way. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PCBPR-E01` | A spec is confronted with no map built. | Pending, saying no map is loaded — even when a rule carries the single-source mark. | The governed code and the targets named by rule code are only reachable through the map. Without it nothing was looked at, and neither Pass nor Fail can be said about an import nobody read. |
| `PCBPR-E02` | A marked rule demands an import, and every file the spec governs (`specifies`) is gone from disk or unreadable. | Pending, saying the linked code could not be read. | A map older than the tree can point at code that was moved or deleted. Failing would charge a missing import on a file that is not there; passing would approve a proof nobody saw. The next map build points the spec at the real files. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PCBPR-I01` | The claim and the target must be on the SAME rule line. An assertion floating in a surrounding paragraph binds the demand to prose, not to a catalogued rule. | places the relation claim outside the rule row and verifies nothing is charged |
| `PCBPR-I02` | The import is charged on the GOVERNED CODE, never on the test. A test that exercises the real function without importing the unit itself would be a false positive. | resolves the charge through the `specifies` edge and verifies the code file is the one read |
| `PCBPR-I03` | A satisfied charge does not swallow a pending suspicion: a file may carry one marked rule that passed and another in prose that nobody confronts. | confronts a spec with both and verifies the suspicion still surfaces |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PCBPR-X01` | Does not hunt for duplicated concepts across the project. | That would demand sweeping the cartesian product of the layers, or judgement. The gate acts only on what the spec DECLARED — the same economy as `dependency-honored`, which charges only the symbol promised in backticks. Duplication nobody named is not reached; the gate stops it from COMING BACK once named. |
| `PCBPR-X02` | Does not compare the VALUES on the two sides. | It asks whether the proof reaches across, not whether the two sides agree today. Reading and comparing both definitions would demand a language-specific extractor and judgement about what "the same" means; the import is what makes the divergence impossible to keep silently. |
| `PCBPR-X03` | Does not decide whether an unmarked prose claim blocks. | The verdict is emitted and the `blocking` setting of the gate in the project's Structure decides. It is born informative — visible without barring — and promoting it is a project decision, not the gate's. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `KindSpec` | core — jurisdiction starts from the node's kind, and `EdgeSpecifies` leads to the governed code |
| DEP2 | `internal/config/config.go` | `CodeLengthPattern` | core — the rule-code shape comes from the project's declaration, so the regex is compiled per call and not frozen in a global |
| DEP3 | `internal/i18n/i18n.go` | `T` | core — the verdicts speak the project's language |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
