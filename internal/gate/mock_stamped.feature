# language: en
# @anchors
#   ref: MCSTM
#   updated_at: 2026-09-19
#   layer: feature

@MCSTM
Feature: MockStamped — the double carries the mark of the snippet it replaces, and the gate recomputes it

  @MCSTM-B01 @unit-level
  Scenario: An artifact that is not a test leaves without a verdict
    Given a node whose kind is spec, code or feature
    When the gate confronts it
    Then it returns Skip, because only a test declares doubles

  @MCSTM-B02 @unit-level
  Scenario: A stamp that matches the module today passes
    Given a test doubling a governed module
    And a stamp whose hash is the hash of the anchored snippet as it stands today
    When the gate confronts it
    Then it returns Pass

  @MCSTM-B03 @unit-level
  Scenario: A snippet that changed since the stamp was written fails
    Given a stamp taken from an older version of the module, with one parameter fewer
    When the gate confronts the test
    Then it returns Fail, and the verdict shows the value the contract has today

  @MCSTM-B04 @unit-level
  Scenario: The stamp is immune to displacement
    Given a module with two new comment lines inserted above the anchored snippet
    And a stamp written before the insertion
    When the gate confronts the test
    Then it returns Pass, because the anchor is searched by content and moving a snippet
      does not change its contract

  @MCSTM-B05 @unit-level
  Scenario: An anchor that vanished fails with its own message
    Given the anchored function was renamed in the real module
    When the gate confronts the test
    Then it returns Fail naming the anchor as renamed or removed, because it is a finding
      and not a tool error — the double is certainly out of date

  @MCSTM-B06 @unit-level
  Scenario: An anchor occurring more than once is ambiguous and fails
    Given the anchor line appears twice in the real module
    When the gate confronts the test
    Then it returns Fail naming the ambiguity, because a stamp pointing at one of the two
      proves nothing and the gate reports rather than guessing

  @MCSTM-B07 @unit-level
  Scenario: The declared line count delimits the window
    Given a module whose last body line changed
    And one stamp covering two lines from the anchor and another covering nine
    When the gate confronts each of them
    Then the two-line window passes and the nine-line window fails, because the reach is
      declared in plain sight instead of guessed by a per-language parser

  @MCSTM-B08 @unit-level
  Scenario: Without the dialect declared the gate goes quiet
    Given a project that declares no double-detection dialect
    When the gate confronts a test carrying a stamp
    Then it returns Skip naming what has to be declared, because adopting the stamp is
      the project's decision

  @MCSTM-B09 @unit-level
  Scenario: A double of a module the project does not govern is not charged
    Given a test doubling a third-party library with no stamp
    When the gate confronts it
    Then it returns Skip, because the drift this gate pursues is the neighbour changing,
      and a locked dependency is not the neighbour

  @MCSTM-B10 @unit-level
  Scenario: A stamp whose module no longer exists is a finding, not a crash
    Given a stamp pointing at a file that is not on disk
    When the gate confronts the test
    Then it returns Fail naming the file as not found

  @MCSTM-B11 @unit-level
  Scenario: The absence of a stamp on a governed double is accused
    Given a test doubling a module the project governs, with no stamp at all
    When the gate confronts it
    Then it returns Fail naming the unstamped module, because a divergent stamp accuses
      while an absent one is the silence that lets the mechanism protect only whoever
      already chose to be protected

  @MCSTM-B12 @unit-level
  Scenario: A stamp present and correct satisfies both charges
    Given a test doubling a governed module
    And the stamp is present and its hash matches the module today
    When the gate confronts it
    Then it returns Pass, satisfying the charge of absence and the charge of correspondence
      in one verdict

  @MCSTM-B13 @unit-level
  Scenario: A dialect regex that does not compile fails loudly
    Given a project declaring a double-detection pattern that is not a valid regex
    When the gate confronts a test
    Then it returns Fail saying the pattern does not compile, because silencing it would
      make the gate sweep zero doubles and report green

  @MCSTM-B14 @unit-level
  Scenario: A dialect regex with no capture group fails
    Given a project declaring a pattern that matches the call and captures nothing
    When the gate confronts a test
    Then it returns Fail saying a capture group is required, because without it the gate
      cannot know which module was doubled

  @MCSTM-B15 @unit-level
  Scenario: Another ecosystem's dialect is charged the same way
    Given a project declaring the Python patch dialect
    And a Python test doubling a governed module with no stamp
    When the gate confronts it
    Then it returns Fail naming the unstamped module, the same charge in another dialect

  @MCSTM-I01 @unit-level
  Scenario: The gate recomputes the hash instead of validating the stamp's format
    Given a stamp that is well formed and was correct when written
    And the real module changed afterwards
    When the gate confronts the test
    Then it returns Fail, because a stamp nobody confronts would certify itself once its
      author regenerated it to match their own mock

  @MCSTM-I02 @unit-level
  Scenario: The double and the stamp are tied by the module's path
    Given a double written through a project alias with no extension
    And a stamp written with the module's path on disk, with extension
    When the gate confronts the test
    Then the charge is satisfied, because both describe the same file and a per-ecosystem
      alias resolver would be needed to see it otherwise

  @MCSTM-I03 @unit-level
  Scenario: A configuration fault fails and a project decision skips
    Given a pattern that does not compile, a pattern with no capture group, and no pattern at all
    When the gate confronts a test under each of them
    Then the first two return Fail and the third returns Skip, because the difference is
      whether somebody CHOSE the silence

  @MCSTM-X01 @unit-level
  Scenario: The gate does not interpret the code of the stamped module
    Given a module whose anchored snippet changed only in a constant inside the body
    When the gate confronts the test
    Then it returns Fail, because it hashes text — a signature extractor would miss this
      and would need to exist once per language

  @MCSTM-X02 @unit-level
  Scenario: The gate carries no built-in dialect for detecting doubles
    Given a project in an ecosystem where a double is not written as a detectable call
    And no double-detection dialect declared
    When the gate confronts a test that doubles a governed module
    Then it returns Skip rather than Pass, because reporting green over what was never
      checked is the worst possible failure in a measuring device

  @MCSTM-X03 @unit-level
  Scenario: The gate does not skip the absence of a stamp to accommodate legacy code
    Given a test that doubles a governed module with no stamp, in a project that predates the gate
    When the gate confronts it
    Then it returns Fail, because designing for the migration case would turn it into a
      permanent property of the framework — legacy is handled by a non-blocking gate
      during adoption, not by a looser ruler

  @MCSTM-X04 @unit-level
  Scenario: The gate does not guarantee cryptographic strength
    Given a stamp whose hash is the truncated digest of the snippet
    When the gate recomputes it
    Then the comparison is over the truncated value, because the stamp is read by a human
      and there is no adversary forging a collision against their own test

  @MCSTM-B16 @unit-level
  Scenario: A module with a dot in its name is matched to its stamp
    Given a test doubling `@/src/stores/auth.store`
    And a stamp on `src/stores/auth.store.ts`
    When the gate confronts the test
    Then the double counts as stamped

  @MCSTM-B17 @unit-level
  Scenario: Changing a module checks the doubles stamped against it
    Given a test whose stamp points at a module
    And another test with no stamp on that module
    When the module is the changed file of a check
    Then the stamping test enters the check, and the other does not
    And a drift names the module file of the stamp

  @MCSTM-B18 @unit-level
  Scenario: The author of a change refreshes the stamps and gets the doubles to adjust
    Given two members of a module stamped by a test
    And one of them changed
    When the author runs the stamp refresh for the module
    Then only the stamp of the changed member is listed, with the old and the new hash
    And its hash is updated, and the gate passes again
    And a stamp whose anchor is gone is listed and not refreshed
