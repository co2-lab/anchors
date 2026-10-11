<!-- @anchors
  code: GDLYG
  layer: desenho
-->

# The guidance layer — the first canonical layer: how the project says it is worked on

> IMPLEMENTED on 2026-10-11 (units `GDGTG`; `CNFGO-B68`/`B69`, `EDSTD-B24`, `GVGDG-B28`,
> `CNGDC-B10`, `DFGTD-B04`). Measured: on clones of MIF and jokenpo the layer is EMPTY — neither has
> a README, CONTRIBUTING, AGENTS or CLAUDE at its root (theirs live in subfolders) — and no
> verdict of theirs changes; both declare `regime: declarativo` on dao/infra/presentation/
> domain-types, so removing the name fallback (D04) changes nothing for them either. In this
> repository the layer holds README.md and CONTRIBUTING.md: `guidance-linked` passes both,
> `header-valid` skips them. A live AI judgment (Claude haiku, through `anchors judge`) on a
> project `init` seeded passed its three files; on one whose CONTRIBUTING.md was told to commit
> with `--no-verify` and to use an undeclared `controllers` layer it failed that file naming
> both passages, and passed the other two. The verdicts stayed on the edgeless nodes, held on
> the next check (✓2 ✗1), and a change to anchors.yaml asked the three again.
>
> APPROVED on 2026-10-11 (the user: "criar a primeira camada canônica … contemplando
> contributing.md readme.md claude.md agent.md etc, e pra essa camada um gate de julgamento para
> ver se ela está respeitando as recomendações do anchors"). Decisions D01–D05 below.

## The problem

Whoever arrives at a project — a person or an agent — learns how to work in it from a handful of
files at its root: `README.md`, `CONTRIBUTING.md`, and, for agents, `AGENTS.md` and `CLAUDE.md`.
An agent reads its file before anything else, and does what it says.

Nothing governs those files. `anchors init` seeds `CONTRIBUTING.md`, `AGENTS.md` and `CLAUDE.md`
once, from the configuration of that day; then the layers change, a gate starts to block, the
workflow moves to github mode — and the guide still describes the project as it was. Worse, a
hand edit can tell agents the opposite of what the gates ask ("commit with `--no-verify` when the
hook is slow", "tests can come after the code"), and every agent that reads it follows the file,
not the gates. Today no gate reads a word of it.

## What exists, and why it is not this

- Anchors has **no canonical layer**. Every layer is the project's: it declares them in
  `layers:`, and Anchors only asks that they be kept apart (Anchors proposes no structure).
- The only fixed layer names are a leftover: `header-valid` exempts `dao`, `infra`,
  `presentation` and `domain-types` from `code:`/`ref:` when they declare no `regime`. That is a
  structure guessed from names; this design does not follow it, and its removal is a separate
  decision (GDLYG-D04).
- `anchors guide` holds Anchors' rules for writing each artifact; nothing holds what a project's
  guidance files must say.

## The design

- **One canonical layer, `guidance`.** Anchors knows it without the project declaring it: the
  files at the root that say how the project is worked on — `README.md`, `CONTRIBUTING.md`,
  `AGENTS.md`, `CLAUDE.md` (GDLYG-D02). It is canonical because it is not about how
  the project's CODE is arranged — Anchors still proposes no code structure —, but about the
  working agreement, which every project under Anchors has, wherever its code lives.
- **What makes it canonical.** It enters every project's Structure as a default, with its
  pattern, kind `doc`, regime `declarativo` and tag `guidance`. A project that declares a layer
  named `guidance` replaces it (another pattern, another file list); one that does not want it
  turns it off (GDLYG-D03). Its files carry no `@anchors` header: a header at the top of a README
  is noise for every reader, and the layer is known by its pattern.
- **A guide for it: `anchors guide guidance`.** What a guidance file must say under Anchors, in
  one place: the order spec → feature → test → code; the layers the project declares and that
  they are kept apart; where the work comes from (`anchors next`) and how a task is framed
  (`anchors work`); which gates bar a commit; watching long commands and agents with
  `anchors monitor`; never `--no-verify` — a gate that fails is fixed or waived in the commit
  message with its reason. And what each file is for: `CONTRIBUTING.md` holds the agreement,
  `AGENTS.md`/`CLAUDE.md` send the agent to it, `README.md` sends a newcomer to it.
- **A mechanical gate, `guidance-linked`** (informative): `CONTRIBUTING.md` exists, and each of
  the other guidance files present points to it. What can be measured is measured, before any
  judgment is spent.
- **A judgment gate, `guidance-honored`** (an AI judges, as `rule-fulfilled` does): each
  guidance file is read against `anchors guide guidance` and against the project's own
  `anchors.yaml`. It fails a passage that contradicts them — an instruction to skip a gate, a
  layer named that the project does not declare, a blocking gate left out, an order of work
  other than the anchor's — naming the passage and proposing the text that replaces it. It does
  not judge style or completeness of prose. It is asked when a guidance file changes or
  `anchors.yaml` changes, never on every commit.
- **The seed follows the guide.** `init` keeps seeding `CONTRIBUTING.md`, `AGENTS.md` and
  `CLAUDE.md` from the configuration; the seeded text is what the guide asks, so a fresh project
  passes both gates.

## The phases

| Phase | What | Proof |
| --- | --- | --- |
| `GDLYG-W01` | **The layer.** The canonical `guidance` layer in every Structure, replaceable and switchable; no header asked of its files; `guidance-linked`. | On clones of MIF and jokenpo: the map holds their guidance files in the layer, no header gate asks them anything, and `guidance-linked` names what is missing. |
| `GDLYG-W02` | **The judgment.** `anchors guide guidance`; `guidance-honored`, asked on a change to a guidance file or to `anchors.yaml`. | A seeded project passes; a CONTRIBUTING.md telling to commit with `--no-verify`, and one naming a layer the project does not declare, each fail with the passage named. |

## What NOT to do

- Propose a code structure through this layer. It governs the working agreement, never where
  code lives.
- Judge prose. The judgment asks for contradictions with the guide and the configuration, not
  for style or length.
- Rewrite a project's files. `init` seeds what is missing; the gates say what contradicts; the
  edit is the project's.

## Decisions

| Code | Question | Decided |
| --- | --- | --- |
| `GDLYG-D01` | The layer's name. | `guidance` (the user, 2026-10-11). |
| `GDLYG-D02` | Which files. | `README.md`, `CONTRIBUTING.md`, `AGENTS.md`, `CLAUDE.md` (the user, 2026-10-11). |
| `GDLYG-D03` | On by default, or only when seeded. | In every project; `guidance: false` turns it off (the user, 2026-10-11). |
| `GDLYG-D04` | The fixed names `dao`/`infra`/`presentation`/`domain-types` in `header-valid`. | Removed: `regime: declarativo` alone decides, measured on clones first (the user, 2026-10-11). |
| `GDLYG-D05` | Does `guidance-honored` bar? | It informs (the user, 2026-10-11). |

## Open Decisions

None.
