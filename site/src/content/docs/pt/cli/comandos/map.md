---
title: "anchors map & impact"
description: "Construção do grafo de dependências, análise da onda de impacto e recodificação."
---

O mapa de dependências do Anchors (`anchors.graph.yaml`) modela todo o repositório como um Grafo Acíclico Dirigido (DAG).

---

## 1. Gerenciando o Grafo

```bash
# Reconstrói o grafo de dependências a partir dos arquivos em disco
anchors map build

# Exibe o resumo do grafo e métricas dos nós
anchors map show

# Checa se há ciclos de dependência no projeto
anchors map circular
```

---

## 2. Análise da Onda de Impacto (`anchors impact`)

Antes de alterar uma spec ou entidade central de domínio, calcule o raio de alcance:

```bash
anchors impact src/domain/user.spec.md
```

Saída:
- **Dependentes Diretos**: Serviços e handlers que importam ou usam esta unidade.
- **Impacto Downstream**: Testes de integração e contratos que ficarão **obsoletos (stale)**.
- **Escopo Recomendado de Testes**: Conjunto mínimo de testes a rodar para validar a alteração.

---

## 3. Renomeação Atômica de Códigos (`anchors recode`)

Precisa renomear ou reorganizar o prefixo de um módulo de `USER` para `ACC`?

```bash
anchors recode USER-B01 ACC-B01
```

O comando renomeia o código atomicamente na spec, na feature, nos testes, no código fonte e no grafo de dependências de uma só vez.
