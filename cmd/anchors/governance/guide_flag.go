package governance

// flagGuide is the ruler for FEATURE FLAGS — the scenarios a flag's value opens. It lives
// beside `guide_product.go` because the two axes are the same shape: something that lives
// outside the target tree, which the spec points at from its own rule line.
const flagGuide = `# Feature flag guide (the scenarios a flag's value opens)

## The problem

A feature flag multiplies the code's paths without multiplying the spec.
` + "`if flag(\"new-checkout\")`" + ` creates two behaviours, and the spec describes ONE — or,
worse, describes both mixed into a sentence that does not say which holds when.

The cost shows up in three places, and all three are silences no other gate sees:

- **in review** — nobody knows which branch is current and which is going to die;
- **in the test** — the ON path gets covered and the OFF path gets no proof;
- **in removal** — the "temporary" flag turns three years old, and deleting it becomes
  archaeology.

## Where it lives

    flags/<name>.flag.md

Outside the target tree, on the same precedent as ` + "`product/`" + `: a flag governs rules
across several units, so it belongs to none of them.

## The shape

    <!-- @anchors
      code: CHKUT
    -->
    # Flag: new-checkout

    ## Scenarios

    | Scenario | When the value | Then |
    | --- | --- | --- |
    | ` + "`CHKUT-G01`" + ` | ` + "`= \"off\"`" + ` | the old checkout answers, and the new one is not called |
    | ` + "`CHKUT-G02`" + ` | ` + "`= \"on\"`" + ` | the new checkout answers, with the same contract |
    | ` + "`CHKUT-G03`" + ` | ` + "`>= 50`" + ` | the rollout percentage decides per user, stably |
    | ` + "`CHKUT-G04`" + ` | ` + "`absent`" + ` | the declared default holds — never an error |

The letter is ` + "`G`" + `, and each scenario is a rule code like any other: it can be cited,
tested and tracked.

## The ABSENT case

` + "`CHKUT-G04`" + ` is the one most often forgotten and the one that breaks hardest. A flag
that does not answer — flag service down, a new environment, a local test — has to have a
declared behaviour. Without it, what holds is whatever the library defaults to, and that
is discovered in production.

Measured against the market's own tools: every one of them can EXPRESS this case
(Flagsmith calls it ` + "`Is Not Set`" + `), and not one of them REQUIRES it. The
` + "`flag-scenarios-complete`" + ` gate is what closes that gap.

A flag that genuinely cannot be absent (read from a local constant, say) waives it:

    @no-absent: read from a local constant, never missing

` + "`@no-absent`" + ` and not ` + "`@TBD`" + `, because the two assert different things: this
one will NEVER have the case, which is a permanent waiver and not a debt.

## The grammar of the condition

FIXED, and deliberately so: prose covers any case and lets nothing be confronted. The
operator set is the union of what the market's tools actually offer — LaunchDarkly,
Unleash and Flagsmith — collapsed onto one operator per MEANING, because the three
disagree only on spelling:

    =  !=  >  >=  <  <=          (` + "`eq`, `NUM_GTE`, `Greater Than Inclusive`…)" + `
    contains  not contains
    starts with  ends with
    matches                       (regex)
    in  not in
    before  after                 (dates)
    rollout                       (percentage)
    absent  present               (` + "`Is Set` / `Is Not Set`" + `)

Whichever spelling your tool taught you is accepted; they all arrive as the same
statement. What is NOT accepted is prose, and the refusal is the point.

Two operators are deliberately absent: ` + "`segmentMatch`" + ` (LaunchDarkly) and
` + "`Modulo`" + ` (Flagsmith). Both evaluate against a service Anchors has no access to and
must not pretend to have — a segment lives in LaunchDarkly's database, not in the
repository. Write the scenario as the VALUE it produces, which is what the code branches
on anyway.

## Tying a rule to a scenario

The SPEC declares that its rule only holds under a scenario:

    ### CRED-V01 — validates the limit before submitting   @gated-by CHKUT-G02

The direction is the same as ` + "`@realizes`" + `, and for the same reason: whoever knows
about the dependency is the spec. A flag file listing its dependents would be an index,
and an index is wrong as of the next spec somebody writes without updating it.

## The gates

| gate | question | verdict |
| --- | --- | --- |
| ` + "`flag-scenario-grammar`" + ` | is the condition written in the grammar? | fails |
| ` + "`flag-scenarios-complete`" + ` | does the flag declare the ABSENT case? | fails |
| ` + "`flag-scenario-exists`" + ` | does the cited scenario exist? | fails |
| ` + "`flag-covered`" + ` | does EVERY scenario have a green test? | informs |

` + "`flag-covered`" + ` charges per SCENARIO, not per flag: the flag whose ON path is tested
and whose OFF path is not passes a per-flag rule while leaving exactly the branch that
will break unproven. It informs rather than blocks because the scenario written today and
tested next commit is ordinary work, not a defect.

## What Anchors does NOT do

- **Read the flag's real value.** It has no access to the flag service, must not have, and
  the value changes per user and per minute. What it confronts is the DECLARED scenario.
- **Dictate the library.** LaunchDarkly, Unleash, an ` + "`if`" + ` in a config file — Anchors
  recognises the declared scenario, not the call.
`
