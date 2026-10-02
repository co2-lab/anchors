---
title: "Gate: scenario-letter-declared"
description: "Verifica se a letra usada no código do cenário existe no vocabulário declarado do projeto."
---

> **Identificador do Gate:** `scenario-letter-declared` / `letra-cenario-declarada`  
> **Código Interno:** `SCLTR`  
> **Categoria:** [A Unidade e Estrutura](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` `Feature` |
| **Alvos Avaliados (`on`)** | `feature` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Verifica se a letra usada no código do cenário existe no vocabulário declarado do projeto.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede que alguém invente letras aleatórias em códigos (ex: AUTH-X01) sem padronização.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Valida a letra do código contra a lista de letras permitidas configurada em rule_types no anchors.yaml.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** A letra pertence ao vocabulário oficial (ex: B, V, E, S, DS, VR).
- **`✗ Fail` (Reprovado):** Código usando letra não declarada no vocabulário.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: scenario-letter-declared
    on: [feature]
    check: scenario-letter-declared
    blocking: true
    measures: "as letras dos códigos estão no vocabulário declarado" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante imediato.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Use uma das letras canônicas do projeto ou declare a nova letra em rule_types no anchors.yaml.
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
