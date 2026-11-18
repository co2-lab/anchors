---
title: Camadas do Projeto
description: "O que são camadas no Anchors, como funcionam as camadas regidas e reconhecidas, e quais gates se aplicam a cada uma."
---

Se você nunca usou o **Anchors**, imagine que está construindo uma casa. Você não começa assentando tijolos aleatoriamente em qualquer canto do terreno: primeiro você tem uma **planta**. A planta define onde fica a fundação, onde passam os canos de água, onde fica a fiação elétrica e onde ficam as paredes.

No Anchors, as **Camadas (Layers)** são exatamente a **planta da casa** do seu software. Elas dizem ao desenvolvedor e à inteligência artificial **onde cada peça de código deve morar**, **quem pode conversar com quem** e **qual nível de rigor é exigido para cada arquivo**.

Sem uma planta clara, uma IA que gera código sofre de amnésia estrutural: ela inventa pastas novas a cada prompt, importa banco de dados dentro de componentes visuais, duplica regras de negócio e quebra a arquitetura silenciosamente.

---

## 1. O que é uma Camada?

Uma **Camada** é um agrupamento arquitetural de arquivos com o mesmo papel no sistema. Em vez de tratar todos os arquivos do repositório como iguais, o Anchors classifica cada arquivo por meio de padrões de diretório e extensão (`pattern`) definidos no arquivo [`anchors.yaml`](/pt/docs/anchors-yaml///).

Cada camada possui propriedades bem definidas:
- **`pattern`**: A máscara de arquivos (glob) que pertencem a ela (ex: `src/services/**/*.ts`).
- **`kind`**: O tipo de artefato (`code`, `spec`, `feature`, `test`, `doc`, `guide`, `plan`, `product`, `flag`).
- **`regime`**: Se a camada é comportamental (`regra`) ou apenas estrutural (`declarativo`).
- **`tags`**: Palavras-chave usadas pelos [gates](/pt/docs/concepts/gates-e-vereditos/) para saber como avaliar a camada.
- **`depends_on`** / **`governs`**: Regras de importação e fronteira arquitetural.

---

## 2. Os Dois Grandes Tipos de Camadas

Esta é a distinção mais importante que você precisa entender no Anchors: **nem todo arquivo precisa de uma especificação completa de regras de negócio**. O Anchors separa as camadas em dois grandes grupos:

```
┌────────────────────────────────────────────────────────┐
│                   CAMADAS NO ANCHORS                   │
└──────────────────────────────────┬─────────────────────┘
                                   │
         ┌─────────────────────────┴─────────────────────────┐
         ▼                                                   ▼
┌─────────────────────────────────┐ ┌─────────────────────────────────┐
│        CAMADAS REGIDAS          │ │       CAMADAS RECONHECIDAS      │
│     (Regime Comportamental)     │ │       (Regime Declarativo)      │
├─────────────────────────────────┤ ├─────────────────────────────────┤
│ • Contêm regras de negócio      │ │ • Não originam regras próprias  │
│ • Exigem a Unidade completa      │ │ • Apenas conectam ou declaram   │
│   (Spec + Feature + Teste)      │ │ • Identificação por layer:      │
│ • Ex: usecase, service, domain  │ │ • Ex: infra, dao, types, doc    │
└─────────────────────────────────┘ └─────────────────────────────────┘
```

### A. Camadas Regidas (Governed Layers)
São as camadas onde vive a **inteligência e o comportamento** do seu sistema. Uma camada regida contém regras de validação, fluxos de decisão, regras financeiras ou processamento de dados.

- **Regime**: `regime: regra` (ou comportamental).
- **Exigência**: Toda unidade de código dessa camada precisa obrigatoriamente de uma [Unidade](/pt/docs/concepts/unidade/) completa:
  1. Uma **Spec** (`.spec.md`) catalogando suas regras de negócio com [códigos de identidade](/pt/docs/concepts/rastreabilidade-e-codigos/).
  2. Uma **Feature** (`.feature`) com cenários em Gherkin cobrindo as regras.
  3. Um **Teste** (`_test.go`, `.test.ts`) comprovando a execução dos cenários.
  4. O **Código** implementando as funções correspondentes.
- **Exemplos comuns**:
  - `usecase`: Casos de uso da aplicação.
  - `service`: Serviços de negócio.
  - `domain`: Entidades e modelos ricos de domínio.
  - `comando`: Comandos executáveis da CLI ou handlers de API.

### B. Camadas Reconhecidas (Recognized / Declarative Layers)
São as camadas técnicas que **não tomam decisões de negócio**. Elas existem para conectar coisas: adaptadores de banco de dados, clientes HTTP, definições de tipos primitivos, scripts de infraestrutura ou arquivos de documentação.

- **Regime**: `regime: declarativo`.
- **Exigência**: Elas **NÃO** exigem especificação funcional nem unidade de cenários. Exigir spec e cenários de um DAO que apenas roda `SELECT * FROM users` forçaria o desenvolvedor a inventar "regras falsas" para enganar a ferramenta.
- **Identidade mínima**: No cabeçalho `@anchors`, basta declarar a que camada o arquivo pertence (`layer: infra` ou `layer: dao`).
- **Exemplos comuns**:
  - `infra`: Código de CDK, Terraform ou conexões com cloud.
  - `dao` / `repository`: Consultas puras a banco de dados.
  - `domain-types`: Definições puras de structs/interfaces e DTOs.
  - `doc`: Documentações de uso, manuais e guias de leitura.

> [!TIP]
> **A regra de ouro:** Se o arquivo decide algo sobre o negócio do projeto, a camada é **Regida**. Se o arquivo apenas executa uma instrução técnica ou declara uma estrutura sem decidir, a camada é **Reconhecida**.

---

## 3. Camadas Especiais de Governança

Além do código tradicional da aplicação, o Anchors governa os próprios artefatos de orientação através de camadas dedicadas:

| Camada Especial | Padrão Típico | Kind | Papel no Projeto |
| --- | --- | --- | --- |
| **Doutrina** | `{CONCEPT,QUALITY,...}.md` | `guide` | Manuais de governança e princípios inegociáveis do projeto. |
| **Produto** | `product/*.doctrine.md` | `product` | [Doutrina de Produto](/pt/docs/concepts/doutrina-de-produto/): regras de negócio transversais que cortam múltiplas telas e unidades. |
| **Flags** | `flags/*.flag.md` | `flag` | [Feature Flags](/pt/docs/concepts/feature-flags/): cenários de ativação/desativação condicional sem duplicar specs. |
| **Desenho** | `DESIGN-*.md` | `doc` | Registros arquiteturais prévios (RFCs, decisões técnicas, trade-offs). |
| **Planos** | `PLANNING.md`, `plans/*.plan.md` | `plan` | Sequenciamento de fases e semeadura de specs futuras. |

---

## 4. Regimes de Teste e Superfícies por Camada

Uma camada regida é testada de acordo com o **regime de teste** adequado à sua natureza:

| Regime Canônico | O que valida | Onde mora (Superfície Típica) | Camadas Recomendadas |
| --- | --- | --- | --- |
| `unit` | Regras puras em memória, cálculos, isolamento | Arquivo de teste unitário ao lado do código | `usecase`, `service`, `domain`, `regra` |
| `integration` | Comunicação com banco real, filas, adaptadores | Pasta de integração (`tests/integration/`) | `dao`, `infra`, `api-client` |
| `e2e` | Fluxo completo do usuário de ponta a ponta | Roteiros executáveis (`.yaml`, Cypress, Playwright) | `interface`, `screen`, `comando` |
| `vr` | Aparência visual e ausência de regressão de tela | Imagens de baseline (`.png`) | `screen`, `componente`, `presentation` |

O gate [`test-level-codes`](/pt/docs/gates///test-level-codes/) garante que cenários de um regime não tentem testar coisas de outro (ex: proibir testes unitários de realizarem asserções de regressão visual).

---

## 5. Regras de Fronteira: Quem Pode Acessar Quem?

Um dos maiores benefícios de declarar camadas é impedir o **acoplamento caótico**. O Anchors fornece dois gates automáticos para proteger suas fronteiras:

1. **[`layer-boundary`](/pt/docs/gates///layer-boundary/)**: Garante que uma camada não importe quem ela está proibida de importar. Por exemplo:
   - A camada de `domain` nunca pode importar `infra` ou `web`.
   - A camada de `usecase` pode importar `domain`, mas não pode importar `controller`.
2. **[`sibling-guard`](/pt/docs/gates///sibling-guard/)**: Garante que módulos vizinhos dentro da mesma camada não violem convenções de parâmetros ou façam chamadas proibidas por caminhos indiretos.

---

## 6. Como Declarar Camadas no `anchors.yaml`

Abaixo está um exemplo real de declaração de camadas:

```yaml
layers:
  # Camada Regida: Casos de Uso (exige spec, feature e teste unitário)
  usecase:
    pattern: "src/usecases/**/*.ts"
    kind: code
    regime: regra
    tags: [usecase, negocio]
    depends_on: [domain]

  # Camada Reconhecida: Infraestrutura (dispensa spec; identidade por layer:)
  infra:
    pattern: "src/infra/**/*.ts"
    kind: code
    regime: declarativo
    tags: [infra, adaptador]

  # Camada Especial: Doutrina de Produto
  produto:
    pattern: "product/*.doctrine.md"
    kind: product
    tags: [doutrina, negocio]

  # Camada de Documentação
  doc:
    pattern: "docs/**/*.md"
    kind: doc
    tags: [documentacao]
```

---

## 7. Mapeamento de Gates por Camada

A tabela a seguir orienta quais [gates](/pt/docs/gates///) devem ser ligados dependendo da camada e do tipo de arquivo:

| Família de Gates | Camadas Regidas (`regra`) | Camadas Reconhecidas (`declarativo`) | Telas e Apresentação | Documentos e Guias |
| --- | :---: | :---: | :---: | :---: |
| **Unidade e Estrutura**<br>([`unit-complete`](/pt/docs/gates///unit-complete/), [`has-code`](/pt/docs/gates///has-code/), [`spec-sections`](/pt/docs/gates///spec-sections/)) | **Obrigatório** (`blocking: true`) | *Dispensado* | **Obrigatório** se houver regras visuais | *Dispensado* |
| **Identidade e Rastreabilidade**<br>([`code-cataloged`](/pt/docs/gates///code-cataloged/), [`feature-test-match`](/pt/docs/gates///feature-test-match/), [`scenario-asserts`](/pt/docs/gates///scenario-asserts/)) | **Obrigatório** (`blocking: true`) | *Dispensado* | **Obrigatório** | *Dispensado* |
| **Fronteiras Arquiteturais**<br>([`layer-boundary`](/pt/docs/gates///layer-boundary/), [`sibling-guard`](/pt/docs/gates///sibling-guard/)) | **Obrigatório** (`blocking: true`) | **Obrigatório** (`blocking: true`) | **Obrigatório** (`blocking: true`) | *Dispensado* |
| **Qualidade da Prova**<br>([`tests-pass`](/pt/docs/gates///tests-pass/), [`line-coverage`](/pt/docs/gates///line-coverage/), [`mutation-score`](/pt/docs/gates///mutation-score/)) | **Recomendado** | Conforme suíte de integração | Conforme suíte E2E | *Dispensado* |
| **Apresentação e Telas**<br>([`presentation-exhaustive`](/pt/docs/gates///presentation-exhaustive/), [`vr-baseline`](/pt/docs/gates///vr-baseline/)) | *Não se aplica* | *Não se aplica* | **Obrigatório** | *Não se aplica* |
| **Segurança e Higiene**<br>([`no-secret-leaked`](/pt/docs/gates///no-secret-leaked/), [`dependency-vulnerable`](/pt/docs/gates///dependency-vulnerable/), [`spellcheck`](/pt/docs/gates///spellcheck/)) | **Obrigatório** | **Obrigatório** | **Obrigatório** | **Obrigatório** |
| **Conformidade de Cabeçalho**<br>([`header-valid`](/pt/docs/gates///header-valid/), [`updated-at-atual`](/pt/docs/gates///updated-at-atual/)) | Exige `code:` ou `ref:` | Exige `layer:` | Exige `code:` ou `ref:` | Exige `layer:` ou `code:` |

Para consultar todos os gates disponíveis em detalhes, acesse o [Catálogo Completo de Gates](/pt/docs/gates///).
