<!-- anchors:generated from doct/comportamento.md.tmpl — DO NOT EDIT: run `anchors docs build` -->


# Comportamento

Todos os cenários do sistema. Cada um leva à unidade que o define.

Um cenário descreve o que o sistema faz numa situação — vem da feature, e é o mesmo que o
teste prova.

## gate

- [An identifier in the wrong language is accused, and an English one passes](camadas/gate.md#cdlng-b01--an-identifier-in-the-wrong-language-is-accused-and-an-english-one-passes) `CDLNG-B01`

- [The verdict returns the word that accused](camadas/gate.md#cdlng-b02--the-verdict-returns-the-word-that-accused) `CDLNG-B02`

- [Only a DECLARATION is the subject](camadas/gate.md#cdlng-b03--only-a-declaration-is-the-subject) `CDLNG-B03`

- [The declarations are found in every form the language offers](camadas/gate.md#cdlng-b04--the-declarations-are-found-in-every-form-the-language-offers) `CDLNG-B04`

- [Deciding one word is separate from deciding a whole identifier](camadas/gate.md#cdlng-b05--deciding-one-word-is-separate-from-deciding-a-whole-identifier) `CDLNG-B05`

- [A short word does not count](camadas/gate.md#cdlng-i01--a-short-word-does-not-count) `CDLNG-I01`

- [The gate does not read comments](camadas/gate.md#cdlng-x01--the-gate-does-not-read-comments) `CDLNG-X01`

- [The gate does not read user-facing text](camadas/gate.md#cdlng-x02--the-gate-does-not-read-user-facing-text) `CDLNG-X02`

- [The gate does not use a dictionary to decide the language](camadas/gate.md#cdlng-x03--the-gate-does-not-use-a-dictionary-to-decide-the-language) `CDLNG-X03`

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

