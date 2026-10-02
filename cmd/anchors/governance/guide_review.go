package governance

// reviewGuide é a régua de quem REVISA um PR. Ele existe porque o `claim` entrega cards
// em `ready-to-review` e não dizia o que conferir — o revisor caía no PR sem saber o que
// é dele e o que já foi medido por script.
//
// A divisão é a que separa as duas coisas: o que se confronta por SCRIPT roda como check
// do PR (mapa em dia, gates, cards declarados); o que exige JULGAMENTO é do revisor. Um
// review que repete o que o check já mediu gasta o revisor no que a máquina faz melhor —
// e deixa passar justamente o que só ele veria.
//
// O PASSO ZERO existe porque a primeira versão deste guia dizia "os checks já
// confrontaram" — e isso PRESSUPÕE que eles rodaram. Medido no PR #66 do app de referência: o
// evento `pull_request` não disparou (o PR estava com conflito, e o GitHub não roda
// workflow de PR em PR conflitado, sem avisar), nenhum check do Anchors existiu, e o PR
// ficou com a cara de PR limpo. O guia mandava o revisor NÃO conferir exatamente o que
// não havia sido conferido.
//
// O review é a última fronteira: ele é o que segura o que o pipeline deixou passar,
// inclusive o pipeline não ter acontecido. Por isso o guia começa perguntando se a
// fronteira anterior existiu — e diz como acioná-la, porque os workflows têm
// `workflow_dispatch` justamente para isso.
const reviewGuide = `# Review guide (the ruler for whoever reviews a PR)

## Step zero: do the checks EXIST?

Before any judgment, look at whether the checks ran. Not whether they passed — whether
they HAPPENED.

A check that did not run does not show up as a failure: it does not show up. The PR looks just
like an approved PR, and the platform does not distinguish "passed" from "never existed".

It happened: a PR with a merge conflict does not trigger a ` + "`pull_request`" + ` workflow on
GitHub, silently. No Anchors gate ran, and nothing flagged it.

    gh pr checks <n>

If the Anchors checks are not in the list, TRIGGER them before reviewing:

    gh workflow run anchors-gates.yml --ref <PR-branch>

If the PR has a conflict, the conflict is the first problem — resolve it first, because
while it exists the checks will not run on their own.

Do not approve a PR whose checks never ran. You are the boundary that remains when the
automation fails silently, and approving here certifies work that nobody confronted.

## What is NOT yours

Once the checks HAVE RUN, they have already confronted, by script:

- the committed map matches the repository;
- the blocking gates pass (` + "`anchors check --all`" + `);
- the cards the PR closes are declared in the body;
- the pipelines are up to date.

If any of them failed, the PR goes back to whoever opened it — it is not review work.
Rechecking that by hand spends you on what the machine does better.

## What IS yours

What requires JUDGMENT, and which therefore no script reaches:

### Does the spec decide what it needed to decide?

A gate checks that the spec has the sections and that the rules have codes. None checks
whether the rule ANSWERS the question the code will ask. A spec that says "the system
must be fast" passes every gate and decides nothing — and whoever implements it will
invent the threshold.

### Did the spec walk its variations, or stop at the first draft?

The gates prove that code, scenario and test agree with the spec — never that the spec
said enough. Most defects that reach users are a rule nobody wrote, and they come in
families: the case reported, and its siblings that arrive one by one. Read the spec as
` + "`anchors guide spec`" + ` ("Writing it") asks it to be written, and look for what is missing:
an input no rule mentions, an effect with no written limit, a variation of a rule — an
empty input, another system, another actor at the same time — that changes the behaviour
and is not a rule. Several rules saying one thing about different cases are an INVARIANT
waiting to be written.

### Does the rule state the INTENT, or the mechanism?

A rule that describes how the code does it ("the body ends at the first closer") is
proven by any test of that code, right or wrong: the unit stays green and wrong
together. The rule states what must hold ("the body is everything the test runs"). And a
FIX carries its rule: the defect is a rule that was not written, so the PR that fixes it
adds the rule — and the invariant it is one case of — with its scenario and its test.

### Does the code realize the rule, or only cite it?

The ` + "`regra-cumprida`" + ` gate already asks that of an AI, and its verdict is in the
map. Your job is different: to confront the JUDGMENT with the excerpt. A verdict of
"pass" over a marker in a generic place (top of the file, an import) is the defect the
automation gets wrong most often.

### Is each fix its own commit, and is each bug marked?

A correction folded into a ` + "`feat`" + ` disappears from the changelog. And a ` + "`Bug:`" + ` footer
is a claim that the defect shipped: on a correction of work that never left this PR it lists
a bug no user ever saw, and without it a real bug fix is lost among the fixes. The hook checks
the footer's format; only you know whether the defect shipped.

### What did the checks leave AWAITING JUDGMENT?

The pipeline does not judge, and that is deliberate: judgment asks for a credential,
answers slowly, and — what settles it — is JUDGMENT, exactly the work review exists to do.

So it DIRECTS. The check prints, in CI, the list of what is awaiting: the target, the
question the gate declares, and the guide holding the ruler. That list is yours to work
through, and it is the moment to do it — the PR has the diff in front of you and the
author is still there.

Measured in the reference app: 59 of 85 specs had never received a verdict, with CI
running on every PR. The earlier version printed only the count ("3 targets awaiting
judgment"), and a count is not an address — whoever reviewed knew there was work and not
which, where, or what to ask.

### What is marked TO REVIEW?

A gate that declares ` + "`review:`" + ` marks its targets to review, apart from its verdict: the
agent's judgment of what it wrote still finds real bugs, and a second look is another
thing. ` + "`anchors review --pending`" + ` lists them with the question; a target stays there until a
review is recorded at its current revision, and a change to it puts it back.

Record what you looked at, with who you are:

    anchors review <target> --gate <g> --by human:<you>
    anchors review <target> --gate <g> --by agent:<vendor>/<model> --findings "<report>"

The product is FINDINGS, not a stamp: each one — what, WHERE (file:line), why — goes to the
issue, and becomes a failing test and a fix, or is dismissed with its reason. There is no
` + "`waived`" + `: a review that did not happen is still due.

### Was what sits under ` + "`@TBD`" + ` judged, or rubber-stamped?

A spec that declares ` + "`@TBD: code`" + ` states that the code does not exist yet — and that is the
normal flow, because the spec is born first. But the judgment gates ask about
code, and faced with ` + "`@TBD`" + ` the question has no subject.

The temptation is to give a PASS to unblock. Such a PASS sits in the map looking like a real
verification, and it is worse than an open pending item: it states that someone looked.

Check what the verdict SAYS. An honest judgment about an absent target names the
absence ("there is no excerpt to confront: the spec declares @TBD and the code does not exist");
a stamp of convenience states that the code realizes the rule. The second is a finding.

And check whether the ` + "`@TBD`" + ` is true: if the code exists and the spec still declares it
as yet to be developed, the declaration has gone stale — and every gate that reads it starts
waiving what it should be demanding.

### Does the test PROVE, or only execute?

Coverage says the line ran. A test that calls the function and asserts nothing about
the result gives 100% coverage and proves nothing at all. Read the assertions, not the
number.

### What did the PR change without saying?

A diff the author does not mention in the description is where what they did not notice
changing lives. Compare the description with what was actually touched.

### Does the spec change change the DIRECTION?

If the PR alters a spec or plan, the revision (` + "`{CODIGO}-R000N`" + `) says what
changed and why. The question is whether that was a correction of FORM — and if it was not, whether
someone decided. The gate checks that the revision EXISTS; whether it is honest is for you to see.

And whom the revision FORGOT. A rule shares vocabulary with its siblings, and it is that
vocabulary a revision changes — not only the text of the rule it rewrites. The
` + "`revision-orphans`" + ` gate names the sibling that speaks of what was rewritten and was not
mentioned; declaring ` + "`Checked:`" + ` asserts that somebody read it, never that it is
correct. Whether reading actually happened, and what it found, is yours.

### Do two rules of the same unit contradict each other?

The gates confront each rule ON ITS OWN: the sections, the code, the scenario, the unit.
Two rules asserting opposite things both pass, because each one is well-formed.

Measured in the reference app, and it is what makes this yours: a spec whose title said
the choice lives on the DEVICE and whose body still said there is nowhere to store it; and
another where one rule listed the device token inside the backup scope and its sibling
named that very classification as the wrong one. Both had complete units and a green
suite. Both surfaced MONTHS later, as decisions that had to escalate.

When a revision is what changed the meaning, ` + "`revision-orphans`" + ` catches it. When there was never a
revision — the two rules were born contradicting each other — nothing does.

### Does the test say WHICH requirement it proves?

Semantic traceability runs on the scenario code appearing in a test case that passed. A
test that exercises the right behaviour without naming the code leaves the requirement
reported as unproven — and the next person writes the test that already exists.

In Go the identifier cannot carry the hyphen, so the code goes in a subtest name
(a subtest whose name opens with the code); elsewhere, in the case name itself.

## YOU DO NOT MOVE THE CARD

Approving, rejecting, commenting — yes. Touching the state label — no.

Who moves it is the pipeline, and it moves by FACT: ` + "`ready-to-review`" + ` when the checks
pass, ` + "`ready-to-test`" + ` when the PR is MERGED. A card in ` + "`ready-to-test`" + ` with the
PR still open tells the board that the work landed when it did not — and
` + "`ready-to-test`" + ` is the end of Anchors' jurisdiction, so nobody confronts it any more.

It happened: two independent reviews ran in parallel over the same PR. The
first approved and moved the card by hand. The second found a real defect the
first had not covered, and did not touch the state — "it is already ` + "`ready-to-test`" + ` from the other
review". The second acted rightly; the first created the fact that blocked it.

If your review rejects, the card stays where it is and the author sees the verdict. There is no
label to write: the wrong state is more expensive than the late state, because the
late one corrects itself at the pipeline's next event.

## YOUR VERDICT IS A LINE ON THE PR

The platform's "approve" does not reach the pipeline: the agents on one machine share an
account, and GitHub refuses an approval of your own PR. What the pipeline reads is a line
of its own, at the start of a line, in a comment on the PR:

    anchors-review: approved by <you>
    anchors-review: rejected by <you>

` + "`<you>`" + ` is the name in your ` + "`anchors-owner:`" + ` line on the card. The line sets the PR's
` + "`anchors/review`" + ` status: green on ` + "`approved`" + `, red on ` + "`rejected`" + `, and pending until
it is posted. A project that requires that status cannot merge before your verdict; one
that merges anyway gets a comment on the card saying the PR merged without the review
outcome, naming you.

Only the reviewer the claim assigned counts, and only a line posted after the assignment —
a verdict from an earlier round does not approve this one. Changed your mind? Post a new
line; the last one wins. Inside a code block the line is an example, not a verdict.

## Pontos de conformidade

The list exists so that nothing is missed by FORGETTING. Each point distils a question
the prose above already explains — read the section when a point is not obvious, because
the point is the reminder, not the argument.

They are NOT a substitute for the checks: everything under "What is NOT yours" the script
already confronted, and redoing it by hand spends you on what the machine does better.
These are the ones no script reaches.

### Before starting (tag: any PR)

- REV-CK1: the PR's checks RAN, and their absence was not read as approval — a PR with
  "no checks reported" has nothing confronted, and a green mark from a check that did not
  run is the worst possible pair
- REV-CK2: what is under "What is NOT yours" came back GREEN — if a blocking gate failed,
  the PR goes back to whoever opened it, and reviewing it by hand is spent effort

### The spec (tag: spec)

- REV-CK3: each rule ANSWERS the question the code will ask — a rule saying "the system
  must be fast" passes every gate and decides nothing, and whoever implements it will
  invent the threshold
- REV-CK4: the judgment of ` + "`regra-cumprida`" + ` was confronted against the EXCERPT it marks —
  a "pass" over a marker in a generic place (top of file, an import) is the defect the
  automation gets wrong most often
- REV-CK5: every target the checks listed as AWAITING JUDGMENT was judged — the pipeline
  does not judge (it needs a credential, answers slowly, and judgment is precisely the
  work review exists for); it prints the list, with the target, the question and the
  guide, and an unjudged target stays unjudged in silence
- REV-CK6: what sits under ` + "`@TBD`" + ` was judged as ABSENCE, not rubber-stamped as realized —
  and the ` + "`@TBD`" + ` is still true: if the code now exists and the spec still declares it, the
  declaration aged
- REV-CK7: a rule that contradicts a SIBLING rule of the same unit was reported — the gates
  confront each rule on its own, and two rules asserting opposite things both pass
- REV-CK16: the spec walked its variations — every input has a rule, every effect its written
  limit, and each variation that changes the behaviour (an empty input, another system,
  another actor at the same time) is a rule; several rules saying one thing are an invariant
- REV-CK17: each rule states the INTENT, not the mechanism — a rule that describes the code
  is proven by any test of that code, right or wrong
- REV-CK18: a fix carries its rule — the PR that fixes a defect adds the rule it broke, and
  the invariant it is one case of, with scenario and test, not the code change alone

### The revision (tag: spec, plan)

- REV-CK8: the revision ` + "`{CODIGO}-R000N`" + ` says what changed AND why — the gate checks it
  exists; whether it is honest is yours
- REV-CK9: a revision that changes DIRECTION was decided by whoever plans, not by whoever
  implements — a correction of form is the author's; a change of direction is not
- REV-CK10: the revision named every sibling rule that speaks of what it rewrote — either
  as revised too, or as read and still valid (` + "`Checked:`" + `)

### The test (tag: test)

- REV-CK11: the assertions confront the RESULT, not only that the line ran — a test that
  calls the function and asserts nothing about the outcome gives 100% coverage and proves
  nothing
- REV-CK12: the test names the scenario code it proves, so the map can see it — a test
  that proves a requirement without naming it leaves the requirement reported as unproven

### The PR itself (tag: any PR)

- REV-CK13: everything the diff touched is mentioned in the description — a change the
  author does not mention is where what they did not notice changing lives
- REV-CK14: the card's state label was NOT moved by hand — who moves it is the pipeline,
  and it moves by FACT; a wrong state costs more than a late one, because the late one
  corrects itself at the next event
- REV-CK15: the verdict was posted as the line ` + "`anchors-review: approved|rejected by <you>`" + `
  on the PR — a review that ends in prose alone leaves ` + "`anchors/review`" + ` pending, and the
  merge records the PR as merged without the review outcome

## When you finish

Approving is not "I found nothing": it is to state that you LOOKED at what is yours. If you looked and
there is nothing to say, approve — the review that only glances is worse than none, because
it creates the impression that someone checked.

Found something that does not belong to this PR? Record it instead of only commenting:

    anchors escalate "<what is wrong>" --about <file> --card <the PR's card>
`
