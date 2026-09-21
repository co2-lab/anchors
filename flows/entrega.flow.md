<!-- @anchors
  code: ENTRG
  updated_at: 2026-09-21
-->
# Entrega — do trabalho pronto ao PR aberto

> **Código**: `ENTRG`

O que quem TRABALHA faz depois de a tarefa fechar. Vai até o PR e para ali — o que vem
depois é do pipeline, e tem fluxo próprio (`PIPEL`).

São dois fluxos e não um porque os donos são diferentes: aqui decide uma pessoa, lá decide
uma máquina. Juntá-los esconderia de quem espera onde a bola está — e "está comigo" e
"está com o CI" pedem coisas opostas de quem lê.

## Montagem

### ENTRG-P01 — commitar o trabalho

Encaixa: `ACPRC` (o pre-commit, que dispara sozinho)

Não é comando digitado: o gate roda no `git commit`, sobre o que está em stage.

Resultados:
- `ACPRC-R01` COMMIT PASSOU → `ENTRG-P03`
- `ACPRC-R04` IGNORADO → `ENTRG-P03`
- `ACPRC-R02` BARRADO POR GATE → `ENTRG-P02`
- `ACPRC-R03` BARRADO POR FALTA DE MAPA → `ENTRG-P04`

### ENTRG-P02 — corrigir o que o gate apontou

Encaixa: `ACHCK` (`anchors check --changed`)

Resultados:
- `ACHCK-R01` PROMOVÍVEL → `ENTRG-P01`
- `ACHCK-R02` BARRADO → `ENTRG-P02`

### ENTRG-P03 — abrir o PR

O trabalho está commitado. O que vem daqui é do pipeline.

Resultados:
- `ENTRG-P05` (o PR está aberto) → `ENTRG-P05`

### ENTRG-P04 — pôr o arquivo novo no mapa

Encaixa: `ACMAP` (`anchors map build`)

Resultados:
- `ACMAP-R01` MAPA EM DIA → `ENTRG-P01`
- `ACMAP-R02` CARIMBO PERDIDO → `ENTRG-P01`

### ENTRG-P05 — a bola está com o pipeline

Daqui em diante quem decide é o CI. O fluxo `PIPEL` descreve o que acontece.

> @terminal
