<!-- @anchors
  code: INCHN
  updated_at: 2026-09-19
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
answer. One needs the GRAPH, because the question crosses the triad — a feature against
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
| `INCHN-B15` | An executable test script steps aside too, by a different path: its format belongs to the runner, and its identity is in the file name. |
| `INCHN-B16` | A guide with no compliance-points section, or with the section and no item in it, fails — the AI judgment gate would otherwise fall back on vague heuristics. |

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
