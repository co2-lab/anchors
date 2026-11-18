# Larder — auditoria de cobertura: cada passo de execução tem gate?

O Anchors prega "fechar todas as pontas" (`QUALITY.md` §5.1): nada de execução ou
propagação pode ficar sem validação. Esta auditoria aplica esse princípio **ao
próprio fluxo** — percorre a cascata passo a passo e, em cada transição, pergunta:
*que gate valida isto?* Um passo sem gate é uma **ponta de cobertura** (PC).

Para cada ponta, o gate que a fecha.

---

## A cascata auditada

```
[1] plano criado
     │  PC-1?
[2] plano → spec              validate-spec ✅
[3] spec → código             architecture, review-IA ✅
[4] spec → feature            validate-feature ✅
[5] feature → teste           traceability, test/coverage ✅
[6] teste → execução          testes passam ✅
[7] spec/feature → doc        PC-2?
 ⟂  cada agente → mapa        PC-3?  (lei de manutenção)
 ⟂  issue → plano             PC-5?
 ⟂  opt-out aplicado          — SEM gate por design (decisão do usuário)
```

Os passos [2]–[6] já têm gate (o miolo spec→código→feature→teste está coberto). As
pontas estão nas **bordas** (entrada e saída) e nos **passos transversais**.

---

## PC-1 — O plano nasce sem gate de entrada

**Passo:** `[1] plano criado`. Todo artefato downstream tem gate, mas o plano — a
origem — é validado por quê ao nascer? Se o plano está incoerente, incompleto ou
inexequível, tudo propaga errado a partir dele. O primeiro passo não tem gate.

**O gate que fecha:** **gate de validação de plano** (dono: Qualidade, régua:
Planejamento). Confronta o plano recém-criado:
- está completo (tem intenção, decomposição em specs, ordem)?
- é coerente (as specs que semeia fazem sentido para a intenção)?
- é exequível (o plano é realizável como está)?
Medidor: **review de IA**. Coerência, completude e exequibilidade de um plano não
são computáveis por script — é julgamento. (Se uma spec do plano cair numa camada
inexistente, o gate de estrutura pega isso depois, no seu próprio passo — não é
trabalho do gate de plano.) Regime: comportamental.

---

## PC-2 — A doc é gerada mas não confrontada (frescor)

**Passo:** `[7] spec/feature → doc`. A doc é âncora de consumo; dissemos que sua
única obrigação é **frescor** (`CONCEPT.md` §2). Mas há um gate de frescor rodando?
Sem ele, a doc drifta em silêncio — o elo terminal fica sem régua, e o stakeholder
lê uma mentira.

**O gate que fecha:** **gate de frescor da doc** (dono: Qualidade). Medidor: **review
de IA**. A sincronia (a aresta doc→fonte ficou stale?) é só o **gatilho** — diz
*quando* revisar. Mas validar que a doc *está boa* — reflete corretamente o que a
spec/feature descreve, está clara para o stakeholder, não mente — é **julgamento**,
não computável: a doc é prosa para humanos, e comparar prosa com a fonte é exatamente
o que script não faz e a IA faz. Então: a sincronia dispara (determinística), o
review de IA confronta (a doc de fato reflete a fonte agora?). Sem isso, a doc drifta
em silêncio e o stakeholder lê uma mentira.

---

## PC-3 — A atualização do mapa não é auto-validada

**Passo transversal:** cada agente, ao criar/mover/remover arquivo, **atualiza o
mapa** (a lei de manutenção, `TRACEABILITY.md` §4). Mas "lei" sem gate é só
intenção: e se um agente esquecer de registrar uma aresta? Vira **dependência
invisível** (o órfão de dependência). O gate "mapa fiel" (`TRACEABILITY.md` §7)
*existe* — a ponta é garantir que ele **roda a cada execução**, não só quando alguém
lembra.

**O gate que fecha:** **gate mapa-fiel incremental, disparado a cada execução de
agente** (dono: Qualidade, régua: Rastreabilidade). Medidor: **determinístico — um
script que percorre o trecho do mapa que o agente tocou e valida as anotações**. Não
precisa de IA nem de estado externo (só lê o mapa contra o disco), e é **incremental**
(só o que o ato mexeu, não o mapa inteiro):
- as arestas que este ato deveria ter criado/movido/removido estão no mapa?
- as arestas tocadas apontam para arquivos que existem? (senão → aresta morta)
- as anotações (rev, carimbo de validação, tipo) das arestas tocadas estão bem formadas?
Roda após cada agente para garantir que a lei de manutenção foi cumprida — o que
torna a lei **verificável**, não só declarada. Barato, local, reproduzível. (A
**integridade global** do mapa — o grafo inteiro consistente, sem ilhas — é papel do
validador de saúde, §5.2 de `QUALITY.md`, não deste gate incremental.)

---

## PC-4 — Opt-out aplicado: SEM gate, por design

**Passo transversal:** um `@no-test` ou `--no-block` é aplicado (opt-out honesto,
`CONCEPT.md` §5.1). Cogitou-se um gate que auditasse opt-outs velhos.

**Decisão: não há gate.** O opt-out é decisão **do usuário** — ele valida se ainda
faz sentido e reage quando não for como esperava. Um gate que o cobrasse seria o
framework arbitrando uma decisão que é do operador, contra o próprio espírito ("o
Anchors detecta e apresenta, não arbitra"). O opt-out já é honesto por ser explícito,
registrado e datado; o rastro basta para o usuário revisar *quando quiser*. Não é uma
ponta aberta — é um passo deliberadamente sem gate, porque a validação é humana.

---

## PC-5 — A conversão issue→plano: laço fechado pelo validador de saúde

**Passo transversal:** o usuário converte uma issue em plano (`CONCEPT.md` §5): a
issue encerra apontando o plano, o plano nasce com o propósito. O laço bidirecional
(issue↔plano) existe como referência — mas o plano pode ser concluído sem resolver a
dessincronia que o originou. O laço fica aberto.

**Não é um gate pontual — é sistêmico.** A conversão em si não precisa de validação
no momento; o laço não fecha por um gate num passo. É uma ponta **sistêmica** (o
estado do ecossistema: "há laços issue→plano que nunca fecharam?"), pega pelo
**validador de saúde do ecossistema** (`QUALITY.md` §5.2) — o meta-gate de visão
global (o `anchors doctor`), que varre laços abertos junto com as outras pontas
sistêmicas. Um plano de origem-issue só está genuinamente resolvido quando a aresta
que gerou a issue voltou a in-sync; o validador de saúde é quem confronta isso, na
varredura global, não no momento da conversão.

---

## Síntese: onde cada ponta fecha

**Gates de fluxo (locais — no momento do passo):**

| ponta | passo | gate | medidor | dono / régua |
|---|---|---|---|---|
| **PC-1** | plano criado | validação de plano | review de IA | Qualidade / Planejamento |
| **PC-2** | doc gerada | frescor da doc | review de IA (sincronia dispara) | Qualidade |
| **PC-3** | mapa atualizado | mapa-fiel incremental (script percorre o trecho tocado) | determin. | Qualidade / Rastreabilidade |

**Ponta sistêmica (visão global — validador de saúde, `QUALITY.md` §5.2):**

| ponta | passo | fecha por |
|---|---|---|
| **PC-5** | issue → plano (laço) | validador de saúde do ecossistema (`anchors doctor`) — varre laços que nunca fecharam + integridade global do mapa + pilares frouxos |

**Sem gate, por design:**

| ponta | passo | por quê |
|---|---|---|
| **PC-4** | opt-out aplicado | decisão do usuário; o framework não arbitra |

**Padrão:** pontas **locais** fecham por gates de fluxo no momento do passo (bordas =
review de IA, transversais = determinístico). Pontas **sistêmicas** (laços, integração
entre pilares, integridade global) só a visão global do **validador de saúde** pega —
é a encarnação executável da maturidade (`CONCEPT.md` §1.1), materializada num CLI
(`anchors doctor`/`status`). E há um passo deliberadamente sem gate (opt-out). O que a
auditoria revelou não foram furos de conceito, mas **passos que ainda não tinham gate
nomeado** — agora cada um tem seu destino: gate de fluxo, validador de saúde, ou
nenhum-por-design.

O meta-gate de completude (`QUALITY.md` §5.1) detecta as pontas *locais*; o validador
de saúde (§5.2) detecta as *sistêmicas*. Esta auditoria é os dois rodados à mão sobre
o Larder.
