# Contributing to Anchors

Anchors governs its own repository. Every rule the tool enforces on a project is enforced
here too, by the same binary, through the same `anchors.yaml`. This guide is for anyone
changing Anchors — people and AI agents alike — and describes how that loop works day to
day.

## Ground rules

- **Everything is written in English**: identifiers, comments, specs, features, tests,
  doctrine documents, commit messages. The one exception is the translation catalog,
  `internal/i18n/locales/*.json`. The `code-language` gate enforces the part a compiler
  can see (identifiers); review holds the rest.
- **User-facing text goes through i18n.** A new message is a key in all three locales:
  `en.json`, `pt-BR.json` and `es.json`. Use `i18n.T("key", args...)`, never a bare string.
- **The workflow mode is `manual`** (see `anchors.yaml`). `anchors check` reports and
  stamps the map, but writes nothing to `issues/`. Add `--record-issues` when you do want
  the files.

## Setting up

Anchors is a single Go module (see `go.mod` for the Go version).

```sh
go test ./...                                   # the whole suite
go build -o ~/go/bin/anchors ./cmd/anchors      # a local build, reported as "dev"
anchors install-hooks                           # pre-commit, commit-msg, pre-push
```

The hooks call the `anchors` on your PATH, so that binary is the one checking your commits.
A build made from a tag should carry its version:

```sh
go build -ldflags "-X main.version=0.1.196 -X main.commit=$(git rev-parse --short HEAD)" \
  -o ~/go/bin/anchors ./cmd/anchors
```

If you change `touch`, `verify` or anything else the hooks run, install your build BEFORE
committing. Otherwise the old binary runs against your change.

## The unit of work: spec, feature, test, code

A unit lives in one directory, as four files with the same stem:

```
internal/gate/mock_stamped.spec.md    the rules (the anchor: spec comes first)
internal/gate/mock_stamped.feature    one scenario per rule, tagged with the rule's code
internal/gate/mock_stamped_test.go    the proof: each test names the rule it proves
internal/gate/mock_stamped.go         the code
```

- The spec declares its code (`code: MCSTM`) and numbers its rules: `MCSTM-B01`
  (behaviour), `-I01` (invariant), `-X01` (constraint), and so on.
- A feature scenario carries the rule's tag, for example `@MCSTM-B17 @unit-level`.
- A test proves a rule by naming it: `t.Run("MCSTM-B17: …", …)`. That name is how the map
  links execution back to the rule.
- A changed rule gets its scenario and its test in the same commit.

`anchors guide spec`, `anchors guide feature` and `anchors guide test` hold the full rulers.

## The daily loop

```sh
# 1. change the code, and the spec/feature/test that govern it
go test ./...                            # green before anything else
anchors map build                        # REQUIRED when you add a governed file
anchors check --changed a.go,b_test.go   # the gates over your change (comma-separated)
git add <the files you changed>          # explicit paths — never `git add -A`
git commit                               # the hooks run the gates and date the headers
```

- **`anchors map build`** registers new files. The pre-commit refuses a governed file that
  is not in the map, because no gate would confront it.
- **`--changed` takes a comma-separated list.** Separate arguments are not read as more files.
- **The pre-commit dates the headers.** It bumps the `updated_at` of the @anchors header of
  every staged file with a real change (`touch.pre_commit`, on by default), and prints
  which ones. `anchors touch` does the same outside a commit.
- **Never `git add -A` or `git add .`.** Other sessions and agents may be editing the same
  tree. Stage exactly what you changed.
- **Never `--no-verify`.** If a hook blocks, fix the cause.

## Testing discipline

A green test is not yet proof. Before calling a rule proven:

1. **Break the code the test claims to protect, and watch the test fail.** Remove the line
   that implements the rule, or invert a condition, or drop a case. A mutant that survives
   means the rule looks proven without being proven, which is worse than an open gap: the
   gap shows up in the report, and this does not.
2. **Restore it and watch the test pass.**
3. **Keep the mutant compilable.** If the mutated code does not build, the test run proves
   nothing. Check the output for `[build failed]`, not only for `--- FAIL`.
4. **Don't pipe test output through `tail` or `grep` and trust what's left.** Failures and
   build errors get filtered out with the noise.

The four ways a test looks like proof and is not: an assertion that matches by accident,
a function tested in isolation (not on the real path), an invariant where only one half
is checked, and a field that appears as input but is never checked in the output.

Tests of git behaviour (`touch`, hooks, ingestion) run against a real temporary repository
(`t.TempDir()` + `git init`), never against the working tree.

## Commits

- **Conventional Commits**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`,
  `build`, `ci`, `chore`, `revert`, with an optional scope: `fix(touch): …`. There is no
  `merge` type; use `chore(merge): …`.
- **The subject line is at most 100 characters.** The detail goes in the body.
- **The body says why, with the measurement.** Say what failed, where, and how you know the
  fix works. The history is how the next person understands a rule.
- **AI agents** end the message with their `Co-Authored-By:` line.

## Releasing

Every release is a tag on `main`. The `release-cli` workflow builds and publishes it with
goreleaser.

```sh
# 1. the commit is in, the suite is green on the COMMITTED state
git log -1 --format=%s                  # check it IS the commit you mean to tag
git tag -a v0.1.197 -m "v0.1.197 — what it brings"
git push origin main && git push origin v0.1.197
```

- **Tag only after checking the commit went in.** A hook that blocked leaves `HEAD` on the
  previous commit. A tag there publishes the old code under the new version, and the tag
  and release then have to be deleted.
- **Projects move to the new version** by raising `min_version` in their `anchors.yaml`.
  When a workflow template changed (`internal/initx/workflows/`), they also re-seed their
  `.github/workflows/anchors-*`: `anchors doctor --fix`, or `initx.SemeiaWorkflows`.
- **A behaviour change that reaches agents** (a new hook action, a new command they should
  use) goes in the project's `notifications.md`, which `anchors next` prints on top.

## Known traps in this repository

- **Rule letters live in three places.** Adding a rule letter (`B`, `I`, `X`, `G`, …) means
  changing:
  - `internal/config/config.go` (`DefaultRuleLetters`);
  - `cmd/anchors/ops/new_templates.go` (the section catalog);
  - `internal/testsig/code.go` (`ruleLetters`, a copy by value).

  The guard tests are `TestRuleLetters_naoDivergeDoConfig` and
  `TestCatalogoNaoUsaLetraForaDasCanonicas`.
- **Header-like text inside code.** An @anchors header is only the block at the top of the
  file (first 10 lines). Don't write test fixtures that rely on the tool ignoring a
  header-shaped string somewhere else. Keep them in constants, as `touch_test.go` does.
- **`anchors.yaml` is loaded with `KnownFields`.** A new config key needs a field in
  `internal/config/config.go`, or every project that declares it fails to load.
- **The spec describes the real code.** Prove each claim a spec makes with something that
  runs (a test, a probe), not by reading a parameter's name.

## For AI agents

- Read the guide for the artifact you are touching (`anchors guide <spec|feature|test|review>`)
  before writing it.
- Work from the gates' output: `anchors check --changed <files>` for your change,
  `anchors check --all` for the full picture. In manual mode no issue files pile up, so the
  check's output is the list.
- A defect in the tool is fixed here, with a test and a release — not worked around in the
  project where it showed up.
