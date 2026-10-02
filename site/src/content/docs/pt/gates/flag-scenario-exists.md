---
title: "Gate: flag-scenario-exists"
description: "Garante que uma regra que cita @gated-by aponta para uma flag e cenário que realmente existem."
---

> **Identificador do Gate:** `flag-scenario-exists` / `cenario-de-flag-existe`  
> **Código Interno:** `FLSCF`  
> **Categoria:** [Doutrina e Feature Flags](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que uma regra que cita @gated-by aponta para uma flag e cenário que realmente existem.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Evita links quebrados para feature flags que foram removidas ou digitadas errado.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Verifica se o identificador em `@gated-by FLG.ESTADO` resolve para um arquivo de flag existente.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** A flag e o estado citado existem.
- **`✗ Fail` (Reprovado):** Referência a flag ou estado inexistente.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Specs sem anotação @gated-by.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: flag-scenario-exists
    on: [spec]
    check: flag-scenario-exists
    blocking: true
    measures: "a anotação @gated-by aponta para flag e cenário existentes" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante imediato.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Crie o arquivo da flag ou corrija a referência em `@gated-by`.
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
