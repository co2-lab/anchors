<!-- anchors:generated from doct/comportamento.md.tmpl — DO NOT EDIT: run `anchors docs build` -->


# Comportamento

Todos os cenários do sistema. Cada um leva à unidade que o define.

Um cenário descreve o que o sistema faz numa situação — vem da feature, e é o mesmo que o
teste prova.

## gate

- [An artifact that is not a spec leaves without a verdict](camadas/gate.md#cdctc-b01--an-artifact-that-is-not-a-spec-leaves-without-a-verdict) `CDCTC-B01`

- [An exported symbol the spec never names fails, and the verdict names it](camadas/gate.md#cdctc-b02--an-exported-symbol-the-spec-never-names-fails-and-the-verdict-names-it) `CDCTC-B02`

- [What the spec already catalogues is never accused](camadas/gate.md#cdctc-b03--what-the-spec-already-catalogues-is-never-accused) `CDCTC-B03`

- [A no-rule marker with a written reason waives the symbol](camadas/gate.md#cdctc-b04--a-no-rule-marker-with-a-written-reason-waives-the-symbol) `CDCTC-B04`

- [A bare no-rule marker does not waive](camadas/gate.md#cdctc-b05--a-bare-no-rule-marker-does-not-waive) `CDCTC-B05`

- [A spec cataloguing every exported symbol passes](camadas/gate.md#cdctc-b06--a-spec-cataloguing-every-exported-symbol-passes) `CDCTC-B06`

- [With no code linked the gate leaves without a verdict](camadas/gate.md#cdctc-b07--with-no-code-linked-the-gate-leaves-without-a-verdict) `CDCTC-B07`

- [Without a declared export pattern the gate skips and says so](camadas/gate.md#cdctc-b08--without-a-declared-export-pattern-the-gate-skips-and-says-so) `CDCTC-B08`

- [With the pattern declared the gate confronts for real in any language](camadas/gate.md#cdctc-b09--with-the-pattern-declared-the-gate-confronts-for-real-in-any-language) `CDCTC-B09`

- [The declared dialect family also supplies the pattern](camadas/gate.md#cdctc-b10--the-declared-dialect-family-also-supplies-the-pattern) `CDCTC-B10`

- [The waiver holds in the comment block above the symbol](camadas/gate.md#cdctc-i01--the-waiver-holds-in-the-comment-block-above-the-symbol) `CDCTC-I01`

- [The waiver does not leak between symbols](camadas/gate.md#cdctc-i02--the-waiver-does-not-leak-between-symbols) `CDCTC-I02`

- [The gate never approves a language it cannot read](camadas/gate.md#cdctc-i03--the-gate-never-approves-a-language-it-cannot-read) `CDCTC-I03`

- [The gate does not judge whether the rule describes the symbol well](camadas/gate.md#cdctc-x01--the-gate-does-not-judge-whether-the-rule-describes-the-symbol-well) `CDCTC-X01`

- [The gate does not decide which symbols deserve a rule](camadas/gate.md#cdctc-x02--the-gate-does-not-decide-which-symbols-deserve-a-rule) `CDCTC-X02`

- [The gate knows no language, the project declares what is public](camadas/gate.md#cdctc-x03--the-gate-knows-no-language-the-project-declares-what-is-public) `CDCTC-X03`

- [The gate does not charge the absence of code](camadas/gate.md#cdctc-x04--the-gate-does-not-charge-the-absence-of-code) `CDCTC-X04`

- [An identifier in the wrong language is accused, and an English one passes](camadas/gate.md#cdlng-b01--an-identifier-in-the-wrong-language-is-accused-and-an-english-one-passes) `CDLNG-B01`

- [The verdict returns the word that accused](camadas/gate.md#cdlng-b02--the-verdict-returns-the-word-that-accused) `CDLNG-B02`

- [Only a DECLARATION is the subject](camadas/gate.md#cdlng-b03--only-a-declaration-is-the-subject) `CDLNG-B03`

- [The declarations are found in every form the language offers](camadas/gate.md#cdlng-b04--the-declarations-are-found-in-every-form-the-language-offers) `CDLNG-B04`

- [Deciding one word is separate from deciding a whole identifier](camadas/gate.md#cdlng-b05--deciding-one-word-is-separate-from-deciding-a-whole-identifier) `CDLNG-B05`

- [A short word does not count](camadas/gate.md#cdlng-i01--a-short-word-does-not-count) `CDLNG-I01`

- [The gate does not read comments](camadas/gate.md#cdlng-x01--the-gate-does-not-read-comments) `CDLNG-X01`

- [The gate does not read user-facing text](camadas/gate.md#cdlng-x02--the-gate-does-not-read-user-facing-text) `CDLNG-X02`

- [The gate does not use a dictionary to decide the language](camadas/gate.md#cdlng-x03--the-gate-does-not-use-a-dictionary-to-decide-the-language) `CDLNG-X03`

- [A status emitted and not declared is accused by number](camadas/gate.md#csdcn-b01--a-status-emitted-and-not-declared-is-accused-by-number) `CSDCN-B01`

- [A status declared and emitted by no path is accused as a phantom](camadas/gate.md#csdcn-b02--a-status-declared-and-emitted-by-no-path-is-accused-as-a-phantom) `CSDCN-B02`

- [A faithful table passes](camadas/gate.md#csdcn-b03--a-faithful-table-passes) `CSDCN-B03`

- [The 500 of the top-level try/catch is not charged](camadas/gate.md#csdcn-b04--the-500-of-the-top-level-trycatch-is-not-charged) `CSDCN-B04`

- [The 5xx range covers, the 4xx range does not](camadas/gate.md#csdcn-b05--the-5xx-range-covers-the-4xx-range-does-not) `CSDCN-B05`

- [A status that lives only in a comment is not emitted](camadas/gate.md#csdcn-b06--a-status-that-lives-only-in-a-comment-is-not-emitted) `CSDCN-B06`

- [Without the contract section there is nothing to confront](camadas/gate.md#csdcn-b07--without-the-contract-section-there-is-nothing-to-confront) `CSDCN-B07`

- [Code that returns no status is skipped](camadas/gate.md#csdcn-b08--code-that-returns-no-status-is-skipped) `CSDCN-B08`

- [A literal status passed to a local helper counts as emitted](camadas/gate.md#csdcn-b09--a-literal-status-passed-to-a-local-helper-counts-as-emitted) `CSDCN-B09`

- [With a dynamic status the phantom side goes quiet and the literals still count](camadas/gate.md#csdcn-b10--with-a-dynamic-status-the-phantom-side-goes-quiet-and-the-literals-still-count) `CSDCN-B10`

- [Without a declared dialect the verdict is Pending](camadas/gate.md#csdcn-b11--without-a-declared-dialect-the-verdict-is-pending) `CSDCN-B11`

- [An explicit opt-out of the http_status field is honoured](camadas/gate.md#csdcn-b12--an-explicit-opt-out-of-the-http-status-field-is-honoured) `CSDCN-B12`

- [The lexicon comes from the project's dialect, not from the gate](camadas/gate.md#csdcn-i01--the-lexicon-comes-from-the-projects-dialect-not-from-the-gate) `CSDCN-I01`

- [A dialect declared by hand teaches the gate its own lexicon](camadas/gate.md#csdcn-i02--a-dialect-declared-by-hand-teaches-the-gate-its-own-lexicon) `CSDCN-I02`

- [A named constant is worth the number it means](camadas/gate.md#csdcn-i03--a-named-constant-is-worth-the-number-it-means) `CSDCN-I03`

- [The gate does not demand the generic ranges](camadas/gate.md#csdcn-x01--the-gate-does-not-demand-the-generic-ranges) `CSDCN-X01`

- [The gate does not judge when each status is right](camadas/gate.md#csdcn-x02--the-gate-does-not-judge-when-each-status-is-right) `CSDCN-X02`

- [The gate does not charge the phantom side under a dynamic status](camadas/gate.md#csdcn-x03--the-gate-does-not-charge-the-phantom-side-under-a-dynamic-status) `CSDCN-X03`

- [A mandatory document that does not exist fails](camadas/gate.md#dcrqd-b01--a-mandatory-document-that-does-not-exist-fails) `DCRQD-B01`

- [A document that exists and does not mention the unit fails](camadas/gate.md#dcrqd-b02--a-document-that-exists-and-does-not-mention-the-unit-fails) `DCRQD-B02`

- [A mention by the identity code counts as documented](camadas/gate.md#dcrqd-b03--a-mention-by-the-identity-code-counts-as-documented) `DCRQD-B03`

- [A mention by the file name also counts](camadas/gate.md#dcrqd-b04--a-mention-by-the-file-name-also-counts) `DCRQD-B04`

- [Satisfying one of two duties is not enough](camadas/gate.md#dcrqd-b05--satisfying-one-of-two-duties-is-not-enough) `DCRQD-B05`

- [Without a declaration nothing is charged](camadas/gate.md#dcrqd-b06--without-a-declaration-nothing-is-charged) `DCRQD-B06`

- [A layer with no trigger is not charged](camadas/gate.md#dcrqd-b07--a-layer-with-no-trigger-is-not-charged) `DCRQD-B07`

- [Aggregated, the verdict is one per document](camadas/gate.md#dcrqd-b08--aggregated-the-verdict-is-one-per-document) `DCRQD-B08`

- [The duty starts from the spec, not from the code](camadas/gate.md#dcrqd-i01--the-duty-starts-from-the-spec-not-from-the-code) `DCRQD-I01`

- [The layer used is the UNIT's, not the node's](camadas/gate.md#dcrqd-i02--the-layer-used-is-the-units-not-the-nodes) `DCRQD-I02`

- [Without a map the aggregated verdict is skipped](camadas/gate.md#dcrqd-i03--without-a-map-the-aggregated-verdict-is-skipped) `DCRQD-I03`

- [The gate does not understand the document's content](camadas/gate.md#dcrqd-x01--the-gate-does-not-understand-the-documents-content) `DCRQD-X01`

- [The gate does not decide which documents are mandatory](camadas/gate.md#dcrqd-x02--the-gate-does-not-decide-which-documents-are-mandatory) `DCRQD-X02`

- [A reference that brings the passage it announces passes](camadas/gate.md#dscdc-b01--a-reference-that-brings-the-passage-it-announces-passes) `DSCDC-B01`

- [A reference that only points is accused](camadas/gate.md#dscdc-b02--a-reference-that-only-points-is-accused) `DSCDC-B02`

- [A quotation counts in any written tradition](camadas/gate.md#dscdc-b03--a-quotation-counts-in-any-written-tradition) `DSCDC-B03`

- [A path inside a code fence is an example, not a reference](camadas/gate.md#dscdc-b04--a-path-inside-a-code-fence-is-an-example-not-a-reference) `DSCDC-B04`

- [The spec citing its own path is identifying itself](camadas/gate.md#dscdc-b05--the-spec-citing-its-own-path-is-identifying-itself) `DSCDC-B05`

- [A rule code is not a revision](camadas/gate.md#dscdc-b06--a-rule-code-is-not-a-revision) `DSCDC-B06`

- [Only the spec is charged](camadas/gate.md#dscdc-b07--only-the-spec-is-charged) `DSCDC-B07`

- [A revision cited with an explanation on the same line passes](camadas/gate.md#dscdc-b08--a-revision-cited-with-an-explanation-on-the-same-line-passes) `DSCDC-B08`

- [With no map the confrontation is skipped](camadas/gate.md#dscdc-b09--with-no-map-the-confrontation-is-skipped) `DSCDC-B09`

- [The ruler matches structure, never vocabulary](camadas/gate.md#dscdc-i01--the-ruler-matches-structure-never-vocabulary) `DSCDC-I01`

- [The on-line explanation escape belongs to the revision, not to the path](camadas/gate.md#dscdc-i02--the-on-line-explanation-escape-belongs-to-the-revision-not-to-the-path) `DSCDC-I02`

- [The verdict names the line and shows what it says](camadas/gate.md#dscdc-i03--the-verdict-names-the-line-and-shows-what-it-says) `DSCDC-I03`

- [The gate does not judge whether the accompanying content is faithful](camadas/gate.md#dscdc-x01--the-gate-does-not-judge-whether-the-accompanying-content-is-faithful) `DSCDC-X01`

- [The gate marks and does not block](camadas/gate.md#dscdc-x02--the-gate-marks-and-does-not-block) `DSCDC-X02`

- [Measuring explanation errs on the permissive side](camadas/gate.md#dscdc-x03--measuring-explanation-errs-on-the-permissive-side) `DSCDC-X03`

- [An artifact that is not a spec leaves without a verdict](camadas/gate.md#dmdcd-b01--an-artifact-that-is-not-a-spec-leaves-without-a-verdict) `DMDCD-B01`

- [A spec without the domain section is failed](camadas/gate.md#dmdcd-b02--a-spec-without-the-domain-section-is-failed) `DMDCD-B02`

- [A section opened and left empty is failed](camadas/gate.md#dmdcd-b03--a-section-opened-and-left-empty-is-failed) `DMDCD-B03`

- [A row filled only with placeholders is not a declaration](camadas/gate.md#dmdcd-b06--a-row-filled-only-with-placeholders-is-not-a-declaration) `DMDCD-B06`

- [An entry with no owner is failed, and the verdict names it](camadas/gate.md#dmdcd-b04--an-entry-with-no-owner-is-failed-and-the-verdict-names-it) `DMDCD-B04`

- [An entry whose owner is named passes](camadas/gate.md#dmdcd-b07--an-entry-whose-owner-is-named-passes) `DMDCD-B07`

- [A waiver with a written reason silences the gate](camadas/gate.md#dmdcd-b05--a-waiver-with-a-written-reason-silences-the-gate) `DMDCD-B05`

- [A bare waiver, with no reason, does not waive](camadas/gate.md#dmdcd-i01--a-bare-waiver-with-no-reason-does-not-waive) `DMDCD-I01`

- [A non-answer in the owner column is not an owner](camadas/gate.md#dmdcd-i02--a-non-answer-in-the-owner-column-is-not-an-owner) `DMDCD-I02`

- [The gate does not judge whether the declared domain is correct](camadas/gate.md#dmdcd-x01--the-gate-does-not-judge-whether-the-declared-domain-is-correct) `DMDCD-X01`

- [The gate does not read the code to check the validation exists](camadas/gate.md#dmdcd-x02--the-gate-does-not-read-the-code-to-check-the-validation-exists) `DMDCD-X02`

- [An artifact that is not a test leaves without a verdict](camadas/gate.md#evfrv-b01--an-artifact-that-is-not-a-test-leaves-without-a-verdict) `EVFRV-B01`

- [Without a built map the gate stays quiet](camadas/gate.md#evfrv-b02--without-a-built-map-the-gate-stays-quiet) `EVFRV-B02`

- [A test with no execution stamp is skipped, not failed](camadas/gate.md#evfrv-b03--a-test-with-no-execution-stamp-is-skipped-not-failed) `EVFRV-B03`

- [A test whose closure is intact passes](camadas/gate.md#evfrv-b04--a-test-whose-closure-is-intact-passes) `EVFRV-B04`

- [The passing verdict says what it checked against](camadas/gate.md#evfrv-b05--the-passing-verdict-says-what-it-checked-against) `EVFRV-B05`

- [A test whose dependency advanced a revision fails](camadas/gate.md#evfrv-b06--a-test-whose-dependency-advanced-a-revision-fails) `EVFRV-B06`

- [The failing verdict names the culprit](camadas/gate.md#evfrv-b07--the-failing-verdict-names-the-culprit) `EVFRV-B07`

- [The failing verdict states the fix](camadas/gate.md#evfrv-b08--the-failing-verdict-states-the-fix) `EVFRV-B08`

- [A test whose own file changed is reported separately from its closure](camadas/gate.md#evfrv-b09--a-test-whose-own-file-changed-is-reported-separately-from-its-closure) `EVFRV-B09`

- [The culprit list is truncated at five and the remainder counted](camadas/gate.md#evfrv-b10--the-culprit-list-is-truncated-at-five-and-the-remainder-counted) `EVFRV-B10`

- [Absence of proof and expired proof are never the same finding](camadas/gate.md#evfrv-i01--absence-of-proof-and-expired-proof-are-never-the-same-finding) `EVFRV-I01`

- [A test that never ran is never approved either](camadas/gate.md#evfrv-i02--a-test-that-never-ran-is-never-approved-either) `EVFRV-I02`

- [Truncation never hides the size of the problem](camadas/gate.md#evfrv-i03--truncation-never-hides-the-size-of-the-problem) `EVFRV-I03`

- [The gate does not charge the absence of a green test](camadas/gate.md#evfrv-x01--the-gate-does-not-charge-the-absence-of-a-green-test) `EVFRV-X01`

- [The gate does not run the test nor judge whether the change broke it](camadas/gate.md#evfrv-x02--the-gate-does-not-run-the-test-nor-judge-whether-the-change-broke-it) `EVFRV-X02`

- [The gate does not read the project's configuration](camadas/gate.md#evfrv-x03--the-gate-does-not-read-the-projects-configuration) `EVFRV-X03`

- [An artifact that is not code leaves without a verdict](camadas/gate.md#lybnl-b01--an-artifact-that-is-not-code-leaves-without-a-verdict) `LYBNL-B01`

- [Content matching a forbidden pattern fails, naming line and reason](camadas/gate.md#lybnl-b02--content-matching-a-forbidden-pattern-fails-naming-line-and-reason) `LYBNL-B02`

- [A rule scoped to a layer charges only that layer](camadas/gate.md#lybnl-b03--a-rule-scoped-to-a-layer-charges-only-that-layer) `LYBNL-B03`

- [A rule with no layer holds for all code](camadas/gate.md#lybnl-b04--a-rule-with-no-layer-holds-for-all-code) `LYBNL-B04`

- [Severity warn records without failing, and the default is error](camadas/gate.md#lybnl-b05--severity-warn-records-without-failing-and-the-default-is-error) `LYBNL-B05`

- [A waiver with a written reason on the line waives that line](camadas/gate.md#lybnl-b06--a-waiver-with-a-written-reason-on-the-line-waives-that-line) `LYBNL-B06`

- [The waiver also holds in the comment on the line above](camadas/gate.md#lybnl-b07--the-waiver-also-holds-in-the-comment-on-the-line-above) `LYBNL-B07`

- [A bare marker with no reason does not waive](camadas/gate.md#lybnl-b08--a-bare-marker-with-no-reason-does-not-waive) `LYBNL-B08`

- [With no boundary declared the verdict is Pending, never Pass](camadas/gate.md#lybnl-b09--with-no-boundary-declared-the-verdict-is-pending-never-pass) `LYBNL-B09`

- [An invalid forbid pattern fails visibly](camadas/gate.md#lybnl-b10--an-invalid-forbid-pattern-fails-visibly) `LYBNL-B10`

- [The pattern is matched against the whole file, catching a multi-line import](camadas/gate.md#lybnl-b11--the-pattern-is-matched-against-the-whole-file-catching-a-multi-line-import) `LYBNL-B11`

- [The same rule is expressible in six language dialects](camadas/gate.md#lybnl-i01--the-same-rule-is-expressible-in-six-language-dialects) `LYBNL-I01`

- [A single-line import of the same shape is still caught](camadas/gate.md#lybnl-i02--a-single-line-import-of-the-same-shape-is-still-caught) `LYBNL-I02`

- [The waiver holds on any line of the matched stretch](camadas/gate.md#lybnl-i03--the-waiver-holds-on-any-line-of-the-matched-stretch) `LYBNL-I03`

- [A line anchor keeps holding per line](camadas/gate.md#lybnl-i04--a-line-anchor-keeps-holding-per-line) `LYBNL-I04`

- [The gate does not decide which boundaries exist](camadas/gate.md#lybnl-x01--the-gate-does-not-decide-which-boundaries-exist) `LYBNL-X01`

- [The gate does not parse the language, it matches text](camadas/gate.md#lybnl-x02--the-gate-does-not-parse-the-language-it-matches-text) `LYBNL-X02`

- [The gate does not judge whether the boundary is the right one to draw](camadas/gate.md#lybnl-x03--the-gate-does-not-judge-whether-the-boundary-is-the-right-one-to-draw) `LYBNL-X03`

- [A rule marked at both declared scopes passes](camadas/gate.md#mrprm-b01--a-rule-marked-at-both-declared-scopes-passes) `MRPRM-B01`

- [A rule missing from one end fails, and the verdict names the empty scope](camadas/gate.md#mrprm-b02--a-rule-missing-from-one-end-fails-and-the-verdict-names-the-empty-scope) `MRPRM-B02`

- [Two markings on the same side do not satisfy the gate](camadas/gate.md#mrprm-b03--two-markings-on-the-same-side-do-not-satisfy-the-gate) `MRPRM-B03`

- [A rule left over at one end fails and is named](camadas/gate.md#mrprm-b04--a-rule-left-over-at-one-end-fails-and-is-named) `MRPRM-B04`

- [Total absence of the prefix is not approval](camadas/gate.md#mrprm-b05--total-absence-of-the-prefix-is-not-approval) `MRPRM-B05`

- [A declaration with no prefix returns Pending](camadas/gate.md#mrprm-b06--a-declaration-with-no-prefix-returns-pending) `MRPRM-B06`

- [With no scopes declared the ruler is the count](camadas/gate.md#mrprm-b07--with-no-scopes-declared-the-ruler-is-the-count) `MRPRM-B07`

- [The marking crosses language](camadas/gate.md#mrprm-b08--the-marking-crosses-language) `MRPRM-B08`

- [Ignored directories never count towards parity](camadas/gate.md#mrprm-b09--ignored-directories-never-count-towards-parity) `MRPRM-B09`

- [The failing verdict names the rule and the empty scope](camadas/gate.md#mrprm-i01--the-failing-verdict-names-the-rule-and-the-empty-scope) `MRPRM-I01`

- [What was not measured is never approved](camadas/gate.md#mrprm-i02--what-was-not-measured-is-never-approved) `MRPRM-I02`

- [The gate does not read what each end actually does](camadas/gate.md#mrprm-x01--the-gate-does-not-read-what-each-end-actually-does) `MRPRM-X01`

- [The gate does not decide which rules live at two ends](camadas/gate.md#mrprm-x02--the-gate-does-not-decide-which-rules-live-at-two-ends) `MRPRM-X02`

- [Files outside the text extension list are not read](camadas/gate.md#mrprm-x03--files-outside-the-text-extension-list-are-not-read) `MRPRM-X03`

- [A node that carries the trigger and is absent from the demanded file fails](camadas/gate.md#obhnb-b01--a-node-that-carries-the-trigger-and-is-absent-from-the-demanded-file-fails) `OBHNB-B01`

- [A node that carries the trigger and does appear passes](camadas/gate.md#obhnb-b02--a-node-that-carries-the-trigger-and-does-appear-passes) `OBHNB-B02`

- [A node without the trigger contracts no obligation](camadas/gate.md#obhnb-b03--a-node-without-the-trigger-contracts-no-obligation) `OBHNB-B03`

- [A waiver exempts only when it carries a written reason](camadas/gate.md#obhnb-b04--a-waiver-exempts-only-when-it-carries-a-written-reason) `OBHNB-B04`

- [A project with no declared obligation is skipped](camadas/gate.md#obhnb-b05--a-project-with-no-declared-obligation-is-skipped) `OBHNB-B05`

- [An acknowledged debt with a written when yields Pending](camadas/gate.md#obhnb-b06--an-acknowledged-debt-with-a-written-when-yields-pending) `OBHNB-B06`

- [A bare debt marker keeps failing](camadas/gate.md#obhnb-b07--a-bare-debt-marker-keeps-failing) `OBHNB-B07`

- [Waiver and debt stay distinct](camadas/gate.md#obhnb-b08--waiver-and-debt-stay-distinct) `OBHNB-B08`

- [The failing verdict offers the three ways out](camadas/gate.md#obhnb-b09--the-failing-verdict-offers-the-three-ways-out) `OBHNB-B09`

- [The token is derived through the declared form](camadas/gate.md#obhnb-i01--the-token-is-derived-through-the-declared-form) `OBHNB-I01`

- [A glob that matches no file produces no violation](camadas/gate.md#obhnb-i02--a-glob-that-matches-no-file-produces-no-violation) `OBHNB-I02`

- [The node's own identified_as wins over the automatic form](camadas/gate.md#obhnb-i03--the-nodes-own-identified-as-wins-over-the-automatic-form) `OBHNB-I03`

- [The gate does not decide which obligations exist](camadas/gate.md#obhnb-x01--the-gate-does-not-decide-which-obligations-exist) `OBHNB-X01`

- [The gate does not understand what the destination does with the token](camadas/gate.md#obhnb-x02--the-gate-does-not-understand-what-the-destination-does-with-the-token) `OBHNB-X02`

- [A declaration written in the body is not read](camadas/gate.md#obhnb-x03--a-declaration-written-in-the-body-is-not-read) `OBHNB-X03`

- [An artifact that is not a spec leaves without a verdict](camadas/gate.md#opqsp-b01--an-artifact-that-is-not-a-spec-leaves-without-a-verdict) `OPQSP-B01`

- [Whoever OPENED the section is confronted by its content](camadas/gate.md#opqsp-b02--whoever-opened-the-section-is-confronted-by-its-content) `OPQSP-B02`

- [An open item bars the spec](camadas/gate.md#opqsp-b03--an-open-item-bars-the-spec) `OPQSP-B03`

- [A section closed honestly releases the spec](camadas/gate.md#opqsp-b04--a-section-closed-honestly-releases-the-spec) `OPQSP-B04`

- [An item marked as resolved does not block](camadas/gate.md#opqsp-b05--an-item-marked-as-resolved-does-not-block) `OPQSP-B05`

- [A question with no code is charged](camadas/gate.md#opqsp-b06--a-question-with-no-code-is-charged) `OPQSP-B06`

- [The count of pending decisions reads the project's own lexicon](camadas/gate.md#opqsp-b07--the-count-of-pending-decisions-reads-the-projects-own-lexicon) `OPQSP-B07`

- [Prose is not an item](camadas/gate.md#opqsp-i01--prose-is-not-an-item) `OPQSP-I01`

- [The section boundary is respected](camadas/gate.md#opqsp-i02--the-section-boundary-is-respected) `OPQSP-I02`

- [Filling in what the question BECOMES does not close the question](camadas/gate.md#opqsp-i03--filling-in-what-the-question-becomes-does-not-close-the-question) `OPQSP-I03`

- [The gate does not judge whether the question is good](camadas/gate.md#opqsp-x01--the-gate-does-not-judge-whether-the-question-is-good) `OPQSP-X01`

- [A spec with no section is a pending item, and the verdict teaches the way out](camadas/gate.md#opqsp-x02--a-spec-with-no-section-is-a-pending-item-and-the-verdict-teaches-the-way-out) `OPQSP-X02`

- [A limit received from the caller passes](camadas/gate.md#pgnhn-b01--a-limit-received-from-the-caller-passes) `PGNHN-B01`

- [A limit hidden in a default value is accused](camadas/gate.md#pgnhn-b02--a-limit-hidden-in-a-default-value-is-accused) `PGNHN-B02`

- [The NAME bounds the promise](camadas/gate.md#pgnhn-b03--the-name-bounds-the-promise) `PGNHN-B03`

- [Sibling functions that paginate are the proof by asymmetry](camadas/gate.md#pgnhn-b04--sibling-functions-that-paginate-are-the-proof-by-asymmetry) `PGNHN-B04`

- [A waiver with a written reason leaves the report](camadas/gate.md#pgnhn-b05--a-waiver-with-a-written-reason-leaves-the-report) `PGNHN-B05`

- [The verdict offers the way out](camadas/gate.md#pgnhn-b06--the-verdict-offers-the-way-out) `PGNHN-B06`

- [Without a declared dialect the verdict is undetermined](camadas/gate.md#pgnhn-b07--without-a-declared-dialect-the-verdict-is-undetermined) `PGNHN-B07`

- [The ruler is agnostic across languages](camadas/gate.md#pgnhn-i01--the-ruler-is-agnostic-across-languages) `PGNHN-I01`

- [Where the construct is not recognised, the gate stays silent](camadas/gate.md#pgnhn-i02--where-the-construct-is-not-recognised-the-gate-stays-silent) `PGNHN-I02`

- [A cursor with no loop does not count as pagination](camadas/gate.md#pgnhn-i03--a-cursor-with-no-loop-does-not-count-as-pagination) `PGNHN-I03`

- [A provider prefix in the name does not hide the promise](camadas/gate.md#pgnhn-i04--a-provider-prefix-in-the-name-does-not-hide-the-promise) `PGNHN-I04`

- [The gate does not invent a cursor the provider does not offer](camadas/gate.md#pgnhn-x01--the-gate-does-not-invent-a-cursor-the-provider-does-not-offer) `PGNHN-X01`

- [The gate does not measure performance or page size](camadas/gate.md#pgnhn-x02--the-gate-does-not-measure-performance-or-page-size) `PGNHN-X02`

- [A source that lives only in the prose is failed](camadas/gate.md#psdpl-b01--a-source-that-lives-only-in-the-prose-is-failed) `PSDPL-B01`

- [The verdict names which source and where its adapter lives](camadas/gate.md#psdpl-b02--the-verdict-names-which-source-and-where-its-adapter-lives) `PSDPL-B02`

- [With the owning plan declared in needs the gate passes](camadas/gate.md#psdpl-b03--with-the-owning-plan-declared-in-needs-the-gate-passes) `PSDPL-B03`

- [Every source of the line is confronted on its own](camadas/gate.md#psdpl-b04--every-source-of-the-line-is-confronted-on-its-own) `PSDPL-B04`

- [The source name matches the adapter regardless of case](camadas/gate.md#psdpl-b05--the-source-name-matches-the-adapter-regardless-of-case) `PSDPL-B05`

- [A source whose adapter nobody seeds is not charged](camadas/gate.md#psdpl-b06--a-source-whose-adapter-nobody-seeds-is-not-charged) `PSDPL-B06`

- [The plan that seeds the adapter is not charged for itself](camadas/gate.md#psdpl-b07--the-plan-that-seeds-the-adapter-is-not-charged-for-itself) `PSDPL-B07`

- [A plan with no source line returns Skip](camadas/gate.md#psdpl-b08--a-plan-with-no-source-line-returns-skip) `PSDPL-B08`

- [An artifact that is not a plan returns Skip](camadas/gate.md#psdpl-b09--an-artifact-that-is-not-a-plan-returns-skip) `PSDPL-B09`

- [What was not measured is never approved](camadas/gate.md#psdpl-i01--what-was-not-measured-is-never-approved) `PSDPL-I01`

- [A seeded file off the naming pattern owns nothing](camadas/gate.md#psdpl-i02--a-seeded-file-off-the-naming-pattern-owns-nothing) `PSDPL-I02`

- [The gate does not confront the order of the phases](camadas/gate.md#psdpl-x01--the-gate-does-not-confront-the-order-of-the-phases) `PSDPL-X01`

- [The gate does not demand a needs pointing at nothing](camadas/gate.md#psdpl-x02--the-gate-does-not-demand-a-needs-pointing-at-nothing) `PSDPL-X02`

- [The gate does not interpret what the source is for](camadas/gate.md#psdpl-x03--the-gate-does-not-interpret-what-the-source-is-for) `PSDPL-X03`

- [An artifact that is not a plan leaves without a verdict](camadas/gate.md#prhnp-b01--an-artifact-that-is-not-a-plan-leaves-without-a-verdict) `PRHNP-B01`

- [A plan with no companion progress file is skipped, not failed](camadas/gate.md#prhnp-b02--a-plan-with-no-companion-progress-file-is-skipped-not-failed) `PRHNP-B02`

- [The skip for a missing companion says how to create it](camadas/gate.md#prhnp-b03--the-skip-for-a-missing-companion-says-how-to-create-it) `PRHNP-B03`

- [A ticked item whose file does not exist is failed](camadas/gate.md#prhnp-b04--a-ticked-item-whose-file-does-not-exist-is-failed) `PRHNP-B04`

- [An open item whose file already exists is failed](camadas/gate.md#prhnp-b05--an-open-item-whose-file-already-exists-is-failed) `PRHNP-B05`

- [A spec the plan seeds and the progress does not list is failed](camadas/gate.md#prhnp-b06--a-spec-the-plan-seeds-and-the-progress-does-not-list-is-failed) `PRHNP-B06`

- [A checkbox item promising no file at all is failed](camadas/gate.md#prhnp-b07--a-checkbox-item-promising-no-file-at-all-is-failed) `PRHNP-B07`

- [The ticked-but-absent finding is reported first](camadas/gate.md#prhnp-b08--the-ticked-but-absent-finding-is-reported-first) `PRHNP-B08`

- [An item in prose citing no path is not charged](camadas/gate.md#prhnp-b09--an-item-in-prose-citing-no-path-is-not-charged) `PRHNP-B09`

- [A progress that agrees with the disk on every item passes](camadas/gate.md#prhnp-b10--a-progress-that-agrees-with-the-disk-on-every-item-passes) `PRHNP-B10`

- [The verdict names each offending path](camadas/gate.md#prhnp-b11--the-verdict-names-each-offending-path) `PRHNP-B11`

- [The companion's path has one definition, derived from the scanner](camadas/gate.md#prhnp-i01--the-companions-path-has-one-definition-derived-from-the-scanner) `PRHNP-I01`

- [A seed is matched by path, never by the item's text](camadas/gate.md#prhnp-i02--a-seed-is-matched-by-path-never-by-the-items-text) `PRHNP-I02`

- [A spec mentioned in the plan's prose is not a seed](camadas/gate.md#prhnp-i03--a-spec-mentioned-in-the-plans-prose-is-not-a-seed) `PRHNP-I03`

- [A template file is never a seeded spec](camadas/gate.md#prhnp-i04--a-template-file-is-never-a-seeded-spec) `PRHNP-I04`

- [The gate does not charge the existence of the progress file](camadas/gate.md#prhnp-x01--the-gate-does-not-charge-the-existence-of-the-progress-file) `PRHNP-X01`

- [The gate does not judge the content of an item beyond the path](camadas/gate.md#prhnp-x02--the-gate-does-not-judge-the-content-of-an-item-beyond-the-path) `PRHNP-X02`

- [The gate does not put the progress file into the map](camadas/gate.md#prhnp-x03--the-gate-does-not-put-the-progress-file-into-the-map) `PRHNP-X03`

- [A spec whose rules the code ignores is accused, and the verdict names them](camadas/gate.md#rlimr-b01--a-spec-whose-rules-the-code-ignores-is-accused-and-the-verdict-names-them) `RLIMR-B01`

- [A rule waived with a written reason closes the account](camadas/gate.md#rlimr-b02--a-rule-waived-with-a-written-reason-closes-the-account) `RLIMR-B02`

- [A unit that predates the practice is a pending item, not a failure](camadas/gate.md#rlimr-b03--a-unit-that-predates-the-practice-is-a-pending-item-not-a-failure) `RLIMR-B03`

- [Declaring the requirement turns the pending item into a failure](camadas/gate.md#rlimr-b04--declaring-the-requirement-turns-the-pending-item-into-a-failure) `RLIMR-B04`

- [A spec with no linked code is not this gate's subject](camadas/gate.md#rlimr-b05--a-spec-with-no-linked-code-is-not-this-gates-subject) `RLIMR-B05`

- [A waiver with no named rule covers every rule of the spec](camadas/gate.md#rlimr-b06--a-waiver-with-no-named-rule-covers-every-rule-of-the-spec) `RLIMR-B06`

- [Requiring the marking never punishes whoever already marks](camadas/gate.md#rlimr-i01--requiring-the-marking-never-punishes-whoever-already-marks) `RLIMR-I01`

- [The identity survives a rename](camadas/gate.md#rlimr-i02--the-identity-survives-a-rename) `RLIMR-I02`

- [The gate does not judge whether the implementation is correct](camadas/gate.md#rlimr-x01--the-gate-does-not-judge-whether-the-implementation-is-correct) `RLIMR-X01`

- [The gate does not demand a mark on EVERY rule](camadas/gate.md#rlimr-x02--the-gate-does-not-demand-a-mark-on-every-rule) `RLIMR-X02`

- [A letter that is not declared in the vocabulary fails](camadas/gate.md#rltyr-b01--a-letter-that-is-not-declared-in-the-vocabulary-fails) `RLTYR-B01`

- [A declared letter under a claimed section passes](camadas/gate.md#rltyr-b02--a-declared-letter-under-a-claimed-section-passes) `RLTYR-B02`

- [A section cataloguing rules under a title no letter claims fails](camadas/gate.md#rltyr-b03--a-section-cataloguing-rules-under-a-title-no-letter-claims-fails) `RLTYR-B03`

- [The same letter claimed by two terms is a conflict in the vocabulary](camadas/gate.md#rltyr-b04--the-same-letter-claimed-by-two-terms-is-a-conflict-in-the-vocabulary) `RLTYR-B04`

- [With no vocabulary declared the gate confronts the canonical letters](camadas/gate.md#rltyr-b05--with-no-vocabulary-declared-the-gate-confronts-the-canonical-letters) `RLTYR-B05`

- [A heading that is the rule code itself is not a category section](camadas/gate.md#rltyr-b06--a-heading-that-is-the-rule-code-itself-is-not-a-category-section) `RLTYR-B06`

- [A section that only cites other sections' codes claims no letter](camadas/gate.md#rltyr-b07--a-section-that-only-cites-other-sections-codes-claims-no-letter) `RLTYR-B07`

- [A section that defines a code in the first table cell is charged](camadas/gate.md#rltyr-b08--a-section-that-defines-a-code-in-the-first-table-cell-is-charged) `RLTYR-B08`

- [A section declared as rule-cataloguing and filled without a code is Pending](camadas/gate.md#rltyr-b09--a-section-declared-as-rule-cataloguing-and-filled-without-a-code-is-pending) `RLTYR-B09`

- [A section whose table already carries the code is not charged](camadas/gate.md#rltyr-b10--a-section-whose-table-already-carries-the-code-is-not-charged) `RLTYR-B10`

- [A declared section outside sections_require_code is not charged](camadas/gate.md#rltyr-b11--a-declared-section-outside-sections-require-code-is-not-charged) `RLTYR-B11`

- [A project that does not use sections_require_code changes no behaviour](camadas/gate.md#rltyr-b12--a-project-that-does-not-use-sections-require-code-changes-no-behaviour) `RLTYR-B12`

- [A spec with no rule code at all is not this gate's problem](camadas/gate.md#rltyr-i01--a-spec-with-no-rule-code-at-all-is-not-this-gates-problem) `RLTYR-I01`

- [The verdict names the letter and where to declare it](camadas/gate.md#rltyr-i02--the-verdict-names-the-letter-and-where-to-declare-it) `RLTYR-I02`

- [A filled section with no code is the gap where the scenario loses its anchor](camadas/gate.md#rltyr-i03--a-filled-section-with-no-code-is-the-gap-where-the-scenario-loses-its-anchor) `RLTYR-I03`

- [The gate does not decide which letters exist](camadas/gate.md#rltyr-x01--the-gate-does-not-decide-which-letters-exist) `RLTYR-X01`

- [Without a declared vocabulary only the letter is charged](camadas/gate.md#rltyr-x02--without-a-declared-vocabulary-only-the-letter-is-charged) `RLTYR-X02`

- [The gate does not judge whether the letter suits the rule](camadas/gate.md#rltyr-x03--the-gate-does-not-judge-whether-the-letter-suits-the-rule) `RLTYR-X03`

- [The gate charges traceability, not format](camadas/gate.md#rltyr-x04--the-gate-charges-traceability-not-format) `RLTYR-X04`

- [A requirement no scenario tags is failed and named](camadas/gate.md#sfmsp-b01--a-requirement-no-scenario-tags-is-failed-and-named) `SFMSP-B01`

- [A requirement that has a scenario is not accused](camadas/gate.md#sfmsp-b02--a-requirement-that-has-a-scenario-is-not-accused) `SFMSP-B02`

- [With every requirement tagged the gate passes](camadas/gate.md#sfmsp-b03--with-every-requirement-tagged-the-gate-passes) `SFMSP-B03`

- [A code merely cited contracts no obligation](camadas/gate.md#sfmsp-b04--a-code-merely-cited-contracts-no-obligation) `SFMSP-B04`

- [A per-requirement waiver with a written reason waives](camadas/gate.md#sfmsp-b05--a-per-requirement-waiver-with-a-written-reason-waives) `SFMSP-B05`

- [A bare per-requirement waiver does not waive](camadas/gate.md#sfmsp-b06--a-bare-per-requirement-waiver-does-not-waive) `SFMSP-B06`

- [A whole-spec waiver drags the waiver to every requirement](camadas/gate.md#sfmsp-b07--a-whole-spec-waiver-drags-the-waiver-to-every-requirement) `SFMSP-B07`

- [Without the whole-spec waiver the same requirements keep failing](camadas/gate.md#sfmsp-b08--without-the-whole-spec-waiver-the-same-requirements-keep-failing) `SFMSP-B08`

- [A bare whole-spec waiver drags nothing](camadas/gate.md#sfmsp-b09--a-bare-whole-spec-waiver-drags-nothing) `SFMSP-B09`

- [A spec with no feature returns Skip](camadas/gate.md#sfmsp-b10--a-spec-with-no-feature-returns-skip) `SFMSP-B10`

- [Requirements are looked for across every linked feature](camadas/gate.md#sfmsp-b11--requirements-are-looked-for-across-every-linked-feature) `SFMSP-B11`

- [An artifact that is not a spec returns Skip](camadas/gate.md#sfmsp-b12--an-artifact-that-is-not-a-spec-returns-skip) `SFMSP-B12`

- [Every waiver requires a written reason](camadas/gate.md#sfmsp-i01--every-waiver-requires-a-written-reason) `SFMSP-I01`

- [Each gate accuses one thing](camadas/gate.md#sfmsp-i02--each-gate-accuses-one-thing) `SFMSP-I02`

- [The gate does not judge whether the scenario proves the requirement](camadas/gate.md#sfmsp-x01--the-gate-does-not-judge-whether-the-scenario-proves-the-requirement) `SFMSP-X01`

- [A cited code produces no accusation](camadas/gate.md#sfmsp-x02--a-cited-code-produces-no-accusation) `SFMSP-X02`

- [The gate does not confront feature against test](camadas/gate.md#sfmsp-x03--the-gate-does-not-confront-feature-against-test) `SFMSP-X03`

- [An artifact that is not a spec leaves without a verdict](camadas/gate.md#trcmt-b01--an-artifact-that-is-not-a-spec-leaves-without-a-verdict) `TRCMT-B01`

- [A recognised layer leaves without a verdict](camadas/gate.md#trcmt-b02--a-recognised-layer-leaves-without-a-verdict) `TRCMT-B02`

- [Without a map the verdict is undetermined](camadas/gate.md#trcmt-b03--without-a-map-the-verdict-is-undetermined) `TRCMT-B03`

- [A spec with the three pieces linked passes](camadas/gate.md#trcmt-b04--a-spec-with-the-three-pieces-linked-passes) `TRCMT-B04`

- [A spec missing a piece is failed, and the verdict names which](camadas/gate.md#trcmt-b05--a-spec-missing-a-piece-is-failed-and-the-verdict-names-which) `TRCMT-B05`

- [The layer may waive a piece for every spec in it](camadas/gate.md#trcmt-b06--the-layer-may-waive-a-piece-for-every-spec-in-it) `TRCMT-B06`

- [The unit may waive a piece in its own spec, with a written reason](camadas/gate.md#trcmt-b07--the-unit-may-waive-a-piece-in-its-own-spec-with-a-written-reason) `TRCMT-B07`

- [The test is reached in two hops, through the feature](camadas/gate.md#trcmt-i01--the-test-is-reached-in-two-hops-through-the-feature) `TRCMT-I01`

- [Waiving the test while the feature carries a scenario is a contradiction](camadas/gate.md#trcmt-i02--waiving-the-test-while-the-feature-carries-a-scenario-is-a-contradiction) `TRCMT-I02`

- [Waiving the test demands saying where the proof is, and the place must exist](camadas/gate.md#trcmt-i03--waiving-the-test-demands-saying-where-the-proof-is-and-the-place-must-exist) `TRCMT-I03`

- [A waiver covers only the piece it declares](camadas/gate.md#trcmt-i04--a-waiver-covers-only-the-piece-it-declares) `TRCMT-I04`

- [The gate does not confront whether the pieces MATCH one another](camadas/gate.md#trcmt-x01--the-gate-does-not-confront-whether-the-pieces-match-one-another) `TRCMT-X01`

- [The gate does not judge the QUALITY of any piece](camadas/gate.md#trcmt-x02--the-gate-does-not-judge-the-quality-of-any-piece) `TRCMT-X02`

- [An artifact that is not a spec leaves without a verdict](camadas/gate.md#vlanv-b01--an-artifact-that-is-not-a-spec-leaves-without-a-verdict) `VLANV-B01`

- [A value of a closed set with no anchor is failed](camadas/gate.md#vlanv-b02--a-value-of-a-closed-set-with-no-anchor-is-failed) `VLANV-B02`

- [The verdict names the unanchored value](camadas/gate.md#vlanv-b03--the-verdict-names-the-unanchored-value) `VLANV-B03`

- [A value whose anchor carries rule key and value passes](camadas/gate.md#vlanv-b04--a-value-whose-anchor-carries-rule-key-and-value-passes) `VLANV-B04`

- [An anchor that asserts one value while the line says another is failed](camadas/gate.md#vlanv-b05--an-anchor-that-asserts-one-value-while-the-line-says-another-is-failed) `VLANV-B05`

- [The verdict of a lying anchor shows both sides of the divergence](camadas/gate.md#vlanv-b06--the-verdict-of-a-lying-anchor-shows-both-sides-of-the-divergence) `VLANV-B06`

- [Lying anchors are reported before the unanchored ones](camadas/gate.md#vlanv-b07--lying-anchors-are-reported-before-the-unanchored-ones) `VLANV-B07`

- [Without a declared value anchor pattern the gate skips](camadas/gate.md#vlanv-b08--without-a-declared-value-anchor-pattern-the-gate-skips) `VLANV-B08`

- [The skip names the setting that enables the gate](camadas/gate.md#vlanv-b09--the-skip-names-the-setting-that-enables-the-gate) `VLANV-B09`

- [A declaration that opens no list is not a closed set](camadas/gate.md#vlanv-b10--a-declaration-that-opens-no-list-is-not-a-closed-set) `VLANV-B10`

- [An anchor pattern with a single capture group does not enable the gate](camadas/gate.md#vlanv-i01--an-anchor-pattern-with-a-single-capture-group-does-not-enable-the-gate) `VLANV-I01`

- [A line carrying an anchor is never read as the end of the list](camadas/gate.md#vlanv-i02--a-line-carrying-an-anchor-is-never-read-as-the-end-of-the-list) `VLANV-I02`

- [With no built map the verdict is pending](camadas/gate.md#vlanv-i03--with-no-built-map-the-verdict-is-pending) `VLANV-I03`

- [The gate does not judge whether the value is a good one](camadas/gate.md#vlanv-x01--the-gate-does-not-judge-whether-the-value-is-a-good-one) `VLANV-X01`

- [The gate does not accuse a line whose literal it cannot read](camadas/gate.md#vlanv-x02--the-gate-does-not-accuse-a-line-whose-literal-it-cannot-read) `VLANV-X02`

- [The gate does not decide what an anchor or a public symbol looks like](camadas/gate.md#vlanv-x03--the-gate-does-not-decide-what-an-anchor-or-a-public-symbol-looks-like) `VLANV-X03`

