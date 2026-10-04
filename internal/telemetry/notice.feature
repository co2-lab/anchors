# language: en
# @anchors
#   code: NTFTA
#   ref: TLNTT
#   updated_at: 2026-10-03
#   layer: feature

@TLNTT
Feature: TelemetryNotice — the telemetry notice reaches whoever did not ask for it, once, and says how to turn it off

  @TLNTT-B01 @unit-level
  Scenario: The notice is shown once per project root
    Given a project root that never showed the notice
    When the notice is called twice for that root
    Then the first call writes the notice
    And the second call writes nothing

  @TLNTT-B02 @unit-level
  Scenario: The text says what is sent, what is not, and how to turn it off
    Given a project root that never showed the notice
    When the notice is written
    Then it mentions decision events, says what it does not send including file content
    And it names both "ANCHORS_TELEMETRY=off" and "telemetry: off"
    And it says it appears "once per project on this machine"

  @TLNTT-B03 @unit-level
  Scenario: The marker is written under the project's unversioned anchors directory
    Given a project root that never showed the notice
    When the notice is written
    Then the file ".anchors/telemetry-noticed" exists under that root

  @TLNTT-B04 @unit-level
  Scenario: The notice is written in the project's language
    Given the languages "en", "pt-BR" and "es"
    When the notice is shown in each
    Then each says, in its language, that it appears once per project on this machine, and each shows "ANCHORS_TELEMETRY=off"

  @TLNTT-I01 @unit-level
  Scenario: After the notice is written it counts as already shown
    Given a project root that never showed the notice
    When the notice is written
    Then asking whether the notice was already shown for that root answers yes

  @TLNTT-X01 @unit-level
  Scenario: The notice does not consult the opt-out itself
    Given the environment variable set to "off"
    When the notice is called for a fresh root
    Then the notice is still written, because deciding whether to show it is the caller's duty

  @TLNTT-E01 @unit-level
  Scenario: An unwritable marker never blocks and the notice shows again
    Given a project root where the ".anchors" directory cannot be created
    When the notice is called twice
    Then both calls write the notice and neither fails
