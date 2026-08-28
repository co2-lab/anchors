# Larder — simulação dos 4 fluxos + pontas abertas

Exercício: rodar o ciclo de vida de ponta a ponta pelos 6 pilares, sobre a app
[Larder](./STRUCTURE.md), observando **interações**, **inter-relações** e **pontas
abertas** (onde o framework não fecha). Ver o [diagrama](./diagram.svg).

Atores são **modelos de IA** — o Anchors foi pensado para ser operado por agentes.

---

## FLUXO 1 — Nova feature de UI: "adicionar item à despensa"

Caminho feliz, começa no Planejamento.

1. **Planejamento** (Planner-AI) — intenção "adicionar item pela tela". O plano
   *semeia* `AddPantryItemScreen.spec.md` (e opcionalmente a usecase spec, se
   quiser reforçar). Modo semente.
2. **Spec** (Spec-AI, Feature-AI) — a spec da tela nasce; Estrutura provê o template
   de screen. Cada estado/regra ganha código `APIT-S01`, `APIT-A01`. Regime
   comportamental → testável por Gherkin.
3. **Rastreabilidade** (Mapper-AI) — o mapa ganha nós e arestas: `plano→spec`,
   `spec→screen.tsx`, `spec→feature`. Identidade `APIT-*` amarra tudo.
4. **Propagação** (Propagator-AI) — spec rev 1. Onda **local**: alcança derivados.
   A usecase spec pode ficar stale (propagação sob demanda) se a tela demandar.
5. **Qualidade** — gates sobre o caminho de impacto: validate-spec, validate-feature,
   architecture, test/coverage (determinísticos) + Review-AI (coerência, spec
   completa). Dupla saída se reprova: issue + bloqueio.

✅ Fecha. Pontas: **PA-1** (aresta plano→spec), **PA-6** (regime misto).

---

## FLUXO 2 — Backend + infra: "listar receitas possíveis"

Estressa regime declarativo e `@no-test`.

1. **Planejamento** — semeia `listRecipes.handler.spec.md`. A cadeia
   handler→usecase→repository→entity→infra **não** está toda no plano.
2. **Spec** — handler spec (evento, não rota). Propaga sob demanda: usecase
   ("receita possível = todos os ingredientes na despensa"), repository (query),
   e a **infra** spec do índice DynamoDB — **declarativa**.
3. **Rastreabilidade** — a infra spec tem identidade e arestas como as outras.
4. **Propagação** — onda desce até a infra. A infra spec propaga para o `.tf`
   (conformidade), não para Gherkin.
5. **Qualidade** — handler/usecase: teste comportamental. Infra: teste de
   **conformidade** (o GSI existe/está conforme?) OU `@no-test` se o time não quer
   testar infra ainda.

✅ Fecha. Pontas: **PA-2** (criação de nós sob demanda), **PA-3** (conformidade como
medidor), **PA-4** (`@no-test`).

---

## FLUXO 3 — Mudança de régua: nova regra de arquitetura

Estressa onda global, maturação, buraco de cobertura — e a recomendação de plano.

1. Intenção: "repository nunca chama outro repository". Como é mudança de **alto
   grau** (o guide rege todas as specs de camada), o caminho **recomendado** é um
   **plano** (`PLANNING.md` §4) — não editar o guide direto. O plano vem acompanhado
   da revisão do guide + as revisões das N specs de repository afetadas, mapeando a
   árvore de antemão.
   > O direto (editar o guide) também funcionaria — onda global reativa —, mas é mais
   > caro. O plano dirige a propagação em vez de descobrir o impacto na marcha.
2. **Propagação** — as revisões semeadas pelo plano propagam normalmente (loop curto,
   âncora por âncora), cobrindo a árvore. Sem onda global às cegas.
3. **Qualidade** — o gate architecture roda sobre cada repository. Os que violam →
   dupla saída (issue + bloqueio). Se muitos violam, o gate nasce informativo
   (maturação) e é promovido quando o time zera. `--no-block` pontual se precisar
   mergear uma correção antes.

✅ Fecha. O plano é a porta eficiente para mudança de régua.

---

## FLUXO 4 — Código alterado fora do fluxo

Estressa propagação inversa, issue stale, decisão humana.

1. Dev edita `pantryRepository.ts` direto (otimiza query), sem tocar a spec. Rev do
   código (o `to` da aresta) avança.
2. **Propagação inversa** — `spec→repo.ts` fica stale pela ponta de baixo. Onda
   sobe até a spec. NÃO auto-resolve → issue `stale`.
3. **Decisão humana** (fluxo bidirecional): ou a spec ficou velha (atualiza), ou o
   código divergiu (corrige). O framework registra, não arbitra.

✅ Fecha.

---

## Pontas abertas encontradas

Agrupadas por tema. Cada uma é candidata a refino futuro — não bug de doutrina, mas
costura onde falta mecanismo.

### Tema A — como o grafo CRESCE ✅ RESOLVIDO

A premissa "o framework descreve percorrer, não criar" estava errada. Resolução:

- **O mapa é uma âncora auto-referente** — suas dependências são todos os documentos
  contidos nele; não precisa de lista sobre si mesmo (`TRACEABILITY.md` §4).
- **A lei de manutenção**: todo agente que cria/move/remove um arquivo **atualiza o
  mapa no mesmo ato**. A criação do mapa é distribuída, não de um dono
  (`TRACEABILITY.md` §4).
- **A onda cria e mapeia**: o agente de propagação disparado **cria** os artefatos
  faltantes e registra seus nós/arestas. A aresta nasce do ato de criação
  (`PROPAGATION.md` §6).
- **PA-1** (aresta plano→spec): o agente de planejamento, ao gerar a spec, registra
  `plano→spec`. ✅
- **PA-2** (propagação cria nós): o agente de propagação cria o artefato faltante e o
  mapeia. ✅
- **Materialização**: um **watcher inteligente** que conhece o mapa + a inter-relação
  dos pilares sabe qual arquivo mudou e qual agente chamar (`PROPAGATION.md` §6).

### Tema B — opt-outs unificados ✅ RESOLVIDO

- **PA-4** — nasce o conceito transversal **opt-out honesto** (`CONCEPT.md` §5.1):
  dispensa explícita, registrada e datada de uma exigência de validação. Três
  invariantes: explícito, registrado (dispensa a exigência, nunca o registro),
  datado+justificado (auditável). É a porta *legítima*; o buraco de cobertura é o
  oposto ponto a ponto (implícito/não-registrado/invisível).
- Instâncias amarradas à lei comum: **`@no-test`** (`SPEC.md` §6, dispensa o teste),
  **`--no-block`** (`QUALITY.md` §7, dispensa o bloqueio pontual), **maturação
  informativa** (`QUALITY.md` §7, dispensa o rigor). Mesmo conceito, coisas
  diferentes.
- O buraco de cobertura (`QUALITY.md` §5.1) agora tem critério: uma ponta sem gate
  só é aceitável se houver opt-out honesto cobrindo-a. ✅

### Tema C — o loop de volta ✅ RESOLVIDO

Duas correções ao modelo (o Fluxo 3 estava mal desenhado):

- **Mudança de régua é recomendação de plano, não regra.** Editar um guide direto
  também propaga e chega ao resultado certo — mas é *caro* (onda global às cegas).
  Um plano que traz as revisões de spec e mapeia a árvore é *mais eficiente*. O
  Anchors **recomenda** plano para mudanças de alto grau, não exige (`PLANNING.md`
  §4). O grau do nó sinaliza *quando o plano compensa*.
- **A issue é intervenção do usuário — 3 rotas** (`CONCEPT.md` §5): resolver ele
  mesmo, delegar a um agente, ou **converter em plano**. A conversão **nunca é
  automática** — é escolha do operador. Rastro: a issue encerra ("convertida no plano
  00XX"), o plano nasce com propósito ("aberto para resolver a issue 00XX"). Não viola
  "done não pode mentir": a issue foi *encaminhada*, não resolvida — débito
  transferido com rastro. ✅
- O Planejamento vira **origem E destino de reentrada** — o fluxo é cíclico.

### Tema D — refinos localizados ✅ RESOLVIDO

- **PA-3 — Conformidade NÃO é 3º medidor** (análise por 5 casos). Ela dissolve em
  determinístico/julgamento conforme o *critério esteja fixo* — Dockerfile↔spec e
  GSI↔query são determinísticos (forma↔forma); "3 réplicas = alta disponibilidade?"
  é julgamento se o critério é vago. O que era novo e valia nomear: o eixo
  **gate local vs. gate com estado externo** (`QUALITY.md` §3) — comparar contra o
  repo vs. contra a nuvem real (drift de infra). Ortogonal ao tipo de medidor. ✅
- **PA-6 — Regime é por requisito, via tag** (`SPEC.md` §6). Como nível e prioridade,
  cada requisito carrega sua tag de regime. "Spec mista" (entity: campos declarativos
  + invariantes comportamentais) = requisitos de regimes diferentes, não um 3º
  regime. Reusa o mecanismo de tag. ✅
