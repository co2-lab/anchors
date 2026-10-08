<!-- @anchors
  code: DPDCD
  layer: desenho
-->

# Duplicate declarations — each gate confronts the repeats of what it controls

> APPROVED on 2026-10-07 (on by default). One mechanism in the engine; each gate that controls a kind of
> declaration says what its occurrences are; a key declared twice in the same file is a
> finding of that gate. On by default, switched off per gate in `anchors.yaml`.

## The problem

The screen preset wrote `PRBOE-B01` twice — once as a rule heading, once as the first cell of
the loading table — and no gate saw it (fixed in v0.1.297, `NWARN-B21`). The same happens by
hand: a block copied and its number left as it was. The code then names two rules with
different texts; the scenario and the test that cite it prove one of them, without saying
which; a revision answers for a rule that is two.

It is not only rule codes. Most readers of declarations keep the first occurrence and drop the
rest without a word (a `seen` map), or overwrite the first with the last (a `map` by key):
two rows for the same environment variable, the same HTTP status declared twice, two
`@used-by:` flags on one symbol, two flag scenarios with the same code. Whatever the gate
confronts, it confronts one of the two, and the other drifts unseen.

A new gate for duplicates would be one gate knowing every kind of declaration. The gate that
already reads a declaration is the one that knows what makes two of them the same.

## The design

### The mechanism, once

- A gate may register an **occurrence reader**: `func(content, node, root, graph, cfg)
  []Occurrence`, an `Occurrence` being a key, the line it is declared on and, when the
  declaration lives in another file than the node judged, that file.
- After the gate runs on a node — and only when it did not skip it — the engine counts the
  keys. A key declared more than once makes the gate fail on that node, naming the key and
  every line; the gate's own finding, if any, is kept beside it.
- The validation is on by default. `duplicates: false` on the gate in `anchors.yaml` switches
  it off; the gate itself is untouched.
- The verdict is the gate's: a duplicate fails a blocking gate as any finding does, and only
  informs on an informative one. A gate that already treats repeats as a divergence
  (`scenario-identity`, for its migration) keeps that verdict.
- The site page of each gate with a reader says what it counts as a duplicate, and how to
  switch it off.

### Only declarations, each with one owner

A citation repeats by nature — a rule cited by ten scenarios, a message cited by five
errors, a testID queried by every flow. Only a declaration counts, and each kind of
declaration is counted by one gate, so a repeat is reported once.

| Gate | What counts as declared twice | Key | Reader to adapt |
| --- | --- | --- | --- |
| `rule-types` | a rule code defined twice in one file — heading, first cell of a rule row, bold bullet —; every letter (B, V, P, E, M, S…), specs and doctrines alike | `CODE-Xnn` | `definedRequirements` / `defineRuleCaptureRE`, keeping every occurrence; outside the Rule uses section, the open decisions, alias and retired lines; a heading and its own table row in the same section are one definition |
| `spec-sections` | a section title repeated under the same parent — every reader of sections sees only the first | catalog section key + parent | `headingTitles` |
| `scenario-identity` | a scenario code tag on two scenarios of one feature (already measured, now with lines; still a divergence) | `CODE-Xnn#NN` | `parseFeatureScenarios` |
| `phase-ordered` | a phase code twice in a plan (already fails, now with lines) | `CODE-Wnn` | `PlanPhases` |
| `revision-orphans` | a revision code defined twice in a spec or plan | `CODE-Rnnnn` | `RevisionsOf` |
| `flag-scenario-grammar` | a flag scenario code twice in a flag file | `CODE-Gnn` | `flagx.ParseContent` (has lines) |
| `open-questions-resolved` | an open question code twice, resolved rows included | `Qnn` | `openItems` |
| `env-declared` | an environment variable in two rows | variable name | `sectionRows` of Environment Variables |
| `contract-status-declared` | the same HTTP status in two rows of the output contract (concrete statuses only) | status | the section's status rows |
| `dependency-honored` | the same `DEPn` in two rows of the Dependencies table | `DEPn` | the spec's Dependencies table |
| `domain-declared` | the same entry in two rows of the Domain | entry name | `domainLines` |
| `used-by-declared` | two `@used-by:` flags on one symbol | symbol | `scan.UsedByIn` (has lines) |
| `testid-consistent` | a testID in two rows of the spec's inventory (exposing one id in two render branches of the code is not counted) | testID | `declaredTestIDs` |
| `examples-match` | an identical row in one outline's Examples | scenario + row values | `exampleTables` (has lines) |

Left out on purpose: citations of any kind; the code's own repeats (an import twice, a
handle in two branches, overloads); measured data (coverage, mutation, test runs); and
duplicates across files (the same code owned by two specs is `identity`'s question, not
this one's).

## The phases

| Phase | What | Proof |
| --- | --- | --- |
| `DPDCD-W01` | **The mechanism.** `Occurrence`, the reader registry, the engine's count after a gate runs, `duplicates:` on the gate, the message, the site pages. One reader to prove it: `rule-types`, rule codes defined twice. | The screen spec born with `-B01` twice (the bug of v0.1.297) fails `rule-types` naming both lines; `duplicates: false` silences it. Measured on clones of jokenpo and MIF: how many specs already have a repeat, read one by one. |
| `DPDCD-W02` | **The spec catalogue.** `spec-sections`, `revision-orphans`, `open-questions-resolved`, `domain-declared`, `dependency-honored`, `env-declared`, `contract-status-declared`. | Each with its test; measured on the clones. |
| `DPDCD-W03` | **Features, plans, flags, code.** `scenario-identity` and `phase-ordered` moved onto the mechanism; `flag-scenario-grammar`, `used-by-declared`, `testid-consistent`, `examples-match`. | Each with its test; measured on the clones; the peers told what each repeat found. |

## What NOT to do

- Count a citation. A key cited many times is the norm, and a reader that cannot tell a
  definition from a citation stays out until it can.
- Report one repeat from two gates. Each kind of declaration has one owner in the table.
- Silence a repeat by keeping the first. A reader that drops repeats is the hole this closes.
- Make the switch a waiver per line. A repeat is fixed in the file; `duplicates: false` is
  for a project that decides the gate's notion of "the same" does not hold for it.

## Decisions

| Code | Question | Decided |
| --- | --- | --- |
| `DPDCD-D01` | A new gate for duplicates, or a validation in each gate? | In each gate that controls a declaration, switchable per gate (the user, 2026-10-07). |
| `DPDCD-D02` | On by default, or opt-in? | On by default; `duplicates: false` switches a gate's validation off. Measured on the peers before each release, and they are told before a version that would flood them. |

## Open Decisions

none
