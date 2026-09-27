# language: en
# @anchors
#   ref: ISLFS
#   updated_at: 2026-09-26
#   layer: feature

@ISLFS
Feature: IssueLifecycle — a divergence recorded so it survives the session, with its state as a folder

  @ISLFS-B01 @unit-level
  Scenario: The key is stable across dates and distinct per gate
    Given a violation of "spec-sections" on "features/x/A.spec.md" detected on 2026-08-07 and again on 2027-01-01
    And the same target violating "spec-has-code"
    When their keys are computed
    Then the two detections share a key and the other gate has a different one

  @ISLFS-B02 @unit-level
  Scenario: The file name is the date and the key, with no slash
    Given a violation of "spec-sections" on "features/x/A.spec.md" detected on 2026-08-07
    When its file name is computed
    Then it starts with "2026-08-07--violation--" and holds no slash

  @ISLFS-B03 @unit-level
  Scenario: The body names the kind, the target, the gate and the detail
    Given a violation of "spec-sections" on "features/x/A.spec.md" with the detail "falta a seção Regras"
    When it is opened
    Then its body holds "VIOLATION", the target, the gate and the detail

  @ISLFS-B04 @unit-level
  Scenario: A new issue is opened in todo
    Given a project with no issues
    When a violation is opened
    Then it is created in "todo/"

  @ISLFS-B05 @unit-level
  Scenario: The same issue is not opened twice
    Given a violation already open in "todo/"
    When the same violation is opened with another date
    Then nothing is created, the state reported is todo, and "todo/" holds one issue

  @ISLFS-B06 @unit-level
  Scenario: An assumed debt is born in future and shows when it is due
    Given a violation declared as a debt due "`lgpd-eliminacao` — na etapa de código da Fase 1"
    When it is opened in future, opened again, and resolved
    Then it is born in "future/", the second open creates nothing, it ends in "done/", and its body says "When it will be paid" and "ASSUMED debt"

  @ISLFS-B07 @unit-level
  Scenario: A decision explains how to close it
    Given an open decision on "b.spec.md"
    When its body is rendered
    Then it says how to close it and what not to do

  @ISLFS-B08 @unit-level
  Scenario: Resolving moves a live issue to done, and only once
    Given a violation open in "todo/", and another one in "doing/"
    When each is resolved, and the first is resolved again
    Then both are in "done/" and the second resolution of the first does nothing

  @ISLFS-B09 @unit-level
  Scenario: A new finding reopens the issue and keeps the old report
    Given a violation open with the report "Laudo A"
    When it is reopened with the report "Laudo B", and then with "Laudo B" again
    Then the issue in "todo/" holds both reports, and "Laudo B" appears once

  @ISLFS-B10 @unit-level
  Scenario: Issues are listed by owner, and no owner means the agent
    Given a violation of the agent, a decision of the user, and an old issue file with no owner line
    When the open issues are listed by owner and the old file's owner is read
    Then each owner gets only its issue and the old file is the agent's

  @ISLFS-B11 @unit-level
  Scenario: Reassigning hands the issue over with its reason
    Given an agent's violation, and an old issue with the label "alvo (regido)" and no owner line
    When each is reassigned to the user with a reason
    Then both are the user's, the violation is still a violation with the reason recorded, and the old issue gained an owner line

  @ISLFS-B12 @unit-level
  Scenario: A full check closes the violations it no longer reproduces
    Given open violations of the gates "header-conforms" (reproduced) and "header-conforme" (renamed), and a violation in doing
    When the issues are reconciled with the reproduced key
    Then the renamed-gate and the doing violations are closed and the reproduced one stays open

  @ISLFS-B13 @unit-level
  Scenario: Issues are files unless GitHub is configured
    Given no backend configured
    When GitHub is configured for "acme/x" and then files again
    Then the backend is files, then the repository "acme/x", then files

  @ISLFS-I01 @unit-level
  Scenario: A resolved issue is moved, not copied, and never resurrected
    Given a violation open in "todo/"
    When it is resolved and then opened again
    Then "todo/" is empty, "done/" holds it, and opening again creates nothing and reports done

  @ISLFS-X01 @unit-level
  Scenario: Reconciling spares decisions and the user's violations
    Given an open decision and an open violation owned by the user, neither reproduced
    When the issues are reconciled
    Then both stay open

  @ISLFS-E01 @unit-level
  Scenario: An issue that cannot be written is reported
    Given a project where "issues" is a file, not a folder
    When a violation is opened
    Then an error is returned and nothing is reported created

  @ISLFS-E02 @unit-level
  Scenario: Reassigning a missing issue is reported
    Given no issue named "nope.md" in "todo/"
    When it is reassigned to the user
    Then an error is returned
