# anchors — CLI

CLI único do framework Anchors, em Go. Exercita o ciclo de vida do Anchors e é
validado contra uma prova de conceito real. Decisões de arquitetura em [`DECISIONS.md`](./DECISIONS.md).

## Estado

Em construção. Primeiro comando funcionando: **`anchors map build`** — a lacuna
central (o mapa de dependências, Tema A da simulação).

## Build & uso

```sh
go build -o anchors ./cmd/anchors

# construir o mapa de um projeto
./anchors map build --root /caminho/do/projeto
# → escreve <root>/anchors.graph.yaml
```

## Validado contra uma prova de conceito real

`anchors map build --root <projeto>` produz (lendo só texto, sem parsear código) um
grafo de porte substancial, com milhares de nós e arestas distribuídas entre as
relações `specifies` (spec → código, co-location), `covered-by` (spec → feature,
co-location), `tested-by` (feature → teste, co-location) e `references` (por código
de cenário — identidade).

As arestas conferem com a realidade (ex.: a trinca `LoginScreen.spec.md` →
`.tsx`/`.feature`, e `.feature` → `.test.tsx`).

## Arquitetura

```
cmd/anchors/        # comandos (cobra): root, map
internal/
  scan/             # percorre o repo lendo TEXTO; extrai código de cenário (regex)
  mapx/             # modelo do grafo, build (co-location + identidade), store (YAML)
  config/           # tabela de marcadores de comentário por linguagem (D4)
```

Princípios (DECISIONS.md): CLI único (D1); só lê texto, nunca parseia código (D2);
Go pelo mérito de mercado (D3); marcadores de comentário configuráveis (D4); lógica
dos validadores reimplementada internamente, não subprocessos (D5).

## Roadmap (caminho crítico do ciclo de vida)

- [x] `guide` — a ponte IA↔Anchors (a inversão): a IA (no CLI dela) roda
      `anchors guide`, lê o playbook embutido (o fluxo planejar → especificar →
      mapear → implementar → testar → confrontar + comandos + o que reportar) e
      opera o Anchors como ferramenta. Sem MCP, sem protocolo novo — a IA já tem
      shell; o guia vale para qualquer cliente de IA e casa com a versão do binário.
      Subcomandos imprimem as réguas de cada artefato — a doutrina UNIVERSAL,
      agnóstica de stack (o projeto especializa no seu guide da camada 'guide'):
      `anchors guide plan` (fase PLANEJAR), `guide spec` (origem da verdade),
      `guide code` (dois eixos + dependência unidirecional + promoção sob demanda),
      `guide feature` (cenários auto-contidos + cobre o catálogo da spec) e
      `guide test` (pirâmide + testar o código real, não o mock). Destilados dos
      guides reais de uma prova de conceito, separando o universal do dialeto. NB: o arquivo do
      subcomando 'test' é `guide_tests.go` (o sufixo `_test.go` é reservado pelo Go).
- [x] `init` — configura o projeto (anchors.yaml) por P&R; inferência determinística
      (sem IA) + perguntas só para as decisões humanas (co-location, granularidade
      das camadas, governs guide↔tag)
- [x] `map build` — construir o mapa (Rastreabilidade)
- [x] `impact <arquivo>` — análise de impacto (Propagação): CONSULTA, duas
      dimensões — ↓ propaga (desce, para em `@noPropagation`) e ↑ valida (sobe, não
      re-propaga; divergência viraria issue). Validado numa prova de conceito real,
      com árvores de impacto de tamanhos variados (de poucas arestas a centenas).
- [x] `check` — o pipeline de gates (Qualidade). INCREMENTAL por padrão
      (`--changed <arquivo>`, sobre o caminho de impacto) ou `--all` (foto completa).
      Gates declarados no anchors.yaml (seção `gates:`), internos (CLI lê texto) ou
      externos (invoca comando, lê exit code). Perfil de vereditos + issues; exit 1
      se fail bloqueante. Validado numa prova de conceito real: `--all` achou
      dezenas de issues reais (a maioria specs incompletas, o restante código sem spec).
      A saída é ESPELHADA em `.anchors/check-all.txt` (ou `check-changed.txt`), com
      cabeçalho de quando/HEAD/árvore: o `--all` custa minutos, e reler não pode
      exigir re-executar (D6). `--only-issues` enxuga a tabela para os gates com
      achado ou pendência, mantendo o total dos omitidos no rodapé (D7);
      `--show-drift` lista TODAS as pendências com endereço (D8).
- [x] `doctor` (alias `status`) — o validador de saúde do ecossistema (QUALITY §5.2):
      visão GLOBAL, caça pontas SISTÊMICAS que o check não vê — integridade do mapa
      (arestas mortas, nós fantasma), órfãos (código sem spec, identidade ausente),
      camadas frouxas, buracos de cobertura de gate. Apresenta e registra, NÃO
      bloqueia. Validado numa prova de conceito real: achou identidades ausentes,
      um kind sem gate, e vários guides sem governo declarado (info).
- [x] `map show` — consulta o mapa: `<arquivo>` (vizinhança ↑regido/↓propaga),
      `--orphans` (nós sem arestas), `--stats` (nós por kind, arestas por tipo).
- [x] `watch` — o WATCHER (PROPAGATION §6), em BACKGROUND. `start` daemoniza
      (re-exec com Setsid, PID em `.anchors/watch.pid`, log em `.anchors/watch.log`)
      e retorna o terminal na hora. Subcomandos: `status`, `stop`, `pause`/`resume`
      (via flag `.anchors/watch.paused`), `logs`. O loop observa (fsnotify), filtra
      pelos patterns da Estrutura, e a cada mudança encadeia classifica→impacto→check
      incremental, REPORTANDO no log (não grava carimbo nem abre issue ainda).
      Validado numa prova de conceito real: start não-bloqueante, status/pause/resume/stop ok;
      cascata reporta corretamente a propagação nos dois sentidos + gates.
- [x] **fila de tasks** — o watcher deixou de só reportar: agora ENFILEIRA. A cada
      mudança, classifica e escreve uma task em `.anchors/tasks/` com o próximo passo
      sugerido (spec→implement, feature→test, …). Comandos: `anchors queue` (lista),
      `anchors next` (puxa+reivindica, claim atômico via rename — dois terminais
      nunca pegam a mesma task), `anchors done <id>` (move p/ `.anchors/done/`). É o
      que desacopla "mudou" de "alguém trabalha": a conversa da IA fica livre; o
      worker (subagente em bg, ou sessão bloqueante com aval do usuário) puxa da fila.
      Validado numa prova de conceito real: watcher enfileira ao salvar spec → next
      → done, ciclo completo.
- [x] **loop check→carimbo→issue** — o `check` deixou de só reportar: agora REGISTRA.
      (1) CARIMBA o mapa — grava `Stamp{from_rev,to_rev,verdict}` nas arestas cujas
      duas pontas foram confrontadas (`mapx.StampEdges`), o que destrava a detecção de
      stale. (2) ABRE issues — uma `violation` por fail bloqueante em `issues/todo/`
      (pkg `internal/issue`, kinds stale/conflict/violation, estados todo/doing/done
      por pasta, imutáveis — CONCEPT §5). Idempotente: reconfrontar não duplica nem
      ressuscita issue em doing/done. Opt-out honesto: `--no-record` só reporta.
      Validado numa prova de conceito real: `check --all` carimbou centenas de arestas
      + abriu dezenas de issues reais; editar uma spec fez suas arestas voltarem a
      stale sozinhas; re-check não duplicou.
- [x] `stale` — lista as arestas stale (nunca validadas vs. drift de rev), agora com
      significado real porque o carimbo existe.
- [x] **issue fecha sozinha (loop completo abre+fecha)** — a issue tem uma Key ESTÁVEL
      no tempo (kind+gate+aresta, sem data), então o `check` liga o pass de hoje à
      issue de ontem: fail bloqueante ABRE, pass RESOLVE (move todo/|doing/ → done/).
      A issue deixou de mentir sobre o estado. Validado numa prova de conceito real:
      spec quebrada abre issue; corrigir a spec fecha sozinha no check seguinte.
- [x] **watcher reage a pastas novas** — antes o fsnotify só observava os dirs do walk
      inicial; um arquivo numa pasta nascida após o `watch start` (ex.: plans/) nunca
      gerava evento. Agora, ao ver Create de um dir, o watcher o registra
      recursivamente E varre os arquivos que já nasceram nele (cobre a corrida
      Create-dir/Create-arquivo). Validado numa prova de conceito real no pior caso
      (mkdir + arquivo sem pausa): a spec dentro da pasta nova enfileira `implement`.
- [x] **higiene da fila** — `anchors drop <id>` descarta uma task (lixo: triage,
      obsoleta) sem arquivar; `anchors reclaim` devolve a pending as tasks claimed
      órfãs (worker morto). `anchors queue` sinaliza quantas estão claimed/triage com
      a dica de comando. Fecha o atrito da fila poluída (achado no exercício C).
- [ ] `issue` (subcomandos take/resolve manuais — mover entre todo/doing/done) e
      `plan` (registro e origem — a 3ª rota: issue→plano).
- [x] **falso-positivo de propagação (#5 de C) — diagnosticado na raiz.** A aresta
      espúria (spec propaga para um teste sem relação) era SINTOMA de códigos de
      cenário DUPLICADOS entre unidades. Novo check `identidade-duplicada` no `doctor`:
      só kinds-donos (spec/feature/test/code) contam (doc/guide/plan apenas
      referenciam), agrupa por UNIDADE (stem co-locado, não diretório), e separa
      cross-domain (⚠ colisão real que cruza features) de intra-feature (ℹ lib+screen
      deliberado). O `doctor` passou a marcar cada item por severidade (grupo é ⚠ se
      houver qualquer ⚠). Achou várias colisões reais latentes numa prova de conceito
      real. Corrigir os códigos é trabalho de produto; o framework agora TORNA
      VISÍVEL o que estava oculto.
- [x] **`anchors code` — gerador/validador de identidade (#7 de C).** Fecha a causa
      do #5: a IA agora obtém um código ÚNICO antes de escrever (não colide como o
      `SPCR` do exercício). `internal/code` porta a doutrina do SPEC_GUIDE de uma
      prova de conceito real: Camada 2 (compressão do nome — consoantes 2+2 / vogais
      / X-padding, resultado quase idêntico ao original sem depender de silabação) e
      Camada 3 (resolução de colisão determinística). `anchors code <nome>` sugere
      livre contra o mapa; `--check <cód>` valida (exit 1 se tomado). Camada 1
      (prefixo de módulo) LIGA na Estrutura: layer ganha `code_prefix`, e um caminho
      que cai nela usa o prefixo do módulo. Validado na prática: `code Spacer` →
      SPCA (evitou o SPCR ocupado).
- [x] **presets de estrutura por stack no `init`** — catálogo de 17 estruturas
      consagradas (`internal/initx/presets.go`), destilado de um estudo das convenções
      por linguagem/framework (docs oficiais + comunidade): node-ts/nest, express,
      nextjs, angular, nuxt, expo-rn, spring, dotnet-clean, laravel, rails, phoenix,
      flutter, cpp, python-lib, django, fastapi, go, rust. O `init` OFERECE um preset
      (menu huh) que preenche as layers de código; presets modulares DEDUZEM o prefixo
      de identidade de cada módulo real (`DeduceModulePrefixes`, com unicidade —
      auth/audit → AT/AD), ligando a Camada 1 do `anchors code` à Estrutura. Lógica
      pura testada (catálogo bem-formado, ToLayers, dedução determinística+única);
      validado gerando o anchors.yaml de um preset (expo-rn) de ponta a ponta.

- [x] **gate de julgamento por IA (QUALITY §5.2)** — o segundo medidor, o subjetivo.
      Guides impõem regras não-scriptáveis ("esta tela quebra em atomic design?");
      isso é julgamento, não regex. Um gate com `measures: judgment` (+ `guide`, `ask`,
      e `tags:` para escopar) NÃO é computado pelo CLI: o `check` marca os alvos como
      `⏳ julgamento` e enfileira uma task `judge`. A IA que opera o Anchors lê o guide,
      confronta o alvo, e reporta com `anchors judge <alvo> --gate <g> --verdict
      pass|fail --reason ...`. O CLI faz a MESMA dupla saída de um gate determinístico:
      carimba a aresta guide→alvo no mapa (`StampEdge`) e abre/resolve a issue. Como o
      carimbo leva a rev do alvo, o veredito de IA ENVELHECE (fica stale) se o alvo
      mudar — mesmo anti-drift. O CLI nunca invoca IA (mantém D3/inversão): é só o
      verbo + a contabilidade; quem julga é a IA, em qualquer cliente. Validado numa
      prova de conceito real: um gate de julgamento enfileirou dezenas de telas;
      julguei uma delas (fail com razão real → issue; depois pass → issue resolvida
      + carimbo ok). `anchors judge
      --pending` lista o que aguarda. Node ganhou `tags` (para gates escopados por tag).
      O `--reason` é o LAUDO COMPLETO (markdown multi-linha), não uma frase de veredito:
      como a IA já leu o guide e o alvo inteiros para julgar, ela emite numa passada só
      cada não-conformidade com o-quê/onde/porquê/como-corrigir — isso vira o corpo da
      issue, sem reprocessar o alvo depois (economia de token). O cabeçalho fixo da
      issue é suprimido quando o laudo já traz seus próprios `##`. Playbook e `judge
      --help` instruem a IA a entregar o laudo, não o veredito seco.

- [x] **`anchors governs`** — o mapa responde "quem cada guide rege e quantos". Sem
      arg: o quadro (guide → nº de nós regidos direto, + total = tamanho de uma
      auditoria por guide). Com um guide: os arquivos regidos, por kind. Dimensiona e
      fatia a auditoria de julgamento por guide, e expõe redundância (numa prova de
      conceito real: vários guides de frontend regem os MESMOS arquivos, em número
      substancial — candidatos a afinar por tag). Query pura `Governs`/`GovernanceSummary`
      no mapx, testada.
- [x] **guide-sem-governo virou Warn** — "um guide que não rege ninguém não é guide,
      é doc". O `doctor` subiu esse achado de Info→Warn e a mensagem ensina a saída:
      declarar a regra `governs` (a tag que ele rege) ou reclassificar como doc. (Numa
      prova de conceito real: vários guides sem governança declarada — regiam
      test/feature na prática, mas o anchors.yaml nunca escreveu a regra. Débito de
      Estrutura, não dos guides.) NB: fica no `doctor` (saúde do grafo), não no `check`
      (gates por alvo, que só veem conteúdo, não o grafo).

- [x] **pontos de conformidade + meta-guide (julgamento menos heurístico)** — um guide
      passa a destilar suas regras numa seção `## Pontos de conformidade` (itens CK1,
      CK2…), e a IA julga o alvo CONTRA CADA PONTO, não contra a prosa "no olho" —
      focado, objetivo, reproduzível. Quando o guide rege targets diferentes (telas E
      componentes), os pontos são AGRUPADOS POR TARGET (`### Para <camada>`, códigos
      SCR-CK1/CMP-CK1…) e a IA aplica só o grupo da camada do alvo + "Para todos". Três
      peças: (1) `anchors guide guide` — o meta-guide (como escrever um guide, com a
      seção obrigatória e o agrupamento por target); (2) checker DETERMINÍSTICO
      `guide-has-checklist` — a PRESENÇA da seção é computável, então um gate barato
      reprova o guide sem checklist antes de qualquer julgamento (validado numa prova
      de conceito real: todos os guides existentes reprovam; adicionar a seção faz
      passar); (3) o playbook de JULGAR
      instrui a aplicar os CKs da camada do alvo. Divisão limpa: a presença da checklist
      é determinística (check); a qualidade de cada ponto e o veredito são julgamento.

- [x] **confiança no entregável: sinais de teste ingeridos (as 4 iniciativas)** —
      o Anchors deixa de só verificar que o teste EXISTE e passa a medir sua QUALIDADE,
      SEM rodar o runner (consome o artefato que o projeto gera — mantém D3). Pacote
      `internal/testsig`: (1) **execução** — parse de JUnit XML (`--junit`), grava
      passou/falhou por nó de teste (gate `tests-pass`); (2) **cobertura de linha** —
      parse de lcov (`--lcov`), % por arquivo de código (gate `line-coverage`, limiar
      70%); (3) **cobertura por CENÁRIO** (o diferencial) — cruza os códigos de cenário
      da spec com os que aparecem em casos que PASSARAM: cada requisito tem teste verde?
      (gate `scenario-coverage`); (4) **ajuda por stack** — cada preset tem um
      `CoverageHint` (o comando de coverage + formatos por stack: jest→lcov+jest-junit,
      go→coverprofile+go-junit-report, pytest→--cov+--junitxml…), que o `init` mostra.
      Comandos `anchors ingest` (grava os sinais no mapa) e `anchors coverage` (lê:
      cenários sem prova + linha < limiar). O sinal leva a rev do nó → fica STALE se o
      arquivo mudar (não confia em teste de versão antiga). Validado numa prova de
      conceito real: ingeri JUnit+lcov reais de uma spec real → cobertura por cenário
      revelou requisitos declarados SEM teste verde, invisíveis a coverage de linha;
      os 3 gates reprovam no `check`.

- [x] **sinais de teste amarrados ao CICLO DE VIDA** — os gates de cobertura deixaram
      de ser config manual e passaram a nascer com o projeto e a disparar no fluxo:
      (1) `SuggestNext("test")` virou `verify-tests` — quando um teste muda, o watcher
      enfileira uma task que instrui RODAR a suíte + `anchors ingest` (junit/lcov) +
      `anchors check`, não só "confronte"; (2) `initx.DefaultGates` — o `init` semeia os
      gates padrão conforme os artefatos escolhidos (spec→spec-completa/tem-codigo;
      test→tests-green/line-coverage; spec+test→scenario-coverage; guide→guide-checklist),
      todos INFORMATIVOS (maturação §7), com confirmação; (3) o playbook (`anchors
      guide`) descreve o passo verify-tests e aponta o CoverageHint do stack. Validado
      numa prova de conceito real: mudar um teste enfileira `verify-tests` com a
      instrução de ingerir; os gates de ciclo (tests-green, scenario-coverage)
      reprovam de verdade após ingest.

- [x] **cobertura do DIFF + delta (o bug que entra AGORA)** — cobertura absoluta não
      basta: um arquivo 90% coberto pode ter a SUA mudança nos 10% descobertos. Duas
      perguntas novas, sobre a mudança: (1) **o que mudei está coberto?** — `anchors
      coverage --diff <ref> --lcov <c>` cruza o `git diff --unified=0` (linhas
      mudadas, parser de unified diff em `testsig`; fallback `--diff-file` p/ quem não
      usa git — agnóstico) com as linhas cobertas do lcov; só linhas INSTRUMENTADAS
      contam; exit 1 se abaixo do limiar. (2) **a cobertura caiu?** — `anchors coverage
      --delta`: o `Signal` preserva a cobertura anterior (`PrevLineCoverage`) ao
      ingerir por cima, então o próprio grafo é o baseline (zero armazenamento novo); o
      gate `coverage-delta` reprova regressão. Validado numa prova de conceito real:
      inseri uma linha real numa tela existente → `--diff` apontou a linha sem teste
      (0%, exit 1); duas ingestões sucessivas com queda de cobertura → `--delta`
      pegou a regressão (exit 1). Ligado ao ciclo: o passo
      verify-tests do playbook manda rodar `--diff`/`--delta`, e `DefaultGates` semeia
      `coverage-delta`.

- [x] **relatório de confiança em docs/ (mergeando camadas)** — `anchors report` gera
      `docs/anchors-test-report.md`, um painel de confiança do entregável consolidado
      do GRAFO (não reparseia o JUnit; lê os Signal ingeridos): execução por CAMADA,
      cobertura por cenário (requisitos provados), cobertura de linha, regressões. Apps
      com suites em camadas distintas (unit/integration/e2e, por exemplo) MERGEIAM num
      painel único — `anchors ingest --junit r.xml --layer <camada>` acumula por camada no
      grafo (não sobrescreve; `Signal.ByLayer`), e o report soma e rotula. É um DOC
      TERMINAL (consumo humano, versionado no git para dar história) — NÃO vira âncora
      nem é regido (o Anchors não governa o próprio output). Honesto sobre o que não
      mediu: distingue "sem prova" (ingerido, teste não cobre) de "não medida" (nenhuma
      ingestão a tocou) — ausência de prova ≠ prova de ausência. Validado numa prova
      de conceito real: ingeri camadas distintas (unit + e2e) → painel soma e rotula
      os resultados corretamente por camada.
- [x] **`report` é uma FAMÍLIA por perspectiva** — não um painel único. Subcomandos,
      cada um um recorte de fontes que já medimos (nenhum inventa dado): `report tests`
      (o painel de confiança), `quality` (vereditos dos gates + débito), `structure`
      (camadas/governança/colisões/órfãos), `config` (estado do anchors.yaml e o que
      falta), `issues` (issues todo/doing/done + tasks da fila — o débito ADOTADO),
      `inconsistencies` (a lista infinita a arrumar — tudo que os validators acham e
      ainda NÃO virou issue). A distinção issues≠inconsistências é deliberada: issue =
      débito reconhecido e rastreado; inconsistência = débito detectado e solto (vira
      issue quando triado). `report all` gera todos + um índice em docs/anchors/.
      Validado numa prova de conceito real: `report all` → todos os relatórios com
      dados reais (quality: dezenas de reprovações bloqueantes; inconsistencies:
      dezenas encontradas pelo health).

## Notas de implementação (achados)

- O código de cenário é extraído por regex idêntico à gramática do
  `scripts/lib/scenario-code.ts` de uma prova de conceito real (fonte única da
  identidade, TRACEABILITY §3).
- `map build` hoje cobre as origens `convention` (co-location) e `inferred` (por
  código de cenário). A origem `declared` (aresta escrita à mão) e o carimbo de
  validação por aresta (`stamp`) ainda não são preenchidos — entram com os comandos
  de propagação e gates.

## A Estrutura declara; o CLI lê

O `map build` NÃO hardcoda nada sobre a estrutura do projeto. Ele lê o
`anchors.yaml` (a Estrutura de Projeto — o grafo virtual, STRUCTURE §2.1):

- **`layers`** — cada camada tem um `pattern` (glob que a reconhece), um `kind`, e
  **`tags`** (rótulos de agrupamento transversais). O CLI classifica cada arquivo
  pela layer cujo pattern casa (mais específico vence).
- **`derived`** — a co-location: os templates (`{{dir}}/{{name}}.spec.md`…) que
  ligam a trinca de um target. A âncora casa por `kind`.
- **`governs`** — a dimensão vertical, **sempre por tag**: `{from: <guide>,
  governs: <tag>}`. O guide rege os nós de todas as layers com aquela tag. O escopo
  vem dos patterns das layers (DRY — nenhum glob duplicado). Sem produto cartesiano:
  cada guide toca só sua tag. A precisão (SCREEN só telas vs. frontend inteiro) é
  uma decisão declarativa de tags, não código.

## Furos fechados vs. abertos

Fechados ao mover a Estrutura para o `anchors.yaml` + modelo de tag:
- ✅ **references infladas** — a regra por código de cenário agora só liga
  cross-target e não-duplicado, com o tipo da relação (`tested-by`), não references.
- ✅ **convenção hardcoded / só a stack da prova de conceito** — tudo vem dos
  `layers`/`derived` do `anchors.yaml`; o CLI é agnóstico.
- ✅ **dimensão vertical (`governs`)** — existe, por tag, com escopo derivado das
  layers.

Ainda abertos (a refinar):
- **`rev` = content-hash (sha256), não git.** Revisitar se "revisão" deve vir do git.
- **carimbo de validação (`stamp`) não preenchido** — entra com os gates.
- **`governs` só cobre guides→código/spec.** plano→spec e a spec-de-arquitetura
  como régua ainda não estão declarados. Precisam do `plan`/da camada arquitetura.
- **precisão das tags na prova de conceito** — hoje vários guides frontend regem
  todo o código mobile (correto, mas grosso); afinar com tags `screen`/`component`
  é ajuste de YAML.
