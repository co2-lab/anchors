<!-- @anchors
  code: INCHN
  updated_at: 2026-10-03
  layer: gate
-->
# InternalChecks — the registry that routes a declared check name to a function

> **Code**: `INCHN`

## Overview

A project declares its gates in the Structure, and a gate that the CLI answers itself
says only a NAME — the check it wants. Something has to turn that name into a function,
and that something is this unit: the registry of internal checkers, plus the small
checkers whose whole answer is reading text.

**Why it is a registry and not a switch.** The checkers do not all have the same
signature, and the difference is not style. One reads only the content of the target.
One also needs the project ROOT, because it has to invoke the version control system to
answer. One needs the GRAPH, because the question crosses the unit — a feature against
the test linked to it. One needs the declaring GATE itself, because the question is
parameterised: a generic gate only knows what to look for after reading its own
configuration, and a project declares several instances of it. Four registries, and the
routing tries them in that order.

**The measured defect this shape exists to prevent**, and it is the one that costs most:
a name that does not resolve must answer PENDING, never Pass. A checker that is declared
and does not exist has not measured anything, and "I did not measure" is neither "it is
clean" nor "it is dirty". Approving here would stamp green over a verification that never
ran, and the project would read coverage where there is none. The same reasoning drives
the aggregate path: a batch or project scope checker that does not resolve answers
Pending too, naming the check that failed to route.

**Two routing paths, and the difference is not a detail.** The per-node path READS THE
TARGET FILE and fails when the read fails — a checker of content with no content has
nothing to answer. The aggregate path does NOT: its scope is the SET, so there is no one
file to read. It hands the checker an empty node and lets it orient itself by root and
configuration. Trying to read a file there would return a read error, and the gate would
fail over a file that never existed.

**What separates this unit from its neighbours.** The engine decides WHICH gates apply to
which node and what the whole run concludes. The rule unit decides how a verification is
named and whether it was waived. This one decides only WHICH FUNCTION answers, and holds
the small checkers whose entire ruler is the text in front of them: the file is not
empty, the file carries an identity code, the file carries a conformant header, the guide
distils its rules into verifiable points.

**The grammar of codes belongs to the project, not to the engine.** The letters that name
a rule type are declared in the Structure, and every pattern in this package that depends
on them has to be reconfigured together. A pattern added without registering it there
stays frozen on the canonical letters — and a scenario written with a letter the project
declared becomes invisible to that gate, which then reports green over what it never
looked at.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the declared check name | any name, resolved or not | — (an unresolved name is a case, not an error) | this unit: what does not resolve is answered undetermined, never approval |
| the target's content | the bytes of the file, text or binary | — (a binary is a case: the text checkers step aside) | this unit: it reads the file and decides, per checker, what it can charge |
| the target node | any node of the map, or an EMPTY node for the aggregate path | — | the engine, which routes by the declared scope |
| the project root | a path, used by the checkers that invoke tooling | — | the caller, which knows where the project lives |
| the map and the Structure | a built graph and a loaded configuration, or nothing | — (absence is handled by each relational checker) | the engine, which passes what it has |
| the rule-type letters | the vocabulary the project declared, or the canonical default | — | the engine, which reconfigures the grammar before running |

## Effects

| Effect | Description |
| --- | --- |
| `INCHN-B01` | A declared check name routes to the function registered under it. |
| `INCHN-B02` | A name that does NOT resolve answers undetermined, never approval — a checker that never ran has not measured anything. |
| `INCHN-B03` | The routing tries the relational registry first, then the one that needs the root, then the pure content one — so a checker that grows a dependency changes registry without changing name. |
| `INCHN-B04` | The per-node path reads the target file, and a read that fails is a failure: a checker of content with no content has nothing to answer. |
| `INCHN-B05` | The aggregate path does NOT read any file: its scope is the set, so the checker receives an empty node and orients itself by the root and the configuration. |
| `INCHN-B06` | A name that does not resolve in the aggregate path answers undetermined too, and the report names the check that failed to route. |
| `INCHN-B07` | An aggregate checker that needs the DECLARING GATE receives it, because a parameterised gate only knows what to look for after reading its own configuration. |
| `INCHN-B08` | `SetRuleLetters` reconfigures every pattern of the package that depends on the project's rule-type vocabulary, together. |
| `INCHN-B09` | A file that is empty or only whitespace fails the emptiness ruler. |
| `INCHN-B10` | A file carrying a scenario code passes the identity ruler, and one carrying none fails it. |
| `INCHN-B11` | A governed file with no identity block at all fails the header ruler. |
| `INCHN-B12` | A governed file whose header carries ownership OR reference passes; one carrying only a layer does not. |
| `INCHN-B13` | A file of a RECOGNIZED layer passes the header ruler with the layer alone, because it has neither an owning spec nor a sibling to reference. |
| `INCHN-B14` | A BINARY file steps aside from the header ruler: there is no comment syntax in an image, and charging one would bar every visual baseline commit. |
| `INCHN-B33` | A governed file whose `@anchors` block stands below the top without `@fixed-header: <why>` fails the header ruler saying the block is not read as the header, and how to fix it. |
| `INCHN-B34` | `scenario-coverage` counts a rule proven only when every scenario the spec's features declare for it is proven — each variant `#NN` a scenario of its own —, and names each variant left unproven: a skipped `#02` is not proven by the green `#01` beside it, nor by a proof recorded for the rule alone. (`checkScenarioCoverage`, `specScenarios`) |
| `INCHN-B35` | `header-valid` reads the header's layer as the project declares it: a file whose header names a layer the Structure declares `regime: declarativo` has its identity in `layer:` alone, whatever the layer is called and whether or not the node carries a regime of its own; a layer declared with another regime still asks for `code:` or `ref:`. (`checkHeaderConforms`, `isRecognizedLayerCfg`) |
| `INCHN-B36` | `line-coverage` and `coverage-delta` skip a file a coverage report listed with no instrumentable line at its revision — there is nothing to cover —, and answer a divergence for one a whole run of its suite left out of the report — a tool omits a file with no instrumentable line, and also one outside what it collects —; a file never listed nor omitted is pending, never measured. (`coverageAbsence`) |
| `INCHN-B37` | `mutation-score` skips a file the mutation tool listed at its current revision with no mutant — an alias, a re-export: it was measured and has nothing to mutate —; a file never listed, or listed at another revision, stays pending. |
| `INCHN-B38` | `line-coverage` holds a code file to the floor of the first glob of the gate's `coverage_floors` that matches it, in name order, and otherwise to its `min_coverage`, 70% when undeclared; a file below a glob's floor fails naming the glob and its reason. (`CoverageFloorFor`) |
| `INCHN-B39` | A header line is read in every comment dialect the map reads — `//`, `#`, `--`, `<!--` and a block comment's ` * ` (`config.HeaderLinePrefix`): a `-- ref: CODE` header has its identity. |
| `INCHN-B40` | `header-valid` reads the header as the map does — the `@anchors` block at the top (`scan.AnchorsHeader`) —, so an `@anchors` further down, in a string or an example, is neither the header nor its identity; and a guide, a document or a test support file, which belong to no unit, have their identity in `layer:` alone. |
| `INCHN-B15` | An executable test script steps aside too, by a different path: its format belongs to the runner, and its identity is in the file name. |
| `INCHN-B16` | A guide with no compliance-points section, or with the section and no item in it, fails — the AI judgment gate would otherwise fall back on vague heuristics. |
| `INCHN-B17` | When the project says how its tests are written, `scenario-coverage` counts a scenario as written in a file the source lists tests in only when a test TITLE cites its code; without a source, or in a file the source lists no test in, a code anywhere in the file outside comments counts. |
| `INCHN-B18` | A support file is not judged by `tests-pass` (Skip, saying why), and it does not count as a test that names a scenario for `scenario-coverage`. |
| `INCHN-B19` | `non-empty` passes a feature only when it declares a scenario, and a scenario opens with any keyword of the official Gherkin table, in any language and synonyms included; the examples table of an outline is not a scenario. |
| `INCHN-B20` | `scenario-coverage` charges only the requirements the spec DEFINES: a code the spec merely cites in its prose is never charged, and a defined requirement with no proven scenario still fails, named. |
| `INCHN-B21` | `scenario-coverage` tells a requirement no test names apart from one a test names but no ingested execution proved, even when no execution was ingested at all, and the verdict says which of the two each one is; a requirement an ingested execution proved is not charged. |
| `INCHN-B22` | A spec whose layer dispenses `tested-by` is skipped by `scenario-coverage`, saying so, as `unit-complete` does; a layer without that opt-out is still charged. |
| `INCHN-B23` | `mutation-score` passes a file whose score reaches the acceptable threshold — the threshold itself included — and fails one below it, naming how many mutants survived and the threshold. A file where no mutant ran does not apply: all ignored says so with the count, and none covered says the line coverage owns it. |
| `INCHN-B24` | With no mutation signal ingested, or with one measured at another revision of the file below the floor or under load, `mutation-score` is pending, saying what to ingest or that the signal is stale; a stale score that met the floor, measured without load, does not block — mutation is not remeasured on every change —, and says how to measure it again. The revision is the mutation's own. |
| `INCHN-B25` | A score between the acceptable and the desirable threshold is a divergence the project's threshold accepts — it informs and never bars; to make it bar, the project raises the acceptable threshold —, naming both ranges and how far the desirable one is; at or above the desirable it passes clean; with no desirable threshold, or one not above the acceptable, the acceptable threshold alone decides. |
| `INCHN-B26` | The mutation thresholds are the ones the ingested report declares; the engine's default of 70% applies only when the report declares none. |
| `INCHN-B27` | With the mutation score measured per scope, the verdict follows the ISOLATED score, never the full one, and the report names both and their delta — a large delta read as coupling to dependents, a small one as a missing assertion, never both at once; with no scopes the total score decides. |
| `INCHN-B28` | A mutation scope measured at an older revision than the file's decides nothing and is left out of the report; a scope carrying no revision stamp still counts. |
| `INCHN-B29` | Outside a git repository `updated-at-atual` skips, naming the missing repository instead of blaming the file as uncommitted; inside one, a new file dated today passes and a wrong date still fails. |
| `INCHN-B30` | `spec-sections` fails a section title written in another language of the catalogue than the project's, naming the expected title, unless the gate declares `enforce_section_language: false`; a title outside the catalogue, or a run with no configuration, is not charged for its language. |
| `INCHN-B31` | The skeletons `anchors new` emits are born conforming: the spec passes the header and spec-sections rulers, and the feature and the test pass the header ruler. |
| `INCHN-B32` | When the share of a file's mutants killed by the time limit is above the gate's `timeout_ceiling` (0.2 by default), its mutation score is Pending as measured under load, whatever the score, and the verdict says to measure again first with fewer workers and with the test cache off — the usual causes are the tool's own parallelism and a time limit taken from a cached run —, and then with no time limit to learn how long a mutant takes; below the ceiling the score decides, and a failure says how many were killed by the time limit. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `INCHN-I01` | A check name that does not resolve NEVER approves. Approving would stamp green over a verification that never ran, and the project would read coverage where there is none. | routes an unknown name through both paths and verifies neither approves |
| `INCHN-I02` | Every registered name is reachable through exactly one of the routing paths. A name registered in no reachable registry is a gate that is accepted in silence and measures nothing. | walks every registered name and verifies the routing resolves it |
| `INCHN-I03` | The compliance ruler is recognised in EVERY language of the catalogue, not only the one the engine was written in. A project seeded in another language would otherwise be born failing a guide its own tooling had just written. | writes the section title in a non-default language and verifies it is recognised |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `INCHN-X01` | Does not decide WHICH checks a project runs. | That is declared in the Structure. A registry that ran what nobody asked for would charge a project for a ruler it never adopted. |
| `INCHN-X02` | Does not invoke external tooling. | These checkers answer by reading TEXT. The ones that shell out belong to the external path of the engine, and mixing them would make a registry lookup depend on what is installed on the machine. |
| `INCHN-X03` | Does not judge whether the text it reads is GOOD. | The rulers here are presence and shape — the file is not empty, it carries an identity, it carries a header. Whether the content is right is judgment, and a deterministic checker that attempted it would fail by a criterion it cannot measure. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `INCHN-E01` | A test the map lists is no longer on disk when `scenario-coverage` looks for the tests that name each scenario code. | That test names nothing; the tests still on disk are read, so a code one of them names is still counted as written. | The map can be older than the tree (a test deleted since the last build): a missing file names no code, and one stale node must not make the codes the other tests name look untested. <!-- @resilient: a stale map node is expected between builds, and the next map build removes it --> |
| `INCHN-E02` | The project's tests source fails, or answers outside its contract, when `scenario-coverage` looks for the tests that name each code | Fail, naming the source's error | Which scenarios have a test cannot be told; calling them untested would send the reader to write tests that may exist |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `CodeLengthPattern` | core — the shape of an identity code is declared by the project, not frozen by the engine |
| DEP2 | `internal/mapx/model.go` | `Node` | core — the target the checker confronts, and its kind |
| DEP3 | `internal/i18n/i18n.go` | `AllTranslations` | core — the compliance ruler is recognised in every language of the catalogue |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
