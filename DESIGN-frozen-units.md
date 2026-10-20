<!-- @anchors
  code: FRUNF
  layer: desenho
-->

# Frozen units — a part of the project kept as it is, by decision, until someone thaws it

> IMPLEMENTED on 2026-10-11 (units `FRSCF`, `FRGTF`; `FRZEX-B14`..`B19`; `RPSCR-B53`..`B55`).
> Measured on a clone of MIF: `anchors freeze --reason … --unit ILPNI --rule EXFME-X01 --in
> product/exclusao-de-conta.doctrine.md` recorded `EXCAE-F01` and marked both places, and
> `anchors freeze EXCAE-F01 --rule RSPNP-X01` added a third; after
> `map build`, `anchors coverage` lost the two frozen rules from the unproven ones and the
> unit's code from the mutation list (1934 → 1929 survivors, its 5); `check --changed` over the
> unit's four files answered `frozen` for them, counted in a column of their own, with no
> failure; a live agent session editing the frozen code received the freeze's code and reason
> before the edit, the edit went through, and the agent asked before keeping it; `thaw
> EXCAE-F01` gave every file back byte for byte. The agent hook answers in 0.6 s on MIF's map.
>
> APPROVED on 2026-10-11. While building, a team sometimes keeps something out on purpose: a
> feature that will not ship now, a flag left off until later, a section of the doctrine whose
> time has not come. Today it either fails gates, stales proofs and fills the mutation and
> coverage reports, or someone deletes it to make the noise stop; and an agent that meets it
> takes it for unfinished work and evolves it without an order. A frozen part keeps its
> place, keeps quiet, and says to whoever reaches it that it is frozen, and why.

## What exists, and why it is not this

- `anchors freeze` / `thaw` stop the **whole project**: `enabled: false` in anchors.yaml, a
  ruleset blocking push and merge on the remote, an issue with the reason. It is the brake.
- `@dormant: <reason>` after a navigation flag keeps a navigation call out of the map's edges
  (NCGNV-B08). It speaks of one call.

The frozen part is the same idea as the brake, at the scale of a unit: the same verbs, the
same mandatory reason, a smaller scope.

## The design

- **A freeze is a thing with a code, recorded in the doctrine** (the user, 2026-10-11: "assim
  o congelamento teria um código e poderia ser rastreado por toda a aplicação"; "o registro
  dos congelamentos fica nos arquivos de doutrina"). A freeze is an entry of a doctrine
  document, coded with the freeze letter `F` — `QLDTQ-F01` —, carrying its reason: the
  doctrine is where the project says what holds and what does not, and a freeze is a
  decision of that kind. Every mark carries the code, not the reason: one freeze can cover a
  feature that crosses several units, files and doctrine sections; `anchors freeze show
  QLDTQ-F01` lists every place it covers, and `anchors thaw QLDTQ-F01` removes every mark and
  the entry at once.
- **Without a doctrine, a freeze is not traceable.** A project with no doctrine document
  marks with the reason itself — `frozen: <reason>` —: the part is quiet and warns all the
  same, and nothing ties one mark to another.
- **What can be frozen, and how it is marked** — always with a reason in the record, never
  with a date (the user, 2026-10-11):

  | Scope | Where | How |
  | --- | --- | --- |
  | a whole unit | its spec's header | `frozen: QLDTQ-F01` |
  | one rule of a unit | the rule's line | `@frozen: QLDTQ-F01`, as `@realizes` is written |
  | one file — code, test | its header | `frozen: QLDTQ-F01` |
  | a whole doctrine document | its header | `frozen: QLDTQ-F01` |
  | a section of a doctrine document | its heading's line, or a doctrine rule's line | `@frozen: QLDTQ-F01` |

  With no doctrine, the mark carries the reason in place of the code.

- **What a frozen part covers.** A frozen unit covers its spec and every file derived from it
  — code, feature, tests, captures. A frozen rule covers its scenarios. A frozen doctrine
  document or section covers its rules, and every spec rule that realizes one of them
  (`@realizes`): the doctrine is the unit's reference, and a frozen reference freezes what
  answers it.
- **The commands.** `anchors freeze --reason "…" --unit <CODE> --rule <CODE-X01> --file
  <path>` (any of them, repeatable) creates a freeze: its entry in a doctrine document (`--in
  <document>` when the project has more than one; the entries live under a section of their
  own), and its marks; `anchors freeze QLDTQ-F01 --unit …` adds a scope to an existing one;
  `anchors freeze show [QLDTQ-F01]` lists the freezes, or every place one covers; `anchors thaw
  QLDTQ-F01` removes its marks and its entry. With no doctrine, `freeze` writes the reason in
  the marks, and `thaw --unit|--rule|--file` removes them. `anchors freeze` with no scope and no
  code stays the project's brake, as it is.
- **The traces hold.** A mark naming a freeze code no doctrine declares, and a freeze entry no
  mark carries, are findings of a gate (`freeze-consistent`): a freeze is either somewhere,
  with its reason, or nowhere.
- **A new rule letter, `F`.** It goes into the three places a letter lives — the canonical
  letters, the section catalog and the test-signal reader's copy —, and a project that
  declares its own `rule_types` adds it to use coded freezes.
- **What a frozen part stops producing.**
  - Every gate's verdict on it is *frozen (reason)* — neither failure nor pending —, counted
    in a column of its own in `anchors check`.
  - It leaves the stale and expired-evidence reports, the mutation and coverage reports, and
    the selection of `anchors test` and `anchors mutation`.
  - `anchors status` lists what is frozen, with each reason, so nothing frozen is forgotten
    for lack of a date.
- **How it speaks to whoever reaches it** — it warns, it does not block (the user, 2026-10-11):
  1. `anchors work`, `next` and the guides, reaching a frozen part, say *FROZEN: <reason> —
     do not change it without an order*.
  2. A commit that changes a frozen part is told so by a gate, with the reason; the gate
     informs and does not bar.
  3. The agent hook (`install-hooks --agent`) gains a check before the agent edits a file: a
     frozen file, or a file of a frozen unit, gets the same warning in the agent's context,
     and the edit goes on.
- **Gates stay generic.** The marks are flags and header fields; the scope is read from the
  map (a unit's derived files, a doctrine rule's realizers), never from a list of paths.

## The phases

| Phase | What | Proof |
| --- | --- | --- |
| `FRUNF-W01` | **The freeze and its scope.** The freeze letter `F` and the doctrine entries; the scan reads `frozen:` in headers and `@frozen:` on rule and heading lines, by code or by reason; the map carries which freeze covers what, down a unit's derived files and from a doctrine section to the spec rules that realize it; `freeze-consistent` confronts marks and records. | One freeze over a unit, a rule of another unit, a file, a doctrine document and a doctrine section: the map names every part it covers, and nothing else; a mark with an unknown code and a record with no mark are each a finding. |
| `FRUNF-W02` | **Quiet.** Gates give *frozen* on a frozen part, counted apart; the stale, evidence, mutation and coverage reports and the run selections leave it out; `status` lists it. | On a clone of MIF, one feature frozen: its failures, stale proofs and mutation entries leave the reports, and `check` counts them as frozen. |
| `FRUNF-W03` | **Speaking.** `freeze` with scopes, `freeze show`, `thaw` by code; `work`, `next` and the guides warn; a commit touching a frozen part is told; the agent hook warns before an edit. | An agent session on a clone: editing a file of a frozen unit shows the freeze's code and reason, and the edit is not blocked; `thaw` removes every mark of the code. |
| `FRUNF-W04` | **Folders and work.** `frozen_paths` in anchors.yaml, applied by the map and confronted by `freeze-consistent`; `freeze --path`, `thaw`, `show`; `next` and the cold start leave frozen work out; `impact` lists frozen parts apart; the progress leaves a frozen plan out; the freezes section in the project's language. | On a clone of MIF, the income-tax folders frozen by path: their gates answer frozen, `next` hands out none of their tasks, `impact` of a shared file lists them apart. |

## Freezing a folder, and leaving the work of a frozen part out (2026-10-20)

> IMPLEMENTED on 2026-10-20 (`FRSCF-B05`, `FRGTF-B06`, `FRZEX-B20`..`B22`, `WRQUW-B22`,
> `TSQUT-B16`, `MPCTI-B06`, `DMLDD-B14`, `LDCML-B09`, `CNFGO-B76`). Measured on a clone of MIF:
> `freeze --in product/planos.doctrine.md --path` over the income-tax features folder and its
> flows recorded `PLNAP-F01` under `## Congelamentos` and wrote the block; after `map build`,
> 141 tests and 29 specs answered frozen in evidence-fresh, tests-pass and scenario-coverage;
> `impact` of a shared file listed 57 frozen parts apart; `thaw PLNAP-F01` gave both files
> back.

Asked from MIF: the income-tax module is off for the year, out of the plan's scope, spread
over 29 units, a features folder and a folder of flows. Frozen unit by unit, its gates went
quiet, but `anchors next` still handed out its cards, the progress counted it, and `impact`
listed its files as work to redo.

- **A folder is frozen in anchors.yaml**, as a scope of a freeze:

  ```yaml
  frozen_paths:
    - freeze: PLNAP-F01          # the freeze's code, or the reason itself with no doctrine
      paths:
        - "apps/mobile/src/features/ir/**"
        - "apps/mobile/.maestro/screens/ir/**"
  ```

  The record stays the doctrine's entry (D05, D06); the yaml adds where it applies. The map
  marks every file a glob matches as frozen, as its own header's mark would — a unit's owner
  matched by a glob freezes the unit, wherever its files are —; a file's own mark wins. The
  gates still read the map, never the globs. A glob that matches no file is a finding of
  `freeze-consistent`, as a mark with an unknown code is. `anchors freeze [CODE] --path
  <glob>` writes the scope, `anchors thaw CODE` removes it with the marks, `freeze show` lists
  it.
- **The work of a frozen part is not handed out.** `anchors next` leaves in the queue, and does
  not claim, a task whose file is frozen: it waits there for the thaw, quiet, not gone. The
  cold start does not count a spec a plan seeds inside a frozen folder as missing, and does
  not seed a frozen plan. `anchors work`, asked by a person for a frozen part, still warns
  and goes on (D02). In github mode the claim runs in the pipeline, which reads no map: the
  card of a frozen unit is still handed out there, with the warning.
- **`anchors impact` lists frozen parts apart**: what a change reaches that is frozen is named
  under its own heading, not among the parts to redo.
- **The progress leaves a frozen plan out.** A demand whose plan is frozen — its header, or a
  frozen folder — is marked frozen and leaves the totals. A task names no unit, so a task of a
  live plan is cancelled with its reason (`progress cancel`), as before.
- **The freezes section follows the project's language** — `## Freezes`, `## Congelamentos`
  —, and a document already holding it under another language's title keeps it.

## What NOT to do

- Freeze without a reason. A frozen part with no reason is indistinguishable from a forgotten
  one; the reason lives once, in the record, and the marks carry its code.
- Block. The user decided the freeze warns: the order to change it comes from a person, and
  the tools say so instead of standing in the way.
- Let a frozen part disappear. It stays in the map, listed in `status`; it is quiet, not
  gone.

## Decisions

| Code | Question | Decided |
| --- | --- | --- |
| `FRUNF-D01` | What is it called? | `freeze`/`thaw` with a scope — the project's brake keeps the verbs with no scope (the user, 2026-10-11). |
| `FRUNF-D02` | Does it block an edit? | No — the agent hook and the commit warn, with the reason (the user, 2026-10-11). |
| `FRUNF-D03` | Does a freeze expire? | No date (the user, 2026-10-11); `status` lists every frozen part so none is forgotten. |
| `FRUNF-D04` | Can the doctrine be frozen? | Yes — a whole document or a section, and what realizes it follows (the user, 2026-10-11). |
| `FRUNF-D05` | Is a freeze an entity with a code? | Yes — recorded once with its reason and traceable everywhere it applies (the user, 2026-10-11). |
| `FRUNF-D06` | Where does a freeze's record live? | In the doctrine: an entry of a doctrine document, coded with the letter `F`; with no doctrine, the marks carry the reason and are not traceable (the user, 2026-10-11). |
| `FRUNF-D07` | Can a folder be frozen? | Yes — by globs in anchors.yaml under a freeze's code; the map applies them (the user, 2026-10-20). |
| `FRUNF-D08` | Does `next` hand out the work of a frozen part? | No — the task waits in the queue for the thaw (the user, 2026-10-20). |
| `FRUNF-D09` | Does `impact` list frozen parts? | Apart, under their own heading (the user, 2026-10-20). |
| `FRUNF-D10` | Does the progress count a frozen plan? | No — marked frozen, out of the totals (the user, 2026-10-20). |

## Open Decisions

None.
