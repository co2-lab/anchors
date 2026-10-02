---
title: "Gate: plan-change-justified"
description: "Garante que um plano ou spec modificado declare no diff por que a mudança aconteceu."
---

> **Identificador do Gate:** `plan-change-justified` / `plano-alterado-justificado`  
> **Código Interno:** `PCJPL`  
> **Categoria:** [Planejamento e Progresso](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Planos` `Spec` |
| **Alvos Avaliados (`on`)** | `plan` `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que um plano ou spec modificado declare no diff por que a mudança aconteceu.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

A deriva silenciosa é perigosa: alterar um plano sem justificativa apaga a inconsistência sem deixar rastro do motivo.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Analisa o git diff e confere se há nota explicativa justificando a alteração da especificação.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** A alteração no arquivo acompanha justificativa de mudança.
- **`✗ Fail` (Reprovado):** Alteração estrutural de plano/spec sem justificativa no commit/diff.
- **`~ Indeterminado/Pending`:** Sem git para inspecionar diff.
- **`Skip` (Dispensado):** Criação de arquivos novos.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: plan-change-justified
    on: [plan, spec]
    check: plan-change-justified
    blocking: true
    measures: "um plano/spec que mudou declara por quê" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante para equipes que precisam de auditoria rigorosa de mudanças.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Adicione uma nota ou entrada de revisão explicando a razão da alteração na spec ou plano.
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
