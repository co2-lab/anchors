<!-- @anchors
  code: CNFGO
  updated_at: 2026-09-26
  layer: config
-->
# Config — loads the project's anchors.yaml, refuses what it cannot honour, and answers every setting with its default

> **Code**: `CNFGO`

## Overview

Every command of Anchors starts by reading the project's `anchors.yaml`. This unit is that reading: it turns
the file into the project's configuration, and it is the one place that sees the whole file at once, so it is
where a wrong declaration has to be caught.

The governing decision is that a declaration the tool cannot honour is a load error, never silence. An
unknown key, a gate enum with a typo, a workflow mode nobody knows, a pattern that does not compile or a
language that is not translated would each load as "nothing declared" and turn a protection off without a
word. Refusing at load names the field and the line before any gate answers for it. When the refusal is an
unknown key, the message also names the hypotheses the reader cannot see alone: the key may be misspelled,
the binary may be older than the file, or the file may be in an older format whose key was renamed, in which
case the fix is to migrate, not to update the binary.

The second decision is that silence keeps the behaviour a project already had. A gate the framework ships
(a canonical gate) may be declared by name only, and its omitted fields are filled from the framework's
declaration while every field the project wrote wins. An undeclared severity never blocks, an undeclared
phase list runs everywhere, an absent `enabled` never freezes, and the workflow, branch and approval
settings each have one default kept in one place.

The unit also answers the project's own vocabulary (section titles, placeholder words, rule letters, scenario
tags, identity code lengths), selects the test suites a command asked for, and resolves the derived-file
patterns a spec governs.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the configuration file | a YAML file whose every key is a field the tool knows | unknown keys, invalid YAML, a file that cannot be read | this unit: the load fails naming the cause |
| gate declarations | gates with distinct IDs (the name stands in when there is no ID) and scope, cost and phase values from their closed lists | a repeated ID, an enum value outside its list | this unit: the load fails naming the gate and the value |
| the workflow block | an absent block, `local`, `manual`, or `github` with an `owner/name` repository and at least one label | GitHub fields in local or manual mode, an unknown mode | this unit: the load fails naming what is missing or extra |
| declared patterns | dialect and derived regular expressions that compile | a pattern that does not compile | this unit: the load fails naming the field |
| the language | a supported language, or none | a language the catalog does not translate | this unit: the load fails listing the supported ones |
| identity code lengths | lengths from 2 to 8 | lengths outside that range | this unit: the load fails naming the length |
| suite filters | any names, in any case and spacing | nothing is outside: an unknown name is reported, not refused | this unit: it labels each missing name by its axis |

## Effects

### Loading the file

| Effect | Description |
| --- | --- |
| `CNFGO-B01` | A key the configuration does not know, at the top level or inside a list item, is a load error that names the key. |
| `CNFGO-B02` | The error for an unknown key adds two hypotheses: the key is misspelled, or the binary is older than the file, with the command that updates it. |
| `CNFGO-B03` | When the file's format is older than the current one (a file with no `version:` is format 1) and the migration registry says the key was renamed, the error advises running the migration and does not advise updating the binary. |
| `CNFGO-B04` | The migration advice needs both conditions: a file already in the current format, a key the registry does not know as renamed, or no registry at all gets the general hint. |
| `CNFGO-B05` | An error that is not about an unknown key (a YAML syntax error) carries no version hint. |
| `CNFGO-B06` | Every gate has a distinct ID, the name standing in when no ID is declared; two gates with the same ID fail the load naming both. |
| `CNFGO-B07` | A gate's `scope`, `cost` and each `when` phase must come from their closed lists; a value outside fails the load naming the value and the valid ones. |
| `CNFGO-B08` | With no workflow block, or mode `local` or `manual`, the workflow is local; a `repo` or `labels` declared in those modes fails the load (`ManualMode`). |
| `CNFGO-B09` | Mode `github` fails the load without a `repo`, with a `repo` not in `owner/name` form, or without at least one label (`GitHubMode`). |
| `CNFGO-B10` | An unknown workflow mode fails the load, saying there is no fallback between modes. |
| `CNFGO-B11` | A declared dialect or derived pattern that does not compile fails the load naming the field (list patterns name their index), and a valid pattern loads. |
| `CNFGO-B12` | The declared language becomes the language of every message from the load on; an unsupported language fails the load listing the supported ones. |
| `CNFGO-B13` | A code length outside 2 to 8 fails the load; valid lengths become the lengths the engine recognizes and are handed to the code generator's hook when one is registered (`SetCodeLengths`, `SetSlotsHook`). |

### Canonical gate declarations

| Effect | Description |
| --- | --- |
| `CNFGO-B14` | A gate the framework's catalog knows inherits, from the catalog's declaration, every field the project omitted, severity included (`SetCanonicalGateResolver`). |
| `CNFGO-B15` | A field the project declared wins over the catalog's, including an explicit `blocking: false` over a blocking canonical gate. |
| `CNFGO-B16` | `run` and `check` are one pair: when the project declares either, neither is inherited. |
| `CNFGO-B17` | A gate the catalog does not know loads exactly as declared. |

### Gate fields

| Effect | Description |
| --- | --- |
| `CNFGO-B18` | A gate whose severity nobody declared does not block (`IsBlocking`). |
| `CNFGO-B19` | The scope is one run per target unless `batch` or `project` is declared; in a full scan the gate uses `scope_full` only when it is `batch` or `project`, and its ordinary scope otherwise (`EffectiveScope`, `ScopeForScan`). |
| `CNFGO-B20` | A gate with no phases runs in every phase, and an unnamed phase admits every gate; a gate with phases runs only in those it lists (`RunsIn`). |
| `CNFGO-B21` | A gate takes part in both perspectives (change and all) unless `skip_on` lists the perspective (`SkipsOn`). |
| `CNFGO-B22` | The mutation report format is read from the `mutation-score` gate only, trimmed and lower-cased, and is the Mutation Testing Elements format when none is declared (`MutationFormat`). |
| `CNFGO-B23` | A gate checks the language of section titles unless it declares that check off (`ChecksSectionLanguage`). |

### Workflow defaults

| Effect | Description |
| --- | --- |
| `CNFGO-B24` | The integration branch is `main` when undeclared; the protected branches are the declared list, else the integration branch and `main` (only `main` when they are the same) (`IntegrationBranchOrDefault`, `ProtectedBranchesOrDefault`). |
| `CNFGO-B25` | A pull request needs one approval when the number is undeclared; a declared number, zero included, is kept (`RequiredApprovalsOrDefault`). |
| `CNFGO-B26` | A stale pipeline and an ingest run outside the test command only warn, unless the project declares that they block (`PipelineVelhoBarra`, `IngestManualBarra`). |

### Freeze

| Effect | Description |
| --- | --- |
| `CNFGO-B27` | Only an explicit `enabled: false` freezes the project; an absent value, `true`, or a configuration that did not load does not (`Frozen`). |
| `CNFGO-B28` | The freeze reason is shown trimmed; a missing or blank reason is replaced by a message that asks for `freeze_reason` (`FreezeReasonText`). |

### Project vocabulary

| Effect | Description |
| --- | --- |
| `CNFGO-B29` | A section title comes from the layer's rename, then the project's, then the framework's default; a blank rename does not count (`SectionTitle`). |
| `CNFGO-B30` | The placeholder words are the declared ones, trimmed and without blanks, or the templates' single default marker word when none remain (`Placeholders`). |
| `CNFGO-B31` | The rule letters are the declared rule types' letters, upper-cased, one character each, without repeats, in declaration order; with none valid they are the canonical `SRVAXBNMDEIQFG` (`RuleLetters`). |
| `CNFGO-B32` | A scenario tag maps to every letter whose rule type declares it, ignoring case and surrounding spaces; a tag no rule type declares is reported as unknown (`TagLetters`). |
| `CNFGO-B33` | The code length pattern, placed after one character class, matches exactly the declared length when there is one, and exactly the range when the lengths are contiguous (`CodeLengthPattern`). |
| `CNFGO-B41` | A rule type catalogues a section when the section's title is one it declares as requiring a code, ignoring case and surrounding spaces (`RequiresCodeIn`). |

### Suite selection

| Effect | Description |
| --- | --- |
| `CNFGO-B34` | With no filter, every declared suite is selected; a suite with no workspace or no scope is not excluded by an axis nobody filtered (`SelecionaSuites`). |
| `CNFGO-B35` | The layer, workspace and scope filters each keep the suites that match one of their names, and the axes combine as an intersection. |
| `CNFGO-B36` | The selected suites keep the order of the file, not the order of the request. |
| `CNFGO-B37` | Names match ignoring case and surrounding spaces. |
| `CNFGO-B38` | A requested name that no suite declares comes back labelled with its axis, all of them together; names that exist but never together are an empty selection, not a missing name. |
| `CNFGO-B39` | The declared layers, workspaces and scopes are listed once each, in file order, leaving out the empty ones (`DeclaredLayers`, `DeclaredWorkspaces`, `DeclaredScopes`). |

### Derived file patterns

| Effect | Description |
| --- | --- |
| `CNFGO-B40` | When the derived files declare `patterns`, those patterns become the code patterns, the other derived files are kept, and `patterns` never becomes a derived file of its own; without it the declared files are used as they are (`PadroesDe`). |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CNFGO-I01` | What the configuration writer saves, the loader reads back with the same layers and gates (`Save`). | saves a loaded configuration and loads the written file again |
| `CNFGO-I02` | Refusing unknown keys never refuses a known one: a file that uses only declared keys loads complete. | loads a file with layers, governs, gates and touch settings and counts what came back |
| `CNFGO-I03` | Every reader of the configuration answers its default on a configuration or workflow that is absent, never a failure (`GitHubMode`, `ManualMode`, `RouteRegistry`). | calls the readers on absent values and checks each default |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CNFGO-X01` | Does not infer the GitHub repository from the git remote. | In a fork, inferring would send issue writes to the wrong repository, and a write in the wrong place is not undone by a revert. |
| `CNFGO-X02` | Does not validate suite layer, workspace or scope names against a fixed list. | The test vocabulary belongs to the project; a fixed list would be the framework deciding it. |

## Errors

Most failures of the load are the refusals stated as behaviours; the rows below catalogue them and add the two I/O failures.

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CNFGO-E01` | The configuration file cannot be read. | The load fails with the read error and no configuration. | A command must not run on a configuration it never read. |
| `CNFGO-E02` | The path the configuration is saved to cannot be written. | Saving returns the write error. | A save that failed silently would leave the project believing its file changed. |
| `CNFGO-E03` | REF[CNFGO-B01]: an unknown key is the refusal B01 states: the load fails naming the key | — | — |
| `CNFGO-E04` | REF[CNFGO-B06]: a repeated gate ID is the refusal B06 states: the load fails naming both gates | — | — |
| `CNFGO-E05` | REF[CNFGO-B07]: a gate enum outside its list is the refusal B07 states | — | — |
| `CNFGO-E06` | REF[CNFGO-B09]: an incomplete GitHub workflow is the refusal B09 states | — | — |
| `CNFGO-E07` | REF[CNFGO-B10]: an unknown workflow mode is the refusal B10 states | — | — |
| `CNFGO-E08` | REF[CNFGO-B11]: a pattern that does not compile is the refusal B11 states | — | — |
| `CNFGO-E09` | REF[CNFGO-B12]: an unsupported language is the refusal B12 states | — | — |
| `CNFGO-E10` | REF[CNFGO-B13]: a code length out of range is the refusal B13 states | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/i18n/i18n.go` | `T`, `Set` | core — the language of every message, set at load |
| DEP2 | `internal/config/min_version.go` | `validarMinVersion` | config — the minimum version check the load runs |
| DEP3 | `internal/config/patterns.go` | `Padroes` | config — a pattern declared as one string or a list |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
