# language: en
# @anchors
#   ref: SGINA
#   updated_at: 2026-09-27
#   layer: feature

@SGINA
Feature: SignalIngestion — hanging the runner's results on the map's nodes

  @SGINA-B01 @unit-level
  Scenario: A test node sums its layers and records its revisions
    Given a test node A.test.tsx at r1 that composes a util at u1
    When 2 passed are ingested with no layer named and then 1 failed under layer e2e
    Then the node holds unit 2 passed and e2e 1 failed, totals 2 passed and 1 failed, revision r1, the util at u1 in its closure, and the ingestion date

  @SGINA-B02 @unit-level
  Scenario: Re-ingesting a layer replaces only that layer
    Given a test node with unit 2 passed and e2e 1 failed
    When unit is ingested again with 5 passed
    Then the node holds unit 5 passed and e2e still 1 failed, totals 5 passed and 1 failed

  @SGINA-B03 @unit-level
  Scenario: A full run records the proven rules and erases the lost ones
    Given the A spec declares V01, V02 and V03, and the cases naming V01 and V02 passed
    When the run is ingested, and later a run where nothing passed is ingested
    Then the spec proves V01 and V02 after the first, and nothing after the second

  @SGINA-B04 @unit-level
  Scenario: One suite never erases another suite's proof
    Given a backend spec proven by the backend suite
    When the mobile suite is ingested, and then the backend suite again without the proof
    Then the backend spec keeps its proof after the mobile suite, and loses it after its own suite stopped proving it

  @SGINA-B05 @unit-level
  Scenario: A suite that proves nothing leaves the union
    Given a spec whose B01 is proven only by the mobile suite
    When the mobile suite is ingested again with nothing proven
    Then the spec proves nothing and holds no mobile entry

  @SGINA-B06 @unit-level
  Scenario: The union is as fresh as its oldest contributor
    Given the mobile suite proved B01 at rev1, the spec moved to rev2, and the backend suite proved B02 at rev2
    When the node's freshness is read, and again after the mobile suite ran at rev2
    Then it reads stale first and fresh after, and an entry with no recorded revision keeps it stale

  @SGINA-B07 @unit-level
  Scenario: A partial run changes only what it saw
    Given the money and member specs fully proven
    When a partial run that saw only member's B01 passing and B02 failing is ingested
    Then money keeps both proofs and member keeps only B01, and a later partial run seeing only money's B01 leaves money's B02 proven

  @SGINA-B08 @unit-level
  Scenario: External suites are dropped and the union recomputed
    Given a spec proven by an external report at rev1 and by the repository suite at rev2, and a code file covered only by an external report
    When the external suites are dropped
    Then external/lcov.info and external/report.xml are named, the spec is fresh with only the repository proof, and the code file has no coverage left

  @SGINA-B09 @unit-level
  Scenario: Coverage with no suite records the percentage and the baseline
    Given the A code file with no signal
    When 3 of 5 lines are ingested, and then 4 of 5
    Then it reads 60 percent with no baseline after the first, and 80 percent with a baseline of 60 after the second

  @SGINA-B10 @unit-level
  Scenario: Suite coverage is the union of lines
    Given the unit suite covering lines 1 and 2 of 4, and the integration suite covering lines 3 and 4
    When both are ingested
    Then the file reads 4 of 4, 100 percent; two suites covering the same lines count them once; and a suite ingested again replaces its lines

  @SGINA-B11 @unit-level
  Scenario: Totals-only suites fall back to the best suite
    Given the unit suite reporting 3 of 10 and the integration suite 5 of 10, with no line detail
    When both are ingested
    Then the file reads 5 of 10

  @SGINA-B12 @unit-level
  Scenario: A suite measured before the edit stays out of the union
    Given the integration suite measured 76 lines of the old text, and the file was edited
    When the unit suite covering 73 of 75 lines of the new text is ingested
    Then the file reads 73 of 75 and the node reads stale

  @SGINA-B13 @unit-level
  Scenario: A report older than the file is kept without a revision
    Given an integration report written before the file's current text
    When it is ingested along with a fresh unit report covering 3 of 4 lines
    Then the integration entry has no revision, the file reads 3 of 4, and the node reads stale

  @SGINA-B14 @unit-level
  Scenario: Line ranges round-trip
    Given the lines 12, 1, 2, 3, 9, 13, 5, 4 and 20
    When they are written as ranges and read back
    Then the ranges read 1-5,9,12-13,20 and give back nine lines

  @SGINA-B15 @unit-level
  Scenario: Mutation totals and the per-scope measurement
    Given the A code file at r1
    When a mutation report of 8 killed and 2 survived with thresholds 60 and 80 is ingested under scope isolated
    Then the file holds 8 killed, 2 survived, the thresholds 60 and 80, and the isolated scope holds the same numbers at r1

  @SGINA-B16 @unit-level
  Scenario: A signal goes stale when its file moves
    Given a code file whose coverage was ingested at r1
    When the file moves to r2
    Then its signal is stale, while a signal with no recorded revision never is

  @SGINA-B17 @unit-level
  Scenario: Paths match at a path boundary in either direction
    Given the node apps/mobile/src/A.tsx and the node a/b.go
    When the report paths /abs/apps/mobile/src/A.tsx, src/A.tsx and xa/b.go are compared
    Then the first two match apps/mobile/src/A.tsx and xa/b.go does not match a/b.go

  @SGINA-B18 @unit-level
  Scenario: Report paths resolve by the report's folder, or stay unowned
    Given Label.tsx in both the landing and the mobile workspace
    When src/atoms/Label.tsx is resolved for a report under the landing folder and for a report at the root
    Then it goes to the landing's Label.tsx first and is ambiguous, with no owner, second

  @SGINA-B19 @unit-level
  Scenario: A report with no instrumented line leaves no percentage behind
    Given the A code file ingested at 4 of 5 lines, with and without a suite
    When a report of 0 of 0 lines is ingested for it
    Then it reads 0 lines at 0 percent

  @SGINA-B20 @unit-level
  Scenario: Among several report paths matching a node, the exact one, then the closest in length, is always chosen
    Given reports where src/A.test.tsx and src/A.tsx are each matched by several paths
    When the execution, coverage and mutation are ingested fifty times over
    Then the test node always takes the exact path's counts, and the code node the closest path's coverage and, on a length tie, the first path's mutants

  @SGINA-I01 @unit-level
  Scenario: The proven rules are the sorted union of the suites
    Given suite one proving B02 and B01 and suite two proving B03 and B01
    When both are ingested
    Then the spec's proven rules are B01, B02 and B03

  @SGINA-X01 @unit-level
  Scenario: Each measurement lands only on its kind of node
    Given a spec, a code file and a test file all matched by the same report path
    When execution, coverage and mutation are ingested for that path
    Then only the test file holds the execution and only the code file holds coverage and mutation

  @SGINA-B21 @unit-level
  Scenario: A mutation ingestion records the timed-out apart
    Given a mutation report where some killed mutants timed out
    When it is ingested
    Then the node records the killed and, apart, how many of them timed out
