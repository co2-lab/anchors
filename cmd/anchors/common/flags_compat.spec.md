<!-- @anchors
  code: FLALF
  updated_at: 2026-09-26
  layer: comando
-->
# FlagAliases — a renamed command flag keeps answering to its old name

> **Code**: `FLALF`

## Overview

Command flags were renamed (Portuguese names became English ones). Scripts, hooks and habits still pass
the old names, and breaking them on an upgrade would punish exactly the people who automated the tool.
This unit lets a command declare that an old flag name is an alias of the current one: the old name is
still accepted, it is hidden from the help so nobody learns it anew, and it carries a deprecation note
that points at the current name.

After parsing, the command resolves its aliases: the value typed under the old name is copied to the
current flag, unless the current name was typed too, in which case the current name wins. A value the
current flag cannot hold is reported under the name the person actually typed.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the current flag name | a flag already declared on the command | a name the command does not declare | this unit: declaring the alias stops the program at start, naming the flag |
| the old flag name | a name not yet declared on the command | a name already in use | the command author (the flag library rejects a duplicate declaration) |
| the alias pairs to resolve | pairs of current and old names | pairs where either flag does not exist on the command | this unit: such a pair is ignored |

## Effects

| Effect | Description |
| --- | --- |
| `FLALF-B01` | `AliasDeFlag` and `ResolveAliases`: a value passed under the old name reaches the current flag once the aliases are resolved. |
| `FLALF-B02` | `AliasDeFlag`: the alias of a switch flag is itself a switch, so the old name can be passed without a value; the alias of any other flag takes a text value. |
| `FLALF-B03` | The old name is hidden from the help and marked deprecated with a note that points at the current name. |
| `FLALF-B04` | `ResolveAliases`: when both the old and the current names are passed, the current name's value is kept. |
| `FLALF-B05` | A pair to resolve whose old or current flag does not exist on the command is ignored without error. |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `FLALF-X01` | Does not copy anything until the command asks for the aliases to be resolved after parsing. | Only after parsing is it known which names were typed; copying earlier would overwrite the current name's value. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `FLALF-E01` | The value typed under the old name cannot be held by the current flag. | An error that names the OLD flag. | The person typed the old name; an error about a flag they never typed would send them looking in the wrong place. |
| `FLALF-E02` | An alias is declared for a current flag the command does not have. | The program stops at start with a message naming the missing flag and the command. | It is a programming error in the command's wiring, and it must surface the first time the command is built, not when a user passes the flag. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | — | the command-line flag library | external — command and flag parsing |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
