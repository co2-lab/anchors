# language: en
# @anchors
#   code: PRFTR
#   ref: PRCRN
#   layer: feature

@PRCRN
Feature: ProcessRuns — the record of the project's long processes, and the operating system's view of them

  @PRCRN-B01 @unit-level
  Scenario: A record is written whole and read back as written
    Given two records and a file that does not read as a record
    When they are saved and listed
    Then both are read back as written, oldest first, and the other file is skipped

  @PRCRN-B02 @unit-level
  Scenario: Finishing a run records how it ended
    Given a running record
    When it is finished with an exit code, and another with none
    Then the first is finished with its exit, and the second ended with none

  @PRCRN-B03 @unit-level
  Scenario: Pruning keeps the runs going and the latest ended ones
    Given a running record, three recent ended ones and one ended long ago
    When the records are pruned to two within a week
    Then the running one and the two latest ended ones are kept

  @PRCRN-B04 @unit-level
  Scenario: The process table is read as the system prints it
    Given CPU times in the forms ps prints, a ps listing and an lsof working-directory listing
    When they are read
    Then each gives its seconds, its processes and their folders

  @PRCRN-B05 @unit-level
  Scenario: A process belongs to the project by its folder or its command line
    Given processes in the root, under it, in a sibling folder sharing its prefix, and with no folder
    When their belonging is asked
    Then the first two belong, the sibling does not, and the last belongs only when its command names the root

  @PRCRN-B06 @unit-level
  Scenario: The process tree gives descendants, ancestors and the tree's CPU time
    Given a parent with a child and a grandchild
    When the tree is asked
    Then the descendants, the ancestors and the summed CPU time are given

  @PRCRN-B07 @unit-level
  Scenario: Listing the system's processes includes the process asking
    Given the running test process
    When the system's processes are listed
    Then its own process is among them

  @PRCRN-B08 @unit-level
  Scenario: The system gives a process's folder, whether it lives, and the load
    Given the running test process
    When its folder, its life and the load are asked, and the life of process 0
    Then its folder is the current one (none on Windows), it lives, process 0 does not, and the load is given where the system has one

  @PRCRN-I01 @unit-level
  Scenario: A reader never sees half a record
    Given a record saved over an existing one
    When it is read back
    Then it is whole and no temporary file is left

  @PRCRN-E01 @unit-level
  Scenario: Finishing a run that has no record fails
    Given no record
    When a run is finished
    Then the error is returned and nothing is written
