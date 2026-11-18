---
title: Propagação e a Onda de Impacto
description: "Entenda como o Anchors detecta o que ficou desatualizado (stale) e guia o sistema até que todas as partes voltem a ser coerentes."
---

Em um sistema complexo, nenhuma alteração vive isolada. Quando você altera uma regra na especificação de um produto:
- O cenário na feature precisa ser atualizado.
- O teste correspondente precisa ser reescrito ou ajustado.
- A função no código fonte precisa ser alterada.
- As outras unidades que chamavam essa função precisam ser verificadas.

Se uma IA ou um desenvolvedor alterar a spec e esquecer de atualizar o código, o sistema entra em **dessincronia**.

No Anchors, o mecanismo que resolve isso é chamado de **Propagação (Propagation)**.

---

## 1. O que é a Propagação?

A **Propagação** é o motor que faz uma alteração num ponto qualquer do repositório percorrer o organismo do software através do [Grafo de Dependências](/pt/docs/concepts/grafo-e-mapa/).

Ela funciona como uma **onda em um lago**: quando você joga uma pedra na água (faz uma alteração), a onda se espalha em círculos concêntricos até tocar todas as margens afetadas.

```
       1. PEDRA NA ÁGUA
       (Você altera a Spec)
              │
              ▼
       2. ONDA PRIMÁRIA
       (A Feature ligada fica STALE / Desatualizada)
              │
              ▼
       3. ONDA SECUNDÁRIA
       (O Teste e o Código ficam STALE)
              │
              ▼
       4. ONDA TERCIÁRIA
       (Módulos dependentes são reavaliados)
              │
              ▼
       5. SISTEMA EM REPOUSO
       (Todos os gates passam; coerência restaurada)
```

---

## 2. O Conceito de `Stale` (Obsoleto / Defasado)

Um nó no grafo é considerado **stale** quando ele próprio ou algum dos nós dos quais ele depende foi alterado mais recentemente do que ele.

Por exemplo:
- Se você alterou o arquivo `Usuario.spec.md` hoje às 14:00, mas o arquivo `Usuario_test.go` foi modificado pela última vez ontem, o teste está **stale** em relação à spec.
- O gate [`evidence-fresh`](/pt/docs/gates//evidence-fresh/) entra em ação e avisa: *"O resultado deste teste não vale mais, porque o código ou a spec mudaram depois que ele rodou"*.

Isso impede o golpe mais comum em pipelines de CI: testes antigos passando com sucesso sobre um código que já mudou de comportamento.

---

## 3. A Disciplina Spec-First na Propagação

O fluxo natural de propagação do Anchors sempre nasce na **especificação**:

1. **Altere a Spec primeiro**: Declare a nova regra ou altere uma regra existente com seu [código de identidade](/pt/docs/concepts/rastreabilidade-e-codigos/).
2. **Propague para a Feature**: Atualize ou adicione o cenário em Gherkin.
3. **Propague para o Código e o Teste**: Implemente a nova lógica e os testes que a provam.
4. **Rode a checagem**: Execute `anchors check`. O Anchors verificará todo o caminho de impacto. Se algum dependente quebrou, o CLI lista exatamente qual arquivo ficou para trás.

---

## 4. Por que a Propagação Salva Horas de Trabalho de IA?

Quando você usa agentes autônomos de IA para codificar sem o Anchors, a IA frequentemente entra em loops infinitos: ela mexe num arquivo, quebra outro que ela não conhecia, tenta consertar o segundo e quebra o primeiro novamente.

Com o Anchors:
- O agente sabe a **ordem exata** em que deve trabalhar (da spec para o teste, e do teste para o código).
- O agente recebe do `anchors check` a lista exata dos arquivos que ficaram stale.
- O agente só para quando a onda de propagação terminar e todos os gates estiverem verdes (`✓`).

---

## 5. Próximos Passos

- [O Grafo e o Mapa de Dependências](/pt/docs/concepts/grafo-e-mapa/): A fiação por onde a onda se propaga.
- [O Fluxo de Trabalho](/pt/docs/workflow/): Como executar o ciclo de propagação no dia a dia.
- [Gate evidence-fresh](/pt/docs/gates//evidence-fresh/): Como o Anchors garante que os testes não estão defasados.
