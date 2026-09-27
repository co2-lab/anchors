# language: en
# @anchors
#   ref: INQSN
#   updated_at: 2026-09-26
#   layer: feature

@INQSN
Feature: InitQuestions — describe the human decisions of `init` so an agent can answer them without the terminal UI, and judge every answer

  @INQSN-B01 @unit-level
  Scenario: The questions come in the order of the terminal UI
    Given an inference proposal with no configuration
    When the questions are built
    Then their identifiers are preset, header, artifacts, gates, colocation, layers, workflow, repo, labels, governs in that order

  @INQSN-B02 @unit-level
  Scenario: Every question carries what the agent needs to decide
    Given an inference proposal with no configuration
    When the questions are built
    Then every question has an identifier, a text, a type and a why
    And every single-choice question has options

  @INQSN-B03 @unit-level
  Scenario: The defaults come from the inference
    Given a colocated proposal with the code layers "api-code" and "web-code" and a detected spec
    When the questions are built
    Then the colocation default is yes
    And the layers question offers and defaults to "api-code" and "web-code"
    And the artifacts default includes "spec"
    And with no proposal the colocation default is no and the layers are empty

  @INQSN-B04 @unit-level
  Scenario: The preset and the work-queue mode have their choices and defaults
    Given the preset catalog "go" and "nextjs"
    When the questions are built
    Then the preset offers "nenhum", "go", "nextjs" with "nenhum" by default
    And the workflow offers "local", "manual", "github" with "local" by default

  @INQSN-B05 @unit-level
  Scenario: Every question gets a verdict and unanswered ones take the default
    Given the questions and an answer only for the header
    When the answers are validated
    Then there is one verdict per question in the questions' order
    And the header is not marked as default while the gates is

  @INQSN-B06 @unit-level
  Scenario: An answer outside the options is refused with the accepted values
    Given the answer "preset-que-nao-existe" for the preset and "code" for the artifacts
    When the answers are validated
    Then the preset and the artifacts verdicts are refused
    And each detail lists the accepted values after "aceitos:"
    And a free label "whatever" for the labels question is accepted

  @INQSN-B07 @unit-level
  Scenario: The github mode requires repository and labels
    Given the answer "github" for the workflow and no repository nor labels
    When the answers are validated
    Then the repo and the labels verdicts are refused
    And with the repository "acme/exemplo" and the label "anchors" the whole set is accepted

  @INQSN-B08 @unit-level
  Scenario: A repository outside the github mode is refused
    Given the answer "local" for the workflow and the repository "owner/nome"
    When the answers are validated
    Then the repo verdict is refused

  @INQSN-B09 @unit-level
  Scenario: One refused answer refuses the whole set
    Given a set whose only invalid answer is the preset
    When the set is judged
    Then the whole set is refused

  @INQSN-B10 @unit-level
  Scenario: An answer given empty is not the default
    Given the artifacts answered with an empty list
    When the answers are validated
    Then the artifacts verdict holds the empty list and is not marked as default
    And with no artifacts answer the verdict is marked as default

  @INQSN-I01 @unit-level
  Scenario: No answer goes missing from the verdict
    Given the questions and a single explicit answer
    When the answers are validated
    Then every question has its verdict entry
    And only the unanswered questions are marked as default

  @INQSN-X01 @unit-level
  Scenario: An invalid answer is not corrected
    Given the answer "preset-que-nao-existe" for the preset
    When the answers are validated
    Then the preset verdict keeps "preset-que-nao-existe" as its value and is refused

  @INQSN-B11 @unit-level
  Scenario: The artifacts question offers the artifact options, code included
    Given the questions of a project
    When the artifacts are answered with spec and code
    Then the answer is accepted
    And the artifacts options are the artifact names init offers
    And a project where inference found code, specs and tests has code, spec and test pre-checked, in that order, on every call

  @INQSN-B12 @unit-level
  Scenario: The question texts, their reasons and the refusal details are in the project's language
    Given the project language set to English, then Portuguese, then Spanish
    When the questions are asked and an unknown preset is answered
    Then the preset question, its reason and the refusal detail are written in that language
    And no question shows a bare catalog key
