# language: en
# @anchors
#   ref: NGSTI
#   updated_at: 2026-10-03
#   layer: feature

@NGSTI
Feature: Ingest — binds the test and log signals the project produced to the nodes of the map

  @NGSTI-B01 @unit-level
  Scenario: The three reports in one pass reach the test, spec and code nodes
    Given a project with a login spec, code and test, and an execution report with one passing, one failing and one skipped case, a line coverage report of three of four lines and a mutation report with two killed and one surviving mutant
    When the three reports are ingested in one pass
    Then the login test node counts one passed, one failed and one skipped
    And the login spec has LOGIN-B01 as its one proven code, and the login code node has three of four lines covered, two mutants killed and one surviving

  @NGSTI-B02 @unit-level
  Scenario: A JUnit report that matches no test node warns
    Given an execution report whose cases name a test file outside the map
    When it is ingested
    Then it warns that no test file matched

  @NGSTI-B03 @unit-level
  Scenario: A project's declared code length governs how JUnit case names are read
    Given a project declaring six-letter codes, with a spec coded LOGINX and a passing case named LOGINX-B01
    When the execution report is ingested
    Then LOGINX-B01 is proven on the spec

  @NGSTI-B04 @unit-level
  Scenario: Manual ingestion warns and proceeds by default
    Given a project declaring a test suite and not refusing manual ingestion
    When an ingestion not run by the test commands starts
    Then it warns of a manual ingestion on the error stream and does not refuse

  @NGSTI-B05 @unit-level
  Scenario: Ingestion run by anchors test never complains
    Given a project that refuses manual ingestion
    When an ingestion run by the test commands starts
    Then it is not refused

  @NGSTI-B06 @unit-level
  Scenario: A project with no declared suite or no config is not asked to use anchors test
    Given a project refusing manual ingestion but declaring no suite, and a directory with no configuration
    When a manual ingestion starts in each
    Then neither is refused

  @NGSTI-B07 @unit-level
  Scenario: The suite key is the report path from the root, or the file name for a report outside the repository
    Given the report /tmp/junit.xml ingested from two checkouts in different places, and a report at apps/mobile/junit.xml inside the repository
    When their suite keys are computed
    Then both external ones are keyed external/junit.xml and the internal one apps/mobile/junit.xml

  @NGSTI-B08 @unit-level
  Scenario: A full run of a suite inside the repository drops the suites ingested from outside it
    Given a map holding the proof of an execution report ingested from outside the repository
    When a full execution report inside the repository is ingested
    Then it says the proof of one external suite was dropped and no external suite remains on the spec

  @NGSTI-B09 @unit-level
  Scenario: A partial run or another external report drops no external suite
    Given a map holding the proof of the external suite external/report.xml
    When a partial run inside the repository and another external report are ingested
    Then the external suite is kept in both cases, while a full in-repository run drops it

  @NGSTI-B10 @unit-level
  Scenario: When two report paths name the same node the path equal to the node id wins
    Given a report naming the node web/src/x.tsx both as src/x.tsx with one value and as web/src/x.tsx with another
    When its paths are resolved fifty times
    Then the node always receives the value of the path web/src/x.tsx

  @NGSTI-B11 @unit-level
  Scenario: An lcov entry for a file edited after the report was written is marked as predating it
    Given a coverage report written after one file and before another, and an entry for a file not on disk
    When the entries are checked against the report
    Then only the file edited after the report is marked as predating it

  @NGSTI-B12 @unit-level
  Scenario: Log occurrences are bound to the spec that declares the failure, stamped with its revision
    Given a log with three occurrences of LOGIN-E01 and one each of LOGIN-E02 and LOGIN-E03, which the login spec declares
    When the logs are ingested
    Then the login spec carries three failures stamped with its revision, LOGIN-E01 with three occurrences

  @NGSTI-B13 @unit-level
  Scenario: Failure codes no spec declares are reported after the log ingestion
    Given a log with an occurrence of GHOST-E07, which no spec declares
    When the logs are ingested
    Then the summary names GHOST-E07 as a code no spec declares

  @NGSTI-I01 @unit-level
  Scenario: Ingesting the same logs again replaces the earlier occurrences instead of adding to them
    Given a log already ingested once, leaving three occurrences of LOGIN-E01
    When the same logs are ingested again
    Then the login spec still carries three failures and LOGIN-E01 still counts three

  @NGSTI-X01 @unit-level
  Scenario: A failure code no spec declares is bound to no spec
    Given a log with an occurrence of GHOST-E07, which no spec declares
    When the logs are ingested
    Then no node of the map carries a GHOST-E07 failure

  @NGSTI-E01 @unit-level
  Scenario: Ingest with no report flag refuses with the usage
    Given a project with a map
    When ingest runs with no report and no log option
    Then it refuses naming the report options

  @NGSTI-E02 @unit-level
  Scenario: Manual ingestion is refused when the project declares manual ingestion blocks
    Given a project declaring a test suite and that manual ingestion must be refused
    When a manual ingestion starts
    Then it is refused with a message naming anchors test

  @NGSTI-E03 @unit-level
  Scenario: Log ingestion without declared log paths is refused
    Given a project that declares no log paths
    When the logs are ingested
    Then it refuses naming logs.paths

  @NGSTI-E04 @unit-level
  Scenario: A missing or malformed report fails naming the report's format
    Given execution, coverage and mutation report paths that do not exist
    When each is ingested
    Then each fails naming the JUnit, lcov or mutation parse

  @NGSTI-E05 @unit-level
  Scenario: Ingesting a report without a map fails and asks for the map build
    Given a directory with no map
    When an execution report is ingested there
    Then it fails with a hint to run the map build

  @NGSTI-B14 @unit-level
  Scenario: A test file's run time is the sum of its cases
    Given a JUnit report whose test file has cases timed 0.5 and 0.75 seconds
    When it is ingested
    Then the test node records 1.25 seconds under the report's suite

  @NGSTI-B17 @unit-level
  Scenario: A spec created after the map build keeps its first proof
    Given a unit whose spec, code and test were created after the last map build
    When a run's report proving its rule is ingested
    Then the spec is in the map with the rule proven

  @NGSTI-B16 @unit-level
  Scenario: A run's proofs are stamped with the tree's revs
    Given a spec edited after the last map build
    When a run's report is ingested, and the map is rebuilt
    Then the spec's proof is fresh, while a manual ingestion of the same report leaves it stale

  @NGSTI-B15 @unit-level
  Scenario: A report's signals are kept under its path from the root
    Given a report inside the repository and one outside it
    When the key of each is read
    Then the first is its path from the root and the second is external/ and its file name

  @NGSTI-B18 @unit-level
  Scenario: A suite's whole coverage run marks the files it left out
    Given a suite declaring its lcov report over the source files, and a file of types the report does not list
    When the report is ingested as a whole run, and again as a partial one in another project
    Then the whole run marks the file of types as omitted and says so, and the partial run marks nothing

  @NGSTI-B19 @unit-level
  Scenario: A JUnit case with no file finds its test by its class's folder and the one file defining its test
    Given a Go project whose src/handlers holds probes_test.go defining TestReady and other_test.go defining TestOther
    And a JUnit report with no file, its classes the import path of src/handlers, with TestReady, a subtest of it, TestNowhere, and TestReady of a missing package
    When the report is ingested
    Then one test file matched, probes_test.go has the two TestReady cases passed, and other_test.go got nothing

  @NGSTI-B20 @unit-level
  Scenario: A targeted file that gave no mutant is measured with nothing to mutate
    Given a targeted mutation run over src/login.ts whose report, in the format, lists no file
    When the report is ingested
    Then src/login.ts is measured at its revision with no mutant

