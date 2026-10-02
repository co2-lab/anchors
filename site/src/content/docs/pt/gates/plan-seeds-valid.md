---
title: "Gate: plan-seeds-valid"
description: "Garante que as specs semeadas pelo plano miram camadas governadas válidas."
---

> **Identificador do Gate:** `plan-seeds-valid` / `sementes-de-plano-validas`  
> **Código Interno:** `PSVPL`  
> **Categoria:** [Planejamento e Progresso](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Planos` `doc` |
| **Alvos Avaliados (`on`)** | `plan` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que as specs semeadas pelo plano miram camadas governadas válidas.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede um plano de semear código em pastas inexistentes ou fora da planta da casa.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Confere os caminhos das specs semeadas no plano contra as regras de camadas do anchors.yaml.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todas as specs semeadas miram camadas governadas válidas.
- **`✗ Fail` (Reprovado):** Spec semeada mirando camada inexistente ou proibida.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: plan-seeds-valid
    on: [plan]
    check: plan-seeds-valid
    blocking: true
    measures: "as specs semeadas pelo plano existem nas camadas certas" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante ao semear novas fases.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Ajuste o caminho da spec no plano para casar com uma camada declarada no `anchors.yaml`.
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
