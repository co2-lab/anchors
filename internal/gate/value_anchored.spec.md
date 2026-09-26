<!-- @anchors
  code: VLANV
  updated_at: 2026-09-26
  layer: gate
-->
# ValueAnchored — a replicated key is declared where it is used, and every copy carries the same value

> **Code**: `VLANV`

> **VLANV-R0001:** the gate stopped charging an anchor on EVERY value of every closed set.
> It now follows DECLARATIONS: a key that is replicated in the code is declared in a comment,
> and the gate confronts the code line below it, every other copy of the key, and — when the
> key is a rule that declares its value in the spec — the spec. Decided by the user: anchoring
> is for values that live in more than one place, and it is the way those places are linked.
>
> **Revises:** `B01`, `B02`, `B03`, `B04`, `B05`, `B06`, `I02`, `X01`, `X02`
> **Checked:** `B07`, `B08`, `I01`, `I03`, `X03`

## Overview

A value that lives in more than one place has no address. Changing a colour means changing it
everywhere, and nothing says where "everywhere" is: the copies are plain literals, and the one
that was forgotten keeps compiling.

The gate turns the replicated value into a symbol. Whoever replicates a value declares its KEY
in a comment, and the code line right below carries the value:

    // @code-reference-[COLOR-SUCCESS]-[#1F8A5B]
    success: '#1F8A5B',

Three confrontations follow from the declaration:

- **LOCAL** — the next code line below it (comment and blank lines skipped) contains the
  declared value: the declaration does not lie about its own line.
- **ACROSS** — every declaration of the same key declares the same value. Without this, a
  change reaches one file — value and declaration together — and every other copy stays
  behind while each file is locally consistent.
- **SPEC** — when the key is a rule code whose defining line in the spec declares a value,
  every declaration is confronted with it. The rule is then the source: change it there, and
  every place still carrying the old value is reported.

The key may therefore have a source in the spec, or be only a replicated reference whose
copies are the truth.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | a code node | specs, features, tests | this unit: skips anything that is not code — declarations live in code |
| the declaration pattern | a regex with TWO capture groups (key, value), declared by the project | a regex with fewer than two groups | this unit: a declaration that asserts no value cannot be checked, and the gate goes quiet instead of approving |
| the map | a built graph | no graph | this unit: without the map the other copies cannot be seen, and the verdict is pending |
| a rule's declared value | a table cell that is ENTIRELY one backticked token, on the rule's defining row | backticked identifiers in headings or prose cells | this unit: reads only whole cells of table rows |

## Effects

| Effect | Description |
| --- | --- |
| `VLANV-B01` | A declaration whose next code line contains the declared value passes. |
| `VLANV-B02` | A declaration whose next code line does not contain the declared value fails, and the verdict shows the key, the declared value and the line. |
| `VLANV-B03` | Comment and blank lines between a declaration and the code are skipped, so declarations may be stacked above one line. |
| `VLANV-B04` | Declarations of the same key with different values fail, in every file that holds one, listing every place and value. |
| `VLANV-B05` | A declaration whose key is a rule declaring a value in the spec, and that disagrees with it, fails, naming the spec and the value it declares. |
| `VLANV-B06` | Without a declared pattern the gate skips and names the setting that enables it. |
| `VLANV-B07` | A declaration with no code line below it fails: it annotates nothing. |
| `VLANV-B08` | Only an anchor that is the whole content of a comment line is a declaration; an anchor inside prose (a comment explaining the syntax) or inside code is a mention, and is not charged. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `VLANV-I01` | A pattern with fewer than two capture groups is treated as not declared. | declares a single-group pattern and verifies the gate does not judge |
| `VLANV-I02` | Every declaration of the project is indexed once per map, not once per file — reading the whole repository per code file is the shape of cost that once made `docs-fresh` 97% of a check. | the index is cached by the map it was built from |
| `VLANV-I03` | With no built map the verdict is never approval. | runs the gate with no graph and verifies it does not approve |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `VLANV-X01` | A rule whose defining line declares no value is not charged against the spec. | Most rules are prose. Reading any backticked identifier as a value would charge every declaration pointing at an ordinary rule. |
| `VLANV-X02` | A literal nobody declared is not charged. | Anchoring is for REPLICATED keys, and whoever replicates declares. Charging every value of every set was the first version, and it was replaced. |
| `VLANV-X03` <!-- @no-scenario: architectural boundary delegating comment syntax and token patterns to project dialect in anchors.yaml --> | Does not decide the anchor's syntax. | The project declares it (`derived.value_anchor`), because it belongs to how the codebase writes. The gate only requires that the anchor stand alone on its comment line (`VLANV-B08`). |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `ValueAnchor` | core — the declaration shape is declared by the project |
| DEP2 | `internal/mapx/model.go` | `KindCode` | core — the kind routes the jurisdiction |
| DEP3 | `internal/gate/spec_feature_match.go` | `defineRuleCaptureRE` | gate — the same reading of "the line that defines a rule" |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
