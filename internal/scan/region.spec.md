<!-- @anchors
  code: SRRGS
  updated_at: 2026-09-26
  layer: scan
-->
# SourceRegion — the declared interval of a rule code in source, and the composition of test scripts

> **Code**: `SRRGS`

## Overview

The other forms of a unit's identity are delimited by the syntax they live in: the rule's line in the
spec, the scenario in the feature, the case in the test. In source code a code marker said where a
requirement starts and nothing said where it ends, so a file with fourteen requirements answered
fourteen times the question "which requirement changed". A region closes that gap: an opening marker
`#region [CODE]` and a closing marker `#endregion [CODE]`, the pair editors already fold, delimit the
interval, and each region gets its own content revision.

The closing marker repeats the code because with an anonymous close a swapped nesting is
undetectable — the counts of opens and closes match and the measured interval is the wrong one, in
silence. With the code, the close is confronted name against name. The unit extracts the regions of a
text and the pairing defects, which the region pairing gate reports. Regions are optional: a file with
none is valid and keeps its whole-file revision.

The unit also reads the COMPOSITION of a test script: the other scripts it runs (`runFlow:` and its
`file:` form), resolved relative to the script's own directory. This is what makes a shared utility
script a dependency of the hundreds of scripts that run it.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the text scanned for regions | any text, in any language and comment syntax; markers are recognised by the `#region` / `#endregion` annotation, not by the comment prefix | — | this unit: a text without markers yields nothing |
| the code between brackets | a rule code of the configured code length followed by `-` and a rule suffix | a marker without a bracketed code, which does not open a region | this unit, through the marker grammar |
| the script content | YAML text of a test script | — | this unit: text without a composition yields nothing |
| the script path | the script's root-relative, slash-separated path | — | the scan passes the path it walked |

## Effects

| Effect | Description |
| --- | --- |
| `SRRGS-B01` | `Regioes`: Regions pair as a stack: an inner region closes before the outer one, comes out first, and each carries the lines of its opening and closing markers. |
| `SRRGS-B02` | A region spans exactly the lines its markers enclose, however deep its indentation and however long the rest of the file. |
| `SRRGS-B03` | A region's revision changes when a line inside it (markers included) changes, and does not change when a line outside it changes. |
| `SRRGS-B04` | A close that names another code is reported and still closes the open region, so the regions after it pair normally and no cascade of derived defects follows. |
| `SRRGS-B05` | A close without a code closes the region on top of the stack without a defect. |
| `SRRGS-B06` | A text with no region markers yields no region and no defect. |
| `SRRGS-B07` | `ComposeRefs`: The composed scripts are the `runFlow:` and `file:` targets ending in `.yaml` or `.yml`, resolved relative to the script's directory; `runScript:` targets are not composition. |
| `SRRGS-B08` | A composed script cited more than once is listed once, and a reference to the script itself is not a dependency. |
| `SRRGS-B09` | A script with no composition yields no dependency at all. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `SRRGS-I01` | One pairing defect produces exactly one reported defect, whatever follows it. | a swapped close followed by a well-formed region yields one defect and both regions |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `SRRGS-X01` | Composition is read only from the literal paths the script declares; nothing is inferred. | The edge must be declared to be trusted; an inferred composition would create dependencies nobody wrote. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `SRRGS-E01` | A region is opened and never closed. | A defect of kind `sem-fecho` naming the opening line and code. | An open interval has no end, and a revision measured to the end of the file would be the wrong one. |
| `SRRGS-E02` | A close appears with no open region. | A defect of kind `fecho-orfao` naming the line and the code it carries. | The close delimits nothing, and ignoring it would hide a marker the author thought was working. |
| `SRRGS-E03` | A close names a code other than the open region's. | A defect of kind `fecho-trocado` naming the line, the open code and the code found. | A swapped nesting measures the wrong interval in silence; the name-against-name confrontation is why the close repeats the code. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `CodeLengthPattern` | config — the configured code length in the marker grammar |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
