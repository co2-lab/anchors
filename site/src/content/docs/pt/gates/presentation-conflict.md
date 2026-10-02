---
title: "Gate: presentation-conflict"
description: "Garante que um prop e uma condição não levam a duas aparências conflitantes."
---

> **Identificador do Gate:** `presentation-conflict` / `apresentacao-sem-conflito`  
> **Código Interno:** `PRSNT`  
> **Categoria:** [Apresentação e Telas](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `UI` `Telas` `Componentes` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que um prop e uma condição não levam a duas aparências conflitantes.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede regras ambíguas que dizem ao mesmo tempo 'botão fica vermelho' e 'botão fica desabilitado cinza'.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Analisa as matrizes de decisão visual procurando sobreposição de condições contraditórias.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Nenhum conflito de aparência encontrado.
- **`✗ Fail` (Reprovado):** Mesma condição conduz a duas aparências visuais incompatíveis.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: presentation-conflict
    on: [spec]
    check: presentation-conflict
    blocking: true
    measures: "um prop e uma condição não levam a duas aparências" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante para design systems e componentes visuais.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Defina a ordem de precedência clara entre os estados visuais na spec.
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
