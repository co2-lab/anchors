---
title: "Gate: count-honored"
description: "Garante que asserções numéricas escritas na spec batem com os números reais no código."
---

> **Identificador do Gate:** `count-honored` / `contagem-honrada`  
> **Código Interno:** `CNHNC`  
> **Categoria:** [Uso de Regras e Contratos](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` |
| **Alvos Avaliados (`on`)** | `spec` `code` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Informativo no início (`blocking: false`) |

---

## 🎯 O que este gate mede?

Garante que asserções numéricas escritas na spec batem com os números reais no código.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Se a spec diz 'o lote processa no máximo 50 itens', o código não pode ter uma constante `MAX_BATCH = 100`.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Extrai números e limites citados nas regras da spec e confere com as constantes no código correspondente.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Os limites numéricos batem entre spec e código.
- **`✗ Fail` (Reprovado):** Divergência numérica entre a spec e a constante no código.
- **`~ Indeterminado/Pending`:** Código pendente.
- **`Skip` (Dispensado):** Specs sem asserções numéricas.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: count-honored
    blocking: false
    measures: "asserções numéricas da spec batem com o código" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Informativo no início; promova a bloqueante ao estabilizar constantes de negócio.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Sincronize o número na spec ou ajuste o valor da constante no código.
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
