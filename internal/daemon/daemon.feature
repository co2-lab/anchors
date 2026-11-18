# language: en
# @anchors
#   code: DMFTD
#   ref: DMSTD
#   updated_at: 2026-10-03
#   layer: feature

@DMSTD
Feature: DaemonState — the background watcher's state files: PID, pause flag, log and meta

  @DMSTD-B01 @unit-level
  Scenario: The state files live in the project's state folder
    Given the project root "/proj"
    When the watcher's paths are computed
    Then they are "/proj/.anchors/watch.pid", "watch.log", "watch.paused" and "watch.meta"

  @DMSTD-B02 @unit-level
  Scenario: Running answers the live PID and 0 otherwise
    Given a project with no PID file, then the PID of this live process, then the text "garbage"
    When the running watcher is asked for each
    Then the answers are 0, this process's PID, and 0

  @DMSTD-B03 @unit-level
  Scenario: A PID file of an exited process is removed
    Given a PID file holding the PID of a process that has exited
    When the running watcher is asked
    Then the answer is 0 and the PID file no longer exists

  @DMSTD-B04 @unit-level
  Scenario: Stopping with no watcher is refused
    Given a project with no PID file
    When the watcher is stopped
    Then the error says "not running"

  @DMSTD-B05 @unit-level
  Scenario: Stopping a running watcher terminates it and removes the PID file
    Given a PID file holding the PID of a live child process
    When the watcher is stopped
    Then the child ends by a signal and the PID file no longer exists

  @DMSTD-B06 @unit-level
  Scenario: Pause is the existence of the flag file
    Given a fresh project
    When the watcher is paused and then resumed
    Then it reads not paused, then paused, then not paused

  @DMSTD-B07 @unit-level
  Scenario: Cleanup removes the PID file and the pause flag
    Given a paused project with a PID file
    When the state is cleaned up
    Then neither the PID file nor the pause flag exists

  @DMSTD-B08 @unit-level
  Scenario: The meta records the start moment and the root
    Given a project with no meta file
    When the meta is read, then written for 2026-09-26 10:30 UTC over "/proj" and read again
    Then the first read is empty and the second is "started=2026-09-26T10:30:00Z" and "root=/proj" on two lines

  @DMSTD-I01 @unit-level
  Scenario: The PID file of an exited process does not outlive the check
    Given a PID file holding the PID of a process that has exited
    When the running watcher is asked
    Then the answer is 0 and the PID file is gone

  @DMSTD-E02 @unit-level
  Scenario: A state folder that cannot be created fails the PID write
    Given a project root where ".anchors" is a file, not a folder
    When the PID is written
    Then an error is returned
