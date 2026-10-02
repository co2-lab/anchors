---
title: "anchors flow & queue"
description: "Operação de watchers em background, fila de tarefas para IAs e entrega de etapas."
---

O Anchors permite desenvolvimento assistido por IA contínuo e ordenado através de uma fila reativa de tarefas.

---

## 1. O Ciclo do Watcher e da Fila

1. **`anchors watch`**: Roda em background. Quando um arquivo é salvo (ex: `order.spec.md`), o watcher percebe o que falta e enfileira a próxima ação (`criar-feature`, `implementar-codigo`, `rodar-testes`).
2. **`anchors queue`**: Lista as tarefas pendentes na fila.
3. **`anchors next`**: O agente de IA chama este comando para puxar o próximo card sem depender de prompting manual.
4. **`anchors work <alvo>`**: Emite o prompt contextual com as instruções exatas daquela etapa.
5. **`anchors deliver`**: O agente registra a entrega; o Anchors valida o resultado com gates.
6. **`anchors done <id>`**: Arquiva o card concluído no histórico.

---

## 2. Tratando Bloqueios e Escalonamento

- **`anchors escalate`**: Se a IA encontra uma ambiguidade na regra, ela estaciona o card em `needs-user` e abre uma issue.
- **`anchors unblock`**: Quando o desenvolvedor responde a dúvida, este comando destrava o card.
- **`anchors decided <card>`**: Retoma o fluxo com a decisão humana registrada.
