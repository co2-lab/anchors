<!-- @anchors
  code: DCFRD
  updated_at: 2026-09-26
  layer: gate
-->
# DocsFresh — the compiled document has to reflect the spec

> **Code**: `DCFRD`

## Overview

Confronts the compiled documentation against the specs that feed it: **the page is there
and it reads well, but is it still saying what the spec says today?**

The documentation is GENERATED. The templates reference passages of the specs, and the
build command produces the pages. The content lives in the spec and only there — the copy
that lives in the compiled output is derived, and derived ages.

Without a gate the failure mode is silent and known: someone changes a rule in the spec,
forgets to recompile, and the documentation goes on asserting the OLD rule. Nobody sees it,
because the page is there, well formed, with real content. That is worse than an empty
page — an obviously incomplete document sends the reader looking for the source; an
out-of-date document convinces.

**INFORMATIVE by explicit decision**: a gate that blocks something a single command fixes
on its own spends a reviewer's attention on machine work. The pipeline runs the build on
merge; this gate exists so the author sees it before.

**Anchored on the SPEC** because the spec is what changes. The compiled document is not a
node of the map — it is generated, and charging a review of a compiler's output would be
charging a review of a compiler's output — so there is nowhere to start but the source.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the node kind | a spec, the only thing that feeds the templates | code, test, feature or guide | this unit: a kind outside the set leaves without a verdict |
| the template directory | an existing directory, or none | — (its absence is a case, not an error) | this unit: a project without the directory never opted into the mechanism |
| the map | a built graph | a nil graph, which the compiler cannot walk | the gate engine, which only routes this gate when the map is built |
| the specs on disk | what the map declares and the disk holds | a spec the map knows and the disk does not have | `internal/doct/doct.go`, which stops on the mismatch and hands back the error |

## Effects

| Effect | Description |
| --- | --- |
| `DCFRD-B01` | A compiled document that no longer matches what the templates would produce now FAILS — the spec changed and nobody recompiled. |
| `DCFRD-B02` | The verdict NAMES the stale documents and where they live, so the author does not have to diff the whole output directory. |
| `DCFRD-B03` | A compiled document that matches what the templates produce passes. |
| `DCFRD-B04` | An artifact that is not a spec leaves without a verdict: only the source of the documentation is confronted. |
| `DCFRD-B05` | A project with no template directory is not charged — it never opted into the compiled-documentation mechanism. |
| `DCFRD-B06` | A declared template whose document was never produced counts as stale by absence, not as satisfied. |
| `DCFRD-B07` | A document written by hand, with no generation marker, is NOT charged: the build refuses to overwrite it, and demanding a command that changes nothing is a warning nobody can act on. |
| `DCFRD-B08` | A template that cannot be compiled fails, carrying the compiler's own error — a broken template is a defect, not a reason to go quiet. |
| `DCFRD-B09` | A project whose specs cannot be read leaves without a verdict, carrying the reason: the gate could not look, and a gate that could not look must not approve. |
| `DCFRD-B10` | The comparison happens IN MEMORY. The gate never writes the compiled document. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DCFRD-I01` | The gate never repairs what it points at. If it compiled the document, the second run would always pass and the defect would only surface for whoever cloned the repository. | confronts a stale project and verifies that nothing was written to the output directory |
| `DCFRD-I02` | The charge starts from the SPEC and never from the compiled document. The document is a compiler output and not a node of the map; charging a review of it would charge a review of a build artifact. | confronts the spec and verifies the verdict lands there |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DCFRD-X01` | Does not judge whether the compiled document is GOOD. | The ruler is equality with what the templates produce now. Whether the template selects the right passages is a decision of whoever wrote it, and one this gate has no way to weigh. |
| `DCFRD-X02` | Does not run the build, even knowing the fix. | A gate that alters files makes the result depend on having run before. The check runs in a hook and in CI, and the second execution would always pass. |
| `DCFRD-X03` | Does not charge documents written by hand. | The build refuses to overwrite a page with no generation marker. Charging an update of a file the compiler will not write would send the author to run a command that changes nothing. |
| `DCFRD-X04` | Does not confront the compiled document as an artifact of its own. | It is generated output, not a declared unit. Giving it a node would demand a spec, a feature and a test for a build product. |

## Errors / Failures

| Rule | Condition | Effect |
| --- | --- | --- |
| `DCFRD-E01` | Underlying I/O or parsing failure | Returns Skip or Pending with error description | @no-scenario: error paths are handled by returning early verdict without panic @resilient: returns early without panic |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/doct/doct.go` | `Stale` | core — the comparison between disk and what the templates produce now belongs to the compiler, not to the gate |
| DEP2 | `internal/mapx/model.go` | `KindSpec` | core — the charge starts from the spec, which is the source that changes |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
