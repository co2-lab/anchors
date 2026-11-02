<!-- @anchors
  code: DLPLY
  layer: desenho
-->

# Dialect per layer — a polyglot project declares each lexicon once, and each layer names its own

> APPROVED on 2026-11-02 (the user: "D01. recusa / D02. sim / D03. sim"), from SS (a monorepo
> with TypeScript apps and a Go backend). The user's direction: the dialect is declared per
> layer, the project only lists the dialects it uses, and each layer points at one.

## The problem

`dialect:` is one block per project. Every gate reads it (`cfg.DialectFor()`, 58 call sites over
29 files), so a project can speak one language. SS is moving its backend to Go in the same
repository as its TypeScript apps (`apps/`, `packages/` in TS; `services/` in Go), because the
gates must keep confronting app spec ↔ OpenAPI contract ↔ Go handler ↔ test. With one lexicon,
the TS patterns read the Go code, match nothing, and the gates pass green over code they never
saw — the failure the dialect exists to prevent.

## The design

- **The project lists its dialects, each declared once, by name:**

  ```yaml
  dialects:
    ts:
      family: ts
      exported_func: "..."
      tests: { pattern: "..." }
    go:
      family: go
      exported_func: "..."
      tests: { pattern: "..." }

  layers:
    screen:      { pattern: "apps/**/screens/*.tsx", kind: code, dialect: ts }
    handler:     { pattern: "services/**/handler.go", kind: code, dialect: go }
    service-go:  { pattern: "services/**/service/*.go", kind: code, dialect: go }
    unit-test-go: { pattern: "services/**/*_test.go", kind: test, dialect: go }
  ```

  A dialect is what it declares, as today (DLCTI-B01): nothing of one dialect fills a field
  another left empty. Every Go layer names `go` and repeats nothing.
- **A file's dialect is its layer's.** A file already belongs to exactly one layer: two patterns
  matching the same file are settled by `priority`, else by the existing tie-break, and `check`
  says where the heuristic decided. So the case of a file matching two layers of different
  dialects is the case the classification already settles: the layer that wins it gives its
  dialect too. No second matching rule, no globs of its own.
- **The single-dialect form stays.** `dialect: { ... }` keeps meaning "the one dialect of the
  project", which every layer has. Nothing changes for a project with one language: no
  migration is needed, and `anchors migrate` does not rewrite it.
- **Several dialects, and a layer with none: the load refuses.** With `dialects:` declaring more
  than one, a `code` or `test` layer that names no dialect refuses the configuration, naming the
  layer. A default would read a new Go layer with the TS lexicon and pass it green; a refusal
  costs one line. Layers that hold Anchors' own formats — `spec`, `feature`, `doc`, `guide`,
  `plan` — need no dialect: a spec is Markdown and a feature Gherkin in every layer (DLCTI-B45).
  With a single entry in `dialects:`, it is every layer's, as `dialect:` is.
- **What is about no file stays the project's.** `gherkin_language` (every feature is written
  and read in one language) and the naming conventions (`set_promise`, `set_slice`) are read
  from the project, not from a dialect. A question a gate asks with no file in hand is answered
  by the project-level fields only; one that needs a lexicon must say which file it reads.
- **Every reader asks with the file.** `DialectFor()` becomes `DialectAt(path)`, resolved through
  the file's layer. The gates already hold the path (`n.ID`, the test's path, the target). A
  file outside every layer has no dialect; the readers that reach such files today are listed
  during the implementation, and each either skips them or says it could not read them — never
  reads them with some layer's lexicon.
- **Imports are resolved with the importer's dialect** — its extensions, its aliases, scoped ones
  included (DLCTI-B33). A file of another dialect is reached only where that resolution lands on
  it, and is then read with its own. A TS screen never imports a Go handler: the bridge between
  them is the contract (OpenAPI), which the gates already confront from both sides.
- **The tests are listed per dialect.** `project_tests.go` reads one source today. Each dialect's
  `tests` source (a pattern, or a script) lists the test files of the layers that name it, and
  the results are united. A dialect with no source leaves its tests unread, and the gates answer
  for them as for a project that declares none.
- **A waiver holds for its dialect.** A dialect's `opt_out` waives the field for the files of the
  layers that name it, and no other.

## What NOT to do

- Scope dialects by glob of their own (`dialect.scoped: [{files: ...}]`). It is a second way of
  saying which files are which, beside the layers, and the two drift: a file in a Go layer read
  with a TS glob's lexicon. The layer already says what a file is; it says what it is written in.
- Default a layer with no dialect to the first or the "main" one. It is the silent green again.
- Merge two dialects' fields for one file. One file, one lexicon.

## The SS draft

The rules SS drafted in `internal/config/dialect.spec.md` (DLCTI-B37–B47, I03, X02, D02, in the
working tree, uncommitted) scope by path. Under this design:

- kept, reworded for the layer: B37 (an entry is what it declares), B40 (the Gherkin language is
  the project's), B43 (a question about a file is answered by its dialect), B44 (imports with the
  importer's dialect), B45 (Anchors' own formats), B46 (opt-out per dialect), B47 (tests per
  dialect), I03 (no file read with another dialect's lexicon), X02 (adding a dialect changes only
  the layers that name it);
- dropped: B38, B39, B42 (glob matching, path normalization, a pattern's `files` inside an
  entry): the layer's classification replaces them;
- replaced: B41 by the load's refusals here (a layer naming an undeclared dialect; a `code` or
  `test` layer with none when several are declared; `gherkin_language` inside a dialect);
- D02 stays open (below).

## Cost

- `config`: `dialects:`, `Layer.Dialect`, `DialectAt(path)`, the load's refusals.
- The 58 call sites move to `DialectAt(path)`; the few with no file in hand read the project's
  fields.
- `project_tests.go`: one source per dialect.
- `init`/`migrate`/doctor: offer a family's example per dialect (D02).
- One probe at SS before the release: the TS layers read exactly as today (X02), the Go layers
  read with the Go lexicon.

## Decisions

| Code | Question | Decided |
| --- | --- | --- |
| `DLPLY-D01` | A `code`/`test` layer with no dialect when several are declared: refuse, or default? | Refuse the load, naming the layer (the user, 2026-11-02). |
| `DLPLY-D02` | Do `init`, `migrate` and the doctor offer a family's example for a new dialect, and how do they find the layers of a language with no dialect yet? | The doctor names the layers whose files' extensions no declared dialect's family reads, and offers that family's example; `init` writes one dialect per family it finds (the user, 2026-11-02). |
| `DLPLY-D03` | Do `gherkin_language` and the naming conventions stay project-level? | Yes: they speak of features and names, not of a language's syntax (the user, 2026-11-02). |
