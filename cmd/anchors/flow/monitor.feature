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
    Given the agent running a background script, a background jest, a foreground jest, a short command, an anchors test, and a command outside a project
    When the hook is called before each
    Then only the two background commands are recorded, by their tool calls' ids, the jest named jest

  @MNCMD-B05 @unit-level
  Scenario: After a command, the hook records its output file
    Given a recorded background command and a foreground one
    When the hook is called after each, with the answers the agent got
    Then the first has its output file, and the second left no record

  @MNCMD-B06 @unit-level
  Scenario: Anchors' long commands record their own run
    Given anchors test in a project, anchors status in it, and anchors test in a folder with no anchors.yaml
    When each begins and ends with exit 2
    Then only the first is recorded, ended with exit 2

  @MNCMD-B07 @unit-level
  Scenario: The monitor says the progress summary when it changes, when asked
    Given a project with a demand whose stage has an open task
    When the monitor reads with no flag, with --summary twice, with summary-every past, after a task is done, after a hand edit twice, and with a block turning it on, silenced by --summary=false
    Then nothing; the summary, then only the heartbeat; the summary again; the summary with what it was; one warning; nothing, and with the block alone the summary

  @MNCMD-B08 @unit-level
  Scenario: At the action level the monitor says only what calls for a reaction
    Given Anchors' check passed, a wrapped jest passed and a go test failed, with nothing running
    When the monitor reads at the action level, twice, and with an unknown level
    Then the jest and the failure go out and the check does not, idle goes out once, and the unknown level is refused, named

  @MNCMD-B09 @unit-level
  Scenario: A task in progress with no activity from its holder is said
    Given a task started 20 minutes ago by a, and no run of a
    When the monitor reads for the first time, then having watched for 30 minutes, again, after a command of a, with --idle off, and with an idle that does not read
    Then the first tick says nothing; then the task is said quiet once, then nothing, then moving again; off says nothing, and the bad value is refused, named

  @MNCMD-B11 @unit-level
  Scenario: A holder's activity in another worktree of the repository counts
    Given a repository with a worktree, a task held by a for an hour, and a run of a recorded in the worktree
    When the monitor of the main tree reads, having watched for an hour
    Then the task is not said quiet

  @MNCMD-B10 @unit-level
  Scenario: The agent's hook records a subagent as a run
    Given the payloads a real session sent: a subagent starting, a background command of its own, and the subagent stopping
    When the hook receives each, and a subagent start outside a project
    Then the subagent is a run of kind agent with its transcript, the command carries its id, the stop closes it with its last message and its task's name, and outside a project nothing is recorded

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
