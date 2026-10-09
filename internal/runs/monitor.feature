# language: en
# @anchors
#   code: MNFTR
#   ref: MNTRS
#   layer: feature

@MNTRS
Feature: RunMonitor — reading the project's runs tick by tick: what started, stalled, died, ended, and what runs

  @MNTRS-B01 @unit-level
  Scenario: A runner under the project that no record accounts for is a run found running
    Given an npx jest with a node worker under the project, a jest of another folder, the monitor itself, an anchors run wrapper, and a shell running a recorded anchors check
    When a tick reads them
    Then one run starts — the npx jest —, recorded by the process table, and the others are no run

  @MNTRS-B02 @unit-level
  Scenario: A quiet run is stalled once, and moving again is said
    Given a run whose CPU, output and processes do not change past its quiet time, a child of it ending meanwhile
    When ticks read it, and then a new process appears under it
    Then one stalled line is given, and then a line saying it moves again

  @MNTRS-B03 @unit-level
  Scenario: A run that ended is said once, by its exit, its summary, or nothing
    Given runs with an exit record, with a summary in their output, with an output cut short, with a launcher that never recorded the end, and seen only by the process table
    When their processes are gone
    Then each ends by its exit, its summary, as died with its last lines, as died, and as ended with its exit unknown
    And a command of the agent the monitor never saw running is said to have ended before the monitor saw it

  @MNTRS-B04 @unit-level
  Scenario: A command of the agent takes the process that runs its program
    Given a record of the agent for "cd app && npx jest --ci" and the jest process under the project
    When a tick reads them
    Then the record takes the process — or, for a maestro flow, the maestro process its runner recognizes —, and a command no process runs ends once its output stops growing

  @MNTRS-B05 @unit-level
  Scenario: A report that landed is said with its failures
    Given a report there before the monitor's memory began, and one modified after
    When a tick reads them
    Then only the second is a line, with its tests, failures and first failure

  @MNTRS-B06 @unit-level
  Scenario: The heartbeat says what runs, what ended, and the load
    Given one running run, one stalled, one ended since the last heartbeat, and a load above the CPUs
    When the heartbeat is due, and then with nothing running
    Then the line names them, counts the ended one and marks the load HIGH, and then says NOTHING

  @MNTRS-B07 @unit-level
  Scenario: A re-armed monitor says what happened meanwhile, once
    Given a monitor that saw a run, wrote its memory and stopped
    When the run dies while nobody watches, and a monitor with the memory reloaded reads twice
    Then the death is its first line, and the second reading says nothing of it

  @MNTRS-B08 @unit-level
  Scenario: Progress is sparse
    Given a run whose output has progress marks
    When ticks read it within the progress interval and past it
    Then one progress line is given past the interval only

  @MNTRS-B09 @unit-level
  Scenario: The settings are the defaults, then the project's block
    Given no block, then a block with timing, a runner replacing jest, a report and bounds, then values that do not read
    When the settings are read
    Then the defaults, then the project's values first, then an error naming each bad key

  @MNTRS-B10 @unit-level
  Scenario: An output's end gives its verdict, its summary and its last lines
    Given outputs ending with a failure beside a pass, with a pass, and with neither
    When their verdicts are read
    Then failed, passed and nothing, with the summary line and the last lines

  @MNTRS-B11 @unit-level
  Scenario: Reports are read by their globs
    Given a JUnit report with a failure
    When the reports are read, and then read again knowing its time
    Then it is parsed the first time and given with its time alone the second

  @MNTRS-B12 @unit-level
  Scenario: The built-in runners recognize the common runners
    Given command lines of the common runners and one of none
    When their runner is asked
    Then each is recognized with its kind, and the last has none

  @MNTRS-B13 @unit-level
  Scenario: A command is recognized and shown by its first line
    Given a heredoc file edit that mentions jest in its body, and a long first line
    When its runner and its brief are read
    Then it is no jest run, and it shows its first line cut, with the lines it left out

  @MNTRS-I01 @unit-level
  Scenario: Every run that ended is said once
    Given runs ending in each way across ticks and a reloaded memory
    When the end lines are counted
    Then each run has exactly one

  @MNTRS-E01 @unit-level
  Scenario: A monitor value that does not read is an error naming its key
    Given a monitor block whose heartbeat does not read
    When the settings are read
    Then the error names monitor.heartbeat
