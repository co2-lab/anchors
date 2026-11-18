# language: en
# @anchors
#   code: MDFMP
#   ref: MDCMP
#   layer: feature

@MDCMP
Feature: MapDeps — the dependency tree of a file

  @MDCMP-B01 @unit-level
  Scenario: The tree goes down what a file uses, and up who uses it, a cycle marked once
    Given a screen using the tokens, which use the colors, which use the tokens
    When the tree is printed down from the screen, and up one level from the tokens
    Then the cycle is marked once, and up shows the screen and the colors

  @MDCMP-B02 @unit-level
  Scenario: A file is named by its own code, its unit's code or its path
    Given the same map
    When a file is named by its code, its unit's code, its path, and a name nobody owns
    Then each resolves to its file, and the last to none

  @MDCMP-B03 @unit-level
  Scenario: The resources of a kind and the files that reach them
    Given two files reaching the transactions table and one reaching the Stripe API
    When the db resources are listed
    Then transactions is listed with both files, and Stripe is not

  @MDCMP-B04 @unit-level
  Scenario: --flat prints the closure as sorted paths, one per line
    Given a handler that uses a service, which uses a client, and the client uses the service back
    When the closure is asked down, up from the client, and down with depth 1
    Then each file of the closure appears once, sorted, the start included, and depth 1 stops at the service
