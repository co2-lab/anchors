# language: en
# @anchors
#   ref: MGSTM
#   updated_at: 2026-09-28
#   layer: feature

@MGSTM
Feature: MigrationSteps — the registered steps, one per format, take any project from format 1 to the current format

  @MGSTM-B01 @unit-level
  Scenario: Format 2 renames the Portuguese keys in their own files
    Given the registered step producing format 2
    When its key renames are read
    Then the map renames "gerado_por", "code_declarado" and "julgamentos" to "generated_by", "code_declared" and "judgments"
    And the configuration renames "trinca_opcional" to "triad_optional"

  @MGSTM-B02 @unit-level
  Scenario: Format 2 renames the Portuguese gate names wherever they are stored
    Given a format 1 map with a judgment "gate: regra-cumprida"
    And a format 1 configuration with "name: regra-cumprida", "id: trinca-completa" and "check: sem-duplicacao"
    When both are migrated to format 2
    Then the map has "gate: rule-fulfilled"
    And the configuration has "name: rule-fulfilled", "id: triad-complete" and "check: no-duplication"

  @MGSTM-B03 @unit-level
  Scenario: Format 3 renames the eight remaining gate names in the configuration
    Given a format 2 configuration with "name: regra-implementada" and "check: header-conforme"
    When it is migrated to format 3
    Then it has "name: rule-implemented" and "check: header-valid"

  @MGSTM-B05 @unit-level
  Scenario: Format 5 renames code letters, not keys
    Given the registered step that produces format 5, and a configuration on format 4
    When the step is read and the configuration is migrated to 5
    Then the step renames plan F to W, flow P to T and R to O, action R to O, no key, and the file only gets version 5

  @MGSTM-B04 @unit-level
  Scenario: Format 4 renames the four keys that lied about what they hold
    Given a format 3 configuration with "auto_judgment", "triad_optional", "requires_code" and "rule_marking"
    When it is migrated to format 4
    Then it has "enable_auto_judgment", "optional_triad_edges", "sections_require_code" and "rule_marking_policy" and none of the old keys

  @MGSTM-I01 @unit-level
  Scenario: The chain from format 1 to the current format has no hole
    Given the registered steps
    When the steps from format 1 to the current format are asked for
    Then there is one step per format from 2 to the current one, each with its reason

  @MGSTM-I02 @unit-level
  Scenario: A key renamed by two formats ends under its latest name
    Given a format 1 configuration with "trinca_opcional: [tested-by]"
    When it is migrated to the current format
    Then it has "optional_triad_edges: [tested-by]" and neither "trinca_opcional" nor "triad_optional"

  @MGSTM-X01 @unit-level
  Scenario: Format 3 leaves the map's gate and the configuration's id alone
    Given a format 2 map with "gate: regra-implementada" and a format 2 configuration with "id: regra-implementada"
    When both are migrated to format 3
    Then both still say "regra-implementada"

  @MGSTM-X02 @unit-level
  Scenario: Format 4 renames nothing in the map
    Given a format 3 map with the keys "auto_judgment" and "rule_marking"
    When it is migrated to format 4
    Then both keys are unchanged and only the version line changed
