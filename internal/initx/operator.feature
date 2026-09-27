# language: en
# @anchors
#   ref: OPDTP
#   updated_at: 2026-09-26
#   layer: feature

@OPDTP
Feature: OperatorDetection — tell whether a person or an AI is running `init`, and whether the discovery phase is still to be done

  @OPDTP-B01 @unit-level
  Scenario: A known agent variable makes the operator an AI even with a terminal
    Given an environment where "CLAUDE_CODE_ENTRYPOINT" is set
    And a terminal is attached
    When the operator is detected
    Then the operator is an AI

  @OPDTP-B02 @unit-level
  Scenario: The generic AI_AGENT variable makes the operator an AI even with a terminal
    Given an environment where only "AI_AGENT" is set
    And a terminal is attached
    When the operator is detected
    Then the operator is an AI
    And no tool name is given

  @OPDTP-B03 @unit-level
  Scenario: No terminal and no agent variable is an AI
    Given an empty environment
    And no terminal is attached
    When the operator is detected
    Then the operator is an AI

  @OPDTP-B04 @unit-level
  Scenario: A terminal and no agent variable is a person
    Given an empty environment
    And a terminal is attached
    When the operator is detected
    Then the operator is a person

  @OPDTP-B05 @unit-level
  Scenario: The tool is named after the known variable that is set
    Given an environment where "CURSOR_TRACE_ID" is set
    When the tool name is asked
    Then the name is "Cursor"
    And with an empty environment the name is empty

  @OPDTP-B06 @unit-level
  Scenario: Several known variables resolve by the first variable name in order
    Given an environment where both "CLAUDECODE" and "AIDER_MODEL" are set
    When the tool name is asked
    Then the name is "Aider", because "AIDER_MODEL" sorts first

  @OPDTP-B07 @unit-level
  Scenario: The discovery phase is due only when nothing is found and nothing is described
    Given a project folder with no project description
    When the phase is checked with an empty proposal, with no proposal, and with proposals that found a code folder, a spec, a feature or a test
    Then the phase is due for the empty and the missing proposal
    And it is not due for any proposal that found something
    And once "PROJECT.md" is written it is not due even with an empty proposal

  @OPDTP-B08 @unit-level
  Scenario: The project description is recognised under its three spellings
    Given a project folder holding only "project.md", or only "Project.md"
    When the project description is looked up
    Then it is found in both folders

  @OPDTP-B09 @unit-level
  Scenario: Only the tools with a stable command line get a command
    Given environments naming Claude Code, Gemini CLI, Aider, Cursor, and no tool
    When the command to open the AI is built
    Then Claude Code opens with "claude", Gemini CLI with "gemini", Aider with "aider --message"
    And Cursor and no tool get no command

  @OPDTP-B10 @unit-level
  Scenario: The discovery prompt points at the guide and back at init
    Given the discovery prompt
    When its text is read
    Then it tells to run "anchors guide project"
    And it ends by pointing at "anchors init"

  @OPDTP-I01 @unit-level
  Scenario: The whole prompt is one argument of the command
    Given each tool that has a command
    When the command to open it is built
    Then its last argument equals the discovery prompt exactly, with no shell quoting added

  @OPDTP-X01 @unit-level
  Scenario: No command is offered for a tool without a stable command line
    Given an environment naming Cursor or Codex
    When the command to open the AI is built
    Then there is no command
