# language: en
# @anchors
#   code: RCCRF
#   ref: RCCRS
#   updated_at: 2026-11-10
#   layer: feature

@RCCRS
Feature: RecordingCurrent — a recorded response is of the system the project speaks to today

  @RCCRS-B01 @unit-level
  Scenario: A recording of another version fails
    Given a system declaring the version it speaks, and recordings of that version and of an older one
    When the gate confronts each
    Then the current one passes and the older one fails, saying both versions

  @RCCRS-B02 @unit-level
  Scenario: A recording past its system's age diverges
    Given a system whose recordings live 90 days, and a recording made 120 days ago
    When the gate confronts it, and again with no max age declared
    Then the first diverges saying its date and age, and the second passes

  @RCCRS-B03 @unit-level
  Scenario: A file with no recording is skipped
    Given a test file with no recording mark
    When the gate confronts it
    Then it is skipped

  @RCCRS-E01 @unit-level
  Scenario: A mark that cannot be placed fails
    Given a mark missing its date, and one naming an undeclared system
    When the gate confronts each
    Then each fails, naming the mark or the system
