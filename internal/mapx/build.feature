# language: en
# @anchors
#   code: BLFTA
#   ref: GRBLG
#   updated_at: 2026-10-04
#   layer: feature

@GRBLG
Feature: GraphBuild — projecting the declared structure onto the scanned files

  @GRBLG-B01 @unit-level
  Scenario: Each scanned file becomes a node with its layer's tags and regime
    Given a code file of layer screen, whose layer carries the tag frontend and the regime behavioural, and a date of 2026-09-01 for it
    When the graph is built
    Then its node carries its path, kind, revision, layer screen, date 2026-09-01, the tag frontend and the regime behavioural

  @GRBLG-B02 @unit-level
  Scenario: The declared identity wins over cited codes and over the anchor
    Given a file whose header declares MTENX while its text first cites DTAXX-B11, and a test declaring DECLR next to an anchor declaring ANCOR
    When their identities are decided
    Then the first is MTENX and the test is DECLR

  @GRBLG-B03 @unit-level
  Scenario: A derived file takes its sibling anchor's identity
    Given a GoLive spec declaring GLCGL, with GoLive code, test and feature beside it, and an Outra test in the same directory
    When the identities of the derived files are decided
    Then the GoLive code, test and feature are GLCGL, and the Outra test gets nothing from it

  @GRBLG-B04 @unit-level
  Scenario: With no header and no anchor, the first code's root is the identity
    Given a file with no header whose text carries ABCDX-B01, and a file with neither header nor codes
    When their identities are decided
    Then the first is ABCDX and the second is empty

  @GRBLG-B05 @unit-level
  Scenario: A vendored file has no local identity
    Given the vendored anchors-board pipeline, whose header says FNDTN
    When the graph is built
    Then its node is marked vendored and carries no code

  @GRBLG-B06 @unit-level
  Scenario: Only a header marks the identity as declared
    Given a spec whose header declares LOGIX and a test with no header that cites LOGIX-A01
    When the graph is built
    Then the spec's node is marked declared and the test's node is not

  @GRBLG-B07 @unit-level
  Scenario: The pieces of one unit are linked
    Given Login code, spec, feature and test in one directory, with templates for spec, feature and test
    When the graph is built
    Then the spec specifies the code, the spec is covered by the feature, and the feature is tested by the test

  @GRBLG-B08 @unit-level
  Scenario: With the code as anchor, the relations still go down from the spec
    Given a configuration whose anchor is the code, and a handler with its spec and feature beside it
    When the graph is built
    Then the spec is covered by the feature, found from the handler as a derived file

  @GRBLG-B09 @unit-level
  Scenario: Without a feature, the spec is tested by the test
    Given a Login spec anchor with its code and test and no feature template
    When the graph is built
    Then the spec specifies the code, the spec is tested by the test, and the code does not specify the spec

  @GRBLG-B10 @unit-level
  Scenario: A layer override reaches the spec through the layer it declares
    Given a spec of file layer spec declaring layer screen, and an override for screen that expects a tsx code file
    When the graph is built
    Then the spec specifies the tsx code file

  @GRBLG-B11 @unit-level
  Scenario: A code override replaces the templates of the kinds it declares, and the others fall back to the default
    Given a TSCTY spec whose code override names only packages/*/tsconfig.json, while the default templates name a ts file, a feature and a test
    When the graph is built
    Then the spec specifies each packages tsconfig.json and not the ts file, and is covered by its default feature, tested by its default test
    And with the override also naming config/tsconfig.feature and an empty test list, the spec is covered by that feature and no test is linked

  @GRBLG-B12 @unit-level
  Scenario: A bracketed directory is literal and a template wildcard expands
    Given a SeloClient unit inside the directory app/selo/[slug], and a spec whose template expands to two tsconfig.json files
    When the graph is built
    Then the SeloClient unit is linked, and the spec specifies both tsconfig.json files

  @GRBLG-B13 @unit-level
  Scenario: A shared own code links spec and test across directories only
    Given the ORDER spec in one directory and a test in another carrying ORDER-B01, and a spec citing ORDER-B01 whose own identity is CITES
    When the graph is built
    Then the ORDER spec is tested by that test through an inferred relation, and the citing spec is not

  @GRBLG-B14 @unit-level
  Scenario: A guide governs the layers of its tag only, and never itself
    Given a spec guide ruling tag spec, a frontend guide ruling tag frontend, and a guide layer that also carries tag spec
    When the graph is built
    Then the spec guide governs the spec and the frontend guide the screen code, the frontend guide does not govern the spec, and no guide governs itself

  @GRBLG-B15 @unit-level
  Scenario: A dependency row becomes a relation to an existing file
    Given a Login spec whose dependency table names an auth store, a useAuth hook and a file that does not exist
    When the relations are built
    Then there are two depends-on relations from the spec, carrying their method and row code, and none to the missing file

  @GRBLG-B16 @unit-level
  Scenario: A seed path is exact and a bare name must be unique
    Given a plan seeding Tela.spec.md by name, apps/y/Outra.spec.md by path, and Ambigua.spec.md which exists twice
    When the relations are built
    Then the plan seeds apps/x/Tela.spec.md and apps/y/Outra.spec.md, and neither Ambigua spec

  @GRBLG-B17 @unit-level
  Scenario: A need links only an existing plan
    Given plan 0002 that needs plan 0001, which exists, and plan 0009, which does not
    When the graph is built
    Then plan 0002 needs plan 0001 and has no relation to plan 0009

  @GRBLG-B18 @unit-level
  Scenario: Realized doctrine and flag scenarios resolve by unit code
    Given a spec realizing LIMIT-R03 and NOPE-R01, gated by CHKUT-G02, with a doctrine declaring LIMIT and a flag declaring CHKUT
    When the graph is built
    Then the spec realizes the doctrine and is gated by the flag, each carrying the local and target rule, and nothing links for NOPE

  @GRBLG-B19 @unit-level
  Scenario: Nodes and relations are sorted
    Given files arriving out of order
    When the graph is built
    Then the nodes are in path order and the relations in source, target and type order

  @GRBLG-B20 @unit-level
  Scenario: A rebuild keeps the stamps and judgments of surviving relations
    Given a previous graph whose spec-to-test relation carries a stamp and a my-gate judgment, and whose other relation changed type
    When a rebuilt graph takes over from it
    Then the surviving relation keeps its stamp and judgment, and the relation that changed type keeps nothing

  @GRBLG-B21 @unit-level
  Scenario: A rebuild keeps a signal only for an unchanged file
    Given a previous graph where two test files carried signals, and one of them changed revision since
    When a rebuilt graph takes over from it
    Then the unchanged file keeps its signal and the changed one has none

  @GRBLG-I01 @unit-level
  Scenario: The build does not depend on the order of the files
    Given the Login files and guides
    When the graph is built from them in one order and in the reverse order
    Then the two graphs are identical

  @GRBLG-X01 @unit-level
  Scenario: With no dates given, nodes carry no date
    Given the Login files and no dates
    When the graph is built
    Then every node's date is empty

  @GRBLG-X02 @unit-level
  Scenario: No relation points to a file that was not scanned
    Given declarations naming missing files: a dependency row, a seed path and a need
    When the graph is built
    Then no relation points to any of the missing files

  @GRBLG-B22 @unit-level
  Scenario: A support file becomes a node marked as support
    Given a scanned test file marked as support and one that is not
    When the map is built
    Then the first node is marked as support, the second is not, and both keep their kind

  @GRBLG-I02 @unit-level
  Scenario: The edges are in a total order
    Given a spec with two dependency rows on the same file
    When the graph is built with the rows in either order
    Then both graphs list the edges in the same order

  @GRBLG-B23 @unit-level
  Scenario: Signals are filled from another map at the same revision
    Given a map with one node without signal, one with its own, and one at another revision, and a source map with signals for all three
    When the signals are filled from the source
    Then only the first gets the source's signal

  @GRBLG-B24 @unit-level
  Scenario: A file's own code is its file code, and a ref keeps its unit
    Given a spec with code ARENA, its screen with "code: ARSCR" and "ref: ARENA", and a helper with "code: TOKNS" and no ref
    When the map is built
    Then the screen's unit is ARENA and its file code ARSCR, and the helper's unit and file code are TOKNS


  @GRBLG-B25 @unit-level
  Scenario: A file with a code of its own and a ref links across directories through the unit it refs
    Given a feature and a test in different directories, each with a code of its own and ref LPSTI, the test also citing OTHER-B01, and OTHER's feature elsewhere
    When the map is built
    Then the feature of LPSTI is tested by the test, and OTHER's feature is not
