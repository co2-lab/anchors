<!-- @anchors
  code: TLRCD
  layer: desenho
-->

# Tool result cache — a project's tool is not run again over the same inputs

> APPROVED on 2026-10-21 (the user: "pode seguir com as recomendações, mas D01-D03 deixa tudo opt-in").

## The cost

At MIF, after the gates ran in parallel, a `check --all` is bounded by the project's own
tools: eslint 1m43s, semgrep 26s, spellcheck 10s, typecheck, and the slow ones the CI runs.
Most of a second `check --all` re-runs them over files that did not change, and gets the
same answer.

## The design

- **A gate opts in by naming what its tool reads besides its targets**:

  ```yaml
  - name: eslint
    run: "bash scripts/anchors-eslint.sh {{files}}"
    cache:
      inputs: ["**/eslint.config.*", "pnpm-lock.yaml", "scripts/anchors-eslint.sh"]
  ```

  A tool reads more than the files it is handed — its configuration, its own version, the
  lockfile that pins its plugins, the script that wraps it. Anchors cannot know which, and a
  cache that guesses serves a stale green. With no `cache.inputs`, the tool runs every time.
- **The key is everything the verdict depends on**:
  - the gate as declared: its command, scope and arguments;
  - the targets it is handed, each with its content's hash, read as the run reads it, from
    the tree or from the index under `--index`;
  - the content of every file the `inputs` globs match;
  - the version of Anchors.
- **A hit replays the verdict** — a pass, and a failure only where the gate opts in
  (`cache.failures: true`) — and its report, marked `cached` with the time
  of the run it came from. A miss runs the tool and stores the result under the key.
- **Where it lives**: `.anchors/cache/tools/<gate>.json`, the last few keys of each gate, out
  of git, like the rest of `.anchors/`.
- **`check --no-cache` runs every tool.** In a CI (`CI` set in the environment) the cache is
  neither read nor written unless the gate opts in (`cache.ci: true`).
- **What is never cached**: a gate with no `inputs`; a tool that failed to start (missing,
  timed out), which is no verdict of the files; a judgment gate.

## What NOT to do

- Cache by modification times: a checkout resets them, and an edit restored by hand keeps a
  new time with the old content. The key is contents.
- Infer the inputs: a cache that guesses wrong serves a green over code the tool never read.

## Decisions

| Code | Question | Decided |
| --- | --- | --- |
| `TLRCD-D01` | Opt-in or default? | Opt-in, per gate, by declaring `cache.inputs` (the user, 2026-10-21). |
| `TLRCD-D02` | Is a failure cached? | Only where the gate opts in, `cache.failures: true`; by default a failure runs again (the user, 2026-10-21). |
| `TLRCD-D03` | Does the CI use it? | Only where the gate opts in, `cache.ci: true` (the user, 2026-10-21). |
