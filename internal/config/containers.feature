# language: en
# @anchors
#   ref: CNTNR
#   updated_at: 2026-09-26
#   layer: feature

@CNTNR
Feature: Containers — what runs separately, and which layers run inside each

  @CNTNR-B01 @unit-level
  Scenario: The declared containers come back as written, and a missing config has none
    Given a configuration declaring the containers app, api and an external db, with app talking to api over HTTPS/JSON
    When the containers are asked for, and asked of a missing configuration
    Then the three come back in order with the talk intact, and the missing configuration has none

  @CNTNR-B02 @unit-level
  Scenario: The internal containers are the declared ones without the external, in declared order
    Given the containers app, api and the external db
    When the internal containers are asked for
    Then they are app and api, in that order

  @CNTNR-B03 @unit-level
  Scenario: A layer is found in its container ignoring case and surrounding spaces
    Given the api container declaring the layer " Handler "
    When the container of the layer "handler" is asked for
    Then it is api

  @CNTNR-B04 @unit-level
  Scenario: A layer no container claims has no container, and that is an answer, not an error
    Given containers that do not declare the layer "lambda"
    When the container of "lambda" is asked for
    Then the answer is empty

  @CNTNR-B05 @unit-level
  Scenario: The layers of an external container are still claimed by it
    Given the external db container declaring the layer "table"
    When the container of "table" is asked for
    Then it is db

  @CNTNR-B06 @unit-level
  Scenario: The orphan layers are the given ones no container claims, in the given order
    Given the existing layers screen, lambda, service and infra, of which only screen and service are declared
    When the orphan layers are asked for
    Then they are lambda and infra, in that order

  @CNTNR-I01 @unit-level
  Scenario: A layer is an orphan exactly when it has no container
    Given the layers screen, HANDLER, table, lambda, service and infra
    When each is classified by its container and by the orphan list
    Then every layer with a container is absent from the orphans, and every layer without one is present

  @CNTNR-X01 @unit-level
  Scenario: With no container declared, no layer is placed by guessing: every layer is an orphan
    Given a configuration with no containers block
    When the orphan layers of the layer "a" are asked for
    Then "a" is an orphan
