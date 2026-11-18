# Larder — estrutura de projeto (simulação Anchors)

App fictícia usada para simular os fluxos do ciclo de vida no Anchors. Exercita
todas as famílias de estrutura levantadas (frontend, backend, infra, CLI/lib) sobre
um único sistema coeso.

> **Larder** — "o que tenho na despensa e o que posso cozinhar". O usuário cadastra
> itens na despensa; o sistema calcula quais receitas são possíveis com o estoque.

---

## Camadas declaradas (o gabarito — pilar Estrutura de Projeto)

Ordem de dependência (a Propagação segue esta planta):

```
infra → entity → repository → usecase → handler → hook → store → screen/component
                                                        └── (mobile)
```

Cada camada tem seu **template de spec** (esqueleto comum + bloco de variação) e seu
**regime** (comportamental ou declarativo).

| camada | regime | template (bloco de variação) |
|---|---|---|
| infra | declarativo | recursos / secrets / permissões |
| entity | misto | invariantes + máquina de estados |
| repository | comportamental | queries / índices / acesso a dado |
| usecase | comportamental | regras de negócio |
| handler | comportamental | evento / request / response / status |
| hook | comportamental | efeitos (queries/mutations) |
| store | comportamental | estado global + ações |
| screen | comportamental | estado + navegação |
| component | comportamental | props + variantes + eventos |
| cli command | comportamental | args / flags / exit codes |

---

## Árvore de arquivos

Cada arquivo de interface tem seus **derivados co-localizados** (spec, feature,
teste) — a co-location gera as arestas `convention` do mapa.

```
larder/
├── anchors.graph.yaml            # O MAPA (pilar Rastreabilidade) — nós + arestas + carimbos
├── plan/
│   └── 0001-add-pantry-item.plan.md      # semeia specs (pilar Planejamento)
│   └── 0002-possible-recipes.plan.md
│
├── guides/                       # as RÉGUAS (nós de alto grau de saída → onda global)
│   ├── spec.guide.md             # como toda spec se comporta
│   ├── architecture.spec.md      # spec de arquitetura (rege as camadas)
│   ├── best-practices.guide.md   # régua do gate de review por IA
│   └── lint.guide.md
│
├── infra/                        # regime DECLARATIVO
│   ├── pantry-table.tf
│   ├── pantry-table.spec.md              # declarativa → teste de conformidade / @no-test
│   ├── photos-bucket.tf
│   └── photos-bucket.spec.md
│
├── api/                          # backend (Lambda)
│   ├── entities/
│   │   ├── PantryItem.ts
│   │   ├── PantryItem.spec.md            # regime misto
│   │   ├── Recipe.ts
│   │   └── Recipe.spec.md
│   ├── repositories/
│   │   ├── pantryRepository.ts
│   │   ├── pantryRepository.spec.md
│   │   ├── pantryRepository.feature
│   │   └── pantryRepository.test.ts
│   ├── usecases/
│   │   ├── addPantryItem.ts
│   │   ├── addPantryItem.spec.md
│   │   ├── addPantryItem.feature
│   │   ├── addPantryItem.test.ts
│   │   ├── listPossibleRecipes.ts
│   │   ├── listPossibleRecipes.spec.md
│   │   ├── listPossibleRecipes.feature
│   │   └── listPossibleRecipes.test.ts
│   └── handlers/
│       ├── addPantryItem.handler.ts
│       ├── addPantryItem.handler.spec.md         # evento, não rota
│       ├── addPantryItem.handler.feature
│       ├── listRecipes.handler.ts
│       ├── listRecipes.handler.spec.md
│       └── listRecipes.handler.feature
│
├── mobile/                       # frontend (React Native)
│   ├── hooks/
│   │   ├── usePantry.ts
│   │   ├── usePantry.spec.md
│   │   ├── usePantry.feature
│   │   └── usePantry.test.tsx
│   ├── stores/
│   │   ├── pantry.store.ts
│   │   └── pantry.store.spec.md
│   ├── components/
│   │   ├── PantryItemCard.tsx
│   │   ├── PantryItemCard.spec.md
│   │   ├── PantryItemCard.feature
│   │   └── PantryItemCard.test.tsx
│   └── screens/
│       ├── AddPantryItemScreen.tsx
│       ├── AddPantryItemScreen.spec.md    # códigos APIT-S01, APIT-A01...
│       ├── AddPantryItemScreen.feature
│       └── AddPantryItemScreen.test.tsx
│
├── cli/                          # ferramenta dev
│   ├── seed-pantry.command.ts
│   ├── seed-pantry.command.spec.md        # args / flags / exit codes
│   └── seed-pantry.command.test.ts
│
└── issues/                       # o REGISTRO (pontas-estado)
    ├── _templates/
    │   ├── stale.md
    │   ├── conflict.md
    │   └── violation.md
    ├── todo/
    ├── doing/
    └── done/
```

---

## Identidade de cenário (pilar Rastreabilidade)

O código atravessa as encarnações. Exemplo do requisito "adicionar item exibe
confirmação":

```
AddPantryItemScreen.spec.md   →  ### APIT-A01: Adicionar item
AddPantryItemScreen.feature   →  @acao @APIT-A01 @nivel-e2e @P0
AddPantryItemScreen.test.tsx  →  it('APIT-A01: adiciona item e mostra confirmação')
```

Gramática de fonte única: `PREFIXO(4) - TIPO(1) NUM(2)` — `APIT-A01`.
