# language: en
# @anchors
#   code: EDFNV
#   ref: ENVDC
#   updated_at: 2026-10-03
#   layer: feature

@ENVDC
Feature: EnvDeclared — the environment variables a unit reads are the ones its spec declares

  @ENVDC-B01 @unit-level
  Scenario: Nothing declared and nothing read leaves without a verdict
    Given a spec with no Environment Variables section whose code reads no variable
    When env-declared confronts it
    Then it leaves without a verdict

  @ENVDC-B02 @unit-level
  Scenario: A variable read and not declared is named
    Given Go code reading PAYMENTS_URL and API_KEY, and a spec declaring PAYMENTS_URL
    When env-declared confronts the spec
    Then it names API_KEY as undeclared

  @ENVDC-B03 @unit-level
  Scenario: A variable declared and not read is named, unless deprecated
    Given a spec declaring OLD_HOST, and LEGACY_TOKEN as "yes: use API_KEY", neither read
    When env-declared confronts the spec
    Then it names OLD_HOST and not LEGACY_TOKEN

  @ENVDC-B04 @unit-level
  Scenario: A read in a comment is no read
    Given code that mentions os.Getenv("DEBUG") only in a comment
    When env-declared confronts the spec
    Then DEBUG is not read

  @ENVDC-B05 @unit-level
  Scenario: Every family's reads are recognised when none is declared
    Given TypeScript code reading process.env.LOG_LEVEL in a project with no dialect
    When env-declared confronts the spec declaring LOG_LEVEL
    Then it passes

  @ENVDC-E01 @unit-level
  Scenario: A specified file that cannot be read reads nothing
    Given a spec declaring PORT that specifies a file the disk does not have
    When env-declared confronts the spec
    Then PORT is named as declared and not read
