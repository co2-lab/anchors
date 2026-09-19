<!-- anchors:generated from doct/camadas/gate.md.tmpl — DO NOT EDIT: run `anchors docs build` -->


# Camada: gate




## DMDCD — DomainDeclared — a spec declara o que a unidade ACEITA, e quem barra o inválido

Confronta uma spec contra a pergunta que o resto do framework não faz: **o que esta
unidade aceita de entrada, e quem garante que o inválido nunca chega?**

A lacuna que ele fecha foi medida: 71% das specs de um projeto real tinham seção de
regras ou efeitos, e apenas 14% diziam o que a unidade aceita. O framework inteiro é
construído sobre catalogar EFEITOS — o que a unidade faz —, e todo defeito de borda
encontrado em três rodadas de review adversarial morava no que ninguém tinha declarado.

A distinção que dá razão ao gate: `## Restrições` diz o que a unidade NÃO faz, e empurra
o dever para FORA; `## Domínio` diz o que ela ACEITA, e nomeia QUEM fica com ele.
Escrever mais restrições não fecha nada — cria órfãos, porque cada "não é meu" precisa de
alguém do outro lado.




| Regra | Vale sempre | Como se prova |
| --- | --- | --- |
| `DMDCD-I01` | A dispensa exige RAZÃO. Uma marca de dispensa nua não silencia o gate — o silêncio sem porquê é o que ele existe para impedir. | confronta uma spec com a dispensa sem razão e verifica que o veredito ainda cobra |
| `DMDCD-I02` | O veredito de reprovação NOMEIA o que está errado — qual entrada ficou sem dono, ou que a seção falta. Um gate que reprova sem dizer o quê transfere o trabalho de diagnóstico para quem lê. | confronta uma spec com entrada órfã e verifica que o nome dela aparece no veredito |


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


## OPQSP — OpenQuestions — spec com pergunta em aberto não está pronta para implementar

Confronta uma spec contra as decisões que ela ainda NÃO tomou, e a mantém fora do "pronto"
enquanto houver pergunta em aberto.

A classe de defeito é a AMBIGUIDADE NÃO RESOLVIDA — a mais barata de evitar e a mais cara
de descobrir tarde. O caminho é sempre o mesmo: a spec não decide algo que o código
precisa; quem implementa escolhe uma leitura defensável e segue; a escolha nunca é
confrontada com quem tinha a resposta; o produto sai com a leitura errada. Nenhum outro
gate pega, porque todas as peças existem e se referenciam — o defeito é uma decisão que
ninguém tomou.

O que este gate acrescenta ao conselho "não chute, registre e reporte" é um LUGAR
declarado para o registro. Sem lugar, registrar vira comentário de PR que morre no merge.
Com lugar, a pergunta é um item de trabalho visível, e a spec só fica implementável quando
a seção esvazia.

O ciclo pretendido: quem escreve percebe o que não sabe e escreve na seção; o gate acusa
enquanto houver item; a pergunta é levada a quem decide; a resposta VIRA REGRA, com
código, e o item sai da seção.




| Regra | Vale sempre | Como se prova |
| --- | --- | --- |
| `OPQSP-I01` | Prosa não é item. Texto explicativo dentro da seção não conta como pergunta — senão o autor aprenderia a não explicar nada. | escreve prosa na seção sem item catalogado e verifica que não bloqueia |
| `OPQSP-I02` | A fronteira da seção é respeitada: o que vem depois dela não é lido como pergunta. Sem isso, a spec inteira viraria seção de decisões. | escreve itens numa seção seguinte e verifica que só os da seção contam |
| `OPQSP-I03` | A coluna que diz no que a pergunta VIRA não é a identidade dela, e preenchê-la não fecha a pergunta. São duas coisas: o destino previsto e a resposta dada. | preenche a coluna de destino sem resolver e verifica que ainda bloqueia |


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


## RLIMR — RuleImplemented — a spec cataloga regras, e o código mostra que as realizou

Confronta a spec contra o código na direção que faltava: **a spec não ficou falando
sozinha?**

É o inverso do gate que valida referências. Aquele confere que os códigos CITADOS pelo
código existem na spec; este confere que as regras DECLARADAS na spec ganharam
implementação. Sem ele, uma spec pode declarar cinco regras novas e o código não ganhar
linha nenhuma — com todos os gates verdes, porque a spec existe, o código existe, e os
dois se referenciam pelo cabeçalho.

Medido: uma spec de interface ganhou cinco regras e 98 linhas, e o arquivo correspondente
tinha ZERO ocorrência do assunto. Metade da entrega era código morto declarado como
pronto, e nenhum dos 26 gates perguntou. O defeito só apareceu quando alguém leu spec e
código na mesma passada.

**A régua é a declaração, não a adivinhação.** Exigir toda regra marcada seria falso por
construção — medido contra 592 unidades, daria 3.121 achados, e nem as unidades bem-feitas
passariam: das que marcam o código, nenhuma marca 100%. A razão é boa: uma restrição ("a
unidade NÃO faz Y") é satisfeita pela AUSÊNCIA de código, e ausência não tem onde receber
marca. Mas "ao menos uma" também não serve — separa quem implementou de quem não
implementou e não diz nada sobre as outras quinze regras. Então quem escreve a spec
DECLARA, regra a regra, se ela tem código.




| Regra | Vale sempre | Como se prova |
| --- | --- | --- |
| `RLIMR-I01` | Exigir a marcação nunca pune quem já marca. Quem fez o trabalho antes da exigência não pode reprovar por tê-lo feito. | declara a exigência sobre uma unidade que marca tudo e verifica que ela passa |
| `RLIMR-I02` | A identidade sobrevive à renomeação: código marcado com o nome anterior continua valendo. Perder a marca num rename transformaria estabilidade de identidade em dívida nova. | marca o código com o nome antigo e verifica que a regra segue reconhecida |


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


## TRCMT — TriadComplete — as peças que realizam uma spec EXISTEM

Confronta uma spec de camada REGIDA contra a pergunta mais simples da trinca: **as peças
que a realizam existem?** O código que ela especifica, a feature que a cobre, e o teste
que a prova.

Existe porque os gates relacionais FALHAM ABERTO por construção. Sem teste ligado, o
confronto feature↔teste devolve "nada a confrontar ainda" em vez de reprovar; sem código,
o de dependências idem. O efeito colateral é grave: uma spec sozinha, sem nenhuma
implementação, atravessa TODOS os gates e o pipeline conclui "pode promover" — o verde
certificando trabalho que não existe.

Este gate fecha o buraco pelo lado positivo. Em vez de perguntar "as peças casam?" — o
que exige que elas existam —, pergunta "as peças existem?".




| Regra | Vale sempre | Como se prova |
| --- | --- | --- |
| `TRCMT-I01` | O teste é alcançado em DOIS saltos — spec → feature → teste —, porque quem aponta o teste é a feature. Conferir o teste direto na spec acusaria falta de teste no projeto inteiro. | liga a trinca em dois saltos e verifica que o gate a considera completa |
| `TRCMT-I02` | Dispensar o teste e escrever cenário na feature é CONTRADIÇÃO, e reprova. As duas afirmações não convivem: ou o cenário é real e alguém precisa prová-lo, ou não deveria existir. | declara a dispensa, liga uma feature com cenário, e verifica a reprovação |
| `TRCMT-I03` | Dispensar o teste exige dizer ONDE a prova está, e o lugar tem de existir. Referência órfã reprova. | declara a dispensa apontando um alvo inexistente e verifica a reprovação |
| `TRCMT-I04` | A dispensa vale só para a peça declarada. Dispensar uma nunca dispensa as outras. | declara a dispensa de uma peça e verifica que as demais seguem cobradas |


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



