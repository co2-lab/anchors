---
title: "Gate: rule-uses-declared"
description: "Verifica se cada regra declara explicitamente quais campos, validações e dependências utiliza."
---

> **Identificador do Gate:** `rule-uses-declared` / `regra-declara-uso`  
> **Código Interno:** `RLUSG`  
> **Categoria:** [Uso de Regras e Contratos](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Informativo no início (`blocking: false`) |

---

## 🎯 O que este gate mede?

Verifica se cada regra declara explicitamente quais campos, validações e dependências utiliza.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Permite calcular o impacto de mudanças de campos: se um campo mudar, o Anchors sabe quais regras são afetadas.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Examina as seções de validação e uso de regras na spec procurando a declaração de campos consumidos.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** As regras declaram as variáveis e dependências que utilizam.
- **`✗ Fail` (Reprovado):** Regra complexa sem declaração de dependências de campos.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Specs puramente declarativas.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: rule-uses-declared
    on: [spec]
    check: rule-uses-declared
    blocking: false
    measures: "toda regra diz o que usa" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Informativo no início; promova a bloqueante ao refinar as specs.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Adicione a seção de campos usados na regra dentro da spec.
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
