// @anchors
//   code: GPCGA
//   ref: GVGDG

package governance

// productGuide é a régua da DOUTRINA DE PRODUTO — a regra que atravessa alvos. Este guia
// embutido casa com a versão do binário; a IA o lê antes de escrever qualquer doutrina.
const productGuide = `# Product doctrine guide (the rule that cuts across targets)

Every spec in Anchors has a TARGET: it describes one unit, and co-location ties the four
ends of the unit to the same directory. That is the model's strength — every rule has an
address, and the gate knows where to confront it.

It is also its limit. In a real application many business rules **cut across targets**:
"the credit limit holds for signup, for simulation and for approval" belongs to none of
the three screens — it belongs to the product the three serve.

Without a place for it there were only two ways out, and both are bad:

- DUPLICATE the rule in all three specs, and accept they diverge at the first change;
- pick an ARBITRARY owner, leaving the other two referencing, implicitly, a rule that
  lives somewhere that is not theirs.

Product doctrine is that place.

## Where it lives

    product/<name>.doctrine.md

In the root, outside the target tree, on the same precedent as ` + "`plans/`" + `. The
` + "`product`" + ` kind comes from the PATH: a doctrine born elsewhere is read as a plain
doc, and the spec realizing it points at a file the map does not recognise.

## What it looks like

    <!-- @anchors
      code: LIMIT
      updated_at: 2026-01-15
    -->
    # Credit limit — what holds wherever money moves

    > **Code**: ` + "`LIMIT`" + `

    ## Overview

    Why this rule exists, and what it costs when it is broken. Write the REASON, not
    just the rule: whoever reads it later needs to know what was being defended.

    ## Rules

    ### LIMIT-R03 — the limit is never exceeded
    ### LIMIT-R04 — the limit is reassessed each cycle

The rule grammar is the SAME as the specs': a code, one of the three catalogued forms
(heading, table row, bold bullet), and a letter stating its nature.

## How a spec realizes it

With the ` + "`@realizes`" + ` TAG, on the line of the rule that concretises it:

    ### CRED-V01 — blocks submit above the limit    @realizes LIMIT-R03

The tag, and not a table column, because a catalogued rule has three valid forms and a
column exists only in one of them — with a column, referencing doctrine would force the
spec to change format.

**It is 1 to MANY**: several specs realize the same rule, and one spec realizes several.
That is exactly what the concept exists to allow, and it is why no gate here demands
parity of counts the way ` + "`unit-complete`" + ` does.

## What the spec writes, and what it does NOT

State HERE only what is specific to this unit, and let ` + "`@realizes`" + ` carry the
rest. Copying the doctrine's text back into the spec recreates the defect the whole axis
exists to eliminate: two texts saying the same thing diverge at the first change, and no
gate sees it because both are well-formed.

    WRONG   ### CRED-V01 — the credit limit is never exceeded   @realizes LIMIT-R03
    RIGHT   ### CRED-V01 — disables submit when the amount is above the limit   @realizes LIMIT-R03

## The gates

    plan-doctrine-exists      does the doctrine the PLAN cites exist?
    doctrine-realized         is the doctrine tied to any spec?
    spec-doctrine-exists      does the doctrine the SPEC cites exist?
    doctrine-not-duplicated   does the spec copy the doctrine's text?

` + "`doctrine-realized`" + ` INFORMS and does not fail: a rule decided today to be
implemented next cycle is legitimate work. But it does not go quiet either — a rule
written and never realized is a product decision that never reached the code.

## The way out: ` + "`@TBD`" + `, not ` + "`@no-*`" + `

    @TBD: <reason>       the piece does not exist YET — DEBT, stays visible
    @no-<thing>: <reason>  it will NEVER exist — permanent waiver

On this axis almost every case is the first: the cross-cutting rule is DECIDED before it
is implemented, and that is how product works. A bare marker waives nothing — the reason
is mandatory.

## What doctrine is NOT: the grouping tag

` + "`@feature: checkout`" + ` in a file header is a GROUPING label — it says what slice
the file belongs to, alongside ` + "`@experimental`" + ` and ` + "`@legacy`" + `. Nothing
reads it: no gate, no edge.

The two answer different questions. The tag groups the FILES of a slice; doctrine decides
the RULES that cut across units. A slice may have doctrine, doctrine may span slices, and
neither implies the other. The reasons are written down in the vertical-grouping doctrine
(` + "`VGRUP`" + `).

## Freezes: what is kept as it is, recorded here

A FREEZE keeps a part of the project as it is, by decision, until someone thaws it — and the
doctrine is where it is recorded, because it is a decision of the same kind. Each freeze is
an entry coded with the letter F, with its reason, under a section of its own:

    ## Freezes

    - **LIMIT-F01** — the plan mode ships next quarter

Every place it covers carries its code: ` + "`frozen: LIMIT-F01`" + ` in a header, ` + "`@frozen: LIMIT-F01`" + `
on a rule's or a heading's line. A frozen doctrine rule, section or document freezes the spec
rules that realize it. Write them with ` + "`anchors freeze --reason \"…\" --unit/--rule/--file`" + `, list them
with ` + "`anchors freeze show`" + `, and thaw with ` + "`anchors thaw LIMIT-F01`" + `; the gate
` + "`freeze-consistent`" + ` asks that every mark names a recorded freeze and every freeze marks
something. A freeze entry is not a rule to realize.

## The plan seeds doctrine too

A plan does not promise only the specs of the units; it promises the PRODUCT RULES they
will realize. Cite the path in backticks, like a spec:

    seeds ` + "`product/limit.doctrine.md`" + `

`
