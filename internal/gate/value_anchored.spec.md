<!-- @anchors
  code: VLANV
  updated_at: 2026-09-19
  layer: gate
-->
# ValueAnchored — every value of a closed set points at the rule that justifies it, and the anchor carries the value

> **Code**: `VLANV`

## Overview

Confronts a closed set against the question the symbol-level gate cannot ask: **does every
value have an address, and does that address tell the truth?**

`code-cataloged` charges the SYMBOL. One `@no-rule` written on the line of
`export const JANELAS` releases the whole list at once. But each value of a closed set is a
SEPARATE domain decision — whoever adds `6h` is deciding something, and whoever removes
`24h` is deciding something too. A waiver granted per symbol pays for all of them with one
signature.

What that cost, measured: the `QueryScope` of a real project declares
`['15m','1h','6h','24h']` with a `@no-rule` on the `export` line. The four values were never
confronted one by one, and two screens began offering `5m`, `30m` and `1d` — which the
contract does not accept. `resolveEscopo` falls back to the `1h` default when the value is
not accepted, so the screen displayed `5m` while showing one hour of data. Every gate green.

**The SECOND BRACKET is what separates citing from proving.** An anchor that only points at
the rule rots in silence: someone swaps `'15m'` for `'5m'` and the comment still looks
correct. With the value written inside it, the gate confronts what the anchor ASSERTS against
what the line SAYS — and that comparison does not depend on the language, because it compares
two parts of the comment with the line the comment annotates.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, which routes by the declared `on:` |
| the map | a built graph, or none | — (no graph is a case, not an error) | this unit: with no map the verdict is pending, never approval |
| the value anchor pattern | a regex with TWO capture groups, declared by the project | a regex with fewer than two groups | this unit: an anchor that cannot be checked is only a comment that ages, and the gate goes quiet instead of approving |
| the export pattern | the regex by which this project writes a public symbol | — (absent is a case, not an error) | this unit: without it the gate cannot find the set, and skips |
| the confronted code | the text of the file the spec specifies | — (a spec with no code is a case) | this unit: with no code node reached, the verdict is a skip |

> `Who guarantees` cannot be left empty nor say only "not mine": if nobody guarantees, the
> duty is orphaned — and that is exactly where the invalid input gets through.

## Effects

| Effect | Description |
| --- | --- |
| `VLANV-B01` | An artifact that is not a spec leaves the confrontation without a verdict: the gate reaches the code through the spec, and has no jurisdiction of its own over the other kinds. |
| `VLANV-B02` | A value of a closed set WITHOUT an anchor fails: the decision was made and left no address. |
| `VLANV-B03` | The failing verdict NAMES the unanchored value, so whoever reads it does not have to hunt for which one it was. |
| `VLANV-B04` | A value whose anchor carries the rule key AND the value itself passes: the address exists and it checks out. |
| `VLANV-B05` | An anchor that ASSERTS one value while the line SAYS another fails — the anchor lies, and a lying anchor looks like traceability while pointing at the wrong place. |
| `VLANV-B06` | The verdict of a lying anchor shows BOTH sides of the divergence: what the anchor asserts and what the line says. |
| `VLANV-B07` | The lying anchors are reported BEFORE the unanchored ones: the absent anchor can be seen, the lying one cannot. |
| `VLANV-B08` | Without a declared value anchor pattern the gate SKIPS: it cannot read, so it does not approve. |
| `VLANV-B09` | The skip for a missing pattern names the setting that enables it, so the reader learns HOW to turn the gate on. |
| `VLANV-B10` | A declaration that does not OPEN a list on the same line is not a closed set, and is not charged. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `VLANV-I01` | An anchor pattern with fewer than two capture groups is treated as if it had not been declared. Without the second group the anchor asserts no value, and the confrontation the gate exists for cannot happen. | declares a single-group pattern and verifies that the gate does not judge |
| `VLANV-I02` | A line carrying an anchor is NEVER read as the end of the list. The anchor itself holds two `]`, and treating them as a closing bracket ended the set at the first anchor — the gate left the block and confronted nothing. | confronts an anchored set and verifies the values past the first anchor are still charged |
| `VLANV-I03` | With no built map the verdict is pending, never approval. Approving without being able to look would stamp what was not measured. | runs the gate with no graph and verifies it does not approve |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `VLANV-X01` | Does not judge whether `15m` is a GOOD value, nor whether the rule `WINDOW-001` says what it should. | The ruler is that the decision HAS an address and that the address does NOT LIE. Whether the value belongs in the domain is judgement, and judgement belongs to whoever knows the product. |
| `VLANV-X02` | Does not accuse a line whose literal it cannot read. | The literal extraction is deliberately simple — quotes, double quotes, backticks. A false negative here is better than a mass of false positives, which would train the team to ignore the gate. |
| `VLANV-X03` | Does not decide WHICH comment shape is an anchor, nor what a public symbol looks like. | Both are declared by the project, because they belong to how this codebase writes. A gate that invented the shape would charge a convention nobody adopted. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `ValueAnchor` | core — the anchor shape is declared by the project, not invented by the gate |
| DEP2 | `internal/mapx/model.go` | `KindSpec` | core — the kind is what routes the jurisdiction |
| DEP3 | `internal/mapx/model.go` | `KindSpec` | core — the gate reaches the code through the spec, and needs the node's kind |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
