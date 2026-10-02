---
title: "anchors check & verify"
description: "Como rodar gates de qualidade no terminal e na esteira de CI/CD."
---

O comando `anchors check` é o motor central de verificação do Anchors. Ele avalia os gates de qualidade configurados contra os arquivos alterados, em stage ou em todo o repositório.

---

## 1. Sintaxe e Opções Comuns

```bash
# Verifica arquivos modificados no diretório de trabalho
anchors check

# Verifica arquivos em stage (usado pelo hook pre-commit)
anchors check --staged

# Verifica todos os arquivos regidos do projeto inteiro
anchors check --all

# Roda apenas um gate específico
anchors check --gate=unit-complete

# Roda em modo estrito (falha se houver avisos/warnings)
anchors check --strict
```

### Códigos de Saída (Exit Codes):
- `0`: Todos os gates passaram com sucesso (`OK`).
- `1`: Um ou mais gates bloqueantes retornaram `FAIL`.
- `2`: Erro de configuração ou leitura de arquivos.

---

## 2. Diferença entre `check` e `verify`

- **`anchors check`**: Avalia os **gates internos do Anchors** (completude de unidade, códigos de identidade, fronteiras arquiteturais, paridade de docs).
- **`anchors verify`**: É o **orquestrador da fase completa**. Ele roda `anchors check`, executa as suítes de teste do projeto (`anchors test`), dispara linters externos e checa mutação antes de aprovar uma entrega.

---

## 3. Exemplo na Esteira de CI/CD

```yaml
- name: Avaliar Gates do Anchors
  run: anchors check --all
```
