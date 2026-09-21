<!-- @anchors
  code: ADOCA
  updated_at: 2026-09-21
-->
# Adoção — do diretório ao primeiro plano

> **Código**: `ADOCA`

O fluxo que o `anchors init` abre. Antes dele não há decisão a tomar; todas nascem dele.

O ramo está documentado no guia e é fácil de pular: *"DESCOBRIR — só se o projeto NÃO
EXISTE ainda... pule esta etapa num projeto que já tem código"*. Como fluxo, a entrevista
de cinco etapas é o ÚNICO caminho quando o diretório está vazio — não dá para chegar à
Estrutura sem passar por ela.

## Montagem

### ADOCA-P01 — configurar o projeto

Encaixa: `ACINI` (`anchors init`)

Resultados:
- `ACINI-R01` ESTRUTURA INFERIDA → `ADOCA-P03`
- `ACINI-R02` DIRETÓRIO VAZIO → `ADOCA-P02`

### ADOCA-P02 — descobrir o que o projeto É

Encaixa: `ACGPJ` (`anchors guide project`)

Cinco etapas, uma pergunta por vez. Sem isto, a Estrutura de um diretório vazio seria
adivinhação com cara de configuração.

Resultados:
- `ACGPJ-R01` DESCOBERTA FEITA → `ADOCA-P01` (agora o init tem de onde inferir)

### ADOCA-P03 — construir o mapa

Encaixa: `ACMAP` (`anchors map build`)

O vigia, o `impact` e o `check` precisam que o mapa exista.

Resultados:
- `ACMAP-R01` MAPA EM DIA → `ADOCA-P04`

### ADOCA-P04 — pôr a esteira no ar

Encaixa: `ACWCH` (`anchors watch start`)

Resultados:
- `ACWCH-R01` VIGIA NO AR → `ADOCA-P05`
- `ACWCH-R02` JÁ ESTAVA RODANDO → `ADOCA-P05`

### ADOCA-P05 — pronto para planejar

O projeto tem Estrutura, mapa e esteira. O trabalho começa pelo fluxo `PLANO`.

> @terminal
