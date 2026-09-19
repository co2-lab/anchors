# language: en
# @anchors
#   ref: EXCMX
#   updated_at: 2026-09-19
#   layer: feature

@EXCMX
Feature: ExternalCommand — executes external tools via shell passing targets as positional arguments

  @EXCMX-B01 @unit-level
  Scenario: An external command exiting with status zero returns Pass
    Given an external command that exits with code 0
    When RunExternalArgs executes the command
    Then it returns Pass with empty detail

  @EXCMX-B02 @unit-level
  Scenario: An external command exiting with non-zero status returns Fail with output
    Given an external command that exits with code 1 and writes to stderr or stdout
    When RunExternalArgs executes the command
    Then it returns Fail containing the trimmed execution output

  @EXCMX-B03 @unit-level
  Scenario: An external command failing with empty output reports execution error
    Given an external command that fails without producing any stdout or stderr
    When RunExternalArgs executes the command
    Then it returns Fail with a detail explaining that no output was produced

  @EXCMX-B04 @unit-level
  Scenario: Single node execution delegates to RunExternalArgs
    Given a node with an identifier
    When runExternal is called on that node
    Then it executes RunExternalArgs passing the node ID as single target

  @EXCMX-B05 @unit-level
  Scenario: Placeholder file is rewritten to positional parameter
    Given a command template containing the placeholder file
    When RunExternalArgs prepares the script
    Then it replaces the placeholder with positional parameter one

  @EXCMX-B06 @unit-level
  Scenario: Placeholder files is rewritten to all positional parameters
    Given a command template containing the placeholder files
    When RunExternalArgs prepares the script
    Then it replaces the placeholder with all positional parameters

  @EXCMX-B07 @unit-level
  Scenario: Execution without targets runs once for project scope
    Given an empty list of targets
    When targets are sliced for command line execution
    Then exactly one batch with no targets is returned

  @EXCMX-B08 @unit-level
  Scenario: Targets within budget run in a single batch
    Given a list of targets whose total length fits within budget
    When targets are sliced
    Then all targets are returned in a single batch

  @EXCMX-B09 @unit-level
  Scenario: Targets exceeding budget are partitioned across batches
    Given a large list of targets exceeding the command line budget
    When targets are sliced
    Then they are divided into multiple batches without losing any targets

  @EXCMX-B10 @unit-level
  Scenario: Failure in any batch causes entire execution to fail
    Given targets partitioned into multiple batches where one batch fails
    When RunExternalArgs executes across all batches
    Then it returns Fail and includes the failure details

  @EXCMX-B11 @unit-level
  Scenario: Single target failure detail is truncated at five hundred characters
    Given a single target execution producing more than five hundred characters of output
    When RunExternalArgs finishes execution
    Then the returned detail is truncated to five hundred characters with a truncation marker

  @EXCMX-B12 @unit-level
  Scenario: Batch failure detail is truncated at four thousand characters
    Given multiple targets producing more than four thousand characters of output
    When RunExternalArgs finishes execution
    Then the returned detail is truncated to four thousand characters with a truncation marker

  @EXCMX-B13 @unit-level
  Scenario: Environment variable overrides argv limit
    Given the environment variable ANCHORS_ARGV_MAX set to a positive integer
    When argvLimit is determined
    Then it returns the configured integer value

  @EXCMX-I01 @unit-level
  Scenario: Target paths are passed strictly in argv preventing injection
    Given a target file path containing shell metacharacters
    When RunExternalArgs executes the command referencing positional parameter one
    Then the filename is treated as literal data and no injected command runs

  @EXCMX-I02 @unit-level
  Scenario: Oversized single target is isolated in its own batch
    Given a single target path whose length exceeds the budget ceiling
    When targets are sliced
    Then the target is placed in its own batch without being dropped

  @EXCMX-I03 @unit-level
  Scenario: Platform default argv limits are enforced
    Given no environment override is set
    When argvLimit is computed
    Then it respects the platform budget ceiling

  @EXCMX-X01 @unit-level
  Scenario: The gate does not parse or interpret linter diagnostics
    Given an external command returning raw linter messages
    When RunExternalArgs executes the command
    Then it returns the unparsed output verbatim without structural interpretation

  @EXCMX-X02 @unit-level
  Scenario: The gate does not aggregate cross-file state across partitioned batches
    Given targets split across separate batches
    When RunExternalArgs executes each batch independently
    Then each batch runs in an isolated shell invocation without shared state
