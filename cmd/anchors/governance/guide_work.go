package governance

// workGuide é a régua de quem PEGA UM CARD. Ele responde ao que o card não cabe dizer
// sem se repetir em cada issue — e o custo da repetição não é só ruído: instrução copiada
// CONGELA. Um card criado hoje carrega o texto de hoje, e quando a instrução muda, os
// cards antigos passam a ensinar o errado sem que nada acuse.
//
// Centralizado no binário, o guia acompanha a versão: quem roda `anchors guide work` lê a
// instrução ATUAL, não a do dia em que o card nasceu.
const WorkGuide = `# Work guide (the ruler for whoever picked up a card)

## Before starting: CLAIM the card

` + "`anchors next`" + ` is the first command, and it is not bureaucracy — it is what
records WHO has the work and moves the card to the right column.

    $ export ANCHORS_SESSION=dev1   # a name of its own for EACH agent
    $ anchors next
    · card #223 is yours — moved to in-progress

**Declare ` + "`ANCHORS_SESSION`" + ` first.** The claim records the owner as
` + "`<machine>/<session>`" + `, and ` + "`next`" + ` refuses to claim without it: two agents of the same
user on one machine would otherwise be ONE owner to the board and take each other's cards.

**Picking the card by hand and starting to implement does not work.** What is lost is not the
record: it is the QUEUE. The claim serves ` + "`ready-to-review`" + ` BEFORE ` + "`to-do`" + `
— reviewing takes priority over new work. A card that never enters the review
column makes the next agent find it empty, and they pick up new work instead of
reviewing what is ready.

Measured in the reference project: 38 of 44 PRs opened with the card still in
` + "`to-do`" + `, and the review queue showing ONE item while 46 pieces of work
waited. The ` + "`gates`" + ` gate fails a PR whose card stayed in ` + "`to-do`" + `.

## Before starting

` + "`anchors status`" + ` says where the project is, and ` + "`anchors guide <artifact>`" + ` says
how to write what you are going to write (spec, code, feature, test).

## The order: from RIGHT to LEFT

Finish what is furthest along before picking up something new. Half-finished work delivers
nothing and ages until the context is lost — whoever resumes it pays the cost of
understanding again, and sometimes discovers the decision changed along the way.

## Found something wrong that is NOT this card?

It happens all the time: a config that contradicts the project's doctrine, a path
nobody documented, a file in the wrong place. **No gate sees that** — the gate opens an
issue for what IT detects, and the rest depends on you.

The cheap path is to fix it on the spot and move on. And then the fix vanishes from the history: whoever
comes later does not know that was ever a problem, nor why the solution is that one.

    anchors escalate "<what is wrong>" --about <file> --card <this card>

The finding is born with the label ` + "`anchors:under-<number>`" + `, and the two are delivered in the SAME
PR: you already have the context in hand, and separating them would make one of the two wait for no reason.

### When NOT to use it

If the fix is trivial AND it is in the file you are already editing, fix it and record the
revision in the file itself (` + "`{CODIGO}-R0001: o que mudou e por quê`" + `). Opening a card
to change one word is bureaucracy.

### READ THE OPEN DECISIONS BEFORE OPENING ANOTHER

You are not the only one working. Before escalating, look at what is already awaiting an answer:

    gh issue list --label anchors:needs-user --state open --limit 100

If your discovery is **the same subject** as one that already exists, comment on it. The discovery
is preserved all the same, and whoever decides gets one question with more evidence instead of two
similar questions.

If it is a **different subject**, escalate. Two decisions about the same file do exist — one
about the value scale, another about where the preference lives are not the same thing.

WHY THIS IS YOURS AND NOT THE COMMAND'S: only someone who read both texts knows whether they are the same
question. Anchors could compare the ` + "`--about`" + `, or match words in the title — and it would be wrong
in both: the same file does not mean the same subject, and two agents describe the same
finding with different words (in different languages, as the case may be). The judgment is the
work; automating it would produce a wrong warning wearing the face of a ruler.

MEASURED in the reference project: 11 escalations in one hour, 6 in the next. A triage
regrouped the 23 open ones and concluded they were **seven decisions** — the rest was the same
subject, seen from different angles. Nobody decides by reading 23 reports to find 7
questions, and while the queue grows like that, your escalation waits too.

#### If the subject is the same, TIE it instead of only commenting

Commenting preserves the discovery. Tying makes one answer unblock them all at once:

    gh issue edit <your card> --add-label anchors:under-<the decision that already exists>

The label is filterable, so whoever decides sees the question's real weight (the filter
` + "`--label anchors:under-<n>`" + ` lists everything it holds), and whoever answers knows exactly what they
released. Without the tie, each card has to be found again and unblocked one by one — and what
is not found again stays stuck after the decision has already come out.

If the change **impacts the project's direction** — or if you are in doubt —, do not make it:

    anchors escalate "<what needs to change>" --about <file> --for-user

That becomes a decision for whoever planned, and the card stops until it comes out. Interpreting the impact
is yours: you are the one with the context of what you found.

If what is wrong is the PIPELINE or the TOOL — a seeded workflow computes the wrong thing, an
` + "`anchors`" + ` command picks the wrong card, a gate misreads a file — it is not a decision: there is nothing to
choose, and the fix lives where you do not edit (` + "`.github/workflows/anchors-*`" + `, Anchors itself):

    anchors escalate "<what is wrong, with the measurement>" --about <file> --bug [--blocking]

` + "`--blocking`" + ` when your card cannot go on until it is fixed; without it the card goes on.

## Before committing

` + "`anchors check`" + ` blocks while there is a pending judgment — it is work for THIS
commit, and whoever touched the file is the one with the context to answer:

    anchors judge --pending
    anchors judge <target> --gate <g> --verdict pass|fail --reason "..."

## When opening the PR: Anchors writes the closing lines

    anchors pr-body

It prints the lines that close the card you claimed AND the findings born under
it (` + "`anchors:under-<n>`" + `). Paste them into the PR body.

You do NOT need to know the platform's syntax — and it is precisely that which gets wrong
silently. GitHub recognizes the closing word only in ENGLISH: writing "Fecha
#44" in a Portuguese project is ignored with no error at all, the PR merges, and the card stays
open. It happened: a PR said "Fecha #44, #49, #50" and all three stayed open.

The link is declared in Anchors' vocabulary — the card you claimed and the findings
under it. The word the platform understands is DERIVED from that, and changes with the
platform, not with the project's language.

## YOU DO NOT CLOSE THE CARD

Who closes it is the MERGE, through the ` + "`Closes #N`" + ` line that ` + "`anchors pr-body`" + ` wrote. Closing by
hand looks like tidying — the work is ready, the PR is open, the card is "done for" — and
it breaks the flow in a way that does not show:

- the card leaves ` + "`ready-to-review`" + ` BEFORE anyone reviews, and the review queue empties;
- the claim looks for a review first (that is the right-to-left order) and finds
  nothing, so it hands out NEW work;
- the PR keeps waiting, and the longer it waits the more expensive it gets to review.

Measured: a project accumulated 25 green PRs awaiting review with ZERO cards in
` + "`ready-to-review`" + `. The cards had been closed one minute BEFORE the PR was opened, and the
board said there was nothing to review while 25 pieces of work waited.

The rule is the same as for the state: the card moves by FACT, and the fact of "delivered" is the merge.

## After opening the PR: the card is not delivered

Opening the PR does not close the card, and **pushing a commit is not a stopping point**. The CI runs
after the push, and its result is part of your work — not news somebody brings you.

    anchors work review          (shows what is missing in the PR that is yours)
    gh pr checks <n> --watch     (BLOCKS until the CI finishes, and returns the verdict)

` + "`--watch`" + ` exists for that: it waits for the process, you do not. Without it, the only way
to know the result is to ask again later — and "awaiting the new round" is a
sentence that ends the turn without delivering anything. The card stays ` + "`in-progress`" + `, with your name
on it, and whoever resumes it pays the cost of understanding all over again.

### The CI failed

Red is work for THIS card, not a new card. Read the failure, fix it, push and
**wait again** — as many rounds as needed. There are only two legitimate exits from the
cycle:

- the CI went green and the card advanced on the board; or
- the failure demands a decision that is not yours, and then it becomes an escalation
  (` + "`anchors escalate ... --for-user`" + `), with the card explicitly stopped.

A third thing is NOT an exit: reporting the diagnosis and stopping. The correct diagnosis is
half the work; the other half is the verdict of the check you triggered.

## When finishing the round: the REPORT has a format

Whoever reads your report decides whether to continue, to review, or to answer a question. And what
decides that is not the narrative of what you did — it is the state: where the card is, whether the CI's
verdict was read, and what is waiting on a person.

    anchors task-status

It discovers what the machine knows (the card and its state, the PR and the checks, what was not
pushed, the stopped decisions) and leaves TWO gaps, which are yours:

- **What I proved** — the rules the suite confronts, and what the mutation killed. Tests that
  pass are not proof; proof is the mutant that died.
- **What was left out** — nothing, or what you left and why. Reduced scope is the decision
  of whoever asked, not yours: if something did not go in, this is where they find out.

Without a format, each round reports what the agent found important — and what gets omitted first
is precisely the state. A report that says "I fixed it and pushed, awaiting the new round"
is CORRECT and insufficient: it does not say the card stayed ` + "`in-progress`" + ` with your name, nor that
the verdict of the check you triggered was not read by anyone.

## The spec is born before the code

That is the normal flow: the spec is the anchor. While the pieces do not exist, declare what is missing:

    > **@TBD: code,feature,test** — the pieces are the phase in progress.

` + "`@TBD`" + ` is *to be developed*, and it expires on its own: when the piece appears in the map, the gate
goes back to confronting it. It is different from ` + "`@no-test`" + `, which states "this unit does NOT
NEED a test" — permanent, and which would erase the demand forever.
`
