<!-- anchors:generated from doct/comportamento.md.tmpl — inputs:896acc52f7490751 — DO NOT EDIT: run `anchors docs build` -->


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

- [Non-specification artifacts skip confrontation](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-B01`

- [Confronting without a map graph returns pending](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-B02`

- [Confronting with an empty identity universe returns pending](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-B03`

- [Citations resolving to existing units in the map pass](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-B04`

- [Citations pointing to non-existent units fail](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-B05`

- [Citations matching the specification's own identity code pass](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-B06`

- [Multiple orphaned requirement citations are reported sorted](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-B07`

- [Identity ownership is resolved from graph nodes and header metadata](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-B08`

- [Non-requirement tokens are ignored](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-B09`

- [Self-references to a specification's own requirements never fail](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-I01`

- [Missing map or empty identity universe returns pending rather than pass](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-I02`

- [Unresolvable external citations always produce a blocking fail verdict](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-I03`

- [Implementation correctness of referenced requirements is not evaluated](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-X01`

- [Non-specification artifacts are not inspected by this gate](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-X02`

- [External requirement citations are not mandatory](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-X03`

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

- [An artifact that is not a spec leaves without a verdict](camadas/gate.md#dephn--dependencyhonored--methods-promised-in-the-dependency-table-are-consumed-in-code) `DEPHN-B01`

- [Without a relational map the verdict is undetermined](camadas/gate.md#dephn--dependencyhonored--methods-promised-in-the-dependency-table-are-consumed-in-code) `DEPHN-B02`

- [A spec declaring no confrontable symbols leaves without a verdict](camadas/gate.md#dephn--dependencyhonored--methods-promised-in-the-dependency-table-are-consumed-in-code) `DEPHN-B03`

- [A spec governing no code leaves the verdict undetermined](camadas/gate.md#dephn--dependencyhonored--methods-promised-in-the-dependency-table-are-consumed-in-code) `DEPHN-B04`

- [When every promised symbol appears in governed code, the gate passes](camadas/gate.md#dephn--dependencyhonored--methods-promised-in-the-dependency-table-are-consumed-in-code) `DEPHN-B05`

- [When a promised symbol is absent from governed code, the gate fails](camadas/gate.md#dephn--dependencyhonored--methods-promised-in-the-dependency-table-are-consumed-in-code) `DEPHN-B06`

- [When an absent symbol resembles an identifier in code, the verdict suggests the rename](camadas/gate.md#dephn--dependencyhonored--methods-promised-in-the-dependency-table-are-consumed-in-code) `DEPHN-B07`

- [Symbols appearing only in comments do not fulfill the promise](camadas/gate.md#dephn--dependencyhonored--methods-promised-in-the-dependency-table-are-consumed-in-code) `DEPHN-B08`

- [Prose descriptions in dependency methods are never treated as contracts](camadas/gate.md#dephn--dependencyhonored--methods-promised-in-the-dependency-table-are-consumed-in-code) `DEPHN-I01`

- [Symbol presence is matched strictly on token word boundaries](camadas/gate.md#dephn--dependencyhonored--methods-promised-in-the-dependency-table-are-consumed-in-code) `DEPHN-I02`

- [Near-symbol rename suggestions are strictly conservative](camadas/gate.md#dephn--dependencyhonored--methods-promised-in-the-dependency-table-are-consumed-in-code) `DEPHN-I03`

- [The gate performs static textual confrontation without runtime execution](camadas/gate.md#dephn--dependencyhonored--methods-promised-in-the-dependency-table-are-consumed-in-code) `DEPHN-X01`

- [The gate does not interpret dependency semantics or parameter signatures](camadas/gate.md#dephn--dependencyhonored--methods-promised-in-the-dependency-table-are-consumed-in-code) `DEPHN-X02`

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

- [A compiled document that no longer matches the spec fails](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-B01`

- [The verdict names the stale documents](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-B02`

- [A compiled document that matches the templates passes](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-B03`

- [An artifact that is not a spec leaves without a verdict](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-B04`

- [A project with no template directory is not charged](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-B05`

- [A template whose document was never produced counts as stale](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-B06`

- [A document written by hand is not charged](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-B07`

- [A template that cannot be compiled fails](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-B08`

- [A project whose specs cannot be read leaves without a verdict](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-B09`

- [The comparison happens in memory](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-B10`

- [The gate never repairs what it points at](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-I01`

- [The charge starts from the spec and never from the compiled document](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-I02`

- [The gate does not judge whether the compiled document is good](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-X01`

- [The gate does not run the build even knowing the fix](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-X02`

- [The gate does not charge documents written by hand](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-X03`

- [The compiled document is not confronted as an artifact of its own](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec) `DCFRD-X04`

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

- [An external command exiting with status zero returns Pass](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-B01`

- [An external command exiting with non-zero status returns Fail with output](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-B02`

- [An external command failing with empty output reports execution error](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-B03`

- [Single node execution delegates to RunExternalArgs](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-B04`

- [Placeholder file is rewritten to positional parameter](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-B05`

- [Placeholder files is rewritten to all positional parameters](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-B06`

- [Execution without targets runs once for project scope](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-B07`

- [Targets within budget run in a single batch](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-B08`

- [Targets exceeding budget are partitioned across batches](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-B09`

- [Failure in any batch causes entire execution to fail](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-B10`

- [Single target failure detail is truncated at five hundred characters](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-B11`

- [Batch failure detail is truncated at four thousand characters](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-B12`

- [Environment variable overrides argv limit](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-B13`

- [Target paths are passed strictly in argv preventing injection](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-I01`

- [Oversized single target is isolated in its own batch](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-I02`

- [Platform default argv limits are enforced](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-I03`

- [The gate does not parse or interpret linter diagnostics](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-X01`

- [The gate does not aggregate cross-file state across partitioned batches](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-X02`

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

- [A gate reaches only the kinds it declares](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B01`

- [A gate that names labels reaches only the nodes carrying one](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B02`

- [One excluded label is enough to keep a node out](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B03`

- [Exclusion wins over the positive label filter](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B04`

- [A gate demanding a mark reaches only the targets that carry it](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B05`

- [An unreadable target does not apply](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B06`

- [A gate with no applicable target does not run at all](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B07`

- [A gate whose required binary is absent steps aside](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B08`

- [A waiver by target spares one node and confronts the rest](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B09`

- [An aggregate gate runs once and reports against the scope](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B10`

- [A judgment gate asks for judgment when nothing has answered it](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B11`

- [A judgment gate reads the stamp an earlier judgement left](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B12`

- [A judgement recorded as waived becomes Skip and never Pass](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B13`

- [A gate declaring neither a command nor a check is undetermined](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B14`

- [Only the pending item that says a decision is still to take bars promotion](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B15`

- [Only the obligations gate produces assumed debt](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B16`

- [The engine reconfigures the code grammar before running anything](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B17`

- [The plain entry point runs with the map and no Structure](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B18`

- [The entry point that carries the Structure hands it to the checkers](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B19`

- [The entry point that knows the sweep kind honours the full-sweep scope](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B20`

- [The entry point that honours a waiver by target keeps the gate running](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B21`

- [A vendored file is out of every internal ruler and still reached by an external command](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B22`

- [A waiver by target never removes the gate from the list](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-I01`

- [Stepping aside, not measuring and failing are three different answers](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-I02`

- [A waived target leaves the failure tally without leaving the report](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-I03`

- [The reported target of an aggregate gate is the scope](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-I04`

- [The engine does not decide whether a target is correct](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-X01`

- [The engine does not compute the verdict of a judgment gate](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-X02`

- [The engine invents neither a map nor a Structure nor a waiver](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-X03`

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

- [A declared name routes to the function registered under it](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B01`

- [A name that does not resolve answers undetermined](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B02`

- [The routing tries the relational registry before the simpler ones](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B03`

- [The per-node path reads the target file, and a failed read is a failure](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B04`

- [The aggregate path reads no file at all](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B05`

- [An unresolved name in the aggregate path is undetermined too](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B06`

- [An aggregate checker that needs the declaring gate receives it](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B07`

- [Setting the rule letters reconfigures every dependent pattern together](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B08`

- [A file that is empty or only whitespace fails the emptiness ruler](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B09`

- [The identity ruler charges the presence of a scenario code](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B10`

- [A governed file with no identity block fails the header ruler](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B11`

- [A governed file passes with ownership or with reference, never with layer alone](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B12`

- [A file of a recognised layer passes with the layer alone](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B13`

- [A binary file steps aside from the header ruler](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B14`

- [An executable test script steps aside by a different path](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B15`

- [A guide without compliance points fails](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B16`

- [An unresolved name never approves, on either path](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-I01`

- [Every registered name is reachable through exactly one routing path](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-I02`

- [The compliance ruler is recognised in every language of the catalogue](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-I03`

- [The registry does not decide which checks a project runs](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-X01`

- [The registry does not invoke external tooling](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-X02`

- [The registry does not judge whether the text is good](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-X03`

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

- [A generated stamp passes the gate, and a change to its snippet fails it](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one) `MKSTP-B01`

- [An ambiguous specifier resolves to the test's workspace, or is skipped](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one) `MKSTP-B02`

- [A stamp per factory key covers only that export](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one) `MKSTP-B03`

- [An automock gets one stamp over the whole module](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one) `MKSTP-B04`

- [A third-party double is not stamped](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one) `MKSTP-B05`

- [The header stays out of a whole-module stamp](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one) `MKSTP-B06`

- [An existing stamp is never rewritten](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one) `MKSTP-I01`

- [A repeated line never anchors a stamp](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one) `MKSTP-I02`

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

- [A module with a dot in its name is matched to its stamp](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B16`

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

- [The PlanPhases function extracts catalogued phase codes](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-B01`

- [Non-plan artifacts skip phase ordering confrontation](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-B02`

- [Plans without phase headings skip confrontation](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-B03`

- [Plans with phase-like sections lacking codes return Pending](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-B04`

- [Plans declaring valid backward phase dependencies pass](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-B05`

- [Plans with duplicate phase codes fail](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-B06`

- [Phases depending on uncatalogued phase codes fail](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-B07`

- [Phases depending on themselves fail](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-B08`

- [Phases depending on future phases fail](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-B09`

- [Specifications declaring existing phase dependencies pass](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-B10`

- [Specifications declaring missing phase dependencies fail](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-B11`

- [Artifacts declaring valid parents pass](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-B12`

- [Artifacts declaring invalid parents, self-parenting, or parent cycles fail](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-B13`

- [Phase dependencies are strictly acyclic and backward-directed](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-I01`

- [Phase and parent targets must exist in the map](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-I02`

- [Parent chains are cycle-free and bounded](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-I03`

- [Phase detection identifies level-three sections regardless of language](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-I04`

- [Small plans are not required to catalog phases](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-X01`

- [Phase duration and calendar timing are not verified](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-X02`

- [Both artifact codes and phase codes are accepted as parents](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-X03`

- [A raw skeleton fails confrontation](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-B01`

- [A header field whose value is the marker fails](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-B02`

- [A table cell holding only the marker fails](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-B03`

- [A title or body line opening with the marker fails](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-B04`

- [The verdict names what was left behind](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-B05`

- [An artifact with every marker replaced passes](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-B06`

- [A section written on purpose to list pending work is not accused](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-I01`

- [The gate does not judge the quality of replacement text](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-X01`

- [A marker in running prose is not accused](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-X02`

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

- [Non-plan artifacts skip confrontation](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B01`

- [Confronting without a map graph returns pending](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B02`

- [A plan with neither revisions nor revisers skips confrontation](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B03`

- [Declaring a revision target that does not exist in the map fails](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B04`

- [A revising plan receives a pending reminder when the target lacks a top notice](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B05`

- [The pending reminder on a revising plan clears once the target carries the notice](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B06`

- [A revised plan lacking a top revision notice fails](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B07`

- [A revised plan placing the revision notice after line 40 fails](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B08`

- [A revised plan with top notice but no section amendment markers returns pending](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B09`

- [A revised plan with top notice and marked section amendments passes](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B10`

- [Markdown alerts and metadata directives are both accepted as valid markers](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B11`

- [Revision notices must be placed within the first 40 lines](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-I01`

- [Missing section amendment markers yield pending rather than failure](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-I02`

- [The pending reminder clears once the revised plan is notified](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-I03`

- [Language neutrality allows markdown alerts and metadata directives](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-X01`

- [Prose description quality accompanying revision markers is not evaluated](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-X02`

- [Section amendment markers are not demanded on unrevised plans](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-X03`

- [Non-plan artifacts skip confrontation](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-B01`

- [Missing project configuration returns Pending](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-B02`

- [Plans without seeded specifications skip confrontation](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-B03`

- [Template specification references are ignored as templates](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-B04`

- [Bare specification file names without directory paths are ignored](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-B05`

- [Informal path abbreviations without real top-level directories are ignored](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-B06`

- [Seeded specifications targeting valid governed layers pass](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-B07`

- [Multiple target source file extensions resolve the governed layer](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-B08`

- [Seeded specifications targeting declarative layers fail](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-B09`

- [Seeded specifications matching no declared layer in a real directory fail](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-B10`

- [Multiple seed defects across declarative and undeclared layers are aggregated](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-B11`

- [Plan seed validation applies exclusively to plan artifacts](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-I01`

- [Declarative layers reject specification seeding](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-I02`

- [Seed validation evaluates structural layer validity without requiring file existence](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-I03`

- [Casual prose citations and templates are not treated as seeded paths](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-I04`

- [Existing file presence is not required for seeded specifications](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-X01`

- [Plan progress synchronization is not evaluated by this gate](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-X02`

- [Specification content and scenarios within seeded files are not verified](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers) `PSVPL-X03`

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

- [Returns an empty list of Promotable gates when profile contains no gates](camadas/gate.md#prgtp--promotablegates--identifies-clean-informative-gates-ready-for-promotion-to-blocking) `PRGTP-B01`

- [An informative gate with passes and zero failures is included](camadas/gate.md#prgtp--promotablegates--identifies-clean-informative-gates-ready-for-promotion-to-blocking) `PRGTP-B02`

- [An informative gate with failures is excluded from promotion](camadas/gate.md#prgtp--promotablegates--identifies-clean-informative-gates-ready-for-promotion-to-blocking) `PRGTP-B03`

- [An informative gate with zero passes is excluded as having no data](camadas/gate.md#prgtp--promotablegates--identifies-clean-informative-gates-ready-for-promotion-to-blocking) `PRGTP-B04`

- [A blocking gate is excluded from promotion suggestions](camadas/gate.md#prgtp--promotablegates--identifies-clean-informative-gates-ready-for-promotion-to-blocking) `PRGTP-B05`

- [Returned Promotable populates Gate with the declared gate name](camadas/gate.md#prgtp--promotablegates--identifies-clean-informative-gates-ready-for-promotion-to-blocking) `PRGTP-B06`

- [Returned Promotable records Passou equal to passed node count](camadas/gate.md#prgtp--promotablegates--identifies-clean-informative-gates-ready-for-promotion-to-blocking) `PRGTP-B07`

- [PromotableGates returns candidates sorted deterministically by name](camadas/gate.md#prgtp--promotablegates--identifies-clean-informative-gates-ready-for-promotion-to-blocking) `PRGTP-B08`

- [Multiple clean informative gates are all collected](camadas/gate.md#prgtp--promotablegates--identifies-clean-informative-gates-ready-for-promotion-to-blocking) `PRGTP-B09`

- [An informative gate is candidate if and only if non-blocking zero failures and positive passes](camadas/gate.md#prgtp--promotablegates--identifies-clean-informative-gates-ready-for-promotion-to-blocking) `PRGTP-I01`

- [A gate with zero passes is never classified as clean](camadas/gate.md#prgtp--promotablegates--identifies-clean-informative-gates-ready-for-promotion-to-blocking) `PRGTP-I02`

- [Does not automatically promote gates or modify anchors yaml](camadas/gate.md#prgtp--promotablegates--identifies-clean-informative-gates-ready-for-promotion-to-blocking) `PRGTP-X01`

- [Does not evaluate gate execution results directly from disk](camadas/gate.md#prgtp--promotablegates--identifies-clean-informative-gates-ready-for-promotion-to-blocking) `PRGTP-X02`

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

- [A reference that does not match the sibling spec fails](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B01`

- [The failing verdict names both sides of the divergence](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B02`

- [A reference equal to the sibling spec passes](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B03`

- [A spec is not confronted](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B04`

- [An artifact with no reference declared leaves without a verdict](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B05`

- [With no sibling spec the gate goes quiet](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B06`

- [A sibling spec that declares no identity counts as no sibling](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B07`

- [The sibling is found by the name convention](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B08`

- [A test file lands on the same sibling its code does](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B09`

- [An intermediate extension is dropped from the stem](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B10`

- [The reference is read whatever the comment syntax of the language](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B11`

- [The accepted identity length comes from the project's Structure](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B12`

- [A leading dot is not a stem separator](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B13`

- [The ruler is the sibling on disk, never the map](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-I01`

- [The gate never repairs what it points at](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-I02`

- [The gate does not charge the absence of the reference field](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-X01`

- [The gate does not charge the absence of the sibling spec](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-X02`

- [The gate does not consult the map to resolve the reference](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-X03`

- [The gate does not judge whether the spec describes the unit well](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-X04`

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

- [Non-spec artifacts skip confrontation](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B01`

- [A spec with no revision has nothing to confront](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B02`

- [A revision naming no rule cannot be confronted](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B03`

- [A revision naming a rule the spec does not define fails](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B04`

- [A sibling sharing vocabulary and left unmentioned is reported](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B05`

- [Every vocabulary-sharing sibling accounted for passes](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B06`

- [Checked clears the accusation without asserting correctness](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B07`

- [A revised rule is never its own orphan](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-I01`

- [Terms shared by the whole unit do not discriminate](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-X01`

- [A revision the branch added whose number the base already uses moves to the next free number](camadas/gate.md#rvrnr--revisionrenumber--the-revisions-a-branch-added-move-to-a-free-number-when-the-base-took-theirs) `RVRNR-B01`

- [After a rebase, only the branch's revision moves, not the base's with the same number](camadas/gate.md#rvrnr--revisionrenumber--the-revisions-a-branch-added-move-to-a-free-number-when-the-base-took-theirs) `RVRNR-B02`

- [When one added revision collides, every added revision of that code moves in order](camadas/gate.md#rvrnr--revisionrenumber--the-revisions-a-branch-added-move-to-a-free-number-when-the-base-took-theirs) `RVRNR-B03`

- [A revision the branch added whose number the base does not use stays](camadas/gate.md#rvrnr--revisionrenumber--the-revisions-a-branch-added-move-to-a-free-number-when-the-base-took-theirs) `RVRNR-B04`

- [A citation is rewritten only on a line the branch added](camadas/gate.md#rvrnr--revisionrenumber--the-revisions-a-branch-added-move-to-a-free-number-when-the-base-took-theirs) `RVRNR-B05`

- [The branch adding the same number twice is refused, naming the code](camadas/gate.md#rvrnr--revisionrenumber--the-revisions-a-branch-added-move-to-a-free-number-when-the-base-took-theirs) `RVRNR-B06`

- [A revision the base has is never renumbered, even with its explanation edited](camadas/gate.md#rvrnr--revisionrenumber--the-revisions-a-branch-added-move-to-a-free-number-when-the-base-took-theirs) `RVRNR-I01`

- [Rewriting is one pass, and a chain of renames never cascades](camadas/gate.md#rvrnr--revisionrenumber--the-revisions-a-branch-added-move-to-a-free-number-when-the-base-took-theirs) `RVRNR-I02`

- [Does not rewrite a revision cited without its unit code](camadas/gate.md#rvrnr--revisionrenumber--the-revisions-a-branch-added-move-to-a-free-number-when-the-base-took-theirs) `RVRNR-X01`

- [An artifact that is not a screen leaves without a verdict, and says why](camadas/gate.md#rtdcl--routedeclared--a-screen-declares-how-one-arrives-and-names-its-neighbours) `RTDCL-B01`

- [A screen with no named route fails](camadas/gate.md#rtdcl--routedeclared--a-screen-declares-how-one-arrives-and-names-its-neighbours) `RTDCL-B02`

- [A navigation row carrying a generic term fails](camadas/gate.md#rtdcl--routedeclared--a-screen-declares-how-one-arrives-and-names-its-neighbours) `RTDCL-B03`

- [A screen with a named route and concrete neighbours passes](camadas/gate.md#rtdcl--routedeclared--a-screen-declares-how-one-arrives-and-names-its-neighbours) `RTDCL-B04`

- [Route and navigation are recognised in either declared language](camadas/gate.md#rtdcl--routedeclared--a-screen-declares-how-one-arrives-and-names-its-neighbours) `RTDCL-B05`

- [The header's declared layer wins over the node's tags](camadas/gate.md#rtdcl--routedeclared--a-screen-declares-how-one-arrives-and-names-its-neighbours) `RTDCL-I01`

- [A generic term in prose is not accused](camadas/gate.md#rtdcl--routedeclared--a-screen-declares-how-one-arrives-and-names-its-neighbours) `RTDCL-I02`

- [No layer other than screen is charged for a route](camadas/gate.md#rtdcl--routedeclared--a-screen-declares-how-one-arrives-and-names-its-neighbours) `RTDCL-X01`

- [The declared route is not confronted against the real router](camadas/gate.md#rtdcl--routedeclared--a-screen-declares-how-one-arrives-and-names-its-neighbours) `RTDCL-X02`

- [Non-spec artifacts skip confrontation](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-B01`

- [Specifications without declared routes skip confrontation](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-B02`

- [Specifications with declared routes return Pending when route registry is unconfigured](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-B03`

- [Invalid route registry glob patterns return Pending](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-B04`

- [Route registry files containing zero registered routes return Pending](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-B05`

- [Declared screen route matching a component navigation prop passes](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-B06`

- [Declared screen route matching a navigation stack parameter type entry passes](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-B07`

- [Declared backend route matching an HTTP resource registration passes](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-B08`

- [HTTP method verb prefixes are stripped when evaluating declared routes](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-B09`

- [Route paths match regardless of leading slash differences between specification and code](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-B10`

- [Custom route pattern regex matches custom route registration patterns](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-B11`

- [Custom route pattern regex that is malformed or lacks capture groups returns Pending](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-B12`

- [Declared routes missing from all registered route definitions fail](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-B13`

- [Route presence is validated only for specifications](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-I01`

- [The gate never approves route existence without inspecting registry files](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-I02`

- [Route matching is slash-normalized between specification and code](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-I03`

- [Missing routes produce a blocking Fail verdict](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-I04`

- [The gate does not require every specification to declare a route](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-X01`

- [Route parameter schemas and payload contracts are not evaluated by this gate](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-X02`

- [Route access permissions and authentication middlewares are outside evaluation scope](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-X03`

- [Building an identifier joins the gate and the rule](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B01`

- [A gate with a single verification gains no separator](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B02`

- [The identifier decomposes into gate and rule](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B03`

- [The rule half is empty when the identifier carries only a gate](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B04`

- [A waiver with no reason is refused](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B05`

- [A waiver with no rule name is refused](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B06`

- [A path is refused as the waiver target](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B07`

- [A target marker with nothing after it is refused](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B08`

- [The waiver accepts both granularities](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B09`

- [A waiver that declares targets does not hold for the whole gate](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B10`

- [A waiver by target spares the named codes and confronts the rest](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B11`

- [A waiver with no declared target holds for every code](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B12`

- [An artifact with no code is not reached by a waiver restricted to targets](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B13`

- [Each target carries its own reason](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B14`

- [The commit message declares waivers that survive in the history](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B15`

- [A commit marker whose reason is blank is refused](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B16`

- [Two waivers merge instead of forcing a choice](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-B17`

- [A waiver never reaches what nobody waived](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-I01`

- [Every accepted waiver carries a written reason](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-I02`

- [A refusal always reaches the caller as an error](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-I03`

- [The unit does not accept a path as the waiver target](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-X01`

- [The unit does not decide whether a rule passes](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-X02`

- [The unit reads neither files nor the map](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it) `RLUEX-X03`

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

- [Non-feature artifacts skip confrontation](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes) `SCASS-B01`

- [Genuinely assertive outcome steps pass](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes) `SCASS-B02`

- [Tautological outcome steps citing only the code fail](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes) `SCASS-B03`

- [Tautological variations wrapped in linking words fail](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes) `SCASS-B04`

- [Outcome steps citing a code with substantive assertion pass](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes) `SCASS-B05`

- [Multiple tautological outcome steps are reported sorted and deduplicated](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes) `SCASS-B06`

- [Outcome steps across recognized dialect keywords are enforced](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes) `SCASS-B07`

- [Empty and comment lines are ignored](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes) `SCASS-B08`

- [Non-outcome steps citing codes are ignored](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes) `SCASS-B09`

- [Up to two residual content words is classified as a tautology](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes) `SCASS-I01`

- [Truly assertive outcome steps are never flagged as tautologies](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes) `SCASS-I02`

- [Language recognition covers all supported dialect alternatives](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes) `SCASS-I03`

- [Prose style and semantic elegance are not graded](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes) `SCASS-X01`

- [Setup and trigger steps are not inspected for code citations](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes) `SCASS-X02`

- [Scenario or test presence is not enforced by this gate](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes) `SCASS-X03`

- [Two scenarios sharing one code are reported](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-B01`

- [The report names the repeated code](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-B02`

- [The report says how many scenarios share the code](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-B03`

- [The report teaches the way out with the project's own code](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-B04`

- [The suffix gives each scenario its own identity](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-B05`

- [Distinct codes pass](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-B06`

- [The verdict is Pending and never a failure](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-B07`

- [An artifact that is not a feature leaves without a verdict](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-B08`

- [A feature with no coded scenario leaves without a verdict](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-B09`

- [Several repeated codes are reported together in a stable order](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-B10`

- [Long titles are shortened in the report](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-B11`

- [Grouping is by the complete code, suffix included](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-I01`

- [The message is deterministic](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-I02`

- [The gate does not judge whether the two scenarios describe different behaviours](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-X01`

- [The gate does not look across features](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-X02`

- [The gate does not charge the absence of a code on a scenario](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-X03`

- [The gate does not renumber the scenarios](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-X04`

- [An artifact that is not a feature leaves without a verdict](camadas/gate.md#scltr--scenarioletterdeclared--the-letter-of-a-scenario-code-exists-in-the-vocabulary) `SCLTR-B01`

- [With no declared vocabulary the gate leaves without a verdict](camadas/gate.md#scltr--scenarioletterdeclared--the-letter-of-a-scenario-code-exists-in-the-vocabulary) `SCLTR-B02`

- [A feature carrying no scenario code leaves without a verdict](camadas/gate.md#scltr--scenarioletterdeclared--the-letter-of-a-scenario-code-exists-in-the-vocabulary) `SCLTR-B03`

- [Every letter inside the vocabulary passes](camadas/gate.md#scltr--scenarioletterdeclared--the-letter-of-a-scenario-code-exists-in-the-vocabulary) `SCLTR-B04`

- [A letter outside the vocabulary is undetermined, not a failure](camadas/gate.md#scltr--scenarioletterdeclared--the-letter-of-a-scenario-code-exists-in-the-vocabulary) `SCLTR-B05`

- [The verdict names the letters that are outside and the codes carrying them](camadas/gate.md#scltr--scenarioletterdeclared--the-letter-of-a-scenario-code-exists-in-the-vocabulary) `SCLTR-B06`

- [Codes sharing one unknown letter are grouped into a single line](camadas/gate.md#scltr--scenarioletterdeclared--the-letter-of-a-scenario-code-exists-in-the-vocabulary) `SCLTR-B07`

- [The scan is over the shape of a code, never over the vocabulary](camadas/gate.md#scltr--scenarioletterdeclared--the-letter-of-a-scenario-code-exists-in-the-vocabulary) `SCLTR-I01`

- [The code-length pattern is read at every call](camadas/gate.md#scltr--scenarioletterdeclared--the-letter-of-a-scenario-code-exists-in-the-vocabulary) `SCLTR-I02`

- [A valid letter is never named in the verdict](camadas/gate.md#scltr--scenarioletterdeclared--the-letter-of-a-scenario-code-exists-in-the-vocabulary) `SCLTR-I03`

- [The gate does not choose between declaring and remapping](camadas/gate.md#scltr--scenarioletterdeclared--the-letter-of-a-scenario-code-exists-in-the-vocabulary) `SCLTR-X01`

- [The tags accompanying a code are not judged](camadas/gate.md#scltr--scenarioletterdeclared--the-letter-of-a-scenario-code-exists-in-the-vocabulary) `SCLTR-X02`

- [Non-feature artifacts skip confrontation](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-B01`

- [Missing or empty rule types in configuration skip confrontation](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-B02`

- [Configuration without tag mappings skips confrontation](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-B03`

- [Features without coded scenarios skip confrontation](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-B04`

- [Codes lacking a recognized rule letter are ignored](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-B05`

- [Unregistered scenario tags are ignored](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-B06`

- [Scenarios whose classification tags align with code letters pass](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-B07`

- [A tag mapped to multiple rule letters passes if any matches](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-B08`

- [A multi-coded scenario passes when a tag matches any code](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-B09`

- [A scenario tag disagreeing with the code rule letter returns Pending](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-B10`

- [The Pending verdict cites details of the type divergence](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-B11`

- [Multiple mismatch findings are sorted deterministically](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-B12`

- [Classification alignment is evaluated only on feature artifacts](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-I01`

- [Absence of tag mappings prevents speculative enforcement](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-I02`

- [Mismatched scenario types return Pending rather than Fail](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-I03`

- [Secondary requirement codes prevent false mismatch reporting](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-I04`

- [The gate does not mandate classification tags on every scenario](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-X01`

- [The gate does not decide whether tag or code is erroneous](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-X02`

- [Specifications and test files are not inspected or modified](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-X03`

- [Tags are permitted to map across multiple rule letters](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-X04`

- [An artifact that is not code leaves without a verdict](camadas/gate.md#sbgrd--siblingguard--sibling-functions-treat-the-same-parameter-consistently) `SBGRD-B01`

- [Without a declared dialect the verdict is undetermined](camadas/gate.md#sbgrd--siblingguard--sibling-functions-treat-the-same-parameter-consistently) `SBGRD-B02`

- [Fewer than three siblings on the same parameter is left alone](camadas/gate.md#sbgrd--siblingguard--sibling-functions-treat-the-same-parameter-consistently) `SBGRD-B03`

- [The sibling that does not guard is accused](camadas/gate.md#sbgrd--siblingguard--sibling-functions-treat-the-same-parameter-consistently) `SBGRD-B04`

- [When every sibling guards, nothing is accused](camadas/gate.md#sbgrd--siblingguard--sibling-functions-treat-the-same-parameter-consistently) `SBGRD-B05`

- [When no sibling guards, nothing is accused either](camadas/gate.md#sbgrd--siblingguard--sibling-functions-treat-the-same-parameter-consistently) `SBGRD-B06`

- [The verdict names the function and the parameter](camadas/gate.md#sbgrd--siblingguard--sibling-functions-treat-the-same-parameter-consistently) `SBGRD-B07`

- [A waiver with a written reason silences the accusation](camadas/gate.md#sbgrd--siblingguard--sibling-functions-treat-the-same-parameter-consistently) `SBGRD-B08`

- [The gate never judges what the guard does](camadas/gate.md#sbgrd--siblingguard--sibling-functions-treat-the-same-parameter-consistently) `SBGRD-I01`

- [The three conservatism conditions hold together](camadas/gate.md#sbgrd--siblingguard--sibling-functions-treat-the-same-parameter-consistently) `SBGRD-I02`

- [The gate does not invent what an exported function looks like](camadas/gate.md#sbgrd--siblingguard--sibling-functions-treat-the-same-parameter-consistently) `SBGRD-X01`

- [A single function in isolation is not accused](camadas/gate.md#sbgrd--siblingguard--sibling-functions-treat-the-same-parameter-consistently) `SBGRD-X02`

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

- [Non-test artifacts skip confrontation](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-B01`

- [Confrontation without a dependency graph returns Pending](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-B02`

- [A test with no linked feature skips confrontation](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-B03`

- [A test whose linked feature cannot be read skips confrontation](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-B04`

- [A test whose linked feature declares no scenario codes skips confrontation](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-B05`

- [A test containing a declared scenario code passes](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-B06`

- [A single declared scenario code is sufficient to pass](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-B07`

- [A test containing none of the linked feature scenario codes fails](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-B08`

- [A test with transposed or misspelled codes fails exact substring confrontation](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-B09`

- [A failing verdict names the linked feature and lists expected codes](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-B10`

- [Scenario codes appearing anywhere in test content satisfy traceability](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-B11`

- [Traceability is strictly charged on test artifacts](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-I01`

- [Tests without a linked feature are never charged](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-I02`

- [The acceptance threshold requires only one scenario code present](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-I03`

- [Missing graph structure produces Pending rather than approving](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-I04`

- [The gate does not enforce one-to-one scenario coverage](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-X01`

- [Test execution results are not inspected by this gate](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-X02`

- [Scenario codes are not required in standalone test files](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-X03`

- [Test assertion semantics and quality are not evaluated](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-X04`

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

- [A piece declared TO BE DEVELOPED leaves the verdict undetermined](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist) `TRCMT-B08`

- [The gate does not judge the QUALITY of any piece](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist) `TRCMT-X02`

- [Non-spec artifacts skip confrontation](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-B01`

- [Specifications without cited triggers skip confrontation](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-B02`

- [Cited triggers return Pending when no vocabulary is declared](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-B03`

- [Non-trigger key-value citations are ignored](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-B04`

- [Natural language mentions without backticks are ignored](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-B05`

- [Valid declared trigger values pass confrontation](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-B06`

- [Valid declared obligation names pass confrontation](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-B07`

- [An undeclared trigger value fails with a nearest-match suggestion](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-B08`

- [An undeclared trigger without close match lists declared triggers](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-B09`

- [An undeclared obligation name fails confrontation](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-B10`

- [Duplicate citations of triggers or obligations are deduplicated](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-B11`

- [Multiple vocabulary errors are sorted deterministically](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-B12`

- [The list of declared triggers is capped at four](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-B13`

- [Trigger confrontation applies exclusively to specification artifacts](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-I01`

- [Missing compliance packs produce Pending rather than Pass](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-I02`

- [Trigger keys are limited to recognized compliance predicates](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-I03`

- [Undeclared triggers produce a blocking Fail verdict](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-I04`

- [The gate does not require every specification to cite triggers](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-X01`

- [Code and test artifacts are not checked for trigger citations](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-X02`

- [Code obligation implementation is not evaluated by this gate](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-X03`

- [Unquoted prose text is not evaluated as symbol citations](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-X04`

- [A declaration matching its code line passes](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B01`

- [A declaration whose code line says another value fails](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B02`

- [Comment and blank lines between declaration and code are skipped](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B03`

- [Copies of the same key with different values are reported](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B04`

- [The value a rule declares in the spec is the source](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B05`

- [Without a declared pattern the gate skips and says how to enable it](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B06`

- [A declaration with nothing below annotates nothing](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B07`

- [An anchor inside prose is a mention, not a declaration](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B08`

- [A prose rule declares no value](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-X01`

- [A literal nobody declared is not charged](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-X02`

- [Artifacts that are not feature files leave with verdict Skip](camadas/gate.md#vrbsv--vrbaseline--ensures-visual-regression-scenarios-have-captured-reference-baseline-images) `VRBSV-B01`

- [Feature files declaring no visual regression scenarios leave with verdict Skip](camadas/gate.md#vrbsv--vrbaseline--ensures-visual-regression-scenarios-have-captured-reference-baseline-images) `VRBSV-B02`

- [Visual regression scenarios having matching baseline images pass](camadas/gate.md#vrbsv--vrbaseline--ensures-visual-regression-scenarios-have-captured-reference-baseline-images) `VRBSV-B03`

- [Visual regression scenarios lacking baseline images fail naming missing codes](camadas/gate.md#vrbsv--vrbaseline--ensures-visual-regression-scenarios-have-captured-reference-baseline-images) `VRBSV-B04`

- [Baseline image matching supports naming variants](camadas/gate.md#vrbsv--vrbaseline--ensures-visual-regression-scenarios-have-captured-reference-baseline-images) `VRBSV-B05`

- [Reads visual regime tag from project configuration](camadas/gate.md#vrbsv--vrbaseline--ensures-visual-regression-scenarios-have-captured-reference-baseline-images) `VRBSV-B06`

- [Falls back to default visual regime tag when configuration is missing](camadas/gate.md#vrbsv--vrbaseline--ensures-visual-regression-scenarios-have-captured-reference-baseline-images) `VRBSV-B07`

- [Multiple missing baseline scenario codes are sorted alphabetically](camadas/gate.md#vrbsv--vrbaseline--ensures-visual-regression-scenarios-have-captured-reference-baseline-images) `VRBSV-B08`

- [Only scenario codes containing visual regression suffix are matched](camadas/gate.md#vrbsv--vrbaseline--ensures-visual-regression-scenarios-have-captured-reference-baseline-images) `VRBSV-B09`

- [Every visual regression scenario declared must correspond to a baseline image](camadas/gate.md#vrbsv--vrbaseline--ensures-visual-regression-scenarios-have-captured-reference-baseline-images) `VRBSV-I01`

- [Visual regime tag mapping treats map key as tag in feature](camadas/gate.md#vrbsv--vrbaseline--ensures-visual-regression-scenarios-have-captured-reference-baseline-images) `VRBSV-I02`

- [Does not evaluate baseline staleness using disk modification timestamps](camadas/gate.md#vrbsv--vrbaseline--ensures-visual-regression-scenarios-have-captured-reference-baseline-images) `VRBSV-X01`

- [Does not fail commits based on git commit dates of baseline images](camadas/gate.md#vrbsv--vrbaseline--ensures-visual-regression-scenarios-have-captured-reference-baseline-images) `VRBSV-X02`

