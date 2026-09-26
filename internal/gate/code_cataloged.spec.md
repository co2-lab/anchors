<!-- @anchors
  code: CDCTC
  updated_at: 2026-09-26
  layer: gate
-->
# CodeCataloged — what the code EXPORTS must be in the spec, or waived in the code

> **Code**: `CDCTC`

## Overview

Confronts a spec against the code it governs, starting from the OTHER side: **this symbol
is public — does it have a rule?**

It is the inverse of `rule-implemented`. That one starts from the spec and asks "does this
rule have code?"; this one starts from the code. Without both, the divergence escapes on
one of the sides — measured in a real project: **a spec catalogued 2 rules for 7 exported
functions, and no gate asked about the remaining 5.**

**Noise is not an argument for not building** — and this line has already been wrong here.
The earlier version said that charging a spec of every symbol would produce hundreds of
legitimate-but-useless findings, and a gate that accuses everything is switched off. The
fear was right; the conclusion was not. Compare the two outcomes: a granular gate
switched off by the project does not protect, **and the project KNOWS**; a gate too coarse
and left on does not protect either, **and it reports GREEN**. The outcome is the same;
what changes is the honesty. And there is a decisive asymmetry: a noisy gate is
CALIBRATABLE by whoever uses it, while a gate that is too coarse cannot be sharpened by
the project — the decision was taken inside and there is no way to recover it. When in
doubt between granular-with-noise and coarse-with-silence, the default is GRANULAR.

The same principle governs the language: the export pattern used to be TypeScript syntax
built in, so in a Go, Python or Ruby project it matched zero symbols and the gate reported
VERDE — stamping approval, with `blocking: true`, over what it had never read. **Without
knowing how to read, the gate goes quiet; it never approves.**

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the artifact's kind | spec; anything else leaves without a verdict | — (another kind is a case, not an error) | this unit: the ruler starts from the spec that governs the code |
| the map | a built graph, or none | — (with no map there is no way to reach the code) | this unit: with no map it does not approve |
| the governed code | the file the spec `specifies`, read from disk | — (no linked code is a case, and belongs to another gate) | the map, through the `specifies` edge |
| the export pattern | a regular expression the project declared, with at least one capture group | a pattern that does not compile, or has no capture group — which SKIPS, it is never assumed | this unit: not knowing how to read is a reason to go quiet, never to approve |

## Effects

| Effect | Description |
| --- | --- |
| `CDCTC-B01` | An artifact that is not a spec leaves without a verdict: the ruler starts from the spec that governs the code. |
| `CDCTC-B02` | An exported symbol the spec never names FAILS, and the verdict names the orphan and its LINE — the name alone would make the reader hunt for the symbol. |
| `CDCTC-B03` | What the spec already catalogues is never accused. |
| `CDCTC-B04` | `@no-rule: <reason>` on the symbol waives it: not every export deserves a rule, and the waiver is what makes the noise manageable without lying. |
| `CDCTC-B05` | A BARE marker, with no written reason, does not waive — it would be a silent way to quiet the gate, and the trace that a decision was taken would vanish. |
| `CDCTC-B06` | A spec cataloguing every exported symbol passes. |
| `CDCTC-B07` | With no code linked the gate leaves without a verdict: the absence belongs to `triad-complete`, and accusing it in both places would duplicate the debt. |
| `CDCTC-B08` | Without a declared export pattern the gate SKIPS and says it skipped, naming how to enable it. It never approves what it cannot read. |
| `CDCTC-B09` | With the pattern declared the gate confronts for real, in any language — the project's own `export_detect` is the ruler. |
| `CDCTC-B10` | The declared dialect family also supplies the pattern: a Go project needs only name its family. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CDCTC-I01` | The waiver holds in the COMMENT BLOCK above the symbol, not only on its own line. The declaration is documentation: whoever writes it puts the explanation alongside, and the explanation rarely fits on one line. Looking only at `i-1` made the gate ignore the declaration and keep accusing — measured in a file where it sat on the second line of a two-line comment. | confronts the declaration inline, one line above, in a two-line block, in a block with paragraphs and in a doc comment, and verifies it holds in all of them |
| `CDCTC-I02` | The waiver does NOT leak between symbols. A symbol with no comment above it does not inherit another symbol's declaration — inheriting would let one marker exempt the whole file, which is the opposite of what it is. | puts the declaration on the first of two symbols and verifies the second does not carry it |
| `CDCTC-I03` | Green over what was never read is the worst possible failure in a measuring instrument. Not knowing how to read is a reason to go quiet and SAY SO, never to approve — and never silently, because the bias of the built-in ecosystem would hide inside that silence. | confronts a Go file with no declared pattern and verifies it does not return Pass |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CDCTC-X01` | Does not judge whether the catalogued rule DESCRIBES the symbol well. | The search is coarse on purpose: it asks whether the spec names the symbol. Judging whether the rule says the right thing about it is judgement, and judgement belongs to another class of gate. |
| `CDCTC-X02` | Does not decide which symbols deserve a rule. | That is the project's call, and the waiver is where it records it — with a written reason. Deciding here would take away exactly the calibration that makes a granular gate usable. |
| `CDCTC-X03` | Does not know any language: whoever declares what is public is the project. | Recognising a public symbol depends on the language, and Anchors does not presume. The built-in TypeScript pattern is only ever a SUGGESTION to a project that has not declared its own, never a silent default. |
| `CDCTC-X04` | Does not charge the absence of code — that belongs to `triad-complete`. | Accusing the same debt in two gates would duplicate the finding, and whoever fixed one would still see the other. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `KindSpec` | core — the gate needs the node's KIND to know whether it has jurisdiction |
| DEP2 | `internal/mapx/model.go` | `EdgeSpecifies` | core — the governed code is reached through the edge the spec declares |
| DEP3 | `internal/config/config.go` | `Config` | core — the export pattern and the dialect family come from the project's Structure |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
