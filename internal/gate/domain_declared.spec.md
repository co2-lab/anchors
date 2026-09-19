<!-- @anchors
  code: DMDCD
  updated_at: 2026-09-19
  layer: gate
-->
# DomainDeclared — a spec declara o que a unidade ACEITA, e quem barra o inválido

> **Código**: `DMDCD`

## Visão Geral

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

## Domínio

| Entrada | Aceita | Fora do domínio | Quem garante |
| --- | --- | --- | --- |
| o artefato confrontado | qualquer nó do mapa | — (o gate não escolhe o alvo) | o motor de gates, que roteia pelo `on:` declarado |
| o conteúdo do artefato | qualquer texto, inclusive vazio | — (texto ausente é um caso, não um erro) | esta unidade: conteúdo sem a seção é REPROVAÇÃO, não exceção |
| a seção de domínio | o título em qualquer idioma do catálogo, e as grafias que o projeto usa | título que o catálogo não nomeia | esta unidade, pelo vocabulário de títulos aceitos |

> `Quem garante` não pode ficar vazio nem dizer só "não é meu": se ninguém garante, o
> dever é órfão — e é exatamente aí que a entrada inválida passa. É a mesma exigência que
> este gate faz das specs que confronta, aplicada a ele próprio.

## Efeitos

| Efeito | Descrição |
| --- | --- |
| `DMDCD-B01` | Artefato que não é spec sai do confronto sem veredito: o gate não tem jurisdição sobre código, teste ou guia. |
| `DMDCD-B02` | Spec SEM a seção de domínio REPROVA. A ausência não é silêncio: é a afirmação não feita. |
| `DMDCD-B03` | Spec com a seção ABERTA e VAZIA reprova também — abrir o título sem declarar nada é o mesmo furo com aparência de conformidade. |
| `DMDCD-B04` | Cada entrada declarada precisa nomear QUEM garante. Entrada sem dono reprova, e o veredito nomeia quais ficaram órfãs. |
| `DMDCD-B05` | A dispensa é DECLARADA e com razão escrita. Quem não tem entrada externa registra isso na spec, e o gate se cala — mas fica o rastro de que alguém olhou. |
| `DMDCD-B06` | Linha preenchida só com marcador de pendência não é declaração: o molde intocado não afirma nada. |
| `DMDCD-B07` | Entrada cujo dono está nomeado passa — é o outro lado da mesma régua, e o que a torna satisfazível. |

## Invariantes

| Regra | Vale sempre | Como se prova |
| --- | --- | --- |
| `DMDCD-I01` | A dispensa exige RAZÃO. Uma marca de dispensa nua não silencia o gate — o silêncio sem porquê é o que ele existe para impedir. | confronta uma spec com a dispensa sem razão e verifica que o veredito ainda cobra |
| `DMDCD-I02` | O veredito de reprovação NOMEIA o que está errado — qual entrada ficou sem dono, ou que a seção falta. Um gate que reprova sem dizer o quê transfere o trabalho de diagnóstico para quem lê. | confronta uma spec com entrada órfã e verifica que o nome dela aparece no veredito |

## Restrições

| Regra | Limite | Por quê |
| --- | --- | --- |
| `DMDCD-X01` | Não julga se a entrada declarada está CERTA — só se ela existe e tem dono. | Se o conjunto de valores aceitos corresponde ao domínio real é julgamento, e julgamento é de outra classe de gate. Aqui a régua é a PRESENÇA da declaração, que é determinística. |
| `DMDCD-X02` | Não confronta o código para conferir se a validação existe de fato. | Esta camada lê TEXTO. Cruzar a declaração com a implementação é trabalho do gate relacional, que tem o mapa; fazê-lo aqui duplicaria a régua em dois lugares que divergiriam. |

## Dependências

| Cód | Arquivo | Método | Camada |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `KindSpec` | núcleo — o gate precisa do KIND do nó para saber se tem jurisdição |
| DEP2 | `internal/config/config.go` | `Config` | núcleo — o léxico de títulos aceitos vem da Estrutura do projeto |

## Decisões em aberto

| Código | Pergunta | Quem decide | Vira |
| --- | --- | --- | --- |

nenhuma
