<!-- anchors:generated from doct/comportamento.md.tmpl — inputs:9ce67bf70732a73b — DO NOT EDIT: run `anchors docs build` -->


# Comportamento

Todos os cenários do sistema. Cada um leva à unidade que o define.

Um cenário descreve o que o sistema faz numa situação — vem da feature, e é o mesmo que o
teste prova.

## apoio

- [A card needs every configured label and the asked state](layers/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode) `BRCRB-B01`

- [The owner is the last ownership comment](layers/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode) `BRCRB-B02`

- [An escalated card goes only to whoever decides the product](layers/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode) `BRCRB-B03`

- [The agent resumes its own work under way](layers/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode) `BRCRB-B04`

- [A unit's card is found by the code in its title](layers/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode) `BRCRB-B05`

- [An open Anchors card is returned by number](layers/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode) `BRCRB-B06`

- [Board queries do not name the repository, other calls do](layers/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode) `BRCRB-B07`

- [Asking for work dispatches the claim pipeline](layers/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode) `BRCRB-B08`

- [Empty labels are refused](layers/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode) `BRCRB-E01`

- [A repository that is not owner/name is refused](layers/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode) `BRCRB-E02`

- [An empty code is refused](layers/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode) `BRCRB-E03`

- [A code no title holds is refused](layers/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode) `BRCRB-E04`

- [An ambiguous code is refused naming every card](layers/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode) `BRCRB-E05`

- [A closed card or a non-Anchors issue is refused by number](layers/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode) `BRCRB-E06`

- [An empty comment is refused](layers/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode) `BRCRB-E07`

- [Another agent's pending run does not count](layers/apoio.md#clwtc--claimwait--asking-the-claim-pipeline-for-a-card-and-waiting-bounded-for-the-answer) `CLWTC-B01`

- [No dispatch while the agent's run is pending](layers/apoio.md#clwtc--claimwait--asking-the-claim-pipeline-for-a-card-and-waiting-bounded-for-the-answer) `CLWTC-B02`

- [An older finished run is not the answer](layers/apoio.md#clwtc--claimwait--asking-the-claim-pipeline-for-a-card-and-waiting-bounded-for-the-answer) `CLWTC-B03`

- [The card is returned as soon as it arrives](layers/apoio.md#clwtc--claimwait--asking-the-claim-pipeline-for-a-card-and-waiting-bounded-for-the-answer) `CLWTC-B04`

- [A finished run ends the wait after one more look](layers/apoio.md#clwtc--claimwait--asking-the-claim-pipeline-for-a-card-and-waiting-bounded-for-the-answer) `CLWTC-B05`

- [A cancelled run is asked again, bounded](layers/apoio.md#clwtc--claimwait--asking-the-claim-pipeline-for-a-card-and-waiting-bounded-for-the-answer) `CLWTC-B06`

- [The wait is bounded and does not dispatch twice](layers/apoio.md#clwtc--claimwait--asking-the-claim-pipeline-for-a-card-and-waiting-bounded-for-the-answer) `CLWTC-B07`

- [An unreadable run list fails before dispatching](layers/apoio.md#clwtc--claimwait--asking-the-claim-pipeline-for-a-card-and-waiting-bounded-for-the-answer) `CLWTC-E01`

- [gh gets no input](layers/apoio.md#ghrng--ghrunner--running-gh-for-the-board-so-that-a-failure-names-its-cause-and-never-waits-for-input) `GHRNG-B01`

- [Exit code 4 is explained as not authenticated](layers/apoio.md#ghrng--ghrunner--running-gh-for-the-board-so-that-a-failure-names-its-cause-and-never-waits-for-input) `GHRNG-B02`

- [Other failures get no authentication hint](layers/apoio.md#ghrng--ghrunner--running-gh-for-the-board-so-that-a-failure-names-its-cause-and-never-waits-for-input) `GHRNG-B03`

- [A heading's anchor follows the GitHub convention and keeps accents](layers/apoio.md#dclnd--doclinks--the-anchors-links-sizes-and-layer-arrows-the-documentation-templates-are-given) `DCLND-B01`

- [A rule's link points at the rule on a small layer and at its unit on a big one](layers/apoio.md#dclnd--doclinks--the-anchors-links-sizes-and-layer-arrows-the-documentation-templates-are-given) `DCLND-B02`

- [A scenario's link falls back to the unit, then to the bare page, on a big layer](layers/apoio.md#dclnd--doclinks--the-anchors-links-sizes-and-layer-arrows-the-documentation-templates-are-given) `DCLND-B03`

- [A selection's size counts units, rules, lines and scenarios, and is remembered per selection](layers/apoio.md#dclnd--doclinks--the-anchors-links-sizes-and-layer-arrows-the-documentation-templates-are-given) `DCLND-B04`

- [A selection that matches no spec has size zero and is not big](layers/apoio.md#dclnd--doclinks--the-anchors-links-sizes-and-layer-arrows-the-documentation-templates-are-given) `DCLND-B05`

- [Arrows between layers aggregate dependency edges between specs and carry their count](layers/apoio.md#dclnd--doclinks--the-anchors-links-sizes-and-layer-arrows-the-documentation-templates-are-given) `DCLND-B06`

- [The arrows are ordered by source layer then target layer](layers/apoio.md#dclnd--doclinks--the-anchors-links-sizes-and-layer-arrows-the-documentation-templates-are-given) `DCLND-B07`

- [A diagram node identifier has only ASCII letters, digits and underscores](layers/apoio.md#dclnd--doclinks--the-anchors-links-sizes-and-layer-arrows-the-documentation-templates-are-given) `DCLND-B08`

- [Every generated link resolves on a small layer and on a big one](layers/apoio.md#dclnd--doclinks--the-anchors-links-sizes-and-layer-arrows-the-documentation-templates-are-given) `DCLND-I01`

- [The bigness asked by a template has no threshold of its own](layers/apoio.md#dclnd--doclinks--the-anchors-links-sizes-and-layer-arrows-the-documentation-templates-are-given) `DCLND-X01`

- [A layer's page is in the folder of its template, or the project language's folder](layers/apoio.md#dclnd--doclinks--the-anchors-links-sizes-and-layer-arrows-the-documentation-templates-are-given) `DCLND-B09`

- [A container carries the specs of the layers it declares, in layer then code order](layers/apoio.md#c4cnc--c4containers--the-declared-containers-each-with-the-units-that-run-in-it-for-the-architecture-page) `C4CNC-B01`

- [An external container gets no level 3](layers/apoio.md#c4cnc--c4containers--the-declared-containers-each-with-the-units-that-run-in-it-for-the-architecture-page) `C4CNC-B02`

- [A layer no container declares is named as orphan](layers/apoio.md#c4cnc--c4containers--the-declared-containers-each-with-the-units-that-run-in-it-for-the-architecture-page) `C4CNC-B03`

- [No configuration means no containers and no orphans](layers/apoio.md#c4cnc--c4containers--the-declared-containers-each-with-the-units-that-run-in-it-for-the-architecture-page) `C4CNC-B04`

- [Every spec layer is held by a container or named orphan](layers/apoio.md#c4cnc--c4containers--the-declared-containers-each-with-the-units-that-run-in-it-for-the-architecture-page) `C4CNC-I01`

- [The containers are only the declared ones](layers/apoio.md#c4cnc--c4containers--the-declared-containers-each-with-the-units-that-run-in-it-for-the-architecture-page) `C4CNC-X01`

- [The spec's content enters the compiled page](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B01`

- [The compiled page opens with the generated marker, its template path and a stamp](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B02`

- [A handwritten page is never overwritten](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B03`

- [A dry run writes nothing](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B04`

- [A template in a subfolder is compiled](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B05`

- [The layer comes from the spec's header](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B06`

- [Only layers that have specs are offered](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B07`

- [A section runs to the next heading of its own level](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B08`

- [A spec is split into its separate rules](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B09`

- [A rule's title drops the HTML comment on its heading](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B10`

- [A missing page or a page with a different stamp is stale](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B11`

- [Editing the template makes its page stale](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B12`

- [The header's date is not part of the stamp](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B13`

- [A handwritten page is never reported stale](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B14`

- [The specs no template reaches are uncovered](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B15`

- [Generated and handwritten markers are recognised](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B16`

- [A page with the unhashed marker is stale only when its body differs](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B17`

- [Right after a build nothing is stale](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-I01`

- [Asking for stale pages writes nothing](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-X01`

- [A spec file the map does not list is not loaded](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-X02`

- [A wrong selection filter fails with a message that helps fix it](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-E01`

- [A failing selection fails the build without writing the page](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-E02`

- [A broken template fails the build](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-E03`

- [A spec in the map but not on disk fails the compiler](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-E04`

- [A project without templates fails the build and has nothing stale or uncovered](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-E05`

- [Asking one spec by an unknown code fails](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-E06`

- [A layer with files and no spec is said, not an error](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B18`

- [The compiler reads through its source](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B19`

- [A compiler with a source compiles the templates it holds](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B20`

- [The pages out of date, compiled without writing](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B21`

- [A section asked by its title in one language is found under another](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B22`

- [The rules are read in the three catalogued forms, the spec's own codes at their first definition](layers/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content) `DTCDC-B23`

- [Every known kind has a title, items to cover and a trap](layers/apoio.md#dcknd--dockinds--what-each-kind-of-project-documentation-must-answer-told-to-the-agent-that-writes-it) `DCKND-B01`

- [The kind is found ignoring case and surrounding spaces](layers/apoio.md#dcknd--dockinds--what-each-kind-of-project-documentation-must-answer-told-to-the-agent-that-writes-it) `DCKND-B02`

- [An unknown kind gets a minimal instruction titled by the kind or the path](layers/apoio.md#dcknd--dockinds--what-each-kind-of-project-documentation-must-answer-told-to-the-agent-that-writes-it) `DCKND-B03`

- [The duty text lists title and path, why, items and trap in order](layers/apoio.md#dcknd--dockinds--what-each-kind-of-project-documentation-must-answer-told-to-the-agent-that-writes-it) `DCKND-B04`

- [The known kinds and the instructions are the same set](layers/apoio.md#dcknd--dockinds--what-each-kind-of-project-documentation-must-answer-told-to-the-agent-that-writes-it) `DCKND-I01`

- [No documentation kind is refused](layers/apoio.md#dcknd--dockinds--what-each-kind-of-project-documentation-must-answer-told-to-the-agent-that-writes-it) `DCKND-X01`

- [Either measure passing its cut-off makes the selection big](layers/apoio.md#dcoxx--doclayout--the-one-decision-of-whether-a-documentation-page-shows-everything-or-summarizes) `DCOXX-B01`

- [The default cut-off is 20 units or 2000 lines](layers/apoio.md#dcoxx--doclayout--the-one-decision-of-whether-a-documentation-page-shows-everything-or-summarizes) `DCOXX-B02`

- [The summary sentence exists only for a big selection and names its numbers](layers/apoio.md#dcoxx--doclayout--the-one-decision-of-whether-a-documentation-page-shows-everything-or-summarizes) `DCOXX-B03`

- [A selection exactly at the limit is not big](layers/apoio.md#dcoxx--doclayout--the-one-decision-of-whether-a-documentation-page-shows-everything-or-summarizes) `DCOXX-I01`

- [The layout never splits a layer into pages](layers/apoio.md#dcoxx--doclayout--the-one-decision-of-whether-a-documentation-page-shows-everything-or-summarizes) `DCOXX-X01`

- [The three fixed templates each have a name, a body and what they answer](layers/apoio.md#dcscd--docscaffolds--the-starting-templates-anchors-docs-init-proposes-one-page-per-question-and-per-layer) `DCSCD-B01`

- [A layer's page template is named after the layer](layers/apoio.md#dcscd--docscaffolds--the-starting-templates-anchors-docs-init-proposes-one-page-per-question-and-per-layer) `DCSCD-B02`

- [Init writes the fixed templates and one page per layer, each opening with its purpose](layers/apoio.md#dcscd--docscaffolds--the-starting-templates-anchors-docs-init-proposes-one-page-per-question-and-per-layer) `DCSCD-B03`

- [An edited template is kept unless forced](layers/apoio.md#dcscd--docscaffolds--the-starting-templates-anchors-docs-init-proposes-one-page-per-question-and-per-layer) `DCSCD-B04`

- [The small layer page carries each scenario under its own heading, and the index links it](layers/apoio.md#dcscd--docscaffolds--the-starting-templates-anchors-docs-init-proposes-one-page-per-question-and-per-layer) `DCSCD-B05`

- [The big layer page summarizes each unit without scenario steps](layers/apoio.md#dcscd--docscaffolds--the-starting-templates-anchors-docs-init-proposes-one-page-per-question-and-per-layer) `DCSCD-B06`

- [The architecture page follows the C4 model](layers/apoio.md#dcscd--docscaffolds--the-starting-templates-anchors-docs-init-proposes-one-page-per-question-and-per-layer) `DCSCD-B07`

- [No container declared is said on the architecture page](layers/apoio.md#dcscd--docscaffolds--the-starting-templates-anchors-docs-init-proposes-one-page-per-question-and-per-layer) `DCSCD-B08`

- [The skeleton init writes compiles](layers/apoio.md#dcscd--docscaffolds--the-starting-templates-anchors-docs-init-proposes-one-page-per-question-and-per-layer) `DCSCD-I01`

- [Init writes no compiled page](layers/apoio.md#dcscd--docscaffolds--the-starting-templates-anchors-docs-init-proposes-one-page-per-question-and-per-layer) `DCSCD-X01`

- [A templates folder that cannot be created stops init with the error](layers/apoio.md#dcscd--docscaffolds--the-starting-templates-anchors-docs-init-proposes-one-page-per-question-and-per-layer) `DCSCD-E01`

- [The templates are named and written in the project's language](layers/apoio.md#dcscd--docscaffolds--the-starting-templates-anchors-docs-init-proposes-one-page-per-question-and-per-layer) `DCSCD-B09`

- [A project with API units gets the OpenAPI template](layers/apoio.md#dcscd--docscaffolds--the-starting-templates-anchors-docs-init-proposes-one-page-per-question-and-per-layer) `DCSCD-B05`

- [Scenarios and outlines open in any dialect, and an examples table does not](layers/apoio.md#gsrgh--gherkinscenarioreader--the-scenarios-of-a-units-feature-with-their-steps-for-the-documentation) `GSRGH-B01`

- [The first code-shaped tag is the code and the others are tags](layers/apoio.md#gsrgh--gherkinscenarioreader--the-scenarios-of-a-units-feature-with-their-steps-for-the-documentation) `GSRGH-B02`

- [The body is the scenario's non-blank lines until the next scenario](layers/apoio.md#gsrgh--gherkinscenarioreader--the-scenarios-of-a-units-feature-with-their-steps-for-the-documentation) `GSRGH-B03`

- [A Background heading after a scenario does not leak into it](layers/apoio.md#gsrgh--gherkinscenarioreader--the-scenarios-of-a-units-feature-with-their-steps-for-the-documentation) `GSRGH-B04`

- [The feature is found by the map edge, in either direction, not by name](layers/apoio.md#gsrgh--gherkinscenarioreader--the-scenarios-of-a-units-feature-with-their-steps-for-the-documentation) `GSRGH-B05`

- [The scenarios of a selection follow the spec selection and its errors](layers/apoio.md#gsrgh--gherkinscenarioreader--the-scenarios-of-a-units-feature-with-their-steps-for-the-documentation) `GSRGH-B06`

- [Every scenario read carries its spec's code](layers/apoio.md#gsrgh--gherkinscenarioreader--the-scenarios-of-a-units-feature-with-their-steps-for-the-documentation) `GSRGH-I01`

- [The body is kept verbatim](layers/apoio.md#gsrgh--gherkinscenarioreader--the-scenarios-of-a-units-feature-with-their-steps-for-the-documentation) `GSRGH-X01`

- [A linked feature missing on disk contributes no scenario and no error](layers/apoio.md#gsrgh--gherkinscenarioreader--the-scenarios-of-a-units-feature-with-their-steps-for-the-documentation) `GSRGH-E01`

- [The three house styles say the same thing](layers/apoio.md#flgrf--flaggrammar--the-fixed-grammar-of-a-feature-flag-scenarios-condition) `FLGRF-B01`

- [The operand loses one layer of quotes](layers/apoio.md#flgrf--flaggrammar--the-fixed-grammar-of-a-feature-flag-scenarios-condition) `FLGRF-B02`

- [The absent and present cases compare against nothing, in any language](layers/apoio.md#flgrf--flaggrammar--the-fixed-grammar-of-a-feature-flag-scenarios-condition) `FLGRF-B03`

- [The longest operator wins](layers/apoio.md#flgrf--flaggrammar--the-fixed-grammar-of-a-feature-flag-scenarios-condition) `FLGRF-B04`

- [Only absent and present need no operand](layers/apoio.md#flgrf--flaggrammar--the-fixed-grammar-of-a-feature-flag-scenarios-condition) `FLGRF-B05`

- [Service-dependent operators and prose are refused](layers/apoio.md#flgrf--flaggrammar--the-fixed-grammar-of-a-feature-flag-scenarios-condition) `FLGRF-X01`

- [An empty condition is refused](layers/apoio.md#flgrf--flaggrammar--the-fixed-grammar-of-a-feature-flag-scenarios-condition) `FLGRF-E01`

- [An operator without an operand is refused](layers/apoio.md#flgrf--flaggrammar--the-fixed-grammar-of-a-feature-flag-scenarios-condition) `FLGRF-E02`

- [A loose word is refused with the accepted forms](layers/apoio.md#flgrf--flaggrammar--the-fixed-grammar-of-a-feature-flag-scenarios-condition) `FLGRF-E03`

- [The G rows of the table are the scenarios](layers/apoio.md#flprf--flagparse--reading-a-projects-flag-files-and-the-scenarios-their-values-open) `FLPRF-B01`

- [The flag's code and name](layers/apoio.md#flprf--flagparse--reading-a-projects-flag-files-and-the-scenarios-their-values-open) `FLPRF-B02`

- [Each scenario records its line](layers/apoio.md#flprf--flagparse--reading-a-projects-flag-files-and-the-scenarios-their-values-open) `FLPRF-B03`

- [A refused condition becomes a finding](layers/apoio.md#flprf--flagparse--reading-a-projects-flag-files-and-the-scenarios-their-values-open) `FLPRF-B04`

- [The flag says whether it declares the absent case](layers/apoio.md#flprf--flagparse--reading-a-projects-flag-files-and-the-scenarios-their-values-open) `FLPRF-B05`

- [Flags load in a stable order, and none without a folder](layers/apoio.md#flprf--flagparse--reading-a-projects-flag-files-and-the-scenarios-their-values-open) `FLPRF-B06`

- [Scenarios are indexed by code](layers/apoio.md#flprf--flagparse--reading-a-projects-flag-files-and-the-scenarios-their-values-open) `FLPRF-B07`

- [A flags folder or flag file that cannot be read is an error](layers/apoio.md#flprf--flagparse--reading-a-projects-flag-files-and-the-scenarios-their-values-open) `FLPRF-E02`

- [Actions then flows are read in a stable order, and no flows give no graph](layers/apoio.md#flblf--flowbuild--assembling-the-work-flow-graph-from-the-projects-flow-and-action-files) `FLBLF-B01`

- [Every coded heading is a state of its file](layers/apoio.md#flblf--flowbuild--assembling-the-work-flow-graph-from-the-projects-flow-and-action-files) `FLBLF-B02`

- [Exits belong to the state above them](layers/apoio.md#flblf--flowbuild--assembling-the-work-flow-graph-from-the-projects-flow-and-action-files) `FLBLF-B03`

- [A line with two codes routes a result](layers/apoio.md#flblf--flowbuild--assembling-the-work-flow-graph-from-the-projects-flow-and-action-files) `FLBLF-B04`

- [Terminal is declared, not inferred](layers/apoio.md#flblf--flowbuild--assembling-the-work-flow-graph-from-the-projects-flow-and-action-files) `FLBLF-B05`

- [The piece a step fits is read in any language](layers/apoio.md#flblf--flowbuild--assembling-the-work-flow-graph-from-the-projects-flow-and-action-files) `FLBLF-B06`

- [A suggested reaction is recorded as a suggestion](layers/apoio.md#flblf--flowbuild--assembling-the-work-flow-graph-from-the-projects-flow-and-action-files) `FLBLF-B07`

- [A routed result keeps the prose around its two codes as its condition](layers/apoio.md#flblf--flowbuild--assembling-the-work-flow-graph-from-the-projects-flow-and-action-files) `FLBLF-B08`

- [A flows folder, actions folder or flow file that cannot be read is an error](layers/apoio.md#flblf--flowbuild--assembling-the-work-flow-graph-from-the-projects-flow-and-action-files) `FLBLF-E01`

- [From a state, only its declared exits](layers/apoio.md#flmdf--flowmodel--the-questions-a-work-flow-graph-answers-what-comes-next-and-what-is-broken) `FLMDF-B01`

- [A state is found by its code](layers/apoio.md#flmdf--flowmodel--the-questions-a-work-flow-graph-answers-what-comes-next-and-what-is-broken) `FLMDF-B02`

- [A result is a code whose letter is O](layers/apoio.md#flmdf--flowmodel--the-questions-a-work-flow-graph-answers-what-comes-next-and-what-is-broken) `FLMDF-B05`

- [States keep the file's order and flows are listed once](layers/apoio.md#flmdf--flowmodel--the-questions-a-work-flow-graph-answers-what-comes-next-and-what-is-broken) `FLMDF-B03`

- [The entry is the first declared step, never a result](layers/apoio.md#flmdf--flowmodel--the-questions-a-work-flow-graph-answers-what-comes-next-and-what-is-broken) `FLMDF-B04`

- [An action's title is its file name](layers/apoio.md#flmdf--flowmodel--the-questions-a-work-flow-graph-answers-what-comes-next-and-what-is-broken) `FLMDF-B06`

- [A step nobody arrives at is unreachable](layers/apoio.md#flmdf--flowmodel--the-questions-a-work-flow-graph-answers-what-comes-next-and-what-is-broken) `FLMDF-B07`

- [A result no flow routes is unhandled](layers/apoio.md#flmdf--flowmodel--the-questions-a-work-flow-graph-answers-what-comes-next-and-what-is-broken) `FLMDF-B08`

- [An exit to a state that does not exist is dangling](layers/apoio.md#flmdf--flowmodel--the-questions-a-work-flow-graph-answers-what-comes-next-and-what-is-broken) `FLMDF-B09`

- [A project without flows answers nothing to every question](layers/apoio.md#flmdf--flowmodel--the-questions-a-work-flow-graph-answers-what-comes-next-and-what-is-broken) `FLMDF-X01`

- [The supported languages and the default](layers/apoio.md#incta--i18ncatalog--every-message-a-person-reads-in-the-projects-language-with-a-fallback-that-never-goes-blank) `INCTA-B01`

- [An empty language is the default](layers/apoio.md#incta--i18ncatalog--every-message-a-person-reads-in-the-projects-language-with-a-fallback-that-never-goes-blank) `INCTA-B02`

- [Messages come out in the current language with their arguments](layers/apoio.md#incta--i18ncatalog--every-message-a-person-reads-in-the-projects-language-with-a-fallback-that-never-goes-blank) `INCTA-B03`

- [A key missing in the current language falls back to English](layers/apoio.md#incta--i18ncatalog--every-message-a-person-reads-in-the-projects-language-with-a-fallback-that-never-goes-blank) `INCTA-B04`

- [A key missing everywhere resolves to itself](layers/apoio.md#incta--i18ncatalog--every-message-a-person-reads-in-the-projects-language-with-a-fallback-that-never-goes-blank) `INCTA-B05`

- [A key's values across languages come back once each](layers/apoio.md#incta--i18ncatalog--every-message-a-person-reads-in-the-projects-language-with-a-fallback-that-never-goes-blank) `INCTA-B06`

- [A key resolves in a given language without changing the current one](layers/apoio.md#incta--i18ncatalog--every-message-a-person-reads-in-the-projects-language-with-a-fallback-that-never-goes-blank) `INCTA-B07`

- [A written title gives back its key and language](layers/apoio.md#incta--i18ncatalog--every-message-a-person-reads-in-the-projects-language-with-a-fallback-that-never-goes-blank) `INCTA-B08`

- [Every language has a catalog with the same keys](layers/apoio.md#incta--i18ncatalog--every-message-a-person-reads-in-the-projects-language-with-a-fallback-that-never-goes-blank) `INCTA-I01`

- [An unsupported language is refused with the options](layers/apoio.md#incta--i18ncatalog--every-message-a-person-reads-in-the-projects-language-with-a-fallback-that-never-goes-blank) `INCTA-E01`

- [Without a declared log path nothing is scanned](layers/apoio.md#lgscl--logscan--finding-the-occurrences-of-declared-failures-in-the-projects-logs) `LGSCL-B01`

- [The format does not matter](layers/apoio.md#lgscl--logscan--finding-the-occurrences-of-declared-failures-in-the-projects-logs) `LGSCL-B02`

- [Look-alikes of a code are not captured](layers/apoio.md#lgscl--logscan--finding-the-occurrences-of-declared-failures-in-the-projects-logs) `LGSCL-B03`

- [An undeclared failure is reported apart](layers/apoio.md#lgscl--logscan--finding-the-occurrences-of-declared-failures-in-the-projects-logs) `LGSCL-B04`

- [The delimiter removes the ambiguity](layers/apoio.md#lgscl--logscan--finding-the-occurrences-of-declared-failures-in-the-projects-logs) `LGSCL-B05`

- [The alias only serves lines that carry no code](layers/apoio.md#lgscl--logscan--finding-the-occurrences-of-declared-failures-in-the-projects-logs) `LGSCL-B06`

- [Each occurrence records its time window](layers/apoio.md#lgscl--logscan--finding-the-occurrences-of-declared-failures-in-the-projects-logs) `LGSCL-B07`

- [Occurrences are sorted by rule](layers/apoio.md#lgscl--logscan--finding-the-occurrences-of-declared-failures-in-the-projects-logs) `LGSCL-B08`

- [A spec's failure rules are listed in every form](layers/apoio.md#lgscl--logscan--finding-the-occurrences-of-declared-failures-in-the-projects-logs) `LGSCL-B09`

- [An invalid alias or timestamp pattern fails the scan](layers/apoio.md#lgscl--logscan--finding-the-occurrences-of-declared-failures-in-the-projects-logs) `LGSCL-E01`

- [The renames of the steps crossed](layers/apoio.md#mgltr--codeletters--rewriting-the-letter-of-the-codes-of-one-kind-of-unit) `MGLTR-B01`

- [A phase, a step and a result get their new letters](layers/apoio.md#mgltr--codeletters--rewriting-the-letter-of-the-codes-of-one-kind-of-unit) `MGLTR-B02`

- [Every other code stays](layers/apoio.md#mgltr--codeletters--rewriting-the-letter-of-the-codes-of-one-kind-of-unit) `MGLTR-B03`

- [The rewrite says what it replaced](layers/apoio.md#mgltr--codeletters--rewriting-the-letter-of-the-codes-of-one-kind-of-unit) `MGLTR-B04`

- [Rewriting twice changes nothing more](layers/apoio.md#mgltr--codeletters--rewriting-the-letter-of-the-codes-of-one-kind-of-unit) `MGLTR-B05`

- [The format is the top-level version line, and its absence means format 1](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-B01`

- [A file already at the target is left untouched](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-B02`

- [Old map keys are renamed where they are keys and keep their values](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-B03`

- [Each rename applies only to its own file](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-B04`

- [A value is renamed only under its key and when it matches whole](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-B05`

- [Steps are applied in order, each on the result of the previous](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-B06`

- [The version is raised even when nothing else changes](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-B07`

- [A map or a configuration without a version gets one after its comment header](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-B08`

- [The result counts each rename by its old form](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-B09`

- [A dry run reports without writing](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-B10`

- [A version line with a trailing comment is read and raised keeping the comment](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-B11`

- [Migrating twice changes nothing the second time](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-I01`

- [A file the YAML parser would refuse is still migrated](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-X01`

- [A missing file is an error and nothing is written](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-E01`

- [A version that is not a number is an error and nothing is written](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-E02`

- [A hole in the chain leaves the file untouched](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-E03`

- [A CRLF file is read and migrated as it is](layers/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key) `MGFLM-B12`

- [Format 2 renames the Portuguese keys in their own files](layers/apoio.md#mgstm--migrationsteps--the-registered-steps-one-per-format-take-any-project-from-format-1-to-the-current-format) `MGSTM-B01`

- [Format 2 renames the Portuguese gate names wherever they are stored](layers/apoio.md#mgstm--migrationsteps--the-registered-steps-one-per-format-take-any-project-from-format-1-to-the-current-format) `MGSTM-B02`

- [Format 3 renames the eight remaining gate names in the configuration](layers/apoio.md#mgstm--migrationsteps--the-registered-steps-one-per-format-take-any-project-from-format-1-to-the-current-format) `MGSTM-B03`

- [Format 5 renames code letters, not keys](layers/apoio.md#mgstm--migrationsteps--the-registered-steps-one-per-format-take-any-project-from-format-1-to-the-current-format) `MGSTM-B05`

- [Format 4 renames the four keys that lied about what they hold](layers/apoio.md#mgstm--migrationsteps--the-registered-steps-one-per-format-take-any-project-from-format-1-to-the-current-format) `MGSTM-B04`

- [The chain from format 1 to the current format has no hole](layers/apoio.md#mgstm--migrationsteps--the-registered-steps-one-per-format-take-any-project-from-format-1-to-the-current-format) `MGSTM-I01`

- [A key renamed by two formats ends under its latest name](layers/apoio.md#mgstm--migrationsteps--the-registered-steps-one-per-format-take-any-project-from-format-1-to-the-current-format) `MGSTM-I02`

- [Format 3 leaves the map's gate and the configuration's id alone](layers/apoio.md#mgstm--migrationsteps--the-registered-steps-one-per-format-take-any-project-from-format-1-to-the-current-format) `MGSTM-X01`

- [Format 4 renames nothing in the map](layers/apoio.md#mgstm--migrationsteps--the-registered-steps-one-per-format-take-any-project-from-format-1-to-the-current-format) `MGSTM-X02`

- [Format 6 renames the triad to the unit](layers/apoio.md#mgstm--migrationsteps--the-registered-steps-one-per-format-take-any-project-from-format-1-to-the-current-format) `MGSTM-B06`

- [Steps registered out of order are kept sorted](layers/apoio.md#mscmg--migrationstepchain--one-step-per-format-version-chained-in-order-and-a-hole-in-the-chain-is-an-error) `MSCMG-B01`

- [The steps between two formats come in ascending order](layers/apoio.md#mscmg--migrationstepchain--one-step-per-format-version-chained-in-order-and-a-hole-in-the-chain-is-an-error) `MSCMG-B02`

- [A file already at the target needs no step](layers/apoio.md#mscmg--migrationstepchain--one-step-per-format-version-chained-in-order-and-a-hole-in-the-chain-is-an-error) `MSCMG-B03`

- [A step may rename code letters by kind](layers/apoio.md#mscmg--migrationstepchain--one-step-per-format-version-chained-in-order-and-a-hole-in-the-chain-is-an-error) `MSCMG-B05`

- [A key some step renames is reported as renamed](layers/apoio.md#mscmg--migrationstepchain--one-step-per-format-version-chained-in-order-and-a-hole-in-the-chain-is-an-error) `MSCMG-B04`

- [The listed steps are a copy of the registry](layers/apoio.md#mscmg--migrationstepchain--one-step-per-format-version-chained-in-order-and-a-hole-in-the-chain-is-an-error) `MSCMG-I01`

- [Only the steps inside the interval are returned](layers/apoio.md#mscmg--migrationstepchain--one-step-per-format-version-chained-in-order-and-a-hole-in-the-chain-is-an-error) `MSCMG-X01`

- [A hole in the chain is an error naming the missing format](layers/apoio.md#mscmg--migrationstepchain--one-step-per-format-version-chained-in-order-and-a-hole-in-the-chain-is-an-error) `MSCMG-E01`

- [A target beyond the last step is an error naming the first missing format](layers/apoio.md#mscmg--migrationstepchain--one-step-per-format-version-chained-in-order-and-a-hole-in-the-chain-is-an-error) `MSCMG-E02`

- [Every known role presents itself](layers/apoio.md#agrlg--agentroles--who-is-who-in-a-project-and-the-capabilities-each-role-carries) `AGRLG-B01`

- [Only the product owner and the architect decide the product](layers/apoio.md#agrlg--agentroles--who-is-who-in-a-project-and-the-capabilities-each-role-carries) `AGRLG-B02`

- [The structure is the architect's](layers/apoio.md#agrlg--agentroles--who-is-who-in-a-project-and-the-capabilities-each-role-carries) `AGRLG-B03`

- [QA, reviewers and dev execute and do not decide](layers/apoio.md#agrlg--agentroles--who-is-who-in-a-project-and-the-capabilities-each-role-carries) `AGRLG-B04`

- [Each reviewing role has its own lens](layers/apoio.md#agrlg--agentroles--who-is-who-in-a-project-and-the-capabilities-each-role-carries) `AGRLG-B05`

- [Typed roles are recognised by name and abbreviation](layers/apoio.md#agrlg--agentroles--who-is-who-in-a-project-and-the-capabilities-each-role-carries) `AGRLG-B06`

- [A role lists its capabilities in alphabetical order](layers/apoio.md#agrlg--agentroles--who-is-who-in-a-project-and-the-capabilities-each-role-carries) `AGRLG-B07`

- [An unknown role can do nothing](layers/apoio.md#agrlg--agentroles--who-is-who-in-a-project-and-the-capabilities-each-role-carries) `AGRLG-X01`

- [The settings live in the local state folder](layers/apoio.md#ussts--usersettings--the-agents-local-settings-the-declared-role-kept-out-of-git) `USSTS-B01`

- [A missing settings file is not an error](layers/apoio.md#ussts--usersettings--the-agents-local-settings-the-declared-role-kept-out-of-git) `USSTS-B02`

- [The legacy field has three states](layers/apoio.md#ussts--usersettings--the-agents-local-settings-the-declared-role-kept-out-of-git) `USSTS-B03`

- [The legacy yes grants only the product decision](layers/apoio.md#ussts--usersettings--the-agents-local-settings-the-declared-role-kept-out-of-git) `USSTS-B04`

- [A declared role wins over the legacy field](layers/apoio.md#ussts--usersettings--the-agents-local-settings-the-declared-role-kept-out-of-git) `USSTS-B05`

- [What is saved is what is loaded](layers/apoio.md#ussts--usersettings--the-agents-local-settings-the-declared-role-kept-out-of-git) `USSTS-B06`

- [The saved file explains itself](layers/apoio.md#ussts--usersettings--the-agents-local-settings-the-declared-role-kept-out-of-git) `USSTS-B07`

- [Typed answers are read in both languages](layers/apoio.md#ussts--usersettings--the-agents-local-settings-the-declared-role-kept-out-of-git) `USSTS-B08`

- [The description names the role or asks for one](layers/apoio.md#ussts--usersettings--the-agents-local-settings-the-declared-role-kept-out-of-git) `USSTS-B09`

- [A settings file that is not YAML fails naming it](layers/apoio.md#ussts--usersettings--the-agents-local-settings-the-declared-role-kept-out-of-git) `USSTS-E01`

- [Every form someone would try turns telemetry off](layers/apoio.md#tlcnt--telemetryconfig--the-opt-out-is-easy-to-find-easy-to-get-right-and-the-environment-overrides-the-file) `TLCNT-B01`

- [With nothing declared telemetry is on](layers/apoio.md#tlcnt--telemetryconfig--the-opt-out-is-easy-to-find-easy-to-get-right-and-the-environment-overrides-the-file) `TLCNT-B02`

- [The environment overrides the file in both directions](layers/apoio.md#tlcnt--telemetryconfig--the-opt-out-is-easy-to-find-easy-to-get-right-and-the-environment-overrides-the-file) `TLCNT-B03`

- [Without the environment the file decides](layers/apoio.md#tlcnt--telemetryconfig--the-opt-out-is-easy-to-find-easy-to-get-right-and-the-environment-overrides-the-file) `TLCNT-B04`

- [A value means the same in the environment and in the file](layers/apoio.md#tlcnt--telemetryconfig--the-opt-out-is-easy-to-find-easy-to-get-right-and-the-environment-overrides-the-file) `TLCNT-I01`

- [A word outside the closed list keeps telemetry on](layers/apoio.md#tlcnt--telemetryconfig--the-opt-out-is-easy-to-find-easy-to-get-right-and-the-environment-overrides-the-file) `TLCNT-X01`

- [A new event keeps its name and attributes](layers/apoio.md#tlevt--telemetryevent--a-decision-event-has-a-name-from-a-closed-vocabulary-a-caller-stamped-instant-and-attributes) `TLEVT-B01`

- [An event without attributes carries an empty set](layers/apoio.md#tlevt--telemetryevent--a-decision-event-has-a-name-from-a-closed-vocabulary-a-caller-stamped-instant-and-attributes) `TLEVT-B02`

- [The instant comes from the caller's clock](layers/apoio.md#tlevt--telemetryevent--a-decision-event-has-a-name-from-a-closed-vocabulary-a-caller-stamped-instant-and-attributes) `TLEVT-B03`

- [The vocabulary is the five declared names](layers/apoio.md#tlevt--telemetryevent--a-decision-event-has-a-name-from-a-closed-vocabulary-a-caller-stamped-instant-and-attributes) `TLEVT-I01`

- [The system clock is never read](layers/apoio.md#tlevt--telemetryevent--a-decision-event-has-a-name-from-a-closed-vocabulary-a-caller-stamped-instant-and-attributes) `TLEVT-X01`

- [The notice is shown once per project root](layers/apoio.md#tlntt--telemetrynotice--the-telemetry-notice-reaches-whoever-did-not-ask-for-it-once-and-says-how-to-turn-it-off) `TLNTT-B01`

- [The text says what is sent, what is not, and how to turn it off](layers/apoio.md#tlntt--telemetrynotice--the-telemetry-notice-reaches-whoever-did-not-ask-for-it-once-and-says-how-to-turn-it-off) `TLNTT-B02`

- [The marker is written under the project's unversioned anchors directory](layers/apoio.md#tlntt--telemetrynotice--the-telemetry-notice-reaches-whoever-did-not-ask-for-it-once-and-says-how-to-turn-it-off) `TLNTT-B03`

- [The notice is written in the project's language](layers/apoio.md#tlntt--telemetrynotice--the-telemetry-notice-reaches-whoever-did-not-ask-for-it-once-and-says-how-to-turn-it-off) `TLNTT-B04`

- [After the notice is written it counts as already shown](layers/apoio.md#tlntt--telemetrynotice--the-telemetry-notice-reaches-whoever-did-not-ask-for-it-once-and-says-how-to-turn-it-off) `TLNTT-I01`

- [The notice does not consult the opt-out itself](layers/apoio.md#tlntt--telemetrynotice--the-telemetry-notice-reaches-whoever-did-not-ask-for-it-once-and-says-how-to-turn-it-off) `TLNTT-X01`

- [An unwritable marker never blocks and the notice shows again](layers/apoio.md#tlntt--telemetrynotice--the-telemetry-notice-reaches-whoever-did-not-ask-for-it-once-and-says-how-to-turn-it-off) `TLNTT-E01`

- [A disabled configuration yields no emitter that still accepts calls](layers/apoio.md#tlemt--telemetryemitter--decision-events-leave-as-otlp-logs-never-block-the-work-and-never-carry-who-uses-the-product) `TLEMT-B01`

- [A blank endpoint means Honeycomb's logs endpoint](layers/apoio.md#tlemt--telemetryemitter--decision-events-leave-as-otlp-logs-never-block-the-work-and-never-carry-who-uses-the-product) `TLEMT-B02`

- [An event is sent as one log record whose body is the event name](layers/apoio.md#tlemt--telemetryemitter--decision-events-leave-as-otlp-logs-never-block-the-work-and-never-carry-who-uses-the-product) `TLEMT-B03`

- [Emitting to a dead collector neither blocks nor fails](layers/apoio.md#tlemt--telemetryemitter--decision-events-leave-as-otlp-logs-never-block-the-work-and-never-carry-who-uses-the-product) `TLEMT-B04`

- [Flushing waits for the send in flight](layers/apoio.md#tlemt--telemetryemitter--decision-events-leave-as-otlp-logs-never-block-the-work-and-never-carry-who-uses-the-product) `TLEMT-B05`

- [Flushing has its own deadline independent of the HTTP client](layers/apoio.md#tlemt--telemetryemitter--decision-events-leave-as-otlp-logs-never-block-the-work-and-never-carry-who-uses-the-product) `TLEMT-B06`

- [Each request is a JSON POST carrying the configured headers](layers/apoio.md#tlemt--telemetryemitter--decision-events-leave-as-otlp-logs-never-block-the-work-and-never-carry-who-uses-the-product) `TLEMT-B07`

- [With NoCodes a unit or rule code never leaves in an attribute](layers/apoio.md#tlemt--telemetryemitter--decision-events-leave-as-otlp-logs-never-block-the-work-and-never-carry-who-uses-the-product) `TLEMT-B08`

- [A value that is not vocabulary becomes its type description](layers/apoio.md#tlemt--telemetryemitter--decision-events-leave-as-otlp-logs-never-block-the-work-and-never-carry-who-uses-the-product) `TLEMT-I01`

- [The resource carries only the product name and version](layers/apoio.md#tlemt--telemetryemitter--decision-events-leave-as-otlp-logs-never-block-the-work-and-never-carry-who-uses-the-product) `TLEMT-X01`

- [An error status from the collector is discarded](layers/apoio.md#tlemt--telemetryemitter--decision-events-leave-as-otlp-logs-never-block-the-work-and-never-carry-who-uses-the-product) `TLEMT-E02`

## comando

- [Without an agent name, a workflow or the tracker client there are no cards](layers/comando.md#agcrg--agentcards--the-cards-this-agent-owns-and-the-card-a-pull-request-declares-read-from-the-tracker) `AGCRG-B01`

- [The tracker is asked for the open cards of the repository owned by this agent](layers/comando.md#agcrg--agentcards--the-cards-this-agent-owns-and-the-card-a-pull-request-declares-read-from-the-tracker) `AGCRG-B02`

- [A card line becomes a card whose state loses the anchors prefix](layers/comando.md#agcrg--agentcards--the-cards-this-agent-owns-and-the-card-a-pull-request-declares-read-from-the-tracker) `AGCRG-B03`

- [Malformed card lines are skipped](layers/comando.md#agcrg--agentcards--the-cards-this-agent-owns-and-the-card-a-pull-request-declares-read-from-the-tracker) `AGCRG-B04`

- [A card with no state kept as the last line](layers/comando.md#agcrg--agentcards--the-cards-this-agent-owns-and-the-card-a-pull-request-declares-read-from-the-tracker) `AGCRG-B05`

- [The issue title keeps the first line of the reason](layers/comando.md#agcrg--agentcards--the-cards-this-agent-owns-and-the-card-a-pull-request-declares-read-from-the-tracker) `AGCRG-B06`

- [A long title is cut to 70 with an ellipsis](layers/comando.md#agcrg--agentcards--the-cards-this-agent-owns-and-the-card-a-pull-request-declares-read-from-the-tracker) `AGCRG-B07`

- [The issue number is read from the last segment of the URL](layers/comando.md#agcrg--agentcards--the-cards-this-agent-owns-and-the-card-a-pull-request-declares-read-from-the-tracker) `AGCRG-B08`

- [The card of a pull request is the first closing keyword at a line start](layers/comando.md#agcrg--agentcards--the-cards-this-agent-owns-and-the-card-a-pull-request-declares-read-from-the-tracker) `AGCRG-B09`

- [The pull request reference is normalised before asking the tracker](layers/comando.md#agcrg--agentcards--the-cards-this-agent-owns-and-the-card-a-pull-request-declares-read-from-the-tracker) `AGCRG-B10`

- [Every returned card has a number and a state without prefix](layers/comando.md#agcrg--agentcards--the-cards-this-agent-owns-and-the-card-a-pull-request-declares-read-from-the-tracker) `AGCRG-I01`

- [Ownership is selected by the question sent to the tracker](layers/comando.md#agcrg--agentcards--the-cards-this-agent-owns-and-the-card-a-pull-request-declares-read-from-the-tracker) `AGCRG-X01`

- [A failing tracker client gives no cards](layers/comando.md#agcrg--agentcards--the-cards-this-agent-owns-and-the-card-a-pull-request-declares-read-from-the-tracker) `AGCRG-E01`

- [A failing tracker client gives no pull request card](layers/comando.md#agcrg--agentcards--the-cards-this-agent-owns-and-the-card-a-pull-request-declares-read-from-the-tracker) `AGCRG-E02`

- [The not-governed error names the path](layers/comando.md#cmclc--commoncli--the-contract-every-command-shares-how-a-path-becomes-a-node-how-not-governed-is-signalled-what-the-binary-says-it-is) `CMCLC-B01`

- [The not-governed exit code is 3](layers/comando.md#cmclc--commoncli--the-contract-every-command-shares-how-a-path-becomes-a-node-how-not-governed-is-signalled-what-the-binary-says-it-is) `CMCLC-B02`

- [A root-relative path is kept, cleaned](layers/comando.md#cmclc--commoncli--the-contract-every-command-shares-how-a-path-becomes-a-node-how-not-governed-is-signalled-what-the-binary-says-it-is) `CMCLC-B03`

- [A path absent under the root is resolved from the working directory](layers/comando.md#cmclc--commoncli--the-contract-every-command-shares-how-a-path-becomes-a-node-how-not-governed-is-signalled-what-the-binary-says-it-is) `CMCLC-B04`

- [A node exists only by its exact identifier](layers/comando.md#cmclc--commoncli--the-contract-every-command-shares-how-a-path-becomes-a-node-how-not-governed-is-signalled-what-the-binary-says-it-is) `CMCLC-B05`

- [A task slug drops only the last extension](layers/comando.md#cmclc--commoncli--the-contract-every-command-shares-how-a-path-becomes-a-node-how-not-governed-is-signalled-what-the-binary-says-it-is) `CMCLC-B06`

- [An unstamped build says it is a development build](layers/comando.md#cmclc--commoncli--the-contract-every-command-shares-how-a-path-becomes-a-node-how-not-governed-is-signalled-what-the-binary-says-it-is) `CMCLC-B07`

- [Root-relative and absolute names of one file resolve to one node](layers/comando.md#cmclc--commoncli--the-contract-every-command-shares-how-a-path-becomes-a-node-how-not-governed-is-signalled-what-the-binary-says-it-is) `CMCLC-I01`

- [A path that cannot be related to the root is kept as given](layers/comando.md#cmclc--commoncli--the-contract-every-command-shares-how-a-path-becomes-a-node-how-not-governed-is-signalled-what-the-binary-says-it-is) `CMCLC-E01`

- [A list of files is read the same way by every command](layers/comando.md#cmclc--commoncli--the-contract-every-command-shares-how-a-path-becomes-a-node-how-not-governed-is-signalled-what-the-binary-says-it-is) `CMCLC-B08`

- [A value passed under the old name reaches the current flag](layers/comando.md#flalf--flagaliases--a-renamed-command-flag-keeps-answering-to-its-old-name) `FLALF-B01`

- [The alias of a switch flag is a switch](layers/comando.md#flalf--flagaliases--a-renamed-command-flag-keeps-answering-to-its-old-name) `FLALF-B02`

- [The old name is hidden and deprecated](layers/comando.md#flalf--flagaliases--a-renamed-command-flag-keeps-answering-to-its-old-name) `FLALF-B03`

- [The current name wins when both are passed](layers/comando.md#flalf--flagaliases--a-renamed-command-flag-keeps-answering-to-its-old-name) `FLALF-B04`

- [A pair naming flags the command does not have is ignored](layers/comando.md#flalf--flagaliases--a-renamed-command-flag-keeps-answering-to-its-old-name) `FLALF-B05`

- [Nothing is copied before the aliases are resolved](layers/comando.md#flalf--flagaliases--a-renamed-command-flag-keeps-answering-to-its-old-name) `FLALF-X01`

- [A value the current flag cannot hold names the old flag](layers/comando.md#flalf--flagaliases--a-renamed-command-flag-keeps-answering-to-its-old-name) `FLALF-E01`

- [An alias of a flag the command does not have stops the program](layers/comando.md#flalf--flagaliases--a-renamed-command-flag-keeps-answering-to-its-old-name) `FLALF-E02`

- [The agent identity falls back from the session to the user to default](layers/comando.md#arcgn--agentrolecli--who-this-agent-is-and-the-role-it-declared-as-the-commands-show-and-ask-it) `ARCGN-B01`

- [The role list shows every known role](layers/comando.md#arcgn--agentrolecli--who-this-agent-is-and-the-role-it-declared-as-the-commands-show-and-ask-it) `ARCGN-B02`

- [Showing a role says whether it decides the product and shows its lens](layers/comando.md#arcgn--agentrolecli--who-this-agent-is-and-the-role-it-declared-as-the-commands-show-and-ask-it) `ARCGN-B03`

- [An unrecognised answer is echoed back and the question is asked again](layers/comando.md#arcgn--agentrolecli--who-this-agent-is-and-the-role-it-declared-as-the-commands-show-and-ask-it) `ARCGN-B04`

- [Three unrecognised answers give up](layers/comando.md#arcgn--agentrolecli--who-this-agent-is-and-the-role-it-declared-as-the-commands-show-and-ask-it) `ARCGN-B05`

- [A pipe or the null device is not an interactive terminal](layers/comando.md#arcgn--agentrolecli--who-this-agent-is-and-the-role-it-declared-as-the-commands-show-and-ask-it) `ARCGN-B06`

- [Only a role that handles escalated cards decides the product](layers/comando.md#arcgn--agentrolecli--who-this-agent-is-and-the-role-it-declared-as-the-commands-show-and-ask-it) `ARCGN-B07`

- [The role list is the settings catalogue, not a copy](layers/comando.md#arcgn--agentrolecli--who-this-agent-is-and-the-role-it-declared-as-the-commands-show-and-ask-it) `ARCGN-X01`

- [A closed input while asking is an error](layers/comando.md#arcgn--agentrolecli--who-this-agent-is-and-the-role-it-declared-as-the-commands-show-and-ask-it) `ARCGN-E01`

- [An unreadable settings file does not unlock the product decision](layers/comando.md#arcgn--agentrolecli--who-this-agent-is-and-the-role-it-declared-as-the-commands-show-and-ask-it) `ARCGN-E02`

- [The header code is read in any comment style](layers/comando.md#uncdn--unitcodes--the-identity-code-of-a-unit-read-from-a-header-from-the-map-or-from-the-codes-a-file-names) `UNCDN-B01`

- [A code of the wrong length or case is not read](layers/comando.md#uncdn--unitcodes--the-identity-code-of-a-unit-read-from-a-header-from-the-map-or-from-the-codes-a-file-names) `UNCDN-B02`

- [The exact map node answers with its code](layers/comando.md#uncdn--unitcodes--the-identity-code-of-a-unit-read-from-a-header-from-the-map-or-from-the-codes-a-file-names) `UNCDN-B03`

- [A node of the same stem answers when the exact node has no code](layers/comando.md#uncdn--unitcodes--the-identity-code-of-a-unit-read-from-a-header-from-the-map-or-from-the-codes-a-file-names) `UNCDN-B04`

- [Without a map a unit has no code](layers/comando.md#uncdn--unitcodes--the-identity-code-of-a-unit-read-from-a-header-from-the-map-or-from-the-codes-a-file-names) `UNCDN-B05`

- [The codes of a file are the unit's own, once each, in order](layers/comando.md#uncdn--unitcodes--the-identity-code-of-a-unit-read-from-a-header-from-the-map-or-from-the-codes-a-file-names) `UNCDN-B06`

- [No listed code belongs to another unit or repeats](layers/comando.md#uncdn--unitcodes--the-identity-code-of-a-unit-read-from-a-header-from-the-map-or-from-the-codes-a-file-names) `UNCDN-I01`

- [A file that cannot be read is an error](layers/comando.md#uncdn--unitcodes--the-identity-code-of-a-unit-read-from-a-header-from-the-map-or-from-the-codes-a-file-names) `UNCDN-E01`

- [A data state written bare is the unit's own code](layers/comando.md#uncdn--unitcodes--the-identity-code-of-a-unit-read-from-a-header-from-the-map-or-from-the-codes-a-file-names) `UNCDN-B07`

- [The project opt-out is read from the project root](layers/comando.md#tlstt--telemetrysetup--every-command-starts-telemetry-the-same-way-the-opt-outs-first-then-the-notice-then-the-emitter) `TLSTT-B01`

- [The environment opt-out builds nothing](layers/comando.md#tlstt--telemetrysetup--every-command-starts-telemetry-the-same-way-the-opt-outs-first-then-the-notice-then-the-emitter) `TLSTT-B02`

- [The notice is shown once, then telemetry runs quietly](layers/comando.md#tlstt--telemetrysetup--every-command-starts-telemetry-the-same-way-the-opt-outs-first-then-the-notice-then-the-emitter) `TLSTT-B03`

- [The authentication header comes only from the environment](layers/comando.md#tlstt--telemetrysetup--every-command-starts-telemetry-the-same-way-the-opt-outs-first-then-the-notice-then-the-emitter) `TLSTT-B04`

- [The project root is the explicit root or the nearest configured directory above](layers/comando.md#tlstt--telemetrysetup--every-command-starts-telemetry-the-same-way-the-opt-outs-first-then-the-notice-then-the-emitter) `TLSTT-B05`

- [Waiting at exit with no emitter returns at once](layers/comando.md#tlstt--telemetrysetup--every-command-starts-telemetry-the-same-way-the-opt-outs-first-then-the-notice-then-the-emitter) `TLSTT-B06`

- [A declared opt-out holds when the configuration does not load](layers/comando.md#tlstt--telemetrysetup--every-command-starts-telemetry-the-same-way-the-opt-outs-first-then-the-notice-then-the-emitter) `TLSTT-B07`

- [Outside a project there is no notice, no emitter and no mark](layers/comando.md#tlstt--telemetrysetup--every-command-starts-telemetry-the-same-way-the-opt-outs-first-then-the-notice-then-the-emitter) `TLSTT-B08`

- [No emitter without the notice](layers/comando.md#tlstt--telemetrysetup--every-command-starts-telemetry-the-same-way-the-opt-outs-first-then-the-notice-then-the-emitter) `TLSTT-I01`

- [Starting telemetry sends nothing](layers/comando.md#tlstt--telemetrysetup--every-command-starts-telemetry-the-same-way-the-opt-outs-first-then-the-notice-then-the-emitter) `TLSTT-X01`

- [backfill-labels in local mode is refused](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-B01`

- [The open decisions are read in one listing](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-B02`

- [Every decision under an origin card holds it](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-B03`

- [The provenance is recovered from a cited pull request that declares a card](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-B04`

- [Only an open origin card receives the block](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-B05`

- [Precedence between decisions is not inferred](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-B06`

- [A blocked-by label already present is not written again](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-B07`

- [A label is created before the card is edited with it](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-B08`

- [The dry run names what it would do and touches nothing](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-B09`

- [The output counts what was recovered, written and skipped](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-B10`

- [A decision under itself is never blocked by itself](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-I01`

- [Nothing is removed and nothing unsupported is written](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-X01`

- [A failed listing fails the command](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-E01`

- [An unreadable listing fails the command](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-E02`

- [A failed edit is reported and skipped](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-E03`

- [An origin card whose state cannot be read is skipped](layers/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies) `BCLBB-E04`

- [decided without a card is refused](layers/comando.md#dcdde--decided--release-the-card-an-escalation-stopped-for-a-person-once-the-decision-became-a-rule) `DCDDE-B01`

- [decided without a resolution is refused and says why](layers/comando.md#dcdde--decided--release-the-card-an-escalation-stopped-for-a-person-once-the-decision-became-a-rule) `DCDDE-B02`

- [decided in local mode is refused](layers/comando.md#dcdde--decided--release-the-card-an-escalation-stopped-for-a-person-once-the-decision-became-a-rule) `DCDDE-B03`

- [decided refuses while an unblock card is open](layers/comando.md#dcdde--decided--release-the-card-an-escalation-stopped-for-a-person-once-the-decision-became-a-rule) `DCDDE-B04`

- [The release removes needs-user and every blocked-by label at once](layers/comando.md#dcdde--decided--release-the-card-an-escalation-stopped-for-a-person-once-the-decision-became-a-rule) `DCDDE-B05`

- [The comment keeps the resolution and who blocked the card](layers/comando.md#dcdde--decided--release-the-card-an-escalation-stopped-for-a-person-once-the-decision-became-a-rule) `DCDDE-B06`

- [The decisions under the card are labelled manual and closed](layers/comando.md#dcdde--decided--release-the-card-an-escalation-stopped-for-a-person-once-the-decision-became-a-rule) `DCDDE-B07`

- [The output counts the decisions closed](layers/comando.md#dcdde--decided--release-the-card-an-escalation-stopped-for-a-person-once-the-decision-became-a-rule) `DCDDE-B08`

- [Blockers are read first and the card is released before its decisions close](layers/comando.md#dcdde--decided--release-the-card-an-escalation-stopped-for-a-person-once-the-decision-became-a-rule) `DCDDE-I01`

- [Only the decisions under this card are closed](layers/comando.md#dcdde--decided--release-the-card-an-escalation-stopped-for-a-person-once-the-decision-became-a-rule) `DCDDE-X01`

- [An unblock lookup that fails or is unreadable refuses](layers/comando.md#dcdde--decided--release-the-card-an-escalation-stopped-for-a-person-once-the-decision-became-a-rule) `DCDDE-E01`

- [A failed release fails the command](layers/comando.md#dcdde--decided--release-the-card-an-escalation-stopped-for-a-person-once-the-decision-became-a-rule) `DCDDE-E02`

- [An unreadable decision list closes nothing](layers/comando.md#dcdde--decided--release-the-card-an-escalation-stopped-for-a-person-once-the-decision-became-a-rule) `DCDDE-E03`

- [Unreadable labels report no blocker](layers/comando.md#dcdde--decided--release-the-card-an-escalation-stopped-for-a-person-once-the-decision-became-a-rule) `DCDDE-E04`

- [A delivery missing its required parts is refused](layers/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review) `DLVRE-B01`

- [The first delivery of a unit is accepted by its spec](layers/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review) `DLVRE-B02`

- [Prose with commas is one decision, and files split on commas](layers/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review) `DLVRE-B03`

- [Local mode writes the record under changes](layers/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review) `DLVRE-B04`

- [Github mode records on the card given](layers/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review) `DLVRE-B05`

- [Github mode finds the card by the unit's code](layers/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review) `DLVRE-B06`

- [A vendored pipeline in github mode records nothing](layers/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review) `DLVRE-B07`

- [The next step is the review of the unit](layers/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review) `DLVRE-B08`

- [The watcher hint appears only when the watcher is not running](layers/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review) `DLVRE-B09`

- [The zero-decisions note appears only when nothing was declared](layers/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review) `DLVRE-B10`

- [A refused delivery records nothing](layers/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review) `DLVRE-I01`

- [Github mode never falls back to changes](layers/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review) `DLVRE-X01`

- [Github mode without a code for the unit fails with the way out](layers/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review) `DLVRE-E01`

- [Github mode with no open card for the code fails with the way out](layers/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review) `DLVRE-E02`

- [A closed card given with --card is refused](layers/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review) `DLVRE-E03`

- [The confrontation never blocks the delivery](layers/comando.md#dlcnd--deliveryconfront--confront-what-a-delivery-declares-against-the-disk-at-the-moment-it-is-declared) `DLCND-B01`

- [Without git the output says the confrontation did not happen](layers/comando.md#dlcnd--deliveryconfront--confront-what-a-delivery-declares-against-the-disk-at-the-moment-it-is-declared) `DLCND-B02`

- [A committed and untouched declared file is accused](layers/comando.md#dlcnd--deliveryconfront--confront-what-a-delivery-declares-against-the-disk-at-the-moment-it-is-declared) `DLCND-B03`

- [Modified files and files in a new directory are not accused](layers/comando.md#dlcnd--deliveryconfront--confront-what-a-delivery-declares-against-the-disk-at-the-moment-it-is-declared) `DLCND-B04`

- [A renamed file and a file under a subdirectory root are not accused](layers/comando.md#dlcnd--deliveryconfront--confront-what-a-delivery-declares-against-the-disk-at-the-moment-it-is-declared) `DLCND-B08`

- [A tested unit without a mutation signal is warned](layers/comando.md#dlcnd--deliveryconfront--confront-what-a-delivery-declares-against-the-disk-at-the-moment-it-is-declared) `DLCND-B05`

- [A failing informative gate is listed by its first sentence](layers/comando.md#dlcnd--deliveryconfront--confront-what-a-delivery-declares-against-the-disk-at-the-moment-it-is-declared) `DLCND-B06`

- [Without gates, map or delivered node no gate is listed](layers/comando.md#dlcnd--deliveryconfront--confront-what-a-delivery-declares-against-the-disk-at-the-moment-it-is-declared) `DLCND-B07`

- [Not looking and finding nothing are different answers](layers/comando.md#dlcnd--deliveryconfront--confront-what-a-delivery-declares-against-the-disk-at-the-moment-it-is-declared) `DLCND-I01`

- [The confrontation writes nothing](layers/comando.md#dlcnd--deliveryconfront--confront-what-a-delivery-declares-against-the-disk-at-the-moment-it-is-declared) `DLCND-X01`

- [Discard without any card argument is refused](layers/comando.md#dscrd--discard--take-off-the-board-a-card-that-no-longer-makes-sense-without-deleting-it) `DSCRD-B01`

- [Discard with a blank reason is refused](layers/comando.md#dscrd--discard--take-off-the-board-a-card-that-no-longer-makes-sense-without-deleting-it) `DSCRD-B02`

- [Discard in local mode is refused](layers/comando.md#dscrd--discard--take-off-the-board-a-card-that-no-longer-makes-sense-without-deleting-it) `DSCRD-B03`

- [The discard label is created before any card is touched, and an existing label does not stop it](layers/comando.md#dscrd--discard--take-off-the-board-a-card-that-no-longer-makes-sense-without-deleting-it) `DSCRD-B04`

- [A card is labelled, then commented, then closed](layers/comando.md#dscrd--discard--take-off-the-board-a-card-that-no-longer-makes-sense-without-deleting-it) `DSCRD-B05`

- [The comment carries the reason and the way back](layers/comando.md#dscrd--discard--take-off-the-board-a-card-that-no-longer-makes-sense-without-deleting-it) `DSCRD-B06`

- [The leading hash is stripped and a blank argument is skipped](layers/comando.md#dscrd--discard--take-off-the-board-a-card-that-no-longer-makes-sense-without-deleting-it) `DSCRD-B07`

- [A card that is already closed is still discarded](layers/comando.md#dscrd--discard--take-off-the-board-a-card-that-no-longer-makes-sense-without-deleting-it) `DSCRD-B08`

- [Each discarded card is reported on standard output](layers/comando.md#dscrd--discard--take-off-the-board-a-card-that-no-longer-makes-sense-without-deleting-it) `DSCRD-B09`

- [A card is never closed before its label and reason are recorded](layers/comando.md#dscrd--discard--take-off-the-board-a-card-that-no-longer-makes-sense-without-deleting-it) `DSCRD-I01`

- [Discarding never deletes the issue](layers/comando.md#dscrd--discard--take-off-the-board-a-card-that-no-longer-makes-sense-without-deleting-it) `DSCRD-X01`

- [A card whose label fails is named in the error while the others are discarded](layers/comando.md#dscrd--discard--take-off-the-board-a-card-that-no-longer-makes-sense-without-deleting-it) `DSCRD-E01`

- [A card whose reason cannot be recorded is not closed](layers/comando.md#dscrd--discard--take-off-the-board-a-card-that-no-longer-makes-sense-without-deleting-it) `DSCRD-E02`

- [escalate without a reason is refused](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B01`

- [Contradictory exits are refused before any call](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B02`

- [escalate in local mode is refused](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B03`

- [The card of the reviewed pull request is the origin](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B04`

- [The one card in hand is the origin, and two are ambiguous](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B05`

- [A finding with no origin is created and warned](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B06`

- [The title is the exit's prefix and the reason's first line](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B07`

- [The labels follow the exit](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B08`

- [A card written with its hash is the same card](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B17`

- [The origin card and the reviewed pull request are both labels](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B09`

- [The body says why and how to go on for each exit](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B10`

- [A target with open cards is warned about, three at most](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B11`

- [The output names the kind and the address](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B12`

- [A decision stops the origin card and prints the way back](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B13`

- [A blocking bug holds the card by blocked-by alone](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B14`

- [An ordinary card and a non-blocking bug let the origin card go on](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B15`

- [The new card's number is read only from an all-digit address tail](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B16`

- [Without an origin card no other card is touched](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-I01`

- [A bug never carries needs-user](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-X01`

- [A card the platform refuses to create fails the command](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-E01`

- [A card that cannot be stopped is warned about](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-E02`

- [A bug in Anchors is reported to Anchors, once per title](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B18`

- [In local mode the bug goes to Anchors alone](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-B19`

- [The report to Anchors carries nothing of the project](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-X02`

- [A refused report to Anchors leaves a link to file it](layers/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must) `SCLTE-E03`

- [Without a target or a label the board is not asked](layers/comando.md#esdps--escalateduplicate--find-the-open-cards-that-already-deal-with-the-target-of-an-escalation) `ESDPS-B01`

- [The board is searched for open cards with the label and the target](layers/comando.md#esdps--escalateduplicate--find-the-open-cards-that-already-deal-with-the-target-of-an-escalation) `ESDPS-B02`

- [A hit counts when the exact target is in its title or its body](layers/comando.md#esdps--escalateduplicate--find-the-open-cards-that-already-deal-with-the-target-of-an-escalation) `ESDPS-B03`

- [Each card found carries its number and its title](layers/comando.md#esdps--escalateduplicate--find-the-open-cards-that-already-deal-with-the-target-of-an-escalation) `ESDPS-B04`

- [A card about a file that merely contains the target's name is not reported](layers/comando.md#esdps--escalateduplicate--find-the-open-cards-that-already-deal-with-the-target-of-an-escalation) `ESDPS-I01`

- [The lookup only reads the board](layers/comando.md#esdps--escalateduplicate--find-the-open-cards-that-already-deal-with-the-target-of-an-escalation) `ESDPS-X01`

- [A failed or unreadable lookup yields nothing](layers/comando.md#esdps--escalateduplicate--find-the-open-cards-that-already-deal-with-the-target-of-an-escalation) `ESDPS-E01`

- [Local mode reads notifications.md at the project root](layers/comando.md#ntfct--notifications--a-message-to-every-agent-read-from-one-file-and-printed-on-top-of-next) `NTFCT-B01`

- [A file with only comments or nothing prints nothing](layers/comando.md#ntfct--notifications--a-message-to-every-agent-read-from-one-file-and-printed-on-top-of-next) `NTFCT-B02`

- [Github mode reads the file raw from the integration branch](layers/comando.md#ntfct--notifications--a-message-to-every-agent-read-from-one-file-and-printed-on-top-of-next) `NTFCT-B03`

- [A file missing on the platform is silence](layers/comando.md#ntfct--notifications--a-message-to-every-agent-read-from-one-file-and-printed-on-top-of-next) `NTFCT-B04`

- [The block names the file and its source and indents the message](layers/comando.md#ntfct--notifications--a-message-to-every-agent-read-from-one-file-and-printed-on-top-of-next) `NTFCT-B05`

- [The explanatory comment is never printed](layers/comando.md#ntfct--notifications--a-message-to-every-agent-read-from-one-file-and-printed-on-top-of-next) `NTFCT-I01`

- [Github mode does not read the agent's checkout](layers/comando.md#ntfct--notifications--a-message-to-every-agent-read-from-one-file-and-printed-on-top-of-next) `NTFCT-X01`

- [A failed read is reported as one line](layers/comando.md#ntfct--notifications--a-message-to-every-agent-read-from-one-file-and-printed-on-top-of-next) `NTFCT-E01`

- [pr-body in local mode is refused](layers/comando.md#prbdp--prbody--write-the-lines-that-link-a-pull-request-to-its-cards-in-the-platforms-syntax) `PRBDP-B01`

- [The link syntax links and never closes](layers/comando.md#prbdp--prbody--write-the-lines-that-link-a-pull-request-to-its-cards-in-the-platforms-syntax) `PRBDP-B02`

- [The requested cards accept the forms a person writes](layers/comando.md#prbdp--prbody--write-the-lines-that-link-a-pull-request-to-its-cards-in-the-platforms-syntax) `PRBDP-B03`

- [Without requested cards the agent's own card is linked](layers/comando.md#prbdp--prbody--write-the-lines-that-link-a-pull-request-to-its-cards-in-the-platforms-syntax) `PRBDP-B04`

- [pr-body with no card from either source is refused](layers/comando.md#prbdp--prbody--write-the-lines-that-link-a-pull-request-to-its-cards-in-the-platforms-syntax) `PRBDP-B05`

- [Each root drags the open findings born under it](layers/comando.md#prbdp--prbody--write-the-lines-that-link-a-pull-request-to-its-cards-in-the-platforms-syntax) `PRBDP-B06`

- [The lines come out in numeric order](layers/comando.md#prbdp--prbody--write-the-lines-that-link-a-pull-request-to-its-cards-in-the-platforms-syntax) `PRBDP-B07`

- [The only-under switch leaves the roots out](layers/comando.md#prbdp--prbody--write-the-lines-that-link-a-pull-request-to-its-cards-in-the-platforms-syntax) `PRBDP-B08`

- [A card that is both a root and a finding is linked once](layers/comando.md#prbdp--prbody--write-the-lines-that-link-a-pull-request-to-its-cards-in-the-platforms-syntax) `PRBDP-I01`

- [pr-body writes nothing to the platform](layers/comando.md#prbdp--prbody--write-the-lines-that-link-a-pull-request-to-its-cards-in-the-platforms-syntax) `PRBDP-X01`

- [A failed findings lookup contributes nothing and the root is still linked](layers/comando.md#prbdp--prbody--write-the-lines-that-link-a-pull-request-to-its-cards-in-the-platforms-syntax) `PRBDP-E01`

- [The progress file sits beside the plan](layers/comando.md#plprp--planprogress--create-a-plans-progress-file-the-state-that-lives-beside-the-decision-and-outside-the-map) `PLPRP-B01`

- [One section per phase declared in the plan's headers](layers/comando.md#plprp--planprogress--create-a-plans-progress-file-the-state-that-lives-beside-the-decision-and-outside-the-map) `PLPRP-B02`

- [The phase code length follows the project's configuration](layers/comando.md#plprp--planprogress--create-a-plans-progress-file-the-state-that-lives-beside-the-decision-and-outside-the-map) `PLPRP-B03`

- [A plan with no phases gets a note saying what to add](layers/comando.md#plprp--planprogress--create-a-plans-progress-file-the-state-that-lives-beside-the-decision-and-outside-the-map) `PLPRP-B04`

- [An existing progress file is never overwritten](layers/comando.md#plprp--planprogress--create-a-plans-progress-file-the-state-that-lives-beside-the-decision-and-outside-the-map) `PLPRP-B05`

- [new progress creates the file for an existing plan](layers/comando.md#plprp--planprogress--create-a-plans-progress-file-the-state-that-lives-beside-the-decision-and-outside-the-map) `PLPRP-B06`

- [new progress without the plan is refused](layers/comando.md#plprp--planprogress--create-a-plans-progress-file-the-state-that-lives-beside-the-decision-and-outside-the-map) `PLPRP-B07`

- [The suffix is the one the scanner keeps out of the map](layers/comando.md#plprp--planprogress--create-a-plans-progress-file-the-state-that-lives-beside-the-decision-and-outside-the-map) `PLPRP-I01`

- [The progress takes its code from the plan's header](layers/comando.md#plprp--planprogress--create-a-plans-progress-file-the-state-that-lives-beside-the-decision-and-outside-the-map) `PLPRP-X01`

- [A plan without a code is refused](layers/comando.md#plprp--planprogress--create-a-plans-progress-file-the-state-that-lives-beside-the-decision-and-outside-the-map) `PLPRP-E01`

- [A plan that cannot be read is refused](layers/comando.md#plprp--planprogress--create-a-plans-progress-file-the-state-that-lives-beside-the-decision-and-outside-the-map) `PLPRP-E02`

- [The merge driver refuses any count other than three arguments](layers/comando.md#prmrp--progressmerge--the-git-merge-driver-of-a-plans-progress-file-unite-both-sides-done-beats-pending) `PRMRP-B01`

- [The union is written into our side](layers/comando.md#prmrp--progressmerge--the-git-merge-driver-of-a-plans-progress-file-unite-both-sides-done-beats-pending) `PRMRP-B02`

- [A missing base does not stop the merge](layers/comando.md#prmrp--progressmerge--the-git-merge-driver-of-a-plans-progress-file-unite-both-sides-done-beats-pending) `PRMRP-B03`

- [The report counts the done items and what the other side added](layers/comando.md#prmrp--progressmerge--the-git-merge-driver-of-a-plans-progress-file-unite-both-sides-done-beats-pending) `PRMRP-B04`

- [An item done on our side stays done](layers/comando.md#prmrp--progressmerge--the-git-merge-driver-of-a-plans-progress-file-unite-both-sides-done-beats-pending) `PRMRP-I01`

- [The driver applies the progress reader's union](layers/comando.md#prmrp--progressmerge--the-git-merge-driver-of-a-plans-progress-file-unite-both-sides-done-beats-pending) `PRMRP-X01`

- [An unreadable our side is named](layers/comando.md#prmrp--progressmerge--the-git-merge-driver-of-a-plans-progress-file-unite-both-sides-done-beats-pending) `PRMRP-E01`

- [An unreadable other side is named and our side is untouched](layers/comando.md#prmrp--progressmerge--the-git-merge-driver-of-a-plans-progress-file-unite-both-sides-done-beats-pending) `PRMRP-E02`

- [A result that cannot be written fails the merge](layers/comando.md#prmrp--progressmerge--the-git-merge-driver-of-a-plans-progress-file-unite-both-sides-done-beats-pending) `PRMRP-E03`

- [queue lists the live tasks with the hygiene hints](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B01`

- [next claims the next task in local mode](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B02`

- [A cold start seeds a plan that still has work](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B03`

- [A cited spec exists by path or by a unique name](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B04`

- [The seed count is recomputed when printed](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B05`

- [done closes by id and in batch](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B06`

- [drop deletes a task without archiving it](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B07`

- [reclaim respects a live worker unless forced](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B08`

- [The default worker is pid at host](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B09`

- [The board is not claimed without a session](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B10`

- [The board identity falls back to the OS user, and says so](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B11`

- [The agent's own card is resumed before asking the pipeline](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B12`

- [A claim without a card names the run to follow](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B13`

- [The claimed card is printed with its state and owner](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B14`

- [The deliverable follows the card's title](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B15`

- [The documentation duties follow the unit](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B16`

- [Only an agent that does not decide the product is told to escalate](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B17`

- [A card under review asks for a verdict](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B18`

- [Other cards end naming the pull request body command](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B19`

- [The role is asked at most once, and never without a terminal](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-B20`

- [A cold start seeds one plan](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-I01`

- [queue claims nothing](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-X01`

- [The board is not claimed without a repository](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-X02`

- [A failed claim run is an error](layers/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board) `WRQUW-E01`

- [The root holds exactly the seventeen flow commands after registration](layers/comando.md#flrgf--flowregister--attach-the-flow-domains-commands-to-the-root-command-once-each) `FLRGF-B01`

- [The watcher's controls live under watch](layers/comando.md#flrgf--flowregister--attach-the-flow-domains-commands-to-the-root-command-once-each) `FLRGF-B02`

- [No flow command is attached twice](layers/comando.md#flrgf--flowregister--attach-the-flow-domains-commands-to-the-root-command-once-each) `FLRGF-I01`

- [The progress command is not attached to the root](layers/comando.md#flrgf--flowregister--attach-the-flow-domains-commands-to-the-root-command-once-each) `FLRGF-X01`

- [The report has the bug form's sections](layers/comando.md#rpbug--reportbug--report-a-bug-in-anchors-itself-where-it-is-fixed-for-every-project) `RPBUG-B01`

- [An open issue with the same title is told it was seen again](layers/comando.md#rpbug--reportbug--report-a-bug-in-anchors-itself-where-it-is-fixed-for-every-project) `RPBUG-B02`

- [A new bug becomes an issue at co2-lab/anchors](layers/comando.md#rpbug--reportbug--report-a-bug-in-anchors-itself-where-it-is-fixed-for-every-project) `RPBUG-B03`

- [A dry run prints the issue and sends nothing](layers/comando.md#rpbug--reportbug--report-a-bug-in-anchors-itself-where-it-is-fixed-for-every-project) `RPBUG-B04`

- [The minimal case comes from a file or from standard input](layers/comando.md#rpbug--reportbug--report-a-bug-in-anchors-itself-where-it-is-fixed-for-every-project) `RPBUG-B05`

- [It works without an anchors.yaml](layers/comando.md#rpbug--reportbug--report-a-bug-in-anchors-itself-where-it-is-fixed-for-every-project) `RPBUG-B06`

- [An escalation whose reason names the project is not sent upstream](layers/comando.md#rpbug--reportbug--report-a-bug-in-anchors-itself-where-it-is-fixed-for-every-project) `RPBUG-B07`

- [A report that names the project is refused](layers/comando.md#rpbug--reportbug--report-a-bug-in-anchors-itself-where-it-is-fixed-for-every-project) `RPBUG-X01`

- [Without what should happen the report is refused](layers/comando.md#rpbug--reportbug--report-a-bug-in-anchors-itself-where-it-is-fixed-for-every-project) `RPBUG-E01`

- [A refused report leaves the link to file it](layers/comando.md#rpbug--reportbug--report-a-bug-in-anchors-itself-where-it-is-fixed-for-every-project) `RPBUG-E02`

- [An unreadable minimal case is refused](layers/comando.md#rpbug--reportbug--report-a-bug-in-anchors-itself-where-it-is-fixed-for-every-project) `RPBUG-E03`

- [The working tree is read from the root](layers/comando.md#tsstt--taskstatus--discover-what-the-machine-knows-about-the-task-at-hand-so-the-agents-report-does-not-have-to) `TSSTT-B01`

- [A card given by number carries its state, and a closed card reads closed](layers/comando.md#tsstt--taskstatus--discover-what-the-machine-knows-about-the-task-at-hand-so-the-agents-report-does-not-have-to) `TSSTT-B02`

- [Without a number the agent's own card is found on the board](layers/comando.md#tsstt--taskstatus--discover-what-the-machine-knows-about-the-task-at-hand-so-the-agents-report-does-not-have-to) `TSSTT-B03`

- [The decisions waiting on a person are listed](layers/comando.md#tsstt--taskstatus--discover-what-the-machine-knows-about-the-task-at-hand-so-the-agents-report-does-not-have-to) `TSSTT-B04`

- [Running checks are running, not failed](layers/comando.md#tsstt--taskstatus--discover-what-the-machine-knows-about-the-task-at-hand-so-the-agents-report-does-not-have-to) `TSSTT-B05`

- [Only the lock's own reversal counts, reduced to a clean first line](layers/comando.md#tsstt--taskstatus--discover-what-the-machine-knows-about-the-task-at-hand-so-the-agents-report-does-not-have-to) `TSSTT-B06`

- [The command prints the report](layers/comando.md#tsstt--taskstatus--discover-what-the-machine-knows-about-the-task-at-hand-so-the-agents-report-does-not-have-to) `TSSTT-B07`

- [The turn-ended event carries numbers and vocabulary only](layers/comando.md#tsstt--taskstatus--discover-what-the-machine-knows-about-the-task-at-hand-so-the-agents-report-does-not-have-to) `TSSTT-B08`

- [The check classes add up to the total](layers/comando.md#tsstt--taskstatus--discover-what-the-machine-knows-about-the-task-at-hand-so-the-agents-report-does-not-have-to) `TSSTT-I01`

- [Every lookup names the configured repository and the root's branch](layers/comando.md#tsstt--taskstatus--discover-what-the-machine-knows-about-the-task-at-hand-so-the-agents-report-does-not-have-to) `TSSTT-X01`

- [Local mode looks up no card](layers/comando.md#tsstt--taskstatus--discover-what-the-machine-knows-about-the-task-at-hand-so-the-agents-report-does-not-have-to) `TSSTT-X02`

- [Failed and unreadable lookups leave the part absent](layers/comando.md#tsstt--taskstatus--discover-what-the-machine-knows-about-the-task-at-hand-so-the-agents-report-does-not-have-to) `TSSTT-E01`

- [The sections come in the order that decides](layers/comando.md#tsrts--taskstatusreport--the-format-of-the-rounds-report-where-the-task-is-the-verdict-what-is-missing-and-what-comes-next) `TSRTS-B01`

- [The card line drops the code and translates the state](layers/comando.md#tsrts--taskstatusreport--the-format-of-the-rounds-report-where-the-task-is-the-verdict-what-is-missing-and-what-comes-next) `TSRTS-B02`

- [A reversion appears before the verdict, with the way to authorise it](layers/comando.md#tsrts--taskstatusreport--the-format-of-the-rounds-report-where-the-task-is-the-verdict-what-is-missing-and-what-comes-next) `TSRTS-B03`

- [The pull request line never hides missing or mixed checks](layers/comando.md#tsrts--taskstatusreport--the-format-of-the-rounds-report-where-the-task-is-the-verdict-what-is-missing-and-what-comes-next) `TSRTS-B04`

- [The working-tree line says what is not yet shared](layers/comando.md#tsrts--taskstatusreport--the-format-of-the-rounds-report-where-the-task-is-the-verdict-what-is-missing-and-what-comes-next) `TSRTS-B05`

- [The decisions waiting on a person have their own section](layers/comando.md#tsrts--taskstatusreport--the-format-of-the-rounds-report-where-the-task-is-the-verdict-what-is-missing-and-what-comes-next) `TSRTS-B06`

- [The two gaps are always explicit](layers/comando.md#tsrts--taskstatusreport--the-format-of-the-rounds-report-where-the-task-is-the-verdict-what-is-missing-and-what-comes-next) `TSRTS-B07`

- [Uncommitted and unpushed work come first](layers/comando.md#tsrts--taskstatusreport--the-format-of-the-rounds-report-where-the-task-is-the-verdict-what-is-missing-and-what-comes-next) `TSRTS-B08`

- [Without a pull request the step depends on where the card is](layers/comando.md#tsrts--taskstatusreport--the-format-of-the-rounds-report-where-the-task-is-the-verdict-what-is-missing-and-what-comes-next) `TSRTS-B09`

- [The check classes are printed in the user's language](layers/comando.md#tsrts--taskstatusreport--the-format-of-the-rounds-report-where-the-task-is-the-verdict-what-is-missing-and-what-comes-next) `TSRTS-B12`

- [With an open pull request the step follows the checks](layers/comando.md#tsrts--taskstatusreport--the-format-of-the-rounds-report-where-the-task-is-the-verdict-what-is-missing-and-what-comes-next) `TSRTS-B10`

- [A closed card and an idle state have their steps](layers/comando.md#tsrts--taskstatusreport--the-format-of-the-rounds-report-where-the-task-is-the-verdict-what-is-missing-and-what-comes-next) `TSRTS-B11`

- [A running check never reads as passed](layers/comando.md#tsrts--taskstatusreport--the-format-of-the-rounds-report-where-the-task-is-the-verdict-what-is-missing-and-what-comes-next) `TSRTS-I01`

- [Rendering looks nothing up](layers/comando.md#tsrts--taskstatusreport--the-format-of-the-rounds-report-where-the-task-is-the-verdict-what-is-missing-and-what-comes-next) `TSRTS-X01`

- [Unblock with no card or with two cards is refused](layers/comando.md#nblck--unblock--open-the-work-card-a-decision-demanded-linked-to-the-card-stuck-waiting-for-a-person) `NBLCK-B01`

- [Unblock with a blank reason is refused](layers/comando.md#nblck--unblock--open-the-work-card-a-decision-demanded-linked-to-the-card-stuck-waiting-for-a-person) `NBLCK-B02`

- [Unblock in local mode is refused](layers/comando.md#nblck--unblock--open-the-work-card-a-decision-demanded-linked-to-the-card-stuck-waiting-for-a-person) `NBLCK-B03`

- [The link label is created on demand, and an existing one does not stop the command](layers/comando.md#nblck--unblock--open-the-work-card-a-decision-demanded-linked-to-the-card-stuck-waiting-for-a-person) `NBLCK-B04`

- [The new card carries the unblock title and the three labels](layers/comando.md#nblck--unblock--open-the-work-card-a-decision-demanded-linked-to-the-card-stuck-waiting-for-a-person) `NBLCK-B05`

- [The body says what to do, where it came from and how the blocked card returns](layers/comando.md#nblck--unblock--open-the-work-card-a-decision-demanded-linked-to-the-card-stuck-waiting-for-a-person) `NBLCK-B06`

- [The blocked card is told which card it waits for](layers/comando.md#nblck--unblock--open-the-work-card-a-decision-demanded-linked-to-the-card-stuck-waiting-for-a-person) `NBLCK-B07`

- [The output gives the new card and the state of the blocked one](layers/comando.md#nblck--unblock--open-the-work-card-a-decision-demanded-linked-to-the-card-stuck-waiting-for-a-person) `NBLCK-B08`

- [The leading hash of the card is stripped](layers/comando.md#nblck--unblock--open-the-work-card-a-decision-demanded-linked-to-the-card-stuck-waiting-for-a-person) `NBLCK-B09`

- [The blocked card keeps its needs-user label](layers/comando.md#nblck--unblock--open-the-work-card-a-decision-demanded-linked-to-the-card-stuck-waiting-for-a-person) `NBLCK-I01`

- [The temporary body file does not survive the command](layers/comando.md#nblck--unblock--open-the-work-card-a-decision-demanded-linked-to-the-card-stuck-waiting-for-a-person) `NBLCK-X01`

- [A card the platform refuses to create fails the command](layers/comando.md#nblck--unblock--open-the-work-card-a-decision-demanded-linked-to-the-card-stuck-waiting-for-a-person) `NBLCK-E01`

- [A failed comment on the blocked card is a warning](layers/comando.md#nblck--unblock--open-the-work-card-a-decision-demanded-linked-to-the-card-stuck-waiting-for-a-person) `NBLCK-E02`

- [status tells stopped, running with its metadata, and paused](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B01`

- [A second start is refused](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B02`

- [pause needs a running watcher, and resume undoes it](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B03`

- [stop terminates the watcher once](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B04`

- [logs prints the log as is](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B05`

- [run names the configuration or the map it could not load](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B06`

- [Files created after the start become tasks until the signal](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B07`

- [The tree is watched without ignored folders, and new folders are swept](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B08`

- [A governed change queues the next missing piece](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B09`

- [A waived piece is skipped](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B10`

- [What is not work is not queued](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B11`

- [A delivery record triggers the review when the unit closes](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B12`

- [A waived test, a Go test, and a record with no unit do not hold the review](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B13`

- [Any markdown directly under changes is a delivery record](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B14`

- [The task id is stable and path-safe](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B15`

- [A task in the queue is not queued twice](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-I01`

- [A change pending at the signal is handled before the loop returns, and none after](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-I02`

- [Handling a change writes nothing outside the queue](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-X01`

- [A file the watcher sees for the first time enters the map](layers/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue) `WTCHA-B16`

- [On unix a detached child leads its own process group](layers/comando.md#wtdmw--watchdaemon--the-watcher-started-in-the-background-survives-the-terminal-that-started-it) `WTDMW-B01`

- [On Windows a detached child is created in a new process group](layers/comando.md#wtdmw--watchdaemon--the-watcher-started-in-the-background-survives-the-terminal-that-started-it) `WTDMW-B02`

- [Exactly one detachment implementation builds per platform](layers/comando.md#wtdmw--watchdaemon--the-watcher-started-in-the-background-survives-the-terminal-that-started-it) `WTDMW-I01`

- [Detaching does not start the child](layers/comando.md#wtdmw--watchdaemon--the-watcher-started-in-the-background-survives-the-terminal-that-started-it) `WTDMW-X01`

- [The command refuses what it cannot compose](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B01`

- [A derived piece is redirected to its unit](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B02`

- [A waived piece stops the prompt](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B03`

- [A declarative layer has no unit piece](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B04`

- [The heading and the role follow the stage](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B05`

- [The target section names the layer, or says it is unclassified](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B06`

- [The guides to read are the artifact's and the layer's, each once and sorted](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B07`

- [The pieces are listed where the project derives them](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B08`

- [Feature and test stages list the regime tags](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B09`

- [Each stage says what is not its scope](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B10`

- [The procedure follows the stage and the layer](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B11`

- [The gates' demands are listed once, for the stage](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B12`

- [A review confronts the unit's delivery records](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B13`

- [The findings already recorded for the unit are listed](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B14`

- [Tests and unit reviews explain the execution signals](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B15`

- [Producing stages record the delivery, reviews close with a verdict](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B16`

- [The verification runs over the stage's own piece](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B17`

- [The waivers are listed for the stage](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B18`

- [When the ruler does not decide, the open-decisions section is named](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B19`

- [Rule marking follows the project's policy](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-B20`

- [No prompt both records a delivery and closes a review](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-I01`

- [Composing writes nothing](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-X01`

- [The prompt cites one open-decisions title and value and the current names](layers/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares) `WRPRW-X02`

- [A file the map does not know is refused](layers/comando.md#dtaui--audit--the-dossier-of-everything-pending-on-one-file-for-fixing-it-in-one-pass) `DTAUI-B01`

- [Without impact the dossier covers the file alone](layers/comando.md#dtaui--audit--the-dossier-of-everything-pending-on-one-file-for-fixing-it-in-one-pass) `DTAUI-B02`

- [With impact the dossier covers the unit on the impact path](layers/comando.md#dtaui--audit--the-dossier-of-everything-pending-on-one-file-for-fixing-it-in-one-pass) `DTAUI-B03`

- [A gate result shows by verdict with the first line of its detail](layers/comando.md#dtaui--audit--the-dossier-of-everything-pending-on-one-file-for-fixing-it-in-one-pass) `DTAUI-B04`

- [A doctor finding shows by severity when it cites a node in scope](layers/comando.md#dtaui--audit--the-dossier-of-everything-pending-on-one-file-for-fixing-it-in-one-pass) `DTAUI-B05`

- [Only a failing gate and a warning finding count as actionable](layers/comando.md#dtaui--audit--the-dossier-of-everything-pending-on-one-file-for-fixing-it-in-one-pass) `DTAUI-B06`

- [A scope with nothing pending says so](layers/comando.md#dtaui--audit--the-dossier-of-everything-pending-on-one-file-for-fixing-it-in-one-pass) `DTAUI-B07`

- [The target prints first and the impact nodes after it](layers/comando.md#dtaui--audit--the-dossier-of-everything-pending-on-one-file-for-fixing-it-in-one-pass) `DTAUI-B08`

- [The impact nodes print sorted by path, the same on every run](layers/comando.md#dtaui--audit--the-dossier-of-everything-pending-on-one-file-for-fixing-it-in-one-pass) `DTAUI-B09`

- [Nothing outside the audited scope reaches the dossier](layers/comando.md#dtaui--audit--the-dossier-of-everything-pending-on-one-file-for-fixing-it-in-one-pass) `DTAUI-I01`

- [A project without configuration fails loading it](layers/comando.md#dtaui--audit--the-dossier-of-everything-pending-on-one-file-for-fixing-it-in-one-pass) `DTAUI-E01`

- [A project without a map fails pointing at the map build](layers/comando.md#dtaui--audit--the-dossier-of-everything-pending-on-one-file-for-fixing-it-in-one-pass) `DTAUI-E02`

- [Each duty is reported under the norm that originates it](layers/comando.md#cmpln--compliance--the-state-of-each-regulatory-duty-grouped-by-the-norm-that-imposes-it) `CMPLN-B01`

- [Each duty is marked by how many of its subjects comply](layers/comando.md#cmpln--compliance--the-state-of-each-regulatory-duty-grouped-by-the-norm-that-imposes-it) `CMPLN-B02`

- [A duty that no subject complies with and no debt explains warns of a disconnected target](layers/comando.md#cmpln--compliance--the-state-of-each-regulatory-duty-grouped-by-the-norm-that-imposes-it) `CMPLN-B03`

- [Verbose lists the missing nodes and the plain report only hints at them](layers/comando.md#cmpln--compliance--the-state-of-each-regulatory-duty-grouped-by-the-norm-that-imposes-it) `CMPLN-B04`

- [A project without duties says so](layers/comando.md#cmpln--compliance--the-state-of-each-regulatory-duty-grouped-by-the-norm-that-imposes-it) `CMPLN-B05`

- [The embedded packs the project did not adopt are listed](layers/comando.md#cmpln--compliance--the-state-of-each-regulatory-duty-grouped-by-the-norm-that-imposes-it) `CMPLN-B06`

- [The total counts node-duty pairs, and each hint names where its target is declared](layers/comando.md#cmpln--compliance--the-state-of-each-regulatory-duty-grouped-by-the-norm-that-imposes-it) `CMPLN-B07`

- [Every duty in force has its line even when no node is subject](layers/comando.md#cmpln--compliance--the-state-of-each-regulatory-duty-grouped-by-the-norm-that-imposes-it) `CMPLN-I01`

- [A pack missing a required value fails naming it](layers/comando.md#cmpln--compliance--the-state-of-each-regulatory-duty-grouped-by-the-norm-that-imposes-it) `CMPLN-E01`

- [A project without configuration fails loading it](layers/comando.md#cmpln--compliance--the-state-of-each-regulatory-duty-grouped-by-the-norm-that-imposes-it) `CMPLN-E02`

- [A project without a map fails pointing at the map build](layers/comando.md#cmpln--compliance--the-state-of-each-regulatory-duty-grouped-by-the-norm-that-imposes-it) `CMPLN-E03`

- [Bare guide prints the operating playbook](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B01`

- [guide --help lists every subcommand](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B14`

- [Each guide subcommand prints its own guide](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B02`

- [The review and work guides append the autonomy section of the root they are given](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B03`

- [Register adds exactly the four governance commands](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B04`

- [The review guide teaches who counts and which verdict wins](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B05`

- [The review guide says the reviewer does not move the card](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B06`

- [The review guide carries continuous conformance points anchored in its prose](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B07`

- [The work guide teaches the claim as the first step](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B08`

- [The work guide makes the agent wait for the CI verdict](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B09`

- [The work guide forbids closing the card by hand](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B10`

- [The work guide sends a finding through escalate and reads the decision queue first](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B11`

- [The project guide covers the discover phase and the playbook points to it before planning](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B12`

- [The test guide names the instrument per input shape and teaches the stamp refresh](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B13`

- [Every command a guide cites exists in the command tree](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-I01`

- [The work guide uses the real label names](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-I02`

- [The verdict line the review guide teaches is the one the pipeline parses](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-I03`

- [Only the review and work guides depend on the project](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-X01`

- [The guides tell a fix from a bug and ask for the marker](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B15`

- [The changelog guide tells the technical changelog from the product one](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B16`

- [The spec guide asks for four passes and a review](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B17`

- [The report-bug guide says how to tell, report and go on](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B18`

- [The guides recommend visual regression for every state of a visual unit](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B19`

- [The test guide recommends a contract test against the compiled OpenAPI](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B20`

- [The guides tie validations to states and errors to messages](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B21`

- [The spec guide asks the environment variables to be declared](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B22`

- [The spec guide asks a unit that loads data for its four states and the failure of its load](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B23`

- [The header guide names the flags beside the code, and the navigation guide shows the screen's In and Out and the flag on every call](layers/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on) `GVGDG-B24`

- [The board ranks each guide by how many files it governs](layers/comando.md#gvrns--governs--who-each-guide-governs-and-how-many-read-from-the-map) `GVRNS-B01`

- [A map without governance has an empty board](layers/comando.md#gvrns--governs--who-each-guide-governs-and-how-many-read-from-the-map) `GVRNS-B02`

- [The detail of a guide groups the files it governs by kind](layers/comando.md#gvrns--governs--who-each-guide-governs-and-how-many-read-from-the-map) `GVRNS-B03`

- [A file that governs nobody is answered, not refused](layers/comando.md#gvrns--governs--who-each-guide-governs-and-how-many-read-from-the-map) `GVRNS-B04`

- [The guide argument is resolved against the root and the map can be given](layers/comando.md#gvrns--governs--who-each-guide-governs-and-how-many-read-from-the-map) `GVRNS-B05`

- [Only governs edges count as governance](layers/comando.md#gvrns--governs--who-each-guide-governs-and-how-many-read-from-the-map) `GVRNS-I01`

- [A missing map fails pointing at the map build](layers/comando.md#gvrns--governs--who-each-guide-governs-and-how-many-read-from-the-map) `GVRNS-E01`

- [A role that decides the product is told to record its decisions](layers/comando.md#atgdt--autonomyguide--what-an-agent-does-with-what-it-does-not-know-by-the-role-declared-locally) `ATGDT-B01`

- [A declared role that does not decide is told who decides](layers/comando.md#atgdt--autonomyguide--what-an-agent-does-with-what-it-does-not-know-by-the-role-declared-locally) `ATGDT-B02`

- [With no role declared the guide is the closed one](layers/comando.md#atgdt--autonomyguide--what-an-agent-does-with-what-it-does-not-know-by-the-role-declared-locally) `ATGDT-B03`

- [Every profile reads that preparing the environment asks no authorization](layers/comando.md#atgdt--autonomyguide--what-an-agent-does-with-what-it-does-not-know-by-the-role-declared-locally) `ATGDT-B04`

- [A role with a lens reads it](layers/comando.md#atgdt--autonomyguide--what-an-agent-does-with-what-it-does-not-know-by-the-role-declared-locally) `ATGDT-B05`

- [Whoever does not decide is told not to ask, to move on and what not to escalate](layers/comando.md#atgdt--autonomyguide--what-an-agent-does-with-what-it-does-not-know-by-the-role-declared-locally) `ATGDT-B06`

- [The old user_issues flag with no role is named as such, never as an empty role](layers/comando.md#atgdt--autonomyguide--what-an-agent-does-with-what-it-does-not-know-by-the-role-declared-locally) `ATGDT-B07`

- [An unreadable declaration reads as no role](layers/comando.md#atgdt--autonomyguide--what-an-agent-does-with-what-it-does-not-know-by-the-role-declared-locally) `ATGDT-I01`

- [A role that decides the product is not forbidden to ask](layers/comando.md#atgdt--autonomyguide--what-an-agent-does-with-what-it-does-not-know-by-the-role-declared-locally) `ATGDT-X01`

- [An unreadable settings file is named, and the guide stays closed](layers/comando.md#atgdt--autonomyguide--what-an-agent-does-with-what-it-does-not-know-by-the-role-declared-locally) `ATGDT-B08`

- [An empty example code defaults to LOGI](layers/comando.md#spgds--specguide--the-projects-own-spec-guide-instantiated-with-its-dialect-and-a-complete-example) `SPGDS-B01`

- [The guide shows the three rule forms and a complete example](layers/comando.md#spgds--specguide--the-projects-own-spec-guide-instantiated-with-its-dialect-and-a-complete-example) `SPGDS-B02`

- [The section titles come from the catalogue the generator uses](layers/comando.md#spgds--specguide--the-projects-own-spec-guide-instantiated-with-its-dialect-and-a-complete-example) `SPGDS-B03`

- [The rule letters are the project's when declared and the canonical ones otherwise](layers/comando.md#spgds--specguide--the-projects-own-spec-guide-instantiated-with-its-dialect-and-a-complete-example) `SPGDS-B04`

- [The code length is stated only when the project declares it](layers/comando.md#spgds--specguide--the-projects-own-spec-guide-instantiated-with-its-dialect-and-a-complete-example) `SPGDS-B05`

- [The guide starts from the command that generates the skeleton](layers/comando.md#spgds--specguide--the-projects-own-spec-guide-instantiated-with-its-dialect-and-a-complete-example) `SPGDS-B06`

- [A project that declares its rule types is not offered the canonical letters](layers/comando.md#spgds--specguide--the-projects-own-spec-guide-instantiated-with-its-dialect-and-a-complete-example) `SPGDS-X01`

- [The project guide asks for the four passes before the spec is handed over](layers/comando.md#spgds--specguide--the-projects-own-spec-guide-instantiated-with-its-dialect-and-a-complete-example) `SPGDS-B07`

- [A failing command prints its error once and exits 1](layers/comando.md#clmnc--climain--the-entry-point-that-stamps-the-build-identity-prints-a-failure-once-and-turns-it-into-the-exit-code-the-hooks-read) `CLMNC-B01`

- [A file the project does not govern exits with the not-governed code](layers/comando.md#clmnc--climain--the-entry-point-that-stamps-the-build-identity-prints-a-failure-once-and-turns-it-into-the-exit-code-the-hooks-read) `CLMNC-B02`

- [The build version reaches the reported version and the map's generator](layers/comando.md#clmnc--climain--the-entry-point-that-stamps-the-build-identity-prints-a-failure-once-and-turns-it-into-the-exit-code-the-hooks-read) `CLMNC-B03`

- [A renamed key in an older config points at migrate](layers/comando.md#clmnc--climain--the-entry-point-that-stamps-the-build-identity-prints-a-failure-once-and-turns-it-into-the-exit-code-the-hooks-read) `CLMNC-B04`

- [Only the ingested failures whose rule carries no conclusion are listed, with the spec that declares them](layers/comando.md#flrsa--failures--the-observed-failures-that-the-spec-has-not-explained-yet) `FLRSA-B01`

- [The open failures are listed from the most frequent to the least](layers/comando.md#flrsa--failures--the-observed-failures-that-the-spec-has-not-explained-yet) `FLRSA-B02`

- [With all the concluded failures are listed too, with their conclusion](layers/comando.md#flrsa--failures--the-observed-failures-that-the-spec-has-not-explained-yet) `FLRSA-B03`

- [An occurrence measured against another version of the spec is flagged](layers/comando.md#flrsa--failures--the-observed-failures-that-the-spec-has-not-explained-yet) `FLRSA-B04`

- [With nothing open the review says so instead of printing an empty list](layers/comando.md#flrsa--failures--the-observed-failures-that-the-spec-has-not-explained-yet) `FLRSA-B05`

- [What the log ingestion binds to a spec is what the failures review lists](layers/comando.md#flrsa--failures--the-observed-failures-that-the-spec-has-not-explained-yet) `FLRSA-I01`

- [The failures review leaves the map and the specs unchanged](layers/comando.md#flrsa--failures--the-observed-failures-that-the-spec-has-not-explained-yet) `FLRSA-X01`

- [Without a map the failures review fails](layers/comando.md#flrsa--failures--the-observed-failures-that-the-spec-has-not-explained-yet) `FLRSA-E01`

- [A spec whose file can no longer be read is skipped by the review](layers/comando.md#flrsa--failures--the-observed-failures-that-the-spec-has-not-explained-yet) `FLRSA-E02`

- [A project with no flow declared is told so and the build is not an error](layers/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map) `FLWOX-B01`

- [Building the flow writes it into the map and keeps the map nodes](layers/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map) `FLWOX-B02`

- [A result no flow routes is named as a warning and the build still succeeds](layers/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map) `FLWOX-B03`

- [The next command lists the valid exits of a step with the suggestion of each result](layers/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map) `FLWOX-B04`

- [A step that fits an action no file declares is flagged](layers/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map) `FLWOX-B05`

- [A terminal step says the work ends and a step with no exit and no terminal mark says it is stuck](layers/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map) `FLWOX-B06`

- [The text drawing keeps the file order, hides the results and sorts each step exits](layers/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map) `FLWOX-B07`

- [The mermaid drawing emits valid identifiers, native line breaks, the stadium shape for terminals and highlights the entry and the ends](layers/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map) `FLWOX-B08`

- [An unknown mermaid direction falls back to top-down](layers/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map) `FLWOX-B09`

- [What the flow build writes is what show and next read](layers/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map) `FLWOX-I01`

- [Showing and navigating the flow leave the map unchanged](layers/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map) `FLWOX-X01`

- [Building a flow without a map is refused and no map is created](layers/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map) `FLWOX-E01`

- [Asking the exits of a step that is not in the flow graph fails](layers/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map) `FLWOX-E02`

- [Showing a flow no name matches fails listing the available flows](layers/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map) `FLWOX-E03`

- [Showing or navigating a map with no flow asks for the flow build](layers/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map) `FLWOX-E04`

- [A change to a spec propagates down to the code and the test it specifies](layers/comando.md#mpcti--impact--what-a-change-to-one-file-reaches-in-both-directions-of-the-map) `MPCTI-B01`

- [A change to the code is validated up against its spec and its guide](layers/comando.md#mpcti--impact--what-a-change-to-one-file-reaches-in-both-directions-of-the-map) `MPCTI-B02`

- [An empty direction is said explicitly instead of printed as an empty list](layers/comando.md#mpcti--impact--what-a-change-to-one-file-reaches-in-both-directions-of-the-map) `MPCTI-B03`

- [A root-relative, native-separator or absolute argument resolves to the same node id](layers/comando.md#mpcti--impact--what-a-change-to-one-file-reaches-in-both-directions-of-the-map) `MPCTI-B04`

- [A relative argument that does not exist under the root is resolved from the working directory](layers/comando.md#mpcti--impact--what-a-change-to-one-file-reaches-in-both-directions-of-the-map) `MPCTI-B05`

- [The resolved node id always uses forward slashes](layers/comando.md#mpcti--impact--what-a-change-to-one-file-reaches-in-both-directions-of-the-map) `MPCTI-I01`

- [The impact query leaves the map unchanged](layers/comando.md#mpcti--impact--what-a-change-to-one-file-reaches-in-both-directions-of-the-map) `MPCTI-X01`

- [A file that is not a node of the map is refused](layers/comando.md#mpcti--impact--what-a-change-to-one-file-reaches-in-both-directions-of-the-map) `MPCTI-E01`

- [Without a map the impact query fails and asks for the map build](layers/comando.md#mpcti--impact--what-a-change-to-one-file-reaches-in-both-directions-of-the-map) `MPCTI-E02`

- [The three reports in one pass reach the test, spec and code nodes](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B01`

- [A JUnit report that matches no test node warns](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B02`

- [A project's declared code length governs how JUnit case names are read](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B03`

- [Manual ingestion warns and proceeds by default](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B04`

- [Ingestion run by anchors test never complains](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B05`

- [A project with no declared suite or no config is not asked to use anchors test](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B06`

- [The suite key is the report path from the root, or the file name for a report outside the repository](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B07`

- [A full run of a suite inside the repository drops the suites ingested from outside it](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B08`

- [A partial run or another external report drops no external suite](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B09`

- [When two report paths name the same node the path equal to the node id wins](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B10`

- [An lcov entry for a file edited after the report was written is marked as predating it](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B11`

- [Log occurrences are bound to the spec that declares the failure, stamped with its revision](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B12`

- [Failure codes no spec declares are reported after the log ingestion](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B13`

- [Ingesting the same logs again replaces the earlier occurrences instead of adding to them](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-I01`

- [A failure code no spec declares is bound to no spec](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-X01`

- [Ingest with no report flag refuses with the usage](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-E01`

- [Manual ingestion is refused when the project declares manual ingestion blocks](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-E02`

- [Log ingestion without declared log paths is refused](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-E03`

- [A missing or malformed report fails naming the report's format](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-E04`

- [Ingesting a report without a map fails and asks for the map build](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-E05`

- [A test file's run time is the sum of its cases](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B14`

- [A spec created after the map build keeps its first proof](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B17`

- [A run's proofs are stamped with the tree's revs](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B16`

- [A report's signals are kept under its path from the root](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B15`

- [A suite's whole coverage run marks the files it left out](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B18`

- [A JUnit case with no file finds its test by its class's folder and the one file defining its test](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B19`

- [A targeted file that gave no mutant is measured with nothing to mutate](layers/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map) `NGSTI-B20`

- [The legacy spelling of the waiver is a waiver all the way to the stamp](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-B01`

- [A reason is required for fail and waived, not for pass](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-B02`

- [Only the exact name of a declared judgment gate is accepted](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-B03`

- [A fail and then a pass stamp the guide edge issue and then ok](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-B04`

- [The review verdict stamps every edge of the target](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-B05`

- [A target not in the map is recorded on its unit's spec](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-B06`

- [A test target not in the map is recorded on its unit's spec, for a Go unit as for a TypeScript one](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-B12`

- [A fail opens an issue, a repeated report changes nothing, a new report reopens it](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-B07`

- [A pass resolves the open issue and a waiver resolves it as a waiver](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-B08`

- [Manual mode stamps the map and prints the report without writing an issue](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-B09`

- [A fail with a patch opens an applicable fix suggestion](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-B10`

- [The verdict closes its judge task and the pending list shrinks](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-B11`

- [A waiver is never stamped ok](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-I01`

- [A waiver is never announced as a pass](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-X01`

- [A judge with no target is refused](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-E01`

- [A judge with no gate is refused](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-E02`

- [A verdict outside the three is refused](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-E03`

- [A gate that is not a judgment gate is refused](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-E05`

- [A project with no map is refused pointing at the build](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-E06`

- [A target whose unit has no piece in the map is refused](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-E07`

- [A map with no configuration beside it is refused](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-E08`

- [An unreadable patch file is refused](layers/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate) `JDGUE-E09`

- [A rebuild keeps the judgment recorded on the guide edge](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-B01`

- [A rebuild keeps the flow graph](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-B02`

- [A lost judgment stamp is reported per gate](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-B03`

- [The edge summary shows every type, the unit's first](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-B04`

- [The layer ambiguity warning is grouped by pair of layers](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-B05`

- [A declared priority silences the layer ambiguity warning](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-B11`

- [Showing the login code lists what governs it and marks it a leaf](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-B06`

- [The orphans are the nodes with no edge](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-B07`

- [The statistics count nodes by kind and edges by type](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-B08`

- [The worklist puts rulers and specs before the code they govern](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-B09`

- [The pending worklist lists only nodes with a failing gate](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-B10`

- [A judgment survives a rebuild of the map](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-I01`

- [The warnings do not fail the build](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-X01`

- [Building with no configuration points at init](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-E01`

- [Showing with no map points at the build](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-E02`

- [Showing a file that is not in the map is refused](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-E03`

- [Showing with no selector is refused](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-E04`

- [The pending worklist with no configuration is refused](layers/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it) `MPCMM-E05`

- [The tree goes down what a file uses, and up who uses it, a cycle marked once](layers/comando.md#mdcmp--mapdeps--the-dependency-tree-of-a-file) `MDCMP-B01`

- [A file is named by its own code, its unit's code or its path](layers/comando.md#mdcmp--mapdeps--the-dependency-tree-of-a-file) `MDCMP-B02`

- [The merged map is written onto our side's file](layers/comando.md#mpmrm--mapmerge--the-git-merge-driver-that-unites-two-versions-of-the-map-instead-of-merging-text) `MPMRM-B01`

- [A node created only on the other branch reaches the merged map](layers/comando.md#mpmrm--mapmerge--the-git-merge-driver-that-unites-two-versions-of-the-map-instead-of-merging-text) `MPMRM-B02`

- [A node both sides have keeps our side's revision](layers/comando.md#mpmrm--mapmerge--the-git-merge-driver-that-unites-two-versions-of-the-map-instead-of-merging-text) `MPMRM-B03`

- [An edge created only on the other branch arrives with its judgment](layers/comando.md#mpmrm--mapmerge--the-git-merge-driver-that-unites-two-versions-of-the-map-instead-of-merging-text) `MPMRM-B04`

- [A side with no judgment on a shared edge does not erase the other side's](layers/comando.md#mpmrm--mapmerge--the-git-merge-driver-that-unites-two-versions-of-the-map-instead-of-merging-text) `MPMRM-B05`

- [The driver reports on the error stream what it kept and what came from the other side](layers/comando.md#mpmrm--mapmerge--the-git-merge-driver-that-unites-two-versions-of-the-map-instead-of-merging-text) `MPMRM-B06`

- [Judgments of different gates on a shared edge are joined, and one gate judged on both sides keeps the latest](layers/comando.md#mpmrm--mapmerge--the-git-merge-driver-that-unites-two-versions-of-the-map-instead-of-merging-text) `MPMRM-B07`

- [The other side's flow states and transitions and the failures of a shared node reach the merged map](layers/comando.md#mpmrm--mapmerge--the-git-merge-driver-that-unites-two-versions-of-the-map-instead-of-merging-text) `MPMRM-B08`

- [Nothing of either side is missing from the merged map](layers/comando.md#mpmrm--mapmerge--the-git-merge-driver-that-unites-two-versions-of-the-map-instead-of-merging-text) `MPMRM-I01`

- [The base version is never read](layers/comando.md#mpmrm--mapmerge--the-git-merge-driver-that-unites-two-versions-of-the-map-instead-of-merging-text) `MPMRM-X01`

- [An unreadable side fails the merge naming the side](layers/comando.md#mpmrm--mapmerge--the-git-merge-driver-that-unites-two-versions-of-the-map-instead-of-merging-text) `MPMRM-E01`

- [A call with two paths is refused](layers/comando.md#mpmrm--mapmerge--the-git-merge-driver-that-unites-two-versions-of-the-map-instead-of-merging-text) `MPMRM-E02`

- [The navigation is printed screen by screen, each with where it leads, sorted](layers/comando.md#mncmp--mapnav--the-apps-navigation-screen-by-screen) `MNCMP-B01`

- [One screen shows where it comes from and where it leads; a name that is no screen is refused](layers/comando.md#mncmp--mapnav--the-apps-navigation-screen-by-screen) `MNCMP-B02`

- [A spec edited after the map build is named as stale](layers/comando.md#mpstm--mapstaleness--names-the-files-of-the-map-whose-content-changed-after-the-map-was-built) `MPSTM-B01`

- [A file removed after the map build is not named](layers/comando.md#mpstm--mapstaleness--names-the-files-of-the-map-whose-content-changed-after-the-map-was-built) `MPSTM-B02`

- [Staleness follows the content, not the modification time](layers/comando.md#mpstm--mapstaleness--names-the-files-of-the-map-whose-content-changed-after-the-map-was-built) `MPSTM-B03`

- [A node with an empty recorded revision is not named](layers/comando.md#mpstm--mapstaleness--names-the-files-of-the-map-whose-content-changed-after-the-map-was-built) `MPSTM-B04`

- [No map names nothing](layers/comando.md#mpstm--mapstaleness--names-the-files-of-the-map-whose-content-changed-after-the-map-was-built) `MPSTM-B05`

- [A freshly rebuilt map is never stale](layers/comando.md#mpstm--mapstaleness--names-the-files-of-the-map-whose-content-changed-after-the-map-was-built) `MPSTM-I01`

- [Asking about staleness leaves the map as it was](layers/comando.md#mpstm--mapstaleness--names-the-files-of-the-map-whose-content-changed-after-the-map-was-built) `MPSTM-X01`

- [A root that cannot be walked names nothing and raises nothing](layers/comando.md#mpstm--mapstaleness--names-the-files-of-the-map-whose-content-changed-after-the-map-was-built) `MPSTM-E01`

- [Codes typed in lower case are renamed in upper case](layers/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project) `RCDEO-B01`

- [The plan counts each file's occurrences by kind](layers/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project) `RCDEO-B02`

- [A recode without the apply switch writes nothing](layers/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project) `RCDEO-B03`

- [Applying rewrites the header and the scenario codes of the spec and the test](layers/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project) `RCDEO-B04`

- [Applying rebuilds the map from the headers](layers/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project) `RCDEO-B05`

- [Applying keeps the judgments, stamps and flow of the previous map](layers/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project) `RCDEO-B06`

- [The plan reports the testIDs and the file renames of the project's dialect](layers/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project) `RCDEO-B07`

- [The plan warns when the files carry testIDs with a prefix other than the dialect's](layers/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project) `RCDEO-B08`

- [A write failure before any file changed fails without saying the project is half converted](layers/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project) `RCDEO-B09`

- [After applying, neither the spec nor the map carries the old code](layers/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project) `RCDEO-I01`

- [The map is rebuilt from the files, not edited as text](layers/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project) `RCDEO-X01`

- [A project with no configuration is refused naming the configuration file](layers/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project) `RCDEO-E01`

- [Recoding a code no file carries is refused](layers/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project) `RCDEO-E02`

- [A write failure after some files changed says the project is half converted](layers/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project) `RCDEO-E03`

- [A recode with a single code is refused](layers/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project) `RCDEO-E04`

- [Registering the map domain makes each of its commands reachable from the root](layers/comando.md#mprgm--mapregister--hangs-the-map-domains-commands-on-the-root-command) `MPRGM-B01`

- [Registering the map domain adds no command outside it](layers/comando.md#mprgm--mapregister--hangs-the-map-domains-commands-on-the-root-command) `MPRGM-X01`

- [The default base is the remote integration branch, else the local one](layers/comando.md#rnmbr--renumber--moves-the-revisions-a-branch-added-when-the-base-already-took-their-number) `RNMBR-B01`

- [Specs changed but not committed, and untracked specs, are examined](layers/comando.md#rnmbr--renumber--moves-the-revisions-a-branch-added-when-the-base-already-took-their-number) `RNMBR-B02`

- [Only the branch's colliding revision moves to the next free number](layers/comando.md#rnmbr--renumber--moves-the-revisions-a-branch-added-when-the-base-already-took-their-number) `RNMBR-B03`

- [Citations move only on the lines the branch added](layers/comando.md#rnmbr--renumber--moves-the-revisions-a-branch-added-when-the-base-already-took-their-number) `RNMBR-B04`

- [A changed binary file is not rewritten](layers/comando.md#rnmbr--renumber--moves-the-revisions-a-branch-added-when-the-base-already-took-their-number) `RNMBR-B05`

- [The dry run writes nothing](layers/comando.md#rnmbr--renumber--moves-the-revisions-a-branch-added-when-the-base-already-took-their-number) `RNMBR-B06`

- [No collision says there is nothing to renumber](layers/comando.md#rnmbr--renumber--moves-the-revisions-a-branch-added-when-the-base-already-took-their-number) `RNMBR-B07`

- [A rewritten file keeps its permission bits](layers/comando.md#rnmbr--renumber--moves-the-revisions-a-branch-added-when-the-base-already-took-their-number) `RNMBR-B08`

- [With files given, only those specs are examined](layers/comando.md#rnmbr--renumber--moves-the-revisions-a-branch-added-when-the-base-already-took-their-number) `RNMBR-B09`

- [The base's revisions keep their meaning after a renumber](layers/comando.md#rnmbr--renumber--moves-the-revisions-a-branch-added-when-the-base-already-took-their-number) `RNMBR-I01`

- [The rewrite is left unstaged for review](layers/comando.md#rnmbr--renumber--moves-the-revisions-a-branch-added-when-the-base-already-took-their-number) `RNMBR-X01`

- [A project with no configuration is refused](layers/comando.md#rnmbr--renumber--moves-the-revisions-a-branch-added-when-the-base-already-took-their-number) `RNMBR-E01`

- [A base with no merge base is refused naming it](layers/comando.md#rnmbr--renumber--moves-the-revisions-a-branch-added-when-the-base-already-took-their-number) `RNMBR-E02`

- [An unreadable spec given by hand is refused naming it](layers/comando.md#rnmbr--renumber--moves-the-revisions-a-branch-added-when-the-base-already-took-their-number) `RNMBR-E03`

- [Pending lists what is to review with its question](layers/comando.md#rvcmr--reviewcommand--record-a-review-with-who-looked-and-what-they-found-or-list-what-is-to-review) `RVCMR-B01`

- [A review with no findings records who looked and opens no issue](layers/comando.md#rvcmr--reviewcommand--record-a-review-with-who-looked-and-what-they-found-or-list-what-is-to-review) `RVCMR-B02`

- [A review with findings records them and opens the issue](layers/comando.md#rvcmr--reviewcommand--record-a-review-with-who-looked-and-what-they-found-or-list-what-is-to-review) `RVCMR-B03`

- [In manual mode the findings write no issue unless asked](layers/comando.md#rvcmr--reviewcommand--record-a-review-with-who-looked-and-what-they-found-or-list-what-is-to-review) `RVCMR-B04`

- [A review record always names who reviewed](layers/comando.md#rvcmr--reviewcommand--record-a-review-with-who-looked-and-what-they-found-or-list-what-is-to-review) `RVCMR-I01`

- [There is no waived](layers/comando.md#rvcmr--reviewcommand--record-a-review-with-who-looked-and-what-they-found-or-list-what-is-to-review) `RVCMR-X01`

- [A wrong target, gate or no reviewer refuses the record](layers/comando.md#rvcmr--reviewcommand--record-a-review-with-who-looked-and-what-they-found-or-list-what-is-to-review) `RVCMR-E01`

- [The live payload says it is live and stamps the read time](layers/comando.md#brsrb--boardserve--the-board-page-served-locally-with-live-state-read-from-the-host-only-when-something-changed) `BRSRB-B01`

- [A read inside the floor does not call the host](layers/comando.md#brsrb--boardserve--the-board-page-served-locally-with-live-state-read-from-the-host-only-when-something-changed) `BRSRB-B02`

- [Past the floor the board sweeps only when something newer exists](layers/comando.md#brsrb--boardserve--the-board-page-served-locally-with-live-state-read-from-the-host-only-when-something-changed) `BRSRB-B03`

- [A failed sweep after a good read serves the good read](layers/comando.md#brsrb--boardserve--the-board-page-served-locally-with-live-state-read-from-the-host-only-when-something-changed) `BRSRB-B04`

- [The full sweep stitches pages, drops pull requests and keeps the owner](layers/comando.md#brsrb--boardserve--the-board-page-served-locally-with-live-state-read-from-the-host-only-when-something-changed) `BRSRB-B05`

- [The comments of open cards follow the cursor and failures yield none](layers/comando.md#brsrb--boardserve--the-board-page-served-locally-with-live-state-read-from-the-host-only-when-something-changed) `BRSRB-B06`

- [The repository comes from the clone and the banner says what is served](layers/comando.md#brsrb--boardserve--the-board-page-served-locally-with-live-state-read-from-the-host-only-when-something-changed) `BRSRB-B07`

- [The server hands out the published page and the live JSON](layers/comando.md#brsrb--boardserve--the-board-page-served-locally-with-live-state-read-from-the-host-only-when-something-changed) `BRSRB-B08`

- [The incremental baseline is the newest update of the served board](layers/comando.md#brsrb--boardserve--the-board-page-served-locally-with-live-state-read-from-the-host-only-when-something-changed) `BRSRB-I01`

- [The page is the pipeline's own board page](layers/comando.md#brsrb--boardserve--the-board-page-served-locally-with-live-state-read-from-the-host-only-when-something-changed) `BRSRB-X01`

- [A refused sweep carries the host's message](layers/comando.md#brsrb--boardserve--the-board-page-served-locally-with-live-state-read-from-the-host-only-when-something-changed) `BRSRB-E01`

- [With nothing to serve the data route answers 502](layers/comando.md#brsrb--boardserve--the-board-page-served-locally-with-live-state-read-from-the-host-only-when-something-changed) `BRSRB-E02`

- [Without a repository the command points at --repo](layers/comando.md#brsrb--boardserve--the-board-page-served-locally-with-live-state-read-from-the-host-only-when-something-changed) `BRSRB-E03`

- [The latest release by default](layers/comando.md#clgcm--anchors-changelog--the-technical-changelog-printed-or-written) `CLGCM-B01`

- [From a tag, or every release](layers/comando.md#clgcm--anchors-changelog--the-technical-changelog-printed-or-written) `CLGCM-B02`

- [What is not released yet](layers/comando.md#clgcm--anchors-changelog--the-technical-changelog-printed-or-written) `CLGCM-B03`

- [Writing the incremental file](layers/comando.md#clgcm--anchors-changelog--the-technical-changelog-printed-or-written) `CLGCM-B04`

- [One file per release](layers/comando.md#clgcm--anchors-changelog--the-technical-changelog-printed-or-written) `CLGCM-B05`

- [The project's template and language](layers/comando.md#clgcm--anchors-changelog--the-technical-changelog-printed-or-written) `CLGCM-B06`

- [No anchors.yaml](layers/comando.md#clgcm--anchors-changelog--the-technical-changelog-printed-or-written) `CLGCM-B07`

- [Nothing to list](layers/comando.md#clgcm--anchors-changelog--the-technical-changelog-printed-or-written) `CLGCM-B08`

- [An unknown start tag](layers/comando.md#clgcm--anchors-changelog--the-technical-changelog-printed-or-written) `CLGCM-E01`

- [A template that cannot be read](layers/comando.md#clgcm--anchors-changelog--the-technical-changelog-printed-or-written) `CLGCM-E02`

- [A broken anchors.yaml](layers/comando.md#clgcm--anchors-changelog--the-technical-changelog-printed-or-written) `CLGCM-E03`

- [A free canonical code is the suggestion](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-B01`

- [A taken canonical is adjusted to a free code naming its owner](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-B02`

- [A path in a layer with a code prefix gets the module prefix](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-B03`

- [A generic basename takes its identity from the parent directory](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-B04`

- [Check answers whether a code is free, ignoring case, and names the owners](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-B05`

- [The list prints one sorted code per line with its folder, the summary kept off stdout](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-B06`

- [The list filters by path prefix and names a filter that matched nothing](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-B07`

- [The JSON list carries each code's folder, file, kind, title and work order fields](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-B08`

- [The title drops the text before the dash](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-B09`

- [The length check accuses only declared codes and proposes the canonical code](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-B10`

- [The unit name drops the artifact suffixes](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-B11`

- [An empty map says no node has an identity](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-B12`

- [The codes come from the map's identity field](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-X01`

- [Without a name or a map the command fails and says what to do](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-E01`

- [The list refuses a project without config or without map](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-E02`

- [A second check in the same process judges only its own map](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-B13`

- [The length check points each divergence to anchors recode, and there is no --fix](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-B14`

- [A name shaped like a code also gets that code's status](layers/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map) `CDCMC-B15`

- [The subject is the first line that is neither blank nor a comment](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B01`

- [Messages git generates pass](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B02`

- [A conventional subject passes, with an optional scope and break mark](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B03`

- [A type with an uppercase letter or outside the closed list is refused](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B04`

- [An empty scope is refused](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B05`

- [A subject over the limit is refused and the diagnosis points to the body](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B06`

- [A subject ending in a period is refused](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B07`

- [Each defect gets its own diagnosis](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B08`

- [A capital letter at the start of the subject is allowed](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B09`

- [The rejection names the subject, teaches the format and lists the types](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B10`

- [An empty or comment-only message passes](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B11`

- [Checks run in order and only the first defect is reported](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B12`

- [A type with nothing after the colon is refused](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B13`

- [The subject limit counts characters, not bytes](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B14`

- [A subject with no space after the colon is refused](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B15`

- [The command only accepts or refuses](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-X01`

- [A message file that cannot be read fails naming the read](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-E01`

- [The Bug footer marks a fix of a defect that shipped](layers/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists) `CMMSC-B16`

- [Build compiles from the tree, even what the map on disk does not know](layers/comando.md#dccmd--docscommand--the-documentation-is-compiled-from-the-specs-through-templates-against-a-map-rebuilt-from-the-tree) `DCCMD-B01`

- [No-map-rebuild compiles against the map on disk and warns](layers/comando.md#dccmd--docscommand--the-documentation-is-compiled-from-the-specs-through-templates-against-a-map-rebuilt-from-the-tree) `DCCMD-B02`

- [A dry run compiles without writing](layers/comando.md#dccmd--docscommand--the-documentation-is-compiled-from-the-specs-through-templates-against-a-map-rebuilt-from-the-tree) `DCCMD-B03`

- [A hand-written page is skipped, kept and named](layers/comando.md#dccmd--docscommand--the-documentation-is-compiled-from-the-specs-through-templates-against-a-map-rebuilt-from-the-tree) `DCCMD-B04`

- [Build with no template says there is nothing to compile](layers/comando.md#dccmd--docscommand--the-documentation-is-compiled-from-the-specs-through-templates-against-a-map-rebuilt-from-the-tree) `DCCMD-B05`

- [Init writes the skeleton once and needs a map](layers/comando.md#dccmd--docscommand--the-documentation-is-compiled-from-the-specs-through-templates-against-a-map-rebuilt-from-the-tree) `DCCMD-B06`

- [Duties answers for the project, a layer or a unit](layers/comando.md#dccmd--docscommand--the-documentation-is-compiled-from-the-specs-through-templates-against-a-map-rebuilt-from-the-tree) `DCCMD-B07`

- [Duties without a declaration teaches the known kinds](layers/comando.md#dccmd--docscommand--the-documentation-is-compiled-from-the-specs-through-templates-against-a-map-rebuilt-from-the-tree) `DCCMD-B08`

- [Build never writes the map](layers/comando.md#dccmd--docscommand--the-documentation-is-compiled-from-the-specs-through-templates-against-a-map-rebuilt-from-the-tree) `DCCMD-X01`

- [Build without a config points at init](layers/comando.md#dccmd--docscommand--the-documentation-is-compiled-from-the-specs-through-templates-against-a-map-rebuilt-from-the-tree) `DCCMD-E01`

- [Duties without a config fails](layers/comando.md#dccmd--docscommand--the-documentation-is-compiled-from-the-specs-through-templates-against-a-map-rebuilt-from-the-tree) `DCCMD-E02`

- [No-map-rebuild without a map says so](layers/comando.md#dccmd--docscommand--the-documentation-is-compiled-from-the-specs-through-templates-against-a-map-rebuilt-from-the-tree) `DCCMD-E03`

- [A blank reason is refused](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-B01`

- [The freeze writes two quoted lines on top and the file still loads](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-B02`

- [The freeze is committed and pushed past refusing hooks](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-B03`

- [No-push leaves the frozen file as a local change](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-B04`

- [In github mode the freeze creates the rule and opens the issue](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-B05`

- [No-ruleset skips the ruleset and still opens the issue](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-B06`

- [Freezing twice changes nothing and shows the reason](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-B07`

- [Each failing layer is a warning and the local brake holds](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-B08`

- [The thaw undoes the three layers](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-B09`

- [Thawing a project that is not frozen is a no-op](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-B10`

- [With nothing on the remote the thaw removes nothing](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-B11`

- [The deprecated Portuguese flag names still work](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-B12`

- [A config that already declares enabled is frozen with a single key and still loads](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-B13`

- [Freeze and thaw give back the file byte for byte](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-I01`

- [The configuration is never reserialized](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-X01`

- [Freeze and thaw refuse a project without config](layers/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers) `FRZEX-E01`

- [The patterns match the map, the compiled docs and plan progress files only](layers/comando.md#gnptg--generatedpaths--the-product-names-the-files-it-derives-so-a-conflict-in-them-is-rebuilt-not-merged) `GNPTG-B01`

- [The default output is one pattern per line](layers/comando.md#gnptg--generatedpaths--the-product-names-the-files-it-derives-so-a-conflict-in-them-is-rebuilt-not-merged) `GNPTG-B02`

- [The re format is one valid alternation](layers/comando.md#gnptg--generatedpaths--the-product-names-the-files-it-derives-so-a-conflict-in-them-is-rebuilt-not-merged) `GNPTG-B03`

- [The dot of a path matches only a literal dot](layers/comando.md#gnptg--generatedpaths--the-product-names-the-files-it-derives-so-a-conflict-in-them-is-rebuilt-not-merged) `GNPTG-B04`

- [Both forms carry the same patterns](layers/comando.md#gnptg--generatedpaths--the-product-names-the-files-it-derives-so-a-conflict-in-them-is-rebuilt-not-merged) `GNPTG-I01`

- [The command only names the derived files](layers/comando.md#gnptg--generatedpaths--the-product-names-the-files-it-derives-so-a-conflict-in-them-is-rebuilt-not-merged) `GNPTG-X01`

- [A directory without anchors.yaml is refused](layers/comando.md#gnptg--generatedpaths--the-product-names-the-files-it-derives-so-a-conflict-in-them-is-rebuilt-not-merged) `GNPTG-E01`

- [Without answers the command only asks](layers/comando.md#inint--initnoninteractive--the-init-that-an-agent-answers-with-flags-it-asks-in-json-and-writes-only-a-complete-valid-set-of-answers) `ININT-B01`

- [The given answers reach the configuration](layers/comando.md#inint--initnoninteractive--the-init-that-an-agent-answers-with-flags-it-asks-in-json-and-writes-only-a-complete-valid-set-of-answers) `ININT-B02`

- [A false flag is a deliberate no](layers/comando.md#inint--initnoninteractive--the-init-that-an-agent-answers-with-flags-it-asks-in-json-and-writes-only-a-complete-valid-set-of-answers) `ININT-B03`

- [The github workflow carries the repository and labels](layers/comando.md#inint--initnoninteractive--the-init-that-an-agent-answers-with-flags-it-asks-in-json-and-writes-only-a-complete-valid-set-of-answers) `ININT-B04`

- [One invalid answer refuses the whole set](layers/comando.md#inint--initnoninteractive--the-init-that-an-agent-answers-with-flags-it-asks-in-json-and-writes-only-a-complete-valid-set-of-answers) `ININT-B05`

- [Defaults writes when asked to](layers/comando.md#inint--initnoninteractive--the-init-that-an-agent-answers-with-flags-it-asks-in-json-and-writes-only-a-complete-valid-set-of-answers) `ININT-B06`

- [Layers prunes the other code layers](layers/comando.md#inint--initnoninteractive--the-init-that-an-agent-answers-with-flags-it-asks-in-json-and-writes-only-a-complete-valid-set-of-answers) `ININT-B08`

- [The success document names the file and the next step](layers/comando.md#inint--initnoninteractive--the-init-that-an-agent-answers-with-flags-it-asks-in-json-and-writes-only-a-complete-valid-set-of-answers) `ININT-B09`

- [Either the whole set is written or nothing is](layers/comando.md#inint--initnoninteractive--the-init-that-an-agent-answers-with-flags-it-asks-in-json-and-writes-only-a-complete-valid-set-of-answers) `ININT-I01`

- [The non-interactive mode never prompts](layers/comando.md#inint--initnoninteractive--the-init-that-an-agent-answers-with-flags-it-asks-in-json-and-writes-only-a-complete-valid-set-of-answers) `ININT-X01`

- [A malformed governs rule is refused with the expected form](layers/comando.md#inint--initnoninteractive--the-init-that-an-agent-answers-with-flags-it-asks-in-json-and-writes-only-a-complete-valid-set-of-answers) `ININT-E01`

- [The governs rules reach the configuration, one rule per tag](layers/comando.md#inint--initnoninteractive--the-init-that-an-agent-answers-with-flags-it-asks-in-json-and-writes-only-a-complete-valid-set-of-answers) `ININT-B10`

- [CONTRIBUTING.md is seeded when absent, and an existing one is left as it is](layers/comando.md#inint--initnoninteractive--the-init-that-an-agent-answers-with-flags-it-asks-in-json-and-writes-only-a-complete-valid-set-of-answers) `ININT-B11`

- [Seeding the gates tells to write gate steps in the project's language](layers/comando.md#inint--initnoninteractive--the-init-that-an-agent-answers-with-flags-it-asks-in-json-and-writes-only-a-complete-valid-set-of-answers) `ININT-B12`

- [An existing config is kept when the overwrite is not confirmed](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B01`

- [A ready repository makes the git step silent](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B02`

- [Without git the step warns and the init goes on](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B03`

- [Declining git names what stays off](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B04`

- [Accepting git leaves a repository with HEAD](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B05`

- [An existing gitignore is kept](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B06`

- [A failed git initialization is reported and names the fix](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B07`

- [The findings report only what exists on disk](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B08`

- [The DISCOVER step is silent when the phase already happened](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B09`

- [An AI operator gets the DISCOVER work order](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B10`

- [A person without a known AI gets the prompt to paste](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B11`

- [A detected AI is opened in the root with the prompt, or declined for the step-by-step](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B12`

- [Non-interactive routes to the JSON mode](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B14`

- [The header guide is seeded in guides/ when the project has no guide directory](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B16`

- [The default gates are offered only when the chosen artifacts have any, and accepted ones are written](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B17`

- [A project with no code, spec, feature or test is announced as new](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B18`

- [The code-layer question is asked only when there are code layers, and a new project is told to declare them later](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B19`

- [Each guide found is asked which tag it governs](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B20`

- [A prompt that cannot run makes the init write nothing](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-I01`

- [Git is never initialized nor committed without a yes](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-X01`

- [The no-terminal error offers the non-interactive mode](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-E01`

- [End of input in line mode is refused, not taken as the defaults](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-E02`

- [--preset is refused with the reason, and nothing is written](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B21`

- [The code-layer question is preceded by the note on layers](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B22`

- [CONTRIBUTING.md is seeded when absent, and an existing one is shown, not touched](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B23`

- [The family's coverage hint is printed](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B24`

- [Accepting the gates tells to write gate steps in the project's language](layers/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave) `INWZN-B25`

- [The hooks go where git looks for them](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-B01`

- [A fresh install writes the three managed hooks, executable](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-B02`

- [Foreign hooks are respected unless forced](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-B03`

- [The hooks anchors wrote, old or new, are updated on reinstall](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-B04`

- [The pre-commit refuses a commit while the remote is frozen](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-B05`

- [A staged set with nothing governed passes the pre-commit](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-B06`

- [A gate failure is deferred to the commit-msg, which blocks it](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-B07`

- [The pre-push refuses a push while the remote is frozen](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-B08`

- [The commit-msg refuses a subject the message check refuses](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-B09`

- [The pre-push refuses a binary older than the remote's minimum version](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-B10`

- [Both merge drivers are registered in git config and .gitattributes](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-B11`

- [Reinstalling never duplicates an attribute line nor damages the user's lines](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-I01`

- [A hook the user wrote is never replaced without --force](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-X01`

- [A directory without anchors.yaml gets no hook](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-E01`

- [Outside a repository the error explains what needs git](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-E02`

- [The pre-push warns when the remote map was written by another version](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-B12`

- [The installed hooks run on a real commit, on every system](layers/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote) `INHKN-B13`

- [Both the map and the config reach the current format](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B01`

- [A missing file is reported and the other still migrates](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B02`

- [A file already current is reported as such](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B03`

- [The renamed keys are listed in alphabetical order with their counts](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B04`

- [A dry run reports without writing](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B05`

- [A real migration asks for the commit](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B06`

- [A second run changes nothing](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-I01`

- [The keys renamed are the migration package's](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-X01`

- [The command reminds of the commit and does not make it](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-X02`

- [Crossing format 5 rewrites the letters of plans, flows and actions](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B07`

- [Crossing format 7 widens the four-character codes and renames the files named by them](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B08`

- [Crossing format 7 gives every governed file a code of its own](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B09`

- [Crossing format 7 carries what each file was measured at to its new revision](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B10`

- [Crossing format 7 records each renamed code in anchors.renames.yaml](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B11`

- [Crossing format 7 turns a file carrying its spec's code into a ref with a code of its own](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B12`

- [Crossing format 7, a file that cannot be written fails the command, and a second run finishes it](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-E02`

- [Crossing format 7 refreshes the stamps a widened code broke](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B13`

- [Crossing format 7 dates every file it rewrote](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B14`

- [Crossing format 7 rewrites the rule codes cited in files the project does not govern](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B15`

- [Crossing format 7 rewrites a code only where it is cited as one, and lists the bare words it left](layers/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format) `MGCMM-B16`

- [An unknown kind is refused](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B01`

- [The name and the output path are required](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B02`

- [A spec is born with a code no unit uses](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B03`

- [A feature or test takes its identity from the sibling spec, or warns it is orphaned](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B04`

- [Code pins the identity](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B05`

- [With and without are validated against the kind's sections](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B06`

- [A preset fixes the sections and their order, and extra sections follow](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B07`

- [A spec for a declarative layer is refused](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B08`

- [Section titles and bodies follow the project lexicon, then its language](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B09`

- [Features and tests follow the project's declared dialect](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B10`

- [The unit regime tag comes from the project, or a visible placeholder](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B11`

- [The artifact is written where out says and never overwritten](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B12`

- [A plan is born with its progress companion](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B13`

- [List-sections prints the menu, with presets only for specs](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B14`

- [The target layer is the one of the unit the artifact describes](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B15`

- [A refused new leaves nothing behind](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-I01`

- [A feature references the spec's identity instead of owning one](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-X01`

- [The help and the unknown-kind refusal name every kind, and --out is mandatory](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B16`

- [The new artifact enters the map at once](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B17`

- [A new artifact that refs its unit is born with a code of its own](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B18`

- [A rule section takes the title its letter lists alone, not one every letter shares](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B19`

- [The screen preset writes the four states of a unit that loads data and the failure of its load, under the titles of States and Errors](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B20`

- [Two sections of one letter never define the same code, and a section of the rules' letter with no title of its own goes inside the rules section the layer names](layers/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler) `NWARN-B21`

- [The catalog holds seven kinds and each is born by new](layers/comando.md#nwtmn--newtemplates--the-catalog-of-artifact-skeletons-which-kinds-exist-their-headers-their-sections-and-the-presets-that-pick-them) `NWTMN-B01`

- [A markdown header is an HTML comment with the identity and placeholders](layers/comando.md#nwtmn--newtemplates--the-catalog-of-artifact-skeletons-which-kinds-exist-their-headers-their-sections-and-the-presets-that-pick-them) `NWTMN-B02`

- [A feature header opens with the Gherkin language line](layers/comando.md#nwtmn--newtemplates--the-catalog-of-artifact-skeletons-which-kinds-exist-their-headers-their-sections-and-the-presets-that-pick-them) `NWTMN-B03`

- [A test header uses the line comment of the output file](layers/comando.md#nwtmn--newtemplates--the-catalog-of-artifact-skeletons-which-kinds-exist-their-headers-their-sections-and-the-presets-that-pick-them) `NWTMN-B04`

- [Every section realizes only canonical rule letters](layers/comando.md#nwtmn--newtemplates--the-catalog-of-artifact-skeletons-which-kinds-exist-their-headers-their-sections-and-the-presets-that-pick-them) `NWTMN-B05`

- [The test body is idiomatic for each known family](layers/comando.md#nwtmn--newtemplates--the-catalog-of-artifact-skeletons-which-kinds-exist-their-headers-their-sections-and-the-presets-that-pick-them) `NWTMN-B06`

- [Every test body carries the scenario code](layers/comando.md#nwtmn--newtemplates--the-catalog-of-artifact-skeletons-which-kinds-exist-their-headers-their-sections-and-the-presets-that-pick-them) `NWTMN-B07`

- [An unknown family gets an instruction, not guessed syntax](layers/comando.md#nwtmn--newtemplates--the-catalog-of-artifact-skeletons-which-kinds-exist-their-headers-their-sections-and-the-presets-that-pick-them) `NWTMN-B08`

- [Each spec preset is an ordered set of catalog sections](layers/comando.md#nwtmn--newtemplates--the-catalog-of-artifact-skeletons-which-kinds-exist-their-headers-their-sections-and-the-presets-that-pick-them) `NWTMN-B09`

- [Product doctrine has its own sections and no layer](layers/comando.md#nwtmn--newtemplates--the-catalog-of-artifact-skeletons-which-kinds-exist-their-headers-their-sections-and-the-presets-that-pick-them) `NWTMN-B10`

- [A unit name becomes snake_case with one separator and whole acronyms](layers/comando.md#nwtmn--newtemplates--the-catalog-of-artifact-skeletons-which-kinds-exist-their-headers-their-sections-and-the-presets-that-pick-them) `NWTMN-B11`

- [The catalog ties rules to what they use](layers/comando.md#nwtmn--newtemplates--the-catalog-of-artifact-skeletons-which-kinds-exist-their-headers-their-sections-and-the-presets-that-pick-them) `NWTMN-B12`

- [The root receives the fifteen operation commands](layers/comando.md#oprgp--opsregister--the-operation-commands-reach-the-cli-through-one-registration-point) `OPRGP-B01`

- [Each operation command is registered exactly once](layers/comando.md#oprgp--opsregister--the-operation-commands-reach-the-cli-through-one-registration-point) `OPRGP-I01`

- [The registration adds commands and nothing else](layers/comando.md#oprgp--opsregister--the-operation-commands-reach-the-cli-through-one-registration-point) `OPRGP-X01`

- [A decision without a date is refused](layers/comando.md#stcms--settingscommand--one-agents-local-decisions-declared-with-a-date-and-kept-out-of-the-projects-configuration) `STCMS-B01`

- [An unknown role or answer is refused naming it](layers/comando.md#stcms--settingscommand--one-agents-local-decisions-declared-with-a-date-and-kept-out-of-the-projects-configuration) `STCMS-B02`

- [Declaring a role records it with the agent and the date](layers/comando.md#stcms--settingscommand--one-agents-local-decisions-declared-with-a-date-and-kept-out-of-the-projects-configuration) `STCMS-B03`

- [The escalated-cards answer is recorded with the agent and the date](layers/comando.md#stcms--settingscommand--one-agents-local-decisions-declared-with-a-date-and-kept-out-of-the-projects-configuration) `STCMS-B04`

- [The role is asked on the terminal when not given](layers/comando.md#stcms--settingscommand--one-agents-local-decisions-declared-with-a-date-and-kept-out-of-the-projects-configuration) `STCMS-B05`

- [An unclear reply is asked again and never assumed](layers/comando.md#stcms--settingscommand--one-agents-local-decisions-declared-with-a-date-and-kept-out-of-the-projects-configuration) `STCMS-B06`

- [Show lists the role's capabilities](layers/comando.md#stcms--settingscommand--one-agents-local-decisions-declared-with-a-date-and-kept-out-of-the-projects-configuration) `STCMS-B07`

- [Show without a role teaches how to declare one](layers/comando.md#stcms--settingscommand--one-agents-local-decisions-declared-with-a-date-and-kept-out-of-the-projects-configuration) `STCMS-B08`

- [A role declaration leaves one source for the escalated-cards question](layers/comando.md#stcms--settingscommand--one-agents-local-decisions-declared-with-a-date-and-kept-out-of-the-projects-configuration) `STCMS-I01`

- [The decisions go to the local settings file](layers/comando.md#stcms--settingscommand--one-agents-local-decisions-declared-with-a-date-and-kept-out-of-the-projects-configuration) `STCMS-X01`

- [A closed input fails instead of assuming no](layers/comando.md#stcms--settingscommand--one-agents-local-decisions-declared-with-a-date-and-kept-out-of-the-projects-configuration) `STCMS-E01`

- [List shows the ids of a state and points pending ones at show](layers/comando.md#sgcms--suggestcommand--the-proposed-fixes-are-listed-shown-applied-or-rejected-and-every-decision-keeps-its-record) `SGCMS-B01`

- [Show prints the reason and the diff, or says it was not found](layers/comando.md#sgcms--suggestcommand--the-proposed-fixes-are-listed-shown-applied-or-rejected-and-every-decision-keeps-its-record) `SGCMS-B02`

- [A dry run only checks the patch](layers/comando.md#sgcms--suggestcommand--the-proposed-fixes-are-listed-shown-applied-or-rejected-and-every-decision-keeps-its-record) `SGCMS-B03`

- [Apply patches the file and then approves with a default reason](layers/comando.md#sgcms--suggestcommand--the-proposed-fixes-are-listed-shown-applied-or-rejected-and-every-decision-keeps-its-record) `SGCMS-B04`

- [Rejecting needs a reason and keeps the record](layers/comando.md#sgcms--suggestcommand--the-proposed-fixes-are-listed-shown-applied-or-rejected-and-every-decision-keeps-its-record) `SGCMS-B05`

- [A failed apply leaves the file and the state as they were](layers/comando.md#sgcms--suggestcommand--the-proposed-fixes-are-listed-shown-applied-or-rejected-and-every-decision-keeps-its-record) `SGCMS-I01`

- [The command decides only suggestions others proposed](layers/comando.md#sgcms--suggestcommand--the-proposed-fixes-are-listed-shown-applied-or-rejected-and-every-decision-keeps-its-record) `SGCMS-X01`

- [A stale patch fails before touching anything](layers/comando.md#sgcms--suggestcommand--the-proposed-fixes-are-listed-shown-applied-or-rejected-and-every-decision-keeps-its-record) `SGCMS-E01`

- [Outside git the failure names the suggestion patch](layers/comando.md#sgcms--suggestcommand--the-proposed-fixes-are-listed-shown-applied-or-rejected-and-every-decision-keeps-its-record) `SGCMS-E02`

- [Only github mode is accepted and the first PR is required](layers/comando.md#sycms--synthesizecommand--two-pull-requests-in-content-conflict-become-one-card-that-asks-for-the-best-of-each-and-every-end-points-at-it) `SYCMS-B01`

- [A leading # on a PR number is dropped](layers/comando.md#sycms--synthesizecommand--two-pull-requests-in-content-conflict-become-one-card-that-asks-for-the-best-of-each-and-every-end-points-at-it) `SYCMS-B02`

- [The card cites both PRs and cards and asks for the best of each](layers/comando.md#sycms--synthesizecommand--two-pull-requests-in-content-conflict-become-one-card-that-asks-for-the-best-of-each-and-every-end-points-at-it) `SYCMS-B03`

- [A one-sided card says the other side is missing and where to look](layers/comando.md#sycms--synthesizecommand--two-pull-requests-in-content-conflict-become-one-card-that-asks-for-the-best-of-each-and-every-end-points-at-it) `SYCMS-B04`

- [The card is labelled for the board and under each origin card](layers/comando.md#sycms--synthesizecommand--two-pull-requests-in-content-conflict-become-one-card-that-asks-for-the-best-of-each-and-every-end-points-at-it) `SYCMS-B05`

- [Each PR is commented and closed and each origin card is pointed at the new card](layers/comando.md#sycms--synthesizecommand--two-pull-requests-in-content-conflict-become-one-card-that-asks-for-the-best-of-each-and-every-end-points-at-it) `SYCMS-B06`

- [Each failed link is a warning and the command succeeds](layers/comando.md#sycms--synthesizecommand--two-pull-requests-in-content-conflict-become-one-card-that-asks-for-the-best-of-each-and-every-end-points-at-it) `SYCMS-B07`

- [The closing line names only the PRs actually closed](layers/comando.md#sycms--synthesizecommand--two-pull-requests-in-content-conflict-become-one-card-that-asks-for-the-best-of-each-and-every-end-points-at-it) `SYCMS-B08`

- [A dry run shows the card and changes nothing](layers/comando.md#sycms--synthesizecommand--two-pull-requests-in-content-conflict-become-one-card-that-asks-for-the-best-of-each-and-every-end-points-at-it) `SYCMS-B09`

- [The card follows the project language](layers/comando.md#sycms--synthesizecommand--two-pull-requests-in-content-conflict-become-one-card-that-asks-for-the-best-of-each-and-every-end-points-at-it) `SYCMS-B10`

- [No PR is closed unless the synthesis card exists](layers/comando.md#sycms--synthesizecommand--two-pull-requests-in-content-conflict-become-one-card-that-asks-for-the-best-of-each-and-every-end-points-at-it) `SYCMS-I01`

- [The card never picks a side](layers/comando.md#sycms--synthesizecommand--two-pull-requests-in-content-conflict-become-one-card-that-asks-for-the-best-of-each-and-every-end-points-at-it) `SYCMS-X01`

- [A card that cannot be opened fails with the host's message](layers/comando.md#sycms--synthesizecommand--two-pull-requests-in-content-conflict-become-one-card-that-asks-for-the-best-of-each-and-every-end-points-at-it) `SYCMS-E01`

- [The plan is the timed files fastest first, then the untimed](layers/comando.md#bdgrn--budgetrun--run-a-suites-files-fastest-first-until-a-time-budget-is-spent) `BDGRN-B01`

- [Batches take what fits and then the untimed one by one](layers/comando.md#bdgrn--budgetrun--run-a-suites-files-fastest-first-until-a-time-budget-is-spent) `BDGRN-B02`

- [A batch still running at the deadline is stopped with its group](layers/comando.md#bdgrn--budgetrun--run-a-suites-files-fastest-first-until-a-time-budget-is-spent) `BDGRN-B03`

- [A budget runs batches fastest first and reports what ran and what was left](layers/comando.md#bdgrn--budgetrun--run-a-suites-files-fastest-first-until-a-time-budget-is-spent) `BDGRN-B04`

- [A mutation budget runs one file per batch and records its time](layers/comando.md#bdgrn--budgetrun--run-a-suites-files-fastest-first-until-a-time-budget-is-spent) `BDGRN-B05`

- [A failed batch does not stop the budget](layers/comando.md#bdgrn--budgetrun--run-a-suites-files-fastest-first-until-a-time-budget-is-spent) `BDGRN-B06`

- [What the budget cannot run is refused before anything runs](layers/comando.md#bdgrn--budgetrun--run-a-suites-files-fastest-first-until-a-time-budget-is-spent) `BDGRN-B07`

- [The plan holds only the suite's own files](layers/comando.md#bdgrn--budgetrun--run-a-suites-files-fastest-first-until-a-time-budget-is-spent) `BDGRN-B08`

- [A budget without a map is refused saying to build it](layers/comando.md#bdgrn--budgetrun--run-a-suites-files-fastest-first-until-a-time-budget-is-spent) `BDGRN-E02`

- [The full sweep confronts every node of the map](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B01`

- [The incremental check confronts only the impact path of the change](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B02`

- [The tests that stamp a changed module enter its impact path](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B03`

- [The pieces of the changed file's unit enter its impact path, for a Go unit as for a TypeScript one](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B70`

- [Changed paths are normalised to the map's form](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B04`

- [Files the project does not govern are recognised as not governed](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B05`

- [An ungoverned file does not taint a batch](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B06`

- [A plan's progress companion is not governed](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B67`

- [Phase and category select the gates charged](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B07`

- [A gate without skip_on runs in both perspectives](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B08`

- [A gate that skips the change perspective runs only on the full sweep](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B09`

- [A gate that skips the full sweep runs only on changed files](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B10`

- [A gate that skips both perspectives is switched off](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B11`

- [A commit message marker with a reason waives a gate](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B12`

- [A waiver in the environment drops the gate and says why](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B13`

- [The deterministic mode drops the judgment gates](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B14`

- [The issue policy follows the workflow mode](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B15`

- [The confronted edges are stamped at the current revisions](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B16`

- [The no-record mode leaves the map untouched](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B17`

- [A blocking failure is filed only when issues are on](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B18`

- [Passes resolve, decisions and debts open in their folders](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B19`

- [The full check closes the violations it did not reproduce](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B20`

- [A record that fails is warned about and the check still reports](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B72`

- [When the check writes no issue it says why](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B75`

- [In github mode the issues go to the board, never to the local folders](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B76`

- [The record summary counts what the record did](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B77`

- [A pending judgment becomes one task in the local queue](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B21`

- [A judge task suggests the review stage, a verb the work command composes](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B69`

- [A queued judgment bars the incremental check only](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B22`

- [The check says how many judgments it queued, and nothing when it queued none](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B73`

- [Stale judge tasks leave the queue](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B23`

- [An incremental check keeps the judgments it did not look at](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B24`

- [A judge task is read back as the gate that queued it even when another gate's name prefixes it](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B68`

- [In github mode the brief is printed without recording and nothing is queued](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B25`

- [The judgment brief names every target](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B26`

- [The judgment brief carries the question and the guide](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B27`

- [The judgment brief groups the targets by gate](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B28`

- [The judgment brief says who judges](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B29`

- [The judgment brief lists ten targets per gate and counts the rest](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B30`

- [With nothing awaiting judgment there is no judgment brief](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B71`

- [Governed files missing from the map make a stale-map warning](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B31`

- [Nodes edited after the map build are warned about on the incremental check](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B32`

- [A map written by another version is warned about](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B33`

- [A map with no writer version raises no warning](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B34`

- [A declared gate with nothing to measure is named](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B35`

- [The local backlog is printed on the full local sweep only](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B36`

- [Governance tips appear on the full sweep only](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B37`

- [With no governance tip the full sweep prints neither a tip nor the pointer to the doctor](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B74`

- [The report is mirrored to a file](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B38`

- [A blocking failure exits 1 with the mirror complete](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B39`

- [The name column fits the longest name](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B40`

- [Each counter column has its own width](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B41`

- [The always-present columns are at least one wide](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B42`

- [Without drift the drift column does not exist](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B43`

- [The indeterminate counter is the skipped and pending less the drift](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B78`

- [An empty drift cell is measured in terminal columns](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B44`

- [A clean gate has nothing pending of any kind](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B45`

- [The default table shows the clean gates](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B46`

- [Only-issues omits the clean gates and counts them](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B47`

- [Show-drift lists every drift item](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B48`

- [Show-drift is not cut by the size of the scan](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B49`

- [Without show-drift only the counter appears](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B50`

- [A small scan lists the reason of each skip](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B51`

- [A large scan does not list the skip reasons](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B52`

- [The legend explains only the symbols used](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B53`

- [A repeated drift reason is written once with its targets](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B54`

- [Distinct drift reasons stay target by target](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B55`

- [The drift heading counts items and gates](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B56`

- [The findings heading counts every kind](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B57`

- [Each failure is listed with its detail, and no finding means no findings heading](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B82`

- [The verdict line says what is still open](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B58`

- [A clean informative gate is named as ready to become blocking](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B83`

- [Occurrences of a detail are printed one per line](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B59`

- [Two occurrences already break](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B60`

- [A single occurrence stays whole](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B61`

- [A glued semicolon does not break](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B62`

- [A long list of items breaks one per line](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B63`

- [A short sentence with commas stays whole](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B64`

- [Long prose with commas stays whole](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B65`

- [A list of paths still breaks](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B66`

- [The time table counts each gate's targets and aligns its columns](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B79`

- [The slowest targets are listed slowest first, and only when a time was recorded](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B80`

- [Times are rounded to what a decision needs](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B81`

- [Declaring a perspective does not change the cost axis](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-I01`

- [The version warning compares names by equality](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-I02`

- [The skip column does not move with or without drift](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-I03`

- [The stale-map warning does not bar the check](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-X01`

- [The check without configuration fails](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-E01`

- [A project with no gate has no pipeline](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-E02`

- [The check without a map points at the map build](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-E03`

- [A waiver without a reason is refused](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-E04`

- [The check with no scope is refused](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-E05`

- [A governed file outside the map bars the check](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-E06`

- [A path on neither disk nor map is an error](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-E07`

- [The check's stamps do not erase what another process wrote meanwhile](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B84`

- [The index flag reads what the commit records](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B85`

- [--index with no scope judges the staged files](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B86`

- [Under --index the date is judged by what the commit records, in every staging state](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B87`

- [Under --index a file the index does not have is not judged](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B88`

- [A check with no phase leaves out a gate declared for manual alone](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B89`

- [The check says how many targets of its scope are to review](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B90`

- [A full sweep names the catalog gates missing over the declared layers](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B91`

- [A gate this run leaves out is still declared](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B92`

- [The catalog line leaves out a gate that would only wait for its premise](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B93`

- [The files check --fix repaired keep their evidence, the line-level signals only when no line moved](layers/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile) `CGPCH-B94`

- [The scenarios of one spec are listed as proven or not](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-B01`

- [A code another unit owns, cited in prose, is not a declared scenario](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-B02`

- [A spec with no unit code keeps every code it declares](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-B03`

- [A spec changed since ingestion flags its signal as stale](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-B04`

- [A spec with no scenario code has nothing to cover](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-B05`

- [The panorama answers by scenario, by line and by mutation](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-B06`

- [Every measured file above the threshold is said with the count, survivors included](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-B07`

- [Each panorama section shows at most fifteen entries and counts the rest](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-B08`

- [The changed lines are crossed with the coverage report](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-B09`

- [A diff with no instrumented line has nothing to cover](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-B10`

- [Changed lines covered below the threshold fail the diff coverage](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-B11`

- [The delta with no drop says so and counts the improvements](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-B12`

- [A file that lost line coverage fails the delta](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-B13`

- [The diff coverage lists the changed files in path order, and resolves a path matching several coverage entries always to the same one](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-B14`

- [Nothing measured is never reported as nothing below the threshold](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-I01`

- [The diff coverage needs no map](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-X01`

- [The coverage command without a map points at the map build](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-E01`

- [A spec outside the map is refused](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-E02`

- [The diff coverage without a coverage report is refused](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-E03`

- [A coverage report that cannot be read fails the diff coverage](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-E04`

- [A git diff outside a repository explains why and offers the diff file](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-E05`

- [A diff file that cannot be read fails the diff coverage](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-E06`

- [The report counts a rule proven by all its variants](layers/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta) `CVCMC-B15`

- [A page left out of date is compiled and staged](layers/comando.md#dcsyn--docssyncforcommit--the-commit-carries-the-pages-its-specs-produce) `DCSYN-B01`

- [With the tree ahead, the page goes straight into the index](layers/comando.md#dcsyn--docssyncforcommit--the-commit-carries-the-pages-its-specs-produce) `DCSYN-B02`

- [Nothing to compile, nothing staged](layers/comando.md#dcsyn--docssyncforcommit--the-commit-carries-the-pages-its-specs-produce) `DCSYN-B03`

- [A page written by hand, or one git ignores, is never touched](layers/comando.md#dcsyn--docssyncforcommit--the-commit-carries-the-pages-its-specs-produce) `DCSYN-X01`

- [Nothing outside docs is written or staged](layers/comando.md#dcsyn--docssyncforcommit--the-commit-carries-the-pages-its-specs-produce) `DCSYN-X02`

- [After the sync the gate finds the pages up to date](layers/comando.md#dcsyn--docssyncforcommit--the-commit-carries-the-pages-its-specs-produce) `DCSYN-I01`

- [A template that does not compile comes back as an error](layers/comando.md#dcsyn--docssyncforcommit--the-commit-carries-the-pages-its-specs-produce) `DCSYN-E01`

- [The diagnosis is printed grouped by the check that found it](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-B01`

- [A group with any warning is marked as a warning](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-B02`

- [A report with no finding says the ecosystem is sound](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-B03`

- [The fix outside the github mode does nothing](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-B04`

- [The fix in github mode seeds the workflow pipelines](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-B05`

- [The fix protects the declared branches and skips one that does not exist](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-B06`

- [The fix disables an approval requirement the author cannot satisfy](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-B07`

- [The fix ensures the state labels and says the board is optional](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-B08`

- [The protection body carries every required field and the approvals](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-B09`

- [The pipelines check in local mode has nothing to check](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-B10`

- [Missing pipelines are named and CI continues by default](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-B11`

- [A project that declared stale pipelines as blocking fails the check](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-B12`

- [The pipelines check reads the project named by the root flag](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-B13`

- [The diagnosis never fails the doctor](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-I01`

- [Running the fix twice changes nothing the second time](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-I02`

- [The fix never creates a board](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-X01`

- [The orphan local queue and delivery records are warned about and kept](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-X02`

- [The doctor without configuration fails](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-E01`

- [The doctor without a map fails](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-E02`

- [The fix without gh refuses before seeding anything](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-E03`

- [The fix with gh not logged in refuses before seeding anything](layers/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment) `HLDCH-E04`

- [The evidence moves to the current content with the reason](layers/comando.md#kpevd--keepevidence--a-change-that-proves-nothing-new-keeps-the-files-evidence) `KPEVD-B01`

- [A rebuilt map's lost evidence comes from HEAD](layers/comando.md#kpevd--keepevidence--a-change-that-proves-nothing-new-keeps-the-files-evidence) `KPEVD-B02`

- [Nothing to keep is said](layers/comando.md#kpevd--keepevidence--a-change-that-proves-nothing-new-keeps-the-files-evidence) `KPEVD-B03`

- [The contract stamps are refreshed under the same declaration](layers/comando.md#kpevd--keepevidence--a-change-that-proves-nothing-new-keeps-the-files-evidence) `KPEVD-B04`

- [A missing reason, file or map fails](layers/comando.md#kpevd--keepevidence--a-change-that-proves-nothing-new-keeps-the-files-evidence) `KPEVD-E01`

- [The suites a later run did not replace come back from HEAD and are carried with the rest](layers/comando.md#kpevd--keepevidence--a-change-that-proves-nothing-new-keeps-the-files-evidence) `KPEVD-B05`

- [The issues in todo and doing are counted, with the user-owned ones apart](layers/comando.md#lcbcl--localbacklog--what-is-still-open-locally-after-a-full-check-said-in-two-lines) `LCBCL-B01`

- [The pending and claimed tasks are counted](layers/comando.md#lcbcl--localbacklog--what-is-still-open-locally-after-a-full-check-said-in-two-lines) `LCBCL-B02`

- [A project with nothing open prints no backlog](layers/comando.md#lcbcl--localbacklog--what-is-still-open-locally-after-a-full-check-said-in-two-lines) `LCBCL-B03`

- [Only the side that has something open gets its line](layers/comando.md#lcbcl--localbacklog--what-is-still-open-locally-after-a-full-check-said-in-two-lines) `LCBCL-B04`

- [User-owned and past-window counts alone do not make a backlog](layers/comando.md#lcbcl--localbacklog--what-is-still-open-locally-after-a-full-check-said-in-two-lines) `LCBCL-I01`

- [Reading the backlog changes no issue and no task](layers/comando.md#lcbcl--localbacklog--what-is-still-open-locally-after-a-full-check-said-in-two-lines) `LCBCL-X01`

- [An issue folder that cannot be listed counts as zero](layers/comando.md#lcbcl--localbacklog--what-is-still-open-locally-after-a-full-check-said-in-two-lines) `LCBCL-E01`

- [The committed map is the one a build of the commit makes](layers/comando.md#mpsyn--mapsyncforcommit--the-commit-carries-the-map-a-build-of-the-commit-makes) `MPSYN-B01`

- [A dated file keeps its proofs](layers/comando.md#mpsyn--mapsyncforcommit--the-commit-carries-the-map-a-build-of-the-commit-makes) `MPSYN-B02`

- [The index, not the tree](layers/comando.md#mpsyn--mapsyncforcommit--the-commit-carries-the-map-a-build-of-the-commit-makes) `MPSYN-B03`

- [An untracked map is left alone](layers/comando.md#mpsyn--mapsyncforcommit--the-commit-carries-the-map-a-build-of-the-commit-makes) `MPSYN-B04`

- [A map that cannot be written gives the error back](layers/comando.md#mpsyn--mapsyncforcommit--the-commit-carries-the-map-a-build-of-the-commit-makes) `MPSYN-E01`

- [A file the commit does not change keeps the proofs HEAD had](layers/comando.md#mpsyn--mapsyncforcommit--the-commit-carries-the-map-a-build-of-the-commit-makes) `MPSYN-B05`

- [With the tree ahead of the commit, the map on disk stays the tree's](layers/comando.md#mpsyn--mapsyncforcommit--the-commit-carries-the-map-a-build-of-the-commit-makes) `MPSYN-B06`

- [What was measured after the last commit survives the next one](layers/comando.md#mpsyn--mapsyncforcommit--the-commit-carries-the-map-a-build-of-the-commit-makes) `MPSYN-B07`

- [The quality domain registers exactly its twelve commands](layers/comando.md#qlcmq--qualitycommands--the-quality-domain-puts-its-twelve-commands-under-the-root-command) `QLCMQ-B01`

- [Each quality command is reachable by its name](layers/comando.md#qlcmq--qualitycommands--the-quality-domain-puts-its-twelve-commands-under-the-root-command) `QLCMQ-B02`

- [No two quality commands share a name](layers/comando.md#qlcmq--qualitycommands--the-quality-domain-puts-its-twelve-commands-under-the-root-command) `QLCMQ-I01`

- [Registering the quality commands prints nothing](layers/comando.md#qlcmq--qualitycommands--the-quality-domain-puts-its-twelve-commands-under-the-root-command) `QLCMQ-X01`

- [A single perspective is written to docs or to the chosen file](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-B01`

- [The all command writes every perspective and an index into docs/anchors](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-B02`

- [Every perspective of a configured project opens with the same header and closes with a footer](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-B03`

- [The tests perspective merges execution by layer and warns about failures and stale signals](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-B04`

- [The tests perspective counts only measured specs and lists the unproven scenarios](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-B05`

- [The tests perspective lists the files below 70% of lines and the coverage regressions](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-B06`

- [The tests perspective with nothing ingested says so in each section](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-B07`

- [The quality perspective gives the verdict per gate and the divergences](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-B08`

- [The structure perspective counts nodes by kind, the governance and the identity findings](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-B09`

- [The configuration perspective lists what anchors.yaml declares and what it misses](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-B10`

- [The issues perspective splits the open issues by who must act, and lists the tasks](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-B11`

- [The inconsistencies perspective lists every health finding by check and the failing gates](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-B12`

- [A finding section takes the findings of its check, caps the list at 25, and is absent when empty](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-B13`

- [Without anchors.yaml the perspectives say what is missing instead of inventing](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-B14`

- [An issue waiting on the user is listed only as the user's](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-I01`

- [Generating the reports leaves the map as it was](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-X01`

- [The reports without a map point at the map build](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-E01`

- [A destination that cannot be written fails the report](layers/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs) `RPRTS-E02`

- [A file is placed in one of four boxes](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-B01`

- [The default takes stale below the minimum and never measured, and each flag opens a side](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-B02`

- [A test file is stale through what it exercises and passes by its own layer](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-B03`

- [A mutation result under load counts as stale and passes at the floor](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-B04`

- [A test file belongs to the suite that ran it](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-B05`

- [Support files and no_signal targets never run](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-B06`

- [The run says what it selected and what it left out](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-B07`

- [A suite without run_changed runs whole](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-B08`

- [A selection that takes nothing runs nothing](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-B09`

- [The selected files run in as few batches as the ceiling allows](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-B10`

- [--all runs whole and does not combine with the other choices](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-B11`

- [A suite with paths is handed only its own files](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-B12`

- [A selective run without a map is refused](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-E01`

- [A file the gate does not confront is not run for it](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-B13`

- [The tests of the rules a changed contract field reaches are taken](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-B14`

- [A file edited since the map was built is read as it is now](layers/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise) `SLCTN-B15`

- [Expired test evidence is listed before the stale edges](layers/comando.md#steds--staleedges--lists-the-confrontation-debt-expired-test-evidence-and-stale-edges) `STEDS-B01`

- [Each expired evidence names why it expired](layers/comando.md#steds--staleedges--lists-the-confrontation-debt-expired-test-evidence-and-stale-edges) `STEDS-B02`

- [Stale edges are split into never validated and drifted](layers/comando.md#steds--staleedges--lists-the-confrontation-debt-expired-test-evidence-and-stale-edges) `STEDS-B03`

- [A map with every edge validated prints the clean message](layers/comando.md#steds--staleedges--lists-the-confrontation-debt-expired-test-evidence-and-stale-edges) `STEDS-B04`

- [An edge stamped at the current revisions is never listed](layers/comando.md#steds--staleedges--lists-the-confrontation-debt-expired-test-evidence-and-stale-edges) `STEDS-I01`

- [Listing the stale edges leaves the map unchanged](layers/comando.md#steds--staleedges--lists-the-confrontation-debt-expired-test-evidence-and-stale-edges) `STEDS-X01`

- [The stale command without a map points at the map build](layers/comando.md#steds--staleedges--lists-the-confrontation-debt-expired-test-evidence-and-stale-edges) `STEDS-E01`

- [With no test named, every test of the map is considered](layers/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change) `CNSTC-B01`

- [The missing stamp is written above the double and totalled](layers/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change) `CNSTC-B02`

- [The dry run says what it would write and writes nothing](layers/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change) `CNSTC-B03`

- [The refresh lists the doubles stamped against the old version with the block change](layers/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change) `CNSTC-B04`

- [The refresh of a file no double is stamped against says so](layers/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change) `CNSTC-B05`

- [A refreshed block with no HEAD version says there is nothing to compare](layers/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change) `CNSTC-B06`

- [The refresh in dry-run lists the doubles and writes nothing](layers/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change) `CNSTC-B07`

- [The block diff lists what left and what came, counting repeats](layers/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change) `CNSTC-B08`

- [Stamping twice writes nothing the second time](layers/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change) `CNSTC-I01`

- [A divergent stamp is left as it is](layers/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change) `CNSTC-X01`

- [The stamp command without configuration fails](layers/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change) `CNSTC-E01`

- [The stamp command without a map points at the map build](layers/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change) `CNSTC-E02`

- [A stamp whose anchor is gone is reported and not refreshed](layers/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change) `CNSTC-E03`

- [A test of the map that cannot be read is reported and skipped](layers/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change) `CNSTC-E04`

- [The refresh handed a test names the modules to refresh](layers/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change) `CNSTC-B09`

- [A directory with no git repository stops at git init](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-B01`

- [Without the git binary status warns and goes on](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-B02`

- [A project with neither PROJECT.md nor configuration is not started](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-B03`

- [A project with PROJECT.md and no configuration is sent to init](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-B04`

- [A configured project with no map is sent to the map build](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-B05`

- [The configuration and the map are counted, and clean informative gates are named](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-B06`

- [The local queue names the first pending step in the cycle's order](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-B07`

- [An assembled local project with no work is sent to the first plan](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-B08`

- [The github queue with missing pipelines stops at the doctor fix](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-B09`

- [The github queue states the pull-request flow and the protected branches](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-B10`

- [The agent's own open cards come before claiming new work](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-B11`

- [Without an agent identity the next step is to claim work](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-B12`

- [A github project with no real work is sent to the first plan](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-B13`

- [Status never names a step beyond the first one missing](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-I01`

- [Status leaves the project as it found it](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-X01`

- [A configuration that does not load fails the status](layers/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step) `PRSTP-E01`

- [The selected suite runs at the root under a header naming it](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-B01`

- [The report this run wrote is ingested into the map](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-B02`

- [A passing suite with no report says nothing was ingested](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-B03`

- [A failing suite is still ingested and stops the rest](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-B04`

- [Without changed files the full command runs](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-B05`

- [The incremental command receives the impact path where it declares it](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-B06`

- [Tests get code and tests, mutation gets only code](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-B07`

- [The target fills the placeholder and is ignored without one](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-B08`

- [A passing run chains coverage, and an empty chain does nothing](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-B09`

- [An impact path with no code file runs nothing](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-B10`

- [A passing run chains the check over the suite's own scope](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-B11`

- [A report stamped by a coarse clock just before the start is this run's](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-B13`

- [A report older than the run is never ingested](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-I01`

- [The declared command runs as declared](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-X01`

- [The suite commands without configuration fail](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-E01`

- [A section with no suite shows how to declare it and fails](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-E02`

- [A filter that names nothing declared is refused with what is declared](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-E03`

- [Declared filters that match no suite together are refused](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-E04`

- [A target placeholder without a target is refused](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-E05`

- [The incremental mode without an incremental command is refused](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-E06`

- [A command line over the ceiling is refused before running](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-E07`

- [A chain naming another command is refused](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-E08`

- [The incremental mode without a map points at the map build](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-E09`

- [An incremental run hands each suite only its own impact files](layers/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map) `STPRS-B12`

- [Only a date inside the header at the top of the file is bumped](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-B01`

- [A real change and a new file are bumped to the date](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-B02`

- [An unchanged file, a date-only change and a file already at the date are not bumped](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-B03`

- [Without the staged flag the candidates are the worktree changes and the untracked files](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-B04`

- [With the staged flag the index is dated and re-staged](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-B05`

- [A staged file with changes outside the index is skipped](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-B06`

- [In a partial commit the real index is dated too](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-B07`

- [The exclude globs of the flag and of the configuration add up](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-B08`

- [The dry run says what it would bump and writes nothing](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-B09`

- [Each bump and each skip is listed with its reason, then the total](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-B10`

- [Without a date the day of the run is written](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-B11`

- [The pre-commit bump is on unless the project turns it off](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-B12`

- [A project root below the repository top dates its own files and nothing outside it](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-B13`

- [Touching twice bumps nothing the second time](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-I01`

- [A changed file with no dated header is neither touched nor listed](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-X01`

- [Outside a git repository touch fails](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-E01`

- [A changed file that cannot be read is skipped and named](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-E02`

- [Named files narrow the touch to them](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-B14`

- [Every named file gets a verdict](layers/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those) `HDTHD-B15`

- [The staged scope is the added, copied, modified and renamed files of the index](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-B01`

- [Nothing staged has nothing to verify](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-B02`

- [Verify hands the files to check in a child process](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-B03`

- [An automatic phase asks check for computable gates and issues only](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-B04`

- [The manual phase, or no phase, asks check for the full report](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-B05`

- [The pre-commit over the index dates the staged files first](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-B06`

- [A project that turned pre-commit dating off gets no dating](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-B07`

- [A staged file with changes outside the index is named, not dated](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-B08`

- [The facade's flags reach check unchanged](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-B09`

- [A project root below the repository top hands check its own staged files, by its own paths](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-B10`

- [Verify never asks check for both the full sweep and a file list](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-I01`

- [Verify holds no verdict of its own](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-X01`

- [The child's not-governed exit stays not-governed](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-E01`

- [Any other failing exit of the child stays a failure](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-E02`

- [Verify with no scope is refused](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-E03`

- [The staged scope outside a git repository is refused](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-E04`

- [A dating failure warns and does not stop the verify](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-E05`

- [Over the index the check reads the index](layers/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline) `VPFVR-B11`

- [A frozen project refuses the command, naming it and the reason](layers/comando.md#clrtc--cliroot--every-command-passes-through-one-root-that-speaks-the-projects-language-and-honours-the-freeze) `CLRTC-B01`

- [The commands that only read, and thaw, run while frozen](layers/comando.md#clrtc--cliroot--every-command-passes-through-one-root-that-speaks-the-projects-language-and-honours-the-freeze) `CLRTC-B02`

- [A subcommand of an allowed command runs while frozen](layers/comando.md#clrtc--cliroot--every-command-passes-through-one-root-that-speaks-the-projects-language-and-honours-the-freeze) `CLRTC-B03`

- [A missing or broken config is not frozen](layers/comando.md#clrtc--cliroot--every-command-passes-through-one-root-that-speaks-the-projects-language-and-honours-the-freeze) `CLRTC-B04`

- [The root flag decides which project is read](layers/comando.md#clrtc--cliroot--every-command-passes-through-one-root-that-speaks-the-projects-language-and-honours-the-freeze) `CLRTC-B05`

- [The project's top-level lang is applied before the command runs](layers/comando.md#clrtc--cliroot--every-command-passes-through-one-root-that-speaks-the-projects-language-and-honours-the-freeze) `CLRTC-B06`

- [The root prints neither the error nor the usage](layers/comando.md#clrtc--cliroot--every-command-passes-through-one-root-that-speaks-the-projects-language-and-honours-the-freeze) `CLRTC-B07`

- [Every command the work guide and the pipelines teach is registered](layers/comando.md#clrtc--cliroot--every-command-passes-through-one-root-that-speaks-the-projects-language-and-honours-the-freeze) `CLRTC-I01`

- [The freeze refusal is written in the project language](layers/comando.md#clrtc--cliroot--every-command-passes-through-one-root-that-speaks-the-projects-language-and-honours-the-freeze) `CLRTC-X01`

- [The language is read from a CRLF file](layers/comando.md#clrtc--cliroot--every-command-passes-through-one-root-that-speaks-the-projects-language-and-honours-the-freeze) `CLRTC-B08`

- [Every command that takes several files reads them the same way](layers/comando.md#clrtc--cliroot--every-command-passes-through-one-root-that-speaks-the-projects-language-and-honours-the-freeze) `CLRTC-I02`

## config

- [A known extension answers with its line-comment prefixes](layers/config.md#cmmrc-b01--a-known-extension-answers-with-its-line-comment-prefixes) `CMMRC-B01`

- [An unknown extension answers with no prefix](layers/config.md#cmmrc-b02--an-unknown-extension-answers-with-no-prefix) `CMMRC-B02`

- [The line comment of a path follows its last extension, ignoring case](layers/config.md#cmmrc-b03--the-line-comment-of-a-path-follows-its-last-extension-ignoring-case) `CMMRC-B03`

- [A markup file gets the opening of a block comment](layers/config.md#cmmrc-b04--a-markup-file-gets-the-opening-of-a-block-comment) `CMMRC-B04`

- [A path with no extension or an unknown one gets the hash comment](layers/config.md#cmmrc-b05--a-path-with-no-extension-or-an-unknown-one-gets-the-hash-comment) `CMMRC-B05`

- [An extension with several prefixes gets the first one declared](layers/config.md#cmmrc-b06--an-extension-with-several-prefixes-gets-the-first-one-declared) `CMMRC-B06`

- [For every extension of the table, the line comment is the table's first prefix](layers/config.md#cmmrc-i01--for-every-extension-of-the-table-the-line-comment-is-the-tables-first-prefix) `CMMRC-I01`

- [The prefix lookup takes the extension as given and does not normalise it](layers/config.md#cmmrc-x01--the-prefix-lookup-takes-the-extension-as-given-and-does-not-normalise-it) `CMMRC-X01`

- [An unknown key in the file is a load error naming the key](layers/config.md#cnfgo-b01--an-unknown-key-in-the-file-is-a-load-error-naming-the-key) `CNFGO-B01`

- [An unknown key gets the misspelling and old-binary hypotheses](layers/config.md#cnfgo-b02--an-unknown-key-gets-the-misspelling-and-old-binary-hypotheses) `CNFGO-B02`

- [A renamed key in an older-format file advises the migration](layers/config.md#cnfgo-b03--a-renamed-key-in-an-older-format-file-advises-the-migration) `CNFGO-B03`

- [The migration advice needs both an older format and a renamed key](layers/config.md#cnfgo-b04--the-migration-advice-needs-both-an-older-format-and-a-renamed-key) `CNFGO-B04`

- [An error that is not an unknown key gets no version hint](layers/config.md#cnfgo-b05--an-error-that-is-not-an-unknown-key-gets-no-version-hint) `CNFGO-B05`

- [Two gates sharing an ID fail the load](layers/config.md#cnfgo-b06--two-gates-sharing-an-id-fail-the-load) `CNFGO-B06`

- [An unknown scope, full scope, cost, phase or perspective value fails the load](layers/config.md#cnfgo-b07--an-unknown-scope-full-scope-cost-phase-or-perspective-value-fails-the-load) `CNFGO-B07`

- [Local and manual modes refuse the GitHub fields](layers/config.md#cnfgo-b08--local-and-manual-modes-refuse-the-github-fields) `CNFGO-B08`

- [GitHub mode requires an owner/name repository and a label](layers/config.md#cnfgo-b09--github-mode-requires-an-ownername-repository-and-a-label) `CNFGO-B09`

- [An unknown workflow mode fails with no fallback](layers/config.md#cnfgo-b10--an-unknown-workflow-mode-fails-with-no-fallback) `CNFGO-B10`

- [The tests source is one source with a pattern that compiles](layers/config.md#cnfgo-b45--the-tests-source-is-one-source-with-a-pattern-that-compiles) `CNFGO-B45`

- [A test level's code filter accepts by allow and refuses by exclude](layers/config.md#cnfgo-b44--a-test-levels-code-filter-accepts-by-allow-and-refuses-by-exclude) `CNFGO-B44`

- [A declared pattern that does not compile fails the load naming the field](layers/config.md#cnfgo-b11--a-declared-pattern-that-does-not-compile-fails-the-load-naming-the-field) `CNFGO-B11`

- [An unsupported language fails the load](layers/config.md#cnfgo-b12--an-unsupported-language-fails-the-load) `CNFGO-B12`

- [Code lengths outside two to eight fail, and valid ones reach the engine](layers/config.md#cnfgo-b13--code-lengths-outside-two-to-eight-fail-and-valid-ones-reach-the-engine) `CNFGO-B13`

- [A file with no code lengths restores the default, whatever an earlier load set](layers/config.md#cnfgo-b42--a-file-with-no-code-lengths-restores-the-default-whatever-an-earlier-load-set) `CNFGO-B42`

- [Every load refusal and the header Save writes are in the project's language](layers/config.md#cnfgo-b43--every-load-refusal-and-the-header-save-writes-are-in-the-projects-language) `CNFGO-B43`

- [A canonical gate inherits every field the project omitted](layers/config.md#cnfgo-b14--a-canonical-gate-inherits-every-field-the-project-omitted) `CNFGO-B14`

- [A field the project declared wins over the canonical one](layers/config.md#cnfgo-b15--a-field-the-project-declared-wins-over-the-canonical-one) `CNFGO-B15`

- [Declaring run or check inherits neither of the pair](layers/config.md#cnfgo-b16--declaring-run-or-check-inherits-neither-of-the-pair) `CNFGO-B16`

- [A gate the catalog does not know loads untouched](layers/config.md#cnfgo-b17--a-gate-the-catalog-does-not-know-loads-untouched) `CNFGO-B17`

- [A gate with no declared severity does not block](layers/config.md#cnfgo-b18--a-gate-with-no-declared-severity-does-not-block) `CNFGO-B18`

- [The scope defaults to one run per target, and the full scan uses scope_full only when it is batch or project](layers/config.md#cnfgo-b19--the-scope-defaults-to-one-run-per-target-and-the-full-scan-uses-scope-full-only-when-it-is-batch-or-project) `CNFGO-B19`

- [A gate with no phases runs in every phase](layers/config.md#cnfgo-b20--a-gate-with-no-phases-runs-in-every-phase) `CNFGO-B20`

- [A gate participates in every perspective unless skip_on excludes it](layers/config.md#cnfgo-b21--a-gate-participates-in-every-perspective-unless-skip-on-excludes-it) `CNFGO-B21`

- [The mutation report format comes from the mutation-score gate, normalized, defaulting to the canonical format](layers/config.md#cnfgo-b22--the-mutation-report-format-comes-from-the-mutation-score-gate-normalized-defaulting-to-the-canonical-format) `CNFGO-B22`

- [Section language is checked unless the gate turns it off](layers/config.md#cnfgo-b23--section-language-is-checked-unless-the-gate-turns-it-off) `CNFGO-B23`

- [The integration branch defaults to main, and a non-main integration branch protects main too](layers/config.md#cnfgo-b24--the-integration-branch-defaults-to-main-and-a-non-main-integration-branch-protects-main-too) `CNFGO-B24`

- [One approval is required unless the project declares another number, zero included](layers/config.md#cnfgo-b25--one-approval-is-required-unless-the-project-declares-another-number-zero-included) `CNFGO-B25`

- [A stale pipeline or a manual ingest blocks only when the project asks](layers/config.md#cnfgo-b26--a-stale-pipeline-or-a-manual-ingest-blocks-only-when-the-project-asks) `CNFGO-B26`

- [Only an explicit enabled false freezes the project](layers/config.md#cnfgo-b27--only-an-explicit-enabled-false-freezes-the-project) `CNFGO-B27`

- [The freeze reason is shown trimmed, and a missing one asks for freeze_reason](layers/config.md#cnfgo-b28--the-freeze-reason-is-shown-trimmed-and-a-missing-one-asks-for-freeze-reason) `CNFGO-B28`

- [A section title comes from the layer, then the project, then the framework](layers/config.md#cnfgo-b29--a-section-title-comes-from-the-layer-then-the-project-then-the-framework) `CNFGO-B29`

- [Placeholder markers default to the templates' marker word](layers/config.md#cnfgo-b30--placeholder-markers-default-to-the-templates-marker-word) `CNFGO-B30`

- [Rule letters come from the declared rule types, or the canonical set](layers/config.md#cnfgo-b31--rule-letters-come-from-the-declared-rule-types-or-the-canonical-set) `CNFGO-B31`

- [A scenario tag maps to every letter that declares it](layers/config.md#cnfgo-b32--a-scenario-tag-maps-to-every-letter-that-declares-it) `CNFGO-B32`

- [The code length pattern matches exactly the declared lengths, contiguous or not](layers/config.md#cnfgo-b33--the-code-length-pattern-matches-exactly-the-declared-lengths-contiguous-or-not) `CNFGO-B33`

- [A rule type catalogues the sections it declares, ignoring case and surrounding spaces](layers/config.md#cnfgo-b41--a-rule-type-catalogues-the-sections-it-declares-ignoring-case-and-surrounding-spaces) `CNFGO-B41`

- [With no filter every suite is selected](layers/config.md#cnfgo-b34--with-no-filter-every-suite-is-selected) `CNFGO-B34`

- [The filter axes intersect](layers/config.md#cnfgo-b35--the-filter-axes-intersect) `CNFGO-B35`

- [Selected suites keep the order of the file](layers/config.md#cnfgo-b36--selected-suites-keep-the-order-of-the-file) `CNFGO-B36`

- [Suite names match ignoring case and surrounding spaces](layers/config.md#cnfgo-b37--suite-names-match-ignoring-case-and-surrounding-spaces) `CNFGO-B37`

- [A name missing from the file is reported with its axis, and an empty combination is not a missing name](layers/config.md#cnfgo-b38--a-name-missing-from-the-file-is-reported-with-its-axis-and-an-empty-combination-is-not-a-missing-name) `CNFGO-B38`

- [The declared vocabulary lists each name once, in file order](layers/config.md#cnfgo-b39--the-declared-vocabulary-lists-each-name-once-in-file-order) `CNFGO-B39`

- [Patterns replace the code template and keep the other derived files](layers/config.md#cnfgo-b40--patterns-replace-the-code-template-and-keep-the-other-derived-files) `CNFGO-B40`

- [What Save writes, Load reads back](layers/config.md#cnfgo-i01--what-save-writes-load-reads-back) `CNFGO-I01`

- [A configuration with every key known keeps loading](layers/config.md#cnfgo-i02--a-configuration-with-every-key-known-keeps-loading) `CNFGO-I02`

- [Every reader answers its default on a nil configuration](layers/config.md#cnfgo-i03--every-reader-answers-its-default-on-a-nil-configuration) `CNFGO-I03`

- [GitHub mode never infers the repository from the git remote](layers/config.md#cnfgo-x01--github-mode-never-infers-the-repository-from-the-git-remote) `CNFGO-X01`

- [Suite names are the project's vocabulary, never a fixed list](layers/config.md#cnfgo-x02--suite-names-are-the-projects-vocabulary-never-a-fixed-list) `CNFGO-X02`

- [A file that cannot be read fails the load with the read error](layers/config.md#cnfgo-e01--a-file-that-cannot-be-read-fails-the-load-with-the-read-error) `CNFGO-E01`

- [Save to a path that cannot be written returns the write error](layers/config.md#cnfgo-e02--save-to-a-path-that-cannot-be-written-returns-the-write-error) `CNFGO-E02`

- [A layer's support globs must be valid](layers/config.md#cnfgo-b46--a-layers-support-globs-must-be-valid) `CNFGO-B46`

- [A gate's timeout ceiling is a share with a default](layers/config.md#cnfgo-b47--a-gates-timeout-ceiling-is-a-share-with-a-default) `CNFGO-B47`

- [A gate's no_signal declares targets with their reason](layers/config.md#cnfgo-b48--a-gates-no-signal-declares-targets-with-their-reason) `CNFGO-B48`

- [A suite's paths say which files it runs](layers/config.md#cnfgo-b49--a-suites-paths-say-which-files-it-runs) `CNFGO-B49`

- [The changelog block has defaults and refuses an unknown mode](layers/config.md#cnfgo-b50--the-changelog-block-has-defaults-and-refuses-an-unknown-mode) `CNFGO-B50`

- [A gate's letters are single letters](layers/config.md#cnfgo-b51--a-gates-letters-are-single-letters) `CNFGO-B51`

- [A gate's invocations compile and name the unit through a group](layers/config.md#cnfgo-b52--a-gates-invocations-compile-and-name-the-unit-through-a-group) `CNFGO-B52`

- [A gate's floor is a percentage](layers/config.md#cnfgo-b53--a-gates-floor-is-a-percentage) `CNFGO-B53`

- [The declared format is read from a CRLF file](layers/config.md#cnfgo-b54--the-declared-format-is-read-from-a-crlf-file) `CNFGO-B54`

- [Each verdict level of a gate takes a state, in order](layers/config.md#cnfgo-b55--each-verdict-level-of-a-gate-takes-a-state-in-order) `CNFGO-B55`

- [A project-wide severity is the default of its blocking gates](layers/config.md#cnfgo-b56--a-project-wide-severity-is-the-default-of-its-blocking-gates) `CNFGO-B56`

- [Coverage floors are percentages, and each one says why](layers/config.md#cnfgo-b57--coverage-floors-are-percentages-and-each-one-says-why) `CNFGO-B57`

- [A field is declared when its path holds a value](layers/config.md#cnfgo-b58--a-field-is-declared-when-its-path-holds-a-value) `CNFGO-B58`

- [A gate's review asks its own question, else the gate's](layers/config.md#cnfgo-b59--a-gates-review-asks-its-own-question-else-the-gates) `CNFGO-B59`

- [A gate relates to the declared layers, and the catalog names what is missing](layers/config.md#cnfgo-b60--a-gate-relates-to-the-declared-layers-and-the-catalog-names-what-is-missing) `CNFGO-B60`

- [A gate's premises are missing or waived, and a waived one is never suggested](layers/config.md#cnfgo-b61--a-gates-premises-are-missing-or-waived-and-a-waived-one-is-never-suggested) `CNFGO-B61`

- [data_states.required makes the data states requirements, and is off by default](layers/config.md#cnfgo-b62--data-statesrequired-makes-the-data-states-requirements-and-is-off-by-default) `CNFGO-B62`

- [A layer may be marked fallible](layers/config.md#cnfgo-b63--a-layer-may-be-marked-fallible) `CNFGO-B63`

- [navigation.entry declares the routes the app opens on](layers/config.md#cnfgo-b64--navigationentry-declares-the-routes-the-app-opens-on) `CNFGO-B64`

- [A gate confronts the repeats of what it declares unless switched off](layers/config.md#cnfgo-b65--a-gate-confronts-the-repeats-of-what-it-declares-unless-switched-off) `CNFGO-B65`

- [The declared containers come back as written, and a missing config has none](layers/config.md#cntnr-b01--the-declared-containers-come-back-as-written-and-a-missing-config-has-none) `CNTNR-B01`

- [The internal containers are the declared ones without the external, in declared order](layers/config.md#cntnr-b02--the-internal-containers-are-the-declared-ones-without-the-external-in-declared-order) `CNTNR-B02`

- [A layer is found in its container ignoring case and surrounding spaces](layers/config.md#cntnr-b03--a-layer-is-found-in-its-container-ignoring-case-and-surrounding-spaces) `CNTNR-B03`

- [A layer no container claims has no container, and that is an answer, not an error](layers/config.md#cntnr-b04--a-layer-no-container-claims-has-no-container-and-that-is-an-answer-not-an-error) `CNTNR-B04`

- [The layers of an external container are still claimed by it](layers/config.md#cntnr-b05--the-layers-of-an-external-container-are-still-claimed-by-it) `CNTNR-B05`

- [The orphan layers are the given ones no container claims, in the given order](layers/config.md#cntnr-b06--the-orphan-layers-are-the-given-ones-no-container-claims-in-the-given-order) `CNTNR-B06`

- [A layer is an orphan exactly when it has no container](layers/config.md#cntnr-i01--a-layer-is-an-orphan-exactly-when-it-has-no-container) `CNTNR-I01`

- [With no container declared, no layer is placed by guessing: every layer is an orphan](layers/config.md#cntnr-x01--with-no-container-declared-no-layer-is-placed-by-guessing-every-layer-is-an-orphan) `CNTNR-X01`

- [A layer that runs no code is never an orphan](layers/config.md#cntnr-b07--a-layer-that-runs-no-code-is-never-an-orphan) `CNTNR-B07`

- [A project that declares no dialect gets only the naming defaults](layers/config.md#dlcti-b01--a-project-that-declares-no-dialect-gets-only-the-naming-defaults) `DLCTI-B01`

- [The family fills every field the project left empty, and a declared field wins](layers/config.md#dlcti-b02--the-family-fills-every-field-the-project-left-empty-and-a-declared-field-wins) `DLCTI-B02`

- [The family name is matched ignoring case](layers/config.md#dlcti-b03--the-family-name-is-matched-ignoring-case) `DLCTI-B03`

- [An unknown family contributes nothing, and the known ones are listed in order](layers/config.md#dlcti-b04--an-unknown-family-contributes-nothing-and-the-known-ones-are-listed-in-order) `DLCTI-B04`

- [The naming conventions apply to any family and a declared one replaces them](layers/config.md#dlcti-b05--the-naming-conventions-apply-to-any-family-and-a-declared-one-replaces-them) `DLCTI-B05`

- [The Gherkin language defaults to English, is found in any case, and one outside the table keeps its code with English keywords](layers/config.md#dlcti-b06--the-gherkin-language-defaults-to-english-is-found-in-any-case-and-one-outside-the-table-keeps-its-code-with-english-keywords) `DLCTI-B06`

- [Every way to open a scenario, in every language, deduplicated and longest first](layers/config.md#dlcti-b07--every-way-to-open-a-scenario-in-every-language-deduplicated-and-longest-first) `DLCTI-B07`

- [Every result keyword, in every language, sorted](layers/config.md#dlcti-b08--every-result-keyword-in-every-language-sorted) `DLCTI-B08`

- [An empty or invalid pattern compiles to nothing](layers/config.md#dlcti-b09--an-empty-or-invalid-pattern-compiles-to-nothing) `DLCTI-B09`

- [The opt-out is read by the YAML field name, ignoring case and spaces](layers/config.md#dlcti-b10--the-opt-out-is-read-by-the-yaml-field-name-ignoring-case-and-spaces) `DLCTI-B10`

- [The set-promise verb is recognised after a provider prefix and never inside a word](layers/config.md#dlcti-b11--the-set-promise-verb-is-recognised-after-a-provider-prefix-and-never-inside-a-word) `DLCTI-B11`

- [The set-slice convention is a query verb opening the name, followed by a slice word](layers/config.md#dlcti-b12--the-set-slice-convention-is-a-query-verb-opening-the-name-followed-by-a-slice-word) `DLCTI-B12`

- [Every pattern a family or a naming default brings compiles](layers/config.md#dlcti-i01--every-pattern-a-family-or-a-naming-default-brings-compiles) `DLCTI-I01`

- [The keywords written for any language are among those every reader recognises](layers/config.md#dlcti-i02--the-keywords-written-for-any-language-are-among-those-every-reader-recognises) `DLCTI-I02`

- [No family brings a collection query, which is the project's to declare](layers/config.md#dlcti-x01--no-family-brings-a-collection-query-which-is-the-projects-to-declare) `DLCTI-X01`

- [The Go family recognises both shapes of error handling](layers/config.md#dlcti-b13--the-go-family-recognises-both-shapes-of-error-handling) `DLCTI-B13`

- [The Go and TS families say how a test is written](layers/config.md#dlcti-b14--the-go-and-ts-families-say-how-a-test-is-written) `DLCTI-B14`

- [The TS family recognises catch with or without its binding](layers/config.md#dlcti-b15--the-ts-family-recognises-catch-with-or-without-its-binding) `DLCTI-B15`

- [The families say what an assertion is, and a project may declare only its own](layers/config.md#dlcti-b16--the-families-say-what-an-assertion-is-and-a-project-may-declare-only-its-own) `DLCTI-B16`

- [The families say how code defines a name](layers/config.md#dlcti-b17--the-families-say-how-code-defines-a-name) `DLCTI-B17`

- [Every examples keyword, in every language, sorted](layers/config.md#dlcti-b18--every-examples-keyword-in-every-language-sorted) `DLCTI-B18`

- [The Go family sees an error in a field and a sentinel error](layers/config.md#dlcti-b19--the-go-family-sees-an-error-in-a-field-and-a-sentinel-error) `DLCTI-B19`

- [Each family reads environment variables its own way](layers/config.md#dlcti-b20--each-family-reads-environment-variables-its-own-way) `DLCTI-B20`

- [The fallible patterns are the project's, and the family's only when the project declares none](layers/config.md#dlcti-b21--the-fallible-patterns-are-the-projects-and-the-familys-only-when-the-project-declares-none) `DLCTI-B21`

- [A family knows its fallible patterns and does not impose them on a project that declared none](layers/config.md#dlcti-b22--a-family-knows-its-fallible-patterns-and-does-not-impose-them-on-a-project-that-declared-none) `DLCTI-B22`

- [import_resolve and navigation_call are read from the dialect](layers/config.md#dlcti-b23--import-resolve-and-navigation-call-are-read-from-the-dialect) `DLCTI-B23`

- [A trigger naming a layer charges every change in that layer, whatever the unit](layers/config.md#dcrqa-b01--a-trigger-naming-a-layer-charges-every-change-in-that-layer-whatever-the-unit) `DCRQA-B01`

- [A trigger naming a unit code charges that unit and not its neighbours in the same layer](layers/config.md#dcrqa-b02--a-trigger-naming-a-unit-code-charges-that-unit-and-not-its-neighbours-in-the-same-layer) `DCRQA-B02`

- [Asked without a unit code, the answer is the layer's alone](layers/config.md#dcrqa-b03--asked-without-a-unit-code-the-answer-is-the-layers-alone) `DCRQA-B03`

- [A documentation with no trigger is never owed by a unit change](layers/config.md#dcrqa-b04--a-documentation-with-no-trigger-is-never-owed-by-a-unit-change) `DCRQA-B04`

- [A trigger matches ignoring case and the spaces around it](layers/config.md#dcrqa-b05--a-trigger-matches-ignoring-case-and-the-spaces-around-it) `DCRQA-B05`

- [Every declared documentation is listed, and a project with no docs block owes none](layers/config.md#dcrqa-b06--every-declared-documentation-is-listed-and-a-project-with-no-docs-block-owes-none) `DCRQA-B06`

- [Naming the unit never removes a documentation the layer alone owes](layers/config.md#dcrqa-i01--naming-the-unit-never-removes-a-documentation-the-layer-alone-owes) `DCRQA-I01`

- [A trigger is matched whole, never as part of a longer name](layers/config.md#dcrqa-x01--a-trigger-is-matched-whole-never-as-part-of-a-longer-name) `DCRQA-X01`

- [Versions are ordered part by part as numbers](layers/config.md#mnvrm-b01--versions-are-ordered-part-by-part-as-numbers) `MNVRM-B01`

- [The tag prefix v is accepted on either side](layers/config.md#mnvrm-b02--the-tag-prefix-v-is-accepted-on-either-side) `MNVRM-B02`

- [With no minimum declared, any running binary satisfies it](layers/config.md#mnvrm-b03--with-no-minimum-declared-any-running-binary-satisfies-it) `MNVRM-B03`

- [A running version that cannot be ordered does not satisfy, and says why](layers/config.md#mnvrm-b04--a-running-version-that-cannot-be-ordered-does-not-satisfy-and-says-why) `MNVRM-B04`

- [A declared minimum that is not MAJOR.MINOR.PATCH is refused, dev included](layers/config.md#mnvrm-b05--a-declared-minimum-that-is-not-majorminorpatch-is-refused-dev-included) `MNVRM-B05`

- [An absent or well-formed minimum is accepted](layers/config.md#mnvrm-b06--an-absent-or-well-formed-minimum-is-accepted) `MNVRM-B06`

- [The refusal names the expected format, an example and what a bad value silences](layers/config.md#mnvrm-b07--the-refusal-names-the-expected-format-an-example-and-what-a-bad-value-silences) `MNVRM-B07`

- [Swapping the two versions always flips the order](layers/config.md#mnvrm-i01--swapping-the-two-versions-always-flips-the-order) `MNVRM-I01`

- [A pre-release is never compared as if it were its final release](layers/config.md#mnvrm-x01--a-pre-release-is-never-compared-as-if-it-were-its-final-release) `MNVRM-X01`

- [A single pattern written as text becomes a list of one](layers/config.md#drptd-b01--a-single-pattern-written-as-text-becomes-a-list-of-one) `DRPTD-B01`

- [A list of patterns is kept whole and in order](layers/config.md#drptd-b02--a-list-of-patterns-is-kept-whole-and-in-order) `DRPTD-B02`

- [An empty list is refused, naming the cause](layers/config.md#drptd-b03--an-empty-list-is-refused-naming-the-cause) `DRPTD-B03`

- [Any shape other than text or a list of text is refused](layers/config.md#drptd-b04--any-shape-other-than-text-or-a-list-of-text-is-refused) `DRPTD-B04`

- [Written back, one pattern is text and several are a list](layers/config.md#drptd-b05--written-back-one-pattern-is-text-and-several-are-a-list) `DRPTD-B05`

- [What is written back reads back as the same patterns](layers/config.md#drptd-i01--what-is-written-back-reads-back-as-the-same-patterns) `DRPTD-I01`

- [The patterns are kept as written, neither expanded nor checked as globs](layers/config.md#drptd-x01--the-patterns-are-kept-as-written-neither-expanded-nor-checked-as-globs) `DRPTD-X01`

- [The project root is the nearest directory above the start that holds the config](layers/config.md#prrpr-b01--the-project-root-is-the-nearest-directory-above-the-start-that-holds-the-config) `PRRPR-B01`

- [With no project above the start, the start comes back unchanged](layers/config.md#prrpr-b02--with-no-project-above-the-start-the-start-comes-back-unchanged) `PRRPR-B02`

- [With no root given, the root is found by walking up from the working directory](layers/config.md#prrpr-b03--with-no-root-given-the-root-is-found-by-walking-up-from-the-working-directory) `PRRPR-B03`

- [The resolved root is always absolute, even outside any project](layers/config.md#prrpr-i01--the-resolved-root-is-always-absolute-even-outside-any-project) `PRRPR-I01`

- [An explicit root is made absolute and never walked above](layers/config.md#prrpr-x01--an-explicit-root-is-made-absolute-and-never-walked-above) `PRRPR-X01`

- [With no source registered, the default gate names are absent, not a failure](layers/config.md#gtvcg-b01--with-no-source-registered-the-default-gate-names-are-absent-not-a-failure) `GTVCG-B01`

- [With a source registered, the default gate names are the ones it gives](layers/config.md#gtvcg-b02--with-a-source-registered-the-default-gate-names-are-the-ones-it-gives) `GTVCG-B02`

- [The letters of plans, flows and actions](layers/config.md#gtvcg-b03--the-letters-of-plans-flows-and-actions) `GTVCG-B03`

- [The answer always comes from the source registered last](layers/config.md#gtvcg-i01--the-answer-always-comes-from-the-source-registered-last) `GTVCG-I01`

- [The names are not cached: each question asks the registered source again](layers/config.md#gtvcg-x01--the-names-are-not-cached-each-question-asks-the-registered-source-again) `GTVCG-X01`

## doct

- [Every declared variable is listed by name, in any language](layers/doct.md#dcenv-b01--every-declared-variable-is-listed-by-name-in-any-language) `DCENV-B01`

- [A variable two units read is listed once, with both](layers/doct.md#dcenv-b02--a-variable-two-units-read-is-listed-once-with-both) `DCENV-B02`

- [The environment page is seeded and compiled](layers/doct.md#dcenv-b03--the-environment-page-is-seeded-and-compiled) `DCENV-B03`

- [A placeholder row is no variable](layers/doct.md#dcenv-e01--a-placeholder-row-is-no-variable) `DCENV-E01`

- [navigation lists the screen specs and the edges the code's flags declare, each end taken to its screen](layers/doct.md#dcnav-b01--navigation-lists-the-screen-specs-and-the-edges-the-codes-flags-declare-each-end-taken-to-its-screen) `DCNAV-B01`

- [An entry is marked, and a screen no entry reaches is marked unreached; with no entry, none is](layers/doct.md#dcnav-b02--an-entry-is-marked-and-a-screen-no-entry-reaches-is-marked-unreached-with-no-entry-none-is) `DCNAV-B02`

- [ScaffoldNavigation and ScaffoldDependencies are seeded when the app has screens and dependency flags, and compile into the pages](layers/doct.md#dcnav-b03--scaffoldnavigation-and-scaffolddependencies-are-seeded-when-the-app-has-screens-and-dependency-flags-and-compile-into-the-pages) `DCNAV-B03`

- [Each Endpoint row of a spec is an operation, with its parameters, body, responses, errors, security and limits](layers/doct.md#opnap-b01--each-endpoint-row-of-a-spec-is-an-operation-with-its-parameters-body-responses-errors-security-and-limits) `OPNAP-B01`

- [A cited contract is a shared schema built from its Domain, in any language](layers/doct.md#opnap-b02--a-cited-contract-is-a-shared-schema-built-from-its-domain-in-any-language) `OPNAP-B02`

- [The document is written in OpenAPI's order, two spaces deep](layers/doct.md#opnap-b03--the-document-is-written-in-openapis-order-two-spaces-deep) `OPNAP-B03`

- [The compiled document carries the generated marker as a YAML comment and is valid YAML](layers/doct.md#opnap-b04--the-compiled-document-carries-the-generated-marker-as-a-yaml-comment-and-is-valid-yaml) `OPNAP-B04`

- [A compiled YAML is recognised as generated and its freshness checked](layers/doct.md#opnap-b05--a-compiled-yaml-is-recognised-as-generated-and-its-freshness-checked) `OPNAP-B05`

- [A contract cited and not found fails the build, naming it](layers/doct.md#opnap-e01--a-contract-cited-and-not-found-fails-the-build-naming-it) `OPNAP-E01`

- [An Endpoint row with no path fails the build](layers/doct.md#opnap-e02--an-endpoint-row-with-no-path-fails-the-build) `OPNAP-E02`

## gate

- [What is no API unit leaves every gate without a verdict](camadas/gate.md#apisp--apispec--the-coherence-of-an-api-spec-and-its-error-codes-in-the-code) `APISP-B01`

- [Every contract cited resolves to a spec with a Domain](camadas/gate.md#apisp--apispec--the-coherence-of-an-api-spec-and-its-error-codes-in-the-code) `APISP-B02`

- [Every error response is complete and under a declared status](camadas/gate.md#apisp--apispec--the-coherence-of-an-api-spec-and-its-error-codes-in-the-code) `APISP-B03`

- [Every declared error code is emitted by the unit's code](camadas/gate.md#apisp--apispec--the-coherence-of-an-api-spec-and-its-error-codes-in-the-code) `APISP-B04`

- [Sections and columns are read in any language](camadas/gate.md#apisp--apispec--the-coherence-of-an-api-spec-and-its-error-codes-in-the-code) `APISP-B05`

- [A code file with no spec beside it leaves without a verdict](camadas/gate.md#apisp--apispec--the-coherence-of-an-api-spec-and-its-error-codes-in-the-code) `APISP-E01`

- [A code file the spec specifies that cannot be read emits nothing](camadas/gate.md#apisp--apispec--the-coherence-of-an-api-spec-and-its-error-codes-in-the-code) `APISP-E02`

- [A section renamed by the project is found](camadas/gate.md#apisp--apispec--the-coherence-of-an-api-spec-and-its-error-codes-in-the-code) `APISP-B06`

- [Below the floor fails naming the lines](camadas/gate.md#brcov--branchcoverage--the-tests-take-the-branches-the-code-has) `BRCOV-B01`

- [A waived branch is left out](camadas/gate.md#brcov--branchcoverage--the-tests-take-the-branches-the-code-has) `BRCOV-B02`

- [A branch no test reaches is likely dead](camadas/gate.md#brcov--branchcoverage--the-tests-take-the-branches-the-code-has) `BRCOV-B03`

- [Nothing to measure is skipped or pending](camadas/gate.md#brcov--branchcoverage--the-tests-take-the-branches-the-code-has) `BRCOV-B04`

- [Branch coverage reads a file with nothing to cover as line coverage does](camadas/gate.md#brcov--branchcoverage--the-tests-take-the-branches-the-code-has) `BRCOV-B05`

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

- [With no map the gate does not approve](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-E01`

- [A code file missing from disk does not hide the orphans of the other target](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-E02`

- [A declared export pattern with no capture group is named, not called undeclared](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-E03`

- [Every file the spec governs is confronted, not only the first](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code) `CDCTC-B11`

- [An identifier in the wrong language is accused, and an English one passes](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-B01`

- [The verdict returns the word that accused](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-B02`

- [Only a DECLARATION is the subject](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-B03`

- [The declarations are found in every form the language offers](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-B04`

- [Deciding one word is separate from deciding a whole identifier](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-B05`

- [The project's own production code declares no identifier in the wrong language](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-B06`

- [No gate decides by matching prose in the team's language](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages) `CDLNG-B07`

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

- [A spec missing from disk is left out of the declared codes](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units) `CRVCD-E01`

- [Changed and removed fields change, added ones do not](camadas/gate.md#ctrim--contractimpact--a-changed-field-names-the-rules-that-use-it-and-their-tests) `CTRIM-B01`

- [A changed field names its rules, here and in the dependents, and their tests](camadas/gate.md#ctrim--contractimpact--a-changed-field-names-the-rules-that-use-it-and-their-tests) `CTRIM-B02`

- [The gate is divergence with the impact, and passes without](camadas/gate.md#ctrim--contractimpact--a-changed-field-names-the-rules-that-use-it-and-their-tests) `CTRIM-B03`

- [The impacted tests are listed for the selection](camadas/gate.md#ctrim--contractimpact--a-changed-field-names-the-rules-that-use-it-and-their-tests) `CTRIM-B04`

- [The rules the change's own revision revises or checks answer the impact](camadas/gate.md#ctrim--contractimpact--a-changed-field-names-the-rules-that-use-it-and-their-tests) `CTRIM-B05`

- [A rule is answered by a revision added to its own spec, even one not yet in git, or by the changed spec's naming its full code; a short code answers only its own spec's rule](camadas/gate.md#ctrim--contractimpact--a-changed-field-names-the-rules-that-use-it-and-their-tests) `CTRIM-B06`

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

- [Without a built map the confrontation is pending](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-E01`

- [Linked code gone from disk is pending, not a code without status](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-E02`

- [A dialect status pattern that does not compile is pending](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-E03`

- [The API catalog's Responses section is the output contract](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those) `CSDCN-B13`

- [What is no API unit leaves without a verdict](camadas/gate.md#cttst--contracttested--an-api-unit-is-proven-against-the-projects-openapi-document) `CTTST-B01`

- [The feature must carry the contract scenario](camadas/gate.md#cttst--contracttested--an-api-unit-is-proven-against-the-projects-openapi-document) `CTTST-B02`

- [A test of the unit must name the contract code](camadas/gate.md#cttst--contracttested--an-api-unit-is-proven-against-the-projects-openapi-document) `CTTST-B03`

- [The contract test must load the OpenAPI document](camadas/gate.md#cttst--contracttested--an-api-unit-is-proven-against-the-projects-openapi-document) `CTTST-B04`

- [The contract regime tag comes from the project](camadas/gate.md#cttst--contracttested--an-api-unit-is-proven-against-the-projects-openapi-document) `CTTST-B05`

- [Scenario and a contract test that loads the document pass](camadas/gate.md#cttst--contracttested--an-api-unit-is-proven-against-the-projects-openapi-document) `CTTST-B06`

- [A code file with no spec beside it leaves without a verdict](camadas/gate.md#cttst--contracttested--an-api-unit-is-proven-against-the-projects-openapi-document) `CTTST-E01`

- [A test that cannot be read is read by its path](camadas/gate.md#cttst--contracttested--an-api-unit-is-proven-against-the-projects-openapi-document) `CTTST-E02`

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

- [Only files are counted, never directories](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-B14`

- [An unreadable file fails the pattern count naming it](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code) `CNHNC-E03`

- [The imports are read by the dialect's pattern and resolved to files of the map](camadas/gate.md#dcgdp--dependencychain--every-import-flagged-with-the-code-it-uses-and-every-symbol-with-who-uses-it) `DCGDP-B01`

- [dep-declared names each import of a governed file with no flag, and the code it would carry](camadas/gate.md#dcgdp--dependencychain--every-import-flagged-with-the-code-it-uses-and-every-symbol-with-who-uses-it) `DCGDP-B02`

- [dep-honored names each flag whose code is not the one its import resolves to](camadas/gate.md#dcgdp--dependencychain--every-import-flagged-with-the-code-it-uses-and-every-symbol-with-who-uses-it) `DCGDP-B03`

- [used-by-declared names each imported symbol whose flag is missing or wrong, and each flag nobody imports](camadas/gate.md#dcgdp--dependencychain--every-import-flagged-with-the-code-it-uses-and-every-symbol-with-who-uses-it) `DCGDP-B04`

- [The fixers write the dependency and used-by flags, correct a wrong code, and remove a stale flag](camadas/gate.md#dcgdp--dependencychain--every-import-flagged-with-the-code-it-uses-and-every-symbol-with-who-uses-it) `DCGDP-B05`

- [A re-export declares the names it lists, and an inline import brings the member it reads](camadas/gate.md#dcgdp--dependencychain--every-import-flagged-with-the-code-it-uses-and-every-symbol-with-who-uses-it) `DCGDP-B06`

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

- [A specified code file missing from disk is left out of the confrontation](camadas/gate.md#dephn--dependencyhonored--methods-promised-in-the-dependency-table-are-consumed-in-code) `DEPHN-E01`

- [A mandatory document that does not exist fails](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B01`

- [A document that exists and does not mention the unit fails](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B02`

- [A mention by the identity code counts as documented](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B03`

- [A mention by the file name also counts](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B04`

- [Satisfying one of two duties is not enough](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B05`

- [Without a declaration nothing is charged](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B06`

- [A layer with no trigger is not charged](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B07`

- [Aggregated, the verdict is one per document](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B08`

- [In a document with sections only a title of its own documents the unit](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B09`

- [The duty starts from the spec, not from the code](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-I01`

- [The layer used is the UNIT's, not the node's](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-I02`

- [Without a map the aggregated verdict is skipped](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-I03`

- [The gate does not understand the document's content](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-X01`

- [The gate does not decide which documents are mandatory](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-X02`

- [A YAML comment is not a Markdown heading](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed) `DCRQD-B10`

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

- [An artifact that is not a spec is skipped](camadas/gate.md#dccvd--docscovered--every-spec-must-reach-some-page-of-the-compiled-documentation) `DCCVD-B01`

- [A project with no templates directory is skipped](camadas/gate.md#dccvd--docscovered--every-spec-must-reach-some-page-of-the-compiled-documentation) `DCCVD-B02`

- [A spec no template reaches fails, naming the spec](camadas/gate.md#dccvd--docscovered--every-spec-must-reach-some-page-of-the-compiled-documentation) `DCCVD-B03`

- [A spec a template reaches passes](camadas/gate.md#dccvd--docscovered--every-spec-must-reach-some-page-of-the-compiled-documentation) `DCCVD-B04`

- [A template that does not compile skips with no message](camadas/gate.md#dccvd--docscovered--every-spec-must-reach-some-page-of-the-compiled-documentation) `DCCVD-B05`

- [The coverage is computed once per project root and map](camadas/gate.md#dccvd--docscovered--every-spec-must-reach-some-page-of-the-compiled-documentation) `DCCVD-B06`

- [Another spec's orphan status never fails the confronted spec](camadas/gate.md#dccvd--docscovered--every-spec-must-reach-some-page-of-the-compiled-documentation) `DCCVD-I01`

- [A reached spec passes without any page having been built](camadas/gate.md#dccvd--docscovered--every-spec-must-reach-some-page-of-the-compiled-documentation) `DCCVD-X01`

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

- [A TBD marker defers only with a reason and outside backticks](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B01`

- [Without a map the gates that read edges are Pending](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B02`

- [plan-doctrine-exists skips an artifact that is not a plan](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B03`

- [A plan whose cited doctrine exists passes](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B04`

- [A plan citing missing doctrines fails, naming each once, in order](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B05`

- [A missing doctrine cited on a TBD line is a divergence, naming it](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B06`

- [A citation with no directory, or of a template, seeds nothing](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B07`

- [doctrine-realized skips an artifact that is not a doctrine](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B08`

- [A doctrine with no rules is skipped](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B09`

- [A doctrine whose every rule is realized passes](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B10`

- [The unrealized rules of a doctrine fail, and only they are named](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B11`

- [Unrealized rules that are all deferred with TBD are Diverge](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B12`

- [An open question is not a rule to realize](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B13`

- [spec-doctrine-exists skips an artifact that is not a spec](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B14`

- [A spec that declares no realization outside TBD lines is skipped](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B15`

- [A declared rule the map resolved, or present in a resolved doctrine, passes](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B16`

- [A declared rule no resolved doctrine holds fails, naming it](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B17`

- [doctrine-not-duplicated skips an artifact that is not a spec](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B18`

- [A spec with no readable realized doctrine is skipped](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B19`

- [A corpus under four rules is Pending, and four rules are measured](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B20`

- [A spec rule that copies the doctrine rule it realizes fails, naming both and the score](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B21`

- [A near copy with a word changed still fails](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B22`

- [A spec rule with text specific to its unit passes](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B23`

- [A copy on a line deferred with TBD is not charged](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B24`

- [spec-realizes-doctrine skips a non-spec and a project with no configuration](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B25`

- [spec-realizes-doctrine skips a spec whose layer does not demand doctrine](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B26`

- [A rule with no realizes tag fails where the layer demands doctrine, naming it](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B27`

- [A rule that declares what it realizes passes](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B28`

- [A rule deferred with TBD is debt, not failure](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B29`

- [A realizes tag after a blank line declares nothing for the rule above](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B30`

- [The demanding layer is the one of the specified target, by edge or by path](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B31`

- [Only a realizes edge realizes a rule](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-I01`

- [A pair that only shares a rare word is not a copy](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-X01`

- [A resolved doctrine that cannot be read confirms none of its rules](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-E01`

- [Doctrine rules and realizes tags are read at the code length the project declares](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B32`

- [A rule is charged once, and an open question is no rule to realize](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied) `DCTRN-B33`

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

- [A key declared twice turns the gate's verdict into a failure naming the lines, unless switched off or skipped](camadas/gate.md#gtdpg--duplicates--each-gate-confronts-the-repeats-of-what-it-declares) `GTDPG-B01`

- [rule-types counts the rule codes a file defines, not those it cites](camadas/gate.md#gtdpg--duplicates--each-gate-confronts-the-repeats-of-what-it-declares) `GTDPG-B02`

- [The spec catalogue's gates count their declarations](camadas/gate.md#gtdpg--duplicates--each-gate-confronts-the-repeats-of-what-it-declares) `GTDPG-B03`

- [A file in no clone passes](camadas/gate.md#duplc--duplication--no-code-file-holds-a-block-copied-from-somewhere-else) `DUPLC-B01`

- [A file in a clone fails naming the other side, from either side](camadas/gate.md#duplc--duplication--no-code-file-holds-a-block-copied-from-somewhere-else) `DUPLC-B02`

- [A clone inside one file names only the line ranges](camadas/gate.md#duplc--duplication--no-code-file-holds-a-block-copied-from-somewhere-else) `DUPLC-B03`

- [Clones within the declared threshold are reported, over it they fail](camadas/gate.md#duplc--duplication--no-code-file-holds-a-block-copied-from-somewhere-else) `DUPLC-B04`

- [Without a declared threshold any clone fails whatever the exit code](camadas/gate.md#duplc--duplication--no-code-file-holds-a-block-copied-from-somewhere-else) `DUPLC-B05`

- [jscpd runs once per scan](camadas/gate.md#duplc--duplication--no-code-file-holds-a-block-copied-from-somewhere-else) `DUPLC-B06`

- [An absolute path in the report is read relative to the root](camadas/gate.md#duplc--duplication--no-code-file-holds-a-block-copied-from-somewhere-else) `DUPLC-B07`

- [No report leaves the check Pending naming why](camadas/gate.md#duplc--duplication--no-code-file-holds-a-block-copied-from-somewhere-else) `DUPLC-E01`

- [The gate runs a pinned jscpd release](camadas/gate.md#duplc--duplication--no-code-file-holds-a-block-copied-from-somewhere-else) `DUPLC-B08`

- [Nothing declared and nothing read leaves without a verdict](camadas/gate.md#envdc--envdeclared--the-environment-variables-a-unit-reads-are-the-ones-its-spec-declares) `ENVDC-B01`

- [A variable read and not declared is named](camadas/gate.md#envdc--envdeclared--the-environment-variables-a-unit-reads-are-the-ones-its-spec-declares) `ENVDC-B02`

- [A variable declared and not read is named, unless deprecated](camadas/gate.md#envdc--envdeclared--the-environment-variables-a-unit-reads-are-the-ones-its-spec-declares) `ENVDC-B03`

- [A read in a comment is no read](camadas/gate.md#envdc--envdeclared--the-environment-variables-a-unit-reads-are-the-ones-its-spec-declares) `ENVDC-B04`

- [Every family's reads are recognised when none is declared](camadas/gate.md#envdc--envdeclared--the-environment-variables-a-unit-reads-are-the-ones-its-spec-declares) `ENVDC-B05`

- [A specified file that cannot be read reads nothing](camadas/gate.md#envdc--envdeclared--the-environment-variables-a-unit-reads-are-the-ones-its-spec-declares) `ENVDC-E01`

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

- [A row the tests do not run fails](camadas/gate.md#exmch--examplesmatch--every-examples-row-is-a-case-its-tests-run) `EXMCH-B01`

- [A value is a whole token with its case](camadas/gate.md#exmch--examplesmatch--every-examples-row-is-a-case-its-tests-run) `EXMCH-B02`

- [A label column is display text](camadas/gate.md#exmch--examplesmatch--every-examples-row-is-a-case-its-tests-run) `EXMCH-B03`

- [The rows are the coded outline's tables](camadas/gate.md#exmch--examplesmatch--every-examples-row-is-a-case-its-tests-run) `EXMCH-B04`

- [The test's body is read as the assertion gate reads it](camadas/gate.md#exmch--examplesmatch--every-examples-row-is-a-case-its-tests-run) `EXMCH-B05`

- [Nothing to confront is skipped](camadas/gate.md#exmch--examplesmatch--every-examples-row-is-a-case-its-tests-run) `EXMCH-B06`

- [Tests that cannot be listed fail the gate with the reason](camadas/gate.md#exmch--examplesmatch--every-examples-row-is-a-case-its-tests-run) `EXMCH-E01`

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

- [Under --index a command reads the commit's content](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments) `EXCMX-B14`

- [Every failure gate skips an artifact that is not a spec](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B01`

- [A failure rule is read in the heading, table row and bullet forms, not in prose](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B02`

- [A spec that declares no failure is skipped by failure-handled and failure-logged](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B03`

- [Without the dialect patterns the failure gates are Pending](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B04`

- [Without governed code that can be read the failure gates are Pending](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B05`

- [A handling written only in a comment line does not count](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B06`

- [A declared failure with no handling in the governed code fails, naming it](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B07`

- [Any handling path in the governed code passes failure-handled](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B08`

- [A failure marked resilient with a reason is not charged](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B09`

- [A bare resilient marker exempts nothing](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B10`

- [A handling that records nothing fails failure-logged, naming the failure](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B11`

- [A handling that records the occurrence passes failure-logged](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B12`

- [Handling in the code with no failure declared in the spec fails](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B13`

- [Code with no handling is skipped by failure-declared](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B14`

- [Handling in the code and a failure declared in the spec passes](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B15`

- [An Errors section closed with none and a reason passes failure-declared](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B16`

- [The conclusions read the resilient and observing reasons of each failure whole](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B17`

- [A conclusion reason ends at its table cell](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B18`

- [failure-handled and failure-logged never both fail the same spec](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-I01`

- [One handling path answers for every declared failure](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-X01`

- [A failure rule is read at the code length the project declares](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B19`

- [A unit with a fallible source and no declared failure fails, naming the source](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B20`

- [A dependency on a file of a fallible layer is a fallible source](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B21`

- [A fallible call with no handling in its window fails, named by its line](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B22`

- [A fallible call whose result is returned hands its failure to the caller](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared) `FLRAI-B23`

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

- [Minor description drift is a divergence, not a failure](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-X03`

- [An unmapped regime tag does not exempt a scenario](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B18`

- [A linked test gone from disk implements nothing while the others still count](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-E01`

- [Comment markers inside strings or glued to a name keep the line](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B19`

- [The drift verdict is written in the project's language](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B20`

- [A parametrised, focused or skipped test is read by its own title](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B21`

- [A failing tests source fails the gate naming the error](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-E02`

- [A support file linked to a feature is not confronted as its test](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B22`

- [A scenario suffix gives each case of a rule its own identity](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B23`

- [Every test that cites the code is compared, not only the first](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B24`

- [A test title that says the scenario's title and more matches it](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B25`

- [A VR code carries its state, and each state's scenario is its own code](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description) `FTMFT-B26`

- [Only a check with a registered fixer is fixable](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-B01`

- [A stale date on a committed file is rewritten to its last commit date](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-B02`

- [A file with an uncommitted edit takes today's date](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-B03`

- [A date that already matches is left alone and reported as nothing](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-B04`

- [Only fixable gates, the nodes they apply to and files on disk are touched](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-B05`

- [Outside a git repository the file is left untouched](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-B06`

- [A file never committed and not edited is left untouched](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-B07`

- [Only the date changes, the rest of the file is kept byte for byte](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-I01`

- [A header without the date field is never given one](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-X01`

- [A write that fails is reported as not fixed, with the cause](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-E01`

- [The repair detail is written in the project's language](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-B08`

- [The fix writes the missing header from the map](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-B09`

- [The fix gives a file with an identity and no code of its own a code unique in the map](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-B10`

- [FixWithConfig gives the fixers the project's config for that run](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-B11`

- [Each repair says whether it moved lines](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-B12`

- [The date repair answers each file from one snapshot of the repository](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix) `FXIXX-B13`

- [The flag gates skip what is not a flag, and a flag with no scenario](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B01`

- [A flag whose every condition is in the grammar passes](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B02`

- [A condition in prose fails the grammar, naming the scenario](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B03`

- [A flag that declares the absent case passes completeness](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B04`

- [A flag without the absent case fails, naming the waiver](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B05`

- [The absent waiver needs a written reason, outside backticks](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B06`

- [Only a citation of a G code is confronted](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B07`

- [A citation of a declared scenario passes](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B08`

- [Citations of scenarios that do not exist fail, each named once and sorted](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B09`

- [Without a map the governance gate is Pending](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B10`

- [A scenario no rule cites fails governance](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B11`

- [A scenario with a reasoned governance waiver is not charged](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B12`

- [Without a map the coverage gate is Pending](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B13`

- [Coverage is judged per scenario](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B14`

- [A scenario no test names fails as having no test](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B15`

- [A written test not yet proven is told apart: not ingested, or ingested and not green](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B16`

- [A gate that could not measure never answers Pass](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-I01`

- [No waiver is accepted without a written reason](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-X01`

- [A gated-by citation is read at the code length the project declares](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B17`

- [Flags that cannot be read leave the citation Pending, naming the flags folder](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-E02`

- [With a tests source a flag scenario is written only when a title cites it](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B18`

- [A failing tests source fails flag-covered naming the error](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-E03`

- [A flag scenario is green only by its own proof](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested) `FLSCF-B19`

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

- [Each verdict level does what the gate's severity says](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B15`

- [Only the obligations gate produces assumed debt](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B16`

- [The engine reconfigures the code grammar before running anything](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B17`

- [The plain entry point runs with the map and no Structure](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B18`

- [The entry point that carries the Structure hands it to the checkers](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B19`

- [The entry point that knows the sweep kind honours the full-sweep scope](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B20`

- [The entry point that honours a waiver by target keeps the gate running](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B21`

- [A vendored file is out of every internal ruler and still reached by an external command](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B22`

- [A gate declared with run executes the custom runner even when canonical specifies check](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B23`

- [A waiver by target never removes the gate from the list](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-I01`

- [Stepping aside, not measuring and failing are three different answers](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-I02`

- [A waived target leaves the failure tally without leaving the report](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-I03`

- [The reported target of an aggregate gate is the scope](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-I04`

- [The engine does not decide whether a target is correct](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-X01`

- [The engine does not compute the verdict of a judgment gate](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-X02`

- [The engine invents neither a map nor a Structure nor a waiver](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-X03`

- [A target declared with nothing to measure is skipped with its reason](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B24`

- [The gate's reach is offered to the runs that choose what to measure](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B25`

- [The gates read files from the source set](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B26`

- [The gates ask whether a file changed through the source set](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B27`

- [A gate that breaks fails its target and the check goes on](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B28`

- [Under an index source no gate reads the tree](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-I05`

- [A gate presupposing an undeclared field asks nothing](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes) `GTENG-B29`

- [A layer the Estrutura does not have fails](camadas/gate.md#hdlyd--headerlayerdeclared--the-layer-a-header-declares-is-one-the-estrutura-has) `HDLYD-B01`

- [A layer that differs only in case fails naming both](camadas/gate.md#hdlyd--headerlayerdeclared--the-layer-a-header-declares-is-one-the-estrutura-has) `HDLYD-B02`

- [Nothing to confront is skipped](camadas/gate.md#hdlyd--headerlayerdeclared--the-layer-a-header-declares-is-one-the-estrutura-has) `HDLYD-B03`

- [Confronting an artifact that is not a spec skips](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B01`

- [Confronting without a graph returns Pending](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B02`

- [Confronting a spec without a declared code skips](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B03`

- [A testID prefix matching the spec identity code passes](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B04`

- [An orphan testID prefix with code shape fails](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B05`

- [A testID prefix matching another declared unit is accepted as legitimate reuse](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B06`

- [A visual regression baseline with a divergent code fails](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B07`

- [A baseline of one of the unit's rules is read by its unit code](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B11`

- [Short testID prefixes of three letters or fewer pass](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B08`

- [Inconsistency failures cite conflicting acronyms and origin files](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B09`

- [When no orphan testID prefixes or baseline discrepancies exist the gate passes](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B10`

- [Without a map graph the relational gate never approves](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-I01`

- [Cross-unit reuse is never permitted for visual regression baselines](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-I02`

- [Absence of a spec code is never double-charged](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-I03`

- [Common shorthand prefixes of three letters or fewer are not scrutinized](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-X01`

- [The gate does not enforce code presence on specs](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-X02`

- [Components referencing parent screen codes in testIDs are not forbidden](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-X03`

- [A governed file gone from disk is left out and the others are still confronted](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-E01`

- [A spec under a bracketed directory still has its baselines confronted](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-E02`

- [A testID prefix with a code the project renamed to this unit's, or to another unit's, passes](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline) `IDCND-B12`

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

- [A test missing from disk does not hide the codes the other tests name](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-E01`

- [With a tests source a scenario is written only when a title cites it](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B17`

- [A failing tests source fails scenario-coverage naming the error](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-E02`

- [A support file is neither run nor counted as naming a scenario](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B18`

- [A feature scenario is recognised in any Gherkin language](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B19`

- [Scenario coverage charges what the spec defines, not what it cites](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B20`

- [Scenario coverage tells a missing test apart from a test never run](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B21`

- [Scenario coverage honours a layer that dispenses tested-by](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B22`

- [Mutation score passes at the threshold and fails below it naming the survivors](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B23`

- [A missing or stale mutation signal is pending, unless it met the floor](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B24`

- [A score between acceptable and desirable is divergence, not failed](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B25`

- [The mutation thresholds come from the report](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B26`

- [The verdict follows the isolated scope and the report reads the delta](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B27`

- [A scope measured at an older revision decides nothing](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B28`

- [Without a repository updated-at skips instead of blaming the file](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B29`

- [A section title in another language fails unless the gate waives it](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B30`

- [The skeletons anchors new emits are born conforming](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B31`

- [A mutation score measured under load is not trusted](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B32`

- [A header below the top fails the header ruler](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B33`

- [Each variant of a rule must be proven](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B34`

- [The header's layer is read as the project declares it](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B35`

- [A file with nothing to cover is not pending](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B36`

- [A file the mutation tool measured with no mutant has nothing to mutate](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B37`

- [Line coverage is held to the gate's floor, or a glob's floor with its reason](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B38`

- [A header in a double-dash comment has its identity](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B39`

- [A header without a code of its own fails, and a code another file owns fails naming it](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B41`

- [The header is the block at the top, and a file of no unit is identified by its layer](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function) `INCHN-B40`

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

- [A marking inside an unreadable directory counts as absent from its end](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it) `MRPRM-E01`

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

- [A key added to a stamped double gets a stamp](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one) `MKSTP-B07`

- [An existing stamp is never rewritten](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one) `MKSTP-I01`

- [A repeated line never anchors a stamp](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one) `MKSTP-I02`

- [Generator does not refresh a divergent stamp](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one) `MKSTP-X01`

- [A double detector that does not compile stops the generator with an error](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one) `MKSTP-E01`

- [An undeclared double detector stops the generator with an error](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one) `MKSTP-E02`

- [A module missing from disk is skipped and the other doubles are still stamped](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one) `MKSTP-E03`

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

- [Changing a module checks the doubles stamped against it](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B17`

- [The author of a change refreshes the stamps and gets the doubles to adjust](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B18`

- [A stamp with a non-positive line count is not a stamp](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-E01`

- [A test missing from disk is left out of the doubles of a changed module](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-E02`

- [The modules a test's stamps point at](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B19`

- [The stamps that held before a mechanical rewrite are refreshed after it, and a stamp already stale stays stale](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B20`

- [The chain's flags change no contract](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B21`

- [A flag in a block comment inside a line changes no contract](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it) `MCSTM-B22`

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

- [An undeclared tie shape is pending instead of guessing one](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-I02`

- [The verdict counts the loose doubles, not just names them](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-I03`

- [The gate does not check whether the annotated type matches the real module](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-X01`

- [The gate does not cover drift of behaviour](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-X02`

- [The gate carries no built-in tie shape and no built-in ecosystem](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-X03`

- [The gate does not charge third-party library doubles](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-X04`

- [A double with no factory is not charged](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-X05`

- [With no map the verdict is pending, not an external-only skip](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-E02`

- [A factory on the next line is still read, tie and all](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces) `MCTYM-B13`

- [nav-annotated names each navigation call with no flag, and each flag naming another screen than its route](camadas/gate.md#ncgnv--navigationchain--every-navigation-flagged-with-the-screen-it-leads-to-and-the-screens-tables-confronted-with-it) `NCGNV-B01`

- [nav-matches-spec names each Out row no flag answers, and each flag the Out table does not declare](camadas/gate.md#ncgnv--navigationchain--every-navigation-flagged-with-the-screen-it-leads-to-and-the-screens-tables-confronted-with-it) `NCGNV-B02`

- [nav-symmetric names each Out with no matching In, and each In with no matching Out](camadas/gate.md#ncgnv--navigationchain--every-navigation-flagged-with-the-screen-it-leads-to-and-the-screens-tables-confronted-with-it) `NCGNV-B03`

- [nav-reachable fails a screen no entry route reaches, and is pending with no entry declared](camadas/gate.md#ncgnv--navigationchain--every-navigation-flagged-with-the-screen-it-leads-to-and-the-screens-tables-confronted-with-it) `NCGNV-B04`

- [The fixer flags a navigation call whose route names one screen, and leaves a back navigation to the author](camadas/gate.md#ncgnv--navigationchain--every-navigation-flagged-with-the-screen-it-leads-to-and-the-screens-tables-confronted-with-it) `NCGNV-B05`

- [On a line ending in a JSX tag the fixer writes the flag where it renders nothing](camadas/gate.md#ncgnv--navigationchain--every-navigation-flagged-with-the-screen-it-leads-to-and-the-screens-tables-confronted-with-it) `NCGNV-B06`

- [A node that carries the trigger and is absent from the demanded file fails](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B01`

- [A node that carries the trigger and does appear passes](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B02`

- [A node without the trigger contracts no obligation](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B03`

- [A waiver exempts only when it carries a written reason](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B04`

- [A project with no declared obligation is skipped](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B05`

- [An acknowledged debt with a written when yields a divergence](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B06`

- [A bare debt marker keeps failing](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B07`

- [Waiver and debt stay distinct](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B08`

- [The failing verdict offers the three ways out](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-B09`

- [The token is derived through the declared form](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-I01`

- [A glob that matches no file produces no violation](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-I02`

- [The node's own identified_as wins over the automatic form](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-I03`

- [The gate does not decide which obligations exist](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-X01`

- [The gate does not understand what the destination does with the token](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-X02`

- [A declaration written in the body is not read](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-X03`

- [An unreadable destination file is not proof, and the others are still searched](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-E01`

- [A must_appear_in glob that does not parse fails naming it](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit) `OBHNB-E02`

- [The duties in force are the resolved list](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project) `BLGTN-B01`

- [Without packs the duties in force are the inline list](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project) `BLGTN-B02`

- [The pack's duties come first, the inline ones last, with the pack's fields carried over](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project) `BLGTN-B03`

- [A pack duty's reason cites the norm it comes from](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project) `BLGTN-B04`

- [A pack that fails to load keeps the inline duties](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project) `BLGTN-B05`

- [The pack duties are read once per root](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project) `BLGTN-B06`

- [Another set of packs for the same root is read, not served from the cache](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project) `BLGTN-B12`

- [The report gives one status per duty, in the order handed](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project) `BLGTN-B07`

- [A node is a subject only when its header carries the duty's trigger](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project) `BLGTN-B08`

- [A waiver with a reason counts as fulfilled and as waived](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project) `BLGTN-B09`

- [A duty declared pending counts as debt](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project) `BLGTN-B10`

- [Every other subject is missing, and the missing list is sorted](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project) `BLGTN-B11`

- [Every subject is counted exactly once](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project) `BLGTN-I01`

- [The report evaluates only the duties it is handed](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project) `BLGTN-X01`

- [A node whose file cannot be read is not a subject](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project) `BLGTN-E02`

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

- [A spec with no section is a divergence item, and the verdict teaches the way out](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement) `OPQSP-X02`

- [A code the question cites as context does not close it](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement) `OPQSP-B08`

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

- [Plans with phase-like sections lacking codes return Diverge](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-B04`

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

- [With no map a phase or parent reference is not judged](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-E01`

- [A plan missing from disk does not unresolve the phases of the other plans](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-E02`

- [The dependency is read in any supported language](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent) `PHORP-B14`

- [A raw skeleton fails confrontation](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-B01`

- [A header field whose value is the marker fails](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-B02`

- [A table cell holding only the marker fails](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-B03`

- [A title or body line opening with the marker fails](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-B04`

- [The verdict names what was left behind](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-B05`

- [An artifact with every marker replaced passes](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-B06`

- [The marker vocabulary is the project's, and TODO by default](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in) `PLCFL-B07`

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

- [A revision written as a section title is read with the others](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B15`

- [Untouched nodes in the impact radius are never charged](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-I01`

- [Revisions must be sequential from one without gaps](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-I02`

- [Files not existing in HEAD are never charged](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-I03`

- [The gate does not judge whether a correction is directional](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-X01`

- [The gate does not enforce identity code presence](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-X02`

- [Without a changed files list the gate skips confrontation](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-X03`

- [A change only in the header and in renamed codes is mechanical and skips](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed) `PCJPL-B16`

- [Non-plan artifacts skip confrontation](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B01`

- [Confronting without a map graph returns pending](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B02`

- [A plan with neither revisions nor revisers skips confrontation](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B03`

- [Declaring a revision target that does not exist in the map fails](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B04`

- [A revising plan receives a divergence reminder when the target lacks a top notice](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B05`

- [The divergence reminder on a revising plan clears once the target carries the notice](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B06`

- [A revised plan lacking a top revision notice fails](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B07`

- [A revised plan placing the revision notice after line 40 fails](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B08`

- [A revised plan with top notice but no section amendment markers returns divergence](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B09`

- [A revised plan with top notice and marked section amendments passes](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B10`

- [Markdown alerts and metadata directives are both accepted as valid markers](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-B11`

- [Revision notices must be placed within the first 40 lines](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-I01`

- [Missing section amendment markers yield divergence rather than failure](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-I02`

- [The divergence reminder clears once the revised plan is notified](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans) `PLRVP-I03`

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

- [Every source line of the plan is read](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-B10`

- [Every letter and digit of the name counts in the match](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-B11`

- [What was not measured is never approved](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-I01`

- [A seeded file off the naming pattern owns nothing](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-I02`

- [The gate does not confront the order of the phases](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-X01`

- [The gate does not demand a needs pointing at nothing](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-X02`

- [The gate does not interpret what the source is for](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it) `PSDPL-X03`

- [The rows are read by position, and a spec with none is skipped](camadas/gate.md#prsnt--presentationgates--the-presentation-validations-confronted) `PRSNT-B01`

- [Every declared value has an appearance](camadas/gate.md#prsnt--presentationgates--the-presentation-validations-confronted) `PRSNT-B02`

- [One prop and one condition lead to one appearance](camadas/gate.md#prsnt--presentationgates--the-presentation-validations-confronted) `PRSNT-B03`

- [The text shown is a message code, not copy](camadas/gate.md#prsnt--presentationgates--the-presentation-validations-confronted) `PRSNT-B04`

- [What changes is something a test can point at](camadas/gate.md#prsnt--presentationgates--the-presentation-validations-confronted) `PRSNT-B05`

- [Each gate gets a summary counting its verdicts](camadas/gate.md#prflo--profile--the-verdicts-of-a-run-gathered-per-gate-and-per-node) `PRFLO-B01`

- [The summary carries the total time and the most expensive run](camadas/gate.md#prflo--profile--the-verdicts-of-a-run-gathered-per-gate-and-per-node) `PRFLO-B02`

- [Every failed result is listed as a failure](camadas/gate.md#prflo--profile--the-verdicts-of-a-run-gathered-per-gate-and-per-node) `PRFLO-B03`

- [Only a failure of a blocking gate blocks promotion](camadas/gate.md#prflo--profile--the-verdicts-of-a-run-gathered-per-gate-and-per-node) `PRFLO-B04`

- [A pending verdict blocks only when its gate blocks and it impedes](camadas/gate.md#prflo--profile--the-verdicts-of-a-run-gathered-per-gate-and-per-node) `PRFLO-B05`

- [Results awaiting judgement are listed as awaiting judgement](camadas/gate.md#prflo--profile--the-verdicts-of-a-run-gathered-per-gate-and-per-node) `PRFLO-B06`

- [Per-node verdicts include only confronted nodes, sorted](camadas/gate.md#prflo--profile--the-verdicts-of-a-run-gathered-per-gate-and-per-node) `PRFLO-B07`

- [A node is failed when a blocking gate failed on it or left it an impeding pending](camadas/gate.md#prflo--profile--the-verdicts-of-a-run-gathered-per-gate-and-per-node) `PRFLO-B08`

- [Promotion is refused exactly when something blocks](camadas/gate.md#prflo--profile--the-verdicts-of-a-run-gathered-per-gate-and-per-node) `PRFLO-I01`

- [The profile does not decide whether a pending impedes](camadas/gate.md#prflo--profile--the-verdicts-of-a-run-gathered-per-gate-and-per-node) `PRFLO-X01`

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

- [The verdict carries only the directions that found something](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-B12`

- [The companion's path has one definition, derived from the scanner](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-I01`

- [A seed is matched by path, never by the item's text](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-I02`

- [A spec mentioned in the plan's prose is not a seed](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-I03`

- [A template file is never a seeded spec](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-I04`

- [The gate does not charge the existence of the progress file](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-X01`

- [The gate does not judge the content of an item beyond the path](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-X02`

- [The gate does not put the progress file into the map](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk) `PRHNP-X03`

- [The project's own source wins over its family's, and neither declares nothing](camadas/gate.md#prjts--projecttests--the-gates-read-the-projects-tests-through-the-source-the-project-declares) `PRJTS-B01`

- [A pattern reads only the map's test files](camadas/gate.md#prjts--projecttests--the-gates-read-the-projects-tests-through-the-source-the-project-declares) `PRJTS-B02`

- [The tests are read once per scan](camadas/gate.md#prjts--projecttests--the-gates-read-the-projects-tests-through-the-source-the-project-declares) `PRJTS-B03`

- [The tests of a set of files come in file order then line](camadas/gate.md#prjts--projecttests--the-gates-read-the-projects-tests-through-the-source-the-project-declares) `PRJTS-B04`

- [A source error is returned with the declaration](camadas/gate.md#prjts--projecttests--the-gates-read-the-projects-tests-through-the-source-the-project-declares) `PRJTS-B05`

- [Support files are not read as tests](camadas/gate.md#prjts--projecttests--the-gates-read-the-projects-tests-through-the-source-the-project-declares) `PRJTS-B06`

- [Under --index the tests are listed from the commit](camadas/gate.md#prjts--projecttests--the-gates-read-the-projects-tests-through-the-source-the-project-declares) `PRJTS-B07`

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

- [Without a map the demand is left pending, not approved](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-E01`

- [Governed code gone from disk leaves the demand pending](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side) `PCBPR-E02`

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

- [A reference to a code that exists nowhere fails](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B14`

- [Without sibling spec, an existing code still skips](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B15`

- [An inferred identity does not satisfy the reference](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B16`

- [Without a graph, absence is not asserted](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit) `RFRSR-B17`

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

- [feature-spec-match skips a node that is not a feature, and a feature with no coded scenario](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B01`

- [Without a map both gates are Pending](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B02`

- [A feature no spec covers is Pending](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B03`

- [A covering spec that defines no requirement leaves feature-spec-match Pending](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B04`

- [A scenario whose rule the spec no longer declares fails, named as written](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B05`

- [A numbered variant is the same rule](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B06`

- [A data state is defined by its name](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B07`

- [The visual baseline is never charged](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B08`

- [The rules of every covering spec count together](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B09`

- [test-feature-match skips what is not a test, and a test no feature exercises](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B10`

- [A linked feature with no coded scenario leaves test-feature-match Pending](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B11`

- [A code the test names that no scenario declares fails](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B12`

- [A test naming the bare rule where the feature declares its variants](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B13`

- [A revision code is not charged as a rule](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B14`

- [A linked spec or feature that cannot be read contributes nothing](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-E01`

- [Neither gate answers Pass when it had nothing to match against](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-I01`

- [Codes of another unit are never charged](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-X01`

- [A code named only in a comment is not a claim of proof](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-X02`

- [A data state defined with the unit prefix is read at the code length the project declares](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B15`

- [A support file is not confronted as a test](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B16`

- [A test of a declarative unit is not charged a feature](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B17`

- [A test naming a variant its feature does not declare](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario) `RVMTR-B18`

- [The targets of a reviewed gate with no review at their revision are to review](camadas/gate.md#rvdur--reviewsdue--the-targets-of-a-reviewed-gate-that-no-review-covers-at-their-current-revision) `RVDUR-B01`

- [A gate with no review has nothing to review, and a named gate lists its own](camadas/gate.md#rvdur--reviewsdue--the-targets-of-a-reviewed-gate-that-no-review-covers-at-their-current-revision) `RVDUR-B02`

- [A change to a reviewed target makes it to review again](camadas/gate.md#rvdur--reviewsdue--the-targets-of-a-reviewed-gate-that-no-review-covers-at-their-current-revision) `RVDUR-B03`

- [Listing what is to review changes no verdict](camadas/gate.md#rvdur--reviewsdue--the-targets-of-a-reviewed-gate-that-no-review-covers-at-their-current-revision) `RVDUR-I01`

- [The list never blocks](camadas/gate.md#rvdur--reviewsdue--the-targets-of-a-reviewed-gate-that-no-review-covers-at-their-current-revision) `RVDUR-X01`

- [Non-spec artifacts skip confrontation](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B01`

- [A spec with no revision has nothing to confront](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B02`

- [A revision naming no rule cannot be confronted](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B03`

- [A revision naming a rule the spec does not define fails](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B04`

- [A sibling sharing vocabulary and left unmentioned is reported](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B05`

- [Every vocabulary-sharing sibling accounted for passes](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B06`

- [Checked clears the accusation without asserting correctness](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B07`

- [A revised rule is never its own orphan](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-I01`

- [Terms shared by the whole unit do not discriminate](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-X01`

- [A rule's title is its heading, not its usage row](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B08`

- [A revision that revised no rule says so](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so) `RVORP-B09`

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

- [An unreadable registry file leaves the route pending](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry) `RTEXR-E03`

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

- [A unit that predates the practice is a divergence item, not a failure](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-B03`

- [Declaring the requirement turns the divergence item into a failure](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-B04`

- [A spec with no linked code is not this gate's subject](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-B05`

- [A waiver with no named rule covers every rule of the spec](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-B06`

- [Requiring the marking never punishes whoever already marks](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-I01`

- [The identity survives a rename](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-I02`

- [The gate does not judge whether the implementation is correct](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-X01`

- [The gate does not demand a mark on EVERY rule](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-X02`

- [A code file that cannot be read is pending, naming the file](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-E01`

- [With data_states.required, the code cites each data state the spec defines](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them) `RLIMR-B07`

- [A letter that is not declared in the vocabulary fails](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B01`

- [A declared letter under a claimed section passes](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B02`

- [A section cataloguing rules under a title no letter claims fails](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B03`

- [The same letter claimed by two terms is a conflict in the vocabulary](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B04`

- [With no vocabulary declared the gate confronts the canonical letters](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B05`

- [A heading that is the rule code itself is not a category section](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B06`

- [A section that only cites other sections' codes claims no letter](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B07`

- [A section that defines a code in the first table cell is charged](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B08`

- [A section declared as rule-cataloguing and filled without a code is a divergence](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared) `RLTYR-B09`

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

- [The sections are read in any language and by position](camadas/gate.md#rlusg--ruleuses--each-rule-says-what-it-uses-and-what-it-uses-exists) `RLUSG-B01`

- [A uses cell lists backticked or comma-separated items](camadas/gate.md#rlusg--ruleuses--each-rule-says-what-it-uses-and-what-it-uses-exists) `RLUSG-B02`

- [A rule that does not say what it uses fails](camadas/gate.md#rlusg--ruleuses--each-rule-says-what-it-uses-and-what-it-uses-exists) `RLUSG-B03`

- [Nothing to ask about is a skip](camadas/gate.md#rlusg--ruleuses--each-rule-says-what-it-uses-and-what-it-uses-exists) `RLUSG-B04`

- [What a rule uses must exist in the spec](camadas/gate.md#rlusg--ruleuses--each-rule-says-what-it-uses-and-what-it-uses-exists) `RLUSG-B05`

- [A spec whose rules say nothing yet is a skip](camadas/gate.md#rlusg--ruleuses--each-rule-says-what-it-uses-and-what-it-uses-exists) `RLUSG-B06`

- [The fields a rule uses appear in the code the spec governs](camadas/gate.md#rlusg--ruleuses--each-rule-says-what-it-uses-and-what-it-uses-exists) `RLUSG-B07`

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

- [The verdict is a divergence and never a failure](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-B07`

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

- [A scenario whose steps copy another's is reported](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code) `SCIDS-B12`

- [An artifact that is not a feature leaves without a verdict](camadas/gate.md#scltr--scenarioletterdeclared--the-letter-of-a-scenario-code-exists-in-the-vocabulary) `SCLTR-B01`

- [With no declared vocabulary the canonical letters apply](camadas/gate.md#scltr--scenarioletterdeclared--the-letter-of-a-scenario-code-exists-in-the-vocabulary) `SCLTR-B02`

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

- [A scenario tag disagreeing with the code rule letter returns a divergence](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-B10`

- [The Diverge verdict cites details of the type divergence](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-B11`

- [Multiple mismatch findings are sorted deterministically](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-B12`

- [Classification alignment is evaluated only on feature artifacts](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-I01`

- [Absence of tag mappings prevents speculative enforcement](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-I02`

- [Mismatched scenario types return Diverge rather than Fail](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter) `STASC-I03`

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

- [Two test files of one unit in one layer fail](camadas/gate.md#sngtu--singletestperunit--a-unit-has-one-test-file-per-test-layer) `SNGTU-B01`

- [A declared split passes](camadas/gate.md#sngtu--singletestperunit--a-unit-has-one-test-file-per-test-layer) `SNGTU-B02`

- [Not code, no test, no map](camadas/gate.md#sngtu--singletestperunit--a-unit-has-one-test-file-per-test-layer) `SNGTU-B03`

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

- [A rule alias needs no scenario of its own](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B13`

- [An alias that stands for no rule fails](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B14`

- [Every waiver requires a written reason](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-I01`

- [Each gate accuses one thing](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-I02`

- [The gate does not judge whether the scenario proves the requirement](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-X01`

- [A cited code produces no accusation](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-X02`

- [The gate does not confront feature against test](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-X03`

- [Without a built map the confrontation is pending](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-E01`

- [A feature missing from disk does not hide the scenarios of the other features](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-E02`

- [Scenarios tagged with a suffix cover the requirement](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B15`

- [A code in the open-decisions section is no requirement](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B16`

- [With data_states.required, a data state the spec defines is a requirement](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B17`

- [A rule marked @retired is no requirement, no rule of the code, and uses nothing](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario) `SFMSP-B18`

- [What has nothing to confront leaves both gates without a verdict](camadas/gate.md#vtrst--statetransitions--every-change-of-a-visual-unit-is-proven-through-what-the-screen-shows) `VTRST-B01`

- [Every validation is the trigger of a transition](camadas/gate.md#vtrst--statetransitions--every-change-of-a-visual-unit-is-proven-through-what-the-screen-shows) `VTRST-B02`

- [A transition from or to an unknown state does not count](camadas/gate.md#vtrst--statetransitions--every-change-of-a-visual-unit-is-proven-through-what-the-screen-shows) `VTRST-B03`

- [A validation exempted with a reason is not asked](camadas/gate.md#vtrst--statetransitions--every-change-of-a-visual-unit-is-proven-through-what-the-screen-shows) `VTRST-B04`

- [Every error names the message it shows](camadas/gate.md#vtrst--statetransitions--every-change-of-a-visual-unit-is-proven-through-what-the-screen-shows) `VTRST-B05`

- [A code file with no spec beside it leaves without a verdict](camadas/gate.md#vtrst--statetransitions--every-change-of-a-visual-unit-is-proven-through-what-the-screen-shows) `VTRST-E01`

- [Sections nested under another are found](camadas/gate.md#vtrst--statetransitions--every-change-of-a-visual-unit-is-proven-through-what-the-screen-shows) `VTRST-B06`

- [A test with no assertion fails, named by its line and title](camadas/gate.md#thsas--testhasassertion--every-test-asserts-something-in-its-body) `THSAS-B01`

- [The body is the block the test's line opens](camadas/gate.md#thsas--testhasassertion--every-test-asserts-something-in-its-body) `THSAS-B02`

- [A multi-line literal is text, not layout](camadas/gate.md#thsas--testhasassertion--every-test-asserts-something-in-its-body) `THSAS-B03`

- [An empty test is a label only when the gate says so](camadas/gate.md#thsas--testhasassertion--every-test-asserts-something-in-its-body) `THSAS-B04`

- [The source's end bounds the body](camadas/gate.md#thsas--testhasassertion--every-test-asserts-something-in-its-body) `THSAS-B05`

- [Nothing to measure is skipped](camadas/gate.md#thsas--testhasassertion--every-test-asserts-something-in-its-body) `THSAS-B06`

- [Tests that cannot be listed fail the gate with the reason](camadas/gate.md#thsas--testhasassertion--every-test-asserts-something-in-its-body) `THSAS-E01`

- [A test that declares it asserts nothing is left out](camadas/gate.md#thsas--testhasassertion--every-test-asserts-something-in-its-body) `THSAS-B07`

- [A test at a line the content does not have](camadas/gate.md#thsas--testhasassertion--every-test-asserts-something-in-its-body) `THSAS-B08`

- [A closer that opens again goes on with the block](camadas/gate.md#thsas--testhasassertion--every-test-asserts-something-in-its-body) `THSAS-B09`

- [A node that is not a feature is skipped](camadas/gate.md#tlvcd--testlevelcodes--each-scenario-references-only-codes-its-test-level-accepts) `TLVCD-B01`

- [A project with no filter per level is skipped](camadas/gate.md#tlvcd--testlevelcodes--each-scenario-references-only-codes-its-test-level-accepts) `TLVCD-B02`

- [A feature with no filtered level is skipped](camadas/gate.md#tlvcd--testlevelcodes--each-scenario-references-only-codes-its-test-level-accepts) `TLVCD-B03`

- [A level with allow accepts only matching codes](camadas/gate.md#tlvcd--testlevelcodes--each-scenario-references-only-codes-its-test-level-accepts) `TLVCD-B04`

- [A level with exclude refuses matching codes even when allowed](camadas/gate.md#tlvcd--testlevelcodes--each-scenario-references-only-codes-its-test-level-accepts) `TLVCD-B05`

- [Every code of the scenario is confronted without its suffix](camadas/gate.md#tlvcd--testlevelcodes--each-scenario-references-only-codes-its-test-level-accepts) `TLVCD-B06`

- [The failure names each refused code with its level, sorted](camadas/gate.md#tlvcd--testlevelcodes--each-scenario-references-only-codes-its-test-level-accepts) `TLVCD-B07`

- [Accepted codes pass](camadas/gate.md#tlvcd--testlevelcodes--each-scenario-references-only-codes-its-test-level-accepts) `TLVCD-B08`

- [Levels come from the gate's own entries, and two entries add their lists](camadas/gate.md#tlvcd--testlevelcodes--each-scenario-references-only-codes-its-test-level-accepts) `TLVCD-B09`

- [A test that names what its unit defines exercises it](camadas/gate.md#tsrch--testreach--a-test-reaches-the-unit-it-says-it-tests) `TSRCH-B01`

- [An import of the unit's module reaches it](camadas/gate.md#tsrch--testreach--a-test-reaches-the-unit-it-says-it-tests) `TSRCH-B02`

- [A test that defines the unit's names exercises a copy](camadas/gate.md#tsrch--testreach--a-test-reaches-the-unit-it-says-it-tests) `TSRCH-B03`

- [Short names are anyone's](camadas/gate.md#tsrch--testreach--a-test-reaches-the-unit-it-says-it-tests) `TSRCH-B04`

- [The ref's unit must be reached](camadas/gate.md#tsrch--testreach--a-test-reaches-the-unit-it-says-it-tests) `TSRCH-B05`

- [A declared invocation reaches the unit it names](camadas/gate.md#tsrch--testreach--a-test-reaches-the-unit-it-says-it-tests) `TSRCH-B06`

- [Nothing to confront is skipped](camadas/gate.md#tsrch--testreach--a-test-reaches-the-unit-it-says-it-tests) `TSRCH-B07`

- [A declared way to reach the unit passes](camadas/gate.md#tsrch--testreach--a-test-reaches-the-unit-it-says-it-tests) `TSRCH-B08`

- [A member of the test's own type is not a copy of the unit's](camadas/gate.md#tsrch--testreach--a-test-reaches-the-unit-it-says-it-tests) `TSRCH-B09`

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

- [With a tests source a test traces only through its titles](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-B12`

- [A failing tests source fails the gate naming the error](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-E03`

- [A support file is not charged with tracing](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves) `TSTRT-B13`

- [The gate skips what is not a spec, and is Pending without a map](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B01`

- [Without a declared handle attribute the gate skips](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B02`

- [A spec with no readable linked code skips](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B03`

- [A unit that exposes no handle and declares none skips](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B04`

- [A handle exposed, declared and queried passes](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B05`

- [A handle the code exposes and the spec does not declare fails](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B06`

- [A handle the spec declares and the code does not expose fails](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B07`

- [A handle no consumer queries fails when a consumer surface exists](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B08`

- [With no consumer surface the queried end is not charged](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B09`

- [The handle attribute is the one the project declares](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B10`

- [A literal handle counts, also through a derived prop or an object key](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B11`

- [A template handle, or a prefix prop, exposes the prefix as a wildcard](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B12`

- [Only the branches of a conditional handle are handles, never its condition](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B13`

- [The inventory is read only inside the test surface section](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B14`

- [In a table only the first cell is the id; on any other line every quoted id counts](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B15`

- [The attribute's own name is not a declared id](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B16`

- [A wildcard at either end covers the concrete ids it opens](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B17`

- [A handle is queried when a consumer mentions it; a wildcard, when it mentions the prefix](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B18`

- [The e2e flows the project declares are consumers](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B19`

- [The test files beside the spec, and in its sibling folders, are consumers](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B20`

- [One handle is one line of the report, whatever the spelling at each end](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-I01`

- [A handle only a consumer mentions is not charged to the spec](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-X01`

- [A surface declared with a glob is read from the directory before the first wildcard](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B21`

- [A consumer that names only a longer id does not query the shorter one](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B22`

- [The test linked to the spec's feature is a consumer wherever it lives](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B23`

- [The report shows whether the feature describes each handle](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-B24`

- [Linked code that cannot be read is not an end](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it) `TICTS-E01`

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

- [A declared E2E directory missing from disk skips the gate](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-E01`

- [An unreadable flow or source leaves the verdict pending](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code) `TQETS-E02`

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

- [A declared pack that does not load leaves the verdict pending](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary) `TRDCT-E02`

- [An artifact that is not a spec leaves without a verdict](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-B01`

- [A recognised layer leaves without a verdict](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-B02`

- [Without a map the verdict is undetermined](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-B03`

- [A spec with the three pieces linked passes](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-B04`

- [A spec missing a piece is failed, and the verdict names which](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-B05`

- [The layer may waive a piece for every spec in it](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-B06`

- [The unit may waive a piece in its own spec, with a written reason](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-B07`

- [The test is reached in two hops, through the feature](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-I01`

- [Waiving the test while the feature carries a scenario is a contradiction](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-I02`

- [Waiving the test demands saying where the proof is, and the place must exist](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-I03`

- [A waiver covers only the piece it declares](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-I04`

- [The gate does not confront whether the pieces MATCH one another](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-X01`

- [A per-rule waiver in a table row does not waive the unit](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-B09`

- [A piece declared TO BE DEVELOPED leaves the verdict undetermined](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-B08`

- [The gate does not judge the QUALITY of any piece](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-X02`

- [A test missing from disk does not orphan a reference another test resolves](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-E01`

- [A feature missing from disk does not hide the scenario of the covered feature](camadas/gate.md#untcp--unitcomplete--the-pieces-that-realize-a-spec-exist) `UNTCP-E02`

- [A declaration matching its code line passes](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B01`

- [A declaration whose code line says another value fails](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B02`

- [Comment and blank lines between declaration and code are skipped](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B03`

- [Copies of the same key with different values are reported](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B04`

- [The value a rule declares in the spec is the source](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B05`

- [Without a declared pattern the gate skips and says how to enable it](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B06`

- [A declaration with nothing below annotates nothing](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B07`

- [An anchor inside prose is a mention, not a declaration](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-B08`

- [A pattern with fewer than two capture groups is treated as not declared](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-I01`

- [Declarations are indexed once per graph instance](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-I02`

- [With no built map the verdict is never approval](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-I03`

- [A prose rule declares no value](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-X01`

- [A literal nobody declared is not charged](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-X02`

- [A spec missing from disk does not drop the values of the other specs](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-E01`

- [A code file missing from disk does not hide the copies in the other files](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-E02`

- [A declared value anchor with too few groups is named, not called undeclared](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value) `VLANV-E03`

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

- [Nothing to confront leaves every gate without a verdict](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-B01`

- [Every state needs a VR scenario](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-B02`

- [Every VR scenario needs a VR test](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-B03`

- [Every VR scenario needs its baseline image, in any image format](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-B04`

- [A VR scenario of a state the spec does not register](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-B05`

- [A VR scenario of no state](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-B06`

- [A VR test of a scenario the feature does not declare](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-B07`

- [A VR scenario carries the regime tag and the state's code, in either form](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-B08`

- [Which tests are of the unit](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-B09`

- [The State letter comes from the project's rule types](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-B10`

- [vr-baseline accepts other image formats](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-B11`

- [One state never answers for another](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-I01`

- [A code file with no spec beside it leaves without a verdict](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-E01`

- [A unit with no feature has no VR scenario](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-E02`

- [A baseline that cannot be found counts as none](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-E03`

- [A test that cannot be read is read by its path](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-E04`

- [A state with no visual value is exempted with @no-vr and a reason](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-B12`

- [Every message is captured like a state](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-B13`

- [The states registered are those of the States section](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-B14`

- [A @no-vr exempts the code of its own row, never a code its reason mentions](camadas/gate.md#vrstc--vrstatescovered--each-state-of-a-visual-unit-tied-to-its-visual-regression-both-ways) `VRSTC-B15`

## infra

- [The key joins the stage and the normalised unit](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-B01`

- [A pending record lives in the changes folder under its key](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-B02`

- [The header carries the stage, the unit, the date and the agent only when named](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-B03`

- [The record states the intent and the touched files](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-B04`

- [Empty decision and proof sections are still written](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-B05`

- [Saving the same stage and unit again replaces the record](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-B06`

- [The pending list is the sorted markdown files of the changes folder](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-B07`

- [Marking a record reviewed moves it to the history under the same name](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-B08`

- [A second review of the same key keeps the first in the history](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-B09`

- [A line break or a comment close in a field cannot break the header](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-B10`

- [A reviewed record leaves the pending list and stays in the history](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-I01`

- [Two different units never share a key](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-I02`

- [A subfolder of changes is never listed as pending](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-X01`

- [A changes folder that cannot be read is an error](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-E01`

- [A changes folder that cannot be created fails the save](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-E02`

- [Marking a missing record reviewed fails](layers/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage) `CHRCC-E03`

- [A breaking change is never silent](layers/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits) `CHNGL-B01`

- [A feat is a feature](layers/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits) `CHNGL-B02`

- [A fix with a Bug footer is a bug fixed, without it a fix](layers/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits) `CHNGL-B03`

- [The internal types and free-form subjects are left out](layers/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits) `CHNGL-B04`

- [The trailing blocks of footers hold the footers](layers/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits) `CHNGL-B05`

- [Each release holds the commits after the previous tag](layers/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits) `CHNGL-B06`

- [What came after the last tag is unreleased](layers/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits) `CHNGL-B07`

- [The built-in template translates its headings](layers/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits) `CHNGL-B08`

- [A written release is marked by its version](layers/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits) `CHNGL-B09`

- [New releases go on top and the file keeps what it holds](layers/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits) `CHNGL-B10`

- [The unreleased block is replaced](layers/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits) `CHNGL-B11`

- [An empty file gets a title](layers/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits) `CHNGL-B12`

- [A template that does not parse is an error](layers/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits) `CHNGL-E01`

- [A git failure names the command](layers/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits) `CHNGL-E02`

- [The markers of a CRLF changelog are found](layers/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits) `CHNGL-B13`

- [Each scope is mirrored to its own file](layers/infra.md#chlgc--checklog--the-checks-output-mirrored-to-a-file-so-it-can-be-reread-without-re-running) `CHLGC-B01`

- [The header comes first and the output is copied after it](layers/infra.md#chlgc--checklog--the-checks-output-mirrored-to-a-file-so-it-can-be-reread-without-re-running) `CHLGC-B02`

- [A long output is copied whole without hanging](layers/infra.md#chlgc--checklog--the-checks-output-mirrored-to-a-file-so-it-can-be-reread-without-re-running) `CHLGC-B03`

- [A mirror that cannot be opened does not stop the check](layers/infra.md#chlgc--checklog--the-checks-output-mirrored-to-a-file-so-it-can-be-reread-without-re-running) `CHLGC-B04`

- [The header records the command, the moment and the HEAD](layers/infra.md#chlgc--checklog--the-checks-output-mirrored-to-a-file-so-it-can-be-reread-without-re-running) `CHLGC-B05`

- [The tree line reports the real count](layers/infra.md#chlgc--checklog--the-checks-output-mirrored-to-a-file-so-it-can-be-reread-without-re-running) `CHLGC-B06`

- [The changed-files mirror leaves the full snapshot intact](layers/infra.md#chlgc--checklog--the-checks-output-mirrored-to-a-file-so-it-can-be-reread-without-re-running) `CHLGC-I01`

- [A tree that could not be counted is never called clean](layers/infra.md#chlgc--checklog--the-checks-output-mirrored-to-a-file-so-it-can-be-reread-without-re-running) `CHLGC-X01`

- [Generic suffixes are dropped unless they are the whole name](layers/infra.md#cdgnc--codegenerator--the-short-stable-identity-code-suggested-for-a-units-name) `CDGNC-B01`

- [Words are split at separators, camel case and acronyms](layers/infra.md#cdgnc--codegenerator--the-short-stable-identity-code-suggested-for-a-units-name) `CDGNC-B02`

- [A long name takes the initials](layers/infra.md#cdgnc--codegenerator--the-short-stable-identity-code-suggested-for-a-units-name) `CDGNC-B03`

- [A two-word name takes letters from each word](layers/infra.md#cdgnc--codegenerator--the-short-stable-identity-code-suggested-for-a-units-name) `CDGNC-B04`

- [A single word takes consonants before vowels](layers/infra.md#cdgnc--codegenerator--the-short-stable-identity-code-suggested-for-a-units-name) `CDGNC-B05`

- [Short codes are padded with X and existing codes are completed the same way](layers/infra.md#cdgnc--codegenerator--the-short-stable-identity-code-suggested-for-a-units-name) `CDGNC-B06`

- [A module prefix starts the code](layers/infra.md#cdgnc--codegenerator--the-short-stable-identity-code-suggested-for-a-units-name) `CDGNC-B07`

- [A collision varies the last position and keeps the prefix](layers/infra.md#cdgnc--codegenerator--the-short-stable-identity-code-suggested-for-a-units-name) `CDGNC-B08`

- [The module prefix is the initial and the first consonant](layers/infra.md#cdgnc--codegenerator--the-short-stable-identity-code-suggested-for-a-units-name) `CDGNC-B09`

- [The generated length follows the smallest declared length](layers/infra.md#cdgnc--codegenerator--the-short-stable-identity-code-suggested-for-a-units-name) `CDGNC-B10`

- [Generation is deterministic and a resolved code is free](layers/infra.md#cdgnc--codegenerator--the-short-stable-identity-code-suggested-for-a-units-name) `CDGNC-I01`

- [A saturated namespace returns the generated code](layers/infra.md#cdgnc--codegenerator--the-short-stable-identity-code-suggested-for-a-units-name) `CDGNC-X01`

- [Generic file names are recognised in any case](layers/infra.md#cfpcd--codefrompath--the-most-meaningful-unique-code-for-a-unit-given-its-file-path) `CFPCD-B01`

- [A generic file name takes its code from the parent folder](layers/infra.md#cfpcd--codefrompath--the-most-meaningful-unique-code-for-a-unit-given-its-file-path) `CFPCD-B02`

- [Artifact suffixes are dropped from the base name](layers/infra.md#cfpcd--codefrompath--the-most-meaningful-unique-code-for-a-unit-given-its-file-path) `CFPCD-B03`

- [A normal file name gets the generated code when free](layers/infra.md#cfpcd--codefrompath--the-most-meaningful-unique-code-for-a-unit-given-its-file-path) `CFPCD-B04`

- [Same-named units get distinct, deterministic codes](layers/infra.md#cfpcd--codefrompath--the-most-meaningful-unique-code-for-a-unit-given-its-file-path) `CFPCD-I01`

- [The state files live in the project's state folder](layers/infra.md#dmstd--daemonstate--the-background-watchers-state-files-pid-pause-flag-log-and-meta) `DMSTD-B01`

- [Running answers the live PID and 0 otherwise](layers/infra.md#dmstd--daemonstate--the-background-watchers-state-files-pid-pause-flag-log-and-meta) `DMSTD-B02`

- [A PID file of an exited process is removed](layers/infra.md#dmstd--daemonstate--the-background-watchers-state-files-pid-pause-flag-log-and-meta) `DMSTD-B03`

- [Stopping with no watcher is refused](layers/infra.md#dmstd--daemonstate--the-background-watchers-state-files-pid-pause-flag-log-and-meta) `DMSTD-B04`

- [Stopping a running watcher terminates it and removes the PID file](layers/infra.md#dmstd--daemonstate--the-background-watchers-state-files-pid-pause-flag-log-and-meta) `DMSTD-B05`

- [Pause is the existence of the flag file](layers/infra.md#dmstd--daemonstate--the-background-watchers-state-files-pid-pause-flag-log-and-meta) `DMSTD-B06`

- [Cleanup removes the PID file and the pause flag](layers/infra.md#dmstd--daemonstate--the-background-watchers-state-files-pid-pause-flag-log-and-meta) `DMSTD-B07`

- [The meta records the start moment and the root](layers/infra.md#dmstd--daemonstate--the-background-watchers-state-files-pid-pause-flag-log-and-meta) `DMSTD-B08`

- [The PID file of an exited process does not outlive the check](layers/infra.md#dmstd--daemonstate--the-background-watchers-state-files-pid-pause-flag-log-and-meta) `DMSTD-I01`

- [A state folder that cannot be created fails the PID write](layers/infra.md#dmstd--daemonstate--the-background-watchers-state-files-pid-pause-flag-log-and-meta) `DMSTD-E02`

- [A live process is alive and an exited one is not](layers/infra.md#dmrnd--daemonruntime--how-each-platform-probes-and-terminates-the-background-watcher) `DMRND-B01`

- [Termination on Unix-like systems is a catchable SIGTERM](layers/infra.md#dmrnd--daemonruntime--how-each-platform-probes-and-terminates-the-background-watcher) `DMRND-B02`

- [A live process of another user is alive](layers/infra.md#dmrnd--daemonruntime--how-each-platform-probes-and-terminates-the-background-watcher) `DMRND-B03`

- [Stopping leaves no PID file even when the process cleans nothing](layers/infra.md#dmrnd--daemonruntime--how-each-platform-probes-and-terminates-the-background-watcher) `DMRND-I01`

- [Without git on the PATH the answer is no binary](layers/infra.md#gtavg--gitavailability--why-a-git-operation-cannot-happen-named-with-its-fix) `GTAVG-B01`

- [A repository above the root is seen](layers/infra.md#gtavg--gitavailability--why-a-git-operation-cannot-happen-named-with-its-fix) `GTAVG-B02`

- [A folder outside any repository is told apart from one inside](layers/infra.md#gtavg--gitavailability--why-a-git-operation-cannot-happen-named-with-its-fix) `GTAVG-B03`

- [The missing binary is explained with the install fix](layers/infra.md#gtavg--gitavailability--why-a-git-operation-cannot-happen-named-with-its-fix) `GTAVG-B04`

- [The missing repository is explained with the init fix](layers/infra.md#gtavg--gitavailability--why-a-git-operation-cannot-happen-named-with-its-fix) `GTAVG-B05`

- [An available git is not explained](layers/infra.md#gtavg--gitavailability--why-a-git-operation-cannot-happen-named-with-its-fix) `GTAVG-B06`

- [Today is the system date](layers/infra.md#gtmtg--gitmeta--what-git-knows-about-the-files-last-commit-dates-pending-changes-head-dirty-count) `GTMTG-B01`

- [The last commit date is the day of the most recent commit of the file](layers/infra.md#gtmtg--gitmeta--what-git-knows-about-the-files-last-commit-dates-pending-changes-head-dirty-count) `GTMTG-B02`

- [The bulk reader gives each file its most recent commit day](layers/infra.md#gtmtg--gitmeta--what-git-knows-about-the-files-last-commit-dates-pending-changes-head-dirty-count) `GTMTG-B03`

- [A new file in a repository has pending changes, and the question was asked](layers/infra.md#gtmtg--gitmeta--what-git-knows-about-the-files-last-commit-dates-pending-changes-head-dirty-count) `GTMTG-B04`

- [The shortcut says yes after an edit and no after a commit](layers/infra.md#gtmtg--gitmeta--what-git-knows-about-the-files-last-commit-dates-pending-changes-head-dirty-count) `GTMTG-B05`

- [HEAD is the short hash and subject, and unknown without a commit](layers/infra.md#gtmtg--gitmeta--what-git-knows-about-the-files-last-commit-dates-pending-changes-head-dirty-count) `GTMTG-B06`

- [The dirty count counts the modified files of a real repository](layers/infra.md#gtmtg--gitmeta--what-git-knows-about-the-files-last-commit-dates-pending-changes-head-dirty-count) `GTMTG-B07`

- [A tree that could not be counted is not reported clean](layers/infra.md#gtmtg--gitmeta--what-git-knows-about-the-files-last-commit-dates-pending-changes-head-dirty-count) `GTMTG-X01`

- [Outside a repository the pending-change question is not known](layers/infra.md#gtmtg--gitmeta--what-git-knows-about-the-files-last-commit-dates-pending-changes-head-dirty-count) `GTMTG-X02`

- [A file's content at the last commit](layers/infra.md#gtmtg--gitmeta--what-git-knows-about-the-files-last-commit-dates-pending-changes-head-dirty-count) `GTMTG-B08`

- [The commit's own changes, from the index](layers/infra.md#gtmtg--gitmeta--what-git-knows-about-the-files-last-commit-dates-pending-changes-head-dirty-count) `GTMTG-B09`

- [The files with uncommitted changes are read in one status, relative to the root](layers/infra.md#gtmtg--gitmeta--what-git-knows-about-the-files-last-commit-dates-pending-changes-head-dirty-count) `GTMTG-B10`

- [Only an administrator whose protection spares administrators can bypass it](layers/infra.md#aprcp--approvalreachable--the-doctor-says-when-the-required-approval-can-never-be-given-and-how-to-get-out) `APRCP-B01`

- [Each refusal of the bypass names its reason](layers/infra.md#aprcp--approvalreachable--the-doctor-says-when-the-required-approval-can-never-be-given-and-how-to-get-out) `APRCP-B02`

- [Zero required approvals or no configuration reports nothing](layers/infra.md#aprcp--approvalreachable--the-doctor-says-when-the-required-approval-can-never-be-given-and-how-to-get-out) `APRCP-B03`

- [Without the platform CLI the reachability check stays silent](layers/infra.md#aprcp--approvalreachable--the-doctor-says-when-the-required-approval-can-never-be-given-and-how-to-get-out) `APRCP-B04`

- [A required approval the account cannot bypass gives one warning on the repository](layers/infra.md#aprcp--approvalreachable--the-doctor-says-when-the-required-approval-can-never-be-given-and-how-to-get-out) `APRCP-B05`

- [The warning names both ways out](layers/infra.md#aprcp--approvalreachable--the-doctor-says-when-the-required-approval-can-never-be-given-and-how-to-get-out) `APRCP-B06`

- [Turning the requirement off sends a protection with zero approvals](layers/infra.md#aprcp--approvalreachable--the-doctor-says-when-the-required-approval-can-never-be-given-and-how-to-get-out) `APRCP-B07`

- [The warning appears exactly when the bypass is refused](layers/infra.md#aprcp--approvalreachable--the-doctor-says-when-the-required-approval-can-never-be-given-and-how-to-get-out) `APRCP-I01`

- [The reachability check never changes the branch protection](layers/infra.md#aprcp--approvalreachable--the-doctor-says-when-the-required-approval-can-never-be-given-and-how-to-get-out) `APRCP-X01`

- [A refused protection update names the branch and carries the platform output](layers/infra.md#aprcp--approvalreachable--the-doctor-says-when-the-required-approval-can-never-be-given-and-how-to-get-out) `APRCP-E01`

- [Local mode checks nothing](layers/infra.md#ghegt--githubenvironment--the-doctor-warns-before-the-work-starts-about-the-pieces-the-github-flow-silently-needs) `GHEGT-B01`

- [Each missing pipeline gives a warning that names what stops happening](layers/infra.md#ghegt--githubenvironment--the-doctor-warns-before-the-work-starts-about-the-pieces-the-github-flow-silently-needs) `GHEGT-B02`

- [A pipeline without serialization is its own finding, not a missing one](layers/infra.md#ghegt--githubenvironment--the-doctor-warns-before-the-work-starts-about-the-pieces-the-github-flow-silently-needs) `GHEGT-B03`

- [A pipeline behind its template is reported only while it carries the marker](layers/infra.md#ghegt--githubenvironment--the-doctor-warns-before-the-work-starts-about-the-pieces-the-github-flow-silently-needs) `GHEGT-B04`

- [An unreadable branch protection is an unprotected branch](layers/infra.md#ghegt--githubenvironment--the-doctor-warns-before-the-work-starts-about-the-pieces-the-github-flow-silently-needs) `GHEGT-B05`

- [A protection without required reviews is partially protected](layers/infra.md#ghegt--githubenvironment--the-doctor-warns-before-the-work-starts-about-the-pieces-the-github-flow-silently-needs) `GHEGT-B06`

- [Without the platform CLI the branch protection is not asked](layers/infra.md#ghegt--githubenvironment--the-doctor-warns-before-the-work-starts-about-the-pieces-the-github-flow-silently-needs) `GHEGT-B07`

- [The protection is read on the declared integration branch](layers/infra.md#ghegt--githubenvironment--the-doctor-warns-before-the-work-starts-about-the-pieces-the-github-flow-silently-needs) `GHEGT-B08`

- [After the fix seeds the pipelines no pipeline finding remains](layers/infra.md#ghegt--githubenvironment--the-doctor-warns-before-the-work-starts-about-the-pieces-the-github-flow-silently-needs) `GHEGT-I01`

- [The board is never charged](layers/infra.md#ghegt--githubenvironment--the-doctor-warns-before-the-work-starts-about-the-pieces-the-github-flow-silently-needs) `GHEGT-X01`

- [Tests and code without evidence freshness get the suggestion](layers/infra.md#gvopg--governanceopportunities--the-doctor-suggests-the-canonical-gates-and-settings-a-project-has-not-adopted-yet) `GVOPG-B01`

- [Code without the secret gate gets the suggestion](layers/infra.md#gvopg--governanceopportunities--the-doctor-suggests-the-canonical-gates-and-settings-a-project-has-not-adopted-yet) `GVOPG-B02`

- [A secret gate that does not block is suboptimal](layers/infra.md#gvopg--governanceopportunities--the-doctor-suggests-the-canonical-gates-and-settings-a-project-has-not-adopted-yet) `GVOPG-B03`

- [Code without dependency audit or duplication gates gets both suggestions](layers/infra.md#gvopg--governanceopportunities--the-doctor-suggests-the-canonical-gates-and-settings-a-project-has-not-adopted-yet) `GVOPG-B04`

- [Tests without JUnit output are a suboptimal configuration](layers/infra.md#gvopg--governanceopportunities--the-doctor-suggests-the-canonical-gates-and-settings-a-project-has-not-adopted-yet) `GVOPG-B05`

- [A nil map or configuration gives nothing](layers/infra.md#gvopg--governanceopportunities--the-doctor-suggests-the-canonical-gates-and-settings-a-project-has-not-adopted-yet) `GVOPG-B06`

- [The quick hints are the first two opportunities](layers/infra.md#gvopg--governanceopportunities--the-doctor-suggests-the-canonical-gates-and-settings-a-project-has-not-adopted-yet) `GVOPG-B07`

- [Every opportunity is informational](layers/infra.md#gvopg--governanceopportunities--the-doctor-suggests-the-canonical-gates-and-settings-a-project-has-not-adopted-yet) `GVOPG-I01`

- [What the project already adopted is not suggested](layers/infra.md#gvopg--governanceopportunities--the-doctor-suggests-the-canonical-gates-and-settings-a-project-has-not-adopted-yet) `GVOPG-X01`

- [Every applicable undeclared catalog gate is suggested](layers/infra.md#gvopg--governanceopportunities--the-doctor-suggests-the-canonical-gates-and-settings-a-project-has-not-adopted-yet) `GVOPG-B08`

- [A family's fallible patterns are suggested to a project that declared none, and declining silences it](layers/infra.md#gvopg--governanceopportunities--the-doctor-suggests-the-canonical-gates-and-settings-a-project-has-not-adopted-yet) `GVOPG-B09`

- [The report is sorted by check then subject, and counts the map](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B01`

- [Warnings keeps only the warnings](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B02`

- [A node whose file is gone is a ghost](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B03`

- [An edge to an unknown node is dead](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B04`

- [A spec that points at nothing has no realization](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B05`

- [A spec without a code has no identity](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B06`

- [A code owned by units of different domains is a duplicate identity](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B07`

- [Units of the same domain may share a code](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B08`

- [Only specs, features, tests and code own a code](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B09`

- [The files of one unit count as one owner](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B10`

- [A file that declares a shared code is not an owner](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B11`

- [A declared layer with no node of its kind is empty](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B12`

- [A guide that governs nothing is reported](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B13`

- [A spec, feature or test kind that no gate confronts is reported](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B14`

- [An unknown perspective in skip_on names the gate and the value](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B15`

- [A gate whose tool is not on the PATH is reported with its install hint](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B16`

- [A missing git binary is told to install, not to initialize](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B17`

- [A project under no repository is told to initialize one](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B18`

- [A .git in the root or an ancestor, as a folder or a file, is a repository](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B19`

- [In GitHub mode a missing repository names the work queue](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B20`

- [A needs to a plan that does not exist is broken](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B21`

- [A cycle of needs gives one finding naming only the cycle, the same on every run](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B22`

- [A chain in order, or no plan at all, gives nothing](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B23`

- [Tests without results and code without coverage are warnings](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B24`

- [Missing or partial mutation is informational](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B25`

- [A healthy project has no warnings](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-I01`

- [Code without a spec is not reported](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-X01`

- [An ingested report that reached no test is told apart from no report](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B26`

- [A gate declared over a field the project neither declares nor waives is named](layers/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project) `DCTRO-B27`

- [A spec with two open decisions gives one warning carrying the count](layers/infra.md#pndcp--pendingdecisions--the-doctor-lists-the-specs-that-still-hold-open-decisions-the-heaviest-first) `PNDCP-B01`

- [A section closed with none is not pending](layers/infra.md#pndcp--pendingdecisions--the-doctor-lists-the-specs-that-still-hold-open-decisions-the-heaviest-first) `PNDCP-B02`

- [The spec with the most open decisions comes first](layers/infra.md#pndcp--pendingdecisions--the-doctor-lists-the-specs-that-still-hold-open-decisions-the-heaviest-first) `PNDCP-B03`

- [A nil map or a map without open decisions gives nothing](layers/infra.md#pndcp--pendingdecisions--the-doctor-lists-the-specs-that-still-hold-open-decisions-the-heaviest-first) `PNDCP-B04`

- [The doctor's count is the check's count](layers/infra.md#pndcp--pendingdecisions--the-doctor-lists-the-specs-that-still-hold-open-decisions-the-heaviest-first) `PNDCP-I01`

- [The count follows the check's rule, not a reading of its own](layers/infra.md#pndcp--pendingdecisions--the-doctor-lists-the-specs-that-still-hold-open-decisions-the-heaviest-first) `PNDCP-X01`

- [A spec missing on disk is skipped and the others are still reported](layers/infra.md#pndcp--pendingdecisions--the-doctor-lists-the-specs-that-still-hold-open-decisions-the-heaviest-first) `PNDCP-E01`

- [A screen spec without any recommended section is told about all six](layers/infra.md#spscs--specsections--the-doctor-tells-when-most-specs-of-a-layer-lack-a-section-a-gate-needs-to-see) `SPSCS-B01`

- [Exactly half or a minority missing is silent](layers/infra.md#spscs--specsections--the-doctor-tells-when-most-specs-of-a-layer-lack-a-section-a-gate-needs-to-see) `SPSCS-B02`

- [A declared gate that is blind gives a warning](layers/infra.md#spscs--specsections--the-doctor-tells-when-most-specs-of-a-layer-lack-a-section-a-gate-needs-to-see) `SPSCS-B03`

- [An undeclared gate gives an informational recommendation that asks to declare it](layers/infra.md#spscs--specsections--the-doctor-tells-when-most-specs-of-a-layer-lack-a-section-a-gate-needs-to-see) `SPSCS-B04`

- [A title in another language counts](layers/infra.md#spscs--specsections--the-doctor-tells-when-most-specs-of-a-layer-lack-a-section-a-gate-needs-to-see) `SPSCS-B05`

- [The header layer wins over the map's](layers/infra.md#spscs--specsections--the-doctor-tells-when-most-specs-of-a-layer-lack-a-section-a-gate-needs-to-see) `SPSCS-B06`

- [Screen sections are not demanded from other layers](layers/infra.md#spscs--specsections--the-doctor-tells-when-most-specs-of-a-layer-lack-a-section-a-gate-needs-to-see) `SPSCS-B07`

- [Three specs lacking a section give one finding with the numbers](layers/infra.md#spscs--specsections--the-doctor-tells-when-most-specs-of-a-layer-lack-a-section-a-gate-needs-to-see) `SPSCS-B08`

- [A nil map or configuration gives nothing](layers/infra.md#spscs--specsections--the-doctor-tells-when-most-specs-of-a-layer-lack-a-section-a-gate-needs-to-see) `SPSCS-B09`

- [Adding the reported section removes the finding](layers/infra.md#spscs--specsections--the-doctor-tells-when-most-specs-of-a-layer-lack-a-section-a-gate-needs-to-see) `SPSCS-I01`

- [No finding names a spec file](layers/infra.md#spscs--specsections--the-doctor-tells-when-most-specs-of-a-layer-lack-a-section-a-gate-needs-to-see) `SPSCS-X01`

- [Specs missing on disk are left out of the counts](layers/infra.md#spscs--specsections--the-doctor-tells-when-most-specs-of-a-layer-lack-a-section-a-gate-needs-to-see) `SPSCS-E01`

- [The artifact options are spec, feature, test, guide, plan and code, in that order](layers/infra.md#archr--artifactchoice--turns-the-artifacts-the-user-chose-at-init-into-artifact-layers-and-colocation) `ARCHR-B01`

- [Each artifact inference found is pre-checked, and nothing else](layers/infra.md#archr--artifactchoice--turns-the-artifacts-the-user-chose-at-init-into-artifact-layers-and-colocation) `ARCHR-B02`

- [A chosen artifact layer is created even when nothing was detected](layers/infra.md#archr--artifactchoice--turns-the-artifacts-the-user-chose-at-init-into-artifact-layers-and-colocation) `ARCHR-B03`

- [An artifact layer that was not chosen is removed, and code layers are untouched](layers/infra.md#archr--artifactchoice--turns-the-artifacts-the-user-chose-at-init-into-artifact-layers-and-colocation) `ARCHR-B04`

- [A chosen artifact layer that already exists is kept as declared](layers/infra.md#archr--artifactchoice--turns-the-artifacts-the-user-chose-at-init-into-artifact-layers-and-colocation) `ARCHR-B05`

- [The guide and plan layers take the detected directory, otherwise the default pattern](layers/infra.md#archr--artifactchoice--turns-the-artifacts-the-user-chose-at-init-into-artifact-layers-and-colocation) `ARCHR-B06`

- [Colocation is declared from the spec as anchor, and the spec is never a derivative](layers/infra.md#archr--artifactchoice--turns-the-artifacts-the-user-chose-at-init-into-artifact-layers-and-colocation) `ARCHR-B07`

- [No colocation is declared when it is not wanted, when spec is not chosen, or when nothing derives from the spec](layers/infra.md#archr--artifactchoice--turns-the-artifacts-the-user-chose-at-init-into-artifact-layers-and-colocation) `ARCHR-B08`

- [Choosing code or leaving it out creates and removes no layer](layers/infra.md#archr--artifactchoice--turns-the-artifacts-the-user-chose-at-init-into-artifact-layers-and-colocation) `ARCHR-B09`

- [Colocation keeps the rest of derived](layers/infra.md#archr--artifactchoice--turns-the-artifacts-the-user-chose-at-init-into-artifact-layers-and-colocation) `ARCHR-B10`

- [A created test layer takes the project's test pattern](layers/infra.md#archr--artifactchoice--turns-the-artifacts-the-user-chose-at-init-into-artifact-layers-and-colocation) `ARCHR-B11`

- [The colocated test template is the project's](layers/infra.md#archr--artifactchoice--turns-the-artifacts-the-user-chose-at-init-into-artifact-layers-and-colocation) `ARCHR-B12`

- [Each detected code directory becomes a code layer named after its last segment](layers/infra.md#blcnb--buildconfig--builds-the-configuration-that-inference-proposes-as-the-default-for-the-init-questions) `BLCNB-B01`

- [With several detected extensions the pattern lists them as a set](layers/infra.md#blcnb--buildconfig--builds-the-configuration-that-inference-proposes-as-the-default-for-the-init-questions) `BLCNB-B02`

- [A proposed code layer excludes specs, features and test files](layers/infra.md#blcnb--buildconfig--builds-the-configuration-that-inference-proposes-as-the-default-for-the-init-questions) `BLCNB-B03`

- [Colocation is proposed only when detected, with templates only for the detected kinds](layers/infra.md#blcnb--buildconfig--builds-the-configuration-that-inference-proposes-as-the-default-for-the-init-questions) `BLCNB-B04`

- [The test handle is proposed only when inference found one](layers/infra.md#blcnb--buildconfig--builds-the-configuration-that-inference-proposes-as-the-default-for-the-init-questions) `BLCNB-B05`

- [The proposal creates no artifact layer and no governs rule](layers/infra.md#blcnb--buildconfig--builds-the-configuration-that-inference-proposes-as-the-default-for-the-init-questions) `BLCNB-X01`

- [The colocated test template follows the project's test convention](layers/infra.md#blcnb--buildconfig--builds-the-configuration-that-inference-proposes-as-the-default-for-the-init-questions) `BLCNB-B06`

- [The proposal's dialect is the family inference found](layers/infra.md#blcnb--buildconfig--builds-the-configuration-that-inference-proposes-as-the-default-for-the-init-questions) `BLCNB-B07`

- [The test layer's pattern is the project's test convention](layers/infra.md#blcnb--buildconfig--builds-the-configuration-that-inference-proposes-as-the-default-for-the-init-questions) `BLCNB-B08`

- [A new project's dialect carries what its family knows can fail](layers/infra.md#blcnb--buildconfig--builds-the-configuration-that-inference-proposes-as-the-default-for-the-init-questions) `BLCNB-B09`

- [The whole guide is a title, a seeding note and the section](layers/infra.md#cngdc--contributingguide--render-the-projects-contributingmd-from-the-configuration-init-writes) `CNGDC-B01`

- [The section states the order of the work](layers/infra.md#cngdc--contributingguide--render-the-projects-contributingmd-from-the-configuration-init-writes) `CNGDC-B02`

- [Where each piece lives](layers/infra.md#cngdc--contributingguide--render-the-projects-contributingmd-from-the-configuration-init-writes) `CNGDC-B03`

- [The declared code layers, or that there are none](layers/infra.md#cngdc--contributingguide--render-the-projects-contributingmd-from-the-configuration-init-writes) `CNGDC-B04`

- [The queue the daily commands name](layers/infra.md#cngdc--contributingguide--render-the-projects-contributingmd-from-the-configuration-init-writes) `CNGDC-B05`

- [What blocks and what informs](layers/infra.md#cngdc--contributingguide--render-the-projects-contributingmd-from-the-configuration-init-writes) `CNGDC-B06`

- [The project's guide folder](layers/infra.md#cngdc--contributingguide--render-the-projects-contributingmd-from-the-configuration-init-writes) `CNGDC-B07`

- [The guide never names what the configuration does not hold](layers/infra.md#cngdc--contributingguide--render-the-projects-contributingmd-from-the-configuration-init-writes) `CNGDC-I01`

- [Rendering writes nothing to disk](layers/infra.md#cngdc--contributingguide--render-the-projects-contributingmd-from-the-configuration-init-writes) `CNGDC-X01`

- [The code layer names are listed sorted, and only code layers](layers/infra.md#indcn--initdecisions--the-pure-decisions-of-init-over-the-proposed-configuration-code-layers-tags-and-governs-rules) `INDCN-B01`

- [Pruning removes the code layers not kept and never an artifact layer](layers/infra.md#indcn--initdecisions--the-pure-decisions-of-init-over-the-proposed-configuration-code-layers-tags-and-governs-rules) `INDCN-B02`

- [The candidate tags are the union of every layer's tags, deduplicated and sorted](layers/infra.md#indcn--initdecisions--the-pure-decisions-of-init-over-the-proposed-configuration-code-layers-tags-and-governs-rules) `INDCN-B03`

- [One governs rule per guide answered with a tag, ordered by guide, skipping the unanswered and none](layers/infra.md#indcn--initdecisions--the-pure-decisions-of-init-over-the-proposed-configuration-code-layers-tags-and-governs-rules) `INDCN-B04`

- [No answer gives no governs rule](layers/infra.md#indcn--initdecisions--the-pure-decisions-of-init-over-the-proposed-configuration-code-layers-tags-and-governs-rules) `INDCN-B05`

- [After pruning, the code layer names are exactly the kept code layers](layers/infra.md#indcn--initdecisions--the-pure-decisions-of-init-over-the-proposed-configuration-code-layers-tags-and-governs-rules) `INDCN-I01`

- [A project with no artifact chosen is seeded with no gate](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B01`

- [A project with spec, feature and test is born with the gates of each](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B02`

- [The gates that cross spec and feature are seeded only when both are chosen](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B03`

- [Choosing only guides seeds only the guide checklist gate](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B04`

- [The parent gate confronts exactly the chosen artifacts among spec, plan and code](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B05`

- [An existing project is born with its gates informative, except the two blocking by nature](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B06`

- [A new project is born with its gates blocking](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B07`

- [A gate that depends on an ingested signal stays informative even in a new project](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B08`

- [Every judgment gate that asks about code or a test carries the @TBD instruction](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B09`

- [The canonical declaration of a gate is found by name in the catalog of every artifact](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B10`

- [Loading a configuration completes a canonical gate declared by name alone](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B11`

- [The list of seeded gates does not change with the age of the project](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-I01`

- [Every gate of the full catalog has a unique name that is also its id](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-I02`

- [Every canonical name the migration renames a legacy gate to is a default gate](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-I03`

- [No default gate carries a legacy name](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-X01`

- [The gate is of the blocking class](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B12`

- [Choosing every artifact init offers seeds every gate of the catalog](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B13`

- [The gates that run on specs, features and tests are seeded without plans](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B14`

- [The gate names registered for the vocabulary check are the full catalog](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B15`

- [Choosing specs seeds header-valid on every governed kind](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B16`

- [no-duplication is the native duplication check on code files](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B17`

- [Choosing plans adds only gates that run on plans](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-I04`

- [No default gate writes Portuguese into the project](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-X02`

- [The mock dialect judgment presupposes the pattern it asks about](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B18`

- [rule-fulfilled is judged and marked to review](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B19`

- [The catalog carries the checkers that measure the unit's kinds](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B20`

- [Init seeds what relates to the project and has its premise](layers/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog) `DFGTD-B21`

- [The local board page is the published page](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B01`

- [The collect expression comes out of the pipeline](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B02`

- [Both collect shapes are accepted and cut before the redirection](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B03`

- [The active agents are counted from the snapshot](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B04`

- [A live board counts the active agents from now](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B05`

- [The agent chips filter the cards by agent](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B06`

- [A released card is collected with no owner](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B07`

- [With no card waiting for a person the blocked strip is hidden and empty](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B08`

- [The strip orders the waiting cards by how many cards each blocks](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B09`

- [An escalated card in a work state appears in the strip with the state it stopped in](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B10`

- [The strip tells a card waiting for a person from one waiting for decided work](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B11`

- [A hostile title is escaped in the strip](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B12`

- [An escalated card is marked in its column](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B13`

- [The roadmap opens a card's details](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B14`

- [A blocked card says which card blocks it](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B15`

- [Decisions and framing requests are listed apart](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B16`

- [Bugs are listed apart from decisions](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-B17`

- [The extracted expression never carries the redirection](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-I01`

- [The page's data attributes are written and read in pairs](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-I02`

- [A collect step of another shape yields no expression](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-X01`

- [A binary without the board files fails loudly](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-E01`

- [A pipeline that changed shape fails naming the change](layers/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes) `BREXB-E02`

- [With git not installed the state is not installed](layers/infra.md#gtstg--gitstate--classify-the-projects-versioning-before-init-scans-it-and-say-what-to-do-about-it) `GTSTG-B01`

- [With git installed and no repository anywhere the state is not initialised](layers/infra.md#gtstg--gitstate--classify-the-projects-versioning-before-init-scans-it-and-say-what-to-do-about-it) `GTSTG-B02`

- [A repository with a branch reference is ready](layers/infra.md#gtstg--gitstate--classify-the-projects-versioning-before-init-scans-it-and-say-what-to-do-about-it) `GTSTG-B03`

- [A repository with packed references is ready](layers/infra.md#gtstg--gitstate--classify-the-projects-versioning-before-init-scans-it-and-say-what-to-do-about-it) `GTSTG-B04`

- [A repository with no reference at all has no commit yet](layers/infra.md#gtstg--gitstate--classify-the-projects-versioning-before-init-scans-it-and-say-what-to-do-about-it) `GTSTG-B05`

- [A subfolder of an existing repository is ready](layers/infra.md#gtstg--gitstate--classify-the-projects-versioning-before-init-scans-it-and-say-what-to-do-about-it) `GTSTG-B06`

- [A repository marker that is a file is ready](layers/infra.md#gtstg--gitstate--classify-the-projects-versioning-before-init-scans-it-and-say-what-to-do-about-it) `GTSTG-B07`

- [Only the not-initialised and no-commit states have an action to offer](layers/infra.md#gtstg--gitstate--classify-the-projects-versioning-before-init-scans-it-and-say-what-to-do-about-it) `GTSTG-B08`

- [Each unready state has its own warning and ready has none](layers/infra.md#gtstg--gitstate--classify-the-projects-versioning-before-init-scans-it-and-say-what-to-do-about-it) `GTSTG-B09`

- [The seeded ignore list covers what Anchors generates](layers/infra.md#gtstg--gitstate--classify-the-projects-versioning-before-init-scans-it-and-say-what-to-do-about-it) `GTSTG-B10`

- [Not installed and not initialised stay two states with different offers](layers/infra.md#gtstg--gitstate--classify-the-projects-versioning-before-init-scans-it-and-say-what-to-do-about-it) `GTSTG-I01`

- [The seeded ignore list does not guess the stack](layers/infra.md#gtstg--gitstate--classify-the-projects-versioning-before-init-scans-it-and-say-what-to-do-about-it) `GTSTG-X01`

- [Whether git is installed comes from the caller](layers/infra.md#gtstg--gitstate--classify-the-projects-versioning-before-init-scans-it-and-say-what-to-do-about-it) `GTSTG-X02`

- [The comment dialect follows the language family](layers/infra.md#hdgdh--headerguide--render-the-projects-header-guide-in-the-languages-comment-dialect-passing-the-gate-that-init-itself-declares) `HDGDH-B01`

- [The grouping example names the first module, or auth](layers/infra.md#hdgdh--headerguide--render-the-projects-header-guide-in-the-languages-comment-dialect-passing-the-gate-that-init-itself-declares) `HDGDH-B02`

- [The module list appears only when there are modules](layers/infra.md#hdgdh--headerguide--render-the-projects-header-guide-in-the-languages-comment-dialect-passing-the-gate-that-init-itself-declares) `HDGDH-B03`

- [The essentials are always present](layers/infra.md#hdgdh--headerguide--render-the-projects-header-guide-in-the-languages-comment-dialect-passing-the-gate-that-init-itself-declares) `HDGDH-B04`

- [The compliance-points section is always present with five points](layers/infra.md#hdgdh--headerguide--render-the-projects-header-guide-in-the-languages-comment-dialect-passing-the-gate-that-init-itself-declares) `HDGDH-B05`

- [The title falls back to project](layers/infra.md#hdgdh--headerguide--render-the-projects-header-guide-in-the-languages-comment-dialect-passing-the-gate-that-init-itself-declares) `HDGDH-B06`

- [The seeded guide passes the checklist heading in every language](layers/infra.md#hdgdh--headerguide--render-the-projects-header-guide-in-the-languages-comment-dialect-passing-the-gate-that-init-itself-declares) `HDGDH-I01`

- [Rendering writes nothing to disk](layers/infra.md#hdgdh--headerguide--render-the-projects-header-guide-in-the-languages-comment-dialect-passing-the-gate-that-init-itself-declares) `HDGDH-X01`

- [Dependency, build and tool directories are not walked](layers/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm) `INPRN-B01`

- [The presence of specs, features and tests is detected](layers/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm) `INPRN-B02`

- [A markdown file in a plans directory is a plan, even inside a guides directory](layers/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm) `INPRN-B03`

- [The first guides directory is detected with every guide file in it](layers/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm) `INPRN-B04`

- [Every folder holding code is a candidate code directory, with no minimum, ordered by volume](layers/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm) `INPRN-B05`

- [The code extensions are the five most frequent, most frequent first](layers/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm) `INPRN-B06`

- [Colocation is detected when at least three stems pair code with a spec, test or feature](layers/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm) `INPRN-B07`

- [The test handle is the known attribute used most, and only from five uses](layers/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm) `INPRN-B08`

- [Inference carries a proposed configuration with code layers and no artifact layer](layers/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm) `INPRN-B09`

- [Only the first three hundred code files are read in search of the test handle](layers/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm) `INPRN-X01`

- [A root that cannot be walked fails the inference with no proposal](layers/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm) `INPRN-E01`

- [A test named in any known dialect pairs with the code of the same stem](layers/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm) `INPRN-B10`

- [The same tree always gives the same code extensions and code directories, ties included](layers/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm) `INPRN-I02`

- [The language family comes from the root manifest, or from the most frequent extension](layers/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm) `INPRN-B11`

- [The test conventions are the forms the project's tests follow, most followed first](layers/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm) `INPRN-B12`

- [The @TBD instruction forbids pass, orders a waiver naming the absence, and names the piece asked about](layers/infra.md#inctn--initcatalogs--the-language-dialect-catalog-and-the-tbd-instruction-that-init-seeds-into-judgment-gates) `INCTN-B04`

- [A file is a test when its name carries a convention's prefix and suffix](layers/infra.md#inctn--initcatalogs--the-language-dialect-catalog-and-the-tbd-instruction-that-init-seeds-into-judgment-gates) `INCTN-B06`

- [A convention gives its glob and its template](layers/infra.md#inctn--initcatalogs--the-language-dialect-catalog-and-the-tbd-instruction-that-init-seeds-into-judgment-gates) `INCTN-B07`

- [A test file's unit name drops the convention's prefix and suffix](layers/infra.md#inctn--initcatalogs--the-language-dialect-catalog-and-the-tbd-instruction-that-init-seeds-into-judgment-gates) `INCTN-B08`

- [With no convention read, the family default is used](layers/infra.md#inctn--initcatalogs--the-language-dialect-catalog-and-the-tbd-instruction-that-init-seeds-into-judgment-gates) `INCTN-B09`

- [A family's coverage hint names the reports ingest reads](layers/infra.md#inctn--initcatalogs--the-language-dialect-catalog-and-the-tbd-instruction-that-init-seeds-into-judgment-gates) `INCTN-B10`

- [A family's gate script runs in its own language](layers/infra.md#inctn--initcatalogs--the-language-dialect-catalog-and-the-tbd-instruction-that-init-seeds-into-judgment-gates) `INCTN-B11`

- [The @TBD instruction demands checking that the @TBD is still true](layers/infra.md#inctn--initcatalogs--the-language-dialect-catalog-and-the-tbd-instruction-that-init-seeds-into-judgment-gates) `INCTN-I02`

- [No convention of the catalog is shadowed by a shorter one](layers/infra.md#inctn--initcatalogs--the-language-dialect-catalog-and-the-tbd-instruction-that-init-seeds-into-judgment-gates) `INCTN-I03`

- [Every family default is a convention of the catalog](layers/infra.md#inctn--initcatalogs--the-language-dialect-catalog-and-the-tbd-instruction-that-init-seeds-into-judgment-gates) `INCTN-I04`

- [The catalog names no folder](layers/infra.md#inctn--initcatalogs--the-language-dialect-catalog-and-the-tbd-instruction-that-init-seeds-into-judgment-gates) `INCTN-X02`

- [A known agent variable makes the operator an AI even with a terminal](layers/infra.md#opdtp--operatordetection--tell-whether-a-person-or-an-ai-is-running-init-and-whether-the-discovery-phase-is-still-to-be-done) `OPDTP-B01`

- [The generic AI_AGENT variable makes the operator an AI even with a terminal](layers/infra.md#opdtp--operatordetection--tell-whether-a-person-or-an-ai-is-running-init-and-whether-the-discovery-phase-is-still-to-be-done) `OPDTP-B02`

- [No terminal and no agent variable is an AI](layers/infra.md#opdtp--operatordetection--tell-whether-a-person-or-an-ai-is-running-init-and-whether-the-discovery-phase-is-still-to-be-done) `OPDTP-B03`

- [A terminal and no agent variable is a person](layers/infra.md#opdtp--operatordetection--tell-whether-a-person-or-an-ai-is-running-init-and-whether-the-discovery-phase-is-still-to-be-done) `OPDTP-B04`

- [The tool is named after the known variable that is set](layers/infra.md#opdtp--operatordetection--tell-whether-a-person-or-an-ai-is-running-init-and-whether-the-discovery-phase-is-still-to-be-done) `OPDTP-B05`

- [Several known variables resolve by the first variable name in order](layers/infra.md#opdtp--operatordetection--tell-whether-a-person-or-an-ai-is-running-init-and-whether-the-discovery-phase-is-still-to-be-done) `OPDTP-B06`

- [The discovery phase is due only when nothing is found and nothing is described](layers/infra.md#opdtp--operatordetection--tell-whether-a-person-or-an-ai-is-running-init-and-whether-the-discovery-phase-is-still-to-be-done) `OPDTP-B07`

- [The project description is recognised under its three spellings](layers/infra.md#opdtp--operatordetection--tell-whether-a-person-or-an-ai-is-running-init-and-whether-the-discovery-phase-is-still-to-be-done) `OPDTP-B08`

- [Only the tools with a stable command line get a command](layers/infra.md#opdtp--operatordetection--tell-whether-a-person-or-an-ai-is-running-init-and-whether-the-discovery-phase-is-still-to-be-done) `OPDTP-B09`

- [The discovery prompt points at the guide and back at init](layers/infra.md#opdtp--operatordetection--tell-whether-a-person-or-an-ai-is-running-init-and-whether-the-discovery-phase-is-still-to-be-done) `OPDTP-B10`

- [The whole prompt is one argument of the command](layers/infra.md#opdtp--operatordetection--tell-whether-a-person-or-an-ai-is-running-init-and-whether-the-discovery-phase-is-still-to-be-done) `OPDTP-I01`

- [No command is offered for a tool without a stable command line](layers/infra.md#opdtp--operatordetection--tell-whether-a-person-or-an-ai-is-running-init-and-whether-the-discovery-phase-is-still-to-be-done) `OPDTP-X01`

- [Seeding copies every carried pack byte for byte](layers/infra.md#pcsdp--packseeding--copy-the-compliance-packs-carried-in-the-binary-into-the-project-never-over-an-adapted-one) `PCSDP-B01`

- [An adapted pack is preserved, not overwritten](layers/infra.md#pcsdp--packseeding--copy-the-compliance-packs-carried-in-the-binary-into-the-project-never-over-an-adapted-one) `PCSDP-B02`

- [The created and preserved lists come back sorted](layers/infra.md#pcsdp--packseeding--copy-the-compliance-packs-carried-in-the-binary-into-the-project-never-over-an-adapted-one) `PCSDP-B03`

- [The available packs are grouped by domain](layers/infra.md#pcsdp--packseeding--copy-the-compliance-packs-carried-in-the-binary-into-the-project-never-over-an-adapted-one) `PCSDP-B04`

- [Seeding twice is the same as seeding once](layers/infra.md#pcsdp--packseeding--copy-the-compliance-packs-carried-in-the-binary-into-the-project-never-over-an-adapted-one) `PCSDP-I01`

- [Seeding does not filter by the adopted jurisdiction](layers/infra.md#pcsdp--packseeding--copy-the-compliance-packs-carried-in-the-binary-into-the-project-never-over-an-adapted-one) `PCSDP-X01`

- [A folder that cannot be written fails the seeding](layers/infra.md#pcsdp--packseeding--copy-the-compliance-packs-carried-in-the-binary-into-the-project-never-over-an-adapted-one) `PCSDP-E01`

- [The questions come in the order of the terminal UI](layers/infra.md#inqsn--initquestions--describe-the-human-decisions-of-init-so-an-agent-can-answer-them-without-the-terminal-ui-and-judge-every-answer) `INQSN-B01`

- [Every question carries what the agent needs to decide](layers/infra.md#inqsn--initquestions--describe-the-human-decisions-of-init-so-an-agent-can-answer-them-without-the-terminal-ui-and-judge-every-answer) `INQSN-B02`

- [The defaults come from the inference](layers/infra.md#inqsn--initquestions--describe-the-human-decisions-of-init-so-an-agent-can-answer-them-without-the-terminal-ui-and-judge-every-answer) `INQSN-B03`

- [The work-queue mode has its choices and default](layers/infra.md#inqsn--initquestions--describe-the-human-decisions-of-init-so-an-agent-can-answer-them-without-the-terminal-ui-and-judge-every-answer) `INQSN-B04`

- [Every question gets a verdict and unanswered ones take the default](layers/infra.md#inqsn--initquestions--describe-the-human-decisions-of-init-so-an-agent-can-answer-them-without-the-terminal-ui-and-judge-every-answer) `INQSN-B05`

- [An answer outside the options is refused with the accepted values](layers/infra.md#inqsn--initquestions--describe-the-human-decisions-of-init-so-an-agent-can-answer-them-without-the-terminal-ui-and-judge-every-answer) `INQSN-B06`

- [The github mode requires repository and labels](layers/infra.md#inqsn--initquestions--describe-the-human-decisions-of-init-so-an-agent-can-answer-them-without-the-terminal-ui-and-judge-every-answer) `INQSN-B07`

- [A repository outside the github mode is refused](layers/infra.md#inqsn--initquestions--describe-the-human-decisions-of-init-so-an-agent-can-answer-them-without-the-terminal-ui-and-judge-every-answer) `INQSN-B08`

- [One refused answer refuses the whole set](layers/infra.md#inqsn--initquestions--describe-the-human-decisions-of-init-so-an-agent-can-answer-them-without-the-terminal-ui-and-judge-every-answer) `INQSN-B09`

- [An answer given empty is not the default](layers/infra.md#inqsn--initquestions--describe-the-human-decisions-of-init-so-an-agent-can-answer-them-without-the-terminal-ui-and-judge-every-answer) `INQSN-B10`

- [No answer goes missing from the verdict](layers/infra.md#inqsn--initquestions--describe-the-human-decisions-of-init-so-an-agent-can-answer-them-without-the-terminal-ui-and-judge-every-answer) `INQSN-I01`

- [An invalid answer is not corrected](layers/infra.md#inqsn--initquestions--describe-the-human-decisions-of-init-so-an-agent-can-answer-them-without-the-terminal-ui-and-judge-every-answer) `INQSN-X01`

- [The artifacts question offers the artifact options, code included](layers/infra.md#inqsn--initquestions--describe-the-human-decisions-of-init-so-an-agent-can-answer-them-without-the-terminal-ui-and-judge-every-answer) `INQSN-B11`

- [The question texts, their reasons and the refusal details are in the project's language](layers/infra.md#inqsn--initquestions--describe-the-human-decisions-of-init-so-an-agent-can-answer-them-without-the-terminal-ui-and-judge-every-answer) `INQSN-B12`

- [Every declared pipeline has a carried template, a role and its serialization need](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B01`

- [A pipeline is missing only when its file is absent](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B02`

- [A serial pipeline without serialization is flagged](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B03`

- [Seeding writes the missing pipelines and returns them sorted](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B04`

- [Seeding leaves a pipeline without the marker untouched](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B05`

- [The board page is seeded outside the pipelines folder, with the marker](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B06`

- [The board page seeding reports created, updated or unchanged](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B07`

- [An intact pipeline that differs from its template is outdated](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B08`

- [The integration branch, main when none is declared, replaces every marked branch line](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B09`

- [Anchors writes the board columns only up to READY TO TEST](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B10`

- [A per-card label is its prefix followed by the card](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B11`

- [The stale pipeline releases an idle owned card once](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B12`

- [A rejected card is released after the shorter rework window](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B13`

- [The review status follows the assigned reviewer's verdict line](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B14`

- [The review job is wired to verdict comments and feeds the mover](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B15`

- [A merge without the review outcome is said on the card](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B16`

- [A green PR publishes the review as pending](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B17`

- [The claim teaches the reviewer the verdict line](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B18`

- [A verdict releases the reviewer once](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B19`

- [The closing line wins over a reference](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B20`

- [A rejection sends the card back to its author](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-B21`

- [What seeding writes is never outdated for the same configuration](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-I01`

- [The parsed verdict line is the one the review guide teaches](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-I02`

- [A file without the marker is never taken over](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-X01`

- [A pipelines folder that cannot be created fails the seeding](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-E01`

- [A pipeline that cannot be written fails the seeding](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-E02`

- [A template the binary does not carry fails the seeding](layers/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns) `FLWRF-E03`

- [A card whose marker only starts with the key is not the issue's](layers/infra.md#ghigt--githubissues--the-issue-lifecycle-on-the-repositorys-cards-when-the-project-works-on-github) `GHIGT-B01`

- [The title names the gate, the kind and the target](layers/infra.md#ghigt--githubissues--the-issue-lifecycle-on-the-repositorys-cards-when-the-project-works-on-github) `GHIGT-B02`

- [The search covers every state in the configured repository](layers/infra.md#ghigt--githubissues--the-issue-lifecycle-on-the-repositorys-cards-when-the-project-works-on-github) `GHIGT-B03`

- [A new card carries the marker and the labels](layers/infra.md#ghigt--githubissues--the-issue-lifecycle-on-the-repositorys-cards-when-the-project-works-on-github) `GHIGT-B04`

- [An assumed debt's card has no flow label](layers/infra.md#ghigt--githubissues--the-issue-lifecycle-on-the-repositorys-cards-when-the-project-works-on-github) `GHIGT-B05`

- [An open card is left alone](layers/infra.md#ghigt--githubissues--the-issue-lifecycle-on-the-repositorys-cards-when-the-project-works-on-github) `GHIGT-B06`

- [A closed card is reopened with the new report](layers/infra.md#ghigt--githubissues--the-issue-lifecycle-on-the-repositorys-cards-when-the-project-works-on-github) `GHIGT-B07`

- [Resolving closes only an open card](layers/infra.md#ghigt--githubissues--the-issue-lifecycle-on-the-repositorys-cards-when-the-project-works-on-github) `GHIGT-B08`

- [The labels applied are the ones init creates](layers/infra.md#ghigt--githubissues--the-issue-lifecycle-on-the-repositorys-cards-when-the-project-works-on-github) `GHIGT-B09`

- [A gh failure is reported with gh's output](layers/infra.md#ghigt--githubissues--the-issue-lifecycle-on-the-repositorys-cards-when-the-project-works-on-github) `GHIGT-E01`

- [The key is stable across dates and distinct per gate](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-B01`

- [The file name is the date and the key, with no slash](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-B02`

- [The body names the kind, the target, the gate and the detail](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-B03`

- [A new issue is opened in todo](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-B04`

- [The same issue is not opened twice](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-B05`

- [An assumed debt is born in future and shows when it is due](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-B06`

- [A decision explains how to close it](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-B07`

- [Resolving moves a live issue to done, and only once](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-B08`

- [A new finding reopens the issue and keeps the old report](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-B09`

- [Issues are listed by owner, and no owner means the agent](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-B10`

- [Reassigning hands the issue over with its reason](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-B11`

- [A full check closes the violations it no longer reproduces](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-B12`

- [Issues are files unless GitHub is configured](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-B13`

- [A resolved issue is moved, not copied, and never resurrected](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-I01`

- [Reconciling spares decisions and the user's violations](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-X01`

- [An issue that cannot be written is reported](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-E01`

- [Reassigning a missing issue is reported](layers/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder) `ISLFS-E02`

- [A pack file gives its metadata and its obligations](layers/infra.md#obpcb--obligationpack--distributable-sets-of-obligations-from-a-norm-resolved-against-the-project) `OBPCB-B01`

- [A reference is a path or a name under packs](layers/infra.md#obpcb--obligationpack--distributable-sets-of-obligations-from-a-norm-resolved-against-the-project) `OBPCB-B02`

- [Placeholders are replaced by the project's values](layers/infra.md#obpcb--obligationpack--distributable-sets-of-obligations-from-a-norm-resolved-against-the-project) `OBPCB-B03`

- [Packs come back sorted by name](layers/infra.md#obpcb--obligationpack--distributable-sets-of-obligations-from-a-norm-resolved-against-the-project) `OBPCB-B04`

- [A pack of an undeclared jurisdiction is skipped with a warning](layers/infra.md#obpcb--obligationpack--distributable-sets-of-obligations-from-a-norm-resolved-against-the-project) `OBPCB-B05`

- [Global packs and projects without jurisdictions load everything](layers/infra.md#obpcb--obligationpack--distributable-sets-of-obligations-from-a-norm-resolved-against-the-project) `OBPCB-B06`

- [A pack without a name is refused](layers/infra.md#obpcb--obligationpack--distributable-sets-of-obligations-from-a-norm-resolved-against-the-project) `OBPCB-E01`

- [A pack without obligations is refused](layers/infra.md#obpcb--obligationpack--distributable-sets-of-obligations-from-a-norm-resolved-against-the-project) `OBPCB-E02`

- [An invalid or missing pack file is refused](layers/infra.md#obpcb--obligationpack--distributable-sets-of-obligations-from-a-norm-resolved-against-the-project) `OBPCB-E03`

- [Unresolved placeholders refuse the whole load](layers/infra.md#obpcb--obligationpack--distributable-sets-of-obligations-from-a-norm-resolved-against-the-project) `OBPCB-E04`

- [An enqueued task is born pending](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-B01`

- [The same target and step are not enqueued twice](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-B02`

- [Listing gives the live tasks sorted, and nothing without a queue](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-B03`

- [A task whose target was deleted is removed](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-B04`

- [Claiming takes a pending task and records the worker](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-B05`

- [A pending residue of a dead claim is not served and is cleaned](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-B06`

- [A done task moves to the history](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-B07`

- [Dropping deletes without history](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-B08`

- [Old and unstamped claims are returned to pending](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-B09`

- [Forced reclaiming returns even recent claims](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-B10`

- [Recently held claims are counted](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-B11`

- [The next step follows the kind that changed](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-B12`

- [The suggestions of the unit kinds are composable by the work command](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-B13`

- [A task whose target is an absolute path that exists is kept](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-B15`

- [The pending count counts pending and claimed tasks](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-B14`

- [Concurrent workers never claim the same task](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-I01`

- [A recent claim is not reclaimed](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-X01`

- [Marking an unknown task done is refused](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-E01`

- [Dropping an unknown task is refused](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-E02`

- [A corrupted task file is listed as triage, and can be claimed and dropped](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-E03`

- [Enqueuing an ID a live task already holds for another target is refused](layers/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it) `TSQUT-E04`

- [The lower convention derives the lower-case prefix](layers/infra.md#rcdlr--recodedialect--the-projects-own-surfaces-of-a-code-testid-prefixes-and-file-names) `RCDLR-B01`

- [Only quoted, hyphenated testID tokens are rewritten](layers/infra.md#rcdlr--recodedialect--the-projects-own-surfaces-of-a-code-testid-prefixes-and-file-names) `RCDLR-B02`

- [An empty or unchanged prefix rewrites nothing](layers/infra.md#rcdlr--recodedialect--the-projects-own-surfaces-of-a-code-testid-prefixes-and-file-names) `RCDLR-B03`

- [The occurrences of a prefix are counted](layers/infra.md#rcdlr--recodedialect--the-projects-own-surfaces-of-a-code-testid-prefixes-and-file-names) `RCDLR-B04`

- [The testID attributes are counted whatever their prefix](layers/infra.md#rcdlr--recodedialect--the-projects-own-surfaces-of-a-code-testid-prefixes-and-file-names) `RCDLR-B05`

- [A path matches the file patterns with the exact code](layers/infra.md#rcdlr--recodedialect--the-projects-own-surfaces-of-a-code-testid-prefixes-and-file-names) `RCDLR-B06`

- [The code is renamed in the file name only](layers/infra.md#rcdlr--recodedialect--the-projects-own-surfaces-of-a-code-testid-prefixes-and-file-names) `RCDLR-B07`

- [A longer code in a file name is not renamed](layers/infra.md#rcdlr--recodedialect--the-projects-own-surfaces-of-a-code-testid-prefixes-and-file-names) `RCDLR-X01`

- [A malformed source or target code is refused](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-B01`

- [Renaming a code to itself is refused](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-B02`

- [The plan holds every file with the code, sorted, with its new content](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-B03`

- [The declared testID prefix is rewritten and counted apart](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-B04`

- [Files named after the code are planned for renaming](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-B05`

- [A divergent testID prefix is warned about and never rewritten](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-B06`

- [A code that appears nowhere is refused](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-B07`

- [Applying writes each file with its mode and performs the renames](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-B08`

- [Moves go through git inside a repository and are plain renames outside](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-B09`

- [An untracked file inside a repository is moved by a plain rename](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-B10`

- [A target code another unit already owns is refused](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-B11`

- [The malformed-code refusal names the lengths the project declares](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-B12`

- [A git refusal is surfaced and never bypassed](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-X01`

- [Planning writes nothing](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-X02`

- [A file that cannot be written stops the apply, naming it](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-E05`

- [A batch plans many codes in one pass, cited only, renames the files their names carry, and lists the bare words](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-B13`

- [A binary file is renamed with its code and its bytes are never rewritten](layers/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project) `RCPLR-B14`

- [A code is well formed only in the project's lengths and alphabet](layers/infra.md#rcrwr--recoderewrite--renaming-an-identity-code-inside-a-text-on-every-surface-where-it-appears) `RCRWR-B01`

- [The header code is replaced in every header style](layers/infra.md#rcrwr--recoderewrite--renaming-an-identity-code-inside-a-text-on-every-surface-where-it-appears) `RCRWR-B02`

- [Only the old code changes in a ref list](layers/infra.md#rcrwr--recoderewrite--renaming-an-identity-code-inside-a-text-on-every-surface-where-it-appears) `RCRWR-B03`

- [Scenario codes keep their suffixes](layers/infra.md#rcrwr--recoderewrite--renaming-an-identity-code-inside-a-text-on-every-surface-where-it-appears) `RCRWR-B04`

- [A bare mention is replaced and its neighbour kept](layers/infra.md#rcrwr--recoderewrite--renaming-an-identity-code-inside-a-text-on-every-surface-where-it-appears) `RCRWR-B05`

- [The dry run classifies each occurrence](layers/infra.md#rcrwr--recoderewrite--renaming-an-identity-code-inside-a-text-on-every-surface-where-it-appears) `RCRWR-B06`

- [Each listed occurrence carries its line, counted from one](layers/infra.md#rcrwr--recoderewrite--renaming-an-identity-code-inside-a-text-on-every-surface-where-it-appears) `RCRWR-B07`

- [A text without the old code is left unchanged](layers/infra.md#rcrwr--recoderewrite--renaming-an-identity-code-inside-a-text-on-every-surface-where-it-appears) `RCRWR-I01`

- [Longer codes and neighbours are never touched](layers/infra.md#rcrwr--recoderewrite--renaming-an-identity-code-inside-a-text-on-every-surface-where-it-appears) `RCRWR-X01`

- [Outside the governed files only the rule and scenario codes are rewritten](layers/infra.md#rcrwr--recoderewrite--renaming-an-identity-code-inside-a-text-on-every-surface-where-it-appears) `RCRWR-B08`

- [A code is rewritten only where it is cited as a code, and the bare words left are listed](layers/infra.md#rcrwr--recoderewrite--renaming-an-identity-code-inside-a-text-on-every-surface-where-it-appears) `RCRWR-B09`

- [A code inside an identifier with an underscore is not a mention](layers/infra.md#rcrwr--recoderewrite--renaming-an-identity-code-inside-a-text-on-every-surface-where-it-appears) `RCRWR-B10`

- [The sh on PATH is the shell](layers/infra.md#psxsh--shell--the-posix-shell-that-runs-a-projects-commands) `PSXSH-B01`

- [On Windows, the shell beside git](layers/infra.md#psxsh--shell--the-posix-shell-that-runs-a-projects-commands) `PSXSH-B02`

- [A command is sh -c with its arguments](layers/infra.md#psxsh--shell--the-posix-shell-that-runs-a-projects-commands) `PSXSH-B03`

- [No shell is an environment error](layers/infra.md#psxsh--shell--the-posix-shell-that-runs-a-projects-commands) `PSXSH-E01`

- [A command does not carry the hook's repository](layers/infra.md#psxsh--shell--the-posix-shell-that-runs-a-projects-commands) `PSXSH-B04`

- [Words are upper-cased, and one-character and digit-only words are dropped](layers/infra.md#txsmt--textsimilarity--how-close-two-texts-that-should-be-equal-are-weighted-by-what-each-word-discriminates) `TXSMT-B01`

- [A word in every text weighs nothing and a rarer word weighs the log of its rarity](layers/infra.md#txsmt--textsimilarity--how-close-two-texts-that-should-be-equal-are-weighted-by-what-each-word-discriminates) `TXSMT-B02`

- [A single-text corpus falls back to the unweighted count](layers/infra.md#txsmt--textsimilarity--how-close-two-texts-that-should-be-equal-are-weighted-by-what-each-word-discriminates) `TXSMT-B03`

- [Texts differing only in case and punctuation are identical](layers/infra.md#txsmt--textsimilarity--how-close-two-texts-that-should-be-equal-are-weighted-by-what-each-word-discriminates) `TXSMT-B04`

- [Both rulers above the threshold make the pair similar](layers/infra.md#txsmt--textsimilarity--how-close-two-texts-that-should-be-equal-are-weighted-by-what-each-word-discriminates) `TXSMT-B05`

- [Rulers that disagree make the pair borderline](layers/infra.md#txsmt--textsimilarity--how-close-two-texts-that-should-be-equal-are-weighted-by-what-each-word-discriminates) `TXSMT-B06`

- [A shared rare word pulls a low-scoring pair to similar](layers/infra.md#txsmt--textsimilarity--how-close-two-texts-that-should-be-equal-are-weighted-by-what-each-word-discriminates) `TXSMT-B07`

- [Different subjects are divergent](layers/infra.md#txsmt--textsimilarity--how-close-two-texts-that-should-be-equal-are-weighted-by-what-each-word-discriminates) `TXSMT-B08`

- [The reported score is the larger of the two rulers](layers/infra.md#txsmt--textsimilarity--how-close-two-texts-that-should-be-equal-are-weighted-by-what-each-word-discriminates) `TXSMT-B09`

- [Texts that differ only by a number are not identical](layers/infra.md#txsmt--textsimilarity--how-close-two-texts-that-should-be-equal-are-weighted-by-what-each-word-discriminates) `TXSMT-B10`

- [Equal texts without any word are identical](layers/infra.md#txsmt--textsimilarity--how-close-two-texts-that-should-be-equal-are-weighted-by-what-each-word-discriminates) `TXSMT-B11`

- [The rulers fall back together, so one weightless side does not make a borderline](layers/infra.md#txsmt--textsimilarity--how-close-two-texts-that-should-be-equal-are-weighted-by-what-each-word-discriminates) `TXSMT-B12`

- [A verdict prints its English name](layers/infra.md#txsmt--textsimilarity--how-close-two-texts-that-should-be-equal-are-weighted-by-what-each-word-discriminates) `TXSMT-B13`

- [A new suggestion is born in pending with its context and reason](layers/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide) `SGSTS-B01`

- [A suggestion with a patch carries it in a diff block](layers/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide) `SGSTS-B02`

- [A suggestion without a patch says the fix needs a human decision](layers/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide) `SGSTS-B03`

- [Opening the same pending suggestion twice creates it once](layers/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide) `SGSTS-B04`

- [A rejected suggestion is never reopened](layers/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide) `SGSTS-B05`

- [Deciding moves the suggestion and records the reason and the decider](layers/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide) `SGSTS-B06`

- [An automatic decision is marked as the AI's](layers/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide) `SGSTS-B07`

- [Listing a state gives its sorted IDs](layers/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide) `SGSTS-B08`

- [The patch comes out clean for git apply](layers/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide) `SGSTS-B09`

- [A decided suggestion is no longer pending](layers/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide) `SGSTS-I01`

- [A suggestion without an ID is refused](layers/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide) `SGSTS-E01`

- [A decision to an unknown state is refused](layers/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide) `SGSTS-E02`

- [A decision without a reason is refused](layers/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide) `SGSTS-E03`

- [Deciding a suggestion that is not pending is refused](layers/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide) `SGSTS-E04`

- [Asking the patch of a suggestion without one is refused](layers/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide) `SGSTS-E05`

- [A pattern reads each call followed by a literal as a test](layers/infra.md#tstls--testlist--the-projects-tests-read-the-way-the-project-says-they-are-written) `TSTLS-B01`

- [Titles may be quoted three ways, with escapes](layers/infra.md#tstls--testlist--the-projects-tests-read-the-way-the-project-says-they-are-written) `TSTLS-B02`

- [A call without a literal title is left out](layers/infra.md#tstls--testlist--the-projects-tests-read-the-way-the-project-says-they-are-written) `TSTLS-B03`

- [The pattern's tests come in file order then position](layers/infra.md#tstls--testlist--the-projects-tests-read-the-way-the-project-says-they-are-written) `TSTLS-B04`

- [A script's contract output is the list of tests](layers/infra.md#tstls--testlist--the-projects-tests-read-the-way-the-project-says-they-are-written) `TSTLS-B05`

- [Output outside the contract is refused naming what is wrong](layers/infra.md#tstls--testlist--the-projects-tests-read-the-way-the-project-says-they-are-written) `TSTLS-B06`

- [A script's file paths are normalised to the map's form](layers/infra.md#tstls--testlist--the-projects-tests-read-the-way-the-project-says-they-are-written) `TSTLS-B07`

- [A source that declares nothing lists nothing](layers/infra.md#tstls--testlist--the-projects-tests-read-the-way-the-project-says-they-are-written) `TSTLS-B08`

- [A source with both a pattern and a script is refused](layers/infra.md#tstls--testlist--the-projects-tests-read-the-way-the-project-says-they-are-written) `TSTLS-E01`

- [A pattern that does not compile is refused](layers/infra.md#tstls--testlist--the-projects-tests-read-the-way-the-project-says-they-are-written) `TSTLS-E02`

- [A failing script is an error naming why](layers/infra.md#tstls--testlist--the-projects-tests-read-the-way-the-project-says-they-are-written) `TSTLS-E03`

- [The script may say where a test ends](layers/infra.md#tstls--testlist--the-projects-tests-read-the-way-the-project-says-they-are-written) `TSTLS-B09`

- [A pattern scans the files through the reader given](layers/infra.md#tstls--testlist--the-projects-tests-read-the-way-the-project-says-they-are-written) `TSTLS-B10`

- [A rule code of each canonical letter is recognized in a test name](layers/infra.md#rcgrl--rulecodegrammar--the-grammar-that-recognizes-a-scenario-code-in-a-tests-name-in-the-projects-vocabulary) `RCGRL-B01`

- [A rule code with a lowercase slug is recognized with the slug](layers/infra.md#rcgrl--rulecodegrammar--the-grammar-that-recognizes-a-scenario-code-in-a-tests-name-in-the-projects-vocabulary) `RCGRL-B02`

- [Design-system and visual-regression codes are recognized](layers/infra.md#rcgrl--rulecodegrammar--the-grammar-that-recognizes-a-scenario-code-in-a-tests-name-in-the-projects-vocabulary) `RCGRL-B03`

- [Declared rule letters replace the vocabulary](layers/infra.md#rcgrl--rulecodegrammar--the-grammar-that-recognizes-a-scenario-code-in-a-tests-name-in-the-projects-vocabulary) `RCGRL-B04`

- [A declared code length replaces the accepted identity length](layers/infra.md#rcgrl--rulecodegrammar--the-grammar-that-recognizes-a-scenario-code-in-a-tests-name-in-the-projects-vocabulary) `RCGRL-B05`

- [An empty declaration keeps the vocabulary in place](layers/infra.md#rcgrl--rulecodegrammar--the-grammar-that-recognizes-a-scenario-code-in-a-tests-name-in-the-projects-vocabulary) `RCGRL-B06`

- [The default rule letters equal the configuration's default letters](layers/infra.md#rcgrl--rulecodegrammar--the-grammar-that-recognizes-a-scenario-code-in-a-tests-name-in-the-projects-vocabulary) `RCGRL-I01`

- [An identity too long or glued to a longer word is not a code](layers/infra.md#rcgrl--rulecodegrammar--the-grammar-that-recognizes-a-scenario-code-in-a-tests-name-in-the-projects-vocabulary) `RCGRL-X01`

- [A scenario code keeps its variant](layers/infra.md#rcgrl--rulecodegrammar--the-grammar-that-recognizes-a-scenario-code-in-a-tests-name-in-the-projects-vocabulary) `RCGRL-B07`

- [A rule is proven only when each of its scenarios is](layers/infra.md#rcgrl--rulecodegrammar--the-grammar-that-recognizes-a-scenario-code-in-a-tests-name-in-the-projects-vocabulary) `RCGRL-B08`

- [Added lines are recorded under the file of the new-file header](layers/infra.md#dcldf--diffchangedlines--which-lines-of-which-files-a-change-added-read-from-a-unified-diff) `DCLDF-B01`

- [A hunk header sets the starting line of the new side](layers/infra.md#dcldf--diffchangedlines--which-lines-of-which-files-a-change-added-read-from-a-unified-diff) `DCLDF-B02`

- [Path prefixes and a tab-separated timestamp are stripped](layers/infra.md#dcldf--diffchangedlines--which-lines-of-which-files-a-change-added-read-from-a-unified-diff) `DCLDF-B03`

- [A deleted file records nothing](layers/infra.md#dcldf--diffchangedlines--which-lines-of-which-files-a-change-added-read-from-a-unified-diff) `DCLDF-B04`

- [Context lines advance the numbering without being recorded](layers/infra.md#dcldf--diffchangedlines--which-lines-of-which-files-a-change-added-read-from-a-unified-diff) `DCLDF-B05`

- [Git compares the working copy with the current commit or with a reference](layers/infra.md#dcldf--diffchangedlines--which-lines-of-which-files-a-change-added-read-from-a-unified-diff) `DCLDF-B06`

- [An added line that starts with two plus signs is a line, not a header](layers/infra.md#dcldf--diffchangedlines--which-lines-of-which-files-a-change-added-read-from-a-unified-diff) `DCLDF-B07`

- [A very long line in the diff is read like any other](layers/infra.md#dcldf--diffchangedlines--which-lines-of-which-files-a-change-added-read-from-a-unified-diff) `DCLDF-B08`

- [Removals do not shift the new-side numbering](layers/infra.md#dcldf--diffchangedlines--which-lines-of-which-files-a-change-added-read-from-a-unified-diff) `DCLDF-I01`

- [A diff file is read without git](layers/infra.md#dcldf--diffchangedlines--which-lines-of-which-files-a-change-added-read-from-a-unified-diff) `DCLDF-X01`

- [A directory that is not a repository fails the git diff](layers/infra.md#dcldf--diffchangedlines--which-lines-of-which-files-a-change-added-read-from-a-unified-diff) `DCLDF-E01`

- [A missing diff file surfaces the read error](layers/infra.md#dcldf--diffchangedlines--which-lines-of-which-files-a-change-added-read-from-a-unified-diff) `DCLDF-E02`

- [A report with a suites root is read case by case](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-B01`

- [A report whose root is a single suite is read](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-B02`

- [Nested suites are flattened](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-B03`

- [A case without a file takes its suite's file](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-B04`

- [A failure or an error marks the case failed and a skip marks it skipped](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-B05`

- [Only codes of passing cases are proven](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-B06`

- [Every case's codes are seen whatever the outcome](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-B07`

- [Every code in a case name is extracted in the declared vocabulary](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-B08`

- [A file that is not a JUnit report yields no case](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-B09`

- [Every proven code is a seen code](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-I01`

- [A code in the suite or class name proves nothing](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-X01`

- [An unreadable report returns the read error](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-E01`

- [Each case carries its run time](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-B10`

- [The proven and seen codes carry the variant](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-B11`

- [A scenario is proven only by its own passing case](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-I02`

- [A green capture of a state proves the state](layers/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report) `JUIJN-B12`

- [Each record becomes one file's coverage in report order](layers/infra.md#lcinl--lcovingest--line-coverage-per-file-and-the-uncovered-lines-of-a-change-read-from-an-lcov-report) `LCINL-B01`

- [A line with hits is covered and a line without is not](layers/infra.md#lcinl--lcovingest--line-coverage-per-file-and-the-uncovered-lines-of-a-change-read-from-an-lcov-report) `LCINL-B02`

- [Stated totals win over the line entries](layers/infra.md#lcinl--lcovingest--line-coverage-per-file-and-the-uncovered-lines-of-a-change-read-from-an-lcov-report) `LCINL-B03`

- [The uncovered lines of a change are instrumented and not covered](layers/infra.md#lcinl--lcovingest--line-coverage-per-file-and-the-uncovered-lines-of-a-change-read-from-an-lcov-report) `LCINL-B04`

- [The instrumented count ignores changed lines the report did not instrument](layers/infra.md#lcinl--lcovingest--line-coverage-per-file-and-the-uncovered-lines-of-a-change-read-from-an-lcov-report) `LCINL-B05`

- [The percentage is covered over total, and zero for a file without lines](layers/infra.md#lcinl--lcovingest--line-coverage-per-file-and-the-uncovered-lines-of-a-change-read-from-an-lcov-report) `LCINL-B06`

- [A record without its end line is closed by the next record or the end of the report](layers/infra.md#lcinl--lcovingest--line-coverage-per-file-and-the-uncovered-lines-of-a-change-read-from-an-lcov-report) `LCINL-B07`

- [A line entry with a checksum field keeps its hit count](layers/infra.md#lcinl--lcovingest--line-coverage-per-file-and-the-uncovered-lines-of-a-change-read-from-an-lcov-report) `LCINL-B08`

- [Entries before the first source-file line belong to no file](layers/infra.md#lcinl--lcovingest--line-coverage-per-file-and-the-uncovered-lines-of-a-change-read-from-an-lcov-report) `LCINL-B09`

- [Uncovered changed lines never exceed the instrumented changed lines](layers/infra.md#lcinl--lcovingest--line-coverage-per-file-and-the-uncovered-lines-of-a-change-read-from-an-lcov-report) `LCINL-I01`

- [Branch and function entries do not change the line counts](layers/infra.md#lcinl--lcovingest--line-coverage-per-file-and-the-uncovered-lines-of-a-change-read-from-an-lcov-report) `LCINL-X01`

- [A missing report returns the error](layers/infra.md#lcinl--lcovingest--line-coverage-per-file-and-the-uncovered-lines-of-a-change-read-from-an-lcov-report) `LCINL-E01`

- [Branch entries record each branch, taken or not](layers/infra.md#lcinl--lcovingest--line-coverage-per-file-and-the-uncovered-lines-of-a-change-read-from-an-lcov-report) `LCINL-B10`

- [The canonical format is read under each of its names](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-B01`

- [Killed and timed-out mutants count as killed](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-B02`

- [A survivor counts as survived with its line](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-B03`

- [Uncovered mutants are counted apart and stay out of the score](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-B04`

- [Ignored mutants are counted apart and stay out of the score](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-B05`

- [A mutant that failed to compile stays out of the count and the score](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-B06`

- [A file where no mutant ran has no score](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-B07`

- [The thresholds are read from the report](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-B08`

- [File paths are normalized to the map's form](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-B09`

- [A file where nothing ran has no score and keeps the ignored count](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-I01`

- [Thresholds absent from the report stay zero](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-X01`

- [An unknown format is refused naming the accepted ones](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-E01`

- [A report that is not the canonical format is refused](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-E02`

- [A report with no file and no schema version is refused](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-E03`

- [A missing report returns the read error](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-E04`

- [A report in the format with no file had nothing to mutate](layers/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report) `MTINM-B10`

- [The listed files are read into one result per file](layers/infra.md#gring--gremlinsingest--the-mutation-score-per-file-read-from-a-gremlins-report) `GRING-B01`

- [Killed and timed-out mutants count as killed](layers/infra.md#gring--gremlinsingest--the-mutation-score-per-file-read-from-a-gremlins-report) `GRING-B02`

- [A mutant that lived counts as survived with its line](layers/infra.md#gring--gremlinsingest--the-mutation-score-per-file-read-from-a-gremlins-report) `GRING-B03`

- [Not viable, runnable, skipped and unknown statuses stay out of the score](layers/infra.md#gring--gremlinsingest--the-mutation-score-per-file-read-from-a-gremlins-report) `GRING-B04`

- [The status is matched ignoring case and spaces](layers/infra.md#gring--gremlinsingest--the-mutation-score-per-file-read-from-a-gremlins-report) `GRING-B05`

- [The thresholds stay zero](layers/infra.md#gring--gremlinsingest--the-mutation-score-per-file-read-from-a-gremlins-report) `GRING-B06`

- [File paths are normalized as in the canonical reading](layers/infra.md#gring--gremlinsingest--the-mutation-score-per-file-read-from-a-gremlins-report) `GRING-B07`

- [The score is killed over killed plus survived](layers/infra.md#gring--gremlinsingest--the-mutation-score-per-file-read-from-a-gremlins-report) `GRING-I01`

- [A mutant no test covered is counted apart and does not enter the score](layers/infra.md#gring--gremlinsingest--the-mutation-score-per-file-read-from-a-gremlins-report) `GRING-B08`

- [A file where no mutant ran scores 100](layers/infra.md#gring--gremlinsingest--the-mutation-score-per-file-read-from-a-gremlins-report) `GRING-B09`

- [The report's own efficacy figure is not used](layers/infra.md#gring--gremlinsingest--the-mutation-score-per-file-read-from-a-gremlins-report) `GRING-X01`

- [A canonical-format report under the gremlins format is refused naming format](layers/infra.md#gring--gremlinsingest--the-mutation-score-per-file-read-from-a-gremlins-report) `GRING-E01`

- [A report without files is refused](layers/infra.md#gring--gremlinsingest--the-mutation-score-per-file-read-from-a-gremlins-report) `GRING-E02`

## mapa

- [Each scanned file becomes a node with its layer's tags and regime](layers/mapa.md#grblg-b01--each-scanned-file-becomes-a-node-with-its-layers-tags-and-regime) `GRBLG-B01`

- [The declared identity wins over cited codes and over the anchor](layers/mapa.md#grblg-b02--the-declared-identity-wins-over-cited-codes-and-over-the-anchor) `GRBLG-B02`

- [A derived file takes its sibling anchor's identity](layers/mapa.md#grblg-b03--a-derived-file-takes-its-sibling-anchors-identity) `GRBLG-B03`

- [With no header and no anchor, the first code's root is the identity](layers/mapa.md#grblg-b04--with-no-header-and-no-anchor-the-first-codes-root-is-the-identity) `GRBLG-B04`

- [A vendored file has no local identity](layers/mapa.md#grblg-b05--a-vendored-file-has-no-local-identity) `GRBLG-B05`

- [Only a header marks the identity as declared](layers/mapa.md#grblg-b06--only-a-header-marks-the-identity-as-declared) `GRBLG-B06`

- [The pieces of one unit are linked](layers/mapa.md#grblg-b07--the-pieces-of-one-unit-are-linked) `GRBLG-B07`

- [With the code as anchor, the relations still go down from the spec](layers/mapa.md#grblg-b08--with-the-code-as-anchor-the-relations-still-go-down-from-the-spec) `GRBLG-B08`

- [Without a feature, the spec is tested by the test](layers/mapa.md#grblg-b09--without-a-feature-the-spec-is-tested-by-the-test) `GRBLG-B09`

- [A layer override reaches the spec through the layer it declares](layers/mapa.md#grblg-b10--a-layer-override-reaches-the-spec-through-the-layer-it-declares) `GRBLG-B10`

- [A code override replaces the templates of the kinds it declares, and the others fall back to the default](layers/mapa.md#grblg-b11--a-code-override-replaces-the-templates-of-the-kinds-it-declares-and-the-others-fall-back-to-the-default) `GRBLG-B11`

- [A bracketed directory is literal and a template wildcard expands](layers/mapa.md#grblg-b12--a-bracketed-directory-is-literal-and-a-template-wildcard-expands) `GRBLG-B12`

- [A shared own code links spec and test across directories only](layers/mapa.md#grblg-b13--a-shared-own-code-links-spec-and-test-across-directories-only) `GRBLG-B13`

- [A guide governs the layers of its tag only, and never itself](layers/mapa.md#grblg-b14--a-guide-governs-the-layers-of-its-tag-only-and-never-itself) `GRBLG-B14`

- [A dependency row becomes a relation to an existing file](layers/mapa.md#grblg-b15--a-dependency-row-becomes-a-relation-to-an-existing-file) `GRBLG-B15`

- [A seed path is exact and a bare name must be unique](layers/mapa.md#grblg-b16--a-seed-path-is-exact-and-a-bare-name-must-be-unique) `GRBLG-B16`

- [A need links only an existing plan](layers/mapa.md#grblg-b17--a-need-links-only-an-existing-plan) `GRBLG-B17`

- [Realized doctrine and flag scenarios resolve by unit code](layers/mapa.md#grblg-b18--realized-doctrine-and-flag-scenarios-resolve-by-unit-code) `GRBLG-B18`

- [Nodes and relations are sorted](layers/mapa.md#grblg-b19--nodes-and-relations-are-sorted) `GRBLG-B19`

- [A rebuild keeps the stamps and judgments of surviving relations](layers/mapa.md#grblg-b20--a-rebuild-keeps-the-stamps-and-judgments-of-surviving-relations) `GRBLG-B20`

- [A rebuild keeps a signal only for an unchanged file](layers/mapa.md#grblg-b21--a-rebuild-keeps-a-signal-only-for-an-unchanged-file) `GRBLG-B21`

- [The build does not depend on the order of the files](layers/mapa.md#grblg-i01--the-build-does-not-depend-on-the-order-of-the-files) `GRBLG-I01`

- [With no dates given, nodes carry no date](layers/mapa.md#grblg-x01--with-no-dates-given-nodes-carry-no-date) `GRBLG-X01`

- [No relation points to a file that was not scanned](layers/mapa.md#grblg-x02--no-relation-points-to-a-file-that-was-not-scanned) `GRBLG-X02`

- [A support file becomes a node marked as support](layers/mapa.md#grblg-b22--a-support-file-becomes-a-node-marked-as-support) `GRBLG-B22`

- [The edges are in a total order](layers/mapa.md#grblg-i02--the-edges-are-in-a-total-order) `GRBLG-I02`

- [Signals are filled from another map at the same revision](layers/mapa.md#grblg-b23--signals-are-filled-from-another-map-at-the-same-revision) `GRBLG-B23`

- [A file's own code is its file code, and a ref keeps its unit](layers/mapa.md#grblg-b24--a-files-own-code-is-its-file-code-and-a-ref-keeps-its-unit) `GRBLG-B24`

- [A file with a code of its own and a ref links across directories through the unit it refs](layers/mapa.md#grblg-b25--a-file-with-a-code-of-its-own-and-a-ref-links-across-directories-through-the-unit-it-refs) `GRBLG-B25`

- [The @dep and @navigates flags become edges to the file whose own code they name](layers/mapa.md#grblg-b26--the-dep-and-navigates-flags-become-edges-to-the-file-whose-own-code-they-name) `GRBLG-B26`

- [A test that was never ingested has no verdict](layers/mapa.md#evfra-b01--a-test-that-was-never-ingested-has-no-verdict) `EVFRA-B01`

- [The test's own change expires its evidence](layers/mapa.md#evfra-b02--the-tests-own-change-expires-its-evidence) `EVFRA-B02`

- [A composed script that changed is named as the culprit](layers/mapa.md#evfra-b03--a-composed-script-that-changed-is-named-as-the-culprit) `EVFRA-B03`

- [Nothing moved means no verdict](layers/mapa.md#evfra-b04--nothing-moved-means-no-verdict) `EVFRA-B04`

- [A signal with no recorded closure is judged by its own file only](layers/mapa.md#evfra-b05--a-signal-with-no-recorded-closure-is-judged-by-its-own-file-only) `EVFRA-B05`

- [The closure descends transitively with current revisions](layers/mapa.md#evfra-b06--the-closure-descends-transitively-with-current-revisions) `EVFRA-B06`

- [A no-propagation node is in the closure but not walked through](layers/mapa.md#evfra-b07--a-no-propagation-node-is-in-the-closure-but-not-walked-through) `EVFRA-B07`

- [A recorded node that left the graph is not a culprit](layers/mapa.md#evfra-b08--a-recorded-node-that-left-the-graph-is-not-a-culprit) `EVFRA-B08`

- [The closure never contains the test itself](layers/mapa.md#evfra-i01--the-closure-never-contains-the-test-itself) `EVFRA-I01`

- [The closure never climbs to the spec above the test](layers/mapa.md#evfra-x01--the-closure-never-climbs-to-the-spec-above-the-test) `EVFRA-X01`

- [A capture's target enters the closure and is not descended](layers/mapa.md#evfra-b09--a-captures-target-enters-the-closure-and-is-not-descended) `EVFRA-B09`

- [A component whose capture diverged stales the captures of who uses it](layers/mapa.md#evfra-b10--a-component-whose-capture-diverged-stales-the-captures-of-who-uses-it) `EVFRA-B10`

- [A flow that asserts a navigation goes stale when its Out row changes or goes, and a flow that passes through does not](layers/mapa.md#evfra-b11--a-flow-that-asserts-a-navigation-goes-stale-when-its-out-row-changes-or-goes-and-a-flow-that-passes-through-does-not) `EVFRA-B11`

- [The written format and the oldest readable format are both accepted](layers/mapa.md#mpfrm-b01--the-written-format-and-the-oldest-readable-format-are-both-accepted) `MPFRM-B01`

- [A map from a newer binary is refused with the upgrade message](layers/mapa.md#mpfrm-b02--a-map-from-a-newer-binary-is-refused-with-the-upgrade-message) `MPFRM-B02`

- [A map older than the readable range asks for migration](layers/mapa.md#mpfrm-b03--a-map-older-than-the-readable-range-asks-for-migration) `MPFRM-B03`

- [A map with no version is format 1 and asks for migration](layers/mapa.md#mpfrm-b04--a-map-with-no-version-is-format-1-and-asks-for-migration) `MPFRM-B04`

- [The newer-map refusal never names the migration command](layers/mapa.md#mpfrm-b05--the-newer-map-refusal-never-names-the-migration-command) `MPFRM-B05`

- [Exactly format 7 is readable](layers/mapa.md#mpfrm-i01--exactly-format-7-is-readable) `MPFRM-I01`

- [Format 1 is migrated, not read](layers/mapa.md#mpfrm-x01--format-1-is-migrated-not-read) `MPFRM-X01`

- [Changing a spec propagates down to its code, feature and test](layers/mapa.md#imanm-b01--changing-a-spec-propagates-down-to-its-code-feature-and-test) `IMANM-B01`

- [A no-propagation child is reached but the wave stops there](layers/mapa.md#imanm-b02--a-no-propagation-child-is-reached-but-the-wave-stops-there) `IMANM-B02`

- [Changing code is validated upward against its spec and the spec's guide](layers/mapa.md#imanm-b03--changing-code-is-validated-upward-against-its-spec-and-the-specs-guide) `IMANM-B03`

- [Climbing to a shared guide does not reach the sibling unit](layers/mapa.md#imanm-b04--climbing-to-a-shared-guide-does-not-reach-the-sibling-unit) `IMANM-B04`

- [A node with no edges has no impact](layers/mapa.md#imanm-b05--a-node-with-no-edges-has-no-impact) `IMANM-B05`

- [The changed node is never in its own lists, even in a cycle](layers/mapa.md#imanm-i01--the-changed-node-is-never-in-its-own-lists-even-in-a-cycle) `IMANM-I01`

- [The analysis leaves the graph as it was](layers/mapa.md#imanm-x01--the-analysis-leaves-the-graph-as-it-was) `IMANM-X01`

- [A new file enters with the same nodes and edges the full build gives](layers/mapa.md#grinc-b01--a-new-file-enters-with-the-same-nodes-and-edges-the-full-build-gives) `GRINC-B01`

- [The derivation links inside the unit are replaced, removals included](layers/mapa.md#grinc-b02--the-derivation-links-inside-the-unit-are-replaced-removals-included) `GRINC-B02`

- [A new anchor gives its declared code to its derived siblings](layers/mapa.md#grinc-b03--a-new-anchor-gives-its-declared-code-to-its-derived-siblings) `GRINC-B03`

- [What the new file declares reaches the existing files](layers/mapa.md#grinc-b04--what-the-new-file-declares-reaches-the-existing-files) `GRINC-B04`

- [A known, ignored or unclassified file adds nothing](layers/mapa.md#grinc-b05--a-known-ignored-or-unclassified-file-adds-nothing) `GRINC-B05`

- [On disk the addition runs under the lock, and not without a map](layers/mapa.md#grinc-b06--on-disk-the-addition-runs-under-the-lock-and-not-without-a-map) `GRINC-B06`

- [The anchor that may own a new file is read, for the override its header chooses](layers/mapa.md#grinc-b07--the-anchor-that-may-own-a-new-file-is-read-for-the-override-its-header-chooses) `GRINC-B07`

- [A relation an existing file declares toward the new one waits for the full build](layers/mapa.md#grinc-x01--a-relation-an-existing-file-declares-toward-the-new-one-waits-for-the-full-build) `GRINC-X01`

- [A reader that fails changes nothing](layers/mapa.md#grinc-e01--a-reader-that-fails-changes-nothing) `GRINC-E01`

- [A test node sums its layers and records its revisions](layers/mapa.md#sgina-b01--a-test-node-sums-its-layers-and-records-its-revisions) `SGINA-B01`

- [Re-ingesting a layer replaces only that layer](layers/mapa.md#sgina-b02--re-ingesting-a-layer-replaces-only-that-layer) `SGINA-B02`

- [A full run records the proven rules and erases the lost ones](layers/mapa.md#sgina-b03--a-full-run-records-the-proven-rules-and-erases-the-lost-ones) `SGINA-B03`

- [One suite never erases another suite's proof](layers/mapa.md#sgina-b04--one-suite-never-erases-another-suites-proof) `SGINA-B04`

- [A suite that proves nothing leaves the union](layers/mapa.md#sgina-b05--a-suite-that-proves-nothing-leaves-the-union) `SGINA-B05`

- [The union is as fresh as its oldest contributor](layers/mapa.md#sgina-b06--the-union-is-as-fresh-as-its-oldest-contributor) `SGINA-B06`

- [A partial run changes only what it saw](layers/mapa.md#sgina-b07--a-partial-run-changes-only-what-it-saw) `SGINA-B07`

- [External suites are dropped and the union recomputed](layers/mapa.md#sgina-b08--external-suites-are-dropped-and-the-union-recomputed) `SGINA-B08`

- [Coverage with no suite records the percentage and the baseline](layers/mapa.md#sgina-b09--coverage-with-no-suite-records-the-percentage-and-the-baseline) `SGINA-B09`

- [Suite coverage is the union of lines](layers/mapa.md#sgina-b10--suite-coverage-is-the-union-of-lines) `SGINA-B10`

- [Totals-only suites fall back to the best suite](layers/mapa.md#sgina-b11--totals-only-suites-fall-back-to-the-best-suite) `SGINA-B11`

- [A suite measured before the edit stays out of the union](layers/mapa.md#sgina-b12--a-suite-measured-before-the-edit-stays-out-of-the-union) `SGINA-B12`

- [A report older than the file is kept without a revision](layers/mapa.md#sgina-b13--a-report-older-than-the-file-is-kept-without-a-revision) `SGINA-B13`

- [Line ranges round-trip](layers/mapa.md#sgina-b14--line-ranges-round-trip) `SGINA-B14`

- [Mutation totals and the per-scope measurement](layers/mapa.md#sgina-b15--mutation-totals-and-the-per-scope-measurement) `SGINA-B15`

- [A signal goes stale when its file moves](layers/mapa.md#sgina-b16--a-signal-goes-stale-when-its-file-moves) `SGINA-B16`

- [Paths match at a path boundary in either direction](layers/mapa.md#sgina-b17--paths-match-at-a-path-boundary-in-either-direction) `SGINA-B17`

- [Report paths resolve by the report's folder, or stay unowned](layers/mapa.md#sgina-b18--report-paths-resolve-by-the-reports-folder-or-stay-unowned) `SGINA-B18`

- [A report with no instrumented line leaves no percentage behind](layers/mapa.md#sgina-b19--a-report-with-no-instrumented-line-leaves-no-percentage-behind) `SGINA-B19`

- [Among several report paths matching a node, the exact one, then the closest in length, is always chosen](layers/mapa.md#sgina-b20--among-several-report-paths-matching-a-node-the-exact-one-then-the-closest-in-length-is-always-chosen) `SGINA-B20`

- [The proven rules are the sorted union of the suites](layers/mapa.md#sgina-i01--the-proven-rules-are-the-sorted-union-of-the-suites) `SGINA-I01`

- [Each measurement lands only on its kind of node](layers/mapa.md#sgina-x01--each-measurement-lands-only-on-its-kind-of-node) `SGINA-X01`

- [A mutation ingestion records the timed-out apart](layers/mapa.md#sgina-b21--a-mutation-ingestion-records-the-timed-out-apart) `SGINA-B21`

- [An execution ingestion records each test file's run time under its suite](layers/mapa.md#sgina-b22--an-execution-ingestion-records-each-test-files-run-time-under-its-suite) `SGINA-B22`

- [The mutation result keeps its own rev](layers/mapa.md#sgina-b25--the-mutation-result-keeps-its-own-rev) `SGINA-B25`

- [The tree's revs replace the map's](layers/mapa.md#sgina-b24--the-trees-revs-replace-the-maps) `SGINA-B24`

- [A run time Anchors measured is recorded on the node](layers/mapa.md#sgina-b23--a-run-time-anchors-measured-is-recorded-on-the-node) `SGINA-B23`

- [The branches are the union of the fresh suites](layers/mapa.md#sgina-b26--the-branches-are-the-union-of-the-fresh-suites) `SGINA-B26`

- [The lines of mutants no test ran are recorded](layers/mapa.md#sgina-b27--the-lines-of-mutants-no-test-ran-are-recorded) `SGINA-B27`

- [A proven variant is kept on the node that declares its rule](layers/mapa.md#sgina-b28--a-proven-variant-is-kept-on-the-node-that-declares-its-rule) `SGINA-B28`

- [A file the coverage report lists, or leaves out, says so](layers/mapa.md#sgina-b29--a-file-the-coverage-report-lists-or-leaves-out-says-so) `SGINA-B29`

- [A suite that ran none of a file's lines leaves its coverage](layers/mapa.md#sgina-b30--a-suite-that-ran-none-of-a-files-lines-leaves-its-coverage) `SGINA-B30`

- [The suites another measurement holds and this signal does not are adopted, and the ones it has stay](layers/mapa.md#sgina-b24--the-suites-another-measurement-holds-and-this-signal-does-not-are-adopted-and-the-ones-it-has-stay) `SGINA-B24`

- [The lock is a file beside the map with its owner](layers/mapa.md#mplck-b01--the-lock-is-a-file-beside-the-map-with-its-owner) `MPLCK-B01`

- [A second writer waits for the first](layers/mapa.md#mplck-b02--a-second-writer-waits-for-the-first) `MPLCK-B02`

- [An abandoned lock is taken over, a live one is not](layers/mapa.md#mplck-b03--an-abandoned-lock-is-taken-over-a-live-one-is-not) `MPLCK-B03`

- [Parallel processes each changing their part all reach the map](layers/mapa.md#mplck-b04--parallel-processes-each-changing-their-part-all-reach-the-map) `MPLCK-B04`

- [The lock is released when the function returns](layers/mapa.md#mplck-b05--the-lock-is-released-when-the-function-returns) `MPLCK-B05`

- [A lock held past the timeout](layers/mapa.md#mplck-e01--a-lock-held-past-the-timeout) `MPLCK-E01`

- [No map, or a refused change, writes nothing](layers/mapa.md#mplck-e02--no-map-or-a-refused-change-writes-nothing) `MPLCK-E02`

- [A lock that cannot be created names the file](layers/mapa.md#mplck-e03--a-lock-that-cannot-be-created-names-the-file) `MPLCK-E03`

- [A relation never stamped is stale](layers/mapa.md#grmdg-b01--a-relation-never-stamped-is-stale) `GRMDG-B01`

- [A stamp is fresh while both ends keep their revisions](layers/mapa.md#grmdg-b02--a-stamp-is-fresh-while-both-ends-keep-their-revisions) `GRMDG-B02`

- [A mutation scope is stale only against a recorded, different revision](layers/mapa.md#grmdg-b03--a-mutation-scope-is-stale-only-against-a-recorded-different-revision) `GRMDG-B03`

- [The map file's keys are fixed and English](layers/mapa.md#grmdg-i01--the-map-files-keys-are-fixed-and-english) `GRMDG-I01`

- [Empty optional fields are left out](layers/mapa.md#grmdg-i02--empty-optional-fields-are-left-out) `GRMDG-I02`

- [The node kinds and relation types are the fixed vocabulary](layers/mapa.md#grmdg-i03--the-node-kinds-and-relation-types-are-the-fixed-vocabulary) `GRMDG-I03`

- [The stamp's date and verdict play no part in staleness](layers/mapa.md#grmdg-x01--the-stamps-date-and-verdict-play-no-part-in-staleness) `GRMDG-X01`

- [The navigation between screens is an edge type of its own](layers/mapa.md#grmdg-b04--the-navigation-between-screens-is-an-edge-type-of-its-own) `GRMDG-B04`

- [A guide governs only the targets of its governs edges, sorted](layers/mapa.md#grqrg-b01--a-guide-governs-only-the-targets-of-its-governs-edges-sorted) `GRQRG-B01`

- [The governance summary counts governs edges only](layers/mapa.md#grqrg-b02--the-governance-summary-counts-governs-edges-only) `GRQRG-B02`

- [The neighbourhood is one step each way, sorted](layers/mapa.md#grqrg-b03--the-neighbourhood-is-one-step-each-way-sorted) `GRQRG-B03`

- [Orphans are the nodes with no edge, sorted](layers/mapa.md#grqrg-b04--orphans-are-the-nodes-with-no-edge-sorted) `GRQRG-B04`

- [Statistics count nodes and edges by kind and type](layers/mapa.md#grqrg-b05--statistics-count-nodes-and-edges-by-kind-and-type) `GRQRG-B05`

- [Parents come before their children](layers/mapa.md#grqrg-b06--parents-come-before-their-children) `GRQRG-B06`

- [A depends-on edge imposes no order](layers/mapa.md#grqrg-b07--a-depends-on-edge-imposes-no-order) `GRQRG-B07`

- [Cycle members are appended at the end](layers/mapa.md#grqrg-b08--cycle-members-are-appended-at-the-end) `GRQRG-B08`

- [The order does not depend on how nodes are stored](layers/mapa.md#grqrg-i01--the-order-does-not-depend-on-how-nodes-are-stored) `GRQRG-I01`

- [Queries leave the graph as it was](layers/mapa.md#grqrg-x01--queries-leave-the-graph-as-it-was) `GRQRG-X01`

- [A node's unit codes come from its identity edges](layers/mapa.md#grqrg-b09--a-nodes-unit-codes-come-from-its-identity-edges) `GRQRG-B09`

- [A review holds only at the revision it looked at](layers/mapa.md#mprvm-b01--a-review-holds-only-at-the-revision-it-looked-at) `MPRVM-B01`

- [Recording a review replaces the same gate's and keeps the others](layers/mapa.md#mprvm-b02--recording-a-review-replaces-the-same-gates-and-keeps-the-others) `MPRVM-B02`

- [A rebuild keeps the reviews whatever the revision](layers/mapa.md#mprvm-b03--a-rebuild-keeps-the-reviews-whatever-the-revision) `MPRVM-B03`

- [A node holds at most one review per gate](layers/mapa.md#mprvm-i01--a-node-holds-at-most-one-review-per-gate) `MPRVM-I01`

- [Who reviewed is kept as given](layers/mapa.md#mprvm-x01--who-reviewed-is-kept-as-given) `MPRVM-X01`

- [Only relations with both ends confronted are stamped](layers/mapa.md#edstd-b01--only-relations-with-both-ends-confronted-are-stamped) `EDSTD-B01`

- [The verdict is issue when an end failed and ok when both passed](layers/mapa.md#edstd-b02--the-verdict-is-issue-when-an-end-failed-and-ok-when-both-passed) `EDSTD-B02`

- [A waived stamp is kept as it was](layers/mapa.md#edstd-b03--a-waived-stamp-is-kept-as-it-was) `EDSTD-B03`

- [A stamp is fresh until an end moves](layers/mapa.md#edstd-b04--a-stamp-is-fresh-until-an-end-moves) `EDSTD-B04`

- [The date moves only when revisions or verdict change](layers/mapa.md#edstd-b05--the-date-moves-only-when-revisions-or-verdict-change) `EDSTD-B05`

- [An undated stamp takes today's date](layers/mapa.md#edstd-b06--an-undated-stamp-takes-todays-date) `EDSTD-B06`

- [One relation is stamped by its ends, a missing one answers false](layers/mapa.md#edstd-b07--one-relation-is-stamped-by-its-ends-a-missing-one-answers-false) `EDSTD-B07`

- [Stamping a relation by gate records the gate and a judgment](layers/mapa.md#edstd-b08--stamping-a-relation-by-gate-records-the-gate-and-a-judgment) `EDSTD-B08`

- [Stamping a node stamps every relation touching it](layers/mapa.md#edstd-b09--stamping-a-node-stamps-every-relation-touching-it) `EDSTD-B09`

- [A judgment holds only at the revisions it was given and only for its gate](layers/mapa.md#edstd-b10--a-judgment-holds-only-at-the-revisions-it-was-given-and-only-for-its-gate) `EDSTD-B10`

- [A judgment survives the next check round](layers/mapa.md#edstd-b11--a-judgment-survives-the-next-check-round) `EDSTD-B11`

- [One judgment per gate, its date kept when nothing changed](layers/mapa.md#edstd-b12--one-judgment-per-gate-its-date-kept-when-nothing-changed) `EDSTD-B12`

- [The stale relations are listed](layers/mapa.md#edstd-b13--the-stale-relations-are-listed) `EDSTD-B13`

- [Judging keeps a waiver another gate recorded, and the gate that waived replaces its own](layers/mapa.md#edstd-b14--judging-keeps-a-waiver-another-gate-recorded-and-the-gate-that-waived-replaces-its-own) `EDSTD-B14`

- [The same round on the same day gives the same stamps](layers/mapa.md#edstd-i01--the-same-round-on-the-same-day-gives-the-same-stamps) `EDSTD-I01`

- [A round carries only the stamps it changed to the map on disk](layers/mapa.md#edstd-b15--a-round-carries-only-the-stamps-it-changed-to-the-map-on-disk) `EDSTD-B15`

- [The stamp carries the caller's date](layers/mapa.md#edstd-x01--the-stamp-carries-the-callers-date) `EDSTD-X01`

- [A new revision that proves nothing new keeps what was measured](layers/mapa.md#edstd-b16--a-new-revision-that-proves-nothing-new-keeps-what-was-measured) `EDSTD-B16`

- [A declared change keeps what was proven, and the lines only when asked](layers/mapa.md#edstd-b17--a-declared-change-keeps-what-was-proven-and-the-lines-only-when-asked) `EDSTD-B17`

- [A repair carries only the evidence that held at the file's revision before it](layers/mapa.md#edstd-b18--a-repair-carries-only-the-evidence-that-held-at-the-files-revision-before-it) `EDSTD-B18`

- [Saving stamps the current format and the running binary's release](layers/mapa.md#grprg-b01--saving-stamps-the-current-format-and-the-running-binarys-release) `GRPRG-B01`

- [The saved file starts with the fixed comment header](layers/mapa.md#grprg-b02--the-saved-file-starts-with-the-fixed-comment-header) `GRPRG-B02`

- [A save that changes only the writer's release leaves the file untouched](layers/mapa.md#grprg-b03--a-save-that-changes-only-the-writers-release-leaves-the-file-untouched) `GRPRG-B03`

- [A real change rewrites the file with the running release](layers/mapa.md#grprg-b04--a-real-change-rewrites-the-file-with-the-running-release) `GRPRG-B04`

- [A later release restamps an unchanged map, an earlier one does not](layers/mapa.md#grprg-b08--a-later-release-restamps-an-unchanged-map-an-earlier-one-does-not) `GRPRG-B08`

- [Loading a map in an unreadable format is refused](layers/mapa.md#grprg-b05--loading-a-map-in-an-unreadable-format-is-refused) `GRPRG-B05`

- [A saved graph loads back unchanged](layers/mapa.md#grprg-i01--a-saved-graph-loads-back-unchanged) `GRPRG-I01`

- [An unreadable map yields no graph at all](layers/mapa.md#grprg-x01--an-unreadable-map-yields-no-graph-at-all) `GRPRG-X01`

- [Loading a missing file returns the read error](layers/mapa.md#grprg-e01--loading-a-missing-file-returns-the-read-error) `GRPRG-E01`

- [Loading text that is not the map returns the parse error](layers/mapa.md#grprg-e02--loading-text-that-is-not-the-map-returns-the-parse-error) `GRPRG-E02`

- [Saving into a missing directory returns the write error](layers/mapa.md#grprg-e04--saving-into-a-missing-directory-returns-the-write-error) `GRPRG-E04`

- [A map read from bytes is read like one from disk](layers/mapa.md#grprg-b06--a-map-read-from-bytes-is-read-like-one-from-disk) `GRPRG-B06`

- [An empty signal loads as no signal](layers/mapa.md#grprg-b07--an-empty-signal-loads-as-no-signal) `GRPRG-B07`

- [With code as the anchor a code file tests through its own name](layers/mapa.md#tsunt-b01--with-code-as-the-anchor-a-code-file-tests-through-its-own-name) `TSUNT-B01`

- [With the spec as the anchor the code path gives the variables](layers/mapa.md#tsunt-b02--with-the-spec-as-the-anchor-the-code-path-gives-the-variables) `TSUNT-B02`

- [An override of the code's layer places its test elsewhere](layers/mapa.md#tsunt-b03--an-override-of-the-codes-layer-places-its-test-elsewhere) `TSUNT-B03`

- [A glob test template matches the tests it covers, a code template is read literally](layers/mapa.md#tsunt-b04--a-glob-test-template-matches-the-tests-it-covers-a-code-template-is-read-literally) `TSUNT-B04`

- [Only the map's tests are answered, each unit once, in order](layers/mapa.md#tsunt-b05--only-the-maps-tests-are-answered-each-unit-once-in-order) `TSUNT-B05`

## mapx

- [A VR test captures its unit's code file and images](layers/mapx.md#vrcpt-b01--a-vr-test-captures-its-units-code-file-and-images) `VRCPT-B01`

- [What captures nothing gets no edge](layers/mapx.md#vrcpt-b02--what-captures-nothing-gets-no-edge) `VRCPT-B02`

- [The closure of a capture is one level](layers/mapx.md#vrcpt-b03--the-closure-of-a-capture-is-one-level) `VRCPT-B03`

- [A VR code is read whole with its state](layers/mapx.md#vrcpt-b04--a-vr-code-is-read-whole-with-its-state) `VRCPT-B04`

- [The screen's change stales its capture, a component's does not](layers/mapx.md#vrcpt-i01--the-screens-change-stales-its-capture-a-components-does-not) `VRCPT-I01`

- [A contract test captures its API unit, its spec and the OpenAPI document](layers/mapx.md#vrcpt-b05--a-contract-test-captures-its-api-unit-its-spec-and-the-openapi-document) `VRCPT-B05`

- [A capture's closure reaches the uncaptured dependencies, and stops at a captured one](layers/mapx.md#vrcpt-b06--a-captures-closure-reaches-the-uncaptured-dependencies-and-stops-at-a-captured-one) `VRCPT-B06`

- [Parts Used names become composes edges](layers/mapx.md#vrcpt-b07--parts-used-names-become-composes-edges) `VRCPT-B07`

- [The captures a changed file reaches](layers/mapx.md#vrcpt-b08--the-captures-a-changed-file-reaches) `VRCPT-B08`

- [A dependency the code's flags declare reaches the capture's closure and the impact of a change, transitively](layers/mapx.md#vrcpt-b09--a-dependency-the-codes-flags-declare-reaches-the-captures-closure-and-the-impact-of-a-change-transitively) `VRCPT-B09`

## migra

- [A four-character code is widened, keeping it as the prefix](layers/migra.md#mgfcd-b01--a-four-character-code-is-widened-keeping-it-as-the-prefix) `MGFCD-B01`

- [The files of one unit get different code names](layers/migra.md#mgfcd-b02--the-files-of-one-unit-get-different-code-names) `MGFCD-B02`

- [A file's code is new and of five characters](layers/migra.md#mgfcd-b03--a-files-code-is-new-and-of-five-characters) `MGFCD-B03`

- [Only a text file with a comment syntax carries the line](layers/migra.md#mgfcd-b04--only-a-text-file-with-a-comment-syntax-carries-the-line) `MGFCD-B04`

- [The line goes below the header's opener, or in a new header](layers/migra.md#mgfcd-b05--the-line-goes-below-the-headers-opener-or-in-a-new-header) `MGFCD-B05`

- [A header carrying its unit's code turns it into a ref beside a code of its own](layers/migra.md#mgfcd-b06--a-header-carrying-its-units-code-turns-it-into-a-ref-beside-a-code-of-its-own) `MGFCD-B06`

- [The renamed codes are read old to current, a code renamed twice to the last one](layers/migra.md#mgfcd-b07--the-renamed-codes-are-read-old-to-current-a-code-renamed-twice-to-the-last-one) `MGFCD-B07`

- [The header's updated_at is set to the day given, and only in the header](layers/migra.md#mgfcd-b08--the-headers-updated-at-is-set-to-the-day-given-and-only-in-the-header) `MGFCD-B08`

## scan

- [The built-in directories are skipped when nothing is declared](layers/scan.md#scigs-b01--the-built-in-directories-are-skipped-when-nothing-is-declared) `SCIGS-B01`

- [A layer pointing inside a built-in directory re-enables it, a catch-all does not](layers/scan.md#scigs-b02--a-layer-pointing-inside-a-built-in-directory-re-enables-it-a-catch-all-does-not) `SCIGS-B02`

- [A gitignore negation re-enables a built-in directory](layers/scan.md#scigs-b03--a-gitignore-negation-re-enables-a-built-in-directory) `SCIGS-B03`

- [The records Anchors writes are never scanned](layers/scan.md#scigs-b04--the-records-anchors-writes-are-never-scanned) `SCIGS-B04`

- [Editor and system ephemera never become files to scan](layers/scan.md#scigs-b05--editor-and-system-ephemera-never-become-files-to-scan) `SCIGS-B05`

- [A slash anchors a gitignore pattern at the root, no slash matches at any depth](layers/scan.md#scigs-b06--a-slash-anchors-a-gitignore-pattern-at-the-root-no-slash-matches-at-any-depth) `SCIGS-B06`

- [A trailing slash ignores the directory and what is below it, but not a file of that name](layers/scan.md#scigs-b07--a-trailing-slash-ignores-the-directory-and-what-is-below-it-but-not-a-file-of-that-name) `SCIGS-B07`

- [The last matching gitignore rule decides](layers/scan.md#scigs-b08--the-last-matching-gitignore-rule-decides) `SCIGS-B08`

- [Without a loaded ignore set the fixed exclusions still hold](layers/scan.md#scigs-b09--without-a-loaded-ignore-set-the-fixed-exclusions-still-hold) `SCIGS-B09`

- [No declaration re-enables the machinery directories](layers/scan.md#scigs-i01--no-declaration-re-enables-the-machinery-directories) `SCIGS-I01`

- [A nested gitignore does not change what the scan sees](layers/scan.md#scigs-x01--a-nested-gitignore-does-not-change-what-the-scan-sees) `SCIGS-X01`

- [Only the progress suffix marks a progress file](layers/scan.md#prflp-b01--only-the-progress-suffix-marks-a-progress-file) `PRFLP-B01`

- [The companion path replaces the extension of the last segment](layers/scan.md#prflp-b02--the-companion-path-replaces-the-extension-of-the-last-segment) `PRFLP-B02`

- [Two sides ticking neighbouring items keep both ticks](layers/scan.md#prflp-b03--two-sides-ticking-neighbouring-items-keep-both-ticks) `PRFLP-B03`

- [An item only on their side is added after our last item](layers/scan.md#prflp-b04--an-item-only-on-their-side-is-added-after-our-last-item) `PRFLP-B04`

- [The lines that are not items survive the merge](layers/scan.md#prflp-b05--the-lines-that-are-not-items-survive-the-merge) `PRFLP-B05`

- [A promoted tick keeps our indentation](layers/scan.md#prflp-b06--a-promoted-tick-keeps-our-indentation) `PRFLP-B06`

- [The done count reads both tick letters at any indentation](layers/scan.md#prflp-b07--the-done-count-reads-both-tick-letters-at-any-indentation) `PRFLP-B07`

- [Merging a side with itself changes nothing](layers/scan.md#prflp-i01--merging-a-side-with-itself-changes-nothing) `PRFLP-I01`

- [A merge never unticks an item](layers/scan.md#prflp-i02--a-merge-never-unticks-an-item) `PRFLP-I02`

- [An item whose text differs between the sides is kept twice](layers/scan.md#prflp-x01--an-item-whose-text-differs-between-the-sides-is-kept-twice) `PRFLP-X01`

- [A nested region closes before the outer one](layers/scan.md#srrgs-b01--a-nested-region-closes-before-the-outer-one) `SRRGS-B01`

- [A deeply indented one-line region does not swallow the file](layers/scan.md#srrgs-b02--a-deeply-indented-one-line-region-does-not-swallow-the-file) `SRRGS-B02`

- [The region revision ignores changes outside it](layers/scan.md#srrgs-b03--the-region-revision-ignores-changes-outside-it) `SRRGS-B03`

- [A swapped close does not cascade into the following regions](layers/scan.md#srrgs-b04--a-swapped-close-does-not-cascade-into-the-following-regions) `SRRGS-B04`

- [A close without a code closes the open region](layers/scan.md#srrgs-b05--a-close-without-a-code-closes-the-open-region) `SRRGS-B05`

- [A file without regions has no defect](layers/scan.md#srrgs-b06--a-file-without-regions-has-no-defect) `SRRGS-B06`

- [Composition resolves relative to the script's directory](layers/scan.md#srrgs-b07--composition-resolves-relative-to-the-scripts-directory) `SRRGS-B07`

- [Duplicates and self-references are not dependencies](layers/scan.md#srrgs-b08--duplicates-and-self-references-are-not-dependencies) `SRRGS-B08`

- [A script with no composition has no dependency](layers/scan.md#srrgs-b09--a-script-with-no-composition-has-no-dependency) `SRRGS-B09`

- [One pairing defect is reported once](layers/scan.md#srrgs-i01--one-pairing-defect-is-reported-once) `SRRGS-I01`

- [Composition comes only from declared paths](layers/scan.md#srrgs-x01--composition-comes-only-from-declared-paths) `SRRGS-X01`

- [A region never closed is reported](layers/scan.md#srrgs-e01--a-region-never-closed-is-reported) `SRRGS-E01`

- [A close with no open region is reported](layers/scan.md#srrgs-e02--a-close-with-no-open-region-is-reported) `SRRGS-E02`

- [A close naming another code is reported](layers/scan.md#srrgs-e03--a-close-naming-another-code-is-reported) `SRRGS-E03`

- [Only files in a declared layer enter the scan](layers/scan.md#rpscr-b01--only-files-in-a-declared-layer-enter-the-scan) `RPSCR-B01`

- [What the ignore set excludes is not scanned](layers/scan.md#rpscr-b02--what-the-ignore-set-excludes-is-not-scanned) `RPSCR-B02`

- [A nested checkout is not scanned](layers/scan.md#rpscr-b03--a-nested-checkout-is-not-scanned) `RPSCR-B03`

- [A progress companion stays out of the map](layers/scan.md#rpscr-b04--a-progress-companion-stays-out-of-the-map) `RPSCR-B04`

- [An upstream workflow carries no codes](layers/scan.md#rpscr-b05--an-upstream-workflow-carries-no-codes) `RPSCR-B05`

- [The revision ignores line endings but not content](layers/scan.md#rpscr-b06--the-revision-ignores-line-endings-but-not-content) `RPSCR-B06`

- [Priority, then pattern length, then layer name decide the layer](layers/scan.md#rpscr-b07--priority-then-pattern-length-then-layer-name-decide-the-layer) `RPSCR-B07`

- [An exclusion removes a path from its layer](layers/scan.md#rpscr-b08--an-exclusion-removes-a-path-from-its-layer) `RPSCR-B08`

- [A Windows path is classified like its slash form](layers/scan.md#rpscr-b09--a-windows-path-is-classified-like-its-slash-form) `RPSCR-B09`

- [A heuristic decision is reported as an ambiguity, a declared priority is not](layers/scan.md#rpscr-b10--a-heuristic-decision-is-reported-as-an-ambiguity-a-declared-priority-is-not) `RPSCR-B10`

- [The unit's layer comes from the header before the path](layers/scan.md#rpscr-b11--the-units-layer-comes-from-the-header-before-the-path) `RPSCR-B11`

- [Codes cited in comments are not owned, and each code is listed once](layers/scan.md#rpscr-b12--codes-cited-in-comments-are-not-owned-and-each-code-is-listed-once) `RPSCR-B12`

- [The project's rule letters are recognised](layers/scan.md#rpscr-b13--the-projects-rule-letters-are-recognised) `RPSCR-B13`

- [The declared identity and the annotations are recorded](layers/scan.md#rpscr-b14--the-declared-identity-and-the-annotations-are-recorded) `RPSCR-B14`

- [The parent is read only inside the header](layers/scan.md#rpscr-b15--the-parent-is-read-only-inside-the-header) `RPSCR-B15`

- [Needs are plan paths for a plan and phase codes for a spec](layers/scan.md#rpscr-b16--needs-are-plan-paths-for-a-plan-and-phase-codes-for-a-spec) `RPSCR-B16`

- [Only a plan revises](layers/scan.md#rpscr-b17--only-a-plan-revises) `RPSCR-B17`

- [A non-spec file declares dependencies in its header](layers/scan.md#rpscr-b18--a-non-spec-file-declares-dependencies-in-its-header) `RPSCR-B18`

- [A YAML test script depends on the scripts it composes](layers/scan.md#rpscr-b19--a-yaml-test-script-depends-on-the-scripts-it-composes) `RPSCR-B19`

- [A spec's dependency table is read in any catalogue language](layers/scan.md#rpscr-b20--a-specs-dependency-table-is-read-in-any-catalogue-language) `RPSCR-B20`

- [A row with a malformed code is not a dependency](layers/scan.md#rpscr-b21--a-row-with-a-malformed-code-is-not-a-dependency) `RPSCR-B21`

- [The method keeps its backticks](layers/scan.md#rpscr-b22--the-method-keeps-its-backticks) `RPSCR-B22`

- [A declared file resolves through the src fallbacks](layers/scan.md#rpscr-b23--a-declared-file-resolves-through-the-src-fallbacks) `RPSCR-B23`

- [A realizes tag pairs with its rule in the three rule forms](layers/scan.md#rpscr-b24--a-realizes-tag-pairs-with-its-rule-in-the-three-rule-forms) `RPSCR-B24`

- [A tag after a blank line has no owning rule](layers/scan.md#rpscr-b25--a-tag-after-a-blank-line-has-no-owning-rule) `RPSCR-B25`

- [Only a spec declares rule tags](layers/scan.md#rpscr-b26--only-a-spec-declares-rule-tags) `RPSCR-B26`

- [A repeated pair is recorded once](layers/scan.md#rpscr-b27--a-repeated-pair-is-recorded-once) `RPSCR-B27`

- [A gated-by tag names only a flag scenario](layers/scan.md#rpscr-b28--a-gated-by-tag-names-only-a-flag-scenario) `RPSCR-B28`

- [A plan seeds only concrete spec and doctrine paths](layers/scan.md#rpscr-b29--a-plan-seeds-only-concrete-spec-and-doctrine-paths) `RPSCR-B29`

- [Header keys are read only inside the header](layers/scan.md#rpscr-b30--header-keys-are-read-only-inside-the-header) `RPSCR-B30`

- [The classification is stable across runs](layers/scan.md#rpscr-i01--the-classification-is-stable-across-runs) `RPSCR-I01`

- [A guide's dependencies table is not a dependency](layers/scan.md#rpscr-x01--a-guides-dependencies-table-is-not-a-dependency) `RPSCR-X01`

- [A spec does not read a dep header line](layers/scan.md#rpscr-x02--a-spec-does-not-read-a-dep-header-line) `RPSCR-X02`

- [A root that cannot be walked returns the error](layers/scan.md#rpscr-e01--a-root-that-cannot-be-walked-returns-the-error) `RPSCR-E01`

- [A layer file that cannot be read fails the walk](layers/scan.md#rpscr-e02--a-layer-file-that-cannot-be-read-fails-the-walk) `RPSCR-E02`

- [Rule tags follow the declared code length](layers/scan.md#rpscr-b31--rule-tags-follow-the-declared-code-length) `RPSCR-B31`

- [A file in its layer's support list is marked as support](layers/scan.md#rpscr-b32--a-file-in-its-layers-support-list-is-marked-as-support) `RPSCR-B32`

- [Only the given files are read, as the walk reads them](layers/scan.md#rpscr-b33--only-the-given-files-are-read-as-the-walk-reads-them) `RPSCR-B33`

- [The staged walk reads the index, not the tree](layers/scan.md#rpscr-b34--the-staged-walk-reads-the-index-not-the-tree) `RPSCR-B34`

- [The index reader reads what the commit records](layers/scan.md#rpscr-b35--the-index-reader-reads-what-the-commit-records) `RPSCR-B35`

- [The governed files where the tree and the index part](layers/scan.md#rpscr-b36--the-governed-files-where-the-tree-and-the-index-part) `RPSCR-B36`

- [The index reader confronts the tree at each read](layers/scan.md#rpscr-b37--the-index-reader-confronts-the-tree-at-each-read) `RPSCR-B37`

- [A rule is defined in any of the three forms](layers/scan.md#rpscr-b38--a-rule-is-defined-in-any-of-the-three-forms) `RPSCR-B38`

- [A spec's Parts Used names its components](layers/scan.md#rpscr-b39--a-specs-parts-used-names-its-components) `RPSCR-B39`

- [The header's ref names the units the file realizes](layers/scan.md#rpscr-b40--the-headers-ref-names-the-units-the-file-realizes) `RPSCR-B40`

- [The @dep and @no-dep flags of import lines are read, with the symbols each import brings](layers/scan.md#rpscr-b41--the-dep-and-no-dep-flags-of-import-lines-are-read-with-the-symbols-each-import-brings) `RPSCR-B41`

- [Each @used-by flag is read with the symbol declared below it](layers/scan.md#rpscr-b42--each-used-by-flag-is-read-with-the-symbol-declared-below-it) `RPSCR-B42`

- [Each @navigates and @no-nav flag is read with its screens, its rule and its call's line](layers/scan.md#rpscr-b43--each-navigates-and-no-nav-flag-is-read-with-its-screens-its-rule-and-its-calls-line) `RPSCR-B43`

- [A spec's Out rows are read by rule, each with a revision of the row alone](layers/scan.md#rpscr-b44--a-specs-out-rows-are-read-by-rule-each-with-a-revision-of-the-row-alone) `RPSCR-B44`

- [Only a marked workflow is owned upstream](layers/scan.md#upowp-b01--only-a-marked-workflow-is-owned-upstream) `UPOWP-B01`

- [A Windows path under the workflow directory is recognised](layers/scan.md#upowp-b02--a-windows-path-under-the-workflow-directory-is-recognised) `UPOWP-B02`

- [The shared-code flag and prose do not open a header](layers/scan.md#upowp-b03--the-shared-code-flag-and-prose-do-not-open-a-header) `UPOWP-B03`

- [An HTML or block-comment header ends with its comment](layers/scan.md#upowp-b04--an-html-or-block-comment-header-ends-with-its-comment) `UPOWP-B04`

- [A line-comment header ends at the first line that is not a comment](layers/scan.md#upowp-b05--a-line-comment-header-ends-at-the-first-line-that-is-not-a-comment) `UPOWP-B05`

- [A header comment that never closes runs to the end of the file](layers/scan.md#upowp-b06--a-header-comment-that-never-closes-runs-to-the-end-of-the-file) `UPOWP-B06`

- [Nothing after the header's comment belongs to it](layers/scan.md#upowp-i01--nothing-after-the-headers-comment-belongs-to-it) `UPOWP-I01`

- [Only the marker and the directory decide ownership](layers/scan.md#upowp-x01--only-the-marker-and-the-directory-decide-ownership) `UPOWP-X01`

- [Only a header at the top is the header, unless it says why it stands lower](layers/scan.md#upowp-b07--only-a-header-at-the-top-is-the-header-unless-it-says-why-it-stands-lower) `UPOWP-B07`

- [Only a comment whose first word is @anchors opens the header](layers/scan.md#upowp-b08--only-a-comment-whose-first-word-is-anchors-opens-the-header) `UPOWP-B08`

