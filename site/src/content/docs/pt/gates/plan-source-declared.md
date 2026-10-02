---
title: "Gate: plan-source-declared"
description: "Garante que um plano que cita uma fonte externa declara explicitamente quem a constrói."
---

> **Identificador do Gate:** `plan-source-declared` / `fonte-de-plano-declarada`  
> **Código Interno:** `PSDPL`  
> **Categoria:** [Planejamento e Progresso](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Planos` `doc` |
| **Alvos Avaliados (`on`)** | `plan` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que um plano que cita uma fonte externa declara explicitamente quem a constrói.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Evita planos baseados em premissas ou artefatos que ninguém se responsabilizou por gerar.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Verifica se referências a fontes possuem a declaração de builder ou responsável.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todas as fontes possuem construtor declarado.
- **`✗ Fail` (Reprovado):** Fonte citada sem declaração de quem a constrói.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Planos autossuficientes.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: plan-source-declared
    on: [plan]
    check: plan-source-declared
    blocking: true
    measures: "um plano que cita uma fonte declara quem a constrói" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante para planos complexos.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Declare o responsável ou passo de geração da fonte no plano.
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
