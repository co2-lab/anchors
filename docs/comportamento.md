<!-- anchors:generated from doct/comportamento.md.tmpl — DO NOT EDIT: run `anchors docs build` -->


# Comportamento

Todos os cenários do sistema. Cada um leva à unidade que o define.

Um cenário descreve o que o sistema faz numa situação — vem da feature, e é o mesmo que o
teste prova.

## gate

- [An artifact that is not a spec leaves without a verdict](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-B01`

- [An exported symbol the spec never names fails, and the verdict names it](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-B02`

- [What the spec already catalogues is never accused](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-B03`

- [A no-rule marker with a written reason waives the symbol](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-B04`

- [A bare no-rule marker does not waive](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-B05`

- [A spec cataloguing every exported symbol passes](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-B06`

- [With no code linked the gate leaves without a verdict](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-B07`

- [Without a declared export pattern the gate skips and says so](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-B08`

- [With the pattern declared the gate confronts for real in any language](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-B09`

- [The declared dialect family also supplies the pattern](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-B10`

- [The waiver holds in the comment block above the symbol](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-I01`

- [The waiver does not leak between symbols](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-I02`

- [The gate never approves a language it cannot read](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-I03`

- [The gate does not judge whether the rule describes the symbol well](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-X01`

- [The gate does not decide which symbols deserve a rule](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-X02`

- [The gate knows no language, the project declares what is public](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-X03`

- [The gate does not charge the absence of code](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-X04`

- [An identifier in the wrong language is accused, and an English one passes](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-B01`

- [The verdict returns the word that accused](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-B02`

- [Only a DECLARATION is the subject](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-B03`

- [The declarations are found in every form the language offers](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-B04`

- [Deciding one word is separate from deciding a whole identifier](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-B05`

- [A short word does not count](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-I01`

- [The gate does not read comments](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-X01`

- [The gate does not read user-facing text](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-X02`

- [The gate does not use a dictionary to decide the language](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-X03`

- [A status emitted and not declared is accused by number](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-B01`

- [A status declared and emitted by no path is accused as a phantom](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-B02`

- [A faithful table passes](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-B03`

- [The 500 of the top-level try/catch is not charged](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-B04`

- [The 5xx range covers, the 4xx range does not](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-B05`

- [A status that lives only in a comment is not emitted](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-B06`

- [Without the contract section there is nothing to confront](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-B07`

- [Code that returns no status is skipped](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-B08`

- [A literal status passed to a local helper counts as emitted](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-B09`

- [With a dynamic status the phantom side goes quiet and the literals still count](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-B10`

- [Without a declared dialect the verdict is Pending](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-B11`

- [An explicit opt-out of the http_status field is honoured](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-B12`

- [The lexicon comes from the project's dialect, not from the gate](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-I01`

- [A dialect declared by hand teaches the gate its own lexicon](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-I02`

- [A named constant is worth the number it means](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-I03`

- [The gate does not demand the generic ranges](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-X01`

- [The gate does not judge when each status is right](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-X02`

- [The gate does not charge the phantom side under a dynamic status](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-X03`

- [Confronting an artifact that is not a spec skips](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-B01`

- [A spec containing no count declarations skips](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-B02`

- [A declared file count matching the number of files on disk passes](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-B03`

- [A declared file count that differs from disk fails](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-B04`

- [A declared regex pattern counts occurrences across files instead of file count](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-B05`

- [A declared regex occurrence count that differs from reality fails](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-B06`

- [A glob pattern matching zero files fails with a path warning](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-B07`

- [An invalid glob expression fails reporting the glob syntax error](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-B08`

- [An invalid regex pattern fails reporting the regex syntax error](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-B09`

- [Divergent count stated in adjacent prose fails even if marker matches](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-B10`

- [Non-restrictive complements attached to prose labels are confronted as total count](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-B11`

- [Qualifying words attached to prose labels indicate subsets and are not accused](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-B12`

- [Arbitrary numbers in prose without declaration markers are ignored](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-B13`

- [Numerical claims without declaration markers never trigger confrontation](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-I01`

- [Matching markers cannot conceal lying prose](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-I02`

- [Zero matched files indicates path error rather than legitimate zero count](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-I03`

- [The gate does not guess what to count from arbitrary text](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-X01`

- [Specs without count markers are not penalized](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-X02`

- [Prose phrases qualifying subsets are not accused as total count divergences](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-X03`

- [A mandatory document that does not exist fails](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B01`

- [A document that exists and does not mention the unit fails](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B02`

- [A mention by the identity code counts as documented](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B03`

- [A mention by the file name also counts](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B04`

- [Satisfying one of two duties is not enough](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B05`

- [Without a declaration nothing is charged](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B06`

- [A layer with no trigger is not charged](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B07`

- [Aggregated, the verdict is one per document](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B08`

- [The duty starts from the spec, not from the code](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-I01`

- [The layer used is the UNIT's, not the node's](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-I02`

- [Without a map the aggregated verdict is skipped](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-I03`

- [The gate does not understand the document's content](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-X01`

- [The gate does not decide which documents are mandatory](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-X02`

- [A reference that brings the passage it announces passes](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own) `DSCDC-B01`

- [A reference that only points is accused](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own) `DSCDC-B02`

- [A quotation counts in any written tradition](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own) `DSCDC-B03`

- [A path inside a code fence is an example, not a reference](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own) `DSCDC-B04`

- [The spec citing its own path is identifying itself](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own) `DSCDC-B05`

- [A rule code is not a revision](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own) `DSCDC-B06`

- [Only the spec is charged](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own) `DSCDC-B07`

- [A revision cited with an explanation on the same line passes](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own) `DSCDC-B08`

- [With no map the confrontation is skipped](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own) `DSCDC-B09`

- [The ruler matches structure, never vocabulary](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own) `DSCDC-I01`

- [The on-line explanation escape belongs to the revision, not to the path](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own) `DSCDC-I02`

- [The verdict names the line and shows what it says](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own) `DSCDC-I03`

- [The gate does not judge whether the accompanying content is faithful](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own) `DSCDC-X01`

- [The gate marks and does not block](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own) `DSCDC-X02`

- [Measuring explanation errs on the permissive side](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own) `DSCDC-X03`

- [An artifact that is not a spec leaves without a verdict](camadas/gate.md#dmdcd--domaindeclared--the-spec-declares-what-the-unit-accepts-and-who-blocks-the-invalid) `DMDCD-B01`

- [A spec without the domain section is failed](camadas/gate.md#dmdcd--domaindeclared--the-spec-declares-what-the-unit-accepts-and-who-blocks-the-invalid) `DMDCD-B02`

- [A section opened and left empty is failed](camadas/gate.md#dmdcd--domaindeclared--the-spec-declares-what-the-unit-accepts-and-who-blocks-the-invalid) `DMDCD-B03`

- [A row filled only with placeholders is not a declaration](camadas/gate.md#dmdcd--domaindeclared--the-spec-declares-what-the-unit-accepts-and-who-blocks-the-invalid) `DMDCD-B06`

- [An entry with no owner is failed, and the verdict names it](camadas/gate.md#dmdcd--domaindeclared--the-spec-declares-what-the-unit-accepts-and-who-blocks-the-invalid) `DMDCD-B04`

- [An entry whose owner is named passes](camadas/gate.md#dmdcd--domaindeclared--the-spec-declares-what-the-unit-accepts-and-who-blocks-the-invalid) `DMDCD-B07`

- [A waiver with a written reason silences the gate](camadas/gate.md#dmdcd--domaindeclared--the-spec-declares-what-the-unit-accepts-and-who-blocks-the-invalid) `DMDCD-B05`

- [A bare waiver, with no reason, does not waive](camadas/gate.md#dmdcd--domaindeclared--the-spec-declares-what-the-unit-accepts-and-who-blocks-the-invalid) `DMDCD-I01`

- [A non-answer in the owner column is not an owner](camadas/gate.md#dmdcd--domaindeclared--the-spec-declares-what-the-unit-accepts-and-who-blocks-the-invalid) `DMDCD-I02`

- [The gate does not judge whether the declared domain is correct](camadas/gate.md#dmdcd--domaindeclared--the-spec-declares-what-the-unit-accepts-and-who-blocks-the-invalid) `DMDCD-X01`

- [The gate does not read the code to check the validation exists](camadas/gate.md#dmdcd--domaindeclared--the-spec-declares-what-the-unit-accepts-and-who-blocks-the-invalid) `DMDCD-X02`

- [An artifact that is not a test leaves without a verdict](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-B01`

- [Without a built map the gate stays quiet](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-B02`

- [A test with no execution stamp is skipped, not failed](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-B03`

- [A test whose closure is intact passes](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-B04`

- [The passing verdict says what it checked against](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-B05`

- [A test whose dependency advanced a revision fails](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-B06`

- [The failing verdict names the culprit](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-B07`

- [The failing verdict states the fix](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-B08`

- [A test whose own file changed is reported separately from its closure](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-B09`

- [The culprit list is truncated at five and the remainder counted](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-B10`

- [Absence of proof and expired proof are never the same finding](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-I01`

- [A test that never ran is never approved either](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-I02`

- [Truncation never hides the size of the problem](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-I03`

- [The gate does not charge the absence of a green test](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-X01`

- [The gate does not run the test nor judge whether the change broke it](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-X02`

- [The gate does not read the project's configuration](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code) `EVFRV-X03`

- [Non-feature artifacts skip confrontation](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B01`

- [A nil graph returns pending without approving](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B02`

- [A feature declaring no scenarios skips confrontation](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B03`

- [A feature with no linked tests returns pending](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B04`

- [Scenarios belonging to non-test surfaces are skipped](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B05`

- [A scenario code completely absent from tests fails](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B06`

- [A scenario code appearing only in comments fails](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B07`

- [Tests implementing scenario codes with exact titles pass](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B08`

- [Tests with matching codes but drifting descriptions issue a warning](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B09`

- [Test titles containing quotes are parsed without truncation](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B10`

- [Sibling scenario codes in composite test titles are extracted](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B11`

- [Shared test titles verify scenario presence through test body](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B12`

- [Test comments contribute to descriptive match but not code presence](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B13`

- [Scenario codes match with exact word boundaries](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B14`

- [Go t.Run declarations are recognized as valid test titles](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B15`

- [Script comment markers are stripped when verifying code presence](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B16`

- [RootCode returns the root requirement code without scenario sub-index](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B17`

- [Missing scenario code is always a failure](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-I01`

- [Descriptive divergence is always an informative warning](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-I02`

- [Code presence ignores comments while description matching reads them](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-I03`

- [Static analysis does not run tests or inspect execution results](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-X01`

- [Non-unit surfaces are left to their respective gates](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-X02`

- [Minor description drift does not block promotion](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-X03`

- [Confronting an artifact that is not a spec skips](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B01`

- [Confronting without a graph returns Pending](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B02`

- [Confronting a spec without a declared code skips](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B03`

- [A testID prefix matching the spec identity code passes](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B04`

- [An orphan testID prefix with code shape fails](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B05`

- [A testID prefix matching another declared unit is accepted as legitimate reuse](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B06`

- [A visual regression baseline with a divergent code fails](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B07`

- [Short testID prefixes of three letters or fewer pass](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B08`

- [Inconsistency failures cite conflicting acronyms and origin files](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B09`

- [When no orphan testID prefixes or baseline discrepancies exist the gate passes](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B10`

- [Without a map graph the relational gate never approves](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-I01`

- [Cross-unit reuse is never permitted for visual regression baselines](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-I02`

- [Absence of a spec code is never double-charged](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-I03`

- [Common shorthand prefixes of three letters or fewer are not scrutinized](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-X01`

- [The gate does not enforce code presence on specs](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-X02`

- [Components referencing parent screen codes in testIDs are not forbidden](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-X03`

- [An artifact that is not code leaves without a verdict](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-B01`

- [Content matching a forbidden pattern fails, naming line and reason](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-B02`

- [A rule scoped to a layer charges only that layer](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-B03`

- [A rule with no layer holds for all code](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-B04`

- [Severity warn records without failing, and the default is error](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-B05`

- [A waiver with a written reason on the line waives that line](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-B06`

- [The waiver also holds in the comment on the line above](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-B07`

- [A bare marker with no reason does not waive](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-B08`

- [With no boundary declared the verdict is Pending, never Pass](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-B09`

- [An invalid forbid pattern fails visibly](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-B10`

- [The pattern is matched against the whole file, catching a multi-line import](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-B11`

- [The same rule is expressible in six language dialects](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-I01`

- [A single-line import of the same shape is still caught](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-I02`

- [The waiver holds on any line of the matched stretch](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-I03`

- [A line anchor keeps holding per line](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-I04`

- [The gate does not decide which boundaries exist](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-X01`

- [The gate does not parse the language, it matches text](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-X02`

- [The gate does not judge whether the boundary is the right one to draw](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own) `LYBNL-X03`

- [A rule marked at both declared scopes passes](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it) `MRPRM-B01`

- [A rule missing from one end fails, and the verdict names the empty scope](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it) `MRPRM-B02`

- [Two markings on the same side do not satisfy the gate](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it) `MRPRM-B03`

- [A rule left over at one end fails and is named](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it) `MRPRM-B04`

- [Total absence of the prefix is not approval](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it) `MRPRM-B05`

- [A declaration with no prefix returns Pending](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it) `MRPRM-B06`

- [With no scopes declared the ruler is the count](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it) `MRPRM-B07`

- [The marking crosses language](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it) `MRPRM-B08`

- [Ignored directories never count towards parity](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it) `MRPRM-B09`

- [The failing verdict names the rule and the empty scope](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it) `MRPRM-I01`

- [What was not measured is never approved](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it) `MRPRM-I02`

- [The gate does not read what each end actually does](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it) `MRPRM-X01`

- [The gate does not decide which rules live at two ends](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it) `MRPRM-X02`

- [Files outside the text extension list are not read](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it) `MRPRM-X03`

- [An artifact that is not a test leaves without a verdict](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B01`

- [A stamp that matches the module today passes](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B02`

- [A snippet that changed since the stamp was written fails](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B03`

- [The stamp is immune to displacement](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B04`

- [An anchor that vanished fails with its own message](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B05`

- [An anchor occurring more than once is ambiguous and fails](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B06`

- [The declared line count delimits the window](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B07`

- [Without the dialect declared the gate goes quiet](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B08`

- [A double of a module the project does not govern is not charged](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B09`

- [A stamp whose module no longer exists is a finding, not a crash](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B10`

- [The absence of a stamp on a governed double is accused](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B11`

- [A stamp present and correct satisfies both charges](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B12`

- [A dialect regex that does not compile fails loudly](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B13`

- [A dialect regex with no capture group fails](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B14`

- [Another ecosystem's dialect is charged the same way](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B15`

- [The gate recomputes the hash instead of validating the stamp's format](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-I01`

- [The double and the stamp are tied by the module's path](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-I02`

- [A configuration fault fails and a project decision skips](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-I03`

- [The gate does not interpret the code of the stamped module](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-X01`

- [The gate carries no built-in dialect for detecting doubles](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-X02`

- [The gate does not skip the absence of a stamp to accommodate legacy code](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-X03`

- [The gate does not guarantee cryptographic strength](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-X04`

- [A double with no tie fails and the verdict names the loose module](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-B01`

- [A double whose factory carries the declared tie passes](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-B02`

- [The charge is per module, not per file](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-B03`

- [A double with no factory is not charged](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-B04`

- [Without the tie shape declared the gate goes quiet](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-B05`

- [An artifact that is not a test leaves without a verdict](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-B06`

- [A test that doubles nobody leaves without a verdict](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-B07`

- [Another runner of the same ecosystem is recognised the same way](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-B08`

- [A third-party library double is not charged](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-B09`

- [The third-party exemption is not an escape hatch](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-B10`

- [A relative import that resolves in the map is an own module](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-B11`

- [Another ecosystem's dialect is charged the same way](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-B12`

- [What is governed is decided by the graph, never by a prefix list](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-I01`

- [An undeclared tie shape skips instead of guessing one](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-I02`

- [The verdict counts the loose doubles, not just names them](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-I03`

- [The gate does not check whether the annotated type matches the real module](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-X01`

- [The gate does not cover drift of behaviour](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-X02`

- [The gate carries no built-in tie shape and no built-in ecosystem](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-X03`

- [The gate does not charge third-party library doubles](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-X04`

- [A node that carries the trigger and is absent from the demanded file fails](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B01`

- [A node that carries the trigger and does appear passes](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B02`

- [A node without the trigger contracts no obligation](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B03`

- [A waiver exempts only when it carries a written reason](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B04`

- [A project with no declared obligation is skipped](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B05`

- [An acknowledged debt with a written when yields Pending](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B06`

- [A bare debt marker keeps failing](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B07`

- [Waiver and debt stay distinct](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B08`

- [The failing verdict offers the three ways out](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B09`

- [The token is derived through the declared form](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-I01`

- [A glob that matches no file produces no violation](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-I02`

- [The node's own identified_as wins over the automatic form](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-I03`

- [The gate does not decide which obligations exist](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-X01`

- [The gate does not understand what the destination does with the token](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-X02`

- [A declaration written in the body is not read](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-X03`

- [An artifact that is not a spec leaves without a verdict](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement) `OPQSP-B01`

- [Whoever OPENED the section is confronted by its content](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement) `OPQSP-B02`

- [An open item bars the spec](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement) `OPQSP-B03`

- [A section closed honestly releases the spec](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement) `OPQSP-B04`

- [An item marked as resolved does not block](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement) `OPQSP-B05`

- [A question with no code is charged](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement) `OPQSP-B06`

- [The count of pending decisions reads the project's own lexicon](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement) `OPQSP-B07`

- [Prose is not an item](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement) `OPQSP-I01`

- [The section boundary is respected](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement) `OPQSP-I02`

- [Filling in what the question BECOMES does not close the question](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement) `OPQSP-I03`

- [The gate does not judge whether the question is good](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement) `OPQSP-X01`

- [A spec with no section is a pending item, and the verdict teaches the way out](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement) `OPQSP-X02`

- [A limit received from the caller passes](camadas/gate.md#pgnhn--paginationhonored--what-promises-a-set-does-not-return-the-first-page-in-silence) `PGNHN-B01`

- [A limit hidden in a default value is accused](camadas/gate.md#pgnhn--paginationhonored--what-promises-a-set-does-not-return-the-first-page-in-silence) `PGNHN-B02`

- [The NAME bounds the promise](camadas/gate.md#pgnhn--paginationhonored--what-promises-a-set-does-not-return-the-first-page-in-silence) `PGNHN-B03`

- [Sibling functions that paginate are the proof by asymmetry](camadas/gate.md#pgnhn--paginationhonored--what-promises-a-set-does-not-return-the-first-page-in-silence) `PGNHN-B04`

- [A waiver with a written reason leaves the report](camadas/gate.md#pgnhn--paginationhonored--what-promises-a-set-does-not-return-the-first-page-in-silence) `PGNHN-B05`

- [The verdict offers the way out](camadas/gate.md#pgnhn--paginationhonored--what-promises-a-set-does-not-return-the-first-page-in-silence) `PGNHN-B06`

- [Without a declared dialect the verdict is undetermined](camadas/gate.md#pgnhn--paginationhonored--what-promises-a-set-does-not-return-the-first-page-in-silence) `PGNHN-B07`

- [The ruler is agnostic across languages](camadas/gate.md#pgnhn--paginationhonored--what-promises-a-set-does-not-return-the-first-page-in-silence) `PGNHN-I01`

- [Where the construct is not recognised, the gate stays silent](camadas/gate.md#pgnhn--paginationhonored--what-promises-a-set-does-not-return-the-first-page-in-silence) `PGNHN-I02`

- [A cursor with no loop does not count as pagination](camadas/gate.md#pgnhn--paginationhonored--what-promises-a-set-does-not-return-the-first-page-in-silence) `PGNHN-I03`

- [A provider prefix in the name does not hide the promise](camadas/gate.md#pgnhn--paginationhonored--what-promises-a-set-does-not-return-the-first-page-in-silence) `PGNHN-I04`

- [The gate does not invent a cursor the provider does not offer](camadas/gate.md#pgnhn--paginationhonored--what-promises-a-set-does-not-return-the-first-page-in-silence) `PGNHN-X01`

- [The gate does not measure performance or page size](camadas/gate.md#pgnhn--paginationhonored--what-promises-a-set-does-not-return-the-first-page-in-silence) `PGNHN-X02`

- [An artifact reached only by the impact radius is skipped](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B01`

- [An artifact with no identity code is skipped](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B02`

- [A changed plan or spec with no declared revision fails](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B03`

- [A changed plan or spec with a valid revision passes](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B04`

- [Citing the revision of another document does not count](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B05`

- [Revisions with non-sequential numbering fail](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B06`

- [Multiple sequential revisions pass and cite the latest](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B07`

- [Revisions written in various markdown formats pass](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B08`

- [A change already explained by the plan revision mechanism passes](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B09`

- [A newly created untracked file is skipped](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B10`

- [A newly created staged file is skipped](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B11`

- [An existing committed file modified without a revision fails](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B12`

- [Outside a git repository the gate trusts the changed list](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B13`

- [RevisionsOf parses declared revisions in order](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B14`

- [Untouched nodes in the impact radius are never charged](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-I01`

- [Revisions must be sequential from one without gaps](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-I02`

- [Files not existing in HEAD are never charged](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-I03`

- [The gate does not judge whether a correction is directional](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-X01`

- [The gate does not enforce identity code presence](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-X02`

- [Without a changed files list the gate skips confrontation](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-X03`

- [A source that lives only in the prose is failed](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-B01`

- [The verdict names which source and where its adapter lives](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-B02`

- [With the owning plan declared in needs the gate passes](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-B03`

- [Every source of the line is confronted on its own](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-B04`

- [The source name matches the adapter regardless of case](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-B05`

- [A source whose adapter nobody seeds is not charged](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-B06`

- [The plan that seeds the adapter is not charged for itself](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-B07`

- [A plan with no source line returns Skip](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-B08`

- [An artifact that is not a plan returns Skip](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-B09`

- [What was not measured is never approved](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-I01`

- [A seeded file off the naming pattern owns nothing](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-I02`

- [The gate does not confront the order of the phases](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-X01`

- [The gate does not demand a needs pointing at nothing](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-X02`

- [The gate does not interpret what the source is for](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-X03`

- [An artifact that is not a plan leaves without a verdict](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-B01`

- [A plan with no companion progress file is skipped, not failed](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-B02`

- [The skip for a missing companion says how to create it](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-B03`

- [A ticked item whose file does not exist is failed](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-B04`

- [An open item whose file already exists is failed](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-B05`

- [A spec the plan seeds and the progress does not list is failed](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-B06`

- [A checkbox item promising no file at all is failed](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-B07`

- [The ticked-but-absent finding is reported first](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-B08`

- [An item in prose citing no path is not charged](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-B09`

- [A progress that agrees with the disk on every item passes](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-B10`

- [The verdict names each offending path](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-B11`

- [The companion's path has one definition, derived from the scanner](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-I01`

- [A seed is matched by path, never by the item's text](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-I02`

- [A spec mentioned in the plan's prose is not a seed](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-I03`

- [A template file is never a seeded spec](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-I04`

- [The gate does not charge the existence of the progress file](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-X01`

- [The gate does not judge the content of an item beyond the path](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-X02`

- [The gate does not put the progress file into the map](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-X03`

- [An artifact that is not a spec leaves without a verdict](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-B01`

- [A marked rule whose governed code does not import the cited unit fails](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-B02`

- [A citation living only in a comment does not satisfy the charge](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-B03`

- [Governed code that imports the cited unit passes](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-B04`

- [An import through a different alias still matches](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-B05`

- [A declared waiver with a written reason is not charged](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-B06`

- [The owner stamp does not demand an import](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-B07`

- [A relation claimed in prose, with no mark, is reported as a suspicion](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-B08`

- [The mark with no target on the line is reported](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-B09`

- [A target declared by rule code is resolved through the map](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-B10`

- [A rule code that resolves to nothing is reported as unresolved](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-B11`

- [A rule code of the unit itself is not a target](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-B12`

- [A rule with no relation claim leaves without a verdict](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-B13`

- [Imports in other language shapes satisfy the charge](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-B14`

- [The claim and the target must be on the same rule line](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-I01`

- [The import is charged on the governed code, never on the test](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-I02`

- [A satisfied charge does not swallow a pending suspicion](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-I03`

- [The gate does not hunt duplicated concepts across the project](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-X01`

- [The gate does not compare the values on the two sides](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-X02`

- [The gate does not decide whether an unmarked prose claim blocks](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-X03`

- [Confronting an artifact that is neither code nor test skips](camadas/gate.md#rphrg--regionpairhonored--every-opened-source-region-must-close-with-its-own-identity-code) `RPHRG-B01`

- [Code or test artifacts without region markers skip](camadas/gate.md#rphrg--regionpairhonored--every-opened-source-region-must-close-with-its-own-identity-code) `RPHRG-B02`

- [Balanced regions closing with matching identity codes pass](camadas/gate.md#rphrg--regionpairhonored--every-opened-source-region-must-close-with-its-own-identity-code) `RPHRG-B03`

- [An opened region that is never closed fails](camadas/gate.md#rphrg--regionpairhonored--every-opened-source-region-must-close-with-its-own-identity-code) `RPHRG-B04`

- [An end region marker without an opening marker fails as an orphan close](camadas/gate.md#rphrg--regionpairhonored--every-opened-source-region-must-close-with-its-own-identity-code) `RPHRG-B05`

- [An end region marker closing with a different code fails](camadas/gate.md#rphrg--regionpairhonored--every-opened-source-region-must-close-with-its-own-identity-code) `RPHRG-B06`

- [Multiple pairing errors are ordered sequentially by line number](camadas/gate.md#rphrg--regionpairhonored--every-opened-source-region-must-close-with-its-own-identity-code) `RPHRG-B07`

- [Region absence is never charged as a failure](camadas/gate.md#rphrg--regionpairhonored--every-opened-source-region-must-close-with-its-own-identity-code) `RPHRG-I01`

- [Pairing defects produce a blocking Fail verdict](camadas/gate.md#rphrg--regionpairhonored--every-opened-source-region-must-close-with-its-own-identity-code) `RPHRG-I02`

- [End markers must explicitly match opening codes](camadas/gate.md#rphrg--regionpairhonored--every-opened-source-region-must-close-with-its-own-identity-code) `RPHRG-I03`

- [The gate does not mandate region markers across all source files](camadas/gate.md#rphrg--regionpairhonored--every-opened-source-region-must-close-with-its-own-identity-code) `RPHRG-X01`

- [Region markers inside specifications or documentation are ignored](camadas/gate.md#rphrg--regionpairhonored--every-opened-source-region-must-close-with-its-own-identity-code) `RPHRG-X02`

- [Internal code semantics inside regions are not evaluated](camadas/gate.md#rphrg--regionpairhonored--every-opened-source-region-must-close-with-its-own-identity-code) `RPHRG-X03`

- [A spec whose rules the code ignores is accused, and the verdict names them](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-B01`

- [A rule waived with a written reason closes the account](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-B02`

- [A unit that predates the practice is a pending item, not a failure](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-B03`

- [Declaring the requirement turns the pending item into a failure](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-B04`

- [A spec with no linked code is not this gate's subject](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-B05`

- [A waiver with no named rule covers every rule of the spec](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-B06`

- [Requiring the marking never punishes whoever already marks](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-I01`

- [The identity survives a rename](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-I02`

- [The gate does not judge whether the implementation is correct](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-X01`

- [The gate does not demand a mark on EVERY rule](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-X02`

- [A letter that is not declared in the vocabulary fails](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B01`

- [A declared letter under a claimed section passes](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B02`

- [A section cataloguing rules under a title no letter claims fails](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B03`

- [The same letter claimed by two terms is a conflict in the vocabulary](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B04`

- [With no vocabulary declared the gate confronts the canonical letters](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B05`

- [A heading that is the rule code itself is not a category section](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B06`

- [A section that only cites other sections' codes claims no letter](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B07`

- [A section that defines a code in the first table cell is charged](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B08`

- [A section declared as rule-cataloguing and filled without a code is Pending](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B09`

- [A section whose table already carries the code is not charged](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B10`

- [A declared section outside sections_require_code is not charged](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B11`

- [A project that does not use sections_require_code changes no behaviour](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B12`

- [A spec with no rule code at all is not this gate's problem](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-I01`

- [The verdict names the letter and where to declare it](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-I02`

- [A filled section with no code is the gap where the scenario loses its anchor](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-I03`

- [The gate does not decide which letters exist](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-X01`

- [Without a declared vocabulary only the letter is charged](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-X02`

- [The gate does not judge whether the letter suits the rule](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-X03`

- [The gate charges traceability, not format](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-X04`

- [A requirement no scenario tags is failed and named](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B01`

- [A requirement that has a scenario is not accused](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B02`

- [With every requirement tagged the gate passes](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B03`

- [A code merely cited contracts no obligation](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B04`

- [A per-requirement waiver with a written reason waives](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B05`

- [A bare per-requirement waiver does not waive](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B06`

- [A whole-spec waiver drags the waiver to every requirement](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B07`

- [Without the whole-spec waiver the same requirements keep failing](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B08`

- [A bare whole-spec waiver drags nothing](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B09`

- [A spec with no feature returns Skip](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B10`

- [Requirements are looked for across every linked feature](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B11`

- [An artifact that is not a spec returns Skip](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B12`

- [Every waiver requires a written reason](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-I01`

- [Each gate accuses one thing](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-I02`

- [The gate does not judge whether the scenario proves the requirement](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-X01`

- [A cited code produces no accusation](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-X02`

- [The gate does not confront feature against test](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-X03`

- [The gate skips when no test handle attribute is declared](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-B01`

- [The gate skips when no E2E surface is configured](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-B02`

- [The gate skips when no exposed handles exist in code](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-B03`

- [A queried handle that exists in code passes](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-B04`

- [A queried handle missing from code fails](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-B05`

- [The verdict names the missing handle and the querying flows](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-B06`

- [A negative assertion on a missing handle fails](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-B07`

- [An exposed template handle covers a queried instance](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-B08`

- [A regex pattern in a flow matches an exposed prefix head](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-B09`

- [A flow handle with runtime interpolation does not trigger a failure](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-B10`

- [A marked handle defined in a lookup table counts as exposed](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-B11`

- [A template behind a nullish coalescing fallback counts as exposed](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-B12`

- [A suffix composed in a child component counts as exposed](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-B13`

- [A terminal suffix composed from a prop counts as exposed](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-B14`

- [A handle appearing only in test files does not count as exposed](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-B15`

- [Missing configuration never reports a pass](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-I01`

- [Findings are grouped by handle identifier](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-I02`

- [Negative assertions are checked against code presence](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-I03`

- [Handles exposed in code but unused in flows are not accused](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-X01`

- [Flow expressions requiring runtime execution are not evaluated](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-X02`

- [Default test handle attributes are not inferred](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-X03`

- [An artifact that is not a spec leaves without a verdict](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist) `TRCMT-B01`

- [A recognised layer leaves without a verdict](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist) `TRCMT-B02`

- [Without a map the verdict is undetermined](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist) `TRCMT-B03`

- [A spec with the three pieces linked passes](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist) `TRCMT-B04`

- [A spec missing a piece is failed, and the verdict names which](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist) `TRCMT-B05`

- [The layer may waive a piece for every spec in it](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist) `TRCMT-B06`

- [The unit may waive a piece in its own spec, with a written reason](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist) `TRCMT-B07`

- [The test is reached in two hops, through the feature](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist) `TRCMT-I01`

- [Waiving the test while the feature carries a scenario is a contradiction](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist) `TRCMT-I02`

- [Waiving the test demands saying where the proof is, and the place must exist](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist) `TRCMT-I03`

- [A waiver covers only the piece it declares](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist) `TRCMT-I04`

- [The gate does not confront whether the pieces MATCH one another](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist) `TRCMT-X01`

- [The gate does not judge the QUALITY of any piece](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist) `TRCMT-X02`

- [An artifact that is not a spec leaves without a verdict](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-B01`

- [A value of a closed set with no anchor is failed](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-B02`

- [The verdict names the unanchored value](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-B03`

- [A value whose anchor carries rule key and value passes](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-B04`

- [An anchor that asserts one value while the line says another is failed](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-B05`

- [The verdict of a lying anchor shows both sides of the divergence](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-B06`

- [Lying anchors are reported before the unanchored ones](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-B07`

- [Without a declared value anchor pattern the gate skips](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-B08`

- [The skip names the setting that enables the gate](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-B09`

- [A declaration that opens no list is not a closed set](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-B10`

- [An anchor pattern with a single capture group does not enable the gate](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-I01`

- [A line carrying an anchor is never read as the end of the list](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-I02`

- [With no built map the verdict is pending](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-I03`

- [The gate does not judge whether the value is a good one](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-X01`

- [The gate does not accuse a line whose literal it cannot read](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-X02`

- [The gate does not decide what an anchor or a public symbol looks like](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value) `VLANV-X03`

