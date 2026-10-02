---
title: "anchors review"
description: "Registra um segundo olhar sobre o que um agente decidiu, com quem olhou e o que achou."
---

Um gate que declara `review:` marca os seus alvos **a revisar**, à parte de como mede. Um alvo fica a revisar até haver uma review registrada na revisão atual dele; uma mudança no alvo torna a review devida de novo. Reviews informam e nunca bloqueiam.

```bash
# O que está a revisar, por gate, com a pergunta
anchors review --pending [--gate <gate>]

# Uma review que não achou nada
anchors review <alvo> --gate <gate> --by human:ana

# Uma review com achados: o relatório inteiro, cada achado com o quê, onde e por quê
anchors review <alvo> --gate <gate> --by agent:<fornecedor>/<modelo> --findings "<relatório>"
```

| Flag | O que faz |
| --- | --- |
| `--pending` | lista os alvos a revisar |
| `--gate` | o gate desta review — um gate que declara `review:` |
| `--by` | quem revisou, do jeito que se nomeia; o registro guarda e não classifica |
| `--findings` | o relatório de achados; abre a issue do alvo para o gate |
| `--record-issues` | no modo manual, também escreve a issue dos achados |

O produto de uma review são achados, não um carimbo: cada um vira um teste que falha e uma correção, ou é descartado com o motivo na issue. Não há `waived` — uma review que não aconteceu continua devida. Veja [julgamento e review](/pt/docs/concepts/ai-judgment/).
