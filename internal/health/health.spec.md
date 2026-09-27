<!-- @anchors
  code: DCTRO
  updated_at: 2026-09-26
  layer: infra
-->
# Doctor — the global health check that hunts the systemic loose ends of a project

> **Code**: `DCTRO`

## Overview

The gates confront one node against one criterion, incrementally. The doctor is the GLOBAL view: it sweeps
the map, the configuration and the disk together and hunts the loose ends no single gate sees: a map that no
longer matches the disk, specs that nothing realizes, one identity owned by two domains, layers and guides
that govern nothing, kinds no gate confronts, configuration that silently does less than it says, a missing
git, plans that can never start, and verification signals the project never ingests.

It PRESENTS and records, but does not block: it runs on demand, outside the merge path. It detects and
presents; it does not arbitrate. Each finding is a warning (an open end that needs attention) or an
information (an observation, not a problem), with the check that found it, the subject and a text.

This unit is also where the doctor's other checks meet: the GitHub environment (`GHEGT`), the pending
decisions (`PNDCP`), the governance opportunities (`GVOPG`) and the spec sections (`SPSCS`) are run by the
same diagnosis and sorted into the same report.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the map | a built graph of nodes and edges | a nil graph | the caller, which builds or loads the map before diagnosing |
| the configuration | the loaded project configuration, with its layers, gates and workflow | a nil configuration | the caller, which loads the configuration first |
| the project root | the directory the map's paths are relative to | — | the caller |
| the plans' needs | plan paths as written in each plan's header | — | the map, which keeps the declared needs even when they point nowhere |

## Effects

### Report

| Effect | Description |
| --- | --- |
| `DCTRO-B01` | `Diagnose`: The diagnosis runs every check and returns the findings sorted by check and then by subject, with the number of nodes, edges and declared layers. |
| `DCTRO-B02` | The warnings of a report are only its findings of warning severity, in order. |

### Map fidelity

| Effect | Description |
| --- | --- |
| `DCTRO-B03` | A node whose file is missing on disk is a warning `no-fantasma`. |
| `DCTRO-B04` | An edge whose source or target is not a node of the map is a warning `aresta-morta` on "source → target". |

### Orphans

| Effect | Description |
| --- | --- |
| `DCTRO-B05` | A spec that is the source of no edge is a warning `spec-sem-realizacao`: a requirement with no incarnation. |
| `DCTRO-B06` | A spec without an identity code is a warning `identidade-ausente`. |

### Duplicate identity

| Effect | Description |
| --- | --- |
| `DCTRO-B07` | A code owned by units of DIFFERENT domains is a warning `identidade-duplicada` on the code, naming the owning units. A domain is the `features/<name>` segment of the path, or its top directory outside `features/`. |
| `DCTRO-B08` | Units of the same domain sharing one code are deliberate and give no finding. |
| `DCTRO-B09` | Only specs, features, tests and code own a code; documents, guides and plans only reference it and never collide. |
| `DCTRO-B10` | The files of one unit (its spec, feature, test and code, same folder and name) count as one owner. |
| `DCTRO-B11` | A file that declares its codes as shared (`@anchors-shared-code`) is not an owner of them. |

### Structure

| Effect | Description |
| --- | --- |
| `DCTRO-B12` | A declared layer with no node of its kind is a warning `camada-vazia`. |
| `DCTRO-B13` | A guide that is the source of no `governs` edge is a warning `guide-sem-governo`. |
| `DCTRO-B14` | A spec, feature or test kind present in the map that no gate confronts is a warning `kind-sem-gate`; code and documents may legitimately have no gate of their own. |

### Configuration

| Effect | Description |
| --- | --- |
| `DCTRO-B15` | A perspective in a gate's `skip_on` other than `change` or `all` is a warning `skip-on-invalido` on the gate, quoting the value. |
| `DCTRO-B16` | A gate whose required tool is not on the PATH is a warning `ferramenta-ausente` on the gate, carrying its install hint when there is one. |

### Git

| Effect | Description |
| --- | --- |
| `DCTRO-B17` | Without the git binary, the warning `git-ausente` is on `git` and asks to install it, never to initialize a repository. |
| `DCTRO-B18` | With git but no repository, the warning `git-ausente` asks to initialize one, and does not point at the PATH. |
| `DCTRO-B19` | A `.git` in the root or in any ancestor counts as a repository, whether it is a folder or a file (a worktree or submodule pointer). |
| `DCTRO-B20` | In GitHub mode the missing repository's text says the work queue has nowhere to come from. |

### Plans

| Effect | Description |
| --- | --- |
| `DCTRO-B21` | A plan whose `needs` names a plan that does not exist is a warning `needs-quebrado` on it, naming the missing target. |
| `DCTRO-B22` | A cycle of `needs` between plans gives one warning `needs-ciclo`, showing a path through the cycle. |
| `DCTRO-B23` | A chain of needs in order, or a project without plans, gives no finding. |

### Signals

| Effect | Description |
| --- | --- |
| `DCTRO-B24` | Tests in the map with no ingested result at all, and code with no ingested coverage at all, are each a warning `sinal-ausente`. One node with the signal is enough. |
| `DCTRO-B25` | Code with no mutation signal is an informational `sinal-ausente`; code where only some files have it is an informational `sinal-ausente` on the partial mutation, carrying how many of how many. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DCTRO-I01` | A project with nothing wrong gives the doctor no warning: each check that finds nothing adds nothing. | a small project with its files on disk, realized specs, a governing guide, gated kinds, a repository and every signal is diagnosed |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DCTRO-X01` | Code without a spec is never reported. | Utilities, constants and pure libraries need no spec; the direction that matters is the reverse, a spec without code. Reporting every file without a spec would be noise. |

## Errors

none — every missing file, tool, repository or plan the doctor meets is itself a finding it presents; the diagnosis has no failure of its own to report.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `Graph`, `Nodes`, `Edges` | mapa — the map being diagnosed |
| DEP2 | `internal/config/config.go` | `Config`, `Gates`, `Layers` | config — layers, gates and workflow |
| DEP3 | `internal/i18n/i18n.go` | `T` | apoio — localized finding texts |
| DEP4 | `internal/health/github_environment.go` | `checkGitHubEnv` | infra — the GitHub environment checks |
| DEP5 | `internal/health/pending_decision.go` | `checkPendingDecisions` | infra — the open decisions |
| DEP6 | `internal/health/governance.go` | `checkGovernanceOpportunities` | infra — the governance suggestions |
| DEP7 | `internal/health/spec_sections.go` | `checkSpecSections` | infra — the recommended spec sections |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
