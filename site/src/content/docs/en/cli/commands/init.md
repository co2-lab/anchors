---
title: "anchors init & new"
description: "Project setup from the project's own folders and language, the contributing guide, git hooks, and artifact scaffolding."
---

Starting a new project or onboarding an existing codebase into Anchors begins with `anchors init`. It reads the project as it is and asks only the decisions that are yours.

---

## 1. Project Initialization (`anchors init`)

```bash
# Interactive setup
anchors init

# For an agent or a script: first the questions, as JSON…
anchors init --non-interactive
# …then the answers, as flags
anchors init --non-interactive --artifacts=spec,feature,test,code --colocation
# or accept every default inferred from disk
anchors init --non-interactive --defaults
```

`anchors init`:

1. **Proposes no structure.** Anchors does not care how your folders are arranged — only that the layers are kept apart. Every folder that holds code is a candidate layer, named after the folder, and you keep the ones that are layers. Layer kinds such as *entry points, use cases, domain, repositories, infrastructure, presentation* are given as illustration, never as a layout to move files into.
2. **Reads the language as dialect.** The family comes from the manifest at the root (`go.mod`, `package.json`, `pyproject.toml`…) or the most frequent extension, and the test naming from your own test files (`*_test.go`, `*.spec.ts`, `test_*.py`…). It is written to `dialect.family`, so the gates read your tests and comments from the first check.
3. **Seeds the gates related to your layers** — a gate is seeded when a declared layer is of a kind it measures and the fields it presupposes are declared. The rest of the catalog that covers your layers is listed by [`anchors doctor`](/docs/cli/commands/doctor/).
4. **Seeds `CONTRIBUTING.md`** from the configuration it writes: the order of the work (spec → feature → test → code), the declared layers, the daily commands and which gates bar a commit. A `CONTRIBUTING.md` you already have is never touched; the section that would be added is shown instead.
5. **Seeds the header guide** and writes [`anchors.yaml`](/docs/anchors-yaml/).

When it seeds gates with `run:`, it also says how to write a gate step that does more than call one tool: in the project's language (`go run ./tools/gates <gate>`, `node tools/gates.mjs <gate>`, `python -m tools.gates <gate>`), not in shell — a shell script breaks on the platform.

| Flag | What it answers |
| --- | --- |
| `--artifacts` | the anchor kinds the project uses (`spec,feature,test,code,guide,plan`) |
| `--layers` | the code folders to keep as layers |
| `--colocation` | the spec, feature and test sit beside the code |
| `--gates` | seed the default gates (`true`) |
| `--header` | seed the header guide (`true`) |
| `--contributing` | seed `CONTRIBUTING.md` when the project has none (`true`) |
| `--workflow`, `--repo`, `--labels` | where the work queue lives: `local`, `manual` or `github` |
| `--governs GUIDE=tag1,tag2` | which tags a guide governs |

---

## 2. Scaffolding Artifacts (`anchors new`)

Generate an artifact with its `@anchors` header and identity already resolved. `--out` is mandatory: the artifact is born next to the unit it describes.

```bash
# A spec gets a new unique code
anchors new spec Invoice --out src/billing/invoice.spec.md

# A feature or a test gets a ref to the spec
anchors new feature Invoice --out src/billing/invoice.feature --code INVCE

# The sections a kind offers, and its presets
anchors new spec --list-sections
```

---

## 3. Git Hooks (`anchors install-hooks`)

```bash
anchors install-hooks
```

Installs the pre-commit that runs the gates over what is staged, so a commit that fails a blocking gate does not land.
