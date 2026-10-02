---
title: "Gate: revision-orphans"
description: "Identifica regras cujo significado mudou em uma revisão sem que isso tenha sido declarado formalmente."
---

> **Identificador do Gate:** `revision-orphans` / `orfas-de-revisao`  
> **Código Interno:** `RVORP`  
> **Categoria:** [Planejamento e Progresso](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Identifica regras cujo significado mudou em uma revisão sem que isso tenha sido declarado formalmente.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Altera o significado de uma regra sem avisar quebra os testes e clientes que dependiam da semântica anterior.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Detecta mudanças no texto de regras existentes sem nova numeração de versão ou nota de revisão.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Revisões de regras são explícitas e numeradas.
- **`✗ Fail` (Reprovado):** Regra alterada silenciosamente sem declaração de revisão.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: revision-orphans
    on: [spec]
    check: revision-orphans
    blocking: true
    measures: "regras cujo significado mudou são declaradas" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante para especificações maduras.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Incremente o número da regra ou adicione uma nota de revisão na spec.
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
