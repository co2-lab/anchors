<!-- @anchors
  code: OPDTP
  updated_at: 2026-10-08
  layer: infra
-->
# OperatorDetection — tell whether a person or an AI is running `init`, and whether the discovery phase is still to be done

> **Code**: `OPDTP`

## Overview

The discovery phase of a project (see `anchors guide project`) needs an AI conducting an interview, and
what Anchors should SAY at the end of `init` depends on who reads it. For a person, the output is an
instruction: here is what is missing and here is how to start, including opening the AI with the prompt
ready. For an AI, the output is the work order itself: read the guide and conduct the interview with the
user, in this conversation. One text for both fails both.

The unit decides who the operator is from two pieces of evidence the caller passes in: whether there
is a terminal, and the environment. An environment variable of a known agent is an EXPLICIT declaration
and outweighs the terminal; the absence of a terminal is only a hint (a plain pipe produces it too), but
it still means nobody is typing the answers, so the AI text is the more useful of the two. It also names
the detected tool, so the message can say "open Claude Code" instead of "open your AI", and builds the
command that opens that tool with the discovery prompt.

Finally it decides whether the discovery phase has still to happen: only when there is nothing to infer
from (no code, spec, feature or test) and no project description written yet.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| whether there is a terminal | a yes/no measured by the caller | — | the caller, which inspects its own standard streams |
| the environment | any lookup of variable names; none means the process environment | — | this unit: a missing lookup falls back to the process environment |
| the inference proposal | the result of the disk inference, or nothing | — | this unit: a missing proposal counts as a project with nothing inferred |
| the project root | an existing folder | a path that does not exist | the caller |

## Effects

| Effect | Description |
| --- | --- |
| `OPDTP-B01` | A known agent variable in the environment makes the operator an AI (`DetectOperator`), even when there is a terminal. |
| `OPDTP-B02` | The generic agent variable `AI_AGENT` makes the operator an AI, even when there is a terminal, although it names no tool. |
| `OPDTP-B03` | With no terminal and no agent variable, the operator is an AI. |
| `OPDTP-B04` | With a terminal and no agent variable, the operator is a person. |
| `OPDTP-B05` | The tool (`AgentName`) is named after the known agent variable that is set; with none set the name is empty. |
| `OPDTP-B06` | When several known agent variables are set, the one whose variable name sorts first names the tool, so the same environment always gives the same name. |
| `OPDTP-B07` | The discovery phase is due (`PrecisaDescobrir`) only when the inference found no code folder, spec, feature or test AND there is no project description at the root; a missing proposal counts as nothing found. |
| `OPDTP-B08` | The project description is recognised (`HasProjectMD`) under the spellings `PROJECT.md`, `project.md` and `Project.md`. |
| `OPDTP-B09` | The command that opens the AI (`CommandToOpenAI`) exists only for Claude Code, Gemini CLI and Aider (Aider takes the prompt as a message option); for any other tool, or none, there is no command. |
| `OPDTP-B10` | The discovery prompt points at `anchors guide project` for the interview and back at `anchors init` for the next step, instead of summarising the interview. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `OPDTP-I01` | The command to open the AI is a list of arguments whose last element is the whole discovery prompt, verbatim — never a shell line the prompt was pasted into. | builds the command for each supported tool and compares the last argument with the prompt |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `OPDTP-X01` | Does not offer a command for a recognised tool whose command-line invocation is not stable (Cursor, Codex). | Offering a command that does not exist is worse than offering nothing. |

## Errors

none — a missing description file is the "not written yet" answer and a missing environment lookup falls back to the process environment; nothing the unit reads can fail in a way it has to report.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
