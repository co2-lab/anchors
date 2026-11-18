# language: en
# @anchors
#   code: SCFTS
#   ref: DCSCD
#   updated_at: 2026-10-03
#   layer: feature

@DCSCD
Feature: DocScaffolds — the starting templates `anchors docs init` proposes, one page per question and per layer

  @DCSCD-B01 @unit-level
  Scenario: The three fixed templates each have a name, a body and what they answer
    Given the proposed scaffolds
    When each is inspected
    Then they are the architecture, behaviour and rules templates
    And each has a template name, a non-empty body and a non-empty sentence of what it answers

  @DCSCD-B02 @unit-level
  Scenario: A layer's page template is named after the layer
    Given the layer "infra"
    When its page scaffold is built in English
    Then its name is "layers/infra.md.tmpl"

  @DCSCD-B03 @unit-level
  Scenario: Init writes the fixed templates and one page per layer, each opening with its purpose
    Given a project whose specs are in the layer "infra"
    When init runs
    Then it writes the three fixed templates and "camadas/infra.md.tmpl"
    And each file starts with a template comment carrying what the page answers

  @DCSCD-B04 @unit-level
  Scenario: An edited template is kept unless forced
    Given a template the team edited after a first init
    When init runs again without force
    Then the edited template is unchanged and it is reported as skipped
    And when init runs with force the template is rewritten

  @DCSCD-B05 @unit-level
  Scenario: The small layer page carries each scenario under its own heading, and the index links it
    Given a small layer whose spec has a feature
    When the skeleton is built
    Then the layer page carries the scenario's steps under "#### GLCGL-B01 — item sem artefato não passa"
    And the behaviour index links to that heading's anchor on the layer page

  @DCSCD-B06 @unit-level
  Scenario: The big layer page summarizes each unit without scenario steps
    Given a layer that the layout considers big
    When the skeleton is built
    Then the layer page carries the layout's summary sentence, the unit's overview and its rule list
    And it does not carry the scenarios' steps

  @DCSCD-B07 @unit-level
  Scenario: The architecture page follows the C4 model
    Given containers app and api talking over HTTPS/JSON and SQL/TLS, and an external database
    When the skeleton is built
    Then the architecture page names both protocols and the database's description
    And it has a level 3 section for app and api and none for the database

  @DCSCD-B08 @unit-level
  Scenario: No container declared is said on the architecture page
    Given a project that declares no container
    When the skeleton is built
    Then the architecture page says that no container is declared

  @DCSCD-I01 @unit-level
  Scenario: The skeleton init writes compiles
    Given a project with a spec and a feature
    When init runs and the documentation is built
    Then the build succeeds with one page per template written

  @DCSCD-X01 @unit-level
  Scenario: Init writes no compiled page
    Given a project with specs
    When init runs
    Then the compiled documentation folder does not exist

  @DCSCD-E01 @unit-level
  Scenario: A templates folder that cannot be created stops init with the error
    Given a project root where a file already occupies the templates folder's name
    When init runs
    Then it returns an error and reports no template written

  @DCSCD-B09 @unit-level
  Scenario: The templates are named and written in the project's language
    Given the languages en, pt-BR and es
    When the scaffolds are built in each
    Then they are named architecture, behavior and rules; arquitetura, comportamento and regras; arquitectura, comportamiento and reglas
    And the architecture page opens with the language's heading, and the layer page lives in layers/, camadas/ or capas/ and asks for the language's overview title
    And no template keeps an unfilled text, and a language with no table gets English

  @DCSCD-B05 @unit-level
  Scenario: A project with API units gets the OpenAPI template
    Given a project with a spec that has an Endpoint section, and one without
    When the scaffolds are written
    Then the first project gets openapi.yaml.tmpl and the second does not

