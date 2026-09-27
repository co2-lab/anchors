# language: en
# @anchors
#   ref: WTDMW
#   updated_at: 2026-09-27
#   layer: feature

@WTDMW
Feature: WatchDaemon — the watcher started in the background survives the terminal that started it

  @WTDMW-B01 @unit-level
  Scenario: On unix a detached child leads its own process group
    Given a child process prepared on a unix-like system
    When it is detached and then started
    Then its process group id is its own pid
    And its process group is not the parent's

  @WTDMW-B02 @unit-level
  Scenario: On Windows a detached child is created in a new process group
    Given a child process prepared on Windows
    When it is detached
    Then it is set to be created with the new-process-group flag

  @WTDMW-I01 @unit-level
  Scenario: Exactly one detachment implementation builds per platform
    Given the flow package's source files
    When the files built for a unix-like target and for a Windows target are listed
    Then the unix-like list holds the unix implementation and not the Windows one
    And the Windows list holds the Windows implementation and not the unix one

  @WTDMW-X01 @unit-level
  Scenario: Detaching does not start the child
    Given a child process prepared and not started
    When it is detached
    Then it still has no running process
