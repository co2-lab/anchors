---
title: O Grafo e o Mapa de Dependências
description: "Como o Anchors constrói o mapa de relacionamentos do seu software e calcula o caminho mínimo de impacto quando algo muda."
---

Em um projeto com centenas de arquivos, se você alterar uma regra na especificação de autenticação, quais outros arquivos do sistema serão impactados?
- A tela de login precisa mudar?
- A API de checkout é afetada?
- Quais testes específicos precisam ser reexecutados?

Se você não souber a resposta exata, tem apenas duas escolhas ruins: ou roda a suíte inteira de testes e checagens (o que pode levar dezenas de minutos), ou não roda nada e reza para que nada tenha quebrado em produção.

O Anchors resolve isso mantendo um **Mapa de Dependências (Grafo)** materializado no arquivo `anchors.graph.yaml`.

---

## 1. O que é o Grafo?

O Grafo do Anchors é uma rede de **nós** (arquivos do seu projeto) conectados por **arestas** (relacionamentos intencionais entre eles):

```
┌─────────────────────┐               ┌─────────────────────┐
│   Login.spec.md     │──governs─────▶│      Login.go       │
│      (Nó Spec)      │               │     (Nó Código)     │
└──────────┬──────────┘               └─────────────────────┘
           │
      covered-by
           │
           ▼
┌─────────────────────┐               ┌─────────────────────┐
│    Login.feature    │──tested-by───▶│    Login_test.go    │
│    (Nó Feature)     │               │     (Nó Teste)      │
└─────────────────────┘               └─────────────────────┘
```

### O que são os Nós?
Cada arquivo registrado no projeto vira um **Nó** no grafo, classificado por seu `kind`:
- `spec`, `feature`, `test`, `code`, `doc`, `guide`, `plan`, `product`, `flag`.

### O que são as Arestas?
As arestas são as ligações entre os arquivos, indicando qual papel um tem em relação ao outro:
- **`governs`**: A spec governa o arquivo de código.
- **`covered-by`**: A spec é coberta pela feature.
- **`tested-by`**: A feature é testada pelo arquivo de teste.
- **`depends-on`**: Um módulo de código importa ou depende de outro.
- **`realizes`**: Uma spec realiza uma [Doutrina de Produto](/pt/docs/concepts/doutrina-de-produto/).
- **`gated-by`**: Uma regra de negócio é condicionada por uma [Feature Flag](/pt/docs/concepts/feature-flags/).

---

## 2. Grafo Virtual vs Grafo Material

Esta é uma das ideias mais elegantes do Anchors:

| Tipo | Onde vive | Como funciona |
| --- | --- | --- |
| **Grafo Virtual** | Na declaração de [Camadas](/pt/docs/layers/) do `anchors.yaml` | É a **planta da casa**: diz como os arquivos *deveriam* se relacionar, antes mesmo de você escrever a primeira linha de código. |
| **Grafo Material** | No arquivo `anchors.graph.yaml` | É a **realidade no disco**: o mapa gerado pelo CLI que reflete exatamente quais arquivos existem agora e como estão conectados por código. |

O grafo virtual resolve o "paradoxo do projeto novo": quando você inicia um projeto do zero, o mapa está vazio, mas o Anchors já sabe o que esperar em cada pasta porque a planta das camadas já dita o formato esperado.

---

## 3. O Caminho Mínimo de Impacto

Quando você executa o comando:

```sh
anchors check
```

O Anchors não perde tempo inspecionando todos os milhares de arquivos do repositório a esmo. Ele consulta o grafo e o Git para identificar:

1. **Quais arquivos foram alterados** desde o último commit.
2. **Quem depende desses arquivos** (caminhando pelas arestas do grafo).
3. **Quais nós ficaram desatualizados ([stale](/pt/docs/concepts/propagacao-e-impacto/))**.

Se você alterou apenas a regra de formatação de telefone na spec de cadastro, o Anchors saberá exatamente que apenas o validador de telefone e o seu teste correspondente precisam ser verificados. Todo o restante do sistema permanece verde sem custo computacional.

---

## 4. Como o Mapa é Gerado e Atualizado

O mapa é construído de forma transparente pelo comando:

```sh
anchors map
```

Esse comando lê os arquivos do disco, extrai os cabeçalhos `@anchors`, localiza os códigos de cenário e reconstrói o arquivo `anchors.graph.yaml`.

No dia a dia, quando você roda `anchors check`, o CLI atualiza o mapa em milissegundos antes de rodar os [gates](/pt/docs/concepts/gates-e-vereditos/).

---

## 5. Próximos Passos

- [A Unidade](/pt/docs/concepts/unidade/): Veja as arestas fundamentais que ligam a spec ao teste.
- [Propagação e Impacto](/pt/docs/concepts/propagacao-e-impacto/): Como as alterações viajam pelo grafo.
- [Gate layer-boundary](/pt/docs/gates//layer-boundary/): Como o grafo impede importações proibidas entre camadas.
