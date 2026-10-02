---
title: "Gate: rule-types"
description: "Garante que os prefixos e letras de códigos seguem a declaração de tipos de regras do projeto."
---

> **Identificador do Gate:** `rule-types` / `tipos-de-regra`  
> **Código Interno:** `RLTYR`  
> **Categoria:** [Falhas e Governança](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Todas as Camadas` |
| **Alvos Avaliados (`on`)** | `spec` `feature` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que os prefixos e letras de códigos seguem a declaração de tipos de regras do projeto.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Mantém o vocabulário de letras de regras padronizado em todo o repositório.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Compara as letras de regras usadas nas specs contra o vocabulário `rule_types` do anchors.yaml.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todas as letras de códigos pertencem ao vocabulário oficial.
- **`✗ Fail` (Reprovado):** Letra de código não cadastrada no anchors.yaml.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: rule-types
    blocking: true
    measures: "as letras dos códigos estão no vocabulário declarado" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante para manter a consistência da taxonomia.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Use as letras padrão (B, V, E, S, DS, VR) ou declare a nova letra em `rule_types`.
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
