# Estrutura do Projeto Anchors

> Este documento descreve a **planta da casa** do próprio repositório `anchors`,
> aplicando e materializando o pilar conceitual definido em [`STRUCTURE.md`](./STRUCTURE.md).
> Ele documenta a árvore de diretórios, as camadas declaradas em [`anchors.yaml`](./anchors.yaml),
> as regras de fronteira e a organização modular do CLI e dos pacotes internos.

---

## 1. Visão Geral e Arquitetura

O Anchors é um monorrepo em Go com ferramentas complementares para documentação e simulação:
- **CLI (`cmd/anchors/`)**: Ferramenta de linha de comando única baseada em Cobra, executada por desenvolvedores e agentes de IA.
- **Módulos Internos (`internal/`)**: Pacotes que implementam o núcleo do framework (leitura de texto, grafo de dependências, avaliação de gates de qualidade, telemetria, detecção de ambiente e utilitários).
- **Portal de Documentação (`site/`)**: Aplicação Astro + Starlight que compila o site bilíngue (Português e Inglês) a partir de `site/src/content/docs/`.
- **Scripts de Suporte (`scripts/`)**: Automações auxiliares para validação de licenças, segredos, SBOM, vulnerabilidades e conformidade estática.
- **Simulação (`simulation/larder/`)**: Aplicação fictícia de referência ("Larder") usada para validar e testar de ponta a ponta todos os fluxos do ciclo de vida do Anchors.

---

## 2. Camadas do Framework (`anchors.yaml`)

O repositório é governado pelas camadas declaradas em [`anchors.yaml`](./anchors.yaml):

| Camada | Padrão de Arquivo | Tipo (`kind`) | Tags | Responsabilidade |
|---|---|---|---|---|
| **`doutrina`** | `*.md` | `guide` | `[doutrina, guide]` | Os documentos canônicos que definem o framework (`CONCEPT.md`, `STRUCTURE.md`, `SPEC.md`, etc.). Governam as camadas de código. |
| **`doc`** | `{README,PLANNING,PROJECT_STRUCTURE}.md` | `doc` | `[doc]` | Documentação operacional, visão geral e mapa estrutural do repositório. |
| **`comando`** | `cmd/anchors/**/*.go` | `code` | `[cli, comando]` | Ponto de entrada do usuário/agente (Cobra CLI), dividido por domínios funcionais. |
| **`gate`** | `internal/gate/*.go` | `code` | `[regra, gate]` | Implementação dos gates de qualidade (segurança, integridade, rastreabilidade, trinca). |
| **`config`** | `internal/config/*.go` | `code` | `[nucleo, config]` | Configuração do projeto, marcadores de comentário por linguagem, resolução de raízes. |
| **`scan`** | `internal/scan/*.go` | `code` | `[nucleo, scan]` | Varredura do repositório lendo texto puro; extração de identidades e anotações. |
| **`mapa`** | `internal/mapx/*.go` | `code` | `[nucleo, mapa]` | Grafo de dependências (`anchors.graph.yaml`), construção de arestas e persistência. |
| **`infra`** | `internal/{...}/*.go` | `code` | `[infra]` | Utilitários de apoio: issues, mudanças, telemetria, daemon, i18n, health, etc. |
| **`teste`** | `**/*_test.go` | `test` | `[test]` | Testes de unidade e integração automatizados. |
| **`spec`** | `**/*.spec.md` | `spec` | `[spec]` | Especificações formais de requisitos e contratos. |
| **`feature`** | `**/*.feature` | `feature` | `[feature]` | Cenários executáveis em formato Gherkin. |

---

## 3. Estrutura do CLI (`cmd/anchors/`)

O diretório [`cmd/anchors/`](./cmd/anchors/) organiza os comandos do CLI por **Domínio / Ciclo de Vida**, evitando uma pasta plana de centenas de arquivos e agrupando comandos afins com seus respectivos testes:

```
cmd/anchors/
├── main.go                      # Ponto de entrada (início, saída, flush de telemetria)
├── root.go                      # Comando raiz (anchors) e registro dos domínios
├── root_test.go                 # Testes de integração do comando raiz contra guias e workflows
│
├── common/                      # Utilitários compartilhados entre comandos do CLI
│   ├── path.go                  # Normalização de caminhos (relTo)
│   ├── flags_compat.go          # Tratamento de flags e compatibilidade (aliasDeFlag)
│   ├── cards.go                 # Utilitários de cards e PRs
│   ├── roles.go                 # Perfis e papéis de trabalho
│   ├── telemetry.go             # Inicialização e flush de telemetria
│   └── spec.go                  # Utilitários de spec e headers
│
├── flow/                        # Ciclo de Execução e Motor de Trabalho
│   ├── work.go                  # anchors work: orquestração do trabalho ativo
│   ├── queue.go                 # anchors queue / next / drop / done: fila de tarefas
│   ├── watch*.go                # anchors watch: observação contínua de mudanças
│   ├── progress*.go             # anchors progress / progress merge: estado dos planos
│   ├── task_status*.go          # anchors task-status: telemetria de ciclo de tarefas
│   ├── deliver*.go              # anchors deliver: entrega e confronto com branch
│   ├── escalate*.go             # anchors escalate: escalonamento de bloqueios/dúvidas
│   ├── decided*.go              # anchors decided: resolução de escalonamentos
│   ├── discard*.go              # anchors discard: descarte seguro de trabalho
│   ├── unblock*.go              # anchors unblock: desbloqueio de tarefas
│   └── pr_body*.go              # anchors pr-body: geração de descrições de PR
│
├── quality/                     # Validação, Análise e Diagnóstico de Saúde
│   ├── check*.go                # anchors check: execução de pipeline de gates
│   ├── verify*.go               # anchors verify: verificação local e em tempo de commit
│   ├── doctor*.go               # anchors doctor: diagnóstico da saúde do ecossistema
│   ├── coverage*.go             # anchors coverage: cobertura de requisitos e código
│   ├── status*.go               # anchors status: resumo geral de saúde e tarefas
│   ├── report*.go               # anchors report: geração de relatórios consolidados
│   ├── stale*.go                # anchors stale: detecção de nós e arestas defasadas
│   └── suite*.go                # anchors test / mutation: execução de suítes de teste
│
├── governance/                  # Guias, Conformidade e Regras
│   ├── guide*.go                # anchors guide (plan, spec, code, feature, test...)
│   ├── spec_guide*.go           # anchors spec-guide: renderização de guias de spec
│   ├── governs*.go              # anchors governs: inspeção de regras de governo
│   ├── compliance*.go           # anchors compliance: conformidade de artefatos
│   └── audit*.go                # anchors audit: auditoria de conformidade
│
├── mapcmd/                      # Grafo e Rastreabilidade
│   ├── map*.go                  # anchors map (build, show): construção e visualização do grafo
│   ├── impact*.go               # anchors impact: análise de impacto bidirecional
│   ├── ingest*.go               # anchors ingest: ingestão de artefatos (lcov, junit, etc.)
│   ├── judge*.go                # anchors judge: julgamento de vereditos do grafo
│   └── recode*.go               # anchors recode: renomeação e migração de identidades
│
└── ops/                         # Configuração, Operação e Automação
    ├── init*.go                 # anchors init: inicialização de projetos com anchors.yaml
    ├── new*.go                  # anchors new: criação padronizada de specs, features e planos
    ├── migrate*.go              # anchors migrate: migração de versões de configuração
    ├── freeze*.go               # anchors freeze / thaw: congelamento de segurança
    ├── settings*.go             # anchors settings: perfil e configuração local do agente
    ├── install_hooks*.go        # anchors install-hooks: instalação de hooks git
    ├── commit_msg*.go           # anchors commit-msg: validação de mensagens de commit
    ├── board_serve*.go          # anchors board: servidor de visualização web
    ├── docs*.go                 # anchors docs: compilação de documentação a partir de doct
    ├── telemetry*.go            # anchors telemetry: aviso e configuração de métricas
    ├── code*.go                 # anchors code: geração e checagem de códigos de identidade
    ├── suggest*.go              # anchors suggest: sugestões contextuais para o agente
    ├── synthesize*.go           # anchors synthesize: síntese de discussões e decisões
    └── generated_paths*.go      # anchors generated-paths: caminhos gerados pelo framework
```

---

## 4. Pacotes Internos (`internal/`)

Os pacotes em [`internal/`](./internal/) compõem a biblioteca privada do Anchors, dividida em três estratos fundamentais:

### 4.1 Núcleo do Framework
- **`config`**: Definição da estrutura `Config`, carregamento e validação de `anchors.yaml`, resolução de raízes absolutas e convenções de comentários por linguagem.
- **`scan`**: Varredura de arquivos sem parsear código-fonte (somente texto puro e expressões regulares), extraindo metadados, identificadores de cenário e anotações `@noPropagation`.
- **`mapx`**: Estrutura do grafo de dependências (`anchors.graph.yaml`), nós, arestas (`specifies`, `covered-by`, `tested-by`, `references`, `convention`), montagem e serialização.

### 4.2 Regras de Qualidade e Gates
- **`gate`**: Mecanismo de execução e catálogo de gates (`secret-nao-vazado`, `layer-boundary`, `testid-coerente`, `trinca-completa`, `fase-ordenada`, etc.), gerando diagnósticos e perfis de veredito (`Pass`, `Fail`, `Pending`, `Skip`).

### 4.3 Infraestrutura e Suporte
- **`board`**: Servidor embutido para renderização do quadro Kanban/progresso no navegador.
- **`change`**: Rastreamento de alterações de arquivos e detecção de drift.
- **`checklog`**: Formatação e cabeçalhos espelhados para relatórios de verificação (`check-all.txt`).
- **`code`**: Algoritmos de compressão de nomes e garantia de unicidade de códigos de identidade.
- **`daemon`**: Gerenciamento de processos em segundo plano (PID, ciclo de vida e concorrência).
- **`doct`**: Compilação de documentação a partir de templates markdown.
- **`gitmeta`**: Consulta segura ao Git (HEAD, status de branches, autor e data) com fallbacks para ambientes sem git.
- **`health`**: Métricas e diagnósticos sistêmicos da saúde do repositório (órfãos, colisões, ambiente GitHub).
- **`i18n`**: Suporte internacional a mensagens (Português, Inglês, Espanhol).
- **`initx`**: Assistente de inicialização de projetos (perguntas interativas e geração de workflows).
- **`issue`**: Criação e manipulação de issues locais do Anchors em formato Markdown.
- **`migra`**: Rotinas de migração de esquema de configuração entre versões do framework.
- **`pack`**: Empacotamento e compilação de artefatos.
- **`queue`**: Gerenciamento atômico da fila de tarefas em `.anchors/tasks/`.
- **`recode`**: Propagação de mudanças de código de identidade através de arquivos e grafos.
- **`settings`**: Armazenamento e leitura das configurações locais do agente em `.anchors/settings.yaml`.
- **`similarity`**: Análise de similaridade textual e heurística.
- **`suggestion`**: Geração de sugestões proativas para agentes durante o ciclo de desenvolvimento.
- **`telemetry`**: Emissão assíncrona de telemetria anônima e verificação de opt-out.
- **`testsig`**: Assinatura e catálogo de testes executáveis.

---

## 5. Regras de Fronteira e Dependências

Para evitar ciclos de importação e manter a arquitetura sustentável:
1. **`internal/config`** é a base do núcleo e **nunca** importa `internal/scan`, `internal/mapx` ou `internal/gate`.
2. **`internal/scan`** pode importar `config`, mas não depende de `mapx` nem `gate`.
3. **`internal/mapx`** importa `config` e `scan`.
4. **`internal/gate`** consome `config`, `scan`, `mapx` e os pacotes de infraestrutura necessários para avaliar as regras.
5. Os comandos do CLI em **`cmd/anchors/`** consom `internal/*`, mas nenhum pacote interno pode importar `cmd/anchors`.
6. Subpacotes de `cmd/anchors` não mantêm dependências circulares entre si; funcionalidades transversais são centralizadas em `cmd/anchors/common/` ou delegadas aos pacotes de `internal/`.
