---
title: "Gate: phase-exists"
description: "Garante que as fases citadas nas tarefas e dependências existem no plano."
---

> **Identificador do Gate:** `phase-exists` / `fase-existe`  
> **Código Interno:** `PHORP`  
> **Categoria:** [Planejamento e Progresso](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Planos` `doc` |
| **Alvos Avaliados (`on`)** | `plan` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que as fases citadas nas tarefas e dependências existem no plano.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede apontar dependências para fases inexistentes ou renomeadas.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Cruza referências `phase:` com a lista de identificadores de fases declaradas no plano.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todas as fases citadas existem.
- **`✗ Fail` (Reprovado):** Referência a fase inexistente.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: phase-exists
    on: [plan]
    check: phase-exists
    blocking: true
    measures: "as fases citadas existem no plano" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante para governança de planos.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Crie a fase no plano ou corrija o nome da referência.
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
