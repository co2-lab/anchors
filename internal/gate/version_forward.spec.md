<!-- @anchors
  code: VRFRA
  updated_at: 2026-11-18
  layer: gate
-->
# VersionForward — a change to a product's version moves it forward

> **Code**: `VRFRA`

## Overview

A product declares in `versioning:` (VRSNN) the files that hold its version. This gate reads the
version those sources hold now against the version they held at HEAD, by the product's scheme
(VRSNE). A commit that leaves the version alone passes: a branch keeps one version over many
commits. A commit that changes it must leave it readable, agreed between the sources and above the
one before: never lower, never a different spelling of the same rank.

Whether the version was already released elsewhere is not a question of the commit, which cannot see
the other branches nor what reached a store: the pipeline's stage asks that.

It is a question of the project's sources, never of one file: declared with `scope: project`, it
gives one verdict.

```yaml
gates:
  - name: version-forward
    check: version-forward
    scope: project
    on: [code]
    blocking: true
```

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the products | those of `versioning:` with at least one source | none with a source | this unit: skip (`VRFRA-B04`) |
| the sources now | the files as the tree holds them | a file missing or without its version | this unit: fail (`VRFRA-B02`) |
| the sources at HEAD | the files as the last commit holds them | a file not at HEAD, or without a version there | this unit: a version born now passes (`VRFRA-B03`) |

## Effects

| Effect | Description |
| --- | --- |
| `VRFRA-B01` | The gate passes when, for every product with a source, the version is the one HEAD holds or above it by the product's scheme. |
| `VRFRA-B02` | The gate fails naming the product when its version goes back, when it changes to a version of the same rank (`1.0.0+a` to `1.0.0+b`), when the scheme cannot read it, or when the sources disagree or cannot be read now. |
| `VRFRA-B03` | A version HEAD did not hold — the source is new, or held none there — passes: there is nothing before it to move forward from. |
| `VRFRA-B04` | With no product holding a source, the gate is skipped, saying there is no written version to confront. |

## Errors

none — what the gate cannot read is its verdict (`VRFRA-B02`), not a failure of its own.
