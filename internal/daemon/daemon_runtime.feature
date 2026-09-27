# language: en
# @anchors
#   ref: DMRND
#   updated_at: 2026-09-26
#   layer: feature

@DMRND
Feature: DaemonRuntime — how each platform probes and terminates the background watcher

  @DMRND-B01 @unit-level
  Scenario: A live process is alive and an exited one is not
    Given the PID of this test process and the PID of a child that has already exited
    When each is probed
    Then the first is alive and the second is not

  @DMRND-B02 @unit-level
  Scenario: Termination on Unix-like systems is a catchable SIGTERM
    Given a live child process on a Unix-like system
    When it is terminated
    Then it ends by the signal SIGTERM, not by a kill

  @DMRND-B03 @unit-level
  Scenario: A live process of another user is alive
    Given a Unix-like system, a test not run as root, and pid 1 answering the probe with EPERM
    When pid 1 is probed
    Then it is reported alive

  @DMRND-I01 @unit-level
  Scenario: Stopping leaves no PID file even when the process cleans nothing
    Given a PID file holding the PID of a child that never cleans the state folder
    When the watcher is stopped
    Then the PID file no longer exists
