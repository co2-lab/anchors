# language: en
# @anchors
#   code: LDCFB
#   ref: LDCFS
#   updated_at: 2026-11-19
#   layer: feature

@LDCFS
Feature: Workload — what each kind of run costs, so heavy runs wait for one another

  @LDCFS-B01 @unit-level
  Scenario: load: declares levels over the defaults, tasks by scope in the order written, wait limits and CI
    Given heavy redeclared, a level huge, a wait limit of 10m, ci true, mutation scoped by backend then packages, and test light
    When the levels, tasks and limits are read
    Then heavy and huge cost what is declared and light its default, the first matching scope wins, mutation waits 20m and test 10m

  @LDCFS-E01 @unit-level
  Scenario: A level with no CPU or an unreadable size, a task the queue does not know, an undeclared level, a bad pattern or duration refuse the configuration
    Given a level with 0 CPU, a memory "lots", a task lint, a level enormous, a pattern "[", a wait limit soon, a negative wait limit and a key level
    When the configuration is loaded
    Then it does not load, naming the field
