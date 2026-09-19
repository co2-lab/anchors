<!-- anchors:generated from doct/camadas/gate.md.tmpl — DO NOT EDIT: run `anchors docs build` -->


# Camada: gate




## CDLNG — CodeLanguage — the code does not go back to mixing languages









#### CDLNG-B01 — An identifier in the wrong language is accused, and an English one passes

```gherkin
    Given a file declaring a function named in Portuguese
    And another declaring a function named in English
    When the gate confronts the file
    Then the Portuguese identifier is accused
    And the English one passes without noise
```

#### CDLNG-B02 — The verdict returns the word that accused

```gherkin
    Given a file declaring an identifier in the wrong language
    When the gate confronts it
    Then the verdict names that exact word, so the reader does not hunt the whole file
```

#### CDLNG-B03 — Only a DECLARATION is the subject

```gherkin
    Given a file whose Portuguese words appear outside any declaration
    When the gate confronts it
    Then it accuses nothing, because what declares no identifier is not read
```

#### CDLNG-B04 — The declarations are found in every form the language offers

```gherkin
    Given a file declaring identifiers as function, type, variable and constant
    When the gate confronts it
    Then every one of those declarations is read, not only the most common form
```

#### CDLNG-B05 — Deciding one word is separate from deciding a whole identifier

```gherkin
    Given a compound identifier made of several words
    When the gate decides its language
    Then it breaks the identifier into its words and decides each one
    And that separation is what lets the length floor apply per word instead of to
      the identifier as a whole
```

#### CDLNG-I01 — A short word does not count

```gherkin
    Given a file declaring identifiers below the length floor
    When the gate confronts it
    Then none is accused, because below the floor there is no language to infer and
      accusing there is the noise that costs a gate its credibility
```

#### CDLNG-X01 — The gate does not read comments

```gherkin
    Given a file whose comments are written in the team's own language
    And every identifier is in English
    When the gate confronts it
    Then it returns Pass, because the comment carries the measurement and the why
```

#### CDLNG-X02 — The gate does not read user-facing text

```gherkin
    Given a file whose user-facing strings are written in the project's language
    And every identifier is in English
    When the gate confronts it
    Then it returns Pass, because that text goes through the translation catalog
```

#### CDLNG-X03 — The gate does not use a dictionary to decide the language

```gherkin
    Given a file declaring compound identifiers in English that no common dictionary holds
    When the gate confronts it
    Then none is accused, because the dictionary approach was measured at ninety percent
      false positives — and a gate that wrong is switched off, defending nothing
```


## DCRQD — DocRequired — the aggregated document the unit must feed









#### DCRQD-B01 — A mandatory document that does not exist fails

```gherkin
    Given a project declaring a document as mandatory for this unit's layer
    And that document does not exist on disk
    When the gate confronts the unit
    Then it returns Fail
```

#### DCRQD-B02 — A document that exists and does not mention the unit fails

```gherkin
    Given the mandatory document exists and never names this unit
    When the gate confronts the unit
    Then it returns Fail, because existence alone would approve an empty file created
      to silence the gate
```

#### DCRQD-B03 — A mention by the identity code counts as documented

```gherkin
    Given the mandatory document cites the unit's identity code
    When the gate confronts the unit
    Then it returns Pass
```

#### DCRQD-B04 — A mention by the file name also counts

```gherkin
    Given the mandatory document cites the unit's file name and not its code
    When the gate confronts the unit
    Then it returns Pass, because the document speaks of the unit either way
```

#### DCRQD-B05 — Satisfying one of two duties is not enough

```gherkin
    Given two documents declared mandatory for this unit's layer
    And only one of them mentions the unit
    When the gate confronts the unit
    Then it returns Fail, because each document is charged on its own
```

#### DCRQD-B06 — Without a declaration nothing is charged

```gherkin
    Given a project that declares no mandatory document
    When the gate confronts the unit
    Then it returns Pass, because the ruler is what the project committed to, not what
      one supposes it owes
```

#### DCRQD-B07 — A layer with no trigger is not charged

```gherkin
    Given a mandatory document whose trigger names another layer
    And the document exists and does not mention this unit
    When the gate confronts the unit
    Then it returns Pass, because this layer triggers no duty
```

#### DCRQD-B08 — Aggregated, the verdict is one per document

```gherkin
    Given three units of a triggering layer and one mandatory document
    When the gate runs over the whole project
    Then it reports ONE verdict for that document, not one per unit
```

#### DCRQD-I01 — The duty starts from the spec, not from the code

```gherkin
    Given a unit whose spec and code both exist
    When the gate confronts the project
    Then the charge lands on the spec, because the spec is what declares the unit
```

#### DCRQD-I02 — The layer used is the UNIT's, not the node's

```gherkin
    Given a spec whose node layer is the spec layer
    And the unit it describes belongs to a triggering layer
    When the gate confronts it
    Then the duty is charged, because reading the node's layer would charge every spec
      of the project the same duty, or none
```

#### DCRQD-I03 — Without a map the aggregated verdict is skipped

```gherkin
    Given no graph built
    When the gate runs aggregated
    Then it does not approve, because approving without being able to look would stamp
      what was never measured
```

#### DCRQD-X01 — The gate does not understand the document's content

```gherkin
    Given the mandatory document cites the unit and describes it wrongly
    When the gate confronts the unit
    Then it returns Pass, because the gate separates "not documented" from "documented" —
      judging the quality of the documentation is another ruler
```

#### DCRQD-X02 — The gate does not decide which documents are mandatory

```gherkin
    Given a project whose Structure declares no duty for this layer
    And a document that a reviewer would consider obviously required
    When the gate confronts the unit
    Then it returns Pass, because inventing duties would charge what nobody committed to
```


## DMDCD — DomainDeclared — the spec declares what the unit ACCEPTS, and who blocks the invalid









#### DMDCD-B01 — An artifact that is not a spec leaves without a verdict

```gherkin
    Given a node whose kind is code, test or feature
    When the gate confronts it
    Then it returns Skip, because only a spec has a domain to declare
```

#### DMDCD-B02 — A spec without the domain section is failed

```gherkin
    Given a spec that catalogues rules and never opens the domain section
    And no declared waiver anywhere in it
    When the gate confronts it
    Then it returns Fail
    And the verdict names the waiver marker, so whoever reads it learns the declared way out
```

#### DMDCD-B03 — A section opened and left empty is failed

```gherkin
    Given a spec whose domain section has only the table header, with no data row
    When the gate confronts it
    Then it returns Fail, because opening the title without declaring anything is the same
      gap wearing the appearance of compliance
```

#### DMDCD-B06 — A row filled only with placeholders is not a declaration

```gherkin
    Given a spec whose domain section carries a single row reading "TODO" in every column
    When the gate confronts it
    Then it returns Fail, because the untouched template asserts nothing
```

#### DMDCD-B04 — An entry with no owner is failed, and the verdict names it

```gherkin
    Given a spec declaring the entry "chave" with its accepted values
    And the column that says who guarantees it is left blank
    When the gate confronts it
    Then it returns Fail
    And the verdict names "chave", so the reader does not have to hunt for the orphan
```

#### DMDCD-B07 — An entry whose owner is named passes

```gherkin
    Given a spec declaring the entry "chave" with its accepted values
    And the column that says who guarantees it reads "the interface, before calling"
    When the gate confronts it
    Then it returns Pass
```

#### DMDCD-B05 — A waiver with a written reason silences the gate

```gherkin
    Given a spec with no domain section
    And a waiver marker followed by the reason "receives only typed values from its own code"
    When the gate confronts it
    Then it returns Skip, and the reason stays in the spec as the record that someone looked
```

#### DMDCD-I01 — A bare waiver, with no reason, does not waive

```gherkin
    Given a spec with no domain section
    And a waiver marker with nothing written after it
    When the gate confronts it
    Then it returns Fail, because a waiver with no why is the silence the gate exists to end
```

#### DMDCD-I02 — A non-answer in the owner column is not an owner

```gherkin
    Given a spec declaring the entry "chave"
    And the column that says who guarantees it reads "I do not validate (MTVRX-X04)"
    When the gate confronts it
    Then it returns Fail, because carrying a restriction into the owner column names nobody —
      it is the sentence that creates the orphan
```

#### DMDCD-X01 — The gate does not judge whether the declared domain is correct

```gherkin
    Given a spec declaring the entry "month" as accepting any text, which is wider than the real domain
    And the column that says who guarantees it names the caller
    When the gate confronts it
    Then it returns Pass, because the ruler here is the PRESENCE of the declaration —
      whether the accepted set matches reality is judgment, and judgment belongs to another gate
```

#### DMDCD-X02 — The gate does not read the code to check the validation exists

```gherkin
    Given a spec whose domain section is complete and every entry has a named owner
    And the code that realises it performs no validation at all
    When the gate confronts it
    Then it returns Pass, because this layer reads TEXT — crossing the declaration with the
      implementation belongs to the relational gate, which has the map
```


## OPQSP — OpenQuestions — a spec with an open question is not ready to implement









#### OPQSP-B01 — An artifact that is not a spec leaves without a verdict

```gherkin
    Given a node whose kind is code, feature or test
    When the gate confronts it
    Then it returns Skip, because only a spec has open decisions to demand
```

#### OPQSP-B02 — Whoever OPENED the section is confronted by its content

```gherkin
    Given a spec whose open-decisions section is opened and carries no item
    When the gate confronts it
    Then it returns Pass, because the content is what the ruler reads
```

#### OPQSP-B03 — An open item bars the spec

```gherkin
    Given a spec whose open-decisions section carries one catalogued question
    When the gate confronts it
    Then it returns Fail, because while there is a question the spec does not pass as ready
```

#### OPQSP-B04 — A section closed honestly releases the spec

```gherkin
    Given a spec whose open-decisions section is opened and carries no item
    When the gate confronts it
    Then it returns Pass, because saying "there is no question" differs from not having looked
```

#### OPQSP-B05 — An item marked as resolved does not block

```gherkin
    Given a spec whose question is marked resolved, citing the rule born from it
    When the gate confronts it
    Then it returns Pass, and the question stays in the trail instead of being swept away
```

#### OPQSP-B06 — A question with no code is charged

```gherkin
    Given a spec whose open-decisions section carries a question written without a code
    When the gate confronts it
    Then it returns Fail, because without identity the question is neither a traceable
      item nor survives a rewrite of the spec
```

#### OPQSP-B07 — The count of pending decisions reads the project's own lexicon

```gherkin
    Given a spec whose open-decisions section carries two questions
    And the project names that section with its own wording
    When the number of pending decisions is counted
    Then it answers two, because counting zero over a section named otherwise would
      assert "no pending decision" about a spec full of them
```

#### OPQSP-I01 — Prose is not an item

```gherkin
    Given a spec whose open-decisions section carries explanatory text and no catalogued item
    When the gate confronts it
    Then it returns Pass, because otherwise the author would learn to explain nothing
```

#### OPQSP-I02 — The section boundary is respected

```gherkin
    Given a spec carrying catalogued items in a section that FOLLOWS the open decisions
    And the open-decisions section itself is empty
    When the gate confronts it
    Then it returns Pass, because otherwise the whole spec would read as a section of decisions
```

#### OPQSP-I03 — Filling in what the question BECOMES does not close the question

```gherkin
    Given a spec whose question names the rule it is expected to become
    And the question is not marked resolved
    When the gate confronts it
    Then it returns Fail, because the intended destination and the answer given are two
      different things
```

#### OPQSP-X01 — The gate does not judge whether the question is good

```gherkin
    Given a spec whose only open question is trivial
    When the gate confronts it
    Then it returns Fail all the same, because the ruler is deterministic — an open item
      exists, or it does not; judging the merit of a doubt belongs to another gate
```

#### OPQSP-X02 — A spec with no section is a pending item, and the verdict teaches the way out

```gherkin
    Given a spec with rules catalogued and no open-decisions section at all
    When the gate confronts it
    Then it returns Pending, because the absence does not tell "everything was decided"
      apart from "the section was deleted"
    And the verdict says how to close it: declare that there is no question, or write
      what is not yet decided
```


## PGNHN — PaginationHonored — what promises a SET does not return the first page in silence









#### PGNHN-B01 — A limit received from the caller passes

```gherkin
    Given an exported function whose signature takes the limit as a parameter
    When the gate confronts it
    Then it returns Pass, because the page is deliberate and whoever asked for it
      knows there is more
```

#### PGNHN-B02 — A limit hidden in a default value is accused

```gherkin
    Given an exported function whose name promises the whole set
    And the limit lives in a default value the caller never sees
    When the gate confronts it
    Then it returns Fail, because the hundred-and-first row is never processed and
      nobody is told
```

#### PGNHN-B03 — The NAME bounds the promise

```gherkin
    Given an exported function whose name promises no set at all
    And it returns a partial result
    When the gate confronts it
    Then it returns Pass, because there is no promise to break
```

#### PGNHN-B04 — Sibling functions that paginate are the proof by asymmetry

```gherkin
    Given a module whose sibling functions loop until the cursor is exhausted
    And one exported function returns a single page
    When the gate confronts it
    Then it returns Fail, because the author knew the pattern — the one that does not
      paginate is forgetfulness, not decision
```

#### PGNHN-B05 — A waiver with a written reason leaves the report

```gherkin
    Given an exported function carrying a waiver marker followed by the reason
    When the gate confronts it
    Then it returns Pass, and the function no longer appears in the report
```

#### PGNHN-B06 — The verdict offers the way out

```gherkin
    Given an exported function accused of hiding the limit
    When the gate confronts it
    Then the verdict names the waiver marker, so whoever reads it learns the declared
      way out instead of guessing
```

#### PGNHN-B07 — Without a declared dialect the verdict is undetermined

```gherkin
    Given a project that declares no dialect for its stack
    When the gate confronts any code
    Then it returns Pending, because approving without being able to read the code
      would stamp what was never checked
```

#### PGNHN-I01 — The ruler is agnostic across languages

```gherkin
    Given the same hidden-limit defect written in two different stacks
    And each project declares its own dialect
    When the gate confronts both
    Then both are accused, because the confronted truth — the name promises a set, the
      return is partial — belongs to no language
```

#### PGNHN-I02 — Where the construct is not recognised, the gate stays silent

```gherkin
    Given code whose pattern the declared dialect does not reach
    When the gate confronts it
    Then it accuses nothing, because a false positive here teaches the team to ignore
      the gate
```

#### PGNHN-I03 — A cursor with no loop does not count as pagination

```gherkin
    Given an exported function that returns the cursor and never walks it
    When the gate confronts it
    Then it returns Fail, because the consumer is left with the same slice, now wearing
      the appearance of completeness
```

#### PGNHN-I04 — A provider prefix in the name does not hide the promise

```gherkin
    Given an exported function whose name carries the provider prefix before the promise
    When the gate confronts it
    Then it returns Fail, because what the name says holds wherever it comes from
```

#### PGNHN-X01 — The gate does not invent a cursor the provider does not offer

```gherkin
    Given an exported function calling a dependency with no pagination mechanism
    When the gate confronts it
    Then it returns Pass, because accusing the absence of a mechanism the dependency
      lacks would hand the author a defect that is not theirs
```

#### PGNHN-X02 — The gate does not measure performance or page size

```gherkin
    Given an exported function that exposes its limit and returns a small page
    When the gate confronts it
    Then it returns Pass, because judging whether a hundred is many depends on the
      domain, and that is the project's decision
```


## RLIMR — RuleImplemented — a spec catalogues rules, and the code shows it realized them









#### RLIMR-B01 — A spec whose rules the code ignores is accused, and the verdict names them

```gherkin
    Given a spec cataloguing three rules, of which the code marks only the first
    When the gate confronts it
    Then it returns Fail
    And the verdict names the rule left without realisation, so the reader does not
      have to diff spec against code by hand
```

#### RLIMR-B02 — A rule waived with a written reason closes the account

```gherkin
    Given a spec cataloguing a restriction the code cannot mark, because it is satisfied
      by the ABSENCE of code
    And the rule's row carries a waiver marker followed by the reason
    When the gate confronts it
    Then it returns Pass, because the declaration is the answer the gate asked for
```

#### RLIMR-B03 — A unit that predates the practice is a pending item, not a failure

```gherkin
    Given a spec whose rules carry no mark anywhere in the code
    And the project does not declare that it requires marking
    When the gate confronts it
    Then it returns Pending
    And the verdict NAMES the debt, instead of pretending approval
```

#### RLIMR-B04 — Declaring the requirement turns the pending item into a failure

```gherkin
    Given the same unmarked spec
    And the project declares in its Structure that marking is required
    When the gate confronts it
    Then it returns Fail, because declaring the requirement is the act of saying
      "here the migration is over"
```

#### RLIMR-B05 — A spec with no linked code is not this gate's subject

```gherkin
    Given a spec that no code realises yet
    When the gate confronts it
    Then it returns Skip, because without the piece on the other side there is no
      confrontation to make — and who accuses the absence is the triad gate
```

#### RLIMR-B06 — A waiver with no named rule covers every rule of the spec

```gherkin
    Given a spec whose waiver marker names no rule in particular
    When the gate confronts it
    Then it returns Pass, because the waiver was declared for the unit as a whole
```

#### RLIMR-I01 — Requiring the marking never punishes whoever already marks

```gherkin
    Given a spec whose every rule is marked in the code
    And the project declares that marking is required
    When the gate confronts it
    Then it returns Pass, because whoever did the work before the requirement cannot
      fail for having done it
```

#### RLIMR-I02 — The identity survives a rename

```gherkin
    Given a code marked with the identity the unit carried before being renamed
    When the gate confronts it
    Then it returns Pass, because losing the mark on a rename would turn identity
      stability into new debt
```

#### RLIMR-X01 — The gate does not judge whether the implementation is correct

```gherkin
    Given a spec whose rules are all marked in the code
    And the marked code does something other than what the rule describes
    When the gate confronts it
    Then it returns Pass, because the ruler here is deterministic — the mark exists, or
      the waiver exists with a reason; whether the code honours the rule is judgment
```

#### RLIMR-X02 — The gate does not demand a mark on EVERY rule

```gherkin
    Given a spec whose restrictions are satisfied by the absence of code
    And those rows declare their waiver with a reason
    When the gate confronts it
    Then it returns Pass, because absence has nowhere to receive a comment — demanding
      it would produce thousands of findings and teach the team to ignore the list
```


## TRCMT — TriadComplete — the pieces that realize a spec EXIST









#### TRCMT-B01 — An artifact that is not a spec leaves without a verdict

```gherkin
    Given a node whose kind is code, feature or test
    When the gate confronts it
    Then it returns Skip, because only a spec has a triad to demand
```

#### TRCMT-B02 — A recognised layer leaves without a verdict

```gherkin
    Given a spec whose layer is declared with the declarative regime
    When the gate confronts it
    Then it returns Skip, because a recognised layer has neither spec nor triad by definition
```

#### TRCMT-B03 — Without a map the verdict is undetermined

```gherkin
    Given a spec and no graph built
    When the gate confronts it
    Then it returns Pending, because approving without being able to look would assert
      what was never measured
```

#### TRCMT-B04 — A spec with the three pieces linked passes

```gherkin
    Given a spec linked to its code, to its feature and to the test that proves it
    When the gate confronts it
    Then it returns Pass
```

#### TRCMT-B05 — A spec missing a piece is failed, and the verdict names which

```gherkin
    Given a spec linked to its code and to nothing else
    When the gate confronts it
    Then it returns Fail
    And the verdict names the feature and the test, and says where each one is born
```

#### TRCMT-B06 — The layer may waive a piece for every spec in it

```gherkin
    Given a layer declaring that the test edge is optional
    And a spec of that layer linked to its code and its feature only
    When the gate confronts it
    Then it returns Pass, because the waiver is declared in the Structure, in plain sight
```

#### TRCMT-B07 — The unit may waive a piece in its own spec, with a written reason

```gherkin
    Given a spec carrying a waiver marker for the test, followed by the reason
    And the spec is linked to its code and its feature
    When the gate confronts it
    Then it returns Pass, because the decision belongs to the unit and is written where
      whoever reads the spec will see it
```

#### TRCMT-I01 — The test is reached in two hops, through the feature

```gherkin
    Given a spec linked to a feature, and that feature linked to the test
    And no edge going straight from the spec to the test
    When the gate confronts it
    Then it returns Pass, because who points at the test is the FEATURE — checking it
      straight on the spec would report a missing test across the whole project
```

#### TRCMT-I02 — Waiving the test while the feature carries a scenario is a contradiction

```gherkin
    Given a spec whose test is waived
    And a linked feature carrying one scenario
    When the gate confronts it
    Then it returns Fail, because either the scenario is real and someone must prove it,
      or it should not exist
```

#### TRCMT-I03 — Waiving the test demands saying where the proof is, and the place must exist

```gherkin
    Given a spec whose test waiver points at a target that no file realises
    When the gate confronts it
    Then it returns Fail, because an orphan reference proves nothing
```

#### TRCMT-I04 — A waiver covers only the piece it declares

```gherkin
    Given a spec that waives the feature and is linked to its code only
    When the gate confronts it
    Then it returns Fail naming the test, because waiving one piece never waives the others
```

#### TRCMT-X01 — The gate does not confront whether the pieces MATCH one another

```gherkin
    Given a spec whose linked feature describes a behaviour the test does not prove
    And the three pieces exist and are linked
    When the gate confronts it
    Then it returns Pass, because matching is the work of the relational gates — this one
      exists precisely because they fail open when the piece is absent
```

#### TRCMT-X02 — The gate does not judge the QUALITY of any piece

```gherkin
    Given a spec linked to a feature with no scenarios and to an empty test
    When the gate confronts it
    Then it returns Pass, because the ruler here is EXISTENCE — confronting the content
      belongs to another gate, and mixing the two would fail by a criterion this one
      cannot measure
```



