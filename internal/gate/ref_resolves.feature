# language: en
# @anchors
#   ref: RFRSR
#   updated_at: 2026-09-26
#   layer: feature

@RFRSR
Feature: RefResolves — the reference points at the spec that really describes the unit

  @RFRSR-B01 @unit-level
  Scenario: A reference that does not match the sibling spec fails
    Given a code file whose header references an identity from before a refactoring
    And the spec co-located with it declares a different identity
    When the gate confronts it
    Then it returns Fail, because a wrong reference looks like traceability and attributes
      the unit to the wrong spec

  @RFRSR-B02 @unit-level
  Scenario: The failing verdict names both sides of the divergence
    Given a code file whose reference diverges from its sibling spec
    When the gate confronts it
    Then the message carries what is written and what the sibling declares, so the reader
      does not have to open two files to see the pair

  @RFRSR-B03 @unit-level
  Scenario: A reference equal to the sibling spec passes
    Given a code file whose reference is the identity its sibling spec declares
    When the gate confronts it
    Then it returns Pass

  @RFRSR-B04 @unit-level
  Scenario: A spec is not confronted
    Given a spec node that declares its own identity
    When the gate confronts it
    Then it returns Skip, because a spec owns an identity instead of citing one

  @RFRSR-B05 @unit-level
  Scenario: An artifact with no reference declared leaves without a verdict
    Given a code file whose header declares no reference at all
    And a sibling spec that exists and declares an identity
    When the gate confronts it
    Then it returns Skip, because charging the absence is the header gate's ruler

  @RFRSR-B06 @unit-level
  Scenario: With no sibling spec the gate goes quiet
    Given a code file with a reference and no spec co-located with it
    When the gate confronts it
    Then it returns Skip, because the missing piece is the triad gate's charge and two
      gates on one defect become noise

  @RFRSR-B07 @unit-level
  Scenario: A sibling spec that declares no identity counts as no sibling
    Given a code file with a reference
    And a spec co-located with it whose header carries no identity
    When the gate confronts it
    Then it returns Skip, because there is nothing to compare the reference against

  @RFRSR-B08 @unit-level
  Scenario: The sibling is found by the name convention
    Given a code file and a spec sharing the stem and the directory
    When the gate confronts it
    Then the spec is the one compared against, because the convention co-locates them

  @RFRSR-B09 @unit-level
  Scenario: A test file lands on the same sibling its code does
    Given a test file whose name carries the test suffix of its language
    And the spec named after the bare stem declares a matching identity
    When the gate confronts it
    Then it returns Pass, because the test suffix is stripped before the stem is computed

  @RFRSR-B10 @unit-level
  Scenario: An intermediate extension is dropped from the stem
    Given a file whose name carries a kind between the stem and the extension
    And the spec named after the bare stem declares a matching identity
    When the gate confronts it
    Then it returns Pass, because the stem stops at the first dot

  @RFRSR-B11 @unit-level
  Scenario: The reference is read whatever the comment syntax of the language
    Given three artifacts declaring the same reference under different line markers
    When the gate confronts each of them
    Then all three are read, because the header syntax is the language's and not the gate's

  @RFRSR-B12 @unit-level
  Scenario: The accepted identity length comes from the project's Structure
    Given a project that declares an identity length of its own
    When the gate confronts an artifact whose reference has that length
    Then the reference is recognised, because the pattern is built per call and not frozen
      at process start

  @RFRSR-B13 @unit-level
  Scenario: A leading dot is not a stem separator
    Given a hidden file whose name begins with a dot
    And a spec named only by the spec suffix sitting beside it
    When the gate confronts it
    Then it returns Skip, because truncating at a leading dot would leave an empty stem and
      attribute every hidden file to that one spec

  @RFRSR-B14 @unit-level
  Scenario: A reference to a code that exists nowhere fails
    Given a code file whose reference points to a code declared nowhere in the project
    And no sibling spec exists on disk
    And a project graph is provided
    When the gate confronts it
    Then it returns Fail, naming the invented code

  @RFRSR-B15 @unit-level
  Scenario: Without sibling spec, an existing code still skips
    Given an artifact whose reference matches a declared code in the graph
    And no sibling spec exists on disk
    When the gate confronts it
    Then it returns Skip, because an infra file with no spec is legitimate as long as the code is real

  @RFRSR-B16 @unit-level
  Scenario: An inferred identity does not satisfy the reference
    Given an artifact referencing an identity that only appears as inferred in the graph
    And no sibling spec exists on disk
    When the gate confronts it
    Then it returns Fail, because an inferred identity owns nothing

  @RFRSR-B17 @unit-level
  Scenario: Without a graph, absence is not asserted
    Given an artifact with a reference and no sibling spec on disk
    And no project graph is provided to the gate
    When the gate confronts it
    Then it returns Skip, because without a graph absence cannot be asserted

  @RFRSR-I01 @unit-level
  Scenario: The ruler is the sibling on disk, never the map
    Given a divergent pair and no graph built
    When the gate confronts it
    Then it still returns Fail, because the convention is legible from the filesystem alone

  @RFRSR-I02 @unit-level
  Scenario: The gate never repairs what it points at
    Given a code file whose reference diverges from its sibling spec
    When the gate confronts it
    Then the file on disk is unchanged, because a gate that fixed the defect would pass on
      its second run

  @RFRSR-X01 @unit-level
  Scenario: The gate does not charge the absence of the reference field
    Given a code file with a sibling spec and no reference field whatsoever
    When the gate confronts it
    Then it returns Skip, because two gates on one defect produce two messages for one fix

  @RFRSR-X02 @unit-level
  Scenario: The gate does not charge the absence of the sibling spec
    Given a code file whose unit has no spec anywhere in the project
    When the gate confronts it
    Then it returns Skip, because the missing piece of a triad is the triad gate's charge

  @RFRSR-X03 @unit-level
  Scenario: The gate does not consult the map to resolve the reference
    Given a graph that contradicts what the sibling spec on disk declares
    When the gate confronts the artifact
    Then the verdict follows the file on disk, because a stale build is exactly the defect
      this gate exists to catch

  @RFRSR-X04 @unit-level
  Scenario: The gate does not judge whether the spec describes the unit well
    Given a code file whose reference matches its sibling spec
    And that spec describes behaviour the code does not have
    When the gate confronts it
    Then it returns Pass, because the ruler here is identity — judging the content belongs
      to another gate
