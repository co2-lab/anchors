# language: en
# @anchors
#   code: DCFTC
#   ref: DCRQA
#   updated_at: 2026-10-03
#   layer: feature

@DCRQA
Feature: DocsRequired — which aggregate documentation a change of a unit obliges to touch

  @DCRQA-B01 @unit-level
  Scenario: A trigger naming a layer charges every change in that layer, whatever the unit
    Given an OpenAPI document triggered by the layer "lambdas"
    When the documents owed by a change in "lambdas" are asked for, with no code and with the code "ANYUN"
    Then the OpenAPI document is owed both times

  @DCRQA-B02 @unit-level
  Scenario: A trigger naming a unit code charges that unit and not its neighbours in the same layer
    Given a schema document triggered by the unit code "DTSTD"
    When the documents owed by the units DTSTD and GLCGL of the layer "infra" are asked for
    Then DTSTD owes the schema and GLCGL owes nothing

  @DCRQA-B03 @unit-level
  Scenario: Asked without a unit code, the answer is the layer's alone
    Given a schema document triggered only by the unit code "DTSTD", and a document with a blank trigger
    When the documents owed by the layer "infra" are asked for without a code
    Then nothing is owed

  @DCRQA-B04 @unit-level
  Scenario: A documentation with no trigger is never owed by a unit change
    Given an architecture document with no trigger
    When the documents owed by several layers and codes are asked for
    Then the architecture document is never among them

  @DCRQA-B05 @unit-level
  Scenario: A trigger matches ignoring case and the spaces around it
    Given a document triggered by " Lambdas " and " dtstd"
    When a change in the layer "LAMBDAS", and one of the unit "DTSTD", are asked about
    Then both are triggered

  @DCRQA-B06 @unit-level
  Scenario: Every declared documentation is listed, and a project with no docs block owes none
    Given a project declaring three documents, a project with no docs block and no configuration at all
    When every declared document, and the documents owed by a change, are asked for
    Then the first lists its three documents, and the other two list and owe nothing

  @DCRQA-I01 @unit-level
  Scenario: Naming the unit never removes a documentation the layer alone owes
    Given the declared documents and the layers lambdas, infra and none
    When each layer is asked with no code and with the codes DTSTD and OTHER
    Then every document owed without a code is still owed with each code

  @DCRQA-X01 @unit-level
  Scenario: A trigger is matched whole, never as part of a longer name
    Given a document triggered by "lambdas" and "DTSTD"
    When the layers "lambda" and "lambdas-edge" and the codes "DTST" and "DTSTDX" are asked about
    Then none of them triggers it
