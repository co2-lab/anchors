---
title: "Gate: rule-uses-implemented"
description: "Verifica se os campos que a regra diz usar aparecem e são consumidos no código governado."
---

> **Identificador do Gate:** `rule-uses-implemented` / `uso-de-regra-implementado`  
> **Código Interno:** `RLIMR`  
> **Categoria:** [Uso de Regras e Contratos](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Informativo no início (`blocking: false`) |

---

## 🎯 O que este gate mede?

Verifica se os campos que a regra diz usar aparecem e são consumidos no código governado.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Evita regras de papel: a spec diz que usa o campo, mas o código nunca o lê.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Verifica no código fonte correspondente se os identificadores dos campos declarados na regra são referenciados.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Os campos declarados na spec aparecem no código governado.
- **`✗ Fail` (Reprovado):** Campos declarados na spec não existem no código.
- **`~ Indeterminado/Pending`:** Código ainda não implementado (@TBD).
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: rule-uses-implemented
    on: [spec]
    check: rule-uses-implemented
    blocking: false
    measures: "os campos que a regra usa aparecem no código" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Informativo.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Implemente a leitura do campo no código fonte ou remova o campo não utilizado da spec.
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
