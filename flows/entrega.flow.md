<!-- @anchors
  code: ENTRG
  updated_at: 2026-09-28
-->
# Entrega — do trabalho pronto ao PR aberto

> **Código**: `ENTRG`

O que quem TRABALHA faz depois de a tarefa fechar. Vai até o PR e para ali — o que vem
depois é do pipeline, e tem fluxo próprio (`PIPEL`).

São dois fluxos e não um porque os donos são diferentes: aqui decide uma pessoa, lá decide
uma máquina. Juntá-los esconderia de quem espera onde a bola está — e "está comigo" e
"está com o CI" pedem coisas opostas de quem lê.

## Montagem

### ENTRG-T01 — commitar o trabalho

Encaixa: `ACPRC` (o pre-commit, que dispara sozinho)

Não é comando digitado: o gate roda no `git commit`, sobre o que está em stage.

Resultados:
- `ACPRC-O01` COMMIT PASSOU → `ENTRG-T03`
- `ACPRC-O04` IGNORADO → `ENTRG-T03`
- `ACPRC-O02` BARRADO POR GATE → `ENTRG-T02`
- `ACPRC-O03` BARRADO POR FALTA DE MAPA → `ENTRG-T04`

### ENTRG-T02 — corrigir o que o gate apontou

Encaixa: `ACHCK` (`anchors check --changed`)

Resultados:
- `ACHCK-O01` PROMOVÍVEL → `ENTRG-T01`
- `ACHCK-O02` BARRADO → `ENTRG-T02`

### ENTRG-T03 — abrir o PR

O trabalho está commitado. O que vem daqui é do pipeline.

Resultados:
- `ENTRG-T05` (o PR está aberto) → `ENTRG-T05`

### ENTRG-T04 — pôr o arquivo novo no mapa

Encaixa: `ACMAP` (`anchors map build`)

Resultados:
- `ACMAP-O01` MAPA EM DIA → `ENTRG-T01`
- `ACMAP-O02` CARIMBO PERDIDO → `ENTRG-T01`

### ENTRG-T05 — a bola está com o pipeline

Daqui em diante quem decide é o CI. O fluxo `PIPEL` descreve o que acontece.

> @terminal
