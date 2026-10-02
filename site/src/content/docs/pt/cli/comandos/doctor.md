---
title: "anchors doctor & audit"
description: "Diagnóstico do ecossistema, verificação de maturidade e auditoria de arquivos."
---

Enquanto o `check` avalia gates pontuais, o `anchors doctor` fornece um **raio-X global de todo o ecossistema**.

---

## 1. Executando `anchors doctor`

```bash
# Roda diagnóstico de saúde padrão
anchors doctor

# Roda com detalhamento diagnóstico completo
anchors doctor --verbose
```

### O que o `doctor` Detecta:
1. **Artefatos Órfãos**: Specs sem código implementado ou testes que verificam regras já deletadas.
2. **Colisão de Códigos de Identidade**: Códigos duplicados (ex: `AUTH-B01`) declarados em dois arquivos diferentes.
3. **Deriva Arquitetural**: Divergências entre as fronteiras de camadas declaradas em `anchors.yaml` e as importações reais.
4. **Buracos de Cobertura**: Camadas regidas que caíram abaixo do nível de maturidade exigido.

---

## 2. Auditoria Pontual com `anchors audit`

Quando um arquivo ou módulo específico falha em vários gates, use o `audit` para gerar a lista de correção:

```bash
anchors audit src/services/billing/invoice.spec.md
```

Ele gera um dossiê unificado com cada cenário ausente, asserção quebrada, dependência obsoleta e violação de camada daquela unidade.
