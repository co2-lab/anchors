# Guia de Gates Externos por Ecossistema

> Este guia documenta como configurar os **gates externos de qualidade, segurança e governança** no Anchors para diferentes ecossistemas e plataformas.
> 
> **Filosofia do Anchors:** O framework é estritamente agnóstico. Ele governa a **semântica da qualidade** (o que medir, quando medir, escopo de lote vs projeto e se bloqueia o avanço). A **ferramenta concreta** de execução (`run:`) pertence à configuração do projeto (`anchors.yaml`).

---

## 1. Visão Geral dos Gates Externos

| Gate | Categoria | Escopo | Momento (`when`) | Custo | Propósito |
|---|---|---|---|---|---|
| [`no-secret-leaked`](#1-no-secret-leaked) | `security` | `batch` | `pre-commit`, `ci` | `fast` | Impede que credenciais e chaves entrem no repositório |
| [`dependency-vulnerable`](#2-dependency-vulnerable) | `security` | `project` | `pre-push`, `ci` | `fast`/`slow` | Audita vulnerabilidades (CVEs) conhecidas nas dependências |
| [`sbom-generated`](#3-sbom-generated) | `provenance` | `project` | `ci` | `slow` | Gera inventário SBOM (Software Bill of Materials) |
| [`no-duplication`](#4-no-duplication) | `quality` | `project` | `ci` | `slow` | Detecta código duplicado (copy-paste) entre arquivos |
| [`spellcheck`](#5-spellcheck) | `style` | `batch` | `pre-commit`, `ci` | `fast` | Verifica grafia em código, testes, specs e documentação |
| [`license-compatible`](#6-license-compatible) | `legal` | `project` | `ci` | `fast` | Assegura ausência de dependências com copyleft forte ou incompatível |
| [`circular`](#7-circular) | `architecture` | `project` | `pre-push`, `ci` | `slow` | Detecta dependências circulares entre módulos ou pacotes |
| [`deadcode`](#8-deadcode) | `maintenance` | `project` | `ci` | `slow` | Identifica exports, funções e arquivos órfãos sem consumidor |
| [`no-test-proof-real`](#9-no-test-proof-real) | `quality` | `node` | `pre-push`, `ci` | `slow` | Gate de IA: verifica se o teste de fato exercita o comportamento ou só cita o código |

---

## 2. Catálogo e Opções por Plataforma

---

### 1. `no-secret-leaked`
* **Objetivo**: Garantir que chaves de API, senhas, tokens de acesso ou certificados nunca sejam commitados.
* **Por que nasce bloqueante**: Segredo vazado no git não tem volta simples — rotacionar credencial é caro e o histórico preserva o valor.

#### Opções de Ferramentas
* **Universal / Recomendada (Binário Autônomo em Go)**:
  - **`gitleaks`**: Opera diretamente no git (staged ou commit range). Não exige nenhum runtime instalado.
  - Instalação: `brew install gitleaks` ou download via GitHub Releases.
  - Exemplo de configuração:
    ```yaml
    - name: no-secret-leaked
      on: [code, test, doc]
      scope: batch
      scope_full: project
      run: "gitleaks git --no-banner --redact -v"
      needs_tool: gitleaks
      install_hint: "brew install gitleaks"
      blocking: true
      when: [pre-commit, ci]
      cost: fast
    ```
* **Alternativas por Ecossistema**:
  - **Universal**: `trufflehog` (`brew install trufflehog`)
  - **Python**: `detect-secrets` (`pip install detect-secrets`)
  - **Shell / AWS**: `git-secrets` (`brew install git-secrets`)

---

### 2. `dependency-vulnerable`
* **Objetivo**: Auditar as dependências de terceiros contra bases de dados de vulnerabilidades (OSV, NVD, GitHub Advisory).
* **Comportamento recomendado**: Informativo inicialmente (`blocking: false`), promovido a bloqueante quando o passivo de CVEs transitivas estiver controlado.

#### Opções de Ferramentas
* **Universal / Recomendada (Binário Autônomo em Go)**:
  - **`osv-scanner`** (Google): Suporta nativamente lockfiles de Go (`go.mod`), Node (`yarn.lock`, `pnpm-lock.yaml`, `package-lock.json`), Python (`requirements.txt`, `poetry.lock`), Rust (`Cargo.lock`), Java (`pom.xml`), PHP (`composer.lock`), Ruby (`Gemfile.lock`), etc.
  - Instalação: `brew install osv-scanner` ou `go install github.com/google/osv-scanner/cmd/osv-scanner@latest`.
  - Exemplo:
    ```yaml
    - name: dependency-vulnerable
      on: [code]
      scope: project
      run: "osv-scanner scan source -r ."
      needs_tool: osv-scanner
      install_hint: "brew install osv-scanner"
      blocking: false
      when: [ci]
      cost: fast
    ```
* **Universal Alternativa**:
  - **`trivy`** (Aqua Security): `trivy fs --security-checks vuln .`
* **Nativas por Ecossistema**:
  - **Go**: `govulncheck ./...` (`go install golang.org/x/vuln/cmd/govulncheck@latest`)
  - **Node.js / TS**: `pnpm audit --prod`, `npm audit --omit=dev`, `yarn audit`
  - **Python**: `pip-audit` (`pip install pip-audit`)
  - **Rust**: `cargo-audit` (`cargo install cargo-audit`)
  - **.NET / C#**: `dotnet list package --vulnerable`

---

### 3. `sbom-generated`
* **Objetivo**: Gerar o inventário padronizado de componentes de software (CycloneDX ou SPDX) para conformidade, auditoria e cadeia de suprimentos.

#### Opções de Ferramentas
* **Universal / Recomendada (Binário Autônomo em Go)**:
  - **`syft`** (Anchore): Inspeciona o diretório e reconhece automaticamente pacotes de todas as principais linguagens.
  - Instalação: `brew install syft`.
  - Exemplo:
    ```yaml
    - name: sbom-generated
      on: [code]
      scope: project
      run: "syft scan dir:. -o cyclonedx-json=sbom.json -q"
      needs_tool: syft
      install_hint: "brew install syft"
      blocking: false
      when: [ci]
      cost: slow
    ```
* **Nativas por Ecossistema**:
  - **Go**: `cyclonedx-gomod app -json -output sbom.json`
  - **Node.js / TS**: `npx @cyclonedx/cyclonedx-npm --output-file sbom.json`
  - **Python**: `cyclonedx-py requirements -o sbom.json`
  - **Rust**: `cargo-cyclonedx` (`cargo install cargo-cyclonedx`)
  - **Java / Kotlin**: Plugins Maven (`cyclonedx-maven-plugin`) ou Gradle (`cyclonedx-gradle-plugin`)

---

### 4. `no-duplication`
* **Objetivo**: Detectar trechos de código idênticos ou com alta similaridade (copy-paste) que deveriam ser abstraídos ou consolidados.
* **Comportamento recomendado**: Informativo (`blocking: false`), com limiar configurado na ferramenta do projeto.

#### Opções de Ferramentas
* **Multi-linguagem (Binário Nativo em C++)**:
  - **`pmd cpd`**: Analisa Java, Go, Python, C#, PHP, JavaScript, TypeScript, C/C++, Ruby, Kotlin, Swift.
  - Instalação: `brew install pmd`.
  - Exemplo:
    ```yaml
    - name: no-duplication
      on: [code]
      scope: project
      run: "pmd cpd --minimum-tokens 70 --dir . --language typescript"
      needs_tool: pmd
      install_hint: "brew install pmd"
    ```
* **Node.js / Multi-linguagem via npm**:
  - **`jscpd`**: Suporta mais de 150 linguagens.
  - Exemplo:
    ```yaml
    - name: no-duplication
      on: [code]
      scope: project
      run: "npx --yes jscpd . --reporters console --silent"
      needs_tool: npx
      install_hint: "instale Node.js"
    ```
* **Nativas por Ecossistema**:
  - **Go**: `dupl` (`go install github.com/mibk/dupl@latest`) → `dupl -t 70`
  - **Python**: `pylint --disable=all --enable=similarities .`
  - **Rust**: `duplication-detector` ou `flcl`
  - **.NET / C#**: `dotnet-sonarscanner` ou `Simian`

---

### 5. `spellcheck`
* **Objetivo**: Eliminar erros de digitação em identificadores de código (camelCase, snake_case), comentários, mensagens de erro, specs e documentação.

#### Opções de Ferramentas
* **Universal / Recomendada (Binário Autônomo em Rust)**:
  - **`typos`**: Ultra-rápido, escrito em Rust, sem dependência de runtime. Reconhece camelCase, kebab-case e snake_case nativamente.
  - Instalação: `brew install typos` ou `cargo install typos-cli`.
  - Exemplo:
    ```yaml
    - name: spellcheck
      on: [code, test, spec]
      scope: batch
      scope_full: project
      run: "typos"
      needs_tool: typos
      install_hint: "brew install typos"
      blocking: true
      when: [pre-commit, ci]
      cost: fast
    ```
* **Node.js / TS (Ecossistema Web)**:
  - **`cspell`**: Muito customizável via `cspell.config.yaml` e dicionários de termos técnicos.
  - Instalação: `pnpm add -D cspell` ou `npm install -D cspell`.
  - Exemplo:
    ```yaml
    - name: spellcheck
      on: [code, test, spec]
      scope: batch
      scope_full: project
      run: "npx cspell --no-progress --no-summary {{files}}"
      needs_tool: npx
      install_hint: "instale Node.js e cspell"
    ```
* **Python**:
  - **`codespell`**: `pip install codespell` → `codespell -q 3`

---

### 6. `license-compatible`
* **Objetivo**: Evitar contaminação do projeto por dependências de terceiros com licenças incompatíveis (como copyleft forte AGPL/SSPL/GPL em projetos proprietários) ou cláusulas restritivas (Commons Clause, Non-Commercial).

#### Opções de Ferramentas
* **Go**:
  - **`go-licenses`**: Analisa o grafo compilado de dependências do módulo Go.
  - Instalação: `go install github.com/google/go-licenses@latest`.
  - Exemplo:
    ```yaml
    - name: license-compatible
      on: [code]
      scope: project
      run: "go-licenses check ./... --disallowed_types=forbidden,restricted"
      needs_tool: go-licenses
      install_hint: "go install github.com/google/go-licenses@latest"
    ```
* **Node.js / TS**:
  - **`license-checker`**: Varre `node_modules` filtrando apenas dependências de produção.
  - Exemplo:
    ```yaml
    - name: license-compatible
      on: [code]
      scope: project
      run: "npx license-checker --production --onlyAllow 'MIT;Apache-2.0;BSD-2-Clause;BSD-3-Clause;ISC'"
      needs_tool: npx
    ```
* **Python**:
  - **`pip-licenses`**: `pip-licenses --from=mixed --format=json`
* **Rust**:
  - **`cargo-deny`**: `cargo-deny check bans licenses` (utilitário completo de governança para Rust)

---

### 7. `circular`
* **Objetivo**: Evitar dependências circulares entre pacotes ou módulos internos, que dificultam manutenção, quebram inicialização de Singletons/módulos e impedem desmembramento de arquitetura.

#### Opções de Ferramentas
* **Go**:
  - Nativo! O compilador Go (`go build`) e ferramentas como `go vet` barram ciclos de importação por especificação da linguagem.
  - Para pacotes internos em monorepos: `golangci-lint run --enable depguard`.
* **Node.js / TypeScript**:
  - **`madge`**: Gera grafo de dependências via AST e detecta ciclos.
  - Exemplo:
    ```yaml
    - name: circular
      on: [code]
      scope: project
      run: "npx madge --circular --extensions ts,tsx src/"
      needs_tool: npx
      when: [pre-push, ci]
      cost: slow
    ```
* **Python**:
  - **`import-linter`**: `pip install import-linter` → `lint-imports`
* **Java / Kotlin**:
  - **`ArchUnit`** ou **`JDepend`**: Testes arquiteturais integrados ao JUnit.
* **Rust**:
  - O sistema de módulos do Rust (`mod` e `use`) com visibilidade estrita controla ciclos entre crates.

---

### 8. `deadcode`
* **Objetivo**: Identificar código morto, funções nunca chamadas, tipos ou exportações órfãs e dependências declaradas mas não importadas.

#### Opções de Ferramentas
* **Go**:
  - **`deadcode`** (Oficial do Go): Análise de acessibilidade a partir do `main` ou pacotes exportados.
  - Instalação: `go install golang.org/x/tools/cmd/deadcode@latest`.
  - Exemplo:
    ```yaml
    - name: deadcode
      on: [code]
      scope: project
      run: "deadcode ./..."
      needs_tool: deadcode
      install_hint: "go install golang.org/x/tools/cmd/deadcode@latest"
    ```
* **Node.js / TypeScript**:
  - **`knip`**: Especializado em monorepos TS/JS, acha exports não usados, arquivos órfãos e pacotes desnecessários.
  - Exemplo:
    ```yaml
    - name: deadcode
      on: [code]
      scope: project
      run: "npx knip"
      needs_tool: npx
      blocking: false
    ```
* **Python**:
  - **`vulture`**: `pip install vulture` → `vulture src/`
* **Rust**:
  - **`cargo-udeps`**: `cargo install cargo-udeps` → `cargo +nightly udeps`

---

### 9. `no-test-proof-real` (Gate Canônico de IA)
* **Objetivo**: Onde uma spec dispensou teste unitário com `@no-test` apontando um teste externo como prova, este gate submete à IA a pergunta: *"O teste apontado realmente EXERCITA o comportamento descrito ou apenas CITA o código/símbolo?"*.
* **Status**: 100% agnóstico e sem ferramenta de runtime externa — roda via `anchors judge` diretamente pelo motor do Anchors.

---

## 3. Matriz Resumida de Ferramentas por Ecossistema

| Gate | Universal (Binário Autônomo) | Go | Node / TypeScript | Python | Rust |
|---|---|---|---|---|---|
| **`no-secret-leaked`** | `gitleaks`, `trufflehog` | `gitleaks` | `gitleaks` | `detect-secrets` | `gitleaks` |
| **`dependency-vulnerable`** | `osv-scanner`, `trivy` | `govulncheck` | `pnpm audit` / `osv-scanner` | `pip-audit` | `cargo-audit` |
| **`sbom-generated`** | `syft`, `trivy sbom` | `syft` / `cyclonedx-gomod` | `syft` / `@cyclonedx/cyclonedx-npm`| `cyclonedx-py` | `cargo-cyclonedx` |
| **`no-duplication`** | `pmd cpd` | `dupl` | `jscpd` | `pylint (similarities)`| `flcl` / `pmd cpd` |
| **`spellcheck`** | `typos` | `typos` | `typos` ou `cspell` | `typos` ou `codespell`| `typos` |
| **`license-compatible`**| — | `go-licenses` | `license-checker` | `pip-licenses` | `cargo-deny` |
| **`circular`** | — | Compilador Go / `go vet` | `madge` | `import-linter` | Compilador Rust |
| **`deadcode`** | — | `deadcode` (x/tools) | `knip` | `vulture` | `cargo-udeps` |
