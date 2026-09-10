---
title: The workflow
description: How a working day happens with Anchors — from claiming work to merge.
---

The pillars say **what** Anchors stands for. This page says **what you do**, in
order, on a normal working day.

It applies whether you work alone or in a team. Where the two differ, it's
marked.

## The board is the state

All work lives in a **card** — an issue with the `anchors` label. Its state is
another label, and the sequence is always the same:

```
to-do → in-progress → ready-to-review → in-review → ready-to-test → ...
└────────── agent's remit ──────────────┘ └──── delivery pipelines' remit ────┘
```

Two things matter about that split:

- **up to `ready-to-review`, whoever works moves the card.** It's yours while
  you hold it.
- **from `ready-to-test` onward, the delivery pipelines move it.** Anchors goes
  no further — it knows nothing about your release process.

Beyond state, two labels cut across any column:

| label | meaning |
| --- | --- |
| `anchors:precisa-do-usuario` | the card awaits a **human decision** and is handed to no one until it comes |
| `anchors:sob-<n>` | this card is a **finding** born while someone worked on card `<n>` |

## 1. Ask for work — don't pick it

```sh
gh workflow run anchors-claim.yml -f agent=<machine>/<session>
```

You **ask**, the pipeline **decides**. This isn't ceremony: it's what stops two
agents from taking the same card.

The reason is technical and has no way around it. If each one claimed directly,
two could read "unowned" before either wrote, and both would believe they own it
— GitHub's API offers no compare-and-swap. The claim pipeline is serialized, so
"is there a free card?" and "assign it" happen with nobody in between.

**Priority runs right to left.** Finish what's furthest along before starting
something new: half-done work delivers nothing, occupies a reviewer, and ages
until the context of whoever wrote it is lost.

And **your** card comes first. If you already have one in flight, the claim hands
it back instead of giving you another.

## 2. Read the rule before writing

```sh
anchors status          # where the project is, and the next step
anchors guide work      # the rule for whoever took a card
anchors guide spec      # (or code/feature/test) how to write the artifact
```

`guide work` answers three things you'll need: the board's order, what to do with
a finding that is **not** your card, and what `anchors check` demands before the
commit.

Don't memorize it — the guide lives in the binary precisely so you don't have to.

## 3. Work, then confront before committing

```sh
anchors check --changed <file>
```

`check` runs the gates over what you touched. It has three outcomes, and telling
them apart saves a lot of frustration:

| symbol | means | blocks? |
| --- | --- | --- |
| `✓` | passed | — |
| `✗ [BLOQUEIA]` | a blocking gate failed | **yes** |
| `✗ [informativo]` | found something that doesn't stop delivery | no |
| `~` | **indeterminate** — the gate had nothing to confront | no |

`~` is the most misunderstood. It isn't failure: it's the gate saying the
artifact on the other side of the relation doesn't exist. A coverage gate on a
file with no test doesn't fail — it doesn't measure.

If a blocking gate fails and the failure is **deliberate**, declare it in the
commit message:

```
[skip-<rule>@<CODE>: why the failure is acceptable here]
```

This doesn't turn the gate off: it records the waiver with the reason, where
whoever comes next will read it.

## 4. Found something that isn't this card?

It happens constantly. A config contradicting the doctrine, a path nobody
documented, a file in the wrong place.

**No gate sees this** — a gate opens issues for what IT detects; the rest is on
you.

```sh
anchors escalate "<what's wrong>" --sobre <file> --card <this card>
```

The finding is born with the `anchors:sob-<n>` label, and both ship in the
**same PR**: you already have the context in hand.

### When NOT to escalate

If the fix is trivial **and** it's in the file you're already editing, fix it and
record the revision in the file itself (`{CODE}-R0001: what changed and why`).
Opening a card to change one word is bureaucracy.

### When the decision isn't yours

If the change **affects the project's direction** — or if you're unsure:

```sh
anchors escalate "<what must change>" --sobre <file> --para-usuario
```

That becomes a decision for whoever planned it, and **the card stops until it
comes**. Judging the impact is your call: you're the one holding the context of
what you found.

## 5. Open the PR — without inventing the syntax

```sh
anchors pr-body
```

It prints the lines that close the card you took **and** the findings born under
it. Paste them into the PR body.

You don't need to know the platform's syntax — and that's exactly what gets
gotten wrong in silence. GitHub only recognizes the closing keyword in
**English**: writing "Fecha #44" in a Portuguese project is ignored with no
error, the PR merges, and the card stays open.

This actually happened: a PR said "Fecha #44, #49, #50" and all three stayed
open.

## 5.1 Wait for the verdict — pushing is not delivering

Opening the PR does not close the card, and **pushing a commit is not a stopping
point**. CI runs after the push, and its result is part of your work.

```sh
gh pr checks <n> --watch
```

`--watch` **blocks** until CI finishes: it waits on the process so you don't have
to. Without it, the only way to learn the result is to ask again later — and
"waiting for the next run" is a sentence that ends the turn delivering nothing.
The card stays `in-progress`, with your name on it.

If CI fails, the red is work on **this** card, not a new one: read the failure,
fix it, push, and wait again. There are exactly two ways out of the loop:

- CI went green and the card moved on the board; or
- the failure needs a decision that isn't yours, and it becomes an escalation
  (`anchors escalate ... --for-user`), with the card explicitly parked.

Reporting the diagnosis and stopping is **not** a third way out. The right
diagnosis is half the work; the other half is the verdict of the check you fired.

## 5.2 The round's report has a format

```sh
anchors task-status
```

Whoever reads your report decides whether to continue, to review, or to answer a
question — and what decides that is not the narrative of what you did. It's the
**state**: where the card is, whether the CI verdict was read, and what is waiting
on a person.

The command discovers what the machine knows (the card and its state, the PR and
its checks, what hasn't been pushed, the decisions parked in `needs-user`) and
leaves **two gaps**, which are yours:

- **What I proved** — the rules the suite confronts, and what mutation killed.
  Passing tests are not proof; proof is the mutation that died.
- **What I left out** — nothing, or what you left and why. Cutting scope is the
  requester's call: if something didn't make it, this is where they find out.

Without a format, each round reports whatever the agent found important, and what
gets omitted first is precisely the state.

## 6. Review

The claim hands out `ready-to-review` cards **before** `to-do` ones — reviewing
comes before starting something new.

```sh
anchors guide review
```

### Step zero: do the checks EXIST?

Before any judgment:

```sh
gh pr checks <n>
```

A check that didn't run **doesn't show up as a failure — it doesn't show up**.
The PR looks exactly like an approved one, and the platform doesn't distinguish
"passed" from "never existed".

It happened: a PR with a merge conflict doesn't trigger `pull_request` workflows
on GitHub, silently. No gate ran, and nothing flagged it.

If the Anchors checks aren't listed, trigger them before reviewing:

```sh
gh workflow run anchors-gates.yml --ref <PR-branch>
```

**Don't approve a PR whose checks never ran.** You are the boundary that remains
when automation fails silently.

### What's yours, and what isn't

What the script already confronted is **not** review work: the map being current,
the blocking gates, the cards declared in the body. Re-checking by hand spends
you on what the machine does better.

What's yours is what requires judgment:

- **Does the spec decide what needed deciding?** A spec saying "the system should
  be fast" passes every gate and decides nothing.
- **Does the code fulfill the rule, or merely cite it?** A `pass` verdict over a
  marker in a generic place (file top, import) is the mistake automation makes
  most often.
- **Does the test PROVE, or just execute?** Read the assertions, not the number.
- **What did the PR change without saying?** A diff the author doesn't mention is
  where what they didn't notice lives.

### The three outcomes

| what you found | what to do |
| --- | --- |
| nothing | move the card to `ready-to-test` |
| **execution** defect (wrong marker, test missing a case) | **fix it yourself**, in the same PR, and return the card to `ready-to-review` releasing ownership |
| **understanding** defect (the spec was misread, the approach doesn't work) | **send it back**: card to `to-do`, ownership to whoever implemented it |

The difference between the last two isn't size, it's nature. Fixing an
understanding defect would hide that the author misread it, and the same mistake
returns on their next card.

> **Whoever fixes doesn't approve their own fix.** By fixing, you became the
> author of that passage.
>
> **Team of one:** the rule presupposes a second pair of eyes, and sometimes
> there isn't one. When the card returns to the SAME agent, the claim **warns and
> hands it over** — blocking wouldn't produce the missing reviewer. The warning
> stays on the card, so whoever reads it later knows the second review wasn't
> independent.

## 7. Merge

The card closes via `Closes #N` in the body, and the pipeline moves it to
`ready-to-test` — the end of Anchors' remit.

From there, delivery is your project's process.
