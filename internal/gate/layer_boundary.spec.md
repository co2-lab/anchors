<!-- @anchors
  code: LYBNL
  updated_at: 2026-09-26
  layer: gate
-->
# LayerBoundary — a layer does not reach what is not its own

> **Code**: `LYBNL`

## Overview

Confronts a code file against the architecture the project DECLARED: **the layers exist on
paper, but do they respect each other?**

Anchors already declared layers — and did not confront whether they hold. Declaring
`screens/`, `hooks/`, `repositories/` and never verifying that the screen does not talk
straight to the repository is drawing the architecture and not defending it: `layers:`
becomes documentation.

The gap surfaced because every project solved it alone. A real project kept a
**366-line shell script with 15 hand-written architectural rules** — all of the same
shape: "files matching THIS pattern must not contain THAT one". Fifteen instances of a
single mechanism, reimplemented because the framework did not offer it.

The rule is declared in `anchors.yaml` and is language-agnostic — **the project writes the
pattern, the engine knows no dialect**: a `layer` says to whom it applies, `forbid` what
must not appear, `because` the reason that keeps it from being ritual, and `severity`
the maturation of the rule. The honest opt-out is `@allow-boundary: <reason>` on the line:
acknowledged debt stays visible and dated in the code, instead of becoming an exception in
a distant list nobody revisits.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the artifact's kind | code; anything else leaves without a verdict | — (another kind is a case, not an error) | this unit: only code carries an import to forbid |
| the declared boundaries | what the Structure declares, or nothing | — (no declaration is a case, and it is Pending, not Pass) | the Structure: what is not declared is not charged |
| the `forbid` pattern | any regular expression the project writes | a pattern that does not compile — which FAILS, it is not ignored | this unit: an invalid pattern is a visible config error |
| the unit's layer | the classification the scan already applies | — (a file outside every layer simply matches no scoped rule) | the scan, from the `layers:` the project declared |

## Effects

| Effect | Description |
| --- | --- |
| `LYBNL-B01` | An artifact that is not code leaves without a verdict: there is no import to forbid in a spec or a feature. |
| `LYBNL-B02` | Content matching a forbidden pattern FAILS, and the verdict names the LINE and the REASON — a prohibition with no motive turns into ritual. |
| `LYBNL-B03` | A rule scoped to a layer charges only that layer: the same import in a hook is legitimate. |
| `LYBNL-B04` | A rule with no `layer` holds for ALL code — that is how a global prohibition is declared (raw clock, literal colour, console log). |
| `LYBNL-B05` | `severity: warn` records without failing; the default severity is `error`. |
| `LYBNL-B06` | `@allow-boundary: <reason>` on the line waives THAT line — acknowledged debt stays visible and dated where it lives. |
| `LYBNL-B07` | The waiver also holds in the comment ON THE LINE ABOVE: in many languages an import has nowhere to carry a readable end-of-line comment, and demanding it inline would push the author not to declare at all. |
| `LYBNL-B08` | A BARE marker, with no written reason, does not waive — it would be a silent way to quiet the gate. |
| `LYBNL-B09` | With no boundary declared the verdict is Pending, never Pass: Anchors does not know the project's architecture, and pretending it checked is worse than saying what is missing. |
| `LYBNL-B10` | An invalid `forbid` pattern FAILS visibly: swallowed in silence it would switch the rule off with nobody knowing. |
| `LYBNL-B11` | The pattern is matched against the WHOLE file, so a target spread over several lines is caught — a multi-line import included. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `LYBNL-I01` | The engine knows no language. The same architectural rule is expressible in the import dialect of TypeScript, Python, Go, Java, Rust or Ruby, because the one who writes the pattern is the project. | writes the same rule in six dialects and verifies violation and legitimate import in each |
| `LYBNL-I02` | Line-to-line matching under-reported. Measured in a reference app: of the 7 screens importing `Modal` from react-native the gate accused ONE — the only one with the import on a single line. The other 6 wrap because they pass 100 columns, and escaped: the screen was accused for FORMATTING, not for being different. | confronts a multi-line import and verifies the failure |
| `LYBNL-I03` | The waiver holds on ANY line of the matched stretch. In an import the formatter wrapped, the natural marking sits on the `from` line; demanding it on the first line of the match would require the author to know where the regex started matching. | puts the waiver on the `from` line of a multi-line import and verifies Pass |
| `LYBNL-I04` | A project that WANTS to anchor the pattern to a single line still can: `^` and `$` keep holding per line, because `(?s)` changes the `.`, not the meaning of the anchors. | confronts a line-anchored pattern against a match in the middle of a line and verifies Pass |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `LYBNL-X01` | Does not decide WHICH boundaries exist. | The architecture is the project's decision, declared in the Structure. A gate that invented boundaries would charge what nobody committed to — and with no declaration it says Pending, not Pass. |
| `LYBNL-X02` | Does not parse the language: it does not read imports, it matches TEXT. | Understanding the import graph of every language would tie the engine to a set of ecosystems. The project writes the pattern in its own dialect, and that is what makes the gate agnostic. |
| `LYBNL-X03` | Does not judge whether the forbidden thing is architecturally wrong. | The ruler here is the DECLARATION. Whether the boundary is the right one to draw is design judgement, and judgement belongs to whoever writes `anchors.yaml`. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Boundary` | core — the rules, their scope, reason and severity are declared in the Structure |
| DEP2 | `internal/scan/scan.go` | `Classify` | core — the layer of the node comes from the same classification the scan uses |
| DEP3 | `internal/mapx/model.go` | `KindCode` | core — the gate needs the node's KIND to know whether it has jurisdiction |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
