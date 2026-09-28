<!-- @anchors
  code: ADOCA
  updated_at: 2026-09-28
-->
# Adoção — do diretório ao primeiro plano

> **Código**: `ADOCA`

O fluxo que o `anchors init` abre. Antes dele não há decisão a tomar; todas nascem dele.

O ramo está documentado no guia e é fácil de pular: *"DESCOBRIR — só se o projeto NÃO
EXISTE ainda... pule esta etapa num projeto que já tem código"*. Como fluxo, a entrevista
de cinco etapas é o ÚNICO caminho quando o diretório está vazio — não dá para chegar à
Estrutura sem passar por ela.

## Montagem

### ADOCA-T01 — configurar o projeto

Encaixa: `ACINI` (`anchors init`)

Resultados:
- `ACINI-O01` ESTRUTURA INFERIDA → `ADOCA-T03`
- `ACINI-O02` DIRETÓRIO VAZIO → `ADOCA-T02`

### ADOCA-T02 — descobrir o que o projeto É

Encaixa: `ACGPJ` (`anchors guide project`)

Cinco etapas, uma pergunta por vez. Sem isto, a Estrutura de um diretório vazio seria
adivinhação com cara de configuração.

Resultados:
- `ACGPJ-O01` DESCOBERTA FEITA → `ADOCA-T01` (agora o init tem de onde inferir)

### ADOCA-T03 — construir o mapa

Encaixa: `ACMAP` (`anchors map build`)

O vigia, o `impact` e o `check` precisam que o mapa exista.

Resultados:
- `ACMAP-O01` MAPA EM DIA → `ADOCA-T04`

### ADOCA-T04 — pôr a esteira no ar

Encaixa: `ACWCH` (`anchors watch start`)

Resultados:
- `ACWCH-O01` VIGIA NO AR → `ADOCA-T05`
- `ACWCH-O02` JÁ ESTAVA RODANDO → `ADOCA-T05`

### ADOCA-T05 — pronto para planejar

O projeto tem Estrutura, mapa e esteira. O trabalho começa pelo fluxo `PLANO`.

> @terminal
