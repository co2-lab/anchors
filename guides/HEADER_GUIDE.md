# Header guide — project

> The block of markings at the top of EVERY file in this project. Seeded by
> `anchors init`; it is the built-in ruler (`anchors guide header`) instantiated for the
> stack here. Mandatory: a file without this header is invisible to what Anchors does
> best.

## The block, in this stack's dialect

In CODE/test/feature (it references the spec's unit):

```
// @anchors
//   ref: LGNN             # references the owning unit (the spec); it is NOT ownership
//   updated_at: 2026-08-08 # day of the last change (the gate checks it against git)
//   layer: screen         # the Structure's layer (usually deduced from the path)
//   @feature: auth
```

In the SPEC (the OWNER of the identity):

```
// @anchors
//   code: LGNN            # the spec OWNS the code
//   updated_at: 2026-08-08
```

## The markings

- `code:` — OWNERSHIP of the identity (the SPEC is the owner). `ref:` — REFERENCE
  (code/feature/test point at the spec's unit; it can be multiple: `ref: A, B`). Every
  file needs one of the two. Generate the code with `anchors code <name>`.
- `updated_at:` — the day of the last change. Whoever changes it updates it; the
  `updated-at-atual` gate checks against git (year-month-day only) and
  `anchors check --fix` corrects it. Do NOT invent the date — let it match the commit.
- `layer:` — the layer; usually deduced from the path, declare it only to override.
- `@feature: <name>` — a free GROUPING label for the vertical module. Nothing reads it:
  no gate, no edge. The vertical axis that IS confronted is product doctrine —
  `product/<name>.doctrine.md`, which the spec points at with `@realizes` (see
  `anchors guide product`).
- `@noPropagation`, `@anchors-shared-code` — honest opt-outs (always with the why
  alongside).

## Rules

- Always at the TOP of the file.
- `code` is the mandatory minimum (`header-valid` gate).
- `updated_at` matches the day of the last commit (`updated-at-atual` gate; `--fix`
  repairs it).
- An opt-out always carries a why alongside.

## Compliance points

- CK1: the header block is the FIRST thing in the file, before any other content — a
  header further down is not read, and the file is invisible to the map.
- CK2: the file declares `code:` or `ref:`, never neither — `code:` only in the spec that
  owns the identity, `ref:` in the code, feature and test that point at it.
- CK3: a file declares `code:` only if it is the spec — code, feature and test that
  declare `code:` claim an ownership that is not theirs, and two owners of one identity
  is an ambiguity the map cannot resolve.
- CK4: `updated_at` matches the day of the last commit that touched the file, and was not
  written by hand to a convenient date.
- CK5: every opt-out marking (`@noPropagation`, `@anchors-shared-code`) carries a written
  reason alongside — a bare marker silences the gate without saying why, which is the
  silence the framework exists to end.
- CK6: `layer:` appears only when it OVERRIDES what the path deduces — declaring the
  layer the path already gives adds a second source of truth for the same fact.

_(Full and universal ruler: `anchors guide header`.)_
