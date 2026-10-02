---
title: "Gate: coverage-delta"
description: "Garante que a alteração atual não reduziu a cobertura de linhas em relação ao baseline."
---

> **Identificador do Gate:** `coverage-delta` / `delta-de-cobertura`  
> **Código Interno:** `INCHN`  
> **Categoria:** [Prova e Execução](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` `code` |
| **Alvos Avaliados (`on`)** | `code` |
| **Tipo de Verificação** | `Interno Determinístico` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que a alteração atual não reduziu a cobertura de linhas em relação ao baseline.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede a regressão de cobertura: uma nova PR não pode diminuir a porcentagem que o projeto já havia conquistado.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Compara a cobertura atual com o baseline registrado no snapshot anterior do projeto.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Delta de cobertura >= 0.
- **`✗ Fail` (Reprovado):** A cobertura caiu em relação ao commit base.
- **`~ Indeterminado/Pending`:** Sem baseline anterior para comparar.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: coverage-delta
    on: [code]
    check: coverage-delta
    blocking: true
    measures: "a cobertura não caiu com esta alteração" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante no CI para evitar degradação de testes.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Adicione testes para o novo código para compensar ou manter a cobertura geral.
2. Reexecute a verificação no terminal:
   ```sh
   anchors check
   ```
3. Se o gate suportar correção automática, você pode tentar o comando:
   ```sh
   anchors check --fix
   ```

---

## 🔗 Conceitos Relacionados

- [Guia Completo de Camadas](/pt/docs/layers/): Entenda quais camadas exigem este gate.
- [A Unidade (The Unit)](/pt/docs/concepts/unidade/): A relação entre Spec, Feature, Teste e Código.
- [Gates e Vereditos](/pt/docs/concepts/gates-e-vereditos/): A mecânica completa de avaliação do Anchors.
- [Catálogo Completo de Gates](/pt/docs/gates//): Retornar ao índice pesquisável de todos os gates.
