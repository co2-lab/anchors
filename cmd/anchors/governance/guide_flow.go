package governance

// flowGuide é a régua do TRABALHO DIRIGIDO POR FLUXO — ações como peças, fluxos como
// montagem. Este guia embutido casa com a versão do binário.
const flowGuide = `# Flow guide (work driven by shape, not by memory)

A document that says "what to do in case X" is INPUT: it requires whoever works to
REMEMBER to consult it, CHOOSE right among the options, and NOT SKIP a step. Three
chances to err per round, forever.

The problem was measured in a real project: a 232-line guide cataloguing eight
approaches for getting unstuck, with the two-attempts rule written on line 11 — and a
dossier that opens by warning "five passes, each one re-aimed the previous one's
target... two already sent a round to the wrong site".

A FLOW inverts the burden. Instead of "here are eight approaches, choose", it answers
"from here, the valid exits are these three". The rule stops depending on memory and
becomes the only way out.

It is the same inversion the gates already perform on the artifact ("remember to write
the test" became 'triad-complete'), applied to the PROCESS.

## Two artifacts, and the difference is the whole idea

  flows/actions/<name>.action.md    an ACTION — the puzzle piece
  flows/<name>.flow.md              a FLOW — the assembly

An ACTION declares WHAT IT DOES and WHICH RESULTS it offers. It deliberately does not
know who comes next, and that is what makes it reusable: 'map build' is a step of the
worker cycle AND of adoption, written once.

A FLOW does not redraw the work — it FITS actions together and says where each result
goes. What it adds is the LINK, and that is where rules stop depending on memory.

## Writing an ACTION

    <!-- @anchors
      code: ACHCK
      updated_at: 2026-01-15
    -->
    # Action: ` + "`anchors check`" + ` — confront what was written

    > **Code**: ` + "`ACHCK`" + `

    ## Command

        anchors check --changed <file>

    ## Results

    ### ACHCK-R01 — PROMOTABLE: no blocking gate failed

    ### ACHCK-R02 — BARRED: a blocking gate failed

Each result gets a code (` + "`-R01`" + `, ` + "`-R02`" + `…) and a SHORT NAME in caps
before the colon — that name is what the arrow carries in the diagram, so it has to be
readable on its own.

**RUN THE COMMAND before writing the results.** Measured while writing the first action:
'check' answers "is not governed by the Structure — nothing to confront", a fourth
outcome that neither passes nor fails, and that nobody would have written from memory.
An action whose results were imagined describes a command that does not exist.

## Writing a FLOW

    ### WORKR-P04 — confront

    Fits: ` + "`ACHCK`" + `

    Results:
    - ` + "`ACHCK-R01`" + ` PROMOTABLE → ` + "`WORKR-P06`" + `
    - ` + "`ACHCK-R02`" + ` BARRED → ` + "`WORKR-P03`" + ` (back to writing)

The ` + "`Fits:`" + ` line names the piece. Each result line routes one outcome to the next step.
The association is POSITIONAL: what comes below a step belongs to it.

A step with no exit declares ` + "`> @terminal`" + `. It is DECLARED and never deduced
from "has no exit": a step with no exit may be the end of the work or an oversight, and
only the author knows which.

## What the shape guarantees

A rule that used to be prose becomes topology:

  "NEVER close the task with a blocking gate red"
      → the BARRED result has no link to 'done'. What has no door is not skipped.

  "approve the scope BEFORE writing the plan file"
      → the draft step has no exit to the file. The only path passes through approval.

  "do NOT judge by the whole prose by eye"
      → reading the guide is the only path to the verdict.

## The commands

    anchors flow build                    scan flows/ and write the graph into the map
    anchors flow show <flow>              draw it (--mermaid for boxes and arrows)
    anchors flow next <step>              the VALID exits from here — the one that drives

'flow build' reports the RESULTS no flow handles. It is the finding the puzzle shape
makes possible and a linear drawing hides: an action declares everything it can answer,
and a result nobody routes is a hole — whoever gets it improvises, which is the whole
problem.

It is a FINDING and not a failure: a project may legitimately not cover every branch.
What it must not do is not know.

## Where the flow lives

In the map ('anchors.graph.yaml'), under its own 'flow:' key — one graph, two contents
that do not mix. The map links ARTIFACTS and is what 'impact' walks; a process state is
nobody's artifact, and crossing it would have 'impact' answer that changing a file
affects a decision.

The DIAGRAM is derived at runtime, never stored: a versioned diagram ages against the
flow it describes, and this repository already measured what that costs.
`
