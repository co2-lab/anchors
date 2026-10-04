# language: en
# @anchors
#   code: STFTH
#   ref: USSTS
#   updated_at: 2026-10-03
#   layer: feature

@USSTS
Feature: UserSettings — the agent's local settings: the declared role, kept out of git

  @USSTS-B01 @unit-level
  Scenario: The settings live in the local state folder
    Given the project root "/proj"
    When the settings path is computed
    Then it is "/proj/.anchors/settings.yaml"

  @USSTS-B02 @unit-level
  Scenario: A missing settings file is not an error
    Given a project with no settings file
    When the settings are loaded
    Then no error is returned, nothing is decided, and escalated cards are not handled

  @USSTS-B03 @unit-level
  Scenario: The legacy field has three states
    Given settings never asked, declared no, and declared yes
    When each is asked whether it is decided and handles escalated cards
    Then the answers are no/no, yes/no and yes/yes

  @USSTS-B04 @unit-level
  Scenario: The legacy yes grants only the product decision
    Given settings with no role and the legacy field set to yes
    When the capabilities are asked
    Then it decides the product and is decided, but cannot write plans or review

  @USSTS-B05 @unit-level
  Scenario: A declared role wins over the legacy field
    Given settings with the role dev and the legacy field set to yes
    When it is asked whether it decides the product
    Then the answer is no

  @USSTS-B06 @unit-level
  Scenario: What is saved is what is loaded
    Given a project with no state folder
    When settings with the legacy yes, the agent "maquina/sessao" and the date 2026-09-08 are saved and loaded
    Then the loaded settings handle escalated cards, with the same agent and date

  @USSTS-B07 @unit-level
  Scenario: The saved file explains itself
    Given settings saved with the legacy no
    When the file is read
    Then its header says it does not go to git, mentions the gitignore, the escalated cards, product-owner, architect and "anchors settings role"

  @USSTS-B08 @unit-level
  Scenario: Typed answers are read in both languages
    Given the answers "s", "sim", "y", "yes", " Sim ", "n", "nao", "não", "NÃO", "talvez" and ""
    When they are parsed
    Then the first five are yes, the next four are no, and the last two are no answer

  @USSTS-B09 @unit-level
  Scenario: The description names the role or asks for one
    Given settings with the role product-owner and the agent "maq/1", and settings with only the legacy yes
    When they are described
    Then the first reads "Product Owner · maq/1" and the second points to "anchors settings role"

  @USSTS-E01 @unit-level
  Scenario: A settings file that is not YAML fails naming it
    Given a settings file holding "role: [unclosed"
    When the settings are loaded
    Then an error naming the settings file is returned
