---
title: "anchors freeze & thaw"
description: "Como congelar o repositório para releases e descongelar para desenvolvimento."
---

Quando o projeto atinge uma release candidate ou marco de entrega, é necessário impedir modificações acidentais, mantendo apenas comandos de inspeção e auditoria liberados.

---

## 1. Congelando o Projeto (`anchors freeze`)

```bash
# Congela o projeto com justificativa de auditoria
anchors freeze --reason="Validação da Release Candidate v2.0.0"
```

Quando congelado:
- Comandos de escrita e mutação (`new`, `work`, `deliver`, `map build`) são **bloqueados**.
- Comandos somente-leitura (`status`, `doctor`, `guide`, `coverage`, `version`) **continuam funcionando**.

---

## 2. Descongelando (`anchors thaw`)

```bash
anchors thaw
```

Retoma as atividades normais de desenvolvimento e desbloqueia a fila de tarefas.
