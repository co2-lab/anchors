# Flow-driven work — o trabalho DIRIGIDO por fluxo

> Documento de DESENHO, anterior à implementação. A funcionalidade nasce ISOLADA; as
> ligações com o resto do Anchors vêm depois.

## O problema, medido em projeto real

No `galaxy` há 58 problemas abertos e um guia de 232 linhas (`guides/destravar.md`) que
cataloga **oito abordagens** para destravar, cada uma com o seu "quando esta é a certa", e
uma regra explícita na linha 11 ("a regra dos DOIS ATAQUES").

E mesmo assim, o dossiê do problema 142 abre com este aviso:

> *"LEIA A ÚLTIMA RODADA PRIMEIRO. Este dossiê tem cinco passagens e cada uma re-apontou o
> alvo da anterior — o título já foi dado por errado e depois confirmado... pelo menos dois
> já mandaram uma rodada para o sítio errado."*

**Cinco rodadas, duas no alvo errado**, com a informação certa escrita e disponível.

A causa não é falta de documento — é a NATUREZA do documento. Um catálogo de "o que fazer
caso X" é INSUMO: exige que quem trabalha lembre de consultar, escolha certo entre oito, e
não pule etapa. São três oportunidades de errar, a cada rodada, para sempre.

Um fluxo inverte a carga. Em vez de *"aqui estão oito abordagens, escolha"*, ele diz
*"você está no estado MEDIDO com dois ataques gastos nesta abordagem; as saídas válidas
daqui são 3, 4 ou 8"*. A IA não precisa LEMBRAR da regra dos dois ataques — a regra é a
única saída disponível.

É a mesma inversão que o Anchors já faz com os gates: em vez de "lembre de escrever o
teste", o `triad-complete` torna a ausência visível. O que falta é a mesma inversão para o
PROCESSO, e não só para o artefato.

## O que já existe, e por que não basta

`anchors work <etapa> --for <alvo>` compõe, em 34ms, um procedimento por ação: o que ler,
em que ordem, o que não é escopo, os passos, os gates que virão, e até a árvore de decisão
sobre onde uma regra nova deve morar. Camadas customizam com `work:`.

Falta a ele UMA coisa, e é a que muda tudo: **o procedimento é uma lista linear, sem
estado**. Ele não sabe que você já tentou medir duas vezes, não tem ramo condicional, e
responde igual na primeira e na quinta rodada.

## O desenho

### O fluxo é um GRAFO PRÓPRIO — e não entra no mapa de artefatos

Decidido: grafo, e não banco. Medido neste repositório, compor um procedimento custa
**34ms** e o `map build` inteiro (548 nós, 232K) custa **562ms**. Uma engine resolveria o
que já custa milissegundos, e traria estado binário fora do git, diff ilegível em PR e uma
dependência num CLI que se vende como binário único.

Mas o grafo do FLUXO é separado do grafo do MAPA, e a separação não é organização — é
correção. O mapa liga ARTEFATOS (spec→código→feature→teste), e é sobre ele que o `impact`
e a Propagação caminham. Um estado de fluxo ("MEDIDO", "DOIS ATAQUES GASTOS") não é
artefato de ninguém: misturá-los faria o `impact` atravessar DECISÕES como se fossem
dependências, e responder que mudar um arquivo afeta um estado de processo.

Há uma segunda razão, mecânica: `mapx.Build` parte de um grafo VAZIO e o preenche varrendo
o disco (três chamadas no código, todas assim). O que não vem de arquivo some no build
seguinte — e o fluxo precisa de nós que não são arquivos.

Então: MESMO ARQUIVO, CHAVE SEPARADA.

    anchors.graph.yaml
      nodes:   os ARTEFATOS      (quem preenche: `map build`)
      edges:   as dependências
      flow:    os FLUXOS         (quem preenche: `flow build`)

Um arquivo só porque é um grafo só, de conteúdos distintos — dois arquivos obrigariam quem
lê a saber de antemão em qual olhar. Chave separada porque os conteúdos não podem se
misturar, pela razão do parágrafo acima.

E a chave está NA STRUCT do mapa, não fora: o `Save` serializa a struct inteira, e o
comentário do `Load` já documenta o custo de não estar — *"os campos que ele reconhece
carregam, os que não reconhece somem, e a próxima gravação escreve o que sobrou. É assim
que se perde dado sem nada acusar."* Um fluxo que não fosse campo seria apagado no
`map build` seguinte, em silêncio.

Por isso o `map build` PRESERVA a chave `flow` (como já preserva os carimbos de
validação), e o `flow build` é quem a preenche.

### O arquivo de fluxo

    flows/<nome>.flow.md

Markdown com header `@anchors`, e serve às duas exigências de uma vez: é a FONTE do grafo
de fluxo, e é a DOCUMENTAÇÃO que quem trabalha lê. Uma fonte só — ter duas é o defeito que
o `docs-covered` existe para pegar.

Os ESTADOS são catalogados dentro dele, com código, na mesma gramática das regras de spec:

    ### DSTRV-N03 — MEDIDO: o termo foi isolado e há número

    Saídas:
    - `DSTRV-N05` quando a medição aponta um alvo concreto
    - `DSTRV-N04` quando dois ataques já falharam nesta abordagem

A aresta `flows-to` liga estado a estado, e vive no `anchors.flow.yaml`.

### A visualização é RUNTIME

`anchors flow show <nome>` monta o diagrama na hora, a partir do grafo de fluxo. Nada de
mermaid guardado em arquivo: um diagrama versionado envelhece em relação ao fluxo que
descreve, e o `docs-fresh` já mostrou o custo disso neste repositório.

## O que NÃO entra nesta primeira rodada

A funcionalidade nasce ISOLADA, e isso é decisão sua. Não entram agora:

- onde o ESTADO de cada trabalho em curso é guardado (o "em que nó estou");
- a ligação com `anchors work`, `next`, `queue` e o watcher;
- o gate que confronta se quem trabalhou seguiu o fluxo.

Essas são as LIGAÇÕES, e elas dependem de o fluxo existir primeiro. A pergunta "onde o
estado vive" tem resposta diferente para um card, para um problema e para um PR, e
respondê-la agora fixaria a errada.

## O que entra

### O ciclo de comandos, espelhando o do mapa

O fluxo tem o seu próprio ciclo, e a simetria com `anchors map` é deliberada: quem já
conhece o mapa não aprende um vocabulário novo.

    anchors flow build            varre `flows/*.flow.md` e escreve `anchors.flow.yaml`
    anchors flow show <nome>      o diagrama, montado em RUNTIME a partir do grafo
    anchors flow next <nome> --from <estado>
                                  as saídas VÁLIDAS daqui — é o comando que dirige

O `build` é separado do `map build` pela razão do parágrafo acima: são grafos de
conteúdos diferentes, e um `impact` que atravessasse estados de processo responderia
errado. A MECÂNICA (varrer, montar nós e arestas, serializar YAML versionado) é a mesma e
se reaproveita; o que não se mistura é o grafo.

`flow next` é o que faz a funcionalidade ser dirigida por fluxo, e não mais um documento:
ele responde "daqui, os passos válidos são estes", em vez de listar tudo e pedir que
alguém escolha.

### A lista

1. `flows/<nome>.flow.md` — o formato, com estados catalogados e saídas condicionais
2. `anchors flow build` → `anchors.flow.yaml` (o mapa de artefatos NÃO muda)
3. `anchors flow show <nome>` — a visualização em runtime
4. `anchors flow next <nome> --from <estado>` — as saídas válidas de um estado
5. os gates do fluxo: todo estado é alcançável? toda saída leva a um estado que existe?
   há estado terminal?
6. `anchors guide flow` e `anchors new flow`
7. um fluxo DEFAULT do Anchors, para o projeto sobrescrever

## Decisões em aberto

| Código | Pergunta | Quem decide | Vira |
| --- | --- | --- | --- |
| `Q01` | O fluxo default do Anchors é o do CICLO (plan→spec→code→feature→test→review) ou o de DESTRAVAR (o do galaxy)? O primeiro exercita a ligação com o `work`; o segundo ataca o problema que originou o pedido. | você | o primeiro `flows/*.flow.md` semeado pelo `init` |
