<!-- @anchors
  code: INWZN
  updated_at: 2026-10-01
  layer: comando
-->
# InitWizard — the interactive init walks a person from an unconfigured directory to a reviewed anchors.yaml, and writes nothing on answers nobody gave

> **Code**: `INWZN`

## Overview

`init` scans the project without AI, proposes a configuration, and confirms it with the person
through questions; the bulk is inferred, and the questions cover only human decisions. It
proposes no project structure: the code layers it offers are the project's own folders that
hold code, and the layer kinds it names (entry points, use cases, domain, repositories,
infrastructure, presentation) are illustration of keeping layers apart, never a layout to move
files into. The language is dialect only — how tests are named, how the test runner emits the
reports ingest reads. Three
files make the wizard, and this spec states how they behave together: the command and its
questions (init), the git step that runs first (init git), and the DISCOVER step that runs
right after the scan (init discover).

The git step comes first because much of Anchors depends on a repository with a commit. A
ready repository passes silently; a machine without git gets a warning and the init goes on;
otherwise the person is offered to initialize the repository or to make the first commit.
Declining is legitimate and names what stays off; accepting initializes only what is
missing, seeds a `.gitignore` only when none exists, and makes the first commit; a failure is
reported with its cause and does not stop the init.

The DISCOVER step exists because the init infers from the disk, and an empty project has
nothing to infer from. When the project has neither code, specs nor a PROJECT.md, an AI
operating the CLI gets a work order (run the project guide, interview the user in the
conversation, write PROJECT.md and INSIGHTS.md, then init again); a person gets the same
instruction, is offered to open the AI tool detected on the machine with the interview prompt,
or gets the prompt to paste.

The rule the three files share: a prompt that cannot run — no terminal, or an input that
ended — is not an answer. Every prompt failure is recorded, and the init then writes nothing
at all (no configuration, no guide, no repository, no commit) and fails with an error that
names the non-interactive mode, the way out for agents, pipes and CI. The one exception is the
overwrite confirmation of an existing configuration: without a yes, the init stops and keeps
the file, which is a successful no.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the operator | a person at a terminal, or an AI detected from its environment | a caller with no terminal and no answers | this unit: fails naming `--non-interactive` |
| the answers | the person's choices at each prompt | a prompt that could not run or whose input ended | this unit: records the failure and writes nothing |
| the git state | not installed, not initialized, initialized without a commit, ready | — | the init package detects it; this unit acts on it |
| the project | an empty directory, or one with code, specs or PROJECT.md | — | the init package infers; this unit decides whether DISCOVER is due |

## Effects

| Effect | Description |
| --- | --- |
| `INWZN-B01` | An existing `anchors.yaml` is overwritten only on an explicit yes; otherwise the init stops successfully and the file stays byte for byte. |
| `INWZN-B02` | In a repository with a commit, the git step prints nothing and asks nothing. |
| `INWZN-B03` | On a machine without git, the git step warns that git is not installed, asks nothing, and the init goes on. |
| `INWZN-B04` | Declining the git offer creates no repository, names what stays off without git, and the init goes on. |
| `INWZN-B05` | Accepting the git offer leaves a repository with HEAD: it initializes the repository only when missing, seeds a `.gitignore`, stages everything and makes the first commit. |
| `INWZN-B06` | An existing `.gitignore` is never overwritten by the git step. |
| `INWZN-B07` | A git initialization that fails is reported with its cause and does not stop the init; a missing author identity is reported with the `git config` commands that fix it. |
| `INWZN-B08` | After the scan, the init prints what it found on disk (specs, features, tests by the pattern they follow, the language family, guides, code directories, co-location), and only what it found. |
| `INWZN-B09` | When the project has code, specs or a PROJECT.md, the DISCOVER step prints nothing and the init goes on. |
| `INWZN-B10` | In a project with nothing to infer from, an AI operator gets the DISCOVER work order and the init goes on. |
| `INWZN-B11` | A person with no known AI tool gets the step-by-step with the whole interview prompt, wrapped at the width without losing a word, and the init goes on. |
| `INWZN-B12` | A person with a detected AI tool is offered to open it: accepting runs the tool in the project root with the interview prompt as one argument, without a shell; declining prints the step-by-step. |
| `INWZN-B14` | `--non-interactive` routes the command to the JSON mode instead of the prompts. |
| `INWZN-B16` | Accepting the header guide seeds `HEADER_GUIDE.md` in the project's guide directory, or in `guides/` when the project has none. |
| `INWZN-B17` | The default gates are offered only when the chosen artifacts yield at least one, and accepting them writes them into `anchors.yaml`. |
| `INWZN-B18` | A project with no code, spec, feature or test is announced as new before the questions; a project with any of them is not. |
| `INWZN-B19` | The code-layer question is asked only when the configuration has code layers; a new project without any is told to declare them once they exist. |
| `INWZN-B20` | Each guide found on disk gets a question asking which tag it governs. |
| `INWZN-B21` | `--preset`, in either mode, is refused before anything runs, with an error saying Anchors proposes no structure and init reads the project's own folders; nothing is written. |
| `INWZN-B22` | Before the code-layer question the init says that each candidate is a folder of the project, that Anchors proposes no structure, and gives layer kinds as illustration. |
| `INWZN-B23` | Accepting the contributing guide seeds `CONTRIBUTING.md` from the written configuration when the project has none; an existing one stays byte for byte, and the init prints the section that would be added. |
| `INWZN-B24` | When the language family has a coverage hint, the init prints how its tests emit the reports `anchors ingest` reads. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `INWZN-I01` | When any prompt of the init cannot run, the init writes no configuration and no guide, creates no repository and no commit, and fails naming `--non-interactive`. | runs the init with no terminal outside git, in a repository with code, in an empty repository, and with an input that ended, checking the tree after each |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `INWZN-X01` | The init never initializes a repository nor commits without an explicit yes. | Git state is the person's; a commit nobody accepted would carry an author and content nobody chose. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `INWZN-E01` | A prompt cannot run because there is no terminal. | The error keeps the caller's context and shows both calls of the non-interactive mode, the one that asks and the one that answers with an example. | An agent or a pipe needs the way out, not only the diagnosis. |
| `INWZN-E02` | The input ends while the prompts run in line mode. | The end of input counts as a failed prompt, and the init refuses to write. | A line-mode prompt returns its default at end of input with no error, which once wrote a configuration with no artifact and no gate under a success message. |
| `INWZN-E03` | REF[INWZN-I01]: every prompt failure ends in the same refusal to write | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/initx` | `Infer`, `DetectGit`, `AvisoGit`, `OfferAction`, `PrecisaDescobrir`, `DetectOperator`, `AgentName`, `CommandToOpenAI`, `PromptDescobrir`, `CoverageHint`, `RenderContributing`, `ContributingSection` | apoio — the inference, the git states, the operator, the dialect and the contributing guide |
| DEP2 | `cmd/anchors/ops/init_non_interactive.go` | `runInitNonInteractive` | comando — ININT |
| DEP3 | `internal/config/config.go` | `Save`, `DefaultFile` | config — the written file |
| DEP4 | `cmd/anchors/governance` | `RenderSpecGuide` | comando — the seeded spec guide |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
