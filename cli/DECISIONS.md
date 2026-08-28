# Anchors CLI — decisões e levantamento

Registro das decisões de arquitetura do CLI do Anchors e o levantamento de comandos
que as embasou. O objetivo desta fase: **um CLI único que exercita o ciclo de vida
do Anchors, validado contra uma prova de conceito real** — antes de investir numa IDE.

---

## Decisões

### D1 — Um CLI único, não uma coleção de scripts
Todo o comportamento vive num só programa (`anchors`), com subcomandos. Não é um
wrapper que dispara scripts externos. A lógica dos validadores existentes (da prova
de conceito) é **referência**, reimplementada como módulos internos do CLI, unificada
sobre um único modelo de rastreabilidade — não dependência.

### D2 — O CLI só lê TEXTO; nunca parseia código
O Anchors não entende a sintaxe de nenhuma linguagem. Ele lê:
- **artefatos de texto** (spec, feature, doc — markdown/gherkin);
- **anotações que vivem em comentários** dentro do código (código de cenário, tags de
  regime/nível, marcações de âncora).
Para achar as anotações, o CLI só precisa saber **qual é a sintaxe de comentário de
cada linguagem** (`//`, `#`, `--`, `/* */`, `<!-- -->`…). Nada de AST, nada de
semântica. Ler texto + tabela de marcadores + filesystem/regex.

### D3 — Linguagem: Go
Escolhida **pelo mérito de mercado**, não por alinhamento com alguma ferramenta específica:
- **binário único estático, multiplataforma**, sem runtime — o usuário baixa e roda;
- **neutro de ecossistema** — instala via brew/scoop/curl, não via um gerenciador de
  pacote de linguagem (um CLI via `npm` sinalizaria "isto é para gente de JS", ruim
  para adoção fora do JS);
- padrão-ouro de CLI de infra (terraform, gh, hugo, kubectl).
O Anchors é para ser bom **independente** do que já existe — pensando em quem vai
adotar, inclusive quem não usa Node nem nenhuma ferramenta correlata que por acaso
também seja em Go (esse eventual alinhamento seria bônus lateral, não a razão).

### D4 — Marcadores de comentário: tabela configurável + defaults embutidos
Defaults por extensão embutidos (`.ts`/`.go` → `//`; `.py`/`.rb` → `#`; `.sql` →
`--`; `.html` → `<!-- -->`…), e o projeto pode estender/sobrescrever. Agnóstico de
verdade: cobre as linguagens comuns out-of-the-box, acomoda as exóticas.

### D5 — Gates: lógica reimplementada, não subprocessos
Os validadores da prova de conceito (validate-specs, check-arch, cobertura
cenário↔teste) viram módulos internos do CLI, unificados. O `scenario-code.ts`
(fonte única do código de cenário) é a referência da identidade — vira parte do
modelo de rastreabilidade do CLI.

### D6 — O `check` espelha a própria saída num arquivo
Toda execução de `check` grava o relatório em `.anchors/check-all.txt` (com `--all`)
ou `.anchors/check-changed.txt` (incremental), além de imprimi-lo na tela.

O motivo é medido, não estético: o `check --all` de um projeto real leva minutos.
Sem o espelho, toda pergunta sobre o resultado — "o que reprovou?", "quais gates
passaram?" — obriga a rodar tudo de novo só para reler o que já saiu. Num terminal
com scroll limitado, ou numa sessão de agente onde a saída rolou para fora do
contexto, a informação simplesmente se perde. Foi o que aconteceu numa varredura de
663 unidades: o mesmo `--all` rodou várias vezes seguidas sem nada ter mudado no
código entre elas.

Três decisões dentro dela:

- **Separado por escopo.** O `--changed` roda a cada commit (pre-commit) sobre um ou
  dois arquivos; o `--all` é a foto completa e cara. Num arquivo só, o primeiro
  commit depois de uma varredura apagaria justamente a foto que custou minutos.
- **Cabeçalho com HEAD e estado da árvore.** Sem eles o arquivo mente por omissão:
  uma leitura de ontem parece a de agora. Com árvore suja, o relatório descreve um
  estado que não está em commit nenhum, e quem relê precisa saber que não vai
  reencontrá-lo pelo hash.
- **Falhar ao escrever não derruba o check.** O espelho é conveniência; abortar a
  varredura porque o disco encheu seria trocar o essencial pelo acessório.

Fica em `.anchors/` (estado efêmero, junto do daemon), não em `issues/`: é cache de
leitura, não conteúdo de projeto.

### D7 — A tabela de gates alinha por coluna, e o gate limpo só some sob flag
A largura do nome vem do maior nome presente (era `%-20s` fixo, e cinco gates passam
disso — `handler-ddb-inline-passivo` tem 26); cada coluna de contador tem a largura
do maior número DAQUELA coluna; e a célula do `⚠` é reservada mesmo vazia.

O `⚠` reservado é o ponto não-óbvio. Omiti-lo nas linhas sem drift empurrava o `~`
para a esquerda só nelas — a mesma coluna existindo em dois lugares na mesma tabela,
que é pior que nenhum alinhamento: o olho desce a lista comparando números que não
estão na mesma vertical. O branco ali não é desperdício, é o que segura a coluna.

Largura por coluna, e não uma única para todas: numa varredura real o `~` chega a 582
e o `✗` fica em 0 ou 1. Uma largura comum obrigaria a coluna dos fails a reservar três
casas para nada.

`--only-issues` omite os gates que passaram em tudo e não deixaram pendência, com o
total no rodapé. É OPT-IN: por padrão a tabela cheia continua, porque ela é a prova de
que os 49 gates rodaram — esconder isso por padrão trocaria a prova por brevidade. O
rodapé existe pelo mesmo motivo: sem o número, a saída enxuta pareceria uma varredura
menor, e não a mesma varredura com a listagem curta.

A coluna do `⚠` existe ou não pela TABELA INTEIRA, não por linha: se algum gate tem
drift ela aparece em todas (vazia vira branco, que é o que segura a coluna seguinte na
vertical); se nenhum tem, ela não existe em lugar nenhum — reservá-la deixaria um
buraco no meio de todas as linhas sem nada que o justificasse, e a tabela sem drift é
o caso comum do `--changed`, não a exceção.

A legenda dos símbolos lista SÓ os que a tabela usou. Uma legenda fixa explicaria `⚠`
e `⏳` em varreduras que não os têm, e explicar coluna ausente é o que ensina a pular a
legenda inteira. A distinção que ela precisa carregar é entre `✗` e `~`: falha é o
gate ter confrontado e divergido; indeterminado é ele não ter tido o que confrontar —
sem isso, `~582` parece um débito de 582 itens, quando é a medida do que não se aplica.

Uma entrada por LINHA, não todas na mesma: com cinco símbolos a linha única passa de
100 colunas e quebra sozinha num terminal estreito — o oposto do que uma legenda serve.
Empilhadas, elas também se comparam entre si, que é como se lê uma legenda.

### D8 — O drift tem bloco de detalhe próprio, e não herda o corte de ruído do `~`
O `⚠` foi separado do `~` porque é a única categoria acionável do balde (D-anterior),
mas os dois compartilhavam a lista de motivos — e essa lista é suprimida acima de 40
resultados, para o `--all` não despejar milhares de linhas de "não se aplica".

O efeito era o pior possível: a tabela anunciava `⚠407` e o detalhe abaixo não trazia
nenhum. Medido no app de referência: 2.430 drifts contados, zero listados. Contar uma categoria como
acionável e depois escondê-la é pior que não a separar — o número vira uma acusação sem
endereço, e o endereço é justamente a parte sobre a qual se pode agir.

Agora o drift tem bloco próprio, sob `--show-drift` — e quando pedido, lista TODAS,
sem teto.

O teto que existiu por um momento (25 itens, com "… e mais 2405") devolvia o problema
que a lista existe para resolver: quem pede os endereços quer agir sobre eles, e 2405
sem endereço é o mesmo que nenhum. Ou a lista é completa, ou o contador da tabela já
bastava.

Por isso o detalhe é OPT-IN e o contador é o padrão: numa varredura grande são milhares
de linhas, e despejá-las sem ninguém pedir enterraria as issues logo abaixo — que são o
que barra a promoção. A tabela continua anunciando `⚠407`, e a flag é o caminho do
número para os endereços. O `~` mantém a supressão por volume: ali o contador basta
mesmo, porque "não se aplica" não é acionável.

---

## Levantamento: o que a prova de conceito já prova vs. o que falta

Cruzamento de dois levantamentos (doutrina dos pilares × scripts reais da prova de
conceito).

**A prova de conceito tem os GATES (a metade "confronta"), mas não o MAPA nem a
PROPAGAÇÃO (a metade "movimento").** Os 4 validadores existem e funcionam, mas são
**fragmentados** — cada um relê a árvore inteira, sem um índice de rastreabilidade
unificado. O furo é exatamente o Tema A da simulação: o **mapa de dependências**, e
sobre ele o impacto, a propagação e o doctor.

Gates que a prova de conceito já prova (viram módulos internos do CLI):
- validação de spec (`validate-specs.ts`), de feature (`validate-features.ts`)
- cobertura cenário↔teste (`validate-test-coverage.ts`)
- arquitetura/camadas (`check-arch.sh`, 15 regras)
- execução de testes (jest mobile + backend), lint/spellcheck/format/circular
- identidade de cenário (`scenario-code.ts` — fonte única)

Lacunas (o CLI implementa novo):
- **mapa de dependências** persistente e consultável (`anchors.graph.yaml`)
- **análise de impacto** (dado o diff, o que revalidar)
- **propagação** (a onda; hoje os 4 validadores não propagam, só validam)
- **doctor/status** (validador de saúde; hoje não há panorama)
- **scaffolding** de spec/feature/teste; **issues** (pastas todo/doing/done)

---

## Inventário de comandos (candidatos)

Caminho crítico para exercitar uma volta do ciclo (origem → verdade → cola →
movimento → régua → registro → saúde):

| # | comando | pilar | fonte |
|---|---|---|---|
| 1 | `anchors plan new` / `execute` | Planejamento | novo |
| 2 | `anchors validate-spec` | Spec | reimplementa POC |
| 3 | `anchors map build` / `update` | Rastreabilidade | **novo — lacuna central** |
| 4 | `anchors impact <arq>` / `propagate` | Propagação | **novo — lacuna** |
| 5 | `anchors check` (pipeline de gates) | Qualidade | reimplementa POC |
| 6 | `anchors issue list/open/move/convert` | Issues | novo |
| 7 | `anchors doctor` / `status` | Validador de saúde | **novo — lacuna** |

Secundários: `classify`, `map show`, `stale`, `gate <nome>`, `gate promote`,
`spec no-test`, `orphans`, `code check`, `gate doc-freshness`. (Muitos [V] se
unificam sob `anchors gate <nome>` com manifesto declarativo — QUALITY §4.)

Convenção de tipo: validação (confronta/gate) · operação (cria/altera) · consulta
(lê/reporta).

---

## Primeiro alvo

**`anchors map build`** — gerar o `anchors.graph.yaml` a partir da prova de conceito
real, via co-location (nomes de arquivo) + código de cenário (regex sobre os
artefatos), sem parsear código. É a lacuna central; destrava impacto, propagação e
doctor. Prova o Tema A na prática. (A ser confirmado antes de codar o esqueleto.)
