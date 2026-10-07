# language: en
# @anchors
#   code: CHFTA
#   ref: CGPCH
#   updated_at: 2026-10-07
#   layer: feature

@CGPCH
Feature: CheckGatePipeline — confronts the map's nodes against the declared gates, records the verdicts and reports the profile

  @CGPCH-B01 @unit-level
  Scenario: The full sweep confronts every node of the map
    Given a map with two code files and one spec, and four declared gates
    When the check runs over everything without recording
    Then the header reads "check --all — 3 nodes, 4 gates"

  @CGPCH-B02 @unit-level
  Scenario: The incremental check confronts only the impact path of the change
    Given a map where the spec a.spec.md specifies a.go, and b.go is unrelated
    When the check runs over the changed spec
    Then it confronts 2 nodes, and b.go is not confronted

  @CGPCH-B03 @unit-level
  Scenario: The tests that stamp a changed module enter its impact path
    Given a module, a test whose contract stamp points at it and a test with no stamp
    When the impact path of the module is resolved
    Then the stamping test is in it and the other test is not

  @CGPCH-B70 @unit-level
  Scenario: The pieces of the changed file's unit enter its impact path, for a Go unit as for a TypeScript one
    Given a map with pkg/foo.go, pkg/foo_test.go, pkg/foo.spec.md, pkg/foo.feature, pkg/bar.go, web/x.ts, web/x.test.ts, web/x.spec.md, py/baz.py and py/baz_test.py, and no edge
    When the impact path of each of pkg/foo_test.go, pkg/foo.go, pkg/foo.spec.md, web/x.test.ts, py/baz_test.py and py/baz.py is computed
    Then each reaches the other pieces of its own unit, and pkg/bar.go never enters

  @CGPCH-B04 @unit-level
  Scenario: Changed paths are normalised to the map's form
    Given an absolute path under the root, a path starting with ./ and a path with ..
    When the changed paths are normalised
    Then they read src/a.go, b.go and d.go

  @CGPCH-B05 @unit-level
  Scenario: Files the project does not govern are recognised as not governed
    Given an empty map, a package.json and an issue file under issues/
    When the check selects the nodes for each of them
    Then each one is reported as not governed

  @CGPCH-B06 @unit-level
  Scenario: An ungoverned file does not taint a batch
    Given a governed file in the map, package.json and yarn.lock
    When the check selects the nodes for the governed file with package.json, and then for package.json with yarn.lock
    Then the first batch confronts one node with the scope "--changed (2 files)", and the second is not governed, naming "package.json (and 1 more)"

  @CGPCH-B67 @unit-level
  Scenario: A plan's progress companion is not governed
    Given an empty map, a plan under plans/ and its -progress.md companion, with plans a governed layer
    When the nodes are selected for each of them
    Then the progress companion is not governed and the plan is a governed file outside the map

  @CGPCH-B07 @unit-level
  Scenario: Phase and category select the gates charged
    Given a gate declared for pre-push only, and no gate of category nothing-declares-this
    When the check runs in the pre-commit phase, and then for that category
    Then the pre-push gate is not charged, and the empty slice prints "no gate to run for this slice"

  @CGPCH-B08 @unit-level
  Scenario: A gate without skip_on runs in both perspectives
    Given a gate with no skip_on declaration
    When the gates are filtered for changed files and for the full sweep
    Then the gate is kept in both

  @CGPCH-B09 @unit-level
  Scenario: A gate that skips the change perspective runs only on the full sweep
    Given a gate declaring skip_on change and a gate with no declaration
    When the gates are filtered for changed files and for the full sweep
    Then only the undeclared gate is kept for changed files, and both for the full sweep

  @CGPCH-B10 @unit-level
  Scenario: A gate that skips the full sweep runs only on changed files
    Given a gate declaring skip_on all
    When the gates are filtered for the full sweep and for changed files
    Then it is dropped from the full sweep and kept for changed files

  @CGPCH-B11 @unit-level
  Scenario: A gate that skips both perspectives is switched off
    Given a gate declaring skip_on change and all
    When the gates are filtered for each perspective
    Then it is dropped from both

  @CGPCH-B12 @unit-level
  Scenario: A commit message marker with a reason waives a gate
    Given a commit message carrying "[skip-code-flagged: the linter is broken upstream]"
    When the check runs with that commit message file
    Then the code-flagged gate does not run and the header counts 3 gates

  @CGPCH-B13 @unit-level
  Scenario: A waiver in the environment drops the gate and says why
    Given ANCHORS_SKIP_RULES set to "code-flagged=the linter is broken upstream"
    When the check runs over everything
    Then it prints "○ waived: code-flagged — the linter is broken upstream" and counts 3 gates

  @CGPCH-B14 @unit-level
  Scenario: The deterministic mode drops the judgment gates
    Given a project whose only gate is a judgment gate
    When the check runs in deterministic mode
    Then it prints "no deterministic gate to run"

  @CGPCH-B15 @unit-level
  Scenario: The issue policy follows the workflow mode
    Given the local, manual and github modes, in and out of CI, with and without the record-issues flag
    When the issue policy is decided
    Then local always writes, github writes only in CI or when asked, manual only when asked

  @CGPCH-B16 @unit-level
  Scenario: The confronted edges are stamped at the current revisions
    Given a spec at revision s1 specifying code at revision r1
    When the check runs over the changed spec
    Then the edge carries a stamp validated from s1 to r1

  @CGPCH-B17 @unit-level
  Scenario: The no-record mode leaves the map untouched
    Given a project with a map
    When the check runs over everything without recording
    Then the map file is byte for byte the same

  @CGPCH-B18 @unit-level
  Scenario: A blocking failure is filed only when issues are on
    Given a blocking failure of docs-fresh on a.spec.md
    When the check records it with issues off, and then with issues on
    Then the map is saved both times, and an issue exists only the second time

  @CGPCH-B19 @unit-level
  Scenario: Passes resolve, decisions and debts open in their folders
    Given an open violation of header-conforms on a.ts, an open decision on b.spec.md, an assumed debt and a plain pending
    When the check records the pass, the decision, the debt and the pending
    Then the violation is done, the decision is in todo owned by the user, the debt is in future and the pending opened nothing

  @CGPCH-B20 @unit-level
  Scenario: The full check closes the violations it did not reproduce
    Given an open violation that this run's results do not reproduce
    When the check records a full run, and then an incremental one
    Then the full run closes it and the incremental one leaves it in todo

  @CGPCH-B72 @unit-level
  Scenario: A record that fails is warned about and the check still reports
    Given a project whose map lives in a folder the check cannot write
    When the incremental check runs and records, once with the project's own map and once with the locked one
    Then the successful record raises no warning, and the failed one warns "failed to record" on the error output while the verdict is still printed

  @CGPCH-B75 @unit-level
  Scenario: When the check writes no issue it says why
    Given a project in manual mode, and a project in github mode run locally
    When each check records with issues off
    Then the first says no issue was written because of the manual mode, and the second says the board issues are left to CI

  @CGPCH-B76 @unit-level
  Scenario: In github mode the issues go to the board, never to the local folders
    Given a project in github mode with a label, and a board command that logs its calls
    When the full sweep runs locally, and then asked to record issues
    Then the local run leaves the board untouched, the recording run calls the board, and no issue folder is written

  @CGPCH-B77 @unit-level
  Scenario: The record summary counts what the record did
    Given records that open a blocking failure, open a decision, resolve a decision, and close a violation a full check did not reproduce
    When each record is written
    Then each summary counts the issues opened and resolved, the closing is announced only when a violation was closed, and the debt line appears only with a debt

  @CGPCH-B21 @unit-level
  Scenario: A pending judgment becomes one task in the local queue
    Given a local project with a judgment gate over a.go
    When the check runs over everything
    Then the queue holds one task judge-rule-kept-a of the judgment kind

  @CGPCH-B69 @unit-level
  Scenario: A judge task suggests the review stage, a verb the work command composes
    Given a local project with a judgment gate over a.go, and a queue holding a judge task with the legacy judge verb
    When the check runs over everything, and the queued judgments are counted
    Then the task suggests "review", which the work command accepts, its reason reads "anchors judge a.go --gate rule-kept", and the legacy task counts as a judgment while a watch review task does not

  @CGPCH-B22 @unit-level
  Scenario: A queued judgment bars the incremental check only
    Given a local project whose queue holds a judge task for a.go
    When the check runs over everything, and then over the changed a.go
    Then the full sweep exits 0 and the incremental check exits 1 saying targets await judgment

  @CGPCH-B73 @unit-level
  Scenario: The check says how many judgments it queued, and nothing when it queued none
    Given a judgment gate over one code file in local mode
    When the full sweep runs twice
    Then the first run says "1 target(s) awaiting AI judgment" and the second, which queued nothing new, does not

  @CGPCH-B23 @unit-level
  Scenario: Stale judge tasks leave the queue
    Given judge tasks for a, b (both from the check) and c (from a person), and a run that enqueued only a
    When the stale judgments are dropped
    Then the tasks of a and c stay and the task of b is gone

  @CGPCH-B24 @unit-level
  Scenario: An incremental check keeps the judgments it did not look at
    Given a queued judge task for b and an incremental run that judged only a
    When the judgments are enqueued incrementally, and then on a full sweep
    Then b stays after the incremental run and is dropped after the full sweep

  @CGPCH-B68 @unit-level
  Scenario: A judge task is read back as the gate that queued it even when another gate's name prefixes it
    Given the judgment gates review and review-deep, and a live review-deep task over src/x.go
    When the stale judgments are dropped after a run where review-deep judged src/x.go
    Then the review-deep task stays in the queue

  @CGPCH-B25 @unit-level
  Scenario: In github mode the brief is printed without recording and nothing is queued
    Given a github-mode project with a judgment gate over a.go
    When the check runs over everything without recording, and then recording
    Then the brief names the gate, its guide, its question and a.go, and the local queue stays empty

  @CGPCH-B26 @unit-level
  Scenario: The judgment brief names every target
    Given three targets awaiting two judgment gates
    When the judgment brief is printed
    Then each of the three targets is named

  @CGPCH-B27 @unit-level
  Scenario: The judgment brief carries the question and the guide
    Given a judgment gate with the guide SPEC.md and a question
    When the judgment brief is printed
    Then the question and SPEC.md are in the brief

  @CGPCH-B28 @unit-level
  Scenario: The judgment brief groups the targets by gate
    Given two targets of the same gate and one of another
    When the judgment brief is printed
    Then the shared question appears once and both gates are named

  @CGPCH-B29 @unit-level
  Scenario: The judgment brief says who judges
    Given targets awaiting judgment
    When the judgment brief is printed
    Then it says the pipeline does NOT judge and names REV-CK5

  @CGPCH-B30 @unit-level
  Scenario: The judgment brief lists ten targets per gate and counts the rest
    Given twenty-five targets awaiting the same gate
    When the judgment brief is printed
    Then ten targets are listed and it prints "… and 15 more"

  @CGPCH-B71 @unit-level
  Scenario: With nothing awaiting judgment there is no judgment brief
    Given a project in github mode whose gates judge nothing
    When the full sweep runs without recording
    Then no judgment brief is printed

  @CGPCH-B31 @unit-level
  Scenario: Governed files missing from the map make a stale-map warning
    Given a map holding a.spec.md, older on disk than the file, and then a new b.spec.md
    When the stale-map warning is computed before and after b.spec.md exists
    Then there is no warning before, and a STALE warning that does not say changed after

  @CGPCH-B32 @unit-level
  Scenario: Nodes edited after the map build are warned about on the incremental check
    Given a map whose node a.go carries a revision older than its content
    When the check runs over the changed a.go, and then over everything
    Then the incremental run warns that the map is OLDER and names a.go, and the full sweep does not

  @CGPCH-B33 @unit-level
  Scenario: A map written by another version is warned about
    Given a map written by 0.1.9 and a running binary 0.1.8
    When the version warning is computed
    Then it names 0.1.9 and 0.1.8 and mentions map build

  @CGPCH-B34 @unit-level
  Scenario: A map with no writer version raises no warning
    Given a map with no recorded writer version
    When the version warning is computed
    Then it is empty

  @CGPCH-B35 @unit-level
  Scenario: A declared gate with nothing to measure is named
    Given a gate declared on features in a project with no feature
    When the check runs over everything
    Then the gate never-applies is listed under the gates with nothing to measure

  @CGPCH-B36 @unit-level
  Scenario: The local backlog is printed on the full local sweep only
    Given projects with an open issue in todo, in local and in github mode
    When the check runs over everything and over a changed file
    Then only the local full sweep prints "still open locally"

  @CGPCH-B37 @unit-level
  Scenario: Governance tips appear on the full sweep only
    Given a project with code and none of the canonical protection gates
    When the check runs over everything, and then over a changed file
    Then the full sweep prints the governance tips and the incremental check does not

  @CGPCH-B74 @unit-level
  Scenario: With no governance tip the full sweep prints neither a tip nor the pointer to the doctor
    Given a project of specs only, about which the tips have nothing to say
    When the full sweep runs
    Then neither a governance tip nor the pointer to the doctor is printed

  @CGPCH-B38 @unit-level
  Scenario: The report is mirrored to a file
    Given a project with a map and gates
    When the check runs over everything
    Then it prints "output mirrored to .anchors/"

  @CGPCH-B39 @unit-level
  Scenario: A blocking failure exits 1 with the mirror complete
    Given a blocking gate that fails on a.go
    When the check runs over everything
    Then it exits 1 and the mirror holds "✗ blocked — 1 blocking gate(s) failed"

  @CGPCH-B40 @unit-level
  Scenario: The name column fits the longest name
    Given the names eslint, handler-ddb-inline-passivo and circular, and then only short names
    When the name width is computed
    Then it fits handler-ddb-inline-passivo, and is at least 18 for the short names

  @CGPCH-B41 @unit-level
  Scenario: Each counter column has its own width
    Given one gate with 1116 passes and another with 1 fail and 582 skips
    When the column widths are computed
    Then pass is 4 wide, fail 1 and skip 3

  @CGPCH-B42 @unit-level
  Scenario: The always-present columns are at least one wide
    Given one gate with every counter at zero
    When the column widths are computed
    Then pass, fail, skip and judgment are 1 wide and drift is 0

  @CGPCH-B43 @unit-level
  Scenario: Without drift the drift column does not exist
    Given two gates, neither with drift
    When the profile table is printed
    Then no line holds more than the two-space separator between the fail counter and the skip counter

  @CGPCH-B78 @unit-level
  Scenario: The indeterminate counter is the skipped and pending less the drift
    Given a gate with one skip and two pending results, one of them drift
    When the profile table is printed
    Then its line reads "⚠1  ~2"

  @CGPCH-B44 @unit-level
  Scenario: An empty drift cell is measured in terminal columns
    Given a drift column three digits wide
    When an empty drift cell is built
    Then it is 4 columns wide

  @CGPCH-B45 @unit-level
  Scenario: A clean gate has nothing pending of any kind
    Given gates with passes only, and with a fail, drift, a skip, a pending item or a judgment
    When cleanliness is decided
    Then only the gate with passes only is clean

  @CGPCH-B46 @unit-level
  Scenario: The default table shows the clean gates
    Given a clean gate and a gate with a failure
    When the profile table is printed without only-issues
    Then the clean gate is listed and no omission footer appears

  @CGPCH-B47 @unit-level
  Scenario: Only-issues omits the clean gates and counts them
    Given two clean gates and a gate with a failure
    When the profile table is printed with only-issues
    Then only the failing gate is listed and the footer reads "2 gate(s) with nothing to report"

  @CGPCH-B48 @unit-level
  Scenario: Show-drift lists every drift item
    Given 120 drift items of one gate
    When the profile table is printed with show-drift
    Then the first, middle and last targets are listed and nothing says "… and"

  @CGPCH-B49 @unit-level
  Scenario: Show-drift is not cut by the size of the scan
    Given a scan of more than forty results with one drift item
    When the profile table is printed with show-drift
    Then the drift block and its target AlertSheet.tsx are printed

  @CGPCH-B50 @unit-level
  Scenario: Without show-drift only the counter appears
    Given one drift item on single-target
    When the profile table is printed without show-drift
    Then single-target is not listed and the table shows ⚠1

  @CGPCH-B51 @unit-level
  Scenario: A small scan lists the reason of each skip
    Given one skip with the reason "not a screen" on n.ts
    When the profile table is printed
    Then it prints "~ 1 indeterminate — not a failure" and the reason under "~ g @ n.ts"

  @CGPCH-B52 @unit-level
  Scenario: A large scan does not list the skip reasons
    Given more than forty skips with a reason
    When the profile table is printed
    Then no indeterminate block is printed

  @CGPCH-B53 @unit-level
  Scenario: The legend explains only the symbols used
    Given a table with no drift, and a table with a judgment pending
    When the profile tables are printed
    Then the first legend has no drift line, and the second has the awaiting-judgment line

  @CGPCH-B54 @unit-level
  Scenario: A repeated drift reason is written once with its targets
    Given fifty drift items of mutation-score with the same reason
    When the drift listing is printed
    Then the reason appears once, followed by "50 target(s):" and every target on its own line

  @CGPCH-B55 @unit-level
  Scenario: Distinct drift reasons stay target by target
    Given three drift items of feature-test-match, each with its own reason
    When the drift listing is printed
    Then each reason is printed and no target count heading appears

  @CGPCH-B56 @unit-level
  Scenario: The drift heading counts items and gates
    Given four drift items across three gates
    When the drift listing is printed
    Then it reads "4 pending item(s) in 3 gate(s)"

  @CGPCH-B57 @unit-level
  Scenario: The findings heading counts every kind
    Given one blocking failure, two informative failures and one divergence
    When the profile table is printed
    Then it counts 3 failure(s), 1 blocking, 2 informative and 1 divergence(s), pointing at show-drift

  @CGPCH-B82 @unit-level
  Scenario: Each failure is listed with its detail, and no finding means no findings heading
    Given a blocking failure with the detail "broken on line 3" and an informative failure with none, and then a profile of passes only
    When the profile tables are printed
    Then the first lists both failures with their marks and the detail under the blocking one, and the second has no findings heading

  @CGPCH-B58 @unit-level
  Scenario: The verdict line says what is still open
    Given profiles with an informative failure, a divergence, a skip, only passes, and a blocking failure
    When the profile table is printed
    Then each verdict line names what is open, and the blocked one reads "✗ blocked — 1 blocking gate(s) failed (+1 informative finding(s))"

  @CGPCH-B83 @unit-level
  Scenario: A clean informative gate is named as ready to become blocking
    Given an informative gate that passed three nodes and failed none, and then the same gate declared blocking
    When the maturation reminder is printed for each
    Then the first names the gate as ready to become blocking, and the second says nothing

  @CGPCH-B59 @unit-level
  Scenario: Occurrences of a detail are printed one per line
    Given the detail "line 8: hex colour; line 10: hex colour; line 12: hex colour"
    When it is indented for the report
    Then it takes three lines, each indented

  @CGPCH-B60 @unit-level
  Scenario: Two occurrences already break
    Given the detail "line 8: error; line 10: error"
    When its occurrences are broken
    Then it takes two lines

  @CGPCH-B61 @unit-level
  Scenario: A single occurrence stays whole
    Given a one-sentence detail with no separator
    When its occurrences are broken
    Then it stays one line

  @CGPCH-B62 @unit-level
  Scenario: A glued semicolon does not break
    Given a detail holding "a;b;c"
    When its occurrences are broken
    Then it stays one line

  @CGPCH-B63 @unit-level
  Scenario: A long list of items breaks one per line
    Given a detail listing thirty file paths separated by commas
    When it is indented for the report
    Then no line is much longer than the break threshold and the last path is kept

  @CGPCH-B64 @unit-level
  Scenario: A short sentence with commas stays whole
    Given the sentence "the spec exists, the code exists, and the two reference each other"
    When its occurrences are broken
    Then it stays one line

  @CGPCH-B65 @unit-level
  Scenario: Long prose with commas stays whole
    Given a long sentence whose clauses are separated by commas
    When its occurrences are broken
    Then it stays one line

  @CGPCH-B66 @unit-level
  Scenario: A list of paths still breaks
    Given twenty paths joined by commas
    When its occurrences are broken
    Then it takes more than one line

  @CGPCH-B79 @unit-level
  Scenario: The time table counts each gate's targets and aligns its columns
    Given a gate with targets of every verdict, a gate with five targets and a gate with one
    When the time table is printed
    Then the first counts 15 targets, the five are right-aligned to its width, the names are padded to the longest, and the single target reads "1 target" with no worst time

  @CGPCH-B80 @unit-level
  Scenario: The slowest targets are listed slowest first, and only when a time was recorded
    Given targets of 500ms, 300ms and 50ms whose gate names sort the other way, and then profiles with no target and with no time recorded
    When the time tables are printed
    Then the first lists the targets slowest first, and the others leave the list out

  @CGPCH-B81 @unit-level
  Scenario: Times are rounded to what a decision needs
    Given the durations 1.23456789s, 1.234567ms and 1.234µs
    When they are rounded for the table
    Then they read 1.23s, 1.2ms and 1µs

  @CGPCH-I01 @unit-level
  Scenario: Declaring a perspective does not change the cost axis
    Given a slow gate that skips the change perspective and a fast gate
    When the gates are filtered for the full sweep with and without skip-slow
    Then both run without skip-slow, and skip-slow removes the slow one

  @CGPCH-I02 @unit-level
  Scenario: The version warning compares names by equality
    Given the versions dev and dev, dev and 0.1.9, 0.1.9 and dev
    When the version warning is computed for each pair
    Then only the equal pair raises no warning

  @CGPCH-I03 @unit-level
  Scenario: The skip column does not move with or without drift
    Given a gate with drift and a gate without
    When the profile table is printed
    Then the skip symbol is in the same column on both lines

  @CGPCH-X01 @unit-level
  Scenario: The stale-map warning does not bar the check
    Given a map whose node a.go carries a revision older than its content
    When the check runs over the changed a.go without recording
    Then it warns and returns no error

  @CGPCH-E01 @unit-level
  Scenario: The check without configuration fails
    Given a directory with no anchors.yaml
    When the check runs over everything
    Then it fails naming "load config"

  @CGPCH-E02 @unit-level
  Scenario: A project with no gate has no pipeline
    Given a configuration with no gates section
    When the check runs over everything
    Then it fails saying "no gate declared"

  @CGPCH-E03 @unit-level
  Scenario: The check without a map points at the map build
    Given a project with configuration and no map
    When the check runs over everything
    Then it fails naming "anchors map build"

  @CGPCH-E04 @unit-level
  Scenario: A waiver without a reason is refused
    Given a commit message carrying "[skip-code-flagged: ]"
    When the check runs with that commit message file
    Then it fails with "invalid --skip-rule"

  @CGPCH-E05 @unit-level
  Scenario: The check with no scope is refused
    Given a map and no changed file
    When the nodes are selected without the full-sweep flag
    Then it fails with "provide --changed"

  @CGPCH-E06 @unit-level
  Scenario: A governed file outside the map bars the check
    Given an empty map and a new file under a governed layer
    When the nodes are selected for that file
    Then it fails saying the file is GOVERNED, and not as not-governed

  @CGPCH-E07 @unit-level
  Scenario: A path on neither disk nor map is an error
    Given an empty map and a path that does not exist
    When the nodes are selected for that path
    Then it fails, and not as not-governed

  @CGPCH-B84 @unit-level
  Scenario: The check's stamps do not erase what another process wrote meanwhile
    Given a check that read the map, and an ingestion that wrote a signal to the map on disk after it
    When the check records its stamps
    Then the map on disk has the check's stamp and still has the ingestion's signal

  @CGPCH-B85 @unit-level
  Scenario: The index flag reads what the commit records
    Given a project outside a repository
    When the check runs with --index
    Then it fails saying the git index could not be read

  @CGPCH-B86 @unit-level
  Scenario: --index with no scope judges the staged files
    Given a project with nothing staged, then with one file staged with an old date
    When the check runs with --index and no scope
    Then it first says nothing is staged, then judges the staged file

  @CGPCH-B87 @unit-level
  Scenario: Under --index the date is judged by what the commit records, in every staging state
    Given a file edited and not staged, and later the same file partly staged with an old date, at the repository's top and below it
    When the check runs over the tree and with --index, by --changed and by --all
    Then over the tree the unstaged edit fails, with --index it does not, and the staged change with the old date fails

  @CGPCH-B88 @unit-level
  Scenario: Under --index a file the index does not have is not judged
    Given a map on disk with a node for a file that was never committed and is gone from the tree
    When the check runs over everything, over the tree and with --index
    Then over the tree the gate cannot read the file, and with --index the node is not judged

  @CGPCH-B89 @unit-level
  Scenario: A check with no phase leaves out a gate declared for manual alone
    Given a gate declared for manual alone and a gate with no phase
    When the gates are filtered for a check with no phase, then for --phase manual
    Then the first check keeps only the gate with no phase
    And the manual check keeps both

  @CGPCH-B90 @unit-level
  Scenario: The check says how many targets of its scope are to review
    Given a reviewed gate with two specs to review, one of them in the check's scope
    When the review line is printed
    Then it says one target is to review and how to list them
    And with nothing due it prints nothing

  @CGPCH-B91 @unit-level
  Scenario: A full sweep names the catalog gates missing over the declared layers
    Given a project whose layers two undeclared catalog gates cover
    When the catalog line is printed
    Then it names both and points to the doctor
    And with every applicable gate declared it prints nothing

  @CGPCH-B92 @unit-level
  Scenario: A gate this run leaves out is still declared
    Given a project whose features the catalog covers
    When check --all runs with feature-test-match declared as manual-only, and again with it undeclared
    Then the catalog line counts one gate fewer when it is declared

  @CGPCH-B93 @unit-level
  Scenario: The catalog line leaves out a gate that would only wait for its premise
    Given a catalog gate presupposing derived.test_handle over a declared layer
    When the catalog line is printed for a project that does not declare it
    Then the gate is not named

  @CGPCH-B94 @unit-level
  Scenario: The files check --fix repaired keep their evidence, the line-level signals only when no line moved
    Given a map with two proven files, one repaired on its own line, one with a line inserted
    When the repaired files' evidence is kept
    Then both proofs are at the new content with the declaration, and only the first carries its line-level signals
