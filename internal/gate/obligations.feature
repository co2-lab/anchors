# language: en
# @anchors
#   ref: BLGTN
#   updated_at: 2026-09-26
#   layer: feature

@BLGTN
Feature: Obligations — the duties in force, resolved from packs and config, and their status across the project

  @BLGTN-B01 @unit-level
  Scenario: The duties in force are the resolved list
    Given a project adopting the lgpd pack with two duties, and one inline duty "local"
    When the duties in force are asked from outside the package
    Then the answer is the same three duties the gate resolves

  @BLGTN-B02 @unit-level
  Scenario: Without packs the duties in force are the inline list
    Given no configuration, and in turn a configuration with the inline duty "local" and no pack
    When the duties in force are resolved
    Then the first gives no duty and the second gives exactly "local"

  @BLGTN-B03 @unit-level
  Scenario: The pack's duties come first, the inline ones last, with the pack's fields carried over
    Given a project adopting the lgpd pack with "lgpd-erasure" and "lgpd-log", and the inline duty "local"
    When the duties in force are resolved
    Then the order is lgpd-erasure, lgpd-log, local, and lgpd-erasure keeps its trigger "carries: pii", its target purge.ts and its identification

  @BLGTN-B04 @unit-level
  Scenario: A pack duty's reason cites the norm it comes from
    Given a pack duty with the reason "personal data must be erasable", the authority "Lei 13.709/2018" and the article "Art. 18, VI"
    When the duties in force are resolved
    Then its reason is "personal data must be erasable (Lei 13.709/2018, Art. 18, VI)", and a duty with no reason nor article has "Lei 13.709/2018" as its reason

  @BLGTN-B05 @unit-level
  Scenario: A pack that fails to load keeps the inline duties
    Given a project adopting a pack file that does not exist, and the inline duty "local"
    When the duties in force are resolved
    Then the error output names the pack failure, and the duties in force are exactly "local"

  @BLGTN-B06 @unit-level
  Scenario: The pack duties are read once per root
    Given a project whose lgpd pack was resolved once, and whose pack file was then removed
    When the duties in force are resolved again for the same root
    Then the two pack duties are still returned

  @BLGTN-B07 @unit-level
  Scenario: The report gives one status per duty, in the order handed
    Given the duties "pii-purgavel" and "untriggered", both targeting purge.ts
    When the compliance report runs
    Then it returns two statuses, pii-purgavel then untriggered, each with the target purge.ts

  @BLGTN-B08 @unit-level
  Scenario: A node is a subject only when its header carries the duty's trigger
    Given six model specs carrying "pii", one spec without the attribute, and the duty "untriggered" with no trigger
    When the compliance report runs
    Then pii-purgavel has 6 subjects and untriggered has none

  @BLGTN-B09 @unit-level
  Scenario: A waiver with a reason counts as fulfilled and as waived
    Given a model spec carrying "pii" that waives pii-purgavel with the reason "shared data"
    When the compliance report runs
    Then that node counts once as waived and once among the 3 fulfilled

  @BLGTN-B10 @unit-level
  Scenario: A duty declared pending counts as debt
    Given a model spec carrying "pii" that declares pii-purgavel pending for "phase 3"
    When the compliance report runs
    Then the duty's debt is 1

  @BLGTN-B11 @unit-level
  Scenario: Every other subject is missing, and the missing list is sorted
    Given the model specs Zeta and Alpha carrying "pii" and absent from purge.ts
    When the compliance report runs
    Then the missing list is Alpha then Zeta

  @BLGTN-I01 @unit-level
  Scenario: Every subject is counted exactly once
    Given six subjects: two that appear in purge.ts, one waived, one in debt and two missing
    When the compliance report runs
    Then fulfilled 3, debt 1 and missing 2 add up to the 6 subjects

  @BLGTN-X01 @unit-level
  Scenario: The report evaluates only the duties it is handed
    Given a configuration adopting the lgpd pack, and the single handed duty "only"
    When the compliance report runs
    Then it returns exactly one status, for "only"

  @BLGTN-E02 @unit-level
  Scenario: A node whose file cannot be read is not a subject
    Given a map node models/Gone.spec.md with no file on disk
    When the compliance report runs
    Then it is not among the 6 subjects nor among the missing
