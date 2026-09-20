# `product/` — a regra que atravessa alvos

> Documento de DESENHO, anterior à implementação. Registra o problema medido, o desenho
> proposto, e as decisões que faltam tomar antes de escrever código.

## O problema

Toda spec do Anchors tem um ALVO: ela descreve uma unidade, e a co-locação amarra as
quatro pontas da tríade ao mesmo diretório. Isso é a força do modelo — cada regra tem
endereço, e o gate sabe onde confrontar.

E é também o limite. Numa aplicação real, muita regra de negócio **atravessa alvos**:
"o limite de crédito vale para o cadastro, para a simulação e para a aprovação" não
pertence a nenhuma das três telas — pertence ao produto que as três servem.

Hoje só há duas saídas, e as duas são ruins:

- **duplicar** a regra nas três specs, e aceitar que elas divirjam na primeira mudança;
- **escolher uma dona** arbitrária, e deixar as outras duas referenciando implicitamente
  uma regra que mora num lugar que não é o delas.

O Anchors já RECONHECE que o eixo vertical existe: a tag `@feature: auth` marca o módulo
vertical, "mesmo quando espalhado por camadas" (HEADER_GUIDE). Mas é um rótulo solto —
não há artefato que defina a funcionalidade, não há regras próprias dela, não há gate.

## O desenho

Uma pasta `product/` no root, no mesmo precedente do `plans/`: artefato transversal, fora
da árvore de alvos, com guide próprio. Dentro dela, `product/<nome>.doctrine.md`.

### O nome

`features/` foi a primeira ideia e colide: `feature` já é o `.feature` do Gherkin neste
vocabulário — `kind: feature`, aresta `covered-by`, gate `feature-not-empty`, 17 usos de
`KindFeature` no motor. E o `.feature` é padrão EXTERNO, que não é nosso para renomear:
quem muda é a ponta nova.

`product/` é melhor que as alternativas que eu havia proposto
(`capability`, `behavior`), por três razões:

- **diz de ONDE a regra vem**, não o que ela é. A regra transversal existe porque alguém
  decidiu o produto — e é isso que a distingue da regra local, que existe porque a
  unidade precisa dela para funcionar.
- **segue o padrão do `plans/`**: uma pasta nomeia um LUGAR, não um tipo. `capability/`
  teria sido o único diretório do projeto nomeado por um conceito abstrato.
- **`behavior` carregava colisão**: `B` é 543 das 841 regras deste repositório, e
  *"a regra `CRED-B03` realiza o behavior `LIMIT`"* confundiria em toda revisão.

### O sufixo do arquivo: `.doctrine.md`

`rules` colide: `## Rules` já é a seção onde a spec cataloga suas regras locais.
Tecnicamente seria outro namespace (seção não é kind), mas o termo já tem lugar, e o
critério aqui é SINERGIA, não ausência de conflito.

`decisions` foi considerado e DESCARTADO, e não por colisão — por INVERSÃO. No
vocabulário atual, `## Open Decisions` é a seção do que ainda NÃO foi decidido: a
pergunta aberta, com código `-Q`, que o gate `open-questions-resolved` cobra que alguém
responda. E `DECISIONS.md` no root registra decisões de arquitetura já tomadas.

Usar `decisions` para a regra de produto faria a mesma palavra significar coisas opostas
na mesma spec: *"a regra realiza a decision `LIMIT-R03`"* (firme, vigente) convivendo com
*"isto está em Open Decisions"* (pendente, por decidir).

`.doctrine.md` ganha por uma sinergia que nenhuma outra candidata tem: **o `anchors.yaml`
já chama sua camada de réguas de `doutrina`**, com a razão escrita — *"não são 'docs':
são a fonte das regras que o CLI implementa, e por isso governam o código"*. É
exatamente o papel do artefato novo, um nível acima.

Fica coerente de ponta a ponta: a doutrina do FRAMEWORK rege o código do framework; a
doutrina do PRODUTO rege as specs do produto. Mesma palavra, mesmo papel, escopos
diferentes — a mesma relação que `## Rules` (local) tem com o que este artefato faz
(transversal).

Também considerados: `.policy.md` (livre e legível, mas sem raiz no vocabulário
existente) e `.contract.md` (colisão forte — `Contract` aqui significa ASSINATURA).

**Já ocupados e por isso fora:** `domain` (seção de spec E tag de camada), `contract`
(seção de spec), `feature` (o Gherkin).

### O artefato

`product/<nome>.doctrine.md`, com `kind: product` (o oitavo), header `@anchors` com código
próprio, e regras catalogadas no mesmo formato das specs — mesma gramática de código,
mesmas letras.

### A aresta

`realizes`: a REGRA DA SPEC aponta para a REGRA DE PRODUTO que ela concretiza.

Direção importa. A spec é quem sabe que está realizando uma regra transversal; o
artefato de produto não pode listar seus realizadores sem virar um índice que envelhece a
cada spec nova — o mesmo defeito que o mapa existe para não ter.

### Os gates

| gate | pergunta | veredito |
| --- | --- | --- |
| `product-realized` | toda regra de produto é realizada por ao menos uma spec? | regra de produto que ninguém realiza é decisão escrita e não implementada |
| `realizes-resolves` | toda `realizes` aponta para uma regra que EXISTE? | o `ref-resolves` do eixo vertical |
| `product-not-duplicated` | a spec COPIA o texto da regra em vez de referenciá-la? | é o defeito que o conceito existe para eliminar |

## As decisões, tomadas

### D1 — Referência obrigatória? CONFIGURÁVEL POR CAMADA, opcional por padrão

Medido neste repositório: das 841 regras, a esmagadora maioria é local à unidade.
`SBGRD-B01` ("artefato que não é código sai sem veredito") não pertence a produto nenhum
— é mecânica de um gate. Obrigatório universal forçaria inventar doutrina guarda-chuva
para calar o gate, que é o vício que o `placeholder_filled` existe para pegar: preencher
o campo sem que o valor signifique nada.

Numa aplicação de produto a proporção se inverte — quase toda regra de tela serve a uma
decisão de produto, e a que não serve é suspeita. Quem sabe qual é o caso não é o
Anchors; é a Estrutura.

    doctrine:
      require_for: [screen, handler]   # vazio (padrão) = nenhuma camada exige

### D2 — Onde vive a referência? TAG NA REGRA

    ### CRED-V01 — limite de crédito respeitado    @realizes LIMIT-R03

Sobrevive a mudança de formato de tabela, e vale nas TRÊS formas de regra catalogada
(cabeçalho, linha de tabela, bullet-negrito) — uma coluna só existiria na do meio, e a
spec teria de trocar de formato para poder referenciar.

### D3 — Tem tríade? NÃO, E NÃO TEM PARIDADE

É 1 para MUITOS: uma regra de doutrina é realizada por várias specs, e é justamente isso
que o conceito existe para permitir.

A consequência importa para os gates: nenhum deles pode cobrar paridade de contagem.
`triad-complete` conta 1:1:1 entre spec, feature e teste; `product-realized` conta
"≥ 1 realizador", que é uma pergunta diferente. Quem prova a regra é a tríade de cada
spec que a realiza — duplicar teste aqui seria provar duas vezes a mesma coisa.

### D4 — A tag `@feature` migra? SIM, e ela hoje NÃO FAZ NADA

Verificado antes de decidir: `@feature` aparece em `internal/initx/header_guide.go` e
`cmd/anchors/governance/guide_header.go`, sempre em TEXTO DE GUIA. Nenhum parser a lê —
`internal/scan/` e `internal/mapx/` não a mencionam.

É promessa documentada e não cumprida: o guia diz que ela "agrupa arquivos da mesma fatia
de domínio, mesmo espalhados por camadas", e nada no motor agrupa coisa alguma.

Então não é migração de funcionalidade — é o artefato de produto CUMPRINDO o que a tag
prometia. `@product: auth` passa a apontar para `product/auth.doctrine.md`, que existe e
é confrontado.

Como nada lê a tag antiga, a migração é de TEXTO (os dois guias), e o passo de formato
em `internal/migra/` serve para reescrever os headers de projetos que já a adotaram —
sem ele, quem escreveu `@feature: auth` fica com um rótulo órfão.

## A implementação

1. `kind: product` (o oitavo) e a pasta `product/` na Estrutura semeada pelo `init`
2. aresta `realizes`, lida da tag `@realizes` pelo `scan`
3. guide: `anchors guide product` + `product/PRODUCT_GUIDE.md` semeado
4. `anchors new product <nome>` — o esqueleto com seções próprias
5. os três gates: `product-realized`, `realizes-resolves`, `product-not-duplicated`
6. `doctrine.require_for` na config, com o gate `spec-realizes-product` por camada
7. migração `@feature` → `@product` (texto dos guias + passo de formato)
