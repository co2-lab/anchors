<!-- @anchors
  code: ACHCK
  updated_at: 2026-09-21
-->
# Ação: `anchors check` — confrontar o que foi escrito

> **Código**: `ACHCK`

Roda os gates sobre o que mudou (`--changed`) ou sobre o projeto inteiro (`--all`).

Uma AÇÃO é a peça do quebra-cabeça: ela declara o que faz e quais RESULTADOS oferece.
Quem encaixa os resultados em próximos passos é o fluxo — a ação não sabe quem vem
depois, e é isso que a torna reusável em mais de um.

## Comando

    anchors check --changed <arquivo>     incremental, o que a mudança tocou
    anchors check --all                   o projeto inteiro

## Resultados

### ACHCK-R01 — PROMOVÍVEL: nenhum gate bloqueante reprovou

Pode haver achado informativo aberto — informativo não barra promoção. E não some
sozinho: é onde vive o que os bloqueantes não confrontam.

### ACHCK-R02 — BARRADO: um gate bloqueante reprovou

Cada gate que reprova gera uma issue. Enquanto o bloqueante estiver vermelho, o trabalho
não avança — é o que o `done` não deve contornar.

### ACHCK-R03 — JULGAMENTO PENDENTE: há gate que nenhum script computa

O `check` não computa esses: marca o alvo com `⏳` e enfileira. O veredito é de uma IA,
contra os pontos de conformidade do guia.

### ACHCK-R04 — FORA DA ESTRUTURA: o alvo não casa camada nenhuma

Nem passou nem reprovou: não havia o que confrontar. Medido ao escrever este documento —
`check --changed internal/flowx/build.go` respondeu *"is not governed by the Structure
(matches no layer in `layers:`) — nothing to confront"*.

É resultado legítimo e acionável: ou o arquivo pertence a uma camada que falta declarar,
ou ele não devia estar ali. Confundi-lo com "passou" é o silêncio que os gates existem
para acabar.

### ACHCK-R05 — MAPA DESATUALIZADO: o mapa é mais velho que os arquivos

O `check` confronta a FOTO que o mapa tem, e ela envelheceu. O trabalho pode estar certo e
o veredito, errado.
