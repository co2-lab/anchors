<!-- @anchors
  code: VGRUP
  updated_at: 2026-09-20
-->
# Vertical grouping — the `@feature` tag, what it is for, and what it is not

> **Code**: `VGRUP`

## Overview

A repository is organised by LAYER: screens here, handlers there, models over there. That
is what the Structure declares, and it is what the map traverses.

But the work almost never arrives by layer. It arrives by SLICE: "checkout is broken",
"onboarding is being redone", "everything about billing needs review" — and the files that
answer for a slice are scattered across every layer.

`@feature: <name>` is the label for that slice. It says what a file belongs to, so a
person or a tool can gather the vertical cut without having to know, in advance, which
layers it touches.

**It is a GROUPING label, not a rule.** This distinction is the reason this doctrine
exists, and it was written after the tag was left, for a while, next to the product
doctrine axis without saying which was which — reading it, someone would take the tag for
the mechanism that verifies cross-cutting rules, write it, and wait for a gate that never
comes.

## Rules

### VGRUP-R01 — the tag says BELONGING, never a rule    @TBD: asserted by the header guides; no unit catalogues it as its own rule yet

`@feature: checkout` asserts "this file is part of the checkout slice". It asserts nothing
about behaviour, forbids nothing, and demands nothing. It is of the same nature as
`@experimental` and `@legacy`: a label whoever reads the file uses to know what it is near.

### VGRUP-R02 — a rule that cuts across units is DOCTRINE, never a tag    @TBD: the doctrine axis realizes it in practice; no spec declares it

When a rule holds for several units — "the credit limit applies to signup, simulation and
approval" — it belongs in `product/<name>.doctrine.md`, and the specs point at it with
`@realizes`. That axis has a file, an identity code, catalogued rules, map edges and
gates that confront them.

The two answer different questions, and confusing them is expensive: the tag groups the
FILES of a slice; the doctrine decides the RULES that cut across units. A slice may have
doctrine, doctrine may span slices, and neither implies the other.

### VGRUP-R03 — the tag is DECLARED but not yet confronted    @TBD: a statement of the current state — it stops being true the day the mechanism lands

Measured on 2026-09-20: `@feature` appears in the header guides and no parser reads it —
`internal/scan/` and `internal/mapx/` do not mention it. Node tags come from the layer
(`layers.tags`), never from the file header.

This is stated here rather than left implicit because a declared-and-unread tag is
indistinguishable, to whoever writes it, from one that works. Whoever uses `@feature`
today gets a label for people to read, and nothing more.

### VGRUP-R04 — the mechanism grows on FILTER and VIEW, never on charging    @TBD: the filter does not exist yet; this rule is the decision taken in advance

What the tag is worth developing for is gathering: `anchors check --feature checkout`,
`anchors map show --feature billing`, a slice view that crosses layers. That is what
turns a label into a tool.

What it must NOT become is a gate that demands the tag. A label that becomes an
obligation stops being an honest label: it gets filled in to silence the charge, and the
grouping fills up with files nobody actually classified — which is worse than no grouping,
because it looks like one.

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `VGRUP-X01` @TBD: holds today because nothing reads the tag | The tag does not define a layer, and never overrides one. | The layer comes from the Structure (the path, or `layer:` in the header). A slice cuts across layers by definition — letting it decide the layer would break exactly what it exists to represent. |
| `VGRUP-X02` @TBD: holds today because nothing reads the tag | The tag carries no rule, requirement or waiver. | Whatever needs to be confronted has a place: a rule in a spec, a cross-cutting rule in doctrine, a waiver in `@no-*` or `@TBD`. Smuggling any of those into a grouping label would hide them from every gate. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |
| `VGRUP-Q01` | Does the tag reach the map as a node tag, or does the filter read the header on demand? The first makes the slice traversable by the graph (and `--feature` becomes an ordinary filter); the second keeps the map free of something no gate confronts. The answer decides whether the grouping is part of the model or a convenience of the CLI. | whoever owns the map's surface | either a node-tag field fed by the header, or a scanner local to the filter |
