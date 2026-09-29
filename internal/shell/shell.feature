# language: en
# @anchors
#   ref: PSXSH
#   updated_at: 2026-09-29
#   layer: feature

@PSXSH
Feature: Shell — the POSIX shell that runs a project's commands

  @PSXSH-B01 @unit-level
  Scenario: The sh on PATH is the shell
    Given a system with sh on PATH
    When the shell is looked up
    Then it is that sh

  @PSXSH-B02 @unit-level
  Scenario: On Windows, the shell beside git
    Given Windows with no sh on PATH and git installed with its shell under bin or usr/bin
    When the shell is looked up
    Then it is git's sh.exe

  @PSXSH-B03 @unit-level
  Scenario: A command is sh -c with its arguments
    Given a script and two arguments
    When the command is built and run
    Then the arguments reach the script as $0 and $1

  @PSXSH-E01 @unit-level
  Scenario: No shell is an environment error
    Given no sh on PATH, and on Windows no shell beside git
    When the shell is looked up
    Then the error says no POSIX shell was found and how to get one
