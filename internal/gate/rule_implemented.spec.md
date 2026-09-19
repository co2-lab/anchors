<!-- @anchors
  code: RLIMR
  updated_at: 2026-09-19
  layer: gate
-->
# RuleImplemented — a spec cataloga regras, e o código mostra que as realizou

> **Código**: `RLIMR`

## Visão Geral

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

## Domínio

| Entrada | Aceita | Fora do domínio | Quem garante |
| --- | --- | --- | --- |
| o artefato confrontado | qualquer nó do mapa | — (o gate não escolhe o alvo) | o motor de gates, que roteia pelo `on:` declarado |
| a spec | com regras catalogadas, ou sem nenhuma | — (spec sem regra é um caso, não um erro) | esta unidade: sem regra não há o que cobrar |
| o código ligado | o que a aresta de realização apontar, ou nenhum | — | esta unidade: sem código ligado o assunto não existe |
| a dispensa por regra | a marca na linha da regra, com razão escrita | marca nua, sem razão | esta unidade: dispensa sem porquê não dispensa |
| a exigência de marcação | declarada na Estrutura do projeto | — (omitida vale como não exigida) | a Estrutura: enquanto o projeto não declara, a dívida é pendência e não reprovação |

## Efeitos

| Efeito | Descrição |
| --- | --- |
| `RLIMR-B01` | Spec cujas regras não aparecem no código é ACUSADA, e o veredito nomeia quais ficaram sem realização. |
| `RLIMR-B02` | Regra dispensada com razão escrita fecha a conta: a declaração vale como resposta, e o que ela dispensa deixa de ser cobrado. |
| `RLIMR-B03` | Unidade anterior à prática vira PENDÊNCIA, não reprovação — enquanto o projeto não declara que exige a marcação. |
| `RLIMR-B04` | Declarada a exigência na Estrutura, a pendência vira reprovação: é o ato de dizer "aqui a migração acabou". |
| `RLIMR-B05` | Spec sem código ligado não é assunto deste gate: sem a peça do outro lado não há confronto a fazer. |
| `RLIMR-B06` | A dispensa pode nomear a regra que ela cobre, ou valer para todas quando não nomeia nenhuma. |

## Invariantes

| Regra | Vale sempre | Como se prova |
| --- | --- | --- |
| `RLIMR-I01` | Exigir a marcação nunca pune quem já marca. Quem fez o trabalho antes da exigência não pode reprovar por tê-lo feito. | declara a exigência sobre uma unidade que marca tudo e verifica que ela passa |
| `RLIMR-I02` | A identidade sobrevive à renomeação: código marcado com o nome anterior continua valendo. Perder a marca num rename transformaria estabilidade de identidade em dívida nova. | marca o código com o nome antigo e verifica que a regra segue reconhecida |

## Restrições

| Regra | Limite | Por quê |
| --- | --- | --- |
| `RLIMR-X01` | Não julga se a implementação está CERTA — só se ela existe e se declara. | Se o código cumpre o que a regra diz é julgamento, e julgamento é de outra classe de gate. A régua aqui é determinística: a marca existe, ou a dispensa existe com razão. |
| `RLIMR-X02` | Não exige marca de TODA regra. | Restrição é satisfeita pela ausência de código, e ausência não tem onde receber comentário. Cobrar as 3.121 ocorrências que a régua ingênua produziria treinaria a equipe a ignorar a lista inteira. |

## Dependências

| Cód | Arquivo | Método | Camada |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `Graph` | núcleo — o código ligado se alcança pela aresta, não por convenção de nome |
| DEP2 | `internal/config/config.go` | `Config` | núcleo — a exigência de marcação é declarada na Estrutura |

## Decisões em aberto

| Código | Pergunta | Quem decide | Vira |
| --- | --- | --- | --- |

nenhuma
