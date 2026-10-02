---
title: "The Anchors CLI"
description: "The official command-line toolchain for AI pair programming and spec-first continuous governance."
---

The **Anchors CLI** (`anchors`) is a fast, standalone binary written in Go that enforces continuous governance for AI-assisted software development.

Rather than trying to embed AI inside the tool, Anchors is designed from the ground up as **the tool that AI agents operate**:
- The AI does not need to memorize the entire framework; it queries the binary (`anchors guide`), learns what needs to be done, and executes commands.
- It is completely client-agnostic: it works with Claude Code, Cursor, Windsurf, Copilot, Gemini CLI, or human terminal developers.
- It operates strictly on text and filesystem contracts — no hidden runtime servers or heavy dependencies.

---

## 🚀 Quick Navigation

- [**Installation Guide**](/docs/cli/installation/) — Install via Homebrew, shell installer, Go toolchain, Windows, Docker, and CI/CD.
- [**Commands Reference**](/docs/cli/commands/) — Complete searchable index of all 55 CLI commands.
- [**The Workflow**](/docs/workflow/) — How to operate the daily cycle from spec to merge.
- [**The anchors.yaml**](/docs/anchors-yaml/) — Configuration schema and layer definitions.

---

## 🛠️ Essential Commands at a Glance

```bash
# 1. Initialize a new project or configure an existing one
anchors init

# 2. Check project health and structural maturity
anchors doctor

# 3. Build and query the dependency graph
anchors map build
anchors impact src/services/auth/login.spec.md

# 4. Run automated quality gates against staged files or whole project
anchors check
anchors check --all

# 5. Start background AI task coordination
anchors watch
```

---

## 🧭 Deep Dive Command Guides

- [**anchors check & verify**](/docs/cli/commands/check/) — Pipeline evaluation, strict modes, and gate filtering.
- [**anchors doctor & audit**](/docs/cli/commands/doctor/) — Ecosystem health, orphaned artifacts, and single-file audits.
- [**anchors map & impact**](/docs/cli/commands/map/) — Dependency DAG analysis, blast radius calculation, and code renaming.
- [**anchors init & new**](/docs/cli/commands/init/) — Project setup wizard, layer discovery, and artifact scaffolding.
- [**anchors flow & queue**](/docs/cli/commands/flow/) — Background watcher, task claiming, and autonomous stage deliveries.
- [**anchors freeze & thaw**](/docs/cli/commands/freeze/) — Release baseline locking and security freezes.
