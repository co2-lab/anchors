<!-- @anchors
  code: GRBLG
  updated_at: 2026-10-07
  layer: mapa
-->
# GraphBuild — projecting the declared structure onto the scanned files: one node per file, and the relations between them

> **Code**: `GRBLG`

## Overview

The map is the projection of the declared structure (the configuration) onto the real files. This unit
takes the scanned files and the configuration and assembles the graph: one node per file, and the
relations between them. It knows nothing of any language or file name — every decision comes from the
configuration, and the dates of each file come from the caller, so the unit neither reads the disk nor
calls version control.

A node's identity is decided in a fixed order. What the author DECLARES in the header wins: it is where
the author says whose file this is. A derived artefact (a test, a feature, the code) with no header takes
the declared identity of its anchor — the spec with the same stem in the same directory — because a test
proves a unit and does not own one; reading its text would take test data for a declaration. Only when
neither exists is the identity inferred from the first rule code in the text. A vendored copy that Anchors
seeded and still owns gets no local identity at all.

Relations come from four sources:
- Co-location. From each anchor file, the configured path templates give where its derived files live,
  and the unit is linked: the spec specifies the code (every file it matches), is covered by the feature,
  which is tested by the test; with no feature, the spec is tested by the test directly. A template can be
  overridden per layer of the UNIT (the layer the spec declares, not the layer its own file matched), or
  replaced whole for one identity code. The anchor's directory and name are data and are matched
  literally; only the template's own wildcards expand.
- Scenario codes. A spec or feature and a test that carry the same rule code, in different directories,
  are linked by an inferred tested-by relation — but only for codes whose root is the file's own identity.
  Citing a sibling unit's rule is not declaring it.
- Governance. Each governing rule links its guide to every file of every layer carrying the rule's tag,
  never to the guide itself.
- Declarations. A spec's dependency table, a plan's seeds and needs, and a spec's realized doctrine and
  flag scenarios become relations — only when the target exists. A relation to nothing would be an anchor
  that lies.

Finally the graph is sorted, so the committed file does not change between runs, and a rebuild keeps what
the previous graph already knew: the stamps and judgments of the relations that survived, and the signals
of the nodes whose content did not change.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the scanned files | every file the scan matched to a layer, with its kind, revision, header and codes | — | the scan decides which files exist and their kinds |
| the configuration | layers, derived templates (possibly absent), governing rules | templates naming variables other than dir, name, ext and module | the configuration loader; an absent derived section simply yields no co-location |
| the dates per path | a map from path to last-change date, or nothing | — | the caller, from version control; with nothing, dates stay empty |
| the previous graph | a graph loaded from disk, or nothing | — | the caller; with nothing, nothing is carried over |

## Effects

### Nodes and identity

| Effect | Description |
| --- | --- |
| `GRBLG-B01` | Every scanned file becomes one node carrying its path, kind, revision and layer, the date the caller gave for that path, and its layer's tags and regime. |
| `GRBLG-B02` | The identity declared in a file's header wins over everything else, including an anchor's identity and codes cited earlier in the text. |
| `GRBLG-B03` | A derived file with no header takes the declared identity of the anchor with the same stem in the same directory; another unit in that directory does not receive it. |
| `GRBLG-B04` | With neither header nor anchor, the identity is the root (`RuleRoot`) of the first rule code in the text, and empty when there is none. |
| `GRBLG-B05` | A vendored file that Anchors still owns has no local identity, whatever its text or header says. |
| `GRBLG-B06` | A node is marked as having a declared identity only when its header declares one. |

### Co-location

| Effect | Description |
| --- | --- |
| `GRBLG-B07` | The files of one unit, named after the anchor stripped of its artifact suffix (`StemOfAnchor`), are linked by the templates: the spec specifies the code, the spec is covered by the feature, and the feature is tested by the test. |
| `GRBLG-B08` | When the anchor is the code, the spec is found from it as a derived file and the relations keep the same direction, from the spec down. |
| `GRBLG-B09` | With no feature, the spec is tested by the test directly. |
| `GRBLG-B10` | A per-layer override applies when its layer is the UNIT's layer — the layer a spec declares in its header wins over the layer its own file matched — and replaces the templates of the layers it lists. |
| `GRBLG-B11` | A per-code override, for the anchor with that identity, replaces the templates of the kinds it declares — an empty list declaring none — and every kind it does not declare keeps the default template (after the per-layer overrides). |
| `GRBLG-B12` | The anchor's directory and name are matched literally (`GlobEscape`), brackets included, while a template's own wildcards expand to every existing match, each of them linked. |

### Scenario and declared relations

| Effect | Description |
| --- | --- |
| `GRBLG-B13` | A spec or feature and a test in different directories that carry the same rule code of their own unit are linked by an inferred tested-by relation; files in the same directory, pairs already linked by co-location, and codes cited from another unit are not. |
| `GRBLG-B14` | A governing rule links its guide to every file of every layer carrying the rule's tag, never to the guide itself, and never to files of layers without the tag. |
| `GRBLG-B15` | Each row of a spec's dependency table becomes a depends-on relation to the named file, carrying the method and the row's code, when that file exists. |
| `GRBLG-B16` | A plan's seed that names a path links to that exact path only; a bare file name links to the one file with that name, and to nothing when several share it. |
| `GRBLG-B17` | A plan's need becomes a needs relation only when the needed plan exists. |
| `GRBLG-B18` | A realized doctrine rule or a flag scenario is resolved by its unit code to the doctrine or flag file declaring that code, carrying the local rule and the target rule; an unknown unit code links nothing. |

### Order and rebuild

| Effect | Description |
| --- | --- |
| `GRBLG-B19` | The nodes are sorted by path and the relations by source, target and type. |
| `GRBLG-B20` | A rebuild (`PreserveStamps`) carries over the stamp and the judgments of every relation that survived with the same type, source and target. |
| `GRBLG-B21` | A rebuild carries over a node's signal, and its declarations of kept evidence, only when the node's revision did not change. |
| `GRBLG-B22` | A support file becomes a node marked as support, keeping its kind and its edges like any other node. |
| `GRBLG-B24` | A file's own header `code:` is its `file_code`. A file that also `ref:`s a unit is not that unit's owner: its unit code comes from the anchor beside it, as without a `code:`, and its identity is not declared; a file with a `code:` and no `ref:` — a spec, or a file with no unit around it — owns its unit, and both codes are the same. |
| `GRBLG-B25` | Linking the incarnations of a scenario code across directories, a file's units are the ones it `ref:`s; only a file that refs none has its own `code:` as its unit — a test with a code of its own and `ref: X` is linked to X's feature by X's codes, and never by a code it only cites. |
| `GRBLG-B26` | The code's flags become edges to the file whose own code they name: a dependency flag is a declared `depends-on` carrying the import's symbols, a navigation flag a `navigates-to` carrying the rule that triggers it; a code no file owns makes no edge, and a waiver makes none. (`flagEdges`) |
| `GRBLG-B23` | Filling signals from another map gives each node that has none the signal that map holds for the same file at the same revision, with its declarations of kept evidence; a node that has a signal keeps it, and a file at another revision gets nothing. (`FillSignals`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `GRBLG-I01` | The same files and configuration always build the same graph, whatever order the files arrive in. | builds from the files in two orders and compares the graphs |
| `GRBLG-I02` | The edges are in a total order — from, to, type, dependency code, method, origin —, so two graphs with the same edges are written the same whatever order they were built in. | builds the same two dependency rows of one spec on one file in both orders and compares |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `GRBLG-X01` | The build does not ask version control or the disk for anything: a file's date is the one the caller gave, and empty when none was given. | Keeping version control out keeps the map buildable from any list of files, and testable. |
| `GRBLG-X02` | No relation points to a file that was not scanned. | A relation to a missing file is an anchor that lies; the gates that confront the declaration report the missing target instead. |

## Errors

none — every branch that skips a relation is normal flow (a target that does not exist, a tag no layer carries, an ambiguous name), stated as behaviours B13–B18 and X02; assembling the graph in memory cannot fail.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/scan/scan.go` | `File` | scan — the scanned files, their headers, codes and declarations |
| DEP2 | `internal/config/config.go` | `Config` | config — layers, derived templates and governing rules |
| DEP3 | `internal/mapx/model.go` | `Graph`, `Node`, `Edge` | mapa — the graph assembled |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
