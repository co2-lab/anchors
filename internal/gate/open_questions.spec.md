<!-- @anchors
  code: OPQSP
  updated_at: 2026-09-19
  layer: gate
-->
# OpenQuestions — spec com pergunta em aberto não está pronta para implementar

> **Código**: `OPQSP`

## Visão Geral

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

## Domínio

| Entrada | Aceita | Fora do domínio | Quem garante |
| --- | --- | --- | --- |
| o artefato confrontado | qualquer nó do mapa | — (o gate não escolhe o alvo) | o motor de gates, que roteia pelo `on:` declarado |
| a seção de decisões | o título em qualquer variação que o catálogo e o projeto nomeiam | título que nenhum dos dois nomeia | esta unidade, pelo vocabulário de títulos aceitos |
| o conteúdo da seção | itens catalogados, prosa, ou nada | — (seção ausente é um caso, não um erro) | esta unidade: quem não abriu a seção não é cobrado |
| o léxico do título | declarado na Estrutura do projeto | — (omitido cai no do framework) | a Estrutura, com o catálogo do framework como piso |

## Efeitos

| Efeito | Descrição |
| --- | --- |
| `OPQSP-B01` | Artefato que não é spec sai sem veredito: só a spec tem decisão em aberto a cobrar. |
| `OPQSP-B02` | Quem ABRIU a seção é confrontado pelo conteúdo dela: item em aberto barra, seção fechada libera. |
| `OPQSP-B03` | Item em aberto BARRA: enquanto houver pergunta, a spec não passa por pronta. |
| `OPQSP-B04` | Seção fechada honestamente — aberta e sem item — libera. Dizer "não há pergunta" é diferente de não ter olhado. |
| `OPQSP-B05` | Item marcado como RESOLVIDO não bloqueia: a pergunta fica no rastro, e o que a fechou é a regra que nasceu dela. |
| `OPQSP-B06` | Cada pergunta precisa de CÓDIGO. Sem identidade ela não vira item rastreável nem sobrevive a uma reescrita da spec. |
| `OPQSP-B07` | `OpenDecisions` CONTA as decisões pendentes de uma spec, para quem precisa do número em vez do veredito — é o que permite reportar a pendência como ponta sistêmica, no mesmo estatuto de um sinal ausente. A contagem lê o léxico do projeto pela mesma via do confronto: contar zero numa spec cuja seção se chama outra coisa afirmaria "não há decisão pendente" sobre uma spec cheia delas, que é o silêncio que esta unidade existe para eliminar. |

## Invariantes

| Regra | Vale sempre | Como se prova |
| --- | --- | --- |
| `OPQSP-I01` | Prosa não é item. Texto explicativo dentro da seção não conta como pergunta — senão o autor aprenderia a não explicar nada. | escreve prosa na seção sem item catalogado e verifica que não bloqueia |
| `OPQSP-I02` | A fronteira da seção é respeitada: o que vem depois dela não é lido como pergunta. Sem isso, a spec inteira viraria seção de decisões. | escreve itens numa seção seguinte e verifica que só os da seção contam |
| `OPQSP-I03` | A coluna que diz no que a pergunta VIRA não é a identidade dela, e preenchê-la não fecha a pergunta. São duas coisas: o destino previsto e a resposta dada. | preenche a coluna de destino sem resolver e verifica que ainda bloqueia |

## Restrições

| Regra | Limite | Por quê |
| --- | --- | --- |
| `OPQSP-X01` | Não julga se a pergunta é BOA nem se a resposta é certa. | A régua é determinística: existe item em aberto, ou não existe. Avaliar o mérito de uma dúvida é julgamento, e julgamento é de outra classe de gate. |
| `OPQSP-X02` | Não REPROVA a spec que não tem a seção — registra a pendência e diz como fechá-la. | A ausência não distingue "tudo foi decidido" de "a seção foi apagada", e as duas pedem coisas opostas. Reprovar seria tratar migração como defeito; calar seria o silêncio que o gate existe para eliminar. O veredito fica indeterminado e ENSINA a saída: fechar com a declaração de que não há pergunta, ou escrever o que não se decidiu. |

## Dependências

| Cód | Arquivo | Método | Camada |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | núcleo — o título da seção pode ser do léxico do projeto, e o gate lê a Estrutura para saber |

## Decisões em aberto

| Código | Pergunta | Quem decide | Vira |
| --- | --- | --- | --- |

nenhuma
