---
title: "Gate: no-test-prova-real"
description: "Avalia se a justificativa apontada para uma dispensa @no-test é real e legítima."
---

> **Identificador do Gate:** `no-test-prova-real`  
> **Código Interno:** `RLUEX`  
> **Categoria:** [Julgamento por IA](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Julgamento por IA (ask)` |
| **Modo Recomendado** | Informativo no início (`blocking: false`) |

---

## 🎯 O que este gate mede?

Avalia se a justificativa apontada para uma dispensa @no-test é real e legítima.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede desenvolvedores de usarem `@no-test: testado em outro lugar` quando na verdade ninguém testou.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

A IA inspeciona o alvo apontado pela dispensa e valida se a prova realmente ocorre lá.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Veredito `pass` confirmando que a prova é real.
- **`✗ Fail` (Reprovado):** Veredito `fail` acusando que a justificativa é falsa ou insuficiente.
- **`~ Indeterminado/Pending`:** Julgamento pendente.
- **`Skip` (Dispensado):** Specs sem anotação @no-test.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: no-test-prova-real
    on: [spec]
    ask: "a prova apontada pelo @no-test exercita mesmo o comportamento?"
    blocking: false 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Informativo; use para auditar dispensas de teste.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Escreva o teste correspondente ou aponte a referência exata da suíte onde a prova ocorre.
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
