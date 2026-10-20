<!-- @anchors
  code: LTNGD
  layer: desenho
-->

# Late ingestion — a report is stamped with what its run measured, not with today's tree

> APPROVED on 2026-10-20 (the user: "aprovada sua proposição"), from MIF's T159.

## The defect

`anchors ingest` stamps a proof with the revisions the map holds at ingestion: the test's
`AtRev`, the revision of every file of its closure (`ClosureRev`), each spec's
`ProvenRevBySuite`. A report ingested right after its run is right. A report ingested later
is not: it reads as if it had run against today's code, and a screen that changed after the
run keeps a fresh proof. `evidence-fresh` cannot see it — the stamp it compares against was
written after the change. Nothing reads the time of the run: Maestro's JUnit carries none,
and the one "older than the file" check (`markPredating`) is lcov's.

## The design

- **The run is placed in time.** A manual ingestion is the case: `anchors test` ingests the
  run it just made, and keeps its behaviour. The time is `--ran-at` (a time, or the commit the
  run measured), else the report's own `timestamp`, else the report file's modification time.
- **What the run measured is the commit of that time** — the last commit at or before it
  (`git rev-list -1 --before`), or the one `--ran-at` names.
- **Each file the proof stands on is confronted with its content at that commit**, by the
  evidence the map already uses (`scan.EvidenceOf`: comments and the header are no evidence).
  The same evidence today → stamped with today's revision, the proof fresh. Another evidence,
  or no file at that commit → stamped `changed-after-run`, which no revision ever equals: the
  proof is stale, naming the file. A file whose modification time is not after the run is
  unchanged and is not read. Without git, or with no commit before the run, every file changed
  after the run is stamped changed.
- **The files are the test, its closure, and the specs whose scenarios it proves.** A spec
  that changed after the run stamps its suite's proof `changed-after-run`.
- **A code renamed after the run is remapped, not lost** (`--alias OLD=NEW`, repeatable, and
  `--alias-file` with one per line): a unit code (`PHA1=PHA1P`) renames each of its scenarios
  and paths, a scenario code renames that one. The report is read with the new codes — case
  names and file paths —, and so is each file's content at the run's commit before it is
  confronted: a rename changes no behaviour, so a file whose only change is the rename keeps
  the proof fresh.
- **The signal records it**: `ran_at`, the time of the run, and `remapped`, the aliases that
  touched the report.

## What NOT to do

- Compare times instead of contents: the rename itself rewrites the files after the run, and
  every remapped proof would stale; a checkout resets modification times.
- Trust a report older than the tree because it is green: that is the defect.

## Decisions

| Code | Question | Decided |
| --- | --- | --- |
| `LTNGD-D01` | Is a late report judged by time or by content? | By content at the run's commit; time only says which commit (2026-10-20). |
| `LTNGD-D02` | Is a code rename evidence? | No — under an alias the old code reads as the new one before the confrontation (2026-10-20). |
| `LTNGD-D03` | Does `anchors test` change? | No — it ingests the run it just made (2026-10-20). |
