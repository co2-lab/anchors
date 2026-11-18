---
title: Julgamento por IA
description: "Entenda como funcionam os gates que fazem perguntas semânticas, a gravação de vereditos com anchors judge e a diferença entre determinismo e avaliação semântica."
---

A maior parte dos [gates](/pt/docs/concepts/gates-e-vereditos/) do Anchors é **estritamente determinística**:
- Um arquivo tem ou não tem cabeçalho.
- O teste passou ou reprovou.
- O código de cenário existe ou não existe.

Essas perguntas são respondidas instantaneamente por analisadores léxicos, regexes ou pelo compilador.

No entanto, existem perguntas essenciais de qualidade que **nenhum script simples de regex consegue responder sozinho**:
- *"O código escrito realmente REALIZA a regra de negócio descrita na spec, ou é só uma casca vazia?"*
- *"A justificativa que o desenvolvedor deu para dispensar um teste (`@no-test`) é legítima e honesta?"*

Para essas perguntas, o Anchors introduz os **Gates de Julgamento por IA**.

---

## 1. Como Funciona um Gate de Julgamento

Em vez de usar `check:` (verificador interno) ou `run:` (comando de shell), o gate declara uma propriedade **`ask:`** no [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: regra-cumprida
    on: [spec]
    ask: "o trecho marcado no código realmente REALIZA o que a regra descreve?"
    blocking: false
    when: [ci]
```

Os três gates canônicos de julgamento por IA são:

| Gate | O que ele pergunta |
| --- | --- |
| [`regra-cumprida`](/pt/docs/gates//regra-cumprida/) | O código marcado entre as tags da regra realmente satisfaz a especificação? |
| [`no-test-prova-real`](/pt/docs/gates//no-test-prova-real/) | A prova alternativa apontada pela dispensa `@no-test` exercita de fato o comportamento? |
| [`mock-detect-cobre-o-dialeto`](/pt/docs/gates//mock-detect-cobre-o-dialeto/) | O padrão de regex declarado pelo projeto captura todos os dublês e mocks existentes no código? |

---

## 2. A Gravação do Veredito com `anchors judge`

Diferente de sistemas que chamam LLMs em tempo real em todo `git commit` (o que seria lento, caro e não reproduzível), o Anchors opera em duas etapas:

1. **A avaliação**: Um agente de IA (ou um revisor humano sênior) analisa a pergunta e os arquivos envolvidos.
2. **A gravação**: O veredito é gravado explicitamente no repositório com o comando `anchors judge`:

```sh
anchors judge src/services/Login.spec.md \
  --gate regra-cumprida \
  --verdict pass \
  --reason "A função valida as 3 tentativas e atualiza o campo LockUntil corretamente conforme AUTH-B01"
```

Uma vez gravado o laudo com o hash do arquivo analisado, o gate passa a rodar em milissegundos nos commits subsequentes. Se o arquivo de código ou a spec forem alterados, o laudo fica [stale (desatualizado)](/pt/docs/concepts/propagacao-e-impacto/) e exige um novo julgamento.

---

## 3. Os Três Valores de Julgamento

O comando `anchors judge` aceita três vereditos possíveis:

| Veredito | Significado | Quando usar |
| :---: | :--- | --- |
| **`pass`** | **Aprovado** | O código ou justificativa cumpre com fidelidade o que foi pedido. |
| **`fail`** | **Reprovado** | O código não realiza a regra, ou a justificativa é rasa/inadequada. |
| **`dispensado`** | **Não se aplica** | O arquivo declara que a implementação ainda está por fazer (`@TBD: code`). |

### Por que o `dispensado` é fundamental?
Se existissem apenas `pass` e `fail`, uma spec recém-escrita cujo código ainda não foi implementado mentiria em qualquer escolha:
- Se marcasse `pass`, estaria aprovando um código inexistente.
- Se marcasse `fail`, estaria gerando um alerta vermelho para um trabalho que acabou de ser planejado honestamente.

O veredito `dispensado` permite que você declare que o trabalho está pendente sem quebrar o pipeline.

---

## 4. Julgamento e Review

Um **julgamento** responde à pergunta de um gate com um veredito que carimba o mapa. Quando o agente que escreveu o código julga o próprio trabalho, um `pass` garante pouco — mas um `fail` ainda acha bugs reais. Uma **review** é outra coisa: um segundo olhar sobre o que um agente decidiu, cujo produto são **achados**, não um carimbo.

Qualquer gate pode declarar `review:`, independente de como mede, então um gate pode ser julgado *e* revisado. Um alvo fica **a revisar** até haver uma review registrada na revisão atual dele; uma mudança no alvo torna a review devida de novo. Reviews informam e nunca bloqueiam.

```yaml
gates:
  - name: rule-fulfilled
    review:
      ask: "Cada trecho marcado faz o que a sua regra diz?"
```

```bash
anchors review --pending                                   # o que está a revisar, e a pergunta
anchors review src/pay.go --gate rule-fulfilled --by human:ana
anchors review src/pay.go --gate rule-fulfilled --by agent:<fornecedor>/<modelo> \
  --findings "### 1. B02 estorna duas vezes (src/pay.go:41) — ..."
```

O `--by` registra quem revisou, do jeito que se nomeia. Os achados abrem a issue do alvo, como um julgamento que reprova: cada um vira um teste que falha e uma correção, ou é descartado com o motivo. Não há `waived` — uma review que não aconteceu continua devida. O `anchors check` mostra `🔍 N alvo(s) a revisar` numa linha própria.

---

## 5. Próximos Passos

- [Gates e Vereditos](/pt/docs/concepts/gates-e-vereditos/): A mecânica geral de execução de gates.
- [Maturidade e Saúde](/pt/docs/concepts/maturidade-e-saude/): Como gates de IA e determinísticos compõem a qualidade.
- [Referência do CLI](/pt/docs/cli//): Comandos completos do `anchors judge`.
