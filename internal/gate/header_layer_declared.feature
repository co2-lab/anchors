# language: en
# @anchors
#   ref: HDLYD
#   updated_at: 2026-09-28
#   layer: feature

@HDLYD
Feature: HeaderLayerDeclared — the layer a header declares is one the Estrutura has

  @HDLYD-B01 @unit-level
  Scenario: A layer the Estrutura does not have fails
    Given headers declaring landing-feature, landing-component, and a body line declaring another layer
    When each is confronted with an Estrutura that has landing-feature and api
    Then the first passes, the second fails naming it and the declared layers, and the body line is not read

  @HDLYD-B02 @unit-level
  Scenario: A layer that differs only in case fails naming both
    Given a header declaring Landing-Feature
    When it is confronted
    Then it fails naming Landing-Feature and landing-feature

  @HDLYD-B03 @unit-level
  Scenario: Nothing to confront is skipped
    Given a header with no layer, a TODO layer, and a project with no layers
    When each is confronted
    Then each is skipped
