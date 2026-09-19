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
// confrontaram" — e isso PRESSUPÕE que eles rodaram. Medido no PR #66 do blue-eyes: o
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

### Does the code realize the rule, or only cite it?

The ` + "`regra-cumprida`" + ` gate already asks that of an AI, and its verdict is in the
map. Your job is different: to confront the JUDGMENT with the excerpt. A verdict of
"pass" over a marker in a generic place (top of the file, an import) is the defect the
automation gets wrong most often.

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

## When you finish

Approving is not "I found nothing": it is to state that you LOOKED at what is yours. If you looked and
there is nothing to say, approve — the review that only glances is worse than none, because
it creates the impression that someone checked.

Found something that does not belong to this PR? Record it instead of only commenting:

    anchors escalate "<what is wrong>" --about <file> --card <the PR's card>
`
