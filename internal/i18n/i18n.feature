# language: en
# @anchors
#   code: INFTA
#   ref: INCTA
#   updated_at: 2026-10-03
#   layer: feature

@INCTA
Feature: I18nCatalog — every message a person reads, in the project's language, with a fallback that never goes blank

  @INCTA-B01 @unit-level
  Scenario: The supported languages and the default
    Given the message catalog
    When the supported languages and the default are read
    Then they are "pt-BR", "en" and "es", and the default is "en"

  @INCTA-B02 @unit-level
  Scenario: An empty language is the default
    Given the current language "pt-BR"
    When the empty language is set
    Then no error is returned and the current language is "en"

  @INCTA-B03 @unit-level
  Scenario: Messages come out in the current language with their arguments
    Given the languages pt-BR, en and es in turn
    When "freeze.done" and "freeze.blocked.reason" with "o plano 0002 quebrou" are resolved
    Then "freeze.done" holds "CONGELADO", "FROZEN" and "CONGELADO", and every reason holds "o plano 0002 quebrou"

  @INCTA-B04 @unit-level
  Scenario: A key missing in the current language falls back to English
    Given the current language "es" whose catalog lacks "freeze.done"
    When "freeze.done" is resolved
    Then the English text holding "FROZEN" is returned

  @INCTA-B05 @unit-level
  Scenario: A key missing everywhere resolves to itself
    Given the current language "pt-BR"
    When the key "nao.existe.esta.chave" is resolved
    Then the text is "nao.existe.esta.chave"

  @INCTA-B06 @unit-level
  Scenario: A key's values across languages come back once each
    Given the key "section.title.states", which reads "Estados" in pt-BR and es and "States" in en
    When all its translations are asked
    Then they are "Estados" then "States"

  @INCTA-B07 @unit-level
  Scenario: A key resolves in a given language without changing the current one
    Given the current language "en"
    When "section.title.rules" is resolved in "pt-BR", and a missing key is resolved in "es"
    Then the first is "Regras", the second is empty, and the current language is still "en"

  @INCTA-B08 @unit-level
  Scenario: A written title gives back its key and language
    Given the titles "  visão geral ", "Reglas", "Coisas do projeto" and the empty title
    When their keys are looked up
    Then they give "section.title.overview" in "pt-BR", "section.title.rules" in "es", and nothing for the last two

  @INCTA-I01 @unit-level
  Scenario: Every language has a catalog with the same keys
    Given the catalogs of every supported language
    When their keys are compared with the default's
    Then every catalog is non-empty and holds exactly the same keys

  @INCTA-E01 @unit-level
  Scenario: An unsupported language is refused with the options
    Given the current language "en"
    When the language "klingon" is set
    Then an error lists "pt-BR", "en" and "es", and the current language is still "en"
