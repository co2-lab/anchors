---
title: "Gate: identity-consistent"
description: "Garante que a identidade da spec bate com o testID exposto e a imagem de baseline visual."
---

> **Identificador do Gate:** `identity-consistent` / `identidade-consistente`  
> **Código Interno:** `IDCND`  
> **Categoria:** [Apresentação e Telas](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `UI` `Telas` `Componentes` |
| **Alvos Avaliados (`on`)** | `spec` `code` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que a identidade da spec bate com o testID exposto e a imagem de baseline visual.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede divergência de nomenclatura entre o design system, o código e os testes visuais.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Cruza o identificador da unidade com o atributo testID exposto no código e o nome do arquivo PNG de baseline.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Identidade idêntica na spec, no código e no baseline visual.
- **`✗ Fail` (Reprovado):** Divergência entre o testID exposto e a identidade da spec.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: identity-consistent
    blocking: true
    measures: "identidade da spec bate com testID e baseline visual" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante em UI.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Alinhe o `testID` no componente com a identidade declarada na spec.
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
