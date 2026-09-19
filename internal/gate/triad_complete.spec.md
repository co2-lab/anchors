<!-- @anchors
  code: TRCMT
  updated_at: 2026-09-19
  layer: gate
-->
# TriadComplete — as peças que realizam uma spec EXISTEM

> **Código**: `TRCMT`

## Visão Geral

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

## Domínio

| Entrada | Aceita | Fora do domínio | Quem garante |
| --- | --- | --- | --- |
| o artefato confrontado | qualquer nó do mapa | — (o gate não escolhe o alvo) | o motor de gates, que roteia pelo `on:` declarado |
| o mapa | um grafo construído, ou nenhum | — (mapa ausente é um caso, não um erro) | esta unidade: sem mapa o veredito é indeterminado, nunca aprovação |
| a camada do alvo | camada regida, reconhecida, ou nenhuma declarada | — | esta unidade, pelo regime declarado na Estrutura |
| a dispensa por unidade | a marca com razão escrita ao lado | marca nua, sem razão | esta unidade: dispensa sem porquê não dispensa |

## Efeitos

| Efeito | Descrição |
| --- | --- |
| `TRCMT-B01` | Artefato que não é spec sai sem veredito: só a spec tem trinca a cobrar. |
| `TRCMT-B02` | Camada RECONHECIDA (regime declarativo) sai sem veredito: ela não tem spec nem trinca por definição. |
| `TRCMT-B03` | Sem mapa o veredito é INDETERMINADO. Aprovar sem poder olhar seria afirmar o que não se mediu. |
| `TRCMT-B04` | Spec com as três peças ligadas passa. |
| `TRCMT-B05` | Spec a que falta alguma peça reprova, e o veredito NOMEIA quais faltam e onde cada uma nasce. |
| `TRCMT-B06` | A camada pode dispensar uma peça em bloco, declarado na Estrutura. |
| `TRCMT-B07` | A unidade pode dispensar uma peça na própria spec, com razão escrita — é a granularidade que a dispensa por camada não alcança. |

## Invariantes

| Regra | Vale sempre | Como se prova |
| --- | --- | --- |
| `TRCMT-I01` | O teste é alcançado em DOIS saltos — spec → feature → teste —, porque quem aponta o teste é a feature. Conferir o teste direto na spec acusaria falta de teste no projeto inteiro. | liga a trinca em dois saltos e verifica que o gate a considera completa |
| `TRCMT-I02` | Dispensar o teste e escrever cenário na feature é CONTRADIÇÃO, e reprova. As duas afirmações não convivem: ou o cenário é real e alguém precisa prová-lo, ou não deveria existir. | declara a dispensa, liga uma feature com cenário, e verifica a reprovação |
| `TRCMT-I03` | Dispensar o teste exige dizer ONDE a prova está, e o lugar tem de existir. Referência órfã reprova. | declara a dispensa apontando um alvo inexistente e verifica a reprovação |
| `TRCMT-I04` | A dispensa vale só para a peça declarada. Dispensar uma nunca dispensa as outras. | declara a dispensa de uma peça e verifica que as demais seguem cobradas |

## Restrições

| Regra | Limite | Por quê |
| --- | --- | --- |
| `TRCMT-X01` | Não confronta se as peças CASAM entre si — só se existem. | Casar é o trabalho dos gates relacionais. Este existe justamente porque eles falham aberto quando a peça não existe; fazer os dois aqui duplicaria a régua. |
| `TRCMT-X02` | Não julga a QUALIDADE de nenhuma peça. | Um teste vazio satisfaz este gate, e é correto: a régua aqui é a EXISTÊNCIA. Quem confronta o conteúdo é outro gate, e confundir os dois faria este reprovar por motivo que não sabe medir. |

## Dependências

| Cód | Arquivo | Método | Camada |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `Graph` | núcleo — as peças são arestas, e sem o grafo não há o que olhar |
| DEP2 | `internal/config/config.go` | `Config` | núcleo — o regime da camada e a dispensa em bloco são declarados na Estrutura |

## Decisões em aberto

| Código | Pergunta | Quem decide | Vira |
| --- | --- | --- | --- |

nenhuma
