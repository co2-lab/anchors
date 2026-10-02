---
title: "Gate: flag-scenarios-complete"
description: "Garante que todos os valores possíveis da feature flag possuem cenários documentados."
---

> **Identificador do Gate:** `flag-scenarios-complete` / `cenarios-de-flag-completos`  
> **Código Interno:** `FLSCF`  
> **Categoria:** [Doutrina e Feature Flags](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `flags` |
| **Alvos Avaliados (`on`)** | `flag` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que todos os valores possíveis da feature flag possuem cenários documentados.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Uma flag booleana não pode documentar apenas o caso ON e esquecer o comportamento do OFF.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Confere se todos os estados possíveis daquele tipo de flag (ex: ON e OFF) estão presentes.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todos os ramos da flag estão catalogados.
- **`✗ Fail` (Reprovado):** Falta documentar um dos estados da flag.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: flag-scenarios-complete
    on: [flag]
    check: flag-scenarios-complete
    blocking: true
    measures: "todos os valores da flag possuem cenários documentados" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante para governança de feature flags.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Adicione o cenário para os estados que faltam (ex: `### OFF — Comportamento desativado`).
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
