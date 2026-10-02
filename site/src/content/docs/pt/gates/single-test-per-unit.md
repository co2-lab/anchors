---
title: "Gate: single-test-per-unit"
description: "Garante que uma unidade tem um único arquivo de teste por camada de teste (ou declara divisão com @split-test)."
---

> **Identificador do Gate:** `single-test-per-unit` / `unico-teste-por-unidade`  
> **Código Interno:** `SNGTU`  
> **Categoria:** [Prova e Execução](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` |
| **Alvos Avaliados (`on`)** | `code` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que uma unidade tem um único arquivo de teste por camada de teste (ou declara divisão com @split-test).

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Evita testes duplicados ou espalhados em vários arquivos que testam a mesma unidade sem critério.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Conta quantos arquivos de teste apontam para a mesma unidade e valida a presença de @split-test se houver mais de um.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Existe exatamente um arquivo de teste por camada para a unidade, ou a divisão foi justificada.
- **`✗ Fail` (Reprovado):** Múltiplos arquivos de teste encontrados para a mesma unidade sem declaração @split-test.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: single-test-per-unit
    on: [code]
    check: single-test-per-unit
    blocking: true
    measures: "uma unidade tem um arquivo de teste por camada de teste" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante para manter a organização dos testes.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Una os arquivos de teste ou declare @split-test no cabeçalho justificando a divisão.
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
