# language: en
# @anchors
#   code: MNCFT
#   ref: MNCMD
#   layer: feature

@MNCMD
Feature: MonitorCommand — `anchors monitor`, `anchors monitor run` and the agent's hook: watching the project's long processes and recording them

  @MNCMD-B01 @unit-level
  Scenario: The monitor reads once, or until every run ended
    Given a project with a run recorded as ended
    When the monitor reads once, and then until done with nothing running
    Then the first prints the run's end and writes its memory, and the second says nothing runs and exits

  @MNCMD-B02 @unit-level
  Scenario: The timing comes from the defaults, the block and the flags
    Given a project whose monitor block sets the heartbeat
    When the monitor runs with a flag that does not read, and with a block that does not read
    Then each fails naming the value

  @MNCMD-B03 @unit-level
  Scenario: A wrapped command's run is recorded with its output and its exit
    Given a command that prints a line and exits with code 3
    When it runs through anchors monitor run
    Then its output passes through, its record holds its process, its output file and exit 3, and the command ends with code 3

  @MNCMD-B04 @unit-level
  Scenario: Before a command, the hook records the agent's long ones
    Given the agent running a background command, a short one, an anchors test, and a command outside a project
    When the hook is called before each
    Then only the background command is recorded, by its tool call's id

  @MNCMD-B05 @unit-level
  Scenario: After a command, the hook records its output file or its exit
    Given a recorded background command and a recorded foreground one
    When the hook is called after each, with the answers the agent got
    Then the first has its output file, and the second ended with its exit code

  @MNCMD-B06 @unit-level
  Scenario: Anchors' long commands record their own run
    Given anchors test in a project, anchors status in it, and anchors test in a folder with no anchors.yaml
    When each begins and ends with exit 2
    Then only the first is recorded, ended with exit 2

  @MNCMD-I01 @unit-level
  Scenario: The agent's hook never blocks the agent
    Given invalid input, another tool, a command outside a project and a valid call
    When the hook runs on each
    Then each prints nothing and succeeds

  @MNCMD-E01 @unit-level
  Scenario: A timing flag that does not read fails naming it
    Given the flag every set to "soon"
    When the monitor runs
    Then it fails naming "soon"
