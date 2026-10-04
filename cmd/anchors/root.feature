# language: en
# @anchors
#   code: RTFTR
#   ref: CLRTC
#   updated_at: 2026-10-03
#   layer: feature

@CLRTC
Feature: CliRoot — every command passes through one root that speaks the project's language and honours the freeze

  @CLRTC-B01 @unit-level
  Scenario: A frozen project refuses the command, naming it and the reason
    Given a project frozen for "rotate the key"
    When generated-paths runs there
    Then it is refused saying `generated-paths` will not run and "rotate the key"

  @CLRTC-B02 @unit-level
  Scenario: The commands that only read, and thaw, run while frozen
    Given a frozen project as the working directory
    When thaw, freeze, status, doctor, guide, help, completion, coverage and impact are checked
    Then none is refused

  @CLRTC-B03 @unit-level
  Scenario: A subcommand of an allowed command runs while frozen
    Given a frozen project as the working directory
    When guide work is checked
    Then it is not refused

  @CLRTC-B04 @unit-level
  Scenario: A missing or broken config is not frozen
    Given a project with no config, one with broken YAML, and one with an unknown key
    When generated-paths is checked in each
    Then none is refused

  @CLRTC-B05 @unit-level
  Scenario: The root flag decides which project is read
    Given a frozen project and a project that is not frozen
    When generated-paths is checked with --root pointing at each
    Then the first is refused and the second is not

  @CLRTC-B06 @unit-level
  Scenario: The project's top-level lang is applied before the command runs
    Given configs with lang pt-BR, with lang es and broken YAML, with a nested lang, with lang xx, and none
    When the language is applied
    Then it becomes pt-BR, es, and stays en for the last three

  @CLRTC-B07 @unit-level
  Scenario: The root prints neither the error nor the usage
    Given a frozen project
    When the real root executes generated-paths with its output captured
    Then the command fails
    And nothing was printed by the root

  @CLRTC-I01 @unit-level
  Scenario: Every command the work guide and the pipelines teach is registered
    Given the registered commands
    When the work guide and every pipeline workflow are read
    Then every `anchors <command>` they cite is registered

  @CLRTC-X01 @unit-level
  Scenario: The freeze refusal is written in the project language
    Given a frozen project with lang pt-BR, and a process in English
    When the real root executes generated-paths
    Then the refusal reads "o projeto está CONGELADO"

  @CLRTC-B08 @unit-level
  Scenario: The language is read from a CRLF file
    Given a configuration whose lang line ends in CRLF
    When the language is read
    Then it is the one declared

  @CLRTC-I02 @unit-level
  Scenario: Every command that takes several files reads them the same way
    Given every registered command
    When the ones whose usage takes several files are confronted with the mark of the shared reading
    Then each carries it
