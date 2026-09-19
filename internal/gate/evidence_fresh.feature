# language: en
# @anchors
#   ref: EVFRV
#   updated_at: 2026-09-19
#   layer: feature

@EVFRV
Feature: EvidenceFresh — the score of this test holds against TODAY's code

  @EVFRV-B01 @unit-level
  Scenario: An artifact that is not a test leaves without a verdict
    Given a node whose kind is code, spec or plan
    When the gate confronts it
    Then it returns Skip, because only a test carries an execution score

  @EVFRV-B02 @unit-level
  Scenario: Without a built map the gate stays quiet
    Given a test node and no graph built
    When the gate confronts it
    Then it returns Skip, because there is no closure to walk and it will not approve
      what it could not look at

  @EVFRV-B03 @unit-level
  Scenario: A test with no execution stamp is skipped, not failed
    Given a test that carries no record of ever having run
    When the gate confronts it
    Then it returns Skip, because a score that was never written cannot have expired

  @EVFRV-B04 @unit-level
  Scenario: A test whose closure is intact passes
    Given a test stamped at the revision it ran on
    And every dependency it recorded still sits at the revision the run measured
    When the gate confronts it
    Then it returns Pass

  @EVFRV-B05 @unit-level
  Scenario: The passing verdict says what it checked against
    Given a test whose intact closure holds one dependency
    When the gate confronts it
    Then the verdict states the size of the closure it walked, because that is the difference
      between "nobody looked" and "I looked and it stands"

  @EVFRV-B06 @unit-level
  Scenario: A test whose dependency advanced a revision fails
    Given a test stamped at the revision it ran on
    And a dependency that has since moved to a newer revision
    When the gate confronts it
    Then it returns Fail, because the score is still written and stopped holding

  @EVFRV-B07 @unit-level
  Scenario: The failing verdict names the culprit
    Given a test whose dependency has since moved to a newer revision
    When the gate confronts it
    Then the verdict names that dependency, so whoever fixes it knows what moved underneath

  @EVFRV-B08 @unit-level
  Scenario: The failing verdict states the fix
    Given a test whose dependency has since moved to a newer revision
    When the gate confronts it
    Then the verdict says to run the test again, because accusing without saying what to do
      transfers the work to the reader

  @EVFRV-B09 @unit-level
  Scenario: A test whose own file changed is reported separately from its closure
    Given a test whose own revision has moved since the run that stamped it
    When the gate confronts it
    Then the verdict reports the test's own file as changed, apart from any closure finding

  @EVFRV-B10 @unit-level
  Scenario: The culprit list is truncated at five and the remainder counted
    Given a test whose stamped closure holds twenty dependencies and all of them moved
    When the gate confronts it
    Then the verdict lists five of them and states how many others there are

  @EVFRV-I01 @unit-level
  Scenario: Absence of proof and expired proof are never the same finding
    Given a test that never ran
    When the gate confronts it
    Then it does not return Fail, because the absent proof is a different debt with a
      different fix — run it the first time, not revalidate it

  @EVFRV-I02 @unit-level
  Scenario: A test that never ran is never approved either
    Given a test that never ran
    When the gate confronts it
    Then it does not return Pass, because approving would state a freshness nobody measured

  @EVFRV-I03 @unit-level
  Scenario: Truncation never hides the size of the problem
    Given a test whose stamped closure holds twenty dependencies and all of them moved
    When the gate confronts it
    Then the verdict counts what it did not list, so the reader still learns how far it spread

  @EVFRV-X01 @unit-level
  Scenario: The gate does not charge the absence of a green test
    Given a test that carries no record of ever having run
    When the gate confronts it
    Then it does not return Fail, because charging the missing test is the coverage gate's ruler

  @EVFRV-X02 @unit-level
  Scenario: The gate does not run the test nor judge whether the change broke it
    Given a test whose dependency moved by a change that could not affect the behaviour
    When the gate confronts it
    Then it returns Fail all the same, because the gate measures whether the evidence still
      covers the current code — the cheap fix settles it for real instead of by opinion

  @EVFRV-X03 @unit-level
  Scenario: The gate does not read the project's configuration
    Given a test with an intact closure and no configuration supplied at all
    When the gate confronts it
    Then it returns Pass, because the confronted truth lives in the map — a ruler depending
      on settings could be turned off by a default nobody chose
